// system/meta/emergence/rules.go

package emergence

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/meta/field"
)

// RuleEngine 涌现规则引擎
type RuleEngine struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		threshold     float64 // 规则触发阈值
		minConfidence float64 // 最小置信度
		maxRules      int     // 最大规则数量
	}

	// 规则状态
	state struct {
		rules      map[string]*EmergenceRule    // 活跃规则
		patterns   map[string]PatternState      // 模式状态
		properties map[string]*EmergentProperty // 属性映射
		history    []RuleEvent                  // 规则历史
	}

	// 依赖项
	detector  *PatternDetector
	generator *PropertyGenerator
	field     *field.UnifiedField
}

// EmergenceRule 涌现规则
type EmergenceRule struct {
	ID            string          // 规则ID
	Name          string          // 规则名称
	Type          string          // 规则类型
	Conditions    []RuleCondition // 触发条件
	Actions       []RuleAction    // 执行动作
	Priority      int             // 优先级
	Confidence    float64         // 置信度
	Created       time.Time       // 创建时间
	LastTriggered time.Time       // 最后触发时间
	TriggerCount  int             // 触发次数
}

// RuleCondition 规则条件
type RuleCondition struct {
	Type     string      // 条件类型
	Target   string      // 目标对象
	Operator string      // 操作符
	Value    interface{} // 比较值
	Weight   float64     // 权重
}

// RuleAction 规则动作
type RuleAction struct {
	Type       string                 // 动作类型
	Target     string                 // 目标对象
	Operation  string                 // 操作类型
	Parameters map[string]interface{} // 参数
}

// PatternState 模式状态
type PatternState struct {
	Pattern    *EmergentPattern
	Active     bool
	Duration   time.Duration
	Strength   float64
	LastUpdate time.Time
	Properties map[string]float64 // 状态属性
	Energy     float64            // 能量值
	Timestamp  time.Time          // 时间戳
}

// RuleEvent 规则事件
type RuleEvent struct {
	Timestamp time.Time
	RuleID    string
	Type      string
	Success   bool
	Details   map[string]interface{}
}

// PatternCorrelation 模式关联
type PatternCorrelation struct {
	Source    PatternState
	Target    PatternState
	Strength  float64
	Effect    float64
	Threshold float64
}

// NewRuleEngine 创建新的规则引擎
func NewRuleEngine(
	detector *PatternDetector,
	generator *PropertyGenerator,
	field *field.UnifiedField) *RuleEngine {

	re := &RuleEngine{
		detector:  detector,
		generator: generator,
		field:     field,
	}

	// 初始化配置
	re.config.threshold = 0.7
	re.config.minConfidence = 0.65
	re.config.maxRules = 1000

	// 初始化状态
	re.state.rules = make(map[string]*EmergenceRule)
	re.state.patterns = make(map[string]PatternState)
	re.state.properties = make(map[string]*EmergentProperty)
	re.state.history = make([]RuleEvent, 0)

	return re
}

// ProcessRules 处理涌现规则
func (re *RuleEngine) ProcessRules() error {
	re.mu.Lock()
	defer re.mu.Unlock()

	// 更新模式状态
	if err := re.updatePatternStates(); err != nil {
		return err
	}

	// 评估并触发规则
	triggeredRules, err := re.evaluateRules()
	if err != nil {
		return err
	}

	// 执行规则动作
	if err := re.executeRules(triggeredRules); err != nil {
		return err
	}

	// 学习新规则
	if err := re.learnNewRules(); err != nil {
		return err
	}

	return nil
}

// updatePatternStates 更新模式状态
func (re *RuleEngine) updatePatternStates() error {
	// 获取当前模式
	patterns, err := re.detector.Detect()
	if err != nil {
		return model.WrapError(err, model.ErrCodeOperation, "failed to detect patterns")
	}

	// 更新状态
	currentTime := time.Now()
	for _, pattern := range patterns {
		state, exists := re.state.patterns[pattern.ID]
		if !exists {
			// 新模式
			state = PatternState{
				Pattern:    &pattern,
				Active:     true,
				Duration:   0,
				Strength:   pattern.Strength,
				LastUpdate: currentTime,
			}
		} else {
			// 更新现有模式
			state.Active = true
			state.Duration += currentTime.Sub(state.LastUpdate)
			state.Strength = pattern.Strength
			state.LastUpdate = currentTime
		}
		re.state.patterns[pattern.ID] = state
	}

	// 处理消失的模式
	for id, state := range re.state.patterns {
		if state.LastUpdate != currentTime {
			state.Active = false
			re.state.patterns[id] = state
		}
	}

	return nil
}

