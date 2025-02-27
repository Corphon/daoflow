//system/evolution/mutation/handler.go

package mutation

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/common"
	"github.com/Corphon/daoflow/system/evolution/pattern"
	"github.com/Corphon/daoflow/system/types"
)

// MutationHandler 突变处理器
type MutationHandler struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		responseThreshold float64 // 响应阈值
		maxRetries        int     // 最大重试次数
		stabilityTarget   float64 // 稳定性目标
		adaptiveResponse  bool    // 自适应响应
	}

	// 处理状态
	state struct {
		active     map[string]*MutationResponse // 活跃响应
		history    []ResponseEvent              // 响应历史
		strategies map[string]*ResponseStrategy // 响应策略
	}

	// 依赖项
	detector  *MutationDetector
	analyzer  common.PatternAnalyzer
	generator *pattern.PatternGenerator
}

// 确保实现接口
var _ common.MutationHandler = (*MutationHandler)(nil)

// MutationResponse 突变响应
type MutationResponse struct {
	ID         string            // 响应ID
	MutationID string            // 对应突变ID
	Strategy   *ResponseStrategy // 使用的策略
	Actions    []ResponseAction  // 响应动作
	Status     string            // 当前状态
	Progress   float64           // 进度
	StartTime  time.Time         // 开始时间
	LastUpdate time.Time         // 最后更新时间
	Retries    int               // 重试次数
}

// ResponseStrategy 响应策略
type ResponseStrategy struct {
	ID         string                 // 策略ID
	Type       string                 // 策略类型
	Conditions []ResponseCondition    // 触发条件
	Actions    []ActionTemplate       // 动作模板
	Priority   int                    // 优先级
	Success    float64                // 成功率
	Parameters map[string]interface{} // 策略参数
}

// ResponseCondition 响应条件
type ResponseCondition struct {
	Type     string      // 条件类型
	Target   string      // 目标对象
	Operator string      // 操作符
	Value    interface{} // 比较值
	Weight   float64     // 权重
}

// ActionTemplate 动作模板
type ActionTemplate struct {
	Type        string                 // 动作类型
	Parameters  map[string]interface{} // 参数模板
	Constraints []ActionConstraint     // 执行约束
	Timeout     time.Duration          // 超时时间
}

// ActionConstraint 动作执行约束
type ActionConstraint struct {
	Type      string        // 约束类型
	Target    string        // 约束目标
	Threshold interface{}   // 阈值
	Duration  time.Duration // 持续时间
}

// ResponseAction 响应动作
type ResponseAction struct {
	ID         string                 // 动作ID
	Type       string                 // 动作类型
	Parameters map[string]interface{} // 实际参数
	Status     string                 // 执行状态
	Result     interface{}            // 执行结果
	StartTime  time.Time              // 开始时间
	EndTime    time.Time              // 结束时间
}

// ResponseEvent 响应事件
type ResponseEvent struct {
	Timestamp  time.Time
	ResponseID string
	Type       string
	Status     string
	Details    map[string]interface{}
	Error      error
}

// ---------------------------------------------------------------------------------------------
// AdjustParameter 调整系统参数
func (mh *MutationHandler) AdjustParameter(target string, params map[string]interface{}) error {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	// 创建调整动作
	action := &ResponseAction{
		ID:   generateActionID(),
		Type: "parameter_adjustment",
		Parameters: map[string]interface{}{
			"target": target,
			"params": params,
		},
		StartTime: time.Now(),
	}

	// 执行调整
	if err := mh.executeSystemAction(action); err != nil {
		action.Status = "failed"
		action.EndTime = time.Now()
		return err
	}

	// 记录成功
	action.Status = "completed"
	action.EndTime = time.Now()

	return nil
}

// Optimize 执行系统优化
func (mh *MutationHandler) Optimize(params map[string]interface{}) error {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	// 创建优化动作
	action := &ResponseAction{
		ID:         generateActionID(),
		Type:       "system_optimization",
		Parameters: params,
		StartTime:  time.Now(),
	}

	// 执行优化
	if err := mh.executeSystemAction(action); err != nil {
		action.Status = "failed"
		action.EndTime = time.Now()
		return err
	}

	// 记录成功
	action.Status = "completed"
	action.EndTime = time.Now()

	return nil
}

