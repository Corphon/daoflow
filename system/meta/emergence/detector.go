//system/meta/emergence/detector.go

package emergence

import (
    "math"
    "sync"
    "time"

    "github.com/Corphon/daoflow/meta/field"
    "github.com/Corphon/daoflow/model"
    "github.com/Corphon/daoflow/system/types"
)

// PatternDetector 模式检测器
type PatternDetector struct {
    mu sync.RWMutex

    // 基础配置
    config struct {
        sensitivity     float64         // 检测灵敏度
        timeWindow     time.Duration   // 检测时间窗口
        minConfidence  float64         // 最小置信度
    }

    // 检测状态
    state struct {
        activePatterns  map[string]*EmergentPattern // 活跃模式
        history        []DetectionEvent            // 检测历史
        lastUpdate     time.Time                   // 最后更新时间
    }

    // 场引用
    field *field.UnifiedField
}

// EmergentPattern 涌现模式
type EmergentPattern struct {
    ID          string                // 模式标识
    Type        string                // 模式类型
    Components  []PatternComponent    // 组成成分
    Properties  map[string]float64    // 模式属性
    Strength    float64               // 模式强度
    Stability   float64               // 模式稳定性
    Formation   time.Time             // 形成时间
    LastUpdate  time.Time             // 最后更新时间
}

// PatternComponent 模式组件
type PatternComponent struct {
    Type      string                 // 组件类型
    Weight    float64                // 权重
    Role      string                 // 角色
    State     map[string]float64     // 状态
}

// DetectionEvent 检测事件
type DetectionEvent struct {
    Timestamp  time.Time
    PatternID  string
    Type       string
    Confidence float64
    Changes    []StateChange
}

// StateChange 状态变化
type StateChange struct {
    Component  string
    Before     map[string]float64
    After      map[string]float64
    Delta      float64
}

// NewPatternDetector 创建新的模式检测器
func NewPatternDetector(field *field.UnifiedField) *PatternDetector {
    pd := &PatternDetector{
        field: field,
    }

    // 初始化配置
    pd.config.sensitivity = 0.75
    pd.config.timeWindow = 10 * time.Minute
    pd.config.minConfidence = 0.65

    // 初始化状态
    pd.state.activePatterns = make(map[string]*EmergentPattern)
    pd.state.history = make([]DetectionEvent, 0)
    pd.state.lastUpdate = time.Now()

    return pd
}

// Detect 执行模式检测
func (pd *PatternDetector) Detect() ([]EmergentPattern, error) {
    pd.mu.Lock()
    defer pd.mu.Unlock()

    // 获取场状态
    fieldState, err := pd.field.GetState()
    if err != nil {
        return nil, model.WrapError(err, model.ErrCodeOperation, "failed to get field state")
    }

    // 检测新模式
    newPatterns := pd.detectNewPatterns(fieldState)

    // 更新现有模式
    pd.updateExistingPatterns(fieldState)

    // 移除消失的模式
    pd.removeVanishedPatterns()

    // 记录检测事件
    pd.recordDetectionEvent(newPatterns)

    // 返回当前活跃的模式
    return pd.getActivePatterns(), nil
}

// detectNewPatterns 检测新模式
func (pd *PatternDetector) detectNewPatterns(state *field.FieldState) []EmergentPattern {
    newPatterns := make([]EmergentPattern, 0)

    // 检测元素组合模式
    elementPatterns := pd.detectElementPatterns(state)
    newPatterns = append(newPatterns, elementPatterns...)

    // 检测能量分布模式
    energyPatterns := pd.detectEnergyPatterns(state)
    newPatterns = append(newPatterns, energyPatterns...)

    // 检测量子态模式
    quantumPatterns := pd.detectQuantumPatterns(state)
    newPatterns = append(newPatterns, quantumPatterns...)

    return newPatterns
}

// detectElementPatterns 检测元素组合模式
func (pd *PatternDetector) detectElementPatterns(state *field.FieldState) []EmergentPattern {
    patterns := make([]EmergentPattern, 0)

    // 获取元素状态
    elements := state.GetElements()
    if len(elements) < 2 {
        return patterns
    }

    // 分析元素组合
    combinations := generateElementCombinations(elements)
    for _, combo := range combinations {
        // 检查组合是否形成模式
        if pattern := pd.analyzeElementCombination(combo); pattern != nil {
            patterns = append(patterns, *pattern)
        }
    }

    return patterns
}

