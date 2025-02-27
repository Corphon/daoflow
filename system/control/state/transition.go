//system/control/state/transition.go

package state

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
	"github.com/Corphon/daoflow/model"
)

const (
	maxTransitionHistory = 1000 // 最大转换历史记录数
)

// StateTransition 状态转换管理器
type StateTransition struct {
	mu sync.RWMutex

	ID         string                 // 转换ID
	SourceID   string                 // 源状态ID
	TargetID   string                 // 目标状态ID
	Type       string                 // 转换类型
	Timestamp  time.Time              // 转换时间
	Properties map[string]interface{} // 转换属性

	// 基础配置
	config struct {
		maxTransitions  int           // 最大转换数量
		timeoutDuration time.Duration // 转换超时时间
		retryLimit      int           // 重试限制
		atomicMode      bool          // 原子模式
	}

	// 转换状态
	state struct {
		transitions map[string]*Transition    // 活跃转换
		pending     map[string]*PendingChange // 待处理变更
		history     []TransitionRecord        // 转换历史
		rollbacks   map[string]*RollbackPlan  // 回滚计划
	}
}

// Transition 状态转换
type Transition struct {
	ID          string        // 转换ID
	Type        string        // 转换类型
	SourceState string        // 源状态ID
	TargetState string        // 目标状态ID
	Changes     []StateChange // 状态变更
	Status      string        // 转换状态
	StartTime   time.Time     // 开始时间
	EndTime     time.Time     // 结束时间
	Retries     int           // 重试次数
}

// StateChange 状态变更
type StateChange struct {
	Type       string      // 变更类型
	Target     string      // 变更目标
	OldValue   interface{} // 原值
	NewValue   interface{} // 新值
	Timestamp  time.Time   // 变更时间
	Reversible bool        // 是否可逆
}

// PendingChange 待处理变更
type PendingChange struct {
	ID           string    // 变更ID
	TransitionID string    // 关联转换ID
	Priority     int       // 优先级
	Dependencies []string  // 依赖变更
	Status       string    // 变更状态
	Created      time.Time // 创建时间
}

// TransitionRecord 转换记录
type TransitionRecord struct {
	TransitionID string                 // 转换ID
	Type         string                 // 记录类型
	Status       string                 // 转换状态
	Details      map[string]interface{} // 详细信息
	Timestamp    time.Time              // 记录时间
}

// RollbackPlan 回滚计划
type RollbackPlan struct {
	TransitionID string         // 转换ID
	Steps        []RollbackStep // 回滚步骤
	Status       string         // 计划状态
	Created      time.Time      // 创建时间
}

// RollbackStep 回滚步骤
type RollbackStep struct {
	ChangeID   string                 // 变更ID
	Action     string                 // 回滚动作
	Parameters map[string]interface{} // 回滚参数
	Status     string                 // 步骤状态
}

// -----------------------------------------
// NewStateTransition 创建新的状态转换管理器
func NewStateTransition() *StateTransition {
	st := &StateTransition{}

	// 初始化配置
	st.config.maxTransitions = 100
	st.config.timeoutDuration = 30 * time.Second
	st.config.retryLimit = 3
	st.config.atomicMode = true

	// 初始化状态
	st.state.transitions = make(map[string]*Transition)
	st.state.pending = make(map[string]*PendingChange)
	st.state.history = make([]TransitionRecord, 0)
	st.state.rollbacks = make(map[string]*RollbackPlan)

	return st
}

// BeginTransition 开始状态转换
func (st *StateTransition) BeginTransition(
	sourceState, targetState string) (*Transition, error) {

	st.mu.Lock()
	defer st.mu.Unlock()

	// 检查转换数量限制
	if len(st.state.transitions) >= st.config.maxTransitions {
		return nil, model.WrapError(nil, model.ErrCodeLimit, "max transitions reached")
	}

	// 创建新转换
	transition := &Transition{
		ID:          core.GenerateID(),
		Type:        determineTransitionType(sourceState, targetState),
		SourceState: sourceState,
		TargetState: targetState,
		Changes:     make([]StateChange, 0),
		Status:      "initiated",
		StartTime:   time.Now(),
	}

	// 存储转换
	st.state.transitions[transition.ID] = transition

	// 创建回滚计划
	st.createRollbackPlan(transition)

	return transition, nil
}

