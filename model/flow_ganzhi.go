// model/flow_ganzhi.go

package model

import (
	"math"
	"sync"

	"github.com/Corphon/daoflow/core"
)

// GanZhiConstants 干支常数
const (
	MaxStemEnergy   = 10.0 // 天干最大能量
	MaxBranchEnergy = 10.0 // 地支最大能量
	CycleLength     = 60   // 六十甲子周期
)

// HeavenlyStem 天干
type HeavenlyStem uint8

const (
	Jia  HeavenlyStem = iota // 甲
	Yi                       // 乙
	Bing                     // 丙
	Ding                     // 丁
	Wu                       // 戊
	Ji                       // 己
	Geng                     // 庚
	Xin                      // 辛
	Ren                      // 壬
	Gui                      // 癸
)

// EarthlyBranch 地支
type EarthlyBranch uint8

const (
	Zi        EarthlyBranch = iota // 子
	Chou                           // 丑
	Yin                            // 寅
	Mao                            // 卯
	Chen                           // 辰
	Si                             // 巳
	Wu_Branch                      // 午
	Wei                            // 未
	Shen                           // 申
	You                            // 酉
	Xu                             // 戌
	Hai                            // 亥
)

// GanZhiFlow 干支模型
type GanZhiFlow struct {
	*BaseFlowModel // 继承基础模型

	// 干支状态 - 内部使用
	state struct {
		stems    map[HeavenlyStem]*StemState
		branches map[EarthlyBranch]*BranchState
		cycle    int
		harmony  float64
	}

	// 内部组件 - 使用 core 层功能
	components struct {
		stemFields   map[HeavenlyStem]*core.Field         // 天干场
		branchFields map[EarthlyBranch]*core.Field        // 地支场
		stemStates   map[HeavenlyStem]*core.QuantumState  // 天干量子态
		branchStates map[EarthlyBranch]*core.QuantumState // 地支量子态
		cycleManager *core.CycleManager                   // 周期管理器
		harmonizer   *core.Harmonizer                     // 和谐器
	}

	mu sync.RWMutex
}

// StemState 天干状态
type StemState struct {
	Energy    float64
	Phase     Phase
	Element   Element // 关联五行
	Polarity  Nature  // 阴阳属性
	Relations map[EarthlyBranch]float64
}

// BranchState 地支状态
type BranchState struct {
	Energy    float64
	Phase     Phase
	Element   Element // 关联五行
	Polarity  Nature  // 阴阳属性
	Relations map[HeavenlyStem]float64
}

// NewGanZhiFlow 创建干支模型
func NewGanZhiFlow() *GanZhiFlow {
	// 创建基础模型
	base := NewBaseFlowModel(ModelGanZhi, (MaxStemEnergy*10 + MaxBranchEnergy*12))

	// 创建干支模型
	flow := &GanZhiFlow{
		BaseFlowModel: base,
	}

	// 初始化状态
	flow.initializeStates()

	// 初始化组件
	flow.initializeComponents()

	return flow
}

// initializeStates 初始化状态
func (f *GanZhiFlow) initializeStates() {
	// 初始化天干状态
	f.state.stems = make(map[HeavenlyStem]*StemState)
	for i := Jia; i <= Gui; i++ {
		f.state.stems[i] = &StemState{
			Energy:    MaxStemEnergy / 10,
			Phase:     PhaseNone,
			Element:   f.getStemElement(i),
			Polarity:  f.getStemPolarity(i),
			Relations: make(map[EarthlyBranch]float64),
		}
	}

	// 初始化地支状态
	f.state.branches = make(map[EarthlyBranch]*BranchState)
	for i := Zi; i <= Hai; i++ {
		f.state.branches[i] = &BranchState{
			Energy:    MaxBranchEnergy / 12,
			Phase:     PhaseNone,
			Element:   f.getBranchElement(i),
			Polarity:  f.getBranchPolarity(i),
			Relations: make(map[HeavenlyStem]float64),
		}
	}
}

