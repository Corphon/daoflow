// core/quantum.go

package core

import (
	"errors"
	"fmt"
	"math"
	"math/cmplx"
	"sync"
	"time"
)

// QuantumState 表示一个量子态系统
// 包含概率幅度、相位以及其他量子特性
type QuantumState struct {
	mu             sync.RWMutex
	probability    float64      // 概率幅度 (0-1)
	phase          float64      // 量子相位 (0-2π)
	energy         float64      // 能量水平
	entropy        float64      // 系统熵
	amplitude      []complex128 // 改为私有
	phaseVariation float64      // 相位变化率
}

// QuantumPattern 常量 - 量子态演化模式
const (
	PatternIntegrate QuantumPattern = "integrate" // 整合模式
	PatternSplit     QuantumPattern = "split"     // 分裂模式
	PatternCycle     QuantumPattern = "cycle"     // 循环模式
	PatternBalance   QuantumPattern = "balance"   // 平衡模式
)

// PatternType 场的模式类型
type PatternType string

const (
	PatternNone      PatternType = ""          // 无模式
	PatternStable    PatternType = "stable"    // 稳定模式
	PatternChaos     PatternType = "chaos"     // 混沌模式
	PatternOscillate PatternType = "oscillate" // 振荡模式 - 替换原来的 cycle
	PatternSpiral    PatternType = "spiral"    // 螺旋模式
)

// QuantumPattern 量子演化模式
type QuantumPattern string

const (

	// 量子态常量
	MaxProbability = 1.0
	MinProbability = 0.0
	DefaultPhase   = 0.0
	TwoPi          = 2 * math.Pi
	DefaultEnergy  = 1.0
	DefaultEntropy = 0.0
)

// QuantumField 量子场接口
type QuantumField interface {
	// 基础操作
	Initialize() error
	Reset() error
	Update(state *QuantumState) error

	// 获取状态
	GetState() *QuantumState
	GetPhase() float64
	GetCoherence() float64
	GetEntropy() float64

	// 场操作
	Evolve(pattern PatternType) error
	Transform(targetState *QuantumState) error

	// 量子操作
	Entangle(other *QuantumState) error
	Decohere() error
	Measure() (float64, error)
}

// QuantumSystem 量子系统
type QuantumSystem struct {
	mu sync.RWMutex

	// 系统状态
	states map[string]*QuantumState // 量子态集合
	field  QuantumField             // 量子场

	// 系统属性
	entanglement float64 // 系统整体纠缠度
	coherence    float64 // 系统整体相干性
	energy       float64 // 系统总能量

	// 配置
	config *QuantumConfig

	// 缓存
	cache struct {
		lastUpdate time.Time
		metrics    map[string]float64
	}
}

// ----------------------------------------------------
// NewQuantumSystem 创建新的量子系统
func NewQuantumSystem(config *QuantumConfig) *QuantumSystem {
	if config == nil {
		config = DefaultQuantumConfig()
	}

	return &QuantumSystem{
		states: make(map[string]*QuantumState),
		field:  NewQuantumField(ScalarField),
		config: config,
	}
}

// GetStates 获取所有量子态
func (qs *QuantumSystem) GetStates() []*QuantumState {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	states := make([]*QuantumState, 0, len(qs.states))
	for _, state := range qs.states {
		states = append(states, state)
	}
	return states
}

// GetPhaseVariation 获取相位变化率
func (qs *QuantumState) GetPhaseVariation() float64 {
	qs.mu.RLock()
	defer qs.mu.RUnlock()
	return qs.phaseVariation
}

