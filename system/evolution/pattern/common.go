// system/evolution/pattern/common.go

package pattern

import (
    "math"
    
    "github.com/Corphon/daoflow/meta/emergence"
    "github.com/Corphon/daoflow/model"
    "github.com/Corphon/daoflow/system/types"
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
    complexity = (componentComplexity * 0.4 + 
                 structuralComplexity * 0.3 + 
                 dynamicComplexity * 0.3)

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
    coherence := (temporalCoherence * 0.4 + 
                 spatialCoherence * 0.3 + 
                 quantumCoherence * 0.3)

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

// extractDynamicFeatures 提取动态特征
func extractDynamicFeatures(pattern emergence.EmergentPattern) map[string]float64 {
    features := make(map[string]float64)

    // 1. 能量特征
    features["energy"] = calculateEnergyFeatures(pattern)
    
    // 2. 演化特征
    features["evolution"] = calculateEvolutionFeatures(pattern)
    
    // 3. 稳定性特征
    features["stability"] = calculateStabilityFeatures(pattern)
    
    // 4. 适应性特征
    features["adaptability"] = calculateAdaptabilityFeatures(pattern)

    return features
}

// determinePatternType 确定模式类型
func determinePatternType(pattern emergence.EmergentPattern) string {
    // 1. 分析模式特征
    features := extractFeatureVector(pattern)
    
    // 2. 计算类型概率
    probabilities := calculateTypeProbs(features)
    
    // 3. 选择最可能的类型
    patternType := selectMostProbableType(probabilities)

    return patternType
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
        "energy":      0.3,
        "evolution":   0.3,
        "stability":   0.2,
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

// 演化特征计算
func calculateEvolutionFeatures(pattern emergence.EmergentPattern) float64 {
    // 演化速率
    evolutionRate := calculateEvolutionRate(pattern)
    
    // 演化方向性
    directionality := calculateEvolutionDirectionality(pattern)
    
    // 演化可预测性
    predictability := calculateEvolutionPredictability(pattern)
    
    return (evolutionRate*0.4 + directionality*0.3 + predictability*0.3)
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

// 标准化函数
func normalizeQuantumValue(value float64) float64 {
    return math.Max(0, math.Min(1, value))
}

func normalizeEnergy(value float64) float64 {
    return math.Max(0, math.Min(1, value/maxEnergyLevel))
}

// 常量定义
const (
    maxEnergyLevel = 1000.0
    minCoherence = 0.1
    maxCoherence = 0.99
)
