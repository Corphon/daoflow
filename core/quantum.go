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
func (qs *QuantumState) Initialize() {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	qs.probability = MaxProbability
	qs.phase = DefaultPhase
	qs.energy = DefaultEnergy
	qs.entropy = DefaultEntropy
}

// SetProbability 设置概率幅度
func (qs *QuantumState) SetProbability(p float64) {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	qs.probability = math.Max(MinProbability, math.Min(MaxProbability, p))
	// 更新熵
	qs.updateEntropy()
}

// GetProbability 获取概率幅度
func (qs *QuantumState) GetProbability() float64 {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	return qs.probability
}

// SetPhase 设置量子相位
func (qs *QuantumState) SetPhase(phase float64) {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	// 确保相位在 [0, 2π) 范围内
	qs.phase = math.Mod(phase, TwoPi)
	if qs.phase < 0 {
		qs.phase += TwoPi
	}
}

// GetPhase 获取量子相位
func (qs *QuantumState) GetPhase() float64 {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	return qs.phase
}

// SetEnergy 设置能量水平
func (qs *QuantumState) SetEnergy(energy float64) {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	qs.energy = math.Max(0, energy)
	// 更新熵
	qs.updateEntropy()
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
func (qs *QuantumState) Evolve(pattern QuantumPattern) {
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
		// 默认模式: 小幅度演化
		qs.phase += math.Pi / 16
		qs.probability = math.Max(0.1, qs.probability)
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
}

// Entangle 量子纠缠
// 将当前量子态与另一个量子态纠缠
func (qs *QuantumState) Entangle(other *QuantumState) error {
	if other == nil {
		return ErrInvalidQuantumState
	}

	qs.mu.Lock()
	defer qs.mu.Unlock()

	// 计算纠缠态的概率和相位
	avgProb := (qs.probability + other.GetProbability()) / 2
	avgPhase := math.Mod((qs.phase+other.GetPhase())/2, TwoPi)

	// 更新当前量子态
	qs.probability = avgProb
	qs.phase = avgPhase

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

// Reset 重置量子态到初始状态
func (qs *QuantumState) Reset() {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	qs.probability = MaxProbability
	qs.phase = DefaultPhase
	qs.energy = DefaultEnergy
	qs.entropy = DefaultEntropy
}
