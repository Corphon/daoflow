// model/flow_bagua.go

package model

import (
	"math"
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
)

// BaGuaConstants 八卦常数
const (
	MaxTrigramEnergy = 12.5 // 每个卦象最大能量
	ChangeThreshold  = 0.2  // 变化阈值
	ResonanceRate    = 0.08 // 共振率
)

// Trigram 卦象
type Trigram uint8

const (
	Qian Trigram = iota // 乾 ☰
	Dui                 // 兑 ☱
	Li                  // 离 ☲
	Zhen                // 震 ☳
	Xun                 // 巽 ☴
	Kan                 // 坎 ☵
	Gen                 // 艮 ☶
	Kun                 // 坤 ☷
)

// BaGuaFlow 八卦模型
type BaGuaFlow struct {
	*BaseFlowModel // 继承基础模型

	// 八卦状态 - 内部使用
	state struct {
		trigrams  map[Trigram]*TrigramState
		resonance float64
		harmony   float64
		changes   []Change
	}

	// 内部组件 - 使用 core 层功能
	components struct {
		fields     map[Trigram]*core.Field        // 卦象场
		states     map[Trigram]*core.QuantumState // 量子态
		resonator  *core.Resonator                // 共振器
		correlator *core.Correlator               // 关联器
	}

	mu sync.RWMutex
}

// TrigramState 卦象状态
type TrigramState struct {
	Energy     float64
	Lines      [3]bool // 三爻状态
	Resonance  float64
	Relations  map[Trigram]float64
	LastChange time.Time
}

// Change 变化记录
type Change struct {
	From      Trigram
	To        Trigram
	Type      ChangeType
	Timestamp time.Time
}

// ChangeType 变化类型
type ChangeType uint8

const (
	NoChange       ChangeType = iota
	NaturalChange             // 自然变化
	ForcedChange              // 强制变化
	ResonantChange            // 共振变化
)

// NewBaGuaFlow 创建八卦模型
func NewBaGuaFlow() *BaGuaFlow {
	// 创建基础模型
	base := NewBaseFlowModel(ModelBaGua, MaxTrigramEnergy*8)

	// 创建八卦模型
	flow := &BaGuaFlow{
		BaseFlowModel: base,
	}

	// 初始化状态
	flow.state.trigrams = make(map[Trigram]*TrigramState)
	flow.initializeTrigrams()

	// 初始化组件
	flow.initializeComponents()

	return flow
}

// initializeTrigrams 初始化卦象
func (f *BaGuaFlow) initializeTrigrams() {
	trigrams := []Trigram{Qian, Dui, Li, Zhen, Xun, Kan, Gen, Kun}
	for _, tri := range trigrams {
		f.state.trigrams[tri] = &TrigramState{
			Energy:     MaxTrigramEnergy / 8,
			Lines:      [3]bool{},
			Resonance:  0,
			Relations:  make(map[Trigram]float64),
			LastChange: time.Now(),
		}
	}
}

// initializeComponents 初始化组件
func (f *BaGuaFlow) initializeComponents() {
	// 初始化场和量子态
	f.components.fields = make(map[Trigram]*core.Field)
	f.components.states = make(map[Trigram]*core.QuantumState)

	for tri := range f.state.trigrams {
		f.components.fields[tri] = core.NewField(core.ScalarField, 3)
		f.components.states[tri] = core.NewQuantumState()
	}

	// 初始化共振器和关联器
	f.components.resonator = core.NewResonator()
	f.components.correlator = core.NewCorrelator()
}

// Start 启动模型
func (f *BaGuaFlow) Start() error {
	if err := f.BaseFlowModel.Start(); err != nil {
		return err
	}

	return f.initializeBaGua()
}

// initializeBaGua 初始化八卦
func (f *BaGuaFlow) initializeBaGua() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 初始化场
	for _, field := range f.components.fields {
		if err := field.Initialize(); err != nil {
			return WrapError(err, ErrCodeOperation, "failed to initialize field")
		}
	}

	// 初始化量子态
	for _, state := range f.components.states {
		if err := state.Initialize(); err != nil {
			return WrapError(err, ErrCodeOperation, "failed to initialize quantum state")
		}
	}

	// 初始化共振器
	return f.components.resonator.Initialize()
}

// Transform 执行八卦转换
func (f *BaGuaFlow) Transform(pattern TransformPattern) error {
	if err := f.BaseFlowModel.Transform(pattern); err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	switch pattern {
	case PatternBalance:
		return f.balanceTrigrams()
	case PatternForward:
		return f.naturalTransform()
	case PatternReverse:
		return f.resonantTransform()
	default:
		return f.adaptiveTransform()
	}
}

// balanceTrigrams 平衡卦象
func (f *BaGuaFlow) balanceTrigrams() error {
	// 计算总能量
	totalEnergy := 0.0
	for _, state := range f.state.trigrams {
		totalEnergy += state.Energy
	}

	// 平均分配能量
	averageEnergy := totalEnergy / 8
	for tri, state := range f.state.trigrams {
		state.Energy = averageEnergy
		if err := f.components.states[tri].SetEnergy(averageEnergy); err != nil {
			return err
		}
	}

	return f.updateTrigramStates()
}

