//model/analyzer.go

package model

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

const (
	trendThreshold           = 0.1 // 趋势判断阈值
	stateTransitionThreshold = 0.1 // 状态转换阈值
	minTransitionInterval    = 1.0 // 最小转换间隔(秒)
	lowValueThreshold        = 0.2 // 低值阈值
	highValueThreshold       = 0.8 // 高值阈值
)

// Analyzer 模型分析器
type Analyzer struct {
	mu sync.RWMutex

	// 配置
	config struct {
		SampleRate    float64       // 采样率
		WindowSize    time.Duration // 窗口大小
		MaxPatterns   int           // 最大模式数
		MinConfidence float64       // 最小置信度
	}

	// 分析缓存
	cache struct {
		patterns  []FlowPattern // 模式缓存
		metrics   ModelMetrics  // 指标缓存
		anomalies []Anomaly     // 异常缓存
	}

	// 分析状态
	status struct {
		lastAnalysis  time.Time // 最后分析时间
		totalAnalyzed int       // 总分析次数
	}
}

// StatePredictor 状态预测器
type StatePredictor struct {
	history []ModelState
}

// PatternMetrics 模式指标
type PatternMetrics struct {
	Frequency  float64       // 出现频率
	Strength   float64       // 模式强度
	Confidence float64       // 置信度
	Duration   time.Duration // 持续时间
	Energy     float64       // 能量水平
	Stability  float64       // 稳定性
}

// TimeSeriesPoint 时间序列点
type TimeSeriesPoint struct {
	Timestamp time.Time
	Value     float64
	Type      string
	Metadata  map[string]interface{}
}

// TimeSeries 时间序列
type TimeSeries struct {
	ID        string
	Points    []TimeSeriesPoint
	Type      string
	StartTime time.Time
	EndTime   time.Time
}

type Span interface {
	GetID() string
	GetStartTime() time.Time
	GetEndTime() time.Time
	GetDuration() time.Duration
	GetMetrics() map[string]float64
}

// -------------------------------------------------------
// NewAnalyzer 创建新的模型分析器
func NewAnalyzer() *Analyzer {
	a := &Analyzer{}

	// 初始化配置
	a.config.SampleRate = 0.1           // 默认采样率10%
	a.config.WindowSize = 1 * time.Hour // 默认1小时窗口
	a.config.MaxPatterns = 100          // 最多保存100个模式
	a.config.MinConfidence = 0.6        // 最小置信度0.6

	// 初始化缓存
	a.cache.patterns = make([]FlowPattern, 0)
	a.cache.metrics = ModelMetrics{}
	a.cache.anomalies = make([]Anomaly, 0)

	// 初始化状态
	a.status.lastAnalysis = time.Now()
	a.status.totalAnalyzed = 0

	return a
}

// DetectPatterns 检测模型模式
func (a *Analyzer) DetectPatterns(spans interface{}) []FlowPattern {
	a.mu.Lock()
	defer a.mu.Unlock()

	patterns := make([]FlowPattern, 0)
	timespan := time.Since(a.status.lastAnalysis)

	// 根据时间间隔动态调整采样率
	if timespan > a.config.WindowSize {
		// 如果间隔超过窗口大小，增加采样率以获取更多数据
		a.config.SampleRate = math.Min(1.0, a.config.SampleRate*1.5)
	} else {
		// 恢复默认采样率
		a.config.SampleRate = 0.1
	}

	// 1. 提取时间序列特征
	timeSeries := extractTimeSeries(spans)

	// 2. 检测基本模式
	for _, series := range timeSeries {
		// 检测周期性模式
		if pattern := detectCyclicPattern(series); pattern != nil {
			patterns = append(patterns, *pattern)
		}

		// 检测趋势性模式
		if pattern := detectTrendPattern(series); pattern != nil {
			patterns = append(patterns, *pattern)
		}

		// 检测状态转换模式
		if pattern := detectTransitionPattern(series); pattern != nil {
			patterns = append(patterns, *pattern)
		}
	}

	// 3. 过滤和合并模式
	patterns = filterPatterns(patterns, a.config.MinConfidence)
	patterns = mergePatterns(patterns)

	// 4. 更新缓存
	if len(patterns) > a.config.MaxPatterns {
		patterns = patterns[:a.config.MaxPatterns]
	}
	a.cache.patterns = patterns
	a.status.lastAnalysis = time.Now()
	a.status.totalAnalyzed++

	return patterns
}

