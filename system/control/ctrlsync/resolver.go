//system/control/ctrlsync/resolver.go

package sync

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/types"
)

// Resolver 解决器
type Resolver struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		resolutionTimeout time.Duration // 解决超时
		maxAttempts       int           // 最大尝试次数
		minConfidence     float64       // 最小置信度
		autoResolve       bool          // 自动解决
	}

	// 解决状态
	state struct {
		conflicts   map[string]*Conflict   // 冲突列表
		strategies  map[string]*Strategy   // 策略列表
		resolutions map[string]*Resolution // 解决方案
		metrics     ResolutionMetrics      // 解决指标
	}
}

// Conflict 冲突信息
type Conflict struct {
	ID         string     // 冲突ID
	Type       string     // 冲突类型
	Status     string     // 冲突状态
	Resolution string     // 解决方案
	Priority   int        // 优先级
	Parties    []Party    // 冲突方
	Resources  []Resource // 相关资源
	Created    time.Time  // 创建时间
	Updated    time.Time  // 更新时间
}

// Party 冲突方
type Party struct {
	ID           string        // 参与方ID
	Type         string        // 参与方类型
	Role         string        // 参与角色
	Position     interface{}   // 立场信息
	Requirements []Requirement // 需求列表
}

// Resource 相关资源
type Resource struct {
	ID           string       // 资源ID
	Type         string       // 资源类型
	State        string       // 资源状态
	Constraints  []Constraint // 资源约束
	Dependencies []string     // 依赖资源
}

// Requirement 需求信息
type Requirement struct {
	ID          string       // 需求ID
	Type        string       // 需求类型
	Priority    int          // 优先级
	Constraints []Constraint // 需求约束
	Flexibility float64      // 灵活度
}

// Strategy 解决策略
type Strategy struct {
	ID         string      // 策略ID
	Type       string      // 策略类型
	Priority   int         // 优先级
	Conditions []Condition // 应用条件
	Actions    []Action    // 策略动作
	Success    float64     // 成功率
}

// Condition 应用条件
type Condition struct {
	Type     string      // 条件类型
	Target   string      // 目标对象
	Value    interface{} // 条件值
	Operator string      // 操作符
	Weight   float64     // 权重
}

// Action 策略动作
type Action struct {
	Type       string                 // 动作类型
	Target     string                 // 目标对象
	Operation  string                 // 操作类型
	Parameters map[string]interface{} // 操作参数
}

// Resolution 解决方案
type Resolution struct {
	ID         string           // 方案ID
	ConflictID string           // 冲突ID
	Type       string           // 方案类型
	Status     string           // 方案状态
	Steps      []ResolutionStep // 解决步骤
	Results    ResolutionResult // 解决结果
	Created    time.Time        // 创建时间
}

// ResolutionStep 解决步骤
type ResolutionStep struct {
	ID        string    // 步骤ID
	Type      string    // 步骤类型
	Action    Action    // 执行动作
	Status    string    // 步骤状态
	StartTime time.Time // 开始时间
	EndTime   time.Time // 结束时间
}

// ResolutionResult 解决结果
type ResolutionResult struct {
	Success    bool               // 是否成功
	Confidence float64            // 置信度
	Impact     map[string]float64 // 影响评估
	Feedback   []string           // 反馈信息
}

// ResolutionMetrics 解决指标
type ResolutionMetrics struct {
	TotalConflicts int                 // 总冲突数
	ResolvedCount  int                 // 已解决数
	SuccessRate    float64             // 成功率
	AverageTime    time.Duration       // 平均耗时
	History        []types.MetricPoint // 历史指标
}

// NewResolver 创建新的解决器
func NewResolver() *Resolver {
	r := &Resolver{}

	// 初始化配置
	r.config.resolutionTimeout = 30 * time.Second
	r.config.maxAttempts = 3
	r.config.minConfidence = 0.75
	r.config.autoResolve = true

	// 初始化状态
	r.state.conflicts = make(map[string]*Conflict)
	r.state.strategies = make(map[string]*Strategy)
	r.state.resolutions = make(map[string]*Resolution)
	r.state.metrics = ResolutionMetrics{
		History: make([]types.MetricPoint, 0),
	}

	return r
}

