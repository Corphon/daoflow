//system/control/state/validator.go

package state

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
	"github.com/Corphon/daoflow/model"
)

const (
	maxMetricsHistory = 1000
)

// StateValidator 状态验证器
type StateValidator struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		validationTimeout time.Duration // 验证超时时间
		strictMode        bool          // 严格模式
		maxRules          int           // 最大规则数
		cacheEnabled      bool          // 启用缓存
	}

	// 验证状态
	state struct {
		rules   map[string]*ValidationRule  // 验证规则
		cache   map[string]*ValidationCache // 验证缓存
		results []ValidationResult          // 验证结果
		metrics ValidationMetrics           // 验证指标
	}
}

// ValidationRule 验证规则
type ValidationRule struct {
	ID        string        // 规则ID
	Name      string        // 规则名称
	Type      string        // 规则类型
	Target    string        // 验证目标
	Condition RuleCondition // 验证条件
	Priority  int           // 优先级
	Enabled   bool          // 是否启用
}

// RuleCondition 规则条件
type RuleCondition struct {
	Type       string                 // 条件类型
	Expression string                 // 条件表达式
	Parameters map[string]interface{} // 条件参数
	ErrorMsg   string                 // 错误消息
}

// ValidationCache 验证缓存
type ValidationCache struct {
	StateID     string          // 状态ID
	RuleResults map[string]bool // 规则结果
	Timestamp   time.Time       // 缓存时间
	TTL         time.Duration   // 生存时间
}

// ValidationResult 验证结果
type ValidationResult struct {
	ID        string                 // 结果ID
	StateID   string                 // 状态ID
	RuleID    string                 // 规则ID
	Valid     bool                   // 是否有效
	Details   map[string]interface{} // 详细信息
	Timestamp time.Time              // 验证时间
}

// ValidationMetrics 验证指标
type ValidationMetrics struct {
	TotalValidations int           // 总验证次数
	SuccessRate      float64       // 成功率
	AverageLatency   time.Duration // 平均延迟
	CacheHitRate     float64       // 缓存命中率
	History          []MetricPoint // 历史指标
}

// MetricPoint 指标点
type MetricPoint struct {
	Timestamp time.Time
	Values    map[string]float64
}

// NewStateValidator 创建新的状态验证器
func NewStateValidator() *StateValidator {
	sv := &StateValidator{}

	// 初始化配置
	sv.config.validationTimeout = 5 * time.Second
	sv.config.strictMode = true
	sv.config.maxRules = 100
	sv.config.cacheEnabled = true

	// 初始化状态
	sv.state.rules = make(map[string]*ValidationRule)
	sv.state.cache = make(map[string]*ValidationCache)
	sv.state.results = make([]ValidationResult, 0)
	sv.state.metrics = ValidationMetrics{
		History: make([]MetricPoint, 0),
	}

	return sv
}

// RegisterRule 注册验证规则
func (sv *StateValidator) RegisterRule(rule *ValidationRule) error {
	if rule == nil {
		return model.WrapError(nil, model.ErrCodeValidation, "nil rule")
	}

	sv.mu.Lock()
	defer sv.mu.Unlock()

	// 验证规则
	if err := sv.validateRule(rule); err != nil {
		return err
	}

	// 检查规则数量限制
	if len(sv.state.rules) >= sv.config.maxRules {
		return model.WrapError(nil, model.ErrCodeLimit, "max rules reached")
	}

	// 存储规则
	sv.state.rules[rule.ID] = rule

	return nil
}

// validateRule 验证规则基本有效性
func (sv *StateValidator) validateRule(rule *ValidationRule) error {
	// 验证基本字段
	if rule.ID == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty rule ID")
	}

	if rule.Type == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty rule type")
	}

	// 验证条件
	if err := sv.validateRuleCondition(rule.Condition); err != nil {
		return err
	}

	// 验证优先级
	if rule.Priority < 0 {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid priority")
	}

	return nil
}

// validateRuleCondition 验证规则条件
func (sv *StateValidator) validateRuleCondition(condition RuleCondition) error {
	if condition.Type == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty condition type")
	}

	// 验证表达式
	if condition.Expression == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty condition expression")
	}

	// 验证参数
	if condition.Parameters == nil {
		return model.WrapError(nil, model.ErrCodeValidation, "nil condition parameters")
	}

	return nil
}

