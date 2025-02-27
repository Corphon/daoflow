//system/meta/resonance/amplifier.go

package resonance

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/Corphon/daoflow/system/meta/emergence"
	"github.com/Corphon/daoflow/system/meta/field"
	"github.com/Corphon/daoflow/system/types"
)

type CouplingTrends struct {
	StrengthTrend float64
	PhaseTrend    float64
	EnergyTrend   float64
	Stability     float64
	Prediction    field.PredictedState
	History       []types.MetricPoint
}

// ResonanceAmplifier 共振放大器
type ResonanceAmplifier struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		minAmplitude          float64       // 最小放大幅度
		maxAmplitude          float64       // 最大放大幅度
		decayRate             float64       // 衰减率
		feedbackRatio         float64       // 反馈比率
		AmplificationInterval time.Duration // 放大间隔
	}

	// 放大状态
	state struct {
		activeResonances map[string]*ResonanceState // 活跃共振
		history          []AmplificationEvent       // 放大历史
		lastUpdate       time.Time                  // 最后更新时间
	}

	// 依赖项
	field     *field.UnifiedField
	detector  *emergence.PatternDetector
	generator *emergence.PropertyGenerator
}

// ResonanceState 共振状态
type ResonanceState struct {
	ID         string           // 共振ID
	Type       string           // 共振类型
	Source     *ResonanceSource // 共振源
	Target     *ResonanceTarget // 共振目标
	Amplitude  float64          // 当前幅度
	Phase      float64          // 当前相位
	Frequency  float64          // 频率
	Energy     float64          // 能量
	Coherence  float64          // 相干度
	StartTime  time.Time        // 开始时间
	LastUpdate time.Time        // 最后更新时间
	Duration   time.Duration    // 持续时间
}

// ResonanceSource 共振源
type ResonanceSource struct {
	Type       string                     // 源类型
	ID         string                     // 源ID
	Pattern    *emergence.EmergentPattern // 关联模式
	Properties map[string]float64         // 源属性
}

// ResonanceTarget 共振目标
type ResonanceTarget struct {
	Type       string                     // 目标类型
	ID         string                     // 目标ID
	Pattern    *emergence.EmergentPattern // 关联模式
	Properties map[string]float64         // 目标属性
}

// AmplificationEvent 放大事件
type AmplificationEvent struct {
	Timestamp   time.Time
	ResonanceID string
	Type        string
	OldState    *ResonanceState
	NewState    *ResonanceState
	Changes     map[string]float64
}

// --------------------------------------------------------------
// NewResonanceAmplifier 创建新的共振放大器
func NewResonanceAmplifier(
	field *field.UnifiedField,
	detector *emergence.PatternDetector,
	generator *emergence.PropertyGenerator) *ResonanceAmplifier {

	ra := &ResonanceAmplifier{
		field:     field,
		detector:  detector,
		generator: generator,
	}

	// 初始化配置
	ra.config.minAmplitude = 0.1
	ra.config.maxAmplitude = 10.0
	ra.config.decayRate = 0.05
	ra.config.feedbackRatio = 0.2
	ra.config.AmplificationInterval = 100 * time.Millisecond

	// 初始化状态
	ra.state.activeResonances = make(map[string]*ResonanceState)
	ra.state.history = make([]AmplificationEvent, 0)
	ra.state.lastUpdate = time.Now()

	return ra
}

// Amplify 执行共振放大
func (ra *ResonanceAmplifier) Amplify() error {
	ra.mu.Lock()
	defer ra.mu.Unlock()

	// 检测新的共振
	newResonances, err := ra.detectResonances()
	if err != nil {
		return err
	}

	// 更新现有共振
	if err := ra.updateResonances(); err != nil {
		return err
	}

	// 应用放大效应
	if err := ra.applyAmplification(newResonances); err != nil {
		return err
	}

	// 处理反馈
	if err := ra.processFeedback(); err != nil {
		return err
	}

	return nil
}

