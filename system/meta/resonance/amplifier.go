//system/resonance/amplifier.go

package resonance

import (
    "math"
    "sync"
    "time"

    "github.com/Corphon/daoflow/meta/field"
    "github.com/Corphon/daoflow/model"
    "github.com/Corphon/daoflow/system/types"
)

// ResonanceAmplifier 共振放大器
type ResonanceAmplifier struct {
    mu sync.RWMutex

    // 基础配置
    config struct {
        gainFactor      float64         // 增益系数
        maxAmplitude    float64         // 最大幅度
        feedbackRatio   float64         // 反馈比例
        stabilityThreshold float64      // 稳定性阈值
    }

    // 放大状态
    state struct {
        resonances    map[string]*AmplifiedResonance  // 被放大的共振
        feedback     []FeedbackLoop                  // 反馈回路
        metrics      AmplifierMetrics                // 放大器指标
    }

    // 场引用
    field *field.UnifiedField
}

// AmplifiedResonance 被放大的共振
type AmplifiedResonance struct {
    ID           string                // 共振ID
    SourceType   string                // 源类型
    Frequency    float64               // 频率
    Phase        float64               // 相位
    Amplitude    float64               // 振幅
    Gain         float64               // 增益
    Stability    float64               // 稳定性
    Created      time.Time             // 创建时间
    LastUpdate   time.Time             // 最后更新时间
}

// FeedbackLoop 反馈回路
type FeedbackLoop struct {
    ID           string                // 回路ID
    Source       string                // 源端点
    Target       string                // 目标端点
    Ratio        float64               // 反馈比例
    Delay        time.Duration         // 反馈延迟
    Active       bool                  // 是否激活
}

// AmplifierMetrics 放大器指标
type AmplifierMetrics struct {
    TotalGain     float64              // 总增益
    Efficiency    float64              // 效率
    Stability     float64              // 稳定性
    History       []MetricPoint        // 历史指标
}

// MetricPoint 指标点
type MetricPoint struct {
    Timestamp    time.Time
    Values       map[string]float64
}

// NewResonanceAmplifier 创建共振放大器
func NewResonanceAmplifier(field *field.UnifiedField) *ResonanceAmplifier {
    ra := &ResonanceAmplifier{
        field: field,
    }

    // 初始化配置
    ra.config.gainFactor = 1.5
    ra.config.maxAmplitude = 10.0
    ra.config.feedbackRatio = 0.2
    ra.config.stabilityThreshold = 0.7

    // 初始化状态
    ra.state.resonances = make(map[string]*AmplifiedResonance)
    ra.state.feedback = make([]FeedbackLoop, 0)
    ra.state.metrics = AmplifierMetrics{
        History: make([]MetricPoint, 0),
    }

    return ra
}

// Amplify 放大共振
func (ra *ResonanceAmplifier) Amplify(resonance *AmplifiedResonance) error {
    ra.mu.Lock()
    defer ra.mu.Unlock()

    // 验证共振参数
    if err := ra.validateResonance(resonance); err != nil {
        return err
    }

    // 计算增益
    gain := ra.calculateGain(resonance)

    // 应用增益
    if err := ra.applyGain(resonance, gain); err != nil {
        return err
    }

    // 添加反馈
    if err := ra.addFeedback(resonance); err != nil {
        return err
    }

    // 更新状态
    ra.updateState(resonance)

    return nil
}

// Tune 调谐共振
func (ra *ResonanceAmplifier) Tune(id string, frequency, phase float64) error {
    ra.mu.Lock()
    defer ra.mu.Unlock()

    resonance, exists := ra.state.resonances[id]
    if !exists {
        return model.WrapError(nil, model.ErrCodeNotFound, "resonance not found")
    }

    // 调整频率
    if err := ra.tuneFrequency(resonance, frequency); err != nil {
        return err
    }

    // 调整相位
    if err := ra.tunePhase(resonance, phase); err != nil {
        return err
    }

    // 重新计算稳定性
    stability := ra.calculateStability(resonance)
    resonance.Stability = stability

    // 更新时间戳
    resonance.LastUpdate = time.Now()

    return nil
}

// AddFeedbackLoop 添加反馈回路
func (ra *ResonanceAmplifier) AddFeedbackLoop(loop FeedbackLoop) error {
    ra.mu.Lock()
    defer ra.mu.Unlock()

    // 验证回路
    if err := ra.validateFeedbackLoop(loop); err != nil {
        return err
    }

    // 添加回路
    ra.state.feedback = append(ra.state.feedback, loop)

    return nil
}

// 内部方法

func (ra *ResonanceAmplifier) validateResonance(resonance *AmplifiedResonance) error {
    if resonance == nil {
        return model.WrapError(nil, model.ErrCodeValidation, "nil resonance")
    }

    if resonance.Frequency <= 0 {
        return model.WrapError(nil, model.ErrCodeValidation, "invalid frequency")
    }

    if math.Abs(resonance.Phase) > math.Pi {
        return model.WrapError(nil, model.ErrCodeValidation, "invalid phase")
    }

    return nil
}

