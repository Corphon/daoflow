// model/flow_wuxing.go

package model

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
)

var (
	defaultFlow *WuXingFlow
	onceLock    sync.Once
)

// WuXingConstants 五行常数
const (
	MaxWuXingElementEnergy = 20.0 // 每个元素最大能量
	CycleThreshold         = 0.3  // 循环阈值
	FlowRate               = 0.05 // 流动率
	ConstraintRatio        = 0.8  // 克制作用损失比例
	ConstraintCost         = 0.2  // 克制作用消耗比例
	GeneratingFactor       = 1.2  // 相生系数
	ConstrainingFactor     = 0.8  // 相克系数
)

// Phase 相位类型
type WuXingElementPhase uint8

const (
	PhaseNull   WuXingElementPhase = iota // 无相位
	PhaseWaxing                           // 旺相
	PhaseStable                           // 相平
	PhaseWaning                           // 衰相
	PhaseWeak                             // 弱相

)

// PhaseConstants 相位常数
const (
	WaxingThreshold = 0.8 // 旺相阈值
	StableThreshold = 0.6 // 相平阈值
	WaningThreshold = 0.4 // 衰相阈值
	WeakThreshold   = 0.2 // 弱相阈值
)

// WuXingElement 五行元素
type WuXingElement uint8

const (
	Wood WuXingElement = iota
	Fire
	Earth
	Metal
	Water
)

// WuXingElementState 五行元素状态结构
type WuXingElementState struct {
	// 基础属性
	Metal   float64 // 金
	Wood    float64 // 木
	Water   float64 // 水
	Fire    float64 // 火
	Earth   float64 // 土
	Balance float64 // 平衡度

	// 状态信息
	Type      string                    // 元素类型
	Energy    float64                   // 能量值
	Phase     WuXingElementPhase        // 相位
	Flow      float64                   // 流动性
	Relations map[WuXingElement]float64 // 关系网络

	// 扩展属性
	Properties map[string]float64 // 属性集
	Timestamp  time.Time          // 时间戳
}

// ---------------------------------------------
// NewWuXingElementState 创建新的五行元素状态
func NewWuXingElementState() *WuXingElementState {
	return &WuXingElementState{
		Relations:  make(map[WuXingElement]float64),
		Properties: make(map[string]float64),
		Timestamp:  time.Now(),
	}
}

// WuXingFlow 五行模型
type WuXingFlow struct {
	*BaseFlowModel // 继承基础模型

	// 五行状态 - 内部使用
	state struct {
		WuXingElements map[WuXingElement]*WuXingElementState
		cycle          CycleType
		strength       float64
	}

	// 内部组件 - 使用 core 层功能
	components struct {
		fields      map[WuXingElement]*core.Field        // 元素场
		states      map[WuXingElement]*core.QuantumState // 量子态
		network     *core.EnergyNetwork                  // 能量网络
		interaction *core.Interaction                    // 元素交互
	}

	mu sync.RWMutex
}

// CycleType 循环类型
type CycleType uint8

const (
	NoCycle           CycleType = iota
	GeneratingCycle             // 生
	ConstrainingCycle           // 克
	RebellionCycle              // 反逆
)

// WuXingElementRelation 五行关系
type WuXingElementRelation struct {
	Factor       float64
	RelationType string
}

// GetWuXingRelation 获取两个元素间的五行关系
func GetWuXingRelation(type1, type2 string) WuXingElementRelation {
	relations := map[string]map[string]WuXingElementRelation{
		"Wood": {
			"Fire":  {Factor: GeneratingFactor, RelationType: "generating"},
			"Earth": {Factor: ConstrainingFactor, RelationType: "controlling"},
		},
		"Fire": {
			"Earth": {Factor: GeneratingFactor, RelationType: "generating"},
			"Metal": {Factor: ConstrainingFactor, RelationType: "controlling"},
		},
		"Earth": {
			"Metal": {Factor: GeneratingFactor, RelationType: "generating"},
			"Water": {Factor: ConstrainingFactor, RelationType: "controlling"},
		},
		"Metal": {
			"Water": {Factor: GeneratingFactor, RelationType: "generating"},
			"Wood":  {Factor: ConstrainingFactor, RelationType: "controlling"},
		},
		"Water": {
			"Wood": {Factor: GeneratingFactor, RelationType: "generating"},
			"Fire": {Factor: ConstrainingFactor, RelationType: "controlling"},
		},
	}

	if rel, ok := relations[type1][type2]; ok {
		return rel
	}
	return WuXingElementRelation{Factor: 0, RelationType: "neutral"}
}

