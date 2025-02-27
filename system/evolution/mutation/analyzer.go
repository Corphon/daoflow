//system/evolution/mutation/analyzer.go

package mutation

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/Corphon/daoflow/model"
)

const (
	maxMetricsHistory = 1000
)

// MutationAnalyzer 突变分析器
type MutationAnalyzer struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		analysisDepth        int           // 分析深度
		correlationThreshold float64       // 相关性阈值
		patternWindow        time.Duration // 模式窗口
		predictionHorizon    time.Duration // 预测周期
	}

	// 分析状态
	state struct {
		analyses    map[string]*model.MutationAnalysis // 分析结果
		patterns    map[string]*MutationPattern        // 突变模式
		predictions []MutationPrediction               // 预测结果
		metrics     AnalysisMetrics                    // 分析指标
	}

	// 依赖项
	detector *MutationDetector
	handler  *MutationHandler
}

// MutationPattern 突变模式
type MutationPattern struct {
	ID         string             // 模式ID
	Signature  []PatternFeature   // 特征签名
	Frequency  float64            // 发生频率
	Conditions map[string]float64 // 触发条件
	Timeline   []PatternEvent     // 时间线
}

// PatternFeature 模式特征
type PatternFeature struct {
	Type       string      // 特征类型
	Value      interface{} // 特征值
	Importance float64     // 重要性
}

// PatternEvent 模式事件
type PatternEvent struct {
	Time time.Time
	Type string
	Data map[string]interface{}
}

// MutationPrediction 突变预测
type MutationPrediction struct {
	ID          string                // 预测ID
	PatternID   string                // 模式ID
	Probability float64               // 发生概率
	TimeFrame   time.Duration         // 时间框架
	Conditions  []PredictionCondition // 预测条件
	Created     time.Time             // 创建时间
}

// PredictionCondition 预测条件
type PredictionCondition struct {
	Type      string      // 条件类型
	PatternID string      // 模式ID
	Expected  interface{} // 预期值
	Tolerance float64     // 容差
}

// AnalysisMetrics 分析指标
type AnalysisMetrics struct {
	Accuracy    map[string]float64 // 准确率指标
	Coverage    float64            // 覆盖率
	Latency     time.Duration      // 分析延迟
	Performance []PerformancePoint // 性能指标
}

// PerformancePoint 性能指标点
type PerformancePoint struct {
	Time    time.Time
	Metrics map[string]float64
}

// -----------------------------------------------------
// NewMutationAnalyzer 创建新的突变分析器
func NewMutationAnalyzer(detector *MutationDetector, handler *MutationHandler) *MutationAnalyzer {
	ma := &MutationAnalyzer{
		detector: detector,
		handler:  handler,
	}

	// 初始化配置
	ma.config.analysisDepth = 3
	ma.config.correlationThreshold = 0.7
	ma.config.patternWindow = 24 * time.Hour
	ma.config.predictionHorizon = 12 * time.Hour

	// 初始化状态
	ma.state.analyses = make(map[string]*model.MutationAnalysis)
	ma.state.patterns = make(map[string]*MutationPattern)
	ma.state.predictions = make([]MutationPrediction, 0)
	ma.state.metrics = AnalysisMetrics{
		Accuracy:    make(map[string]float64),
		Performance: make([]PerformancePoint, 0),
	}

	return ma
}

// Analyze 执行突变分析
func (ma *MutationAnalyzer) Analyze() error {
	ma.mu.Lock()
	defer ma.mu.Unlock()

	// 获取最新突变
	mutations, err := ma.detector.GetRecentMutations(ma.config.patternWindow)
	if err != nil {
		return err
	}

	// 分析突变模式
	patterns := ma.analyzePatterns(mutations)

	// 执行因果分析
	analyses := ma.analyzeCausality(mutations)

	// 生成预测
	predictions := ma.generatePredictions(patterns)

	// 更新状态
	ma.updateState(patterns, analyses, predictions)

	// 更新指标
	ma.updateMetrics()

	return nil
}

// analyzePatterns 分析突变模式
func (ma *MutationAnalyzer) analyzePatterns(
	mutations []*Mutation) map[string]*MutationPattern {

	patterns := make(map[string]*MutationPattern)

	// 按时间分组分析
	timeGroups := ma.groupMutationsByTime(mutations)

	for _, group := range timeGroups {
		// 提取模式特征
		features := ma.extractPatternFeatures(group)

		// 识别模式
		if pattern := ma.identifyPattern(features); pattern != nil {
			patterns[pattern.ID] = pattern
		}
	}

	return patterns
}

