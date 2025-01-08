// system/adaptation.go

package system

import (
    "math"
    "sync"
    "time"
    "context"

    "github.com/Corphon/daoflow/model"
)

// AdaptationConstants 适应常数
const (
    MaxAdaptLevel    = 100.0  // 最大适应等级
    LearnRate       = 0.1    // 学习率
    MemoryDecay     = 0.05   // 记忆衰减率
    PatternThreshold = 0.6   // 模式识别阈值
)

// AdaptationSystem 适应系统
type AdaptationSystem struct {
    mu sync.RWMutex

    // 关联系统
    evolution *EvolutionSystem
    integrate *model.IntegrateFlow

    // 适应状态
    state struct {
        Level      float64                // 适应等级
        Patterns   map[string]*Pattern    // 已识别模式
        Memory     *AdaptiveMemory        // 适应性记忆
        Strategies map[string]*Strategy   // 适应策略
    }

    // 学习控制
    learning struct {
        rate     float64
        momentum float64
        history  []LearningRecord
    }

    monitor *AdaptationMonitor
    ctx     context.Context
    cancel  context.CancelFunc
}

// Pattern 模式结构
type Pattern struct {
    ID        string
    Type      PatternType
    Strength  float64
    Frequency float64
    LastSeen  time.Time
    Features  map[string]float64
}

// PatternType 模式类型
type PatternType uint8

const (
    PatternStable PatternType = iota  // 稳定模式
    PatternCyclic                     // 循环模式
    PatternEmergent                   // 涌现模式
    PatternChaotic                    // 混沌模式
)

// AdaptiveMemory 适应性记忆
type AdaptiveMemory struct {
    ShortTerm  []MemoryUnit  // 短期记忆
    LongTerm   []MemoryUnit  // 长期记忆
    Capacity   int           // 记忆容量
    Threshold  float64       // 转换阈值
}

// MemoryUnit 记忆单元
type MemoryUnit struct {
    Pattern   *Pattern
    Weight    float64
    Age       time.Duration
    Context   map[string]float64
}

// Strategy 适应策略
type Strategy struct {
    ID         string
    Type       StrategyType
    Conditions []Condition
    Actions    []Action
    Success    float64    // 成功率
    Cost       float64    // 能量消耗
}

// StrategyType 策略类型
type StrategyType uint8

const (
    StrategyReactive StrategyType = iota  // 反应式策略
    StrategyProactive                     // 主动式策略
    StrategyHybrid                        // 混合式策略
)

// NewAdaptationSystem 创建适应系统
func NewAdaptationSystem(ctx context.Context, es *EvolutionSystem, 
    integrate *model.IntegrateFlow) *AdaptationSystem {
    
    ctx, cancel := context.WithCancel(ctx)
    
    as := &AdaptationSystem{
        evolution: es,
        integrate: integrate,
        ctx:       ctx,
        cancel:    cancel,
    }

    // 初始化状态
    as.state.Level = 1.0
    as.state.Patterns = make(map[string]*Pattern)
    as.state.Memory = newAdaptiveMemory()
    as.state.Strategies = make(map[string]*Strategy)

    // 初始化学习参数
    as.learning.rate = LearnRate
    as.learning.momentum = 0.9

    as.monitor = NewAdaptationMonitor(as)
    
    go as.runAdaptation()
    return as
}

// runAdaptation 运行适应过程
func (as *AdaptationSystem) runAdaptation() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-as.ctx.Done():
            return
        case <-ticker.C:
            as.adapt()
        }
    }
}

// adapt 执行适应
func (as *AdaptationSystem) adapt() {
    as.mu.Lock()
    defer as.mu.Unlock()

    // 获取当前系统状态
    systemState := as.integrate.GetSystemState()
    evolutionStatus := as.evolution.GetEvolutionStatus()

    // 模式识别
    patterns := as.recognizePatterns(systemState)

    // 更新记忆
    as.updateMemory(patterns)

    // 策略选择与执行
    strategy := as.selectStrategy(systemState, evolutionStatus)
    if strategy != nil {
        as.executeStrategy(strategy)
    }

    // 学习与优化
    as.learn(patterns, strategy)
}