// extractTimeSeries 从数据中提取时间序列
func extractTimeSeries(spans interface{}) []TimeSeries {
	series := make([]TimeSeries, 0)

	switch s := spans.(type) {
	case []*Span:
		// 按指标类型分组
		metrics := groupMetricsByType(s)

		// 为每种指标创建时间序列
		for metricType, points := range metrics {
			ts := TimeSeries{
				ID:        generateTimeSeriesID(metricType),
				Points:    points,
				Type:      metricType,
				StartTime: points[0].Timestamp,
				EndTime:   points[len(points)-1].Timestamp,
			}
			series = append(series, ts)
		}
	}

	return series
}

// groupMetricsByType 按指标类型分组
func groupMetricsByType(spans []*Span) map[string][]TimeSeriesPoint {
	metrics := make(map[string][]TimeSeriesPoint)

	for _, span := range spans {
		// 使用接口方法访问
		spanMetrics := (*span).GetMetrics()
		startTime := (*span).GetStartTime()
		spanID := (*span).GetID()
		duration := (*span).GetDuration()

		for metricType, value := range spanMetrics {
			point := TimeSeriesPoint{
				Timestamp: startTime,
				Value:     value,
				Type:      metricType,
				Metadata: map[string]interface{}{
					"spanID":   spanID,
					"duration": duration,
				},
			}
			metrics[metricType] = append(metrics[metricType], point)
		}
	}

	// 对每个指标序列按时间排序
	for _, points := range metrics {
		sort.Slice(points, func(i, j int) bool {
			return points[i].Timestamp.Before(points[j].Timestamp)
		})
	}

	return metrics
}

// CalculateMetrics 计算模型指标
func (a *Analyzer) CalculateMetrics(spans interface{}) ModelMetrics {
	a.mu.Lock()
	defer a.mu.Unlock()

	metrics := ModelMetrics{}

	// 基础指标
	metrics.Basic.TotalSpans = countSpans(spans)
	metrics.Basic.ErrorRate = calculateErrorRate(spans)
	metrics.Basic.Latency = calculateLatency(spans)

	// 能量指标
	metrics.Energy.Total = calculateTotalEnergy(spans)
	metrics.Energy.Average = calculateAverageEnergy(spans)
	metrics.Energy.Variance = calculateEnergyVariance(spans)

	// 状态指标
	metrics.State.Transitions = countStateTransitions(spans)
	metrics.State.Stability = calculateStateStability(spans)
	metrics.State.Uptime = calculateUptime(spans)

	// 性能指标
	metrics.Performance.Throughput = calculateThroughput(spans)
	metrics.Performance.Latency = calculateLatency(spans)
	metrics.Performance.ErrorRate = calculateErrorRate(spans)
	metrics.Performance.QPS = calculateQPS(spans)

	// 更新缓存
	a.cache.metrics = metrics

	return metrics
}

// 辅助函数
func countSpans(spans interface{}) int {
	if spanArray, ok := spans.([]*Span); ok {
		return len(spanArray)
	}
	return 0
}

func calculateUptime(spans interface{}) float64 {
	if spanArray, ok := spans.([]*Span); ok && len(spanArray) > 0 {
		first := (*spanArray[0]).GetStartTime()
		last := (*spanArray[len(spanArray)-1]).GetEndTime()
		return last.Sub(first).Seconds()
	}
	return 0
}

func calculateQPS(spans interface{}) float64 {
	if spanArray, ok := spans.([]*Span); ok && len(spanArray) > 0 {
		duration := (*spanArray[len(spanArray)-1]).GetEndTime().Sub((*spanArray[0]).GetStartTime())
		if duration > 0 {
			return float64(len(spanArray)) / duration.Seconds()
		}
	}
	return 0
}

// DetectAnomalies 检测模型异常
func (a *Analyzer) DetectAnomalies(spans interface{}) []Anomaly {
	a.mu.Lock()
	defer a.mu.Unlock()

	anomalies := make([]Anomaly, 0)
	metrics := a.cache.metrics

	// 1. 检测能量异常
	if energyAnomalies := detectEnergyAnomalies(spans, metrics.Energy); len(energyAnomalies) > 0 {
		anomalies = append(anomalies, energyAnomalies...)
	}

	// 2. 检测状态异常
	if stateAnomalies := detectStateAnomalies(spans, metrics.State); len(stateAnomalies) > 0 {
		anomalies = append(anomalies, stateAnomalies...)
	}

	// 3. 检测性能异常
	if perfAnomalies := detectPerformanceAnomalies(spans, metrics.Performance); len(perfAnomalies) > 0 {
		anomalies = append(anomalies, perfAnomalies...)
	}

	// 更新缓存
	a.cache.anomalies = anomalies

	return anomalies
}

