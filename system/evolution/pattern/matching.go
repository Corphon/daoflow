//system/evolution/pattern/matching.go

package pattern

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
	"github.com/Corphon/daoflow/system/meta/emergence"
	"github.com/Corphon/daoflow/system/meta/resonance"
	"github.com/Corphon/daoflow/system/types"
)

// EvolutionMatcher 演化匹配器
type EvolutionMatcher struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		matchThreshold float64 // 匹配阈值
		evolutionDepth int     // 演化深度
		adaptiveBias   float64 // 自适应偏差
		contextWeight  float64 // 上下文权重
	}

	// 匹配状态
	state struct {
		matches      map[string]*EvolutionMatch    // 当前匹配
		trajectories map[string]*EvolutionPath     // 演化轨迹
		context      *MatchingContext              // 匹配上下文
		patterns     map[string]*RecognizedPattern // 模式集合
		metrics      struct {                      // 指标
			activityLevel float64
			energyLevel   float64
			stability     float64
			changeRate    float64
		}
	}

	// 依赖项
	recognizer *PatternRecognizer
	matcher    *resonance.PatternMatcher
}

// EvolutionMatch 演化匹配
type EvolutionMatch struct {
	ID         string             // 匹配ID
	SourceID   string             // 源模式ID
	TargetID   string             // 目标模式ID
	Similarity float64            // 相似度
	Evolution  []EvolutionStep    // 演化步骤
	Context    map[string]float64 // 上下文因素
	StartTime  time.Time          // 开始时间
	LastUpdate time.Time          // 最后更新时间
}

// EvolutionPath 演化轨迹
type EvolutionPath struct {
	ID          string             // 轨迹ID
	Steps       []PathStep         // 轨迹步骤
	Properties  map[string]float64 // 轨迹属性
	Probability float64            // 概率
	Created     time.Time          // 创建时间
}

// EvolutionStep 演化步骤
type EvolutionStep struct {
	Timestamp time.Time
	Type      string             // 步骤类型
	Changes   map[string]float64 // 变化量
	Energy    float64            // 能量变化
	Stability float64            // 稳定性
}

// PathStep 轨迹步骤
type PathStep struct {
	Pattern     *RecognizedPattern // 相关模式
	Transition  string             // 转换类型
	Probability float64            // 转换概率
	Context     map[string]float64 // 步骤上下文
}

// MatchingContext 匹配上下文
type MatchingContext struct {
	Time        time.Time          // 当前时间
	Environment map[string]float64 // 环境因素
	History     []ContextState     // 历史状态
	Bias        map[string]float64 // 偏差项
}

// ContextState 上下文状态
type ContextState struct {
	Timestamp time.Time
	Factors   map[string]float64
	Influence float64
}

// ------------------------------------------------------------
// NewEvolutionMatcher 创建新的演化匹配器
func NewEvolutionMatcher(
	recognizer *PatternRecognizer,
	config *types.EvolutionConfig) (*EvolutionMatcher, error) {
	if recognizer == nil {
		return nil, fmt.Errorf("nil pattern recognizer")
	}
	if config == nil {
		return nil, fmt.Errorf("nil evolution config")
	}

	em := &EvolutionMatcher{
		recognizer: recognizer,
	}

	// 初始化配置
	em.config.matchThreshold = config.MatchThreshold
	em.config.evolutionDepth = config.EvolutionDepth
	em.config.adaptiveBias = config.AdaptiveBias
	em.config.contextWeight = config.ContextWeight

	// 初始化状态
	em.state.matches = make(map[string]*EvolutionMatch)
	em.state.trajectories = make(map[string]*EvolutionPath)
	em.state.context = &MatchingContext{
		Time:        time.Now(),
		Environment: make(map[string]float64),
		History:     make([]ContextState, 0),
		Bias:        make(map[string]float64),
	}

	return em, nil
}

// Match 执行演化匹配
func (em *EvolutionMatcher) Match() error {
	em.mu.Lock()
	defer em.mu.Unlock()

	// 更新上下文
	em.updateContext()

	// 获取当前模式
	patterns := em.recognizer.GetPatterns()

	// 执行匹配
	matches := em.matchPatterns(patterns)

	// 更新演化轨迹
	em.updateTrajectories(matches)

	// 预测演化方向
	em.predictEvolution()

	return nil
}

