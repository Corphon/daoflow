//system/meta/resonance/cross.go

package resonance

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/meta/emergence"
	"github.com/Corphon/daoflow/system/meta/field"
	"github.com/Corphon/daoflow/system/types"
)

const (
	maxTransitionHistory = 1000
)

type LayerResonance struct {
	ID          string
	SourceLayer string
	TargetLayer string
	Strength    float64
	Created     time.Time
}

// CrossResonance 跨层共振处理器
type CrossResonance struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		resonanceThreshold float64 // 共振阈值
		couplingStrength   float64 // 耦合强度
		interactionRange   float64 // 交互范围
		decayFactor        float64 // 衰减因子
		coherenceThreshold float64 // 相干度阈值
		minCoherence       float64 // 最小相干度
	}

	// 共振状态
	state struct {
		layers      map[string]*ResonanceLayer  // 共振层
		bridges     map[string]*ResonanceBridge // 层间桥接
		transitions []TransitionEvent           // 转换事件
	}

	// 依赖项
	field     *field.UnifiedField
	amplifier *ResonanceAmplifier
	matcher   *PatternMatcher
}

// ResonanceLayer 共振层
type ResonanceLayer struct {
	ID         string                                // 层ID
	Level      int                                   // 层级
	Patterns   map[string]*emergence.EmergentPattern // 层内模式
	Properties map[string]float64                    // 层属性
	Energy     float64                               // 层能量
	Coherence  float64                               // 相干度
	Created    time.Time                             // 创建时间
}

// ResonanceBridge 共振桥接
type ResonanceBridge struct {
	ID          string          // 桥接ID
	SourceLayer string          // 源层ID
	TargetLayer string          // 目标层ID
	Strength    float64         // 桥接强度
	Phase       float64         // 相位差
	Channels    []BridgeChannel // 传输通道
	Active      bool            // 是否活跃
	Created     time.Time       // 创建时间
}

// BridgeChannel 桥接通道
type BridgeChannel struct {
	Type      string  // 通道类型
	Capacity  float64 // 传输容量
	Load      float64 // 当前负载
	Direction int     // 传输方向
	State     string  // 通道状态
}

// TransitionEvent 转换事件
type TransitionEvent struct {
	Timestamp   time.Time
	SourceLayer string
	TargetLayer string
	Type        string
	Pattern     *emergence.EmergentPattern
	Energy      float64
	Success     bool
}

// -----------------------------------------
// NewCrossResonance 创建新的跨层共振处理器
func NewCrossResonance(
	field *field.UnifiedField,
	amplifier *ResonanceAmplifier,
	matcher *PatternMatcher) *CrossResonance {

	cr := &CrossResonance{
		field:     field,
		amplifier: amplifier,
		matcher:   matcher,
	}

	// 初始化配置
	cr.config.resonanceThreshold = 0.7
	cr.config.couplingStrength = 0.5
	cr.config.interactionRange = 2.0
	cr.config.decayFactor = 0.1
	cr.config.coherenceThreshold = 0.3
	cr.config.minCoherence = 0.2

	// 初始化状态
	cr.state.layers = make(map[string]*ResonanceLayer)
	cr.state.bridges = make(map[string]*ResonanceBridge)
	cr.state.transitions = make([]TransitionEvent, 0)

	return cr
}

// Process 处理跨层共振
func (cr *CrossResonance) Process() error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	// 更新层状态
	if err := cr.updateLayers(); err != nil {
		return err
	}

	// 处理层间共振
	if err := cr.processInterLayerResonance(); err != nil {
		return err
	}

	// 管理桥接
	if err := cr.manageBridges(); err != nil {
		return err
	}

	// 处理转换
	if err := cr.handleTransitions(); err != nil {
		return err
	}

	return nil
}

