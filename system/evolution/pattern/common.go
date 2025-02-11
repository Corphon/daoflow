// system/evolution/pattern/common.go

package pattern

import (
	"math"
	"time"

	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/meta/emergence"
)

// 常量定义
const (
	maxEnergyLevel = 1000.0
	minCoherence   = 0.1
	maxCoherence   = 0.99
)

// calculatePatternComplexity 计算模式复杂度
func calculatePatternComplexity(pattern *RecognizedPattern) float64 {
	if pattern == nil {
		return 0
	}

	complexity := 0.0

	// 1. 组件复杂度
	componentComplexity := calculateComponentComplexity(pattern.Signature.Components)

	// 2. 结构复杂度
	structuralComplexity := calculateStructuralComplexity(pattern.Signature.Structure)

	// 3. 动态复杂度
	dynamicComplexity := calculateDynamicComplexity(pattern.Signature.Dynamics)

	// 综合复杂度计算
	complexity = (componentComplexity*0.4 +
		structuralComplexity*0.3 +
		dynamicComplexity*0.3)

	return normalizeComplexity(complexity)
}

// calculatePatternCoherence 计算模式相干性
func calculatePatternCoherence(pattern *RecognizedPattern) float64 {
	if pattern == nil {
		return 0
	}

	// 1. 时间相干性
	temporalCoherence := calculateTemporalCoherence(pattern.Evolution)

	// 2. 空间相干性
	spatialCoherence := calculateSpatialCoherence(pattern.Signature)

	// 3. 量子相干性
	quantumCoherence := calculateQuantumCoherence(pattern)

	// 综合相干性计算
	coherence := (temporalCoherence*0.4 +
		spatialCoherence*0.3 +
		quantumCoherence*0.3)

	return normalizeCoherence(coherence)
}

// extractStructuralFeatures 提取结构特征
func extractStructuralFeatures(pattern emergence.EmergentPattern) map[string]interface{} {
	features := make(map[string]interface{})

	// 1. 拓扑特征
	features["topology"] = extractTopologyFeatures(pattern)

	// 2. 连接特征
	features["connectivity"] = extractConnectivityFeatures(pattern)

	// 3. 对称特征
	features["symmetry"] = extractSymmetryFeatures(pattern)

	// 4. 层级特征
	features["hierarchy"] = extractHierarchyFeatures(pattern)

	return features
}

// 修改extractHierarchyFeatures函数
func extractHierarchyFeatures(pattern emergence.EmergentPattern) map[string]float64 {
	hierarchy := make(map[string]float64)

	// 构建组件层次图
	levels := make(map[string]int)
	maxLevel := 0

	// 基于组件关系确定层次
	for _, comp := range pattern.Components {
		level := 0
		signatureComp1 := convertToSignatureComponent(comp)
		for _, other := range pattern.Components {
			signatureComp2 := convertToSignatureComponent(other)
			if calculateComponentRelation(signatureComp1, signatureComp2) > 0.8 {
				level++
			}
		}
		levels[comp.Role] = level
		if level > maxLevel {
			maxLevel = level
		}
	}

	// ...其余代码不变...
	return hierarchy
}

// convertToSignatureComponent 将PatternComponent转换为SignatureComponent
func convertToSignatureComponent(comp emergence.PatternComponent) SignatureComponent {
	return SignatureComponent{
		Type:        comp.Type,
		Properties:  comp.Properties,
		Weight:      comp.Weight,
		Role:        comp.Role,
		Connections: make([]ComponentConnection, 0), // 暂时为空
	}
}

// calculateComponentRelation 计算组件关系强度
func calculateComponentRelation(c1, c2 SignatureComponent) float64 {
	// 基础关系强度
	baseStrength := math.Min(c1.Weight, c2.Weight)

	// 类型相关性调整
	typeAdjustment := 1.0
	if c1.Type == c2.Type {
		typeAdjustment = 1.2 // 同类型加强
	}

	// 角色关系调整
	roleAdjustment := 1.0
	switch {
	case c1.Type == "element" && c2.Type == "element":
		relation := model.GetWuXingRelation(c1.Role, c2.Role)
		roleAdjustment = relation.Factor

	case c1.Type == "energy" && c2.Type == "energy":
		// 能量梯度关系
		energyDiff := math.Abs(c1.Weight - c2.Weight)
		roleAdjustment = 1.0 / (1.0 + energyDiff)

	case c1.Type == "quantum" && c2.Type == "quantum":
		// 量子关联
		if c1.Properties != nil && c2.Properties != nil {
			ent1 := c1.Properties["entanglement"]
			ent2 := c2.Properties["entanglement"]
			roleAdjustment = math.Sqrt(ent1 * ent2)
		}
	}

	return baseStrength * typeAdjustment * roleAdjustment
}

// calculateHierarchyBalance 计算层级平衡度
func calculateHierarchyBalance(levels map[string]int) float64 {
	if len(levels) == 0 {
		return 0
	}

	// 计算每层节点数
	levelCounts := make(map[int]int)
	for _, level := range levels {
		levelCounts[level]++
	}

	// 计算层级分布的方差
	mean := float64(len(levels)) / float64(len(levelCounts))
	variance := 0.0
	for _, count := range levelCounts {
		diff := float64(count) - mean
		variance += diff * diff
	}
	variance /= float64(len(levelCounts))

	// 归一化到[0,1]区间
	return 1.0 / (1.0 + variance)
}

// extractDynamicFeatures 提取动态特征
func extractDynamicFeatures(pattern emergence.EmergentPattern) map[string]float64 {
	features := make(map[string]float64)

	// 1. 能量特征
	features["energy"] = calculateEnergyFeatures(pattern)

	// 2. 演化特征
	evolutionFeatures := calculateEvolutionFeatures(pattern)
	for k, v := range evolutionFeatures {
		features[k] = v
	}

	// 3. 稳定性特征
	features["stability"] = calculateStabilityFeatures(pattern)

	// 4. 适应性特征
	features["adaptability"] = calculateAdaptabilityFeatures(pattern)

	return features
}

// calculateEvolutionFeatures 计算演化特征
func calculateEvolutionFeatures(pattern emergence.EmergentPattern) map[string]float64 {
	features := make(map[string]float64)
	// 演化速率
	features["rate"] = calculateEvolutionRate(pattern)
	// 演化方向性
	features["directionality"] = calculateEvolutionDirectionality(pattern)
	// 演化可预测性
	features["predictability"] = calculateEvolutionPredictability(pattern)
	return features
}

// calculateEvolutionRate 计算演化速率
func calculateEvolutionRate(pattern emergence.EmergentPattern) float64 {
	if len(pattern.Evolution) < 2 {
		return 0
	}
	// 计算状态变化率
	changes := 0.0
	for i := 1; i < len(pattern.Evolution); i++ {
		diff := calculateStateDifference(
			convertPatternState(pattern.Evolution[i-1]),
			convertPatternState(pattern.Evolution[i]))
		changes += diff
	}
	// 归一化速率
	timeSpan := pattern.Evolution[len(pattern.Evolution)-1].Timestamp.Sub(
		pattern.Evolution[0].Timestamp).Seconds()
	if timeSpan > 0 {
		return math.Min(1.0, changes/timeSpan)
	}
	return 0
}

// calculateEvolutionDirectionality 计算演化方向性
func calculateEvolutionDirectionality(pattern emergence.EmergentPattern) float64 {
	if len(pattern.Evolution) < 3 {
		return 0.5
	}
	// 计算方向一致性
	consistency := 0.0
	prevDirection := 0.0
	for i := 1; i < len(pattern.Evolution)-1; i++ {
		// 计算相邻状态的变化方向
		diff1 := calculateStateDifference(
			convertPatternState(pattern.Evolution[i-1]),
			convertPatternState(pattern.Evolution[i]))
		diff2 := calculateStateDifference(
			convertPatternState(pattern.Evolution[i]),
			convertPatternState(pattern.Evolution[i+1]))
		// 方向相似度
		direction := diff2 - diff1
		if i > 1 {
			// 计算方向一致性
			consistency += math.Cos(math.Atan2(direction, prevDirection))
		}
		prevDirection = direction
	}
	return (consistency/float64(len(pattern.Evolution)-2) + 1) / 2
}

