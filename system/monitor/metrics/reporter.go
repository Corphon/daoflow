// system/monitor/metrics/reporter.go

package metrics

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "time"

    "github.com/Corphon/daoflow/system/types"
)

// Reporter 指标报告器
type Reporter struct {
    mu sync.RWMutex

    // 基础配置
    config struct {
        Interval    time.Duration          // 报告间隔
        Format      string                 // 报告格式
        OutputPath  string                 // 输出路径
        Thresholds  map[string]float64     // 报告阈值
        Filters     []string               // 指标过滤器
    }

    // 数据源
    collector *Collector
    
    // 报告缓存
    cache struct {
        lastReport  types.Report
        lastUpdate  time.Time
        history     []types.Report
    }

    // 订阅者
    subscribers map[string]ReportSubscriber

    // 状态
    status struct {
        isRunning  bool
        errors     []error
        lastError  time.Time
    }
}

// ReportSubscriber 报告订阅者接口
type ReportSubscriber interface {
    OnReport(report types.Report) error
    GetID() string
}

// Report 报告结构
type Report struct {
    // 基本信息
    ID        string    `json:"id"`
    Timestamp time.Time `json:"timestamp"`
    Period    string    `json:"period"`

    // 系统概览
    Summary struct {
        Status    string  `json:"status"`
        Health    float64 `json:"health"`
        Issues    int     `json:"issues"`
    } `json:"summary"`

    // 详细指标
    Metrics types.MetricsData `json:"metrics"`

    // 趋势分析
    Trends struct {
        Energy     []float64 `json:"energy"`
        Field      []float64 `json:"field"`
        Coherence  []float64 `json:"coherence"`
    } `json:"trends"`

    // 告警信息
    Alerts []types.Alert `json:"alerts"`

    // 建议actions
    Recommendations []string `json:"recommendations"`
}

// NewReporter 创建新的报告器
func NewReporter(collector *Collector, config types.MetricsConfig) *Reporter {
    r := &Reporter{
        collector:    collector,
        subscribers: make(map[string]ReportSubscriber),
    }

    // 初始化配置
    r.config.Interval = config.ReportInterval
    r.config.Format = config.ReportFormat
    r.config.OutputPath = config.OutputPath
    r.config.Thresholds = config.Thresholds
    r.config.Filters = config.Filters

    return r
}

// Start 启动报告器
func (r *Reporter) Start(ctx context.Context) error {
    r.mu.Lock()
    if r.status.isRunning {
        r.mu.Unlock()
        return types.NewSystemError(types.ErrRuntime, "reporter already running", nil)
    }
    r.status.isRunning = true
    r.mu.Unlock()

    // 启动报告循环
    go r.reportLoop(ctx)

    return nil
}

// Stop 停止报告器
func (r *Reporter) Stop() error {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.status.isRunning = false
    return nil
}

// Subscribe 订阅报告
func (r *Reporter) Subscribe(subscriber ReportSubscriber) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    id := subscriber.GetID()
    if _, exists := r.subscribers[id]; exists {
        return types.NewSystemError(types.ErrExists, "subscriber already exists", nil)
    }

    r.subscribers[id] = subscriber
    return nil
}

// Unsubscribe 取消订阅
func (r *Reporter) Unsubscribe(id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.subscribers[id]; !exists {
        return types.NewSystemError(types.ErrNotFound, "subscriber not found", nil)
    }

    delete(r.subscribers, id)
    return nil
}

// reportLoop 报告循环
func (r *Reporter) reportLoop(ctx context.Context) {
    ticker := time.NewTicker(r.config.Interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            if err := r.generateAndSendReport(ctx); err != nil {
                r.handleError(err)
            }
        }
    }
}

// generateAndSendReport 生成并发送报告
func (r *Reporter) generateAndSendReport(ctx context.Context) error {
    // 获取当前指标
    metrics := r.collector.GetCurrentMetrics()
    history := r.collector.GetMetricsHistory()

    // 生成报告
    report := r.generateReport(metrics, history)

    // 缓存报告
    r.cacheReport(report)

    // 发送给订阅者
    r.notifySubscribers(report)

    // 保存报告
    if err := r.saveReport(report); err != nil {
        return err
    }

    return nil
}