// updateLayers 更新层状态
func (cr *CrossResonance) updateLayers() error {
	// 获取当前模式
	patterns, err := cr.matcher.GetActivePatterns()
	if err != nil {
		return err
	}

	// 对模式进行分层
	layers := cr.stratifyPatterns(patterns)

	// 更新层状态
	for level, layerPatterns := range layers {
		layerID := fmt.Sprintf("layer_%d", level)

		if layer, exists := cr.state.layers[layerID]; exists {
			// 更新现有层
			cr.updateExistingLayer(layer, layerPatterns)
		} else {
			// 创建新层
			cr.createNewLayer(layerID, level, layerPatterns)
		}
	}

	// 移除空层
	cr.removeEmptyLayers()

	return nil
}

// updateExistingLayer 更新现有层
func (cr *CrossResonance) updateExistingLayer(layer *ResonanceLayer, patterns []*emergence.EmergentPattern) {
	// 更新模式集合
	layer.Patterns = make(map[string]*emergence.EmergentPattern)
	for _, pattern := range patterns {
		layer.Patterns[pattern.ID] = pattern
	}

	// 更新层属性
	layer.Energy = calculateLayerEnergy(patterns)
	layer.Coherence = calculateLayerCoherence(patterns)
}

// createNewLayer 创建新层
func (cr *CrossResonance) createNewLayer(id string, level int, patterns []*emergence.EmergentPattern) {
	layer := &ResonanceLayer{
		ID:         id,
		Level:      level,
		Patterns:   make(map[string]*emergence.EmergentPattern),
		Properties: make(map[string]float64),
		Created:    time.Now(),
	}

	// 添加模式
	for _, pattern := range patterns {
		layer.Patterns[pattern.ID] = pattern
	}

	// 初始化层属性
	layer.Energy = calculateLayerEnergy(patterns)
	layer.Coherence = calculateLayerCoherence(patterns)

	cr.state.layers[id] = layer
}

// removeEmptyLayers 移除空层
func (cr *CrossResonance) removeEmptyLayers() {
	for id, layer := range cr.state.layers {
		if len(layer.Patterns) == 0 {
			delete(cr.state.layers, id)
		}
	}
}

// 辅助函数
// calculateLayerEnergy 计算层能量
func calculateLayerEnergy(patterns []*emergence.EmergentPattern) float64 {
	energy := 0.0
	for _, pattern := range patterns {
		if value, exists := pattern.Properties["energy"]; exists {
			energy += value
		}
	}
	return energy
}

func calculateLayerCoherence(patterns []*emergence.EmergentPattern) float64 {
	if len(patterns) == 0 {
		return 0
	}

	coherence := 0.0
	for _, pattern := range patterns {
		coherence += pattern.GetStructureCoherence()
	}
	return coherence / float64(len(patterns))
}

// processInterLayerResonance 处理层间共振
func (cr *CrossResonance) processInterLayerResonance() error {
	// 遍历相邻层对
	for id1, layer1 := range cr.state.layers {
		for id2, layer2 := range cr.state.layers {
			if id1 == id2 || !cr.areLayersAdjacent(layer1, layer2) {
				continue
			}

			// 检查层间共振
			if resonance := cr.detectLayerResonance(layer1, layer2); resonance != nil {
				// 处理共振效应
				if err := cr.handleLayerResonance(resonance); err != nil {
					continue
				}
			}
		}
	}

	return nil
}

// areLayersAdjacent 判断两层是否相邻
func (cr *CrossResonance) areLayersAdjacent(layer1, layer2 *ResonanceLayer) bool {
	// 检查层级差是否为1
	levelDiff := math.Abs(float64(layer1.Level - layer2.Level))
	return levelDiff == 1
}