// RegisterConflict 注册冲突
func (r *Resolver) RegisterConflict(conflict *Conflict) error {
	if conflict == nil {
		return model.WrapError(nil, model.ErrCodeValidation, "nil conflict")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// 验证冲突
	if err := r.validateConflict(conflict); err != nil {
		return err
	}

	// 存储冲突
	r.state.conflicts[conflict.ID] = conflict

	// 如果开启自动解决，尝试解决冲突
	if r.config.autoResolve {
		go r.resolveConflict(conflict)
	}

	return nil
}

// resolveConflict 自动解决冲突
func (r *Resolver) resolveConflict(conflict *Conflict) {
	// 创建解决方案
	resolution := &Resolution{
		ID:         generateResolutionID(),
		ConflictID: conflict.ID,
		Type:       determineResolutionType(conflict),
		Status:     "initiated",
		Steps:      make([]ResolutionStep, 0),
		Created:    time.Now(),
	}

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), r.config.resolutionTimeout)
	defer cancel()

	// 尝试解决
	attempts := 0
	for attempts < r.config.maxAttempts {
		select {
		case <-ctx.Done():
			r.handleResolutionFailure(resolution, "timeout")
			return
		default:
			// 执行解决方案
			if err := r.executeResolution(resolution); err != nil {
				attempts++
				if attempts >= r.config.maxAttempts {
					r.handleResolutionFailure(resolution, "max_attempts")
					return
				}
				time.Sleep(time.Second * time.Duration(attempts)) // 退避重试
				continue
			}

			// 检查置信度
			if resolution.Results.Confidence < r.config.minConfidence {
				attempts++
				if attempts >= r.config.maxAttempts {
					r.handleResolutionFailure(resolution, "low_confidence")
					return
				}
				time.Sleep(time.Second * time.Duration(attempts)) // 退避重试
				continue
			}

			// 解决成功
			r.handleResolutionSuccess(resolution)
			return
		}
	}
}

// determineResolutionType 确定解决方案类型
func determineResolutionType(conflict *Conflict) string {
	// 基于冲突类型和特征确定解决方案类型
	switch conflict.Type {
	case "resource_conflict":
		if len(conflict.Resources) > 0 {
			return "resource_allocation"
		}
		return "resource_negotiation"

	case "state_conflict":
		return "state_reconciliation"

	case "requirement_conflict":
		// 检查是否所有参与方都有需求
		hasRequirements := true
		for _, party := range conflict.Parties {
			if len(party.Requirements) == 0 {
				hasRequirements = false
				break
			}
		}
		if hasRequirements {
			return "requirement_negotiation"
		}
		return "mediation"

	default:
		return "general_resolution"
	}
}

// executeResolution 执行解决方案
func (r *Resolver) executeResolution(resolution *Resolution) error {
	// 获取冲突信息
	conflict, exists := r.state.conflicts[resolution.ConflictID]
	if !exists {
		return model.WrapError(nil, model.ErrCodeNotFound, "conflict not found")
	}

	// 选择合适的策略
	var selectedStrategy *Strategy
	highestPriority := -1

	for _, strategy := range r.state.strategies {
		if !r.isStrategyApplicable(strategy, conflict) {
			continue
		}
		if strategy.Priority > highestPriority {
			selectedStrategy = strategy
			highestPriority = strategy.Priority
		}
	}

	if selectedStrategy == nil {
		return model.WrapError(nil, model.ErrCodeNotFound, "no applicable strategy")
	}

	// 执行策略动作
	for _, action := range selectedStrategy.Actions {
		step := ResolutionStep{
			ID:        fmt.Sprintf("step_%d", len(resolution.Steps)),
			Type:      action.Type,
			Action:    action,
			Status:    "pending",
			StartTime: time.Now(),
		}

		// 执行动作
		if err := r.executeAction(&step, conflict); err != nil {
			step.Status = "failed"
			resolution.Steps = append(resolution.Steps, step)
			return err
		}

		step.Status = "completed"
		step.EndTime = time.Now()
		resolution.Steps = append(resolution.Steps, step)
	}

	// 更新解决结果
	resolution.Results = ResolutionResult{
		Success:    true,
		Confidence: r.calculateConfidence(resolution),
		Impact:     r.assessImpact(resolution),
		Feedback:   []string{"Resolution completed successfully"},
	}

	return nil
}

// 内部辅助方法

// isStrategyApplicable 检查策略是否适用
func (r *Resolver) isStrategyApplicable(strategy *Strategy, conflict *Conflict) bool {
	for _, condition := range strategy.Conditions {
		if !r.evaluateCondition(condition, conflict) {
			return false
		}
	}
	return true
}

