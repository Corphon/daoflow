// model/flow_integrate.go

package model

import (
    "math"
    "sync"
    "time"

    "github.com/Corphon/daoflow/core"
)

// IntegrateConstants 集成常数
const (
    SystemCapacity    = 2000.0      // 系统总容量
    IntegrationCycle  = time.Minute // 集成周期
    BalanceThreshold  = 0.15        // 平衡阈值
    ResonanceMinimum = 0.3         // 最小共振阈值
    SystemLayers     = 4           // 系统层数(阴阳、五行、八卦、干支)
)

// FlowSystem 流系统状态
type FlowSystem struct {
    Energy     float64   // 系统能量
    Entropy    float64   // 系统熵
    Harmony    float64   // 和谐度
    Balance    float64   // 平衡度
    Coherence  float64   // 相干度
    Phase      float64   // 系统相位
}

// IntegrateFlow 集成模型
type IntegrateFlow struct {
    *BaseFlowModel
    
    // 子系统
    yinyang  *YinYangFlow
    wuxing   *WuXingFlow
    bagua    *BaGuaFlow
    ganzhi   *GanZhiFlow
    
    // 系统状态
    system   FlowSystem
    
    // 核心组件
    coreField    *core.UnifiedField    // 统一场
    coreQuantum  *core.QuantumSystem   // 量子系统
    
    // 状态追踪
    stateHistory []SystemState
    transitions  chan StateTransition
}

// SystemState 系统状态
type SystemState struct {
    Timestamp  time.Time
    System     FlowSystem
    YinYang    ModelState
    WuXing     ModelState
    BaGua      ModelState
    GanZhi     ModelState
}

// NewIntegrateFlow 创建集成流模型
func NewIntegrateFlow() *IntegrateFlow {
    // 创建子系统
    yy := NewYinYangFlow()
    wx := NewWuXingFlow()
    bg := NewBaGuaFlow(wx)
    gz := NewGanZhiFlow(wx)
    
    iflow := &IntegrateFlow{
        BaseFlowModel: NewBaseFlowModel(ModelIntegrate, SystemCapacity),
        yinyang:      yy,
        wuxing:       wx,
        bagua:        bg,
        ganzhi:       gz,
        coreField:    core.NewUnifiedField(SystemLayers),
        coreQuantum:  core.NewQuantumSystem(SystemLayers),
        transitions:  make(chan StateTransition, 100),
    }
    
    // 初始化系统状态
    iflow.system = FlowSystem{
        Energy:    SystemCapacity * 0.5,
        Entropy:   0,
        Harmony:   1.0,
        Balance:   1.0,
        Coherence: 1.0,
        Phase:     0,
    }
    
    go iflow.runIntegration()
    return iflow
}

// runIntegration 运行系统集成
func (if *IntegrateFlow) runIntegration() {
    ticker := time.NewTicker(IntegrationCycle)
    defer ticker.Stop()

    for {
        select {
        case <-if.done:
            return
        case <-ticker.C:
            if.integrate()
        case transition := <-if.transitions:
            if.handleTransition(transition)
        }
    }
}

// integrate 执行系统集成
func (if *IntegrateFlow) integrate() {
    if.mu.Lock()
    defer if.mu.Unlock()

    // 收集子系统状态
    states := if.collectSystemStates()
    
    // 计算量子态演化
    if.evolveQuantumStates(states)
    
    // 更新统一场
    if.updateUnifiedField(states)
    
    // 计算系统特性
    if.calculateSystemProperties(states)
    
    // 进行能量再分配
    if.redistributeEnergy()
    
    // 记录状态
    if.recordState(states)
}

// collectSystemStates 收集系统状态
func (if *IntegrateFlow) collectSystemStates() SystemState {
    return SystemState{
        Timestamp: time.Now(),
        System:    if.system,
        YinYang:   if.yinyang.GetState(),
        WuXing:    if.wuxing.GetState(),
        BaGua:     if.bagua.GetState(),
        GanZhi:    if.ganzhi.GetState(),
    }
}