// calculateEvolutionPredictability 计算演化可预测性
func calculateEvolutionPredictability(pattern emergence.EmergentPattern) float64 {
	if len(pattern.Evolution) < 3 {
		return 0.5
	}
	// 使用简单的时间序列分析
	predictions := make([]float64, len(pattern.Evolution)-2)
	actuals := make([]float64, len(pattern.Evolution)-2)
	for i := 2; i < len(pattern.Evolution); i++ {
		// 基于前两个状态预测
		predicted := pattern.Evolution[i-2].Pattern.Properties["energy"] +
			(pattern.Evolution[i-1].Pattern.Properties["energy"] -
				pattern.Evolution[i-2].Pattern.Properties["energy"])
		actual := pattern.Evolution[i].Pattern.Properties["energy"]
		predictions[i-2] = predicted
		actuals[i-2] = actual
	}
	// 计算预测准确度
	error := 0.0
	for i := range predictions {
		if actuals[i] != 0 {
			error += math.Abs(predictions[i]-actuals[i]) / actuals[i]
		}
	}
	return 1.0 - math.Min(1.0, error/float64(len(predictions)))
}

// determinePatternType 确定模式类型
func determinePatternType(pattern emergence.EmergentPattern) string {
	// 1. 分析模式特征
	features := extractFeatureVector(&pattern)

	// 2. 计算类型概率
	probabilities := calculateTypeProbs(features)

	// 3. 选择最可能的类型
	patternType := selectMostProbableType(probabilities)

	return patternType
}

// extractFeatureVector 提取特征向量
func extractFeatureVector(pattern *emergence.EmergentPattern) map[string]float64 {
	features := make(map[string]float64)

	// 基本特征
	features["strength"] = pattern.Strength
	features["stability"] = pattern.Stability

	// 结构特征
	features["complexity"] = pattern.GetStructureComplexity()
	features["coherence"] = pattern.GetStructureCoherence()

	// 动态特征
	dynamic := extractDynamicFeatures(*pattern)
	for k, v := range dynamic {
		features[k] = v
	}

	return features
}

// calculateInitialStability 计算初始稳定性
func calculateInitialStability(pattern emergence.EmergentPattern) float64 {
	// 1. 组件稳定性
	componentStability := calculateComponentsStability(pattern.Components)

	// 2. 结构稳定性
	structuralStability := calculateStructuralStability(pattern)

	// 3. 能量稳定性
	energyStability := calculateEnergyStability(pattern)

	// 加权平均
	return (componentStability*0.4 + structuralStability*0.3 + energyStability*0.3)
}

// calculateTypeProbs 计算类型概率
func calculateTypeProbs(features map[string]float64) map[string]float64 {
	probs := make(map[string]float64)

	// 基于特征计算各类型概率
	probs["resonance"] = calculateResonanceProb(features)
	probs["field"] = calculateFieldProb(features)
	probs["quantum"] = calculateQuantumProb(features)
	probs["element"] = calculateElementProb(features)

	// 归一化概率
	total := 0.0
	for _, p := range probs {
		total += p
	}
	if total > 0 {
		for k := range probs {
			probs[k] /= total
		}
	}

	return probs
}

// 计算共振类型概率
func calculateResonanceProb(features map[string]float64) float64 {
	// 共振类型特征权重
	weights := map[string]float64{
		"coherence": 0.4, // 相干性权重
		"frequency": 0.3, // 频率权重
		"stability": 0.3, // 稳定性权重
	}

	prob := 0.0
	for feat, weight := range weights {
		if value, exists := features[feat]; exists {
			prob += value * weight
		}
	}

	return math.Max(0, math.Min(1, prob))
}

// 计算场类型概率
func calculateFieldProb(features map[string]float64) float64 {
	weights := map[string]float64{
		"strength":   0.4, // 场强权重
		"uniformity": 0.3, // 均匀性权重
		"coupling":   0.3, // 耦合性权重
	}

	prob := 0.0
	for feat, weight := range weights {
		if value, exists := features[feat]; exists {
			prob += value * weight
		}
	}

	return math.Max(0, math.Min(1, prob))
}

// 计算量子类型概率
func calculateQuantumProb(features map[string]float64) float64 {
	weights := map[string]float64{
		"entanglement": 0.4, // 纠缠度权重
		"coherence":    0.3, // 相干性权重
		"purity":       0.3, // 纯度权重
	}

	prob := 0.0
	for feat, weight := range weights {
		if value, exists := features[feat]; exists {
			prob += value * weight
		}
	}

	return math.Max(0, math.Min(1, prob))
}

// 计算元素类型概率
func calculateElementProb(features map[string]float64) float64 {
	weights := map[string]float64{
		"energy":    0.4, // 能量权重
		"stability": 0.3, // 稳定性权重
		"polarity":  0.3, // 极性权重
	}

	prob := 0.0
	for feat, weight := range weights {
		if value, exists := features[feat]; exists {
			prob += value * weight
		}
	}

	return math.Max(0, math.Min(1, prob))
}

// selectMostProbableType 选择最可能类型
func selectMostProbableType(probs map[string]float64) string {
	maxProb := 0.0
	maxType := "unknown"

	for t, p := range probs {
		if p > maxProb {
			maxProb = p
			maxType = t
		}
	}

	// 概率太低时返回unknown
	if maxProb < 0.3 {
		return "unknown"
	}

	return maxType
}

// 辅助函数

func calculateComponentComplexity(components []SignatureComponent) float64 {
	if len(components) == 0 {
		return 0
	}

	complexity := 0.0
	for _, comp := range components {
		// 组件内部复杂度
		internalComplexity := calculateInternalComplexity(comp)
		// 组件关系复杂度
		relationalComplexity := calculateRelationalComplexity(comp)

		complexity += (internalComplexity + relationalComplexity) * comp.Weight
	}

	return complexity / float64(len(components))
}

// calculateInternalComplexity 计算组件内部复杂度
func calculateInternalComplexity(comp SignatureComponent) float64 {
	// 基础复杂度
	baseComplexity := 0.3 // 基础分值

	// 属性复杂度
	propertyComplexity := float64(len(comp.Properties)) * 0.1

	// 类型复杂度
	typeComplexity := 0.0
	switch comp.Type {
	case "quantum":
		typeComplexity = 0.4 // 量子组件最复杂
	case "field":
		typeComplexity = 0.3 // 场组件次之
	case "element":
		typeComplexity = 0.2 // 元素组件再次
	case "energy":
		typeComplexity = 0.1 // 能量组件最简单
	}

	return baseComplexity + propertyComplexity + typeComplexity
}

// calculateRelationalComplexity 计算组件关系复杂度
func calculateRelationalComplexity(comp SignatureComponent) float64 {
	// 关联数量复杂度
	connectionComplexity := float64(len(comp.Connections)) * 0.2

	// 关系类型复杂度
	relationComplexity := 0.0
	for _, conn := range comp.Connections {
		switch conn.Type {
		case "quantum_entanglement":
			relationComplexity += 0.4 // 量子纠缠最复杂
		case "field_coupling":
			relationComplexity += 0.3 // 场耦合次之
		case "energy_transfer":
			relationComplexity += 0.2 // 能量传输再次
		case "element_interaction":
			relationComplexity += 0.1 // 元素相互作用最简单
		}
	}

	if len(comp.Connections) > 0 {
		relationComplexity /= float64(len(comp.Connections))
	}

	return connectionComplexity + relationComplexity
}

func calculateStructuralComplexity(structure map[string]interface{}) float64 {
	complexity := 0.0

	// 分析结构的层次性
	if hierarchy, ok := structure["hierarchy"].(float64); ok {
		complexity += hierarchy * 0.3
	}

	// 分析结构的连通性
	if connectivity, ok := structure["connectivity"].(float64); ok {
		complexity += connectivity * 0.3
	}

	// 分析结构的对称性
	if symmetry, ok := structure["symmetry"].(float64); ok {
		complexity += (1 - symmetry) * 0.4 // 越不对称越复杂
	}

	return complexity
}

func calculateDynamicComplexity(dynamics map[string]float64) float64 {
	complexity := 0.0
	weights := map[string]float64{
		"energy":       0.3,
		"evolution":    0.3,
		"stability":    0.2,
		"adaptability": 0.2,
	}

	for metric, weight := range weights {
		if value, ok := dynamics[metric]; ok {
			complexity += value * weight
		}
	}

	return complexity
}

func normalizeComplexity(value float64) float64 {
	return math.Max(0, math.Min(1, value))
}

func normalizeCoherence(value float64) float64 {
	return math.Max(0, math.Min(1, value))
}

// 时间相关计算
func calculateTemporalCoherence(evolution []PatternState) float64 {
	if len(evolution) < 2 {
		return 1.0 // 单一状态视为完全相干
	}

	coherence := 0.0
	totalWeight := 0.0
	decayFactor := 0.95 // 时间衰减因子

	// 计算状态转换的连续性
	for i := 1; i < len(evolution); i++ {
		weight := math.Pow(decayFactor, float64(len(evolution)-i))
		stateDiff := calculateStateDifference(evolution[i-1], evolution[i])
		coherence += (1.0 - stateDiff) * weight
		totalWeight += weight
	}

	return coherence / totalWeight
}