func (ra *ResonanceAmplifier) calculateGain(resonance *AmplifiedResonance) float64 {
    // 基础增益
    baseGain := ra.config.gainFactor * resonance.Amplitude

    // 考虑稳定性
    stabilityFactor := math.Sqrt(resonance.Stability)
    
    // 反馈调整
    feedbackFactor := ra.calculateFeedbackFactor(resonance)

    // 最终增益
    gain := baseGain * stabilityFactor * feedbackFactor

    // 限制最大值
    return math.Min(gain, ra.config.maxAmplitude)
}

func (ra *ResonanceAmplifier) applyGain(resonance *AmplifiedResonance, gain float64) error {
    // 检查稳定性
    if resonance.Stability < ra.config.stabilityThreshold {
        return model.WrapError(nil, model.ErrCodeValidation, "resonance not stable enough")
    }

    // 应用增益
    resonance.Gain = gain
    resonance.Amplitude *= gain

    // 更新场状态
    if err := ra.updateFieldState(resonance); err != nil {
        return err
    }

    return nil
}

func (ra *ResonanceAmplifier) calculateFeedbackFactor(resonance *AmplifiedResonance) float64 {
    totalFeedback := 0.0
    activeFeedback := 0

    // 计算所有活跃反馈的贡献
    for _, loop := range ra.state.feedback {
        if !loop.Active {
            continue
        }
        if loop.Target == resonance.ID {
            totalFeedback += loop.Ratio
            activeFeedback++
        }
    }

    if activeFeedback == 0 {
        return 1.0
    }

    // 平均反馈因子
    return 1.0 + (totalFeedback / float64(activeFeedback))
}

func (ra *ResonanceAmplifier) updateFieldState(resonance *AmplifiedResonance) error {
    // 更新场能量
    energy := resonance.Amplitude * resonance.Amplitude
    if err := ra.field.AddEnergy(energy); err != nil {
        return err
    }

    // 更新场相位
    if err := ra.field.SetPhase(resonance.Phase); err != nil {
        return err
    }

    return nil
}

func (ra *ResonanceAmplifier) calculateStability(resonance *AmplifiedResonance) float64 {
    // 振幅稳定性
    amplitudeStability := 1.0 - (resonance.Amplitude / ra.config.maxAmplitude)
    
    // 相位稳定性
    phaseStability := math.Cos(resonance.Phase)
    
    // 增益稳定性
    gainStability := 1.0 - (resonance.Gain / (2 * ra.config.gainFactor))

    // 综合稳定性
    return (amplitudeStability + math.Abs(phaseStability) + gainStability) / 3.0
}

func (ra *ResonanceAmplifier) validateFeedbackLoop(loop FeedbackLoop) error {
    if loop.Source == "" || loop.Target == "" {
        return model.WrapError(nil, model.ErrCodeValidation, "invalid source or target")
    }

    if loop.Ratio < 0 || loop.Ratio > 1 {
        return model.WrapError(nil, model.ErrCodeValidation, "invalid feedback ratio")
    }

    return nil
}

func (ra *ResonanceAmplifier) updateState(resonance *AmplifiedResonance) {
    // 更新共振状态
    ra.state.resonances[resonance.ID] = resonance

    // 更新指标
    ra.updateMetrics()
}

func (ra *ResonanceAmplifier) updateMetrics() {
    metrics := &ra.state.metrics
    
    // 计算总增益
    totalGain := 0.0
    for _, res := range ra.state.resonances {
        totalGain += res.Gain
    }
    metrics.TotalGain = totalGain

    // 计算效率
    metrics.Efficiency = ra.calculateEfficiency()

    // 计算稳定性
    metrics.Stability = ra.calculateOverallStability()

    // 记录指标点
    point := MetricPoint{
        Timestamp: time.Now(),
        Values: map[string]float64{
            "total_gain":  metrics.TotalGain,
            "efficiency": metrics.Efficiency,
            "stability": metrics.Stability,
        },
    }
    metrics.History = append(metrics.History, point)

    // 限制历史记录长度
    if len(metrics.History) > maxMetricsHistory {
        metrics.History = metrics.History[1:]
    }
}

// 帮助函数

func (ra *ResonanceAmplifier) calculateEfficiency() float64 {
    if len(ra.state.resonances) == 0 {
        return 1.0
    }

    totalEfficiency := 0.0
    for _, res := range ra.state.resonances {
        // 计算单个共振的效率
        efficiency := res.Amplitude / (res.Gain * res.Amplitude)
        totalEfficiency += efficiency
    }

    return totalEfficiency / float64(len(ra.state.resonances))
}

func (ra *ResonanceAmplifier) calculateOverallStability() float64 {
    if len(ra.state.resonances) == 0 {
        return 1.0
    }

    totalStability := 0.0
    for _, res := range ra.state.resonances {
        totalStability += res.Stability
    }

    return totalStability / float64(len(ra.state.resonances))
}

const (
    maxMetricsHistory = 1000
)
