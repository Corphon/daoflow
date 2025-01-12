// system/monitor/metrics/collector.go

package metrics

import (
    "context"
    "sync"
    "time"

    "github.com/Corphon/daoflow/system/types"
)

type GenerationMetrics struct {
    TotalGenerated  int
    SuccessRate     float64
    AverageScore    float64
    ComplexityDist  map[float64]int
    Evolution       []types.MetricPoint
}

// Collector 指标收集器
type Collector struct {
    mu sync.RWMutex

    // 基础配置
    config types.MetricsConfig

    // 当前指标
    current types.MetricsData

    // 指标历史
    history []types.MetricsData

    // 采集状态
    status struct {
        isRunning bool
        lastRun   time.Time
        errors    []error
    }

    // 采集器组件
    collectors map[types.MetricType]MetricCollector

    // 通知通道
    notifications chan types.Alert
}

// MetricCollector 具体指标收集器接口
type MetricCollector interface {
    Collect(context.Context) (interface{}, error)
    Type() types.MetricType
    Validate() error
}

// NewCollector 创建新的指标收集器
func NewCollector(config types.MetricsConfig) *Collector {
    return &Collector{
        config:        config,
        collectors:    make(map[types.MetricType]MetricCollector),
        notifications: make(chan types.Alert, 100),
        history:      make([]types.MetricsData, 0, config.HistorySize),
    }
}

// Start 启动指标收集
func (c *Collector) Start(ctx context.Context) error {
    c.mu.Lock()
    if c.status.isRunning {
        c.mu.Unlock()
        return types.NewSystemError(types.ErrRuntime, "collector already running", nil)
    }
    c.status.isRunning = true
    c.mu.Unlock()

    // 启动收集循环
    go c.collectLoop(ctx)

    return nil
}

// Stop 停止指标收集
func (c *Collector) Stop() error {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.status.isRunning = false
    return nil
}

// RegisterCollector 注册指标收集器
func (c *Collector) RegisterCollector(collector MetricCollector) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    // 验证收集器
    if err := collector.Validate(); err != nil {
        return types.WrapError(err, types.ErrInvalid, "invalid collector")
    }

    // 注册收集器
    metricType := collector.Type()
    if _, exists := c.collectors[metricType]; exists {
        return types.NewSystemError(types.ErrExists, "collector already registered", nil)
    }

    c.collectors[metricType] = collector
    return nil
}

// GetCurrentMetrics 获取当前指标
func (c *Collector) GetCurrentMetrics() types.MetricsData {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    return c.current
}

// GetMetricsHistory 获取指标历史
func (c *Collector) GetMetricsHistory() []types.MetricsData {
    c.mu.RLock()
    defer c.mu.RUnlock()

    // 创建副本以避免数据竞争
    history := make([]types.MetricsData, len(c.history))
    copy(history, c.history)
    return history
}

// collectLoop 指标收集循环
func (c *Collector) collectLoop(ctx context.Context) {
    ticker := time.NewTicker(c.config.Interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            if err := c.collect(ctx); err != nil {
                c.handleError(err)
            }
        }
    }
}

// collect 执行一次完整的指标收集
func (c *Collector) collect(ctx context.Context) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    // 创建新的指标数据
    metrics := types.MetricsData{
        Timestamp: time.Now(),
    }

    // 收集各类指标
    for typ, collector := range c.collectors {
        data, err := collector.Collect(ctx)
        if err != nil {
            c.status.errors = append(c.status.errors, err)
            continue
        }

        // 根据指标类型存储数据
        switch typ {
        case types.MetricEnergy:
            metrics.System.Energy = data.(float64)
        case types.MetricField:
            metrics.System.Field = data.(types.FieldState)
        case types.MetricQuantum:
            metrics.System.Quantum = data.(types.QuantumState)
        case types.MetricEmergence:
            metrics.System.Emergence = data.(types.EmergentProperty)
        }
    }

    // 更新当前指标
    c.current = metrics

    // 添加到历史记录
    c.history = append(c.history, metrics)

    // 如果历史记录超过限制，删除最旧的记录
    if len(c.history) > c.config.HistorySize {
        c.history = c.history[1:]
    }

    // 检查指标阈值并生成告警
    c.checkThresholds(metrics)

    c.status.lastRun = time.Now()
    return nil
}

// checkThresholds 检查指标阈值
func (c *Collector) checkThresholds(metrics types.MetricsData) {
    // 检查系统能量
    if metrics.System.Energy < c.config.Thresholds["min_energy"] {
        c.notify(types.Alert{
            Level:   types.SeverityWarning,
            Type:    "energy_low",
            Message: "System energy below minimum threshold",
            Time:    time.Now(),
        })
    }

    // 检查场强度
    if field := metrics.System.Field; field.Strength > c.config.Thresholds["max_field_strength"] {
        c.notify(types.Alert{
            Level:   types.SeverityWarning,
            Type:    "field_high",
            Message: "Field strength exceeds maximum threshold",
            Time:    time.Now(),
        })
    }

    // 检查量子相干性
    if quantum := metrics.System.Quantum; quantum.Coherence < c.config.Thresholds["min_coherence"] {
        c.notify(types.Alert{
            Level:   types.SeverityWarning,
            Type:    "coherence_low",
            Message: "Quantum coherence below minimum threshold",
            Time:    time.Now(),
        })
    }
}

// handleError 处理收集过程中的错误
func (c *Collector) handleError(err error) {
    c.notify(types.Alert{
        Level:   types.SeverityError,
        Type:    "collector_error",
        Message: err.Error(),
        Time:    time.Now(),
    })
}

// notify 发送告警通知
func (c *Collector) notify(alert types.Alert) {
    select {
    case c.notifications <- alert:
    default:
        // 如果通道已满，记录错误但不阻塞
    }
}