// Transform 执行系统转换
func (mh *MutationHandler) Transform(params map[string]interface{}) error {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	// 创建转换动作
	action := &ResponseAction{
		ID:         generateActionID(),
		Type:       "system_transform",
		Parameters: params,
		StartTime:  time.Now(),
	}

	// 执行转换
	if err := mh.executeSystemAction(action); err != nil {
		action.Status = "failed"
		action.EndTime = time.Now()
		return err
	}

	// 记录成功
	action.Status = "completed"
	action.EndTime = time.Now()

	return nil
}

// executeSystemAction 执行系统动作
func (mh *MutationHandler) executeSystemAction(action *ResponseAction) error {
	// 记录开始执行
	mh.recordResponseEvent(&MutationResponse{
		ID:      core.GenerateID(),
		Actions: []ResponseAction{*action},
	}, "action_start", map[string]interface{}{
		"action_type": action.Type,
		"parameters":  action.Parameters,
	})

	// TODO: 实际的系统操作执行
	// 这里需要调用具体的系统操作实现

	// 更新结束时间
	action.EndTime = time.Now()

	// 记录完成
	mh.recordResponseEvent(&MutationResponse{
		ID:      core.GenerateID(),
		Actions: []ResponseAction{*action},
	}, "action_complete", map[string]interface{}{
		"duration": time.Since(action.StartTime).Milliseconds(),
	})

	return nil
}

// 辅助函数
func generateActionID() string {
	return fmt.Sprintf("act_%d", time.Now().UnixNano())
}

// GetCurrentState 获取当前系统状态
func (mh *MutationHandler) GetCurrentState() (*model.SystemState, error) {
	mh.mu.RLock()
	defer mh.mu.RUnlock()

	// 计算当前系统的各项指标
	metrics := mh.GetHandlingMetrics()
	stability := mh.calculateSystemStability()

	currentState := &model.SystemState{
		// 基础状态指标
		Energy:     float64(len(mh.state.active)) / float64(maxHistoryLength), // 活跃度作为能量
		Entropy:    1.0 - stability,                                           // 系统混乱度为稳定性的补值
		Harmony:    metrics.SuccessRate,                                       // 成功率作为和谐度
		Balance:    stability,                                                 // 稳定性作为平衡度
		Phase:      determineSystemPhase(stability, metrics.SuccessRate),      // 根据状态确定相位
		Timestamp:  time.Now(),
		Properties: make(map[string]interface{}),
	}

	// 扩展属性
	currentState.Properties["mutation_count"] = len(mh.state.active)
	currentState.Properties["strategy_count"] = len(mh.state.strategies)
	currentState.Properties["response_rate"] = mh.calculateResponseRate()
	currentState.Properties["success_rate"] = metrics.SuccessRate
	currentState.Properties["stability"] = stability
	currentState.Properties["average_latency"] = metrics.AverageLatency.Milliseconds()
	currentState.Properties["total_handled"] = metrics.TotalHandled
	currentState.Properties["last_handled"] = metrics.LastHandled

	return currentState, nil
}

// determineSystemPhase 根据系统状态确定相位
func determineSystemPhase(stability float64, successRate float64) model.Phase {
	// 根据稳定性和成功率综合判断系统相位
	if stability >= 0.8 && successRate >= 0.8 {
		return model.PhaseYang // 系统处于良性循环
	} else if stability <= 0.2 || successRate <= 0.2 {
		return model.PhaseYin // 系统处于恶性循环
	}
	return model.PhaseNone // 系统处于中性状态
}

// 辅助方法
func (mh *MutationHandler) calculateResponseRate() float64 {
	if len(mh.state.history) == 0 {
		return 0
	}

	totalResponses := 0
	for _, event := range mh.state.history {
		if event.Type == "response_executed" {
			totalResponses++
		}
	}

	return float64(totalResponses) / float64(len(mh.state.history))
}

func (mh *MutationHandler) calculateSuccessRate() float64 {
	if len(mh.state.history) == 0 {
		return 0
	}

	successCount := 0
	for _, event := range mh.state.history {
		// 判断是否为成功状态
		if event.Status == "completed" && event.Error == nil {
			successCount++
		}
	}

	return float64(successCount) / float64(len(mh.state.history))
}

