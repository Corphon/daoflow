//system/evolution/adaptation/learning.go

package adaptation

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
	"github.com/Corphon/daoflow/system/evolution/pattern"
	"github.com/Corphon/daoflow/system/types"
)

const (
	maxModelHistory = 100
)

// AdaptiveLearning 适应性学习系统
type AdaptiveLearning struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		learningRate    float64 // 学习率
		memoryCapacity  int     // 记忆容量
		explorationRate float64 // 探索率
		decayFactor     float64 // 衰减因子
	}

	// 学习状态
	state struct {
		knowledge          map[string]*KnowledgeUnit // 知识单元
		experiences        []LearningExperience      // 学习经验
		models             map[string]*LearningModel // 学习模型
		statistics         LearningStatistics        // 学习统计
		prevKnowledgeCount int                       // 上次知识数量
	}

	// 依赖项
	strategy *AdaptationStrategy
	matcher  *pattern.EvolutionMatcher
}

// KnowledgeUnit 知识单元
type KnowledgeUnit struct {
	ID           string            // 单元ID
	Type         string            // 知识类型
	Content      interface{}       // 知识内容
	Metadata     KnowledgeMetadata // 元数据
	Connections  []KnowledgeLink   // 知识关联
	ValidationFn func() bool       // 验证函数
	Created      time.Time         // 创建时间
}

// KnowledgeMetadata 知识元数据
type KnowledgeMetadata struct {
	Source     string    // 知识来源
	Confidence float64   // 置信度
	Usage      int       // 使用次数
	LastAccess time.Time // 最后访问
	Tags       []string  // 标签
}

// KnowledgeLink 知识关联
type KnowledgeLink struct {
	TargetID string                 // 目标ID
	Type     string                 // 关联类型
	Strength float64                // 关联强度
	Context  map[string]interface{} // 关联上下文
}

// LearningExperience 学习经验
type LearningExperience struct {
	ID        string                 // 经验ID
	Type      string                 // 经验类型
	Scenario  string                 // 场景描述
	Action    LearningAction         // 执行动作
	Result    LearningResult         // 执行结果
	Feedback  float64                // 反馈值
	Timestamp time.Time              // 记录时间
	Context   map[string]interface{} // 上下文信息
}

// ExperienceResult 经验结果
type ExperienceResult struct {
	Status  string                 // 执行状态
	Data    map[string]interface{} // 结果数据
	Metrics map[string]float64     // 结果指标
	Error   error                  // 错误信息
}

// LearningAction 学习动作
type LearningAction struct {
	Type       string                 // 动作类型
	Parameters map[string]interface{} // 动作参数
	Context    map[string]interface{} // 执行上下文
}

// LearningResult 学习结果
type LearningResult struct {
	Status   string             // 执行状态
	Outcome  interface{}        // 执行结果
	Metrics  map[string]float64 // 结果指标
	Duration time.Duration      // 执行时长
}

// LearningModel 学习模型
type LearningModel struct {
	ID          string                 // 模型ID
	Type        string                 // 模型类型
	Parameters  map[string]interface{} // 模型参数
	State       ModelState             // 模型状态
	Performance ModelPerformance       // 性能指标
}

// ModelState 模型状态
type ModelState struct {
	Version       int                // 版本号
	TrainingData  []TrainingItem     // 训练数据
	Weights       map[string]float64 // 模型权重
	LastUpdate    time.Time          // 最后更新
	LastLoss      float64            // 最后损失值
	Gradients     map[string]float64 // 梯度信息
	PrevGradients map[string]float64 // 前一次梯度(用于动量计算)
}

// ModelPerformance 模型性能
type ModelPerformance struct {
	Accuracy float64            // 准确率
	Loss     float64            // 损失值
	History  []PerformancePoint // 历史表现
	Details  TrainingDetails    // 训练细节
}

// PerformancePoint 性能记录点
type PerformancePoint struct {
	Time    time.Time          // 记录时间
	Metrics map[string]float64 // 性能指标
	Details struct {           // 详细信息
		BatchSize  int     // 批次大小
		Iterations int     // 迭代次数
		Duration   float64 // 训练时长
	}
}

// TrainingItem 训练项
type TrainingItem struct {
	Input  map[string]interface{} // 输入数据
	Output interface{}            // 期望输出
	Weight float64                // 样本权重
}

// LearningStatistics 学习统计
type LearningStatistics struct {
	TotalExperiences int                // 总经验数
	SuccessRate      float64            // 成功率
	KnowledgeGrowth  float64            // 知识增长率
	ModelAccuracy    map[string]float64 // 模型准确率
}

// PatternCondition 模式条件
type PatternCondition struct {
	Type   string      // 条件类型
	Key    string      // 条件键
	Value  interface{} // 条件值
	Weight float64     // 条件权重
}

// PatternOutcome 模式结果
type PatternOutcome struct {
	Type    string             // 结果类型
	Metrics map[string]float64 // 指标数据
	Weight  float64            // 结果权重
}

// TrainingDetails 训练详情
type TrainingDetails struct {
	BatchSize  int     // 批次大小
	Iterations int     // 迭代次数
	Duration   float64 // 训练时长
}

// RulePattern 规则模式
type RulePattern struct {
	Type       string
	Target     string
	Condition  RuleCondition
	Action     RuleAction
	Confidence float64
	Frequency  float64
}

// ParameterPattern 参数模式
type ParameterPattern struct {
	Type       string                 // 参数类型
	Parameters map[string]interface{} // 参数值
	Weight     float64                // 权重
}

// ExperiencePattern 添加Success字段
type ExperiencePattern struct {
	Type       string
	Confidence float64
	Frequency  float64
	Context    map[string]interface{}
	Conditions []PatternCondition
	Outcomes   []PatternOutcome
	Success    bool
}

// --------------------------------------------------------------------

// NewAdaptiveLearning 创建新的适应性学习系统
func NewAdaptiveLearning(matcher *pattern.EvolutionMatcher, config *types.AdaptationConfig) (*AdaptiveLearning, error) {
	if matcher == nil {
		return nil, fmt.Errorf("nil evolution matcher")
	}
	if config == nil {
		return nil, fmt.Errorf("nil adaptation config")
	}

	al := &AdaptiveLearning{
		matcher: matcher,
	}

	// 初始化配置和状态
	// ...

	return al, nil
}

// Learn 执行学习过程
func (al *AdaptiveLearning) Learn() error {
	al.mu.Lock()
	defer al.mu.Unlock()

	// 收集学习经验
	if err := al.collectExperiences(); err != nil {
		return err
	}

	// 更新知识库
	if err := al.updateKnowledge(); err != nil {
		return err
	}

	// 训练模型
	if err := al.trainModels(); err != nil {
		return err
	}

	// 应用学习成果
	if err := al.applyLearning(); err != nil {
		return err
	}

	// 更新统计信息
	al.updateStatistics()

	return nil
}

// updateStatistics 更新学习统计信息
func (al *AdaptiveLearning) updateStatistics() {
	stats := &al.state.statistics

	// 更新基础统计
	stats.TotalExperiences = len(al.state.experiences)

	// 计算成功率
	successCount := 0
	for _, exp := range al.state.experiences {
		if exp.Result.Status == "success" {
			successCount++
		}
	}
	if stats.TotalExperiences > 0 {
		stats.SuccessRate = float64(successCount) / float64(stats.TotalExperiences)
	}

	// 计算知识增长率
	currentKnowledge := len(al.state.knowledge)
	if al.state.prevKnowledgeCount > 0 {
		stats.KnowledgeGrowth = float64(currentKnowledge-al.state.prevKnowledgeCount) /
			float64(al.state.prevKnowledgeCount)
	}
	al.state.prevKnowledgeCount = currentKnowledge

	// 更新模型准确率
	for id, model := range al.state.models {
		stats.ModelAccuracy[id] = model.Performance.Accuracy
	}
}

// collectExperiences 收集学习经验
func (al *AdaptiveLearning) collectExperiences() error {
	// 获取最新策略执行结果
	results, err := al.strategy.GetRecentResults()
	if err != nil {
		return err
	}

	// 转换为学习经验
	for _, result := range results {
		experience := al.createExperience(result)
		al.addExperience(experience)
	}

	return nil
}

