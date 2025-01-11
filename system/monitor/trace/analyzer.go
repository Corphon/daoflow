// system/monitor/trace/analyzer.go

package trace

import (
    "context"
    "sync"
    "time"

    "github.com/Corphon/daoflow/system/types"
)

// TraceAnalysis 追踪分析结果
type TraceAnalysis struct {
    ID          string                 `json:"id"`
    Timestamp   time.Time             `json:"timestamp"`
    TraceID     TraceID               `json:"trace_id"`
    Duration    time.Duration         `json:"duration"`
    SpanCount   int                   `json:"span_count"`
    Patterns    []Pattern             `json:"patterns"`
    Bottlenecks []Bottleneck         `json:"bottlenecks"`
    Metrics     map[string]float64    `json:"metrics"`
    Anomalies   []Anomaly            `json:"anomalies"`
}

// Pattern 追踪模式
type Pattern struct {
    Type        string    `json:"type"`
    Confidence  float64   `json:"confidence"`
    Frequency   int       `json:"frequency"`
    Duration    time.Duration `json:"duration"`
    SpanIDs     []SpanID  `json:"span_ids"`
}

// Bottleneck 性能瓶颈
type Bottleneck struct {
    SpanID      SpanID    `json:"span_id"`
    Type        string    `json:"type"`
    Severity    float64   `json:"severity"`
    Duration    time.Duration `json:"duration"`
    Impact      float64   `json:"impact"`
}

// Anomaly 异常情况
type Anomaly struct {
    Type        string    `json:"type"`
    SpanID      SpanID    `json:"span_id"`
    Timestamp   time.Time `json:"timestamp"`
    Score       float64   `json:"score"`
    Description string    `json:"description"`
}

// Analyzer 追踪分析器
type Analyzer struct {
    mu sync.RWMutex

    // 配置
    config struct {
        WindowSize       int           // 分析窗口大小
        UpdateInterval   time.Duration // 更新间隔
        AnomalyThreshold float64      // 异常阈值
        PatternMinFreq   int          // 模式最小频率
    }

    // 数据源
    tracker  *Tracker
    recorder *Recorder

    // 分析缓存
    cache struct {
        traces    map[TraceID]*TraceAnalysis
        patterns  []Pattern
        anomalies []Anomaly
    }

    // 分析状态
    status struct {
        isRunning    bool
        lastAnalysis time.Time
        errors       []error
    }
}

// NewAnalyzer 创建新的分析器
func NewAnalyzer(tracker *Tracker, recorder *Recorder, config types.TraceConfig) *Analyzer {
    a := &Analyzer{
        tracker:  tracker,
        recorder: recorder,
    }

    // 初始化配置
    a.config.WindowSize = config.AnalysisWindowSize
    a.config.UpdateInterval = config.AnalysisInterval
    a.config.AnomalyThreshold = config.AnomalyThreshold
    a.config.PatternMinFreq = config.PatternMinFreq

    // 初始化缓存
    a.cache.traces = make(map[TraceID]*TraceAnalysis)

    return a
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

    // 启动分析循环
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
                a.recordError(err)
            }
        }
    }
}

// analyze 执行分析
func (a *Analyzer) analyze(ctx context.Context) error {
    // 获取追踪数据
    traces := a.getTracesInWindow()

    // 分析模式
    patterns := a.analyzePatterns(traces)

    // 检测瓶颈
    bottlenecks := a.detectBottlenecks(traces)

    // 检测异常
    anomalies := a.detectAnomalies(traces)

    // 生成分析结果
    analysis := &TraceAnalysis{
        ID:         generateID(),
        Timestamp:  time.Now(),
        Patterns:   patterns,
        Bottlenecks: bottlenecks,
        Anomalies:  anomalies,
    }

    // 缓存结果
    a.cacheAnalysis(analysis)

    return nil
}

// analyzePatterns 分析追踪模式
func (a *Analyzer) analyzePatterns(traces map[TraceID]*TraceAnalysis) []Pattern {
    patterns := make([]Pattern, 0)
    
    // 构建模式图
    graph := newPatternGraph()
    for _, trace := range traces {
        graph.addTrace(trace)
    }

    // 提取频繁模式
    frequentPatterns := graph.findFrequentPatterns(a.config.PatternMinFreq)
    
    // 评估模式可信度
    for _, p := range frequentPatterns {
        confidence := a.calculatePatternConfidence(p, traces)
        if confidence > 0.7 { // 可信度阈值
            patterns = append(patterns, Pattern{
                Type:       p.Type,
                Confidence: confidence,
                Frequency:  p.Frequency,
                Duration:   p.Duration,
                SpanIDs:   p.SpanIDs,
            })
        }
    }

    return patterns
}

