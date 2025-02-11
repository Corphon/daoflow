//system/evolution/pattern/recognition.go

package pattern

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/Corphon/daoflow/system/common"
	"github.com/Corphon/daoflow/system/meta/emergence"
	"github.com/Corphon/daoflow/system/meta/resonance"
)

const (
	maxHistoryLength = 1000
)

// PatternRecognizer 模式识别器
type PatternRecognizer struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		minConfidence float64 // 最小置信度
		learningRate  float64 // 学习率
		memoryDepth   int     // 记忆深度
		adaptiveRate  bool    // 是否使用自适应学习率
	}

	// 识别状态
	state struct {
		patterns   map[string]*RecognizedPattern // 已识别模式
		memories   []PatternMemory               // 模式记忆
		statistics PatternStatistics             // 统计信息
	}

	mutationAnalyzer common.PatternAnalyzer        // 使用接口而不是具体类型
	detector         *emergence.PatternDetector    // 模式检测器
	matcher          *resonance.PatternMatcher     // 模式匹配器
	amplifier        *resonance.ResonanceAmplifier // 共振放大器
}

// PatternSignature 模式特征
type PatternSignature struct {
	Components []SignatureComponent   // 组成成分
	Structure  map[string]interface{} // 结构特征
	Dynamics   map[string]float64     // 动态特征
	Context    map[string]string      // 上下文信息
}

// ComponentConnection 组件连接
type ComponentConnection struct {
	Type       string             // 连接类型
	Target     string             // 目标组件ID
	Strength   float64            // 连接强度
	Properties map[string]float64 // 连接属性
}

// SignatureComponent 特征组件
type SignatureComponent struct {
	Type        string                // 组件类型
	Properties  map[string]float64    // 组件属性
	Weight      float64               // 权重
	Role        string                // 角色
	Connections []ComponentConnection // 组件连接
}

// PatternMemory 模式记忆
type PatternMemory struct {
	Timestamp    time.Time
	Pattern      *RecognizedPattern
	Context      map[string]interface{}
	Associations []string
}

// PatternStatistics 模式统计
type PatternStatistics struct {
	TotalPatterns  int
	ActivePatterns int
	Recognition    map[string]float64 // 识别率统计
	Accuracy       map[string]float64 // 准确率统计
	Evolution      []StatPoint        // 演化趋势
}

// StatPoint 统计点
type StatPoint struct {
	Timestamp time.Time
	Metrics   map[string]float64
}

// NewPatternRecognizer 创建新的模式识别器
func NewPatternRecognizer(
	detector *emergence.PatternDetector,
	matcher *resonance.PatternMatcher,
	amplifier *resonance.ResonanceAmplifier) *PatternRecognizer {

	pr := &PatternRecognizer{
		detector:  detector,
		matcher:   matcher,
		amplifier: amplifier,
	}

	// 初始化配置
	pr.config.minConfidence = 0.75
	pr.config.learningRate = 0.1
	pr.config.memoryDepth = 100
	pr.config.adaptiveRate = true

	// 初始化状态
	pr.state.patterns = make(map[string]*RecognizedPattern)
	pr.state.memories = make([]PatternMemory, 0)
	pr.state.statistics = PatternStatistics{
		Recognition: make(map[string]float64),
		Accuracy:    make(map[string]float64),
		Evolution:   make([]StatPoint, 0),
	}

	return pr
}

// Recognize 执行模式识别
func (pr *PatternRecognizer) Recognize() error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	// 获取当前模式
	patterns, err := pr.detector.Detect()
	if err != nil {
		return err
	}

	// 识别新模式
	newPatterns := pr.recognizeNewPatterns(patterns)

	// 更新现有模式
	pr.updateExistingPatterns(patterns)

	// 构建模式记忆
	pr.buildPatternMemory(newPatterns)

	// 更新统计信息
	pr.updateStatistics()

	return nil
}