// handleLayerResonance 处理层间共振
func (cr *CrossResonance) handleLayerResonance(resonance *LayerResonance) error {
	// 获取源层和目标层
	sourceLayer := cr.state.layers[resonance.SourceLayer]
	targetLayer := cr.state.layers[resonance.TargetLayer]

	if sourceLayer == nil || targetLayer == nil {
		return fmt.Errorf("invalid layer reference")
	}

	// 计算能量转移
	transferEnergy := resonance.Strength * math.Min(sourceLayer.Energy, targetLayer.Energy)

	// 创建桥接
	bridge := &ResonanceBridge{
		ID:          generateBridgeID(),
		SourceLayer: resonance.SourceLayer,
		TargetLayer: resonance.TargetLayer,
		Strength:    resonance.Strength,
		Phase:       calculatePhaseDifference(sourceLayer, targetLayer),
		Channels:    make([]BridgeChannel, 0),
		Active:      true,
		Created:     time.Now(),
	}

	// 添加传输通道
	bridge.Channels = append(bridge.Channels, BridgeChannel{
		Type:      "energy",
		Capacity:  transferEnergy,
		Load:      0,
		Direction: 1,
		State:     "open",
	})

	// 保存桥接
	cr.state.bridges[bridge.ID] = bridge

	return nil
}

// generateBridgeID 复用已有的ID生成模式
func generateBridgeID() string {
	return fmt.Sprintf("bridge_%d", time.Now().UnixNano())
}

// calculatePhaseDifference 计算两层之间的相位差
func calculatePhaseDifference(layer1, layer2 *ResonanceLayer) float64 {
	// 计算各层的平均相位
	phase1 := calculateLayerPhase(layer1)
	phase2 := calculateLayerPhase(layer2)

	// 计算相位差并归一化到[-π,π]区间
	diff := phase1 - phase2
	return normalizePhase(diff)
}

// normalizePhase 将相位标准化到[-π,π]区间
func normalizePhase(phase float64) float64 {
	// 将相位标准化到 [-π, π] 区间
	for phase > math.Pi {
		phase -= 2 * math.Pi
	}
	for phase < -math.Pi {
		phase += 2 * math.Pi
	}
	return phase
}

// calculateLayerPhase 计算层的平均相位
func calculateLayerPhase(layer *ResonanceLayer) float64 {
	if len(layer.Patterns) == 0 {
		return 0
	}

	totalPhase := 0.0
	count := 0

	for _, pattern := range layer.Patterns {
		if phase, exists := pattern.Properties["phase"]; exists {
			totalPhase += phase
			count++
		}
	}

	if count > 0 {
		return totalPhase / float64(count)
	}
	return 0
}

// manageBridges 管理桥接
func (cr *CrossResonance) manageBridges() error {
	// 更新现有桥接
	for id, bridge := range cr.state.bridges {
		if valid := cr.validateBridge(bridge); !valid {
			delete(cr.state.bridges, id)
			continue
		}

		// 更新桥接状态
		if err := cr.updateBridge(bridge); err != nil {
			continue
		}
	}

	// 检测新的桥接机会
	newBridges := cr.detectNewBridges()
	for _, bridge := range newBridges {
		cr.state.bridges[bridge.ID] = bridge
	}

	return nil
}

// validateBridge 验证桥接有效性
func (cr *CrossResonance) validateBridge(bridge *ResonanceBridge) bool {
	// 检查源层和目标层是否存在
	sourceLayer := cr.state.layers[bridge.SourceLayer]
	targetLayer := cr.state.layers[bridge.TargetLayer]
	if sourceLayer == nil || targetLayer == nil {
		return false
	}

	// 检查桥接是否超时
	if time.Since(bridge.Created) > types.MaxBridgeAge {
		return false
	}

	// 检查桥接强度是否足够
	if bridge.Strength < cr.config.resonanceThreshold {
		return false
	}

	return true
}

// updateBridge 更新桥接状态
func (cr *CrossResonance) updateBridge(bridge *ResonanceBridge) error {
	// 更新桥接强度
	sourceLayer := cr.state.layers[bridge.SourceLayer]
	targetLayer := cr.state.layers[bridge.TargetLayer]

	coupling := cr.calculateLayerCoupling(sourceLayer, targetLayer)
	bridge.Strength = coupling

	// 更新传输通道
	for i := range bridge.Channels {
		channel := &bridge.Channels[i]
		if err := cr.updateChannel(channel, sourceLayer, targetLayer); err != nil {
			channel.State = "blocked"
		}
	}

	return nil
}

