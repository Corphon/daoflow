//system/meta/emergence/generator.go

package emergence

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/meta/field"
)

const (
	maxHistoryLength    = 1000
	maxEvolutionHistory = 100
)

// PropertyGenerator 属性生成器
type PropertyGenerator struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		evolutionRate float64 // 演化速率
		complexity    float64 // 复杂度阈值
		stability     float64 // 稳定性要求
		minEnergy     float64 // 最小能量要求
		minStability  float64 // 最小稳定性要求
		minCoherence  float64 // 最小相干性要求
	}

	// 生成状态
	state struct {
		properties map[string]*EmergentProperty // 当前属性
		potential  []PotentialProperty          // 潜在属性
		history    []GenerationEvent            // 生成历史
	}

	// 依赖项
	detector *PatternDetector    // 模式检测器
	field    *field.UnifiedField // 统一场
}

// EmergentProperty 涌现属性
type EmergentProperty struct {
	ID         string              // 属性ID
	Name       string              // 属性名称
	Type       string              // 属性类型
	Value      float64             // 属性值
	Components []PropertyComponent // 组成成分
	Stability  float64             // 稳定性
	Evolution  []PropertyState     // 演化历史
	Created    time.Time           // 创建时间
	Updated    time.Time           // 更新时间
}

// PropertyComponent 属性组件
type PropertyComponent struct {
	PatternID string  // 关联模式ID
	Weight    float64 // 权重
	Role      string  // 作用角色
	Influence float64 // 影响度
}

// PropertyState 属性状态
type PropertyState struct {
	Timestamp time.Time
	Value     float64
	Stability float64
}

// PotentialProperty 潜在属性
type PotentialProperty struct {
	Type         string        // 属性类型
	Probability  float64       // 出现概率
	Requirements []string      // 所需条件
	TimeFrame    time.Duration // 预计时间框架
	Energy       float64       // 所需能量
}

// GenerationEvent 生成事件
type GenerationEvent struct {
	Timestamp  time.Time
	PropertyID string
	Type       string
	Old        *EmergentProperty
	New        *EmergentProperty
	Changes    map[string]float64
}

// ----------------------------------------------
// NewPropertyGenerator 创建新的属性生成器
func NewPropertyGenerator(detector *PatternDetector, field *field.UnifiedField) *PropertyGenerator {
	pg := &PropertyGenerator{
		detector: detector,
		field:    field,
	}

	// 初始化配置
	pg.config.evolutionRate = 0.1
	pg.config.complexity = 0.65
	pg.config.stability = 0.75
	pg.config.minEnergy = 0.3
	pg.config.minStability = 0.4
	pg.config.minCoherence = 0.5

	// 初始化状态
	pg.state.properties = make(map[string]*EmergentProperty)
	pg.state.potential = make([]PotentialProperty, 0)
	pg.state.history = make([]GenerationEvent, 0)

	return pg
}

// Generate 生成新属性
func (pg *PropertyGenerator) Generate() error {
	pg.mu.Lock()
	defer pg.mu.Unlock()

	// 获取当前模式
	patterns, err := pg.detector.Detect()
	if err != nil {
		return model.WrapError(err, model.ErrCodeOperation, "failed to detect patterns")
	}

	// 分析潜在属性
	potentials := pg.analyzePotentials(patterns)
	pg.state.potential = potentials

	// 生成新属性
	for _, potential := range potentials {
		if pg.shouldGenerate(potential) {
			if err := pg.generateProperty(potential, patterns); err != nil {
				return err
			}
		}
	}

	// 更新现有属性
	pg.updateProperties(patterns)

	return nil
}

// analyzePotentials 分析潜在属性
func (pg *PropertyGenerator) analyzePotentials(patterns []EmergentPattern) []PotentialProperty {
	potentials := make([]PotentialProperty, 0)

	// 分析模式组合
	combinations := pg.analyzePatternCombinations(patterns)

	// 评估每个组合的潜在属性
	for _, combo := range combinations {
		potential := pg.evaluatePotential(combo)
		if potential != nil {
			potentials = append(potentials, *potential)
		}
	}

	// 按概率排序
	sortPotentialsByProbability(potentials)

	return potentials
}

