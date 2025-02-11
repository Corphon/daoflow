//system/evolution/pattern/generation.go

package pattern

import (
	"fmt"
	"math"
	"math/rand/v2"
	"sort"
	"sync"
	"time"

	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/meta/emergence"
)

// PatternGenerator 模式生成器
type PatternGenerator struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		generationRate float64 // 生成率
		mutationRate   float64 // 变异率
		complexityBias float64 // 复杂度偏好
		energyBalance  float64 // 能量平衡因子
	}

	// 生成状态
	state struct {
		templates  map[string]*GenerationTemplate // 生成模板
		candidates []*PatternCandidate            // 候选模式
		history    []GenerationEvent              // 生成历史
		metrics    GenerationMetrics              // 生成指标
	}

	// 依赖项
	recognizer *PatternRecognizer
	matcher    *EvolutionMatcher
}

// GenerationTemplate 生成模板
type GenerationTemplate struct {
	ID          string                 // 模板ID
	Type        string                 // 模板类型
	Structure   TemplateStructure      // 结构定义
	Constraints []GenerationConstraint // 生成约束
	Properties  map[string]Range       // 属性范围
	Success     float64                // 成功率
	UsageCount  int                    // 使用次数
}

// GenerationStructure 生成结构定义
type GenerationStructure struct {
	Components []ComponentSpec // 组件规格
	Relations  []RelationSpec  // 关系规格
	Dynamics   DynamicsSpec    // 动态规格
}

// TemplateStructure 模板结构
type TemplateStructure struct {
	Components []ComponentSpec // 组件规格
	Relations  []RelationSpec  // 关系规格
	Dynamics   DynamicsSpec    // 动态规格
}

// ComponentSpec 组件规格
type ComponentSpec struct {
	Type       string           // 组件类型
	Required   bool             // 是否必需
	Properties map[string]Range // 属性范围
	Quantity   Range            // 数量范围
}

// RelationSpec 关系规格
type RelationSpec struct {
	Source   string // 源组件类型
	Target   string // 目标组件类型
	Type     string // 关系类型
	Strength Range  // 强度范围
}

// DynamicsSpec 动态规格
type DynamicsSpec struct {
	TimeScale Range           // 时间尺度
	Evolution []EvolutionRule // 演化规则
	Stability Range           // 稳定性范围
}

// EvolutionRule 演化规则
type EvolutionRule struct {
	Type       string             // 规则类型
	Target     string             // 作用目标
	Condition  RuleCondition      // 触发条件
	Effect     map[string]float64 // 效果参数
	TimeWindow Range              // 时间窗口
	Priority   int                // 优先级
}

// RuleCondition 规则条件
type RuleCondition struct {
	Property  string  // 属性
	Operator  string  // 运算符
	Threshold float64 // 阈值
	Tolerance float64 // 容差
}

// Range 数值范围
type Range struct {
	Min       float64
	Max       float64
	Preferred float64
}

// GenerationConstraint 生成约束
type GenerationConstraint struct {
	Type      string      // 约束类型
	Target    string      // 约束目标
	Condition string      // 约束条件
	Value     interface{} // 约束值
}

// PatternCandidate 候选模式
type PatternCandidate struct {
	ID         string                     // 候选ID
	Template   string                     // 使用的模板
	Pattern    *emergence.EmergentPattern // 生成的模式
	Score      float64                    // 评分
	Generation int                        // 生成代数
	Created    time.Time                  // 创建时间
}

// GenerationEvent 生成事件
type GenerationEvent struct {
	Timestamp  time.Time
	TemplateID string
	PatternID  string
	Success    bool
	Score      float64
	Changes    map[string]interface{}
}

// GenerationMetrics 生成指标
type GenerationMetrics struct {
	TotalGenerated int
	SuccessRate    float64
	AverageScore   float64
	ComplexityDist map[float64]int
	Evolution      []MetricPoint
}

// MetricPoint 指标点
type MetricPoint struct {
	Timestamp time.Time
	Metrics   map[string]float64
}