// updateChannel 更新通道状态
func (cr *CrossResonance) updateChannel(
	channel *BridgeChannel,
	sourceLayer *ResonanceLayer,
	targetLayer *ResonanceLayer) error {

	switch channel.Type {
	case "energy":
		// 检查能量传输条件
		if sourceLayer.Energy < channel.Capacity {
			return fmt.Errorf("insufficient energy")
		}

		// 更新负载
		maxLoad := math.Min(sourceLayer.Energy, channel.Capacity)
		channel.Load = maxLoad * cr.config.couplingStrength

		// 检查通道状态
		if channel.Load > 0 {
			channel.State = "active"
		} else {
			channel.State = "idle"
		}

	case "coherence":
		// 检查相干性传输
		coherenceMatch := math.Abs(sourceLayer.Coherence - targetLayer.Coherence)
		if coherenceMatch > cr.config.resonanceThreshold {
			return fmt.Errorf("coherence mismatch")
		}

		channel.Load = (sourceLayer.Coherence + targetLayer.Coherence) / 2
		channel.State = "active"

	default:
		return fmt.Errorf("unknown channel type")
	}

	return nil
}

// detectNewBridges 检测新桥接机会
func (cr *CrossResonance) detectNewBridges() []*ResonanceBridge {
	bridges := make([]*ResonanceBridge, 0)

	// 遍历相邻层对
	for id1, layer1 := range cr.state.layers {
		for id2, layer2 := range cr.state.layers {
			if id1 == id2 || !cr.areLayersAdjacent(layer1, layer2) {
				continue
			}

			// 检查是否已存在桥接
			if cr.hasBridge(layer1.ID, layer2.ID) {
				continue
			}

			// 检测层间共振
			if resonance := cr.detectLayerResonance(layer1, layer2); resonance != nil {
				bridge := &ResonanceBridge{
					ID:          generateBridgeID(),
					SourceLayer: resonance.SourceLayer,
					TargetLayer: resonance.TargetLayer,
					Strength:    resonance.Strength,
					Phase:       calculatePhaseDifference(layer1, layer2),
					Channels:    make([]BridgeChannel, 0),
					Active:      true,
					Created:     time.Now(),
				}
				bridges = append(bridges, bridge)
			}
		}
	}

	return bridges
}

// hasBridge 检查两层之间是否已存在桥接
func (cr *CrossResonance) hasBridge(sourceID, targetID string) bool {
	// 检查所有桥接
	for _, bridge := range cr.state.bridges {
		// 检查双向桥接
		if (bridge.SourceLayer == sourceID && bridge.TargetLayer == targetID) ||
			(bridge.SourceLayer == targetID && bridge.TargetLayer == sourceID) {
			return true
		}
	}
	return false
}

// handleTransitions 处理转换
func (cr *CrossResonance) handleTransitions() error {
	for _, bridge := range cr.state.bridges {
		if !bridge.Active {
			continue
		}

		// 检查转换条件
		if transitions := cr.detectTransitions(bridge); len(transitions) > 0 {
			// 执行转换
			for _, transition := range transitions {
				if err := cr.executeTransition(transition); err != nil {
					continue
				}
			}
		}
	}

	return nil
}

// detectTransitions 检测桥接上的转换
func (cr *CrossResonance) detectTransitions(bridge *ResonanceBridge) []*TransitionEvent {
	transitions := make([]*TransitionEvent, 0)

	sourceLayer := cr.state.layers[bridge.SourceLayer]
	targetLayer := cr.state.layers[bridge.TargetLayer]

	if sourceLayer == nil || targetLayer == nil {
		return transitions
	}

	// 检查每个活跃通道
	for _, channel := range bridge.Channels {
		if channel.State != "active" {
			continue
		}

		// 基于通道类型检测转换
		switch channel.Type {
		case "energy":
			// 能量转换条件
			if channel.Load >= channel.Capacity {
				// 找出可能转换的模式
				for _, pattern := range sourceLayer.Patterns {
					if cr.isPatternTransferable(pattern, sourceLayer, targetLayer) {
						transitions = append(transitions, &TransitionEvent{
							Timestamp:   time.Now(),
							SourceLayer: bridge.SourceLayer,
							TargetLayer: bridge.TargetLayer,
							Type:        "energy_transfer",
							Pattern:     pattern,
							Energy:      channel.Load,
							Success:     false,
						})
					}
				}
			}

		case "coherence":
			// 相干性转换
			if channel.Load > cr.config.resonanceThreshold {
				// 检查相干性匹配的模式
				for _, pattern := range sourceLayer.Patterns {
					if cr.checkCoherenceMatch(pattern, sourceLayer, targetLayer) {
						transitions = append(transitions, &TransitionEvent{
							Timestamp:   time.Now(),
							SourceLayer: bridge.SourceLayer,
							TargetLayer: bridge.TargetLayer,
							Type:        "coherence_transfer",
							Pattern:     pattern,
							Energy:      channel.Load,
							Success:     false,
						})
					}
				}
			}
		}
	}

	return transitions
}

