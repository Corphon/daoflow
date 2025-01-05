// model/flow_ganzhi.go

package model

import (
    "math"
    "sync"
    "time"

    "github.com/Corphon/daoflow/core"
)

// GanZhiConstants 天干地支常数
const (
    GanCount     = 10    // 天干数量
    ZhiCount     = 12    // 地支数量
    CycleLength  = 60    // 六十甲子周期
    PhaseAngle   = 30.0  // 相位角度
    BaseEnergy   = 100.0 // 基础能量值
)

// GanZhiFlow 天干地支模型
type GanZhiFlow struct {
    *BaseFlowModel
    
    // 天干属性
    gans map[Gan]*GanAttributes
    
    // 地支属性
    zhis map[Zhi]*ZhiAttributes
    
    // 组合状态
    combinations map[int]*Combination
    
    // 关联系统
    wuxing *WuXingFlow
    
    // 周期控制
    cycle struct {
        current    int
        timestamp time.Time
        duration  time.Duration
    }
}

// GanAttributes 天干属性
type GanAttributes struct {
    Element   WuXingPhase // 五行属性
    Nature    Nature      // 阴阳属性
    Energy    float64     // 能量值
    Position  float64     // 位置角度
    Velocity  float64     // 角速度
}

// ZhiAttributes 地支属性
type ZhiAttributes struct {
    MainElement  WuXingPhase   // 主气五行
    SubElements  []WuXingPhase // 余气五行
    Nature       Nature        // 阴阳属性
    Energy       float64       // 能量值
    Position     float64       // 位置角度
    Momentum     float64       // 角动量
}

// Combination 干支组合
type Combination struct {
    Gan         Gan
    Zhi         Zhi
    Energy      float64
    Resonance   float64
    Phase       float64
    LastUpdate  time.Time
}

// NewGanZhiFlow 创建天干地支流模型
func NewGanZhiFlow(wx *WuXingFlow) *GanZhiFlow {
    gz := &GanZhiFlow{
        BaseFlowModel: NewBaseFlowModel(ModelGanZhi, BaseEnergy*float64(GanCount+ZhiCount)),
        gans:         make(map[Gan]*GanAttributes),
        zhis:         make(map[Zhi]*ZhiAttributes),
        combinations: make(map[int]*Combination),
        wuxing:      wx,
    }
    
    gz.initializeGan()
    gz.initializeZhi()
    gz.initializeCombinations()
    
    go gz.runCycle()
    return gz
}

// initializeGan 初始化天干
func (gz *GanZhiFlow) initializeGan() {
    configs := map[Gan]struct {
        element   WuXingPhase
        nature    Nature
        position  float64
    }{
        GanJia:  {Wood, NatureYang, 0},
        GanYi:   {Wood, NatureYin, 36},
        GanBing: {Fire, NatureYang, 72},
        GanDing: {Fire, NatureYin, 108},
        GanWu:   {Earth, NatureYang, 144},
        GanJi:   {Earth, NatureYin, 180},
        GanGeng: {Metal, NatureYang, 216},
        GanXin:  {Metal, NatureYin, 252},
        GanRen:  {Water, NatureYang, 288},
        GanGui:  {Water, NatureYin, 324},
    }
    
    for gan, config := range configs {
        gz.gans[gan] = &GanAttributes{
            Element:   config.element,
            Nature:    config.nature,
            Energy:    BaseEnergy,
            Position:  config.position,
            Velocity:  2 * math.Pi / float64(CycleLength * time.Hour),
        }
    }
}

// initializeZhi 初始化地支
func (gz *GanZhiFlow) initializeZhi() {
    configs := map[Zhi]struct {
        main     WuXingPhase
        sub      []WuXingPhase
        nature   Nature
        position float64
    }{
        ZhiZi:   {Water, []WuXingPhase{Water}, NatureYang, 0},
        ZhiChou: {Earth, []WuXingPhase{Earth, Metal, Water}, NatureYin, 30},
        ZhiYin:  {Wood, []WuXingPhase{Wood, Fire, Earth}, NatureYang, 60},
        ZhiMao:  {Wood, []WuXingPhase{Wood}, NatureYin, 90},
        ZhiChen: {Earth, []WuXingPhase{Earth, Water, Wood}, NatureYang, 120},
        ZhiSi:   {Fire, []WuXingPhase{Fire, Earth, Metal}, NatureYin, 150},
        ZhiWu:   {Fire, []WuXingPhase{Fire}, NatureYang, 180},
        ZhiWei:  {Earth, []WuXingPhase{Earth, Fire, Wood}, NatureYin, 210},
        ZhiShen: {Metal, []WuXingPhase{Metal, Water, Earth}, NatureYang, 240},
        ZhiYou:  {Metal, []WuXingPhase{Metal}, NatureYin, 270},
        ZhiXu:   {Earth, []WuXingPhase{Earth, Fire, Metal}, NatureYang, 300},
        ZhiHai:  {Water, []WuXingPhase{Water, Wood}, NatureYin, 330},
    }
    
    for zhi, config := range configs {
        gz.zhis[zhi] = &ZhiAttributes{
            MainElement: config.main,
            SubElements: config.sub,
            Nature:     config.nature,
            Energy:     BaseEnergy,
            Position:   config.position,
            Momentum:   BaseEnergy * config.position,
        }
    }
}

