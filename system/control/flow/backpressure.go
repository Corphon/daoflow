// system/control/flow/backpressure.go

package flow

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/Corphon/daoflow/model"
)

// BackpressureManager 背压管理器
type BackpressureManager struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		sampleInterval    time.Duration // 采样间隔
		windowSize        time.Duration // 窗口大小
		pressureThreshold float64       // 压力阈值
		recoveryFactor    float64       // 恢复因子
	}

	// 背压状态
	state struct {
		pressures  map[string]*Pressure        // 压力状态
		thresholds map[string]*Threshold       // 阈值配置
		monitors   map[string]*PressureMonitor // 压力监控
		metrics    PressureMetrics             // 压力指标
	}

	// 依赖项
	scheduler *FlowScheduler
	balancer  *FlowBalancer
}

// Pressure 压力状态
type Pressure struct {
	ID         string    // 压力ID
	Source     string    // 压力源
	Level      float64   // 压力等级
	Type       string    // 压力类型
	Status     string    // 当前状态
	Trend      string    // 变化趋势
	LastUpdate time.Time // 最后更新
}

// Threshold 阈值配置
type Threshold struct {
	ID       string            // 阈值ID
	Target   string            // 目标对象
	Limits   []PressureLimit   // 压力限制
	Actions  []ThresholdAction // 阈值动作
	Priority int               // 优先级
}

// PressureLimit 压力限制
type PressureLimit struct {
	Type     string        // 限制类型
	Value    float64       // 限制值
	Duration time.Duration // 持续时间
	Action   string        // 触发动作
}

// ThresholdAction 阈值动作
type ThresholdAction struct {
	Type         string                 // 动作类型
	Parameters   map[string]interface{} // 动作参数
	Cooldown     time.Duration          // 冷却时间
	LastExecuted time.Time              // 最后执行
}

// PressureMonitor 压力监控
type PressureMonitor struct {
	ID         string            // 监控ID
	Target     string            // 监控目标
	Samples    []PressureSample  // 压力采样
	Statistics MonitorStatistics // 监控统计
	Status     string            // 监控状态
}

// PressureSample 压力采样
type PressureSample struct {
	Timestamp time.Time
	Value     float64
	Tags      map[string]string
}

// MonitorStatistics 监控统计
type MonitorStatistics struct {
	Average  float64 // 平均值
	Peak     float64 // 峰值
	Trend    float64 // 趋势
	Variance float64 // 方差
}

// PressureMetrics 压力指标
type PressureMetrics struct {
	TotalPressure     float64       // 总压力
	AverageLevel      float64       // 平均等级
	ThresholdBreaches int           // 阈值突破次数
	Recoveries        int           // 恢复次数
	History           []MetricPoint // 历史指标
}

// MetricPoint 指标点
type MetricPoint struct {
	Timestamp time.Time
	Values    map[string]float64
}

// ------------------------------------------------
// NewBackpressureManager 创建新的背压管理器
func NewBackpressureManager(
	scheduler *FlowScheduler,
	balancer *FlowBalancer) *BackpressureManager {

	bm := &BackpressureManager{
		scheduler: scheduler,
		balancer:  balancer,
	}

	// 初始化配置
	bm.config.sampleInterval = 100 * time.Millisecond
	bm.config.windowSize = 5 * time.Second
	bm.config.pressureThreshold = 0.8
	bm.config.recoveryFactor = 0.7

	// 初始化状态
	bm.state.pressures = make(map[string]*Pressure)
	bm.state.thresholds = make(map[string]*Threshold)
	bm.state.monitors = make(map[string]*PressureMonitor)
	bm.state.metrics = PressureMetrics{
		History: make([]MetricPoint, 0),
	}

	return bm
}

// Monitor 监控背压
func (bm *BackpressureManager) Monitor() error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	// 采样压力数据
	if err := bm.samplePressures(); err != nil {
		return err
	}

	// 分析压力状态
	if err := bm.analyzePressures(); err != nil {
		return err
	}

	// 执行阈值动作
	if err := bm.executeThresholdActions(); err != nil {
		return err
	}

	// 更新指标
	bm.updateMetrics()

	return nil
}

// RegisterThreshold 注册阈值配置
func (bm *BackpressureManager) RegisterThreshold(threshold *Threshold) error {
	if threshold == nil {
		return model.WrapError(nil, model.ErrCodeValidation, "nil threshold")
	}

	bm.mu.Lock()
	defer bm.mu.Unlock()

	// 验证阈值配置
	if err := bm.validateThreshold(threshold); err != nil {
		return err
	}

	// 存储阈值配置
	bm.state.thresholds[threshold.ID] = threshold

	return nil
}

