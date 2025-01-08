// system/evolution.go

package system

import (
    "math"
    "sync"
    "time"
    "context"

    "github.com/Corphon/daoflow/model"
    "github.com/Corphon/daoflow/core"
)

// EvolutionConstants 演化常数
const (
    MaxEvolutionLevel = 100.0   // 最大演化等级
    BaseAdaptRate    = 0.15    // 基础适应率
    MutationRate     = 0.05    // 变异率
    EnergyThreshold  = 0.7     // 能量阈值
    MinStability     = 0.3     // 最小稳定度
)

// EvolutionSystem 演化系统
type EvolutionSystem struct {
    mu sync.RWMutex

    // 集成模型
    integrate *model.IntegrateFlow
    
    // 演化状态
    state struct {
        Level      float64            // 演化等级
        Stability  float64            // 稳定度
        Diversity  float64            // 多样性
        Fitness    float64            // 适应度
        Emergence  map[string]float64 // 涌现属性
    }
    
    // 适应策略
    strategy struct {
        adaptationRate float64
        mutationProb  float64
        energyPolicy  EnergyPolicy
    }
    
    // 演化链追踪
    evolution struct {
        chain    []EvolutionNode
        branches map[string][]EvolutionBranch
    }

    // 监控
    monitor *EvolutionMonitor
    
    // 上下文控制
    ctx    context.Context
    cancel context.CancelFunc
}

// EnergyPolicy 能量策略
type EnergyPolicy struct {
    Distribution []float64  // 能量分配比例
    Threshold    float64    // 触发阈值
    Priority     []string   // 优先级序列
}

// EvolutionNode 演化节点
type EvolutionNode struct {
    ID        string
    Level     float64
    State     model.SystemState
    Timestamp time.Time
    Changes   []StateChange
}

// EvolutionBranch 演化分支
type EvolutionBranch struct {
    ParentID  string
    ChildID   string
    Type      BranchType
    Cause     string
    Strength  float64
}

// StateChange 状态改变
type StateChange struct {
    Field     string
    OldValue  float64
    NewValue  float64
    Cause     string
}

// BranchType 分支类型
type BranchType uint8

const (
    BranchEvolution BranchType = iota // 常规演化
    BranchMutation                    // 变异
    BranchEmergence                   // 涌现
    BranchFusion                      // 融合
)

// NewEvolutionSystem 创建演化系统
func NewEvolutionSystem(ctx context.Context, integrate *model.IntegrateFlow) *EvolutionSystem {
    ctx, cancel := context.WithCancel(ctx)
    
    es := &EvolutionSystem{
        integrate: integrate,
        ctx:      ctx,
        cancel:   cancel,
    }
    
    es.state.Level = 1.0
    es.state.Emergence = make(map[string]float64)
    es.evolution.branches = make(map[string][]EvolutionBranch)
    
    es.initializeStrategy()
    es.monitor = NewEvolutionMonitor(es)
    
    go es.runEvolution()
    return es
}

// initializeStrategy 初始化策略
func (es *EvolutionSystem) initializeStrategy() {
    es.strategy.adaptationRate = BaseAdaptRate
    es.strategy.mutationProb = MutationRate
    es.strategy.energyPolicy = EnergyPolicy{
        Distribution: []float64{0.3, 0.3, 0.2, 0.2}, // 阴阳、五行、八卦、干支
        Threshold:    EnergyThreshold,
        Priority:     []string{"stability", "diversity", "fitness"},
    }
}

// runEvolution 运行演化过程
func (es *EvolutionSystem) runEvolution() {
    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-es.ctx.Done():
            return
        case <-ticker.C:
            es.evolve()
        }
    }
}

// evolve 执行演化
func (es *EvolutionSystem) evolve() {
    es.mu.Lock()
    defer es.mu.Unlock()

    // 获取当前系统状态
    currentState := es.integrate.GetSystemState()
    
    // 计算系统指标
    stability := es.calculateStability(currentState)
    diversity := es.calculateDiversity(currentState)
    fitness := es.calculateFitness(stability, diversity)
    
    // 更新状态
    es.state.Stability = stability
    es.state.Diversity = diversity
    es.state.Fitness = fitness
    
    // 检查是否需要演化
    if es.shouldEvolve(fitness) {
        es.performEvolution(currentState)
    }
    
    // 检查涌现现象
    es.checkEmergence(currentState)
}