// detectEnergyAnomalies 检测能量异常
func detectEnergyAnomalies(spans interface{}, energy Energy) []Anomaly {
	anomalies := make([]Anomaly, 0)

	// 检查平均能量异常
	avgEnergy := calculateAverageEnergy(spans)
	if math.Abs(avgEnergy-energy.Average) > energy.Variance*2 {
		anomalies = append(anomalies, Anomaly{
			Type:      "energy",
			Subtype:   "average",
			Severity:  math.Abs(avgEnergy-energy.Average) / energy.Average,
			Value:     avgEnergy,
			Expected:  energy.Average,
			Threshold: energy.Variance * 2,
			Time:      time.Now(),
		})
	}

	// 检查能量波动异常
	currentVariance := calculateEnergyVariance(spans)
	if currentVariance > energy.Variance*3 {
		anomalies = append(anomalies, Anomaly{
			Type:      "energy",
			Subtype:   "variance",
			Severity:  currentVariance / energy.Variance,
			Value:     currentVariance,
			Expected:  energy.Variance,
			Threshold: energy.Variance * 3,
			Time:      time.Now(),
		})
	}

	return anomalies
}

// detectStateAnomalies 检测状态异常
func detectStateAnomalies(spans interface{}, state State) []Anomaly {
	anomalies := make([]Anomaly, 0)

	// 检查状态转换频率异常
	transitions := countStateTransitions(spans)
	expectedTransitions := float64(state.Transitions) * state.Stability
	if float64(transitions) > expectedTransitions*2 {
		anomalies = append(anomalies, Anomaly{
			Type:      "state",
			Subtype:   "transitions",
			Severity:  float64(transitions) / expectedTransitions,
			Value:     float64(transitions),
			Expected:  expectedTransitions,
			Threshold: expectedTransitions * 2,
			Time:      time.Now(),
		})
	}

	// 检查稳定性异常
	stability := calculateStateStability(spans)
	if stability < state.Stability/2 {
		anomalies = append(anomalies, Anomaly{
			Type:      "state",
			Subtype:   "stability",
			Severity:  (state.Stability - stability) / state.Stability,
			Value:     stability,
			Expected:  state.Stability,
			Threshold: state.Stability / 2,
			Time:      time.Now(),
		})
	}

	return anomalies
}

// detectPerformanceAnomalies 检测性能异常
func detectPerformanceAnomalies(spans interface{}, perf Performance) []Anomaly {
	anomalies := make([]Anomaly, 0)

	// 检查吞吐量异常
	throughput := calculateThroughput(spans)
	if throughput < perf.Throughput/2 {
		anomalies = append(anomalies, Anomaly{
			Type:      "performance",
			Subtype:   "throughput",
			Severity:  (perf.Throughput - throughput) / perf.Throughput,
			Value:     throughput,
			Expected:  perf.Throughput,
			Threshold: perf.Throughput / 2,
			Time:      time.Now(),
		})
	}

	// 检查延迟异常
	latency := calculateLatency(spans)
	if latency > perf.Latency*2 {
		anomalies = append(anomalies, Anomaly{
			Type:      "performance",
			Subtype:   "latency",
			Severity:  latency / perf.Latency,
			Value:     latency,
			Expected:  perf.Latency,
			Threshold: perf.Latency * 2,
			Time:      time.Now(),
		})
	}

	// 检查错误率异常
	errorRate := calculateErrorRate(spans)
	if errorRate > perf.ErrorRate*2 {
		anomalies = append(anomalies, Anomaly{
			Type:      "performance",
			Subtype:   "error_rate",
			Severity:  errorRate / perf.ErrorRate,
			Value:     errorRate,
			Expected:  perf.ErrorRate,
			Threshold: perf.ErrorRate * 2,
			Time:      time.Now(),
		})
	}

	return anomalies
}

// NewStatePredictor 创建状态预测器
func NewStatePredictor() *StatePredictor {
	return &StatePredictor{
		history: make([]ModelState, 0),
	}
}