// detectResonances 检测新的共振
func (ra *ResonanceAmplifier) detectResonances() ([]*ResonanceState, error) {
	// 获取当前模式
	patterns, err := ra.detector.Detect()
	if err != nil {
		return nil, err
	}

	resonances := make([]*ResonanceState, 0)

	// 分析模式间的共振可能
	for i, pattern1 := range patterns {
		for j := i + 1; j < len(patterns); j++ {
			pattern2 := patterns[j]

			// 检查是否存在共振关系
			if resonance := ra.checkResonance(pattern1, pattern2); resonance != nil {
				resonances = append(resonances, resonance)
			}
		}
	}

	return resonances, nil
}

// checkResonance 检查共振关系
func (ra *ResonanceAmplifier) checkResonance(
	pattern1, pattern2 emergence.EmergentPattern) *ResonanceState {

	// 计算模式间的相互作用
	interaction := ra.calculateInteraction(pattern1, pattern2)

	// 判断是否达到共振条件
	if !ra.isResonanceConditionMet(interaction) {
		return nil
	}

	// 创建共振状态
	resonance := &ResonanceState{
		ID:         generateResonanceID(),
		Type:       determineResonanceType(interaction),
		Source:     createResonanceSource(pattern1),
		Target:     createResonanceTarget(pattern2),
		StartTime:  time.Now(),
		LastUpdate: time.Now(),
	}

	// 初始化共振参数
	ra.initializeResonanceParameters(resonance, interaction)

	return resonance
}

// determineResonanceType 确定共振类型
func determineResonanceType(interaction float64) string {
	if interaction > 0.8 {
		return "strong_resonance"
	} else if interaction > 0.5 {
		return "medium_resonance"
	}
	return "weak_resonance"
}

// createResonanceSource 创建共振源
func createResonanceSource(pattern emergence.EmergentPattern) *ResonanceSource {
	return &ResonanceSource{
		Type:    pattern.Type,
		ID:      pattern.ID,
		Pattern: &pattern,
		Properties: map[string]float64{
			"energy":    pattern.Properties["energy"],
			"coherence": pattern.GetStructureCoherence(),
			"strength":  pattern.Strength,
		},
	}
}

// createResonanceTarget 创建共振目标
func createResonanceTarget(pattern emergence.EmergentPattern) *ResonanceTarget {
	return &ResonanceTarget{
		Type:    pattern.Type,
		ID:      pattern.ID,
		Pattern: &pattern,
		Properties: map[string]float64{
			"energy":    pattern.Properties["energy"],
			"coherence": pattern.GetStructureCoherence(),
			"strength":  pattern.Strength,
		},
	}
}

// initializeResonanceParameters 初始化共振参数
func (ra *ResonanceAmplifier) initializeResonanceParameters(resonance *ResonanceState, interaction float64) {
	resonance.Amplitude = interaction * ra.config.minAmplitude
	resonance.Phase = 0.0
	resonance.Frequency = 1.0
	resonance.Energy = resonance.Amplitude * resonance.Amplitude / 2.0
	resonance.Coherence = 1.0
}

// updateResonances 更新现有共振
func (ra *ResonanceAmplifier) updateResonances() error {
	currentTime := time.Now()

	for id, resonance := range ra.state.activeResonances {
		// 检查共振是否仍然有效
		if !ra.validateResonance(resonance) {
			delete(ra.state.activeResonances, id)
			continue
		}

		// 保存旧状态用于记录变化
		oldState := copyResonanceState(resonance)

		// 更新共振参数
		if err := ra.updateResonanceParameters(resonance); err != nil {
			continue
		}

		// 应用衰减
		ra.applyDecay(resonance)

		// 更新时间信息
		resonance.Duration = currentTime.Sub(resonance.StartTime)
		resonance.LastUpdate = currentTime

		// 记录状态变化
		ra.recordAmplificationEvent(resonance, oldState)
	}

	return nil
}

// validateResonance 验证共振是否有效
func (ra *ResonanceAmplifier) validateResonance(resonance *ResonanceState) bool {
	// 检查时间有效性
	if time.Since(resonance.LastUpdate) > types.MaxResonanceAge {
		return false
	}

	// 检查能量是否足够
	if resonance.Energy < ra.config.minAmplitude {
		return false
	}

	// 检查相干性是否足够
	if resonance.Coherence < 0.1 {
		return false
	}

	return true
}

