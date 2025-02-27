package adaptation

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/evolution/mutation"
	"github.com/Corphon/daoflow/system/evolution/pattern"
	"github.com/Corphon/daoflow/system/types"
)

const (
	maxHistoryLength = 1000 // 最大历史记录长度
	maxRules         = 100  // 最大规则数量
)

// AdaptationStrategy 适应策略管理器
type AdaptationStrategy struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		strategyUpdateInterval time.Duration // 策略更新间隔
		maxStrategies          int           // 最大策略数
		minEffectiveness       float64       // 最小有效性
		adaptiveThreshold      float64       // 自适应阈值
	}

	// 策略状态
	state struct {
		strategies map[string]*Strategy     // 当前策略
		rules      map[string]*StrategyRule // 策略规则
		history    []StrategyEvent          // 策略历史
		metrics    StrategyMetrics          // 策略指标
	}

	// 依赖项
	patternMatcher  *pattern.EvolutionMatcher
	mutationHandler *mutation.MutationHandler
}

// Strategy 适应策略
type Strategy struct {
	ID            string                 // 策略ID
	Type          string                 // 策略类型
	Priority      int                    // 优先级
	Rules         []string               // 关联规则
	Parameters    map[string]interface{} // 策略参数
	Conditions    []StrategyCondition    // 触发条件
	Actions       []StrategyAction       // 执行动作
	Effectiveness float64                // 有效性评分
	Created       time.Time              // 创建时间
	LastUsed      time.Time              // 最后使用时间
}

// StrategyRule 策略规则
type StrategyRule struct {
	ID        string        // 规则ID
	Name      string        // 规则名称
	Type      string        // 规则类型
	Target    string        // 规则目标
	Condition RuleCondition // 规则条件
	Action    RuleAction    // 规则动作
	Weight    float64       // 规则权重
	Enabled   bool          // 是否启用
}

// StrategyCondition 策略条件
type StrategyCondition struct {
	Type     string      // 条件类型
	Target   string      // 目标对象
	Operator string      // 操作符
	Value    interface{} // 比较值
	Priority int         // 优先级
}

// StrategyAction 策略动作
type StrategyAction struct {
	Type       string                 // 动作类型
	Target     string                 // 目标对象
	Operation  string                 // 操作类型
	Parameters map[string]interface{} // 动作参数
	Timeout    time.Duration          // 超时时间
}

// RuleCondition 规则条件
type RuleCondition struct {
	Expression string                 // 条件表达式
	Parameters map[string]interface{} // 条件参数
	Threshold  float64                // 阈值
}

// RuleAction 规则动作
type RuleAction struct {
	Function   string                 // 执行函数
	Parameters map[string]interface{} // 执行参数
	ResultType string                 // 结果类型
}

// StrategyEvent 策略事件
type StrategyEvent struct {
	Timestamp  time.Time
	StrategyID string
	Type       string
	Status     string
	Details    map[string]interface{}
}

// StrategyMetrics 策略指标
type StrategyMetrics struct {
	TotalExecutions int
	SuccessRate     float64
	AverageLatency  time.Duration
	Effectiveness   map[string]float64
	History         []MetricPoint
}

// MetricPoint 指标点
type MetricPoint struct {
	Timestamp time.Time
	Values    map[string]float64
}

// NewAdaptationStrategy 创建新的适应策略管理器
func NewAdaptationStrategy(matcher *pattern.EvolutionMatcher, handler *mutation.MutationHandler) (*AdaptationStrategy, error) {
	if matcher == nil {
		return nil, fmt.Errorf("nil evolution matcher")
	}
	if handler == nil {
		return nil, fmt.Errorf("nil mutation handler")
	}

	as := &AdaptationStrategy{
		patternMatcher:  matcher,
		mutationHandler: handler,
	}

	// 初始化配置
	as.config.strategyUpdateInterval = time.Hour
	as.config.maxStrategies = 100
	as.config.minEffectiveness = 0.5
	as.config.adaptiveThreshold = 0.7

	// 初始化状态
	as.state.strategies = make(map[string]*Strategy)
	as.state.rules = make(map[string]*StrategyRule)
	as.state.history = make([]StrategyEvent, 0)
	as.state.metrics = StrategyMetrics{
		Effectiveness: make(map[string]float64),
		History:       make([]MetricPoint, 0),
	}

	return as, nil
}

