// model/flow_yinyang.go

package model

import (
    "math"
    "sync"
    "time"

    "github.com/Corphon/daoflow/core"
)

// YinYangConstants 阴阳常数
const (
    MaxImbalance    = 0.3   // 最大失衡度 (30%)
    NeutralPoint    = 0.5   // 中性点
    CycleInterval   = 12.0  // 周期间隔(小时)
    TransformThreshold = 0.8 // 转化阈值
)

// YinYangFlow 阴阳模型
type YinYangFlow struct {
    *BaseFlowModel
    
    // 阴阳特性
    yinRatio  float64 // 阴性比例 (0-1)
    yangRatio float64 // 阳性比例 (0-1)
    
    // 周期控制
    cyclePeriod float64       // 周期长度
    phaseOffset float64       // 相位偏移
    lastCycle   time.Time     // 上次周期时间
    
    // 波动特性
    waveAmplitude float64    // 波动幅度
    waveFrequency float64    // 波动频率
    damping      float64     // 阻尼系数
}

// NewYinYangFlow 创建阴阳流模型
func NewYinYangFlow() *YinYangFlow {
    yy := &YinYangFlow{
        BaseFlowModel: NewBaseFlowModel(ModelYinYang, 100.0),
        yinRatio:     0.5,  // 初始平衡
        yangRatio:    0.5,
        cyclePeriod:  CycleInterval * float64(time.Hour),
        phaseOffset:  0,
        waveAmplitude: 0.1, // 10%波动
        waveFrequency: 2 * math.Pi / CycleInterval,
        damping:      0.05,
    }
    
    go yy.runCycle()
    return yy
}

// runCycle 运行阴阳周期
func (yy *YinYangFlow) runCycle() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-yy.done:
            return
        case <-ticker.C:
            yy.updateCycle()
        }
    }
}

// updateCycle 更新阴阳周期
func (yy *YinYangFlow) updateCycle() {
    yy.mu.Lock()
    defer yy.mu.Unlock()
    
    now := time.Now()
    elapsed := now.Sub(yy.lastCycle).Hours()
    
    // 使用简谐运动模型计算阴阳变化
    // x(t) = A * e^(-γt) * cos(ωt + φ)
    amplitude := yy.waveAmplitude * math.Exp(-yy.damping*elapsed)
    phase := yy.waveFrequency*elapsed + yy.phaseOffset
    oscillation := amplitude * math.Cos(phase)
    
    // 更新阴阳比例
    baseRatio := NeutralPoint + oscillation
    yy.yinRatio = math.Max(0, math.Min(1, baseRatio))
    yy.yangRatio = 1 - yy.yinRatio
    
    // 应用物理效应
    yy.applyPhysicalEffects()
    
    yy.lastCycle = now
}

// applyPhysicalEffects 应用物理效应
func (yy *YinYangFlow) applyPhysicalEffects() {
    // 能量守恒: E = Ek + Ep
    // Ek: 动能 (阳) 
    // Ep: 势能 (阴)
    totalEnergy := yy.state.Energy
    
    // 计算动势能分配
    kineticEnergy := totalEnergy * yy.yangRatio
    potentialEnergy := totalEnergy * yy.yinRatio
    
    // 更新能量系统
    yy.energy.TransformEnergy(map[core.EnergyType]float64{
        core.KineticEnergy:   kineticEnergy,
        core.PotentialEnergy: potentialEnergy,
    })
    
    // 更新场强度
    fieldStrength := yy.calculateFieldStrength()
    yy.field.SetStrength(fieldStrength)
    
    // 更新物理特性
    yy.physics.ApplyYinYangTransformation(yy.yinRatio)
}

// calculateFieldStrength 计算场强度
func (yy *YinYangFlow) calculateFieldStrength() float64 {
    // 基于阴阳比例的场强度计算
    // 使用双曲正切函数使场强度在边界处平滑
    imbalance := math.Abs(yy.yinRatio - NeutralPoint)
    normalizedImbalance := imbalance / MaxImbalance
    
    // tanh函数将值映射到(-1,1)区间
    fieldStrength := math.Tanh(normalizedImbalance)
    
    return fieldStrength * yy.state.Energy
}

// Transform 实现阴阳转化
func (yy *YinYangFlow) Transform(pattern TransformPattern) error {
    yy.mu.Lock()
    defer yy.mu.Unlock()
    
    if pattern.TransformRatio > TransformThreshold {
        // 阴阳互转
        yy.yinRatio, yy.yangRatio = yy.yangRatio, yy.yinRatio
        
        // 能量转换
        currentEnergy := yy.state.Energy
        yy.state.Energy = currentEnergy * pattern.TransformRatio
        
        // 相位调整
        yy.phaseOffset = math.Pi - yy.phaseOffset
        
        // 更新物理效应
        yy.applyPhysicalEffects()
    }
    
    return nil
}

// Interact 实现阴阳相互作用
func (yy *YinYangFlow) Interact(other FlowModel) error {
    yy.mu.Lock()
    defer yy.mu.Unlock()
    
    // 获取对方状态
    otherState := other.GetState()
    
    // 计算相互作用强度
    interactionStrength := yy.calculateInteractionStrength(otherState)
    
    // 能量交换
    energyTransfer := yy.calculateEnergyTransfer(other, interactionStrength)
    
    // 更新状态
    yy.state.Energy += energyTransfer
    
    // 记录交互
    yy.interactions[other.GetModelType().String()] = InteractionRecord{
        Timestamp: time.Now(),
        Target:    other.GetModelType(),
        Effect:    energyTransfer,
        Duration:  time.Second,
    }
    
    return nil
}

// calculateInteractionStrength 计算相互作用强度
func (yy *YinYangFlow) calculateInteractionStrength(otherState ModelState) float64 {
    // 基于阴阳相性计算作用强度
    natureCompatibility := 1.0
    if otherState.Nature == NatureYin && yy.yinRatio > yy.yangRatio {
        natureCompatibility = 1.5 // 阴阴相应增强
    } else if otherState.Nature == NatureYang && yy.yangRatio > yy.yinRatio {
        natureCompatibility = 1.5 // 阳阳相应增强
    }
    
    // 考虑能量差异
    energyDiff := math.Abs(yy.state.Energy - otherState.Energy)
    energyFactor := 1.0 / (1.0 + energyDiff/100.0)
    
    return natureCompatibility * energyFactor
}
