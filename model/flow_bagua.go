// model/flow_bagua.go

package model

import (
    "math"
    "time"

    "github.com/Corphon/daoflow/core"
)

// BaGuaConstants 八卦常数
const (
    TrigramCount = 8         // 八卦数量
    LayerCount   = 3         // 三才层数
    BaseRate     = 0.125     // 基础转换率 (1/8)
    SyncRate     = 0.2       // 同步率
    PhaseShift   = math.Pi/4 // 相位偏移
)

// BaGuaTrigram 八卦卦象
type BaGuaTrigram uint8

const (
    Qian BaGuaTrigram = iota // 乾
    Dui                      // 兑
    Li                       // 离
    Zhen                     // 震
    Xun                      // 巽
    Kan                      // 坎
    Gen                      // 艮
    Kun                      // 坤
)

// BaGuaFlow 八卦模型
type BaGuaFlow struct {
    *BaseFlowModel

    // 依赖的五行模型
    wuxing *WuXingFlow

    // 八卦状态
    trigrams map[BaGuaTrigram]*TrigramState

    // 八卦场
    fields map[BaGuaTrigram]*core.Field

    // 八卦量子态
    states map[BaGuaTrigram]*core.QuantumState
}

// TrigramState 卦象状态
type TrigramState struct {
    Trigram    BaGuaTrigram
    Energy     float64
    Yao        [3]bool     // 三爻状态
    Nature     Nature      // 阴阳属性
    Element    WuXingPhase // 对应五行
}

// NewBaGuaFlow 创建八卦模型
func NewBaGuaFlow(wuxing *WuXingFlow) *BaGuaFlow {
    base := NewBaseFlowModel(ModelBaGua, 800.0)
    
    bg := &BaGuaFlow{
        BaseFlowModel: base,
        wuxing:        wuxing,
        trigrams:      make(map[BaGuaTrigram]*TrigramState),
        fields:        make(map[BaGuaTrigram]*core.Field),
        states:        make(map[BaGuaTrigram]*core.QuantumState),
    }

    // 初始化八卦
    bg.initializeTrigrams()
    
    // 初始化场
    bg.initializeFields()
    
    // 初始化量子态
    bg.initializeQuantumStates()

    return bg
}

// initializeTrigrams 初始化八卦
func (bg *BaGuaFlow) initializeTrigrams() {
    // 初始化能量
    baseEnergy := bg.capacity / TrigramCount

    // 八卦定义及其属性
    trigramDefs := []struct {
        trigram BaGuaTrigram
        yao     [3]bool
        nature  Nature
        element WuXingPhase
    }{
        {Qian, [3]bool{true, true, true}, NatureYang, Metal},
        {Dui, [3]bool{true, true, false}, NatureYin, Metal},
        {Li, [3]bool{true, false, true}, NatureYang, Fire},
        {Zhen, [3]bool{false, false, true}, NatureYang, Wood},
        {Xun, [3]bool{true, false, false}, NatureYin, Wood},
        {Kan, [3]bool{false, true, false}, NatureYin, Water},
        {Gen, [3]bool{false, true, true}, NatureYang, Earth},
        {Kun, [3]bool{false, false, false}, NatureYin, Earth},
    }

    for _, def := range trigramDefs {
        bg.trigrams[def.trigram] = &TrigramState{
            Trigram: def.trigram,
            Energy:  baseEnergy,
            Yao:     def.yao,
            Nature:  def.nature,
            Element: def.element,
        }
    }
}

// initializeFields 初始化场
func (bg *BaGuaFlow) initializeFields() {
    for trigram := range bg.trigrams {
        bg.fields[trigram] = core.NewField()
        bg.fields[trigram].SetPhase(float64(trigram) * PhaseShift)
    }
}

// initializeQuantumStates 初始化量子态
func (bg *BaGuaFlow) initializeQuantumStates() {
    for trigram := range bg.trigrams {
        bg.states[trigram] = core.NewQuantumState()
        bg.states[trigram].Initialize()
    }
}

// Transform 八卦转换
func (bg *BaGuaFlow) Transform(pattern TransformPattern) error {
    bg.mu.Lock()
    defer bg.mu.Unlock()

    if !bg.running {
        return NewModelError(ErrCodeOperation, "model not running", nil)
    }

    switch pattern {
    case PatternNormal:
        bg.naturalTransform()
    case PatternForward:
        bg.forwardTransform()
    case PatternReverse:
        bg.reverseTransform()
    case PatternBalance:
        bg.balanceTransform()
    case PatternMutate:
        bg.mutateTransform()
    default:
        return NewModelError(ErrCodeOperation, "invalid transform pattern", nil)
    }

    // 同步五行模型
    bg.synchronizeWithWuXing()
    
    // 更新量子态
    bg.updateQuantumStates()
    
    // 更新场
    bg.updateFields()
    
    // 更新状态
    bg.updateModelState()

    return nil
}

