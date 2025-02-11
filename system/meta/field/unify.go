//system/meta/field/unify.go

package field

import (
	"math"
	"math/cmplx"
	"sync"
	"time"

	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/types"
)

// 使用现有常量
const maxHistorySize = types.MaxHistoryLength

// WuXingElement 元素数据结构
type WuXingElement struct {
	Type       string                     // 元素类型
	Energy     float64                    // 元素能量
	Properties map[string]float64         // 元素属性
	History    []model.WuXingElementState // 状态历史
	Position   struct {                   // 位置信息
		X int
		Y int
	}
}

// UnifiedField 统一场结构
type UnifiedField struct {
	mu sync.RWMutex

	// 使用model层定义的CoreState
	core model.CoreState

	// 复用model层的五行模型
	wuxing *model.WuXingFlow

	// 复用model层的阴阳模型
	yinyang *model.YinYangFlow

	// 场组件
	components struct {
		scalar  *FieldTensor
		vector  *FieldTensor
		metric  *FieldTensor
		quantum *FieldTensor
	}

	// 统一特性(meta层特有的高层抽象)
	properties struct {
		symmetry   string             // 对称性类型
		invariants []float64          // 不变量
		topology   FieldTopology      // 拓扑结构
		dimension  int                // 维度
		Properties map[string]float64 // 动态属性映射
	}

	// 场耦合关系
	couplings map[string]*FieldCoupling

	// 添加元素管理
	WuXingElements []*WuXingElement // 五行元素集合

	// 添加状态字段
	state struct {
		History  []UnifiedState // 状态历史记录
		Strength float64        // 当前场强度
		Phase    float64        // 当前相位
		Energy   float64
	}
}

// UnifiedState 统一场状态
type UnifiedState struct {
	Time           time.Time
	Energy         float64
	Symmetry       string
	YinYang        YinYangState
	WuXingElements model.WuXingElementState
	Metrics        UnifiedMetrics
}

// YinYangState 阴阳状态
type YinYangState struct {
	Yin     float64
	Yang    float64
	Harmony float64
	Phase   model.Phase
}

// UnifiedMetrics 统一度量
type UnifiedMetrics struct {
	Strength  float64 // 场强度
	Coherence float64 // 相干性
	Stability float64 // 稳定性
	Harmony   float64 // 和谐度
	Phase     float64 // 相位
}

// FieldTopology 场拓扑结构
type FieldTopology struct {
	Dimension    int
	Connectivity float64
	Curvature    float64
	Holes        int
	Genus        int
}

// 常量定义
const (
	generatingFactor        = 0.3
	controllingFactor       = 0.2
	minWuXingElementEnergy  = 0.0
	maxWuXingElementEnergy  = 1.0
	maxWuXingElementHistory = 1000
	evolutionTimeStep       = time.Second / 100
)

// 结构体定义
type WuXingElementInteraction struct {
	energyDelta float64
	influences  map[string]float64
}

// WuXingElementRelation 元素关系
type WuXingElementRelation struct {
	factor       float64 // 关系因子
	relationType string  // 关系类型
}

type FieldPattern struct {
	Type           string
	Strength       float64
	WuXingElements []string
	Properties     map[string]float64
}

// PredictedState 预测状态
type PredictedState struct {
	Time           time.Time          // 预测时间
	Energy         float64            // 能量值
	Phase          float64            // 相位
	Strength       float64            // 耦合强度
	WuXingElements map[string]float64 // 元素状态
	Properties     map[string]float64 // 属性集
}

type EvolutionPrediction struct {
	StartTime  time.Time
	Duration   time.Duration
	States     []PredictedState
	Confidence float64
}

// NewUnifiedField 创建新的统一场
func NewUnifiedField(dimension int) (*UnifiedField, error) {
	if dimension < 3 {
		return nil, model.WrapError(nil, model.ErrCodeValidation,
			"dimension must be at least 3")
	}

	uf := &UnifiedField{
		couplings: make(map[string]*FieldCoupling),
	}

	// 初始化场组件
	if err := uf.initComponents(dimension); err != nil {
		return nil, err
	}

	// 初始化统一特性
	uf.initProperties(dimension)

	// 初始化五行属性
	uf.initWuXingElements()

	// 初始化阴阳属性
	uf.initYinYang()

	// 初始化Properties
	uf.properties.Properties = make(map[string]float64)

	return uf, nil
}

