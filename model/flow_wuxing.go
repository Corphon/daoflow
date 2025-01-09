// model/flow_wuxing.go

package model

import (
	"math"
	"sync"

	"github.com/Corphon/daoflow/core"
)

// WuXingConstants 五行常数
const (
	MaxElementEnergy = 20.0 // 每个元素最大能量
	CycleThreshold   = 0.3  // 循环阈值
	FlowRate         = 0.05 // 流动率
	ConstraintRatio  = 0.8  // 克制作用损失比例
	ConstraintCost   = 0.2  // 克制作用消耗比例
)

// Phase 相位类型
type ElementPhase uint8

const (
	PhaseNull   ElementPhase = iota // 无相位
	PhaseWaxing                     // 旺相
	PhaseStable                     // 相平
	PhaseWaning                     // 衰相
	PhaseWeak                       // 弱相
)

// PhaseConstants 相位常数
const (
	WaxingThreshold = 0.8 // 旺相阈值
	StableThreshold = 0.6 // 相平阈值
	WaningThreshold = 0.4 // 衰相阈值
	WeakThreshold   = 0.2 // 弱相阈值
)

// Element 五行元素
type Element uint8

const (
	Wood Element = iota
	Fire
	Earth
	Metal
	Water
)

// WuXingFlow 五行模型
type WuXingFlow struct {
	*BaseFlowModel // 继承基础模型

	// 五行状态 - 内部使用
	state struct {
		elements map[Element]*ElementState
		cycle    CycleType
		strength float64
	}

	// 内部组件 - 使用 core 层功能
	components struct {
		fields      map[Element]*core.Field        // 元素场
		states      map[Element]*core.QuantumState // 量子态
		network     *core.EnergyNetwork            // 能量网络
		interaction *core.Interaction              // 元素交互
	}

	mu sync.RWMutex
}

// ElementState 元素状态
type ElementState struct {
	Energy    float64
	Phase     ElementPhase
	Flow      float64
	Relations map[Element]float64
}

// CycleType 循环类型
type CycleType uint8

const (
	NoCycle           CycleType = iota
	GeneratingCycle             // 生
	ConstrainingCycle           // 克
	RebellionCycle              // 反逆
)

// NewWuXingFlow 创建五行模型
func NewWuXingFlow() *WuXingFlow {
	// 创建基础模型
	base := NewBaseFlowModel(ModelWuXing, MaxElementEnergy*5)

	// 创建五行模型
	flow := &WuXingFlow{
		BaseFlowModel: base,
	}

	// 初始化状态
	flow.state.elements = make(map[Element]*ElementState)
	flow.initializeElements()

	// 初始化组件
	flow.initializeComponents()

	return flow
}

// initializeElements 初始化元素
func (f *WuXingFlow) initializeElements() {
	elements := []Element{Wood, Fire, Earth, Metal, Water}
	for _, elem := range elements {
		f.state.elements[elem] = &ElementState{
			Energy:    MaxElementEnergy / 5,
			Phase:     PhaseNull,
			Flow:      0,
			Relations: make(map[Element]float64),
		}
	}
}

// initializeComponents 初始化组件
func (f *WuXingFlow) initializeComponents() {
	// 初始化场
	f.components.fields = make(map[Element]*core.Field)
	f.components.states = make(map[Element]*core.QuantumState)

	for elem := range f.state.elements {
		f.components.fields[elem] = core.NewField(core.ScalarField, 3)
		f.components.states[elem] = core.NewQuantumState()
	}

	// 初始化能量网络
	f.components.network = core.NewEnergyNetwork()
	f.components.interaction = core.NewInteraction()
}

// Start 启动模型
func (f *WuXingFlow) Start() error {
	if err := f.BaseFlowModel.Start(); err != nil {
		return err
	}

	// 初始化场和量子态
	return f.initializeWuXing()
}

// initializeWuXing 初始化五行
func (f *WuXingFlow) initializeWuXing() error {
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

	// 初始化能量网络
	return f.components.network.Initialize()
}