// NewPatternGenerator 创建新的模式生成器
func NewPatternGenerator(
	recognizer *PatternRecognizer,
	matcher *EvolutionMatcher) *PatternGenerator {

	pg := &PatternGenerator{
		recognizer: recognizer,
		matcher:    matcher,
	}

	// 初始化配置
	pg.config.generationRate = 0.3
	pg.config.mutationRate = 0.1
	pg.config.complexityBias = 0.4
	pg.config.energyBalance = 0.7

	// 初始化状态
	pg.state.templates = make(map[string]*GenerationTemplate)
	pg.state.candidates = make([]*PatternCandidate, 0)
	pg.state.history = make([]GenerationEvent, 0)
	pg.state.metrics = GenerationMetrics{
		ComplexityDist: make(map[float64]int),
		Evolution:      make([]MetricPoint, 0),
	}

	return pg
}

// Generate 生成新模式
func (pg *PatternGenerator) Generate() error {
	pg.mu.Lock()
	defer pg.mu.Unlock()

	// 选择模板
	template := pg.selectTemplate()
	if template == nil {
		return model.WrapError(nil, model.ErrCodeOperation, "no suitable template")
	}

	// 生成候选模式
	candidates := pg.generateCandidates(template)

	// 评估候选模式
	evaluated := pg.evaluateCandidates(candidates)

	// 选择最佳候选
	selected := pg.selectBestCandidates(evaluated)

	// 优化选中的模式
	optimized := pg.optimizePatterns(selected)

	// 更新生成指标
	pg.updateMetrics(optimized)

	return nil
}

// selectTemplate 选择合适的生成模板
func (pg *PatternGenerator) selectTemplate() *GenerationTemplate {
	if len(pg.state.templates) == 0 {
		return nil
	}

	// 计算每个模板的得分
	scores := make(map[string]float64)
	for id, template := range pg.state.templates {
		// 基础分数
		score := template.Success

		// 使用频率调整
		usageRate := float64(template.UsageCount) / float64(len(pg.state.history)+1)
		score *= (1.0 - usageRate*0.5) // 降低过度使用的模板分数

		// 复杂度偏好
		if complexity, ok := template.Properties["complexity"]; ok {
			score *= 1.0 + (complexity.Preferred-0.5)*pg.config.complexityBias
		}

		scores[id] = score
	}

	// 选择最高分的模板
	var bestTemplate *GenerationTemplate
	bestScore := -1.0

	for id, score := range scores {
		if score > bestScore {
			bestScore = score
			bestTemplate = pg.state.templates[id]
		}
	}

	return bestTemplate
}

// RegisterTemplate 注册生成模板
func (pg *PatternGenerator) RegisterTemplate(template *GenerationTemplate) error {
	if template == nil {
		return model.WrapError(nil, model.ErrCodeValidation, "nil template")
	}

	pg.mu.Lock()
	defer pg.mu.Unlock()

	// 验证模板
	if err := pg.validateTemplate(template); err != nil {
		return err
	}

	// 存储模板
	pg.state.templates[template.ID] = template

	return nil
}

// generateCandidates 生成候选模式
func (pg *PatternGenerator) generateCandidates(
	template *GenerationTemplate) []*PatternCandidate {

	candidates := make([]*PatternCandidate, 0)

	// 生成多个候选
	for i := 0; i < maxCandidates; i++ {
		// 构建基础模式
		pattern := pg.buildBasePattern(template)

		// 应用变异
		if rand.Float64() < pg.config.mutationRate {
			pattern = pg.mutatePattern(pattern)
		}

		// 创建候选
		candidate := &PatternCandidate{
			ID:         generateCandidateID(),
			Template:   template.ID,
			Pattern:    pattern,
			Generation: 0,
			Created:    time.Now(),
		}

		candidates = append(candidates, candidate)
	}

	return candidates
}

// buildBasePattern 构建基础模式
func (pg *PatternGenerator) buildBasePattern(template *GenerationTemplate) *emergence.EmergentPattern {
	// 创建基础模式
	pattern := &emergence.EmergentPattern{
		ID:         generatePatternID(),
		Type:       template.Type,
		Formation:  time.Now(),
		Components: make([]emergence.PatternComponent, 0),
		Properties: make(map[string]float64),
	}

	// 1. 构建组件
	for _, spec := range template.Structure.Components {
		// 确定组件数量
		count := int(spec.Quantity.Min + rand.Float64()*(spec.Quantity.Max-spec.Quantity.Min))
		for i := 0; i < count; i++ {
			comp := buildComponent(spec)
			pattern.Components = append(pattern.Components, comp)
		}
	}

	// 2. 构建关系
	for _, relation := range template.Structure.Relations {
		establishRelation(pattern, relation)
	}

	// 3. 应用动态规则
	for _, rule := range template.Structure.Dynamics.Evolution {
		applyEvolutionRule(pattern, rule)
	}

	return pattern
}

