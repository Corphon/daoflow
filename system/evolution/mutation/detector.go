//system/evolution/mutation/detector.go

package mutation

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/common" // 引入共享接口包
	"github.com/Corphon/daoflow/system/evolution/pattern"
)

// MutationMetrics 突变指标
type MutationMetrics struct {
	// 基础请求统计
	TotalRequests  int           // 总请求数
	ResponseCount  int           // 响应数
	SuccessCount   int           // 成功数
	ErrorCount     int           // 错误数
	AverageLatency time.Duration // 平均延迟

	// 性能指标
	ThroughputRate float64 // 吞吐率
	ErrorRate      float64 // 错误率
	SuccessRate    float64 // 成功率

	// 时间窗口
	WindowStart time.Time // 统计开始时间
	WindowEnd   time.Time // 统计结束时间
}

// MutationDetector 突变检测器
type MutationDetector struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		detectionThreshold float64       // 检测阈值
		timeWindow         time.Duration // 时间窗口
		sensitivity        float64       // 灵敏度
		stabilityFactor    float64       // 稳定性因子
	}

	// 使用接口而不是具体类型
	patternAnalyzer common.PatternAnalyzer
	recognizer      pattern.PatternRecognizer
	currentState    *model.SystemState

	// 检测状态
	state struct {
		mutations    map[string]*Mutation         // 已检测突变
		observations []MutationObservation        // 观察记录
		baselines    map[string]*MutationBaseline // 基准线
		metrics      MutationMetrics              // 突变指标
	}
}

// MutationSource 突变源
type MutationSource struct {
	PatternID string                 // 相关模式ID
	Location  string                 // 发生位置
	Context   map[string]interface{} // 上下文信息
	Energy    float64                // 能量水平
}

// MutationChange 突变变化
type MutationChange struct {
	Property  string      // 变化属性
	OldValue  interface{} // 原值
	NewValue  interface{} // 新值
	Delta     float64     // 变化量
	Timestamp time.Time   // 变化时间
}

// MutationObservation 突变观察
type MutationObservation struct {
	Timestamp time.Time
	PatternID string
	Metrics   map[string]float64
	Anomalies []string
}

// MutationBaseline 突变基准线
type MutationBaseline struct {
	PatternID  string
	Metrics    map[string]BaselineMetric
	LastUpdate time.Time
	Confidence float64
}

// BaselineMetric 基准度量
type BaselineMetric struct {
	Mean    float64
	StdDev  float64
	Bounds  [2]float64 // [min, max]
	History []float64
}

// NewMutationDetector 创建新的突变检测器
func NewMutationDetector(analyzer common.PatternAnalyzer) *MutationDetector {

	md := &MutationDetector{
		patternAnalyzer: analyzer,
	}

	// 初始化配置
	md.config.detectionThreshold = 0.75
	md.config.timeWindow = 10 * time.Minute
	md.config.sensitivity = 0.8
	md.config.stabilityFactor = 0.6

	// 初始化状态
	md.state.mutations = make(map[string]*Mutation)
	md.state.observations = make([]MutationObservation, 0)
	md.state.baselines = make(map[string]*MutationBaseline)

	return md
}

// Detect 执行突变检测
func (md *MutationDetector) Detect() error {
	md.mu.Lock()
	defer md.mu.Unlock()

	// 获取当前模式
	patterns, err := md.recognizer.GetPatterns()
	if err != nil {
		return err
	}

	// 更新观察记录
	md.updateObservations(patterns)

	// 更新基准线
	md.updateBaselines()

	// 检测突变
	mutations := md.detectMutations(patterns)

	// 验证突变
	validated := md.validateMutations(mutations)

	// 更新突变状态
	md.updateMutations(validated)

	return nil
}

// updateMutations 方法
func (md *MutationDetector) updateMutations(mutations []*Mutation) {
	currentTime := time.Now()

	for _, mutation := range mutations {
		// 更新现有突变或添加新突变
		if existing, exists := md.state.mutations[mutation.ID]; exists {
			// 更新现有突变
			existing.Status = mutation.Status
			existing.Changes = mutation.Changes
			existing.Severity = mutation.Severity
			existing.Probability = mutation.Probability
			existing.LastUpdate = currentTime
		} else {
			// 添加新突变
			mutation.LastUpdate = currentTime
			md.state.mutations[mutation.ID] = mutation
		}
	}

	// 清理过期突变
	md.cleanupExpiredMutations()
}

