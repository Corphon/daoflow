// api/lifecycle_api.go

package api

import (
    "context"
    "sync"
    "time"

    "github.com/Corphon/daoflow/system"
)

// LifecycleState 系统生命周期状态
type LifecycleState string

const (
    StateUninitialized LifecycleState = "uninitialized" // 未初始化
    StateInitializing  LifecycleState = "initializing"  // 初始化中
    StateReady         LifecycleState = "ready"         // 就绪
    StateRunning       LifecycleState = "running"       // 运行中
    StatePaused        LifecycleState = "paused"        // 已暂停
    StateStopping      LifecycleState = "stopping"      // 停止中
    StateStopped       LifecycleState = "stopped"       // 已停止
    StateError         LifecycleState = "error"         // 错误状态
)

// LifecycleStatus 生命周期状态信息
type LifecycleStatus struct {
    State       LifecycleState     `json:"state"`        // 当前状态
    StartTime   time.Time          `json:"start_time"`   // 启动时间
    Uptime      time.Duration      `json:"uptime"`       // 运行时长
    ErrorCount  int                `json:"error_count"`  // 错误计数
    LastError   string             `json:"last_error"`   // 最后错误
    Metrics     map[string]float64 `json:"metrics"`      // 核心指标
}

// LifecycleAPI 生命周期管理API
type LifecycleAPI struct {
    mu sync.RWMutex

    system  *system.SystemCore    // 系统核心引用
    state   LifecycleState       // 当前状态
    metrics *LifecycleMetrics    // 状态指标
    events  chan LifecycleEvent  // 事件通道

    // 配置选项
    opts *Options

    // 上下文控制
    ctx    context.Context
    cancel context.CancelFunc
}

// LifecycleMetrics 生命周期指标
type LifecycleMetrics struct {
    StartTime     time.Time
    ErrorCount    int
    LastError     error
    StateChanges  int
    HealthChecks  int
    FailureRate   float64
}

// LifecycleEvent 生命周期事件
type LifecycleEvent struct {
    Type      string          `json:"type"`       // 事件类型
    State     LifecycleState  `json:"state"`      // 相关状态
    Timestamp time.Time       `json:"timestamp"`  // 发生时间
    Details   interface{}     `json:"details"`    // 详细信息
}

// NewLifecycleAPI 创建生命周期API实例
func NewLifecycleAPI(system *system.SystemCore, opts *Options) *LifecycleAPI {
    ctx, cancel := context.WithCancel(context.Background())

    api := &LifecycleAPI{
        system:  system,
        state:   StateUninitialized,
        metrics: &LifecycleMetrics{},
        events:  make(chan LifecycleEvent, 100),
        opts:    opts,
        ctx:     ctx,
        cancel:  cancel,
    }

    // 启动监控协程
    go api.monitor()

    return api
}

// Initialize 初始化系统
func (l *LifecycleAPI) Initialize() error {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.state != StateUninitialized {
        return NewError(ErrInvalidState, "system already initialized")
    }

    l.setState(StateInitializing)

    // 执行系统初始化
    if err := l.system.Initialize(); err != nil {
        l.recordError(err)
        l.setState(StateError)
        return err
    }

    l.setState(StateReady)
    return nil
}

// Start 启动系统
func (l *LifecycleAPI) Start() error {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.state != StateReady && l.state != StatePaused {
        return NewError(ErrInvalidState, "system not ready to start")
    }

    // 执行系统启动
    if err := l.system.Start(); err != nil {
        l.recordError(err)
        l.setState(StateError)
        return err
    }

    l.metrics.StartTime = time.Now()
    l.setState(StateRunning)
    return nil
}

// Stop 停止系统
func (l *LifecycleAPI) Stop() error {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.state != StateRunning && l.state != StatePaused {
        return NewError(ErrInvalidState, "system not running")
    }

    l.setState(StateStopping)

    // 执行系统停止
    if err := l.system.Stop(); err != nil {
        l.recordError(err)
        l.setState(StateError)
        return err
    }

    l.setState(StateStopped)
    return nil
}

// Pause 暂停系统
func (l *LifecycleAPI) Pause() error {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.state != StateRunning {
        return NewError(ErrInvalidState, "system not running")
    }

    // 执行系统暂停
    if err := l.system.Pause(); err != nil {
        l.recordError(err)
        return err
    }

    l.setState(StatePaused)
    return nil
}

// Resume 恢复系统运行
func (l *LifecycleAPI) Resume() error {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.state != StatePaused {
        return NewError(ErrInvalidState, "system not paused")
    }

    // 执行系统恢复
    if err := l.system.Resume(); err != nil {
        l.recordError(err)
        return err
    }

    l.setState(StateRunning)
    return nil
}

// GetStatus 获取系统状态
func (l *LifecycleAPI) GetStatus() (*LifecycleStatus, error) {
    l.mu.RLock()
    defer l.mu.RUnlock()

    var uptime time.Duration
    if !l.metrics.StartTime.IsZero() {
        uptime = time.Since(l.metrics.StartTime)
    }

    status := &LifecycleStatus{
        State:      l.state,
        StartTime:  l.metrics.StartTime,
        Uptime:     uptime,
        ErrorCount: l.metrics.ErrorCount,
    }

    if l.metrics.LastError != nil {
        status.LastError = l.metrics.LastError.Error()
    }

    // 获取核心指标
    status.Metrics = map[string]float64{
        "state_changes": float64(l.metrics.StateChanges),
        "health_checks": float64(l.metrics.HealthChecks),
        "failure_rate": l.metrics.FailureRate,
    }

    return status, nil
}

// Subscribe 订阅生命周期事件
func (l *LifecycleAPI) Subscribe() (<-chan LifecycleEvent, error) {
    return l.events, nil
}

// monitor 内部监控协程
func (l *LifecycleAPI) monitor() {
    ticker := time.NewTicker(time.Second * 30)
    defer ticker.Stop()

    for {
        select {
        case <-l.ctx.Done():
            return
        case <-ticker.C:
            l.checkHealth()
        }
    }
}

// setState 更新状态并发送事件
func (l *LifecycleAPI) setState(state LifecycleState) {
    l.state = state
    l.metrics.StateChanges++

    // 发送状态变更事件
    l.events <- LifecycleEvent{
        Type:      "state_change",
        State:     state,
        Timestamp: time.Now(),
    }
}

// recordError 记录错误
func (l *LifecycleAPI) recordError(err error) {
    l.metrics.ErrorCount++
    l.metrics.LastError = err
    
    // 计算失败率
    total := float64(l.metrics.HealthChecks)
    if total > 0 {
        l.metrics.FailureRate = float64(l.metrics.ErrorCount) / total
    }
}

// checkHealth 健康检查
func (l *LifecycleAPI) checkHealth() {
    l.metrics.HealthChecks++

    // 执行系统健康检查
    if err := l.system.HealthCheck(); err != nil {
        l.recordError(err)
        l.setState(StateError)
    }
}

// Close 关闭API
func (l *LifecycleAPI) Close() error {
    l.cancel()
    close(l.events)
    return nil
}