// findNaturalChange 寻找自然变化卦
func (f *BaGuaFlow) findNaturalChange(tri Trigram) Trigram {
	// 八卦自然变化规律：
	// 乾(☰) <-> 坤(☷)  阳极生阴，阴极生阳
	// 兑(☱) <-> 艮(☶)  泽山相辅
	// 离(☲) <-> 坎(☵)  火水相济
	// 震(☳) <-> 巽(☴)  雷风相荡
	naturalChangePairs := map[Trigram]Trigram{
		Qian: Kun,  // 乾 -> 坤
		Kun:  Qian, // 坤 -> 乾
		Dui:  Gen,  // 兑 -> 艮
		Gen:  Dui,  // 艮 -> 兑
		Li:   Kan,  // 离 -> 坎
		Kan:  Li,   // 坎 -> 离
		Zhen: Xun,  // 震 -> 巽
		Xun:  Zhen, // 巽 -> 震
	}

	// 获取能量状态
	currentState := f.state.trigrams[tri]
	targetState := f.state.trigrams[naturalChangePairs[tri]]

	// 检查能量差异是否足够触发变化
	energyDiff := math.Abs(currentState.Energy - targetState.Energy)
	if energyDiff > ChangeThreshold*MaxTrigramEnergy {
		return naturalChangePairs[tri]
	}

	return tri // 如果不满足变化条件，返回原卦
}

// changeTrigram 改变卦象
func (f *BaGuaFlow) changeTrigram(from, to Trigram, changeType ChangeType) error {
	fromState := f.state.trigrams[from]
	toState := f.state.trigrams[to]

	// 计算能量转换
	transferEnergy := fromState.Energy * 0.5 // 转移一半能量

	// 更新能量
	fromState.Energy -= transferEnergy
	toState.Energy += transferEnergy

	// 更新量子态
	if err := f.components.states[from].SetEnergy(fromState.Energy); err != nil {
		return err
	}
	if err := f.components.states[to].SetEnergy(toState.Energy); err != nil {
		return err
	}

	// 更新场
	if err := f.components.fields[from].Update(fromState.Energy); err != nil {
		return err
	}
	if err := f.components.fields[to].Update(toState.Energy); err != nil {
		return err
	}

	// 记录变化
	change := Change{
		From:      from,
		To:        to,
		Type:      changeType,
		Timestamp: time.Now(),
	}
	f.state.changes = append(f.state.changes, change)

	// 更新最后变化时间
	fromState.LastChange = time.Now()
	toState.LastChange = time.Now()

	// 更新关联关系
	fromState.Relations[to] += ResonanceRate
	toState.Relations[from] += ResonanceRate

	// 更新共振状态
	if err := f.components.resonator.Update(); err != nil {
		return err
	}

	// 更新关联状态
	if err := f.components.correlator.Update(); err != nil {
		return err
	}

	return nil
}

// naturalTransform 自然变化
func (f *BaGuaFlow) naturalTransform() error {
	for tri, state := range f.state.trigrams {
		// 检查能量阈值
		if state.Energy > ChangeThreshold*MaxTrigramEnergy {
			// 寻找相应变化卦
			targetTri := f.findNaturalChange(tri)
			if targetTri != tri {
				if err := f.changeTrigram(tri, targetTri, NaturalChange); err != nil {
					return err
				}
			}
		}
	}
	return f.updateTrigramStates()
}

// findResonatingPairs finds pairs of trigrams that are in resonance
func (f *BaGuaFlow) findResonatingPairs() [][2]Trigram {
	pairs := make([][2]Trigram, 0)
	f.mu.RLock()
	defer f.mu.RUnlock()

	// Get all possible trigram combinations
	trigrams := []Trigram{Qian, Dui, Li, Zhen, Xun, Kan, Gen, Kun}

	// Check each pair of trigrams for resonance
	for i := 0; i < len(trigrams); i++ {
		for j := i + 1; j < len(trigrams); j++ {
			tri1 := trigrams[i]
			tri2 := trigrams[j]

			// Check resonance conditions:
			// 1. Both trigrams must have sufficient energy
			// 2. Their resonance relationship must be above the threshold
			state1 := f.state.trigrams[tri1]
			state2 := f.state.trigrams[tri2]

			if state1.Energy >= ChangeThreshold*MaxTrigramEnergy &&
				state2.Energy >= ChangeThreshold*MaxTrigramEnergy &&
				state1.Relations[tri2] >= ResonanceRate {
				pairs = append(pairs, [2]Trigram{tri1, tri2})
			}
		}
	}

	return pairs
}