// cleanupExpiredMutations 清理过期突变
func (md *MutationDetector) cleanupExpiredMutations() {
	currentTime := time.Now()
	expirationTime := currentTime.Add(-md.config.timeWindow)

	for id, mutation := range md.state.mutations {
		if mutation.LastUpdate.Before(expirationTime) {
			delete(md.state.mutations, id)
		}
	}
}

// cleanupObservations 清理过期观察记录
func (md *MutationDetector) cleanupObservations() {
	currentTime := time.Now()
	expirationTime := currentTime.Add(-md.config.timeWindow)

	// 保留未过期的观察记录
	validObservations := make([]MutationObservation, 0)
	for _, observation := range md.state.observations {
		if !observation.Timestamp.Before(expirationTime) {
			validObservations = append(validObservations, observation)
		}
	}

	md.state.observations = validObservations
}

// updateObservations 更新观察记录
func (md *MutationDetector) updateObservations(patterns []*pattern.RecognizedPattern) {
	currentTime := time.Now()

	// 更新系统状态
	md.currentState = &model.SystemState{
		Properties: make(map[string]interface{}),
		Timestamp:  currentTime,
	}

	// 收集系统状态属性
	systemMetrics := md.collectSystemMetrics(patterns)
	for key, value := range systemMetrics {
		md.currentState.Properties[key] = value
	}
	// 创建新的观察记录
	for _, pat := range patterns {
		observation := MutationObservation{
			Timestamp: currentTime,
			PatternID: pat.ID,
			Metrics:   md.collectMetrics(pat),
			Anomalies: make([]string, 0),
		}

		// 检查异常
		anomalies := md.checkAnomalies(pat, observation.Metrics)
		observation.Anomalies = anomalies

		md.state.observations = append(md.state.observations, observation)
	}

	// 清理过期观察
	md.cleanupObservations()
}

// updateBaselines 更新基准线
func (md *MutationDetector) updateBaselines() {
	currentTime := time.Now()

	// 获取观察数据，使用配置的时间窗口
	observations := md.getRecentObservations(md.config.timeWindow)

	// 更新每个模式的基准线
	for patternID := range md.collectPatternIDs(observations) {
		// 获取模式的观察数据
		patternObs := md.filterObservationsByPattern(observations, patternID)

		// 计算新的基准线
		baseline := md.calculateBaseline(patternObs)

		if baseline != nil {
			baseline.LastUpdate = currentTime
			md.state.baselines[patternID] = baseline
		}
	}
}

// detectMutations 检测突变
func (md *MutationDetector) detectMutations(
	patterns []*pattern.RecognizedPattern) []*Mutation {

	mutations := make([]*Mutation, 0)

	for _, pat := range patterns {
		// 获取基准线
		baseline := md.state.baselines[pat.ID]
		if baseline == nil {
			continue
		}

		// 检查突变条件
		if changes := md.checkMutationConditions(pat, baseline); len(changes) > 0 {
			// 创建突变记录
			mutation := &Mutation{
				ID:         generateMutationID(),
				Type:       md.determineMutationType(changes),
				Source:     createMutationSource(pat),
				Changes:    changes,
				DetectedAt: time.Now(),
				Status:     "detected",
			}

			// 计算严重程度和概率
			mutation.Severity = md.calculateMutationSeverity(mutation)
			mutation.Probability = md.calculateMutationProbability(mutation)

			mutations = append(mutations, mutation)
		}
	}

	return mutations
}

// validateMutations 验证突变
func (md *MutationDetector) validateMutations(
	mutations []*Mutation) []*Mutation {

	validated := make([]*Mutation, 0)

	for _, mutation := range mutations {
		// 检查突变有效性
		if md.isMutationValid(mutation) {
			validated = append(validated, mutation)
		}
	}

	return validated
}

// GetActiveMutations returns currently active mutations
func (md *MutationDetector) GetActiveMutations() ([]*Mutation, error) {
	md.mu.RLock()
	defer md.mu.RUnlock()

	currentTime := time.Now()
	cutoffTime := currentTime.Add(-md.config.timeWindow)

	active := make([]*Mutation, 0)
	for _, mutation := range md.state.mutations {
		// 检查突变是否在活跃时间窗口内且状态为活跃
		if !mutation.LastUpdate.Before(cutoffTime) && mutation.Status == "detected" {
			active = append(active, mutation)
		}
	}

	return active, nil
}

// GetObservations 获取指定时间窗口内的观察记录
func (md *MutationDetector) GetObservations(window time.Duration) []MutationObservation {
	md.mu.RLock()
	defer md.mu.RUnlock()
	return md.getRecentObservations(window)
}