// calculateActivityLevel 计算活跃度 (PatternRecognizer专用版本)
func calculateActivityLevelForRecognizer(patterns map[string]*RecognizedPattern) float64 {
	if len(patterns) == 0 {
		return 0
	}

	activeCount := 0
	for _, pattern := range patterns {
		if pattern.Active {
			activeCount++
		}
	}

	return float64(activeCount) / float64(len(patterns))
}

// updateStatistics 更新统计信息
func (pr *PatternRecognizer) updateStatistics() {
	// 创建新的统计点
	point := StatPoint{
		Timestamp: time.Now(),
		Metrics: map[string]float64{
			"recognition_rate": calculateRecognitionRate(pr.state.patterns),
			"accuracy":         calculateAccuracy(pr.state.patterns),
			"activity":         calculateActivityLevelForRecognizer(pr.state.patterns), // 修改这里
			"stability":        calculateAverageStability(pr.state.patterns),
		},
	}

	// 更新历史趋势
	pr.state.statistics.Evolution = append(pr.state.statistics.Evolution, point)
	if len(pr.state.statistics.Evolution) > maxHistoryLength {
		pr.state.statistics.Evolution = pr.state.statistics.Evolution[1:]
	}

	// 更新识别率统计
	pr.state.statistics.Recognition["total"] = float64(len(pr.state.patterns))
	pr.state.statistics.Recognition["active"] = calculateActiveCount(pr.state.patterns)

	// 更新准确率统计
	pr.state.statistics.Accuracy["confidence"] = point.Metrics["accuracy"]
	pr.state.statistics.Accuracy["stability"] = point.Metrics["stability"]
}

// calculateActivityLevel 计算活跃度(为RecognizedPattern专门实现的版本)
func calculateActivityLevel(em *EvolutionMatcher) float64 {
	if len(em.state.patterns) == 0 {
		return 0
	}

	activeCount := 0
	for _, pattern := range em.state.patterns {
		if pattern.Active {
			activeCount++
		}
	}

	return float64(activeCount) / float64(len(em.state.patterns))
}

// 辅助计算函数
// calculateRecognitionRate 计算识别率
func calculateRecognitionRate(patterns map[string]*RecognizedPattern) float64 {
	if len(patterns) == 0 {
		return 0
	}

	recognizedCount := 0
	for _, pattern := range patterns {
		if pattern.Confidence > 0.5 { // 识别阈值
			recognizedCount++
		}
	}

	return float64(recognizedCount) / float64(len(patterns))
}

func calculateAccuracy(patterns map[string]*RecognizedPattern) float64 {
	if len(patterns) == 0 {
		return 0
	}
	highConfidenceCount := 0
	for _, p := range patterns {
		if p.Confidence >= 0.75 {
			highConfidenceCount++
		}
	}
	return float64(highConfidenceCount) / float64(len(patterns))
}

func calculateAverageStability(patterns map[string]*RecognizedPattern) float64 {
	if len(patterns) == 0 {
		return 0
	}
	totalStability := 0.0
	for _, p := range patterns {
		totalStability += p.Stability
	}
	return totalStability / float64(len(patterns))
}

func calculateActiveCount(patterns map[string]*RecognizedPattern) float64 {
	activeCount := 0
	for _, p := range patterns {
		if p.Active {
			activeCount++
		}
	}
	return float64(activeCount)
}

// recognizeNewPatterns 识别新模式
func (pr *PatternRecognizer) recognizeNewPatterns(
	patterns []emergence.EmergentPattern) []*RecognizedPattern {

	newPatterns := make([]*RecognizedPattern, 0)

	for _, pattern := range patterns {
		// 检查是否是新模式
		if pr.isKnownPattern(pattern) {
			continue
		}

		// 提取模式特征
		signature := pr.extractSignature(pattern)

		// 评估模式
		confidence := pr.evaluatePattern(pattern, signature)
		if confidence < pr.config.minConfidence {
			continue
		}

		// 创建新的识别模式
		recognized := &RecognizedPattern{
			ID:          generatePatternID(),
			Type:        determinePatternType(pattern),
			Signature:   signature,
			Confidence:  confidence,
			Stability:   calculateInitialStability(pattern),
			FirstSeen:   time.Now(),
			LastSeen:    time.Now(),
			Occurrences: 1,
			Evolution:   make([]PatternState, 0),
		}

		// 添加到已识别模式
		pr.state.patterns[recognized.ID] = recognized
		newPatterns = append(newPatterns, recognized)
	}

	return newPatterns
}

