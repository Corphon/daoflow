//system/evolution/mutation/detector.go

package mutation

import (
    "sync"
    "time"

    "github.com/Corphon/daoflow/system/common" // 引入共享接口包
)

// MutationDetector 突变检测器
type MutationDetector struct {
    mu sync.RWMutex

    // 基础配置
    config struct {
        detectionThreshold float64      // 检测阈值
        timeWindow        time.Duration // 时间窗口
        sensitivity      float64       // 灵敏度
        stabilityFactor  float64       // 稳定性因子
    }

    // 使用接口而不是具体类型
    patternAnalyzer common.PatternAnalyzer
    
    // 检测状态
    state struct {
        mutations    map[string]*Mutation       // 已检测突变
        observations []MutationObservation      // 观察记录
        baselines   map[string]*MutationBaseline // 基准线
    }
}

// Mutation 突变
type Mutation struct {
    ID          string                // 突变ID
    Type        string                // 突变类型
    Source      *MutationSource       // 突变源
    Changes     []MutationChange      // 变化列表
    Severity    float64               // 严重程度
    Probability float64               // 发生概率
    DetectedAt  time.Time            // 检测时间
    Status      string                // 当前状态
}

// MutationSource 突变源
type MutationSource struct {
    PatternID   string                // 相关模式ID
    Location    string                // 发生位置
    Context     map[string]interface{} // 上下文信息
    Energy      float64               // 能量水平
}

// MutationChange 突变变化
type MutationChange struct {
    Property    string                // 变化属性
    OldValue    interface{}           // 原值
    NewValue    interface{}           // 新值
    Delta       float64               // 变化量
    Timestamp   time.Time             // 变化时间
}

// MutationObservation 突变观察
type MutationObservation struct {
    Timestamp   time.Time
    PatternID   string
    Metrics     map[string]float64
    Anomalies   []string
}

// MutationBaseline 突变基准线
type MutationBaseline struct {
    PatternID   string
    Metrics     map[string]BaselineMetric
    LastUpdate  time.Time
    Confidence  float64
}

// BaselineMetric 基准度量
type BaselineMetric struct {
    Mean        float64
    StdDev      float64
    Bounds      [2]float64  // [min, max]
    History     []float64
}

// NewMutationDetector 创建新的突变检测器
func NewMutationDetector(analyzer common.PatternAnalyzer) *MutationDetector {
    
    md := &MutationDetector{
        patternAnalyzer: analyzer,
    }

    // 初始化配置
    md.config.detectionThreshold = 0.75
    md.config.timeWindow = 10 * time.Minute
    md.config.sensitivity = 0.8
    md.config.stabilityFactor = 0.6

    // 初始化状态
    md.state.mutations = make(map[string]*Mutation)
    md.state.observations = make([]MutationObservation, 0)
    md.state.baselines = make(map[string]*MutationBaseline)

    return md
}

// Detect 执行突变检测
func (md *MutationDetector) Detect() error {
    md.mu.Lock()
    defer md.mu.Unlock()

    // 获取当前模式
    patterns, err := md.recognizer.GetPatterns()
    if err != nil {
        return err
    }

    // 更新观察记录
    md.updateObservations(patterns)

    // 更新基准线
    md.updateBaselines()

    // 检测突变
    mutations := md.detectMutations(patterns)

    // 验证突变
    validated := md.validateMutations(mutations)

    // 更新突变状态
    md.updateMutations(validated)

    return nil
}

// updateObservations 更新观察记录
func (md *MutationDetector) updateObservations(patterns []*pattern.RecognizedPattern) {
    currentTime := time.Now()

    // 创建新的观察记录
    for _, pat := range patterns {
        observation := MutationObservation{
            Timestamp: currentTime,
            PatternID: pat.ID,
            Metrics:   md.collectMetrics(pat),
            Anomalies: make([]string, 0),
        }

        // 检查异常
        anomalies := md.checkAnomalies(pat, observation.Metrics)
        observation.Anomalies = anomalies

        md.state.observations = append(md.state.observations, observation)
    }

    // 清理过期观察
    md.cleanupObservations()
}

// updateBaselines 更新基准线
func (md *MutationDetector) updateBaselines() {
    currentTime := time.Now()

    // 获取观察数据
    observations := md.getRecentObservations(md.config.timeWindow)

    // 更新每个模式的基准线
    for patternID := range md.collectPatternIDs(observations) {
        // 获取模式的观察数据
        patternObs := md.filterObservationsByPattern(observations, patternID)
        
        // 计算新的基准线
        baseline := md.calculateBaseline(patternObs)
        
        if baseline != nil {
            baseline.LastUpdate = currentTime
            md.state.baselines[patternID] = baseline
        }
    }
}

// detectMutations 检测突变
func (md *MutationDetector) detectMutations(
    patterns []*pattern.RecognizedPattern) []*Mutation {
    
    mutations := make([]*Mutation, 0)

    for _, pat := range patterns {
        // 获取基准线
        baseline := md.state.baselines[pat.ID]
        if baseline == nil {
            continue
        }

        // 检查突变条件
        if changes := md.checkMutationConditions(pat, baseline); len(changes) > 0 {
            // 创建突变记录
            mutation := &Mutation{
                ID:         generateMutationID(),
                Type:      md.determineMutationType(changes),
                Source:    createMutationSource(pat),
                Changes:   changes,
                DetectedAt: time.Now(),
                Status:    "detected",
            }

            // 计算严重程度和概率
            mutation.Severity = md.calculateMutationSeverity(mutation)
            mutation.Probability = md.calculateMutationProbability(mutation)

            mutations = append(mutations, mutation)
        }
    }

    return mutations
}

// validateMutations 验证突变
func (md *MutationDetector) validateMutations(
    mutations []*Mutation) []*Mutation {
    
    validated := make([]*Mutation, 0)

    for _, mutation := range mutations {
        // 检查突变有效性
        if md.isMutationValid(mutation) {
            validated = append(validated, mutation)
        }
    }

    return validated
}

// 辅助函数

func (md *MutationDetector) collectMetrics(
    pattern *pattern.RecognizedPattern) map[string]float64 {
    
    metrics := make(map[string]float64)
    
    // 收集关键指标
    metrics["energy"] = pattern.Energy
    metrics["stability"] = pattern.Stability
    metrics["complexity"] = calculatePatternComplexity(pattern)
    metrics["coherence"] = calculatePatternCoherence(pattern)
    
    return metrics
}

func (md *MutationDetector) checkAnomalies(
    pattern *pattern.RecognizedPattern,
    metrics map[string]float64) []string {
    
    anomalies := make([]string, 0)
    
    baseline := md.state.baselines[pattern.ID]
    if baseline == nil {
        return anomalies
    }

    // 检查每个指标
    for name, value := range metrics {
        if baseMetric, ok := baseline.Metrics[name]; ok {
            if md.isAnomaly(value, baseMetric) {
                anomalies = append(anomalies, name)
            }
        }
    }

    return anomalies
}

func generateMutationID() string {
    return fmt.Sprintf("mut_%d", time.Now().UnixNano())
}