// initComponents 初始化场组件
func (uf *UnifiedField) initComponents(dimension int) error {
	// 创建标量场
	scalar := NewFieldTensor(dimension, 1)
	uf.components.scalar = scalar

	// 创建向量场
	vector := NewFieldTensor(dimension, 2)
	uf.components.vector = vector

	// 创建度规场
	metric := NewFieldTensor(dimension, 2)
	uf.components.metric = metric

	// 创建量子场
	quantum := NewFieldTensor(dimension, 1)
	uf.components.quantum = quantum

	return nil
}

// initProperties 初始化统一特性
func (uf *UnifiedField) initProperties(dimension int) {
	uf.properties.dimension = dimension
	uf.properties.symmetry = "undefined"
	uf.properties.invariants = make([]float64, 0)

	// 初始化拓扑结构
	uf.properties.topology = FieldTopology{
		Dimension:    dimension,
		Connectivity: 1.0,
		Curvature:    0.0,
		Holes:        0,
		Genus:        0,
	}
}

// initWuXingElements 初始化五行属性
func (uf *UnifiedField) initWuXingElements() {
	// 使用model层的WuXingFlow初始化
	uf.wuxing = model.NewWuXingFlow()
}

// initYinYang 初始化阴阳属性
func (uf *UnifiedField) initYinYang() {
	// 使用model层的YinYangFlow初始化
	uf.yinyang = model.NewYinYangFlow()
}

// Evolve 演化统一场
func (uf *UnifiedField) Evolve() error {
	uf.mu.Lock()
	defer uf.mu.Unlock()

	// 更新场组件
	if err := uf.evolveComponents(); err != nil {
		return err
	}

	// 更新耦合关系
	if err := uf.evolveCouplings(); err != nil {
		return err
	}

	// 更新五行属性
	uf.evolveWuXingElements()

	// 更新阴阳属性
	uf.evolveYinYang()

	// 计算新的统一状态
	state := uf.calculateUnifiedState()

	// 记录状态
	uf.recordState(state)

	return nil
}

// evolveYinYang 演化阴阳属性
func (uf *UnifiedField) evolveYinYang() {
	// 计算阴阳比例
	yinRatio := uf.yinyang.GetState().Energy / (uf.yinyang.GetState().Energy + uf.wuxing.GetState().Energy)

	// 根据阴阳比例选择转换模式
	var pattern model.TransformPattern
	switch {
	case yinRatio > 0.618: // 黄金分割比
		pattern = model.PatternReverse // 阳转阴
	case yinRatio < 0.382: // 1 - 0.618
		pattern = model.PatternForward // 阴转阳
	default:
		pattern = model.PatternBalance // 保持平衡
	}

	// 执行转换
	uf.yinyang.Transform(pattern)
}

// calculateUnifiedState 计算统一状态
func (uf *UnifiedField) calculateUnifiedState() UnifiedState {
	return UnifiedState{
		Time:     time.Now(),
		Energy:   uf.core.EnergyState.GetTotalEnergy(),
		Symmetry: uf.properties.symmetry,
		YinYang: YinYangState{
			Yin:     uf.yinyang.GetState().YinEnergy,
			Yang:    uf.yinyang.GetState().YangEnergy,
			Harmony: uf.yinyang.GetState().Harmony,
			Phase:   uf.yinyang.GetState().Phase,
		},
		WuXingElements: model.WuXingElementState{
			Metal:   uf.wuxing.GetWuXingElementEnergy("Metal"),
			Wood:    uf.wuxing.GetWuXingElementEnergy("Wood"),
			Water:   uf.wuxing.GetWuXingElementEnergy("Water"),
			Fire:    uf.wuxing.GetWuXingElementEnergy("Fire"),
			Earth:   uf.wuxing.GetWuXingElementEnergy("Earth"),
			Balance: uf.wuxing.GetState().Balance,
		},
		Metrics: UnifiedMetrics{
			Strength: uf.core.FieldState.GetStrength(),
		},
	}
}

// recordState 记录统一状态
func (uf *UnifiedField) recordState(state UnifiedState) {
	uf.state.History = append(uf.state.History, state)

	// 限制历史记录长度
	if len(uf.state.History) > maxHistorySize {
		uf.state.History = uf.state.History[1:]
	}
}