// NewWuXingFlow 创建五行模型
func NewWuXingFlow() *WuXingFlow {
	// 创建基础模型
	base := NewBaseFlowModel(ModelWuXing, MaxWuXingElementEnergy*5)

	// 创建五行模型
	flow := &WuXingFlow{
		BaseFlowModel: base,
	}

	// 初始化状态
	flow.state.WuXingElements = make(map[WuXingElement]*WuXingElementState)
	flow.initializeWuXingElements()

	// 初始化组件
	flow.initializeComponents()

	return flow
}

// initializeWuXingElements 初始化元素
func (f *WuXingFlow) initializeWuXingElements() {
	WuXingElements := []WuXingElement{Wood, Fire, Earth, Metal, Water}
	for _, elem := range WuXingElements {
		f.state.WuXingElements[elem] = &WuXingElementState{
			Energy:    MaxWuXingElementEnergy / 5,
			Phase:     PhaseNull,
			Flow:      0,
			Relations: make(map[WuXingElement]float64),
		}
	}
}

// initializeComponents 初始化组件
func (f *WuXingFlow) initializeComponents() {
	// 初始化场
	f.components.fields = make(map[WuXingElement]*core.Field)
	f.components.states = make(map[WuXingElement]*core.QuantumState)

	for elem := range f.state.WuXingElements {
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
	var source, target WuXingElement
	sequence := []WuXingElement{Wood, Fire, Earth, Metal, Water}

	// 检查相邻元素的能量差异
	for i := 0; i < len(sequence); i++ {
		current := sequence[i]
		next := sequence[(i+1)%len(sequence)]

		diff := f.state.WuXingElements[current].Energy - f.state.WuXingElements[next].Energy
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

	return f.updateWuXingElementStates()
}

// isGeneratingPair 检查是否为相生关系
func (f *WuXingFlow) isGeneratingPair(source, target WuXingElement) bool {
	// 相生关系：木生火、火生土、土生金、金生水、水生木
	generating := map[WuXingElement]WuXingElement{
		Wood:  Fire,
		Fire:  Earth,
		Earth: Metal,
		Metal: Water,
		Water: Wood,
	}
	return generating[source] == target
}

// isConstrainingPair 检查是否为相克关系
func (f *WuXingFlow) isConstrainingPair(source, target WuXingElement) bool {
	// 相克关系：木克土、土克水、水克火、火克金、金克木
	constraining := map[WuXingElement]WuXingElement{
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
		return f.balanceWuXingElements()
	default:
		return f.naturalTransform()
	}
}

// generateTransform 相生转换
func (f *WuXingFlow) generateTransform() error {
	sequence := []WuXingElement{Wood, Fire, Earth, Metal, Water}

	for i := 0; i < len(sequence); i++ {
		current := sequence[i]
		next := sequence[(i+1)%len(sequence)]

		// 计算转换量
		transferAmount := f.state.WuXingElements[current].Energy * FlowRate

		// 执行能量转换
		if err := f.transferEnergy(current, next, transferAmount); err != nil {
			return err
		}
	}

	f.state.cycle = GeneratingCycle
	return f.updateWuXingElementStates()
}

// determinePhase 确定元素相位******
func (f *WuXingFlow) determinePhase(elem WuXingElement) WuXingElementPhase {
	state := f.state.WuXingElements[elem]
	energyRatio := state.Energy / MaxWuXingElementEnergy

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
func (f *WuXingFlow) applyConstraint(source, target WuXingElement, strength float64) error {
	// 获取源和目标元素的状态
	sourceState := f.state.WuXingElements[source]
	targetState := f.state.WuXingElements[target]

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
	constraints := map[WuXingElement]WuXingElement{
		Wood:  Earth,
		Earth: Water,
		Water: Fire,
		Fire:  Metal,
		Metal: Wood,
	}

	for source, target := range constraints {
		// 计算克制强度
		strength := f.state.WuXingElements[source].Energy * FlowRate

		// 执行克制作用
		if err := f.applyConstraint(source, target, strength); err != nil {
			return err
		}
	}

	f.state.cycle = ConstrainingCycle
	return f.updateWuXingElementStates()
}

// balanceWuXingElements 平衡元素
func (f *WuXingFlow) balanceWuXingElements() error {
	// 计算总能量
	totalEnergy := 0.0
	for _, state := range f.state.WuXingElements {
		totalEnergy += state.Energy
	}

	// 平均分配能量
	averageEnergy := totalEnergy / 5
	for elem, state := range f.state.WuXingElements {
		state.Energy = averageEnergy
		if err := f.components.states[elem].SetEnergy(averageEnergy); err != nil {
			return err
		}
	}

	f.state.cycle = NoCycle
	return f.updateWuXingElementStates()
}

// transferEnergy 转移能量
func (f *WuXingFlow) transferEnergy(from, to WuXingElement, amount float64) error {
	// 更新能量
	f.state.WuXingElements[from].Energy -= amount
	f.state.WuXingElements[to].Energy += amount

	// 更新量子态
	if err := f.components.states[from].SetEnergy(f.state.WuXingElements[from].Energy); err != nil {
		return err
	}
	if err := f.components.states[to].SetEnergy(f.state.WuXingElements[to].Energy); err != nil {
		return err
	}

	// 更新能量网络
	return f.components.network.UpdateFlow(string(from), string(to), amount)
}

// updateWuXingElementStates 更新元素状态
func (f *WuXingFlow) updateWuXingElementStates() error {
	totalStrength := 0.0

	// 更新每个元素的状态
	for elem, state := range f.state.WuXingElements {
		// 更新场
		if err := f.components.fields[elem].Update(state.Energy); err != nil {
			return err
		}

		// 更新相位
		state.Phase = f.determinePhase(elem)

		// 计算元素强度
		strength := f.calculateWuXingElementStrength(elem)
		totalStrength += strength

		// 更新相互作用
		for other := range f.state.WuXingElements {
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

// calculateWuXingElementStrength 计算元素强度
func (f *WuXingFlow) calculateWuXingElementStrength(elem WuXingElement) float64 {
	state := f.state.WuXingElements[elem]
	field := f.components.fields[elem]

	// 综合考虑能量、场强和量子态
	energyFactor := state.Energy / MaxWuXingElementEnergy
	fieldStrength := field.GetStrength()
	quantumCoherence := f.components.states[elem].GetCoherence()

	return (energyFactor + fieldStrength + quantumCoherence) / 3
}

// calculatePhaseFactor 计算相位影响因子
func (f *WuXingFlow) calculatePhaseFactor(phase1, phase2 WuXingElementPhase) float64 {
	// 相位匹配度计算
	if phase1 == phase2 {
		return 1.0
	}

	// 相位差异度计算
	phaseDiff := math.Abs(float64(phase1) - float64(phase2))
	return 1.0 - (phaseDiff / 4.0) // 4是最大相位差
}

// calculateRelation 计算元素关系强度
func (f *WuXingFlow) calculateRelation(elem1, elem2 WuXingElement) float64 {
	if elem1 == elem2 {
		return 0
	}

	// 获取元素状态
	state1 := f.state.WuXingElements[elem1]
	state2 := f.state.WuXingElements[elem2]

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
	energyFactor := math.Min(state1.Energy, state2.Energy) / MaxWuXingElementEnergy

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
	for _, state := range f.state.WuXingElements {
		totalEnergy += state.Energy
	}

	// 按比例分配增量到各元素
	for elem, state := range f.state.WuXingElements {
		ratio := state.Energy / totalEnergy
		state.Energy += delta * ratio

		// 更新量子态
		if err := f.components.states[elem].SetEnergy(state.Energy); err != nil {
			return err
		}
	}

	return f.updateWuXingElementStates()
}

// validateWuXingElement 验证元素状态
func (f *WuXingFlow) validateWuXingElement(elem WuXingElement) error {
	state, exists := f.state.WuXingElements[elem]
	if !exists {
		return NewModelError(ErrCodeInvalid, "invalid WuXingElement", nil)
	}

	if state.Energy < 0 || state.Energy > MaxWuXingElementEnergy {
		return NewModelError(ErrCodeRange, "WuXingElement energy out of range", nil)
	}

	return nil
}

// WuXingElementFromString 从字符串转换为WuXingElement类型
func WuXingElementFromString(s string) (WuXingElement, bool) {
	switch s {
	case "Wood":
		return Wood, true
	case "Fire":
		return Fire, true
	case "Earth":
		return Earth, true
	case "Metal":
		return Metal, true
	case "Water":
		return Water, true
	default:
		return Wood, false
	}
}

// GetWuXingElementEnergy 获取元素能量
func (f *WuXingFlow) GetWuXingElementEnergy(WuXingElement string) float64 {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if elem, ok := WuXingElementFromString(WuXingElement); ok {
		if state, exists := f.state.WuXingElements[elem]; exists {
			return state.Energy
		}
	}
	return 0
}

// GetBalance 获取五行平衡度
func (f *WuXingFlow) GetBalance() float64 {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var total float64
	for _, state := range f.state.WuXingElements {
		total += state.Energy
	}

	// 计算能量分布的均匀度
	balance := 1.0
	expected := total / float64(len(f.state.WuXingElements))
	for _, state := range f.state.WuXingElements {
		diff := math.Abs(state.Energy - expected)
		balance *= (1 - diff/total)
	}

	return balance
}

// GeneratingWuXingElements 获取生成关系的元素
func GeneratingWuXingElements(WuXingElementType string) []string {
	// 五行相生关系
	generating := map[string][]string{
		"Wood":  {"Fire"},  // 木生火
		"Fire":  {"Earth"}, // 火生土
		"Earth": {"Metal"}, // 土生金
		"Metal": {"Water"}, // 金生水
		"Water": {"Wood"},  // 水生木
	}

	if WuXingElements, exists := generating[WuXingElementType]; exists {
		return WuXingElements
	}
	return []string{}
}

// ConstrainingWuXingElements 获取克制关系的元素
func ConstrainingWuXingElements(WuXingElementType string) []string {
	// 五行相克关系
	constraining := map[string][]string{
		"Wood":  {"Earth"}, // 木克土
		"Earth": {"Water"}, // 土克水
		"Water": {"Fire"},  // 水克火
		"Fire":  {"Metal"}, // 火克金
		"Metal": {"Wood"},  // 金克木
	}

	if WuXingElements, exists := constraining[WuXingElementType]; exists {
		return WuXingElements
	}
	return []string{}
}

// String 返回五行元素的字符串表示
func (we WuXingElement) String() string {
	switch we {
	case Wood:
		return "Wood"
	case Fire:
		return "Fire"
	case Earth:
		return "Earth"
	case Metal:
		return "Metal"
	case Water:
		return "Water"
	default:
		return "Unknown"
	}
}

// GetEnergy 获取五行元素能量
func (we *WuXingElement) GetEnergy() float64 {
	if flow := defaultWuXingFlow(); flow != nil {
		return flow.GetWuXingElementEnergy(we.String())
	}
	return 0
}

// GetProperties 获取五行元素属性
func (we *WuXingElement) GetProperties() map[string]float64 {
	properties := make(map[string]float64)

	// 基本属性
	properties["phase"] = float64(defaultWuXingFlow().determinePhase(*we))
	properties["flow"] = defaultWuXingFlow().state.WuXingElements[*we].Flow

	// 关系属性
	for other := Wood; other <= Water; other++ {
		if *we != other {
			relation := defaultWuXingFlow().calculateRelation(*we, other)
			properties[fmt.Sprintf("relation_%s", other.String())] = relation
		}
	}

	return properties
}

// defaultWuXingFlow 获取默认五行模型单例
func defaultWuXingFlow() *WuXingFlow {
	// 单例模式获取或创建WuXingFlow实例
	onceLock.Do(func() {
		defaultFlow = NewWuXingFlow()
	})
	return defaultFlow
}

// GetType 获取五行元素类型
func (we *WuXingElement) GetType() string {
	return we.String() // 复用已有的String()方法
}
