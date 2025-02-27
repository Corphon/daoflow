//system/meta/emergence/detector.go

package emergence

import (
	"context"
	"fmt"
	"math"
	"math/cmplx"
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/meta/field"
)

// PatternDetector 模式检测器
type PatternDetector struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		sensitivity       float64       // 检测灵敏度
		timeWindow        time.Duration // 检测时间窗口
		minConfidence     float64       // 最小置信度
		patternThreshold  float64       // 模式阈值
		maxElementEnergy  float64       // 最大元素能量
		maxClusterRadius  float64       // 最大聚集半径
		maxEnergyLevel    float64       // 最大能量级别
		DetectionInterval time.Duration // 检测间隔
	}

	// 检测状态
	state struct {
		activePatterns map[string]*EmergentPattern // 活跃模式
		history        []DetectionEvent            // 检测历史
		lastUpdate     time.Time                   // 最后更新时间
	}

	// 场引用
	field *field.UnifiedField
}

// EmergentPattern 涌现模式
type EmergentPattern struct {
	ID         string             // 模式标识
	Type       string             // 模式类型
	Components []PatternComponent // 组成成分
	Properties map[string]float64 // 模式属性
	Strength   float64            // 模式强度
	Stability  float64            // 模式稳定性
	Energy     float64            // 模式能量
	Formation  time.Time          // 形成时间
	Evolution  []PatternState     // 演化历史
	LastUpdate time.Time          // 最后更新时间
}

// PatternComponent 模式组件
type PatternComponent struct {
	// 场引用
	ID         string             // 组件ID
	Type       string             // 组件类型
	Weight     float64            // 权重
	Role       string             // 角色
	State      map[string]float64 // 状态
	Properties map[string]float64 // 属性
}

// DetectionEvent 检测事件
type DetectionEvent struct {
	Timestamp  time.Time
	PatternID  string
	Type       string
	Confidence float64
	Changes    []StateChange
}

// StateChange 状态变化
type StateChange struct {
	Component string
	Before    map[string]float64
	After     map[string]float64
	Delta     float64
}

// EnergyCluster 能量聚集
type EnergyCluster struct {
	Center   core.Point
	Radius   float64
	Energy   float64
	Gradient float64
	Elements []string
}

// EnergyFlow 能量流动
type EnergyFlow struct {
	Source    core.Point
	Target    core.Point
	Rate      float64
	Direction float64
	Intensity float64
}

// QuantumEntanglement 量子纠缠结构
type QuantumEntanglement struct {
	Strength     float64
	Participants []string
	Duration     time.Duration
	Phase        float64
}

// QuantumCoherence 量子相干结构
type QuantumCoherence struct {
	Amplitude   float64
	Phase       float64
	Stability   float64
	Decoherence float64
}

// ------------------------------------------------------------------
// NewPatternDetector 创建新的模式检测器
func NewPatternDetector(field *field.UnifiedField) *PatternDetector {
	pd := &PatternDetector{
		field: field,
	}

	// 初始化配置
	pd.config.sensitivity = 0.75
	pd.config.timeWindow = 10 * time.Minute
	pd.config.minConfidence = 0.65
	pd.config.patternThreshold = 0.5
	pd.config.maxElementEnergy = 20.0
	pd.config.maxClusterRadius = 5.0
	pd.config.maxEnergyLevel = 100.0
	pd.config.DetectionInterval = 5 * time.Second

	// 初始化状态
	pd.state.activePatterns = make(map[string]*EmergentPattern)
	pd.state.history = make([]DetectionEvent, 0)
	pd.state.lastUpdate = time.Now()

	return pd
}

// Detect 执行模式检测
func (pd *PatternDetector) Detect() ([]EmergentPattern, error) {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	// 获取场状态
	fieldState, err := pd.field.GetState()
	if err != nil {
		return nil, model.WrapError(err, model.ErrCodeOperation, "failed to get field state")
	}

	// 检测新模式
	newPatterns := pd.detectNewPatterns(fieldState)

	// 更新现有模式
	pd.updateExistingPatterns(fieldState)

	// 移除消失的模式
	pd.removeVanishedPatterns()

	// 记录检测事件
	pd.recordDetectionEvent(newPatterns)

	// 返回当前活跃的模式
	return pd.getActivePatterns(), nil
}

// removeVanishedPatterns 移除消失的模式
func (pd *PatternDetector) removeVanishedPatterns() {
	currentTime := time.Now()
	timeout := pd.config.timeWindow

	// 遍历现有模式
	for id, pattern := range pd.state.activePatterns {
		// 检查模式是否超时
		if currentTime.Sub(pattern.LastUpdate) > timeout {
			delete(pd.state.activePatterns, id)
		}
		// 检查模式强度
		if pattern.Strength < pd.config.sensitivity {
			delete(pd.state.activePatterns, id)
		}
	}
}

// getActivePatterns 获取当前活跃的模式
func (pd *PatternDetector) getActivePatterns() []EmergentPattern {
	patterns := make([]EmergentPattern, 0, len(pd.state.activePatterns))
	for _, pattern := range pd.state.activePatterns {
		patterns = append(patterns, *pattern)
	}
	return patterns
}