// GetLearningRate 获取当前学习率
func (al *AdaptiveLearning) GetLearningRate() float64 {
	al.mu.RLock()
	defer al.mu.RUnlock()

	// 返回当前学习率
	return al.config.learningRate
}

// UpdateLearningRate 更新学习率
func (al *AdaptiveLearning) UpdateLearningRate(baseRate float64) {
	al.mu.Lock()
	defer al.mu.Unlock()

	// 基于性能和经验调整学习率
	accuracy := 0.0
	for _, model := range al.state.models {
		accuracy += model.Performance.Accuracy
	}
	if len(al.state.models) > 0 {
		accuracy /= float64(len(al.state.models))
	}

	// 动态调整学习率
	if accuracy > 0.8 {
		baseRate *= 0.9 // 高准确度时降低学习率
	} else if accuracy < 0.5 {
		baseRate *= 1.1 // 低准确度时提高学习率
	}

	// 应用衰减因子
	al.config.learningRate = baseRate * al.config.decayFactor
}

// createExperience 创建学习经验
func (al *AdaptiveLearning) createExperience(event StrategyEvent) LearningExperience {
	experience := LearningExperience{
		ID:        fmt.Sprintf("exp_%d", time.Now().UnixNano()),
		Type:      "strategy_execution",
		Timestamp: event.Timestamp,
		Context:   make(map[string]interface{}),
		Result: LearningResult{ // 修改这里
			Status:   event.Status,
			Outcome:  event.Details,
			Metrics:  make(map[string]float64),
			Duration: time.Since(event.Timestamp),
		},
	}

	// 提取上下文信息
	if strategy, exists := al.strategy.state.strategies[event.StrategyID]; exists {
		experience.Context["strategy_type"] = strategy.Type
		experience.Context["strategy_params"] = strategy.Parameters
		experience.Context["effectiveness"] = strategy.Effectiveness
	}

	return experience
}

// updateKnowledge 更新知识库
func (al *AdaptiveLearning) updateKnowledge() error {
	// 分析新经验
	patterns := al.analyzeExperiences()

	// 提取知识
	for _, pattern := range patterns {
		knowledge := al.extractKnowledge(pattern)
		if knowledge != nil {
			al.integrateKnowledge(knowledge)
		}
	}

	// 验证知识有效性
	al.validateKnowledge()

	return nil
}

// analyzeExperiences 分析经验模式
func (al *AdaptiveLearning) analyzeExperiences() []ExperiencePattern {
	patterns := make([]ExperiencePattern, 0)

	// 提取最近的经验样本
	recentExperiences := al.state.experiences
	if len(recentExperiences) == 0 {
		return patterns
	}

	// 分组分析
	groupedExperiences := groupExperiencesByType(recentExperiences)
	for expType, experiences := range groupedExperiences {
		// 分析成功模式
		if pattern := analyzeSuccessPattern(experiences); pattern != nil {
			patterns = append(patterns, *pattern)
		}

		// 分析失败模式
		if pattern := analyzeFailurePattern(experiences); pattern != nil {
			patterns = append(patterns, *pattern)
		}

		// 分析适应模式
		if pattern := analyzeAdaptationPattern(expType, experiences); pattern != nil {
			patterns = append(patterns, *pattern)
		}
	}

	return patterns
}

// analyzeSuccessPattern 分析成功模式
func analyzeSuccessPattern(experiences []LearningExperience) *ExperiencePattern {
	if len(experiences) == 0 {
		return nil
	}

	pattern := &ExperiencePattern{
		Type:       "success",
		Confidence: calculatePatternConfidence(experiences),
		Frequency:  calculateSuccessFrequency(experiences),
		Context:    extractCommonContext(experiences),
		Conditions: extractSuccessConditions(experiences),
		Outcomes:   extractPositiveOutcomes(experiences),
	}

	// 验证模式有效性
	if !isValidPattern(pattern) {
		return nil
	}

	return pattern
}

// calculateSuccessFrequency 计算成功频率
func calculateSuccessFrequency(experiences []LearningExperience) float64 {
	if len(experiences) == 0 {
		return 0
	}

	successCount := 0
	for _, exp := range experiences {
		if exp.Result.Status == "success" {
			successCount++
		}
	}

	return float64(successCount) / float64(len(experiences))
}

// extractCommonContext 提取共同上下文
func extractCommonContext(experiences []LearningExperience) map[string]interface{} {
	if len(experiences) == 0 {
		return nil
	}

	// 以第一个经验的上下文作为基准
	common := make(map[string]interface{})
	for k, v := range experiences[0].Context {
		common[k] = v
	}

	// 保留共同的上下文项
	for _, exp := range experiences[1:] {
		for k, v := range common {
			if expVal, exists := exp.Context[k]; !exists || expVal != v {
				delete(common, k)
			}
		}
	}

	return common
}

// extractSuccessConditions 提取成功条件
func extractSuccessConditions(experiences []LearningExperience) []PatternCondition {
	conditions := make([]PatternCondition, 0)

	// 分析前置条件
	for _, exp := range experiences {
		if exp.Result.Status == "success" {
			for k, v := range exp.Context {
				if isSignificantCondition(k, v, experiences) {
					conditions = append(conditions, PatternCondition{
						Type:   "context",
						Key:    k,
						Value:  v,
						Weight: calculateConditionWeight(k, v, experiences),
					})
				}
			}
		}
	}

	return mergeSimilarConditions(conditions)
}

// mergeSimilarConditions 合并相似条件
func mergeSimilarConditions(conditions []PatternCondition) []PatternCondition {
	if len(conditions) <= 1 {
		return conditions
	}

	// 按Key分组的条件Map
	grouped := make(map[string][]PatternCondition)
	for _, cond := range conditions {
		grouped[cond.Key] = append(grouped[cond.Key], cond)
	}

	// 合并结果
	merged := make([]PatternCondition, 0)
	for key, group := range grouped {
		if len(group) == 1 {
			merged = append(merged, group[0])
			continue
		}

		// 计算组内平均权重
		totalWeight := 0.0
		for _, cond := range group {
			totalWeight += cond.Weight
		}
		avgWeight := totalWeight / float64(len(group))

		// 使用最高权重的值
		bestCond := group[0]
		for _, cond := range group[1:] {
			if cond.Weight > bestCond.Weight {
				bestCond = cond
			}
		}

		// 创建合并后的条件
		mergedCond := PatternCondition{
			Type:   bestCond.Type,
			Key:    key,
			Value:  bestCond.Value,
			Weight: avgWeight,
		}
		merged = append(merged, mergedCond)
	}

	return merged
}

// extractPositiveOutcomes 提取正向结果
func extractPositiveOutcomes(experiences []LearningExperience) []PatternOutcome {
	outcomes := make([]PatternOutcome, 0)

	// 分析成功经验的结果
	for _, exp := range experiences {
		if exp.Result.Status == "success" {
			if metrics := extractSignificantMetrics(exp.Result.Metrics); len(metrics) > 0 {
				outcomes = append(outcomes, PatternOutcome{
					Type:    "metrics",
					Metrics: metrics,
					Weight:  calculateOutcomeWeight(exp),
				})
			}
		}
	}

	return mergeRelatedOutcomes(outcomes)
}

// extractSignificantMetrics 提取显著指标
func extractSignificantMetrics(metrics map[string]float64) map[string]float64 {
	if len(metrics) == 0 {
		return nil
	}

	significant := make(map[string]float64)

	// 计算均值和标准差
	mean := calculateMetricsMean(metrics)
	stdDev := calculateMetricsStdDev(metrics, mean)

	// 提取显著指标(超过1个标准差)
	for key, value := range metrics {
		if math.Abs(value-mean) > stdDev {
			significant[key] = value
		}
	}

	return significant
}

// calculateOutcomeWeight 计算结果权重
func calculateOutcomeWeight(exp LearningExperience) float64 {
	// 基础权重
	weight := 1.0

	// 根据时间衰减调整
	age := time.Since(exp.Timestamp).Hours()
	timeDecay := math.Exp(-age / 24.0) // 24小时衰减
	weight *= timeDecay

	// 根据结果显著性调整
	if metrics := exp.Result.Metrics; len(metrics) > 0 {
		significance := calculateMetricsSignificance(metrics)
		weight *= significance
	}

	return weight
}