// DotProduct 计算两个量子态的内积
func (qs *QuantumState) DotProduct(other *QuantumState) (complex128, error) {
	if other == nil {
		return 0, fmt.Errorf("other quantum state cannot be nil")
	}

	qs.mu.RLock()
	defer qs.mu.RUnlock()

	other.mu.RLock()
	defer other.mu.RUnlock()

	if len(qs.amplitude) != len(other.amplitude) {
		return 0, fmt.Errorf("quantum states must have the same dimension")
	}

	var result complex128
	for i := range qs.amplitude {
		result += qs.amplitude[i] * other.amplitude[i]
	}
	return result, nil
}

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
	qs.amplitude = make([]complex128, 1)
	qs.amplitude[0] = complex(1, 0) // 初始化为基态

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

	// 验证amplitude
	if qs.amplitude == nil || len(qs.amplitude) == 0 {
		return fmt.Errorf("invalid amplitude: nil or empty")
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

// GetEntropy 获取量子态的熵
func (qs *QuantumState) GetEntropy() float64 {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	// von Neumann熵计算
	if len(qs.amplitude) == 0 {
		return 0
	}

	// 计算概率密度矩阵的特征值
	var probabilities []float64
	totalProb := 0.0
	for _, amp := range qs.amplitude {
		prob := math.Pow(cmplx.Abs(amp), 2)
		probabilities = append(probabilities, prob)
		totalProb += prob
	}

	// 归一化概率
	if totalProb > 0 {
		for i := range probabilities {
			probabilities[i] /= totalProb
		}
	}

	// 计算熵 S = -Tr(ρ ln ρ)
	entropy := 0.0
	for _, p := range probabilities {
		if p > 0 {
			entropy -= p * math.Log(p)
		}
	}

	// 归一化到[0,1]区间
	maxEntropy := math.Log(float64(len(qs.amplitude)))
	if maxEntropy > 0 {
		entropy /= maxEntropy
	}

	return entropy
}

// Evolve 量子态演化
func (qs *QuantumState) Evolve(pattern QuantumPattern) error {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	// 保存初始值用于验证
	initialPhase := qs.phase
	initialProb := qs.probability

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

	// 根据状态变化更新振幅
	phaseDiff := qs.phase - initialPhase

	// 使用概率差异来调整能量
	probDiff := qs.probability - initialProb
	qs.energy *= (1 + probDiff)

	// 更新振幅以反映状态变化
	if len(qs.amplitude) > 0 {
		qs.amplitude[0] *= complex(math.Cos(phaseDiff), math.Sin(phaseDiff))
		// 调整振幅大小以匹配新的概率
		currentAmp := cmplx.Abs(qs.amplitude[0])
		if currentAmp > 0 {
			qs.amplitude[0] *= complex(math.Sqrt(qs.probability)/currentAmp, 0)
		}
	}

	// 更新熵
	qs.updateEntropy()

	return nil
}

// Collapse 量子态坍缩
func (qs *QuantumState) Collapse() error {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	// 根据当前概率决定坍缩后的状态
	if qs.probability >= 0.5 {
		qs.probability = MaxProbability
	} else {
		qs.probability = MinProbability
	}

	qs.phase = DefaultPhase

	// 更新振幅为对应的本征态
	qs.amplitude = make([]complex128, 1)
	if qs.probability == MaxProbability {
		qs.amplitude[0] = complex(1, 0)
	} else {
		qs.amplitude[0] = complex(0, 0)
	}

	qs.updateEntropy()
	return nil
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

// CalculatePurity 计算量子态的纯度
func (qs *QuantumState) CalculatePurity() float64 {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	// 纯度计算 |⟨ψ|ψ⟩|²
	return qs.probability * qs.probability
}

// GetEntanglement 获取量子纠缠度
func (qs *QuantumState) GetEntanglement() float64 {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	// 基于振幅和相位计算纠缠度
	// 纠缠度与量子态的叠加程度相关
	phaseContribution := math.Cos(qs.phase)
	amplitudeContribution := qs.probability * qs.probability

	// 归一化到[0,1]区间
	entanglement := (phaseContribution + 1.0) * amplitudeContribution / 2.0
	return math.Max(0, math.Min(1, entanglement))
}

// GetAmplitude 获取量子态振幅
func (qs *QuantumState) GetAmplitude() []complex128 {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	result := make([]complex128, len(qs.amplitude))
	copy(result, qs.amplitude)
	return result
}

// SetAmplitude 设置量子态振幅
func (qs *QuantumState) SetAmplitude(newAmplitude []complex128) error {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	if newAmplitude == nil {
		return fmt.Errorf("amplitude cannot be nil")
	}

	qs.amplitude = make([]complex128, len(newAmplitude))
	copy(qs.amplitude, newAmplitude)
	return nil
}

// GetAmplitudeValue 获取振幅绝对值
func (qs *QuantumState) GetAmplitudeValue() float64 {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	// 振幅是概率的平方根
	return math.Sqrt(qs.probability)
}

// quantumField 量子场实现
type quantumField struct {
	mu sync.RWMutex

	state     *QuantumState
	fieldType FieldType
	pattern   PatternType

	// 量子缓存
	cache struct {
		lastState  *QuantumState
		lastPhase  float64
		lastEnergy float64
		lastUpdate time.Time
	}
}

// NewQuantumField 创建量子场
func NewQuantumField(fieldType FieldType) QuantumField {
	qf := &quantumField{
		state:     NewQuantumState(),
		fieldType: fieldType,
		pattern:   PatternNone,
	}
	qf.cache.lastState = NewQuantumState()
	qf.cache.lastUpdate = time.Now()
	return qf
}

// Initialize 初始化
func (qf *quantumField) Initialize() error {
	qf.mu.Lock()
	defer qf.mu.Unlock()

	if err := qf.state.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize quantum state: %w", err)
	}

	// 重置缓存
	qf.cache.lastState = NewQuantumState()
	if err := qf.cache.lastState.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize cache state: %w", err)
	}

	qf.cache.lastPhase = qf.state.GetPhase()
	qf.cache.lastEnergy = qf.state.GetEnergy()
	qf.cache.lastUpdate = time.Now()

	return nil
}