// determineTransitionType 确定转换类型
func determineTransitionType(sourceState, targetState string) string {
	// 基于状态组合确定转换类型
	switch {
	case sourceState == "inactive" && targetState == "active":
		return "activation"
	case sourceState == "active" && targetState == "inactive":
		return "deactivation"
	case sourceState == "active" && targetState == "error":
		return "failure"
	case sourceState == "error" && targetState == "active":
		return "recovery"
	case sourceState == targetState:
		return "update"
	default:
		return "transition"
	}
}

// createRollbackPlan 创建回滚计划
func (st *StateTransition) createRollbackPlan(transition *Transition) {
	plan := &RollbackPlan{
		TransitionID: transition.ID,
		Steps:        make([]RollbackStep, 0),
		Status:       "pending",
		Created:      time.Now(),
	}

	// 预创建回滚步骤槽位
	for i := 0; i < len(transition.Changes); i++ {
		plan.Steps = append(plan.Steps, RollbackStep{
			ChangeID:   fmt.Sprintf("%s_step_%d", transition.ID, i),
			Action:     "revert",
			Parameters: make(map[string]interface{}),
			Status:     "pending",
		})
	}

	// 存储回滚计划
	st.state.rollbacks[transition.ID] = plan
}

// CommitTransition 提交状态转换
func (st *StateTransition) CommitTransition(transitionID string) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	transition, exists := st.state.transitions[transitionID]
	if !exists {
		return model.WrapError(nil, model.ErrCodeNotFound, "transition not found")
	}

	// 验证转换状态
	if err := st.validateTransitionStatus(transition); err != nil {
		return err
	}

	// 提交所有变更
	if err := st.commitChanges(transition); err != nil {
		// 执行回滚
		st.rollbackTransition(transition)
		return err
	}

	// 更新转换状态
	transition.Status = "completed"
	transition.EndTime = time.Now()

	// 记录转换完成 - 将 Transition 转换为 TransitionData
	st.recordTransition(&TransitionData{
		ID:        transition.ID,
		SourceID:  transition.SourceState,
		TargetID:  transition.TargetState,
		Type:      transition.Type,
		Timestamp: time.Now(),
		Properties: map[string]interface{}{
			"status":   "completed",
			"end_time": transition.EndTime,
			"duration": transition.EndTime.Sub(transition.StartTime),
		},
	}, "completed", nil)

	// 清理相关资源
	st.cleanupTransition(transition)

	return nil
}

// rollbackTransition 执行转换回滚
func (st *StateTransition) rollbackTransition(transition *Transition) error {
	// 获取回滚计划
	plan, exists := st.state.rollbacks[transition.ID]
	if !exists {
		return model.WrapError(nil, model.ErrCodeNotFound, "rollback plan not found")
	}

	// 逆序执行回滚步骤
	for i := len(plan.Steps) - 1; i >= 0; i-- {
		step := plan.Steps[i]
		if step.Status != "pending" {
			continue
		}

		// 执行回滚动作
		if err := st.executeRollbackStep(transition, &step); err != nil {
			// 记录回滚失败
			st.recordTransition(
				&TransitionData{
					ID:        transition.ID,
					SourceID:  transition.SourceState,
					TargetID:  transition.TargetState,
					Type:      transition.Type,
					Timestamp: time.Now(),
					Properties: map[string]interface{}{
						"status": "completed",
					},
				},
				"completed",
				nil,
			)
			return err
		}

		step.Status = "completed"
	}

	// 更新回滚计划状态
	plan.Status = "completed"

	// 更新转换状态
	transition.Status = "rolled_back"

	return nil
}

// executeRollbackStep 执行回滚步骤
func (st *StateTransition) executeRollbackStep(transition *Transition, step *RollbackStep) error {
	// 查找对应的变更
	var change *StateChange
	for i := range transition.Changes {
		if transition.Changes[i].Type+fmt.Sprint(i) == step.ChangeID {
			change = &transition.Changes[i]
			break
		}
	}

	if change == nil || !change.Reversible {
		return model.WrapError(nil, model.ErrCodeValidation, "change not found or not reversible")
	}

	// 执行回滚操作
	temp := change.OldValue
	change.OldValue = change.NewValue
	change.NewValue = temp
	change.Timestamp = time.Now()

	return nil
}

// cleanupTransition 清理转换相关资源
func (st *StateTransition) cleanupTransition(transition *Transition) {
	// 删除活跃转换记录
	delete(st.state.transitions, transition.ID)

	// 清理相关的待处理变更
	for id, pending := range st.state.pending {
		if pending.TransitionID == transition.ID {
			delete(st.state.pending, id)
		}
	}

	// 清理回滚计划
	delete(st.state.rollbacks, transition.ID)

	// 如果历史记录超过限制，清理最旧的记录
	if len(st.state.history) > st.config.maxTransitions {
		st.state.history = st.state.history[1:]
	}
}

