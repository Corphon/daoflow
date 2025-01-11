// system/monitor/metrics/analyzer.go

package metrics

import (
    "context"
    "math"
    "sync"
    "time"

    "github.com/Corphon/daoflow/system/types"
)

// Analyzer 指标分析器
type Analyzer struct {
    mu sync.RWMutex

    // 基础配置
    config struct {
        WindowSize      int               // 分析窗口大小
        UpdateInterval  time.Duration     // 更新间隔
        Thresholds     map[string]float64 // 分析阈值
        Patterns       []types.Pattern    // 模式匹配规则
    }

    // 数据源
    collector *Collector
    
    // 分析结果缓存
    cache struct {
        lastAnalysis types.AnalysisResult
        history      []types.AnalysisResult
        patterns     map[string][]float64
    }

    // 分析状态
    status struct {
        isRunning    bool
        lastRun      time.Time
        errors       []error
    }
}

// NewAnalyzer 创建新的分析器
func NewAnalyzer(collector *Collector, config types.MetricsConfig) *Analyzer {
    return &Analyzer{
        collector: collector,
        cache: struct {
            lastAnalysis types.AnalysisResult
            history      []types.AnalysisResult
            patterns     map[string][]float64
        }{
            patterns: make(map[string][]float64),
        },
    }
}

// Start 启动分析器
func (a *Analyzer) Start(ctx context.Context) error {
    a.mu.Lock()
    if a.status.isRunning {
        a.mu.Unlock()
        return types.NewSystemError(types.ErrRuntime, "analyzer already running", nil)
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

    a.status.isRunning = false
    return nil
}

// analysisLoop 分析循环
func (a *Analyzer) analysisLoop(ctx context.Context) {
    ticker := time.NewTicker(a.config.UpdateInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            if err := a.analyze(ctx); err != nil {
                a.handleError(err)
            }
        }
    }
}

// analyze 执行指标分析
func (a *Analyzer) analyze(ctx context.Context) error {
    // 获取最新指标数据
    metrics := a.collector.GetCurrentMetrics()
    history := a.collector.GetMetricsHistory()

    // 创建分析结果
    result := types.AnalysisResult{
        Timestamp: time.Now(),
        Metrics:   metrics,
    }

    // 执行各类分析
    a.analyzeEvolutionTrends(&result, history)
    a.analyzeFieldDynamics(&result, history)
    a.analyzeQuantumStates(&result, history)
    a.analyzeEmergentPatterns(&result, history)
    a.predictFutureStates(&result, history)

    // 生成洞察和建议
    a.generateInsights(&result)

    // 缓存结果
    a.cacheResult(result)

    return nil
}

// analyzeEvolutionTrends 分析演化趋势
func (a *Analyzer) analyzeEvolutionTrends(result *types.AnalysisResult, history []types.MetricsData) {
    if len(history) < 2 {
        return
    }

    // 计算演化速度
    evolutionSpeeds := make([]float64, len(history)-1)
    for i := 1; i < len(history); i++ {
        current := history[i].System.Evolution.Level
        previous := history[i-1].System.Evolution.Level
        timeDiff := history[i].Timestamp.Sub(history[i-1].Timestamp).Seconds()
        evolutionSpeeds[i-1] = (current - previous) / timeDiff
    }

    // 分析趋势稳定性
    stability := a.calculateStability(evolutionSpeeds)
    result.EvolutionAnalysis = types.EvolutionAnalysis{
        AverageSpeed: a.calculateAverage(evolutionSpeeds),
        Stability:    stability,
        Trend:       a.determineTrend(evolutionSpeeds),
    }
}

// analyzeFieldDynamics 分析场动力学
func (a *Analyzer) analyzeFieldDynamics(result *types.AnalysisResult, history []types.MetricsData) {
    if len(history) < 2 {
        return
    }

    // 计算场强度变化
    fieldStrengths := make([]float64, len(history))
    for i, metrics := range history {
        fieldStrengths[i] = metrics.System.Field.Strength
    }

    // 分析场的稳定性和波动
    result.FieldAnalysis = types.FieldAnalysis{
        Stability:    a.calculateFieldStability(fieldStrengths),
        Oscillation:  a.detectOscillations(fieldStrengths),
        Coherence:    a.calculateFieldCoherence(history),
    }
}