// calculateStateDifference 计算两个状态之间的差异
func calculateStateDifference(state1, state2 PatternState) float64 {
	// 1. 强度差异
	strengthDiff := math.Abs(state1.Pattern.Strength - state2.Pattern.Strength)

	// 2. 相位差异
	phase1 := state1.Pattern.Properties["phase"]
	phase2 := state2.Pattern.Properties["phase"]
	phaseDiff := normalizePhase(phase1 - phase2)

	// 3. 能量差异
	energy1 := state1.Pattern.Properties["energy"]
	energy2 := state2.Pattern.Properties["energy"]
	energyDiff := math.Abs(energy1 - energy2)

	// 4. 相干性差异
	coherence1 := state1.Pattern.Properties["coherence"]
	coherence2 := state2.Pattern.Properties["coherence"]
	coherenceDiff := math.Abs(coherence1 - coherence2)

	// 加权平均差异
	weights := map[string]float64{
		"strength":  0.3,
		"phase":     0.3,
		"energy":    0.2,
		"coherence": 0.2,
	}

	totalDiff := strengthDiff*weights["strength"] +
		phaseDiff*weights["phase"] +
		energyDiff*weights["energy"] +
		coherenceDiff*weights["coherence"]

	return math.Min(1.0, totalDiff)
}

// normalizePhase 将相位标准化到[-π,π]区间
func normalizePhase(phase float64) float64 {
	// 将相位标准化到 [-π, π] 区间
	for phase > math.Pi {
		phase -= 2 * math.Pi
	}
	for phase < -math.Pi {
		phase += 2 * math.Pi
	}
	return phase
}

// 空间相关计算
func calculateSpatialCoherence(signature PatternSignature) float64 {
	// 计算组件间的空间关联度
	componentCoherence := calculateComponentCoherence(signature.Components)

	// 计算结构的空间一致性
	structuralCoherence := calculateStructuralCoherence(signature.Structure)

	// 计算场的空间分布一致性
	fieldCoherence := calculateFieldCoherence(signature.Dynamics)

	return (componentCoherence*0.4 + structuralCoherence*0.3 + fieldCoherence*0.3)
}

// calculateComponentCoherence 计算组件相干性
func calculateComponentCoherence(components []SignatureComponent) float64 {
	if len(components) < 2 {
		return 1.0
	}

	coherence := 0.0
	pairs := 0

	// 计算组件对之间的相干性
	for i := 0; i < len(components)-1; i++ {
		for j := i + 1; j < len(components); j++ {
			// 计算组件关系
			relation := calculateComponentRelation(components[i], components[j])
			coherence += relation
			pairs++
		}
	}

	if pairs == 0 {
		return 0
	}

	return coherence / float64(pairs)
}

// calculateStructuralCoherence 计算结构相干性
func calculateStructuralCoherence(structure map[string]interface{}) float64 {
	// 检查关键结构特征
	coherence := 0.0
	count := 0.0

	// 层级一致性
	if hierarchy, ok := structure["hierarchy"].(float64); ok {
		coherence += hierarchy
		count++
	}

	// 对称性贡献
	if symmetry, ok := structure["symmetry"].(float64); ok {
		coherence += symmetry
		count++
	}

	// 连通性贡献
	if connectivity, ok := structure["connectivity"].(float64); ok {
		coherence += connectivity
		count++
	}

	if count == 0 {
		return 0
	}

	return coherence / count
}

// calculateFieldCoherence 计算场相干性
func calculateFieldCoherence(dynamics map[string]float64) float64 {
	// 提取关键动态特征
	var phaseCoherence, amplitudeCoherence, energyCoherence float64

	if phase, ok := dynamics["phase"]; ok {
		phaseCoherence = math.Cos(phase) // 相位相干性
	}

	if amplitude, ok := dynamics["amplitude"]; ok {
		amplitudeCoherence = math.Min(amplitude, 1.0) // 振幅相干性
	}

	if energy, ok := dynamics["energy"]; ok {
		energyCoherence = 1.0 / (1.0 + math.Abs(energy-0.5)) // 能量相干性
	}

	return (phaseCoherence + amplitudeCoherence + energyCoherence) / 3.0
}

// 量子相关计算
func calculateQuantumCoherence(pattern *RecognizedPattern) float64 {
	// 1. 计算量子态纯度
	purity := calculateQuantumPurity(pattern)

	// 2. 计算退相干度
	decoherence := calculateDecoherenceFactor(pattern)

	// 3. 计算量子纠缠度
	entanglement := calculateEntanglementDegree(pattern)

	return (purity*0.4 + (1-decoherence)*0.3 + entanglement*0.3)
}

// calculateEntanglementDegree 计算纠缠度
func calculateEntanglementDegree(pattern *RecognizedPattern) float64 {
	if len(pattern.Evolution) < 2 {
		return 0
	}

	entanglement := 0.0
	totalWeight := 0.0
	decayFactor := 0.9

	// 计算历史状态之间的量子纠缠度
	for i := 1; i < len(pattern.Evolution); i++ {
		weight := math.Pow(decayFactor, float64(i))

		// 计算相邻状态间的纠缠度
		state1 := pattern.Evolution[i-1].Pattern.Properties
		state2 := pattern.Evolution[i].Pattern.Properties

		// 计算量子态的相关性
		phase1 := state1["phase"]
		phase2 := state2["phase"]
		phaseDiff := normalizePhase(phase1 - phase2)

		// 使用相位差和态重叠计算纠缠度
		overlap := math.Cos(phaseDiff)
		stateEntanglement := math.Abs(overlap)

		entanglement += stateEntanglement * weight
		totalWeight += weight
	}

	if totalWeight == 0 {
		return 0
	}

	return entanglement / totalWeight
}

// 特征提取相关
func extractTopologyFeatures(pattern emergence.EmergentPattern) map[string]float64 {
	topology := make(map[string]float64)

	// 连通性分析
	topology["connectivity"] = calculateConnectivity(pattern)

	// 环路分析
	topology["cycles"] = detectCycles(pattern)

	// 层级深度
	topology["depth"] = calculateHierarchyDepth(pattern)

	// 分支因子
	topology["branching_factor"] = calculateBranchingFactor(pattern)

	return topology
}

// calculateConnectivity 计算连通度
func calculateConnectivity(pattern emergence.EmergentPattern) float64 {
	if len(pattern.Components) < 2 {
		return 1.0
	}

	// 计算实际连接数
	connections := 0.0
	for i := 0; i < len(pattern.Components)-1; i++ {
		for j := i + 1; j < len(pattern.Components); j++ {
			if hasConnection(pattern.Components[i], pattern.Components[j]) {
				connections++
			}
		}
	}

	// 计算最大可能连接数
	maxConnections := float64(len(pattern.Components)*(len(pattern.Components)-1)) / 2.0

	return connections / maxConnections
}

// hasConnection 检查两个组件之间是否存在连接
func hasConnection(c1, c2 emergence.PatternComponent) bool {
	// 基于组件类型检查
	switch {
	case c1.Type == "quantum" && c2.Type == "quantum":
		// 量子纠缠连接
		if c1.Properties != nil && c2.Properties != nil {
			ent1 := c1.Properties["entanglement"]
			ent2 := c2.Properties["entanglement"]
			// 纠缠度大于阈值认为存在连接
			if math.Sqrt(ent1*ent2) > 0.5 {
				return true
			}
		}

	case c1.Type == "field" && c2.Type == "field":
		// 场耦合连接
		if c1.Properties != nil && c2.Properties["coupling"] > 0.5 {
			return true
		}

	case c1.Type == "element" && c2.Type == "element":
		// 五行相生相克关系
		relation := model.GetWuXingRelation(c1.Role, c2.Role)
		if relation.Factor > 0 {
			return true
		}

	case c1.Type == "energy" && c2.Type == "energy":
		// 能量梯度连接
		if math.Abs(c1.Weight-c2.Weight) < 0.3 {
			return true
		}
	}

	return false
}

// detectCycles 检测环路
func detectCycles(pattern emergence.EmergentPattern) float64 {
	// 构建邻接矩阵
	n := len(pattern.Components)
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}

	// 填充邻接矩阵
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if hasConnection(pattern.Components[i], pattern.Components[j]) {
				adj[i][j] = true
				adj[j][i] = true
			}
		}
	}

	// 统计环路数
	cycles := countCycles(adj)
	return float64(cycles) / float64(n)
}