// establishRelation 建立组件关系
func establishRelation(pattern *emergence.EmergentPattern, relation RelationSpec) {
	// 遍历组件建立关系
	for i := range pattern.Components {
		if pattern.Components[i].Type != relation.Source {
			continue
		}

		for j := range pattern.Components {
			if i == j || pattern.Components[j].Type != relation.Target {
				continue
			}

			// 生成关系强度
			strength := relation.Strength.Min +
				rand.Float64()*(relation.Strength.Max-relation.Strength.Min)

			// 添加关系属性
			if pattern.Components[i].Properties == nil {
				pattern.Components[i].Properties = make(map[string]float64)
			}
			if pattern.Components[j].Properties == nil {
				pattern.Components[j].Properties = make(map[string]float64)
			}

			// 设置关联强度
			key := fmt.Sprintf("relation_%s_%d", relation.Type, j)
			pattern.Components[i].Properties[key] = strength
			key = fmt.Sprintf("relation_%s_%d", relation.Type, i)
			pattern.Components[j].Properties[key] = strength
		}
	}
}

// applyEvolutionRule 应用演化规则
func applyEvolutionRule(pattern *emergence.EmergentPattern, rule EvolutionRule) {
	// 根据规则类型应用不同的演化效果
	switch rule.Type {
	case "energy_transfer":
		applyEnergyTransfer(pattern, rule)
	case "phase_coupling":
		applyPhaseCoupling(pattern, rule)
	case "property_change":
		applyPropertyChange(pattern, rule)
	}
}

// 辅助函数
func applyEnergyTransfer(pattern *emergence.EmergentPattern, rule EvolutionRule) {
	// 能量转移规则
	if energy, ok := rule.Effect["energy"]; ok {
		for i := range pattern.Components {
			if pattern.Components[i].Type == rule.Target {
				pattern.Components[i].Properties["energy"] *= energy
			}
		}
	}
}

func applyPhaseCoupling(pattern *emergence.EmergentPattern, rule EvolutionRule) {
	// 相位耦合规则
	if coupling, ok := rule.Effect["coupling"]; ok {
		for i := range pattern.Components {
			if pattern.Components[i].Type == rule.Target {
				phase := pattern.Components[i].Properties["phase"]
				pattern.Components[i].Properties["phase"] =
					normalizePhase(phase * coupling)
			}
		}
	}
}

func applyPropertyChange(pattern *emergence.EmergentPattern, rule EvolutionRule) {
	// 属性变化规则
	for prop, value := range rule.Effect {
		for i := range pattern.Components {
			if pattern.Components[i].Type == rule.Target {
				pattern.Components[i].Properties[prop] *= value
			}
		}
	}
}

// mutatePattern 变异模式
func (pg *PatternGenerator) mutatePattern(pattern *emergence.EmergentPattern) *emergence.EmergentPattern {
	mutated := *pattern // 复制原模式

	// 1. 组件变异
	for i := range mutated.Components {
		if rand.Float64() < 0.3 { // 30%变异概率
			mutateComponent(&mutated.Components[i])
		}
	}

	// 2. 属性变异
	for key := range mutated.Properties {
		if rand.Float64() < 0.2 { // 20%变异概率
			mutated.Properties[key] *= 0.8 + rand.Float64()*0.4 // 变异范围±20%
		}
	}

	// 3. 能量平衡
	balanceEnergy(&mutated)

	return &mutated
}

// 辅助函数
func buildComponent(spec ComponentSpec) emergence.PatternComponent {
	comp := emergence.PatternComponent{
		Type:       spec.Type,
		Properties: make(map[string]float64),
	}

	// 设置属性
	for prop, rng := range spec.Properties {
		comp.Properties[prop] = rng.Min + rand.Float64()*(rng.Max-rng.Min)
	}

	return comp
}