// updateExistingPatterns 更新现有模式
func (pr *PatternRecognizer) updateExistingPatterns(
	patterns []emergence.EmergentPattern) {

	for id, recognized := range pr.state.patterns {
		// 查找匹配的当前模式
		matched := false
		for _, pattern := range patterns {
			if pr.isPatternMatch(recognized, pattern) {
				// 更新模式状态
				pr.updatePatternState(recognized, pattern)
				matched = true
				break
			}
		}

		// 处理未匹配的模式
		if !matched {
			// 检查是否应该保留模式
			if pr.shouldRetainPattern(recognized) {
				// 降低置信度
				recognized.Confidence *= (1 - pr.config.learningRate)
			} else {
				// 移除模式
				delete(pr.state.patterns, id)
			}
		}
	}
}

// isPatternMatch 检查模式是否匹配
func (pr *PatternRecognizer) isPatternMatch(recognized *RecognizedPattern, pattern emergence.EmergentPattern) bool {
	// 1. 类型匹配
	if recognized.Type != pattern.Type {
		return false
	}

	// 2. 特征相似度
	signature := pr.extractSignature(pattern)
	similarity := calculateSignatureSimilarity(recognized.Signature, signature)

	// 3. 时间关联性
	timeDiff := time.Since(recognized.LastSeen)
	timeCorrelation := math.Exp(-timeDiff.Hours() / 24.0) // 24小时衰减

	return similarity*timeCorrelation >= pr.config.minConfidence
}

// updatePatternState 更新模式状态
func (pr *PatternRecognizer) updatePatternState(recognized *RecognizedPattern, pattern emergence.EmergentPattern) {
	// 更新基本信息
	recognized.LastSeen = time.Now()
	recognized.Occurrences++
	recognized.Active = true

	// 更新模式状态
	state := PatternState{
		Pattern:    &pattern,
		Active:     true,
		Duration:   time.Since(recognized.FirstSeen),
		LastUpdate: time.Now(),
		Properties: pattern.Properties,
	}
	recognized.Evolution = append(recognized.Evolution, state)

	// 更新置信度
	newConfidence := pr.evaluatePattern(pattern, recognized.Signature)
	recognized.Confidence = recognized.Confidence*(1-pr.config.learningRate) +
		newConfidence*pr.config.learningRate

	// 更新稳定性
	recognized.Stability = calculatePatternStability(recognized)
}

// shouldRetainPattern 判断是否应该保留模式
func (pr *PatternRecognizer) shouldRetainPattern(pattern *RecognizedPattern) bool {
	// 1. 检查置信度
	if pattern.Confidence < pr.config.minConfidence*0.5 {
		return false
	}

	// 2. 检查活跃度
	if !pattern.Active {
		inactiveDuration := time.Since(pattern.LastSeen)
		if inactiveDuration > 24*time.Hour {
			return false
		}
	}

	// 3. 检查历史稳定性
	if pattern.Stability < 0.3 {
		return false
	}

	return true
}

// buildPatternMemory 构建模式记忆
func (pr *PatternRecognizer) buildPatternMemory(newPatterns []*RecognizedPattern) {
	memory := PatternMemory{
		Timestamp:    time.Now(),
		Pattern:      nil,
		Context:      make(map[string]interface{}),
		Associations: make([]string, 0),
	}

	// 记录新模式
	for _, pattern := range newPatterns {
		memory.Pattern = pattern
		memory.Context = pr.extractContext(pattern)
		memory.Associations = pr.findAssociations(pattern)

		pr.state.memories = append(pr.state.memories, memory)
	}

	// 限制记忆深度
	if len(pr.state.memories) > pr.config.memoryDepth {
		pr.state.memories = pr.state.memories[1:]
	}
}

