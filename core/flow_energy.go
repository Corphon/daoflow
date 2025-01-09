// core/flow_energy.go

package core

import (
	"errors"
	"fmt"
	"math"
	"sync"
)

var (
	ErrExceedCapacity     = errors.New("energy exceeds system capacity")
	ErrInvalidParameter   = errors.New("invalid parameter value")
	ErrInsufficientEnergy = errors.New("insufficient energy for conversion")
)

// EnergyType 能量类型
type EnergyType uint8

const (
	PotentialEnergy EnergyType = iota // 势能 - 储存能量
	KineticEnergy                     // 动能 - 运动能量
	ThermalEnergy                     // 热能 - 热力学能量
	FieldEnergy                       // 场能 - 场相关能量
)

// FlowTypeMap 能量类型名称映射
var FlowTypeMap = map[EnergyType]string{
	PotentialEnergy: "potential",
	KineticEnergy:   "kinetic",
	ThermalEnergy:   "thermal",
	FieldEnergy:     "field",
}

// 能量系统常量
const (
	MinEnergy       = 0.0    // 最小能量
	MaxEnergy       = 1000.0 // 最大能量
	EntropyFactor   = 0.01   // 熵增因子
	DissipationRate = 0.05   // 能量耗散率
	DefaultBalance  = 1.0    // 默认平衡度
)

// EnergySystem 能量系统
type EnergySystem struct {
	mu sync.RWMutex

	// 能量分量
	potential float64 // 势能储存
	kinetic   float64 // 动能储存
	thermal   float64 // 热能储存
	field     float64 // 场能储存

	// 系统特性
	entropy  float64 // 系统熵
	capacity float64 // 能量容量
	balance  float64 // 能量平衡度

	// 转换效率矩阵
	conversionEfficiency map[EnergyType]map[EnergyType]float64
}

// NewEnergySystem 创建能量系统
func NewEnergySystem(capacity float64) *EnergySystem {
	if capacity <= 0 {
		capacity = MaxEnergy
	}

	es := &EnergySystem{
		capacity:             math.Min(capacity, MaxEnergy),
		balance:              DefaultBalance,
		conversionEfficiency: make(map[EnergyType]map[EnergyType]float64),
	}

	es.initConversionEfficiency()
	return es
}

// initConversionEfficiency 初始化能量转换效率
func (es *EnergySystem) initConversionEfficiency() {
	types := []EnergyType{PotentialEnergy, KineticEnergy, ThermalEnergy, FieldEnergy}

	for _, from := range types {
		es.conversionEfficiency[from] = make(map[EnergyType]float64)
		for _, to := range types {
			switch {
			case from == to:
				es.conversionEfficiency[from][to] = 1.0
			case from == PotentialEnergy && to == KineticEnergy:
				es.conversionEfficiency[from][to] = 0.9
			case from == KineticEnergy && to == ThermalEnergy:
				es.conversionEfficiency[from][to] = 0.85
			case from == ThermalEnergy && to == FieldEnergy:
				es.conversionEfficiency[from][to] = 0.8
			default:
				es.conversionEfficiency[from][to] = 0.75
			}
		}
	}
}

// Convert 能量转换
func (es *EnergySystem) Convert(from, to EnergyType, amount float64) (float64, error) {
	es.mu.Lock()
	defer es.mu.Unlock()

	if amount <= 0 {
		return 0, ErrInvalidParameter
	}

	// 检查源能量是否足够
	if es.getEnergy(from) < amount {
		return 0, ErrInsufficientEnergy
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

	// 更新平衡度
	es.calculateBalance()

	return converted, nil
}

// TransformEnergy 能量形态转换
func (es *EnergySystem) TransformEnergy(energyMap map[EnergyType]float64) error {
	es.mu.Lock()
	defer es.mu.Unlock()

	// 计算总能量确保守恒
	var totalEnergy float64
	for _, amount := range energyMap {
		if amount < 0 {
			return ErrInvalidParameter
		}
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

// GetEnergyState 获取能量状态
func (es *EnergySystem) GetEnergyState() map[string]float64 {
	es.mu.RLock()
	defer es.mu.RUnlock()

	return map[string]float64{
		"potential": es.potential,
		"kinetic":   es.kinetic,
		"thermal":   es.thermal,
		"field":     es.field,
		"total":     es.GetTotalEnergy(),
		"entropy":   es.entropy,
		"balance":   es.balance,
		"capacity":  es.capacity,
	}
}

// GetEnergy 获取指定类型的能量
func (es *EnergySystem) GetEnergy(typ EnergyType) float64 {
	es.mu.RLock()
	defer es.mu.RUnlock()
	return es.getEnergy(typ)
}

// GetBalance 获取能量平衡度
func (es *EnergySystem) GetBalance() float64 {
	es.mu.RLock()
	defer es.mu.RUnlock()
	return es.balance
}

// getEnergy 内部获取能量方法（无锁）
func (es *EnergySystem) getEnergy(typ EnergyType) float64 {
	switch typ {
	case PotentialEnergy:
		return es.potential
	case KineticEnergy:
		return es.kinetic
	case ThermalEnergy:
		return es.thermal
	case FieldEnergy:
		return es.field
	default:
		return 0
	}
}

// increaseEnergy 增加特定类型的能量（无锁）
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

// decreaseEnergy 减少特定类型的能量（无锁）
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

// getTotalEnergy 获取总能量（无锁）
func (es *EnergySystem) GetTotalEnergy() float64 {
	return es.potential + es.kinetic + es.thermal + es.field
}

// calculateBalance 计算能量平衡度（无锁）
func (es *EnergySystem) calculateBalance() {
	totalEnergy := es.GetTotalEnergy()
	if totalEnergy == 0 {
		es.balance = DefaultBalance
		return
	}

	// 计算能量分布的标准差
	mean := totalEnergy / 4
	variance := (math.Pow(es.potential-mean, 2) +
		math.Pow(es.kinetic-mean, 2) +
		math.Pow(es.thermal-mean, 2) +
		math.Pow(es.field-mean, 2)) / 4

	// 平衡度 = 1 / (1 + 标准差/总能量)
	es.balance = 1 / (1 + math.Sqrt(variance)/totalEnergy)
}

// String 返回能量系统的字符串表示
func (es *EnergySystem) String() string {
	state := es.GetEnergyState()
	return fmt.Sprintf("EnergySystem{total: %.2f, capacity: %.2f, balance: %.2f, entropy: %.2f}",
		state["total"], state["capacity"], state["balance"], state["entropy"])
}