// evolveComponents 演化场组件
func (uf *UnifiedField) evolveComponents() error {
	// 演化标量场
	if err := uf.evolveScalarField(); err != nil {
		return err
	}

	// 演化向量场
	if err := uf.evolveVectorField(); err != nil {
		return err
	}

	// 演化度规场
	if err := uf.evolveMetricField(); err != nil {
		return err
	}

	// 演化量子场
	if err := uf.evolveQuantumField(); err != nil {
		return err
	}

	return nil
}

// 辅助函数获取时间步长的秒数
func getEvolutionStepSeconds(step time.Duration) float64 {
	return step.Seconds()
}

// evolveScalarField 演化标量场
func (uf *UnifiedField) evolveScalarField() error {
	// 应用标量场方程
	field := uf.components.scalar

	// 计算拉普拉斯算子
	laplacian := calculateLaplacian(field)

	stepSeconds := getEvolutionStepSeconds(evolutionTimeStep)
	// 应用波动方程
	for i := 0; i < field.dimension; i++ {
		for j := 0; j < field.dimension; j++ {
			value, _ := field.GetComponent([]int{i, j})
			// 将float64转换为complex128
			term := complex(laplacian[i][j]*stepSeconds, 0)
			newValue := value + term
			field.SetComponent([]int{i, j}, newValue)
		}
	}

	return nil
}

// evolveVectorField 演化向量场
func (uf *UnifiedField) evolveVectorField() error {
	field := uf.components.vector
	curl := calculateCurl(field)
	div := calculateDivergence(field)
	stepSeconds := getEvolutionStepSeconds(evolutionTimeStep)

	for i := 0; i < field.dimension; i++ {
		for j := 0; j < field.dimension; j++ {
			value, _ := field.GetComponent([]int{i, j})
			// 将float64转为complex128
			term := complex((curl[i][j]-div[i][j])*stepSeconds, 0)
			newValue := value + term
			field.SetComponent([]int{i, j}, newValue)
		}
	}
	return nil
}

// evolveMetricField 演化度规场
func (uf *UnifiedField) evolveMetricField() error {
	field := uf.components.metric
	ricci := calculateRicciTensor(field)
	energyMomentum := calculateEnergyMomentumTensor(field)
	stepSeconds := getEvolutionStepSeconds(evolutionTimeStep)

	for i := 0; i < field.dimension; i++ {
		for j := 0; j < field.dimension; j++ {
			value, _ := field.GetComponent([]int{i, j})
			// 将float64转为complex128
			term := complex((ricci[i][j]-0.5*energyMomentum[i][j])*stepSeconds, 0)
			newValue := value + term
			field.SetComponent([]int{i, j}, newValue)
		}
	}
	return nil
}

// evolveQuantumField 演化量子场
func (uf *UnifiedField) evolveQuantumField() error {
	field := uf.components.quantum
	hamiltonian := calculateHamiltonian(field)
	stepSeconds := getEvolutionStepSeconds(evolutionTimeStep)

	for i := 0; i < field.dimension; i++ {
		for j := 0; j < field.dimension; j++ {
			value, _ := field.GetComponent([]int{i, j})
			// 使用复数直接计算,使用转换后的秒数
			newValue := value + hamiltonian[i][j]*complex(0, -stepSeconds)
			field.SetComponent([]int{i, j}, newValue)
		}
	}
	return nil
}

// evolveCouplings 演化耦合关系
func (uf *UnifiedField) evolveCouplings() error {
	for _, coupling := range uf.couplings {
		// 使用正确的方法名Update
		if err := coupling.Update(); err != nil {
			return err
		}
	}
	return nil
}

// evolveWuXingElements 演化五行元素
func (uf *UnifiedField) evolveWuXingElements() {
	uf.mu.Lock()
	defer uf.mu.Unlock()

	// 更新各元素状态
	for i, WuXingElement := range uf.WuXingElements {
		// 计算元素间相互作用
		interactions := uf.calculateWuXingElementInteractions(i)

		// 更新元素能量
		WuXingElement.Energy += interactions.energyDelta

		// 更新元素属性
		WuXingElement.Properties = uf.updateWuXingElementProperties(WuXingElement, interactions)

		// 应用五行相生相克规则
		uf.applyWuXingRules(WuXingElement, interactions)

		// 记录状态变化
		uf.recordWuXingElementState(WuXingElement)
	}

	// 更新整体场态
	uf.updateFieldState()
}