// mergeRelatedOutcomes 合并相关结果
func mergeRelatedOutcomes(outcomes []PatternOutcome) []PatternOutcome {
	if len(outcomes) <= 1 {
		return outcomes
	}

	// 按类型分组
	grouped := make(map[string][]PatternOutcome)
	for _, outcome := range outcomes {
		grouped[outcome.Type] = append(grouped[outcome.Type], outcome)
	}

	// 合并每组结果
	merged := make([]PatternOutcome, 0)
	for _, group := range grouped {
		if len(group) == 1 {
			merged = append(merged, group[0])
			continue
		}

		// 合并指标和权重
		mergedOutcome := PatternOutcome{
			Type:    group[0].Type,
			Metrics: mergeMetrics(group),
			Weight:  calculateAverageWeight(group),
		}
		merged = append(merged, mergedOutcome)
	}

	return merged
}

// 辅助函数
func calculateMetricsMean(metrics map[string]float64) float64 {
	total := 0.0
	for _, v := range metrics {
		total += v
	}
	return total / float64(len(metrics))
}

func calculateMetricsStdDev(metrics map[string]float64, mean float64) float64 {
	varSum := 0.0
	for _, v := range metrics {
		diff := v - mean
		varSum += diff * diff
	}
	return math.Sqrt(varSum / float64(len(metrics)))
}

func calculateMetricsSignificance(metrics map[string]float64) float64 {
	if len(metrics) == 0 {
		return 0
	}
	mean := calculateMetricsMean(metrics)
	stdDev := calculateMetricsStdDev(metrics, mean)
	return math.Min(1.0, stdDev/mean)
}

func mergeMetrics(outcomes []PatternOutcome) map[string]float64 {
	merged := make(map[string]float64)
	weights := make(map[string]float64)

	for _, outcome := range outcomes {
		for k, v := range outcome.Metrics {
			merged[k] += v * outcome.Weight
			weights[k] += outcome.Weight
		}
	}

	// 归一化
	for k := range merged {
		if weights[k] > 0 {
			merged[k] /= weights[k]
		}
	}

	return merged
}

func calculateAverageWeight(outcomes []PatternOutcome) float64 {
	total := 0.0
	for _, outcome := range outcomes {
		total += outcome.Weight
	}
	return total / float64(len(outcomes))
}

// 辅助函数
func isSignificantCondition(key string, value interface{}, experiences []LearningExperience) bool {
	successCount := 0
	totalCount := 0

	for _, exp := range experiences {
		if v, exists := exp.Context[key]; exists && v == value {
			if exp.Result.Status == "success" {
				successCount++
			}
			totalCount++
		}
	}

	return totalCount > 0 && float64(successCount)/float64(totalCount) >= 0.7
}

func calculateConditionWeight(key string, value interface{}, experiences []LearningExperience) float64 {
	successCount := 0
	totalCount := 0

	for _, exp := range experiences {
		if v, exists := exp.Context[key]; exists && v == value {
			if exp.Result.Status == "success" {
				successCount++
			}
			totalCount++
		}
	}

	return float64(successCount) / float64(totalCount)
}

// analyzeFailurePattern 分析失败模式
func analyzeFailurePattern(experiences []LearningExperience) *ExperiencePattern {
	if len(experiences) == 0 {
		return nil
	}

	pattern := &ExperiencePattern{
		Type:       "failure",
		Confidence: calculatePatternConfidence(experiences),
		Frequency:  calculateFailureFrequency(experiences),
		Context:    extractCommonContext(experiences),
		Conditions: extractFailureConditions(experiences),
		Outcomes:   extractNegativeOutcomes(experiences),
	}

	// 验证模式有效性
	if !isValidPattern(pattern) {
		return nil
	}

	return pattern
}

// calculateFailureFrequency 计算失败频率
func calculateFailureFrequency(experiences []LearningExperience) float64 {
	if len(experiences) == 0 {
		return 0
	}

	failureCount := 0
	for _, exp := range experiences {
		if exp.Result.Status == "failure" {
			failureCount++
		}
	}

	return float64(failureCount) / float64(len(experiences))
}

// extractFailureConditions 提取失败条件
func extractFailureConditions(experiences []LearningExperience) []PatternCondition {
	conditions := make([]PatternCondition, 0)

	// 分析失败前置条件
	for _, exp := range experiences {
		if exp.Result.Status == "failure" {
			for k, v := range exp.Context {
				if isSignificantCondition(k, v, experiences) {
					conditions = append(conditions, PatternCondition{
						Type:   "context",
						Key:    k,
						Value:  v,
						Weight: calculateConditionWeight(k, v, experiences),
					})
				}
			}
		}
	}

	return mergeSimilarConditions(conditions)
}

// extractNegativeOutcomes 提取负面结果
func extractNegativeOutcomes(experiences []LearningExperience) []PatternOutcome {
	outcomes := make([]PatternOutcome, 0)

	// 分析失败经验的结果
	for _, exp := range experiences {
		if exp.Result.Status == "failure" {
			if metrics := extractSignificantMetrics(exp.Result.Metrics); len(metrics) > 0 {
				outcomes = append(outcomes, PatternOutcome{
					Type:    "metrics",
					Metrics: metrics,
					Weight:  calculateOutcomeWeight(exp),
				})
			}
		}
	}

	return mergeRelatedOutcomes(outcomes)
}

// analyzeAdaptationPattern 分析适应模式
func analyzeAdaptationPattern(expType string, experiences []LearningExperience) *ExperiencePattern {
	if len(experiences) == 0 {
		return nil
	}

	pattern := &ExperiencePattern{
		Type:       "adaptation",
		Confidence: calculateAdaptationConfidence(experiences),
		Frequency:  calculateAdaptationFrequency(experiences, expType),
		Context:    extractAdaptationContext(experiences),
		Conditions: extractAdaptationConditions(experiences),
		Outcomes:   extractAdaptationOutcomes(experiences),
	}

	// 验证模式有效性
	if !isValidPattern(pattern) {
		return nil
	}

	return pattern
}

// calculateAdaptationConfidence 计算适应置信度
func calculateAdaptationConfidence(experiences []LearningExperience) float64 {
	if len(experiences) == 0 {
		return 0
	}

	// 计算适应成功率
	successCount := 0
	totalCount := 0
	for _, exp := range experiences {
		if isAdaptationSuccess(exp) {
			successCount++
		}
		totalCount++
	}

	return float64(successCount) / float64(totalCount)
}

// calculateAdaptationFrequency 计算适应频率
func calculateAdaptationFrequency(experiences []LearningExperience, expType string) float64 {
	if len(experiences) == 0 {
		return 0
	}

	adaptCount := 0
	for _, exp := range experiences {
		if exp.Type == expType && isAdaptiveAction(exp) {
			adaptCount++
		}
	}

	return float64(adaptCount) / float64(len(experiences))
}

// extractAdaptationContext 提取适应上下文
func extractAdaptationContext(experiences []LearningExperience) map[string]interface{} {
	context := make(map[string]interface{})

	// 提取环境状态
	context["environment"] = extractEnvironmentState(experiences)

	// 提取触发条件
	context["triggers"] = extractTriggerConditions(experiences)

	// 提取适应策略
	context["strategies"] = extractAdaptationStrategies(experiences)

	return context
}

// extractAdaptationConditions 提取适应条件
func extractAdaptationConditions(experiences []LearningExperience) []PatternCondition {
	conditions := make([]PatternCondition, 0)

	// 分析适应前置条件
	for _, exp := range experiences {
		if isAdaptiveAction(exp) {
			for k, v := range exp.Context {
				if isSignificantCondition(k, v, experiences) {
					conditions = append(conditions, PatternCondition{
						Type:   "adaptation",
						Key:    k,
						Value:  v,
						Weight: calculateConditionWeight(k, v, experiences),
					})
				}
			}
		}
	}

	return mergeSimilarConditions(conditions)
}