// Reset 重置
func (qf *quantumField) Reset() error {
	return qf.Initialize()
}

// Update 更新状态
func (qf *quantumField) Update(state *QuantumState) error {
	qf.mu.Lock()
	defer qf.mu.Unlock()

	// 保存当前状态到缓存
	qf.cache.lastState = qf.state
	qf.cache.lastPhase = qf.state.GetPhase()
	qf.cache.lastEnergy = qf.state.GetEnergy()
	qf.cache.lastUpdate = time.Now()

	// 创建新状态并复制值
	newState := NewQuantumState()
	newState.SetPhase(state.GetPhase())
	newState.SetProbability(state.GetProbability())
	newState.SetEnergy(state.GetEnergy())
	newState.SetAmplitude(state.GetAmplitude())

	qf.state = newState

	return nil
}

// GetState 获取状态
func (qf *quantumField) GetState() *QuantumState {
	qf.mu.RLock()
	defer qf.mu.RUnlock()
	return qf.state
}

// GetPhase 获取相位
func (qf *quantumField) GetPhase() float64 {
	qf.mu.RLock()
	defer qf.mu.RUnlock()
	return qf.state.GetPhase()
}

// GetCoherence 获取相干性
func (qf *quantumField) GetCoherence() float64 {
	qf.mu.RLock()
	defer qf.mu.RUnlock()
	return qf.state.GetCoherence()
}

// GetEntropy 获取熵
func (qf *quantumField) GetEntropy() float64 {
	qf.mu.RLock()
	defer qf.mu.RUnlock()
	return qf.state.GetEntropy()
}

// Evolve 演化
func (qf *quantumField) Evolve(pattern PatternType) error {
	qf.mu.Lock()
	defer qf.mu.Unlock()

	var qPattern QuantumPattern
	// 映射PatternType到QuantumPattern
	switch pattern {
	case PatternStable:
		qPattern = "balance"
	case PatternChaos:
		qPattern = "split"
	case PatternOscillate:
		qPattern = "cycle"
	case PatternSpiral:
		qPattern = "integrate"
	default:
		qPattern = "integrate"
	}

	// 更新模式
	qf.pattern = pattern

	// 执行量子态演化
	if err := qf.state.Evolve(qPattern); err != nil {
		return fmt.Errorf("quantum evolution failed: %w", err)
	}

	// 更新相干性和熵
	qf.updateCoherenceAndEntropy()

	return nil
}

// 然后修改实现
func (qf *quantumField) Transform(target *QuantumState) error {
	if target == nil {
		return fmt.Errorf("target state cannot be nil")
	}

	qf.mu.Lock()
	defer qf.mu.Unlock()

	// 保存当前状态到缓存
	qf.cache.lastState = qf.state
	qf.cache.lastPhase = qf.state.GetPhase()
	qf.cache.lastEnergy = qf.state.GetEnergy()
	qf.cache.lastUpdate = time.Now()

	// 计算转换梯度
	phaseGrad := (target.GetPhase() - qf.state.GetPhase()) / 2
	energyGrad := (target.GetEnergy() - qf.state.GetEnergy()) / 2
	probGrad := (target.GetProbability() - qf.state.GetProbability()) / 2

	// 应用转换
	if err := qf.state.SetPhase(qf.state.GetPhase() + phaseGrad); err != nil {
		return fmt.Errorf("failed to set phase: %w", err)
	}

	if err := qf.state.SetEnergy(qf.state.GetEnergy() + energyGrad); err != nil {
		return fmt.Errorf("failed to set energy: %w", err)
	}

	if err := qf.state.SetProbability(qf.state.GetProbability() + probGrad); err != nil {
		return fmt.Errorf("failed to set probability: %w", err)
	}

	// 更新相干性和熵
	qf.updateCoherenceAndEntropy()

	return nil
}

