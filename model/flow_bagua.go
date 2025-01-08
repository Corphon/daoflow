// model/flow_bagua.go

package model

import (
    "math"
    "sync"
    "time"

    "github.com/Corphon/daoflow/core"
)

// BaGuaConstants 八卦常数
const (
    TrigramLines    = 3     // 卦象层数
    OctantAngle     = 45.0  // 八分角度
    BasePotential   = 8.0   // 基础势能
    ResonanceRate   = 0.15  // 共振率
    MutationThreshold = 0.7 // 变卦阈值
)

// Trigram 卦象
type Trigram uint8

const (
    Qian Trigram = iota // 乾 ☰
    Kun                 // 坤 ☷
    Zhen               // 震 ☳
    Xun                // 巽 ☴
    Kan                // 坎 ☵
    Li                 // 离 ☲
    Gen                // 艮 ☶
    Dui                // 兑 ☱
)

// TrigramAttributes 卦象属性
type TrigramAttributes struct {
    Lines      [TrigramLines]bool // 爻线
    Direction  float64            // 方位角度
    Element    WuXingPhase        // 关联五行
    Nature     Nature             // 阴阳属性
    Energy     float64            // 能量值
}

// BaGuaFlow 八卦模型
type BaGuaFlow struct {
    *BaseFlowModel
    mu          sync.RWMutex
    
    // 卦象系统
    trigrams    map[Trigram]*TrigramAttributes
    current     Trigram
    
    // 关联系统
    wuxing      *WuXingFlow
    
    // 核心组件
    coreField   *core.Field
    coreState   *core.QuantumState
}

// NewBaGuaFlow 创建八卦流模型
func NewBaGuaFlow(wx *WuXingFlow) *BaGuaFlow {
    bg := &BaGuaFlow{
        BaseFlowModel: NewBaseFlowModel(ModelBaGua, 800.0), // 8卦*100能量
        trigrams:     make(map[Trigram]*TrigramAttributes),
        wuxing:       wx,
        coreField:    core.NewField(),
        coreState:    core.NewQuantumState(),
    }
    
    bg.initializeTrigrams()
    return bg
}

// initializeTrigrams 初始化卦象
func (bg *BaGuaFlow) initializeTrigrams() {
    // 卦象配置
    configs := map[Trigram]struct {
        lines     [TrigramLines]bool
        direction float64
        element   WuXingPhase
        nature    Nature
    }{
        Qian: {[TrigramLines]bool{true, true, true}, 0, Metal, NatureYang},
        Kun:  {[TrigramLines]bool{false, false, false}, 180, Earth, NatureYin},
        Zhen: {[TrigramLines]bool{false, false, true}, 90, Wood, NatureYang},
        Xun:  {[TrigramLines]bool{true, true, false}, 270, Wood, NatureYin},
        Kan:  {[TrigramLines]bool{false, true, false}, 0, Water, NatureYin},
        Li:   {[TrigramLines]bool{true, false, true}, 180, Fire, NatureYang},
        Gen:  {[TrigramLines]bool{false, true, true}, 45, Earth, NatureYang},
        Dui:  {[TrigramLines]bool{true, false, false}, 315, Metal, NatureYin},
    }
    
    // 初始化每个卦象
    for trigram, config := range configs {
        bg.trigrams[trigram] = &TrigramAttributes{
            Lines:     config.lines,
            Direction: config.direction,
            Element:   config.element,
            Nature:    config.nature,
            Energy:    100.0, // 初始能量均匀分布
        }
    }
}

// GetCurrentTrigram 获取当前卦象
func (bg *BaGuaFlow) GetCurrentTrigram() Trigram {
    bg.mu.RLock()
    defer bg.mu.RUnlock()
    return bg.current
}

// GetTrigramAttributes 获取卦象属性
func (bg *BaGuaFlow) GetTrigramAttributes(t Trigram) *TrigramAttributes {
    bg.mu.RLock()
    defer bg.mu.RUnlock()
    return bg.trigrams[t]
}

// Transform 实现状态转换
func (bg *BaGuaFlow) Transform(pattern TransformPattern) error {
    bg.mu.Lock()
    defer bg.mu.Unlock()

    // 计算量子态演化
    bg.coreState.Evolve(time.Second)
    
    // 更新场强度
    probability := bg.coreState.GetProbability()
    bg.coreField.SetStrength(probability)
    
    // 计算新卦象
    newTrigram := bg.calculateNextTrigram(pattern)
    
    // 更新能量分布
    bg.redistributeEnergy(newTrigram)
    
    // 同步五行系统
    if bg.wuxing != nil {
        currentAttrs := bg.trigrams[bg.current]
        newAttrs := bg.trigrams[newTrigram]
        bg.wuxing.AdjustPhaseEnergy(currentAttrs.Element, -currentAttrs.Energy * 0.5)
        bg.wuxing.AdjustPhaseEnergy(newAttrs.Element, newAttrs.Energy * 0.5)
    }
    
    bg.current = newTrigram
    return nil
}

// calculateNextTrigram 计算下一个卦象
func (bg *BaGuaFlow) calculateNextTrigram(pattern TransformPattern) Trigram {
    // 基于量子态概率分布计算转换
    prob := bg.coreState.GetProbability()
    if prob > MutationThreshold {
        // 阴阳变换规则
        attrs := bg.trigrams[bg.current]
        lines := attrs.Lines
        for i := 0; i < TrigramLines; i++ {
            if math.Rand.Float64() < ResonanceRate {
                lines[i] = !lines[i]
            }
        }
        return bg.findTrigramByLines(lines)
    }
    return bg.current
}

// redistributeEnergy 重新分配能量
func (bg *BaGuaFlow) redistributeEnergy(newTrigram Trigram) {
    oldAttrs := bg.trigrams[bg.current]
    newAttrs := bg.trigrams[newTrigram]
    
    transferEnergy := oldAttrs.Energy * ResonanceRate
    oldAttrs.Energy -= transferEnergy
    newAttrs.Energy += transferEnergy
    
    // 更新总能量
    bg.state.Energy = bg.calculateTotalEnergy()
}

// calculateTotalEnergy 计算总能量
func (bg *BaGuaFlow) calculateTotalEnergy() float64 {
    var total float64
    for _, attrs := range bg.trigrams {
        total += attrs.Energy
    }
    return total
}

// findTrigramByLines 根据爻线找到对应卦象
func (bg *BaGuaFlow) findTrigramByLines(lines [TrigramLines]bool) Trigram {
    for t, attrs := range bg.trigrams {
        if attrs.Lines == lines {
            return t
        }
    }
    return bg.current // 如果未找到匹配则保持当前卦象
}