// extractAdaptationOutcomes 提取适应结果
func extractAdaptationOutcomes(experiences []LearningExperience) []PatternOutcome {
	outcomes := make([]PatternOutcome, 0)

	// 分析适应结果
	for _, exp := range experiences {
		if isAdaptiveAction(exp) && isAdaptationSuccess(exp) {
			if metrics := extractSignificantMetrics(exp.Result.Metrics); len(metrics) > 0 {
				outcomes = append(outcomes, PatternOutcome{
					Type:    "adaptation",
					Metrics: metrics,
					Weight:  calculateOutcomeWeight(exp),
				})
			}
		}
	}

	return mergeRelatedOutcomes(outcomes)
}

// 辅助函数
func isAdaptiveAction(exp LearningExperience) bool {
	return exp.Action.Type == "adaptation"
}

func isAdaptationSuccess(exp LearningExperience) bool {
	return exp.Result.Status == "success" && isAdaptiveAction(exp)
}

func extractEnvironmentState(experiences []LearningExperience) map[string]float64 {
	state := make(map[string]float64)
	count := make(map[string]int)

	for _, exp := range experiences {
		for k, v := range exp.Context {
			if val, ok := v.(float64); ok {
				state[k] += val
				count[k]++
			}
		}
	}

	// 计算平均值
	for k := range state {
		if count[k] > 0 {
			state[k] /= float64(count[k])
		}
	}

	return state
}

func extractTriggerConditions(experiences []LearningExperience) []string {
	triggers := make(map[string]bool)
	for _, exp := range experiences {
		if trigger, ok := exp.Context["trigger"].(string); ok {
			triggers[trigger] = true
		}
	}

	result := make([]string, 0)
	for trigger := range triggers {
		result = append(result, trigger)
	}
	return result
}

func extractAdaptationStrategies(experiences []LearningExperience) []string {
	strategies := make(map[string]bool)
	for _, exp := range experiences {
		if strategy, ok := exp.Context["strategy"].(string); ok {
			strategies[strategy] = true
		}
	}

	result := make([]string, 0)
	for strategy := range strategies {
		result = append(result, strategy)
	}
	return result
}

// 辅助函数
func calculatePatternConfidence(experiences []LearningExperience) float64 {
	total := 0.0
	for _, exp := range experiences {
		if exp.Result.Status == "success" {
			total += 1.0
		}
	}
	return total / float64(len(experiences))
}

func isValidPattern(pattern *ExperiencePattern) bool {
	return pattern.Confidence >= 0.3 && // 最小置信度
		len(pattern.Conditions) > 0 && // 至少有一个条件
		len(pattern.Outcomes) > 0 // 至少有一个结果
}

// extractKnowledge 从经验模式提取知识
func (al *AdaptiveLearning) extractKnowledge(pattern ExperiencePattern) *KnowledgeUnit {
	knowledge := &KnowledgeUnit{
		ID:      generateKnowledgeID(),
		Type:    pattern.Type,
		Content: pattern,
		Metadata: KnowledgeMetadata{
			Source:     "experience_analysis",
			Confidence: pattern.Confidence,
			Usage:      0,
			LastAccess: time.Now(),
			Tags:       []string{pattern.Type, "auto_generated"},
		},
		Created: time.Now(),
	}

	// 添加验证函数
	knowledge.ValidationFn = func() bool {
		return validatePatternKnowledge(pattern)
	}

	// 建立知识关联
	knowledge.Connections = al.findKnowledgeConnections(pattern)

	return knowledge
}

// validatePatternKnowledge 验证模式知识有效性
func validatePatternKnowledge(pattern ExperiencePattern) bool {
	// 1. 检查基本有效性
	if pattern.Confidence < 0.3 || pattern.Frequency < 0.1 {
		return false
	}

	// 2. 检查条件完整性
	if len(pattern.Conditions) == 0 {
		return false
	}

	// 3. 检查结果有效性
	if len(pattern.Outcomes) == 0 {
		return false
	}

	// 4. 验证上下文
	if pattern.Context == nil {
		return false
	}

	return true
}

// findKnowledgeConnections 查找知识关联
func (al *AdaptiveLearning) findKnowledgeConnections(pattern ExperiencePattern) []KnowledgeLink {
	connections := make([]KnowledgeLink, 0)

	// 遍历现有知识
	for id, existing := range al.state.knowledge {
		// 跳过自身
		if existing.Type == pattern.Type {
			continue
		}

		// 1. 检查条件关联
		if relationScore := compareConditions(pattern.Conditions, existing); relationScore > 0.7 {
			connections = append(connections, KnowledgeLink{
				TargetID: id,
				Type:     "condition_related",
				Strength: relationScore,
				Context: map[string]interface{}{
					"relation_type": "condition",
					"score":         relationScore,
				},
			})
		}

		// 2. 检查结果关联
		if relationScore := compareOutcomes(pattern.Outcomes, existing); relationScore > 0.7 {
			connections = append(connections, KnowledgeLink{
				TargetID: id,
				Type:     "outcome_related",
				Strength: relationScore,
				Context: map[string]interface{}{
					"relation_type": "outcome",
					"score":         relationScore,
				},
			})
		}

		// 3. 检查上下文关联
		if relationScore := compareContexts(pattern.Context, existing); relationScore > 0.7 {
			connections = append(connections, KnowledgeLink{
				TargetID: id,
				Type:     "context_related",
				Strength: relationScore,
				Context: map[string]interface{}{
					"relation_type": "context",
					"score":         relationScore,
				},
			})
		}
	}

	return connections
}

// 辅助函数
func compareConditions(conditions []PatternCondition, knowledge *KnowledgeUnit) float64 {
	if knowledge.Content == nil {
		return 0
	}

	if existingPattern, ok := knowledge.Content.(ExperiencePattern); ok {
		matches := 0
		for _, c1 := range conditions {
			for _, c2 := range existingPattern.Conditions {
				if c1.Key == c2.Key && c1.Value == c2.Value {
					matches++
					break
				}
			}
		}

		totalConditions := math.Max(float64(len(conditions)),
			float64(len(existingPattern.Conditions)))
		if totalConditions > 0 {
			return float64(matches) / totalConditions
		}
	}
	return 0
}

func compareOutcomes(outcomes []PatternOutcome, knowledge *KnowledgeUnit) float64 {
	if knowledge.Content == nil {
		return 0
	}

	if existingPattern, ok := knowledge.Content.(ExperiencePattern); ok {
		matches := 0
		for _, o1 := range outcomes {
			for _, o2 := range existingPattern.Outcomes {
				if compareMetrics(o1.Metrics, o2.Metrics) > 0.8 {
					matches++
					break
				}
			}
		}

		totalOutcomes := math.Max(float64(len(outcomes)),
			float64(len(existingPattern.Outcomes)))
		if totalOutcomes > 0 {
			return float64(matches) / totalOutcomes
		}
	}
	return 0
}

func compareContexts(context1 map[string]interface{}, knowledge *KnowledgeUnit) float64 {
	if knowledge.Content == nil {
		return 0
	}

	if existingPattern, ok := knowledge.Content.(ExperiencePattern); ok {
		matches := 0
		totalKeys := 0

		for k1, v1 := range context1 {
			totalKeys++
			if v2, exists := existingPattern.Context[k1]; exists && v1 == v2 {
				matches++
			}
		}

		for k := range existingPattern.Context {
			if _, exists := context1[k]; !exists {
				totalKeys++
			}
		}

		if totalKeys > 0 {
			return float64(matches) / float64(totalKeys)
		}
	}
	return 0
}

func compareMetrics(m1, m2 map[string]float64) float64 {
	matches := 0
	totalMetrics := 0

	for k1, v1 := range m1 {
		totalMetrics++
		if v2, exists := m2[k1]; exists && math.Abs(v1-v2) < 0.1 {
			matches++
		}
	}

	for k := range m2 {
		if _, exists := m1[k]; !exists {
			totalMetrics++
		}
	}

	if totalMetrics > 0 {
		return float64(matches) / float64(totalMetrics)
	}
	return 0
}

