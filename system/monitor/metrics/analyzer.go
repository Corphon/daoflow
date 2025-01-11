// system/monitor/metrics/analyzer.go

package metrics

import (
    "context"
    "math"
    "sync"
    "time"

    "github.com/Corphon/daoflow/model"
    "github.com/Corphon/daoflow/system/types"
)

// Analyzer 指标分析器
type Analyzer struct {
    mu sync.RWMutex

    // 基础配置
    config types.MetricsConfig

    // 数据源
    collector *Collector
    
    // 分析结果缓存
    cache struct {
        lastAnalysis    *AnalysisResult
        history        []*AnalysisResult
        patterns       map[string][]float64
        modelMetrics   *model.ModelMetrics
    }

    // 分析状态
    status struct {
        isRunning    bool
        lastRun      time.Time
        errors       []error
    }
}

// AnalysisResult 分析结果
type AnalysisResult struct {
    ID        string
    Timestamp time.Time
    
    // 系统指标分析
    SystemMetrics  types.MetricsData
    
    // 模型指标分析
    ModelMetrics   model.ModelMetrics
    
    // 量子分析
    QuantumAnalysis struct {
        Entanglement float64
        Coherence    float64
        Stability    float64
        Phase        float64
    }
    
    // 场分析
    FieldAnalysis struct {
        Strength    float64
        Uniformity  float64
        Coupling    float64
        Resonance   float64
    }
    
    // 涌现分析
    EmergenceAnalysis struct {
        Patterns    []types.EmergentPattern
        Complexity  float64
        Stability   float64
        Potential   float64
    }
    
    // 预测
    Predictions struct {
        NextState       model.ModelState
        EnergyTrend    float64
        FieldEvolution []float64
        EmergenceProb  map[string]float64
    }
    
    // 洞察
    Insights []types.Insight
}

// NewAnalyzer 创建新的分析器
func NewAnalyzer(collector *Collector, config types.MetricsConfig) *Analyzer {
    return &Analyzer{
        collector: collector,
        config:    config,
        cache: struct {
            lastAnalysis    *AnalysisResult
            history        []*AnalysisResult
            patterns       map[string][]float64
            modelMetrics   *model.ModelMetrics
        }{
            patterns:     make(map[string][]float64),
            modelMetrics: &model.ModelMetrics{},
        },
    }
}

// Start 启动分析器
func (a *Analyzer) Start(ctx context.Context) error {
    a.mu.Lock()
    if a.status.isRunning {
        a.mu.Unlock()
        return model.WrapError(nil, model.ErrCodeOperation, "analyzer already running")
    }
    a.status.isRunning = true
    a.mu.Unlock()

    go a.analysisLoop(ctx)
    return nil
}

// Stop 停止分析器
func (a *Analyzer) Stop() error {
    a.mu.Lock()
    defer a.mu.Unlock()

    if !a.status.isRunning {
        return model.WrapError(nil, model.ErrCodeOperation, "analyzer not running")
    }

    a.status.isRunning = false
    return nil
}

// analyze 执行分析
func (a *Analyzer) analyze(ctx context.Context) error {
    // 获取最新指标数据
    metrics := a.collector.GetCurrentMetrics()
    modelMetrics := a.collector.GetModelMetrics()
    history := a.collector.GetMetricsHistory()

    // 创建分析结果
    result := &AnalysisResult{
        ID:           generateAnalysisID(),
        Timestamp:    time.Now(),
        SystemMetrics: metrics,
        ModelMetrics:  modelMetrics,
    }

    // 执行各类分析
    if err := a.analyzeQuantumStates(result); err != nil {
        return model.WrapError(err, model.ErrCodeOperation, "quantum analysis failed")
    }

    if err := a.analyzeFieldDynamics(result); err != nil {
        return model.WrapError(err, model.ErrCodeOperation, "field analysis failed")
    }

    if err := a.analyzeEmergentPatterns(result); err != nil {
        return model.WrapError(err, model.ErrCodeOperation, "emergence analysis failed")
    }

    if err := a.generatePredictions(result); err != nil {
        return model.WrapError(err, model.ErrCodeOperation, "prediction generation failed")
    }

    // 生成洞察
    if err := a.generateInsights(result); err != nil {
        return model.WrapError(err, model.ErrCodeOperation, "insight generation failed")
    }

    // 缓存结果
    a.cacheResult(result)

    return nil
}

