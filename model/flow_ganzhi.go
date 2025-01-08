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
    
    // 关联五行系统
    wuxing    *WuXingFlow
    
    // 天干系统
    gans      map[int]*GanAttributes
    
    // 地支系统
    zhis      map[int]*ZhiAttributes
    
    // 周期控制
    cycle struct {
        current    int
        timestamp time.Time
        duration  time.Duration
    }
    
    // 量子系统
    quantumStates map[int]*core.QuantumState
}

// GanAttributes 天干属性
type GanAttributes struct {
    Element    WuXingPhase     // 五行属性
    Nature     Nature          // 阴阳属性
    Energy     float64         // 能量值
    Position   float64         // 位置角度
    Field      *core.Field    // 关联场
}

// ZhiAttributes 地支属性
type ZhiAttributes struct {
    MainElement   WuXingPhase     // 主气五行
    SubElements   []WuXingPhase   // 余气五行
    Nature        Nature          // 阴阳属性
    Energy        float64         // 能量值
    Position      float64         // 位置角度
    Field         *core.Field    // 关联场
}

// NewGanZhiFlow 创建天干地支流模型
func NewGanZhiFlow(wx *WuXingFlow) *GanZhiFlow {
    gz := &GanZhiFlow{
        BaseFlowModel:  NewBaseFlowModel(ModelGanZhi, BaseEnergy*float64(GanCount+ZhiCount)),
        wuxing:         wx,
        gans:           make(map[int]*GanAttributes),
        zhis:           make(map[int]*ZhiAttributes),
        quantumStates:  make(map[int]*core.QuantumState),
    }
    
    gz.initializeGan()
    gz.initializeZhi()
    gz.initializeQuantumStates()
    
    go gz.runCycle()
    return gz
}

// initializeGan 初始化天干
func (gz *GanZhiFlow) initializeGan() {
    // 天干配置
    configs := []struct {
        element   WuXingPhase
        nature    Nature
        position  float64
    }{
        {Wood, NatureYang, 0},    // 甲
        {Wood, NatureYin, 36},    // 乙
        {Fire, NatureYang, 72},   // 丙
        {Fire, NatureYin, 108},   // 丁
        {Earth, NatureYang, 144}, // 戊
        {Earth, NatureYin, 180},  // 己
        {Metal, NatureYang, 216}, // 庚
        {Metal, NatureYin, 252},  // 辛
        {Water, NatureYang, 288}, // 壬
        {Water, NatureYin, 324},  // 癸
    }
    
    for i, config := range configs {
        gz.gans[i] = &GanAttributes{
            Element:   config.element,
            Nature:    config.nature,
            Energy:    BaseEnergy,
            Position:  config.position,
            Field:     core.NewField(),
        }
    }
}

// initializeZhi 初始化地支
func (gz *GanZhiFlow) initializeZhi() {
    // 地支配置
    configs := []struct {
        main     WuXingPhase
        sub      []WuXingPhase
        nature   Nature
        position float64
    }{
        {Water, []WuXingPhase{Water}, NatureYang, 0},           // 子
        {Earth, []WuXingPhase{Earth, Metal}, NatureYin, 30},    // 丑
        {Wood, []WuXingPhase{Wood, Fire}, NatureYang, 60},      // 寅
        {Wood, []WuXingPhase{Wood}, NatureYin, 90},             // 卯
        {Earth, []WuXingPhase{Earth, Water}, NatureYang, 120},  // 辰
        {Fire, []WuXingPhase{Fire, Earth}, NatureYin, 150},     // 巳
        {Fire, []WuXingPhase{Fire}, NatureYang, 180},           // 午
        {Earth, []WuXingPhase{Earth, Fire}, NatureYin, 210},    // 未
        {Metal, []WuXingPhase{Metal, Earth}, NatureYang, 240},  // 申
        {Metal, []WuXingPhase{Metal}, NatureYin, 270},          // 酉
        {Earth, []WuXingPhase{Earth, Metal}, NatureYang, 300},  // 戌
        {Water, []WuXingPhase{Water, Wood}, NatureYin, 330},    // 亥
    }
    
    for i, config := range configs {
        gz.zhis[i] = &ZhiAttributes{
            MainElement: config.main,
            SubElements: config.sub,
            Nature:     config.nature,
            Energy:     BaseEnergy,
            Position:   config.position,
            Field:      core.NewField(),
        }
    }
}