// validateKnowledge 验证知识有效性
func (al *AdaptiveLearning) validateKnowledge() {
	for id, knowledge := range al.state.knowledge {
		// 跳过新知识
		if time.Since(knowledge.Created) < 24*time.Hour {
			continue
		}

		// 验证知识
		if knowledge.ValidationFn != nil && !knowledge.ValidationFn() {
			// 降低置信度
			knowledge.Metadata.Confidence *= 0.9

			// 如果置信度太低，删除知识
			if knowledge.Metadata.Confidence < 0.3 {
				delete(al.state.knowledge, id)
			}
		}
	}
}

func groupExperiencesByType(experiences []LearningExperience) map[string][]LearningExperience {
	grouped := make(map[string][]LearningExperience)
	for _, exp := range experiences {
		grouped[exp.Type] = append(grouped[exp.Type], exp)
	}
	return grouped
}

// trainModels 训练模型
func (al *AdaptiveLearning) trainModels() error {
	for _, model := range al.state.models {
		// 准备训练数据
		trainingData := al.prepareTrainingData(model)

		// 执行训练
		if err := al.trainModel(model, trainingData); err != nil {
			continue
		}

		// 评估模型性能
		al.evaluateModel(model)
	}

	return nil
}

// prepareTrainingData 准备训练数据
func (al *AdaptiveLearning) prepareTrainingData(model *LearningModel) []TrainingItem {
	trainingData := make([]TrainingItem, 0)

	// 从经验中提取训练样本
	for _, exp := range al.state.experiences {
		if item := convertExperienceToTraining(exp, model.Type); item != nil {
			trainingData = append(trainingData, *item)
		}
	}

	// 从知识库中补充样本
	for _, knowledge := range al.state.knowledge {
		if items := extractTrainingFromKnowledge(knowledge, model.Type); len(items) > 0 {
			trainingData = append(trainingData, items...)
		}
	}

	return trainingData
}

// trainModel 执行模型训练
func (al *AdaptiveLearning) trainModel(model *LearningModel, data []TrainingItem) error {
	if len(data) == 0 {
		return fmt.Errorf("no training data")
	}

	// 更新训练状态
	model.State.Version++
	model.State.TrainingData = data
	model.State.LastUpdate = time.Now()

	// 配置训练参数
	batchSize := calculateBatchSize(len(data))
	iterations := calculateIterations(len(data))

	// 执行训练
	startTime := time.Now()
	for i := 0; i < iterations; i++ {
		batch := selectBatch(data, batchSize)
		if err := trainBatch(model, batch); err != nil {
			return err
		}
		updateModelWeights(model)
	}

	// 记录训练详情
	model.Performance.Details.BatchSize = batchSize
	model.Performance.Details.Iterations = iterations
	model.Performance.Details.Duration = time.Since(startTime).Seconds()

	return nil
}

// evaluateModel 评估模型性能
func (al *AdaptiveLearning) evaluateModel(model *LearningModel) {
	// 更新准确率
	model.Performance.Accuracy = calculateModelAccuracy(model)

	// 更新损失值
	model.Performance.Loss = calculateModelLoss(model)

	// 记录性能历史
	point := PerformancePoint{
		Time: time.Now(),
		Metrics: map[string]float64{
			"accuracy": model.Performance.Accuracy,
			"loss":     model.Performance.Loss,
		},
		Details: model.Performance.Details,
	}

	// 维护历史记录长度
	model.Performance.History = append(model.Performance.History, point)
	if len(model.Performance.History) > maxModelHistory {
		model.Performance.History = model.Performance.History[1:]
	}
}

// 辅助函数
func convertExperienceToTraining(exp LearningExperience, modelType string) *TrainingItem {
	switch modelType {
	case "pattern":
		return convertToPatternTraining(exp)
	case "strategy":
		return convertToStrategyTraining(exp)
	default:
		return nil
	}
}

func extractTrainingFromKnowledge(k *KnowledgeUnit, modelType string) []TrainingItem {
	items := make([]TrainingItem, 0)

	// 根据不同的模型类型处理知识单元
	switch modelType {
	case "pattern":
		if pattern, ok := k.Content.(ExperiencePattern); ok {
			items = append(items, createPatternTrainingItems(pattern)...)
		}

	case "rule":
		if pattern, ok := k.Content.(RulePattern); ok {
			items = append(items, createRuleTrainingItems(pattern)...)
		}

	case "parameter":
		if pattern, ok := k.Content.(ParameterPattern); ok {
			items = append(items, createParameterTrainingItems(pattern)...)
		}
	}

	return items
}

// 针对规则模式创建训练项
func createRuleTrainingItems(pattern RulePattern) []TrainingItem {
	items := make([]TrainingItem, 0)

	// 从规则条件创建训练样本
	input := make(map[string]interface{})
	input["type"] = pattern.Type
	input["target"] = pattern.Target
	input["condition"] = pattern.Condition

	items = append(items, TrainingItem{
		Input:  input,
		Output: pattern.Action,
		Weight: pattern.Confidence,
	})

	return items
}

// 针对参数模式创建训练项
func createParameterTrainingItems(pattern ParameterPattern) []TrainingItem {
	items := make([]TrainingItem, 0)

	input := make(map[string]interface{})
	input["type"] = pattern.Type
	for k, v := range pattern.Parameters {
		input[k] = v
	}

	items = append(items, TrainingItem{
		Input:  input,
		Output: true,
		Weight: pattern.Weight,
	})

	return items
}

// convertToPatternTraining 转换经验到模式训练项
func convertToPatternTraining(exp LearningExperience) *TrainingItem {
	if exp.Type != "pattern" {
		return nil
	}

	// 提取输入特征
	input := make(map[string]interface{})
	for k, v := range exp.Context {
		input[k] = v
	}

	// 添加结果指标作为特征
	for k, v := range exp.Result.Metrics {
		input["metric_"+k] = v
	}

	return &TrainingItem{
		Input:  input,
		Output: exp.Result.Status == "success",
		Weight: calculateExperienceWeight(exp),
	}
}

// convertToStrategyTraining 转换经验到策略训练项
func convertToStrategyTraining(exp LearningExperience) *TrainingItem {
	if exp.Type != "strategy" {
		return nil
	}

	// 提取策略参数作为输入
	input := make(map[string]interface{})
	if params, ok := exp.Context["strategy_params"].(map[string]interface{}); ok {
		for k, v := range params {
			input[k] = v
		}
	}

	// 提取环境状态
	if state, ok := exp.Context["environment"].(map[string]interface{}); ok {
		for k, v := range state {
			input["env_"+k] = v
		}
	}

	return &TrainingItem{
		Input:  input,
		Output: exp.Result.Status == "success",
		Weight: calculateExperienceWeight(exp),
	}
}

// createPatternTrainingItems 从经验模式创建训练项
func createPatternTrainingItems(pattern ExperiencePattern) []TrainingItem {
	items := make([]TrainingItem, 0)

	// 从条件创建正例
	for _, cond := range pattern.Conditions {
		input := make(map[string]interface{})
		input[cond.Key] = cond.Value

		items = append(items, TrainingItem{
			Input:  input,
			Output: true,
			Weight: cond.Weight * pattern.Confidence,
		})
	}

	// 从结果创建训练样本
	for _, outcome := range pattern.Outcomes {
		input := make(map[string]interface{})
		for k, v := range outcome.Metrics {
			input[k] = v
		}

		items = append(items, TrainingItem{
			Input:  input,
			Output: pattern.Type == "success",
			Weight: outcome.Weight * pattern.Confidence,
		})
	}

	return items
}

// 辅助函数
func calculateExperienceWeight(exp LearningExperience) float64 {
	// 基础权重
	weight := 1.0

	// 根据时间衰减调整
	age := time.Since(exp.Timestamp).Hours()
	timeDecay := math.Exp(-age / 24.0) // 24小时衰减
	weight *= timeDecay

	// 根据结果可信度调整
	if metrics := exp.Result.Metrics; len(metrics) > 0 {
		if confidence, ok := metrics["confidence"]; ok {
			weight *= confidence
		}
	}

	return weight
}
func calculateBatchSize(dataSize int) int {
	return min(32, max(1, dataSize/10))
}