// generateReport 生成报告
func (r *Reporter) generateReport(current types.MetricsData, history []types.MetricsData) Report {
    report := Report{
        ID:        fmt.Sprintf("report-%d", time.Now().Unix()),
        Timestamp: time.Now(),
        Period:    r.config.Interval.String(),
    }

    // 设置系统概览
    report.Summary.Status = r.calculateSystemStatus(current)
    report.Summary.Health = r.calculateSystemHealth(current)
    report.Summary.Issues = r.countSystemIssues(current)

    // 设置详细指标
    report.Metrics = current

    // 计算趋势
    report.Trends.Energy = r.calculateTrend(history, "energy")
    report.Trends.Field = r.calculateTrend(history, "field")
    report.Trends.Coherence = r.calculateTrend(history, "coherence")

    // 生成建议
    report.Recommendations = r.generateRecommendations(current, history)

    return report
}

// calculateSystemStatus 计算系统状态
func (r *Reporter) calculateSystemStatus(metrics types.MetricsData) string {
    // 基于指标计算系统状态
    if metrics.System.Health < r.config.Thresholds["critical_health"] {
        return "Critical"
    } else if metrics.System.Health < r.config.Thresholds["warning_health"] {
        return "Warning"
    }
    return "Healthy"
}

// calculateSystemHealth 计算系统健康度
func (r *Reporter) calculateSystemHealth(metrics types.MetricsData) float64 {
    // 综合各项指标计算健康度
    var health float64
    
    // 能量健康度
    energyHealth := metrics.System.Energy / r.config.Thresholds["max_energy"]
    
    // 场健康度
    fieldHealth := metrics.System.Field.Strength / r.config.Thresholds["max_field_strength"]
    
    // 量子相干性健康度
    coherenceHealth := metrics.System.Quantum.Coherence
    
    // 加权平均
    health = (energyHealth*0.4 + fieldHealth*0.3 + coherenceHealth*0.3)
    
    return health
}

// calculateTrend 计算指标趋势
func (r *Reporter) calculateTrend(history []types.MetricsData, metricType string) []float64 {
    trend := make([]float64, len(history))
    
    for i, metrics := range history {
        switch metricType {
        case "energy":
            trend[i] = metrics.System.Energy
        case "field":
            trend[i] = metrics.System.Field.Strength
        case "coherence":
            trend[i] = metrics.System.Quantum.Coherence
        }
    }
    
    return trend
}

// generateRecommendations 生成建议
func (r *Reporter) generateRecommendations(current types.MetricsData, history []types.MetricsData) []string {
    var recommendations []string

    // 基于能量水平的建议
    if current.System.Energy < r.config.Thresholds["min_energy"] {
        recommendations = append(recommendations, "Increase system energy level")
    }

    // 基于场强度的建议
    if current.System.Field.Strength > r.config.Thresholds["max_field_strength"] {
        recommendations = append(recommendations, "Reduce field strength to maintain stability")
    }

    // 基于量子相干性的建议
    if current.System.Quantum.Coherence < r.config.Thresholds["min_coherence"] {
        recommendations = append(recommendations, "Enhance quantum coherence")
    }

    return recommendations
}

// saveReport 保存报告
func (r *Reporter) saveReport(report Report) error {
    if r.config.OutputPath == "" {
        return nil
    }

    // 序列化报告
    data, err := json.MarshalIndent(report, "", "  ")
    if err != nil {
        return types.WrapError(err, types.ErrInternal, "failed to marshal report")
    }

    // TODO: 实现报告存储逻辑
    // 可以存储到文件、数据库或其他存储系统

    return nil
}

// notifySubscribers 通知订阅者
func (r *Reporter) notifySubscribers(report Report) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    for _, subscriber := range r.subscribers {
        go func(s ReportSubscriber) {
            if err := s.OnReport(report); err != nil {
                r.handleError(err)
            }
        }(subscriber)
    }
}

// handleError 处理错误
func (r *Reporter) handleError(err error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.status.errors = append(r.status.errors, err)
    r.status.lastError = time.Now()
}

// cacheReport 缓存报告
func (r *Reporter) cacheReport(report Report) {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.cache.lastReport = report
    r.cache.lastUpdate = time.Now()
    r.cache.history = append(r.cache.history, report)

    // 维护缓存大小
    if len(r.cache.history) > 100 {
        r.cache.history = r.cache.history[1:]
    }
}