// updateResonanceParameters 更新共振参数
func (ra *ResonanceAmplifier) updateResonanceParameters(resonance *ResonanceState) error {
	// 更新振幅
	amplificationFactor := ra.calculateAmplificationFactor(resonance)
	resonance.Amplitude *= amplificationFactor

	// 更新能量
	resonance.Energy = resonance.Amplitude * resonance.Amplitude / 2.0

	// 更新相位
	dt := time.Since(resonance.LastUpdate).Seconds()
	resonance.Phase += 2 * math.Pi * resonance.Frequency * dt
	resonance.Phase = math.Mod(resonance.Phase, 2*math.Pi)

	// 更新相干性
	resonance.Coherence *= math.Exp(-ra.config.decayRate * dt)

	return nil
}

// recordAmplificationEvent 记录放大事件
func (ra *ResonanceAmplifier) recordAmplificationEvent(
	resonance *ResonanceState,
	oldState *ResonanceState) {

	// 创建事件
	event := AmplificationEvent{
		Timestamp:   time.Now(),
		ResonanceID: resonance.ID,
		Type:        resonance.Type,
		OldState:    oldState,
		NewState:    copyResonanceState(resonance),
		Changes: map[string]float64{
			"amplitude": resonance.Amplitude - oldState.Amplitude,
			"energy":    resonance.Energy - oldState.Energy,
			"coherence": resonance.Coherence - oldState.Coherence,
		},
	}

	// 添加到历史记录
	ra.state.history = append(ra.state.history, event)

	// 限制历史记录长度
	if len(ra.state.history) > 1000 {
		ra.state.history = ra.state.history[1:]
	}
}

// applyAmplification 应用放大效应
func (ra *ResonanceAmplifier) applyAmplification(newResonances []*ResonanceState) error {
	// 处理新的共振
	for _, resonance := range newResonances {
		// 计算放大系数
		amplificationFactor := ra.calculateAmplificationFactor(resonance)

		// 应用放大效应
		if err := ra.amplifyResonance(resonance, amplificationFactor); err != nil {
			continue
		}

		// 添加到活跃共振中
		ra.state.activeResonances[resonance.ID] = resonance
	}

	return nil
}

// amplifyResonance 应用放大效应到共振状态
func (ra *ResonanceAmplifier) amplifyResonance(resonance *ResonanceState, factor float64) error {
	// 应用放大系数
	resonance.Amplitude *= factor
	resonance.Energy = resonance.Amplitude * resonance.Amplitude / 2.0

	// 更新源和目标的能量
	if resonance.Source != nil && resonance.Source.Pattern != nil {
		resonance.Source.Properties["energy"] += resonance.Energy * ra.config.feedbackRatio
	}

	if resonance.Target != nil && resonance.Target.Pattern != nil {
		resonance.Target.Properties["energy"] += resonance.Energy * ra.config.feedbackRatio
	}

	// 更新相干性
	resonance.Coherence = math.Min(1.0, resonance.Coherence+factor*0.1)

	// 更新时间
	resonance.LastUpdate = time.Now()

	return nil
}

// processFeedback 处理反馈
func (ra *ResonanceAmplifier) processFeedback() error {
	for _, resonance := range ra.state.activeResonances {
		// 计算反馈强度
		feedback := ra.calculateFeedback(resonance)

		// 应用反馈效应
		if err := ra.applyFeedback(resonance, feedback); err != nil {
			continue
		}
	}

	return nil
}

// calculateFeedback 计算反馈强度
func (ra *ResonanceAmplifier) calculateFeedback(resonance *ResonanceState) float64 {
	// 基础反馈强度
	baseFeedback := resonance.Energy * resonance.Coherence

	// 时间衰减
	timeFactor := math.Exp(-ra.config.decayRate *
		time.Since(resonance.StartTime).Seconds())

	// 相位调制
	phaseMod := math.Cos(resonance.Phase)

	// 综合反馈强度
	feedback := baseFeedback * timeFactor * (1.0 + phaseMod) / 2.0

	return feedback * ra.config.feedbackRatio
}

