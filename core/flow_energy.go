// core/flow_energy.go

package core

import (
    "math"
    "sync"
)

// EnergyType 能量类型
type EnergyType uint8

const (
    PotentialEnergy EnergyType = iota // 势能
    KineticEnergy                     // 动能
    ThermalEnergy                     // 热能
    FieldEnergy                       // 场能
)

// 能量转换常数
const (
    MinEnergy       = 0.0        // 最小能量
    MaxEnergy       = 1000.0     // 最大能量
    EntropyFactor   = 0.01       // 熵增因子
    DissipationRate = 0.05       // 能量耗散率
)

// EnergySystem 能量系统
type EnergySystem struct {
    mu sync.RWMutex
    
    // 能量组成
    potential float64      // 势能储存
    kinetic   float64      // 动能储存
    thermal   float64      // 热能储存
    field     float64      // 场能储存
    
    // 系统特性
    entropy   float64      // 系统熵
    capacity  float64      // 能量容量
    balance   float64      // 能量平衡度
    
    // 转换效率
    conversionEfficiency map[EnergyType]map[EnergyType]float64
}

// NewEnergySystem 创建能量系统
func NewEnergySystem(capacity float64) *EnergySystem {
    es := &EnergySystem{
        capacity: math.Max(MinEnergy, math.Min(capacity, MaxEnergy)),
        conversionEfficiency: make(map[EnergyType]map[EnergyType]float64),
    }
    
    // 初始化转换效率矩阵
    es.initConversionEfficiency()
    return es
}

// initConversionEfficiency 初始化能量转换效率
func (es *EnergySystem) initConversionEfficiency() {
    types := []EnergyType{PotentialEnergy, KineticEnergy, ThermalEnergy, FieldEnergy}
    
    for _, from := range types {
        es.conversionEfficiency[from] = make(map[EnergyType]float64)
        for _, to := range types {
            if from == to {
                es.conversionEfficiency[from][to] = 1.0
            } else {
                // 默认转换效率0.8
                es.conversionEfficiency[from][to] = 0.8
            }
        }
    }
}

// Convert 能量转换
func (es *EnergySystem) Convert(from, to EnergyType, amount float64) float64 {
    es.mu.Lock()
    defer es.mu.Unlock()
    
    if amount <= 0 {
        return 0
    }
    
    // 获取转换效率
    efficiency := es.conversionEfficiency[from][to]
    
    // 计算转换后的能量
    converted := amount * efficiency
    
    // 计算熵增
    entropyIncrease := amount * (1 - efficiency) * EntropyFactor
    es.entropy += entropyIncrease
    
    // 更新能量存储
    es.decreaseEnergy(from, amount)
    es.increaseEnergy(to, converted)
    
    return converted
}

// TransformEnergy 能量形态转换
func (es *EnergySystem) TransformEnergy(energyMap map[EnergyType]float64) error {
    es.mu.Lock()
    defer es.mu.Unlock()
    
    // 计算总能量确保守恒
    var totalEnergy float64
    for _, amount := range energyMap {
        totalEnergy += amount
    }
    
    if totalEnergy > es.capacity {
        return ErrExceedCapacity
    }
    
    // 更新各能量形态
    es.potential = energyMap[PotentialEnergy]
    es.kinetic = energyMap[KineticEnergy]
    es.thermal = energyMap[ThermalEnergy]
    es.field = energyMap[FieldEnergy]
    
    // 计算能量平衡度
    es.calculateBalance()
    
    return nil
}

// 增加特定类型的能量
func (es *EnergySystem) increaseEnergy(typ EnergyType, amount float64) {
    switch typ {
    case PotentialEnergy:
        es.potential += amount
    case KineticEnergy:
        es.kinetic += amount
    case ThermalEnergy:
        es.thermal += amount
    case FieldEnergy:
        es.field += amount
    }
}

// 减少特定类型的能量
func (es *EnergySystem) decreaseEnergy(typ EnergyType, amount float64) {
    switch typ {
    case PotentialEnergy:
        es.potential = math.Max(0, es.potential-amount)
    case KineticEnergy:
        es.kinetic = math.Max(0, es.kinetic-amount)
    case ThermalEnergy:
        es.thermal = math.Max(0, es.thermal-amount)
    case FieldEnergy:
        es.field = math.Max(0, es.field-amount)
    }
}

// calculateBalance 计算能量平衡度
func (es *EnergySystem) calculateBalance() {
    totalEnergy := es.potential + es.kinetic + es.thermal + es.field
    if totalEnergy == 0 {
        es.balance = 1.0
        return
    }
    
    // 计算能量分布的标准差
    mean := totalEnergy / 4
    variance := (
        math.Pow(es.potential-mean, 2) +
        math.Pow(es.kinetic-mean, 2) +
        math.Pow(es.thermal-mean, 2) +
        math.Pow(es.field-mean, 2)
    ) / 4
    
    // 平衡度 = 1 / (1 + 标准差/总能量)
    es.balance = 1 / (1 + math.Sqrt(variance)/totalEnergy)
}

// GetEnergyState 获取能量状态
func (es *EnergySystem) GetEnergyState() map[string]float64 {
    es.mu.RLock()
    defer es.mu.RUnlock()
    
    return map[string]float64{
        "potential": es.potential,
        "kinetic":   es.kinetic,
        "thermal":   es.thermal,
        "field":     es.field,
        "entropy":   es.entropy,
        "balance":   es.balance,
    }
}