// matchPatterns 匹配模式
func (em *EvolutionMatcher) matchPatterns(
	patterns []*RecognizedPattern) []*EvolutionMatch {

	matches := make([]*EvolutionMatch, 0)

	// 对每对模式进行匹配
	for i, source := range patterns {
		for j := i + 1; j < len(patterns); j++ {
			target := patterns[j]

			// 计算演化相似度
			similarity := em.calculateEvolutionSimilarity(source, target)
			if similarity < em.config.matchThreshold {
				continue
			}

			// 创建匹配
			match := em.createMatch(source, target, similarity)
			matches = append(matches, match)
		}
	}

	return matches
}

// updateTrajectories 更新演化轨迹
func (em *EvolutionMatcher) updateTrajectories(matches []*EvolutionMatch) {
	for _, match := range matches {
		// 检查是否存在相关轨迹
		trajectoryID := em.findRelatedTrajectory(match)
		if trajectoryID == "" {
			// 创建新轨迹
			trajectory := em.createTrajectory(match)
			em.state.trajectories[trajectory.ID] = trajectory
		} else {
			// 更新现有轨迹
			em.updateExistingTrajectory(trajectoryID, match)
		}
	}

	// 移除过期轨迹
	em.cleanupTrajectories()
}

// findRelatedTrajectory 查找相关轨迹
func (em *EvolutionMatcher) findRelatedTrajectory(match *EvolutionMatch) string {
	for id, trajectory := range em.state.trajectories {
		if len(trajectory.Steps) > 0 {
			lastStep := trajectory.Steps[len(trajectory.Steps)-1]
			if lastStep.Pattern.ID == match.SourceID {
				return id
			}
		}
	}
	return ""
}

// createTrajectory 创建新轨迹
func (em *EvolutionMatcher) createTrajectory(match *EvolutionMatch) *EvolutionPath {
	trajectory := &EvolutionPath{
		ID:          core.GenerateID(),
		Steps:       make([]PathStep, 0),
		Properties:  make(map[string]float64),
		Probability: 1.0,
		Created:     time.Now(),
	}

	// 添加初始步骤
	step := PathStep{
		Pattern:     em.recognizer.GetPattern(match.SourceID),
		Transition:  "initial",
		Probability: 1.0,
		Context:     match.Context,
	}
	trajectory.Steps = append(trajectory.Steps, step)

	return trajectory
}

// updateExistingTrajectory 更新现有轨迹
func (em *EvolutionMatcher) updateExistingTrajectory(trajectoryID string, match *EvolutionMatch) {
	trajectory := em.state.trajectories[trajectoryID]

	// 添加新步骤
	step := PathStep{
		Pattern:     em.recognizer.GetPattern(match.TargetID),
		Transition:  "evolution",
		Probability: match.Similarity,
		Context:     match.Context,
	}
	trajectory.Steps = append(trajectory.Steps, step)

	// 更新轨迹属性
	trajectory.Properties["length"] = float64(len(trajectory.Steps))
	trajectory.Properties["avgSimilarity"] = em.calculateAverageSimilarity(trajectory)
}

// calculateAverageSimilarity 计算轨迹平均相似度
func (em *EvolutionMatcher) calculateAverageSimilarity(trajectory *EvolutionPath) float64 {
	if len(trajectory.Steps) == 0 {
		return 0
	}

	totalSimilarity := 0.0
	for _, step := range trajectory.Steps {
		totalSimilarity += step.Probability
	}

	return totalSimilarity / float64(len(trajectory.Steps))
}

// cleanupTrajectories 清理过期轨迹
func (em *EvolutionMatcher) cleanupTrajectories() {
	const maxAge = 24 * time.Hour
	now := time.Now()

	for id, trajectory := range em.state.trajectories {
		age := now.Sub(trajectory.Created)
		if age > maxAge {
			delete(em.state.trajectories, id)
		}
	}
}

// predictEvolution 预测演化方向
func (em *EvolutionMatcher) predictEvolution() {
	for _, trajectory := range em.state.trajectories {
		// 分析轨迹模式
		pattern := em.analyzeTrajectoryPattern(trajectory)

		// 预测下一步
		nextStep := em.predictNextStep(trajectory, pattern)
		if nextStep != nil {
			trajectory.Steps = append(trajectory.Steps, *nextStep)
		}

		// 更新概率
		trajectory.Probability = em.calculateTrajectoryProbability(trajectory)
	}
}

// analyzeTrajectoryPattern 分析轨迹模式
func (em *EvolutionMatcher) analyzeTrajectoryPattern(trajectory *EvolutionPath) *TrajectoryPattern {
	if len(trajectory.Steps) < 2 {
		return nil
	}

	lastPattern := convertToEmergentPattern(trajectory.Steps[len(trajectory.Steps)-1].Pattern)

	// 提取演化特征
	features := map[string]float64{
		"direction": calculateEvolutionDirectionality(lastPattern),
		"rate":      calculateEvolutionRate(lastPattern),
		"stability": calculateStabilityFeatures(lastPattern),
	}

	return &TrajectoryPattern{
		Direction: features["direction"],
		Rate:      features["rate"],
		Stability: features["stability"],
	}
}