// initializeComponents 初始化组件
func (f *GanZhiFlow) initializeComponents() {
	// 初始化场
	f.components.stemFields = make(map[HeavenlyStem]*core.Field)
	f.components.branchFields = make(map[EarthlyBranch]*core.Field)

	// 初始化量子态
	f.components.stemStates = make(map[HeavenlyStem]*core.QuantumState)
	f.components.branchStates = make(map[EarthlyBranch]*core.QuantumState)

	// 初始化天干组件
	for stem := range f.state.stems {
		f.components.stemFields[stem] = core.NewField(core.ScalarField, 3)
		f.components.stemStates[stem] = core.NewQuantumState()
	}

	// 初始化地支组件
	for branch := range f.state.branches {
		f.components.branchFields[branch] = core.NewField(core.ScalarField, 3)
		f.components.branchStates[branch] = core.NewQuantumState()
	}

	// 初始化周期管理器和和谐器
	f.components.cycleManager = core.NewCycleManager(CycleLength)
	f.components.harmonizer = core.NewHarmonizer()
}

// Transform 执行干支转换
func (f *GanZhiFlow) Transform(pattern TransformPattern) error {
	if err := f.BaseFlowModel.Transform(pattern); err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	switch pattern {
	case PatternForward:
		return f.cyclicTransform()
	case PatternBalance:
		return f.harmonizeElements()
	default:
		return f.naturalTransform()
	}
}

// getCurrentGanZhi 获取当前干支组合
func (f *GanZhiFlow) getCurrentGanZhi() (HeavenlyStem, EarthlyBranch) {
	// 从当前周期计算天干和地支
	stemIndex := f.state.cycle % 10   // 天干循环：10
	branchIndex := f.state.cycle % 12 // 地支循环：12

	return HeavenlyStem(stemIndex), EarthlyBranch(branchIndex)
}

// updateEnergies 更新干支能量状态
func (f *GanZhiFlow) updateEnergies(stem HeavenlyStem, branch EarthlyBranch) error {
	// 获取当前状态
	stemState := f.state.stems[stem]
	branchState := f.state.branches[branch]

	// 计算能量交互
	stemEnergy := stemState.Energy
	branchEnergy := branchState.Energy

	// 基于五行相生相克关系调整能量
	elementInteraction := f.calculateElementInteraction(stemState.Element, branchState.Element)

	// 基于阴阳属性调整能量
	polarityFactor := f.calculatePolarityFactor(stemState.Polarity, branchState.Polarity)

	// 计算新能量
	newStemEnergy := stemEnergy * (1 + elementInteraction*polarityFactor)
	newBranchEnergy := branchEnergy * (1 + elementInteraction*polarityFactor)

	// 限制能量范围
	newStemEnergy = math.Min(MaxStemEnergy, math.Max(0, newStemEnergy))
	newBranchEnergy = math.Min(MaxBranchEnergy, math.Max(0, newBranchEnergy))

	// 更新状态
	stemState.Energy = newStemEnergy
	branchState.Energy = newBranchEnergy

	return nil
}

// updateQuantumStates 更新量子态
func (f *GanZhiFlow) updateQuantumStates(stem HeavenlyStem, branch EarthlyBranch) error {
	// 更新天干量子态
	if err := f.components.stemStates[stem].SetEnergy(f.state.stems[stem].Energy); err != nil {
		return err
	}

	// 更新地支量子态
	if err := f.components.branchStates[branch].SetEnergy(f.state.branches[branch].Energy); err != nil {
		return err
	}

	// 更新相位关系
	stemPhase := float64(stem) * 2 * math.Pi / 10     // 天干相位
	branchPhase := float64(branch) * 2 * math.Pi / 12 // 地支相位

	// 设置量子态相位
	if err := f.components.stemStates[stem].SetPhase(stemPhase); err != nil {
		return err
	}
	if err := f.components.branchStates[branch].SetPhase(branchPhase); err != nil {
		return err
	}

	return nil
}