// analyzePatternCombinations 分析模式组合
func (pg *PropertyGenerator) analyzePatternCombinations(patterns []EmergentPattern) [][]EmergentPattern {
	combinations := make([][]EmergentPattern, 0)

	// 生成2-3个模式的组合
	for i := 0; i < len(patterns); i++ {
		for j := i + 1; j < len(patterns); j++ {
			// 双模式组合
			combo := []EmergentPattern{patterns[i], patterns[j]}
			combinations = append(combinations, combo)

			// 三模式组合
			for k := j + 1; k < len(patterns); k++ {
				tripleCombo := []EmergentPattern{patterns[i], patterns[j], patterns[k]}
				combinations = append(combinations, tripleCombo)
			}
		}
	}

	return combinations
}

// evaluatePotential 评估潜在属性
func (pg *PropertyGenerator) evaluatePotential(patterns []EmergentPattern) *PotentialProperty {
	// 计算组合特征
	coherence := calculatePatternCoherence(patterns)
	stability := calculatePatternStability(patterns)
	complexity := calculatePatternsComplexity(patterns)

	// 如果特征值太低则忽略
	if coherence < pg.config.complexity ||
		stability < pg.config.stability {
		return nil
	}

	// 创建潜在属性
	potential := &PotentialProperty{
		Type:         determinePotentialType(patterns),
		Probability:  (coherence + stability + complexity) / 3.0,
		Requirements: extractRequirements(patterns),
		TimeFrame:    estimateTimeFrame(stability),
		Energy:       calculateRequiredEnergy(patterns),
	}

	return potential
}

// sortPotentialsByProbability 按概率排序潜在属性
func sortPotentialsByProbability(potentials []PotentialProperty) {
	sort.Slice(potentials, func(i, j int) bool {
		return potentials[i].Probability > potentials[j].Probability
	})
}

// 辅助函数
func calculatePatternCoherence(patterns []EmergentPattern) float64 {
	if len(patterns) == 0 {
		return 0
	}

	coherence := 0.0
	for _, p := range patterns {
		if value, exists := p.Properties["coherence"]; exists {
			coherence += value
		}
	}

	return coherence / float64(len(patterns))
}

func calculatePatternStability(patterns []EmergentPattern) float64 {
	if len(patterns) == 0 {
		return 0
	}

	totalStability := 0.0
	for _, p := range patterns {
		// 结合强度和稳定性
		stability := p.Stability
		if value, exists := p.Properties["stability"]; exists {
			stability = (stability + value) / 2.0
		}
		totalStability += stability
	}

	return totalStability / float64(len(patterns))
}

// calculatePatternComplexity 计算单个模式复杂度
func calculatePatternComplexity(pattern *EmergentPattern) float64 {
	if pattern == nil {
		return 0
	}

	// 基础复杂度 - 组件数量
	baseComplexity := float64(len(pattern.Components)) / 10.0

	// 关系复杂度 - 组件间关联
	relationComplexity := 0.0
	for i := 0; i < len(pattern.Components)-1; i++ {
		for j := i + 1; j < len(pattern.Components); j++ {
			if pattern.Components[i].Type == pattern.Components[j].Type {
				relationComplexity += 0.1
			}
		}
	}

	// 属性复杂度
	propertyComplexity := float64(len(pattern.Properties)) / 10.0

	return (baseComplexity + relationComplexity + propertyComplexity) / 3.0
}

func determinePotentialType(patterns []EmergentPattern) string {
	if len(patterns) == 0 {
		return "unknown"
	}

	// 统计模式类型
	typeCount := make(map[string]int)
	for _, p := range patterns {
		typeCount[p.Type]++
	}

	// 找出最主要的类型
	maxCount := 0
	mainType := "unknown"
	for t, count := range typeCount {
		if count > maxCount {
			maxCount = count
			mainType = t
		}
	}

	return "emergent_" + mainType
}