// analyzeCausality 分析因果关系
func (ma *MutationAnalyzer) analyzeCausality(
	mutations []*Mutation) []*model.MutationAnalysis {

	analyses := make([]*model.MutationAnalysis, 0)

	for _, mutation := range mutations {
		// 创建分析实例
		analysis := &model.MutationAnalysis{
			ID:         generateAnalysisID(),
			MutationID: mutation.ID,
			Created:    time.Now(),
		}

		// 分析因果因素
		analysis.Causes = ma.analyzeCauses(mutation)

		// 分析影响效果
		analysis.Effects = ma.analyzeEffects(mutation)

		// 分析相关性
		analysis.Correlations = ma.findCorrelations(mutation)

		// 评估风险
		analysis.Risk = ma.assessRisk(mutation)

		analyses = append(analyses, analysis)
	}

	return analyses
}

// generatePredictions 生成预测
func (ma *MutationAnalyzer) generatePredictions(
	patterns map[string]*MutationPattern) []MutationPrediction {

	predictions := make([]MutationPrediction, 0)

	for _, pattern := range patterns {
		// 分析模式趋势
		trend := ma.analyzePatternTrend(pattern)

		// 预测未来发生
		if prediction := ma.predictPattern(pattern, trend); prediction != nil {
			predictions = append(predictions, *prediction)
		}
	}

	return predictions
}

// GetRecentMutations returns mutations within the time window
func (md *MutationDetector) GetRecentMutations(window time.Duration) ([]*Mutation, error) {
	md.mu.RLock()
	defer md.mu.RUnlock()

	currentTime := time.Now()
	cutoffTime := currentTime.Add(-window)

	recent := make([]*Mutation, 0)
	for _, mutation := range md.state.mutations {
		if !mutation.DetectedAt.Before(cutoffTime) {
			recent = append(recent, mutation)
		}
	}
	return recent, nil
}

// updateState updates the analyzer state with new data
func (ma *MutationAnalyzer) updateState(patterns map[string]*MutationPattern,
	analyses []*model.MutationAnalysis, predictions []MutationPrediction) {
	ma.state.patterns = patterns
	for _, analysis := range analyses {
		ma.state.analyses[analysis.ID] = analysis
	}
	ma.state.predictions = predictions
}

// groupMutationsByTime groups mutations by time periods
func (ma *MutationAnalyzer) groupMutationsByTime(mutations []*Mutation) [][]*Mutation {
	if len(mutations) == 0 {
		return nil
	}

	groups := make([][]*Mutation, 0)
	currentGroup := []*Mutation{mutations[0]}

	for i := 1; i < len(mutations); i++ {
		if ma.isInSameTimeGroup(mutations[i], currentGroup[0]) {
			currentGroup = append(currentGroup, mutations[i])
		} else {
			groups = append(groups, currentGroup)
			currentGroup = []*Mutation{mutations[i]}
		}
	}

	if len(currentGroup) > 0 {
		groups = append(groups, currentGroup)
	}

	return groups
}

func (ma *MutationAnalyzer) isInSameTimeGroup(m1, m2 *Mutation) bool {
	return m1.DetectedAt.Sub(m2.DetectedAt) <= ma.config.patternWindow
}

// extractPatternFeatures extracts features from a group of mutations
func (ma *MutationAnalyzer) extractPatternFeatures(mutations []*Mutation) []PatternFeature {
	features := make([]PatternFeature, 0)

	// Extract common patterns from mutations
	typeFreq := make(map[string]int)
	severitySum := make(map[string]float64)

	for _, mutation := range mutations {
		typeFreq[mutation.Type]++
		severitySum[mutation.Type] += mutation.Severity
	}

	// Convert to features
	for mutType, freq := range typeFreq {
		features = append(features, PatternFeature{
			Type:       mutType,
			Value:      float64(freq) / float64(len(mutations)),
			Importance: severitySum[mutType] / float64(freq),
		})
	}

	return features
}

// identifyPattern identifies a pattern from features
func (ma *MutationAnalyzer) identifyPattern(features []PatternFeature) *MutationPattern {
	if len(features) == 0 {
		return nil
	}

	pattern := &MutationPattern{
		ID:         fmt.Sprintf("pat_%d", time.Now().UnixNano()),
		Signature:  features,
		Frequency:  calculatePatternFrequency(features),
		Conditions: make(map[string]float64),
		Timeline:   make([]PatternEvent, 0),
	}

	return pattern
}