// detectNewPatterns 检测新模式
func (pd *PatternDetector) detectNewPatterns(state *model.FieldState) []EmergentPattern {
	newPatterns := make([]EmergentPattern, 0)

	// 检测元素组合模式
	elementPatterns := pd.detectElementPatterns(state)
	newPatterns = append(newPatterns, elementPatterns...)

	// 检测能量分布模式
	energyPatterns := pd.detectEnergyPatterns(state)
	newPatterns = append(newPatterns, energyPatterns...)

	// 检测量子态模式
	quantumPatterns := pd.detectQuantumPatterns(state)
	newPatterns = append(newPatterns, quantumPatterns...)

	return newPatterns
}

// detectElementPatterns 检测元素组合模式
func (pd *PatternDetector) detectElementPatterns(state *model.FieldState) []EmergentPattern {
	patterns := make([]EmergentPattern, 0)

	// 获取元素状态
	wuxingElements := state.GetElements()
	if len(wuxingElements) < 2 {
		return patterns
	}

	// 转换WuXingElement为Element
	elements := make([]*model.Element, len(wuxingElements))
	for i, we := range wuxingElements {
		elements[i] = &model.Element{
			Type:       we.String(), // 转换五行枚举为字符串
			Energy:     we.GetEnergy(),
			Properties: we.GetProperties(),
		}
	}

	// 分析元素组合
	combinations := generateElementCombinations(elements)
	for _, combo := range combinations {
		if pattern := pd.analyzeElementCombination(combo); pattern != nil {
			patterns = append(patterns, *pattern)
		}
	}

	return patterns
}

// generateElementCombinations 生成元素组合
func generateElementCombinations(elements []*model.Element) [][]*model.Element {
	combinations := make([][]*model.Element, 0)

	// 生成2个元素的组合
	for i := 0; i < len(elements); i++ {
		for j := i + 1; j < len(elements); j++ {
			combo := []*model.Element{elements[i], elements[j]}
			combinations = append(combinations, combo)
		}
	}

	return combinations
}

// analyzeElementCombination 分析元素组合是否形成模式
func (pd *PatternDetector) analyzeElementCombination(elements []*model.Element) *EmergentPattern {
	// 直接使用model.Element
	interaction := pd.calculateElementInteraction(elements)
	if interaction < pd.config.patternThreshold {
		return nil
	}

	// 创建模式
	pattern := &EmergentPattern{
		ID:         generatePatternID(),
		Type:       "element_combination",
		Strength:   interaction,
		Formation:  time.Now(),
		Components: make([]PatternComponent, len(elements)),
	}

	// 添加组件信息
	for i, elem := range elements {
		pattern.Components[i] = PatternComponent{
			Type:   "element",
			Role:   elem.GetType(),
			Weight: elem.GetEnergy() / pd.config.maxElementEnergy,
		}
	}

	return pattern
}

// generatePatternID 生成唯一的模式ID
func generatePatternID() string {
	return fmt.Sprintf("pat_%d", time.Now().UnixNano())
}

// calculateElementInteraction 计算元素间相互作用强度
func (pd *PatternDetector) calculateElementInteraction(elements []*model.Element) float64 {
	if len(elements) != 2 {
		return 0
	}

	e1, e2 := elements[0], elements[1]

	// 基础相互作用强度
	baseStrength := math.Sqrt(e1.GetEnergy() * e2.GetEnergy())

	// 计算关系强度
	relation := model.GetWuXingRelation(e1.GetType(), e2.GetType())
	relationFactor := relation.Factor

	return baseStrength * relationFactor
}

// calculateDistance 计算元素间距离
func calculateDistance(e1, e2 *model.Element) float64 {
	// 基于能量差的距离
	energyDist := math.Abs(e1.GetEnergy() - e2.GetEnergy())

	// 基于五行关系的调整
	relation := model.GetWuXingRelation(e1.GetType(), e2.GetType()).Factor
	relationDist := 2.0 - relation // 相生=1.2->0.8, 相克=0.8->1.2

	return energyDist * relationDist
}

// calculateSimilarity 计算元素相似度
func calculateSimilarity(e1, e2 *model.Element) float64 {
	// 属性差异
	diffSum := 0.0
	count := 0.0

	// 比较共同属性
	for key, val1 := range e1.GetProperties() {
		if val2, exists := e2.GetProperties()[key]; exists {
			diff := math.Abs(val1 - val2)
			diffSum += diff
			count++
		}
	}

	if count == 0 {
		return 0
	}

	// 归一化相似度
	avgDiff := diffSum / count
	similarity := 1.0 / (1.0 + avgDiff)

	return similarity
}

// detectEnergyPatterns 检测能量分布模式
func (pd *PatternDetector) detectEnergyPatterns(state *model.FieldState) []EmergentPattern {
	patterns := make([]EmergentPattern, 0)

	// 分析能量分布
	energyDist := state.GetEnergyDistribution()

	// 检测能量聚集
	clusters := pd.detectEnergyClusters(energyDist)
	for _, cluster := range clusters {
		if pattern := pd.analyzeEnergyCluster(cluster); pattern != nil {
			patterns = append(patterns, *pattern)
		}
	}

	// 检测能量流动
	flows := pd.detectEnergyFlows(energyDist)
	for _, flow := range flows {
		if pattern := pd.analyzeEnergyFlow(flow); pattern != nil {
			patterns = append(patterns, *pattern)
		}
	}

	return patterns
}