// AddChange 添加状态变更
func (st *StateTransition) AddChange(
	transitionID string,
	change StateChange) error {

	st.mu.Lock()
	defer st.mu.Unlock()

	transition, exists := st.state.transitions[transitionID]
	if !exists {
		return model.WrapError(nil, model.ErrCodeNotFound, "transition not found")
	}

	// 验证变更
	if err := st.validateChange(change); err != nil {
		return err
	}

	// 添加变更
	transition.Changes = append(transition.Changes, change)

	// 更新回滚计划
	st.updateRollbackPlan(transition, change)

	return nil
}

// validateChange 验证变更
func (st *StateTransition) validateChange(change StateChange) error {
	// 验证变更类型
	if change.Type == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty change type")
	}

	// 验证变更目标
	if change.Target == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty change target")
	}

	// 验证值对
	if change.OldValue == nil || change.NewValue == nil {
		return model.WrapError(nil, model.ErrCodeValidation, "nil change values")
	}

	// 验证时间戳
	if change.Timestamp.IsZero() {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid timestamp")
	}

	return nil
}

// updateRollbackPlan 更新回滚计划
func (st *StateTransition) updateRollbackPlan(transition *Transition, change StateChange) {
	// 获取回滚计划
	plan, exists := st.state.rollbacks[transition.ID]
	if !exists {
		return
	}

	// 创建新的回滚步骤
	step := RollbackStep{
		ChangeID: change.Type + fmt.Sprint(len(transition.Changes)-1),
		Action:   "revert",
		Parameters: map[string]interface{}{
			"target":    change.Target,
			"old_value": change.OldValue,
			"new_value": change.NewValue,
		},
		Status: "pending",
	}

	// 添加步骤到计划
	plan.Steps = append(plan.Steps, step)
}

// 辅助函数

func (st *StateTransition) validateTransitionStatus(
	transition *Transition) error {

	if transition.Status == "completed" {
		return model.WrapError(nil, model.ErrCodeValidation, "transition already completed")
	}

	if transition.Status == "failed" {
		return model.WrapError(nil, model.ErrCodeValidation, "transition failed")
	}

	return nil
}

func (st *StateTransition) commitChanges(transition *Transition) error {
	for _, change := range transition.Changes {
		pendingChange := st.createPendingChange(transition, change)

		// 执行变更
		if err := st.executeChange(pendingChange); err != nil {
			return err
		}
	}
	return nil
}

// createPendingChange 创建待处理变更
func (st *StateTransition) createPendingChange(transition *Transition, change StateChange) *PendingChange {
	return &PendingChange{
		ID:           fmt.Sprintf("change_%s_%d", transition.ID, len(transition.Changes)),
		TransitionID: transition.ID,
		Priority:     determinePriority(change.Type),
		Dependencies: findDependencies(transition.Changes, change),
		Status:       "pending",
		Created:      time.Now(),
	}
}

// executeChange 执行变更
func (st *StateTransition) executeChange(pending *PendingChange) error {
	// 检查依赖是否满足
	for _, depID := range pending.Dependencies {
		dep, exists := st.state.pending[depID]
		if !exists || dep.Status != "completed" {
			return model.WrapError(nil, model.ErrCodeDependency,
				fmt.Sprintf("dependency not satisfied: %s", depID))
		}
	}

	// 添加到待处理队列
	st.state.pending[pending.ID] = pending

	// 尝试执行变更
	if st.config.atomicMode {
		// 原子模式：所有依赖必须完成
		pending.Status = "completed"
	} else {
		// 非原子模式：可以并行执行
		go func() {
			// 异步执行变更
			time.Sleep(100 * time.Millisecond) // 模拟执行时间
			pending.Status = "completed"
		}()
	}

	return nil
}

// 内部辅助函数

// determinePriority 确定变更优先级
func determinePriority(changeType string) int {
	switch changeType {
	case "critical":
		return 0
	case "important":
		return 1
	case "normal":
		return 2
	default:
		return 3
	}
}

// findDependencies 查找变更依赖
func findDependencies(changes []StateChange, current StateChange) []string {
	var deps []string
	for i, change := range changes {
		// 如果当前变更的目标依赖于之前的变更
		if isDependentChange(current, change) {
			deps = append(deps, fmt.Sprintf("change_%d", i))
		}
	}
	return deps
}