// samplePressures 采样压力数据
func (bm *BackpressureManager) samplePressures() error {
	currentTime := time.Now()

	for _, monitor := range bm.state.monitors {
		// 获取压力数据
		value, err := bm.getPressureValue(monitor.Target)
		if err != nil {
			continue
		}

		// 创建采样
		sample := PressureSample{
			Timestamp: currentTime,
			Value:     value,
			Tags:      make(map[string]string),
		}

		// 添加采样
		monitor.Samples = append(monitor.Samples, sample)

		// 维护采样窗口
		bm.maintainSampleWindow(monitor)

		// 更新统计
		bm.updateMonitorStatistics(monitor)
	}

	return nil
}

// getPressureValue 获取指定目标的压力值
func (bm *BackpressureManager) getPressureValue(target string) (float64, error) {
	// 根据目标类型获取压力值
	switch {
	case strings.HasPrefix(target, "flow:"):
		// 获取流压力
		flowID := strings.TrimPrefix(target, "flow:")
		return bm.scheduler.GetFlowPressure(flowID)
	case strings.HasPrefix(target, "resource:"):
		// 获取资源压力
		resourceID := strings.TrimPrefix(target, "resource:")
		return bm.balancer.GetResourcePressure(resourceID)
	default:
		return 0, fmt.Errorf("unsupported pressure target: %s", target)
	}
}

// maintainSampleWindow 维护采样窗口
func (bm *BackpressureManager) maintainSampleWindow(monitor *PressureMonitor) {
	windowStart := time.Now().Add(-bm.config.windowSize)

	// 移除过期样本
	validSamples := make([]PressureSample, 0)
	for _, sample := range monitor.Samples {
		if sample.Timestamp.After(windowStart) {
			validSamples = append(validSamples, sample)
		}
	}
	monitor.Samples = validSamples
}

// updateMonitorStatistics 更新监控统计
func (bm *BackpressureManager) updateMonitorStatistics(monitor *PressureMonitor) {
	if len(monitor.Samples) == 0 {
		return
	}

	// 计算统计值
	var sum, peak, sumSquares float64
	for _, sample := range monitor.Samples {
		sum += sample.Value
		sumSquares += sample.Value * sample.Value
		if sample.Value > peak {
			peak = sample.Value
		}
	}

	count := float64(len(monitor.Samples))
	average := sum / count

	// 计算方差
	variance := (sumSquares / count) - (average * average)

	// 计算趋势（使用最新样本和平均值的差异）
	trend := monitor.Samples[len(monitor.Samples)-1].Value - average

	// 更新统计
	monitor.Statistics = MonitorStatistics{
		Average:  average,
		Peak:     peak,
		Trend:    trend,
		Variance: variance,
	}
}

// analyzePressures 分析压力状态
func (bm *BackpressureManager) analyzePressures() error {
	for _, pressure := range bm.state.pressures {
		// 获取相关监控
		monitor := bm.state.monitors[pressure.Source]
		if monitor == nil {
			continue
		}

		// 更新压力等级
		pressure.Level = bm.calculatePressureLevel(monitor)

		// 确定压力趋势
		pressure.Trend = bm.determinePressureTrend(monitor)

		// 更新压力状态
		bm.updatePressureStatus(pressure)
	}

	return nil
}

// calculatePressureLevel 计算压力等级
func (bm *BackpressureManager) calculatePressureLevel(monitor *PressureMonitor) float64 {
	if len(monitor.Samples) == 0 {
		return 0
	}

	// 使用监控统计计算压力等级
	level := monitor.Statistics.Average*0.5 + // 平均值权重50%
		monitor.Statistics.Peak*0.3 + // 峰值权重30%
		math.Max(0, monitor.Statistics.Trend)*0.2 // 上升趋势权重20%

	// 确保在[0,1]范围内
	return math.Max(0, math.Min(1, level))
}

// determinePressureTrend 确定压力趋势
func (bm *BackpressureManager) determinePressureTrend(monitor *PressureMonitor) string {
	if len(monitor.Samples) < 2 {
		return "stable"
	}

	// 基于趋势值判断方向
	trend := monitor.Statistics.Trend
	switch {
	case trend > 0.1:
		return "rising"
	case trend < -0.1:
		return "falling"
	default:
		return "stable"
	}
}

