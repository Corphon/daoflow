// system/monitor/alert/detector.go

package alert

import (
    "context"
    "sync"
    "time"

    "github.com/Corphon/daoflow/system/types"
)

// AlertCondition 告警条件
type AlertCondition struct {
    ID          string                 // 条件ID
    Name        string                 // 条件名称
    Type        string                 // 条件类型
    Metric      string                 // 监控指标
    Operator    string                 // 操作符
    Threshold   float64                // 阈值
    Duration    time.Duration          // 持续时间
    Severity    types.IssueSeverity    // 严重程度
    Labels      map[string]string      // 标签
    Actions     []string               // 触发动作
}

// AlertState 告警状态
type AlertState struct {
    ConditionID string
    Active      bool
    StartTime   time.Time
    LastUpdate  time.Time
    Value       float64
    Count       int
}

// Detector 告警检测器
type Detector struct {
    mu sync.RWMutex

    // 配置
    config struct {
        CheckInterval time.Duration
        BufferSize    int
        MinInterval   time.Duration
        MaxRetries    int
    }

    // 条件管理
    conditions map[string]*AlertCondition
    states     map[string]*AlertState

    // 通知通道
    alertChan chan types.Alert

    // 检测状态
    status struct {
        isRunning   bool
        lastCheck   time.Time
        errorCount  int
        errors     []error
    }

    // 指标源
    metricsSource MetricsSource
}

// MetricsSource 指标数据源接口
type MetricsSource interface {
    GetMetric(name string) (float64, error)
    GetMetrics() (map[string]float64, error)
}

// NewDetector 创建新的检测器
func NewDetector(config types.AlertConfig, metricsSource MetricsSource) *Detector {
    d := &Detector{
        conditions:    make(map[string]*AlertCondition),
        states:       make(map[string]*AlertState),
        alertChan:    make(chan types.Alert, config.BufferSize),
        metricsSource: metricsSource,
    }

    // 设置配置
    d.config.CheckInterval = config.CheckInterval
    d.config.BufferSize = config.BufferSize
    d.config.MinInterval = config.MinInterval
    d.config.MaxRetries = config.MaxRetries

    return d
}

// Start 启动检测器
func (d *Detector) Start(ctx context.Context) error {
    d.mu.Lock()
    if d.status.isRunning {
        d.mu.Unlock()
        return types.NewSystemError(types.ErrRuntime, "detector already running", nil)
    }
    d.status.isRunning = true
    d.mu.Unlock()

    // 启动检测循环
    go d.detectionLoop(ctx)

    return nil
}

// Stop 停止检测器
func (d *Detector) Stop() error {
    d.mu.Lock()
    defer d.mu.Unlock()

    d.status.isRunning = false
    return nil
}

// AddCondition 添加告警条件
func (d *Detector) AddCondition(condition *AlertCondition) error {
    d.mu.Lock()
    defer d.mu.Unlock()

    if _, exists := d.conditions[condition.ID]; exists {
        return types.NewSystemError(types.ErrExists, "condition already exists", nil)
    }

    d.conditions[condition.ID] = condition
    d.states[condition.ID] = &AlertState{
        ConditionID: condition.ID,
        Active:     false,
    }

    return nil
}

// RemoveCondition 移除告警条件
func (d *Detector) RemoveCondition(id string) error {
    d.mu.Lock()
    defer d.mu.Unlock()

    if _, exists := d.conditions[id]; !exists {
        return types.NewSystemError(types.ErrNotFound, "condition not found", nil)
    }

    delete(d.conditions, id)
    delete(d.states, id)

    return nil
}

// detectionLoop 检测循环
func (d *Detector) detectionLoop(ctx context.Context) {
    ticker := time.NewTicker(d.config.CheckInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            if err := d.check(); err != nil {
                d.handleError(err)
            }
        }
    }
}

// check 执行检查
func (d *Detector) check() error {
    d.mu.Lock()
    defer d.mu.Unlock()

    metrics, err := d.metricsSource.GetMetrics()
    if err != nil {
        return err
    }

    for id, condition := range d.conditions {
        state := d.states[id]
        
        // 获取当前值
        value, exists := metrics[condition.Metric]
        if !exists {
            continue
        }

        // 检查条件
        exceeded := d.evaluateCondition(condition, value)
        
        // 更新状态
        if exceeded {
            if !state.Active {
                // 新的告警
                state.Active = true
                state.StartTime = time.Now()
                state.Count = 1
            } else {
                state.Count++
            }
            
            // 检查持续时间
            if time.Since(state.StartTime) >= condition.Duration {
                d.triggerAlert(condition, value)
            }
        } else {
            if state.Active {
                // 恢复告警
                d.resolveAlert(condition)
            }
            state.Active = false
            state.Count = 0
        }
        
        state.LastUpdate = time.Now()
        state.Value = value
    }

    d.status.lastCheck = time.Now()
    return nil
}

// evaluateCondition 评估条件
func (d *Detector) evaluateCondition(condition *AlertCondition, value float64) bool {
    switch condition.Operator {
    case ">":
        return value > condition.Threshold
    case ">=":
        return value >= condition.Threshold
    case "<":
        return value < condition.Threshold
    case "<=":
        return value <= condition.Threshold
    case "==":
        return value == condition.Threshold
    default:
        return false
    }
}

// triggerAlert 触发告警
func (d *Detector) triggerAlert(condition *AlertCondition, value float64) {
    alert := types.Alert{
        ID:        generateAlertID(),
        Type:      condition.Type,
        Level:     condition.Severity,
        Message:   generateAlertMessage(condition, value),
        Source:    condition.Name,
        Time:      time.Now(),
        Status:    "firing",
        Labels:    condition.Labels,
    }

    select {
    case d.alertChan <- alert:
    default:
        // 缓冲区满时记录错误
        d.handleError(types.NewSystemError(types.ErrOverflow, "alert buffer full", nil))
    }
}

// resolveAlert 解除告警
func (d *Detector) resolveAlert(condition *AlertCondition) {
    alert := types.Alert{
        ID:        generateAlertID(),
        Type:      condition.Type,
        Level:     condition.Severity,
        Message:   "Alert resolved: " + condition.Name,
        Source:    condition.Name,
        Time:      time.Now(),
        Status:    "resolved",
        Labels:    condition.Labels,
    }

    select {
    case d.alertChan <- alert:
    default:
        d.handleError(types.NewSystemError(types.ErrOverflow, "alert buffer full", nil))
    }
}

// handleError 处理错误
func (d *Detector) handleError(err error) {
    d.mu.Lock()
    defer d.mu.Unlock()

    d.status.errors = append(d.status.errors, err)
    d.status.errorCount++
}

// GetAlertChannel 获取告警通道
func (d *Detector) GetAlertChannel() <-chan types.Alert {
    return d.alertChan
}

// generateAlertID 生成告警ID
func generateAlertID() string {
    return fmt.Sprintf("alert-%d", time.Now().UnixNano())
}

// generateAlertMessage 生成告警消息
func generateAlertMessage(condition *AlertCondition, value float64) string {
    return fmt.Sprintf("%s: current value %.2f %s threshold %.2f",
        condition.Name, value, condition.Operator, condition.Threshold)
}