func mutateComponent(comp *emergence.PatternComponent) {
	// 属性变异
	for key := range comp.Properties {
		if rand.Float64() < 0.3 {
			comp.Properties[key] *= 0.9 + rand.Float64()*0.2
		}
	}
}

func balanceEnergy(pattern *emergence.EmergentPattern) {
	// 能量守恒
	totalEnergy := 0.0
	for _, comp := range pattern.Components {
		if energy, exists := comp.Properties["energy"]; exists {
			totalEnergy += energy
		}
	}

	// 重新分配能量
	if totalEnergy > 0 {
		factor := pattern.Energy / totalEnergy
		for i := range pattern.Components {
			if energy, exists := pattern.Components[i].Properties["energy"]; exists {
				pattern.Components[i].Properties["energy"] = energy * factor
			}
		}
	}
}

// evaluateCandidates 评估候选模式
func (pg *PatternGenerator) evaluateCandidates(
	candidates []*PatternCandidate) []*PatternCandidate {

	for _, candidate := range candidates {
		// 计算基础分数
		baseScore := pg.calculateBaseScore(candidate.Pattern)

		// 评估复杂度
		complexityScore := pg.evaluateComplexity(candidate.Pattern)

		// 检查能量平衡
		energyScore := pg.checkEnergyBalance(candidate.Pattern)

		// 组合得分
		candidate.Score = pg.combineScores(baseScore, complexityScore, energyScore)
	}

	return candidates
}

// calculateBaseScore 计算基础分数
func (pg *PatternGenerator) calculateBaseScore(pattern *emergence.EmergentPattern) float64 {
	// 基于强度和稳定性
	score := pattern.Strength*0.6 + pattern.Stability*0.4

	// 考虑时间衰减
	age := time.Since(pattern.Formation).Hours()
	decay := math.Exp(-age / 24.0) // 24小时衰减周期

	return score * decay
}

// evaluateComplexity 评估复杂度
func (pg *PatternGenerator) evaluateComplexity(pattern *emergence.EmergentPattern) float64 {
	// 组件复杂度
	componentComplexity := calculateComponentComplexity(
		convertComponents(pattern.Components))

	// 结构复杂度
	structuralComplexity := calculateStructuralComplexity(
		extractStructureMap(pattern)) // 解引用

	// 动态复杂度
	dynamicComplexity := calculateDynamicComplexity(
		extractDynamicFeatures(*pattern)) // 解引用

	return (componentComplexity*0.4 + structuralComplexity*0.3 + dynamicComplexity*0.3)
}

// checkEnergyBalance 检查能量平衡
func (pg *PatternGenerator) checkEnergyBalance(pattern *emergence.EmergentPattern) float64 {
	totalEnergy := 0.0
	maxEnergy := 0.0

	// 计算总能量和最大能量
	for _, comp := range pattern.Components {
		if energy, exists := comp.Properties["energy"]; exists {
			totalEnergy += energy
			if energy > maxEnergy {
				maxEnergy = energy
			}
		}
	}

	// 能量平衡度
	if maxEnergy > 0 {
		balance := 1.0 - (maxEnergy-totalEnergy/float64(len(pattern.Components)))/maxEnergy
		return math.Max(0, math.Min(1, balance))
	}

	return 1.0
}

// combineScores 组合分数
func (pg *PatternGenerator) combineScores(baseScore, complexityScore, energyScore float64) float64 {
	// 使用配置的权重
	weightedScore := baseScore*0.5 + complexityScore*0.3 + energyScore*0.2

	// 应用复杂度偏好
	adjustedScore := weightedScore * (1.0 + (complexityScore-0.5)*pg.config.complexityBias)

	return math.Max(0, math.Min(1, adjustedScore))
}

// 辅助函数
func convertComponents(components []emergence.PatternComponent) []SignatureComponent {
	result := make([]SignatureComponent, len(components))
	for i, comp := range components {
		result[i] = convertToSignatureComponent(comp)
	}
	return result
}