// analyzeQuantumStates 分析量子态
func (a *Analyzer) analyzeQuantumStates(result *types.AnalysisResult, history []types.MetricsData) {
    quantumStates := make([]types.QuantumState, len(history))
    for i, metrics := range history {
        quantumStates[i] = metrics.System.Quantum
    }

    result.QuantumAnalysis = types.QuantumAnalysis{
        Entanglement: a.calculateEntanglement(quantumStates),
        Coherence:    a.calculateQuantumCoherence(quantumStates),
        Stability:    a.calculateQuantumStability(quantumStates),
    }
}

// analyzeEmergentPatterns 分析涌现模式
func (a *Analyzer) analyzeEmergentPatterns(result *types.AnalysisResult, history []types.MetricsData) {
    patterns := make([]types.EmergentPattern, 0)
    
    // 提取涌现模式
    for _, metrics := range history {
        if metrics.System.Emergence.Pattern != nil {
            patterns = append(patterns, metrics.System.Emergence.Pattern)
        }
    }

    // 分析模式特征
    result.EmergenceAnalysis = types.EmergenceAnalysis{
        Patterns:    a.identifyPatterns(patterns),
        Complexity: a.calculateComplexity(patterns),
        Stability:  a.calculatePatternStability(patterns),
    }
}

// predictFutureStates 预测未来状态
func (a *Analyzer) predictFutureStates(result *types.AnalysisResult, history []types.MetricsData) {
    if len(history) < a.config.WindowSize {
        return
    }

    // 使用时间序列分析预测未来状态
    predictions := types.Predictions{
        Energy:     a.predictEnergy(history),
        Field:      a.predictField(history),
        Quantum:    a.predictQuantum(history),
        Emergence:  a.predictEmergence(history),
    }

    result.Predictions = predictions
}

// generateInsights 生成洞察和建议
func (a *Analyzer) generateInsights(result *types.AnalysisResult) {
    insights := make([]types.Insight, 0)

    // 基于演化分析生成洞察
    if result.EvolutionAnalysis.Stability < a.config.Thresholds["min_stability"] {
        insights = append(insights, types.Insight{
            Type:    "evolution_stability",
            Level:   types.SeverityWarning,
            Message: "System evolution showing signs of instability",
            Recommendation: "Consider adjusting evolution parameters",
        })
    }

    // 基于场分析生成洞察
    if result.FieldAnalysis.Oscillation > a.config.Thresholds["max_oscillation"] {
        insights = append(insights, types.Insight{
            Type:    "field_oscillation",
            Level:   types.SeverityWarning,
            Message: "Excessive field oscillations detected",
            Recommendation: "Implement field dampening measures",
        })
    }

    result.Insights = insights
}

// 工具函数

func (a *Analyzer) calculateAverage(values []float64) float64 {
    if len(values) == 0 {
        return 0
    }
    sum := 0.0
    for _, v := range values {
        sum += v
    }
    return sum / float64(len(values))
}

func (a *Analyzer) calculateStability(values []float64) float64 {
    if len(values) < 2 {
        return 1.0
    }
    
    variance := 0.0
    mean := a.calculateAverage(values)
    
    for _, v := range values {
        diff := v - mean
        variance += diff * diff
    }
    
    variance /= float64(len(values))
    return 1.0 / (1.0 + math.Sqrt(variance))
}

func (a *Analyzer) determineTrend(values []float64) string {
    if len(values) < 2 {
        return "stable"
    }
    
    lastValues := values[len(values)-3:]
    increasing := 0
    decreasing := 0
    
    for i := 1; i < len(lastValues); i++ {
        if lastValues[i] > lastValues[i-1] {
            increasing++
        } else if lastValues[i] < lastValues[i-1] {
            decreasing++
        }
    }
    
    if increasing > decreasing {
        return "increasing"
    } else if decreasing > increasing {
        return "decreasing"
    }
    return "stable"
}

// handleError 处理错误
func (a *Analyzer) handleError(err error) {
    a.mu.Lock()
    defer a.mu.Unlock()

    a.status.errors = append(a.status.errors, err)
}

// cacheResult 缓存分析结果
func (a *Analyzer) cacheResult(result types.AnalysisResult) {
    a.mu.Lock()
    defer a.mu.Unlock()

    a.cache.lastAnalysis = result
    a.cache.history = append(a.cache.history, result)

    // 维护缓存大小
    if len(a.cache.history) > 100 {
        a.cache.history = a.cache.history[1:]
    }
}