// evaluateRules 评估规则
func (re *RuleEngine) evaluateRules() ([]*EmergenceRule, error) {
	triggered := make([]*EmergenceRule, 0)

	// 按优先级排序规则
	rules := re.getSortedRules()

	// 评估每个规则
	for _, rule := range rules {
		if re.shouldTriggerRule(rule) {
			triggered = append(triggered, rule)
		}
	}

	return triggered, nil
}

// shouldTriggerRule 判断是否应该触发规则
func (re *RuleEngine) shouldTriggerRule(rule *EmergenceRule) bool {
	// 检查置信度
	if rule.Confidence < re.config.minConfidence {
		return false
	}

	// 评估所有条件
	totalWeight := 0.0
	satisfiedWeight := 0.0

	for _, cond := range rule.Conditions {
		totalWeight += cond.Weight
		if re.evaluateCondition(cond) {
			satisfiedWeight += cond.Weight
		}
	}

	if totalWeight == 0 {
		return false
	}

	// 计算满足度
	satisfaction := satisfiedWeight / totalWeight
	return satisfaction >= re.config.threshold
}

// evaluateCondition 评估条件
func (re *RuleEngine) evaluateCondition(cond RuleCondition) bool {
	switch cond.Type {
	case "pattern":
		return re.evaluatePatternCondition(cond)
	case "property":
		return re.evaluatePropertyCondition(cond)
	case "field":
		return re.evaluateFieldCondition(cond)
	default:
		return false
	}
}

// evaluatePatternCondition 评估模式条件
func (re *RuleEngine) evaluatePatternCondition(cond RuleCondition) bool {
	// 获取模式名称和期望属性
	patternName := cond.Target
	expectedValue := cond.Value.(float64)

	// 在当前模式中查找
	for _, pattern := range re.detector.state.activePatterns {
		if pattern.Type == patternName {
			// 根据比较操作符评估
			switch cond.Operator {
			case "eq":
				return pattern.Strength == expectedValue
			case "gt":
				return pattern.Strength > expectedValue
			case "lt":
				return pattern.Strength < expectedValue
			case "exists":
				return true
			}
		}
	}
	return false
}

// evaluatePropertyCondition 评估属性条件
func (re *RuleEngine) evaluatePropertyCondition(cond RuleCondition) bool {
	// 获取属性值
	if value, exists := re.field.GetPropertyValue(cond.Target); exists {
		expectedValue := cond.Value.(float64)

		// 根据操作符比较
		switch cond.Operator {
		case "eq":
			return value == expectedValue
		case "gt":
			return value > expectedValue
		case "lt":
			return value < expectedValue
		}
	}
	return false
}

// evaluateFieldCondition 评估场条件
func (re *RuleEngine) evaluateFieldCondition(cond RuleCondition) bool {
	// 直接使用GetPropertyValue获取字段值
	if value, exists := re.field.GetPropertyValue(cond.Target); exists {
		return evaluateNumericCondition(value, cond.Value.(float64), cond.Operator)
	}
	return false
}

// evaluateNumericCondition 评估数值条件
func evaluateNumericCondition(actual, expected float64, operator string) bool {
	switch operator {
	case "eq":
		return actual == expected
	case "gt":
		return actual > expected
	case "lt":
		return actual < expected
	default:
		return false
	}
}