func calculateIterations(dataSize int) int {
	return min(1000, max(10, dataSize/32*3))
}

func selectBatch(data []TrainingItem, batchSize int) []TrainingItem {
	batch := make([]TrainingItem, 0, batchSize)
	for i := 0; i < batchSize; i++ {
		idx := rand.Intn(len(data))
		batch = append(batch, data[idx])
	}
	return batch
}

// trainBatch 执行批次训练
func trainBatch(model *LearningModel, batch []TrainingItem) error {
	// 1. 前向传播
	predictions := make([]float64, len(batch))
	for i, item := range batch {
		// 计算预测值
		pred, err := forwardPropagate(model, item.Input)
		if err != nil {
			return fmt.Errorf("forward propagation failed: %v", err)
		}
		predictions[i] = pred
	}

	// 2. 计算损失
	batchLoss := 0.0
	for i, item := range batch {
		expected := getExpectedValue(item.Output)
		loss := calculateItemLoss(predictions[i], expected)
		batchLoss += loss * item.Weight
	}
	batchLoss /= float64(len(batch))

	// 3. 反向传播
	gradients := make(map[string]float64)
	for i, item := range batch {
		itemGrads := backPropagate(model, item.Input, predictions[i],
			getExpectedValue(item.Output))
		// 累积梯度
		for key, grad := range itemGrads {
			gradients[key] += grad * item.Weight
		}
	}

	// 4. 更新模型状态
	model.State.LastLoss = batchLoss
	model.State.Gradients = gradients

	return nil
}

// updateModelWeights 更新模型权重
func updateModelWeights(model *LearningModel) {
	learningRate := 0.01 // 基础学习率

	// 1. 应用动量
	momentum := 0.9
	if model.State.PrevGradients != nil {
		for key := range model.State.Weights {
			model.State.Weights[key] -= learningRate * ((1-momentum)*model.State.Gradients[key] +
				momentum*model.State.PrevGradients[key])
		}
	} else {
		for key := range model.State.Weights {
			model.State.Weights[key] -= learningRate * model.State.Gradients[key]
		}
	}

	// 2. 保存当前梯度
	model.State.PrevGradients = model.State.Gradients

	// 3. L2正则化
	lambda := 0.01
	for key := range model.State.Weights {
		model.State.Weights[key] *= (1 - learningRate*lambda)
	}
}

// calculateModelAccuracy 计算模型准确率
func calculateModelAccuracy(model *LearningModel) float64 {
	if len(model.State.TrainingData) == 0 {
		return 0
	}

	correctCount := 0
	totalCount := 0

	for _, item := range model.State.TrainingData {
		// 获取预测值
		pred, err := forwardPropagate(model, item.Input)
		if err != nil {
			continue
		}

		// 比较预测值和实际值
		expected := getExpectedValue(item.Output)
		if isCorrectPrediction(pred, expected) {
			correctCount++
		}
		totalCount++
	}

	if totalCount == 0 {
		return 0
	}
	return float64(correctCount) / float64(totalCount)
}

// calculateModelLoss 计算模型损失值
func calculateModelLoss(model *LearningModel) float64 {
	if len(model.State.TrainingData) == 0 {
		return 1.0
	}

	totalLoss := 0.0
	totalWeight := 0.0

	for _, item := range model.State.TrainingData {
		// 获取预测值
		pred, err := forwardPropagate(model, item.Input)
		if err != nil {
			continue
		}

		// 计算加权损失
		expected := getExpectedValue(item.Output)
		loss := calculateItemLoss(pred, expected)
		totalLoss += loss * item.Weight
		totalWeight += item.Weight
	}

	if totalWeight == 0 {
		return 1.0
	}
	return totalLoss / totalWeight
}

// 辅助函数
func forwardPropagate(model *LearningModel, input map[string]interface{}) (float64, error) {
	// 转换输入特征为向量
	features := make([]float64, len(model.State.Weights))
	for i, key := range getSortedKeys(model.State.Weights) {
		if val, ok := input[key]; ok {
			if fVal, ok := val.(float64); ok {
				features[i] = fVal
			}
		}
	}

	// 计算加权和
	sum := 0.0
	for i, feature := range features {
		sum += feature * model.State.Weights[getSortedKeys(model.State.Weights)[i]]
	}

	// 应用激活函数(sigmoid)
	return 1.0 / (1.0 + math.Exp(-sum)), nil
}

func backPropagate(model *LearningModel, input map[string]interface{},
	prediction, expected float64) map[string]float64 {
	gradients := make(map[string]float64)

	// 计算输出层梯度
	outputGrad := 2 * (prediction - expected)

	// 计算每个权重的梯度
	for key := range model.State.Weights {
		if val, ok := input[key]; ok {
			if fVal, ok := val.(float64); ok {
				gradients[key] = outputGrad * fVal
			}
		}
	}

	return gradients
}

func calculateItemLoss(prediction, expected float64) float64 {
	diff := prediction - expected
	return diff * diff // 均方误差
}

func getExpectedValue(output interface{}) float64 {
	switch v := output.(type) {
	case float64:
		return v
	case bool:
		if v {
			return 1.0
		}
		return 0.0
	default:
		return 0.0
	}
}

func isCorrectPrediction(prediction, expected float64) bool {
	threshold := 0.5
	return (prediction >= threshold && expected >= threshold) ||
		(prediction < threshold && expected < threshold)
}

func getSortedKeys(m map[string]float64) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// applyLearning 应用学习成果
func (al *AdaptiveLearning) applyLearning() error {
	// 更新策略参数
	if err := al.updateStrategyParameters(); err != nil {
		return err
	}

	// 生成新规则
	if err := al.generateNewRules(); err != nil {
		return err
	}

	// 优化现有规则
	if err := al.optimizeRules(); err != nil {
		return err
	}

	return nil
}

// updateStrategyParameters 更新策略参数
func (al *AdaptiveLearning) updateStrategyParameters() error {
	// 分析经验数据
	patterns := al.analyzeExperiences()

	// 提取成功经验的参数模式
	successParams := extractSuccessParameters(patterns)

	// 更新策略参数
	for _, pattern := range successParams {
		if err := al.strategy.UpdateParameters(pattern.Type, pattern.Parameters); err != nil {
			continue
		}
	}

	return nil
}

// generateNewRules 生成新规则
func (al *AdaptiveLearning) generateNewRules() error {
	// 从经验中提取规则模式
	patterns := al.analyzeRulePatterns()

	// 生成新规则
	for _, pattern := range patterns {
		rule := &StrategyRule{
			ID:        core.GenerateID(), //ID生成函数
			Type:      pattern.Type,
			Target:    pattern.Target,
			Condition: pattern.Condition,
			Action:    pattern.Action,
			Weight:    calculateRuleWeight(pattern),
		}

		// 注册新规则
		if err := al.strategy.RegisterRule(rule); err != nil {
			continue
		}
	}

	return nil
}

// analyzeRulePatterns 分析规则模式
func (al *AdaptiveLearning) analyzeRulePatterns() []RulePattern {
	patterns := make([]RulePattern, 0)

	// 从经验中提取规则模式
	groupedExp := groupExperiencesByType(al.state.experiences)

	for expType, experiences := range groupedExp {
		// 分析成功规则模式
		if pattern := analyzeSuccessRulePattern(experiences); pattern != nil {
			patterns = append(patterns, *pattern)
		}

		// 分析失败规则模式
		if pattern := analyzeFailureRulePattern(experiences); pattern != nil {
			patterns = append(patterns, *pattern)
		}

		// 分析适应规则模式
		if pattern := analyzeAdaptationRulePattern(expType, experiences); pattern != nil {
			patterns = append(patterns, *pattern)
		}
	}

	return patterns
}