// extractContext 提取模式上下文
func (pr *PatternRecognizer) extractContext(pattern *RecognizedPattern) map[string]interface{} {
	context := make(map[string]interface{})

	// 1. 时间上下文
	context["timestamp"] = pattern.LastSeen
	context["duration"] = time.Since(pattern.FirstSeen)

	// 2. 环境上下文
	environment := make(map[string]float64)
	environment["stability"] = pattern.Stability
	environment["confidence"] = pattern.Confidence
	environment["activity"] = boolToFloat64(pattern.Active)
	context["environment"] = environment

	// 3. 演化上下文
	if len(pattern.Evolution) > 0 {
		lastState := pattern.Evolution[len(pattern.Evolution)-1]
		context["current_state"] = lastState
		context["evolution_stage"] = len(pattern.Evolution)
	}

	return context
}

// boolToFloat64 将bool转换为float64
func boolToFloat64(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

// findAssociations 查找相关联的模式
func (pr *PatternRecognizer) findAssociations(pattern *RecognizedPattern) []string {
	associations := make([]string, 0)

	// 遍历所有已识别的模式
	for id, other := range pr.state.patterns {
		if id == pattern.ID {
			continue
		}

		// 1. 类型关联
		if other.Type == pattern.Type {
			associations = append(associations, id)
			continue
		}

		// 2. 特征相似度关联
		similarity := calculateSignatureSimilarity(pattern.Signature, other.Signature)
		if similarity > pr.config.minConfidence {
			associations = append(associations, id)
			continue
		}

		// 3. 演化关联
		if evolutionaryRelated := checkEvolutionaryRelation(pattern, other); evolutionaryRelated {
			associations = append(associations, id)
			continue
		}
	}

	return associations
}

// 辅助函数: 检查演化关联
func checkEvolutionaryRelation(p1, p2 *RecognizedPattern) bool {
	if len(p1.Evolution) == 0 || len(p2.Evolution) == 0 {
		return false
	}

	// 检查时间重叠
	timeOverlap := p1.LastSeen.Sub(p2.FirstSeen) > 0 &&
		p2.LastSeen.Sub(p1.FirstSeen) > 0

	// 检查状态转换
	stateTransition := false
	for _, state1 := range p1.Evolution {
		for _, state2 := range p2.Evolution {
			if calculateStateDifference(state1, state2) < 0.3 {
				stateTransition = true
				break
			}
		}
	}

	return timeOverlap && stateTransition
}

// 辅助函数

func (pr *PatternRecognizer) isKnownPattern(
	pattern emergence.EmergentPattern) bool {

	for _, recognized := range pr.state.patterns {
		if pr.isPatternMatch(recognized, pattern) {
			return true
		}
	}
	return false
}

func (pr *PatternRecognizer) extractSignature(
	pattern emergence.EmergentPattern) PatternSignature {

	signature := PatternSignature{
		Components: make([]SignatureComponent, 0),
		Structure:  make(map[string]interface{}),
		Dynamics:   make(map[string]float64),
		Context:    make(map[string]string),
	}

	// 提取组件特征
	for _, comp := range pattern.Components {
		component := SignatureComponent{
			Type:       comp.Type,
			Properties: make(map[string]float64),
			Weight:     comp.Weight,
			Role:       comp.Role,
		}

		// 复制属性
		for k, v := range comp.Properties {
			component.Properties[k] = v
		}

		signature.Components = append(signature.Components, component)
	}

	// 提取结构特征
	signature.Structure = extractStructuralFeatures(pattern)

	// 提取动态特征
	signature.Dynamics = extractDynamicFeatures(pattern)

	return signature
}

func (pr *PatternRecognizer) evaluatePattern(
	pattern emergence.EmergentPattern,
	signature PatternSignature) float64 {

	// 基础置信度
	baseConfidence := pattern.Strength

	// 结构完整性评分
	structureScore := evaluateStructure(signature.Structure)

	// 动态稳定性评分
	stabilityScore := evaluateStability(signature.Dynamics)

	// 组合评分
	confidence := (baseConfidence + structureScore + stabilityScore) / 3.0

	return confidence
}

// evaluateStructure 评估结构完整性
func evaluateStructure(structure map[string]interface{}) float64 {
	// 1. 拓扑完整性
	topologyScore := 0.0
	if topology, ok := structure["topology"].(map[string]float64); ok {
		topologyScore = evaluateTopology(topology)
	}

	// 2. 连接完整性
	connectivityScore := 0.0
	if connectivity, ok := structure["connectivity"].(map[string]float64); ok {
		connectivityScore = evaluateConnectivity(connectivity)
	}

	// 3. 层级完整性
	hierarchyScore := 0.0
	if hierarchy, ok := structure["hierarchy"].(map[string]float64); ok {
		hierarchyScore = evaluateHierarchy(hierarchy)
	}

	return (topologyScore*0.4 + connectivityScore*0.3 + hierarchyScore*0.3)
}

// evaluateStability 评估动态稳定性
func evaluateStability(dynamics map[string]float64) float64 {
	// 1. 能量稳定性
	energyStability := 0.0
	if energy, exists := dynamics["energy"]; exists {
		energyStability = 1.0 - math.Min(1.0, math.Abs(energy-0.5)*2)
	}

	// 2. 演化稳定性
	evolutionStability := 0.0
	if rate, exists := dynamics["rate"]; exists {
		evolutionStability = 1.0 - math.Min(1.0, rate)
	}

	// 3. 相位稳定性
	phaseStability := 0.0
	if phase, exists := dynamics["phase"]; exists {
		phaseStability = math.Cos(phase) // 相位越接近0/π越稳定
	}

	return (energyStability*0.4 + evolutionStability*0.3 + phaseStability*0.3)
}

// 辅助函数
func evaluateTopology(topology map[string]float64) float64 {
	completeness := 0.0
	count := 0.0

	for key, value := range topology {
		switch key {
		case "connectivity":
			completeness += value
			count++
		case "symmetry":
			completeness += value
			count++
		case "complexity":
			completeness += 1.0 - value // 复杂度越低完整性越高
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return completeness / count
}

func evaluateConnectivity(connectivity map[string]float64) float64 {
	if density, exists := connectivity["density"]; exists {
		return density // 连接密度直接作为完整性度量
	}
	return 0
}

func evaluateHierarchy(hierarchy map[string]float64) float64 {
	if depth, exists := hierarchy["depth"]; exists {
		return 1.0 - math.Exp(-depth) // 层级深度越大完整性越高
	}
	return 0
}
func generatePatternID() string {
	return fmt.Sprintf("pat_%d", time.Now().UnixNano())
}

// GetPatterns 获取已识别的模式
func (pr *PatternRecognizer) GetPatterns() ([]*RecognizedPattern, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	patterns := make([]*RecognizedPattern, 0, len(pr.state.patterns))

	// 只返回活跃且置信度足够的模式
	for _, pattern := range pr.state.patterns {
		if pattern.Active && pattern.Confidence >= pr.config.minConfidence {
			patterns = append(patterns, pattern)
		}
	}

	return patterns, nil
}

// GetPattern 获取指定ID的模式
func (pr *PatternRecognizer) GetPattern(id string) *RecognizedPattern {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	if pattern, exists := pr.state.patterns[id]; exists {
		if pattern.Active && pattern.Confidence >= pr.config.minConfidence {
			return pattern
		}
	}
	return nil
}

// GetActivationLevel 获取模式激活水平
func (rp *RecognizedPattern) GetActivationLevel() float64 {
	if !rp.Active {
		return 0
	}

	// 基础激活度(置信度)
	activation := rp.Confidence

	// 使用频率影响
	usageScore := math.Min(1.0, float64(rp.Occurrences)/100.0)

	// 时间衰减
	age := time.Since(rp.LastSeen).Hours()
	timeDecay := math.Exp(-age / 24.0) // 24小时衰减周期

	// 组合计算
	return (activation*0.5 + usageScore*0.3) * timeDecay
}