// GetBaselines 获取当前基准线集合
func (md *MutationDetector) GetBaselines() map[string]*MutationBaseline {
	md.mu.RLock()
	defer md.mu.RUnlock()

	// 创建基准线的副本
	baselines := make(map[string]*MutationBaseline, len(md.state.baselines))
	for id, baseline := range md.state.baselines {
		baselines[id] = baseline
	}
	return baselines
}

// 辅助函数

func (md *MutationDetector) collectMetrics(
	pattern *pattern.RecognizedPattern) map[string]float64 {

	metrics := make(map[string]float64)

	// 基础指标收集
	metrics["stability"] = pattern.GetStability()
	metrics["strength"] = pattern.GetStrength() // Using GetStrength() instead of Intensity
	metrics["complexity"] = md.calculatePatternComplexity(pattern)
	metrics["coherence"] = md.calculatePatternCoherence(pattern)

	return metrics
}

// calculatePatternComplexity 计算模式复杂度
func (md *MutationDetector) calculatePatternComplexity(pattern *pattern.RecognizedPattern) float64 {
	// 基于强度和稳定性计算复杂度
	return pattern.GetStrength() * pattern.GetStability()
}

// calculatePatternCoherence 计算模式一致性
func (md *MutationDetector) calculatePatternCoherence(pattern *pattern.RecognizedPattern) float64 {
	// 基于稳定性计算相干性
	return pattern.GetStability() * (1 - (1 - pattern.GetStability()))
}

func (md *MutationDetector) checkAnomalies(
	pattern *pattern.RecognizedPattern,
	metrics map[string]float64) []string {

	anomalies := make([]string, 0)

	baseline := md.state.baselines[pattern.ID]
	if baseline == nil {
		return anomalies
	}

	// 检查每个指标
	for name, value := range metrics {
		if baseMetric, ok := baseline.Metrics[name]; ok {
			if md.isAnomaly(value, baseMetric) {
				anomalies = append(anomalies, name)
			}
		}
	}

	return anomalies
}

// getRecentObservations 获取最近的观察记录
func (md *MutationDetector) getRecentObservations(window time.Duration) []MutationObservation {
	currentTime := time.Now()
	cutoffTime := currentTime.Add(-window)

	recent := make([]MutationObservation, 0)
	for _, obs := range md.state.observations {
		if !obs.Timestamp.Before(cutoffTime) {
			recent = append(recent, obs)
		}
	}
	return recent
}

// collectPatternIDs 收集所有模式ID
func (md *MutationDetector) collectPatternIDs(observations []MutationObservation) map[string]bool {
	patterns := make(map[string]bool)
	for _, obs := range observations {
		patterns[obs.PatternID] = true
	}
	return patterns
}

// filterObservationsByPattern 按模式ID过滤观察记录
func (md *MutationDetector) filterObservationsByPattern(
	observations []MutationObservation,
	patternID string) []MutationObservation {

	filtered := make([]MutationObservation, 0)
	for _, obs := range observations {
		if obs.PatternID == patternID {
			filtered = append(filtered, obs)
		}
	}
	return filtered
}

// calculateBaseline 计算基准线
func (md *MutationDetector) calculateBaseline(observations []MutationObservation) *MutationBaseline {
	if len(observations) == 0 {
		return nil
	}

	baseline := &MutationBaseline{
		PatternID:  observations[0].PatternID,
		Metrics:    make(map[string]BaselineMetric),
		LastUpdate: time.Now(),
		Confidence: calculateConfidence(len(observations)),
	}

	// 为每个指标计算基准值
	metrics := collectMetricNames(observations)
	for _, name := range metrics {
		values := collectMetricValues(observations, name)
		baseline.Metrics[name] = calculateMetricBaseline(values)
	}

	return baseline
}

// isAnomaly 检查是否为异常值
func (md *MutationDetector) isAnomaly(value float64, metric BaselineMetric) bool {
	// 使用标准差作为异常判断标准
	threshold := metric.StdDev * md.config.sensitivity

	if value < metric.Mean-threshold || value > metric.Mean+threshold {
		return true
	}
	return false
}

// checkMutationConditions 检查突变条件
func (md *MutationDetector) checkMutationConditions(
	pattern *pattern.RecognizedPattern,
	baseline *MutationBaseline) []MutationChange {

	changes := make([]MutationChange, 0)
	currentMetrics := md.collectMetrics(pattern)

	for name, value := range currentMetrics {
		if baseMetric, ok := baseline.Metrics[name]; ok {
			if md.isSignificantChange(value, baseMetric) {
				changes = append(changes, MutationChange{
					Property:  name,
					OldValue:  baseMetric.Mean,
					NewValue:  value,
					Delta:     value - baseMetric.Mean,
					Timestamp: time.Now(),
				})
			}
		}
	}

	return changes
}