// naturalTransform 自然转换
func (f *WuXingFlow) naturalTransform() error {
	// 获取能量差异最大的相邻元素
	maxDiff := 0.0
	var source, target Element
	sequence := []Element{Wood, Fire, Earth, Metal, Water}

	// 检查相邻元素的能量差异
	for i := 0; i < len(sequence); i++ {
		current := sequence[i]
		next := sequence[(i+1)%len(sequence)]

		diff := f.state.elements[current].Energy - f.state.elements[next].Energy
		if math.Abs(diff) > math.Abs(maxDiff) {
			maxDiff = diff
			if diff > 0 {
				source = current
				target = next
			} else {
				source = next
				target = current
			}
		}
	}

	// 如果能量差异超过阈值，执行自然流动
	if math.Abs(maxDiff) > CycleThreshold {
		// 计算转换量
		transferAmount := math.Abs(maxDiff) * FlowRate

		// 执行能量转换
		if err := f.transferEnergy(source, target, transferAmount); err != nil {
			return err
		}

		// 设置循环类型
		if f.isGeneratingPair(source, target) {
			f.state.cycle = GeneratingCycle
		} else if f.isConstrainingPair(source, target) {
			f.state.cycle = ConstrainingCycle
		}
	} else {
		// 能量差异小，维持当前状态
		f.state.cycle = NoCycle
	}

	return f.updateElementStates()
}

// isGeneratingPair 检查是否为相生关系
func (f *WuXingFlow) isGeneratingPair(source, target Element) bool {
	// 相生关系：木生火、火生土、土生金、金生水、水生木
	generating := map[Element]Element{
		Wood:  Fire,
		Fire:  Earth,
		Earth: Metal,
		Metal: Water,
		Water: Wood,
	}
	return generating[source] == target
}

// isConstrainingPair 检查是否为相克关系
func (f *WuXingFlow) isConstrainingPair(source, target Element) bool {
	// 相克关系：木克土、土克水、水克火、火克金、金克木
	constraining := map[Element]Element{
		Wood:  Earth,
		Earth: Water,
		Water: Fire,
		Fire:  Metal,
		Metal: Wood,
	}
	return constraining[source] == target
}

// Transform 执行五行转换
func (f *WuXingFlow) Transform(pattern TransformPattern) error {
	if err := f.BaseFlowModel.Transform(pattern); err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	switch pattern {
	case PatternForward:
		return f.generateTransform()
	case PatternReverse:
		return f.constrainTransform()
	case PatternBalance:
		return f.balanceElements()
	default:
		return f.naturalTransform()
	}
}

// generateTransform 相生转换
func (f *WuXingFlow) generateTransform() error {
	sequence := []Element{Wood, Fire, Earth, Metal, Water}

	for i := 0; i < len(sequence); i++ {
		current := sequence[i]
		next := sequence[(i+1)%len(sequence)]

		// 计算转换量
		transferAmount := f.state.elements[current].Energy * FlowRate

		// 执行能量转换
		if err := f.transferEnergy(current, next, transferAmount); err != nil {
			return err
		}
	}

	f.state.cycle = GeneratingCycle
	return f.updateElementStates()
}

// determinePhase 确定元素相位******
func (f *WuXingFlow) determinePhase(elem Element) ElementPhase {
	state := f.state.elements[elem]
	energyRatio := state.Energy / MaxElementEnergy

	switch {
	case energyRatio >= WaxingThreshold:
		return PhaseWaxing
	case energyRatio >= StableThreshold:
		return PhaseStable
	case energyRatio >= WaningThreshold:
		return PhaseWaning
	case energyRatio >= WeakThreshold:
		return PhaseWeak
	default:
		return PhaseNull
	}
}

