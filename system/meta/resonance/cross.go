//system/meta/resonance/cross.go

package resonance

import (
    "math"
    "sync"
    "time"

    "github.com/Corphon/daoflow/meta/field"
    "github.com/Corphon/daoflow/meta/emergence"
    "github.com/Corphon/daoflow/model"
    "github.com/Corphon/daoflow/system/types"
)

// CrossResonance 跨层共振处理器
type CrossResonance struct {
    mu sync.RWMutex

    // 基础配置
    config struct {
        resonanceThreshold float64      // 共振阈值
        couplingStrength   float64      // 耦合强度
        interactionRange   float64      // 交互范围
        decayFactor       float64      // 衰减因子
    }

    // 共振状态
    state struct {
        layers      map[string]*ResonanceLayer  // 共振层
        bridges     map[string]*ResonanceBridge // 层间桥接
        transitions []TransitionEvent           // 转换事件
    }

    // 依赖项
    field      *field.UnifiedField
    amplifier  *ResonanceAmplifier
    matcher    *PatternMatcher
}

// ResonanceLayer 共振层
type ResonanceLayer struct {
    ID         string                 // 层ID
    Level      int                    // 层级
    Patterns   map[string]*emergence.EmergentPattern // 层内模式
    Properties map[string]float64     // 层属性
    Energy     float64                // 层能量
    Coherence  float64                // 相干度
    Created    time.Time              // 创建时间
}

// ResonanceBridge 共振桥接
type ResonanceBridge struct {
    ID          string               // 桥接ID
    SourceLayer string               // 源层ID
    TargetLayer string               // 目标层ID
    Strength    float64              // 桥接强度
    Phase       float64              // 相位差
    Channels    []BridgeChannel      // 传输通道
    Active      bool                 // 是否活跃
    Created     time.Time            // 创建时间
}

// BridgeChannel 桥接通道
type BridgeChannel struct {
    Type       string               // 通道类型
    Capacity   float64              // 传输容量
    Load       float64              // 当前负载
    Direction  int                  // 传输方向
    State      string               // 通道状态
}

// TransitionEvent 转换事件
type TransitionEvent struct {
    Timestamp    time.Time
    SourceLayer  string
    TargetLayer  string
    Type         string
    Pattern      *emergence.EmergentPattern
    Energy       float64
    Success      bool
}

// NewCrossResonance 创建新的跨层共振处理器
func NewCrossResonance(
    field *field.UnifiedField,
    amplifier *ResonanceAmplifier,
    matcher *PatternMatcher) *CrossResonance {
    
    cr := &CrossResonance{
        field:     field,
        amplifier: amplifier,
        matcher:   matcher,
    }

    // 初始化配置
    cr.config.resonanceThreshold = 0.7
    cr.config.couplingStrength = 0.5
    cr.config.interactionRange = 2.0
    cr.config.decayFactor = 0.1

    // 初始化状态
    cr.state.layers = make(map[string]*ResonanceLayer)
    cr.state.bridges = make(map[string]*ResonanceBridge)
    cr.state.transitions = make([]TransitionEvent, 0)

    return cr
}

// Process 处理跨层共振
func (cr *CrossResonance) Process() error {
    cr.mu.Lock()
    defer cr.mu.Unlock()

    // 更新层状态
    if err := cr.updateLayers(); err != nil {
        return err
    }

    // 处理层间共振
    if err := cr.processInterLayerResonance(); err != nil {
        return err
    }

    // 管理桥接
    if err := cr.manageBridges(); err != nil {
        return err
    }

    // 处理转换
    if err := cr.handleTransitions(); err != nil {
        return err
    }

    return nil
}

// updateLayers 更新层状态
func (cr *CrossResonance) updateLayers() error {
    // 获取当前模式
    patterns, err := cr.matcher.GetActivePatterns()
    if err != nil {
        return err
    }

    // 对模式进行分层
    layers := cr.stratifyPatterns(patterns)

    // 更新层状态
    for level, layerPatterns := range layers {
        layerID := fmt.Sprintf("layer_%d", level)
        
        if layer, exists := cr.state.layers[layerID]; exists {
            // 更新现有层
            cr.updateExistingLayer(layer, layerPatterns)
        } else {
            // 创建新层
            cr.createNewLayer(layerID, level, layerPatterns)
        }
    }

    // 移除空层
    cr.removeEmptyLayers()

    return nil
}