// isPatternTransferable 判断模式是否可在层间转移
func (cr *CrossResonance) isPatternTransferable(
	pattern *emergence.EmergentPattern,
	sourceLayer *ResonanceLayer,
	targetLayer *ResonanceLayer) bool {

	// 检查能量条件
	if pattern.Properties["energy"] > targetLayer.Energy {
		return false
	}

	// 检查层级差异
	levelDiff := math.Abs(float64(sourceLayer.Level - targetLayer.Level))
	if levelDiff > 1 {
		return false
	}

	// 检查模式类型兼容性
	for _, targetPattern := range targetLayer.Patterns {
		if cr.checkPatternCompatibility(pattern, targetPattern) {
			return true
		}
	}

	return false
}

// checkCoherenceMatch 检查相干性匹配
func (cr *CrossResonance) checkCoherenceMatch(
	pattern *emergence.EmergentPattern,
	sourceLayer *ResonanceLayer,
	targetLayer *ResonanceLayer) bool {

	// 获取相干度
	sourceCoherence := sourceLayer.Coherence
	targetCoherence := targetLayer.Coherence
	patternCoherence := pattern.GetStructureCoherence()

	// 检查相干度匹配
	coherenceDiff := math.Abs(sourceCoherence - targetCoherence)
	if coherenceDiff > cr.config.coherenceThreshold {
		return false
	}

	// 检查相位匹配
	phase1 := sourceLayer.Properties["phase"]
	phase2 := targetLayer.Properties["phase"]
	if math.Abs(phase1-phase2) > math.Pi/2 {
		return false
	}

	return patternCoherence > cr.config.minCoherence
}

// 辅助函数
func (cr *CrossResonance) checkPatternCompatibility(
	p1 *emergence.EmergentPattern,
	p2 *emergence.EmergentPattern) bool {

	// 检查类型兼容
	if p1.Type == p2.Type {
		return true
	}

	// 检查增强关系
	return emergence.CheckEnhancingRelation(p1.Type, p2.Type)
}

// stratifyPatterns 对模式进行分层
func (cr *CrossResonance) stratifyPatterns(
	patterns []*emergence.EmergentPattern) map[int][]*emergence.EmergentPattern {

	layers := make(map[int][]*emergence.EmergentPattern)

	for _, pattern := range patterns {
		// 计算模式层级
		level := cr.calculatePatternLevel(pattern)

		// 添加到对应层
		if _, exists := layers[level]; !exists {
			layers[level] = make([]*emergence.EmergentPattern, 0)
		}
		layers[level] = append(layers[level], pattern)
	}

	return layers
}

// detectLayerResonance 检测层间共振
func (cr *CrossResonance) detectLayerResonance(
	layer1, layer2 *ResonanceLayer) *LayerResonance {

	// 计算层间耦合
	coupling := cr.calculateLayerCoupling(layer1, layer2)
	if coupling < cr.config.resonanceThreshold {
		return nil
	}

	// 检查相位匹配
	if !cr.checkPhaseMatch(layer1, layer2) {
		return nil
	}

	// 创建层间共振
	resonance := &LayerResonance{
		ID:          generateResonanceID(),
		SourceLayer: layer1.ID,
		TargetLayer: layer2.ID,
		Strength:    coupling,
		Created:     time.Now(),
	}

	return resonance
}