// PredictNext 预测下一个状态
func (sp *StatePredictor) PredictNext(metrics ModelMetrics) (ModelState, error) {
	// 根据转换次数预测下一个相位
	var nextPhase ProcessPhase
	switch metrics.State.Transitions % 4 {
	case 0:
		nextPhase = ProcessPhaseInitial
	case 1:
		nextPhase = ProcessPhaseTransform
	case 2:
		nextPhase = ProcessPhaseStable
	case 3:
		nextPhase = ProcessPhaseComplete
	default:
		nextPhase = ProcessPhaseNone
	}

	nextState := ModelState{
		Energy:     metrics.Energy.Total * (1 + metrics.Energy.Average/100),
		Phase:      Phase(nextPhase), // 转换为基础Phase类型
		Nature:     NatureNeutral,
		UpdateTime: time.Now(),
	}
	return nextState, nil
}

// generateTimeSeriesID 生成时间序列ID
func generateTimeSeriesID(metricType string) string {
	return fmt.Sprintf("ts_%s_%d", metricType, time.Now().UnixNano())
}

// detectCyclicPattern 检测周期性模式
func detectCyclicPattern(series TimeSeries) *FlowPattern {
	if len(series.Points) < 4 {
		return nil
	}

	// 计算自相关性来检测周期
	periods := detectPeriods(series.Points)
	if len(periods) == 0 {
		return nil
	}

	// 创建周期性模式
	return &FlowPattern{
		ID:   generatePatternID(),
		Type: "cyclic",
		Metrics: PatternMetrics{
			Frequency:  calculateFrequency(periods),
			Strength:   calculateCyclicStrength(series.Points, periods),
			Confidence: calculateCyclicConfidence(series.Points, periods),
			Duration:   series.EndTime.Sub(series.StartTime),
		},
		Properties: map[string]interface{}{
			"periods": periods,
			"phases":  detectPhases(series.Points, periods[0]),
		},
		Created: time.Now(),
	}
}

// generatePatternID 生成模式ID
func generatePatternID() string {
	return fmt.Sprintf("pattern_%d", time.Now().UnixNano())
}

// calculateFrequency 计算周期频率
func calculateFrequency(periods []float64) float64 {
	if len(periods) == 0 {
		return 0
	}
	// 使用最显著的周期计算频率
	mainPeriod := periods[0]
	if mainPeriod > 0 {
		return 1.0 / mainPeriod
	}
	return 0
}

// calculateCyclicStrength 计算周期强度
func calculateCyclicStrength(points []TimeSeriesPoint, periods []float64) float64 {
	if len(points) < 2 || len(periods) == 0 {
		return 0
	}

	// 计算周期波动的一致性
	mainPeriod := periods[0]
	totalDeviation := 0.0
	cycleCount := 0

	for i := 0; i < len(points)-1; i++ {
		if float64(points[i+1].Timestamp.Sub(points[i].Timestamp).Seconds()) >= mainPeriod {
			deviation := math.Abs(points[i+1].Value - points[i].Value)
			totalDeviation += deviation
			cycleCount++
		}
	}

	if cycleCount > 0 {
		averageDeviation := totalDeviation / float64(cycleCount)
		// 归一化强度到 [0,1] 范围
		return math.Max(0, 1-averageDeviation/math.Max(1e-6, getValueRange(points)))
	}
	return 0
}

// calculateCyclicConfidence 计算周期置信度
func calculateCyclicConfidence(points []TimeSeriesPoint, periods []float64) float64 {
	if len(points) < 4 || len(periods) == 0 {
		return 0
	}

	// 计算周期预测准确性
	mainPeriod := periods[0]
	totalError := 0.0
	predictions := 0

	for i := 0; i < len(points)-int(mainPeriod); i++ {
		predicted := points[i].Value
		actual := points[i+int(mainPeriod)].Value
		error := math.Abs(predicted - actual)
		totalError += error
		predictions++
	}

	if predictions > 0 {
		averageError := totalError / float64(predictions)
		// 归一化置信度到 [0,1] 范围
		return math.Max(0, 1-averageError/math.Max(1e-6, getValueRange(points)))
	}
	return 0
}