// analyzeSuccessRulePattern 分析成功规则模式
func analyzeSuccessRulePattern(experiences []LearningExperience) *RulePattern {
	if len(experiences) == 0 {
		return nil
	}

	// 提取成功规则条件
	condition := RuleCondition{
		Expression: "success_rate > threshold",
		Parameters: extractRuleParameters(experiences, "success"),
		Threshold:  0.7,
	}

	// 提取规则动作
	action := RuleAction{
		Function: "adjust_strategy",
		Parameters: map[string]interface{}{
			"confidence": calculateSuccessConfidence(experiences),
			"direction":  1.0,
		},
		ResultType: "strategy_adjustment",
	}

	return &RulePattern{
		Type:       "success",
		Target:     "strategy",
		Condition:  condition,
		Action:     action,
		Confidence: calculatePatternConfidence(experiences),
		Frequency:  calculateSuccessFrequency(experiences),
	}
}

// analyzeFailureRulePattern 分析失败规则模式
func analyzeFailureRulePattern(experiences []LearningExperience) *RulePattern {
	if len(experiences) == 0 {
		return nil
	}

	// 提取失败规则条件
	condition := RuleCondition{
		Expression: "failure_rate > threshold",
		Parameters: extractRuleParameters(experiences, "failure"),
		Threshold:  0.5,
	}

	// 提取规则动作
	action := RuleAction{
		Function: "adjust_strategy",
		Parameters: map[string]interface{}{
			"confidence": calculateFailureConfidence(experiences),
			"direction":  -1.0,
		},
		ResultType: "strategy_adjustment",
	}

	return &RulePattern{
		Type:       "failure",
		Target:     "strategy",
		Condition:  condition,
		Action:     action,
		Confidence: calculatePatternConfidence(experiences),
		Frequency:  calculateFailureFrequency(experiences),
	}
}

// analyzeAdaptationRulePattern 分析适应规则模式
func analyzeAdaptationRulePattern(expType string, experiences []LearningExperience) *RulePattern {
	if len(experiences) == 0 {
		return nil
	}

	// 提取适应规则条件
	condition := RuleCondition{
		Expression: "adaptation_rate > threshold",
		Parameters: extractRuleParameters(experiences, "adaptation"),
		Threshold:  0.6,
	}

	// 提取规则动作
	action := RuleAction{
		Function: "optimize_strategy",
		Parameters: map[string]interface{}{
			"type":      expType,
			"fitness":   calculateAdaptationFitness(experiences),
			"direction": calculateAdaptationDirection(experiences),
		},
		ResultType: "strategy_optimization",
	}

	return &RulePattern{
		Type:       "adaptation",
		Target:     expType,
		Condition:  condition,
		Action:     action,
		Confidence: calculateAdaptationConfidence(experiences),
		Frequency:  calculateAdaptationFrequency(experiences, expType),
	}
}

// 辅助函数
func extractRuleParameters(experiences []LearningExperience, ruleType string) map[string]interface{} {
	params := make(map[string]interface{})

	// 基础统计参数
	successCount := 0
	failureCount := 0
	totalCount := 0

	// 类型特定的统计
	adaptationCount := 0
	effectiveCount := 0

	for _, exp := range experiences {
		totalCount++

		// 基础统计
		if exp.Result.Status == "success" {
			successCount++
			// 收集成功经验的特定指标
			if metrics, ok := exp.Result.Metrics[ruleType]; ok {
				params[ruleType+"_metrics"] = metrics
			}
		} else {
			failureCount++
		}

		// 根据规则类型收集特定参数
		switch ruleType {
		case "success":
			// 收集成功相关指标
			if effectiveness, ok := exp.Result.Metrics["effectiveness"]; ok {
				params["effectiveness"] = effectiveness
			}

		case "failure":
			// 收集失败原因统计
			if reason, ok := exp.Context["failure_reason"].(string); ok {
				params["failure_patterns"] = append(
					params["failure_patterns"].([]string), reason)
			}

		case "adaptation":
			// 收集适应性指标
			if isAdaptiveAction(exp) {
				adaptationCount++
				if exp.Result.Status == "success" {
					effectiveCount++
				}
			}
		}
	}

	// 设置基础统计参数
	params["success_rate"] = float64(successCount) / float64(totalCount)
	params["failure_rate"] = float64(failureCount) / float64(totalCount)
	params["total_count"] = totalCount

	// 设置类型特定参数
	switch ruleType {
	case "adaptation":
		if adaptationCount > 0 {
			params["adaptation_rate"] = float64(adaptationCount) / float64(totalCount)
			params["adaptation_effectiveness"] = float64(effectiveCount) / float64(adaptationCount)
		}
	}

	return params
}

func calculateSuccessConfidence(experiences []LearningExperience) float64 {
	if len(experiences) == 0 {
		return 0
	}

	successCount := 0
	for _, exp := range experiences {
		if exp.Result.Status == "success" {
			successCount++
		}
	}

	return float64(successCount) / float64(len(experiences))
}

func calculateFailureConfidence(experiences []LearningExperience) float64 {
	if len(experiences) == 0 {
		return 0
	}

	failureCount := 0
	for _, exp := range experiences {
		if exp.Result.Status != "success" {
			failureCount++
		}
	}

	return float64(failureCount) / float64(len(experiences))
}

func calculateAdaptationFitness(experiences []LearningExperience) float64 {
	if len(experiences) == 0 {
		return 0
	}

	totalFitness := 0.0
	for _, exp := range experiences {
		if fitness, ok := exp.Result.Metrics["fitness"]; ok {
			totalFitness += fitness
		}
	}

	return totalFitness / float64(len(experiences))
}

func calculateAdaptationDirection(experiences []LearningExperience) float64 {
	if len(experiences) < 2 {
		return 0
	}

	// 计算适应趋势
	trends := make([]float64, 0)
	for i := 1; i < len(experiences); i++ {
		prev := experiences[i-1].Result.Metrics["fitness"]
		curr := experiences[i].Result.Metrics["fitness"]
		trends = append(trends, curr-prev)
	}

	// 计算平均趋势
	total := 0.0
	for _, trend := range trends {
		total += trend
	}

	return total / float64(len(trends))
}

// optimizeRules 优化规则
func (al *AdaptiveLearning) optimizeRules() error {
	// 获取现有规则
	rules := al.strategy.GetRules()

	for _, rule := range rules {
		// 评估规则效果
		effectiveness := evaluateRuleEffectiveness(rule, al.state.experiences)

		if effectiveness < 0.5 {
			// 尝试优化规则
			optimized := optimizeRule(rule, al.state.experiences)
			if optimized != nil {
				al.strategy.UpdateRule(optimized)
			}
		}
	}

	return nil
}

// 辅助函数
func extractSuccessParameters(patterns []ExperiencePattern) []ParameterPattern {
	params := make([]ParameterPattern, 0)

	for _, pattern := range patterns {
		if pattern.Success && pattern.Confidence > 0.7 {
			params = append(params, ParameterPattern{
				Type:       pattern.Type,
				Parameters: extractParameters(pattern),
				Weight:     pattern.Confidence,
			})
		}
	}

	return params
}

// extractParameters 从模式中提取参数
func extractParameters(pattern ExperiencePattern) map[string]interface{} {
	params := make(map[string]interface{})

	// 从条件中提取参数
	for _, cond := range pattern.Conditions {
		if cond.Type == "parameter" {
			params[cond.Key] = cond.Value
		}
	}

	// 从上下文中提取参数
	for k, v := range pattern.Context {
		if strings.HasPrefix(k, "param_") {
			params[strings.TrimPrefix(k, "param_")] = v
		}
	}

	// 从结果中提取参数调整
	for _, outcome := range pattern.Outcomes {
		if outcome.Type == "parameter_adjustment" {
			for k, v := range outcome.Metrics {
				params[k] = v
			}
		}
	}

	return params
}

// 权重计算
func calculateRuleWeight(pattern RulePattern) float64 {
	// 基础权重
	baseWeight := pattern.Confidence * pattern.Frequency

	// 根据规则类型调整权重
	switch pattern.Type {
	case "success":
		baseWeight *= 1.2 // 成功规则加权
	case "failure":
		baseWeight *= 0.8 // 失败规则减权
	case "adaptation":
		baseWeight *= 1.1 // 适应规则略微加权
	}

	return math.Max(0, math.Min(1, baseWeight))
}