// processInterLayerResonance 处理层间共振
func (cr *CrossResonance) processInterLayerResonance() error {
    // 遍历相邻层对
    for id1, layer1 := range cr.state.layers {
        for id2, layer2 := range cr.state.layers {
            if id1 == id2 || !cr.areLayersAdjacent(layer1, layer2) {
                continue
            }

            // 检查层间共振
            if resonance := cr.detectLayerResonance(layer1, layer2); resonance != nil {
                // 处理共振效应
                if err := cr.handleLayerResonance(resonance); err != nil {
                    continue
                }
            }
        }
    }

    return nil
}

// manageBridges 管理桥接
func (cr *CrossResonance) manageBridges() error {
    // 更新现有桥接
    for id, bridge := range cr.state.bridges {
        if valid := cr.validateBridge(bridge); !valid {
            delete(cr.state.bridges, id)
            continue
        }

        // 更新桥接状态
        if err := cr.updateBridge(bridge); err != nil {
            continue
        }
    }

    // 检测新的桥接机会
    newBridges := cr.detectNewBridges()
    for _, bridge := range newBridges {
        cr.state.bridges[bridge.ID] = bridge
    }

    return nil
}

// handleTransitions 处理转换
func (cr *CrossResonance) handleTransitions() error {
    for _, bridge := range cr.state.bridges {
        if !bridge.Active {
            continue
        }

        // 检查转换条件
        if transitions := cr.detectTransitions(bridge); len(transitions) > 0 {
            // 执行转换
            for _, transition := range transitions {
                if err := cr.executeTransition(transition); err != nil {
                    continue
                }
            }
        }
    }

    return nil
}

// stratifyPatterns 对模式进行分层
func (cr *CrossResonance) stratifyPatterns(
    patterns []*emergence.EmergentPattern) map[int][]*emergence.EmergentPattern {
    
    layers := make(map[int][]*emergence.EmergentPattern)

    for _, pattern := range patterns {
        // 计算模式层级
        level := cr.calculatePatternLevel(pattern)
        
        // 添加到对应层
        if _, exists := layers[level]; !exists {
            layers[level] = make([]*emergence.EmergentPattern, 0)
        }
        layers[level] = append(layers[level], pattern)
    }

    return layers
}

// detectLayerResonance 检测层间共振
func (cr *CrossResonance) detectLayerResonance(
    layer1, layer2 *ResonanceLayer) *LayerResonance {
    
    // 计算层间耦合
    coupling := cr.calculateLayerCoupling(layer1, layer2)
    if coupling < cr.config.resonanceThreshold {
        return nil
    }

    // 检查相位匹配
    if !cr.checkPhaseMatch(layer1, layer2) {
        return nil
    }

    // 创建层间共振
    resonance := &LayerResonance{
        ID:          generateResonanceID(),
        SourceLayer: layer1.ID,
        TargetLayer: layer2.ID,
        Strength:    coupling,
        Created:     time.Now(),
    }

    return resonance
}

// executeTransition 执行转换
func (cr *CrossResonance) executeTransition(transition *TransitionEvent) error {
    // 获取源层和目标层
    sourceLayer := cr.state.layers[transition.SourceLayer]
    targetLayer := cr.state.layers[transition.TargetLayer]
    
    if sourceLayer == nil || targetLayer == nil {
        return model.WrapError(nil, model.ErrCodeValidation, "invalid layers")
    }

    // 执行模式转换
    if err := cr.transferPattern(transition.Pattern, sourceLayer, targetLayer); err != nil {
        return err
    }

    // 更新能量分布
    cr.redistributeEnergy(sourceLayer, targetLayer, transition.Energy)

    // 记录转换事件
    cr.recordTransition(transition)

    return nil
}

// 辅助函数

func (cr *CrossResonance) calculatePatternLevel(pattern *emergence.EmergentPattern) int {
    // 基于模式复杂度和能量计算层级
    complexity := calculatePatternComplexity(pattern)
    energy := pattern.Energy
    
    // 使用对数尺度
    level := int(math.Log2(complexity * energy))
    return math.Max(0, float64(level))
}

func (cr *CrossResonance) calculateLayerCoupling(layer1, layer2 *ResonanceLayer) float64 {
    // 计算层间耦合强度
    energyCoupling := math.Sqrt(layer1.Energy * layer2.Energy)
    coherenceCoupling := (layer1.Coherence + layer2.Coherence) / 2
    
    return energyCoupling * coherenceCoupling * cr.config.couplingStrength
}

func (cr *CrossResonance) recordTransition(transition *TransitionEvent) {
    cr.state.transitions = append(cr.state.transitions, *transition)

    // 限制历史记录长度
    if len(cr.state.transitions) > maxTransitionHistory {
        cr.state.transitions = cr.state.transitions[1:]
    }
}

const (
    maxTransitionHistory = 1000
)

type LayerResonance struct {
    ID          string
    SourceLayer string
    TargetLayer string
    Strength    float64
    Created     time.Time
}