// Entangle 量子纠缠
func (qf *quantumField) Entangle(other *QuantumState) error {
	if other == nil {
		return fmt.Errorf("other quantum state cannot be nil")
	}

	qf.mu.Lock()
	defer qf.mu.Unlock()

	// 保存当前状态到缓存
	qf.cache.lastState = qf.state
	qf.cache.lastPhase = qf.state.GetPhase()
	qf.cache.lastEnergy = qf.state.GetEnergy()
	qf.cache.lastUpdate = time.Now()

	// 计算纠缠态
	newPhase := (qf.state.GetPhase() + other.GetPhase()) / 2
	newEnergy := (qf.state.GetEnergy() + other.GetEnergy()) / 2
	newProb := math.Min(qf.state.GetCoherence(), other.GetCoherence())

	// 应用新状态
	if err := qf.state.SetPhase(newPhase); err != nil {
		return fmt.Errorf("failed to set phase: %w", err)
	}

	if err := qf.state.SetEnergy(newEnergy); err != nil {
		return fmt.Errorf("failed to set energy: %w", err)
	}

	if err := qf.state.SetProbability(newProb); err != nil {
		return fmt.Errorf("failed to set probability: %w", err)
	}

	// 增加熵
	currentEntropy := qf.state.GetEntropy()
	if err := qf.state.SetEnergy(currentEntropy + 0.1); err != nil {
		return fmt.Errorf("failed to update entropy: %w", err)
	}

	// 更新相干性和熵
	qf.updateCoherenceAndEntropy()

	return nil
}

// Decohere 退相干
func (qf *quantumField) Decohere() error {
	qf.mu.Lock()
	defer qf.mu.Unlock()

	// 保存当前状态到缓存
	qf.cache.lastState = qf.state
	qf.cache.lastPhase = qf.state.GetPhase()
	qf.cache.lastEnergy = qf.state.GetEnergy()
	qf.cache.lastUpdate = time.Now()

	// 增加熵
	currentEntropy := qf.state.GetEntropy()
	if err := qf.state.SetEnergy(currentEntropy + 0.2); err != nil {
		return fmt.Errorf("failed to update entropy: %w", err)
	}

	// 降低相干性
	currentCoherence := qf.state.GetCoherence()
	if err := qf.state.SetProbability(currentCoherence * 0.8); err != nil {
		return fmt.Errorf("failed to reduce coherence: %w", err)
	}

	// 更新相干性和熵
	qf.updateCoherenceAndEntropy()

	return nil
}

// Measure 测量
func (qf *quantumField) Measure() (float64, error) {
	qf.mu.Lock()
	defer qf.mu.Unlock()

	// 保存测量前状态到缓存
	qf.cache.lastState = qf.state
	qf.cache.lastPhase = qf.state.GetPhase()
	qf.cache.lastEnergy = qf.state.GetEnergy()
	qf.cache.lastUpdate = time.Now()

	// 测量会导致状态坍缩和退相干
	measurementValue := qf.state.GetPhase()

	// 执行退相干
	if err := qf.Decohere(); err != nil {
		return 0, fmt.Errorf("measurement decoherence failed: %w", err)
	}

	// 状态坍缩 - 概率回归到测量值对应的本征态
	if err := qf.state.Collapse(); err != nil {
		return 0, fmt.Errorf("measurement collapse failed: %w", err)
	}

	// 更新相干性和熵
	qf.updateCoherenceAndEntropy()

	return measurementValue, nil
}

// 内部辅助方法

func (qf *quantumField) evolveStable() error {
	qf.mu.Lock()
	defer qf.mu.Unlock()
	return nil
}

func (qf *quantumField) evolveChaos() error {
	qf.mu.Lock()
	defer qf.mu.Unlock()
	if err := qf.state.SetPhase(math.Sin(qf.state.GetPhase())); err != nil {
		return err
	}
	if err := qf.state.SetEnergy(qf.state.GetEnergy() * 1.1); err != nil {
		return err
	}
	return qf.state.SetEnergy(qf.state.GetEntropy() + 0.1)
}

