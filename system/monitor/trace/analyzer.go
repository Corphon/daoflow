// system/monitor/trace/analyzer.go

package trace

import (
    "context"
    "sync"
    "time"

    "github.com/Corphon/daoflow/model"
    "github.com/Corphon/daoflow/system/types"
)

// TraceAnalysis 追踪分析结果
type TraceAnalysis struct {
    ID          string                 
    Timestamp   time.Time             
    TraceID     TraceID               
    Duration    time.Duration         
    SpanCount   int                   
    
    // 系统层面分析
    Patterns    []types.TracePattern  
    Bottlenecks []types.Bottleneck    
    Metrics     map[string]float64    
    Anomalies   []types.Anomaly      

    // 模型层面分析
    ModelAnalysis struct {
        State       model.ModelState
        Flow        model.FlowState
        Patterns    []model.FlowPattern
        Metrics     model.ModelMetrics
        Anomalies   []model.Anomaly
    }

    // 量子层面分析
    QuantumAnalysis struct {
        Entanglement float64
        Coherence    float64
        Phase        float64
        States      []model.QuantumState
    }

    // 场动力学分析
    FieldAnalysis struct {
        Strength    float64
        Coupling    float64
        Resonance   float64
        Evolution   []model.FieldState
    }
}

// Analyzer 追踪分析器
type Analyzer struct {
    mu sync.RWMutex

    // 配置
    config types.TraceConfig

    // 数据源
    tracker  *Tracker
    recorder *Recorder

    // 分析缓存
    cache struct {
        traces    map[TraceID]*TraceAnalysis
        patterns  []types.TracePattern
        anomalies []types.Anomaly
    }

    // 分析状态
    status struct {
        isRunning    bool
        lastAnalysis time.Time
        errors       []error
    }

    // 模型分析器
    modelAnalyzer *model.Analyzer
}

