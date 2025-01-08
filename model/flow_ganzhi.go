// model/flow_ganzhi.go

package model

import (
    "math"
    "time"

    "github.com/Corphon/daoflow/core"
)

// GanZhiConstants 天干地支常数
const (
    TianGanCount  = 10          // 天干数量
    DiZhiCount    = 12          // 地支数量
    CycleLength   = 60          // 六十甲子周期
    FlowRate      = 0.1         // 流转率
    MixingRate    = 0.15        // 合化率
)

// GanZhiFlow 天干地支模型
type GanZhiFlow struct {
    *BaseFlowModel

    // 依赖的五行模型
    wuxing *WuXingFlow

    // 天干状态
    tiangan map[int]*TianGan
    // 地支状态
    dizhi map[int]*DiZhi

    // 量子场组件
    tianganFields map[int]*core.Field
    dizhiFields  map[int]*core.Field

    // 量子态组件
    tianganStates map[int]*core.QuantumState
    dizhiStates  map[int]*core.QuantumState
}

// TianGan 天干
type TianGan struct {
    Index     int           // 序号
    Name      string        // 名称
    Element   WuXingPhase   // 五行属性
    Nature    Nature        // 阴阳性质
    Energy    float64       // 能量
    Phase     float64       // 相位
}

// DiZhi 地支
type DiZhi struct {
    Index     int           // 序号
    Name      string        // 名称
    Element   WuXingPhase   // 五行属性
    Nature    Nature        // 阴阳性质
    Energy    float64       // 能量
    Phase     float64       // 相位
    Hidden    []WuXingPhase // 藏干（纳音）
}

// NewGanZhiFlow 创建天干地支模型
func NewGanZhiFlow(wuxing *WuXingFlow) *GanZhiFlow {
    base := NewBaseFlowModel(ModelGanZhi, 1200.0)
    
    gz := &GanZhiFlow{
        BaseFlowModel:  base,
        wuxing:         wuxing,
        tiangan:        make(map[int]*TianGan),
        dizhi:         make(map[int]*DiZhi),
        tianganFields:  make(map[int]*core.Field),
        dizhiFields:   make(map[int]*core.Field),
        tianganStates: make(map[int]*core.QuantumState),
        dizhiStates:  make(map[int]*core.QuantumState),
    }

    // 初始化天干地支
    gz.initializeGanZhi()
    
    // 初始化场
    gz.initializeFields()
    
    // 初始化量子态
    gz.initializeQuantumStates()

    return gz
}

// initializeGanZhi 初始化天干地支
func (gz *GanZhiFlow) initializeGanZhi() {
    // 初始化天干
    tianganDefs := []struct {
        name    string
        element WuXingPhase
        nature  Nature
    }{
        {"甲", Wood,  NatureYang},
        {"乙", Wood,  NatureYin},
        {"丙", Fire,  NatureYang},
        {"丁", Fire,  NatureYin},
        {"戊", Earth, NatureYang},
        {"己", Earth, NatureYin},
        {"庚", Metal, NatureYang},
        {"辛", Metal, NatureYin},
        {"壬", Water, NatureYang},
        {"癸", Water, NatureYin},
    }

    // 初始化地支
    dizhiDefs := []struct {
        name    string
        element WuXingPhase
        nature  Nature
        hidden  []WuXingPhase
    }{
        {"子", Water,  NatureYang, []WuXingPhase{Water}},
        {"丑", Earth,  NatureYin,  []WuXingPhase{Earth, Metal, Water}},
        {"寅", Wood,   NatureYang, []WuXingPhase{Wood, Fire, Earth}},
        {"卯", Wood,   NatureYin,  []WuXingPhase{Wood}},
        {"辰", Earth,  NatureYang, []WuXingPhase{Earth, Water, Wood}},
        {"巳", Fire,   NatureYin,  []WuXingPhase{Fire, Earth, Metal}},
        {"午", Fire,   NatureYang, []WuXingPhase{Fire, Earth}},
        {"未", Earth,  NatureYin,  []WuXingPhase{Earth, Fire, Wood}},
        {"申", Metal,  NatureYang, []WuXingPhase{Metal, Water, Earth}},
        {"酉", Metal,  NatureYin,  []WuXingPhase{Metal}},
        {"戌", Earth,  NatureYang, []WuXingPhase{Earth, Fire, Metal}},
        {"亥", Water,  NatureYin,  []WuXingPhase{Water, Wood}},
    }

    // 初始化天干能量
    tianganEnergy := gz.capacity * 0.4 / float64(TianGanCount)
    for i, def := range tianganDefs {
        gz.tiangan[i] = &TianGan{
            Index:   i,
            Name:    def.name,
            Element: def.element,
            Nature:  def.nature,
            Energy:  tianganEnergy,
            Phase:   2 * math.Pi * float64(i) / float64(TianGanCount),
        }
    }

    // 初始化地支能量
    dizhiEnergy := gz.capacity * 0.6 / float64(DiZhiCount)
    for i, def := range dizhiDefs {
        gz.dizhi[i] = &DiZhi{
            Index:   i,
            Name:    def.name,
            Element: def.element,
            Nature:  def.nature,
            Energy:  dizhiEnergy,
            Phase:   2 * math.Pi * float64(i) / float64(DiZhiCount),
            Hidden:  def.hidden,
        }
    }
}