func extractRequirements(patterns []EmergentPattern) []string {
	reqs := make([]string, 0)
	for _, p := range patterns {
		reqs = append(reqs, p.Type)
	}
	return reqs
}

func estimateTimeFrame(stability float64) time.Duration {
	// 基于稳定性估计时间框架
	return time.Duration(float64(time.Hour) * (1.0 / stability))
}

func calculateRequiredEnergy(patterns []EmergentPattern) float64 {
	energy := 0.0
	for _, p := range patterns {
		if value, exists := p.Properties["energy"]; exists {
			energy += value
		}
	}
	return energy
}

// shouldGenerate 判断是否应该生成新属性
func (pg *PropertyGenerator) shouldGenerate(potential PotentialProperty) bool {
	// 检查概率阈值
	if potential.Probability < pg.config.complexity {
		return false
	}

	// 检查能量条件
	fieldEnergy := pg.field.GetEnergy()
	if fieldEnergy < potential.Energy {
		return false
	}

	// 检查要求条件
	for _, req := range potential.Requirements {
		if !pg.checkRequirement(req) {
			return false
		}
	}

	return true
}

// generateProperty 生成新属性
func (pg *PropertyGenerator) generateProperty(
	potential PotentialProperty,
	patterns []EmergentPattern) error {

	// 创建新属性
	property := &EmergentProperty{
		ID:         generatePropertyID(),
		Type:       potential.Type,
		Created:    time.Now(),
		Updated:    time.Now(),
		Components: make([]PropertyComponent, 0),
		Evolution:  make([]PropertyState, 0),
	}

	// 初始化属性值
	initialState, err := pg.calculateInitialState(potential, patterns)
	if err != nil {
		return err
	}

	property.Value = initialState.Value
	property.Stability = initialState.Stability
	property.Evolution = append(property.Evolution, initialState)

	// 建立组件关联
	components := pg.establishComponents(potential, patterns)
	property.Components = components

	// 记录生成事件
	event := GenerationEvent{
		Timestamp:  time.Now(),
		PropertyID: property.ID,
		Type:       "creation",
		New:        property,
	}
	pg.state.history = append(pg.state.history, event)

	// 保存新属性
	pg.state.properties[property.ID] = property

	return nil
}

// establishComponents 建立组件关联
func (pg *PropertyGenerator) establishComponents(potential PotentialProperty, patterns []EmergentPattern) []PropertyComponent {
	components := make([]PropertyComponent, 0)

	// 遍历所需模式建立关联
	for _, pattern := range patterns {
		// 检查模式是否满足要求
		if !containsPattern(potential.Requirements, pattern.Type) {
			continue
		}

		// 创建组件关联
		component := PropertyComponent{
			PatternID: pattern.ID,
			Weight:    calculateComponentWeight(pattern, potential),
			Role:      determineComponentRole(pattern, potential),
		}

		components = append(components, component)
	}

	return components
}

// 辅助函数
func containsPattern(requirements []string, patternType string) bool {
	for _, req := range requirements {
		if req == patternType {
			return true
		}
	}
	return false
}

func calculateComponentWeight(pattern EmergentPattern, potential PotentialProperty) float64 {
	// 基础权重
	baseWeight := pattern.Strength * (1.0 + pattern.Stability) / 2.0

	// 根据潜在属性类型调整权重
	switch potential.Type {
	case "resonance":
		if pattern.Type == "resonance" {
			baseWeight *= 1.2 // 共振类型加强
		}
	case "field":
		if pattern.Type == "field" {
			baseWeight *= 1.1 // 场类型加强
		}
	case "element":
		if pattern.Type == "element" {
			baseWeight *= 1.15 // 元素类型加强
		}
	}

	return baseWeight
}

