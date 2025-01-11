//system/evolution/adaptation/strategy.go

package adaptation

import (
    "sync"
    "time"

    "github.com/Corphon/daoflow/evolution/pattern"
    "github.com/Corphon/daoflow/evolution/mutation"
    "github.com/Corphon/daoflow/meta/field"
    "github.com/Corphon/daoflow/model"
    "github.com/Corphon/daoflow/system/types"
)

// AdaptationStrategy 适应策略管理器
type AdaptationStrategy struct {
    mu sync.RWMutex

    // 基础配置
    config struct {
        strategyUpdateInterval time.Duration // 策略更新间隔
        maxStrategies         int           // 最大策略数
        minEffectiveness      float64       // 最小有效性
        adaptiveThreshold     float64       // 自适应阈值
    }

    // 策略状态
    state struct {
        strategies  map[string]*Strategy     // 当前策略
        rules      map[string]*StrategyRule // 策略规则
        history    []StrategyEvent          // 策略历史
        metrics    StrategyMetrics          // 策略指标
    }

    // 依赖项
    patternMatcher *pattern.EvolutionMatcher
    mutationHandler *mutation.MutationHandler
}

// Strategy 适应策略
type Strategy struct {
    ID            string                // 策略ID
    Type          string                // 策略类型
    Priority      int                   // 优先级
    Rules         []string              // 关联规则
    Parameters    map[string]interface{} // 策略参数
    Conditions    []StrategyCondition   // 触发条件
    Actions       []StrategyAction      // 执行动作
    Effectiveness float64               // 有效性评分
    Created       time.Time             // 创建时间
    LastUsed      time.Time             // 最后使用时间
}

// StrategyRule 策略规则
type StrategyRule struct {
    ID           string                // 规则ID
    Name         string                // 规则名称
    Type         string                // 规则类型
    Target       string                // 规则目标
    Condition    RuleCondition         // 规则条件
    Action       RuleAction            // 规则动作
    Weight       float64               // 规则权重
    Enabled      bool                  // 是否启用
}

// StrategyCondition 策略条件
type StrategyCondition struct {
    Type         string                // 条件类型
    Target       string                // 目标对象
    Operator     string                // 操作符
    Value        interface{}           // 比较值
    Priority     int                   // 优先级
}

// StrategyAction 策略动作
type StrategyAction struct {
    Type         string                // 动作类型
    Target       string                // 目标对象
    Operation    string                // 操作类型
    Parameters   map[string]interface{} // 动作参数
    Timeout      time.Duration         // 超时时间
}

// RuleCondition 规则条件
type RuleCondition struct {
    Expression   string                // 条件表达式
    Parameters   map[string]interface{} // 条件参数
    Threshold    float64               // 阈值
}

// RuleAction 规则动作
type RuleAction struct {
    Function     string                // 执行函数
    Parameters   map[string]interface{} // 执行参数
    ResultType   string                // 结果类型
}

// StrategyEvent 策略事件
type StrategyEvent struct {
    Timestamp    time.Time
    StrategyID   string
    Type         string
    Status       string
    Details      map[string]interface{}
}

// StrategyMetrics 策略指标
type StrategyMetrics struct {
    TotalExecutions  int
    SuccessRate     float64
    AverageLatency  time.Duration
    Effectiveness   map[string]float64
    History         []MetricPoint
}

// MetricPoint 指标点
type MetricPoint struct {
    Timestamp    time.Time
    Values       map[string]float64
}

// NewAdaptationStrategy 创建新的适应策略管理器
func NewAdaptationStrategy(
    matcher *pattern.EvolutionMatcher,
    handler *mutation.MutationHandler) *AdaptationStrategy {
    
    as := &AdaptationStrategy{
        patternMatcher:  matcher,
        mutationHandler: handler,
    }

    // 初始化配置
    as.config.strategyUpdateInterval = 30 * time.Minute
    as.config.maxStrategies = 100
    as.config.minEffectiveness = 0.6
    as.config.adaptiveThreshold = 0.75

    // 初始化状态
    as.state.strategies = make(map[string]*Strategy)
    as.state.rules = make(map[string]*StrategyRule)
    as.state.history = make([]StrategyEvent, 0)
    as.state.metrics = StrategyMetrics{
        Effectiveness: make(map[string]float64),
        History:      make([]MetricPoint, 0),
    }

    return as
}