// GetRecentResults 获取最近的策略执行结果
func (as *AdaptationStrategy) GetRecentResults() ([]StrategyEvent, error) {
	as.mu.RLock()
	defer as.mu.RUnlock()

	results := make([]StrategyEvent, 0)

	// 获取最近的执行记录
	for i := len(as.state.history) - 1; i >= 0; i-- {
		event := as.state.history[i]
		if event.Type == "execution" {
			results = append(results, event)
		}
		// 只返回最近的记录
		if len(results) >= 100 {
			break
		}
	}

	return results, nil
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

// updateRules 更新规则状态
func (as *AdaptationStrategy) updateRules() error {
	as.mu.Lock()
	defer as.mu.Unlock()

	// 遍历规则进行检查和更新
	for id, rule := range as.state.rules {
		// 检查规则有效性
		effectiveness := as.evaluateRuleEffectiveness(rule)

		// 处理低效规则
		if effectiveness < as.config.minEffectiveness {
			// 尝试优化规则
			if optimized := as.optimizeRule(rule); optimized != nil {
				as.state.rules[id] = optimized
				as.recordStrategyEvent(rule, "rule_optimized", map[string]interface{}{
					"old_effectiveness": effectiveness,
					"new_effectiveness": as.evaluateRuleEffectiveness(optimized),
				})
			}
		}
	}

	return nil
}

// evaluateRuleEffectiveness 评估规则有效性
func (as *AdaptationStrategy) evaluateRuleEffectiveness(rule *StrategyRule) float64 {
	// 获取规则的历史记录
	ruleEvents := as.getRuleEvents(rule.ID)

	if len(ruleEvents) == 0 {
		return 0
	}

	// 计算成功率
	successCount := 0
	for _, event := range ruleEvents {
		if event.Status == "success" {
			successCount++
		}
	}

	// 计算有效性得分
	effectiveness := float64(successCount) / float64(len(ruleEvents))

	// 考虑规则权重
	return effectiveness * rule.Weight
}

// optimizeRule 优化规则实现
func (as *AdaptationStrategy) optimizeRule(rule *StrategyRule) *StrategyRule {
	// 创建规则副本
	optimized := *rule

	// 基于历史数据优化条件阈值
	if events := as.getRuleEvents(rule.ID); len(events) > 0 {
		optimized.Condition.Threshold = as.findOptimalThreshold(events)
	}

	// 优化动作参数
	if params := as.optimizeActionParameters(rule); len(params) > 0 {
		optimized.Action.Parameters = params
	}

	// 调整规则权重
	optimized.Weight = as.calculateOptimizedWeight(rule)

	return &optimized
}

// recordStrategyEvent 支持记录规则事件
func (as *AdaptationStrategy) recordStrategyEvent(
	source interface{},
	eventType string,
	details map[string]interface{}) {

	var event StrategyEvent

	// 根据源对象类型设置事件属性
	switch src := source.(type) {
	case *Strategy:
		event = StrategyEvent{
			Timestamp:  time.Now(),
			StrategyID: src.ID,
			Type:       eventType,
			Status:     "completed",
			Details:    details,
		}
	case *StrategyRule:
		event = StrategyEvent{
			Timestamp:  time.Now(),
			StrategyID: src.ID,
			Type:       "rule_" + eventType,
			Status:     "completed",
			Details:    details,
		}
	default:
		return
	}

	as.state.history = append(as.state.history, event)

	// 限制历史记录长度
	if len(as.state.history) > maxHistoryLength {
		as.state.history = as.state.history[1:]
	}
}

// 辅助方法

func (as *AdaptationStrategy) getRuleEvents(ruleID string) []StrategyEvent {
	events := make([]StrategyEvent, 0)
	for _, event := range as.state.history {
		if event.StrategyID == ruleID {
			events = append(events, event)
		}
	}
	return events
}

func (as *AdaptationStrategy) findOptimalThreshold(events []StrategyEvent) float64 {
	// 收集成功事件的阈值
	successThresholds := make([]float64, 0)
	for _, event := range events {
		if event.Status == "success" {
			if threshold, ok := event.Details["threshold"].(float64); ok {
				successThresholds = append(successThresholds, threshold)
			}
		}
	}

	if len(successThresholds) == 0 {
		return 0.5 // 默认阈值
	}

	// 计算平均最优阈值
	total := 0.0
	for _, t := range successThresholds {
		total += t
	}
	return total / float64(len(successThresholds))
}

func (as *AdaptationStrategy) optimizeActionParameters(rule *StrategyRule) map[string]interface{} {
	events := as.getRuleEvents(rule.ID)
	params := make(map[string]interface{})

	// 从成功事件中提取有效参数
	for _, event := range events {
		if event.Status == "success" {
			if eventParams, ok := event.Details["parameters"].(map[string]interface{}); ok {
				for k, v := range eventParams {
					params[k] = v
				}
			}
		}
	}

	return params
}

func (as *AdaptationStrategy) calculateOptimizedWeight(rule *StrategyRule) float64 {
	effectiveness := as.evaluateRuleEffectiveness(rule)
	return math.Max(0.1, math.Min(1.0, effectiveness*1.2)) // 略微提升权重
}

// cleanupStrategies 清理无效策略
func (as *AdaptationStrategy) cleanupStrategies() {
	as.mu.Lock()
	defer as.mu.Unlock()

	// 获取当前时间
	now := time.Now()

	// 清理过期或无效策略
	for id, strategy := range as.state.strategies {
		// 检查有效性条件
		if strategy.Effectiveness < as.config.minEffectiveness {
			// 记录清理事件
			as.recordStrategyEvent(strategy, "strategy_cleaned", map[string]interface{}{
				"reason": "low_effectiveness",
				"value":  strategy.Effectiveness,
			})
			delete(as.state.strategies, id)
			continue
		}

		// 检查最后使用时间
		if now.Sub(strategy.LastUsed) > 24*time.Hour {
			// 记录清理事件
			as.recordStrategyEvent(strategy, "strategy_cleaned", map[string]interface{}{
				"reason":    "expired",
				"last_used": strategy.LastUsed,
			})
			delete(as.state.strategies, id)
		}
	}
}

// updateMetrics 更新策略指标
func (as *AdaptationStrategy) updateMetrics() {
	point := MetricPoint{
		Timestamp: time.Now(),
		Values: map[string]float64{
			"total_strategies":      float64(len(as.state.strategies)),
			"total_rules":           float64(len(as.state.rules)),
			"success_rate":          as.calculateSuccessRate(),
			"average_effectiveness": as.calculateAverageEffectiveness(),
			"execution_latency":     float64(as.state.metrics.AverageLatency.Milliseconds()),
		},
	}

	// 更新历史记录
	as.state.metrics.History = append(as.state.metrics.History, point)

	// 限制历史记录长度
	if len(as.state.metrics.History) > maxHistoryLength {
		as.state.metrics.History = as.state.metrics.History[1:]
	}

	// 更新总体指标
	as.state.metrics.TotalExecutions++
	as.state.metrics.SuccessRate = as.calculateSuccessRate()
}

// 辅助函数
func (as *AdaptationStrategy) calculateSuccessRate() float64 {
	if as.state.metrics.TotalExecutions == 0 {
		return 0
	}

	successCount := 0
	for _, event := range as.state.history {
		if event.Status == "success" {
			successCount++
		}
	}

	return float64(successCount) / float64(as.state.metrics.TotalExecutions)
}

func (as *AdaptationStrategy) calculateAverageEffectiveness() float64 {
	if len(as.state.strategies) == 0 {
		return 0
	}

	total := 0.0
	for _, strategy := range as.state.strategies {
		total += strategy.Effectiveness
	}

	return total / float64(len(as.state.strategies))
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

// removeWorstStrategy 移除最不有效的策略
func (as *AdaptationStrategy) removeWorstStrategy() {
	// 寻找效果最差的策略
	var worstStrategy *Strategy
	worstEffectiveness := math.MaxFloat64

	for _, strategy := range as.state.strategies {
		if strategy.Effectiveness < worstEffectiveness {
			worstStrategy = strategy
			worstEffectiveness = strategy.Effectiveness
		}
	}

	if worstStrategy != nil {
		// 删除该策略
		delete(as.state.strategies, worstStrategy.ID)

		// 记录移除事件
		as.recordStrategyEvent(worstStrategy, "strategy_removed", map[string]interface{}{
			"reason":        "low_effectiveness",
			"effectiveness": worstEffectiveness,
		})
	}
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

// getStrategyEvents 获取策略相关事件
func (as *AdaptationStrategy) getStrategyEvents(strategyID string) []StrategyEvent {
	events := make([]StrategyEvent, 0)

	// 遍历历史记录找到相关事件
	for _, event := range as.state.history {
		if event.StrategyID == strategyID {
			events = append(events, event)
		}
	}

	// 按时间顺序排序
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.Before(events[j].Timestamp)
	})

	// 只返回最近的事件(避免历史数据过多影响评估)
	maxEvents := 100 // 限制事件数量
	if len(events) > maxEvents {
		events = events[len(events)-maxEvents:]
	}

	return events
}

// evaluateStrategy 评估策略有效性
func (as *AdaptationStrategy) evaluateStrategy(strategy *Strategy) float64 {
	// 获取策略相关事件
	events := as.getStrategyEvents(strategy.ID)
	if len(events) == 0 {
		return 0
	}

	// 计算成功率
	successCount := 0
	weightedScore := 0.0
	totalWeight := 0.0

	for _, event := range events {
		if event.Status == "success" {
			successCount++
			// 考虑时间衰减
			age := time.Since(event.Timestamp).Hours()
			weight := math.Exp(-age / 24.0) // 24小时衰减
			weightedScore += weight
			totalWeight += weight
		}
	}

	// 基础有效性得分
	baseScore := float64(successCount) / float64(len(events))

	// 加入时间加权得分
	timeWeightedScore := 0.0
	if totalWeight > 0 {
		timeWeightedScore = weightedScore / totalWeight
	}

	// 综合评分
	return baseScore*0.6 + timeWeightedScore*0.4
}

// optimizeStrategy 优化策略实现
func (as *AdaptationStrategy) optimizeStrategy(strategy *Strategy) error {
	// 获取历史执行数据
	events := as.getStrategyEvents(strategy.ID)
	if len(events) == 0 {
		return fmt.Errorf("no historical data for strategy optimization")
	}

	// 1. 优化参数
	optimizedParams := make(map[string]interface{})
	successEvents := filterSuccessEvents(events)
	if len(successEvents) > 0 {
		optimizedParams = extractOptimalParameters(successEvents)
		strategy.Parameters = optimizedParams
	}

	// 2. 优化条件
	for i := range strategy.Conditions {
		if threshold := findOptimalConditionThreshold(events, strategy.Conditions[i]); threshold > 0 {
			strategy.Conditions[i].Value = threshold
		}
	}

	// 3. 优化动作
	for i := range strategy.Actions {
		if params := optimizeActionParams(events, strategy.Actions[i]); len(params) > 0 {
			strategy.Actions[i].Parameters = params
		}
	}

	// 4. 更新优先级
	strategy.Priority = calculateOptimizedPriority(strategy, events)

	// 记录优化事件
	as.recordStrategyEvent(strategy, "optimized", map[string]interface{}{
		"new_parameters": optimizedParams,
		"new_priority":   strategy.Priority,
	})

	return nil
}

// 辅助函数
func filterSuccessEvents(events []StrategyEvent) []StrategyEvent {
	success := make([]StrategyEvent, 0)
	for _, event := range events {
		if event.Status == "success" {
			success = append(success, event)
		}
	}
	return success
}

func extractOptimalParameters(events []StrategyEvent) map[string]interface{} {
	params := make(map[string]interface{})
	paramValues := make(map[string][]float64)

	// 收集参数值
	for _, event := range events {
		if eventParams, ok := event.Details["parameters"].(map[string]interface{}); ok {
			for k, v := range eventParams {
				if fv, ok := v.(float64); ok {
					paramValues[k] = append(paramValues[k], fv)
				}
			}
		}
	}

	// 计算最优值
	for k, values := range paramValues {
		if len(values) > 0 {
			sort.Float64s(values)
			median := values[len(values)/2]
			mean := calculateMean(values)
			// 使用加权平均
			params[k] = median*0.6 + mean*0.4
		}
	}

	return params
}

func findOptimalConditionThreshold(events []StrategyEvent, condition StrategyCondition) float64 {
	values := make([]float64, 0)
	for _, event := range events {
		if event.Status == "success" {
			if v, ok := event.Details[condition.Target].(float64); ok {
				values = append(values, v)
			}
		}
	}

	if len(values) == 0 {
		return 0
	}

	sort.Float64s(values)
	median := values[len(values)/2]
	mean := calculateMean(values)
	return median*0.6 + mean*0.4
}

func optimizeActionParams(events []StrategyEvent, action StrategyAction) map[string]interface{} {
	params := make(map[string]interface{})
	for _, event := range events {
		if event.Status == "success" && event.Type == action.Type {
			if actionParams, ok := event.Details["action_params"].(map[string]interface{}); ok {
				for k, v := range actionParams {
					params[k] = v
				}
			}
		}
	}
	return params
}

func calculateOptimizedPriority(strategy *Strategy, events []StrategyEvent) int {
	// 基于成功率和影响程度计算优先级
	successRate := float64(0)
	if len(events) > 0 {
		successCount := 0
		for _, event := range events {
			if event.Status == "success" {
				successCount++
			}
		}
		successRate = float64(successCount) / float64(len(events))
	}

	// 优先级范围1-10
	return int(math.Max(1, math.Min(10, float64(strategy.Priority)*successRate*1.5)))
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

// getCurrentState 获取当前系统状态
func (as *AdaptationStrategy) getCurrentState() (*model.SystemState, error) {
	// 从依赖的 Handler 获取系统状态
	state, err := as.mutationHandler.GetCurrentState()
	if err != nil {
		return nil, err
	}
	return state, nil
}

// selectApplicableStrategies 选择适用的策略
func (as *AdaptationStrategy) selectApplicableStrategies(state *model.SystemState) []*Strategy {
	applicable := make([]*Strategy, 0)

	sysState := types.FromModelSystemState(state)
	for _, strategy := range as.state.strategies {
		if as.isStrategyApplicable(strategy, sysState) {
			applicable = append(applicable, strategy)
		}
	}

	return applicable
}

// isStrategyApplicable 检查策略是否适用
func (as *AdaptationStrategy) isStrategyApplicable(strategy *Strategy, state *types.SystemState) bool {
	// 检查每个条件
	for _, condition := range strategy.Conditions {
		if !as.evaluateCondition(condition, state) {
			return false
		}
	}
	return true
}

// evaluateCondition 评估条件
func (as *AdaptationStrategy) evaluateCondition(condition StrategyCondition, state *types.SystemState) bool {
	value, ok := state.Properties[condition.Target]
	if !ok {
		return false
	}

	// 根据操作符评估条件
	switch condition.Operator {
	case ">":
		return value.(float64) > condition.Value.(float64)
	case "<":
		return value.(float64) < condition.Value.(float64)
	case "==":
		return value == condition.Value
	case "!=":
		return value != condition.Value
	default:
		return false
	}
}

// sortStrategiesByPriority 按优先级排序策略
func (as *AdaptationStrategy) sortStrategiesByPriority(strategies []*Strategy) []*Strategy {
	sorted := make([]*Strategy, len(strategies))
	copy(sorted, strategies)

	sort.Slice(sorted, func(i, j int) bool {
		// 优先级高的排在前面
		return sorted[i].Priority > sorted[j].Priority
	})

	return sorted
}

// executeStrategy 执行单个策略
func (as *AdaptationStrategy) executeStrategy(strategy *Strategy, modelState *model.SystemState) error {
	// 记录开始执行
	startTime := time.Now()
	as.recordStrategyEvent(strategy, "execution_start", nil)

	state := types.FromModelSystemState(modelState)

	// 执行每个动作
	for _, action := range strategy.Actions {
		if err := as.executeAction(action, state); err != nil {
			return err
		}
	}

	// 记录执行成功
	as.recordStrategyEvent(strategy, "execution_complete", map[string]interface{}{
		"duration": time.Since(startTime).Milliseconds(),
		"state":    state,
	})

	return nil
}

// executeAction 执行动作
func (as *AdaptationStrategy) executeAction(action StrategyAction, state *types.SystemState) error {
	switch action.Operation {
	case "adjust":
		action.Parameters["current_state"] = state
		return as.adjustSystemParameter(action.Target, action.Parameters)
	case "optimize":
		action.Parameters["system_state"] = state
		return as.optimizeSystem(action.Parameters)
	case "transform":
		action.Parameters["state_info"] = state
		return as.transformSystem(action.Parameters)
	default:
		return fmt.Errorf("unknown action operation: %s", action.Operation)
	}
}

// 辅助方法
func (as *AdaptationStrategy) adjustSystemParameter(target string, params map[string]interface{}) error {
	// 通过 Handler 调整系统参数
	return as.mutationHandler.AdjustParameter(target, params)
}

func (as *AdaptationStrategy) optimizeSystem(params map[string]interface{}) error {
	// 通过 Handler 优化系统
	return as.mutationHandler.Optimize(params)
}

func (as *AdaptationStrategy) transformSystem(params map[string]interface{}) error {
	// 通过 Handler 转换系统
	return as.mutationHandler.Transform(params)
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

// validateCondition 验证策略条件
func (as *AdaptationStrategy) validateCondition(condition StrategyCondition) error {
	if condition.Type == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty condition type")
	}

	if condition.Target == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty condition target")
	}

	if condition.Operator == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty condition operator")
	}

	// 验证操作符
	validOperators := map[string]bool{
		">":  true,
		"<":  true,
		"==": true,
		"!=": true,
		">=": true,
		"<=": true,
	}
	if !validOperators[condition.Operator] {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid operator")
	}

	return nil
}

