//core/engine.go

package core

import (
	"math"
	"sync"
	"time"
)

// Status 定义引擎状态
type Status string

const (
	StatusInitialized Status = "initialized"
	StatusRunning     Status = "running"
	StatusStopped     Status = "stopped"
	StatusError       Status = "error"
)

// Engine 核心引擎
type Engine struct {
	mu sync.RWMutex

	// 基础组件
	components struct {
		field     *Field        // 统一场
		quantum   *QuantumState // 量子态
		energy    *EnergySystem // 能量系统
		resonator *Resonator    // 共振器
	}

	// 引擎状态
	state struct {
		status    string        // 运行状态
		startTime time.Time     // 启动时间
		metrics   EngineMetrics // 引擎指标
	}

	// 系统组件
	energySystem  *EnergySystem  // 能量系统
	quantumSystem *QuantumSystem // 量子系统
	networkSystem *EnergyNetwork // 网络系统
	fieldSystem   *FieldSystem   // 场系统

	// 系统参数
	maxEnergy float64 // 最大能量值

	// 配置信息
	config *Config
}

// Config 引擎配置
type Config struct {
	Base struct {
		MaxEnergy     float64       // 最大能量
		UpdateRate    time.Duration // 更新频率
		MaxGoroutines int           // 最大协程数
	}
	Field   *FieldConfig   // 场配置
	Quantum *QuantumConfig // 量子配置
	Energy  *EnergyConfig  // 能量配置
}

// EngineMetrics 引擎指标
type EngineMetrics struct {
	Uptime        time.Duration      // 运行时间
	EnergyLevel   float64            // 能量水平
	FieldStrength float64            // 场强度
	Coherence     float64            // 相干度
	Performance   map[string]float64 // 性能指标
}

// ----------------------------------------------
// NewEngine 创建新的引擎实例
func NewEngine(cfg *Config) (*Engine, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	e := &Engine{
		config: cfg,
	}

	// 初始化组件
	e.components.field = NewField(ScalarField, cfg.Field.Dimension)
	e.components.quantum = NewQuantumState()
	e.components.energy = NewEnergySystem(cfg.Base.MaxEnergy)
	e.components.resonator = NewResonator()

	// 初始化状态
	e.state.status = "initialized"
	e.state.startTime = time.Now()

	return e, nil
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Base: struct {
			MaxEnergy     float64
			UpdateRate    time.Duration
			MaxGoroutines int
		}{
			MaxEnergy:     1000,
			UpdateRate:    time.Second,
			MaxGoroutines: 100,
		},
		Field:   DefaultFieldConfig(),
		Quantum: DefaultQuantumConfig(),
		Energy:  DefaultEnergyConfig(),
	}
}

// GetTotalEnergy 获取系统总能量
func (e *Engine) GetTotalEnergy() float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// 基础能量系统
	baseEnergy := e.energySystem.GetTotalEnergy()

	// 获取各个子系统的能量
	var totalEnergy float64

	// 量子系统能量
	if e.quantumSystem != nil {
		for _, state := range e.quantumSystem.GetStates() {
			totalEnergy += state.GetEnergy()
		}
	}

	// 网络系统能量
	if e.networkSystem != nil {
		totalEnergy += e.networkSystem.GetTotalEnergy()
	}

	// 场系统能量
	if e.fieldSystem != nil {
		totalEnergy += e.fieldSystem.GetEnergy()
	}

	// 返回归一化的总能量 (0-1范围)
	systemEnergy := (baseEnergy + totalEnergy) / e.maxEnergy
	return math.Max(0.0, math.Min(1.0, systemEnergy))
}

// Initialize 初始化引擎
func (e *Engine) Initialize() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 初始化组件
	e.components.field = NewField(ScalarField, e.config.Field.Dimension)
	e.components.quantum = NewQuantumState()
	e.components.energy = NewEnergySystem(e.config.Base.MaxEnergy)
	e.components.resonator = NewResonator()

	// 更新状态
	e.state.status = string(StatusInitialized)
	e.state.startTime = time.Now()

	return nil
}

// Shutdown 关闭引擎
func (e *Engine) Shutdown() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 关闭组件
	e.components.field = nil
	e.components.quantum = nil
	e.components.energy = nil
	e.components.resonator = nil

	// 更新状态
	e.state.status = string(StatusStopped)

	return nil
}

// Status 获取引擎状态
func (e *Engine) Status() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.state.status
}

// GetState 获取引擎状态
func (e *Engine) GetState() map[string]interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	state := make(map[string]interface{})

	// 能量系统状态
	state["energy"] = map[string]interface{}{
		"total":    e.energySystem.GetTotalEnergy(),
		"balance":  e.energySystem.GetBalance(),
		"capacity": e.energySystem.GetCapacity(),
	}

	// 量子系统状态
	if e.quantumSystem != nil {
		state["quantum"] = map[string]interface{}{
			"coherence":    e.quantumSystem.GetCoherence(),
			"entanglement": e.quantumSystem.GetEntanglement(),
		}
	}

	// 场系统状态
	if e.fieldSystem != nil {
		state["field"] = map[string]interface{}{
			"energy":    e.fieldSystem.GetEnergy(),
			"strength":  e.fieldSystem.GetStrength(),
			"coupling":  e.fieldSystem.GetCoupling(),
			"resonance": e.fieldSystem.GetResonance(),
		}
	}

	return state
}

// GetEnergySystem 获取能量系统
func (e *Engine) GetEnergySystem() *EnergySystem {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.energySystem
}

// GetFieldSystem 获取场系统
func (e *Engine) GetFieldSystem() *FieldSystem {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.fieldSystem
}

// GetQuantumSystem 获取量子系统
func (e *Engine) GetQuantumSystem() *QuantumSystem {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.quantumSystem
}