func (mh *MutationHandler) calculateSystemStability() float64 {
	// 基于活跃突变和响应计算系统稳定性
	if len(mh.state.active) == 0 {
		return 1.0 // 完全稳定
	}

	totalStability := 0.0
	for _, response := range mh.state.active {
		if response.Status == "completed" {
			totalStability += response.Progress
		}
	}

	return totalStability / float64(len(mh.state.active))
}

// NewMutationHandler 创建新的突变处理器
func NewMutationHandler(detector *MutationDetector, config *types.MutationConfig) (*MutationHandler, error) {
	if detector == nil {
		return nil, fmt.Errorf("nil mutation detector")
	}
	if config == nil {
		return nil, fmt.Errorf("nil mutation config")
	}

	mh := &MutationHandler{
		detector: detector,
	}

	// 初始化配置
	mh.config.responseThreshold = config.Handler.ResponseThreshold
	mh.config.maxRetries = config.Handler.MaxRetries
	mh.config.stabilityTarget = config.Handler.StabilityTarget
	mh.config.adaptiveResponse = config.Handler.AdaptiveResponse

	// 初始化状态
	mh.state.active = make(map[string]*MutationResponse)
	mh.state.history = make([]ResponseEvent, 0)
	mh.state.strategies = make(map[string]*ResponseStrategy)

	return mh, nil
}

// Handle 处理突变
func (mh *MutationHandler) Handle() error {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	// 获取待处理的突变
	mutations, err := mh.detector.GetActiveMutations()
	if err != nil {
		return err
	}

	// 为每个突变选择策略
	responses := mh.selectStrategies(mutations)

	// 执行响应
	if err := mh.executeResponses(responses); err != nil {
		return err
	}

	// 监控响应进度
	mh.monitorResponses()

	// 更新响应状态
	mh.updateResponseStatus()

	return nil
}

// RegisterStrategy 注册响应策略
func (mh *MutationHandler) RegisterStrategy(strategy *ResponseStrategy) error {
	if strategy == nil {
		return model.WrapError(nil, model.ErrCodeValidation, "nil strategy")
	}

	mh.mu.Lock()
	defer mh.mu.Unlock()

	// 验证策略
	if err := mh.validateStrategy(strategy); err != nil {
		return err
	}

	// 存储策略
	mh.state.strategies[strategy.ID] = strategy

	return nil
}

// selectStrategies 选择响应策略
func (mh *MutationHandler) selectStrategies(
	mutations []*Mutation) []*MutationResponse {

	responses := make([]*MutationResponse, 0)

	for _, mutation := range mutations {
		// 检查是否已有响应
		if mh.hasActiveResponse(mutation.ID) {
			continue
		}

		// 选择最佳策略
		strategy := mh.selectBestStrategy(mutation)
		if strategy == nil {
			continue
		}

		// 创建响应
		response := mh.createResponse(mutation, strategy)
		responses = append(responses, response)
	}

	return responses
}