// validateAction 验证策略动作
func (as *AdaptationStrategy) validateAction(action StrategyAction) error {
	if action.Type == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty action type")
	}

	if action.Target == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty action target")
	}

	if action.Operation == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty action operation")
	}

	// 验证操作类型
	validOperations := map[string]bool{
		"adjust":    true,
		"optimize":  true,
		"transform": true,
	}
	if !validOperations[action.Operation] {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid operation")
	}

	// 验证超时设置
	if action.Timeout < 0 {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid timeout")
	}

	return nil
}

// UpdateParameters 更新策略参数
func (as *AdaptationStrategy) UpdateParameters(strategyType string, params map[string]interface{}) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	// 查找相应类型的策略
	var targetStrategy *Strategy
	for _, strategy := range as.state.strategies {
		if strategy.Type == strategyType {
			targetStrategy = strategy
			break
		}
	}

	if targetStrategy == nil {
		return fmt.Errorf("strategy type %s not found", strategyType)
	}

	// 验证参数
	if err := as.validateParameters(params); err != nil {
		return err
	}

	// 更新参数
	oldParams := targetStrategy.Parameters
	targetStrategy.Parameters = params

	// 记录更新事件
	as.recordStrategyEvent(targetStrategy, "parameters_updated", map[string]interface{}{
		"old_params": oldParams,
		"new_params": params,
	})

	return nil
}