// evaluateCondition 评估条件
func (r *Resolver) evaluateCondition(condition Condition, conflict *Conflict) bool {
	// 根据条件类型进行评估
	switch condition.Type {
	case "resource":
		// 检查资源相关条件
		return r.evaluateResourceCondition(condition, conflict)

	case "state":
		// 检查状态相关条件
		return r.evaluateStateCondition(condition, conflict)

	case "requirement":
		// 检查需求相关条件
		return r.evaluateRequirementCondition(condition, conflict)

	case "party":
		// 检查参与方相关条件
		return r.evaluatePartyCondition(condition, conflict)

	default:
		return false
	}
}

// executeAction 执行动作
func (r *Resolver) executeAction(step *ResolutionStep, conflict *Conflict) error {
	// 执行前检查
	if step == nil || conflict == nil {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid action parameters")
	}

	// 根据动作类型执行操作
	switch step.Action.Type {
	case "resource_allocation":
		return r.executeResourceAllocation(step, conflict)

	case "state_update":
		return r.executeStateUpdate(step, conflict)

	case "requirement_adjustment":
		return r.executeRequirementAdjustment(step, conflict)

	default:
		return model.WrapError(nil, model.ErrCodeValidation,
			fmt.Sprintf("unknown action type: %s", step.Action.Type))
	}
}

// 辅助方法

func (r *Resolver) evaluateResourceCondition(condition Condition, conflict *Conflict) bool {
	// 检查资源是否存在
	if len(conflict.Resources) == 0 {
		return false
	}

	// 根据操作符比较值
	for _, resource := range conflict.Resources {
		switch condition.Operator {
		case "exists":
			if resource.ID == condition.Target {
				return true
			}
		case "state":
			if resource.ID == condition.Target && resource.State == condition.Value.(string) {
				return true
			}
		case "available":
			if resource.ID == condition.Target && resource.State == "available" {
				return true
			}
		}
	}
	return false
}

func (r *Resolver) evaluateStateCondition(condition Condition, conflict *Conflict) bool {
	switch condition.Operator {
	case "eq":
		return conflict.Status == condition.Value.(string)
	case "ne":
		return conflict.Status != condition.Value.(string)
	case "priority":
		return conflict.Priority >= int(condition.Value.(float64))
	}
	return false
}

func (r *Resolver) evaluateRequirementCondition(condition Condition, conflict *Conflict) bool {
	for _, party := range conflict.Parties {
		for _, req := range party.Requirements {
			if req.Type == condition.Target {
				switch condition.Operator {
				case "flexibility":
					return req.Flexibility >= condition.Value.(float64)
				case "priority":
					return req.Priority >= int(condition.Value.(float64))
				}
			}
		}
	}
	return false
}

func (r *Resolver) evaluatePartyCondition(condition Condition, conflict *Conflict) bool {
	for _, party := range conflict.Parties {
		if party.ID == condition.Target {
			switch condition.Operator {
			case "role":
				return party.Role == condition.Value.(string)
			case "type":
				return party.Type == condition.Value.(string)
			}
		}
	}
	return false
}

func (r *Resolver) executeResourceAllocation(step *ResolutionStep, conflict *Conflict) error {
	params := step.Action.Parameters
	resourceID := params["resource_id"].(string)
	targetState := params["target_state"].(string)

	// 更新资源状态
	for i, resource := range conflict.Resources {
		if resource.ID == resourceID {
			conflict.Resources[i].State = targetState
			return nil
		}
	}

	return model.WrapError(nil, model.ErrCodeNotFound, "resource not found")
}

func (r *Resolver) executeStateUpdate(step *ResolutionStep, conflict *Conflict) error {
	if newState, ok := step.Action.Parameters["new_state"].(string); ok {
		conflict.Status = newState
		conflict.Updated = time.Now()
		return nil
	}
	return model.WrapError(nil, model.ErrCodeValidation, "invalid state parameter")
}

func (r *Resolver) executeRequirementAdjustment(step *ResolutionStep, conflict *Conflict) error {
	params := step.Action.Parameters
	partyID := params["party_id"].(string)
	reqType := params["requirement_type"].(string)
	flexibility := params["flexibility"].(float64)

	// 更新需求灵活度
	for _, party := range conflict.Parties {
		if party.ID == partyID {
			for i, req := range party.Requirements {
				if req.Type == reqType {
					party.Requirements[i].Flexibility = flexibility
					return nil
				}
			}
		}
	}

	return model.WrapError(nil, model.ErrCodeNotFound, "requirement not found")
}

// calculateConfidence 计算置信度
func (r *Resolver) calculateConfidence(resolution *Resolution) float64 {
	successSteps := 0
	for _, step := range resolution.Steps {
		if step.Status == "completed" {
			successSteps++
		}
	}
	if len(resolution.Steps) == 0 {
		return 0
	}
	return float64(successSteps) / float64(len(resolution.Steps))
}