func (qf *quantumField) evolveCycle() error {
	qf.mu.Lock()
	defer qf.mu.Unlock()
	newPhase := qf.state.GetPhase() + math.Pi/4
	if newPhase >= 2*math.Pi {
		newPhase -= 2 * math.Pi
	}
	return qf.state.SetPhase(newPhase)
}

func (qf *quantumField) evolveSpiral() error {
	qf.mu.Lock()
	defer qf.mu.Unlock()
	if err := qf.state.SetPhase(qf.state.GetPhase() + math.Pi/8); err != nil {
		return err
	}
	return qf.state.SetEnergy(qf.state.GetEnergy() * 1.05)
}

func (qf *quantumField) evolveWave() error {
	qf.mu.Lock()
	defer qf.mu.Unlock()
	if err := qf.state.SetPhase(math.Sin(qf.state.GetPhase())); err != nil {
		return err
	}
	return qf.state.SetEnergy(math.Cos(qf.state.GetEnergy()))
}

func (qf *quantumField) evolveField() error {
	qf.mu.Lock()
	defer qf.mu.Unlock()
	if err := qf.state.SetPhase(qf.state.GetPhase() + 0.1); err != nil {
		return err
	}
	return qf.state.SetEnergy(qf.state.GetEnergy() + 0.1)
}

func (qf *quantumField) evolveIntegrated() error {
	qf.mu.Lock()
	defer qf.mu.Unlock()
	if err := qf.evolveCycle(); err != nil {
		return err
	}
	if err := qf.evolveWave(); err != nil {
		return err
	}
	return qf.evolveField()
}

func (qf *quantumField) updateCoherenceAndEntropy() {
	if qf.cache.lastState == nil {
		return
	}

	timeDiff := time.Since(qf.cache.lastUpdate).Seconds()
	phaseDiff := math.Abs(qf.state.GetPhase() - qf.cache.lastPhase)
	energyDiff := math.Abs(qf.state.GetEnergy() - qf.cache.lastEnergy)

	// 相干性随时间和状态变化衰减
	decayFactor := math.Exp(-(phaseDiff + energyDiff) * timeDiff)
	currentCoherence := qf.state.GetCoherence()
	qf.state.SetPhase(currentCoherence * decayFactor)

	// 熵随时间和状态变化增加
	entropyIncrease := 0.1 * (phaseDiff + energyDiff) * timeDiff
	currentEntropy := qf.state.GetEntropy()
	qf.state.SetEnergy(currentEntropy + entropyIncrease)
}

// GetCoherence 获取量子相干性
func (qs *QuantumSystem) GetCoherence() float64 {
	qs.mu.RLock()
	defer qs.mu.RUnlock()
	return qs.coherence
}

// GetEntanglement 获取量子纠缠度
func (qs *QuantumSystem) GetEntanglement() float64 {
	qs.mu.RLock()
	defer qs.mu.RUnlock()
	return qs.entanglement
}

// GetStability 获取量子态稳定性
// 稳定性基于相位一致性和概率幅度的稳定程度
func (qs *QuantumState) GetStability() float64 {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	// 1. 相位稳定性贡献
	phaseStability := 1.0 - math.Abs(math.Sin(qs.phase))

	// 2. 概率稳定性贡献
	// 概率越接近0或1表示状态越确定
	probStability := 1.0 - 2.0*math.Abs(qs.probability-0.5)

	// 3. 能量稳定性贡献
	energyStability := math.Exp(-qs.energy / DefaultEnergy)

	// 4. 熵对稳定性的负面影响
	entropyFactor := 1.0 - qs.entropy

	// 综合计算稳定性
	stability := (phaseStability*0.3 +
		probStability*0.3 +
		energyStability*0.2 +
		entropyFactor*0.2)

	// 确保结果在[0,1]范围内
	return math.Max(0, math.Min(1, stability))
}

// GetMetrics 获取量子态指标
func (qs *QuantumState) GetMetrics() map[string]interface{} {
	qs.mu.RLock()
	defer qs.mu.RUnlock()

	return map[string]interface{}{
		"probability":  qs.probability,
		"phase":        qs.phase,
		"energy":       qs.energy,
		"entropy":      qs.entropy,
		"coherence":    qs.GetCoherence(),
		"entanglement": qs.GetEntanglement(),
		"stability":    qs.GetStability(),
	}
}

// GetState 获取量子状态自身
func (qs *QuantumState) GetState() *QuantumState {
	return qs
}