// calculateStability 计算稳定度
func (es *EvolutionSystem) calculateStability(state model.SystemState) float64 {
    // 使用统计物理学的熵理论
    entropy := state.System.Entropy
    energy := state.System.Energy
    
    // 稳定度 = 1 / (1 + 归一化熵)
    normalizedEntropy := entropy / math.Log(energy+1)
    return 1.0 / (1.0 + normalizedEntropy)
}

// calculateDiversity 计算多样性
func (es *EvolutionSystem) calculateDiversity(state model.SystemState) float64 {
    // 使用信息论的香农多样性指数
    var diversity float64
    totalEnergy := state.System.Energy
    
    // 计算每个子系统的能量占比
    ratios := []float64{
        state.YinYang / totalEnergy,
        state.System.WuXingEnergy / totalEnergy,
        state.System.BaGuaEnergy / totalEnergy,
        state.System.GanZhiEnergy / totalEnergy,
    }
    
    // H = -Σ(pi * log(pi))
    for _, p := range ratios {
        if p > 0 {
            diversity -= p * math.Log2(p)
        }
    }
    
    return diversity / math.Log2(4) // 归一化
}

// calculateFitness 计算适应度
func (es *EvolutionSystem) calculateFitness(stability, diversity float64) float64 {
    // 使用多目标优化理论
    weights := []float64{0.6, 0.4} // 稳定性权重大于多样性
    
    return weights[0]*stability + weights[1]*diversity
}

// shouldEvolve 判断是否应该演化
func (es *EvolutionSystem) shouldEvolve(fitness float64) bool {
    // 基于适应度和随机扰动判断
    threshold := es.strategy.energyPolicy.Threshold
    randomFactor := 1.0 + (rand.Float64()-0.5)*0.2 // ±10%波动
    
    return fitness*randomFactor < threshold
}

// performEvolution 执行演化
func (es *EvolutionSystem) performEvolution(state model.SystemState) {
    // 创建演化节点
    node := EvolutionNode{
        ID:        generateID(),
        Level:     es.state.Level,
        State:     state,
        Timestamp: time.Now(),
    }
    
    // 尝试多种演化路径
    paths := es.calculateEvolutionPaths(state)
    
    // 选择最优路径
    bestPath := es.selectBestPath(paths)
    
    // 应用演化
    es.applyEvolution(bestPath)
    
    // 记录变化
    node.Changes = es.recordChanges(state, bestPath)
    es.evolution.chain = append(es.evolution.chain, node)
}

// checkEmergence 检查涌现现象
func (es *EvolutionSystem) checkEmergence(state model.SystemState) {
    // 使用复杂系统理论检测涌现属性
    patterns := es.detectPatterns(state)
    
    for pattern, strength := range patterns {
        if strength > es.strategy.energyPolicy.Threshold {
            es.state.Emergence[pattern] = strength
            
            // 创建涌现分支
            branch := EvolutionBranch{
                ParentID: es.evolution.chain[len(es.evolution.chain)-1].ID,
                Type:     BranchEmergence,
                Cause:    pattern,
                Strength: strength,
            }
            
            es.evolution.branches[pattern] = append(
                es.evolution.branches[pattern], branch)
        }
    }
}

// detectPatterns 检测模式
func (es *EvolutionSystem) detectPatterns(state model.SystemState) map[string]float64 {
    patterns := make(map[string]float64)
    
    // 检测阴阳平衡模式
    yinYangBalance := es.detectYinYangBalance(state)
    if yinYangBalance > MinStability {
        patterns["yin_yang_balance"] = yinYangBalance
    }
    
    // 检测五行循环模式
    wuxingCycle := es.detectWuXingCycle(state)
    if wuxingCycle > MinStability {
        patterns["wuxing_cycle"] = wuxingCycle
    }
    
    // 检测八卦共振
    baguaResonance := es.detectBaGuaResonance(state)
    if baguaResonance > MinStability {
        patterns["bagua_resonance"] = baguaResonance
    }
    
    // 检测干支合化
    ganzhiHarmony := es.detectGanZhiHarmony(state)
    if ganzhiHarmony > MinStability {
        patterns["ganzhi_harmony"] = ganzhiHarmony
    }
    
    return patterns
}

// GetEvolutionStatus 获取演化状态
func (es *EvolutionSystem) GetEvolutionStatus() map[string]interface{} {
    es.mu.RLock()
    defer es.mu.RUnlock()
    
    return map[string]interface{}{
        "level":     es.state.Level,
        "stability": es.state.Stability,
        "diversity": es.state.Diversity,
        "fitness":   es.state.Fitness,
        "emergence": es.state.Emergence,
    }
}

// Close 关闭演化系统
func (es *EvolutionSystem) Close() error {
    es.cancel()
    return nil
}