// Execute 执行策略管理
func (as *AdaptationStrategy) Execute() error {
    as.mu.Lock()
    defer as.mu.Unlock()

    // 更新策略评估
    if err := as.updateStrategyEvaluation(); err != nil {
        return err
    }

    // 选择和应用策略
    if err := as.applyStrategies(); err != nil {
        return err
    }

    // 更新规则状态
    if err := as.updateRules(); err != nil {
        return err
    }

    // 清理无效策略
    as.cleanupStrategies()

    // 更新指标
    as.updateMetrics()

    return nil
}

// RegisterStrategy 注册新策略
func (as *AdaptationStrategy) RegisterStrategy(strategy *Strategy) error {
    if strategy == nil {
        return model.WrapError(nil, model.ErrCodeValidation, "nil strategy")
    }

    as.mu.Lock()
    defer as.mu.Unlock()

    // 验证策略
    if err := as.validateStrategy(strategy); err != nil {
        return err
    }

    // 检查策略数量限制
    if len(as.state.strategies) >= as.config.maxStrategies {
        // 移除最不有效的策略
        as.removeWorstStrategy()
    }

    // 存储策略
    as.state.strategies[strategy.ID] = strategy

    // 记录事件
    as.recordStrategyEvent(strategy, "registered", nil)

    return nil
}

// updateStrategyEvaluation 更新策略评估
func (as *AdaptationStrategy) updateStrategyEvaluation() error {
    for _, strategy := range as.state.strategies {
        // 评估策略有效性
        effectiveness := as.evaluateStrategy(strategy)
        
        // 更新评分
        strategy.Effectiveness = effectiveness

        // 检查是否需要调整
        if effectiveness < as.config.minEffectiveness {
            // 尝试优化策略
            if err := as.optimizeStrategy(strategy); err != nil {
                continue
            }
        }
    }

    return nil
}

// applyStrategies 应用策略
func (as *AdaptationStrategy) applyStrategies() error {
    // 获取当前系统状态
    state, err := as.getCurrentState()
    if err != nil {
        return err
    }

    // 选择适用的策略
    applicable := as.selectApplicableStrategies(state)

    // 按优先级排序
    sortedStrategies := as.sortStrategiesByPriority(applicable)

    // 执行策略
    for _, strategy := range sortedStrategies {
        if err := as.executeStrategy(strategy, state); err != nil {
            // 记录错误但继续执行其他策略
            as.recordStrategyEvent(strategy, "execution_error", map[string]interface{}{
                "error": err.Error(),
            })
            continue
        }

        // 更新使用时间
        strategy.LastUsed = time.Now()
    }

    return nil
}

// 辅助函数

func (as *AdaptationStrategy) validateStrategy(strategy *Strategy) error {
    if strategy.ID == "" {
        return model.WrapError(nil, model.ErrCodeValidation, "empty strategy ID")
    }

    // 验证条件
    for _, condition := range strategy.Conditions {
        if err := as.validateCondition(condition); err != nil {
            return err
        }
    }

    // 验证动作
    for _, action := range strategy.Actions {
        if err := as.validateAction(action); err != nil {
            return err
        }
    }

    return nil
}

func (as *AdaptationStrategy) recordStrategyEvent(
    strategy *Strategy,
    eventType string,
    details map[string]interface{}) {
    
    event := StrategyEvent{
        Timestamp:  time.Now(),
        StrategyID: strategy.ID,
        Type:      eventType,
        Status:    "completed",
        Details:   details,
    }

    as.state.history = append(as.state.history, event)

    // 限制历史记录长度
    if len(as.state.history) > maxHistoryLength {
        as.state.history = as.state.history[1:]
    }
}

const (
    maxHistoryLength = 1000
)