// updateHarmony 更新和谐度
func (f *GanZhiFlow) updateHarmony() error {
	// 重置和谐器组件
	f.components.harmonizer.Initialize()

	// 计算五行和谐度
	elementHarmony := f.calculateElementHarmony()
	f.components.harmonizer.UpdateComponent("element", elementHarmony)

	// 计算阴阳和谐度
	polarityHarmony := f.calculatePolarityHarmony()
	f.components.harmonizer.UpdateComponent("polarity", polarityHarmony)

	// 计算能量和谐度
	energyHarmony := f.calculateEnergyHarmony()
	f.components.harmonizer.UpdateComponent("energy", energyHarmony)

	// 获取总体和谐度
	f.state.harmony = f.components.harmonizer.GetHarmony()

	return nil
}

// 辅助方法

// calculateElementInteraction 计算五行相互作用
func (f *GanZhiFlow) calculateElementInteraction(elem1, elem2 Element) float64 {
	// 五行相生关系：+0.2
	// 五行相克关系：-0.1
	// 五行相同：+0.1
	// 其他情况：0

	switch {
	case isGenerating(elem1, elem2):
		return 0.2
	case isControlling(elem1, elem2):
		return -0.1
	case elem1 == elem2:
		return 0.1
	default:
		return 0
	}
}

// calculatePolarityFactor 计算阴阳因子
func (f *GanZhiFlow) calculatePolarityFactor(n1, n2 Nature) float64 {
	if n1 == n2 {
		return -0.1 // 同性相斥
	}
	return 0.1 // 异性相吸
}

// calculateElementHarmony 计算五行和谐度
func (f *GanZhiFlow) calculateElementHarmony() float64 {
	elementCount := make(map[Element]float64)
	totalEnergy := 0.0

	// 统计各五行能量
	for _, state := range f.state.stems {
		elementCount[state.Element] += state.Energy
		totalEnergy += state.Energy
	}
	for _, state := range f.state.branches {
		elementCount[state.Element] += state.Energy
		totalEnergy += state.Energy
	}

	// 计算五行分布的均衡程度
	if totalEnergy == 0 {
		return 1.0
	}

	variance := 0.0
	expectedRatio := 1.0 / 5.0

	for _, energy := range elementCount {
		ratio := energy / totalEnergy
		variance += math.Pow(ratio-expectedRatio, 2)
	}

	// 将方差转换为和谐度（0-1范围）
	harmony := 1.0 - math.Sqrt(variance)/0.9 // 0.9是可能的最大方差
	return math.Max(0, math.Min(1, harmony))
}

// calculatePolarityHarmony 计算阴阳和谐度
func (f *GanZhiFlow) calculatePolarityHarmony() float64 {
	var yinEnergy, yangEnergy float64

	// 统计阴阳能量
	for _, state := range f.state.stems {
		if state.Polarity == NatureYang {
			yangEnergy += state.Energy
		} else {
			yinEnergy += state.Energy
		}
	}
	for _, state := range f.state.branches {
		if state.Polarity == NatureYang {
			yangEnergy += state.Energy
		} else {
			yinEnergy += state.Energy
		}
	}

	totalEnergy := yinEnergy + yangEnergy
	if totalEnergy == 0 {
		return 1.0
	}

	// 计算阴阳平衡度
	balance := 1.0 - math.Abs(yinEnergy-yangEnergy)/totalEnergy
	return math.Max(0, math.Min(1, balance))
}