// determineMutationType 确定突变类型
func (md *MutationDetector) determineMutationType(changes []MutationChange) string {
	// 基于变化特征确定类型
	if len(changes) == 0 {
		return "unknown"
	}

	// 分析变化模式
	hasEnergy := false
	hasStructure := false
	for _, change := range changes {
		switch change.Property {
		case "energy", "intensity":
			hasEnergy = true
		case "structure", "pattern":
			hasStructure = true
		}
	}

	// 确定类型
	if hasEnergy && hasStructure {
		return "compound"
	} else if hasEnergy {
		return "energetic"
	} else if hasStructure {
		return "structural"
	}

	return "behavioral"
}

// createMutationSource creates a new mutation source
func createMutationSource(pattern *pattern.RecognizedPattern) *MutationSource {
	return &MutationSource{
		PatternID: pattern.ID,
		Context:   make(map[string]interface{}),
		Energy:    pattern.GetStrength(), // Use GetStrength() instead of Intensity
	}
}

// calculateMutationSeverity 计算突变严重程度
func (md *MutationDetector) calculateMutationSeverity(mutation *Mutation) float64 {
	if len(mutation.Changes) == 0 {
		return 0
	}

	totalSeverity := 0.0
	for _, change := range mutation.Changes {
		severity := math.Abs(change.Delta)
		totalSeverity += severity
	}

	return totalSeverity / float64(len(mutation.Changes))
}

// calculateMutationProbability 计算突变概率
func (md *MutationDetector) calculateMutationProbability(mutation *Mutation) float64 {
	// 基于历史数据和当前状态计算概率
	baseline := md.state.baselines[mutation.Source.PatternID]
	if baseline == nil {
		return 0.5 // 默认概率
	}

	// 结合多个因素
	factors := []float64{
		baseline.Confidence,
		1 - normalizeValue(mutation.Severity, 0, 1),
		md.config.stabilityFactor,
	}

	// 计算加权平均
	totalWeight := 0.0
	weightedSum := 0.0
	for i, factor := range factors {
		weight := 1.0 / float64(i+1) // 递减权重
		weightedSum += factor * weight
		totalWeight += weight
	}

	return weightedSum / totalWeight
}

// isMutationValid 验证突变有效性
func (md *MutationDetector) isMutationValid(mutation *Mutation) bool {
	// 检查基本属性
	if mutation == nil || len(mutation.Changes) == 0 {
		return false
	}

	// 检查阈值
	if mutation.Probability < md.config.detectionThreshold {
		return false
	}

	// 检查时间窗口
	if time.Since(mutation.DetectedAt) > md.config.timeWindow {
		return false
	}

	return true
}

// 辅助函数

func calculateConfidence(sampleSize int) float64 {
	// 基于样本大小计算置信度
	minSamples := 5
	maxSamples := 100

	if sampleSize < minSamples {
		return 0.5
	}
	if sampleSize > maxSamples {
		return 1.0
	}

	return 0.5 + 0.5*float64(sampleSize-minSamples)/float64(maxSamples-minSamples)
}

func collectMetricNames(observations []MutationObservation) []string {
	nameSet := make(map[string]bool)
	for _, obs := range observations {
		for name := range obs.Metrics {
			nameSet[name] = true
		}
	}

	names := make([]string, 0, len(nameSet))
	for name := range nameSet {
		names = append(names, name)
	}
	return names
}

func collectMetricValues(observations []MutationObservation, metricName string) []float64 {
	values := make([]float64, 0)
	for _, obs := range observations {
		if value, ok := obs.Metrics[metricName]; ok {
			values = append(values, value)
		}
	}
	return values
}

func calculateMetricBaseline(values []float64) BaselineMetric {
	if len(values) == 0 {
		return BaselineMetric{}
	}

	// 计算均值
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	// 计算标准差
	sumSquares := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquares += diff * diff
	}
	stdDev := math.Sqrt(sumSquares / float64(len(values)))

	// 计算边界
	min := mean - 2*stdDev
	max := mean + 2*stdDev

	return BaselineMetric{
		Mean:    mean,
		StdDev:  stdDev,
		Bounds:  [2]float64{min, max},
		History: values,
	}
}