// evolveQuantumStates 演化量子态
func (if *IntegrateFlow) evolveQuantumStates(states SystemState) {
    // 构建量子态向量
    stateVector := []float64{
        states.YinYang.Energy / if.system.Energy,
        states.WuXing.Energy / if.system.Energy,
        states.BaGua.Energy / if.system.Energy,
        states.GanZhi.Energy / if.system.Energy,
    }
    
    // 应用量子演化
    if.coreQuantum.Evolve(stateVector, IntegrationCycle)
}

// updateUnifiedField 更新统一场
func (if *IntegrateFlow) updateUnifiedField(states SystemState) {
    // 更新场强度
    fieldStrengths := []float64{
        states.YinYang.Energy / SystemCapacity,
        states.WuXing.Energy / SystemCapacity,
        states.BaGua.Energy / SystemCapacity,
        states.GanZhi.Energy / SystemCapacity,
    }
    
    if.coreField.UpdateStrengths(fieldStrengths)
    
    // 计算场相互作用
    if.coreField.CalculateInteractions()
}

// calculateSystemProperties 计算系统特性
func (if *IntegrateFlow) calculateSystemProperties(states SystemState) {
    // 计算总能量
    if.system.Energy = states.YinYang.Energy +
                      states.WuXing.Energy +
                      states.BaGua.Energy +
                      states.GanZhi.Energy
    
    // 计算系统熵
    if.system.Entropy = if.calculateEntropy(states)
    
    // 计算和谐度
    if.system.Harmony = if.calculateHarmony(states)
    
    // 计算平衡度
    if.system.Balance = if.calculateBalance(states)
    
    // 计算相干度
    if.system.Coherence = if.coreQuantum.GetCoherence()
    
    // 更新系统相位
    if.system.Phase = if.coreQuantum.GetGlobalPhase()
}

// calculateEntropy 计算系统熵
func (if *IntegrateFlow) calculateEntropy(states SystemState) float64 {
    // 使用统计熵公式
    totalEnergy := if.system.Energy
    if totalEnergy == 0 {
        return 0
    }
    
    energies := []float64{
        states.YinYang.Energy,
        states.WuXing.Energy,
        states.BaGua.Energy,
        states.GanZhi.Energy,
    }
    
    var entropy float64
    for _, e := range energies {
        if e > 0 {
            p := e / totalEnergy
            entropy -= p * math.Log(p)
        }
    }
    
    return entropy
}

// calculateHarmony 计算和谐度
func (if *IntegrateFlow) calculateHarmony(states SystemState) float64 {
    // 基于场的相互作用计算和谐度
    return if.coreField.GetHarmony()
}

// calculateBalance 计算平衡度
func (if *IntegrateFlow) calculateBalance(states SystemState) float64 {
    // 计算能量分布的均匀程度
    mean := if.system.Energy / float64(SystemLayers)
    var variance float64
    
    energies := []float64{
        states.YinYang.Energy,
        states.WuXing.Energy,
        states.BaGua.Energy,
        states.GanZhi.Energy,
    }
    
    for _, e := range energies {
        diff := e - mean
        variance += diff * diff
    }
    
    variance /= float64(SystemLayers)
    return 1 / (1 + math.Sqrt(variance))
}

// redistributeEnergy 重新分配能量
func (if *IntegrateFlow) redistributeEnergy() {
    if if.system.Balance < BalanceThreshold {
        // 获取量子态概率分布
        probs := if.coreQuantum.GetProbabilities()
        
        // 按概率分配能量
        totalEnergy := if.system.Energy
        if.yinyang.AdjustEnergy(totalEnergy * probs[0])
        if.wuxing.AdjustEnergy(totalEnergy * probs[1])
        if.bagua.AdjustEnergy(totalEnergy * probs[2])
        if.ganzhi.AdjustEnergy(totalEnergy * probs[3])
    }
}

// recordState 记录状态
func (if *IntegrateFlow) recordState(state SystemState) {
    if.stateHistory = append(if.stateHistory, state)
    if len(if.stateHistory) > 1000 {
        if.stateHistory = if.stateHistory[1:]
    }
}

// GetSystemState 获取系统状态
func (if *IntegrateFlow) GetSystemState() FlowSystem {
    if.mu.RLock()
    defer if.mu.RUnlock()
    return if.system
}
