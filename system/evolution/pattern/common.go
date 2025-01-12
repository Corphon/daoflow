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

// 更多辅助函数...
