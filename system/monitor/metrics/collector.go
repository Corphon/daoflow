// system/monitor/metrics/collector.go

package metrics

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/types"
)

type GenerationMetrics struct {
	TotalGenerated int
	SuccessRate    float64
	AverageScore   float64
	ComplexityDist map[float64]int
	Evolution      []types.MetricPoint
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

	// 采集状态
	state struct {
		totalSamples   int64     // 总采样数
		droppedSamples int64     // 丢弃采样数
		lastCollection time.Time // 最后采集时间
		metrics        struct {
			AverageLatency time.Duration // 平均延迟
			MaxLatency     time.Duration // 最大延迟
			SuccessRate    float64       // 成功率
		}
	}

	// 样本缓冲
	samples []types.MetricsData

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
		history:       make([]types.MetricsData, 0, config.Base.HistorySize),
	}
}

// ---------------------------------------------------
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
func (c *Collector) GetCurrentMetrics() *types.MetricsData {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return &c.current
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
	ticker := time.NewTicker(c.config.Base.Interval)
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

		// 使用指针避免复制锁
		switch typ {
		case types.MetricEnergy:
			metrics.System.Energy = data.(float64)
		case types.MetricField:
			field := data.(*types.FieldState)
			metrics.System.Field = field
		case types.MetricQuantum:
			quantum := data.(*types.QuantumState)
			metrics.System.Quantum = quantum
		case types.MetricEmergence:
			emergence := data.(types.EmergentProperty)
			metrics.System.Emergence = &emergence
		}
	}

	// 使用指针复制
	c.current = metrics
	c.history = append(c.history, metrics)

	if len(c.history) > c.config.Base.HistorySize {
		c.history = c.history[1:]
	}

	// 传递指针避免复制
	c.checkThresholds(metrics)

	c.status.lastRun = time.Now()
	return nil
}

// checkThresholds 检查指标阈值
func (c *Collector) checkThresholds(metrics types.MetricsData) {
	// 检查系统能量
	if metrics.System.Energy < c.config.Base.Thresholds["min_energy"] {
		c.notify(types.Alert{
			Level:   types.AlertLevelWarning,
			Type:    "energy_low",
			Message: "System energy below minimum threshold",
			Time:    time.Now(),
		})
	}

	// 检查场强度
	if field := metrics.System.Field; field.GetStrength() > c.config.Base.Thresholds["max_field_strength"] {
		c.notify(types.Alert{
			Level:   types.AlertLevelWarning,
			Type:    "field_high",
			Message: "Field strength exceeds maximum threshold",
			Time:    time.Now(),
		})
	}

	// 检查量子相干性
	if quantum := metrics.System.Quantum; quantum.GetCoherence() < c.config.Base.Thresholds["min_coherence"] {
		c.notify(types.Alert{
			Level:   types.AlertLevelWarning,
			Type:    "coherence_low",
			Message: "Quantum coherence below minimum threshold",
			Time:    time.Now(),
		})
	}
}