// detectPhases 检测周期相位
func detectPhases(points []TimeSeriesPoint, period float64) []string {
	if len(points) < 2 || period <= 0 {
		return nil
	}

	phases := make([]string, 0)
	periodSeconds := int64(period)

	for i := 0; i < len(points); i++ {
		timeInPeriod := points[i].Timestamp.Unix() % periodSeconds
		phase := ""

		// 将周期分为4个相位
		switch {
		case timeInPeriod < periodSeconds/4:
			phase = "rising"
		case timeInPeriod < periodSeconds/2:
			phase = "peak"
		case timeInPeriod < 3*periodSeconds/4:
			phase = "falling"
		default:
			phase = "trough"
		}

		// 只添加相位变化点
		if len(phases) == 0 || phases[len(phases)-1] != phase {
			phases = append(phases, phase)
		}
	}

	return phases
}

// getValueRange 获取数值范围
func getValueRange(points []TimeSeriesPoint) float64 {
	if len(points) == 0 {
		return 0
	}

	min, max := points[0].Value, points[0].Value
	for _, p := range points {
		if p.Value < min {
			min = p.Value
		}
		if p.Value > max {
			max = p.Value
		}
	}
	return max - min
}

// detectTrendPattern 检测趋势性模式
func detectTrendPattern(series TimeSeries) *FlowPattern {
	if len(series.Points) < 3 {
		return nil
	}

	// 计算趋势特征
	slope, r2 := calculateTrendLine(series.Points)
	if math.Abs(slope) < 0.1 || r2 < 0.6 {
		return nil
	}

	// 创建趋势性模式
	return &FlowPattern{
		ID:   generatePatternID(),
		Type: "trend",
		Metrics: PatternMetrics{
			Frequency:  1.0,
			Strength:   math.Abs(slope),
			Confidence: r2,
			Duration:   series.EndTime.Sub(series.StartTime),
		},
		Properties: map[string]interface{}{
			"slope":     slope,
			"direction": getTrendDirection(slope),
			"stability": r2,
		},
		Created: time.Now(),
	}
}

// getTrendDirection 获取趋势方向
func getTrendDirection(slope float64) string {
	switch {
	case slope > trendThreshold:
		return "rising"
	case slope < -trendThreshold:
		return "falling"
	default:
		return "stable"
	}
}

// detectTransitionPattern 检测状态转换模式
func detectTransitionPattern(series TimeSeries) *FlowPattern {
	if len(series.Points) < 2 {
		return nil
	}

	// 检测状态变化点
	transitions := detectStateTransitions(series.Points)
	if len(transitions) == 0 {
		return nil
	}

	// 创建状态转换模式
	return &FlowPattern{
		ID:   generatePatternID(),
		Type: "transition",
		Metrics: PatternMetrics{
			Frequency:  float64(len(transitions)) / series.EndTime.Sub(series.StartTime).Hours(),
			Strength:   calculateTransitionStrength(transitions),
			Confidence: calculateTransitionConfidence(transitions),
			Duration:   series.EndTime.Sub(series.StartTime),
		},
		Properties: map[string]interface{}{
			"transitions": transitions,
			"states":      detectStates(series.Points, transitions),
		},
		Created: time.Now(),
	}
}

// calculateTransitionStrength 计算转换强度
func calculateTransitionStrength(transitions []TimeSeriesPoint) float64 {
	if len(transitions) == 0 {
		return 0
	}

	// 计算转换幅度的平均值
	totalMagnitude := 0.0
	maxMagnitude := 0.0
	for _, t := range transitions {
		magnitude := math.Abs(t.Value)
		totalMagnitude += magnitude
		if magnitude > maxMagnitude {
			maxMagnitude = magnitude
		}
	}

	// 归一化强度
	if maxMagnitude > 0 {
		return (totalMagnitude / float64(len(transitions))) / maxMagnitude
	}
	return 0
}

// calculateTransitionConfidence 计算转换置信度
func calculateTransitionConfidence(transitions []TimeSeriesPoint) float64 {
	if len(transitions) < 2 {
		return 0
	}

	// 基于转换间隔的一致性计算置信度
	intervals := make([]float64, len(transitions)-1)
	totalInterval := 0.0

	for i := 1; i < len(transitions); i++ {
		interval := transitions[i].Timestamp.Sub(transitions[i-1].Timestamp).Seconds()
		intervals[i-1] = interval
		totalInterval += interval
	}

	// 计算间隔的变异系数
	avgInterval := totalInterval / float64(len(intervals))
	variance := 0.0
	for _, interval := range intervals {
		diff := interval - avgInterval
		variance += diff * diff
	}
	variance /= float64(len(intervals))

	// 一致性越高，变异系数越小，置信度越高
	cv := math.Sqrt(variance) / avgInterval
	return math.Max(0, 1-cv)
}