func evaluateRuleEffectiveness(rule *StrategyRule, experiences []LearningExperience) float64 {
	successCount := 0
	totalCount := 0

	for _, exp := range experiences {
		if isRuleApplicable(rule, exp) {
			if exp.Result.Status == "success" {
				successCount++
			}
			totalCount++
		}
	}

	if totalCount == 0 {
		return 0
	}
	return float64(successCount) / float64(totalCount)
}

// isRuleApplicable 检查规则是否适用于经验
func isRuleApplicable(rule *StrategyRule, exp LearningExperience) bool {
	// 1. 检查目标类型匹配
	if rule.Target != exp.Type {
		return false
	}

	// 2. 检查条件表达式
	switch rule.Condition.Expression {
	case "success_rate > threshold":
		if rate, ok := exp.Context["success_rate"].(float64); ok {
			return rate > rule.Condition.Threshold
		}
	case "failure_rate > threshold":
		if rate, ok := exp.Context["failure_rate"].(float64); ok {
			return rate > rule.Condition.Threshold
		}
	case "adaptation_rate > threshold":
		if rate, ok := exp.Context["adaptation_rate"].(float64); ok {
			return rate > rule.Condition.Threshold
		}
	}

	// 3. 检查自定义参数
	for key, expected := range rule.Condition.Parameters {
		if actual, exists := exp.Context[key]; !exists || actual != expected {
			return false
		}
	}

	return true
}

func optimizeRule(rule *StrategyRule, experiences []LearningExperience) *StrategyRule {
	// 创建规则副本
	optimized := *rule

	// 基于经验优化条件阈值
	if threshold := findOptimalThreshold(rule, experiences); threshold > 0 {
		optimized.Condition.Threshold = threshold
	}

	// 优化动作参数
	if params := optimizeActionParameters(rule, experiences); len(params) > 0 {
		optimized.Action.Parameters = params
	}

	return &optimized
}

// findOptimalThreshold 找到最优阈值
func findOptimalThreshold(rule *StrategyRule, experiences []LearningExperience) float64 {
	if len(experiences) == 0 {
		return 0
	}

	// 收集统计数据
	values := make([]float64, 0)
	for _, exp := range experiences {
		// 基于规则条件类型获取相关值
		switch rule.Condition.Expression {
		case "success_rate > threshold":
			if rate, ok := exp.Context["success_rate"].(float64); ok {
				values = append(values, rate)
			}
		case "failure_rate > threshold":
			if rate, ok := exp.Context["failure_rate"].(float64); ok {
				values = append(values, rate)
			}
		case "adaptation_rate > threshold":
			if rate, ok := exp.Context["adaptation_rate"].(float64); ok {
				values = append(values, rate)
			}
		}
	}

	if len(values) == 0 {
		return 0
	}

	// 寻找最优阈值
	sort.Float64s(values)
	medianIndex := len(values) / 2
	mean := calculateMean(values)
	median := values[medianIndex]

	// 使用加权平均作为最优阈值
	return mean*0.6 + median*0.4
}

// optimizeActionParameters 优化动作参数
func optimizeActionParameters(rule *StrategyRule, experiences []LearningExperience) map[string]interface{} {
	if len(experiences) == 0 {
		return nil
	}

	// 收集成功动作的参数统计
	paramStats := make(map[string][]float64)
	for _, exp := range experiences {
		if exp.Result.Status == "success" && exp.Action.Type == rule.Action.Function {
			for k, v := range exp.Action.Parameters {
				if fv, ok := v.(float64); ok {
					paramStats[k] = append(paramStats[k], fv)
				}
			}
		}
	}

	// 计算最优参数
	optimizedParams := make(map[string]interface{})
	for param, values := range paramStats {
		if len(values) > 0 {
			// 使用加权平均值作为最优参数
			sort.Float64s(values)
			median := values[len(values)/2]
			mean := calculateMean(values)
			optimizedParams[param] = mean*0.7 + median*0.3
		}
	}

	return optimizedParams
}

// 辅助函数
func calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// 辅助函数

func (al *AdaptiveLearning) addExperience(experience LearningExperience) {
	al.state.experiences = append(al.state.experiences, experience)

	// 限制经验数量
	if len(al.state.experiences) > al.config.memoryCapacity {
		al.state.experiences = al.state.experiences[1:]
	}
}

func (al *AdaptiveLearning) integrateKnowledge(knowledge *KnowledgeUnit) {
	// 检查知识是否已存在
	if existing, exists := al.state.knowledge[knowledge.ID]; exists {
		// 合并知识
		al.mergeKnowledge(existing, knowledge)
	} else {
		// 添加新知识
		al.state.knowledge[knowledge.ID] = knowledge
	}
}

// mergeKnowledge 合并知识
func (al *AdaptiveLearning) mergeKnowledge(existing, new *KnowledgeUnit) {
	// 1. 合并元数据
	existing.Metadata.Confidence = (existing.Metadata.Confidence*float64(existing.Metadata.Usage) +
		new.Metadata.Confidence) / float64(existing.Metadata.Usage+1)
	existing.Metadata.Usage++
	existing.Metadata.LastAccess = time.Now()

	// 合并标签
	existing.Metadata.Tags = mergeUniqueTags(existing.Metadata.Tags, new.Metadata.Tags)

	// 2. 更新关联
	existing.Connections = mergeKnowledgeConnections(existing.Connections, new.Connections)

	// 3. 合并内容
	if pattern, ok := new.Content.(ExperiencePattern); ok {
		if existingPattern, ok := existing.Content.(ExperiencePattern); ok {
			existing.Content = mergeExperiencePatterns(existingPattern, pattern)
		}
	}

	// 4. 更新验证函数
	if new.ValidationFn != nil {
		existing.ValidationFn = new.ValidationFn
	}
}

// 辅助函数
func mergeUniqueTags(tags1, tags2 []string) []string {
	tagMap := make(map[string]bool)
	for _, tag := range tags1 {
		tagMap[tag] = true
	}
	for _, tag := range tags2 {
		tagMap[tag] = true
	}

	merged := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		merged = append(merged, tag)
	}
	return merged
}

func mergeKnowledgeConnections(conns1, conns2 []KnowledgeLink) []KnowledgeLink {
	connMap := make(map[string]KnowledgeLink)

	// 处理第一组连接
	for _, conn := range conns1 {
		connMap[conn.TargetID] = conn
	}

	// 合并第二组连接
	for _, conn := range conns2 {
		if existing, exists := connMap[conn.TargetID]; exists {
			// 更新现有连接的强度
			connMap[conn.TargetID] = KnowledgeLink{
				TargetID: conn.TargetID,
				Type:     conn.Type,
				Strength: (existing.Strength + conn.Strength) / 2,
				Context:  mergeContexts(existing.Context, conn.Context),
			}
		} else {
			connMap[conn.TargetID] = conn
		}
	}

	// 转换回切片
	merged := make([]KnowledgeLink, 0, len(connMap))
	for _, conn := range connMap {
		merged = append(merged, conn)
	}
	return merged
}

func mergeExperiencePatterns(p1, p2 ExperiencePattern) ExperiencePattern {
	return ExperiencePattern{
		Type:       p1.Type,
		Confidence: (p1.Confidence + p2.Confidence) / 2,
		Frequency:  (p1.Frequency + p2.Frequency) / 2,
		Context:    mergeContexts(p1.Context, p2.Context),
		Conditions: mergeSimilarConditions(append(p1.Conditions, p2.Conditions...)),
		Outcomes:   mergeRelatedOutcomes(append(p1.Outcomes, p2.Outcomes...)),
		Success:    p1.Success || p2.Success,
	}
}

func mergeContexts(ctx1, ctx2 map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	// 复制第一个上下文
	for k, v := range ctx1 {
		merged[k] = v
	}

	// 合并第二个上下文
	for k, v := range ctx2 {
		if existing, exists := merged[k]; exists {
			// 如果两个值都是数值类型,取平均值
			if f1, ok1 := existing.(float64); ok1 {
				if f2, ok2 := v.(float64); ok2 {
					merged[k] = (f1 + f2) / 2
					continue
				}
			}
		}
		merged[k] = v
	}

	return merged
}

func generateKnowledgeID() string {
	return fmt.Sprintf("know_%d", time.Now().UnixNano())
}