// analyzeQuantumStates 分析量子态
func (a *Analyzer) analyzeQuantumStates(result *AnalysisResult) error {
    quantum := result.ModelMetrics.Quantum
    
    result.QuantumAnalysis.Entanglement = calculateEntanglement(quantum)
    result.QuantumAnalysis.Coherence = calculateCoherence(quantum)
    result.QuantumAnalysis.Stability = calculateQuantumStability(quantum)
    result.QuantumAnalysis.Phase = quantum.Phase

    return nil
}

// analyzeFieldDynamics 分析场动力学
func (a *Analyzer) analyzeFieldDynamics(result *AnalysisResult) error {
    field := result.ModelMetrics.Field
    
    result.FieldAnalysis.Strength = field.Strength
    result.FieldAnalysis.Uniformity = calculateFieldUniformity(field)
    result.FieldAnalysis.Coupling = calculateFieldCoupling(field)
    result.FieldAnalysis.Resonance = calculateResonance(field)

    return nil
}

// analyzeEmergentPatterns 分析涌现模式
func (a *Analyzer) analyzeEmergentPatterns(result *AnalysisResult) error {
    patterns := detectEmergentPatterns(result.SystemMetrics, result.ModelMetrics)
    
    result.EmergenceAnalysis.Patterns = patterns
    result.EmergenceAnalysis.Complexity = calculateComplexity(patterns)
    result.EmergenceAnalysis.Stability = calculatePatternStability(patterns)
    result.EmergenceAnalysis.Potential = calculateEmergencePotential(patterns)

    return nil
}

// generatePredictions 生成预测
func (a *Analyzer) generatePredictions(result *AnalysisResult) error {
    // 预测下一个模型状态
    nextState, err := predictNextState(result.ModelMetrics)
    if err != nil {
        return err
    }
    result.Predictions.NextState = nextState

    // 预测能量趋势
    result.Predictions.EnergyTrend = predictEnergyTrend(result.ModelMetrics)

    // 预测场演化
    result.Predictions.FieldEvolution = predictFieldEvolution(result.ModelMetrics)

    // 预测涌现概率
    result.Predictions.EmergenceProb = predictEmergenceProbabilities(result.EmergenceAnalysis)

    return nil
}

// generateInsights 生成洞察
func (a *Analyzer) generateInsights(result *AnalysisResult) error {
    insights := make([]types.Insight, 0)

    // 基于量子分析生成洞察
    if result.QuantumAnalysis.Coherence < a.config.Thresholds["min_coherence"] {
        insights = append(insights, types.Insight{
            Type:    "quantum_coherence",
            Level:   types.SeverityWarning,
            Message: "Low quantum coherence detected",
            Recommendation: "Consider adjusting quantum parameters",
        })
    }

    // 基于场分析生成洞察
    if result.FieldAnalysis.Stability < a.config.Thresholds["min_field_stability"] {
        insights = append(insights, types.Insight{
            Type:    "field_stability",
            Level:   types.SeverityWarning,
            Message: "Field instability detected",
            Recommendation: "Implement field stabilization measures",
        })
    }

    result.Insights = insights
    return nil
}

// 辅助函数...
func calculateEntanglement(quantum model.QuantumState) float64 {
    // 实现量子纠缠度计算
    return quantum.GetEntanglement()
}

func calculateCoherence(quantum model.QuantumState) float64 {
    // 实现相干性计算
    return quantum.GetCoherence()
}

func calculateQuantumStability(quantum model.QuantumState) float64 {
    // 实现量子态稳定性计算
    phase := quantum.GetPhase()
    amplitude := quantum.GetAmplitude()
    return (1.0 - math.Abs(math.Sin(phase))) * amplitude
}

func calculateFieldUniformity(field model.FieldState) float64 {
    // 实现场均匀性计算
    gradient := field.GetGradient()
    maxGradient := 0.0
    for _, g := range gradient {
        if math.Abs(g) > maxGradient {
            maxGradient = math.Abs(g)
        }
    }
    return 1.0 - (maxGradient / field.GetStrength())
}

func calculateFieldCoupling(field model.FieldState) float64 {
    // 实现场耦合强度计算
    return field.GetCoupling()
}

func calculateResonance(field model.FieldState) float64 {
    // 实现共振强度计算
    return field.GetResonance()
}

func detectEmergentPatterns(systemMetrics types.MetricsData, modelMetrics model.ModelMetrics) []types.EmergentPattern {
    patterns := make([]types.EmergentPattern, 0)
    
    // 检测系统级涌现模式
    systemPatterns := detectSystemPatterns(systemMetrics)
    patterns = append(patterns, systemPatterns...)
    
    // 检测模型级涌现模式
    modelPatterns := detectModelPatterns(modelMetrics)
    patterns = append(patterns, modelPatterns...)
    
    return patterns
}

