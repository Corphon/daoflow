// api/health.go

package api

import (
    "context"
    "sync"
    "time"

    "github.com/Corphon/daoflow/system"
)

// HealthStatus 健康状态
type HealthStatus string

const (
    StatusHealthy    HealthStatus = "healthy"    // 健康
    StatusDegraded   HealthStatus = "degraded"   // 性能降级
    StatusUnhealthy  HealthStatus = "unhealthy"  // 不健康
    StatusCritical   HealthStatus = "critical"   // 严重故障
    StatusRecovering HealthStatus = "recovering" // 恢复中
)

// ComponentHealth 组件健康信息
type ComponentHealth struct {
    Name          string                 `json:"name"`           // 组件名称
    Status        HealthStatus           `json:"status"`         // 健康状态
    LastCheck     time.Time              `json:"last_check"`     // 最后检查时间
    LastFailure   time.Time              `json:"last_failure"`   // 最后故障时间
    ErrorCount    int                    `json:"error_count"`    // 错误计数
    Performance   map[string]float64     `json:"performance"`    // 性能指标
    Dependencies  []string               `json:"dependencies"`   // 依赖组件
    Details       map[string]interface{} `json:"details"`        // 详细信息
}

// HealthCheck 健康检查配置
type HealthCheck struct {
    ID           string        `json:"id"`            // 检查ID
    Component    string        `json:"component"`     // 目标组件
    Type         string        `json:"type"`          // 检查类型
    Interval     time.Duration `json:"interval"`      // 检查间隔
    Timeout      time.Duration `json:"timeout"`       // 超时时间
    Threshold    float64       `json:"threshold"`     // 阈值
    RetryCount   int           `json:"retry_count"`   // 重试次数
    Enabled      bool          `json:"enabled"`       // 是否启用
}

// SystemHealth 系统健康状态
type SystemHealth struct {
    Status        HealthStatus              `json:"status"`         // 整体状态
    Components    map[string]ComponentHealth `json:"components"`    // 组件状态
    StartTime     time.Time                 `json:"start_time"`    // 启动时间
    Uptime        time.Duration             `json:"uptime"`        // 运行时间
    LastIncident  time.Time                 `json:"last_incident"` // 最后故障
    HealthScore   float64                   `json:"health_score"`  // 健康评分
}

// HealthEvent 健康事件
type HealthEvent struct {
    Type         string                 `json:"type"`          // 事件类型
    Component    string                 `json:"component"`     // 相关组件
    Status       HealthStatus           `json:"status"`        // 健康状态
    Timestamp    time.Time              `json:"timestamp"`     // 事件时间
    Description  string                 `json:"description"`   // 描述信息
    Metrics      map[string]float64     `json:"metrics"`      // 相关指标
}

// HealthAPI 健康检查API
type HealthAPI struct {
    mu     sync.RWMutex
    system *system.SystemCore

    // 组件健康状态
    components map[string]*ComponentHealth
    
    // 健康检查配置
    checks    map[string]*HealthCheck
    
    // 检查协程控制
    checkers  map[string]context.CancelFunc
    
    // 事件通道
    events    chan HealthEvent
    
    // 系统启动时间
    startTime time.Time
    
    ctx      context.Context
    cancel   context.CancelFunc
}

// NewHealthAPI 创建健康检查API实例
func NewHealthAPI(sys *system.SystemCore) *HealthAPI {
    ctx, cancel := context.WithCancel(context.Background())
    
    api := &HealthAPI{
        system:     sys,
        components: make(map[string]*ComponentHealth),
        checks:     make(map[string]*HealthCheck),
        checkers:   make(map[string]context.CancelFunc),
        events:     make(chan HealthEvent, 100),
        startTime:  time.Now(),
        ctx:        ctx,
        cancel:     cancel,
    }
    
    // 初始化默认健康检查
    api.initDefaultChecks()
    
    return api
}

// RegisterCheck 注册健康检查
func (h *HealthAPI) RegisterCheck(check *HealthCheck) error {
    h.mu.Lock()
    defer h.mu.Unlock()

    if check.ID == "" {
        return NewError(ErrInvalidCheck, "check ID is required")
    }

    // 保存检查配置
    h.checks[check.ID] = check

    // 如果检查已启用，启动检查协程
    if check.Enabled {
        h.startChecker(check)
    }

    return nil
}

// StartCheck 启动健康检查
func (h *HealthAPI) StartCheck(checkID string) error {
    h.mu.Lock()
    defer h.mu.Unlock()

    check, exists := h.checks[checkID]
    if !exists {
        return NewError(ErrCheckNotFound, "check not found")
    }

    if !check.Enabled {
        check.Enabled = true
        h.startChecker(check)
    }

    return nil
}

// StopCheck 停止健康检查
func (h *HealthAPI) StopCheck(checkID string) error {
    h.mu.Lock()
    defer h.mu.Unlock()

    check, exists := h.checks[checkID]
    if !exists {
        return NewError(ErrCheckNotFound, "check not found")
    }

    if check.Enabled {
        check.Enabled = false
        if cancel, exists := h.checkers[checkID]; exists {
            cancel()
            delete(h.checkers, checkID)
        }
    }

    return nil
}