// detectEnergyClusters 检测能量聚集
func (pd *PatternDetector) detectEnergyClusters(dist map[core.Point]float64) []EnergyCluster {
	clusters := make([]EnergyCluster, 0)
	visited := make(map[core.Point]bool)

	for point, energy := range dist {
		if visited[point] || energy < pd.config.sensitivity {
			continue
		}

		// 寻找聚集中心
		cluster := pd.expandCluster(point, dist, visited)
		if cluster.Energy > pd.config.patternThreshold {
			clusters = append(clusters, cluster)
		}
	}

	return clusters
}

// expandCluster 扩展能量聚集区域
func (pd *PatternDetector) expandCluster(
	center core.Point,
	dist map[core.Point]float64,
	visited map[core.Point]bool) EnergyCluster {

	cluster := EnergyCluster{
		Center:   center,
		Energy:   dist[center],
		Elements: make([]string, 0),
	}

	// 标记中心点已访问
	visited[center] = true

	// 查找相邻点
	neighbors := getNeighborPoints(center)
	for _, p := range neighbors {
		if energy, exists := dist[p]; exists {
			if !visited[p] && energy >= pd.config.sensitivity {
				// 计算到中心的距离
				distance := calculatePointDistance(center, p)
				if distance <= pd.config.maxClusterRadius {
					// 递归扩展
					subCluster := pd.expandCluster(p, dist, visited)
					// 更新聚集特征
					cluster.Energy += subCluster.Energy
					cluster.Radius = math.Max(cluster.Radius, distance)
					cluster.Gradient = (cluster.Energy - energy) / distance
				}
			}
		}
	}

	return cluster
}

// getNeighborPoints 获取相邻点
func getNeighborPoints(p core.Point) []core.Point {
	neighbors := make([]core.Point, 0)
	// 上下左右四个方向
	directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for _, d := range directions {
		neighbor := core.Point{
			X: p.X + d[0],
			Y: p.Y + d[1],
		}
		neighbors = append(neighbors, neighbor)
	}
	return neighbors
}