// naturalTransform 自然转换
func (bg *BaGuaFlow) naturalTransform() {
    for _, state := range bg.trigrams {
        // 量子态演化
        qstate := bg.states[state.Trigram]
        qstate.Evolve(state.Trigram.String())
        
        // 应用量子涨落
        fluctuation := qstate.GetFluctuation()
        state.Energy *= (1 + fluctuation)
        
        // 爻位变换
        bg.transformYao(state)
    }
}

// forwardTransform 顺序转换
func (bg *BaGuaFlow) forwardTransform() {
    sequence := []BaGuaTrigram{Qian, Dui, Li, Zhen, Xun, Kan, Gen, Kun}
    bg.sequentialTransform(sequence)
}

// reverseTransform 逆序转换
func (bg *BaGuaFlow) reverseTransform() {
    sequence := []BaGuaTrigram{Kun, Gen, Kan, Xun, Zhen, Li, Dui, Qian}
    bg.sequentialTransform(sequence)
}

// sequentialTransform 序列转换
func (bg *BaGuaFlow) sequentialTransform(sequence []BaGuaTrigram) {
    for i := 0; i < len(sequence)-1; i++ {
        current := bg.trigrams[sequence[i]]
        next := bg.trigrams[sequence[i+1]]
        
        transferEnergy := current.Energy * BaseRate
        current.Energy -= transferEnergy
        next.Energy += transferEnergy
    }
}

// balanceTransform 平衡转换
func (bg *BaGuaFlow) balanceTransform() {
    averageEnergy := bg.state.Energy / TrigramCount
    
    for _, state := range bg.trigrams {
        state.Energy = averageEnergy
    }
}

// mutateTransform 变异转换
func (bg *BaGuaFlow) mutateTransform() {
    for _, state := range bg.trigrams {
        // 获取量子涨落
        fluctuation := bg.states[state.Trigram].GetFluctuation()
        
        // 变异爻位
        for i := range state.Yao {
            if math.Abs(fluctuation) > 0.5 {
                state.Yao[i] = !state.Yao[i]
            }
        }
        
        // 调整能量
        state.Energy *= (1 + fluctuation)
    }
}

// transformYao 爻位变换
func (bg *BaGuaFlow) transformYao(state *TrigramState) {
    // 基于量子态概率变换爻位
    for i := range state.Yao {
        if bg.states[state.Trigram].GetProbability() > 0.5 {
            state.Yao[i] = !state.Yao[i]
        }
    }
}

// synchronizeWithWuXing 与五行同步
func (bg *BaGuaFlow) synchronizeWithWuXing() {
    if bg.wuxing == nil {
        return
    }

    // 同步能量
    for _, state := range bg.trigrams {
        elementEnergy := bg.wuxing.GetElementEnergy(state.Element)
        syncEnergy := elementEnergy * SyncRate
        state.Energy = (state.Energy + syncEnergy) / 2
    }
}

// updateQuantumStates 更新量子态
func (bg *BaGuaFlow) updateQuantumStates() {
    totalEnergy := bg.state.Energy
    
    for trigram, state := range bg.trigrams {
        probability := state.Energy / totalEnergy
        bg.states[trigram].SetProbability(probability)
        bg.states[trigram].SetPhase(float64(trigram) * PhaseShift)
        bg.states[trigram].Evolve(trigram.String())
    }
}

// updateFields 更新场
func (bg *BaGuaFlow) updateFields() {
    totalEnergy := bg.state.Energy
    
    for trigram, state := range bg.trigrams {
        field := bg.fields[trigram]
        field.SetStrength(state.Energy / totalEnergy)
        field.SetPhase(bg.states[trigram].GetPhase())
        field.Evolve()
    }
}

// updateModelState 更新模型状态
func (bg *BaGuaFlow) updateModelState() {
    // 更新状态属性
    for trigram, state := range bg.trigrams {
        prefix := trigram.String()
        bg.state.Properties[prefix+"Energy"] = state.Energy
        bg.state.Properties[prefix+"Yao"] = state.Yao
    }
    
    bg.state.Phase = PhaseBaGua
    bg.state.UpdateTime = time.Now()
}

// GetTrigramState 获取卦象状态
func (bg *BaGuaFlow) GetTrigramState(trigram BaGuaTrigram) *TrigramState {
    bg.mu.RLock()
    defer bg.mu.RUnlock()
    
    if state, exists := bg.trigrams[trigram]; exists {
        return state
    }
    return nil
}