// initializeFields 初始化场
func (gz *GanZhiFlow) initializeFields() {
    // 初始化天干场
    for i := range gz.tiangan {
        gz.tianganFields[i] = core.NewField()
        gz.tianganFields[i].SetPhase(gz.tiangan[i].Phase)
    }

    // 初始化地支场
    for i := range gz.dizhi {
        gz.dizhiFields[i] = core.NewField()
        gz.dizhiFields[i].SetPhase(gz.dizhi[i].Phase)
    }
}

// initializeQuantumStates 初始化量子态
func (gz *GanZhiFlow) initializeQuantumStates() {
    // 初始化天干量子态
    for i := range gz.tiangan {
        gz.tianganStates[i] = core.NewQuantumState()
        gz.tianganStates[i].Initialize()
    }

    // 初始化地支量子态
    for i := range gz.dizhi {
        gz.dizhiStates[i] = core.NewQuantumState()
        gz.dizhiStates[i].Initialize()
    }
}

// Transform 天干地支转换
func (gz *GanZhiFlow) Transform(pattern TransformPattern) error {
    gz.mu.Lock()
    defer gz.mu.Unlock()

    if !gz.running {
        return NewModelError(ErrCodeOperation, "model not running", nil)
    }

    switch pattern {
    case PatternNormal:
        gz.naturalTransform()
    case PatternForward:
        gz.cycleTransform()
    case PatternReverse:
        gz.reverseCycleTransform()
    case PatternBalance:
        gz.balanceTransform()
    case PatternMutate:
        gz.mutateTransform()
    default:
        return NewModelError(ErrCodeOperation, "invalid transform pattern", nil)
    }

    // 同步五行能量
    gz.synchronizeWithWuXing()
    
    // 更新量子态
    gz.updateQuantumStates()
    
    // 更新场
    gz.updateFields()
    
    // 更新状态
    gz.updateModelState()

    return nil
}

// naturalTransform 自然转换
func (gz *GanZhiFlow) naturalTransform() {
    // 天干自然转换
    for i, gan := range gz.tiangan {
        state := gz.tianganStates[i]
        state.Evolve(gan.Name)
        fluctuation := state.GetFluctuation()
        gan.Energy *= (1 + fluctuation)
        gan.Phase = math.Mod(gan.Phase + FlowRate, 2*math.Pi)
    }

    // 地支自然转换
    for i, zhi := range gz.dizhi {
        state := gz.dizhiStates[i]
        state.Evolve(zhi.Name)
        fluctuation := state.GetFluctuation()
        zhi.Energy *= (1 + fluctuation)
        zhi.Phase = math.Mod(zhi.Phase + FlowRate*0.8, 2*math.Pi)
    }
}

// cycleTransform 循环转换
func (gz *GanZhiFlow) cycleTransform() {
    // 天干循环
    for i := 0; i < TianGanCount-1; i++ {
        energy := gz.tiangan[i].Energy * FlowRate
        gz.tiangan[i].Energy -= energy
        gz.tiangan[i+1].Energy += energy
    }

    // 地支循环
    for i := 0; i < DiZhiCount-1; i++ {
        energy := gz.dizhi[i].Energy * FlowRate
        gz.dizhi[i].Energy -= energy
        gz.dizhi[i+1].Energy += energy
    }
}

// reverseCycleTransform 逆循环转换
func (gz *GanZhiFlow) reverseCycleTransform() {
    // 天干逆循环
    for i := TianGanCount - 1; i > 0; i-- {
        energy := gz.tiangan[i].Energy * FlowRate
        gz.tiangan[i].Energy -= energy
        gz.tiangan[i-1].Energy += energy
    }

    // 地支逆循环
    for i := DiZhiCount - 1; i > 0; i-- {
        energy := gz.dizhi[i].Energy * FlowRate
        gz.dizhi[i].Energy -= energy
        gz.dizhi[i-1].Energy += energy
    }
}

// balanceTransform 平衡转换
func (gz *GanZhiFlow) balanceTransform() {
    // 平衡天干能量
    tianganEnergy := gz.capacity * 0.4 / float64(TianGanCount)
    for _, gan := range gz.tiangan {
        gan.Energy = tianganEnergy
    }

    // 平衡地支能量
    dizhiEnergy := gz.capacity * 0.6 / float64(DiZhiCount)
    for _, zhi := range gz.dizhi {
        zhi.Energy = dizhiEnergy
    }

    // 重置动量
    for _, gan := range gz.tiangan {
        gan.Momentum = Vector3D{}
    }
    for _, zhi := range gz.dizhi {
        zhi.Momentum = Vector3D{}
    }
}

