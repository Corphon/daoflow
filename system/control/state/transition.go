//system/control/state/transition.go

package state

import (
    "sync"
    "time"

    "github.com/Corphon/daoflow/meta/field"
    "github.com/Corphon/daoflow/model"
    "github.com/Corphon/daoflow/system/types"
)

// StateTransition 状态转换管理器
type StateTransition struct {
    mu sync.RWMutex

    // 基础配置
    config struct {
        maxTransitions   int           // 最大转换数量
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
    ID           string                // 转换ID
    Type         string                // 转换类型
    SourceState  string                // 源状态ID
    TargetState  string                // 目标状态ID
    Changes      []StateChange         // 状态变更
    Status       string                // 转换状态
    StartTime    time.Time            // 开始时间
    EndTime      time.Time            // 结束时间
    Retries      int                  // 重试次数
}

// StateChange 状态变更
type StateChange struct {
    Type         string                // 变更类型
    Target       string                // 变更目标
    OldValue     interface{}           // 原值
    NewValue     interface{}           // 新值
    Timestamp    time.Time            // 变更时间
    Reversible   bool                 // 是否可逆
}

// PendingChange 待处理变更
type PendingChange struct {
    ID           string                // 变更ID
    TransitionID string                // 关联转换ID
    Priority     int                   // 优先级
    Dependencies []string              // 依赖变更
    Status       string                // 变更状态
    Created      time.Time            // 创建时间
}

// TransitionRecord 转换记录
type TransitionRecord struct {
    TransitionID string                // 转换ID
    Type         string                // 记录类型
    Status       string                // 转换状态
    Details      map[string]interface{} // 详细信息
    Timestamp    time.Time            // 记录时间
}

// RollbackPlan 回滚计划
type RollbackPlan struct {
    TransitionID string                // 转换ID
    Steps       []RollbackStep         // 回滚步骤
    Status      string                 // 计划状态
    Created     time.Time             // 创建时间
}

// RollbackStep 回滚步骤
type RollbackStep struct {
    ChangeID    string                // 变更ID
    Action      string                // 回滚动作
    Parameters  map[string]interface{} // 回滚参数
    Status      string                // 步骤状态
}

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
        ID:          generateTransitionID(),
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

    // 记录转换完成
    st.recordTransition(transition, "completed", nil)

    // 清理相关资源
    st.cleanupTransition(transition)

    return nil
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

func (st *StateTransition) recordTransition(
    transition *Transition,
    recordType string,
    details map[string]interface{}) {
    
    record := TransitionRecord{
        TransitionID: transition.ID,
        Type:        recordType,
        Status:      transition.Status,
        Details:     details,
        Timestamp:   time.Now(),
    }

    st.state.history = append(st.state.history, record)

    // 限制历史记录数量
    if len(st.state.history) > maxTransitionHistory {
        st.state.history = st.state.history[1:]
    }
}

func generateTransitionID() string {
    return fmt.Sprintf("trans_%d", time.Now().UnixNano())
}

const (
    maxTransitionHistory = 1000
)
