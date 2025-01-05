// core/flow_physics.go

package core

import (
    "fmt"  
    "math"
    "sync"
)

// 定义常量
const (
    G               = 9.80665    // 重力加速度 (m/s²)
    StdTemperature  = 298.15     // 标准温度 (K)
    StdPressure     = 101.325    // 标准大气压 (kPa)
    StdDensity      = 1.0        // 标准密度 (kg/m³) 
    SpecificHeat    = 4186.0     // 比热容 J/(kg·K)
    BoltzmannConst  = 1.380649e-23 // 玻尔兹曼常数
    PlanckConst     = 6.62607015e-34 // 普朗克常数
)

// FlowPhysics 流体物理特性
type FlowPhysics struct {
    mu sync.RWMutex

    // 基本物理属性 
    density     float64 // 密度 (kg/m³)
    viscosity   float64 // 粘性系数 (Pa·s)
    temperature float64 // 温度 (K)
    pressure    float64 // 压力 (kPa)
    volume      float64 // 体积 (m³)
    entropy     float64 // 熵 (J/K)

    // 动力学属性
    velocity    Vector3D // 速度向量 (m/s)
    acceleration Vector3D // 加速度向量 (m/s²)
    angularVelocity float64 // 角速度 (rad/s)
    momentum    Vector3D // 动量向量 (kg·m/s)

    // 能量属性  
    kineticEnergy   float64 // 动能 (J)
    potentialEnergy float64 // 势能 (J)
    thermalEnergy   float64 // 热能 (J)
    totalEnergy     float64 // 总能量 (J)

    // 场属性
    field *ForceField // 力场

    // 阴阳属性
    yinRatio    float64 // 阴性比例 (0-1)
    yangRatio   float64 // 阳性比例 (0-1)
}

// Vector3D 三维向量
type Vector3D struct {
    X, Y, Z float64
}

// ForceField 力场定义
type ForceField struct {
    strength  float64     // 场强度 (N/m²)
    gradient  []float64   // 场梯度
    potential [][]float64 // 势场分布
    
    // 场特性
    dissipation float64   // 耗散率
    coherence   float64   // 相干性
    resonance   float64   // 共振频率
}

// NewFlowPhysics 创建新的流体物理实例
func NewFlowPhysics() *FlowPhysics {
    fp := &FlowPhysics{
        density:     StdDensity,
        viscosity:   0.001,  // 水的粘性系数
        temperature: StdTemperature,
        pressure:    StdPressure,
        volume:      1.0,
        entropy:     0,
        
        velocity:    Vector3D{0, 0, 0},
        acceleration: Vector3D{0, 0, 0},
        angularVelocity: 0,
        momentum:    Vector3D{0, 0, 0},
        
        yinRatio:    0.5, // 初始平衡
        yangRatio:   0.5,
        
        field:      newForceField(),
    }
    
    // 初始化能量
    fp.updateEnergy()
    return fp
}

// newForceField 创建新的力场
func newForceField() *ForceField {
    return &ForceField{
        strength:    1.0,
        gradient:    make([]float64, 3),
        potential:   make([][]float64, 3),
        dissipation: 0.1,
        coherence:   0.8,
        resonance:   1.0,
    }
}

// CalculateReynoldsNumber 计算雷诺数
func (fp *FlowPhysics) CalculateReynoldsNumber(characteristicLength float64) float64 {
    fp.mu.RLock()
    defer fp.mu.RUnlock()
    
    speed := fp.calculateSpeed()
    return (fp.density * speed * characteristicLength) / fp.viscosity
}

// CalculateEntropy 计算熵变
func (fp *FlowPhysics) CalculateEntropy() float64 {
    fp.mu.RLock()
    defer fp.mu.RUnlock()
    
    // 基于玻尔兹曼熵公式简化计算
    return BoltzmannConst * math.Log(fp.volume/fp.density)
}

// updateEnergy 更新所有能量分量
func (fp *FlowPhysics) updateEnergy() {
    // 计算动能: E = 1/2 * m * v²
    speed := fp.calculateSpeed()
    mass := fp.density * fp.volume
    fp.kineticEnergy = 0.5 * mass * math.Pow(speed, 2)
    
    // 计算势能: E = m * g * h
    fp.potentialEnergy = mass * G * fp.velocity.Y
    
    // 计算热能: E = m * c * ΔT
    fp.thermalEnergy = mass * SpecificHeat * (fp.temperature - 273.15)
    
    // 总能量
    fp.totalEnergy = fp.kineticEnergy + fp.potentialEnergy + fp.thermalEnergy
}

// ApplyYinYangTransformation 应用阴阳转化
func (fp *FlowPhysics) ApplyYinYangTransformation(yinRatio float64) error {
    fp.mu.Lock()
    defer fp.mu.Unlock()
    
    if yinRatio < 0 || yinRatio > 1 {
        return fmt.Errorf("invalid yin ratio: %f", yinRatio)
    }
    
    // 更新阴阳比例
    fp.yinRatio = yinRatio
    fp.yangRatio = 1 - yinRatio
    
    // 应用阴阳特性变化
    fp.applyYinCharacteristics(yinRatio)
    fp.applyYangCharacteristics(1 - yinRatio)
    
    // 更新能量状态
    fp.updateEnergy()
    
    return nil
}

// applyYinCharacteristics 应用阴性特征
func (fp *FlowPhysics) applyYinCharacteristics(ratio float64) {
    // 阴性特征: 增加粘性,降温,降压,增加熵
    fp.viscosity *= (1 + 0.5*ratio)
    fp.temperature *= (1 - 0.3*ratio)
    fp.pressure *= (1 - 0.2*ratio)
    fp.entropy *= (1 + 0.4*ratio)
}

// applyYangCharacteristics 应用阳性特征
func (fp *FlowPhysics) applyYangCharacteristics(ratio float64) {
    // 阳性特征: 增加动能,降低粘性,升温,升压
    speed := fp.calculateSpeed()
    fp.velocity = Vector3D{
        X: fp.velocity.X * (1 + 0.3*ratio),
        Y: fp.velocity.Y * (1 + 0.3*ratio),
        Z: fp.velocity.Z * (1 + 0.3*ratio),
    }
    fp.viscosity *= (1 - 0.3*ratio)
    fp.temperature *= (1 + 0.2*ratio)
    fp.pressure *= (1 + 0.2*ratio)
}

// calculateSpeed 计算速度大小
func (fp *FlowPhysics) calculateSpeed() float64 {
    return math.Sqrt(
        math.Pow(fp.velocity.X, 2) +
        math.Pow(fp.velocity.Y, 2) +
        math.Pow(fp.velocity.Z, 2),
    )
}

// GetPhysicsState 获取物理状态
func (fp *FlowPhysics) GetPhysicsState() map[string]float64 {
    fp.mu.RLock()
    defer fp.mu.RUnlock()
    
    return map[string]float64{
        "density": fp.density,
        "temperature": fp.temperature,
        "pressure": fp.pressure,
        "entropy": fp.entropy,
        "total_energy": fp.totalEnergy,
        "yin_ratio": fp.yinRatio,
        "yang_ratio": fp.yangRatio,
    }
}