// ValidateState 验证系统状态
func (sv *StateValidator) ValidateState(state *SystemState) error {
	sv.mu.RLock()
	defer sv.mu.RUnlock()

	startTime := time.Now()

	// 检查缓存
	if sv.config.cacheEnabled {
		if result := sv.checkCache(state.ID); result != nil {
			return sv.processValidationResult(result)
		}
	}

	// 执行验证
	results, err := sv.executeValidation(state)
	if err != nil {
		return err
	}

	// 更新缓存
	if sv.config.cacheEnabled {
		sv.updateCache(state.ID, results)
	}

	// 更新指标
	sv.updateMetrics(startTime)

	// 处理结果
	return sv.processResults(results)
}

// checkCache 检查验证缓存
func (sv *StateValidator) checkCache(stateID string) *ValidationResult {
	cache, exists := sv.state.cache[stateID]
	if !exists {
		return nil
	}

	// 检查缓存是否过期
	if time.Since(cache.Timestamp) > cache.TTL {
		delete(sv.state.cache, stateID)
		return nil
	}

	// 构建验证结果
	result := &ValidationResult{
		ID:        core.GenerateID(),
		StateID:   stateID,
		Valid:     true,
		Timestamp: time.Now(),
		Details: map[string]interface{}{
			"cache_hit": true,
			"rules":     cache.RuleResults,
		},
	}

	return result
}

// processValidationResult 处理验证结果
func (sv *StateValidator) processValidationResult(result *ValidationResult) error {
	// 添加到结果历史
	sv.state.results = append(sv.state.results, *result)

	// 如果验证失败，返回错误
	if !result.Valid {
		return model.WrapError(nil, model.ErrCodeValidation,
			fmt.Sprintf("validation failed for state %s", result.StateID))
	}

	return nil
}

// updateCache 更新验证缓存
func (sv *StateValidator) updateCache(stateID string, results []ValidationResult) {
	// 创建新的缓存记录
	cache := &ValidationCache{
		StateID:     stateID,
		RuleResults: make(map[string]bool),
		Timestamp:   time.Now(),
		TTL:         5 * time.Minute, // 默认缓存时间
	}

	// 记录每个规则的结果
	for _, result := range results {
		cache.RuleResults[result.RuleID] = result.Valid
	}

	// 更新缓存
	sv.state.cache[stateID] = cache
}

// processResults 处理验证结果集
func (sv *StateValidator) processResults(results []ValidationResult) error {
	if len(results) == 0 {
		return nil
	}

	// 添加到结果历史
	sv.state.results = append(sv.state.results, results...)

	// 在严格模式下，任何失败都返回错误
	if sv.config.strictMode {
		for _, result := range results {
			if !result.Valid {
				return model.WrapError(nil, model.ErrCodeValidation,
					fmt.Sprintf("validation failed: %s", result.ID))
			}
		}
	}

	// 维护结果历史大小
	if len(sv.state.results) > maxMetricsHistory {
		sv.state.results = sv.state.results[1:]
	}

	return nil
}

// ValidateTransition 验证状态转换
func (sv *StateValidator) ValidateTransition(
	current, next *SystemState) error {

	sv.mu.RLock()
	defer sv.mu.RUnlock()

	// 验证转换规则
	for _, rule := range sv.state.rules {
		if !rule.Enabled {
			continue
		}

		if err := sv.validateTransitionRule(rule, current, next); err != nil {
			return err
		}
	}

	return nil
}

// validateTransitionRule 验证状态转换规则
func (sv *StateValidator) validateTransitionRule(
	rule *ValidationRule,
	current, next *SystemState) error {

	// 检查规则类型
	if rule.Type != "transition" {
		return nil // 非转换规则，跳过
	}

	// 验证基础条件
	if current == nil || next == nil {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid state objects")
	}

	// 构造验证结果
	result := ValidationResult{
		ID:        core.GenerateID(),
		StateID:   next.ID,
		RuleID:    rule.ID,
		Valid:     true,
		Timestamp: time.Now(),
		Details: map[string]interface{}{
			"source_state": current.ID,
			"target_state": next.ID,
		},
	}

	// 根据条件类型执行验证
	switch rule.Condition.Type {
	case "state_change":
		// 验证状态变更是否合法
		if err := sv.validateStateChange(current, next, rule.Condition); err != nil {
			result.Valid = false
			result.Details["error"] = err.Error()
			sv.state.results = append(sv.state.results, result)
			return err
		}

	case "property_change":
		// 验证属性变更是否符合规则
		if err := sv.validatePropertyChange(current, next, rule.Condition); err != nil {
			result.Valid = false
			result.Details["error"] = err.Error()
			sv.state.results = append(sv.state.results, result)
			return err
		}

	case "component_change":
		// 验证组件变更
		if err := sv.validateComponentChange(current, next, rule.Condition); err != nil {
			result.Valid = false
			result.Details["error"] = err.Error()
			sv.state.results = append(sv.state.results, result)
			return err
		}
	}

	// 记录验证结果
	sv.state.results = append(sv.state.results, result)

	return nil
}