// initializeQuantumStates 初始化量子态
func (gz *GanZhiFlow) initializeQuantumStates() {
    // 为天干和地支创建量子态
    for i := 0; i < GanCount+ZhiCount; i++ {
        state := core.NewQuantumState()
        angle := float64(i) * 2 * math.Pi / float64(GanCount+ZhiCount)
        state.SetPhase(angle)
        gz.quantumStates[i] = state
    }
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
    
    // 更新周期位置
    gz.cycle.current = (gz.cycle.current + 1) % CycleLength
    
    // 更新量子态
    gz.evolveQuantumStates()
    
    // 更新场效应
    gz.updateFields()
    
    // 更新能量分布
    gz.redistributeEnergy()
    
    // 同步五行系统
    gz.synchronizeWithWuXing()
    
    // 更新状态
    gz.updateState()
}

// evolveQuantumStates 演化量子态
func (gz *GanZhiFlow) evolveQuantumStates() {
    for _, state := range gz.quantumStates {
        state.Evolve(time.Hour)
    }
}

// updateFields 更新场效应
func (gz *GanZhiFlow) updateFields() {
    // 更新天干场
    for i, gan := range gz.gans {
        state := gz.quantumStates[i]
        fieldStrength := state.GetProbability()
        gan.Field.SetStrength(fieldStrength)
    }
    
    // 更新地支场
    for i, zhi := range gz.zhis {
        state := gz.quantumStates[i+GanCount]
        fieldStrength := state.GetProbability()
        zhi.Field.SetStrength(fieldStrength)
    }
}

// redistributeEnergy 重新分配能量
func (gz *GanZhiFlow) redistributeEnergy() {
    totalEnergy := gz.state.Energy
    baseEnergy := totalEnergy / float64(GanCount+ZhiCount)
    
    // 分配天干能量
    for _, gan := range gz.gans {
        gan.Energy = baseEnergy * gz.quantumStates[0].GetProbability()
    }
    
    // 分配地支能量
    for _, zhi := range gz.zhis {
        zhi.Energy = baseEnergy * gz.quantumStates[GanCount].GetProbability()
    }
}

// synchronizeWithWuXing 同步五行系统
func (gz *GanZhiFlow) synchronizeWithWuXing() {
    if gz.wuxing == nil {
        return
    }
    
    // 同步天干五行
    for _, gan := range gz.gans {
        energy := gan.Energy * gan.Field.GetStrength()
        gz.wuxing.AdjustPhaseEnergy(gan.Element, energy)
    }
    
    // 同步地支五行
    for _, zhi := range gz.zhis {
        mainEnergy := zhi.Energy * zhi.Field.GetStrength()
        gz.wuxing.AdjustPhaseEnergy(zhi.MainElement, mainEnergy)
        
        // 处理余气
        subEnergy := mainEnergy * 0.5 / float64(len(zhi.SubElements))
        for _, element := range zhi.SubElements {
            gz.wuxing.AdjustPhaseEnergy(element, subEnergy)
        }
    }
}

// updateState 更新状态
func (gz *GanZhiFlow) updateState() {
    // 更新能量分布
    for i, gan := range gz.gans {
        gz.state.Properties[fmt.Sprintf("gan_%d", i)] = gan.Energy
    }
    for i, zhi := range gz.zhis {
        gz.state.Properties[fmt.Sprintf("zhi_%d", i)] = zhi.Energy
    }
    
    // 更新相位
    gz.state.Phase = PhaseWuXing
    
    // 通知观察者
    gz.notifyObservers()
}