// applyConstraint 执行克制作用
func (f *WuXingFlow) applyConstraint(source, target Element, strength float64) error {
	// 获取源和目标元素的状态
	sourceState := f.state.elements[source]
	targetState := f.state.elements[target]

	// 计算实际克制效果
	// 克制效果与源元素能量和目标元素能量的比值有关
	ratio := sourceState.Energy / targetState.Energy
	effectiveStrength := strength * math.Min(1.0, ratio)

	// 计算能量转换
	// 克制作用会降低目标元素的能量，同时消耗源元素的部分能量
	energyLoss := effectiveStrength * 0.8 // 目标损失的能量
	energyCost := effectiveStrength * 0.2 // 源消耗的能量

	// 验证能量约束
	if targetState.Energy-energyLoss < 0 {
		energyLoss = targetState.Energy // 防止能量变为负值
	}
	if sourceState.Energy-energyCost < 0 {
		return NewModelError(ErrCodeOperation, "insufficient energy for constraint", nil)
	}

	// 执行能量调整
	targetState.Energy -= energyLoss
	sourceState.Energy -= energyCost

	// 更新量子态
	if err := f.components.states[target].SetEnergy(targetState.Energy); err != nil {
		return WrapError(err, ErrCodeOperation, "failed to update target quantum state")
	}
	if err := f.components.states[source].SetEnergy(sourceState.Energy); err != nil {
		return WrapError(err, ErrCodeOperation, "failed to update source quantum state")
	}

	// 更新场
	if err := f.components.fields[target].Update(targetState.Energy); err != nil {
		return WrapError(err, ErrCodeOperation, "failed to update target field")
	}
	if err := f.components.fields[source].Update(sourceState.Energy); err != nil {
		return WrapError(err, ErrCodeOperation, "failed to update source field")
	}

	// 更新相互作用关系
	targetState.Relations[source] -= effectiveStrength
	sourceState.Relations[target] += effectiveStrength

	// 记录能量流动
	if err := f.components.network.UpdateFlow(
		string(source),
		string(target),
		-energyLoss, // 使用负值表示克制作用
	); err != nil {
		return WrapError(err, ErrCodeOperation, "failed to update energy network")
	}

	// 更新元素相位
	sourceState.Phase = f.determinePhase(source)
	targetState.Phase = f.determinePhase(target)

	// 更新流动状态
	sourceState.Flow += effectiveStrength
	targetState.Flow -= effectiveStrength

	return nil
}

// constrainTransform 相克转换
func (f *WuXingFlow) constrainTransform() error {
	// 木克土、土克水、水克火、火克金、金克木
	constraints := map[Element]Element{
		Wood:  Earth,
		Earth: Water,
		Water: Fire,
		Fire:  Metal,
		Metal: Wood,
	}

	for source, target := range constraints {
		// 计算克制强度
		strength := f.state.elements[source].Energy * FlowRate

		// 执行克制作用
		if err := f.applyConstraint(source, target, strength); err != nil {
			return err
		}
	}

	f.state.cycle = ConstrainingCycle
	return f.updateElementStates()
}

// balanceElements 平衡元素
func (f *WuXingFlow) balanceElements() error {
	// 计算总能量
	totalEnergy := 0.0
	for _, state := range f.state.elements {
		totalEnergy += state.Energy
	}

	// 平均分配能量
	averageEnergy := totalEnergy / 5
	for elem, state := range f.state.elements {
		state.Energy = averageEnergy
		if err := f.components.states[elem].SetEnergy(averageEnergy); err != nil {
			return err
		}
	}

	f.state.cycle = NoCycle
	return f.updateElementStates()
}

// transferEnergy 转移能量
func (f *WuXingFlow) transferEnergy(from, to Element, amount float64) error {
	// 更新能量
	f.state.elements[from].Energy -= amount
	f.state.elements[to].Energy += amount

	// 更新量子态
	if err := f.components.states[from].SetEnergy(f.state.elements[from].Energy); err != nil {
		return err
	}
	if err := f.components.states[to].SetEnergy(f.state.elements[to].Energy); err != nil {
		return err
	}

	// 更新能量网络
	return f.components.network.UpdateFlow(string(from), string(to), amount)
}

