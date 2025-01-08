// model/flow_integrate.go

package model

import (
    "sync"
    "time"

    "github.com/Corphon/daoflow/core"
)

// IntegrateConstants 集成常数
const (
    IntegrateSyncRate = 0.2     // 同步率
    IntegrateBalance = 0.25     // 平衡系数
    ResonanceThreshold = 0.8    // 共振阈值
)

// IntegrateFlow 集成流模型
type IntegrateFlow struct {
    *BaseFlowModel
    mu sync.RWMutex

    // 子模型
    yinyang *YinYangFlow
    wuxing  *WuXingFlow
    bagua   *BaGuaFlow
    ganzhi  *GanZhiFlow

    // 统一场
    unifiedField *core.Field

    // 量子纠缠态
    entangledState *core.QuantumState

    // 系统状态
    systemState SystemState
}

// SystemState 系统状态
type SystemState struct {
    Energy     float64
    Entropy    float64
    Harmony    float64
    Balance    float64
    Phase      Phase
    Timestamp  time.Time
}

// NewIntegrateFlow 创建集成流模型
func NewIntegrateFlow() *IntegrateFlow {
    base := NewBaseFlowModel(ModelIntegrate, 2000.0)
    
    // 创建子模型
    yinyang := NewYinYangFlow()
    wuxing := NewWuXingFlow()
    bagua := NewBaGuaFlow(wuxing)
    ganzhi := NewGanZhiFlow(wuxing)

    return &IntegrateFlow{
        BaseFlowModel:  base,
        yinyang:       yinyang,
        wuxing:        wuxing,
        bagua:         bagua,
        ganzhi:        ganzhi,
        unifiedField:  core.NewField(),
        entangledState: core.NewQuantumState(),
        systemState:   SystemState{
            Energy:    0,
            Entropy:   0,
            Harmony:   1,
            Balance:   1,
            Phase:     PhaseNone,
            Timestamp: time.Now(),
        },
    }
}

// Start 启动集成模型
func (if *IntegrateFlow) Start() error {
    if.mu.Lock()
    defer if.mu.Unlock()

    if if.running {
        return NewModelError(ErrCodeOperation, "model already started", nil)
    }

    // 启动子模型
    if err := if.yinyang.Start(); err != nil {
        return err
    }
    if err := if.wuxing.Start(); err != nil {
        return err
    }
    if err := if.bagua.Start(); err != nil {
        return err
    }
    if err := if.ganzhi.Start(); err != nil {
        return err
    }

    // 初始化统一场
    if.unifiedField.Initialize()
    
    // 初始化量子纠缠态
    if.entangledState.Initialize()
    
    if.running = true
    return nil
}

// Stop 停止集成模型
func (if *IntegrateFlow) Stop() error {
    if.mu.Lock()
    defer if.mu.Unlock()

    if !if.running {
        return NewModelError(ErrCodeOperation, "model not running", nil)
    }

    // 停止子模型
    if err := if.yinyang.Stop(); err != nil {
        return err
    }
    if err := if.wuxing.Stop(); err != nil {
        return err
    }
    if err := if.bagua.Stop(); err != nil {
        return err
    }
    if err := if.ganzhi.Stop(); err != nil {
        return err
    }

    if.running = false
    return nil
}

// Transform 集成转换
func (if *IntegrateFlow) Transform(pattern TransformPattern) error {
    if.mu.Lock()
    defer if.mu.Unlock()

    if !if.running {
        return NewModelError(ErrCodeOperation, "model not running", nil)
    }

    // 转换子模型
    if err := if.yinyang.Transform(pattern); err != nil {
        return err
    }
    if err := if.wuxing.Transform(pattern); err != nil {
        return err
    }
    if err := if.bagua.Transform(pattern); err != nil {
        return err
    }
    if err := if.ganzhi.Transform(pattern); err != nil {
        return err
    }

    // 同步子模型
    if.synchronizeModels()
    
    // 更新量子态
    if.updateQuantumStates()
    
    // 更新场
    if.updateFields()
    
    // 更新系统状态
    if.updateSystemState()

    return nil
}

// synchronizeModels 同步子模型
func (if *IntegrateFlow) synchronizeModels() {
    // 阴阳与五行同步
    yinYangState := if.yinyang.GetState()
    wuxingState := if.wuxing.GetState()
    
    syncEnergy := math.Min(yinYangState.Energy, wuxingState.Energy) * IntegrateSyncRate
    if.yinyang.AdjustEnergy(syncEnergy)
    if.wuxing.AdjustEnergy(syncEnergy)

    // 八卦与干支同步
    baguaState := if.bagua.GetState()
    ganzhiState := if.ganzhi.GetState()
    
    syncEnergy = math.Min(baguaState.Energy, ganzhiState.Energy) * IntegrateSyncRate
    if.bagua.AdjustEnergy(syncEnergy)
    if.ganzhi.AdjustEnergy(syncEnergy)
}

