// system/evolution/pattern/analyzer.go

package pattern

import (
	"math"
	"time"

	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/common"
	"github.com/Corphon/daoflow/system/meta/field"
	"github.com/Corphon/daoflow/system/types"
)

// PatternAnalyzerImpl 模式分析器实现
type PatternAnalyzerImpl struct {
	// 分析配置
	config struct {
		minConfidence float64            // 最小置信度
		timeDecay     float64            // 时间衰减因子
		weightFactors map[string]float64 // 特征权重
	}

	// 分析状态
	state struct {
		patterns     map[string][]float64 // 模式历史分数
		lastAnalysis time.Time            // 最后分析时间
		metrics      AnalyzerMetrics      // 分析指标
	}
}

// AnalyzerMetrics 分析器指标
type AnalyzerMetrics struct {
	TotalAnalyzed int
	AverageScore  float64
	Accuracy      float64
	History       []types.MetricPoint
}

// NewPatternAnalyzer 创建新的模式分析器
func NewPatternAnalyzer() *PatternAnalyzerImpl {
	pa := &PatternAnalyzerImpl{}

	// 初始化配置
	pa.config.minConfidence = 0.6
	pa.config.timeDecay = 0.95
	pa.config.weightFactors = map[string]float64{
		"strength":  0.3,
		"stability": 0.3,
		"coherence": 0.2,
		"evolution": 0.2,
	}

	// 初始化状态
	pa.state.patterns = make(map[string][]float64)
	pa.state.lastAnalysis = time.Now()
	pa.state.metrics = AnalyzerMetrics{
		History: make([]types.MetricPoint, 0),
	}

	return pa
}

// AnalyzePattern 分析单个模式
func (pa *PatternAnalyzerImpl) AnalyzePattern(p common.SharedPattern) (float64, error) {
	if p == nil {
		return 0, model.WrapError(nil, model.ErrCodeValidation, "nil pattern")
	}

	// 1. 基础特征分析
	baseScore := pa.analyzeBaseFeatures(p)

	// 2. 时间维度分析
	timeScore := pa.analyzeTimeFeatures(p)

	// 3. 稳定性分析
	stabilityScore := pa.analyzeStability(p)

	// 4. 演化趋势分析
	evolutionScore := pa.analyzeEvolution(p)

	// 5. 整合分析结果
	finalScore := pa.integrateScores(map[string]float64{
		"base":      baseScore,
		"time":      timeScore,
		"stability": stabilityScore,
		"evolution": evolutionScore,
	})

	// 更新分析历史
	pa.updateAnalysisHistory(p.GetID(), finalScore)

	return finalScore, nil
}

// ComparePatterns 比较两个模式
func (pa *PatternAnalyzerImpl) ComparePatterns(p1, p2 common.SharedPattern) (float64, error) {
	if p1 == nil || p2 == nil {
		return 0, model.WrapError(nil, model.ErrCodeValidation, "nil pattern(s)")
	}

	// 1. 类型相似度
	typeSimilarity := pa.compareTypes(p1.GetType(), p2.GetType())

	// 2. 强度相似度
	strengthSimilarity := pa.compareStrengths(p1.GetStrength(), p2.GetStrength())

	// 3. 稳定性相似度
	stabilitySimilarity := pa.compareStability(p1.GetStability(), p2.GetStability())

	// 4. 时间关联度
	timeCorrelation := pa.calculateTimeCorrelation(p1.GetTimestamp(), p2.GetTimestamp())

	// 5. 整合比较结果
	similarity := pa.integrateSimilarity(map[string]float64{
		"type":      typeSimilarity,
		"strength":  strengthSimilarity,
		"stability": stabilitySimilarity,
		"time":      timeCorrelation,
	})

	return similarity, nil
}

// 内部分析方法

func (pa *PatternAnalyzerImpl) analyzeBaseFeatures(p common.SharedPattern) float64 {
	// 分析基础特征
	score := 0.0
	weights := pa.config.weightFactors

	// 强度评分
	strengthScore := normalizeValue(p.GetStrength())
	score += strengthScore * weights["strength"]

	// 稳定性评分
	stabilityScore := normalizeValue(p.GetStability())
	score += stabilityScore * weights["stability"]

	return score
}

func (pa *PatternAnalyzerImpl) analyzeTimeFeatures(p common.SharedPattern) float64 {
	// 分析时间特征
	age := time.Since(p.GetTimestamp())

	// 使用时间衰减函数
	timeScore := math.Exp(-float64(age.Hours()) / 24.0) // 24小时特征时间

	return normalizeValue(timeScore)
}

func (pa *PatternAnalyzerImpl) analyzeStability(p common.SharedPattern) float64 {
	stability := p.GetStability()

	// 考虑历史稳定性
	if history, exists := pa.state.patterns[p.GetID()]; exists {
		historicalStability := calculateHistoricalStability(history)
		stability = (stability + historicalStability) / 2
	}

	return normalizeValue(stability)
}