// updatePressureStatus 更新压力状态
func (bm *BackpressureManager) updatePressureStatus(pressure *Pressure) {
	// 基于压力等级和趋势更新状态
	switch {
	case pressure.Level >= bm.config.pressureThreshold:
		if pressure.Trend == "rising" {
			pressure.Status = "critical"
		} else {
			pressure.Status = "high"
		}
	case pressure.Level >= bm.config.pressureThreshold*0.7:
		pressure.Status = "warning"
	default:
		pressure.Status = "normal"
	}

	pressure.LastUpdate = time.Now()
}

// executeThresholdActions 执行阈值动作
func (bm *BackpressureManager) executeThresholdActions() error {
	for _, threshold := range bm.state.thresholds {
		// 检查是否需要执行动作
		if !bm.shouldExecuteActions(threshold) {
			continue
		}

		// 执行所有相关动作
		for _, action := range threshold.Actions {
			if err := bm.executeAction(action); err != nil {
				continue
			}
		}
	}

	return nil
}

// shouldExecuteActions 检查是否需要执行阈值动作
func (bm *BackpressureManager) shouldExecuteActions(threshold *Threshold) bool {
	// 获取目标压力状态
	pressure, exists := bm.state.pressures[threshold.Target]
	if !exists {
		return false
	}

	// 检查每个压力限制
	for _, limit := range threshold.Limits {
		switch limit.Type {
		case "level":
			if pressure.Level >= limit.Value {
				return true
			}
		case "duration":
			if pressure.Status != "normal" &&
				time.Since(pressure.LastUpdate) >= limit.Duration {
				return true
			}
		case "trend":
			if pressure.Trend == "rising" && pressure.Level >= limit.Value {
				return true
			}
		}
	}

	return false
}

// executeAction 执行阈值动作
func (bm *BackpressureManager) executeAction(action ThresholdAction) error {
	// 检查冷却时间
	if time.Since(action.LastExecuted) < action.Cooldown {
		return nil
	}

	// 根据动作类型执行相应操作
	switch action.Type {
	case "throttle":
		return bm.scheduler.ThrottleFlow(
			action.Parameters["flow_id"].(string),
			action.Parameters["rate"].(float64),
		)
	case "scale":
		return bm.balancer.ScaleResource(
			action.Parameters["resource_id"].(string),
			action.Parameters["factor"].(float64),
		)
	case "notify":
		// 实现通知逻辑
		return nil
	default:
		return fmt.Errorf("unsupported action type: %s", action.Type)
	}
}

// 辅助函数

func (bm *BackpressureManager) validateThreshold(threshold *Threshold) error {
	if threshold.ID == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty threshold ID")
	}

	if threshold.Target == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty threshold target")
	}

	return nil
}

func (bm *BackpressureManager) updateMetrics() {
	point := MetricPoint{
		Timestamp: time.Now(),
		Values:    make(map[string]float64),
	}

	// 计算当前指标
	point.Values["total_pressure"] = bm.calculateTotalPressure()
	point.Values["average_level"] = bm.calculateAverageLevel()

	bm.state.metrics.History = append(bm.state.metrics.History, point)

	// 限制历史记录数量
	if len(bm.state.metrics.History) > maxMetricsHistory {
		bm.state.metrics.History = bm.state.metrics.History[1:]
	}
}

// calculateTotalPressure 计算总压力值
func (bm *BackpressureManager) calculateTotalPressure() float64 {
	totalPressure := 0.0
	count := 0

	// 遍历所有压力状态计算总和
	for _, pressure := range bm.state.pressures {
		if pressure != nil {
			totalPressure += pressure.Level
			count++
		}
	}

	// 返回平均压力值
	if count > 0 {
		return totalPressure / float64(count)
	}
	return 0
}

// calculateAverageLevel 计算平均压力等级
func (bm *BackpressureManager) calculateAverageLevel() float64 {
	totalLevel := 0.0
	activeMonitors := 0

	// 遍历所有监控器计算平均等级
	for _, monitor := range bm.state.monitors {
		if len(monitor.Samples) > 0 {
			totalLevel += monitor.Statistics.Average
			activeMonitors++
		}
	}

	// 返回平均等级
	if activeMonitors > 0 {
		return math.Min(1.0, totalLevel/float64(activeMonitors))
	}
	return 0
}