// NewAnalyzer 创建新的分析器
func NewAnalyzer(tracker *Tracker, recorder *Recorder, config types.TraceConfig) *Analyzer {
    return &Analyzer{
        tracker:      tracker,
        recorder:    recorder,
        config:      config,
        modelAnalyzer: model.NewAnalyzer(),
        cache: struct {
            traces    map[TraceID]*TraceAnalysis
            patterns  []types.TracePattern
            anomalies []types.Anomaly
        }{
            traces: make(map[TraceID]*TraceAnalysis),
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
    // 获取追踪数据
    traces := a.getTracesInWindow()

    for traceID, spans := range traces {
        analysis := &TraceAnalysis{
            ID:        generateAnalysisID(),
            Timestamp: time.Now(),
            TraceID:   traceID,
        }

        // 系统层面分析
        if err := a.analyzeSystemTrace(analysis, spans); err != nil {
            return model.WrapError(err, model.ErrCodeOperation, "system analysis failed")
        }

        // 模型层面分析
        if err := a.analyzeModelTrace(analysis, spans); err != nil {
            return model.WrapError(err, model.ErrCodeOperation, "model analysis failed")
        }

        // 量子层面分析
        if err := a.analyzeQuantumTrace(analysis, spans); err != nil {
            return model.WrapError(err, model.ErrCodeOperation, "quantum analysis failed")
        }

        // 场动力学分析
        if err := a.analyzeFieldTrace(analysis, spans); err != nil {
            return model.WrapError(err, model.ErrCodeOperation, "field analysis failed")
        }

        // 缓存分析结果
        a.cacheAnalysis(analysis)
    }

    return nil
}

// analyzeSystemTrace 分析系统层面的追踪
func (a *Analyzer) analyzeSystemTrace(analysis *TraceAnalysis, spans []*Span) error {
    // 检测系统模式
    patterns := a.detectSystemPatterns(spans)
    analysis.Patterns = patterns

    // 检测瓶颈
    bottlenecks := a.detectBottlenecks(spans)
    analysis.Bottlenecks = bottlenecks

    // 计算指标
    metrics := a.calculateSystemMetrics(spans)
    analysis.Metrics = metrics

    // 检测异常
    anomalies := a.detectSystemAnomalies(spans, patterns)
    analysis.Anomalies = anomalies

    return nil
}

// analyzeModelTrace 分析模型层面的追踪
func (a *Analyzer) analyzeModelTrace(analysis *TraceAnalysis, spans []*Span) error {
    modelSpans := a.filterModelSpans(spans)
    if len(modelSpans) == 0 {
        return nil
    }

    // 分析模型状态
    state := a.analyzeModelState(modelSpans)
    analysis.ModelAnalysis.State = state

    // 分析流状态
    flow := a.analyzeModelFlow(modelSpans)
    analysis.ModelAnalysis.Flow = flow

    // 检测模型模式
    patterns := a.modelAnalyzer.DetectPatterns(modelSpans)
    analysis.ModelAnalysis.Patterns = patterns

    // 计算模型指标
    metrics := a.modelAnalyzer.CalculateMetrics(modelSpans)
    analysis.ModelAnalysis.Metrics = metrics

    // 检测模型异常
    anomalies := a.modelAnalyzer.DetectAnomalies(modelSpans)
    analysis.ModelAnalysis.Anomalies = anomalies

    return nil
}

// analyzeQuantumTrace 分析量子层面的追踪
func (a *Analyzer) analyzeQuantumTrace(analysis *TraceAnalysis, spans []*Span) error {
    // 提取量子态相关的跨度
    quantumSpans := a.filterQuantumSpans(spans)
    if len(quantumSpans) == 0 {
        return nil
    }

    // 分析量子纠缠
    entanglement := a.calculateEntanglement(quantumSpans)
    analysis.QuantumAnalysis.Entanglement = entanglement

    // 分析相干性
    coherence := a.calculateCoherence(quantumSpans)
    analysis.QuantumAnalysis.Coherence = coherence

    // 分析相位
    phase := a.calculatePhase(quantumSpans)
    analysis.QuantumAnalysis.Phase = phase

    // 提取量子态序列
    states := a.extractQuantumStates(quantumSpans)
    analysis.QuantumAnalysis.States = states

    return nil
}

// analyzeFieldTrace 分析场动力学追踪
func (a *Analyzer) analyzeFieldTrace(analysis *TraceAnalysis, spans []*Span) error {
    // 提取场相关的跨度
    fieldSpans := a.filterFieldSpans(spans)
    if len(fieldSpans) == 0 {
        return nil
    }

    // 分析场强度
    strength := a.calculateFieldStrength(fieldSpans)
    analysis.FieldAnalysis.Strength = strength

    // 分析场耦合
    coupling := a.calculateFieldCoupling(fieldSpans)
    analysis.FieldAnalysis.Coupling = coupling

    // 分析共振
    resonance := a.calculateResonance(fieldSpans)
    analysis.FieldAnalysis.Resonance = resonance

    // 提取场态演化序列
    evolution := a.extractFieldEvolution(fieldSpans)
    analysis.FieldAnalysis.Evolution = evolution

    return nil
}

// 过滤器方法
func (a *Analyzer) filterModelSpans(spans []*Span) []*Span {
    modelSpans := make([]*Span, 0)
    for _, span := range spans {
        if span.ModelType != model.ModelTypeNone {
            modelSpans = append(modelSpans, span)
        }
    }
    return modelSpans
}

func (a *Analyzer) filterQuantumSpans(spans []*Span) []*Span {
    quantumSpans := make([]*Span, 0)
    for _, span := range spans {
        if _, ok := span.Fields["quantum_state"]; ok {
            quantumSpans = append(quantumSpans, span)
        }
    }
    return quantumSpans
}

func (a *Analyzer) filterFieldSpans(spans []*Span) []*Span {
    fieldSpans := make([]*Span, 0)
    for _, span := range spans {
        if _, ok := span.Fields["field_state"]; ok {
            fieldSpans = append(fieldSpans, span)
        }
    }
    return fieldSpans
}

// 缓存方法
func (a *Analyzer) cacheAnalysis(analysis *TraceAnalysis) {
    a.mu.Lock()
    defer a.mu.Unlock()

    a.cache.traces[analysis.TraceID] = analysis
    a.status.lastAnalysis = analysis.Timestamp
}

// 辅助方法
func (a *Analyzer) calculateEntanglement(spans []*Span) float64 {
    if len(spans) < 2 {
        return 0.0
    }

    var totalEntanglement float64
    pairCount := 0

    // 计算所有量子态对之间的纠缠度
    for i := 0; i < len(spans)-1; i++ {
        for j := i + 1; j < len(spans); j++ {
            state1, ok1 := spans[i].Fields["quantum_state"].(model.QuantumState)
            state2, ok2 := spans[j].Fields["quantum_state"].(model.QuantumState)
            
            if !ok1 || !ok2 {
                continue
            }

            // 计算两个量子态之间的纠缠度
            entanglement := calculatePairEntanglement(state1, state2)
            totalEntanglement += entanglement
            pairCount++
        }
    }

    if pairCount == 0 {
        return 0.0
    }

    return totalEntanglement / float64(pairCount)
}

// calculatePairEntanglement 计算两个量子态之间的纠缠度
func calculatePairEntanglement(state1, state2 model.QuantumState) float64 {
    // 计算态矢量的内积
    overlap := state1.DotProduct(state2)
    
    // 计算纠缠度：使用冯诺依曼熵的简化形式
    // E = -Tr(ρ log₂ ρ)，其中 ρ 是约化密度矩阵
    p := math.Abs(overlap)
    if p == 0 || p == 1 {
        return 0.0
    }
    
    entropy := -p*math.Log2(p) - (1-p)*math.Log2(1-p)
    return entropy
}

func (a *Analyzer) calculateCoherence(spans []*Span) float64 {
    if len(spans) == 0 {
        return 0.0
    }

    var totalCoherence float64
    validSpans := 0

    for _, span := range spans {
        state, ok := span.Fields["quantum_state"].(model.QuantumState)
        if !ok {
            continue
        }

        // 计算相干性：使用态的纯度作为相干性度量
        // C = |⟨ψ|ψ⟩|²
        purity := state.CalculatePurity()
        phaseStability := calculatePhaseStability(state)
        
        // 综合考虑纯度和相位稳定性
        coherence := math.Sqrt(purity * phaseStability)
        totalCoherence += coherence
        validSpans++
    }

    if validSpans == 0 {
        return 0.0
    }

    return totalCoherence / float64(validSpans)
}

func (a *Analyzer) calculatePhase(spans []*Span) float64 {
    if len(spans) == 0 {
        return 0.0
    }

    var phases []float64
    for _, span := range spans {
        state, ok := span.Fields["quantum_state"].(model.QuantumState)
        if !ok {
            continue
        }

        // 获取量子态的相位
        phase := state.GetPhase()
        phases = append(phases, phase)
    }

    if len(phases) == 0 {
        return 0.0
    }

    // 计算平均相位和相位相干性
    avgPhase := calculateAveragePhase(phases)
    phaseCoherence := calculatePhaseCoherence(phases, avgPhase)

    return normalizePhase(avgPhase * phaseCoherence)
}

func (a *Analyzer) calculateFieldStrength(spans []*Span) float64 {
    if len(spans) == 0 {
        return 0.0
    }

    var totalStrength float64
    var weights float64

    for _, span := range spans {
        state, ok := span.Fields["field_state"].(model.FieldState)
        if !ok {
            continue
        }

        // 获取基本场强
        baseStrength := state.GetStrength()

        // 考虑时间衰减
        timeFactor := calculateTimeFactor(span.StartTime, span.EndTime)
        
        // 考虑空间分布
        spatialFactor := calculateSpatialFactor(state.GetDistribution())

        // 计算加权场强
        weightedStrength := baseStrength * timeFactor * spatialFactor
        
        // 累加加权和
        weight := span.Duration.Seconds()
        totalStrength += weightedStrength * weight
        weights += weight
    }

    if weights == 0 {
        return 0.0
    }

    return totalStrength / weights
}

func (a *Analyzer) calculateFieldCoupling(spans []*Span) float64 {
    if len(spans) < 2 {
        return 0.0
    }

    var totalCoupling float64
    var couplingCount int

    // 计算场之间的耦合强度
    for i := 0; i < len(spans)-1; i++ {
        for j := i + 1; j < len(spans); j++ {
            field1, ok1 := spans[i].Fields["field_state"].(model.FieldState)
            field2, ok2 := spans[j].Fields["field_state"].(model.FieldState)
            
            if !ok1 || !ok2 {
                continue
            }

            // 计算场间耦合
            coupling := calculateFieldInteraction(field1, field2)
            
            // 考虑时空关联
            spacetimeFactor := calculateSpacetimeCorrelation(spans[i], spans[j])
            
            totalCoupling += coupling * spacetimeFactor
            couplingCount++
        }
    }

    if couplingCount == 0 {
        return 0.0
    }

    return totalCoupling / float64(couplingCount)
}

func (a *Analyzer) calculateResonance(spans []*Span) float64 {
    if len(spans) == 0 {
        return 0.0
    }

    // 提取场的频率特征
    frequencies := make([]float64, 0)
    amplitudes := make([]float64, 0)

    for _, span := range spans {
        field, ok := spans[i].Fields["field_state"].(model.FieldState)
        if !ok {
            continue
        }

        freq := field.GetFrequency()
        amp := field.GetAmplitude()
        
        frequencies = append(frequencies, freq)
        amplitudes = append(amplitudes, amp)
    }

    if len(frequencies) == 0 {
        return 0.0
    }

    // 计算共振强度
    resonanceStrength := calculateResonanceStrength(frequencies, amplitudes)
    
    // 计算共振稳定性
    resonanceStability := calculateResonanceStability(frequencies)
    
    // 综合评估
    return math.Sqrt(resonanceStrength * resonanceStability)
}

// 辅助函数

func calculatePhaseStability(state model.QuantumState) float64 {
    initialPhase := state.GetPhase()
    phaseVariation := state.GetPhaseVariation()
    return math.Exp(-phaseVariation * phaseVariation)
}

func calculateAveragePhase(phases []float64) float64 {
    if len(phases) == 0 {
        return 0.0
    }
    
    sumSin := 0.0
    sumCos := 0.0
    for _, phase := range phases {
        sumSin += math.Sin(phase)
        sumCos += math.Cos(phase)
    }
    
    return math.Atan2(sumSin, sumCos)
}

func calculatePhaseCoherence(phases []float64, avgPhase float64) float64 {
    if len(phases) == 0 {
        return 0.0
    }
    
    var sumDeviation float64
    for _, phase := range phases {
        deviation := math.Abs(normalizePhase(phase - avgPhase))
        sumDeviation += deviation * deviation
    }
    
    return math.Exp(-sumDeviation / float64(len(phases)))
}

func normalizePhase(phase float64) float64 {
    // 将相位标准化到 [-π, π] 区间
    normalized := math.Mod(phase, 2*math.Pi)
    if normalized > math.Pi {
        normalized -= 2 * math.Pi
    } else if normalized < -math.Pi {
        normalized += 2 * math.Pi
    }
    return normalized
}

func calculateTimeFactor(start, end time.Time) float64 {
    duration := end.Sub(start).Seconds()
    return math.Exp(-duration / 3600.0) // 1小时特征时间
}

func calculateSpatialFactor(distribution []float64) float64 {
    if len(distribution) == 0 {
        return 1.0
    }
    
    // 计算空间分布的均匀性
    var sum, sumSq float64
    for _, value := range distribution {
        sum += value
        sumSq += value * value
    }
    
    mean := sum / float64(len(distribution))
    variance := sumSq/float64(len(distribution)) - mean*mean
    
    return 1.0 / (1.0 + variance)
}

func calculateFieldInteraction(field1, field2 model.FieldState) float64 {
    // 计算场之间的相互作用强度
    overlap := field1.CalculateOverlap(field2)
    strength := math.Sqrt(field1.GetStrength() * field2.GetStrength())
    return overlap * strength
}

func calculateSpacetimeCorrelation(span1, span2 *Span) float64 {
    // 计算时空关联度
    timeCorr := calculateTimeCorrelation(span1.StartTime, span1.EndTime, 
                                       span2.StartTime, span2.EndTime)
    spaceCorr := calculateSpaceCorrelation(span1, span2)
    return math.Sqrt(timeCorr * spaceCorr)
}

func calculateResonanceStrength(frequencies, amplitudes []float64) float64 {
    if len(frequencies) != len(amplitudes) || len(frequencies) == 0 {
        return 0.0
    }
    
    // 计算频率匹配度和振幅增强
    var resonanceSum float64
    for i := 0; i < len(frequencies)-1; i++ {
        for j := i + 1; j < len(frequencies); j++ {
            freqMatch := calculateFrequencyMatch(frequencies[i], frequencies[j])
            ampProduct := amplitudes[i] * amplitudes[j]
            resonanceSum += freqMatch * ampProduct
        }
    }
    
    return resonanceSum / float64(len(frequencies)*(len(frequencies)-1)/2)
}

func calculateFrequencyMatch(f1, f2 float64) float64 {
    // 计算频率匹配度，使用高斯函数
    diff := math.Abs(f1 - f2)
    return math.Exp(-diff * diff / (2.0 * 0.1)) // 0.1是带宽参数
}

func calculateResonanceStability(frequencies []float64) float64 {
    if len(frequencies) < 2 {
        return 1.0
    }
    
    // 计算频率的稳定性
    var sum, sumSq float64
    for _, f := range frequencies {
        sum += f
        sumSq += f * f
    }
    
    mean := sum / float64(len(frequencies))
    variance := sumSq/float64(len(frequencies)) - mean*mean
    
    return 1.0 / (1.0 + variance)
}
