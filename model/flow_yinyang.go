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
    cyclePeriod float64     // 周期长度
    phaseOffset float64     // 相位偏移
    lastCycle   time.Time   // 上次周期时间
    
    // 振动特性
    waveAmplitude float64   // 波动幅度
    waveFrequency float64   // 波动频率
    damping       float64   // 阻尼系数
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
    
    // 初始化状态属性
    yy.state.Properties["yinRatio"] = yy.yinRatio
    yy.state.Properties["yangRatio"] = yy.yangRatio
    yy.state.Phase = PhaseYinYang
    
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
    
    // 使用量子谐振子模型计算阴阳变化
    // ψ(t) = A * e^(-γt) * cos(ωt + φ)
    amplitude := yy.waveAmplitude * math.Exp(-yy.damping*elapsed)
    phase := yy.waveFrequency*elapsed + yy.phaseOffset
    oscillation := amplitude * math.Cos(phase)
    
    // 更新阴阳比例
    baseRatio := NeutralPoint + oscillation
    yy.yinRatio = math.Max(0, math.Min(1, baseRatio))
    yy.yangRatio = 1 - yy.yinRatio
    
    // 更新物理效应
    yy.applyPhysicalEffects()
    
    // 更新状态
    yy.updateState()
    
    yy.lastCycle = now
}

// applyPhysicalEffects 应用物理效应
func (yy *YinYangFlow) applyPhysicalEffects() {
    // 将阴阳比例转换为物理量
    // 使用核心物理系统进行计算
    physicsState := &core.PhysicsState{
        Temperature: yy.yangRatio * 100,  // 阳性对应温度
        Pressure:    yy.yinRatio * 100,   // 阴性对应压力
        Density:     yy.state.Energy / 100,
        Entropy:     yy.calculateEntropy(),
    }
    
    // 应用物理状态
    yy.corePhysics.ApplyState(physicsState)
}

// calculateEntropy 计算系统熵
func (yy *YinYangFlow) calculateEntropy() float64 {
    // 使用信息熵公式: S = -k * (p_yin * ln(p_yin) + p_yang * ln(p_yang))
    if yy.yinRatio == 0 || yy.yangRatio == 0 {
        return 0
    }
    
    k := 1.0 // 玻尔兹曼常数的类比
    entropy := -k * (
        yy.yinRatio * math.Log(yy.yinRatio) +
        yy.yangRatio * math.Log(yy.yangRatio),
    )
    
    return entropy
}

// updateState 更新状态
func (yy *YinYangFlow) updateState() {
    yy.state.Properties["yinRatio"] = yy.yinRatio
    yy.state.Properties["yangRatio"] = yy.yangRatio
    yy.state.Properties["entropy"] = yy.calculateEntropy()
    
    // 根据阴阳比例确定性质
    if math.Abs(yy.yinRatio - yy.yangRatio) < 0.1 {
        yy.state.Nature = NatureBalance
    } else if yy.yinRatio > yy.yangRatio {
        yy.state.Nature = NatureYin
    } else {
        yy.state.Nature = NatureYang
    }
}

// Transform 实现阴阳转化
func (yy *YinYangFlow) Transform(pattern TransformPattern) error {
    yy.mu.Lock()
    defer yy.mu.Unlock()
    
    if pattern.TransformRatio > TransformThreshold {
        // 阴阳互转
        yy.yinRatio, yy.yangRatio = yy.yangRatio, yy.yinRatio
        
        // 相位调整
        yy.phaseOffset = math.Pi - yy.phaseOffset
        
        // 更新物理效应
        yy.applyPhysicalEffects()
        
        // 更新状态
        yy.updateState()
    }
    
    return nil
}

// GetYinYangRatio 获取阴阳比例
func (yy *YinYangFlow) GetYinYangRatio() (float64, float64) {
    yy.mu.RLock()
    defer yy.mu.RUnlock()
    return yy.yinRatio, yy.yangRatio
}