// validateParameters 验证参数有效性
func (as *AdaptationStrategy) validateParameters(params map[string]interface{}) error {
	if params == nil {
		return fmt.Errorf("nil parameters")
	}

	// 验证必需参数
	requiredParams := []string{"weight", "threshold"}
	for _, required := range requiredParams {
		if _, exists := params[required]; !exists {
			return fmt.Errorf("missing required parameter: %s", required)
		}
	}

	return nil
}

// RegisterRule 注册新规则
func (as *AdaptationStrategy) RegisterRule(rule *StrategyRule) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	// 验证规则
	if err := as.validateRule(rule); err != nil {
		return err
	}

	// 检查规则存在性
	if _, exists := as.state.rules[rule.ID]; exists {
		return fmt.Errorf("rule %s already exists", rule.ID)
	}

	// 检查规则数量限制
	if len(as.state.rules) >= maxRules {
		return fmt.Errorf("max rules limit reached")
	}

	// 存储规则
	as.state.rules[rule.ID] = rule

	// 记录事件
	as.recordStrategyEvent(rule, "rule_registered", map[string]interface{}{
		"rule_type": rule.Type,
		"target":    rule.Target,
	})

	return nil
}

// validateRule 验证规则有效性
func (as *AdaptationStrategy) validateRule(rule *StrategyRule) error {
	if rule == nil {
		return fmt.Errorf("nil rule")
	}

	if rule.ID == "" {
		return fmt.Errorf("empty rule ID")
	}

	if rule.Type == "" {
		return fmt.Errorf("empty rule type")
	}

	if rule.Target == "" {
		return fmt.Errorf("empty rule target")
	}

	return nil
}

// GetRules 获取所有规则
func (as *AdaptationStrategy) GetRules() []*StrategyRule {
	as.mu.RLock()
	defer as.mu.RUnlock()

	rules := make([]*StrategyRule, 0, len(as.state.rules))
	for _, rule := range as.state.rules {
		rules = append(rules, rule)
	}

	return rules
}

// UpdateRule 更新规则
func (as *AdaptationStrategy) UpdateRule(rule *StrategyRule) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	// 验证规则
	if err := as.validateRule(rule); err != nil {
		return err
	}

	// 检查规则存在性
	oldRule, exists := as.state.rules[rule.ID]
	if !exists {
		return fmt.Errorf("rule %s not found", rule.ID)
	}

	// 更新规则
	as.state.rules[rule.ID] = rule

	// 记录更新事件
	as.recordStrategyEvent(rule, "rule_updated", map[string]interface{}{
		"old_weight":    oldRule.Weight,
		"new_weight":    rule.Weight,
		"old_condition": oldRule.Condition,
		"new_condition": rule.Condition,
	})

	return nil
}