// countCycles 使用DFS统计环路数
func countCycles(adj [][]bool) int {
	n := len(adj)
	if n == 0 {
		return 0
	}

	visited := make([]bool, n)
	parent := make([]int, n)
	cycleCount := 0

	var dfs func(int, int)
	dfs = func(v int, p int) {
		visited[v] = true
		parent[v] = p

		// 检查所有邻接节点
		for u := 0; u < n; u++ {
			if !adj[v][u] {
				continue
			}

			// 未访问的节点
			if !visited[u] {
				dfs(u, v)
			} else if u != p && u != parent[v] {
				// 发现环路
				cycleCount++
			}
		}
	}

	// 对每个未访问的节点进行DFS
	for i := 0; i < n; i++ {
		if !visited[i] {
			dfs(i, -1)
		}
	}

	// 由于每个环被计数两次,需要除以2
	return cycleCount / 2
}

// calculateHierarchyDepth 计算层级深度
func calculateHierarchyDepth(pattern emergence.EmergentPattern) float64 {
	levels := make(map[string]int)
	maxLevel := 0

	// 基于组件关系构建层级
	for _, comp := range pattern.Components {
		level := calculateComponentLevel(comp, pattern.Components)
		levels[comp.ID] = level
		if level > maxLevel {
			maxLevel = level
		}
	}

	return float64(maxLevel) / 10.0 // 归一化
}

// calculateComponentLevel 计算组件层级
func calculateComponentLevel(comp emergence.PatternComponent, allComps []emergence.PatternComponent) int {
	level := 0

	// 转换为SignatureComponent进行计算
	signatureComp1 := convertToSignatureComponent(comp)

	// 计算与其他组件的关系来确定层级
	for _, other := range allComps {
		signatureComp2 := convertToSignatureComponent(other)
		// 强关联增加层级
		if calculateComponentRelation(signatureComp1, signatureComp2) > 0.8 {
			level++
		}
	}

	return level
}

// calculateBranchingFactor 计算分支因子
func calculateBranchingFactor(pattern emergence.EmergentPattern) float64 {
	if len(pattern.Components) < 2 {
		return 0
	}

	// 统计每个组件的分支数
	branches := 0.0
	for _, comp := range pattern.Components {
		branchCount := countComponentBranches(comp, pattern.Components)
		branches += float64(branchCount)
	}

	// 计算平均分支数
	return branches / float64(len(pattern.Components))
}

// countComponentBranches 统计组件分支数
func countComponentBranches(comp emergence.PatternComponent, allComps []emergence.PatternComponent) int {
	branchCount := 0

	// 遍历所有其他组件检查连接
	for _, other := range allComps {
		if comp.ID == other.ID {
			continue
		}

		// 使用已有的hasConnection函数检查连接
		if hasConnection(comp, other) {
			branchCount++
		}
	}

	return branchCount
}

func extractConnectivityFeatures(pattern emergence.EmergentPattern) map[string]float64 {
	connectivity := make(map[string]float64)

	// 局部连接密度
	connectivity["local_density"] = calculateLocalDensity(pattern)

	// 全局连接强度
	connectivity["global_strength"] = calculateGlobalStrength(pattern)

	// 连接分布均匀度
	connectivity["distribution"] = calculateConnectionDistribution(pattern)

	// 连接稳定性
	connectivity["stability"] = calculateConnectionStability(pattern)

	return connectivity
}

// calculateLocalDensity 计算局部连接密度
func calculateLocalDensity(pattern emergence.EmergentPattern) float64 {
	if len(pattern.Components) < 2 {
		return 0
	}

	// 统计每个组件的局部连接数
	localDensities := make([]float64, len(pattern.Components))
	for i, comp := range pattern.Components {
		connections := 0
		for j, other := range pattern.Components {
			if i != j && hasConnection(comp, other) {
				connections++
			}
		}
		localDensities[i] = float64(connections) / float64(len(pattern.Components)-1)
	}

	// 计算平均局部密度
	totalDensity := 0.0
	for _, density := range localDensities {
		totalDensity += density
	}

	return totalDensity / float64(len(pattern.Components))
}

// calculateGlobalStrength 计算全局连接强度
func calculateGlobalStrength(pattern emergence.EmergentPattern) float64 {
	totalStrength := 0.0
	connections := 0

	// 累加所有连接的强度
	for i := 0; i < len(pattern.Components)-1; i++ {
		for j := i + 1; j < len(pattern.Components); j++ {
			if hasConnection(pattern.Components[i], pattern.Components[j]) {
				strength := calculateComponentRelation(
					convertToSignatureComponent(pattern.Components[i]),
					convertToSignatureComponent(pattern.Components[j]))
				totalStrength += strength
				connections++
			}
		}
	}

	if connections == 0 {
		return 0
	}

	return totalStrength / float64(connections)
}

// calculateConnectionDistribution 计算连接分布均匀度
func calculateConnectionDistribution(pattern emergence.EmergentPattern) float64 {
	// 统计每个组件的连接数
	connectionCounts := make([]int, len(pattern.Components))
	for i, comp := range pattern.Components {
		for j, other := range pattern.Components {
			if i != j && hasConnection(comp, other) {
				connectionCounts[i]++
			}
		}
	}

	// 计算连接数分布的方差
	mean := 0.0
	for _, count := range connectionCounts {
		mean += float64(count)
	}
	mean /= float64(len(connectionCounts))

	variance := 0.0
	for _, count := range connectionCounts {
		diff := float64(count) - mean
		variance += diff * diff
	}
	variance /= float64(len(connectionCounts))

	// 均匀度 = 1 / (1 + 方差)
	return 1.0 / (1.0 + variance)
}

// calculateConnectionStability 计算连接稳定性
func calculateConnectionStability(pattern emergence.EmergentPattern) float64 {
	if len(pattern.Components) < 2 {
		return 1.0
	}

	// 基于组件权重计算连接稳定性
	stability := 0.0
	count := 0

	for i := 0; i < len(pattern.Components)-1; i++ {
		for j := i + 1; j < len(pattern.Components); j++ {
			if hasConnection(pattern.Components[i], pattern.Components[j]) {
				// 连接越强,权重越接近,稳定性越高
				weightDiff := math.Abs(pattern.Components[i].Weight -
					pattern.Components[j].Weight)
				stability += 1.0 / (1.0 + weightDiff)
				count++
			}
		}
	}

	if count == 0 {
		return 0
	}

	return stability / float64(count)
}

func extractSymmetryFeatures(pattern emergence.EmergentPattern) map[string]float64 {
	symmetry := make(map[string]float64)

	// 空间对称性
	symmetry["spatial"] = calculateSpatialSymmetry(pattern)

	// 时间对称性
	symmetry["temporal"] = calculateTemporalSymmetry(pattern)

	// 量子对称性
	symmetry["quantum"] = calculateQuantumSymmetry(pattern)

	// 场对称性
	symmetry["field"] = calculateFieldSymmetry(pattern)

	return symmetry
}

// calculateSpatialSymmetry 计算空间对称性
func calculateSpatialSymmetry(pattern emergence.EmergentPattern) float64 {
	// 1. 组件空间分布对称性
	componentSymmetry := calculateComponentSymmetry(pattern.Components)

	// 2. 拓扑结构对称性
	topologySymmetry := calculateTopologySymmetry(pattern.Components)

	// 3. 属性分布对称性
	propertySymmetry := calculatePropertySymmetry(pattern.Properties)

	// 加权平均
	return componentSymmetry*0.4 + topologySymmetry*0.3 + propertySymmetry*0.3
}

// calculateComponentSymmetry 计算组件对称性
func calculateComponentSymmetry(components []emergence.PatternComponent) float64 {
	n := len(components)
	if n < 2 {
		return 0
	}

	symmetricPairs := 0.0
	totalPairs := float64(n * (n - 1) / 2)

	// 检查组件对的对称性
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			// 类型和权重相近的组件对被认为是对称的
			if components[i].Type == components[j].Type &&
				math.Abs(components[i].Weight-components[j].Weight) < 0.1 {
				symmetricPairs++
			}
		}
	}

	return symmetricPairs / totalPairs
}

// calculateTopologySymmetry 计算拓扑对称性
func calculateTopologySymmetry(components []emergence.PatternComponent) float64 {
	n := len(components)
	if n < 2 {
		return 0
	}

	// 构建距离矩阵
	distances := make([][]float64, n)
	for i := range distances {
		distances[i] = make([]float64, n)
	}

	// 计算组件间距离
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			typeDist := 0.0
			if components[i].Type == components[j].Type {
				typeDist = 1.0
			}
			weightDist := 1.0 - math.Abs(components[i].Weight-components[j].Weight)
			dist := (typeDist + weightDist) / 2.0
			distances[i][j] = dist
			distances[j][i] = dist
		}
	}

	// 检查拓扑对称性
	symmetry := 0.0
	pairs := 0
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			for k := 0; k < n-1; k++ {
				for l := k + 1; l < n; l++ {
					if (i != k || j != l) &&
						math.Abs(distances[i][j]-distances[k][l]) < 0.1 {
						symmetry += 1.0
					}
					pairs++
				}
			}
		}
	}

	if pairs > 0 {
		return symmetry / float64(pairs)
	}
	return 0
}

