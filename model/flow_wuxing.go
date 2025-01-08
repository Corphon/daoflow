// model/flow_wuxing.go

package model

import (
    "math"
    "sync"
    "time"

    "github.com/Corphon/daoflow/core"
)

// WuXingConstants 五行常数
const (
    BaseInteraction = 0.2   // 基础相互作用强度
    CycleStrength   = 0.15  // 周期作用强度
    MaxPhaseEnergy  = 100.0 // 最大相位能量
    BalancePoint    = 20.0  // 平衡点能量值
)

// WuXingPhase 五行相位
type WuXingPhase uint8

const (
    Wood WuXingPhase = iota // 木
    Fire                    // 火
    Earth                   // 土
    Metal                   // 金
    Water                   // 水
)

// PhaseRelation 相位关系
type PhaseRelation struct {
    Source      WuXingPhase
    Target      WuXingPhase
    Type        RelationType
    Strength    float64
}

// RelationType 关系类型
type RelationType uint8

const (
    Generate RelationType = iota // 相生
    Restrict              // 相克
    Counter               // 相泄
)

// WuXingFlow 五行模型
type WuXingFlow struct {
    *BaseFlowModel
    
    // 相位能量
    phaseEnergies map[WuXingPhase]float64
    
    // 相位场效应
    phaseFields map[WuXingPhase]*core.Field
    
    // 关系矩阵
    relations []PhaseRelation
    
    // 量子状态
    quantumStates map[WuXingPhase]*core.QuantumState
}

// NewWuXingFlow 创建五行流模型
func NewWuXingFlow() *WuXingFlow {
    wx := &WuXingFlow{
        BaseFlowModel:  NewBaseFlowModel(ModelWuXing, MaxPhaseEnergy*5),
        phaseEnergies:  make(map[WuXingPhase]float64),
        phaseFields:    make(map[WuXingPhase]*core.Field),
        quantumStates:  make(map[WuXingPhase]*core.QuantumState),
    }
    
    wx.initializePhases()
    wx.initializeRelations()
    
    go wx.runCycle()
    return wx
}

// initializePhases 初始化相位
func (wx *WuXingFlow) initializePhases() {
    // 初始化各相位
    for phase := Wood; phase <= Water; phase++ {
        // 能量初始化
        wx.phaseEnergies[phase] = BalancePoint
        
        // 场初始化
        wx.phaseFields[phase] = core.NewField()
        wx.phaseFields[phase].SetStrength(1.0)
        
        // 量子态初始化
        wx.quantumStates[phase] = core.NewQuantumState()
        wx.quantumStates[phase].SetPhase(float64(phase) * 2 * math.Pi / 5)
    }
    
    // 更新状态属性
    wx.updateStateProperties()
}

// initializeRelations 初始化关系
func (wx *WuXingFlow) initializeRelations() {
    // 相生关系
    generateRelations := []PhaseRelation{
        {Wood, Fire, Generate, BaseInteraction},
        {Fire, Earth, Generate, BaseInteraction},
        {Earth, Metal, Generate, BaseInteraction},
        {Metal, Water, Generate, BaseInteraction},
        {Water, Wood, Generate, BaseInteraction},
    }
    
    // 相克关系
    restrictRelations := []PhaseRelation{
        {Wood, Earth, Restrict, BaseInteraction},
        {Earth, Water, Restrict, BaseInteraction},
        {Water, Fire, Restrict, BaseInteraction},
        {Fire, Metal, Restrict, BaseInteraction},
        {Metal, Wood, Restrict, BaseInteraction},
    }
    
    wx.relations = append(generateRelations, restrictRelations...)
}

// runCycle 运行五行周期
func (wx *WuXingFlow) runCycle() {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-wx.done:
            return
        case <-ticker.C:
            wx.processCycle()
        }
    }
}