// calculatePointDistance 计算两点间距离
func calculatePointDistance(p1, p2 core.Point) float64 {
	dx := float64(p1.X - p2.X)
	dy := float64(p1.Y - p2.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

// analyzeEnergyCluster 分析能量聚集
func (pd *PatternDetector) analyzeEnergyCluster(cluster EnergyCluster) *EmergentPattern {
	return &EmergentPattern{
		ID:       generatePatternID(),
		Type:     "energy_cluster",
		Strength: cluster.Energy,
		Components: []PatternComponent{{
			Type:   "energy",
			Role:   "center",
			Weight: cluster.Energy,
		}},
		Properties: map[string]float64{
			"radius":   cluster.Radius,
			"gradient": cluster.Gradient,
			"density":  cluster.Energy / (math.Pi * cluster.Radius * cluster.Radius),
		},
	}
}

// detectEnergyFlows 检测能量流动
func (pd *PatternDetector) detectEnergyFlows(dist map[core.Point]float64) []EnergyFlow {
	flows := make([]EnergyFlow, 0)

	// 计算能量梯度
	for p1, e1 := range dist {
		for p2, e2 := range dist {
			if gradient := pd.calculateEnergyGradient(p1, e1, p2, e2); gradient > pd.config.sensitivity {
				flows = append(flows, EnergyFlow{
					Source:    p1,
					Target:    p2,
					Rate:      gradient,
					Direction: calculateDirection(p1, p2),
					Intensity: math.Abs(e1 - e2),
				})
			}
		}
	}

	return flows
}

// calculateEnergyGradient 计算能量梯度
func (pd *PatternDetector) calculateEnergyGradient(p1 core.Point, e1 float64, p2 core.Point, e2 float64) float64 {
	// 计算距离
	distance := calculatePointDistance(p1, p2)
	if distance == 0 {
		return 0
	}

	// 计算能量差除以距离得到梯度
	return math.Abs(e2-e1) / distance
}

// calculateDirection 计算方向角度(弧度)
func calculateDirection(p1, p2 core.Point) float64 {
	dx := float64(p2.X - p1.X)
	dy := float64(p2.Y - p1.Y)

	// 使用反正切函数计算角度
	angle := math.Atan2(dy, dx)

	// 确保角度在[0,2π]范围内
	if angle < 0 {
		angle += 2 * math.Pi
	}

	return angle
}

// analyzeEnergyFlow 分析能量流动
func (pd *PatternDetector) analyzeEnergyFlow(flow EnergyFlow) *EmergentPattern {
	return &EmergentPattern{
		ID:       generatePatternID(),
		Type:     "energy_flow",
		Strength: flow.Intensity,
		Components: []PatternComponent{
			{
				Type:   "energy",
				Role:   "source",
				Weight: flow.Rate,
			},
			{
				Type:   "energy",
				Role:   "target",
				Weight: flow.Rate,
			},
		},
		Properties: map[string]float64{
			"rate":      flow.Rate,
			"direction": flow.Direction,
			"intensity": flow.Intensity,
		},
	}
}

// detectQuantumPatterns 检测量子态模式
func (pd *PatternDetector) detectQuantumPatterns(state *model.FieldState) []EmergentPattern {
	patterns := make([]EmergentPattern, 0)

	// 获取量子态信息
	quantumState := state.GetQuantumState()

	// 检测纠缠模式
	entanglements := pd.detectEntanglements(quantumState)
	for _, ent := range entanglements {
		if pattern := pd.analyzeEntanglement(ent); pattern != nil {
			patterns = append(patterns, *pattern)
		}
	}

	// 检测相干模式
	coherences := pd.detectCoherences(quantumState)
	for _, coh := range coherences {
		if pattern := pd.analyzeCoherence(coh); pattern != nil {
			patterns = append(patterns, *pattern)
		}
	}

	return patterns
}

// detectEntanglements 检测量子纠缠模式
func (pd *PatternDetector) detectEntanglements(state *core.QuantumState) []QuantumEntanglement {
	entanglements := make([]QuantumEntanglement, 0)

	// 获取纠缠度
	entanglement := state.GetEntanglement()
	if entanglement > pd.config.sensitivity {
		// 检测到纠缠
		ent := QuantumEntanglement{
			Strength: entanglement,
			Phase:    state.GetPhase(),
			Duration: pd.config.timeWindow,
		}
		entanglements = append(entanglements, ent)
	}

	return entanglements
}

// analyzeEntanglement 分析量子纠缠模式
func (pd *PatternDetector) analyzeEntanglement(ent QuantumEntanglement) *EmergentPattern {
	return &EmergentPattern{
		ID:       generatePatternID(),
		Type:     "quantum_entanglement",
		Strength: ent.Strength,
		Components: []PatternComponent{{
			Type:   "quantum",
			Role:   "entangled_state",
			Weight: ent.Strength,
		}},
		Properties: map[string]float64{
			"phase":    ent.Phase,
			"duration": ent.Duration.Seconds(),
		},
	}
}

// detectCoherences 检测量子相干模式
func (pd *PatternDetector) detectCoherences(state *core.QuantumState) []QuantumCoherence {
	coherences := make([]QuantumCoherence, 0)

	// 获取相干性
	coherence := state.GetCoherence()
	if coherence > pd.config.sensitivity {
		// 获取振幅数组并计算平均模值
		amplitudes := state.GetAmplitude()
		avgAmplitude := 0.0
		if len(amplitudes) > 0 {
			for _, amp := range amplitudes {
				avgAmplitude += cmplx.Abs(amp)
			}
			avgAmplitude /= float64(len(amplitudes))
		}

		// 检测到相干
		coh := QuantumCoherence{
			Amplitude:   avgAmplitude,
			Phase:       state.GetPhase(),
			Stability:   coherence,
			Decoherence: 1 - coherence,
		}
		coherences = append(coherences, coh)
	}

	return coherences
}

// analyzeCoherence 分析量子相干模式
func (pd *PatternDetector) analyzeCoherence(coh QuantumCoherence) *EmergentPattern {
	return &EmergentPattern{
		ID:       generatePatternID(),
		Type:     "quantum_coherence",
		Strength: coh.Stability,
		Components: []PatternComponent{{
			Type:   "quantum",
			Role:   "coherent_state",
			Weight: coh.Amplitude,
		}},
		Properties: map[string]float64{
			"phase":       coh.Phase,
			"amplitude":   coh.Amplitude,
			"decoherence": coh.Decoherence,
		},
	}
}

// updateExistingPatterns 更新现有模式
func (pd *PatternDetector) updateExistingPatterns(state *model.FieldState) {
	for id, pattern := range pd.state.activePatterns {
		// 检查模式是否仍然存在
		if exists := pd.verifyPattern(pattern, state); !exists {
			continue
		}

		// 更新模式属性
		pd.updatePatternProperties(pattern, state)

		// 检查模式稳定性
		if pattern.Stability < pd.config.minConfidence {
			delete(pd.state.activePatterns, id)
			continue
		}

		pattern.LastUpdate = time.Now()
	}
}

// updatePatternProperties 更新模式属性
func (pd *PatternDetector) updatePatternProperties(pattern *EmergentPattern, state *model.FieldState) {
	// 更新模式强度
	pattern.Strength = pd.calculatePatternStrength(pattern, state)

	// 更新各组件状态
	for i, comp := range pattern.Components {
		if newState := pd.getComponentState(comp, state); newState != nil {
			pattern.Components[i].State = newState
		}
	}

	// 计算稳定性
	pattern.Stability = pd.calculatePatternStability(pattern)

	// 更新基本属性
	pattern.Properties = pd.calculatePatternProperties(pattern, state)
}

// verifyPattern 验证模式是否仍然存在
func (pd *PatternDetector) verifyPattern(pattern *EmergentPattern, state *model.FieldState) bool {
	// 检查组件是否仍然存在
	for _, comp := range pattern.Components {
		if !pd.componentExists(comp, state) {
			return false
		}
	}

	// 检查模式强度
	strength := pd.calculatePatternStrength(pattern, state)
	if strength < pd.config.sensitivity {
		return false
	}

	pattern.Strength = strength
	return true
}

// recordDetectionEvent 记录检测事件
func (pd *PatternDetector) recordDetectionEvent(newPatterns []EmergentPattern) {
	event := DetectionEvent{
		Timestamp: time.Now(),
		Changes:   make([]StateChange, 0),
	}

	// 记录新模式
	for _, pattern := range newPatterns {
		change := StateChange{
			Component: pattern.ID,
			After:     pattern.Properties,
		}
		event.Changes = append(event.Changes, change)
	}

	pd.state.history = append(pd.state.history, event)

	// 限制历史记录长度
	if len(pd.state.history) > maxHistoryLength {
		pd.state.history = pd.state.history[1:]
	}
}

// 辅助函数
// componentExists 检查组件是否存在
func (pd *PatternDetector) componentExists(comp PatternComponent, state *model.FieldState) bool {
	switch comp.Type {
	case "element":
		return state.HasElement(comp.Role)
	case "energy":
		return state.HasEnergyLevel(comp.Weight)
	case "quantum":
		// 检查量子态属性
		if qs := state.GetQuantumState(); qs != nil {
			// 逐个检查量子态属性
			for key, expectedValue := range comp.State {
				switch key {
				case "probability":
					if math.Abs(qs.GetProbability()-expectedValue) > 0.1 {
						return false
					}
				case "phase":
					if math.Abs(qs.GetPhase()-expectedValue) > 0.1 {
						return false
					}
				case "coherence":
					if math.Abs(qs.GetCoherence()-expectedValue) > 0.1 {
						return false
					}
				}
			}
			return true
		}
		return false
	default:
		return false
	}
}

func (pd *PatternDetector) calculatePatternStrength(pattern *EmergentPattern, state *model.FieldState) float64 {
	totalStrength := 0.0
	weightSum := 0.0

	for _, comp := range pattern.Components {
		strength := pd.calculateComponentStrength(comp, state)
		totalStrength += strength * comp.Weight
		weightSum += comp.Weight
	}

	if weightSum == 0 {
		return 0
	}

	return totalStrength / weightSum
}

// calculateComponentStrength 计算组件强度
func (pd *PatternDetector) calculateComponentStrength(comp PatternComponent, state *model.FieldState) float64 {
	switch comp.Type {
	case "element":
		// 元素组件强度
		if element := pd.findElement(comp.Role, state); element != nil {
			return element.Energy / pd.config.maxElementEnergy
		}

	case "energy":
		// 能量组件强度
		return state.GetEnergyLevel() / pd.config.maxEnergyLevel

	case "quantum":
		// 量子组件强度
		if quantum := state.GetQuantumState(); quantum != nil {
			return quantum.GetCoherence()
		}

	case "field":
		// 场组件强度
		return state.GetFieldStrength()
	}

	return 0
}

// findElement 查找指定类型的元素
func (pd *PatternDetector) findElement(elementType string, state *model.FieldState) *model.Element {
	wuxingElements := state.GetElements()
	for _, we := range wuxingElements {
		if we.GetType() == elementType {
			// 转换WuXingElement为Element
			return &model.Element{
				Type:       we.GetType(),
				Energy:     we.GetEnergy(),
				Properties: we.GetProperties(),
			}
		}
	}
	return nil
}

// getComponentState 获取组件的状态
func (pd *PatternDetector) getComponentState(comp PatternComponent, state *model.FieldState) map[string]float64 {
	// 根据组件类型获取状态
	switch comp.Type {
	case "element":
		if element := pd.findElement(comp.Role, state); element != nil {
			return element.Properties
		}
	case "energy":
		return map[string]float64{
			"level": state.GetEnergyLevel(),
			"flow":  state.GetEnergyFlow(),
		}
	}
	return nil
}

func (pd *PatternDetector) calculatePatternStability(pattern *EmergentPattern) float64 {
	// 基于组件状态计算稳定性
	stabilitySum := 0.0
	weights := 0.0

	for _, comp := range pattern.Components {
		weight := comp.Weight
		stability := pd.calculateComponentStability(comp)
		stabilitySum += stability * weight
		weights += weight
	}

	if weights > 0 {
		return stabilitySum / weights
	}
	return 0
}

func (pd *PatternDetector) calculatePatternProperties(pattern *EmergentPattern, state *model.FieldState) map[string]float64 {
	props := make(map[string]float64)

	// 基础属性
	props["coherence"] = pd.calculatePatternCoherence(pattern)
	props["complexity"] = calculatePatternComplexity(pattern)
	props["energy"] = pd.calculatePatternEnergy(pattern, state)

	return props
}

// calculatePatternsComplexity 计算模式组合的复杂度
func calculatePatternsComplexity(patterns []EmergentPattern) float64 {
	if len(patterns) == 0 {
		return 0
	}

	// 基础复杂度 - 模式数量
	baseComplexity := float64(len(patterns)) / 10.0

	// 关系复杂度 - 模式间关联
	relationComplexity := 0.0
	for i := 0; i < len(patterns)-1; i++ {
		for j := i + 1; j < len(patterns); j++ {
			if patterns[i].Type == patterns[j].Type {
				relationComplexity += 0.1
			}
		}
	}

	// 属性复杂度
	propertyComplexity := 0.0
	for _, p := range patterns {
		propertyComplexity += float64(len(p.Properties)) / 10.0
	}

	return (baseComplexity + relationComplexity + propertyComplexity) / 3.0
}

// calculateComponentStability 计算组件稳定性
func (pd *PatternDetector) calculateComponentStability(comp PatternComponent) float64 {
	if comp.State == nil {
		return 0
	}

	// 计算状态变量的稳定性
	stateVariance := 0.0
	for _, value := range comp.State {
		stateVariance += math.Abs(value - 0.5) // 偏离中值的程度
	}

	return 1.0 / (1.0 + stateVariance)
}

// calculatePatternCoherence 计算模式相干性
func (pd *PatternDetector) calculatePatternCoherence(pattern *EmergentPattern) float64 {
	// 时间相干性
	timeCoherence := 1.0
	if pattern.Formation.After(time.Time{}) { // 使用 Formation 替代 Created
		age := time.Since(pattern.Formation).Hours()
		timeCoherence = math.Exp(-age / 24.0) // 24小时衰减
	}

	// 空间相干性
	spaceCoherence := pd.calculateSpatialCoherence(pattern.Components)

	// 量子相干性
	quantumCoherence := pd.calculateQuantumCoherence(pattern)

	// 综合计算相干性，各因素权重相等
	return (timeCoherence + spaceCoherence + quantumCoherence) / 3.0
}

// CalculateStructureComplexity 计算结构复杂度(导出方法)
func (pd *PatternDetector) CalculateStructureComplexity(pattern *EmergentPattern) float64 {
	// 组件复杂度
	componentComplexity := float64(len(pattern.Components)) / 10.0

	// 关系复杂度
	relationComplexity := pd.calculateRelationComplexity(pattern)

	// 结构复杂度
	structureComplexity := pd.calculateTopologyComplexity(pattern)

	return (componentComplexity + relationComplexity + structureComplexity) / 3.0
}

// calculatePatternEnergy 计算模式能量
func (pd *PatternDetector) calculatePatternEnergy(pattern *EmergentPattern, state *model.FieldState) float64 {
	totalEnergy := 0.0

	// 累加组件能量
	for _, comp := range pattern.Components {
		if comp.Type == "element" {
			if element := pd.findElement(comp.Role, state); element != nil {
				totalEnergy += element.Energy * comp.Weight
			}
		}
	}

	return totalEnergy / pd.config.maxElementEnergy
}

// calculateSpatialCoherence 计算空间相干性
func (pd *PatternDetector) calculateSpatialCoherence(components []PatternComponent) float64 {
	if len(components) < 2 {
		return 0
	}

	// 计算组件间的空间关联
	coherence := 0.0
	pairs := 0

	for i := 0; i < len(components)-1; i++ {
		for j := i + 1; j < len(components); j++ {
			// 计算组件对的空间关联度
			correlation := pd.calculateComponentCorrelation(components[i], components[j])
			coherence += correlation
			pairs++
		}
	}

	return coherence / float64(pairs)
}

// calculateComponentCorrelation 计算组件间的空间关联度
func (pd *PatternDetector) calculateComponentCorrelation(c1, c2 PatternComponent) float64 {
	// 判断组件类型
	if c1.Type != c2.Type {
		return 0.5 // 不同类型组件的基础关联度
	}

	// 根据组件类型计算关联度
	switch c1.Type {
	case "element":
		// 元素组件关联度基于五行关系
		relation := model.GetWuXingRelation(c1.Role, c2.Role)
		return (relation.Factor + 1.0) / 2.0

	case "energy":
		// 能量组件关联度基于能级差异
		energyDiff := math.Abs(c1.Weight - c2.Weight)
		return 1.0 / (1.0 + energyDiff)

	case "quantum":
		// 量子组件关联度基于量子纠缠
		if c1.Properties != nil && c2.Properties != nil {
			ent1 := c1.Properties["entanglement"]
			ent2 := c2.Properties["entanglement"]
			return math.Sqrt(ent1 * ent2)
		}
	}

	return 0.5 // 默认关联度
}

// calculateQuantumCoherence 计算量子相干性
func (pd *PatternDetector) calculateQuantumCoherence(pattern *EmergentPattern) float64 {
	// 提取量子组件
	quantumComponents := make([]PatternComponent, 0)
	for _, comp := range pattern.Components {
		if comp.Type == "quantum" {
			quantumComponents = append(quantumComponents, comp)
		}
	}

	if len(quantumComponents) == 0 {
		return 0
	}

	// 计算量子相干性
	coherence := 0.0
	for _, comp := range quantumComponents {
		if value, exists := comp.Properties["coherence"]; exists {
			coherence += value
		}
	}

	return coherence / float64(len(quantumComponents))
}

// calculateRelationComplexity 计算关系复杂度
func (pd *PatternDetector) calculateRelationComplexity(pattern *EmergentPattern) float64 {
	relationCount := 0
	relationStrength := 0.0

	for i := 0; i < len(pattern.Components)-1; i++ {
		for j := i + 1; j < len(pattern.Components); j++ {
			// 计算组件间的关系强度
			strength := pd.calculateComponentRelation(
				pattern.Components[i],
				pattern.Components[j],
			)
			if strength > 0 {
				relationCount++
				relationStrength += strength
			}
		}
	}

	if relationCount == 0 {
		return 0
	}

	// 关系复杂度与关系数量和强度相关
	return (float64(relationCount) / float64(len(pattern.Components))) *
		(relationStrength / float64(relationCount))
}

// calculateComponentRelation 计算组件间的关系强度
func (pd *PatternDetector) calculateComponentRelation(c1, c2 PatternComponent) float64 {
	// 基础关系强度
	baseStrength := math.Min(c1.Weight, c2.Weight)

	// 类型相关性调整
	typeAdjustment := 1.0
	if c1.Type == c2.Type {
		typeAdjustment = 1.2 // 同类型加强
	}

	// 角色关系调整
	roleAdjustment := 1.0
	switch {
	case c1.Type == "element" && c2.Type == "element":
		relation := model.GetWuXingRelation(c1.Role, c2.Role)
		roleAdjustment = relation.Factor

	case c1.Type == "energy" && c2.Type == "energy":
		// 能量梯度关系
		energyDiff := math.Abs(c1.Weight - c2.Weight)
		roleAdjustment = 1.0 / (1.0 + energyDiff)

	case c1.Type == "quantum" && c2.Type == "quantum":
		// 量子关联
		if c1.Properties != nil && c2.Properties != nil {
			ent1 := c1.Properties["entanglement"]
			ent2 := c2.Properties["entanglement"]
			roleAdjustment = math.Sqrt(ent1 * ent2)
		}
	}

	return baseStrength * typeAdjustment * roleAdjustment
}

// calculateStructureComplexity 计算结构复杂度
func (pd *PatternDetector) calculateStructureComplexity(pattern *EmergentPattern) float64 {
	// 组件复杂度
	componentComplexity := float64(len(pattern.Components)) / 10.0

	// 关系复杂度
	relationComplexity := pd.calculateRelationComplexity(pattern)

	// 结构复杂度
	structureComplexity := pd.calculateTopologyComplexity(pattern)

	return (componentComplexity + relationComplexity + structureComplexity) / 3.0
}

// calculateTopologyComplexity 计算拓扑复杂度
func (pd *PatternDetector) calculateTopologyComplexity(pattern *EmergentPattern) float64 {
	// 连通性
	connectivity := 0.0
	for i := 0; i < len(pattern.Components)-1; i++ {
		for j := i + 1; j < len(pattern.Components); j++ {
			if pd.calculateComponentRelation(pattern.Components[i], pattern.Components[j]) > 0 {
				connectivity++
			}
		}
	}
	if len(pattern.Components) > 1 {
		connectivity /= float64((len(pattern.Components) * (len(pattern.Components) - 1)) / 2)
	}

	return connectivity
}

// calculateHierarchyComplexity 计算层次复杂度
func (pd *PatternDetector) calculateHierarchyComplexity(pattern *EmergentPattern) float64 {
	// 构建组件层次图
	levels := make(map[string]int)
	maxLevel := 0

	// 基于组件关系确定层次
	for _, comp := range pattern.Components {
		level := 0
		for _, other := range pattern.Components {
			if pd.calculateComponentRelation(comp, other) > 0.8 { // 强关系
				level++
			}
		}
		levels[comp.Role] = level
		if level > maxLevel {
			maxLevel = level
		}
	}

	if maxLevel == 0 {
		return 0
	}

	return float64(maxLevel) / 10.0 // 归一化
}

// calculateSymmetryDegree 计算对称度
func (pd *PatternDetector) calculateSymmetryDegree(pattern *EmergentPattern) float64 {
	if len(pattern.Components) < 2 {
		return 1.0
	}

	// 计算组件对的对称性
	symmetricPairs := 0.0
	totalPairs := 0.0

	for i := 0; i < len(pattern.Components)-1; i++ {
		for j := i + 1; j < len(pattern.Components); j++ {
			totalPairs++
			// 检查组件对是否对称(关系强度相近)
			rel1 := pd.calculateComponentRelation(pattern.Components[i], pattern.Components[j])
			rel2 := pd.calculateComponentRelation(pattern.Components[j], pattern.Components[i])
			if math.Abs(rel1-rel2) < 0.1 {
				symmetricPairs++
			}
		}
	}

	return symmetricPairs / totalPairs
}

// GetStructureComplexity 获取结构复杂度
func (ep *EmergentPattern) GetStructureComplexity() float64 {
	if value, exists := ep.Properties["complexity"]; exists {
		return value
	}
	return defaultDetector.calculateStructureComplexity(ep)
}

// GetStructureCoherence 获取结构相干性
func (ep *EmergentPattern) GetStructureCoherence() float64 {
	if value, exists := ep.Properties["coherence"]; exists {
		return value
	}
	return defaultDetector.calculateStructureCoherence(ep)
}

// calculateStructureCoherence 计算结构相干性
func (pd *PatternDetector) calculateStructureCoherence(pattern *EmergentPattern) float64 {
	// 时间相干性
	timeCoherence := 0.4 // 默认值

	// 空间相干性
	spaceCoherence := pd.calculateSpatialCoherence(pattern.Components)

	// 量子相干性
	quantumCoherence := pd.calculateQuantumCoherence(pattern)

	return (timeCoherence + spaceCoherence + quantumCoherence) / 3.0
}

// GetStructureSymmetry 获取结构对称性
func (ep *EmergentPattern) GetStructureSymmetry() float64 {
	if value, exists := ep.Properties["symmetry"]; exists {
		return value
	}
	return calculateSymmetryDegree(ep)
}

func calculateSymmetryDegree(pattern *EmergentPattern) float64 {
	if pattern == nil || len(pattern.Components) < 2 {
		return 0
	}

	// 计算结构对称性
	symmetry := 0.0

	// 1. 组件对称性
	componentSymmetry := calculateComponentSymmetry(pattern.Components)

	// 2. 拓扑对称性
	topologySymmetry := calculateTopologySymmetry(pattern.Components)

	// 3. 属性对称性
	propertySymmetry := calculatePropertySymmetry(pattern.Properties)

	// 加权平均
	symmetry = componentSymmetry*0.4 + topologySymmetry*0.3 + propertySymmetry*0.3

	return math.Max(0, math.Min(1, symmetry)) // 确保在0-1范围内
}

// 计算组件对称性
func calculateComponentSymmetry(components []PatternComponent) float64 {
	n := len(components)
	if n < 2 {
		return 0
	}

	symmetricPairs := 0.0
	totalPairs := float64(n * (n - 1) / 2)

	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if components[i].Type == components[j].Type &&
				math.Abs(components[i].Weight-components[j].Weight) < 0.1 {
				symmetricPairs++
			}
		}
	}

	return symmetricPairs / totalPairs
}

