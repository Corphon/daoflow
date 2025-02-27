// core/config.go
package core

import "time"

// FieldConfig 场配置
type FieldConfig struct {
	// 基础配置
	Type      FieldType // 场类型
	Dimension int       // 场维度
	GridSize  int       // 网格大小
	Boundary  []float64 // 边界条件

	// 场特性配置
	InitialStrength float64       // 初始强度
	UpdateInterval  time.Duration // 更新间隔
	WaveNumber      float64       // 波数
	Frequency       float64       // 频率

	// 相互作用配置
	Coupling    float64 // 耦合强度
	Interaction float64 // 相互作用强度
}

// QuantumConfig 量子配置
type QuantumConfig struct {
	// 基础配置
	InitialState    []complex128 // 初始量子态
	Dimension       int          // 维度
	MaxEntanglement float64      // 最大纠缠度

	// 量子特性
	CoherenceTime    time.Duration // 相干时间
	DecoherenceRate  float64       // 退相干率
	EntanglementRate float64       // 纠缠率
	UpdateInterval   time.Duration // 更新间隔
}

// EnergyConfig 能量配置
type EnergyConfig struct {
	// 能量限制
	MinEnergy float64 // 最小能量
	MaxEnergy float64 // 最大能量

	// 能量动态
	DissipationRate float64       // 能量耗散率
	ExchangeRate    float64       // 能量交换率
	UpdateInterval  time.Duration // 更新间隔

	// 能量分布
	InitialDistribution map[string]float64 // 初始能量分布
}

//--------------------------------------------------------
// DefaultFieldConfig 返回默认场配置
func DefaultFieldConfig() *FieldConfig {
	return &FieldConfig{
		Type:            ScalarField,
		Dimension:       3,
		GridSize:        DefaultGridSize,
		InitialStrength: 1.0,
		UpdateInterval:  time.Second / 10,
		WaveNumber:      1.0,
		Frequency:       1.0,
		Coupling:        0.5,
		Interaction:     0.5,
	}
}

// DefaultQuantumConfig 返回默认量子配置
func DefaultQuantumConfig() *QuantumConfig {
	return &QuantumConfig{
		Dimension:        3,
		MaxEntanglement:  1.0,
		CoherenceTime:    time.Second,
		DecoherenceRate:  0.1,
		EntanglementRate: 0.1,
		UpdateInterval:   time.Second / 100,
	}
}

// DefaultEnergyConfig 返回默认能量配置
func DefaultEnergyConfig() *EnergyConfig {
	return &EnergyConfig{
		MinEnergy:       0.0,
		MaxEnergy:       1000.0,
		DissipationRate: 0.01,
		ExchangeRate:    0.1,
		UpdateInterval:  time.Second / 10,
		InitialDistribution: map[string]float64{
			"field":   0.4,
			"quantum": 0.3,
			"system":  0.3,
		},
	}
}