// detectEnergyPatterns 检测能量分布模式
func (pd *PatternDetector) detectEnergyPatterns(state *field.FieldState) []EmergentPattern {
    patterns := make([]EmergentPattern, 0)

    // 分析能量分布
    energyDist := state.GetEnergyDistribution()
    
    // 检测能量聚集
    clusters := pd.detectEnergyClusters(energyDist)
    for _, cluster := range clusters {
        if pattern := pd.analyzeEnergyCluster(cluster); pattern != nil {
            patterns = append(patterns, *pattern)
        }
    }

    // 检测能量流动
    flows := pd.detectEnergyFlows(energyDist)
    for _, flow := range flows {
        if pattern := pd.analyzeEnergyFlow(flow); pattern != nil {
            patterns = append(patterns, *pattern)
        }
    }

    return patterns
}

// detectQuantumPatterns 检测量子态模式
func (pd *PatternDetector) detectQuantumPatterns(state *field.FieldState) []EmergentPattern {
    patterns := make([]EmergentPattern, 0)

    // 获取量子态信息
    quantumState := state.GetQuantumState()
    
    // 检测纠缠模式
    entanglements := pd.detectEntanglements(quantumState)
    for _, ent := range entanglements {
        if pattern := pd.analyzeEntanglement(ent); pattern != nil {
            patterns = append(patterns, *pattern)
        }
    }

    // 检测相干模式
    coherences := pd.detectCoherences(quantumState)
    for _, coh := range coherences {
        if pattern := pd.analyzeCoherence(coh); pattern != nil {
            patterns = append(patterns, *pattern)
        }
    }

    return patterns
}

// updateExistingPatterns 更新现有模式
func (pd *PatternDetector) updateExistingPatterns(state *field.FieldState) {
    for id, pattern := range pd.state.activePatterns {
        // 检查模式是否仍然存在
        if exists := pd.verifyPattern(pattern, state); !exists {
            continue
        }

        // 更新模式属性
        pd.updatePatternProperties(pattern, state)

        // 检查模式稳定性
        if pattern.Stability < pd.config.minConfidence {
            delete(pd.state.activePatterns, id)
            continue
        }

        pattern.LastUpdate = time.Now()
    }
}

// verifyPattern 验证模式是否仍然存在
func (pd *PatternDetector) verifyPattern(pattern *EmergentPattern, state *field.FieldState) bool {
    // 检查组件是否仍然存在
    for _, comp := range pattern.Components {
        if !pd.componentExists(comp, state) {
            return false
        }
    }

    // 检查模式强度
    strength := pd.calculatePatternStrength(pattern, state)
    if strength < pd.config.sensitivity {
        return false
    }

    pattern.Strength = strength
    return true
}

// recordDetectionEvent 记录检测事件
func (pd *PatternDetector) recordDetectionEvent(newPatterns []EmergentPattern) {
    event := DetectionEvent{
        Timestamp:  time.Now(),
        Changes:    make([]StateChange, 0),
    }

    // 记录新模式
    for _, pattern := range newPatterns {
        change := StateChange{
            Component: pattern.ID,
            After:     pattern.Properties,
        }
        event.Changes = append(event.Changes, change)
    }

    pd.state.history = append(pd.state.history, event)

    // 限制历史记录长度
    if len(pd.state.history) > maxHistoryLength {
        pd.state.history = pd.state.history[1:]
    }
}

// 辅助函数

func (pd *PatternDetector) componentExists(comp PatternComponent, state *field.FieldState) bool {
    switch comp.Type {
    case "element":
        return state.HasElement(comp.Role)
    case "energy":
        return state.HasEnergyLevel(comp.Weight)
    case "quantum":
        return state.HasQuantumState(comp.State)
    default:
        return false
    }
}

func (pd *PatternDetector) calculatePatternStrength(pattern *EmergentPattern, state *field.FieldState) float64 {
    totalStrength := 0.0
    weightSum := 0.0

    for _, comp := range pattern.Components {
        strength := pd.calculateComponentStrength(comp, state)
        totalStrength += strength * comp.Weight
        weightSum += comp.Weight
    }

    if weightSum == 0 {
        return 0
    }

    return totalStrength / weightSum
}

const (
    maxHistoryLength = 1000 // 最大历史记录长度
)