// analyzeCauses analyzes causal factors of a mutation
func (ma *MutationAnalyzer) analyzeCauses(mutation *Mutation) []model.CausalFactor {
	causes := make([]model.CausalFactor, 0)

	// Analyze mutation source
	if mutation.Source != nil {
		causes = append(causes, model.CausalFactor{
			Type:       "source",
			Source:     mutation.Source.PatternID,
			Weight:     mutation.Source.Energy,
			Confidence: 0.8,
			Evidence:   []string{mutation.Source.Location},
		})
	}

	return causes
}

// analyzeEffects analyzes effects of a mutation
func (ma *MutationAnalyzer) analyzeEffects(mutation *Mutation) []model.Effect {
	effects := make([]model.Effect, 0)

	for _, change := range mutation.Changes {
		effects = append(effects, model.Effect{
			Target:     change.Property,
			Type:       "property_change",
			Magnitude:  change.Delta,
			Duration:   time.Since(change.Timestamp),
			Reversible: true,
		})
	}

	return effects
}

// findCorrelations finds correlations for a mutation
func (ma *MutationAnalyzer) findCorrelations(mutation *Mutation) []model.Correlation {
	correlations := make([]model.Correlation, 0)

	// Find correlations with other mutations
	for id, other := range ma.state.analyses {
		if id == mutation.ID {
			continue
		}

		strength := calculateCorrelationStrength(mutation, other)
		if strength >= ma.config.correlationThreshold {
			correlations = append(correlations, model.Correlation{
				SourceID:   mutation.ID,
				TargetID:   id,
				Type:       "mutation",
				Strength:   strength,
				Direction:  determineCorrelationDirection(mutation, other),
				TimeOffset: other.Created.Sub(mutation.DetectedAt),
			})
		}
	}

	return correlations
}

// assessRisk assesses risk of a mutation
func (ma *MutationAnalyzer) assessRisk(mutation *Mutation) model.RiskAssessment {
	return model.RiskAssessment{
		Level: determineRiskLevel(mutation.Severity),
		Score: calculateRiskScore(mutation),
		Factors: []model.RiskFactor{
			{
				Type:        "severity",
				Impact:      mutation.Severity,
				Probability: mutation.Probability,
				Urgency:     determineUrgency(mutation),
			},
		},
		Mitigation: suggestMitigations(mutation),
	}
}

// analyzePatternTrend analyzes trend of a pattern
func (ma *MutationAnalyzer) analyzePatternTrend(pattern *MutationPattern) float64 {
	if len(pattern.Timeline) < 2 {
		return 0
	}

	// Calculate trend based on frequency changes
	return pattern.Frequency
}

// predictPattern predicts future occurrence of a pattern
func (ma *MutationAnalyzer) predictPattern(pattern *MutationPattern, trend float64) *MutationPrediction {
	if trend <= 0 {
		return nil
	}

	return &MutationPrediction{
		ID:          fmt.Sprintf("pred_%d", time.Now().UnixNano()),
		PatternID:   pattern.ID,
		Probability: calculatePredictionProbability(pattern, trend),
		TimeFrame:   ma.config.predictionHorizon,
		Created:     time.Now(),
	}
}

// calculateAccuracy calculates analysis accuracy
func (ma *MutationAnalyzer) calculateAccuracy() float64 {
	if len(ma.state.predictions) == 0 {
		return 1.0
	}

	total := float64(len(ma.state.predictions))
	accurate := 0.0

	for _, pred := range ma.state.predictions {
		if ma.isPredictionAccurate(pred) {
			accurate++
		}
	}

	return accurate / total
}

// calculateCoverage calculates analysis coverage
func (ma *MutationAnalyzer) calculateCoverage() float64 {
	if len(ma.state.patterns) == 0 {
		return 1.0
	}

	covered := float64(len(ma.state.analyses))
	total := float64(len(ma.state.patterns))

	return covered / total
}

// Helper functions
func calculatePatternFrequency(features []PatternFeature) float64 {
	if len(features) == 0 {
		return 0
	}

	total := 0.0
	for _, f := range features {
		if val, ok := f.Value.(float64); ok {
			total += val
		}
	}
	return total / float64(len(features))
}