// validateComponentChange 验证组件变更
func (sv *StateValidator) validateComponentChange(
	current, next *SystemState,
	condition RuleCondition) error {

	// 获取目标组件
	componentID, ok := condition.Parameters["component_id"].(string)
	if !ok {
		return model.WrapError(nil, model.ErrCodeValidation, "missing component id")
	}

	// 获取组件
	currentComp, existsCurrent := current.Components[componentID]
	nextComp, existsNext := next.Components[componentID]

	// 检查组件存在性
	if !existsCurrent || !existsNext {
		return model.WrapError(nil, model.ErrCodeValidation,
			fmt.Sprintf("component %s not found", componentID))
	}

	// 根据条件表达式验证变更
	switch condition.Expression {
	case "status_transition":
		// 验证状态转换
		allowedStatus, ok := condition.Parameters["allowed_status"].([]string)
		if !ok {
			return model.WrapError(nil, model.ErrCodeValidation,
				"invalid status configuration")
		}

		validTransition := false
		for _, status := range allowedStatus {
			if nextComp.Status == status {
				validTransition = true
				break
			}
		}

		if !validTransition {
			return model.WrapError(nil, model.ErrCodeValidation,
				fmt.Sprintf("invalid component status transition: %s -> %s",
					currentComp.Status, nextComp.Status))
		}

	case "resource_limit":
		// 验证资源限制
		maxResource, ok := condition.Parameters["max_resource"].(float64)
		if !ok {
			return model.WrapError(nil, model.ErrCodeValidation,
				"invalid resource limit configuration")
		}

		if nextComp.ResourceUsage > maxResource {
			return model.WrapError(nil, model.ErrCodeValidation,
				fmt.Sprintf("component resource usage exceeds limit: %.2f > %.2f",
					nextComp.ResourceUsage, maxResource))
		}

	case "dependency_check":
		// 验证依赖完整性
		for depID := range nextComp.Dependencies {
			if _, exists := next.Components[depID]; !exists {
				return model.WrapError(nil, model.ErrCodeValidation,
					fmt.Sprintf("missing dependent component: %s", depID))
			}
		}
	}

	return nil
}

// validateStateChange 验证状态变更
func (sv *StateValidator) validateStateChange(
	current, next *SystemState,
	condition RuleCondition) error {

	// 获取允许的状态转换
	allowedTransitions, ok := condition.Parameters["allowed_transitions"].(map[string][]string)
	if !ok {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid transition configuration")
	}

	// 检查转换是否允许
	if validStates, exists := allowedTransitions[current.Status]; exists {
		for _, validState := range validStates {
			if next.Status == validState {
				return nil
			}
		}
	}

	return model.WrapError(nil, model.ErrCodeValidation,
		fmt.Sprintf("invalid state transition from %s to %s", current.Status, next.Status))
}

// validatePropertyChange 验证属性变更
func (sv *StateValidator) validatePropertyChange(
	current, next *SystemState,
	condition RuleCondition) error {

	property, ok := condition.Parameters["property"].(string)
	if !ok {
		return model.WrapError(nil, model.ErrCodeValidation, "missing property parameter")
	}

	// 获取属性值
	currentVal, existsCurrent := current.Properties[property]
	nextVal, existsNext := next.Properties[property]

	// 检查属性存在性
	if !existsCurrent || !existsNext {
		return model.WrapError(nil, model.ErrCodeValidation,
			fmt.Sprintf("property %s not found", property))
	}

	// 根据条件表达式验证变更
	switch condition.Expression {
	case "no_decrease":
		if sv.compareValues(nextVal, currentVal) < 0 {
			return model.WrapError(nil, model.ErrCodeValidation,
				fmt.Sprintf("property %s cannot decrease", property))
		}
	case "no_increase":
		if sv.compareValues(nextVal, currentVal) > 0 {
			return model.WrapError(nil, model.ErrCodeValidation,
				fmt.Sprintf("property %s cannot increase", property))
		}
	}

	return nil
}