func extractStructureMap(pattern *emergence.EmergentPattern) map[string]interface{} {
	structure := make(map[string]interface{})
	structure["hierarchy"] = calculateHierarchyDepth(*pattern)
	structure["connectivity"] = calculateConnectivity(*pattern)
	structure["symmetry"] = calculateStructuralSymmetry(pattern)
	return structure
}

// selectBestCandidates 选择最佳候选
func (pg *PatternGenerator) selectBestCandidates(
	candidates []*PatternCandidate) []*PatternCandidate {

	// 按分数排序
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})

	// 选择前N个
	selected := candidates
	if len(selected) > maxSelected {
		selected = selected[:maxSelected]
	}

	return selected
}

// optimizePatterns 优化模式
func (pg *PatternGenerator) optimizePatterns(
	patterns []*PatternCandidate) []*PatternCandidate {

	optimized := make([]*PatternCandidate, 0)

	for _, pattern := range patterns {
		// 应用优化规则
		improved := pg.optimizePattern(pattern)

		// 检查优化效果
		if improved.Score > pattern.Score {
			optimized = append(optimized, improved)
		} else {
			optimized = append(optimized, pattern)
		}
	}

	return optimized
}

// optimizePattern 优化单个模式候选
func (pg *PatternGenerator) optimizePattern(candidate *PatternCandidate) *PatternCandidate {
	// 创建副本
	optimized := &PatternCandidate{
		ID:         candidate.ID + "_opt",
		Template:   candidate.Template,
		Pattern:    candidate.Pattern.Clone(),
		Generation: candidate.Generation + 1,
		Score:      candidate.Score,
		Created:    time.Now(),
	}

	// 1. 优化组件权重
	for i := range optimized.Pattern.Components {
		comp := &optimized.Pattern.Components[i]
		// 基于组件使用度调整权重
		usage := calculateComponentUsage(comp)
		comp.Weight = adjustWeight(comp.Weight, usage)
	}

	// 2. 优化属性分布
	optimizeProperties(optimized.Pattern)

	// 3. 能量平衡优化
	optimizeEnergyDistribution(optimized.Pattern)

	// 4. 重新评分
	optimized.Score = pg.evaluatePattern(optimized.Pattern)

	return optimized
}

// evaluatePattern 评估单个模式
func (pg *PatternGenerator) evaluatePattern(pattern *emergence.EmergentPattern) float64 {
	// 计算基础分数
	baseScore := pg.calculateBaseScore(pattern)

	// 评估复杂度
	complexityScore := pg.evaluateComplexity(pattern)

	// 检查能量平衡
	energyScore := pg.checkEnergyBalance(pattern)

	// 组合得分
	return pg.combineScores(baseScore, complexityScore, energyScore)
}

// adjustWeight 调整权重
func adjustWeight(current, usage float64) float64 {
	// 根据使用度调整权重
	adjustment := (usage - 0.5) * 0.1 // ±10%调整
	newWeight := current * (1.0 + adjustment)
	return math.Max(0.1, math.Min(1.0, newWeight))
}

// optimizeProperties 优化属性分布
func optimizeProperties(pattern *emergence.EmergentPattern) {
	// 计算属性均值和方差
	for key := range pattern.Properties {
		values := make([]float64, 0)
		for _, comp := range pattern.Components {
			if v, exists := comp.Properties[key]; exists {
				values = append(values, v)
			}
		}

		if len(values) > 0 {
			// 归一化属性分布
			mean := calculateMean(values)
			variance := calculateVariance(values, mean)
			if variance > 0.25 { // 过大的方差
				normalizePropertyDistribution(pattern, key, mean)
			}
		}
	}
}

// optimizeEnergyDistribution 优化能量分布
func optimizeEnergyDistribution(pattern *emergence.EmergentPattern) {
	totalEnergy := pattern.Energy
	components := len(pattern.Components)
	if components == 0 {
		return
	}

	// 计算理想能量分布
	idealEnergy := totalEnergy / float64(components)

	// 调整组件能量
	for i := range pattern.Components {
		if energy, exists := pattern.Components[i].Properties["energy"]; exists {
			// 向理想值靠拢
			diff := idealEnergy - energy
			pattern.Components[i].Properties["energy"] = energy + diff*0.3 // 30%修正
		}
	}
}