// applyFeedback 应用反馈效应
func (ra *ResonanceAmplifier) applyFeedback(resonance *ResonanceState, feedback float64) error {
	if resonance.Source == nil || resonance.Target == nil {
		return fmt.Errorf("invalid resonance state")
	}

	// 更新源模式
	if resonance.Source.Pattern != nil {
		resonance.Source.Pattern.Strength *= (1.0 + feedback)
		resonance.Source.Properties["energy"] += feedback * resonance.Energy
	}

	// 更新目标模式
	if resonance.Target.Pattern != nil {
		resonance.Target.Pattern.Strength *= (1.0 + feedback)
		resonance.Target.Properties["energy"] += feedback * resonance.Energy
	}

	// 更新共振状态
	resonance.Amplitude *= (1.0 + feedback*0.5)
	resonance.Energy = resonance.Amplitude * resonance.Amplitude / 2.0
	resonance.Coherence = math.Min(1.0, resonance.Coherence*(1.0+feedback*0.1))

	return nil
}

// 辅助函数

// 修改calculateInteraction使用新函数
func (ra *ResonanceAmplifier) calculateInteraction(
	pattern1, pattern2 emergence.EmergentPattern) float64 {

	// 计算模式强度乘积
	baseStrength := pattern1.Strength * pattern2.Strength

	// 计算相位差
	phaseDifference := calculatePatternPhaseDifference(pattern1, pattern2)

	// 计算频率匹配度
	frequencyMatch := calculateFrequencyMatch(pattern1, pattern2)

	return baseStrength * math.Cos(phaseDifference) * frequencyMatch
}

// calculatePatternPhaseDifference 计算模式相位差
func calculatePatternPhaseDifference(pattern1, pattern2 emergence.EmergentPattern) float64 {
	// 从属性中获取相位
	phase1 := 0.0
	if value, exists := pattern1.Properties["phase"]; exists {
		phase1 = value
	}

	phase2 := 0.0
	if value, exists := pattern2.Properties["phase"]; exists {
		phase2 = value
	}

	// 计算相位差并归一化
	diff := phase1 - phase2
	return normalizePhase(diff)
}

// calculateFrequencyMatch 计算频率匹配度
func calculateFrequencyMatch(pattern1, pattern2 emergence.EmergentPattern) float64 {
	// 从属性中获取频率
	freq1 := 1.0
	if value, exists := pattern1.Properties["frequency"]; exists {
		freq1 = value
	}

	freq2 := 1.0
	if value, exists := pattern2.Properties["frequency"]; exists {
		freq2 = value
	}

	// 计算频率比
	freqRatio := math.Min(freq1, freq2) / math.Max(freq1, freq2)
	return freqRatio
}

func (ra *ResonanceAmplifier) isResonanceConditionMet(interaction float64) bool {
	return interaction > ra.config.minAmplitude
}

func (ra *ResonanceAmplifier) calculateAmplificationFactor(resonance *ResonanceState) float64 {
	// 基础放大因子
	factor := resonance.Coherence * resonance.Energy

	// 应用配置限制
	factor = math.Max(ra.config.minAmplitude,
		math.Min(ra.config.maxAmplitude, factor))

	return factor
}

func (ra *ResonanceAmplifier) applyDecay(resonance *ResonanceState) {
	timeSinceUpdate := time.Since(resonance.LastUpdate).Seconds()
	decayFactor := math.Exp(-ra.config.decayRate * timeSinceUpdate)

	resonance.Amplitude *= decayFactor
	resonance.Energy *= decayFactor
}

func generateResonanceID() string {
	return fmt.Sprintf("res_%d", time.Now().UnixNano())
}

func copyResonanceState(state *ResonanceState) *ResonanceState {
	if state == nil {
		return nil
	}

	copy := *state
	return &copy
}

// Start 启动共振放大器
func (ra *ResonanceAmplifier) Start(ctx context.Context) error {
	ra.mu.Lock()
	defer ra.mu.Unlock()

	// 启动放大循环
	go ra.amplificationLoop(ctx)

	return nil
}

// Stop 停止共振放大器
func (ra *ResonanceAmplifier) Stop() error {
	ra.mu.Lock()
	defer ra.mu.Unlock()

	// 清理资源
	return nil
}

// amplificationLoop 放大循环
func (ra *ResonanceAmplifier) amplificationLoop(ctx context.Context) {
	ticker := time.NewTicker(ra.config.AmplificationInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ra.Amplify()
		}
	}
}