// updateQuantumStates 更新量子态
func (if *IntegrateFlow) updateQuantumStates() {
    // 更新纠缠态
    yinYangProb := if.yinyang.GetState().Energy / if.capacity
    wuxingProb := if.wuxing.GetState().Energy / if.capacity
    baguaProb := if.bagua.GetState().Energy / if.capacity
    ganzhiProb := if.ganzhi.GetState().Energy / if.capacity
    
    avgProb := (yinYangProb + wuxingProb + baguaProb + ganzhiProb) / 4
    if.entangledState.SetProbability(avgProb)
    if.entangledState.Evolve("integrate")
}

// updateFields 更新场
func (if *IntegrateFlow) updateFields() {
    // 更新统一场
    totalStrength := (if.yinyang.field.GetStrength() +
                     if.wuxing.field.GetStrength() +
                     if.bagua.field.GetStrength() +
                     if.ganzhi.field.GetStrength()) / 4
                     
    if.unifiedField.SetStrength(totalStrength)
    if.unifiedField.SetPhase(if.entangledState.GetPhase())
    if.unifiedField.Evolve()
}

// updateSystemState 更新系统状态
func (if *IntegrateFlow) updateSystemState() {
    // 计算总能量
    if.systemState.Energy = if.yinyang.GetState().Energy +
                           if.wuxing.GetState().Energy +
                           if.bagua.GetState().Energy +
                           if.ganzhi.GetState().Energy

    // 计算熵
    if.systemState.Entropy = if.calculateSystemEntropy()

    // 计算和谐度
    if.systemState.Harmony = if.calculateSystemHarmony()

    // 计算平衡度
    if.systemState.Balance = if.calculateSystemBalance()

    // 更新时间戳
    if.systemState.Timestamp = time.Now()

    // 更新模型状态
    if.state.Energy = if.systemState.Energy
    if.state.Properties["entropy"] = if.systemState.Entropy
    if.state.Properties["harmony"] = if.systemState.Harmony
    if.state.Properties["balance"] = if.systemState.Balance
    if.state.UpdateTime = if.systemState.Timestamp
}

// calculateSystemEntropy 计算系统熵
func (if *IntegrateFlow) calculateSystemEntropy() float64 {
    if if.systemState.Energy <= 0 {
        return 0
    }

    // 使用量子态计算熵
    return -if.entangledState.GetProbability() * math.Log(if.entangledState.GetProbability())
}

// calculateSystemHarmony 计算系统和谐度
func (if *IntegrateFlow) calculateSystemHarmony() float64 {
    // 基于场强度计算和谐度
    fieldStrength := if.unifiedField.GetStrength()
    return math.Min(1.0, fieldStrength/ResonanceThreshold)
}

// calculateSystemBalance 计算系统平衡度
func (if *IntegrateFlow) calculateSystemBalance() float64 {
    if if.systemState.Energy <= 0 {
        return 1
    }

    // 计算各子系统能量比例的方差
    totalEnergy := if.systemState.Energy
    energyRatios := []float64{
        if.yinyang.GetState().Energy / totalEnergy,
        if.wuxing.GetState().Energy / totalEnergy,
        if.bagua.GetState().Energy / totalEnergy,
        if.ganzhi.GetState().Energy / totalEnergy,
    }

    variance := 0.0
    meanRatio := 0.25 // 理想平均比例
    for _, ratio := range energyRatios {
        diff := ratio - meanRatio
        variance += diff * diff
    }
    variance /= 4

    // 转换为平衡度（0-1）
    return 1 - math.Min(1, variance/IntegrateBalance)
}

// GetSystemState 获取系统状态
func (if *IntegrateFlow) GetSystemState() SystemState {
    if.mu.RLock()
    defer if.mu.RUnlock()
    return if.systemState
}

// Close 关闭集成模型
func (if *IntegrateFlow) Close() error {
    if err := if.Stop(); err != nil {
        return err
    }

    // 关闭子模型
    if err := if.yinyang.Close(); err != nil {
        return err
    }
    if err := if.wuxing.Close(); err != nil {
        return err
    }
    if err := if.bagua.Close(); err != nil {
        return err
    }
    if err := if.ganzhi.Close(); err != nil {
        return err
    }

    return if.BaseFlowModel.Close()
}