// executeRules 执行规则
func (re *RuleEngine) executeRules(rules []*EmergenceRule) error {
	for _, rule := range rules {
		if err := re.executeRule(rule); err != nil {
			// 记录错误但继续执行其他规则
			re.recordRuleEvent(rule, "execution_failed", false, map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		// 更新规则状态
		rule.LastTriggered = time.Now()
		rule.TriggerCount++

		// 记录成功执行
		re.recordRuleEvent(rule, "executed", true, nil)
	}

	return nil
}

// executeRule 执行单个规则
func (re *RuleEngine) executeRule(rule *EmergenceRule) error {
	for _, action := range rule.Actions {
		if err := re.executeAction(action); err != nil {
			return err
		}
	}
	return nil
}

// executeAction 执行动作
func (re *RuleEngine) executeAction(action RuleAction) error {
	switch action.Type {
	case "create_property":
		return re.createProperty(action)
	case "modify_field":
		return re.modifyField(action)
	case "adjust_pattern":
		return re.adjustPattern(action)
	default:
		return model.WrapError(nil, model.ErrCodeValidation,
			fmt.Sprintf("unknown action type: %s", action.Type))
	}
}

// modifyField 修改场
func (re *RuleEngine) modifyField(action RuleAction) error {
	// 获取修改参数
	field, ok := action.Parameters["field"].(string)
	if !ok {
		return model.WrapError(nil, model.ErrCodeValidation, "missing field name")
	}

	value, ok := action.Parameters["value"].(float64)
	if !ok {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid field value")
	}

	// 更新字段值
	return re.field.SetPropertyValue(field, value)
}

// adjustPattern 调整模式
func (re *RuleEngine) adjustPattern(action RuleAction) error {
	// 获取模式参数
	patternID, ok := action.Parameters["pattern_id"].(string)
	if !ok {
		return model.WrapError(nil, model.ErrCodeValidation, "missing pattern id")
	}

	adjustment, ok := action.Parameters["adjustment"].(float64)
	if !ok {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid adjustment value")
	}

	// 调整模式
	pattern, exists := re.state.patterns[patternID]
	if !exists {
		return model.WrapError(nil, model.ErrCodeNotFound, "pattern not found")
	}

	pattern.Strength *= (1 + adjustment)
	return nil
}

// createProperty 创建属性
func (re *RuleEngine) createProperty(action RuleAction) error {
	// 获取属性参数
	name, ok := action.Parameters["name"].(string)
	if !ok {
		return model.WrapError(nil, model.ErrCodeValidation, "missing property name")
	}

	value, ok := action.Parameters["value"].(float64)
	if !ok {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid property value")
	}

	// 创建属性
	property := EmergentProperty{
		ID:      core.GenerateID(),
		Name:    name,
		Value:   value,
		Created: time.Now(),
		Updated: time.Now(),
	}

	// 保存属性
	re.state.properties[property.ID] = &property
	return nil
}

// learnNewRules 学习新规则
func (re *RuleEngine) learnNewRules() error {
	// 分析模式关联
	patterns := re.analyzePatternCorrelations()

	// 从关联中提取规则
	newRules := re.extractRulesFromPatterns(patterns)

	// 验证和添加新规则
	for _, rule := range newRules {
		if err := re.validateAndAddRule(rule); err != nil {
			continue
		}
	}

	return nil
}

// analyzePatternCorrelations 分析模式关联
func (re *RuleEngine) analyzePatternCorrelations() []PatternCorrelation {
	correlations := make([]PatternCorrelation, 0)

	// 遍历活跃模式对
	for id1, state1 := range re.state.patterns {
		for id2, state2 := range re.state.patterns {
			if id1 == id2 {
				continue
			}

			// 计算关联度
			correlation := re.calculatePatternCorrelation(state1, state2)
			if correlation.Strength > re.config.minConfidence {
				correlations = append(correlations, correlation)
			}
		}
	}

	return correlations
}

// calculatePatternCorrelation 计算模式关联度
func (re *RuleEngine) calculatePatternCorrelation(state1, state2 PatternState) PatternCorrelation {
	correlation := PatternCorrelation{
		Source: state1,
		Target: state2,
	}

	// 计算强度相关性
	strengthCorr := calculateStrengthCorrelation(state1.Pattern.Strength, state2.Pattern.Strength)

	// 计算时间相关性
	timeCorr := calculateTimeCorrelation(state1.LastUpdate, state2.LastUpdate)

	// 计算空间相关性
	spaceCorr := re.calculateSpatialCorrelation(state1.Pattern, state2.Pattern)

	// 组合关联度
	correlation.Strength = (strengthCorr*0.4 + timeCorr*0.3 + spaceCorr*0.3)

	// 计算影响效应
	correlation.Effect = calculateCorrelationEffect(state1, state2)

	// 设置阈值
	correlation.Threshold = state1.Pattern.Strength * 0.8 // 设置为源模式强度的80%

	return correlation
}

// calculateSpatialCorrelation 计算模式空间关联度
func (re *RuleEngine) calculateSpatialCorrelation(pattern1, pattern2 *EmergentPattern) float64 {
	if len(pattern1.Components) == 0 || len(pattern2.Components) == 0 {
		return 0
	}

	// 计算组件间的空间关联
	totalCorrelation := 0.0
	pairs := 0

	// 遍历组件对计算关联度
	for _, comp1 := range pattern1.Components {
		for _, comp2 := range pattern2.Components {
			// 基于组件类型计算空间关联
			correlation := re.calculateComponentSpatialCorrelation(comp1, comp2)
			totalCorrelation += correlation
			pairs++
		}
	}

	if pairs == 0 {
		return 0
	}

	return totalCorrelation / float64(pairs)
}

// calculateComponentSpatialCorrelation 计算组件空间关联度
func (re *RuleEngine) calculateComponentSpatialCorrelation(c1, c2 PatternComponent) float64 {
	// 基础关联度
	baseCorr := 0.5

	// 类型相关性
	if c1.Type == c2.Type {
		baseCorr = 0.8
	}

	// 根据组件类型调整关联度
	switch c1.Type {
	case "element":
		// 基于五行关系
		relation := model.GetWuXingRelation(c1.Role, c2.Role)
		return baseCorr * relation.Factor

	case "energy":
		// 基于能量差异
		energyDiff := math.Abs(c1.Weight - c2.Weight)
		return baseCorr / (1.0 + energyDiff)

	case "quantum":
		// 基于量子态关联
		if c1.Properties != nil && c2.Properties != nil {
			phase1 := c1.Properties["phase"]
			phase2 := c2.Properties["phase"]
			phaseDiff := math.Abs(phase1 - phase2)
			return baseCorr * (1.0 - phaseDiff/(2*math.Pi))
		}
	}

	return baseCorr
}

// 辅助函数
func calculateStrengthCorrelation(s1, s2 float64) float64 {
	return 1.0 - math.Abs(s1-s2)/math.Max(s1, s2)
}

func calculateTimeCorrelation(t1, t2 time.Time) float64 {
	dt := math.Abs(t1.Sub(t2).Seconds())
	return math.Exp(-dt / 3600.0) // 1小时特征时间
}

func calculateCorrelationEffect(source, target PatternState) float64 {
	// 基于历史状态计算影响效应
	effect := (target.Pattern.Strength - source.Pattern.Strength) / source.Pattern.Strength
	return math.Max(-0.5, math.Min(0.5, effect)) // 限制在[-0.5,0.5]范围内
}

// extractRulesFromPatterns 从模式关联中提取规则
func (re *RuleEngine) extractRulesFromPatterns(correlations []PatternCorrelation) []*EmergenceRule {
	rules := make([]*EmergenceRule, 0)

	for _, corr := range correlations {
		// 创建条件
		condition := RuleCondition{
			Type:     "pattern",
			Target:   corr.Source.Pattern.Type,
			Operator: "gt",
			Value:    corr.Threshold,
			Weight:   corr.Strength,
		}

		// 创建动作
		action := RuleAction{
			Type:      "modify_field",
			Target:    corr.Target.Pattern.Type,
			Operation: "adjust",
			Parameters: map[string]interface{}{
				"adjustment": corr.Effect,
			},
		}

		// 创建规则
		rule := &EmergenceRule{
			ID:         core.GenerateID(),
			Type:       "correlation",
			Conditions: []RuleCondition{condition},
			Actions:    []RuleAction{action},
			Priority:   int(corr.Strength * 100),
			Confidence: corr.Strength,
		}

		rules = append(rules, rule)
	}

	return rules
}

// validateAndAddRule 验证并添加规则
func (re *RuleEngine) validateAndAddRule(rule *EmergenceRule) error {
	// 1. 验证规则完整性
	if rule.ID == "" || len(rule.Conditions) == 0 || len(rule.Actions) == 0 {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid rule structure")
	}

	// 2. 检查规则是否已存在
	if _, exists := re.state.rules[rule.ID]; exists {
		return model.WrapError(nil, model.ErrCodeDuplicate, "rule already exists")
	}

	// 3. 检查规则数量限制
	if len(re.state.rules) >= re.config.maxRules {
		return model.WrapError(nil, model.ErrCodeLimit, "max rules limit reached")
	}

	// 4. 添加规则
	re.state.rules[rule.ID] = rule

	return nil
}

// recordRuleEvent 记录规则事件
func (re *RuleEngine) recordRuleEvent(
	rule *EmergenceRule,
	eventType string,
	success bool,
	details map[string]interface{}) {

	event := RuleEvent{
		Timestamp: time.Now(),
		RuleID:    rule.ID,
		Type:      eventType,
		Success:   success,
		Details:   details,
	}

	re.state.history = append(re.state.history, event)

	// 限制历史记录长度
	if len(re.state.history) > maxHistoryLength {
		re.state.history = re.state.history[1:]
	}
}

// 辅助函数

func (re *RuleEngine) getSortedRules() []*EmergenceRule {
	rules := make([]*EmergenceRule, 0, len(re.state.rules))
	for _, rule := range re.state.rules {
		rules = append(rules, rule)
	}

	// 按优先级排序
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority > rules[j].Priority
	})

	return rules
}