func (pa *PatternAnalyzerImpl) analyzeEvolution(p common.SharedPattern) float64 {
	// 分析演化趋势
	history := pa.state.patterns[p.GetID()]
	if len(history) < 2 {
		return 0.5 // 默认中性评分
	}

	// 计算趋势
	trend := calculateTrend(history)

	// 计算波动性
	volatility := calculateVolatility(history)

	return normalizeValue((trend + (1 - volatility)) / 2)
}

func (pa *PatternAnalyzerImpl) integrateScores(scores map[string]float64) float64 {
	totalScore := 0.0
	totalWeight := 0.0

	for feature, score := range scores {
		weight := pa.config.weightFactors[feature]
		totalScore += score * weight
		totalWeight += weight
	}

	if totalWeight == 0 {
		return 0
	}

	return normalizeValue(totalScore / totalWeight)
}

// 比较方法

func (pa *PatternAnalyzerImpl) compareTypes(type1, type2 string) float64 {
	if type1 == type2 {
		return 1.0
	}

	// 可以实现更复杂的类型相似度计算
	// 例如基于类型层次结构或特征相似度
	return 0.3
}

func (pa *PatternAnalyzerImpl) compareStrengths(strength1, strength2 float64) float64 {
	// 使用相对差异计算相似度
	diff := math.Abs(strength1 - strength2)
	maxStrength := math.Max(strength1, strength2)

	if maxStrength == 0 {
		return 1.0
	}

	return 1.0 - (diff / maxStrength)
}

func (pa *PatternAnalyzerImpl) compareStability(stability1, stability2 float64) float64 {
	// 类似强度比较
	diff := math.Abs(stability1 - stability2)
	return 1.0 - diff
}

func (pa *PatternAnalyzerImpl) calculateTimeCorrelation(t1, t2 time.Time) float64 {
	// 计算时间差异
	timeDiff := math.Abs(float64(t1.Sub(t2).Hours()))

	// 使用指数衰减函数
	return math.Exp(-timeDiff / 24.0) // 24小时特征时间
}

func (pa *PatternAnalyzerImpl) integrateSimilarity(similarities map[string]float64) float64 {
	return pa.integrateScores(similarities) // 复用分数整合逻辑
}

// 辅助函数

func (pa *PatternAnalyzerImpl) updateAnalysisHistory(patternID string, score float64) {
	history := pa.state.patterns[patternID]
	history = append(history, score)

	// 限制历史长度
	if len(history) > types.MaxHistoryLength {
		history = history[1:]
	}

	pa.state.patterns[patternID] = history

	// 更新指标
	pa.updateMetrics(score)
}

func (pa *PatternAnalyzerImpl) updateMetrics(score float64) {
	metrics := &pa.state.metrics
	metrics.TotalAnalyzed++

	// 更新平均分
	metrics.AverageScore = (metrics.AverageScore*float64(metrics.TotalAnalyzed-1) + score) /
		float64(metrics.TotalAnalyzed)

	// 添加新的指标点
	metrics.History = append(metrics.History, types.MetricPoint{
		Timestamp: time.Now(),
		Values: map[string]float64{
			"score":    score,
			"average":  metrics.AverageScore,
			"accuracy": metrics.Accuracy,
		},
	})

	// 限制历史记录
	if len(metrics.History) > types.MaxMetricsHistory {
		metrics.History = metrics.History[1:]
	}
}

// calculateHistoricalStability 历史稳定性
func calculateHistoricalStability(history []float64) float64 {
	if len(history) < 2 {
		return 1.0
	}

	mean := calculateMean(history)
	variance := calculateVariance(history, mean)
	return 1.0 - math.Min(1.0, variance)
}

// calculateVariance 计算方差
func calculateVariance(values []float64, mean float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sumSquares := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquares += diff * diff
	}
	return sumSquares / float64(len(values))
}

func calculateTrend(history []float64) float64 {
	if len(history) < 2 {
		return 0.5
	}

	// 简单线性回归
	x := make([]float64, len(history))
	for i := range x {
		x[i] = float64(i)
	}

	slope := field.CalculateLinearRegression(x, history)

	// 将斜率映射到[-1,1]区间
	normalized := math.Tanh(slope)

	// 映射到[0,1]区间
	return (normalized + 1.0) / 2.0
}

// calculateVolatility 计算波动性
func calculateVolatility(history []float64) float64 {
	if len(history) < 2 {
		return 0.0
	}

	mean := calculateMean(history)
	variance := calculateVariance(history, mean)
	return math.Min(1.0, variance)
}

func normalizeValue(value float64) float64 {
	return math.Max(0.0, math.Min(1.0, value))
}