func determineComponentRole(pattern EmergentPattern, potential PotentialProperty) string {
	// 基于模式类型和潜在属性类型确定组件角色
	switch {
	case pattern.Type == potential.Type:
		return "core" // 类型匹配作为核心
	case pattern.Type == "resonance":
		return "catalyst" // 共振模式作为催化剂
	case pattern.Type == "field":
		return "container" // 场模式作为容器
	case CheckEnhancingRelation(pattern.Type, potential.Type):
		return "enhancer" // 增强关系
	default:
		return "support" // 默认支持角色
	}
}

// 检查增强关系
func CheckEnhancingRelation(patternType, potentialType string) bool {
	// 检查模式类型是否增强潜在属性类型
	enhancingPairs := map[string][]string{
		"resonance": {"field", "quantum"},
		"field":     {"element", "energy"},
		"element":   {"resonance", "structure"},
	}

	if enhancers, exists := enhancingPairs[patternType]; exists {
		for _, t := range enhancers {
			if t == potentialType {
				return true
			}
		}
	}
	return false
}

// updateProperties 更新现有属性
func (pg *PropertyGenerator) updateProperties(patterns []EmergentPattern) {
	for id, property := range pg.state.properties {
		// 检查属性是否仍然有效
		if valid := pg.validateProperty(property, patterns); !valid {
			delete(pg.state.properties, id)
			continue
		}

		// 更新属性状态
		oldState := copyPropertyState(property)
		if err := pg.evolveProperty(property, patterns); err != nil {
			continue
		}

		// 记录变化
		event := GenerationEvent{
			Timestamp:  time.Now(),
			PropertyID: property.ID,
			Type:       "update",
			Old:        oldState,
			New:        property,
			Changes:    calculatePropertyChanges(oldState, property),
		}
		pg.state.history = append(pg.state.history, event)
	}

	// 限制历史记录长度
	if len(pg.state.history) > maxHistoryLength {
		pg.state.history = pg.state.history[1:]
	}
}

// validateProperty 验证属性是否有效
func (pg *PropertyGenerator) validateProperty(property *EmergentProperty, patterns []EmergentPattern) bool {
	// 检查组件依赖
	for _, comp := range property.Components {
		patternExists := false
		for _, pattern := range patterns {
			if pattern.ID == comp.PatternID {
				patternExists = true
				// 检查模式强度是否足够维持属性
				if pattern.Strength < pg.config.stability {
					return false
				}
				break
			}
		}
		if !patternExists {
			return false
		}
	}

	// 检查能量条件
	requiredEnergy := 0.0
	for _, comp := range property.Components {
		for _, pattern := range patterns {
			if pattern.ID == comp.PatternID {
				if value, exists := pattern.Properties["energy"]; exists {
					requiredEnergy += value * comp.Weight
				}
			}
		}
	}

	fieldEnergy := pg.field.GetEnergy()
	if fieldEnergy < requiredEnergy {
		return false
	}

	// 检查稳定性
	if property.Stability < pg.config.stability {
		return false
	}

	return true
}

// evolveProperty 演化属性
func (pg *PropertyGenerator) evolveProperty(
	property *EmergentProperty,
	patterns []EmergentPattern) error {

	// 计算新状态
	newState, err := pg.calculateNewState(property, patterns)
	if err != nil {
		return err
	}

	// 应用演化
	property.Value = newState.Value
	property.Stability = newState.Stability
	property.Updated = time.Now()
	property.Evolution = append(property.Evolution, newState)

	// 限制演化历史长度
	if len(property.Evolution) > maxEvolutionHistory {
		property.Evolution = property.Evolution[1:]
	}

	return nil
}