// compareValues 比较值大小
func (sv *StateValidator) compareValues(v1, v2 interface{}) int {
	switch val1 := v1.(type) {
	case float64:
		val2, ok := v2.(float64)
		if !ok {
			return 0
		}
		if val1 < val2 {
			return -1
		} else if val1 > val2 {
			return 1
		}
	case int:
		val2, ok := v2.(int)
		if !ok {
			return 0
		}
		if val1 < val2 {
			return -1
		} else if val1 > val2 {
			return 1
		}
	}
	return 0
}

// executeValidation 执行验证
func (sv *StateValidator) executeValidation(
	state *SystemState) ([]ValidationResult, error) {

	results := make([]ValidationResult, 0)

	// 按优先级排序规则
	rules := sv.sortRulesByPriority()

	// 执行每个规则
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		result := sv.evaluateRule(rule, state)
		results = append(results, result)

		// 严格模式下，任何失败都立即返回
		if sv.config.strictMode && !result.Valid {
			return results, sv.createValidationError(result)
		}
	}

	return results, nil
}

// sortRulesByPriority 按优先级排序规则
func (sv *StateValidator) sortRulesByPriority() []*ValidationRule {
	rules := make([]*ValidationRule, 0, len(sv.state.rules))
	for _, rule := range sv.state.rules {
		rules = append(rules, rule)
	}

	// 按优先级降序排序
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority > rules[j].Priority
	})

	return rules
}

// evaluateRule 执行单个规则验证
func (sv *StateValidator) evaluateRule(rule *ValidationRule, state *SystemState) ValidationResult {
	result := ValidationResult{
		ID:        core.GenerateID(),
		StateID:   state.ID,
		RuleID:    rule.ID,
		Valid:     true,
		Timestamp: time.Now(),
		Details:   make(map[string]interface{}),
	}

	// 根据规则类型执行验证
	switch rule.Type {
	case "state":
		if err := sv.validateStateChange(state, state, rule.Condition); err != nil {
			result.Valid = false
			result.Details["error"] = err.Error()
		}
	case "property":
		if err := sv.validatePropertyChange(state, state, rule.Condition); err != nil {
			result.Valid = false
			result.Details["error"] = err.Error()
		}
	case "component":
		if err := sv.validateComponentChange(state, state, rule.Condition); err != nil {
			result.Valid = false
			result.Details["error"] = err.Error()
		}
	}

	return result
}

// createValidationError 创建验证错误
func (sv *StateValidator) createValidationError(result ValidationResult) error {
	return model.WrapError(nil, model.ErrCodeValidation,
		fmt.Sprintf("validation failed: %s - %v",
			result.RuleID,
			result.Details["error"]))
}

// 辅助函数

func (sv *StateValidator) updateMetrics(startTime time.Time) {
	duration := time.Since(startTime)

	point := MetricPoint{
		Timestamp: time.Now(),
		Values: map[string]float64{
			"latency":        float64(duration.Milliseconds()),
			"success_rate":   sv.calculateSuccessRate(),
			"cache_hit_rate": sv.calculateCacheHitRate(),
		},
	}

	sv.state.metrics.History = append(sv.state.metrics.History, point)

	// 限制历史记录数量
	if len(sv.state.metrics.History) > maxMetricsHistory {
		sv.state.metrics.History = sv.state.metrics.History[1:]
	}
}

// calculateSuccessRate 计算验证成功率
func (sv *StateValidator) calculateSuccessRate() float64 {
	if len(sv.state.results) == 0 {
		return 1.0 // 无验证记录时返回完美分数
	}

	successCount := 0
	for _, result := range sv.state.results {
		if result.Valid {
			successCount++
		}
	}

	return float64(successCount) / float64(len(sv.state.results))
}

// calculateCacheHitRate 计算缓存命中率
func (sv *StateValidator) calculateCacheHitRate() float64 {
	if sv.state.metrics.TotalValidations == 0 {
		return 0
	}

	hitCount := 0
	for _, result := range sv.state.results {
		if details, ok := result.Details["cache_hit"].(bool); ok && details {
			hitCount++
		}
	}

	return float64(hitCount) / float64(sv.state.metrics.TotalValidations)
}
