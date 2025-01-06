// api/metrics.go

package api

import (
    "context"
    "sync"
    "time"

    "github.com/Corphon/daoflow/system"
)

// MetricType 指标类型
type MetricType string

const (
    TypeGauge     MetricType = "gauge"      // 仪表盘类型(可增可减)
    TypeCounter   MetricType = "counter"    // 计数器类型(只增)
    TypeHistogram MetricType = "histogram"  // 直方图类型(分布)
    TypeSummary   MetricType = "summary"    // 摘要类型(统计)
    TypeBuffer    MetricType = "buffer"    // 缓冲区指标
)

// MetricValue 指标值
type MetricValue struct {
    Type      MetricType             `json:"type"`       // 指标类型
    Value     float64                `json:"value"`      // 当前值
    Timestamp time.Time              `json:"timestamp"`  // 时间戳
    Labels    map[string]string      `json:"labels"`     // 标签
    Metadata  map[string]interface{} `json:"metadata"`   // 元数据
}

// MetricSeries 指标序列
type MetricSeries struct {
    Name        string                `json:"name"`         // 指标名称
    Type        MetricType            `json:"type"`         // 指标类型
    Description string                `json:"description"`   // 描述
    Unit        string                `json:"unit"`         // 单位
    Labels      map[string]string     `json:"labels"`       // 标签
    Values      []*MetricValue        `json:"values"`       // 历史值
    Statistics  *MetricStats          `json:"statistics"`   // 统计信息
}

// MetricStats 指标统计
type MetricStats struct {
    Min      float64   `json:"min"`       // 最小值
    Max      float64   `json:"max"`       // 最大值
    Sum      float64   `json:"sum"`       // 总和
    Count    int64     `json:"count"`     // 计数
    Average  float64   `json:"average"`   // 平均值
    Variance float64   `json:"variance"`  // 方差
    LastUpdate time.Time `json:"last_update"` // 最后更新时间
}

// MetricsAPI 指标监控API
type MetricsAPI struct {
    mu     sync.RWMutex
    system *system.SystemCore
    
    // 指标存储
    metrics map[string]*MetricSeries
    
    // 聚合设置
    aggregationInterval time.Duration
    retentionPeriod    time.Duration
    
    // 事件通道
    events chan MetricEvent
    
    ctx    context.Context
    cancel context.CancelFunc

    // MetricsAPI 添加缓冲区指标支持
    bufferStats map[string]*BufferStats
}

// MetricEvent 指标事件
type MetricEvent struct {
    Type      string        `json:"type"`       // 事件类型
    Series    string        `json:"series"`     // 指标序列
    Value     *MetricValue  `json:"value"`      // 指标值
    Timestamp time.Time     `json:"timestamp"`  // 事件时间
}

// MetricsConfig 指标配置
type MetricsConfig struct {
    AggregationInterval time.Duration     // 聚合间隔
    RetentionPeriod    time.Duration     // 保留期限
    DefaultLabels      map[string]string // 默认标签
}

// NewMetricsAPI 创建指标API实例
func NewMetricsAPI(sys *system.SystemCore, config MetricsConfig) *MetricsAPI {
    ctx, cancel := context.WithCancel(context.Background())
    
    api := &MetricsAPI{
        system:              sys,
        metrics:            make(map[string]*MetricSeries),
        aggregationInterval: config.AggregationInterval,
        retentionPeriod:    config.RetentionPeriod,
        events:             make(chan MetricEvent, 100),
        ctx:                ctx,
        cancel:             cancel,
    }
    
    go api.aggregate()
    go api.cleanup()
    return api
}

// RegisterMetric 注册新指标
func (m *MetricsAPI) RegisterMetric(name string, typ MetricType, desc string, unit string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if _, exists := m.metrics[name]; exists {
        return NewError(ErrMetricExists, "metric already exists")
    }

    series := &MetricSeries{
        Name:        name,
        Type:        typ,
        Description: desc,
        Unit:        unit,
        Labels:      make(map[string]string),
        Values:      make([]*MetricValue, 0),
        Statistics:  &MetricStats{},
    }

    m.metrics[name] = series
    return nil
}