// calculateEnergyHarmony 计算能量和谐度
func (f *GanZhiFlow) calculateEnergyHarmony() float64 {
	var stemVariance, branchVariance float64
	var totalStemEnergy, totalBranchEnergy float64

	// 计算天干能量方差
	avgStemEnergy := MaxStemEnergy / 2
	for _, state := range f.state.stems {
		totalStemEnergy += state.Energy
		stemVariance += math.Pow(state.Energy-avgStemEnergy, 2)
	}
	stemVariance /= 10 // 天干数量

	// 计算地支能量方差
	avgBranchEnergy := MaxBranchEnergy / 2
	for _, state := range f.state.branches {
		totalBranchEnergy += state.Energy
		branchVariance += math.Pow(state.Energy-avgBranchEnergy, 2)
	}
	branchVariance /= 12 // 地支数量

	// 计算总体和谐度
	totalVariance := (stemVariance + branchVariance) / 2
	maxVariance := math.Pow(MaxStemEnergy/2, 2) // 最大可能方差

	harmony := 1.0 - math.Sqrt(totalVariance)/math.Sqrt(maxVariance)
	return math.Max(0, math.Min(1, harmony))
}

// 辅助函数：判断五行相生
func isGenerating(e1, e2 Element) bool {
	// 木生火、火生土、土生金、金生水、水生木
	generatingPairs := map[Element]Element{
		Wood:  Fire,
		Fire:  Earth,
		Earth: Metal,
		Metal: Water,
		Water: Wood,
	}
	return generatingPairs[e1] == e2
}

// 辅助函数：判断五行相克
func isControlling(e1, e2 Element) bool {
	// 木克土、土克水、水克火、火克金、金克木
	controllingPairs := map[Element]Element{
		Wood:  Earth,
		Earth: Water,
		Water: Fire,
		Fire:  Metal,
		Metal: Wood,
	}
	return controllingPairs[e1] == e2
}

// cyclicTransform 周期性转换
func (f *GanZhiFlow) cyclicTransform() error {
	// 推进周期
	f.state.cycle = (f.state.cycle + 1) % CycleLength

	// 获取当前干支组合
	stem, branch := f.getCurrentGanZhi()

	// 更新能量状态
	if err := f.updateEnergies(stem, branch); err != nil {
		return err
	}

	// 更新量子态
	if err := f.updateQuantumStates(stem, branch); err != nil {
		return err
	}

	return f.updateHarmony()
}

// balanceElements 平衡五行能量
func (f *GanZhiFlow) balanceElements(distribution map[Element]float64) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 计算理想的平均能量
	totalEnergy := 0.0
	for _, energy := range distribution {
		totalEnergy += energy
	}
	idealEnergy := totalEnergy / 5.0 // 五行平均

	// 计算每个五行的能量差异
	adjustments := make(map[Element]float64)
	for element, energy := range distribution {
		adjustments[element] = idealEnergy - energy
	}

	// 应用调整到天干
	for stem, state := range f.state.stems {
		element := state.Element
		adjustment := adjustments[element] / float64(len(f.state.stems))

		// 计算新能量
		newEnergy := state.Energy + adjustment
		// 确保能量在有效范围内
		newEnergy = math.Max(0, math.Min(MaxStemEnergy, newEnergy))

		// 更新能量
		state.Energy = newEnergy
		if err := f.components.stemStates[stem].SetEnergy(newEnergy); err != nil {
			return err
		}
	}

	// 应用调整到地支
	for branch, state := range f.state.branches {
		element := state.Element
		adjustment := adjustments[element] / float64(len(f.state.branches))

		// 计算新能量
		newEnergy := state.Energy + adjustment
		// 确保能量在有效范围内
		newEnergy = math.Max(0, math.Min(MaxBranchEnergy, newEnergy))

		// 更新能量
		state.Energy = newEnergy
		if err := f.components.branchStates[branch].SetEnergy(newEnergy); err != nil {
			return err
		}
	}

	// 更新相关联的量子态
	for stem := range f.state.stems {
		quantumState := f.components.stemStates[stem]
		if err := quantumState.Update(); err != nil {
			return err
		}
	}

	for branch := range f.state.branches {
		quantumState := f.components.branchStates[branch]
		if err := quantumState.Update(); err != nil {
			return err
		}
	}

	return nil
}