// handleError 处理收集过程中的错误
func (c *Collector) handleError(err error) {
	c.notify(types.Alert{
		Level:   types.AlertSeverityToLevel(types.SeverityError),
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

// GetModelMetrics 获取模型指标
func (c *Collector) GetModelMetrics() (model.ModelMetrics, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	metrics := model.ModelMetrics{}

	// 设置能量指标
	metrics.Energy.Total = c.current.System.Energy
	metrics.Energy.Average = c.calculateAverageEnergy()
	metrics.Energy.Variance = c.calculateEnergyVariance()

	// 设置性能指标
	metrics.Performance.Throughput = c.calculateThroughput()
	metrics.Performance.QPS = c.calculateQPS()
	metrics.Performance.ErrorRate = c.calculateErrorRate()

	// 设置状态指标
	metrics.State.Stability = c.calculateStability()
	metrics.State.Transitions = c.countStateTransitions()
	metrics.State.Uptime = c.calculateUptime()

	return metrics, nil
}

// calculateAverageEnergy 计算平均能量
func (c *Collector) calculateAverageEnergy() float64 {
	total := 0.0
	count := 0.0
	for _, data := range c.history {
		total += data.System.Energy
		count++
	}
	if count == 0 {
		return 0
	}
	return total / count
}

// calculateEnergyVariance 计算能量方差
func (c *Collector) calculateEnergyVariance() float64 {
	avg := c.calculateAverageEnergy()
	total := 0.0
	count := 0.0
	for _, data := range c.history {
		diff := data.System.Energy - avg
		total += diff * diff
		count++
	}
	if count == 0 {
		return 0
	}
	return total / count
}

// calculateThroughput 计算吞吐量
func (c *Collector) calculateThroughput() float64 {
	if len(c.history) < 2 {
		return 0
	}
	duration := c.history[len(c.history)-1].Timestamp.Sub(c.history[0].Timestamp)
	return float64(len(c.history)) / duration.Seconds()
}

// calculateQPS 计算每秒查询数
func (c *Collector) calculateQPS() float64 {
	return c.calculateThroughput()
}

// calculateErrorRate 计算错误率
func (c *Collector) calculateErrorRate() float64 {
	total := 0
	errors := 0
	for _, data := range c.history {
		if data.Status == types.StatusError {
			errors++
		}
		total++
	}
	if total == 0 {
		return 0
	}
	return float64(errors) / float64(total)
}

// calculateStability 计算稳定性
func (c *Collector) calculateStability() float64 {
	if len(c.history) < 2 {
		return 1.0
	}

	var variance float64
	prev := c.history[0].System.Energy
	for _, data := range c.history[1:] {
		diff := data.System.Energy - prev
		variance += diff * diff
		prev = data.System.Energy
	}
	variance /= float64(len(c.history) - 1)

	return 1.0 / (1.0 + math.Sqrt(variance))
}

// countStateTransitions 计算状态转换次数
func (c *Collector) countStateTransitions() int {
	if len(c.history) < 2 {
		return 0
	}

	transitions := 0
	prev := c.history[0].Status
	for _, data := range c.history[1:] {
		if data.Status != prev {
			transitions++
		}
		prev = data.Status
	}
	return transitions
}

// calculateUptime 计算运行时间
func (c *Collector) calculateUptime() float64 {
	if len(c.history) == 0 {
		return 0
	}
	duration := time.Since(c.history[0].Timestamp)
	return duration.Seconds()
}

// GetMetricsMap 实现 MetricsSource 接口
func (c *Collector) GetMetricsMap() (map[string]float64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	metrics := make(map[string]float64)

	// 基础系统指标
	metrics["energy"] = c.current.System.Energy
	metrics["field_strength"] = c.current.System.Field.GetStrength()
	metrics["coherence"] = c.current.System.Quantum.GetCoherence()

	// 性能指标
	metrics["collection_rate"] = float64(c.state.totalSamples) / float64(c.state.droppedSamples+1)
	metrics["avg_latency"] = float64(c.state.metrics.AverageLatency.Milliseconds())
	metrics["max_latency"] = float64(c.state.metrics.MaxLatency.Milliseconds())
	metrics["success_rate"] = c.state.metrics.SuccessRate

	// 资源指标
	metrics["memory_usage"] = float64(c.getMemoryUsage())
	metrics["goroutines"] = float64(c.getGoroutineCount())

	return metrics, nil
}

// GetMetric 实现 MetricsSource 接口
func (c *Collector) GetMetric(name string) (float64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 从当前指标中获取值
	switch name {
	case "energy":
		return c.current.System.Energy, nil
	case "field_strength":
		return c.current.System.Field.GetStrength(), nil
	case "coherence":
		return c.current.System.Quantum.GetCoherence(), nil
	default:
		if val, ok := c.current.Custom[name].(float64); ok {
			return val, nil
		}
		return 0, fmt.Errorf("metric %s not found", name)
	}
}

// GetMetrics implements alert.MetricsSource
func (c *Collector) GetMetrics() (map[string]float64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	metrics := make(map[string]float64)

	// 基础系统指标
	metrics["energy"] = c.current.System.Energy
	metrics["field_strength"] = c.current.System.Field.GetStrength()
	metrics["coherence"] = c.current.System.Quantum.GetCoherence()

	// 性能指标
	metrics["collection_rate"] = float64(c.state.totalSamples) / float64(c.state.droppedSamples+1)
	metrics["avg_latency"] = float64(c.state.metrics.AverageLatency.Milliseconds())
	metrics["max_latency"] = float64(c.state.metrics.MaxLatency.Milliseconds())
	metrics["success_rate"] = c.state.metrics.SuccessRate

	// 资源指标
	metrics["memory_usage"] = float64(c.getMemoryUsage())
	metrics["goroutines"] = float64(c.getGoroutineCount())

	return metrics, nil
}

// GetMetricsData 获取收集器指标
func (c *Collector) GetMetricsData() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	metrics := make(map[string]interface{})
	// 基础指标
	metrics["samples_collected"] = c.state.totalSamples
	metrics["samples_dropped"] = c.state.droppedSamples
	metrics["collection_rate"] = float64(c.state.totalSamples) / float64(c.state.droppedSamples+1)
	metrics["last_collection"] = c.state.lastCollection
	metrics["buffer_size"] = len(c.samples)

	// 性能指标
	metrics["avg_latency"] = c.state.metrics.AverageLatency.Milliseconds()
	metrics["max_latency"] = c.state.metrics.MaxLatency.Milliseconds()
	metrics["success_rate"] = c.state.metrics.SuccessRate

	// 资源使用
	metrics["memory_usage"] = c.getMemoryUsage()
	metrics["goroutines"] = c.getGoroutineCount()

	return metrics
}

// getMemoryUsage 获取内存使用情况
func (c *Collector) getMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc / 1024 / 1024 // 转换为MB
}

// getGoroutineCount 获取协程数量
func (c *Collector) getGoroutineCount() int {
	return runtime.NumGoroutine()
}