// calculateCorrelationStrength 计算两个突变之间的相关性强度
func calculateCorrelationStrength(m1 *Mutation, a2 *model.MutationAnalysis) float64 {
	if m1 == nil || a2 == nil {
		return 0
	}

	// 1. 时间相关性 (30%)
	timeCorrelation := calculateTimeCorrelation(m1.DetectedAt, a2.Created)

	// 2. 特征相关性 (40%)
	featureCorrelation := calculateFeatureCorrelation(m1.Changes, a2.Effects)

	// 3. 因果相关性 (30%)
	causalCorrelation := calculateCausalCorrelation(m1.Source, a2.Causes)

	// 计算加权总分
	totalScore := (timeCorrelation * 0.3) +
		(featureCorrelation * 0.4) +
		(causalCorrelation * 0.3)

	return math.Max(0, math.Min(1, totalScore))
}

// determineCorrelationDirection 确定相关性方向
// 返回: -1 (负相关), 0 (无相关), 1 (正相关)
func determineCorrelationDirection(m1 *Mutation, a2 *model.MutationAnalysis) int {
	if m1 == nil || a2 == nil {
		return 0
	}

	// 1. 检查时间顺序
	if m1.DetectedAt.After(a2.Created) {
		return -1
	}

	// 2. 分析变化趋势
	m1Trend := calculateChangeTrend(m1.Changes)
	a2Trend := calculateEffectTrend(a2.Effects)

	// 3. 确定方向
	if m1Trend*a2Trend > 0 {
		return 1 // 同向变化
	} else if m1Trend*a2Trend < 0 {
		return -1 // 反向变化
	}

	return 0 // 无明显相关
}

// 辅助函数

func calculateTimeCorrelation(t1, t2 time.Time) float64 {
	// 计算时间距离，距离越近相关性越强
	timeDiff := t2.Sub(t1).Abs()
	maxDiff := 24 * time.Hour

	if timeDiff > maxDiff {
		return 0
	}
	return 1 - float64(timeDiff)/float64(maxDiff)
}

func calculateFeatureCorrelation(changes []MutationChange, effects []model.Effect) float64 {
	if len(changes) == 0 || len(effects) == 0 {
		return 0
	}

	matchCount := 0
	for _, change := range changes {
		for _, effect := range effects {
			if change.Property == effect.Target {
				matchCount++
				break
			}
		}
	}

	return float64(matchCount) / math.Max(float64(len(changes)), float64(len(effects)))
}

func calculateCausalCorrelation(source *MutationSource, causes []model.CausalFactor) float64 {
	if source == nil || len(causes) == 0 {
		return 0
	}

	totalWeight := 0.0
	matchWeight := 0.0

	for _, cause := range causes {
		totalWeight += cause.Weight
		if cause.Source == source.PatternID {
			matchWeight += cause.Weight
		}
	}

	if totalWeight == 0 {
		return 0
	}
	return matchWeight / totalWeight
}

func calculateChangeTrend(changes []MutationChange) float64 {
	if len(changes) == 0 {
		return 0
	}

	trend := 0.0
	for _, change := range changes {
		trend += change.Delta
	}
	return trend / float64(len(changes))
}

func calculateEffectTrend(effects []model.Effect) float64 {
	if len(effects) == 0 {
		return 0
	}

	trend := 0.0
	for _, effect := range effects {
		trend += effect.Magnitude
	}
	return trend / float64(len(effects))
}

func determineRiskLevel(severity float64) string {
	if severity >= 0.8 {
		return "high"
	} else if severity >= 0.5 {
		return "medium"
	}
	return "low"
}

func calculateRiskScore(m *Mutation) float64 {
	return m.Severity * m.Probability
}

func determineUrgency(m *Mutation) int {
	if m.Severity >= 0.8 {
		return 3 // High
	} else if m.Severity >= 0.5 {
		return 2 // Medium
	}
	return 1 // Low
}