// harmonizeElements 调和五行
func (f *GanZhiFlow) harmonizeElements() error {
	// 计算五行分布
	elementDistribution := make(map[Element]float64)

	// 统计天干五行能量
	for _, state := range f.state.stems {
		elementDistribution[state.Element] += state.Energy
	}

	// 统计地支五行能量
	for _, state := range f.state.branches {
		elementDistribution[state.Element] += state.Energy
	}

	// 调和五行能量
	if err := f.balanceElements(elementDistribution); err != nil {
		return err
	}

	return f.updateHarmony()
}

// calculateInteractions 计算干支相互作用
func (f *GanZhiFlow) calculateInteractions(stem HeavenlyStem, branch EarthlyBranch) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	stemState := f.state.stems[stem]
	branchState := f.state.branches[branch]

	// 计算五行相互作用
	elementEffect := f.calculateElementInteraction(stemState.Element, branchState.Element)

	// 计算阴阳相互作用
	polarityEffect := f.calculatePolarityFactor(stemState.Polarity, branchState.Polarity)

	// 复合作用系数
	interactionFactor := elementEffect * polarityEffect

	// 更新天干状态
	newStemEnergy := stemState.Energy * (1 + interactionFactor)
	newStemEnergy = math.Max(0, math.Min(MaxStemEnergy, newStemEnergy))
	stemState.Energy = newStemEnergy

	// 更新地支状态
	newBranchEnergy := branchState.Energy * (1 + interactionFactor)
	newBranchEnergy = math.Max(0, math.Min(MaxBranchEnergy, newBranchEnergy))
	branchState.Energy = newBranchEnergy

	// 更新量子态
	if err := f.components.stemStates[stem].SetEnergy(newStemEnergy); err != nil {
		return err
	}
	if err := f.components.branchStates[branch].SetEnergy(newBranchEnergy); err != nil {
		return err
	}

	// 更新相位关系
	stemPhase := float64(stem) * 2 * math.Pi / 10
	branchPhase := float64(branch) * 2 * math.Pi / 12

	if err := f.components.stemStates[stem].SetPhase(stemPhase); err != nil {
		return err
	}
	if err := f.components.branchStates[branch].SetPhase(branchPhase); err != nil {
		return err
	}

	return nil
}