// assessImpact 评估影响
func (r *Resolver) assessImpact(resolution *Resolution) map[string]float64 {
	return map[string]float64{
		"success_rate": r.calculateConfidence(resolution),
		"completion":   1.0,
	}
}

// handleResolutionSuccess 处理解决成功
func (r *Resolver) handleResolutionSuccess(resolution *Resolution) {
	r.mu.Lock()
	defer r.mu.Unlock()

	resolution.Status = "completed"
	r.state.resolutions[resolution.ID] = resolution
	r.state.metrics.ResolvedCount++
	r.state.metrics.SuccessRate = float64(r.state.metrics.ResolvedCount) / float64(r.state.metrics.TotalConflicts)

	// 更新指标
	r.updateMetrics()
}

// handleResolutionFailure 处理解决失败
func (r *Resolver) handleResolutionFailure(resolution *Resolution, reason string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	resolution.Status = "failed"
	resolution.Results.Feedback = append(resolution.Results.Feedback,
		fmt.Sprintf("Resolution failed: %s", reason))
	r.state.resolutions[resolution.ID] = resolution

	// 更新指标
	r.updateMetrics()
}

// RegisterStrategy 注册策略
func (r *Resolver) RegisterStrategy(strategy *Strategy) error {
	if strategy == nil {
		return model.WrapError(nil, model.ErrCodeValidation, "nil strategy")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// 验证策略
	if err := r.validateStrategy(strategy); err != nil {
		return err
	}

	// 存储策略
	r.state.strategies[strategy.ID] = strategy

	return nil
}

// Resolve 解决冲突
func (r *Resolver) Resolve(conflictID string) (*Resolution, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 获取冲突
	conflict, exists := r.state.conflicts[conflictID]
	if !exists {
		return nil, model.WrapError(nil, model.ErrCodeNotFound, "conflict not found")
	}

	// 创建解决方案
	resolution := &Resolution{
		ID:         generateResolutionID(),
		ConflictID: conflictID,
		Type:       determineResolutionType(conflict),
		Status:     "initiated",
		Steps:      make([]ResolutionStep, 0),
		Created:    time.Now(),
	}

	// 执行解决方案
	if err := r.executeResolution(resolution); err != nil {
		return nil, err
	}

	// 存储解决方案
	r.state.resolutions[resolution.ID] = resolution

	// 更新指标
	r.updateMetrics()

	return resolution, nil
}

// 辅助函数

func (r *Resolver) validateConflict(conflict *Conflict) error {
	if conflict.ID == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty conflict ID")
	}

	if len(conflict.Parties) == 0 {
		return model.WrapError(nil, model.ErrCodeValidation, "no conflict parties")
	}

	return nil
}

func (r *Resolver) validateStrategy(strategy *Strategy) error {
	if strategy.ID == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty strategy ID")
	}

	if len(strategy.Actions) == 0 {
		return model.WrapError(nil, model.ErrCodeValidation, "no strategy actions")
	}

	return nil
}

func (r *Resolver) updateMetrics() {
	point := types.MetricPoint{
		Timestamp: time.Now(),
		Values: map[string]float64{
			"total_conflicts": float64(len(r.state.conflicts)),
			"resolved_count":  float64(r.state.metrics.ResolvedCount),
			"success_rate":    r.state.metrics.SuccessRate,
		},
	}

	r.state.metrics.History = append(r.state.metrics.History, point)

	// 限制历史记录数量
	if len(r.state.metrics.History) > maxMetricsHistory {
		r.state.metrics.History = r.state.metrics.History[1:]
	}
}

func generateResolutionID() string {
	return fmt.Sprintf("res_%d", time.Now().UnixNano())
}

// ResolveConflict 解决单个冲突
func (r *Resolver) ResolveConflict(conflict *Conflict) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 验证冲突
	if err := r.validateConflict(conflict); err != nil {
		return err
	}

	// 创建解决方案
	resolution := &Resolution{
		ID:         generateResolutionID(),
		ConflictID: conflict.ID,
		Type:       determineResolutionType(conflict),
		Status:     "initiated",
		Steps:      make([]ResolutionStep, 0),
		Created:    time.Now(),
	}

	// 执行解决方案
	if err := r.executeResolution(resolution); err != nil {
		return err
	}

	// 存储解决方案
	r.state.resolutions[resolution.ID] = resolution

	// 更新指标
	r.updateMetrics()

	return nil
}