// checkPhaseMatch 检查层间相位匹配
func (cr *CrossResonance) checkPhaseMatch(layer1, layer2 *ResonanceLayer) bool {
	// 计算相位差
	phaseDiff := calculatePhaseDifference(layer1, layer2)

	// 相位差在阈值范围内认为匹配
	// 使用π/4作为相位匹配阈值
	if math.Abs(phaseDiff) > math.Pi/4 {
		return false
	}

	// 检查相位相干性
	coherence1 := layer1.Coherence
	coherence2 := layer2.Coherence

	if coherence1 < cr.config.coherenceThreshold ||
		coherence2 < cr.config.coherenceThreshold {
		return false
	}

	return true
}

// executeTransition 执行转换
func (cr *CrossResonance) executeTransition(transition *TransitionEvent) error {
	// 获取源层和目标层
	sourceLayer := cr.state.layers[transition.SourceLayer]
	targetLayer := cr.state.layers[transition.TargetLayer]

	if sourceLayer == nil || targetLayer == nil {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid layers")
	}

	// 执行模式转换
	if err := cr.transferPattern(transition.Pattern, sourceLayer, targetLayer); err != nil {
		return err
	}

	// 更新能量分布
	cr.redistributeEnergy(sourceLayer, targetLayer, transition.Energy)

	// 记录转换事件
	cr.recordTransition(transition)

	return nil
}

// transferPattern 转移模式
func (cr *CrossResonance) transferPattern(
	pattern *emergence.EmergentPattern,
	sourceLayer *ResonanceLayer,
	targetLayer *ResonanceLayer) error {

	// 先从源层移除
	delete(sourceLayer.Patterns, pattern.ID)

	// 计算转移后的属性变化
	pattern.Properties["level"] = float64(targetLayer.Level)
	pattern.Properties["transition_count"] =
		pattern.Properties["transition_count"] + 1

	// 添加到目标层
	targetLayer.Patterns[pattern.ID] = pattern

	return nil
}

// redistributeEnergy 重新分配能量
func (cr *CrossResonance) redistributeEnergy(
	sourceLayer *ResonanceLayer,
	targetLayer *ResonanceLayer,
	transferEnergy float64) {

	// 能量守恒
	sourceLayer.Energy -= transferEnergy
	targetLayer.Energy += transferEnergy * cr.config.couplingStrength

	// 更新相干性
	coherence := (sourceLayer.Coherence + targetLayer.Coherence) / 2.0
	sourceLayer.Coherence = coherence
	targetLayer.Coherence = coherence

	// 更新场强度
	sourceLayer.Properties["field_strength"] -= transferEnergy * 0.1
	targetLayer.Properties["field_strength"] += transferEnergy * 0.1
}

// 辅助函数

func (cr *CrossResonance) calculatePatternLevel(pattern *emergence.EmergentPattern) int {
	// 基于模式复杂度和能量计算层级
	complexity := emergence.GetDefaultDetector().CalculateStructureComplexity(pattern)

	// 从Properties获取能量
	energy := 0.0
	if value, exists := pattern.Properties["energy"]; exists {
		energy = value
	}

	// 使用对数尺度,并确保返回整数
	level := int(math.Log2(complexity * energy))
	if level < 0 {
		return 0
	}
	return level
}

func (cr *CrossResonance) calculateLayerCoupling(layer1, layer2 *ResonanceLayer) float64 {
	// 计算层间耦合强度
	energyCoupling := math.Sqrt(layer1.Energy * layer2.Energy)
	coherenceCoupling := (layer1.Coherence + layer2.Coherence) / 2

	return energyCoupling * coherenceCoupling * cr.config.couplingStrength
}

func (cr *CrossResonance) recordTransition(transition *TransitionEvent) {
	cr.state.transitions = append(cr.state.transitions, *transition)

	// 限制历史记录长度
	if len(cr.state.transitions) > maxTransitionHistory {
		cr.state.transitions = cr.state.transitions[1:]
	}
}