// applyResonance applies resonance effects between two trigrams
func (f *BaGuaFlow) applyResonance(tri1, tri2 Trigram) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	state1 := f.state.trigrams[tri1]
	state2 := f.state.trigrams[tri2]

	// Calculate resonance energy transfer
	resonanceEnergy := ResonanceRate * math.Min(state1.Energy, state2.Energy)

	// Update quantum states through the resonator
	if err := f.components.resonator.ApplyResonance(
		f.components.states[tri1],
		f.components.states[tri2],
		resonanceEnergy,
	); err != nil {
		return err
	}

	// Update trigram states
	now := time.Now()
	state1.Resonance += resonanceEnergy
	state2.Resonance += resonanceEnergy
	state1.LastChange = now
	state2.LastChange = now

	// Record the resonant change
	f.state.changes = append(f.state.changes, Change{
		From:      tri1,
		To:        tri2,
		Type:      ResonantChange,
		Timestamp: now,
	})

	return nil
}

// resonantTransform 共振变化
func (f *BaGuaFlow) resonantTransform() error {
	// 更新共振器
	if err := f.components.resonator.Update(); err != nil {
		return err
	}

	// 检查共振条件
	resonatingPairs := f.findResonatingPairs()
	for _, pair := range resonatingPairs {
		if err := f.applyResonance(pair[0], pair[1]); err != nil {
			return err
		}
	}

	return f.updateTrigramStates()
}

// calculateSystemEntropy 计算系统熵
func (f *BaGuaFlow) calculateSystemEntropy() float64 {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var totalEntropy float64
	var totalEnergy float64

	// 收集所有卦象的能量和熵
	for _, state := range f.state.trigrams {
		totalEnergy += state.Energy
	}

	if totalEnergy == 0 {
		return 0
	}

	// 计算每个卦象的归一化熵贡献
	for _, state := range f.state.trigrams {
		if state.Energy > 0 {
			// 计算能量占比
			p := state.Energy / totalEnergy
			// 使用Shannon熵公式
			totalEntropy -= p * math.Log2(p)
		}
	}

	// 归一化熵值到[0,1]范围
	// 对于8个卦象，最大熵值为log2(8)=3
	return totalEntropy / 3.0
}

// calculateHarmony 计算系统和谐度
func (f *BaGuaFlow) calculateHarmony() float64 {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var totalHarmony float64
	var connections float64

	// 检查每对卦象之间的关系
	trigrams := []Trigram{Qian, Dui, Li, Zhen, Xun, Kan, Gen, Kun}

	for i := 0; i < len(trigrams); i++ {
		for j := i + 1; j < len(trigrams); j++ {
			tri1 := trigrams[i]
			tri2 := trigrams[j]

			// 获取两个卦象的状态
			state1 := f.state.trigrams[tri1]
			state2 := f.state.trigrams[tri2]

			// 计算关系强度
			relationStrength := state1.Relations[tri2]
			if relationStrength > 0 {
				// 考虑能量平衡
				energyBalance := 1.0 - math.Abs(state1.Energy-state2.Energy)/MaxTrigramEnergy
				// 考虑共振状态
				resonanceMatch := math.Min(state1.Resonance, state2.Resonance)

				// 综合计算和谐度
				harmony := relationStrength * energyBalance * resonanceMatch
				totalHarmony += harmony
				connections++
			}
		}
	}

	// 如果没有连接，返回0
	if connections == 0 {
		return 0
	}

	// 归一化和谐度到[0,1]范围
	return totalHarmony / connections
}

// adaptiveTransform 适应性变化
func (f *BaGuaFlow) adaptiveTransform() error {
	// 计算系统熵
	entropy := f.calculateSystemEntropy()

	// 根据熵值调整变化
	if entropy > ChangeThreshold {
		return f.naturalTransform()
	}
	return f.resonantTransform()
}

// updateTrigramStates 更新卦象状态
func (f *BaGuaFlow) updateTrigramStates() error {
	// 更新量子态
	for tri, state := range f.components.states {
		if err := state.Update(); err != nil {
			return err
		}
		// 更新卦象能量
		f.state.trigrams[tri].Energy = state.GetEnergy()
	}

	// 更新共振度
	f.state.resonance = f.components.resonator.GetResonance()

	// 更新和谐度
	f.state.harmony = f.calculateHarmony()

	return nil
}

// Close 关闭模型
func (f *BaGuaFlow) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 清理组件
	for tri := range f.components.fields {
		f.components.fields[tri] = nil
		f.components.states[tri] = nil
	}
	f.components.resonator = nil
	f.components.correlator = nil

	return f.BaseFlowModel.Close()
}

// AdjustEnergy 调整八卦能量
func (f *BaGuaFlow) AdjustEnergy(delta float64) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 首先调用基类方法
	if err := f.BaseFlowModel.AdjustEnergy(delta); err != nil {
		return err
	}

	// 计算当前总能量并按比例分配
	totalEnergy := 0.0
	for _, state := range f.state.trigrams {
		totalEnergy += state.Energy
	}

	for tri, state := range f.state.trigrams {
		ratio := state.Energy / totalEnergy
		state.Energy += delta * ratio

		// 更新量子态
		if err := f.components.states[tri].SetEnergy(state.Energy); err != nil {
			return err
		}
	}

	return f.updateTrigramStates()
}