// updateElementStates 更新元素状态
func (f *WuXingFlow) updateElementStates() error {
	totalStrength := 0.0

	// 更新每个元素的状态
	for elem, state := range f.state.elements {
		// 更新场
		if err := f.components.fields[elem].Update(state.Energy); err != nil {
			return err
		}

		// 更新相位
		state.Phase = f.determinePhase(elem)

		// 计算元素强度
		strength := f.calculateElementStrength(elem)
		totalStrength += strength

		// 更新相互作用
		for other := range f.state.elements {
			if elem != other {
				relation := f.calculateRelation(elem, other)
				state.Relations[other] = relation
			}
		}
	}

	// 更新整体强度
	f.state.strength = totalStrength / 5

	return nil
}

// calculateElementStrength 计算元素强度
func (f *WuXingFlow) calculateElementStrength(elem Element) float64 {
	state := f.state.elements[elem]
	field := f.components.fields[elem]

	// 综合考虑能量、场强和量子态
	energyFactor := state.Energy / MaxElementEnergy
	fieldStrength := field.GetStrength()
	quantumCoherence := f.components.states[elem].GetCoherence()

	return (energyFactor + fieldStrength + quantumCoherence) / 3
}

// calculatePhaseFactor 计算相位影响因子
func (f *WuXingFlow) calculatePhaseFactor(phase1, phase2 ElementPhase) float64 {
	// 相位匹配度计算
	if phase1 == phase2 {
		return 1.0
	}

	// 相位差异度计算
	phaseDiff := math.Abs(float64(phase1) - float64(phase2))
	return 1.0 - (phaseDiff / 4.0) // 4是最大相位差
}

// calculateRelation 计算元素关系强度
func (f *WuXingFlow) calculateRelation(elem1, elem2 Element) float64 {
	if elem1 == elem2 {
		return 0
	}

	// 获取元素状态
	state1 := f.state.elements[elem1]
	state2 := f.state.elements[elem2]

	// 基础关系强度
	baseStrength := 0.0

	// 判断关系类型
	if f.isGeneratingPair(elem1, elem2) {
		// 相生关系
		baseStrength = 0.6
	} else if f.isConstrainingPair(elem1, elem2) {
		// 相克关系
		baseStrength = -0.4
	}

	// 考虑能量水平
	energyFactor := math.Min(state1.Energy, state2.Energy) / MaxElementEnergy

	// 考虑相位影响
	phaseFactor := f.calculatePhaseFactor(state1.Phase, state2.Phase)

	// 综合计算关系强度
	return baseStrength * energyFactor * phaseFactor
}

// Close 关闭模型
func (f *WuXingFlow) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 清理组件
	for elem := range f.components.fields {
		f.components.fields[elem] = nil
		f.components.states[elem] = nil
	}
	f.components.network = nil
	f.components.interaction = nil

	return f.BaseFlowModel.Close()
}

// AdjustEnergy 调整五行能量
func (f *WuXingFlow) AdjustEnergy(delta float64) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 首先调用基类方法
	if err := f.BaseFlowModel.AdjustEnergy(delta); err != nil {
		return err
	}

	// 计算当前总能量
	totalEnergy := 0.0
	for _, state := range f.state.elements {
		totalEnergy += state.Energy
	}

	// 按比例分配增量到各元素
	for elem, state := range f.state.elements {
		ratio := state.Energy / totalEnergy
		state.Energy += delta * ratio

		// 更新量子态
		if err := f.components.states[elem].SetEnergy(state.Energy); err != nil {
			return err
		}
	}

	return f.updateElementStates()
}

// validateElement 验证元素状态
func (f *WuXingFlow) validateElement(elem Element) error {
	state, exists := f.state.elements[elem]
	if !exists {
		return NewModelError(ErrCodeInvalid, "invalid element", nil)
	}

	if state.Energy < 0 || state.Energy > MaxElementEnergy {
		return NewModelError(ErrCodeRange, "element energy out of range", nil)
	}

	return nil
}