// predictNextStep 预测下一步
func (em *EvolutionMatcher) predictNextStep(trajectory *EvolutionPath, pattern *TrajectoryPattern) *PathStep {
	if pattern == nil || len(trajectory.Steps) == 0 {
		return nil
	}

	lastStep := trajectory.Steps[len(trajectory.Steps)-1]

	// 基于当前模式预测
	predictedPattern := &RecognizedPattern{
		ID:        core.GenerateID(),
		Type:      lastStep.Pattern.Type,
		Evolution: lastStep.Pattern.Evolution,
	}

	// 按演化趋势调整属性
	for key, value := range lastStep.Pattern.Properties {
		predictedPattern.Properties[key] = value * (1 + pattern.Rate*pattern.Direction)
	}

	return &PathStep{
		Pattern:     predictedPattern,
		Transition:  "predicted",
		Probability: pattern.Stability,
		Context:     lastStep.Context,
	}
}

// calculateTrajectoryProbability 计算轨迹概率
func (em *EvolutionMatcher) calculateTrajectoryProbability(trajectory *EvolutionPath) float64 {
	if len(trajectory.Steps) == 0 {
		return 0
	}

	// 基于步骤概率计算整体概率
	probability := 1.0
	decayFactor := 0.95 // 每步衰减因子

	for i, step := range trajectory.Steps {
		stepWeight := math.Pow(decayFactor, float64(len(trajectory.Steps)-i-1))
		probability *= step.Probability * stepWeight
	}

	// 考虑时间衰减
	age := time.Since(trajectory.Created).Hours()
	timeDecay := math.Exp(-age / 24.0) // 24小时衰减周期

	return probability * timeDecay
}

// TrajectoryPattern 轨迹模式结构
type TrajectoryPattern struct {
	Direction float64 // 演化方向
	Rate      float64 // 演化速率
	Stability float64 // 稳定性
}

// calculateEvolutionSimilarity 计算演化相似度
func (em *EvolutionMatcher) calculateEvolutionSimilarity(
	source, target *RecognizedPattern) float64 {

	// 基础相似度
	baseSimilarity := calculatePatternSimilarity(source, target)

	// 演化特征相似度
	evolutionSimilarity := em.compareEvolutionFeatures(source, target)

	// 上下文相似度
	contextSimilarity := em.calculateContextSimilarity(source, target)

	// 组合相似度
	similarity := (baseSimilarity +
		evolutionSimilarity*(1-em.config.contextWeight) +
		contextSimilarity*em.config.contextWeight) / 3.0

	return similarity
}

// calculatePatternSimilarity 计算模式基础相似度
func calculatePatternSimilarity(source, target *RecognizedPattern) float64 {
	if source == nil || target == nil {
		return 0
	}

	// 1. 类型相似度
	typeSimilarity := 0.0
	if source.Type == target.Type {
		typeSimilarity = 1.0
	}

	// 2. 属性相似度
	propertySimilarity := calculatePropertySimilarity(source.Properties, target.Properties)

	// 3. 结构相似度
	structureSimilarity := calculateStructuralSimilarity(source.Pattern, target.Pattern)

	return (typeSimilarity*0.3 + propertySimilarity*0.3 + structureSimilarity*0.4)
}