func suggestMitigations(m *Mutation) []string {
	if m == nil {
		return []string{"invalid mutation"}
	}

	mitigations := make([]string, 0)

	// 1. 基于严重程度的建议
	switch determineRiskLevel(m.Severity) {
	case "high":
		mitigations = append(mitigations,
			"immediate intervention required",
			"activate emergency response plan",
			"notify system administrators",
		)
	case "medium":
		mitigations = append(mitigations,
			"increase monitoring frequency",
			"prepare contingency measures",
			"analyze root causes",
		)
	case "low":
		mitigations = append(mitigations,
			"monitor regularly",
			"document changes",
			"update baseline metrics",
		)
	}

	// 2. 基于突变类型的具体建议
	switch m.Type {
	case "energetic":
		mitigations = append(mitigations,
			"balance energy distribution",
			"optimize resource allocation",
		)
	case "structural":
		mitigations = append(mitigations,
			"verify system integrity",
			"strengthen weak components",
		)
	case "behavioral":
		mitigations = append(mitigations,
			"adjust behavior parameters",
			"review pattern recognition rules",
		)
	}

	// 3. 基于概率的预防措施
	if m.Probability > 0.7 {
		mitigations = append(mitigations,
			"implement preventive measures",
			"enhance early warning system",
		)
	}

	return mitigations
}

func (ma *MutationAnalyzer) isPredictionAccurate(p MutationPrediction) bool {
	if p.Created.Add(p.TimeFrame).Before(time.Now()) {
		// 预测时间窗口已过，检查实际发生情况

		// 1. 基础准确性检查
		if p.Probability < 0.2 {
			// 低概率预测，默认为准确
			return true
		}

		// 2. 时间窗口评估
		timeDeviation := time.Since(p.Created)
		if timeDeviation > p.TimeFrame*2 {
			// 超出预期时间范围太多
			return false
		}

		// 3. 概率阈值检查
		if p.Probability > 0.8 {
			// 高概率预测需要更严格的验证
			// TODO: 实现实际验证逻辑
			return ma.checkHighProbabilityPrediction(p)
		}

		// 4. 条件满足度检查
		for _, condition := range p.Conditions {
			if !ma.isConditionMet(condition) {
				return false
			}
		}

		return true
	}

	// 预测时间窗口未结束，暂时返回 true
	return true
}

// 辅助函数

func (ma *MutationAnalyzer) checkHighProbabilityPrediction(p MutationPrediction) bool {
	if p.Probability < 0.8 {
		return true // 只验证高概率预测
	}

	// 1. 验证时间窗口合理性 (30%)
	timeScore := calculateTimeScore(p)

	// 2. 验证条件满足度 (40%)
	conditionScore := ma.calculateConditionScore(p)

	// 3. 验证模式稳定性 (30%)
	stabilityScore := calculateStabilityScore(p)

	// 计算加权总分
	totalScore := (timeScore * 0.3) +
		(conditionScore * 0.4) +
		(stabilityScore * 0.3)

	// 高概率预测要求更高的准确度阈值
	return totalScore >= 0.85
}

func (ma *MutationAnalyzer) isConditionMet(condition PredictionCondition) bool {
	if condition.Expected == nil {
		return false
	}

	// 根据不同的条件类型进行验证
	switch condition.Type {
	case "threshold":
		// 检查数值是否在容差范围内
		if expectedVal, ok := condition.Expected.(float64); ok {
			actualVal := ma.getCurrentValue(condition) // 使用成员方法
			tolerance := math.Abs(expectedVal * condition.Tolerance)
			return math.Abs(actualVal-expectedVal) <= tolerance
		}

	case "state":
		// 检查状态是否匹配
		if expectedState, ok := condition.Expected.(string); ok {
			actualState := ma.getCurrentState(condition) // 使用成员方法
			return expectedState == actualState
		}

	case "trend":
		// 检查趋势方向
		if expectedTrend, ok := condition.Expected.(float64); ok {
			actualTrend := ma.getCurrentTrend(condition) // 使用成员方法
			return (expectedTrend * actualTrend) > 0     // 同向为正
		}
	}

	return false
}

// 辅助函数

func calculateTimeScore(p MutationPrediction) float64 {
	elapsed := time.Since(p.Created)
	if elapsed > p.TimeFrame {
		return 0
	}
	return 1 - (float64(elapsed) / float64(p.TimeFrame))
}

func (ma *MutationAnalyzer) calculateConditionScore(p MutationPrediction) float64 {
	if len(p.Conditions) == 0 {
		return 1
	}

	metCount := 0
	for _, condition := range p.Conditions {
		if ma.isConditionMet(condition) {
			metCount++
		}
	}
	return float64(metCount) / float64(len(p.Conditions))
}