// 计算拓扑对称性
func calculateTopologySymmetry(components []PatternComponent) float64 {
	if len(components) < 2 {
		return 0
	}

	// 计算组件对之间的距离矩阵
	n := len(components)
	distances := make([][]float64, n)
	for i := range distances {
		distances[i] = make([]float64, n)
	}

	// 填充距离矩阵
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			// 基于组件权重和类型计算距离
			typeDist := 0.0
			if components[i].Type == components[j].Type {
				typeDist = 1.0
			}
			weightDist := 1.0 - math.Abs(components[i].Weight-components[j].Weight)

			// 综合距离
			dist := (typeDist + weightDist) / 2.0
			distances[i][j] = dist
			distances[j][i] = dist
		}
	}

	// 计算对称度
	symmetry := 0.0
	pairs := 0
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			// 检查(i,j)与其他对称点的距离是否相等
			for k := 0; k < n-1; k++ {
				for l := k + 1; l < n; l++ {
					if (i != k || j != l) && math.Abs(distances[i][j]-distances[k][l]) < 0.1 {
						symmetry += 1.0
					}
					pairs++
				}
			}
		}
	}

	if pairs > 0 {
		return symmetry / float64(pairs)
	}
	return 0
}

// 计算属性对称性
func calculatePropertySymmetry(properties map[string]float64) float64 {
	if len(properties) == 0 {
		return 0
	}

	// 计算属性值的分布对称性
	values := make([]float64, 0, len(properties))
	for _, v := range properties {
		values = append(values, v)
	}

	// 计算属性值的偏度作为对称性指标
	mean := calculateMean(values)
	variance := calculateVariance(values, mean)
	skewness := calculateSkewness(values, mean, variance)

	// 转换为0-1范围
	return 1.0 / (1.0 + math.Abs(skewness))
}