// calculatePropertySymmetry 计算属性对称性
func calculatePropertySymmetry(properties map[string]float64) float64 {
	if len(properties) == 0 {
		return 0
	}

	// 提取属性值
	values := make([]float64, 0, len(properties))
	for _, v := range properties {
		values = append(values, v)
	}

	// 计算属性分布的偏度作为对称性度量
	mean := calculateMean(values)
	variance := calculateVariance(values, mean)
	skewness := calculateSkewness(values, mean, variance)

	// 偏度越小表示分布越对称
	return 1.0 / (1.0 + math.Abs(skewness))
}

// calculateMean 计算平均值
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

// calculateSkewness 计算偏度
func calculateSkewness(values []float64, mean float64, variance float64) float64 {
	if len(values) == 0 || variance == 0 {
		return 0
	}

	stdDev := math.Sqrt(variance)
	sum := 0.0
	for _, v := range values {
		diff := (v - mean) / stdDev
		sum += diff * diff * diff
	}
	return sum / float64(len(values))
}

// calculateTemporalSymmetry 计算时间对称性
func calculateTemporalSymmetry(pattern emergence.EmergentPattern) float64 {
	if len(pattern.Evolution) < 2 {
		return 1.0
	}

	// 计算状态序列的对称性
	symmetry := 0.0
	n := len(pattern.Evolution)
	for i := 0; i < n/2; i++ {
		// 对比前后状态的相似度
		diff := calculateStateDifference(
			convertPatternState(pattern.Evolution[i]),
			convertPatternState(pattern.Evolution[n-1-i]),
		)
		symmetry += 1.0 - diff
	}

	return symmetry / float64(n/2)
}

// calculateQuantumSymmetry 计算量子对称性
func calculateQuantumSymmetry(pattern emergence.EmergentPattern) float64 {
	// 提取量子属性
	phase := 0.0
	if value, exists := pattern.Properties["phase"]; exists {
		phase = value
	}

	// 相位对称性
	phaseSymmetry := math.Cos(phase)

	// 纠缠对称性
	entanglement := 0.0
	if value, exists := pattern.Properties["entanglement"]; exists {
		entanglement = value
	}

	return (math.Abs(phaseSymmetry) + entanglement) / 2.0
}

// calculateFieldSymmetry 计算场对称性
func calculateFieldSymmetry(pattern emergence.EmergentPattern) float64 {
	// 1. 强度分布对称性
	strengthSymmetry := 0.0
	if value, exists := pattern.Properties["field_strength"]; exists {
		strengthSymmetry = 1.0 - math.Abs(value-0.5)*2
	}

	// 2. 梯度对称性
	gradientSymmetry := 0.0
	if value, exists := pattern.Properties["field_gradient"]; exists {
		gradientSymmetry = 1.0 - value
	}

	return (strengthSymmetry + gradientSymmetry) / 2.0
}

// 能量特征计算
func calculateEnergyFeatures(pattern emergence.EmergentPattern) float64 {
	// 基础能量
	baseEnergy := pattern.Energy

	// 量子贡献
	quantumContribution := calculateQuantumEnergyContribution(pattern)

	// 场贡献
	fieldContribution := calculateFieldEnergyContribution(pattern)

	// 综合能量
	totalEnergy := baseEnergy + quantumContribution + fieldContribution

	return normalizeEnergy(totalEnergy)
}

// calculateQuantumEnergyContribution 计算量子能量贡献
func calculateQuantumEnergyContribution(pattern emergence.EmergentPattern) float64 {
	// 检查量子属性
	if value, exists := pattern.Properties["quantum_energy"]; exists {
		return value
	}

	// 计算量子组件的能量贡献
	energy := 0.0
	for _, comp := range pattern.Components {
		if comp.Type == "quantum" {
			// 考虑纠缠和相干性
			if ent, ok := comp.Properties["entanglement"]; ok {
				energy += ent * comp.Weight
			}
			if coh, ok := comp.Properties["coherence"]; ok {
				energy += coh * comp.Weight
			}
		}
	}
	return energy
}

// calculateFieldEnergyContribution 计算场能量贡献
func calculateFieldEnergyContribution(pattern emergence.EmergentPattern) float64 {
	// 检查场属性
	if value, exists := pattern.Properties["field_energy"]; exists {
		return value
	}

	// 计算场组件的能量贡献
	energy := 0.0
	for _, comp := range pattern.Components {
		if comp.Type == "field" {
			// 考虑场强度和耦合
			if strength, ok := comp.Properties["field_strength"]; ok {
				energy += strength * comp.Weight
			}
			if coupling, ok := comp.Properties["coupling"]; ok {
				energy += coupling * comp.Weight
			}
		}
	}
	return energy
}

// 稳定性特征计算
func calculateStabilityFeatures(pattern emergence.EmergentPattern) float64 {
	// 结构稳定性
	structuralStability := calculateStructuralStability(pattern)

	// 动态稳定性
	dynamicStability := calculateDynamicStability(pattern)

	// 量子稳定性
	quantumStability := calculateQuantumStability(pattern)

	return (structuralStability*0.4 + dynamicStability*0.3 + quantumStability*0.3)
}

// 计算结构稳定性
func calculateStructuralStability(pattern emergence.EmergentPattern) float64 {
	// 获取拓扑特征
	topology := extractTopologyFeatures(pattern)
	// 获取连接特征
	connectivity := extractConnectivityFeatures(pattern)

	// 计算拓扑稳定性（连通性越高越稳定）
	topoStability := (topology["connectivity"] + topology["depth"]) / 2.0

	// 计算连接稳定性（分布越均匀越稳定）
	connStability := connectivity["stability"] * connectivity["distribution"]

	// 综合结构稳定性
	return math.Min(1.0, (topoStability*0.6 + connStability*0.4))
}

// 计算动态稳定性
func calculateDynamicStability(pattern emergence.EmergentPattern) float64 {
	// 分析演化特征
	evolution := calculateEvolutionFeatures(pattern)

	// 演化速率越慢越稳定
	rateStability := 1.0 - evolution["rate"]

	// 方向一致性越高越稳定
	directionStability := evolution["directionality"]

	// 综合动态稳定性
	return math.Min(1.0, (rateStability*0.5 + directionStability*0.5))
}

// 计算量子稳定性
func calculateQuantumStability(pattern emergence.EmergentPattern) float64 {
	// 获取量子特征
	quantum := extractQuantumFeatures(pattern)

	// 量子纯度越高越稳定
	purityStability := quantum["purity"]

	// 退相干程度越低越稳定
	decoherenceStability := 1.0 - quantum["decoherence"]

	// 纠缠持久性
	entanglementStability := quantum["entanglement"]

	// 综合量子稳定性
	return math.Min(1.0, (purityStability*0.4 + decoherenceStability*0.4 + entanglementStability*0.2))
}

// extractQuantumFeatures 提取量子特征
func extractQuantumFeatures(pattern emergence.EmergentPattern) map[string]float64 {
	quantum := make(map[string]float64)

	// 从模式属性中提取量子特征
	// 1. 提取纯度
	if value, exists := pattern.Properties["quantum_purity"]; exists {
		quantum["purity"] = value
	} else {
		// 从组件中计算纯度
		purity := 0.0
		count := 0
		for _, comp := range pattern.Components {
			if comp.Type == "quantum" {
				if p, ok := comp.Properties["purity"]; ok {
					purity += p
					count++
				}
			}
		}
		if count > 0 {
			quantum["purity"] = purity / float64(count)
		} else {
			quantum["purity"] = 0.5 // 默认值
		}
	}

	// 2. 提取退相干度
	if value, exists := pattern.Properties["decoherence"]; exists {
		quantum["decoherence"] = value
	} else {
		quantum["decoherence"] = calculateDecoherenceFromComponents(pattern.Components)
	}

	// 3. 提取纠缠度
	if value, exists := pattern.Properties["entanglement"]; exists {
		quantum["entanglement"] = value
	} else {
		quantum["entanglement"] = calculateEntanglementFromComponents(pattern.Components)
	}

	return quantum
}

// calculateDecoherenceFromComponents 从组件计算退相干度
func calculateDecoherenceFromComponents(components []emergence.PatternComponent) float64 {
	decoherence := 0.0
	count := 0

	for _, comp := range components {
		if comp.Type == "quantum" {
			if d, ok := comp.Properties["decoherence"]; ok {
				decoherence += d
				count++
			}
		}
	}

	if count > 0 {
		return decoherence / float64(count)
	}
	return 0.5 // 默认值
}