// updateRelations 更新关系网络
func (f *GanZhiFlow) updateRelations() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 更新天干与地支的关系
	for stem, stemState := range f.state.stems {
		for branch, branchState := range f.state.branches {
			// 计算关系强度因子
			relationFactor := f.calculateRelationFactor(stemState, branchState)

			// 更新干支关系
			stemState.Relations[branch] = relationFactor
			branchState.Relations[stem] = relationFactor

			// 如果关系强度超过阈值，更新量子态的相干性
			if relationFactor > HarmonyThreshold {
				// 调整量子态相位以增加相干性
				stemQuantum := f.components.stemStates[stem]
				branchQuantum := f.components.branchStates[branch]

				stemPhase := stemQuantum.GetPhase()
				branchPhase := branchQuantum.GetPhase()

				// 计算相位调整
				phaseDiff := math.Abs(stemPhase - branchPhase)
				if phaseDiff > math.Pi {
					phaseDiff = 2*math.Pi - phaseDiff
				}

				// 对相位差较大的情况进行调整
				if phaseDiff > math.Pi/2 {
					newPhase := (stemPhase + branchPhase) / 2
					if err := stemQuantum.SetPhase(newPhase); err != nil {
						return err
					}
					if err := branchQuantum.SetPhase(newPhase); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// calculateRelationFactor 计算关系强度因子
func (f *GanZhiFlow) calculateRelationFactor(stemState *StemState, branchState *BranchState) float64 {
	// 基础关系强度
	baseStrength := 0.0

	// 考虑五行相生相克
	elementFactor := f.calculateElementInteraction(stemState.Element, branchState.Element)
	baseStrength += elementFactor * 0.4 // 五行影响权重 40%

	// 考虑阴阳关系
	polarityFactor := f.calculatePolarityFactor(stemState.Polarity, branchState.Polarity)
	baseStrength += polarityFactor * 0.3 // 阴阳影响权重 30%

	// 考虑能量水平
	energyFactor := 1.0 - math.Abs(stemState.Energy/MaxStemEnergy-branchState.Energy/MaxBranchEnergy)
	baseStrength += energyFactor * 0.3 // 能量影响权重 30%

	// 确保关系强度在 [0,1] 范围内
	relationStrength := math.Max(0, math.Min(1, (baseStrength+1)/2))

	return relationStrength
}

// 在干支系统中，自然转换遵循以下原则：
// 1. 计算当前干支组合的相互作用
// 2. 更新干支之间的关系网络
// 3. 维护系统的和谐度
func (f *GanZhiFlow) naturalTransform() error {
	// 获取当前状态
	stem, branch := f.getCurrentGanZhi()

	// 计算相互作用：包括五行相生相克和阴阳平衡
	if err := f.calculateInteractions(stem, branch); err != nil {
		return err
	}

	// 更新关系网络：调整干支之间的关联强度
	if err := f.updateRelations(); err != nil {
		return err
	}

	// 更新系统和谐度
	return f.updateHarmony()
}

// Close 关闭模型
func (f *GanZhiFlow) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 清理组件
	for stem := range f.components.stemFields {
		f.components.stemFields[stem] = nil
		f.components.stemStates[stem] = nil
	}
	for branch := range f.components.branchFields {
		f.components.branchFields[branch] = nil
		f.components.branchStates[branch] = nil
	}

	f.components.cycleManager = nil
	f.components.harmonizer = nil

	return f.BaseFlowModel.Close()
}

// 辅助方法...
func (f *GanZhiFlow) getStemElement(stem HeavenlyStem) Element {
	// 实现天干五行对应关系
	return Wood // 示例返回
}

func (f *GanZhiFlow) getBranchElement(branch EarthlyBranch) Element {
	// 实现地支五行对应关系
	return Wood // 示例返回
}

func (f *GanZhiFlow) getStemPolarity(stem HeavenlyStem) Nature {
	// 实现天干阴阳属性
	return NatureYang // 示例返回
}

func (f *GanZhiFlow) getBranchPolarity(branch EarthlyBranch) Nature {
	// 实现地支阴阳属性
	return NatureYin // 示例返回
}

// AdjustEnergy 调整干支能量
func (f *GanZhiFlow) AdjustEnergy(delta float64) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 首先调用基类方法
	if err := f.BaseFlowModel.AdjustEnergy(delta); err != nil {
		return err
	}

	// 计算天干和地支的总能量
	stemEnergy := 0.0
	branchEnergy := 0.0
	for _, state := range f.state.stems {
		stemEnergy += state.Energy
	}
	for _, state := range f.state.branches {
		branchEnergy += state.Energy
	}
	totalEnergy := stemEnergy + branchEnergy

	// 按比例分配到天干和地支
	deltaStem := delta * (stemEnergy / totalEnergy)
	deltaBranch := delta * (branchEnergy / totalEnergy)

	// 更新天干能量
	for stem, state := range f.state.stems {
		ratio := state.Energy / stemEnergy
		state.Energy += deltaStem * ratio
		if err := f.components.stemStates[stem].SetEnergy(state.Energy); err != nil {
			return err
		}
	}

	// 更新地支能量
	for branch, state := range f.state.branches {
		ratio := state.Energy / branchEnergy
		state.Energy += deltaBranch * ratio
		if err := f.components.branchStates[branch].SetEnergy(state.Energy); err != nil {
			return err
		}
	}

	return nil
}