// calculateMean 计算平均值
func calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// calculateVariance 计算方差
func calculateVariance(values []float64, mean float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sumSquares := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquares += diff * diff
	}
	return sumSquares / float64(len(values))
}

// calculateSkewness 计算偏度
func calculateSkewness(values []float64, mean float64, variance float64) float64 {
	if len(values) == 0 || variance == 0 {
		return 0
	}

	stdDev := math.Sqrt(variance)
	sum := 0.0
	for _, v := range values {
		diff := (v - mean) / stdDev
		sum += diff * diff * diff
	}
	return sum / float64(len(values))
}

var (
	defaultDetector *PatternDetector
	detectorOnce    sync.Once
)

// GetDefaultDetector 获取默认检测器
func GetDefaultDetector() *PatternDetector {
	detectorOnce.Do(func() {
		// 从全局Field实例创建默认检测器
		field := field.GetDefaultField()
		defaultDetector = NewPatternDetector(field)
	})
	return defaultDetector
}

func init() {
	// 确保默认检测器初始化
	GetDefaultDetector()
}

// EmergentPattern Clone 方法
func (ep *EmergentPattern) Clone() *EmergentPattern {
	clone := &EmergentPattern{
		ID:         ep.ID + "_clone",
		Type:       ep.Type,
		Strength:   ep.Strength,
		Energy:     ep.Energy,
		Formation:  ep.Formation,
		LastUpdate: ep.LastUpdate,
		Components: make([]PatternComponent, len(ep.Components)),
		Properties: make(map[string]float64),
	}

	// 复制组件
	for i, comp := range ep.Components {
		clone.Components[i] = comp.Clone()
	}

	// 复制属性
	for k, v := range ep.Properties {
		clone.Properties[k] = v
	}

	return clone
}

// PatternComponent Clone 方法
func (pc *PatternComponent) Clone() PatternComponent {
	clone := PatternComponent{
		ID:         pc.ID,
		Type:       pc.Type,
		Weight:     pc.Weight,
		Role:       pc.Role,
		Properties: make(map[string]float64),
	}

	// 复制属性
	for k, v := range pc.Properties {
		clone.Properties[k] = v
	}

	return clone
}

// Start 启动模式检测器
func (pd *PatternDetector) Start(ctx context.Context) error {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	// 启动模式检测循环
	go pd.detectionLoop(ctx)

	return nil
}

// Stop 停止模式检测器
func (pd *PatternDetector) Stop() error {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	// 清理资源
	return nil
}

// detectionLoop 检测循环
func (pd *PatternDetector) detectionLoop(ctx context.Context) {
	ticker := time.NewTicker(pd.config.DetectionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pd.Detect()
		}
	}
}