// calculateEntanglementFromComponents 从组件计算纠缠度
func calculateEntanglementFromComponents(components []emergence.PatternComponent) float64 {
	entanglement := 0.0
	count := 0

	for _, comp := range components {
		if comp.Type == "quantum" {
			if e, ok := comp.Properties["entanglement"]; ok {
				entanglement += e
				count++
			}
		}
	}

	if count > 0 {
		return entanglement / float64(count)
	}
	return 0.3 // 默认值
}

// 适应性特征计算
func calculateAdaptabilityFeatures(pattern emergence.EmergentPattern) float64 {
	// 响应性
	responsiveness := calculateResponseCapability(pattern)

	// 学习能力
	learningCapability := calculateLearningCapability(pattern)

	// 自组织能力
	selfOrganization := calculateSelfOrganization(pattern)

	return (responsiveness*0.4 + learningCapability*0.3 + selfOrganization*0.3)
}

// calculateSelfOrganization 计算自组织能力
func calculateSelfOrganization(pattern emergence.EmergentPattern) float64 {
	// 计算自组织的稳定性和灵活性
	stability := calculateSelfOrganizationStability(pattern)
	flexibility := calculateSelfOrganizationFlexibility(pattern)

	// 综合自组织能力
	return (stability*0.5 + flexibility*0.5)
}

// calculateSelfOrganizationStability 计算自组织稳定性
func calculateSelfOrganizationStability(pattern emergence.EmergentPattern) float64 {
	// 计算结构稳定性
	structuralStability := calculateStructuralStability(pattern)

	// 计算动态稳定性
	dynamicStability := calculateDynamicStability(pattern)

	// 综合稳定性
	return (structuralStability*0.5 + dynamicStability*0.5)
}

// calculateSelfOrganizationFlexibility 计算自组织灵活性
func calculateSelfOrganizationFlexibility(pattern emergence.EmergentPattern) float64 {
	// 计算响应灵敏度
	sensitivity := calculateResponseSensitivity(pattern)

	// 计算学习能力
	learningCapability := calculateLearningCapability(pattern)

	// 综合灵活性
	return (sensitivity*0.5 + learningCapability*0.5)
}

// calculateResponseCapability 计算响应能力
func calculateResponseCapability(pattern emergence.EmergentPattern) float64 {
	// 计算响应速度和灵敏度
	speed := calculateResponseSpeed(pattern)
	sensitivity := calculateResponseSensitivity(pattern)

	// 综合响应能力
	return (speed*0.5 + sensitivity*0.5)
}

// calculateResponseSensitivity 计算响应灵敏度
func calculateResponseSensitivity(pattern emergence.EmergentPattern) float64 {
	if len(pattern.Evolution) < 2 {
		return 0
	}

	// 计算状态变化的敏感度
	sensitivity := 0.0
	for i := 1; i < len(pattern.Evolution); i++ {
		diff := calculateStateDifference(
			convertPatternState(pattern.Evolution[i-1]),
			convertPatternState(pattern.Evolution[i]))
		sensitivity += diff
	}

	return math.Min(1.0, sensitivity/float64(len(pattern.Evolution)-1))
}

// calculateLearningCapability 计算学习能力
func calculateLearningCapability(pattern emergence.EmergentPattern) float64 {
	// 计算学习速率和准确度
	learningRate := calculateLearningRate(pattern)
	accuracy := calculateLearningAccuracy(pattern)

	// 综合学习能力
	return (learningRate*0.5 + accuracy*0.5)
}

// calculateLearningAccuracy 计算学习准确度
func calculateLearningAccuracy(pattern emergence.EmergentPattern) float64 {
	if len(pattern.Evolution) < 3 {
		return 0.5
	}

	// 使用简单的时间序列分析
	predictions := make([]float64, len(pattern.Evolution)-2)
	actuals := make([]float64, len(pattern.Evolution)-2)

	for i := 2; i < len(pattern.Evolution); i++ {
		// 基于前两个状态预测
		predicted := pattern.Evolution[i-2].Pattern.Properties["energy"] +
			(pattern.Evolution[i-1].Pattern.Properties["energy"] -
				pattern.Evolution[i-2].Pattern.Properties["energy"])

		actual := pattern.Evolution[i].Pattern.Properties["energy"]

		predictions[i-2] = predicted
		actuals[i-2] = actual
	}

	// 计算预测准确度
	error := 0.0
	for i := range predictions {
		if actuals[i] != 0 {
			error += math.Abs(predictions[i]-actuals[i]) / actuals[i]
		}
	}

	return 1.0 - math.Min(1.0, error/float64(len(predictions)))
}

// calculateResponseSpeed 计算响应速度
func calculateResponseSpeed(pattern emergence.EmergentPattern) float64 {
	if len(pattern.Evolution) < 2 {
		return 0
	}

	// 计算状态变化率
	changes := 0.0
	for i := 1; i < len(pattern.Evolution); i++ {
		diff := calculateStateDifference(
			convertPatternState(pattern.Evolution[i-1]),
			convertPatternState(pattern.Evolution[i]))
		changes += diff
	}

	// 归一化速率
	timeSpan := pattern.Evolution[len(pattern.Evolution)-1].Timestamp.Sub(
		pattern.Evolution[0].Timestamp).Seconds()
	if timeSpan > 0 {
		return math.Min(1.0, changes/timeSpan)
	}
	return 0
}

// calculateLearningRate 计算学习速率
func calculateLearningRate(pattern emergence.EmergentPattern) float64 {
	if len(pattern.Evolution) < 2 {
		return 0
	}

	// 计算状态变化率
	changes := 0.0
	for i := 1; i < len(pattern.Evolution); i++ {
		diff := calculateStateDifference(
			convertPatternState(pattern.Evolution[i-1]),
			convertPatternState(pattern.Evolution[i]))
		changes += diff
	}

	// 归一化速率
	timeSpan := pattern.Evolution[len(pattern.Evolution)-1].Timestamp.Sub(
		pattern.Evolution[0].Timestamp).Seconds()
	if timeSpan > 0 {
		return math.Min(1.0, changes/timeSpan)
	}
	return 0
}

// 量子态计算
func calculateQuantumPurity(pattern *RecognizedPattern) float64 {
	if pattern == nil || len(pattern.Evolution) == 0 {
		return 0
	}

	// 获取量子态信息
	state := pattern.Evolution[len(pattern.Evolution)-1]

	// 计算密度矩阵
	densityMatrix := calculateDensityMatrix(state)

	// 计算迹
	purity := calculateMatrixTrace(densityMatrix)

	return normalizeQuantumValue(purity)
}

// calculateDensityMatrix 计算量子态的密度矩阵
func calculateDensityMatrix(state PatternState) [][]complex128 {
	// 基于Properties中的能量和相位构造密度矩阵
	energy := state.Properties["energy"]
	phase := state.Properties["phase"]

	densityMatrix := make([][]complex128, 2)
	for i := range densityMatrix {
		densityMatrix[i] = make([]complex128, 2)
	}

	// 构造简化的密度矩阵
	theta := phase * math.Pi // 将相位转换为角度
	p := energy              // 用能量表示概率

	// 填充密度矩阵元素
	densityMatrix[0][0] = complex(p, 0)
	densityMatrix[0][1] = complex(math.Sqrt(p*(1-p))*math.Cos(theta), math.Sqrt(p*(1-p))*math.Sin(theta))
	densityMatrix[1][0] = complex(math.Sqrt(p*(1-p))*math.Cos(theta), -math.Sqrt(p*(1-p))*math.Sin(theta))
	densityMatrix[1][1] = complex(1-p, 0)

	return densityMatrix
}

// calculateMatrixTrace 计算矩阵的迹
func calculateMatrixTrace(matrix [][]complex128) float64 {
	trace := 0.0
	for i := range matrix {
		trace += real(matrix[i][i])
	}
	return trace
}

// 退相干计算
func calculateDecoherenceFactor(pattern *RecognizedPattern) float64 {
	if len(pattern.Evolution) < 2 {
		return 0
	}

	decoherence := 0.0
	totalWeight := 0.0
	decayFactor := 0.9

	// 计算量子相干性随时间的衰减
	for i := 1; i < len(pattern.Evolution); i++ {
		weight := math.Pow(decayFactor, float64(i))
		stateDiff := calculateQuantumStateDifference(
			pattern.Evolution[i-1],
			pattern.Evolution[i],
		)
		decoherence += stateDiff * weight
		totalWeight += weight
	}

	return normalizeQuantumValue(decoherence / totalWeight)
}

