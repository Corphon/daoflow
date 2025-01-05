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
    Restrict                    // 相克
    Counter                     // 相泄
)

// WuXingFlow 五行模型
type WuXingFlow struct {
    *BaseFlowModel
    
    // 相位能量
    phaseEnergies map[WuXingPhase]float64
    
    // 关系矩阵
    relations []PhaseRelation
    
    // 场效应
    fieldEffects map[WuXingPhase]*FieldEffect
    
    // 周期控制
    cycleControl struct {
        currentPhase WuXingPhase
        cycleTime    time.Duration
        lastCycle    time.Time
    }
}

// FieldEffect 场效应
type FieldEffect struct {
    Strength    float64   // 场强度
    Radius      float64   // 作用范围
    Frequency   float64   // 振动频率
    Phase       float64   // 相位角
}

// NewWuXingFlow 创建五行流模型
func NewWuXingFlow() *WuXingFlow {
    wx := &WuXingFlow{
        BaseFlowModel:  NewBaseFlowModel(ModelWuXing, MaxPhaseEnergy*5),
        phaseEnergies:  make(map[WuXingPhase]float64),
        fieldEffects:   make(map[WuXingPhase]*FieldEffect),
    }
    
    wx.initializePhases()
    wx.initializeRelations()
    
    go wx.runCycle()
    return wx
}

// initializePhases 初始化相位
func (wx *WuXingFlow) initializePhases() {
    // 初始化各相位能量
    for phase := Wood; phase <= Water; phase++ {
        wx.phaseEnergies[phase] = BalancePoint
        
        // 初始化场效应
        wx.fieldEffects[phase] = &FieldEffect{
            Strength:  1.0,
            Radius:    10.0,
            Frequency: 2 * math.Pi / float64(24*time.Hour),
            Phase:     float64(phase) * 2 * math.Pi / 5,
        }
    }
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
    ticker := time.NewTicker(time.Hour)
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

    now := time.Now()
    
    // 更新场效应
    wx.updateFieldEffects(now)
    
    // 处理相互作用
    wx.processInteractions()
    
    // 能量平衡
    wx.balanceEnergies()
    
    wx.cycleControl.lastCycle = now
}

// updateFieldEffects 更新场效应
func (wx *WuXingFlow) updateFieldEffects(t time.Time) {
    elapsed := t.Sub(wx.cycleControl.lastCycle).Seconds()
    
    for phase, effect := range wx.fieldEffects {
        // 使用量子场论的波函数概念
        // ψ(t) = A * e^(-iωt)
        omega := effect.Frequency
        amplitude := effect.Strength
        
        // 计算场强度
        effect.Strength = amplitude * math.Cos(omega*elapsed + effect.Phase)
        
        // 更新相位能量
        energy := wx.phaseEnergies[phase]
        fieldContribution := effect.Strength * CycleStrength
        wx.phaseEnergies[phase] = math.Max(0, math.Min(MaxPhaseEnergy,
            energy + fieldContribution))
    }
}

// processInteractions 处理相互作用
func (wx *WuXingFlow) processInteractions() {
    for _, relation := range wx.relations {
        sourceEnergy := wx.phaseEnergies[relation.Source]
        targetEnergy := wx.phaseEnergies[relation.Target]
        
        // 计算作用强度
        interactionStrength := wx.calculateInteractionStrength(
            sourceEnergy, targetEnergy, relation)
        
        // 应用相互作用
        wx.applyInteraction(relation, interactionStrength)
    }
}

// calculateInteractionStrength 计算作用强度
func (wx *WuXingFlow) calculateInteractionStrength(
    sourceEnergy, targetEnergy float64,
    relation PhaseRelation,
) float64 {
    // 基于能量差计算基础强度
    energyDiff := sourceEnergy - targetEnergy
    baseStrength := relation.Strength * math.Tanh(energyDiff/BalancePoint)
    
    // 考虑场效应
    sourceField := wx.fieldEffects[relation.Source]
    fieldFactor := math.Abs(sourceField.Strength)
    
    return baseStrength * fieldFactor
}

// applyInteraction 应用相互作用
func (wx *WuXingFlow) applyInteraction(
    relation PhaseRelation,
    strength float64,
) {
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

// balanceEnergies 平衡能量
func (wx *WuXingFlow) balanceEnergies() {
    var totalEnergy float64
    for _, energy := range wx.phaseEnergies {
        totalEnergy += energy
    }
    
    // 计算平均能量
    avgEnergy := totalEnergy / 5
    
    // 应用能量平衡
    for phase := range wx.phaseEnergies {
        diff := wx.phaseEnergies[phase] - avgEnergy
        if math.Abs(diff) > BalancePoint {
            adjustment := diff * 0.1 // 渐进调整
            wx.phaseEnergies[phase] -= adjustment
        }
    }
}