// GetComponentHealth 获取组件健康状态
func (h *HealthAPI) GetComponentHealth(component string) (*ComponentHealth, error) {
    h.mu.RLock()
    defer h.mu.RUnlock()

    health, exists := h.components[component]
    if !exists {
        return nil, NewError(ErrComponentNotFound, "component not found")
    }

    return health, nil
}

// GetSystemHealth 获取系统健康状态
func (h *HealthAPI) GetSystemHealth() (*SystemHealth, error) {
    h.mu.RLock()
    defer h.mu.RUnlock()

    health := &SystemHealth{
        Components:   make(map[string]ComponentHealth),
        StartTime:   h.startTime,
        Uptime:      time.Since(h.startTime),
    }

    // 复制组件状态
    for name, comp := range h.components {
        health.Components[name] = *comp
    }

    // 计算系统整体状态
    health.Status = h.calculateSystemStatus()
    health.HealthScore = h.calculateHealthScore()

    return health, nil
}

// Subscribe 订阅健康事件
func (h *HealthAPI) Subscribe() (<-chan HealthEvent, error) {
    return h.events, nil
}

// startChecker 启动检查协程
func (h *HealthAPI) startChecker(check *HealthCheck) {
    checkCtx, cancel := context.WithCancel(h.ctx)
    h.checkers[check.ID] = cancel

    go func() {
        ticker := time.NewTicker(check.Interval)
        defer ticker.Stop()

        for {
            select {
            case <-checkCtx.Done():
                return
            case <-ticker.C:
                h.performCheck(check)
            }
        }
    }()
}

// performCheck 执行健康检查
func (h *HealthAPI) performCheck(check *HealthCheck) {
    h.mu.Lock()
    defer h.mu.Unlock()

    // 获取或创建组件健康状态
    health, exists := h.components[check.Component]
    if !exists {
        health = &ComponentHealth{
            Name:        check.Component,
            Status:      StatusHealthy,
            Performance: make(map[string]float64),
            Details:    make(map[string]interface{}),
        }
        h.components[check.Component] = health
    }

    // 执行检查并更新状态
    ctx, cancel := context.WithTimeout(context.Background(), check.Timeout)
    defer cancel()

    if err := h.system.CheckComponentHealth(ctx, check.Component); err != nil {
        h.handleCheckFailure(health, check, err)
    } else {
        h.handleCheckSuccess(health)
    }
}

// handleCheckFailure 处理检查失败
func (h *HealthAPI) handleCheckFailure(health *ComponentHealth, check *HealthCheck, err error) {
    health.ErrorCount++
    health.LastFailure = time.Now()

    // 根据错误次数决定状态
    if health.ErrorCount >= check.RetryCount {
        health.Status = StatusUnhealthy
        if health.ErrorCount >= check.RetryCount*2 {
            health.Status = StatusCritical
        }
    }

    // 发送事件
    h.events <- HealthEvent{
        Type:        "check_failed",
        Component:   health.Name,
        Status:      health.Status,
        Timestamp:   time.Now(),
        Description: err.Error(),
    }
}

// handleCheckSuccess 处理检查成功
func (h *HealthAPI) handleCheckSuccess(health *ComponentHealth) {
    if health.Status != StatusHealthy {
        health.Status = StatusRecovering
        if health.ErrorCount == 0 {
            health.Status = StatusHealthy
        }
    }
    
    health.LastCheck = time.Now()
    health.ErrorCount = 0
}

// calculateSystemStatus 计算系统整体状态
func (h *HealthAPI) calculateSystemStatus() HealthStatus {
    criticalCount := 0
    unhealthyCount := 0
    
    for _, comp := range h.components {
        switch comp.Status {
        case StatusCritical:
            criticalCount++
        case StatusUnhealthy:
            unhealthyCount++
        }
    }

    if criticalCount > 0 {
        return StatusCritical
    }
    if unhealthyCount > 0 {
        return StatusDegraded
    }
    return StatusHealthy
}

// calculateHealthScore 计算健康评分
func (h *HealthAPI) calculateHealthScore() float64 {
    if len(h.components) == 0 {
        return 100.0
    }

    var totalScore float64
    for _, comp := range h.components {
        switch comp.Status {
        case StatusHealthy:
            totalScore += 100.0
        case StatusDegraded:
            totalScore += 60.0
        case StatusUnhealthy:
            totalScore += 30.0
        case StatusCritical:
            totalScore += 0.0
        case StatusRecovering:
            totalScore += 80.0
        }
    }

    return totalScore / float64(len(h.components))
}

// initDefaultChecks 初始化默认健康检查
func (h *HealthAPI) initDefaultChecks() {
    // 添加默认的系统组件检查
    defaultChecks := []*HealthCheck{
        {
            ID:         "system-core",
            Component:  "system",
            Type:      "core",
            Interval:  time.Second * 30,
            Timeout:   time.Second * 5,
            RetryCount: 3,
            Enabled:   true,
        },
        // 添加其他默认检查...
    }

    for _, check := range defaultChecks {
        h.RegisterCheck(check)
    }
}

// Close 关闭API
func (h *HealthAPI) Close() error {
    h.cancel()
    close(h.events)
    return nil
}