// processCycle 处理五行周期
func (wx *WuXingFlow) processCycle() {
    wx.mu.Lock()
    defer wx.mu.Unlock()

    // 处理量子态演化
    wx.evolveQuantumStates()
    
    // 处理相互作用
    wx.processInteractions()
    
    // 更新场效应
    wx.updateFields()
    
    // 平衡能量
    wx.balanceEnergies()
    
    // 更新状态
    wx.updateStateProperties()
}

// evolveQuantumStates 演化量子态
func (wx *WuXingFlow) evolveQuantumStates() {
    for phase, state := range wx.quantumStates {
        // 应用时间演化算子
        state.Evolve(time.Second)
        
        // 更新相位能量
        probability := state.GetProbability()
        wx.phaseEnergies[phase] = probability * MaxPhaseEnergy
    }
}

// processInteractions 处理相互作用
func (wx *WuXingFlow) processInteractions() {
    for _, relation := range wx.relations {
        sourceEnergy := wx.phaseEnergies[relation.Source]
        targetEnergy := wx.phaseEnergies[relation.Target]
        
        // 计算相互作用强度
        strength := wx.calculateInteractionStrength(sourceEnergy, targetEnergy, relation)
        
        // 应用相互作用
        wx.applyInteraction(relation, strength)
    }
}

// calculateInteractionStrength 计算相互作用强度
func (wx *WuXingFlow) calculateInteractionStrength(
    sourceEnergy, targetEnergy float64,
    relation PhaseRelation,
) float64 {
    // 基于能量差计算基础强度
    energyDiff := sourceEnergy - targetEnergy
    baseStrength := relation.Strength * math.Tanh(energyDiff/BalancePoint)
    
    // 考虑量子态的相干性
    sourceState := wx.quantumStates[relation.Source]
    targetState := wx.quantumStates[relation.Target]
    coherence := core.CalculateCoherence(sourceState, targetState)
    
    return baseStrength * coherence
}

// applyInteraction 应用相互作用
func (wx *WuXingFlow) applyInteraction(relation PhaseRelation, strength float64) {
    switch relation.Type {
    case Generate:
        // 相生：能量转移
        transferAmount := strength * wx.phaseEnergies[relation.Source]
        wx.phaseEnergies[relation.Source] -= transferAmount
        wx.phaseEnergies[relation.Target] += transferAmount
        
    case Restrict:
        // 相克：能量抑制
        suppressAmount := strength * wx.phaseEnergies[relation.Target]
        wx.phaseEnergies[relation.Target] -= suppressAmount
        
    case Counter:
        // 相泄：能量耗散
        dissipateAmount := strength * math.Min(
            wx.phaseEnergies[relation.Source],
            wx.phaseEnergies[relation.Target],
        )
        wx.phaseEnergies[relation.Source] -= dissipateAmount
        wx.phaseEnergies[relation.Target] -= dissipateAmount
    }
}

// updateFields 更新场效应
func (wx *WuXingFlow) updateFields() {
    for phase, energy := range wx.phaseEnergies {
        field := wx.phaseFields[phase]
        
        // 更新场强度
        normalizedEnergy := energy / MaxPhaseEnergy
        field.SetStrength(normalizedEnergy)
        
        // 应用场效应
        wx.corePhysics.ApplyField(field)
    }
}

// updateStateProperties 更新状态属性
func (wx *WuXingFlow) updateStateProperties() {
    // 更新总能量
    var totalEnergy float64
    for _, energy := range wx.phaseEnergies {
        totalEnergy += energy
    }
    wx.state.Energy = totalEnergy
    
    // 更新相位属性
    for phase, energy := range wx.phaseEnergies {
        wx.state.Properties[phase.String()] = energy
    }
    
    // 更新相位
    wx.state.Phase = PhaseWuXing
}

// GetPhaseEnergy 获取相位能量
func (wx *WuXingFlow) GetPhaseEnergy(phase WuXingPhase) float64 {
    wx.mu.RLock()
    defer wx.mu.RUnlock()
    return wx.phaseEnergies[phase]
}
