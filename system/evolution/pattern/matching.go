//system/evolution/pattern/matching.go

package pattern

import (
    "sync"
    "time"

    "github.com/Corphon/daoflow/meta/field"
    "github.com/Corphon/daoflow/meta/emergence"
    "github.com/Corphon/daoflow/meta/resonance"
    "github.com/Corphon/daoflow/model"
    "github.com/Corphon/daoflow/system/types"
)

// EvolutionMatcher 演化匹配器
type EvolutionMatcher struct {
    mu sync.RWMutex

    // 基础配置
    config struct {
        matchThreshold  float64         // 匹配阈值
        evolutionDepth  int             // 演化深度
        adaptiveBias    float64         // 自适应偏差
        contextWeight   float64         // 上下文权重
    }

    // 匹配状态
    state struct {
        matches        map[string]*EvolutionMatch   // 当前匹配
        trajectories   map[string]*EvolutionPath    // 演化轨迹
        context       *MatchingContext              // 匹配上下文
    }

    // 依赖项
    recognizer *PatternRecognizer
    matcher    *resonance.PatternMatcher
}

// EvolutionMatch 演化匹配
type EvolutionMatch struct {
    ID           string                // 匹配ID
    SourceID     string                // 源模式ID
    TargetID     string                // 目标模式ID
    Similarity   float64               // 相似度
    Evolution    []EvolutionStep       // 演化步骤
    Context      map[string]float64    // 上下文因素
    StartTime    time.Time             // 开始时间
    LastUpdate   time.Time             // 最后更新时间
}

// EvolutionPath 演化轨迹
type EvolutionPath struct {
    ID           string                // 轨迹ID
    Steps        []PathStep            // 轨迹步骤
    Properties   map[string]float64    // 轨迹属性
    Probability  float64               // 概率
    Created      time.Time             // 创建时间
}

// EvolutionStep 演化步骤
type EvolutionStep struct {
    Timestamp    time.Time
    Type         string                // 步骤类型
    Changes      map[string]float64    // 变化量
    Energy       float64               // 能量变化
    Stability    float64               // 稳定性
}

// PathStep 轨迹步骤
type PathStep struct {
    Pattern      *RecognizedPattern    // 相关模式
    Transition   string                // 转换类型
    Probability  float64               // 转换概率
    Context      map[string]float64    // 步骤上下文
}

// MatchingContext 匹配上下文
type MatchingContext struct {
    Time         time.Time             // 当前时间
    Environment  map[string]float64    // 环境因素
    History      []ContextState        // 历史状态
    Bias         map[string]float64    // 偏差项
}

// ContextState 上下文状态
type ContextState struct {
    Timestamp    time.Time
    Factors      map[string]float64
    Influence    float64
}

// NewEvolutionMatcher 创建新的演化匹配器
func NewEvolutionMatcher(
    recognizer *PatternRecognizer,
    matcher *resonance.PatternMatcher) *EvolutionMatcher {
    
    em := &EvolutionMatcher{
        recognizer: recognizer,
        matcher:    matcher,
    }

    // 初始化配置
    em.config.matchThreshold = 0.7
    em.config.evolutionDepth = 5
    em.config.adaptiveBias = 0.3
    em.config.contextWeight = 0.4

    // 初始化状态
    em.state.matches = make(map[string]*EvolutionMatch)
    em.state.trajectories = make(map[string]*EvolutionPath)
    em.state.context = &MatchingContext{
        Time:        time.Now(),
        Environment: make(map[string]float64),
        History:     make([]ContextState, 0),
        Bias:        make(map[string]float64),
    }

    return em
}

// Match 执行演化匹配
func (em *EvolutionMatcher) Match() error {
    em.mu.Lock()
    defer em.mu.Unlock()

    // 更新上下文
    em.updateContext()

    // 获取当前模式
    patterns, err := em.recognizer.GetPatterns()
    if err != nil {
        return err
    }

    // 执行匹配
    matches := em.matchPatterns(patterns)

    // 更新演化轨迹
    em.updateTrajectories(matches)

    // 预测演化方向
    em.predictEvolution()

    return nil
}

