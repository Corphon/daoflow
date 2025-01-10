// core/quantum.go

package core

import (
	"errors"
	"fmt"
	"math"
	"sync"
)

// QuantumState 表示一个量子态系统
// 包含概率幅度、相位以及其他量子特性
type QuantumState struct {
	mu          sync.RWMutex
	probability float64 // 概率幅度 (0-1)
	phase       float64 // 量子相位 (0-2π)
	energy      float64 // 能量水平
	entropy     float64 // 系统熵
}

// QuantumPattern 量子演化模式
type QuantumPattern string

const (
	// 演化模式常量
	PatternIntegrate QuantumPattern = "integrate" // 整合模式
	PatternSplit     QuantumPattern = "split"     // 分裂模式
	PatternCycle     QuantumPattern = "cycle"     // 循环模式
	PatternBalance   QuantumPattern = "balance"   // 平衡模式

	// 量子态常量
	MaxProbability = 1.0
	MinProbability = 0.0
	DefaultPhase   = 0.0
	TwoPi          = 2 * math.Pi
	DefaultEnergy  = 1.0
	DefaultEntropy = 0.0
)

// NewQuantumState 创建一个新的量子态
func NewQuantumState() *QuantumState {
	return &QuantumState{
		probability: MaxProbability,
		phase:       DefaultPhase,
		energy:      DefaultEnergy,
		entropy:     DefaultEntropy,
	}
}

// Initialize 初始化量子态
func (qs *QuantumState) Initialize() error {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	// 进行初始化
	qs.probability = MaxProbability
	qs.phase = DefaultPhase
	qs.energy = DefaultEnergy
	qs.entropy = DefaultEntropy

	// 验证初始化状态
	if err := qs.validateState(); err != nil {
		return fmt.Errorf("failed to initialize quantum state: %w", err)
	}

	return nil
}

// validateState 验证量子态
func (qs *QuantumState) validateState() error {
	// 验证概率
	if qs.probability < MinProbability || qs.probability > MaxProbability {
		return fmt.Errorf("invalid probability: %v", qs.probability)
	}

	// 验证相位
	if qs.phase < 0 || qs.phase >= TwoPi {
		return fmt.Errorf("invalid phase: %v", qs.phase)
	}

	// 验证能量
	if qs.energy < 0 {
		return fmt.Errorf("invalid energy: %v", qs.energy)
	}

	// 验证熵
	if qs.entropy < 0 {
		return fmt.Errorf("invalid entropy: %v", qs.entropy)
	}

	return nil
}

// Reset 重置量子态到初始状态
func (qs *QuantumState) Reset() error {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	return qs.Initialize()
}

// SetProbability 设置概率幅度
func (qs *QuantumState) SetProbability(p float64) error {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	if p < MinProbability || p > MaxProbability {
		return fmt.Errorf("probability out of range [%v, %v]: %v", MinProbability, MaxProbability, p)
	}

	qs.probability = p
	qs.updateEntropy()
	return nil
}

// GetProbability 获取概率幅度
func (qs *QuantumState) GetProbability() float64 {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	return qs.probability
}

// SetPhase 设置量子相位
func (qs *QuantumState) SetPhase(phase float64) error {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	// 确保相位在 [0, 2π) 范围内
	phase = math.Mod(phase, TwoPi)
	if phase < 0 {
		phase += TwoPi
	}

	qs.phase = phase
	return nil
}

// GetPhase 获取量子相位
func (qs *QuantumState) GetPhase() float64 {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	return qs.phase
}

// SetEnergy 设置能量水平
func (qs *QuantumState) SetEnergy(energy float64) error {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	if energy < 0 {
		return fmt.Errorf("energy cannot be negative: %v", energy)
	}

	qs.energy = energy
	qs.updateEntropy()
	return nil
}

// GetEnergy 获取能量水平
func (qs *QuantumState) GetEnergy() float64 {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	return qs.energy
}

// GetEntropy 获取系统熵
func (qs *QuantumState) GetEntropy() float64 {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	return qs.entropy
}