// calculateNewState 计算属性的新状态
func (pg *PropertyGenerator) calculateNewState(property *EmergentProperty, patterns []EmergentPattern) (PropertyState, error) {
	// 1. 计算组件影响
	componentEffect := 0.0
	totalWeight := 0.0

	for _, comp := range property.Components {
		// 查找关联的模式
		var pattern *EmergentPattern
		for _, p := range patterns {
			if p.ID == comp.PatternID {
				pattern = &p
				break
			}
		}
		if pattern == nil {
			continue
		}

		// 计算影响
		effect := pattern.Strength * comp.Weight
		componentEffect += effect
		totalWeight += comp.Weight
	}

	if totalWeight > 0 {
		componentEffect /= totalWeight
	}

	// 2. 计算新值
	currentValue := property.Value
	evolutionRate := pg.config.evolutionRate
	newValue := currentValue + (componentEffect-currentValue)*evolutionRate

	// 3. 计算稳定性
	stability := calculateStability(property.Evolution, newValue)

	return PropertyState{
		Timestamp: time.Now(),
		Value:     newValue,
		Stability: stability,
	}, nil
}

// calculateStability 计算稳定性
func calculateStability(history []PropertyState, newValue float64) float64 {
	if len(history) < 2 {
		return 1.0
	}

	// 计算最近值的方差,包括新值
	values := make([]float64, len(history)+1)
	for i, state := range history {
		values[i] = state.Value
	}
	values[len(history)] = newValue

	// 计算平均值
	mean := 0.0
	for _, value := range values {
		mean += value
	}
	mean /= float64(len(values))

	// 计算方差
	variance := 0.0
	for _, value := range values {
		diff := value - mean
		variance += diff * diff
	}
	variance /= float64(len(values))

	return 1.0 / (1.0 + variance)
}

// checkRequirement 检查是否满足要求
func (pg *PropertyGenerator) checkRequirement(req string) bool {
	switch req {
	case "energy":
		return pg.field.GetEnergy() > pg.config.minEnergy
	case "stability":
		return pg.field.GetStability() > pg.config.minStability
	case "coherence":
		return pg.field.GetCoherence() > pg.config.minCoherence
	default:
		return false
	}
}

// calculateInitialState 计算初始状态
func (pg *PropertyGenerator) calculateInitialState(
	potential PotentialProperty,
	patterns []EmergentPattern) (PropertyState, error) {

	// 基于模式计算初始值
	initialValue := 0.0
	totalWeight := 0.0
	for _, pattern := range patterns {
		if weight := pg.calculatePatternWeight(pattern, potential); weight > 0 {
			initialValue += pattern.Strength * weight
			totalWeight += weight
		}
	}

	if totalWeight > 0 {
		initialValue /= totalWeight
	} else {
		// 默认中间值
		initialValue = 0.5
	}

	// 计算初始稳定性 - 使用带模式参数的 calculatePatternStability
	stability := calculatePatternStability(patterns)

	return PropertyState{
		Timestamp: time.Now(),
		Value:     initialValue,
		Stability: stability,
	}, nil
}

// calculatePatternWeight 计算模式权重
func (pg *PropertyGenerator) calculatePatternWeight(pattern EmergentPattern, potential PotentialProperty) float64 {
	// 1. 基础权重 - 基于模式强度和稳定性
	baseWeight := pattern.Strength * (1.0 + pattern.Stability) / 2.0

	// 2. 类型权重 - 基于模式类型和潜在属性的匹配度
	typeWeight := 1.0
	if pattern.Type == potential.Type {
		typeWeight = 1.2 // 类型匹配加权
	}

	// 3. 能量权重 - 基于能量水平的贡献
	energyWeight := 1.0
	if energy, exists := pattern.Properties["energy"]; exists {
		energyRatio := energy / potential.Energy
		energyWeight = math.Min(1.0, energyRatio)
	}

	// 4. 组合权重
	weight := baseWeight * typeWeight * energyWeight

	return weight
}
func generatePropertyID() string {
	return fmt.Sprintf("prop_%d", time.Now().UnixNano())
}

func copyPropertyState(property *EmergentProperty) *EmergentProperty {
	if property == nil {
		return nil
	}

	copy := *property
	return &copy
}

func calculatePropertyChanges(old, new *EmergentProperty) map[string]float64 {
	changes := make(map[string]float64)

	if old != nil && new != nil {
		changes["value"] = new.Value - old.Value
		changes["stability"] = new.Stability - old.Stability
	}

	return changes
}