// calculateWuXingElementInteractions 计算元素间相互作用
func (uf *UnifiedField) calculateWuXingElementInteractions(index int) WuXingElementInteraction {
	WuXingElement := uf.WuXingElements[index]
	interaction := WuXingElementInteraction{
		energyDelta: 0,
		influences:  make(map[string]float64),
	}

	for i, other := range uf.WuXingElements {
		if i == index {
			continue
		}

		// 计算相生相克关系
		relation := getWuXingRelation(WuXingElement.Type, other.Type)

		// 计算相互作用强度
		strength := calculateInteractionStrength(WuXingElement, other)

		// 根据五行关系调整能量交换
		energyExchange := strength * relation.factor

		interaction.energyDelta += energyExchange
		interaction.influences[other.Type] = energyExchange
	}

	return interaction
}

// updateWuXingElementProperties 更新元素属性
func (uf *UnifiedField) updateWuXingElementProperties(
	WuXingElement *WuXingElement,
	interaction WuXingElementInteraction) map[string]float64 {

	newProps := make(map[string]float64)

	// 基础属性更新
	for key, value := range WuXingElement.Properties {
		// 考虑相互作用影响
		influence := interaction.influences[key]

		// 应用量子效应
		quantumFactor := uf.calculateQuantumFactor(WuXingElement)

		// 更新属性值
		newValue := value + influence*quantumFactor

		// 确保属性值在有效范围内
		newProps[key] = clamp(newValue, 0, 1)
	}

	return newProps
}

// calculateQuantumFactor 计算量子影响因子
func (uf *UnifiedField) calculateQuantumFactor(WuXingElement *WuXingElement) float64 {
	// 获取场量子态
	quantumField := uf.components.quantum

	// 基于元素位置获取量子场值
	position := WuXingElement.Position // 添加Position字段到WuXingElement结构
	value, err := quantumField.GetComponent([]int{position.X, position.Y})
	if err != nil {
		return 1.0 // 出错时返回中性因子
	}

	// 计算量子效应
	// 1. 量子相干性影响
	coherence := cmplx.Abs(value)

	// 2. 相位影响
	phase := cmplx.Phase(value)
	phaseEffect := math.Cos(phase)

	// 3. 考虑元素能量对量子效应的影响
	energyFactor := WuXingElement.Energy / model.MaxWuXingElementEnergy

	// 4. 综合量子效应
	quantumEffect := (coherence + math.Abs(phaseEffect) + energyFactor) / 3.0

	// 归一化到[0.5, 1.5]范围
	return 0.5 + quantumEffect
}

// applyWuXingRules 应用五行规则
func (uf *UnifiedField) applyWuXingRules(WuXingElement *WuXingElement, interaction WuXingElementInteraction) {
	// 直接使用五行模型的方法
	WuXingElementEnergy := uf.wuxing.GetWuXingElementEnergy(WuXingElement.Type)

	// 应用相生规则 - 使用model层定义的常量
	for _, rel := range model.GeneratingWuXingElements(WuXingElement.Type) {
		if influence, ok := interaction.influences[rel]; ok {
			WuXingElementEnergy += influence * model.FlowRate
		}
	}

	// 应用相克规则
	for _, rel := range model.ConstrainingWuXingElements(WuXingElement.Type) {
		if influence, ok := interaction.influences[rel]; ok {
			WuXingElementEnergy -= influence * model.ConstraintRatio
		}
	}

	// 确保能量守恒 - 使用model层定义的常量
	WuXingElement.Energy = math.Max(0, math.Min(WuXingElementEnergy, model.MaxWuXingElementEnergy))
}

// recordWuXingElementState 记录元素状态
func (uf *UnifiedField) recordWuXingElementState(WuXingElement *WuXingElement) {
	state := model.WuXingElementState{
		Timestamp:  time.Now(),
		Type:       WuXingElement.Type,
		Energy:     WuXingElement.Energy,
		Properties: WuXingElement.Properties,
	}

	WuXingElement.History = append(WuXingElement.History, state)

	// 限制历史记录长度
	if len(WuXingElement.History) > maxWuXingElementHistory {
		WuXingElement.History = WuXingElement.History[1:]
	}
}