// calculatePropertySimilarity 计算属性相似度
func calculatePropertySimilarity(props1, props2 map[string]float64) float64 {
	if len(props1) == 0 || len(props2) == 0 {
		return 0
	}

	similarity := 0.0
	count := 0.0

	// 遍历所有共同属性
	for key, val1 := range props1 {
		if val2, exists := props2[key]; exists {
			similarity += 1.0 - math.Abs(val1-val2)
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return similarity / count
}

// calculateStructuralSimilarity 计算结构相似度
func calculateStructuralSimilarity(pattern1, pattern2 *emergence.EmergentPattern) float64 {
	if pattern1 == nil || pattern2 == nil {
		return 0
	}

	// 1. 组件复杂度相似度
	comp1 := pattern1.GetStructureComplexity()
	comp2 := pattern2.GetStructureComplexity()
	complexitySimilarity := 1.0 - math.Abs(comp1-comp2)

	// 2. 相干性相似度
	coh1 := pattern1.GetStructureCoherence()
	coh2 := pattern2.GetStructureCoherence()
	coherenceSimilarity := 1.0 - math.Abs(coh1-coh2)

	// 3. 对称性相似度
	sym1 := pattern1.GetStructureSymmetry()
	sym2 := pattern2.GetStructureSymmetry()
	symmetrySimilarity := 1.0 - math.Abs(sym1-sym2)

	return (complexitySimilarity*0.4 + coherenceSimilarity*0.3 + symmetrySimilarity*0.3)
}

// compareEvolutionFeatures 比较演化特征
func (em *EvolutionMatcher) compareEvolutionFeatures(source, target *RecognizedPattern) float64 {
	sourceFeatures := extractEvolutionFeatures(source)
	targetFeatures := extractEvolutionFeatures(target)

	similarity := 0.0
	count := 0.0

	// 比较关键演化特征
	for key, sourceVal := range sourceFeatures {
		if targetVal, exists := targetFeatures[key]; exists {
			similarity += 1.0 - math.Abs(sourceVal-targetVal)
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return similarity / count
}

// calculateContextSimilarity 计算上下文相似度
func (em *EvolutionMatcher) calculateContextSimilarity(source, target *RecognizedPattern) float64 {
	// 1. 时间相关性
	timeDiff := target.LastSeen.Sub(source.LastSeen).Hours()
	timeCorrelation := math.Exp(-timeDiff / 24.0) // 24小时衰减

	// 2. 环境因素相似度
	environmentSimilarity := calculateEnvironmentSimilarity(
		em.state.context.Environment,
		source.Context,
		target.Context)

	// 3. 状态相关性
	stateSimilarity := calculateStateSimilarity(source, target)

	return (timeCorrelation*0.3 + environmentSimilarity*0.4 + stateSimilarity*0.3)
}

// updateContext 更新上下文
func (em *EvolutionMatcher) updateContext() {
	currentTime := time.Now()

	// 更新环境因素
	em.updateEnvironmentFactors()

	// 记录当前状态
	state := ContextState{
		Timestamp: currentTime,
		Factors:   make(map[string]float64),
		Influence: calculateContextInfluence(em.state.context.Environment),
	}

	// 复制当前环境因素
	for k, v := range em.state.context.Environment {
		state.Factors[k] = v
	}

	// 更新历史
	em.state.context.History = append(em.state.context.History, state)
	if len(em.state.context.History) > maxContextHistory {
		em.state.context.History = em.state.context.History[1:]
	}

	// 更新时间
	em.state.context.Time = currentTime
}

// updateEnvironmentFactors 更新环境因素
func (em *EvolutionMatcher) updateEnvironmentFactors() {
	// 基础环境因素
	em.state.context.Environment["time_of_day"] = normalizeTimeOfDay(time.Now())
	em.state.context.Environment["activity_level"] = calculateActivityLevel(em)
	em.state.context.Environment["energy_level"] = calculateSystemEnergy(em)
	em.state.context.Environment["stability"] = calculateSystemStability(em)

	// 动态环境因素
	if len(em.state.context.History) > 0 {
		lastState := em.state.context.History[len(em.state.context.History)-1]
		em.state.context.Environment["change_rate"] = calculateChangeRate(lastState, em.state.context.Environment)
	}
}

// calculateContextInfluence 计算上下文影响度
func calculateContextInfluence(env map[string]float64) float64 {
	weights := map[string]float64{
		"time_of_day":    0.1,
		"activity_level": 0.3,
		"energy_level":   0.3,
		"stability":      0.2,
		"change_rate":    0.1,
	}

	influence := 0.0
	totalWeight := 0.0

	for factor, weight := range weights {
		if value, exists := env[factor]; exists {
			influence += value * weight
			totalWeight += weight
		}
	}

	if totalWeight == 0 {
		return 0.5 // 默认中等影响
	}

	return influence / totalWeight
}

// 辅助函数

func (em *EvolutionMatcher) createMatch(
	source, target *RecognizedPattern,
	similarity float64) *EvolutionMatch {

	match := &EvolutionMatch{
		ID:         generateMatchID(),
		SourceID:   source.ID,
		TargetID:   target.ID,
		Similarity: similarity,
		Evolution:  make([]EvolutionStep, 0),
		Context:    make(map[string]float64),
		StartTime:  time.Now(),
		LastUpdate: time.Now(),
	}

	return match
}

func generateMatchID() string {
	return fmt.Sprintf("match_%d", time.Now().UnixNano())
}

const (
	maxContextHistory = 100
)