// recognizePatterns 识别模式
func (as *AdaptationSystem) recognizePatterns(state model.SystemState) []*Pattern {
    patterns := make([]*Pattern, 0)

    // 使用动态系统理论分析状态
    dynamics := as.analyzeDynamics(state)

    // 识别稳定模式
    if stablePattern := as.detectStablePattern(dynamics); stablePattern != nil {
        patterns = append(patterns, stablePattern)
    }

    // 识别循环模式
    if cyclicPattern := as.detectCyclicPattern(dynamics); cyclicPattern != nil {
        patterns = append(patterns, cyclicPattern)
    }

    // 识别涌现模式
    if emergentPattern := as.detectEmergentPattern(dynamics); emergentPattern != nil {
        patterns = append(patterns, emergentPattern)
    }

    return patterns
}

// analyzeDynamics 分析系统动态特性
func (as *AdaptationSystem) analyzeDynamics(state model.SystemState) map[string]float64 {
    dynamics := make(map[string]float64)

    // 使用混沌理论计算李雅普诺夫指数
    lyapunov := as.calculateLyapunovExponent(state)
    dynamics["chaos_degree"] = lyapunov

    // 使用信息熵评估系统复杂度
    entropy := as.calculateInformationEntropy(state)
    dynamics["complexity"] = entropy

    // 计算系统相位空间轨迹
    trajectory := as.calculatePhaseTrajectory(state)
    dynamics["trajectory_stability"] = trajectory

    return dynamics
}

// calculateLyapunovExponent 计算李雅普诺夫指数
func (as *AdaptationSystem) calculateLyapunovExponent(state model.SystemState) float64 {
    // 获取历史状态序列
    history := as.learning.history
    if len(history) < 2 {
        return 0
    }

    // 计算状态空间中的发散率
    divergence := 0.0
    for i := 1; i < len(history); i++ {
        // 计算相邻状态间的距离
        dist := as.calculateStateDistance(history[i].State, history[i-1].State)
        if dist > 0 {
            divergence += math.Log(dist)
        }
    }

    return divergence / float64(len(history))
}

// updateMemory 更新适应性记忆
func (as *AdaptationSystem) updateMemory(patterns []*Pattern) {
    memory := as.state.Memory

    // 更新短期记忆
    for _, pattern := range patterns {
        unit := MemoryUnit{
            Pattern: pattern,
            Weight:  1.0,
            Age:     0,
            Context: as.getCurrentContext(),
        }
        memory.ShortTerm = append(memory.ShortTerm, unit)
    }

    // 应用记忆衰减
    as.applyMemoryDecay()

    // 转移到长期记忆
    as.consolidateMemory()
}

// applyMemoryDecay 应用记忆衰减
func (as *AdaptationSystem) applyMemoryDecay() {
    for i := range as.state.Memory.ShortTerm {
        as.state.Memory.ShortTerm[i].Weight *= (1 - MemoryDecay)
    }

    for i := range as.state.Memory.LongTerm {
        as.state.Memory.LongTerm[i].Weight *= (1 - MemoryDecay/2)
    }
}

// learn 学习过程
func (as *AdaptationSystem) learn(patterns []*Pattern, strategy *Strategy) {
    // 计算学习目标
    target := as.calculateLearningTarget()

    // 更新策略权重
    if strategy != nil {
        delta := (target - strategy.Success) * as.learning.rate
        strategy.Success += delta + as.learning.momentum*delta
    }

    // 更新模式识别阈值
    as.updatePatternThresholds(patterns)

    // 记录学习历史
    as.learning.history = append(as.learning.history, LearningRecord{
        Patterns: patterns,
        Strategy: strategy,
        Result:   target,
        Time:    time.Now(),
    })
}

// GetAdaptationStatus 获取适应状态
func (as *AdaptationSystem) GetAdaptationStatus() map[string]interface{} {
    as.mu.RLock()
    defer as.mu.RUnlock()

    return map[string]interface{}{
        "level":     as.state.Level,
        "patterns":  len(as.state.Patterns),
        "memory":    len(as.state.Memory.LongTerm),
        "strategies": len(as.state.Strategies),
    }
}

// Close 关闭适应系统
func (as *AdaptationSystem) Close() error {
    as.cancel()
    return nil
}