func calculateComplexity(patterns []types.EmergentPattern) float64 {
    if len(patterns) == 0 {
        return 0.0
    }
    
    totalComplexity := 0.0
    for _, pattern := range patterns {
        totalComplexity += pattern.Complexity
    }
    return totalComplexity / float64(len(patterns))
}

func calculatePatternStability(patterns []types.EmergentPattern) float64 {
    if len(patterns) == 0 {
        return 1.0
    }
    
    totalStability := 0.0
    for _, pattern := range patterns {
        totalStability += pattern.Stability
    }
    return totalStability / float64(len(patterns))
}

func calculateEmergencePotential(patterns []types.EmergentPattern) float64 {
    // 计算涌现潜力
    potential := 0.0
    weights := map[string]float64{
        "complexity": 0.3,
        "stability":  0.3,
        "coupling":   0.4,
    }
    
    for _, pattern := range patterns {
        weightedSum := pattern.Complexity * weights["complexity"] +
                      pattern.Stability * weights["stability"] +
                      pattern.Coupling * weights["coupling"]
        potential += weightedSum
    }
    
    return math.Min(1.0, potential/float64(len(patterns)))
}

func predictNextState(metrics model.ModelMetrics) (model.ModelState, error) {
    // 使用当前指标预测下一个状态
    predictor := model.NewStatePredictor()
    return predictor.PredictNext(metrics)
}

func predictEnergyTrend(metrics model.ModelMetrics) float64 {
    // 预测能量趋势
    currentEnergy := metrics.GetTotalEnergy()
    previousEnergy := metrics.GetPreviousEnergy()
    return (currentEnergy - previousEnergy) / previousEnergy
}

func predictFieldEvolution(metrics model.ModelMetrics) []float64 {
    // 预测场演化序列
    evolution := make([]float64, 10) // 预测未来10个时间步
    currentField := metrics.Field.GetStrength()
    
    for i := range evolution {
        // 简单线性预测示例
        evolution[i] = currentField * (1 + float64(i)*0.1)
    }
    
    return evolution
}

func predictEmergenceProbabilities(analysis struct {
    Patterns    []types.EmergentPattern
    Complexity  float64
    Stability   float64
    Potential   float64
}) map[string]float64 {
    probs := make(map[string]float64)
    
    // 基于当前分析预测各类涌现模式的概率
    for _, pattern := range analysis.Patterns {
        probability := calculatePatternProbability(pattern, analysis.Complexity, analysis.Stability)
        probs[pattern.Type] = probability
    }
    
    return probs
}

func calculatePatternProbability(pattern types.EmergentPattern, complexity, stability float64) float64 {
    // 基于模式特征、复杂度和稳定性计算概率
    baseProbability := pattern.Strength * stability
    adjustedProbability := baseProbability * (1 - complexity/2) // 复杂度越高，概率越低
    
    return math.Max(0, math.Min(1, adjustedProbability))
}

// cacheResult 缓存分析结果
func (a *Analyzer) cacheResult(result *AnalysisResult) {
    a.mu.Lock()
    defer a.mu.Unlock()

    a.cache.lastAnalysis = result
    a.cache.history = append(a.cache.history, result)
    a.cache.modelMetrics = &result.ModelMetrics

    // 维护历史大小
    if len(a.cache.history) > a.config.MaxHistorySize {
        a.cache.history = a.cache.history[1:]
    }
}

// GetLastAnalysis 获取最新分析结果
func (a *Analyzer) GetLastAnalysis() *AnalysisResult {
    a.mu.RLock()
    defer a.mu.RUnlock()
    return a.cache.lastAnalysis
}

// GetAnalysisHistory 获取分析历史
func (a *Analyzer) GetAnalysisHistory(limit int) []*AnalysisResult {
    a.mu.RLock()
    defer a.mu.RUnlock()

    if limit <= 0 || limit > len(a.cache.history) {
        limit = len(a.cache.history)
    }

    history := make([]*AnalysisResult, limit)
    copy(history, a.cache.history[len(a.cache.history)-limit:])
    return history
}

// generateAnalysisID 生成分析ID
func generateAnalysisID() string {
    return fmt.Sprintf("analysis-%d", time.Now().UnixNano())
}