// Evolve 量子态演化
func (qs *QuantumState) Evolve(pattern QuantumPattern) error {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	// 根据不同模式进行演化
	switch pattern {
	case PatternIntegrate:
		// 整合模式: 相位变化较大，概率趋于稳定
		qs.phase += math.Pi / 4
		qs.probability = math.Pow(qs.probability, 0.9)
	case PatternSplit:
		// 分裂模式: 相位变化小，概率波动大
		qs.phase += math.Pi / 8
		qs.probability *= 0.95
	case PatternCycle:
		// 循环模式: 相位均匀变化
		qs.phase += math.Pi / 6
		qs.probability = 0.5 + 0.5*math.Sin(qs.phase)
	case PatternBalance:
		// 平衡模式: 概率趋于平衡态
		qs.phase += math.Pi / 12
		qs.probability = (qs.probability + 0.5) / 2
	default:
		return fmt.Errorf("unknown evolution pattern: %v", pattern)
	}

	// 确保相位在 [0, 2π) 范围内
	qs.phase = math.Mod(qs.phase, TwoPi)
	if qs.phase < 0 {
		qs.phase += TwoPi
	}

	// 确保概率在 [0, 1] 范围内
	qs.probability = math.Max(MinProbability, math.Min(MaxProbability, qs.probability))

	// 更新熵
	qs.updateEntropy()

	return nil
}

// Collapse 量子态坍缩
// 将量子态坍缩到一个确定态
func (qs *QuantumState) Collapse() {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	// 根据当前概率决定坍缩后的状态
	if qs.probability >= 0.5 {
		qs.probability = MaxProbability
	} else {
		qs.probability = MinProbability
	}

	qs.phase = DefaultPhase
	qs.updateEntropy()
}

// updateEntropy 更新系统熵
// 使用 Shannon 熵公式: H = -p*log(p) - (1-p)*log(1-p)
func (qs *QuantumState) updateEntropy() {
	p := qs.probability
	if p == 0 || p == 1 {
		qs.entropy = 0
		return
	}

	q := 1 - p
	qs.entropy = -p*math.Log2(p) - q*math.Log2(q)
}

// String 返回量子态的字符串表示
func (qs *QuantumState) String() string {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	return fmt.Sprintf("QuantumState{probability: %.4f, phase: %.4f, energy: %.4f, entropy: %.4f}",
		qs.probability, qs.phase, qs.energy, qs.entropy)
}

// 错误定义
var (
	ErrInvalidQuantumState = errors.New("invalid quantum state")
)

// GetCoherence 获取量子相干性
// 相干性与概率幅度和相位的稳定性相关
func (qs *QuantumState) GetCoherence() float64 {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	// 相干性与概率幅度和相位的稳定性相关
	// 使用概率和相位计算相干性
	phaseContribution := math.Cos(qs.phase)   // 相位对相干性的贡献
	probabilityContribution := qs.probability // 概率对相干性的贡献

	// 相干性在 [0,1] 范围内
	coherence := (phaseContribution + 1) * probabilityContribution / 2
	return math.Max(0, math.Min(1, coherence))
}

// AddEnergy 增加量子态的能量
func (qs *QuantumState) AddEnergy(deltaE float64) error {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	if deltaE < 0 {
		return fmt.Errorf("energy increment cannot be negative: %v", deltaE)
	}

	// 计算新能量
	newEnergy := qs.energy + deltaE

	// 根据能量变化调整概率幅度
	// 使用指数衰减函数确保概率保持在 [0,1] 范围内
	probabilityDelta := (1 - qs.probability) * (1 - math.Exp(-deltaE/qs.energy))
	newProbability := qs.probability + probabilityDelta

	// 确保概率在有效范围内
	if newProbability > MaxProbability {
		newProbability = MaxProbability
	}
	if newProbability < MinProbability {
		newProbability = MinProbability
	}

	// 更新状态
	qs.energy = newEnergy
	qs.probability = newProbability

	// 更新熵
	qs.updateEntropy()

	return nil
}

// Update 更新量子态
func (qs *QuantumState) Update() error {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	// 根据当前的能量和相位更新概率
	coherence := qs.GetCoherence()
	energyFactor := math.Exp(-qs.energy / DefaultEnergy)

	// 更新概率，考虑能量和相干性的影响
	newProbability := qs.probability*coherence*(1-energyFactor) +
		MinProbability*energyFactor

	// 确保概率在有效范围内
	qs.probability = math.Max(MinProbability, math.Min(MaxProbability, newProbability))

	// 更新相位
	qs.phase = math.Mod(qs.phase+math.Pi/4*coherence, TwoPi)

	// 更新熵
	qs.updateEntropy()

	return nil
}