// calculateQuantumStateDifference 计算两个量子态之间的差异
func calculateQuantumStateDifference(state1, state2 PatternState) float64 {
	// 1. 相位差异
	phase1 := state1.Pattern.Properties["phase"]
	phase2 := state2.Pattern.Properties["phase"]
	phaseDiff := normalizePhase(phase1 - phase2)

	// 2. 纠缠度差异
	entanglement1 := state1.Pattern.Properties["entanglement"]
	entanglement2 := state2.Pattern.Properties["entanglement"]
	entanglementDiff := math.Abs(entanglement1 - entanglement2)

	// 3. 相干性差异
	coherence1 := state1.Pattern.Properties["coherence"]
	coherence2 := state2.Pattern.Properties["coherence"]
	coherenceDiff := math.Abs(coherence1 - coherence2)

	// 加权平均差异
	weights := map[string]float64{
		"phase":        0.4,
		"entanglement": 0.3,
		"coherence":    0.3,
	}

	totalDiff := phaseDiff*weights["phase"] +
		entanglementDiff*weights["entanglement"] +
		coherenceDiff*weights["coherence"]

	return math.Min(1.0, totalDiff)
}

// 标准化函数
func normalizeQuantumValue(value float64) float64 {
	return math.Max(0, math.Min(1, value))
}

func normalizeEnergy(value float64) float64 {
	return math.Max(0, math.Min(1, value/maxEnergyLevel))
}

// calculateStructuralSymmetry 计算结构对称性
func calculateStructuralSymmetry(pattern *emergence.EmergentPattern) float64 {
	if pattern == nil || len(pattern.Components) < 2 {
		return 0
	}

	// 1. 组件对称性
	componentSymmetry := calculateComponentSymmetry(pattern.Components)

	// 2. 拓扑对称性
	topologySymmetry := calculateTopologySymmetry(pattern.Components)

	// 3. 属性对称性
	propertySymmetry := calculatePropertySymmetry(pattern.Properties)

	// 加权平均
	symmetry := componentSymmetry*0.4 + topologySymmetry*0.3 + propertySymmetry*0.3

	return math.Max(0, math.Min(1, symmetry)) // 确保在0-1范围内
}

// calculateComponentUsage 计算组件使用度
func calculateComponentUsage(comp *emergence.PatternComponent) float64 {
	// 基础使用度
	usage := 0.3 // 基础值

	// 基于属性增加使用度
	for _, v := range comp.Properties {
		if v > 0 {
			usage += 0.1
		}
	}

	// 基于角色调整
	switch comp.Role {
	case "core":
		usage *= 1.2
	case "catalyst":
		usage *= 1.1
	}

	return math.Min(1.0, usage)
}

// normalizePropertyDistribution 标准化属性分布
func normalizePropertyDistribution(pattern *emergence.EmergentPattern, key string, mean float64) {
	// 调整参数 - 允许一定程度的波动
	tolerance := 0.2  // 允许20%的波动
	adjustRate := 0.3 // 每次调整30%

	// 遍历组件调整属性值
	for i := range pattern.Components {
		if value, exists := pattern.Components[i].Properties[key]; exists {
			// 计算与均值的偏差
			diff := value - mean
			if math.Abs(diff) > tolerance*mean {
				// 向均值靠拢,但保留一定随机性
				adjustment := diff * adjustRate
				newValue := value - adjustment

				// 确保值在有效范围内[0,1]
				pattern.Components[i].Properties[key] = math.Max(0, math.Min(1, newValue))
			}
		}
	}

	// 更新模式的整体属性
	pattern.Properties[key] = calculateMean(extractValues(pattern.Components, key))
}

// extractValues 提取属性值
func extractValues(components []emergence.PatternComponent, key string) []float64 {
	values := make([]float64, 0)
	for _, comp := range components {
		if v, exists := comp.Properties[key]; exists {
			values = append(values, v)
		}
	}
	return values
}

// convertToEmergentPattern 将RecognizedPattern转换为EmergentPattern
func convertToEmergentPattern(recognized *RecognizedPattern) emergence.EmergentPattern {
	if recognized == nil || recognized.Pattern == nil {
		return emergence.EmergentPattern{}
	}
	return *recognized.Pattern
}

// extractEvolutionFeatures 从RecognizedPattern中提取演化特征
func extractEvolutionFeatures(pattern *RecognizedPattern) map[string]float64 {
	// 使用convertToEmergentPattern将RecognizedPattern转换为EmergentPattern
	emergentPattern := convertToEmergentPattern(pattern)

	// 复用已有的calculateEvolutionFeatures方法
	features := calculateEvolutionFeatures(emergentPattern)

	// 添加额外的识别模式特征
	features["recognition_confidence"] = pattern.Confidence
	features["activation_level"] = pattern.GetActivationLevel()
	features["evolution_stage"] = float64(len(pattern.Evolution))

	// 标准化所有特征值到[0,1]区间
	for k, v := range features {
		features[k] = math.Max(0, math.Min(1, v))
	}

	return features
}

// calculateEnvironmentSimilarity 计算环境相似度
func calculateEnvironmentSimilarity(envBase, env1, env2 map[string]float64) float64 {
	if len(env1) == 0 || len(env2) == 0 {
		return 0
	}

	similarity := 0.0
	count := 0.0

	// 比较环境因素
	for key, baseVal := range envBase {
		if val1, ok1 := env1[key]; ok1 {
			if val2, ok2 := env2[key]; ok2 {
				// 相对于基准环境的变化率
				delta1 := math.Abs(val1 - baseVal)
				delta2 := math.Abs(val2 - baseVal)
				// 变化率的相似度
				similarity += 1.0 - math.Abs(delta1-delta2)/(delta1+delta2+1e-6)
				count++
			}
		}
	}

	if count == 0 {
		return 0
	}
	return similarity / count
}

// calculateStateSimilarity 计算状态相似度
func calculateStateSimilarity(source, target *RecognizedPattern) float64 {
	// 1. 激活状态相似度
	activationSim := 1.0
	if source.Active != target.Active {
		activationSim = 0.0
	}

	// 2. 置信度相似度
	confidenceSim := 1.0 - math.Abs(source.Confidence-target.Confidence)

	// 3. 演化阶段相似度
	evolutionSim := calculateEvolutionStageSimilarity(source, target)

	return (activationSim*0.3 + confidenceSim*0.3 + evolutionSim*0.4)
}

// calculateEvolutionStageSimilarity 计算演化阶段相似度
func calculateEvolutionStageSimilarity(source, target *RecognizedPattern) float64 {
	if len(source.Evolution) == 0 || len(target.Evolution) == 0 {
		return 0
	}

	// 1. 阶段数量差异
	stageDiff := math.Abs(float64(len(source.Evolution) - len(target.Evolution)))
	stageRatio := 1.0 - math.Min(1.0, stageDiff/float64(len(source.Evolution)))

	// 2. 最新状态相似度
	sourceLatest := source.Evolution[len(source.Evolution)-1]
	targetLatest := target.Evolution[len(target.Evolution)-1]
	latestSim := 1.0 - calculateStateDifference(
		convertPatternState(convertLocalPatternState(sourceLatest)),
		convertPatternState(convertLocalPatternState(targetLatest)))

	// 3. 演化趋势相似度
	sourceTrend := calculateEvolutionDirectionality(convertToEmergentPattern(source))
	targetTrend := calculateEvolutionDirectionality(convertToEmergentPattern(target))
	trendSim := 1.0 - math.Abs(sourceTrend-targetTrend)

	return (stageRatio*0.3 + latestSim*0.4 + trendSim*0.3)
}

// 环境因素相关计算函数
// normalizeTimeOfDay 标准化一天中的时间 (0-1)
func normalizeTimeOfDay(t time.Time) float64 {
	hour := float64(t.Hour()) + float64(t.Minute())/60.0
	return hour / 24.0
}

// calculateSystemEnergy 计算系统能量水平
func calculateSystemEnergy(em *EvolutionMatcher) float64 {
	if len(em.state.patterns) == 0 {
		return 0
	}

	totalEnergy := 0.0
	for _, pattern := range em.state.patterns {
		if pattern.Active {
			totalEnergy += pattern.Pattern.Energy
		}
	}

	return math.Min(1.0, totalEnergy/float64(len(em.state.patterns)))
}

// calculateSystemStability 计算系统稳定性
func calculateSystemStability(em *EvolutionMatcher) float64 {
	if len(em.state.patterns) == 0 {
		return 1.0
	}

	totalStability := 0.0
	count := 0
	for _, pattern := range em.state.patterns {
		if pattern.Active {
			totalStability += pattern.Stability
			count++
		}
	}

	if count == 0 {
		return 1.0
	}
	return totalStability / float64(count)
}