// executeResponses 执行响应
func (mh *MutationHandler) executeResponses(responses []*MutationResponse) error {
	for _, response := range responses {
		// 准备执行环境
		context := mh.prepareExecutionContext(response)

		// 执行响应动作
		if err := mh.executeResponseActions(response, context); err != nil {
			// 记录错误但继续执行其他响应
			mh.recordResponseEvent(response, "execution_error", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		// 更新响应状态
		response.LastUpdate = time.Now()
		mh.state.active[response.ID] = response
	}

	return nil
}

// monitorResponses 监控响应进度
func (mh *MutationHandler) monitorResponses() {
	for id, response := range mh.state.active {
		// 检查响应是否超时
		if mh.isResponseTimeout(response) {
			mh.handleResponseTimeout(response)
			continue
		}

		// 更新进度
		progress := mh.calculateResponseProgress(response)
		response.Progress = progress

		// 检查完成状态
		if mh.isResponseComplete(response) {
			mh.finalizeResponse(response)
			delete(mh.state.active, id)
		}
	}
}

// updateResponseStatus 更新响应状态
func (mh *MutationHandler) updateResponseStatus() {
	currentTime := time.Now()

	for _, response := range mh.state.active {
		// 检查动作状态
		allComplete := true
		for _, action := range response.Actions {
			if action.Status != "completed" {
				allComplete = false
				break
			}
		}

		// 更新响应状态
		if allComplete {
			response.Status = "completed"
		} else if response.Retries >= mh.config.maxRetries {
			response.Status = "failed"
		}

		// 记录状态变更
		mh.recordResponseEvent(response, "status_update", map[string]interface{}{
			"status":   response.Status,
			"progress": response.Progress,
		})

		response.LastUpdate = currentTime
	}
}

// HandleMutation implements common.MutationHandler
func (mh *MutationHandler) HandleMutation(mutation common.Mutation) error {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	// 为单个突变选择策略
	strategy := mh.selectBestStrategy(mutation)
	if strategy == nil {
		return fmt.Errorf("no suitable strategy found for mutation %v", mutation.GetID())
	}

	// 创建响应
	response := mh.createResponse(mutation, strategy)

	// 执行响应
	context := mh.prepareExecutionContext(response)
	if err := mh.executeResponseActions(response, context); err != nil {
		return err
	}

	// 更新状态
	response.LastUpdate = time.Now()
	mh.state.active[response.ID] = response

	return nil
}

// GetHandlingMetrics implements common.MutationHandler
func (mh *MutationHandler) GetHandlingMetrics() common.HandlingMetrics {
	mh.mu.RLock()
	defer mh.mu.RUnlock()

	return common.HandlingMetrics{
		TotalHandled:   len(mh.state.history),
		SuccessRate:    mh.calculateSuccessRate(),
		AverageLatency: mh.calculateAverageLatency(),
		LastHandled:    mh.getLastHandledTime(),
	}
}

// 新增辅助方法
func (mh *MutationHandler) calculateAverageLatency() time.Duration {
	if len(mh.state.history) == 0 {
		return 0
	}

	var totalLatency time.Duration
	for _, event := range mh.state.history {
		if details, ok := event.Details["duration"]; ok {
			if duration, ok := details.(int64); ok {
				totalLatency += time.Duration(duration) * time.Millisecond
			}
		}
	}

	return totalLatency / time.Duration(len(mh.state.history))
}

func (mh *MutationHandler) getLastHandledTime() time.Time {
	if len(mh.state.history) == 0 {
		return time.Time{}
	}
	return mh.state.history[len(mh.state.history)-1].Timestamp
}

// 辅助函数

func (mh *MutationHandler) validateStrategy(strategy *ResponseStrategy) error {
	if strategy.ID == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty strategy ID")
	}

	// 验证条件
	for _, condition := range strategy.Conditions {
		if err := mh.validateCondition(condition); err != nil {
			return err
		}
	}

	// 验证动作模板
	for _, template := range strategy.Actions {
		if err := mh.validateActionTemplate(template); err != nil {
			return err
		}
	}

	return nil
}

func (mh *MutationHandler) recordResponseEvent(
	response *MutationResponse,
	eventType string,
	details map[string]interface{}) {

	event := ResponseEvent{
		Timestamp:  time.Now(),
		ResponseID: response.ID,
		Type:       eventType,
		Status:     response.Status,
		Details:    details,
		Error:      nil, // 默认无错误
	}

	// 如果 details 中包含错误信息，则设置错误字段
	if errVal, ok := details["error"]; ok {
		if errStr, ok := errVal.(string); ok {
			event.Error = fmt.Errorf("response error: %s", errStr)
		}
	}

	mh.state.history = append(mh.state.history, event)

	// 限制历史记录长度
	if len(mh.state.history) > maxHistoryLength {
		mh.state.history = mh.state.history[1:]
	}
}

const (
	maxHistoryLength = 1000
)

// selectBestStrategy 选择最佳策略
func (mh *MutationHandler) selectBestStrategy(mutation interface{}) *ResponseStrategy {
	var mutID string
	switch m := mutation.(type) {
	case *Mutation:
		mutID = m.ID
	case common.Mutation:
		mutID = m.GetID()
	default:
		return nil
	}

	var bestStrategy *ResponseStrategy
	var highestScore float64

	for _, strategy := range mh.state.strategies {
		score := mh.evaluateStrategyFit(strategy, mutID)
		if score > highestScore {
			highestScore = score
			bestStrategy = strategy
		}
	}

	return bestStrategy
}

// resolveActionParameters 参数解析
func (mh *MutationHandler) resolveActionParameters(
	templateParams map[string]interface{},
	context map[string]interface{}) map[string]interface{} {

	resolved := make(map[string]interface{})
	for k, v := range templateParams {
		if strVal, ok := v.(string); ok && strVal[0] == '$' {
			if contextVal, exists := context[strVal[1:]]; exists {
				resolved[k] = contextVal
				continue
			}
		}
		resolved[k] = v
	}
	return resolved
}

// executeAction 执行动作
func (mh *MutationHandler) executeAction(action *ResponseAction) error {
	// TODO: 实现具体的动作执行逻辑
	action.Status = "completed"
	action.EndTime = time.Now()
	return nil
}

// hasActiveResponse 检查是否有活跃响应
func (mh *MutationHandler) hasActiveResponse(mutationID string) bool {
	for _, resp := range mh.state.active {
		if resp.MutationID == mutationID {
			return true
		}
	}
	return false
}

// createResponse 创建响应
func (mh *MutationHandler) createResponse(mutation interface{}, strategy *ResponseStrategy) *MutationResponse {
	var mutID string
	switch m := mutation.(type) {
	case *Mutation:
		mutID = m.ID
	case common.Mutation:
		mutID = m.GetID()
	default:
		return nil
	}

	response := &MutationResponse{
		ID:         fmt.Sprintf("resp_%d", time.Now().UnixNano()),
		MutationID: mutID,
		Strategy:   strategy,
		Status:     "pending",
		StartTime:  time.Now(),
		LastUpdate: time.Now(),
		Actions:    make([]ResponseAction, 0),
	}

	return response
}

// prepareExecutionContext 准备执行上下文
func (mh *MutationHandler) prepareExecutionContext(response *MutationResponse) map[string]interface{} {
	return map[string]interface{}{
		"response_id": response.ID,
		"mutation_id": response.MutationID,
		"start_time":  response.StartTime,
	}
}

// executeResponseActions 执行响应动作
func (mh *MutationHandler) executeResponseActions(response *MutationResponse, context map[string]interface{}) error {
	for _, actionTemplate := range response.Strategy.Actions {
		action := ResponseAction{
			ID:         fmt.Sprintf("act_%d", time.Now().UnixNano()),
			Type:       actionTemplate.Type,
			Parameters: mh.resolveActionParameters(actionTemplate.Parameters, context),
			StartTime:  time.Now(),
			Status:     "executing",
		}

		if err := mh.executeAction(&action); err != nil {
			return err
		}

		response.Actions = append(response.Actions, action)
	}
	return nil
}

// isResponseTimeout 检查响应是否超时
func (mh *MutationHandler) isResponseTimeout(response *MutationResponse) bool {
	if len(response.Strategy.Actions) == 0 {
		return false
	}
	maxTimeout := response.Strategy.Actions[0].Timeout
	for _, action := range response.Strategy.Actions[1:] {
		if action.Timeout > maxTimeout {
			maxTimeout = action.Timeout
		}
	}
	return time.Since(response.StartTime) > maxTimeout
}

// handleResponseTimeout 处理响应超时
func (mh *MutationHandler) handleResponseTimeout(response *MutationResponse) {
	response.Status = "timeout"
}

// isResponseComplete 检查响应是否完成
func (mh *MutationHandler) isResponseComplete(response *MutationResponse) bool {
	if len(response.Actions) == 0 {
		return false
	}

	for _, action := range response.Actions {
		if action.Status != "completed" {
			return false
		}
	}
	return true
}

// finalizeResponse 完成响应处理
func (mh *MutationHandler) finalizeResponse(response *MutationResponse) {
	response.Status = "completed"
	response.LastUpdate = time.Now()
	mh.recordResponseEvent(response, "completed", map[string]interface{}{
		"duration": time.Since(response.StartTime).String(),
		"actions":  len(response.Actions),
	})
}

// validateCondition 验证响应条件
func (mh *MutationHandler) validateCondition(condition ResponseCondition) error {
	if condition.Type == "" {
		return fmt.Errorf("empty condition type")
	}
	if condition.Target == "" {
		return fmt.Errorf("empty condition target")
	}
	if condition.Weight < 0 || condition.Weight > 1 {
		return fmt.Errorf("invalid condition weight: %v", condition.Weight)
	}
	return nil
}

// validateActionTemplate 验证动作模板
func (mh *MutationHandler) validateActionTemplate(template ActionTemplate) error {
	if template.Type == "" {
		return fmt.Errorf("empty action type")
	}
	if template.Timeout <= 0 {
		return fmt.Errorf("invalid timeout: %v", template.Timeout)
	}
	return nil
}

// 新增辅助方法
func (mh *MutationHandler) evaluateStrategyFit(strategy *ResponseStrategy, mutationID string) float64 {
	if strategy == nil {
		return 0
	}

	// 1. 基础适应度分数 (30%)
	baseScore := float64(strategy.Priority) / 100.0

	// 2. 历史成功率 (25%)
	successScore := strategy.Success

	// 3. 条件匹配度 (25%)
	var conditionScore float64
	totalWeight := 0.0
	for _, condition := range strategy.Conditions {
		if mh.evaluateConditionMatch(condition, mutationID) {
			conditionScore += condition.Weight
		}
		totalWeight += condition.Weight
	}
	if totalWeight > 0 {
		conditionScore = conditionScore / totalWeight
	}

	// 4. 资源适应度 (20%)
	resourceScore := mh.evaluateResourceFit(strategy, mutationID)

	// 计算加权总分
	totalScore := (baseScore * 0.30) +
		(successScore * 0.25) +
		(conditionScore * 0.25) +
		(resourceScore * 0.20)

	// 归一化到 [0,1] 范围
	return math.Max(0, math.Min(1, totalScore))
}

// evaluateConditionMatch 评估条件匹配度
func (mh *MutationHandler) evaluateConditionMatch(condition ResponseCondition, mutationID string) bool {
	// 获取突变相关的上下文值
	contextValue, exists := mh.getMutationContextValue(mutationID, condition.Target)
	if !exists {
		return false
	}

	// 根据操作符进行比较
	switch condition.Operator {
	case ">":
		return compareValues(contextValue, condition.Value) > 0
	case "<":
		return compareValues(contextValue, condition.Value) < 0
	case "==":
		return compareValues(contextValue, condition.Value) == 0
	case "!=":
		return compareValues(contextValue, condition.Value) != 0
	default:
		return false
	}
}

// evaluateResourceFit 评估资源适应度
func (mh *MutationHandler) evaluateResourceFit(strategy *ResponseStrategy, mutationID string) float64 {
	// 检查系统当前资源状态
	activeResponses := len(mh.state.active)
	if activeResponses >= maxHistoryLength {
		return 0 // 系统负载已满
	}

	// 计算当前资源使用率
	resourceUsage := float64(activeResponses) / float64(maxHistoryLength)

	// 获取突变相关的资源需求
	resourceDemand := 1.0
	if response, exists := mh.state.active[mutationID]; exists {
		// 考虑已有响应的资源占用
		resourceDemand = float64(len(response.Actions)) / float64(maxHistoryLength)
	}

	// 根据策略动作数量和系统负载计算适应度
	actionCount := len(strategy.Actions)
	if actionCount == 0 {
		return 1 // 无动作策略默认资源适应度为1
	}

	// 考虑当前资源使用率、动作数量和突变资源需求
	adjustedUsage := (resourceUsage + resourceDemand) / 2
	return math.Max(0, 1-adjustedUsage) / float64(actionCount)
}

// getMutationContextValue 获取突变上下文值
func (mh *MutationHandler) getMutationContextValue(mutationID string, target string) (interface{}, bool) {
	// 遍历活跃突变
	for _, response := range mh.state.active {
		if response.MutationID == mutationID {
			if value, exists := response.Strategy.Parameters[target]; exists {
				return value, true
			}
		}
	}
	return nil, false
}

// compareValues 比较两个值
func compareValues(a, b interface{}) int {
	// 转换为可比较的类型
	switch v1 := a.(type) {
	case float64:
		if v2, ok := b.(float64); ok {
			if v1 < v2 {
				return -1
			} else if v1 > v2 {
				return 1
			}
			return 0
		}
	case int:
		if v2, ok := b.(int); ok {
			if v1 < v2 {
				return -1
			} else if v1 > v2 {
				return 1
			}
			return 0
		}
	}
	return 0
}

// calculateResponseProgress 计算响应进度
func (mh *MutationHandler) calculateResponseProgress(response *MutationResponse) float64 {
	if len(response.Actions) == 0 {
		return 0
	}

	completed := 0
	for _, action := range response.Actions {
		if action.Status == "completed" {
			completed++
		}
	}

	return float64(completed) / float64(len(response.Actions))
}