// detectStates 检测状态序列
func detectStates(points []TimeSeriesPoint, transitions []TimeSeriesPoint) []string {
	if len(points) == 0 || len(transitions) == 0 {
		return nil
	}

	states := make([]string, 0)
	currentState := "initial"
	states = append(states, currentState)

	transitionIdx := 0
	for i := 1; i < len(points); i++ {
		// 检查是否到达转换点
		if transitionIdx < len(transitions) &&
			points[i].Timestamp.Equal(transitions[transitionIdx].Timestamp) {
			// 确定新状态
			newState := determineState(points[i].Value, transitions[transitionIdx].Value)
			if newState != currentState {
				states = append(states, newState)
				currentState = newState
			}
			transitionIdx++
		}
	}

	return states
}

// determineState 确定状态
func determineState(value, transitionValue float64) string {
	// 基于变化率和当前值确定状态
	if math.Abs(transitionValue) < stateTransitionThreshold {
		if value < lowValueThreshold {
			return "low_stable"
		} else if value > highValueThreshold {
			return "high_stable"
		}
		return "stable"
	} else if transitionValue > 0 {
		if value > highValueThreshold {
			return "saturating"
		}
		return "increasing"
	} else {
		if value < lowValueThreshold {
			return "bottoming"
		}
		return "decreasing"
	}
}

// 辅助函数

// detectPeriods 检测时间序列中的周期
func detectPeriods(points []TimeSeriesPoint) []float64 {
	if len(points) < 4 {
		return nil
	}

	// 计算自相关序列
	maxLag := len(points) / 2
	autocorr := make([]float64, maxLag)
	mean := 0.0

	// 计算均值
	for _, p := range points {
		mean += p.Value
	}
	mean /= float64(len(points))

	// 计算自相关系数
	for lag := 0; lag < maxLag; lag++ {
		sum := 0.0
		count := 0
		for i := 0; i < len(points)-lag; i++ {
			sum += (points[i].Value - mean) * (points[i+lag].Value - mean)
			count++
		}
		if count > 0 {
			autocorr[lag] = sum / float64(count)
		}
	}

	// 查找峰值作为周期候选
	periods := make([]float64, 0)
	for i := 2; i < len(autocorr)-1; i++ {
		if autocorr[i] > autocorr[i-1] && autocorr[i] > autocorr[i+1] {
			timeDiff := points[i].Timestamp.Sub(points[0].Timestamp).Seconds()
			if timeDiff > 0 {
				periods = append(periods, timeDiff)
			}
		}
	}

	// 按周期长度排序
	sort.Float64s(periods)
	return periods
}