func calculateStabilityScore(p MutationPrediction) float64 {
	// 基于预测概率的稳定性评分
	baseScore := p.Probability

	// 考虑时间因素的衰减
	timeDecay := math.Exp(-float64(time.Since(p.Created)) / float64(p.TimeFrame))

	return baseScore * timeDecay
}

func (ma *MutationAnalyzer) getCurrentValue(condition PredictionCondition) float64 {
	// 根据条件类型获取对应的系统指标值
	switch condition.Type {
	case "energy":
		// 获取系统能量水平
		if value, ok := ma.detector.GetCurrentState().Properties["energy"].(float64); ok {
			return value
		}
	case "stability":
		// 获取系统稳定性
		if value, ok := ma.detector.GetCurrentState().Properties["stability"].(float64); ok {
			return value
		}
	case "frequency":
		// 获取模式频率
		if pattern, ok := ma.state.patterns[condition.PatternID]; ok {
			return pattern.Frequency
		}
	}
	return 0
}

func (ma *MutationAnalyzer) getCurrentState(condition PredictionCondition) string {
	// 获取系统当前状态
	state := ma.detector.GetCurrentState()

	switch condition.Type {
	case "phase":
		// 获取系统相位
		return string(state.Phase)
	case "cycle":
		// 获取系统周期类型
		if cycleType, ok := state.Properties["cycle_type"].(string); ok {
			return cycleType
		}
	case "pattern":
		// 获取模式状态
		if pattern, ok := ma.state.patterns[condition.PatternID]; ok {
			if len(pattern.Timeline) > 0 {
				return pattern.Timeline[len(pattern.Timeline)-1].Type
			}
		}
	}
	return ""
}

func (ma *MutationAnalyzer) getCurrentTrend(condition PredictionCondition) float64 {
	// 计算趋势需要历史数据
	history := ma.state.metrics.Performance
	if len(history) < 2 {
		return 0
	}

	// 获取指定指标的最近趋势
	switch condition.Type {
	case "energy":
		return calculateMetricTrend(history, "energy")
	case "stability":
		return calculateMetricTrend(history, "stability")
	case "pattern_frequency":
		if pattern, ok := ma.state.patterns[condition.PatternID]; ok {
			return calculatePatternTrend(pattern.Timeline)
		}
	}
	return 0
}

// 辅助函数：计算指标趋势
func calculateMetricTrend(history []PerformancePoint, metricName string) float64 {
	if len(history) < 2 {
		return 0
	}

	// 获取最近两个数据点
	current := history[len(history)-1].Metrics[metricName]
	previous := history[len(history)-2].Metrics[metricName]

	// 计算变化率
	timeDiff := history[len(history)-1].Time.Sub(history[len(history)-2].Time)
	valueDiff := current - previous

	// 归一化趋势值
	trend := valueDiff / float64(timeDiff.Seconds())

	// 限制趋势范围在 [-1, 1]
	return math.Max(-1, math.Min(1, trend))
}

// 辅助函数：计算模式趋势
func calculatePatternTrend(timeline []PatternEvent) float64 {
	if len(timeline) < 2 {
		return 0
	}

	// 分析最近的事件序列
	eventCount := 0
	for i := len(timeline) - 1; i >= 0 && i >= len(timeline)-5; i-- {
		if timeline[i].Type == "activation" {
			eventCount++
		}
	}

	// 计算趋势
	total := math.Min(5, float64(len(timeline)))
	return (float64(eventCount)/total)*2 - 1 // 映射到 [-1, 1]
}

func calculatePredictionProbability(p *MutationPattern, trend float64) float64 {
	return p.Frequency * (1 + trend)
}

func (ma *MutationAnalyzer) updateMetrics() {
	point := PerformancePoint{
		Time:    time.Now(),
		Metrics: make(map[string]float64),
	}

	// 计算性能指标
	point.Metrics["accuracy"] = ma.calculateAccuracy()
	point.Metrics["coverage"] = ma.calculateCoverage()
	point.Metrics["latency"] = float64(ma.state.metrics.Latency.Milliseconds())

	ma.state.metrics.Performance = append(ma.state.metrics.Performance, point)

	// 限制历史记录长度
	if len(ma.state.metrics.Performance) > maxMetricsHistory {
		ma.state.metrics.Performance = ma.state.metrics.Performance[1:]
	}
}

func generateAnalysisID() string {
	return fmt.Sprintf("ana_%d", time.Now().UnixNano())
}