// matchPatterns 匹配模式
func (em *EvolutionMatcher) matchPatterns(
    patterns []*RecognizedPattern) []*EvolutionMatch {
    
    matches := make([]*EvolutionMatch, 0)

    // 对每对模式进行匹配
    for i, source := range patterns {
        for j := i + 1; j < len(patterns); j++ {
            target := patterns[j]
            
            // 计算演化相似度
            similarity := em.calculateEvolutionSimilarity(source, target)
            if similarity < em.config.matchThreshold {
                continue
            }

            // 创建匹配
            match := em.createMatch(source, target, similarity)
            matches = append(matches, match)
        }
    }

    return matches
}

// updateTrajectories 更新演化轨迹
func (em *EvolutionMatcher) updateTrajectories(matches []*EvolutionMatch) {
    for _, match := range matches {
        // 检查是否存在相关轨迹
        trajectoryID := em.findRelatedTrajectory(match)
        if trajectoryID == "" {
            // 创建新轨迹
            trajectory := em.createTrajectory(match)
            em.state.trajectories[trajectory.ID] = trajectory
        } else {
            // 更新现有轨迹
            em.updateExistingTrajectory(trajectoryID, match)
        }
    }

    // 移除过期轨迹
    em.cleanupTrajectories()
}

// predictEvolution 预测演化方向
func (em *EvolutionMatcher) predictEvolution() {
    for _, trajectory := range em.state.trajectories {
        // 分析轨迹模式
        pattern := em.analyzeTrajectoryPattern(trajectory)
        
        // 预测下一步
        nextStep := em.predictNextStep(trajectory, pattern)
        if nextStep != nil {
            trajectory.Steps = append(trajectory.Steps, *nextStep)
        }

        // 更新概率
        trajectory.Probability = em.calculateTrajectoryProbability(trajectory)
    }
}

// calculateEvolutionSimilarity 计算演化相似度
func (em *EvolutionMatcher) calculateEvolutionSimilarity(
    source, target *RecognizedPattern) float64 {
    
    // 基础相似度
    baseSimilarity := calculatePatternSimilarity(source, target)

    // 演化特征相似度
    evolutionSimilarity := em.compareEvolutionFeatures(source, target)

    // 上下文相似度
    contextSimilarity := em.calculateContextSimilarity(source, target)

    // 组合相似度
    similarity := (baseSimilarity +
                  evolutionSimilarity * (1 - em.config.contextWeight) +
                  contextSimilarity * em.config.contextWeight) / 3.0

    return similarity
}

// updateContext 更新上下文
func (em *EvolutionMatcher) updateContext() {
    currentTime := time.Now()

    // 更新环境因素
    em.updateEnvironmentFactors()

    // 记录当前状态
    state := ContextState{
        Timestamp: currentTime,
        Factors:   make(map[string]float64),
        Influence: calculateContextInfluence(em.state.context.Environment),
    }
    
    // 复制当前环境因素
    for k, v := range em.state.context.Environment {
        state.Factors[k] = v
    }

    // 更新历史
    em.state.context.History = append(em.state.context.History, state)
    if len(em.state.context.History) > maxContextHistory {
        em.state.context.History = em.state.context.History[1:]
    }

    // 更新时间
    em.state.context.Time = currentTime
}

// 辅助函数

func (em *EvolutionMatcher) createMatch(
    source, target *RecognizedPattern,
    similarity float64) *EvolutionMatch {
    
    match := &EvolutionMatch{
        ID:         generateMatchID(),
        SourceID:   source.ID,
        TargetID:   target.ID,
        Similarity: similarity,
        Evolution:  make([]EvolutionStep, 0),
        Context:    make(map[string]float64),
        StartTime:  time.Now(),
        LastUpdate: time.Now(),
    }

    return match
}

func generateMatchID() string {
    return fmt.Sprintf("match_%d", time.Now().UnixNano())
}

const (
    maxContextHistory = 100
)