// calculateTrendLine 计算趋势线
func calculateTrendLine(points []TimeSeriesPoint) (slope float64, r2 float64) {
	if len(points) < 2 {
		return 0, 0
	}

	// 计算时间序列的 x 值(时间间隔)和 y 值
	n := float64(len(points))
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0
	baseTime := points[0].Timestamp.Unix()

	for _, p := range points {
		x := float64(p.Timestamp.Unix() - baseTime)
		y := p.Value

		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	// 计算斜率
	if sumX2 == sumX*sumX/n {
		return 0, 0
	}
	slope = (sumXY - sumX*sumY/n) / (sumX2 - sumX*sumX/n)

	// 计算 R² 值
	meanY := sumY / n
	totalSS := 0.0
	residualSS := 0.0

	for _, p := range points {
		x := float64(p.Timestamp.Unix() - baseTime)
		y := p.Value
		predicted := slope*x + (sumY-slope*sumX)/n

		totalSS += (y - meanY) * (y - meanY)
		residualSS += (y - predicted) * (y - predicted)
	}

	if totalSS == 0 {
		return slope, 0
	}
	r2 = 1 - residualSS/totalSS

	return slope, r2
}

// detectStateTransitions 检测状态转换点
func detectStateTransitions(points []TimeSeriesPoint) []TimeSeriesPoint {
	if len(points) < 2 {
		return nil
	}

	transitions := make([]TimeSeriesPoint, 0)
	lastState := "stable"
	var lastValue float64 = points[0].Value
	var lastTransitionTime time.Time = points[0].Timestamp

	for i := 1; i < len(points); i++ {
		// 计算变化率
		timeDiff := points[i].Timestamp.Sub(points[i-1].Timestamp).Seconds()
		if timeDiff < minTransitionInterval {
			continue
		}

		valueDiff := points[i].Value - lastValue
		changeRate := valueDiff / timeDiff

		// 确定当前状态
		currentState := "stable"
		if math.Abs(changeRate) > stateTransitionThreshold {
			if changeRate > 0 {
				currentState = "increasing"
			} else {
				currentState = "decreasing"
			}
		}

		// 检测状态转换
		if currentState != lastState &&
			points[i].Timestamp.Sub(lastTransitionTime).Seconds() >= minTransitionInterval {
			transitions = append(transitions, TimeSeriesPoint{
				Timestamp: points[i].Timestamp,
				Value:     changeRate,
				Type:      "state_transition",
				Metadata: map[string]interface{}{
					"from_state": lastState,
					"to_state":   currentState,
					"magnitude":  math.Abs(valueDiff),
				},
			})
			lastTransitionTime = points[i].Timestamp
		}

		lastState = currentState
		lastValue = points[i].Value
	}

	return transitions
}

// filterPatterns 根据置信度过滤模式
func filterPatterns(patterns []FlowPattern, minConfidence float64) []FlowPattern {
	filtered := make([]FlowPattern, 0)
	for _, pattern := range patterns {
		if pattern.Metrics.Confidence >= minConfidence {
			filtered = append(filtered, pattern)
		}
	}
	return filtered
}

// mergePatterns 合并相似模式
func mergePatterns(patterns []FlowPattern) []FlowPattern {
	if len(patterns) < 2 {
		return patterns
	}

	// 按类型分组
	typeGroups := make(map[string][]FlowPattern)
	for _, pattern := range patterns {
		typeGroups[pattern.Type] = append(typeGroups[pattern.Type], pattern)
	}

	merged := make([]FlowPattern, 0)
	// 处理每个类型组
	for _, group := range typeGroups {
		// 单个模式直接添加
		if len(group) == 1 {
			merged = append(merged, group[0])
			continue
		}

		// 合并相似模式
		for len(group) > 0 {
			base := group[0]
			similar := make([]FlowPattern, 0)
			remaining := make([]FlowPattern, 0)

			// 查找相似模式
			for _, other := range group[1:] {
				if areSimilarPatterns(base, other) {
					similar = append(similar, other)
				} else {
					remaining = append(remaining, other)
				}
			}

			// 合并相似模式
			if len(similar) > 0 {
				merged = append(merged, mergeSimularPatterns(append([]FlowPattern{base}, similar...)))
			} else {
				merged = append(merged, base)
			}

			group = remaining
		}
	}

	return merged
}

// areSimilarPatterns 判断两个模式是否相似
func areSimilarPatterns(p1, p2 FlowPattern) bool {
	// 相同类型
	if p1.Type != p2.Type {
		return false
	}

	// 时间接近
	if p1.Created.Sub(p2.Created).Hours() > 24 {
		return false
	}

	// 特征相似
	strengthDiff := math.Abs(p1.Metrics.Strength - p2.Metrics.Strength)
	if strengthDiff > 0.3 {
		return false
	}

	confidenceDiff := math.Abs(p1.Metrics.Confidence - p2.Metrics.Confidence)
	return confidenceDiff <= 0.2
}

// mergeSimularPatterns 合并相似模式
func mergeSimularPatterns(patterns []FlowPattern) FlowPattern {
	if len(patterns) == 0 {
		return FlowPattern{}
	}

	base := patterns[0]
	if len(patterns) == 1 {
		return base
	}

	// 计算平均指标
	var totalStrength, totalConfidence float64
	var totalDuration time.Duration
	var maxFrequency float64

	for _, p := range patterns {
		totalStrength += p.Metrics.Strength
		totalConfidence += p.Metrics.Confidence
		totalDuration += p.Metrics.Duration
		if p.Metrics.Frequency > maxFrequency {
			maxFrequency = p.Metrics.Frequency
		}
	}

	count := float64(len(patterns))
	merged := FlowPattern{
		ID:   generatePatternID(),
		Type: base.Type,
		Metrics: PatternMetrics{
			Frequency:  maxFrequency,
			Strength:   totalStrength / count,
			Confidence: totalConfidence / count,
			Duration:   time.Duration(int64(totalDuration) / int64(count)),
		},
		Properties: mergeProperties(patterns),
		Created:    time.Now(),
	}

	return merged
}

// mergeProperties 合并模式属性
func mergeProperties(patterns []FlowPattern) map[string]interface{} {
	merged := make(map[string]interface{})

	// 合并所有模式的属性
	for _, p := range patterns {
		for k, v := range p.Properties {
			if existing, ok := merged[k]; ok {
				// 如果是数值类型，取平均值
				if fv, ok := v.(float64); ok {
					if ef, ok := existing.(float64); ok {
						merged[k] = (ef + fv) / 2
					}
				}
			} else {
				merged[k] = v
			}
		}
	}

	return merged
}

// 能量指标计算
func calculateTotalEnergy(spans interface{}) float64 {
	total := 0.0
	if spanArray, ok := spans.([]*Span); ok {
		for _, span := range spanArray {
			metrics := (*span).GetMetrics()
			if energy, exists := metrics["energy"]; exists {
				total += energy
			}
		}
	}
	return total
}

func calculateAverageEnergy(spans interface{}) float64 {
	total := calculateTotalEnergy(spans)
	if spanArray, ok := spans.([]*Span); ok && len(spanArray) > 0 {
		return total / float64(len(spanArray))
	}
	return 0
}

func calculateEnergyVariance(spans interface{}) float64 {
	avg := calculateAverageEnergy(spans)
	variance := 0.0
	count := 0
	if spanArray, ok := spans.([]*Span); ok {
		for _, span := range spanArray {
			metrics := (*span).GetMetrics()
			if energy, exists := metrics["energy"]; exists {
				diff := energy - avg
				variance += diff * diff
				count++
			}
		}
	}
	if count > 0 {
		return variance / float64(count)
	}
	return 0
}

// 状态指标计算
func countStateTransitions(spans interface{}) int {
	transitions := 0
	if spanArray, ok := spans.([]*Span); ok && len(spanArray) > 1 {
		for i := 1; i < len(spanArray); i++ {
			prev := (*spanArray[i-1]).GetMetrics()["state"]
			curr := (*spanArray[i]).GetMetrics()["state"]
			if prev != curr {
				transitions++
			}
		}
	}
	return transitions
}

func calculateStateStability(spans interface{}) float64 {
	transitions := countStateTransitions(spans)
	if spanArray, ok := spans.([]*Span); ok && len(spanArray) > 1 {
		duration := (*spanArray[len(spanArray)-1]).GetEndTime().Sub((*spanArray[0]).GetStartTime())
		if duration > 0 {
			// 计算单位时间的转换频率，并转换为稳定性指标
			frequency := float64(transitions) / duration.Hours()
			return 1.0 / (1.0 + frequency)
		}
	}
	return 1.0
}

func determineCurrentPhase(spans interface{}) Phase {
	if spanArray, ok := spans.([]*Span); ok && len(spanArray) > 0 {
		lastSpan := spanArray[len(spanArray)-1]
		if phase, exists := (*lastSpan).GetMetrics()["phase"]; exists {
			return Phase(int(phase))
		}
	}
	return PhaseNone
}

// 性能指标计算
func calculateThroughput(spans interface{}) float64 {
	if spanArray, ok := spans.([]*Span); ok && len(spanArray) > 0 {
		duration := (*spanArray[len(spanArray)-1]).GetEndTime().Sub((*spanArray[0]).GetStartTime())
		if duration > 0 {
			return float64(len(spanArray)) / duration.Seconds()
		}
	}
	return 0
}

func calculateLatency(spans interface{}) float64 {
	total := 0.0
	count := 0
	if spanArray, ok := spans.([]*Span); ok {
		for _, span := range spanArray {
			total += (*span).GetDuration().Seconds()
			count++
		}
	}
	if count > 0 {
		return total / float64(count)
	}
	return 0
}

func calculateErrorRate(spans interface{}) float64 {
	errors := 0
	if spanArray, ok := spans.([]*Span); ok {
		for _, span := range spanArray {
			if err, exists := (*span).GetMetrics()["error"]; exists && err > 0 {
				errors++
			}
		}
		if len(spanArray) > 0 {
			return float64(errors) / float64(len(spanArray))
		}
	}
	return 0
}