// calculateChangeRate 计算变化率
func calculateChangeRate(lastState ContextState, currentEnv map[string]float64) float64 {
	if len(lastState.Factors) == 0 {
		return 0
	}

	// 计算环境因素的变化率
	totalChange := 0.0
	count := 0.0

	for key, currentValue := range currentEnv {
		if lastValue, exists := lastState.Factors[key]; exists {
			change := math.Abs(currentValue - lastValue)
			totalChange += change
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return totalChange / count
}

// 辅助函数
func calculateComponentsStability(components []emergence.PatternComponent) float64 {
	if len(components) == 0 {
		return 0
	}

	totalStability := 0.0
	totalWeight := 0.0

	for _, comp := range components {
		weight := comp.Weight
		stability := 1.0 - math.Abs(0.5-weight)*2 // 权重越接近0.5越稳定
		totalStability += stability * weight
		totalWeight += weight
	}

	if totalWeight == 0 {
		return 0
	}
	return totalStability / totalWeight
}

func calculateEnergyStability(pattern emergence.EmergentPattern) float64 {
	if pattern.Energy == 0 {
		return 1.0 // 无能量波动为最稳定状态
	}

	// 能量分布的均匀性
	energyVariance := calculateEnergyVariance(pattern)
	return 1.0 - math.Min(1.0, energyVariance)
}

// calculateEnergyVariance 计算能量方差
func calculateEnergyVariance(pattern emergence.EmergentPattern) float64 {
	// 如果没有组件,返回0
	if len(pattern.Components) == 0 {
		return 0
	}

	// 收集所有组件的能量值
	energies := make([]float64, 0)
	totalEnergy := 0.0

	for _, comp := range pattern.Components {
		if energy, exists := comp.Properties["energy"]; exists {
			energies = append(energies, energy)
			totalEnergy += energy
		}
	}

	if len(energies) == 0 {
		return 0
	}

	// 计算平均能量
	meanEnergy := totalEnergy / float64(len(energies))

	// 计算方差
	variance := 0.0
	for _, energy := range energies {
		diff := energy - meanEnergy
		variance += diff * diff
	}
	variance /= float64(len(energies))

	// 归一化方差到[0,1]区间
	return math.Min(1.0, variance/meanEnergy)
}

// calculateSignatureSimilarity 计算签名相似度
func calculateSignatureSimilarity(sig1, sig2 PatternSignature) float64 {
	// 1. 组件相似度
	componentSimilarity := calculateComponentsSimilarity(sig1.Components, sig2.Components)

	// 2. 结构相似度
	structureSimilarity := calculateStructureMapSimilarity(sig1.Structure, sig2.Structure)

	// 3. 动态特征相似度
	dynamicSimilarity := calculatePropertySimilarity(sig1.Dynamics, sig2.Dynamics)

	// 4. 上下文相似度
	contextSimilarity := calculateContextMapSimilarity(sig1.Context, sig2.Context)

	// 加权平均
	return (componentSimilarity*0.4 +
		structureSimilarity*0.3 +
		dynamicSimilarity*0.2 +
		contextSimilarity*0.1)
}

// calculateComponentsSimilarity 计算组件集合相似度
func calculateComponentsSimilarity(comps1, comps2 []SignatureComponent) float64 {
	if len(comps1) == 0 || len(comps2) == 0 {
		return 0
	}

	totalSimilarity := 0.0
	maxSimilarities := make([]float64, len(comps1))

	// 对每个组件找到最佳匹配
	for i, c1 := range comps1 {
		maxSim := 0.0
		for _, c2 := range comps2 {
			sim := calculateComponentSimilarity(c1, c2)
			if sim > maxSim {
				maxSim = sim
			}
		}
		maxSimilarities[i] = maxSim
	}

	// 计算平均相似度
	for _, sim := range maxSimilarities {
		totalSimilarity += sim
	}

	return totalSimilarity / float64(len(comps1))
}

// calculateComponentSimilarity 计算单个组件相似度
func calculateComponentSimilarity(c1, c2 SignatureComponent) float64 {
	// 1. 类型相似度
	typeSimilarity := 0.0
	if c1.Type == c2.Type {
		typeSimilarity = 1.0
	}

	// 2. 权重相似度
	weightSimilarity := 1.0 - math.Abs(c1.Weight-c2.Weight)

	// 3. 属性相似度
	propertySimilarity := calculatePropertySimilarity(c1.Properties, c2.Properties)

	// 4. 角色相似度
	roleSimilarity := 0.0
	if c1.Role == c2.Role {
		roleSimilarity = 1.0
	} else if c1.Type == "element" && c2.Type == "element" {
		// 元素类型考虑五行关系
		relation := model.GetWuXingRelation(c1.Role, c2.Role)
		roleSimilarity = relation.Factor
	}

	// 5. 连接相似度
	connectionSimilarity := calculateConnectionSimilarity(c1.Connections, c2.Connections)

	// 加权平均计算总相似度
	return (typeSimilarity*0.3 +
		weightSimilarity*0.2 +
		propertySimilarity*0.2 +
		roleSimilarity*0.2 +
		connectionSimilarity*0.1)
}

// calculateConnectionSimilarity 计算连接相似度
func calculateConnectionSimilarity(conns1, conns2 []ComponentConnection) float64 {
	if len(conns1) == 0 && len(conns2) == 0 {
		return 1.0
	}
	if len(conns1) == 0 || len(conns2) == 0 {
		return 0.0
	}

	// 计算共同连接类型的比例
	commonTypes := make(map[string]bool)
	allTypes := make(map[string]bool)

	for _, conn := range conns1 {
		allTypes[conn.Type] = true
	}
	for _, conn := range conns2 {
		allTypes[conn.Type] = true
		if _, exists := allTypes[conn.Type]; exists {
			commonTypes[conn.Type] = true
		}
	}

	return float64(len(commonTypes)) / float64(len(allTypes))
}

// calculateStructureMapSimilarity 计算结构映射相似度
func calculateStructureMapSimilarity(m1, m2 map[string]interface{}) float64 {
	if len(m1) == 0 || len(m2) == 0 {
		return 0
	}

	similarity := 0.0
	count := 0.0

	for key, val1 := range m1 {
		if val2, exists := m2[key]; exists {
			// 根据值的类型计算相似度
			switch v1 := val1.(type) {
			case float64:
				if v2, ok := val2.(float64); ok {
					similarity += 1.0 - math.Abs(v1-v2)
					count++
				}
			case string:
				if v2, ok := val2.(string); ok {
					if v1 == v2 {
						similarity += 1.0
					}
					count++
				}
			}
		}
	}

	if count == 0 {
		return 0
	}
	return similarity / count
}

// calculateContextMapSimilarity 计算上下文映射相似度
func calculateContextMapSimilarity(m1, m2 map[string]string) float64 {
	if len(m1) == 0 || len(m2) == 0 {
		return 0
	}

	matches := 0.0
	total := float64(len(m1))

	for key, val1 := range m1 {
		if val2, exists := m2[key]; exists && val1 == val2 {
			matches++
		}
	}

	return matches / total
}

// calculatePatternStability 计算模式稳定性
func calculatePatternStability(pattern *RecognizedPattern) float64 {
	if pattern == nil {
		return 0
	}

	// 1. 时间稳定性
	timeStability := calculateTimeStability(pattern)

	// 2. 结构稳定性
	structStability := calculateStructuralStability(convertToEmergentPattern(pattern))

	// 3. 动态稳定性
	dynamicStability := calculateDynamicStability(convertToEmergentPattern(pattern))

	// 4. 量子稳定性
	quantumStability := calculateQuantumStability(convertToEmergentPattern(pattern))

	// 加权平均计算总稳定性
	stability := (timeStability*0.3 +
		structStability*0.3 +
		dynamicStability*0.2 +
		quantumStability*0.2)

	return math.Max(0, math.Min(1, stability))
}

// calculateTimeStability 计算时间稳定性
func calculateTimeStability(pattern *RecognizedPattern) float64 {
	if len(pattern.Evolution) == 0 {
		return 1.0
	}

	// 基于出现频率的稳定性
	frequencyStability := math.Min(1.0, float64(pattern.Occurrences)/100.0)

	// 基于持续时间的稳定性
	duration := time.Since(pattern.FirstSeen).Hours()
	durationStability := math.Min(1.0, duration/24.0) // 24小时作为参考

	// 基于历史变化的稳定性
	variationStability := calculateTemporalCoherence(pattern.Evolution)

	return (frequencyStability*0.3 + durationStability*0.3 + variationStability*0.4)
}
