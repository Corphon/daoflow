// system/monitor/metrics/reporter.go

package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Corphon/daoflow/system/types"
)

// Reporter 指标报告器
type Reporter struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		Interval   time.Duration      // 报告间隔
		Format     string             // 报告格式
		OutputPath string             // 输出路径
		Thresholds map[string]float64 // 报告阈值
		Filters    []string           // 指标过滤器
	}

	// 数据源
	collector *Collector

	// 报告缓存
	cache struct {
		lastReport types.Report
		lastUpdate time.Time
		history    []types.Report
	}

	// 订阅者
	subscribers map[string]ReportSubscriber

	// 状态
	status struct {
		isRunning bool
		errors    []error
		lastError time.Time
	}
}

// ReportSubscriber 报告订阅者接口
type ReportSubscriber interface {
	OnReport(report types.Report) error
	GetID() string
}

// ------------------------------------------------------------
// AddSubscriber 添加报告订阅者
func (r *Reporter) AddSubscriber(subscriber ReportSubscriber) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.subscribers[subscriber.GetID()] = subscriber
}

// NewReporter 创建新的报告器
func NewReporter(collector *Collector, config types.MetricsConfig) *Reporter {
	r := &Reporter{
		collector:   collector,
		subscribers: make(map[string]ReportSubscriber),
	}

	// 初始化配置
	r.config.Interval = config.Report.ReportInterval
	r.config.Format = config.Report.ReportFormat
	r.config.OutputPath = config.Report.OutputPath
	r.config.Thresholds = config.Report.Thresholds
	r.config.Filters = config.Report.Filters

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
	// 先检查上下文是否已取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// 获取当前指标
	metrics := r.collector.GetCurrentMetrics()
	historyData := r.collector.GetMetricsHistory()

	// 检查是否已取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// 转换为指针切片
	history := make([]*types.MetricsData, len(historyData))
	for i := range historyData {
		history[i] = &historyData[i]
	}

	// 生成报告
	report := r.generateReport(metrics, history)

	// 检查是否已取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// 缓存报告
	r.cacheReport(report)

	// 发送给订阅者(使用上下文)
	r.notifySubscribers(report)

	// 保存报告(使用上下文)
	if err := r.saveReport(report); err != nil {
		return err
	}

	return nil
}

// generateReport 生成报告
func (r *Reporter) generateReport(current *types.MetricsData, history []*types.MetricsData) types.Report {
	report := types.Report{
		ID:        fmt.Sprintf("report-%d", time.Now().Unix()),
		Timestamp: time.Now(),
		Period:    r.config.Interval.String(),
	}

	// 设置系统概览
	report.Summary.Status = r.calculateSystemStatus(*current) // 解引用
	report.Summary.Health = r.calculateSystemHealth(*current) // 解引用
	report.Summary.Issues = r.countSystemIssues(*current)     // 解引用

	// 设置详细指标
	report.Metrics = *current // 解引用

	// 计算趋势
	report.Trends.Energy = r.calculateTrend(history, "energy")
	report.Trends.Field = r.calculateTrend(history, "field")
	report.Trends.Coherence = r.calculateTrend(history, "coherence")

	// 生成建议
	report.Recommendations = r.generateRecommendations(*current, history) // 解引用

	return report
}

// Reporter已有定义的基础上添加方法:
func (r *Reporter) calculateSystemStatus(metrics types.MetricsData) string {
	// 根据健康度判断状态
	health := r.calculateSystemHealth(metrics)
	if health >= 0.8 {
		return "healthy"
	} else if health >= 0.5 {
		return "warning"
	}
	return "critical"
}

func (r *Reporter) calculateSystemHealth(metrics types.MetricsData) float64 {
	// 计算系统健康度
	healthFactors := map[string]float64{
		"energy":    0.3,
		"field":     0.3,
		"coherence": 0.4,
	}

	health := 0.0
	health += metrics.System.Energy * healthFactors["energy"]
	health += metrics.System.Field.GetStrength() * healthFactors["field"]
	health += metrics.System.Quantum.GetCoherence() * healthFactors["coherence"]

	return math.Min(1.0, health)
}

func (r *Reporter) countSystemIssues(metrics types.MetricsData) int {
	issues := 0

	// 检查各项指标是否超过阈值
	if metrics.System.Energy < r.config.Thresholds["min_energy"] {
		issues++
	}
	if metrics.System.Field.GetStrength() > r.config.Thresholds["max_field_strength"] {
		issues++
	}
	if metrics.System.Quantum.GetCoherence() < r.config.Thresholds["min_coherence"] {
		issues++
	}

	return issues
}

// 修改calculateTrend方法签名使用指针切片
func (r *Reporter) calculateTrend(history []*types.MetricsData, metricType string) []float64 {
	trend := make([]float64, len(history))

	for i, data := range history {
		switch metricType {
		case "energy":
			trend[i] = data.System.Energy
		case "field":
			trend[i] = data.System.Field.GetStrength()
		case "coherence":
			trend[i] = data.System.Quantum.GetCoherence()
		}
	}

	return trend
}

// generateRecommendations 生成建议
func (r *Reporter) generateRecommendations(current types.MetricsData, history []*types.MetricsData) []string {
	var recommendations []string

	// 基于当前值的建议
	if current.System.Energy < r.config.Thresholds["min_energy"] {
		recommendations = append(recommendations, "Increase system energy level")
	}

	if current.System.Field.GetStrength() > r.config.Thresholds["max_field_strength"] {
		recommendations = append(recommendations, "Reduce field strength to maintain stability")
	}

	if current.System.Quantum.GetCoherence() < r.config.Thresholds["min_coherence"] {
		recommendations = append(recommendations, "Enhance quantum coherence")
	}

	// 基于历史趋势的建议
	if len(history) > 1 {
		// 能量趋势
		energyTrend := calculateTrend(history, func(m *types.MetricsData) float64 {
			return m.System.Energy
		})
		if energyTrend < 0 {
			recommendations = append(recommendations, "Energy showing downward trend, consider energy optimization")
		}

		// 场强趋势
		fieldTrend := calculateTrend(history, func(m *types.MetricsData) float64 {
			return m.System.Field.GetStrength()
		})
		if fieldTrend > 0.5 {
			recommendations = append(recommendations, "Field strength increasing rapidly, monitor stability")
		}
	}

	return recommendations
}

// calculateTrend 计算指标趋势
func calculateTrend(history []*types.MetricsData, getter func(*types.MetricsData) float64) float64 {
	if len(history) < 2 {
		return 0
	}

	first := getter(history[0])
	last := getter(history[len(history)-1])
	return (last - first) / float64(len(history)-1)
}

// saveReport 保存报告
func (r *Reporter) saveReport(report types.Report) error {
	if r.config.OutputPath == "" {
		return nil
	}

	// 序列化报告
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return types.WrapError(err, types.ErrInternal, "failed to marshal report")
	}

	// 构建文件路径
	fileName := fmt.Sprintf("report_%s.json", report.ID)
	filePath := filepath.Join(r.config.OutputPath, fileName)

	// 确保目录存在
	if err := os.MkdirAll(r.config.OutputPath, 0755); err != nil {
		return types.WrapError(err, types.ErrStorage, "failed to create output directory")
	}

	// 写入文件
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return types.WrapError(err, types.ErrStorage, "failed to write report file")
	}

	return nil
}

// notifySubscribers 通知订阅者
func (r *Reporter) notifySubscribers(report types.Report) {
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
func (r *Reporter) cacheReport(report types.Report) {
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