// updateFieldState 更新场状态
func (uf *UnifiedField) updateFieldState() {
	// 计算总能量
	totalEnergy := 0.0
	for _, WuXingElement := range uf.WuXingElements {
		totalEnergy += WuXingElement.Energy
	}

	// 更新场强度
	uf.state.Strength = totalEnergy / float64(len(uf.WuXingElements))

	// 更新场相位
	uf.state.Phase = uf.calculateFieldPhase()

	// 更新量子特性
	uf.updateQuantumProperties()

	// 记录场状态
	uf.recordFieldState()
}

// 添加计算场相位的方法
func (uf *UnifiedField) calculateFieldPhase() float64 {
	// 从量子场组件获取相位
	value, err := uf.components.quantum.GetComponent([]int{0, 0})
	if err != nil {
		return 0
	}
	return cmplx.Phase(value)
}

// 添加更新量子特性的方法
func (uf *UnifiedField) updateQuantumProperties() {
	// 更新量子相干性
	coherence := 0.0
	for i := 0; i < uf.components.quantum.dimension; i++ {
		for j := 0; j < uf.components.quantum.dimension; j++ {
			value, _ := uf.components.quantum.GetComponent([]int{i, j})
			coherence += cmplx.Abs(value)
		}
	}
	uf.components.quantum.quantum.coherence = coherence / float64(uf.components.quantum.dimension*uf.components.quantum.dimension)
}

// 添加记录场状态的方法
func (uf *UnifiedField) recordFieldState() {
	state := UnifiedState{
		Time:     time.Now(),
		Energy:   uf.core.EnergyState.GetTotalEnergy(),
		Symmetry: uf.properties.symmetry,
		YinYang: YinYangState{
			Yin:     uf.yinyang.GetState().YinEnergy,
			Yang:    uf.yinyang.GetState().YangEnergy,
			Harmony: uf.yinyang.GetState().Harmony,
			Phase:   uf.yinyang.GetState().Phase,
		},
		WuXingElements: model.WuXingElementState{
			Metal:   uf.wuxing.GetWuXingElementEnergy("Metal"),
			Wood:    uf.wuxing.GetWuXingElementEnergy("Wood"),
			Water:   uf.wuxing.GetWuXingElementEnergy("Water"),
			Fire:    uf.wuxing.GetWuXingElementEnergy("Fire"),
			Earth:   uf.wuxing.GetWuXingElementEnergy("Earth"),
			Balance: uf.wuxing.GetState().Balance,
		},
		Metrics: UnifiedMetrics{
			Strength: uf.state.Strength,
			Phase:    uf.state.Phase,
		},
	}

	uf.state.History = append(uf.state.History, state)
	if len(uf.state.History) > maxHistorySize {
		uf.state.History = uf.state.History[1:]
	}
}

// AnalyzePatterns 分析场模式
func (uf *UnifiedField) AnalyzePatterns() []FieldPattern {
	uf.mu.RLock()
	defer uf.mu.RUnlock()

	patterns := make([]FieldPattern, 0)

	// 分析元素组合模式
	WuXingElementPatterns := uf.detectWuXingElementPatterns()
	patterns = append(patterns, WuXingElementPatterns...)

	// 分析能量分布模式
	energyPatterns := uf.detectEnergyPatterns()
	patterns = append(patterns, energyPatterns...)

	// 分析量子态模式
	quantumPatterns := uf.detectQuantumPatterns()
	patterns = append(patterns, quantumPatterns...)

	return patterns
}

// detectWuXingElementPatterns 检测元素组合模式
func (uf *UnifiedField) detectWuXingElementPatterns() []FieldPattern {
	patterns := make([]FieldPattern, 0)

	// 分析元素组合
	for i, elem1 := range uf.WuXingElements {
		for j := i + 1; j < len(uf.WuXingElements); j++ {
			elem2 := uf.WuXingElements[j]
			// 检查元素间相互作用
			interaction := uf.calculateWuXingElementInteractions(j)
			if interaction.energyDelta > 0.1 {
				pattern := FieldPattern{
					Type:           "WuXingElement_interaction",
					Strength:       interaction.energyDelta,
					WuXingElements: []string{elem1.Type, elem2.Type},
					Properties: map[string]float64{
						"influence": interaction.influences[elem2.Type],
					},
				}
				patterns = append(patterns, pattern)
			}
		}
	}
	return patterns
}