// RecordMetric 记录指标值
func (m *MetricsAPI) RecordMetric(name string, value float64, labels map[string]string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    series, exists := m.metrics[name]
    if !exists {
        return NewError(ErrMetricNotFound, "metric not found")
    }

    metricValue := &MetricValue{
        Type:      series.Type,
        Value:     value,
        Timestamp: time.Now(),
        Labels:    labels,
    }

    series.Values = append(series.Values, metricValue)
    m.updateStats(series, metricValue)

    m.events <- MetricEvent{
        Type:      "metric_recorded",
        Series:    name,
        Value:     metricValue,
        Timestamp: time.Now(),
    }

    return nil
}

// GetMetric 获取指标序列
func (m *MetricsAPI) GetMetric(name string) (*MetricSeries, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    series, exists := m.metrics[name]
    if !exists {
        return nil, NewError(ErrMetricNotFound, "metric not found")
    }

    return series, nil
}

// GetMetricValue 获取指标当前值
func (m *MetricsAPI) GetMetricValue(name string) (*MetricValue, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    series, exists := m.metrics[name]
    if !exists {
        return nil, NewError(ErrMetricNotFound, "metric not found")
    }

    if len(series.Values) == 0 {
        return nil, NewError(ErrMetricNoValue, "no values recorded")
    }

    return series.Values[len(series.Values)-1], nil
}

// QueryMetrics 查询指标数据
func (m *MetricsAPI) QueryMetrics(query map[string]interface{}) ([]*MetricSeries, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    var results []*MetricSeries

    // 实现查询逻辑
    for _, series := range m.metrics {
        if m.matchQuery(series, query) {
            results = append(results, series)
        }
    }

    return results, nil
}

// Subscribe 订阅指标事件
func (m *MetricsAPI) Subscribe() (<-chan MetricEvent, error) {
    return m.events, nil
}

// aggregate 聚合计算协程
func (m *MetricsAPI) aggregate() {
    ticker := time.NewTicker(m.aggregationInterval)
    defer ticker.Stop()

    for {
        select {
        case <-m.ctx.Done():
            return
        case <-ticker.C:
            m.aggregateMetrics()
        }
    }
}

// cleanup 清理过期数据协程
func (m *MetricsAPI) cleanup() {
    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-m.ctx.Done():
            return
        case <-ticker.C:
            m.cleanupMetrics()
        }
    }
}

// updateStats 更新统计信息
func (m *MetricsAPI) updateStats(series *MetricSeries, value *MetricValue) {
    stats := series.Statistics
    stats.LastUpdate = time.Now()
    stats.Count++
    stats.Sum += value.Value

    if stats.Count == 1 {
        stats.Min = value.Value
        stats.Max = value.Value
        stats.Average = value.Value
        return
    }

    if value.Value < stats.Min {
        stats.Min = value.Value
    }
    if value.Value > stats.Max {
        stats.Max = value.Value
    }

    // 更新平均值和方差
    oldAvg := stats.Average
    stats.Average = stats.Sum / float64(stats.Count)
    stats.Variance = (stats.Variance*(float64(stats.Count-1)) + 
        (value.Value-oldAvg)*(value.Value-stats.Average)) / float64(stats.Count)
}

// aggregateMetrics 聚合指标数据
func (m *MetricsAPI) aggregateMetrics() {
    m.mu.Lock()
    defer m.mu.Unlock()

    now := time.Now()
    for _, series := range m.metrics {
        // 实现聚合逻辑
    }
}

// cleanupMetrics 清理过期数据
func (m *MetricsAPI) cleanupMetrics() {
    m.mu.Lock()
    defer m.mu.Unlock()

    cutoff := time.Now().Add(-m.retentionPeriod)
    for _, series := range m.metrics {
        // 清理过期数据
        var validValues []*MetricValue
        for _, value := range series.Values {
            if value.Timestamp.After(cutoff) {
                validValues = append(validValues, value)
            }
        }
        series.Values = validValues
    }
}

// matchQuery 检查指标是否匹配查询条件
func (m *MetricsAPI) matchQuery(series *MetricSeries, query map[string]interface{}) bool {
    // 实现查询匹配逻辑
    return true
}

// Close 关闭API
func (m *MetricsAPI) Close() error {
    m.cancel()
    close(m.events)
    return nil
}