// updateMetrics 更新指标
func (pg *PatternGenerator) updateMetrics(patterns []*PatternCandidate) {
	metrics := pg.state.metrics

	// 更新总数
	metrics.TotalGenerated += len(patterns)

	// 计算成功率
	successCount := 0
	totalScore := 0.0

	for _, pattern := range patterns {
		if pattern.Score >= successThreshold {
			successCount++
		}
		totalScore += pattern.Score
	}

	metrics.SuccessRate = float64(successCount) / float64(len(patterns))
	metrics.AverageScore = totalScore / float64(len(patterns))

	// 记录演化点
	point := MetricPoint{
		Timestamp: time.Now(),
		Metrics: map[string]float64{
			"success_rate":  metrics.SuccessRate,
			"average_score": metrics.AverageScore,
		},
	}

	metrics.Evolution = append(metrics.Evolution, point)
}

// 辅助函数

func (pg *PatternGenerator) validateTemplate(template *GenerationTemplate) error {
	if template.ID == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty template ID")
	}

	// 验证结构
	if err := pg.validateStructure(template.Structure); err != nil {
		return err
	}

	// 验证约束
	if err := pg.validateConstraints(template.Constraints); err != nil {
		return err
	}

	return nil
}

// validateStructure 验证模板结构
func (pg *PatternGenerator) validateStructure(structure TemplateStructure) error {
	// 1. 验证组件规格
	if len(structure.Components) == 0 {
		return model.WrapError(nil, model.ErrCodeValidation, "no components defined")
	}

	for _, comp := range structure.Components {
		// 验证组件类型
		if comp.Type == "" {
			return model.WrapError(nil, model.ErrCodeValidation, "empty component type")
		}

		// 验证属性范围
		for _, rng := range comp.Properties {
			if rng.Min > rng.Max {
				return model.WrapError(nil, model.ErrCodeValidation, "invalid property range")
			}
		}

		// 验证数量范围
		if comp.Quantity.Min > comp.Quantity.Max {
			return model.WrapError(nil, model.ErrCodeValidation, "invalid quantity range")
		}
	}

	// 2. 验证关系规格
	for _, rel := range structure.Relations {
		if rel.Source == "" || rel.Target == "" {
			return model.WrapError(nil, model.ErrCodeValidation, "invalid relation")
		}
		if rel.Strength.Min > rel.Strength.Max {
			return model.WrapError(nil, model.ErrCodeValidation, "invalid strength range")
		}
	}

	// 3. 验证动态规格
	if structure.Dynamics.TimeScale.Min > structure.Dynamics.TimeScale.Max {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid time scale")
	}

	return nil
}

// validateConstraints 验证模板约束
func (pg *PatternGenerator) validateConstraints(constraints []GenerationConstraint) error {
	for _, constraint := range constraints {
		// 验证约束类型
		if constraint.Type == "" {
			return model.WrapError(nil, model.ErrCodeValidation, "empty constraint type")
		}

		// 验证约束目标
		if constraint.Target == "" {
			return model.WrapError(nil, model.ErrCodeValidation, "empty constraint target")
		}

		// 验证约束条件
		if constraint.Condition == "" {
			return model.WrapError(nil, model.ErrCodeValidation, "empty constraint condition")
		}

		// 检查条件和目标的匹配性
		if !pg.isValidConstraintCondition(constraint.Type, constraint.Condition) {
			return model.WrapError(nil, model.ErrCodeValidation,
				"invalid constraint condition for type")
		}
	}

	return nil
}

// isValidConstraintCondition 检查约束条件是否有效
func (pg *PatternGenerator) isValidConstraintCondition(
	constraintType string, condition string) bool {

	validConditions := map[string][]string{
		"numeric":     {"gt", "lt", "eq", "gte", "lte"},
		"categorical": {"eq", "neq"},
		"boolean":     {"true", "false"},
	}

	if conditions, exists := validConditions[constraintType]; exists {
		for _, c := range conditions {
			if c == condition {
				return true
			}
		}
	}

	return false
}

func generateCandidateID() string {
	return fmt.Sprintf("cand_%d", time.Now().UnixNano())
}

const (
	maxCandidates    = 100
	maxSelected      = 10
	successThreshold = 0.7
)
