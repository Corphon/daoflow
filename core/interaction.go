// core/interaction.go

package core

import (
	"math"
	"sync"
)

// InteractionType 相互作用类型
type InteractionType uint8

const (
	NoInteraction      InteractionType = iota
	WeakInteraction                    // 弱相互作用
	StrongInteraction                  // 强相互作用
	FieldInteraction                   // 场相互作用
	QuantumInteraction                 // 量子相互作用
)

// InteractionConstants 相互作用常数
const (
	MinCoupling     = 0.0 // 最小耦合强度
	MaxCoupling     = 1.0 // 最大耦合强度
	DefaultCoupling = 0.5 // 默认耦合强度
)

// Interaction 相互作用
type Interaction struct {
	mu sync.RWMutex

	// 相互作用属性
	interactionType InteractionType // 相互作用类型
	coupling        float64         // 耦合强度
	phase           float64         // 相位差
	strength        float64         // 作用强度

	// 相互作用状态
	state struct {
		energy    float64 // 相互作用能量
		entropy   float64 // 相互作用熵
		coherence float64 // 相干性
	}
}

// -----------------------------------------------
// NewInteraction 创建新的相互作用
func NewInteraction() *Interaction {
	return &Interaction{
		interactionType: NoInteraction,
		coupling:        DefaultCoupling,
		phase:           0,
		strength:        0,
	}
}

// Initialize 初始化相互作用
func (i *Interaction) Initialize() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.interactionType = NoInteraction
	i.coupling = DefaultCoupling
	i.phase = 0
	i.strength = 0

	i.state.energy = 0
	i.state.entropy = 0
	i.state.coherence = 1

	return nil
}

// Update 更新相互作用
func (i *Interaction) Update(state1, state2 *QuantumState) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	// 计算相互作用强度
	energy1 := state1.GetEnergy()
	energy2 := state2.GetEnergy()
	phase1 := state1.GetPhase()
	phase2 := state2.GetPhase()

	// 计算相位差
	i.phase = math.Abs(phase1 - phase2)

	// 计算作用强度
	i.strength = i.coupling * math.Sqrt(energy1*energy2)

	// 更新相互作用能量
	i.state.energy = i.strength * math.Cos(i.phase)

	// 更新相干性
	i.state.coherence = math.Exp(-i.phase * i.phase)

	// 更新熵
	if i.strength > 0 {
		i.state.entropy = -i.strength * math.Log(i.strength)
	} else {
		i.state.entropy = 0
	}

	// 确定相互作用类型
	i.determineInteractionType()

	return nil
}

// determineInteractionType 确定相互作用类型
func (i *Interaction) determineInteractionType() {
	if i.strength < 0.2 {
		i.interactionType = WeakInteraction
	} else if i.strength > 0.8 {
		i.interactionType = StrongInteraction
	} else if i.phase < math.Pi/4 {
		i.interactionType = QuantumInteraction
	} else {
		i.interactionType = FieldInteraction
	}
}

// GetType 获取相互作用类型
func (i *Interaction) GetType() InteractionType {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.interactionType
}

// GetStrength 获取作用强度
func (i *Interaction) GetStrength() float64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.strength
}

// GetEnergy 获取相互作用能量
func (i *Interaction) GetEnergy() float64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.state.energy
}

// GetCoherence 获取相干性
func (i *Interaction) GetCoherence() float64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.state.coherence
}

// SetCoupling 设置耦合强度
func (i *Interaction) SetCoupling(coupling float64) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if coupling < MinCoupling || coupling > MaxCoupling {
		return NewCoreError("invalid coupling strength")
	}

	i.coupling = coupling
	return nil
}

// Reset 重置相互作用
func (i *Interaction) Reset() error {
	return i.Initialize()
}