// initializeCombinations 初始化干支组合
func (gz *GanZhiFlow) initializeCombinations() {
    for i := 0; i < CycleLength; i++ {
        gan := Gan(i % GanCount)
        zhi := Zhi(i % ZhiCount)
        
        gz.combinations[i] = &Combination{
            Gan:        gan,
            Zhi:        zhi,
            Energy:     calculateCombinationEnergy(gz.gans[gan], gz.zhis[zhi]),
            Resonance:  calculateResonance(gz.gans[gan], gz.zhis[zhi]),
            Phase:      0,
            LastUpdate: time.Now(),
        }
    }
}

// calculateCombinationEnergy 计算组合能量
func calculateCombinationEnergy(gan *GanAttributes, zhi *ZhiAttributes) float64 {
    // 使用量子力学的能级公式
    // E = E₀ + ℏω(n + ½)
    baseEnergy := (gan.Energy + zhi.Energy) / 2
    angularFrequency := 2 * math.Pi / CycleLength
    quantumNumber := math.Abs(gan.Position - zhi.Position) / PhaseAngle
    
    return baseEnergy * (1 + angularFrequency*(quantumNumber + 0.5))
}

// calculateResonance 计算共振强度
func calculateResonance(gan *GanAttributes, zhi *ZhiAttributes) float64 {
    // 计算五行相生相克
    elementFactor := calculateElementRelation(gan.Element, zhi.MainElement)
    
    // 计算阴阳调和
    natureFactor := 1.0
    if gan.Nature == zhi.Nature {
        natureFactor = 1.5
    }
    
    // 计算位置谐振
    positionDiff := math.Abs(gan.Position - zhi.Position)
    if positionDiff > 180 {
        positionDiff = 360 - positionDiff
    }
    positionFactor := 1 - (positionDiff / 360)
    
    return (elementFactor + natureFactor + positionFactor) / 3
}

// runCycle 运行干支周期
func (gz *GanZhiFlow) runCycle() {
    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-gz.done:
            return
        case <-ticker.C:
            gz.updateCycle()
        }
    }
}

// updateCycle 更新干支周期
func (gz *GanZhiFlow) updateCycle() {
    gz.mu.Lock()
    defer gz.mu.Unlock()
    
    now := time.Now()
    
    // 更新周期位置
    gz.cycle.current = (gz.cycle.current + 1) % CycleLength
    
    // 更新组合状态
    combination := gz.combinations[gz.cycle.current]
    gan := gz.gans[combination.Gan]
    zhi := gz.zhis[combination.Zhi]
    
    // 更新能量和相位
    combination.Energy = calculateCombinationEnergy(gan, zhi)
    combination.Resonance = calculateResonance(gan, zhi)
    combination.Phase += 2 * math.Pi / CycleLength
    combination.LastUpdate = now
    
    // 影响五行系统
    if gz.wuxing != nil {
        gz.updateWuXing(combination)
    }
}

// updateWuXing 更新五行系统
func (gz *GanZhiFlow) updateWuXing(comb *Combination) {
    gan := gz.gans[comb.Gan]
    zhi := gz.zhis[comb.Zhi]
    
    // 调整天干五行
    gz.wuxing.AdjustElement(gan.Element, int8(comb.Energy/10))
    
    // 调整地支五行
    gz.wuxing.AdjustElement(zhi.MainElement, int8(comb.Energy/10))
    for _, subElement := range zhi.SubElements {
        gz.wuxing.AdjustElement(subElement, int8(comb.Energy/20))
    }
}

// GetCurrentCombination 获取当前干支组合
func (gz *GanZhiFlow) GetCurrentCombination() *Combination {
    gz.mu.RLock()
    defer gz.mu.RUnlock()
    return gz.combinations[gz.cycle.current]
}