func normalizeValue(value float64, min float64, max float64) float64 {
	if value < min {
		return 0
	}
	if value > max {
		return 1
	}
	return (value - min) / (max - min)
}

func (md *MutationDetector) isSignificantChange(value float64, metric BaselineMetric) bool {
	relativeDiff := math.Abs(value-metric.Mean) / metric.Mean
	return relativeDiff > md.config.sensitivity
}

func calculatePatternEnergy(pattern *pattern.RecognizedPattern) float64 {
	if pattern == nil {
		return 0
	}

	// 1. 基础能量 - 来自模式强度
	baseEnergy := pattern.GetStrength()

	// 2. 活跃能量 - 基于激活水平
	activeEnergy := pattern.GetActivationLevel()

	// 3. 稳定性能量 - 基于模式稳定性
	stabilityEnergy := pattern.GetStability()

	// 4. 演化能量 - 基于演化历史
	evolutionEnergy := calculateEvolutionEnergy(pattern)

	// 5. 计算总能量
	// - 基础能量权重 0.3
	// - 活跃能量权重 0.2
	// - 稳定性能量权重 0.3
	// - 演化能量权重 0.2
	totalEnergy := (baseEnergy * 0.3) +
		(activeEnergy * 0.2) +
		(stabilityEnergy * 0.3) +
		(evolutionEnergy * 0.2)

	// 归一化到 [0,1] 范围
	return math.Max(0, math.Min(1, totalEnergy))
}

// calculateEvolutionEnergy 计算演化能量
func calculateEvolutionEnergy(pattern *pattern.RecognizedPattern) float64 {
	if len(pattern.Evolution) == 0 {
		return 0
	}

	// 计算演化趋势
	energy := 0.0
	for i := 1; i < len(pattern.Evolution); i++ {
		curr := pattern.Evolution[i]
		prev := pattern.Evolution[i-1]

		// 分析状态变化
		if value, exists := curr.Properties["energy"]; exists {
			prevValue := prev.Properties["energy"]
			// 计算变化率
			delta := value - prevValue
			energy += delta
		}
	}

	// 归一化演化能量
	return math.Max(0, math.Min(1, energy/float64(len(pattern.Evolution))))
}
func generateMutationID() string {
	return fmt.Sprintf("mut_%d", time.Now().UnixNano())
}

// 2. 添加 GetCurrentState 方法
func (md *MutationDetector) GetCurrentState() *model.SystemState {
	md.mu.RLock()
	defer md.mu.RUnlock()

	if md.currentState == nil {
		return &model.SystemState{
			Properties: make(map[string]interface{}),
		}
	}
	return md.currentState
}

// 辅助方法收集系统指标
func (md *MutationDetector) collectSystemMetrics(
	patterns []*pattern.RecognizedPattern) map[string]interface{} {

	metrics := make(map[string]interface{})

	// 计算系统能量
	totalEnergy := 0.0
	for _, pat := range patterns {
		totalEnergy += calculatePatternEnergy(pat)
	}
	metrics["energy"] = totalEnergy

	// 计算系统稳定性
	stability := md.calculateSystemStability(md.state.metrics)
	metrics["stability"] = stability

	return metrics
}

// 系统稳定性计算方法
func (md *MutationDetector) calculateSystemStability(metrics MutationMetrics) float64 {
	// 基于多个指标加权计算系统稳定性
	weights := map[string]float64{
		"response_rate": 0.3, // 响应率权重
		"success_rate":  0.3, // 成功率权重
		"error_rate":    0.2, // 错误率权重
		"latency":       0.2, // 延迟权重
	}

	stability := 0.0

	if metrics.TotalRequests > 0 {
		// 响应率
		responseRate := float64(metrics.ResponseCount) / float64(metrics.TotalRequests)
		stability += responseRate * weights["response_rate"]

		// 成功率
		successRate := float64(metrics.SuccessCount) / float64(metrics.TotalRequests)
		stability += successRate * weights["success_rate"]

		// 错误率 (反向指标)
		errorRate := float64(metrics.ErrorCount) / float64(metrics.TotalRequests)
		stability += (1 - errorRate) * weights["error_rate"]
	}

	// 延迟评分 (使用反比函数将延迟转换为0-1的得分)
	if metrics.AverageLatency > 0 {
		latencyScore := 1.0 / (1.0 + math.Log1p(float64(metrics.AverageLatency.Milliseconds())))
		stability += latencyScore * weights["latency"]
	}

	// 确保结果在0-1范围内
	return math.Max(0, math.Min(1, stability))
}