// detectEnergyPatterns 检测能量分布模式
func (uf *UnifiedField) detectEnergyPatterns() []FieldPattern {
	patterns := make([]FieldPattern, 0)

	// 计算总能量和平均能量
	totalEnergy := 0.0
	for _, WuXingElement := range uf.WuXingElements {
		totalEnergy += WuXingElement.Energy
	}
	avgEnergy := totalEnergy / float64(len(uf.WuXingElements))

	// 检测能量聚集
	for _, WuXingElement := range uf.WuXingElements {
		if WuXingElement.Energy > avgEnergy*1.5 {
			pattern := FieldPattern{
				Type:           "energy_concentration",
				Strength:       WuXingElement.Energy / avgEnergy,
				WuXingElements: []string{WuXingElement.Type},
				Properties: map[string]float64{
					"energy_ratio": WuXingElement.Energy / totalEnergy,
				},
			}
			patterns = append(patterns, pattern)
		}
	}
	return patterns
}

// detectQuantumPatterns 检测量子态模式
func (uf *UnifiedField) detectQuantumPatterns() []FieldPattern {
	patterns := make([]FieldPattern, 0)

	// 获取量子场值
	value, err := uf.components.quantum.GetComponent([]int{0, 0})
	if err != nil {
		return patterns
	}

	// 计算量子特性
	coherence := cmplx.Abs(value)
	phase := cmplx.Phase(value)

	// 检测量子相干模式
	if coherence > 0.8 {
		pattern := FieldPattern{
			Type:           "quantum_coherence",
			Strength:       coherence,
			WuXingElements: []string{"quantum_field"},
			Properties: map[string]float64{
				"phase":     phase,
				"amplitude": coherence,
			},
		}
		patterns = append(patterns, pattern)
	}
	return patterns
}

// PredictEvolution 预测场演化
func (uf *UnifiedField) PredictEvolution(duration time.Duration) (*EvolutionPrediction, error) {
	uf.mu.RLock()
	defer uf.mu.RUnlock()

	if len(uf.state.History) < minDataPoints {
		return nil, model.WrapError(nil, model.ErrCodeValidation,
			"insufficient historical data for prediction")
	}

	prediction := &EvolutionPrediction{
		StartTime: time.Now(),
		Duration:  duration,
		States:    make([]PredictedState, 0),
	}

	// 使用时间序列分析预测未来状态
	currentState := uf.getCurrentState()
	timeSteps := int(duration / evolutionTimeStep)

	for i := 0; i < timeSteps; i++ {
		nextUnifiedState := uf.predictNextState(currentState)
		nextPredictedState := uf.toPredictedState(nextUnifiedState)
		prediction.States = append(prediction.States, nextPredictedState)
		currentState = nextUnifiedState
	}

	// 计算预测可信度
	prediction.Confidence = uf.calculatePredictionConfidence(prediction.States)

	return prediction, nil
}

// toPredictedState 将UnifiedState转换为PredictedState
func (uf *UnifiedField) toPredictedState(us UnifiedState) PredictedState {
	return PredictedState{
		Time:   us.Time,
		Energy: us.Energy,
		Phase:  float64(us.YinYang.Phase), // 将Phase转换为float64
		WuXingElements: map[string]float64{
			"Metal": us.WuXingElements.Metal,
			"Wood":  us.WuXingElements.Wood,
			"Water": us.WuXingElements.Water,
			"Fire":  us.WuXingElements.Fire,
			"Earth": us.WuXingElements.Earth,
		},
	}
}

// 辅助函数

func getWuXingRelation(type1, type2 string) WuXingElementRelation {
	relations := map[string]map[string]WuXingElementRelation{
		"Wood": {
			"Fire":  {factor: generatingFactor, relationType: "generating"},
			"Earth": {factor: controllingFactor, relationType: "controlling"},
		},
		"Fire": {
			"Earth": {factor: generatingFactor, relationType: "generating"},
			"Metal": {factor: controllingFactor, relationType: "controlling"},
		},
		"Earth": {
			"Metal": {factor: generatingFactor, relationType: "generating"},
			"Water": {factor: controllingFactor, relationType: "controlling"},
		},
		"Metal": {
			"Water": {factor: generatingFactor, relationType: "generating"},
			"Wood":  {factor: controllingFactor, relationType: "controlling"},
		},
		"Water": {
			"Wood": {factor: generatingFactor, relationType: "generating"},
			"Fire": {factor: controllingFactor, relationType: "controlling"},
		},
	}

	if rel, ok := relations[type1][type2]; ok {
		return rel
	}
	return WuXingElementRelation{factor: 0, relationType: "neutral"}
}

