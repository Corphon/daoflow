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
    Zhen                // 震 ☳
    Xun                 // 巽 ☴
    Kan                 // 坎 ☵
    Li                  // 离 ☲
    Gen                 // 艮 ☶
    Dui                 // 兑 ☱
)

// TrigramAttributes 卦象属性
type TrigramAttributes struct {
    Lines      [TrigramLines]bool // 爻线
    Direction  float64           // 方位角度
    Element    WuXingPhase       // 关联五行
    Nature     Nature            // 阴阳属性
    Energy     float64           // 能量值
    Potential  float64           // 势能
}

// BaGuaFlow 八卦模型
type BaGuaFlow struct {
    *BaseFlowModel
    
    // 卦象系统
    trigrams    map[Trigram]*TrigramAttributes
    
    // 场态
    fieldMatrix [][]float64      // 场分布矩阵
    potential   [][]float64      // 势场分布
    
    // 变化记录
    mutations   map[Trigram][]Mutation
    
    // 外部关联
    wuxing     *WuXingFlow    // 五行关联
    yinyang    *YinYangFlow   // 阴阳关联
}

// Mutation 变卦记录
type Mutation struct {
    From      Trigram
    To        Trigram
    Time      time.Time
    Cause     string
    Energy    float64
}

// NewBaGuaFlow 创建八卦流模型
func NewBaGuaFlow(wx *WuXingFlow, yy *YinYangFlow) *BaGuaFlow {
    bg := &BaGuaFlow{
        BaseFlowModel: NewBaseFlowModel(ModelBaGua, 800.0), // 8卦*100能量
        trigrams:      make(map[Trigram]*TrigramAttributes),
        fieldMatrix:   make([][]float64, 8),
        potential:     make([][]float64, 8),
        mutations:     make(map[Trigram][]Mutation),
        wuxing:       wx,
        yinyang:      yy,
    }
    
    // 初始化场矩阵
    for i := range bg.fieldMatrix {
        bg.fieldMatrix[i] = make([]float64, 8)
        bg.potential[i] = make([]float64, 8)
    }
    
    bg.initializeTrigrams()
    go bg.runFieldCalculation()
    return bg
}

// initializeTrigrams 初始化卦象
func (bg *BaGuaFlow) initializeTrigrams() {
    // 定义卦象配置
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
            Energy:    BasePotential,
            Potential: calculateInitialPotential(config.direction),
        }
    }
}

// calculateInitialPotential 计算初始势能
func calculateInitialPotential(direction float64) float64 {
    // 使用余弦函数创建周期性势能分布
    return BasePotential * (1 + math.Cos(direction*math.Pi/180.0))
}

// runFieldCalculation 运行场计算
func (bg *BaGuaFlow) runFieldCalculation() {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-bg.done:
            return
        case <-ticker.C:
            bg.updateField()
        }
    }
}

// updateField 更新场分布
func (bg *BaGuaFlow) updateField() {
    bg.mu.Lock()
    defer bg.mu.Unlock()
    
    // 计算场分布
    for i, trigram := range bg.trigrams {
        // 计算场强度
        fieldStrength := bg.calculateFieldStrength(trigram)
        
        // 更新场矩阵
        x := int(trigram.Direction / OctantAngle)
        y := int(trigram.Energy / BasePotential)
        if x < 8 && y < 8 {
            bg.fieldMatrix[x][y] = fieldStrength
            
            // 更新势场
            bg.potential[x][y] = bg.calculatePotential(trigram)
        }
        
        // 检查变卦条件
        if bg.shouldMutate(trigram) {
            bg.mutate(i)
        }
    }
}

// calculateFieldStrength 计算场强度
func (bg *BaGuaFlow) calculateFieldStrength(attr *TrigramAttributes) float64 {
    // 使用量子场论的波函数叠加
    psi := complex(attr.Energy/100.0, attr.Potential/BasePotential)
    
    // |ψ|² 给出概率密度
    return math.Pow(cmplx.Abs(psi), 2)
}

// calculatePotential 计算势能
func (bg *BaGuaFlow) calculatePotential(attr *TrigramAttributes) float64 {
    // 考虑五行和阴阳影响
    wuxingContribution := 0.0
    if bg.wuxing != nil {
        if energy, err := bg.wuxing.GetElementStrength(attr.Element); err == nil {
            wuxingContribution = float64(energy) / 100.0
        }
    }
    
    yinyangContribution := 0.0
    if bg.yinyang != nil {
        yin, yang := bg.yinyang.GetRatio()
        if attr.Nature == NatureYin {
            yinyangContribution = yin
        } else {
            yinyangContribution = yang
        }
    }
    
    // 合成势能
    return attr.Potential * (1 + wuxingContribution) * (1 + yinyangContribution)
}

// shouldMutate 判断是否应该变卦
func (bg *BaGuaFlow) shouldMutate(attr *TrigramAttributes) bool {
    // 计算变化趋势
    energyRatio := attr.Energy / (BasePotential * 8)
    potentialRatio := attr.Potential / (BasePotential * 2)
    
    // 使用统计力学的相变模型
    transitionProbability := math.Exp(-(energyRatio * potentialRatio))
    
    return transitionProbability > MutationThreshold
}

// mutate 执行变卦
func (bg *BaGuaFlow) mutate(from Trigram) {
    // 寻找最适合的目标卦象
    var bestTarget Trigram
    maxResonance := 0.0
    
    fromAttr := bg.trigrams[from]
    
    for to, toAttr := range bg.trigrams {
        if to == from {
            continue
        }
        
        // 计算共振度
        resonance := bg.calculateResonance(fromAttr, toAttr)
        if resonance > maxResonance {
            maxResonance = resonance
            bestTarget = to
        }
    }
    
    // 记录变卦
    if maxResonance > ResonanceRate {
        mutation := Mutation{
            From:   from,
            To:     bestTarget,
            Time:   time.Now(),
            Cause:  "field resonance",
            Energy: fromAttr.Energy,
        }
        
        bg.mutations[from] = append(bg.mutations[from], mutation)
        
        // 能量转换
        bg.trigrams[bestTarget].Energy += fromAttr.Energy * maxResonance
        fromAttr.Energy *= (1 - maxResonance)
    }
}

// calculateResonance 计算共振度
func (bg *BaGuaFlow) calculateResonance(from, to *TrigramAttributes) float64 {
    // 计算方向共振
    directionDiff := math.Abs(from.Direction - to.Direction)
    if directionDiff > 180 {
        directionDiff = 360 - directionDiff
    }
    directionResonance := 1 - (directionDiff / 180.0)
    
    // 计算能量共振
    energyRatio := math.Min(from.Energy, to.Energy) / math.Max(from.Energy, to.Energy)
    
    // 计算属性共振
    natureResonance := 0.5
    if from.Nature == to.Nature {
        natureResonance = 1.0
    }
    
    // 综合共振度
    return (directionResonance + energyRatio + natureResonance) / 3
}