// detectBottlenecks 检测性能瓶颈
func (a *Analyzer) detectBottlenecks(traces map[TraceID]*TraceAnalysis) []Bottleneck {
    bottlenecks := make([]Bottleneck, 0)

    // 计算span统计信息
    spanStats := make(map[string]struct {
        count    int
        duration time.Duration
        maxTime  time.Duration
    })

    // 收集统计数据
    for _, trace := range traces {
        for _, span := range trace.Spans {
            stats := spanStats[span.Name]
            stats.count++
            stats.duration += span.Duration
            if span.Duration > stats.maxTime {
                stats.maxTime = span.Duration
            }
            spanStats[span.Name] = stats
        }
    }

    // 识别瓶颈
    for name, stats := range spanStats {
        avgDuration := stats.duration / time.Duration(stats.count)
        if avgDuration > a.config.BottleneckThreshold {
            bottlenecks = append(bottlenecks, Bottleneck{
                Type:     "duration",
                Severity: float64(avgDuration) / float64(a.config.BottleneckThreshold),
                Duration: avgDuration,
                Impact:   float64(stats.count) * float64(avgDuration),
            })
        }
    }

    return bottlenecks
}

// detectAnomalies 检测异常
func (a *Analyzer) detectAnomalies(traces map[TraceID]*TraceAnalysis) []Anomaly {
    anomalies := make([]Anomaly, 0)

    // 计算基准统计
    baselineStats := a.calculateBaselineStats(traces)

    // 检测异常
    for traceID, trace := range traces {
        // 检查持续时间异常
        if score := a.calculateDurationAnomaly(trace, baselineStats); score > a.config.AnomalyThreshold {
            anomalies = append(anomalies, Anomaly{
                Type:      "duration",
                SpanID:    trace.SpanID,
                Timestamp: trace.Timestamp,
                Score:     score,
                Description: "Abnormal trace duration detected",
            })
        }

        // 检查模式异常
        if score := a.calculatePatternAnomaly(trace, baselineStats); score > a.config.AnomalyThreshold {
            anomalies = append(anomalies, Anomaly{
                Type:      "pattern",
                SpanID:    trace.SpanID,
                Timestamp: trace.Timestamp,
                Score:     score,
                Description: "Unusual trace pattern detected",
            })
        }
    }

    return anomalies
}

// calculateBaselineStats 计算基准统计信息
func (a *Analyzer) calculateBaselineStats(traces map[TraceID]*TraceAnalysis) map[string]interface{} {
    stats := make(map[string]interface{})

    // 计算平均持续时间
    var totalDuration time.Duration
    for _, trace := range traces {
        totalDuration += trace.Duration
    }
    stats["avgDuration"] = totalDuration / time.Duration(len(traces))

    // 计算其他统计信息
    // ...

    return stats
}

// cacheAnalysis 缓存分析结果
func (a *Analyzer) cacheAnalysis(analysis *TraceAnalysis) {
    a.mu.Lock()
    defer a.mu.Unlock()

    a.cache.traces[analysis.TraceID] = analysis
    a.status.lastAnalysis = time.Now()
}

// recordError 记录错误
func (a *Analyzer) recordError(err error) {
    a.mu.Lock()
    defer a.mu.Unlock()

    a.status.errors = append(a.status.errors, err)
}

// getTracesInWindow 获取分析窗口内的追踪数据
func (a *Analyzer) getTracesInWindow() map[TraceID]*TraceAnalysis {
    a.mu.RLock()
    defer a.mu.RUnlock()

    windowStart := time.Now().Add(-time.Duration(a.config.WindowSize) * time.Second)
    traces := make(map[TraceID]*TraceAnalysis)

    for id, trace := range a.cache.traces {
        if trace.Timestamp.After(windowStart) {
            traces[id] = trace
        }
    }

    return traces
}

// GetAnalysis 获取分析结果
func (a *Analyzer) GetAnalysis(traceID TraceID) (*TraceAnalysis, error) {
    a.mu.RLock()
    defer a.mu.RUnlock()

    analysis, exists := a.cache.traces[traceID]
    if !exists {
        return nil, types.NewSystemError(types.ErrNotFound, "trace analysis not found", nil)
    }

    return analysis, nil
}
