// system/metrics.go

package system

import (
    "math"
    "sync"
    "time"

    "github.com/Corphon/daoflow/model"
)

// MetricsConstants 指标常数
const (
    MaxMetricHistory = 1000   // 最大历史记录数
    SampleRate      = 10     // 采样率(Hz)
    DecayFactor     = 0.95   // 衰减因子
    Smoothing       = 0.1    // 平滑系数
)

// MetricsSystem 指标系统
type MetricsSystem struct {
    mu sync.RWMutex

    // 系统指标
    core struct {
        Evolution      EvolutionMetrics      // 演化指标
        Adaptation     AdaptationMetrics     // 适应指标
        Synchronization SynchronizationMetrics // 同步指标
        Optimization   OptimizationMetrics   // 优化指标
        Emergence      EmergenceMetrics      // 涌现指标
    }

    // 性能指标
    performance struct {
        CPU       float64    // CPU使用率
        Memory    float64    // 内存使用率
        Energy    float64    // 能量效率
        Latency   float64    // 系统延迟
    }

    // 量化指标
    quantum struct {
        Coherence    float64    // 相干度
        Entanglement float64    // 纠缠度
        Phase        float64    // 相位
        Wave        []float64   // 波函数
    }

    // 统计数据
    stats struct {
        History    []MetricPoint  // 历史数据
        Trends     []TrendLine    // 趋势线
        Moments    Moments        // 统计矩
        Spectrum   []float64      // 频谱
    }
}

// EvolutionMetrics 演化指标
type EvolutionMetrics struct {
    Level      float64    // 演化等级
    Speed      float64    // 演化速度
    Direction  Vector3D   // 演化方向
    Entropy    float64    // 演化熵
}

// AdaptationMetrics 适应指标
type AdaptationMetrics struct {
    Fitness    float64    // 适应度
    Learning   float64    // 学习率
    Memory     float64    // 记忆强度
    Response   float64    // 响应速度
}

// SynchronizationMetrics 同步指标
type SynchronizationMetrics struct {
    Coherence  float64    // 相干度
    Phase      float64    // 相位差
    Coupling   float64    // 耦合强度
    Stability  float64    // 稳定性
}

// OptimizationMetrics 优化指标
type OptimizationMetrics struct {
    Objective  float64    // 目标值
    Progress   float64    // 进度
    Efficiency float64    // 效率
    Quality    float64    // 质量
}

// EmergenceMetrics 涌现指标
type EmergenceMetrics struct {
    Complexity float64    // 复杂度
    Novelty    float64    // 新颖度
    Integration float64   // 集成度
    Potential  float64    // 潜力
}

// MetricPoint 度量点
type MetricPoint struct {
    Name      string
    Value     float64
    Tags      map[string]string
    Time      time.Time
}

// TrendLine 趋势线
type TrendLine struct {
    Slope     float64
    Intercept float64
    R2        float64
    Points    []Point2D
}

// Moments 统计矩
type Moments struct {
    Mean      float64
    Variance  float64
    Skewness  float64
    Kurtosis  float64
}

// NewMetricsSystem 创建指标系统
func NewMetricsSystem() *MetricsSystem {
    ms := &MetricsSystem{}
    ms.initializeMetrics()
    return ms
}

// initializeMetrics 初始化指标
func (ms *MetricsSystem) initializeMetrics() {
    // 初始化波函数数组
    ms.quantum.Wave = make([]float64, 100)
    
    // 初始化历史数据
    ms.stats.History = make([]MetricPoint, 0, MaxMetricHistory)
    
    // 初始化频谱数据
    ms.stats.Spectrum = make([]float64, 1024)
}

// CollectMetrics 收集指标
func (ms *MetricsSystem) CollectMetrics(state *SystemState) {
    ms.mu.Lock()
    defer ms.mu.Unlock()

    // 收集系统核心指标
    ms.collectCoreMetrics(state)

    // 收集性能指标
    ms.collectPerformanceMetrics()

    // 收集量子指标
    ms.collectQuantumMetrics()

    // 更新统计数据
    ms.updateStatistics()
}

// collectCoreMetrics 收集核心指标
func (ms *MetricsSystem) collectCoreMetrics(state *SystemState) {
    // 收集演化指标
    ms.core.Evolution = EvolutionMetrics{
        Level:     state.Evolution.Level,
        Speed:     ms.calculateEvolutionSpeed(state),
        Direction: ms.calculateEvolutionDirection(state),
        Entropy:   ms.calculateEvolutionEntropy(state),
    }

    // 收集适应指标
    ms.core.Adaptation = AdaptationMetrics{
        Fitness:   state.Adaptation.Fitness,
        Learning:  state.Adaptation.LearningRate,
        Memory:    ms.calculateMemoryStrength(state),
        Response:  ms.calculateResponseTime(state),
    }

    // 其他核心指标收集...
}

// calculateEvolutionEntropy 计算演化熵
func (ms *MetricsSystem) calculateEvolutionEntropy(state *SystemState) float64 {
    var entropy float64
    totalEnergy := state.Energy

    // 使用香农熵公式
    for _, e := range state.Evolution.EnergyDistribution {
        if e > 0 {
            p := e / totalEnergy
            entropy -= p * math.Log2(p)
        }
    }

    return entropy
}

// updateStatistics 更新统计数据
func (ms *MetricsSystem) updateStatistics() {
    // 更新历史数据
    if len(ms.stats.History) >= MaxMetricHistory {
        ms.stats.History = ms.stats.History[1:]
    }
    
    // 计算趋势
    ms.calculateTrends()
    
    // 计算统计矩
    ms.calculateMoments()
    
    // 计算频谱
    ms.calculateSpectrum()
}

// calculateTrends 计算趋势
func (ms *MetricsSystem) calculateTrends() {
    // 使用最小二乘法计算趋势线
    for name := range ms.groupMetricsByName() {
        points := ms.getMetricPoints(name)
        if len(points) < 2 {
            continue
        }

        trend := ms.calculateLinearRegression(points)
        ms.stats.Trends = append(ms.stats.Trends, trend)
    }
}

// calculateMoments 计算统计矩
func (ms *MetricsSystem) calculateMoments() {
    values := make([]float64, len(ms.stats.History))
    for i, point := range ms.stats.History {
        values[i] = point.Value
    }

    ms.stats.Moments = Moments{
        Mean:     ms.calculateMean(values),
        Variance: ms.calculateVariance(values),
        Skewness: ms.calculateSkewness(values),
        Kurtosis: ms.calculateKurtosis(values),
    }
}

// GetMetricsSnapshot 获取指标快照
func (ms *MetricsSystem) GetMetricsSnapshot() map[string]interface{} {
    ms.mu.RLock()
    defer ms.mu.RUnlock()

    return map[string]interface{}{
        "core": map[string]interface{}{
            "evolution":      ms.core.Evolution,
            "adaptation":     ms.core.Adaptation,
            "synchronization": ms.core.Synchronization,
            "optimization":   ms.core.Optimization,
            "emergence":      ms.core.Emergence,
        },
        "performance": map[string]interface{}{
            "cpu":     ms.performance.CPU,
            "memory":  ms.performance.Memory,
            "energy":  ms.performance.Energy,
            "latency": ms.performance.Latency,
        },
        "quantum": map[string]interface{}{
            "coherence":    ms.quantum.Coherence,
            "entanglement": ms.quantum.Entanglement,
            "phase":        ms.quantum.Phase,
        },
        "stats": map[string]interface{}{
            "moments":  ms.stats.Moments,
            "trends":   ms.stats.Trends,
            "spectrum": ms.stats.Spectrum,
        },
    }
}