// mutateTransform 变异转换
func (gz *GanZhiFlow) mutateTransform() {
    // 使用量子涨落进行随机变异
    for i, gan := range gz.tiangan {
        fluctuation := gz.tianganStates[i].GetFluctuation()
        gan.Energy *= (1 + fluctuation)
    }
    
    for i, zhi := range gz.dizhi {
        fluctuation := gz.dizhiStates[i].GetFluctuation()
        zhi.Energy *= (1 + fluctuation)
    }
}

// synchronizeWithWuXing 与五行同步
func (gz *GanZhiFlow) synchronizeWithWuXing() {
    if gz.wuxing == nil {
        return
    }

    // 同步天干与五行
    for _, gan := range gz.tiangan {
        wuxingEnergy := gz.wuxing.GetElementEnergy(gan.Element)
        resonance := math.Min(gan.Energy, wuxingEnergy) * ResonanceRate
        gan.Energy += resonance
    }

    // 同步地支与五行
    for _, zhi := range gz.dizhi {
        wuxingEnergy := gz.wuxing.GetElementEnergy(zhi.Element)
        resonance := math.Min(zhi.Energy, wuxingEnergy) * ResonanceRate
        zhi.Energy += resonance
    }
}

// updateQuantumStates 更新量子态
func (gz *GanZhiFlow) updateQuantumStates() {
    // 更新天干量子态
    totalTianganEnergy := gz.getTotalTianganEnergy()
    for i, gan := range gz.tiangan {
        probability := gan.Energy / totalTianganEnergy
        gz.tianganStates[i].SetProbability(probability)
        gz.tianganStates[i].Evolve(gan.Name)
    }
    
    // 更新地支量子态
    totalDizhiEnergy := gz.getTotalDizhiEnergy()
    for i, zhi := range gz.dizhi {
        probability := zhi.Energy / totalDizhiEnergy
        gz.dizhiStates[i].SetProbability(probability)
        gz.dizhiStates[i].Evolve(zhi.Name)
    }
}

// updateFields 更新场
func (gz *GanZhiFlow) updateFields() {
    // 更新天干场
    totalTianganEnergy := gz.getTotalTianganEnergy()
    for i, gan := range gz.tiangan {
        field := gz.tianganFields[i]
        field.SetStrength(gan.Energy / totalTianganEnergy)
        field.SetPhase(gz.tianganStates[i].GetPhase())
        field.Evolve()
    }
    
    // 更新地支场
    totalDizhiEnergy := gz.getTotalDizhiEnergy()
    for i, zhi := range gz.dizhi {
        field := gz.dizhiFields[i]
        field.SetStrength(zhi.Energy / totalDizhiEnergy)
        field.SetPhase(gz.dizhiStates[i].GetPhase())
        field.Evolve()
    }
}

// updateModelState 更新模型状态
func (gz *GanZhiFlow) updateModelState() {
    // 更新状态属性
    for i, gan := range gz.tiangan {
        gz.state.Properties["tiangan_"+gan.Name] = gan.Energy
    }
    for i, zhi := range gz.dizhi {
        gz.state.Properties["dizhi_"+zhi.Name] = zhi.Energy
    }
    
    gz.state.Phase = PhaseGanZhi
    gz.state.UpdateTime = time.Now()
}

// GetTianGanEnergy 获取天干能量
func (gz *GanZhiFlow) GetTianGanEnergy(index int) float64 {
    gz.mu.RLock()
    defer gz.mu.RUnlock()
    
    if gan, exists := gz.tiangan[index]; exists {
        return gan.Energy
    }
    return 0
}

// GetDiZhiEnergy 获取地支能量
func (gz *GanZhiFlow) GetDiZhiEnergy(index int) float64 {
    gz.mu.RLock()
    defer gz.mu.RUnlock()
    
    if zhi, exists := gz.dizhi[index]; exists {
        return zhi.Energy
    }
    return 0
}

// getTotalTianganEnergy 获取总天干能量
func (gz *GanZhiFlow) getTotalTianganEnergy() float64 {
    var total float64
    for _, gan := range gz.tiangan {
        total += gan.Energy
    }
    return total
}

// getTotalDizhiEnergy 获取总地支能量
func (gz *GanZhiFlow) getTotalDizhiEnergy() float64 {
    var total float64
    for _, zhi := range gz.dizhi {
        total += zhi.Energy
    }
    return total
}

// Close 关闭干支模型
func (gz *GanZhiFlow) Close() error {
    if err := gz.BaseFlowModel.Close(); err != nil {
        return err
    }
    
    // 清理资源
    gz.tianganFields = nil
    gz.dizhiFields = nil
    gz.tianganStates = nil
    gz.dizhiStates = nil
    
    return nil
}