// isDependentChange 检查变更依赖关系
func isDependentChange(current, previous StateChange) bool {
	// 1. 相同目标的变更
	if current.Target == previous.Target {
		return true
	}

	// 2. 组件依赖关系
	if isComponentChange(current.Type) && isComponentChange(previous.Type) {
		return checkComponentDependency(current.Target, previous.Target)
	}

	// 3. 资源依赖关系
	if isResourceChange(current.Type) && isResourceChange(previous.Type) {
		return checkResourceDependency(current.Target, previous.Target)
	}

	// 4. 状态转换顺序依赖
	if current.Type == "state" && previous.Type == "state" {
		return checkStateTransitionOrder(
			previous.OldValue.(string),
			previous.NewValue.(string),
			current.OldValue.(string),
			current.NewValue.(string))
	}

	// 5. 跨类型依赖
	return checkCrossDependency(current, previous)
}

// 辅助函数

// isComponentChange 检查是否为组件变更
func isComponentChange(changeType string) bool {
	return strings.HasPrefix(changeType, "component_")
}

// isResourceChange 检查是否为资源变更
func isResourceChange(changeType string) bool {
	return strings.HasPrefix(changeType, "resource_")
}

// checkComponentDependency 检查组件依赖关系
func checkComponentDependency(currentTarget, previousTarget string) bool {
	// TODO: 从组件依赖图中查找依赖关系
	// 当前简化实现：检查命名约定
	return strings.HasPrefix(currentTarget, previousTarget+"_") ||
		strings.HasPrefix(previousTarget, currentTarget+"_")
}

// checkResourceDependency 检查资源依赖关系
func checkResourceDependency(currentTarget, previousTarget string) bool {
	// TODO: 从资源依赖配置中查找依赖关系
	// 当前简化实现：检查资源路径
	return strings.HasPrefix(currentTarget, previousTarget+"/") ||
		strings.HasPrefix(previousTarget, currentTarget+"/")
}

// checkStateTransitionOrder 检查状态转换顺序
func checkStateTransitionOrder(prevOld, prevNew, currOld, currNew string) bool {
	// 定义有效的状态转换
	validTransitions := map[string][]string{
		"inactive": {"active"},
		"active":   {"inactive", "error"},
		"error":    {"inactive"},
	}

	// 1. 检查两个转换序列是否连续
	if prevNew == currOld {
		// 检查当前转换是否有效
		if validStates, exists := validTransitions[currOld]; exists {
			for _, validState := range validStates {
				if validState == currNew {
					return true
				}
			}
		}
	}

	// 2. 检查前一个转换是否有效
	if validStates, exists := validTransitions[prevOld]; exists {
		found := false
		for _, validState := range validStates {
			if validState == prevNew {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// 3. 检查当前转换是否有效
	if validStates, exists := validTransitions[currOld]; exists {
		found := false
		for _, validState := range validStates {
			if validState == currNew {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// 4. 检查两个转换组合是否形成有效的状态链
	// 例如：inactive->active->error 是有效的
	chainValid := false
	if prevOld == "inactive" && prevNew == "active" {
		if currOld == "active" && (currNew == "inactive" || currNew == "error") {
			chainValid = true
		}
	}

	return chainValid
}

// checkCrossDependency 检查跨类型依赖
func checkCrossDependency(current, previous StateChange) bool {
	// 定义跨类型依赖规则
	crossDependencies := map[string][]string{
		"component_status":    {"resource_allocation", "state_update"},
		"resource_allocation": {"state_update"},
		"state_update":        {"component_init"},
	}

	// 检查是否存在跨类型依赖
	if deps, exists := crossDependencies[current.Type]; exists {
		for _, dep := range deps {
			if previous.Type == dep {
				return true
			}
		}
	}

	return false
}

func (st *StateTransition) recordTransition(
	transitionData *TransitionData, // 改为接收 TransitionData
	recordType string,
	details map[string]interface{}) {

	record := TransitionRecord{
		TransitionID: transitionData.ID, // 使用 TransitionData 的字段
		Type:         recordType,
		Status:       "completed",
		Details:      details,
		Timestamp:    time.Now(),
	}

	st.state.history = append(st.state.history, record)

	// 限制历史记录数量
	if len(st.state.history) > maxTransitionHistory {
		st.state.history = st.state.history[1:]
	}
}