func calculateInteractionStrength(e1, e2 *WuXingElement) float64 {
	// 基础相互作用强度
	baseStrength := math.Sqrt(e1.Energy * e2.Energy)

	// 距离因子
	distance := calculateWuXingElementDistance(e1, e2)
	distanceFactor := 1.0 / (1.0 + distance)

	// 属性相似度
	similarity := calculatePropertySimilarity(e1.Properties, e2.Properties)

	return baseStrength * distanceFactor * similarity
}

func calculateWuXingElementDistance(e1, e2 *WuXingElement) float64 {
	// 简化的空间距离计算
	return math.Abs(e1.Energy - e2.Energy)
}

func calculatePropertySimilarity(props1, props2 map[string]float64) float64 {
	if len(props1) == 0 || len(props2) == 0 {
		return 0
	}

	similarity := 0.0
	count := 0

	for key, value1 := range props1 {
		if value2, ok := props2[key]; ok {
			similarity += 1.0 - math.Abs(value1-value2)
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return similarity / float64(count)
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// getCurrentState 获取当前统一场状态
func (uf *UnifiedField) getCurrentState() UnifiedState {
	return uf.calculateUnifiedState()
}

// predictNextState 预测下一个状态,返回UnifiedState
func (uf *UnifiedField) predictNextState(current UnifiedState) UnifiedState {
	// 使用时间序列分析预测
	energyTrend := uf.calculateEnergyTrend()
	fieldTrend := uf.calculateFieldTrend()

	return UnifiedState{
		Time:     current.Time.Add(evolutionTimeStep),
		Energy:   current.Energy * (1 + energyTrend),
		Symmetry: current.Symmetry,
		YinYang: YinYangState{
			Yin:     current.YinYang.Yin * (1 + energyTrend),
			Yang:    current.YinYang.Yang * (1 + energyTrend),
			Harmony: current.YinYang.Harmony,
			Phase:   model.FromFloat64(current.YinYang.Phase.ToFloat64() + fieldTrend),
		},
		WuXingElements: model.WuXingElementState{
			Metal:   current.WuXingElements.Metal * (1 + energyTrend),
			Wood:    current.WuXingElements.Wood * (1 + energyTrend),
			Water:   current.WuXingElements.Water * (1 + energyTrend),
			Fire:    current.WuXingElements.Fire * (1 + energyTrend),
			Earth:   current.WuXingElements.Earth * (1 + energyTrend),
			Balance: current.WuXingElements.Balance,
		},
		Metrics: current.Metrics,
	}
}

// calculatePredictionConfidence 计算预测可信度
func (uf *UnifiedField) calculatePredictionConfidence(states []PredictedState) float64 {
	if len(states) == 0 {
		return 0
	}

	// 基于历史数据计算预测准确度
	var totalError float64
	historyLen := len(uf.state.History)
	if historyLen < 2 {
		return 0.5 // 默认置信度
	}

	// 计算历史预测误差
	for i := 1; i < historyLen; i++ {
		predicted := uf.state.History[i-1].Energy
		actual := uf.state.History[i].Energy
		error := math.Abs(predicted-actual) / actual
		totalError += error
	}

	// 转换为置信度
	confidence := 1 - (totalError / float64(historyLen-1))
	return math.Max(0, math.Min(1, confidence))
}

// 辅助函数

func (uf *UnifiedField) calculateEnergyTrend() float64 {
	if len(uf.state.History) < 2 {
		return 0
	}
	latest := uf.state.History[len(uf.state.History)-1]
	previous := uf.state.History[len(uf.state.History)-2]
	return (latest.Energy - previous.Energy) / previous.Energy
}

func (uf *UnifiedField) calculateFieldTrend() float64 {
	if len(uf.state.History) < 2 {
		return 0
	}
	latest := uf.state.History[len(uf.state.History)-1]
	previous := uf.state.History[len(uf.state.History)-2]
	return latest.YinYang.Phase.ToFloat64() - previous.YinYang.Phase.ToFloat64()
}

// GetPropertyValue 获取属性值
func (uf *UnifiedField) GetPropertyValue(name string) (float64, bool) {
	uf.mu.RLock()
	defer uf.mu.RUnlock()

	// 查找基本属性
	switch name {
	case "energy":
		return uf.state.Energy, true
	case "strength":
		return uf.state.Strength, true
	case "phase":
		return uf.state.Phase, true
	}

	// 查找统一特性
	if value, exists := uf.properties.Properties[name]; exists {
		return value, true
	}

	return 0, false
}

// SetPropertyValue 设置属性值
func (uf *UnifiedField) SetPropertyValue(name string, value float64) error {
	uf.mu.Lock()
	defer uf.mu.Unlock()

	// 更新基本属性
	switch name {
	case "energy":
		uf.state.Energy = value
		return nil
	case "strength":
		uf.state.Strength = value
		return nil
	case "phase":
		uf.state.Phase = value
		return nil
	}

	// 更新统一特性
	if uf.properties.Properties == nil {
		uf.properties.Properties = make(map[string]float64)
	}
	uf.properties.Properties[name] = value
	return nil
}

// 辅助方法
func (uf *UnifiedField) calculateStability() float64 {
	// 基于能量波动计算稳定性
	uf.mu.RLock()
	defer uf.mu.RUnlock()

	// 1. 获取能量历史
	if len(uf.state.History) < 2 {
		return 1.0 // 历史记录不足时返回最大稳定性
	}

	// 2. 计算能量波动
	energyVariance := 0.0
	meanEnergy := 0.0

	// 计算平均能量
	for _, state := range uf.state.History {
		meanEnergy += state.Energy
	}
	meanEnergy /= float64(len(uf.state.History))

	// 计算方差
	for _, state := range uf.state.History {
		diff := state.Energy - meanEnergy
		energyVariance += diff * diff
	}
	energyVariance /= float64(len(uf.state.History))

	// 3. 计算稳定性指数 (方差越小越稳定)
	stabilityIndex := 1.0 / (1.0 + math.Sqrt(energyVariance))

	// 4. 考虑量子相干性影响
	coherence := uf.core.FieldState.GetCoherence()

	// 5. 综合计算稳定性 (范围0-1)
	stability := (stabilityIndex + coherence) / 2.0

	return math.Max(0.0, math.Min(1.0, stability))
}

// GetEnergy 获取场的总能量
func (uf *UnifiedField) GetEnergy() float64 {
	uf.mu.RLock()
	defer uf.mu.RUnlock()

	// 获取当前统一状态
	currentState := uf.calculateUnifiedState()

	// 返回总能量
	return currentState.Energy
}

// GetState 替代GetPropertyValue获取状态
func (uf *UnifiedField) GetState() (*model.FieldState, error) {
	uf.mu.RLock()
	defer uf.mu.RUnlock()

	// 获取当前统一状态
	currentState := uf.calculateUnifiedState()

	state := &model.FieldState{
		Energy:   currentState.Energy,
		Elements: make([]*model.WuXingElement, 0),
		Properties: map[string]float64{
			"strength": currentState.Metrics.Strength,
			"harmony":  currentState.YinYang.Harmony,
			"balance":  currentState.WuXingElements.Balance,
		},
		Timestamp: currentState.Time,
	}

	return state, nil
}

var (
	defaultField     *UnifiedField
	defaultFieldOnce sync.Once
)

// GetDefaultField 获取默认统一场实例
func GetDefaultField() *UnifiedField {
	defaultFieldOnce.Do(func() {
		// 创建3维统一场
		field, err := NewUnifiedField(3)
		if err != nil {
			panic(err) // 初始化失败直接panic
		}
		defaultField = field
	})
	return defaultField
}

func init() {
	// 确保默认统一场初始化
	GetDefaultField()
}

// GetStability 获取场的稳定性
func (uf *UnifiedField) GetStability() float64 {
	uf.mu.RLock()
	defer uf.mu.RUnlock()
	return uf.calculateStability()
}

// GetCoherence 获取场的相干性
func (uf *UnifiedField) GetCoherence() float64 {
	uf.mu.RLock()
	defer uf.mu.RUnlock()

	if value, exists := uf.properties.Properties["coherence"]; exists {
		return value
	}

	// 计算相干度
	coherence := 0.0
	if qs := uf.components.quantum; qs != nil {
		coherence = qs.GetCoherence()
	}
	return coherence
}
