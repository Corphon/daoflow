//system/control/state/manager.go

package state

import (
    "sync"
    "time"

    "github.com/Corphon/daoflow/meta/field"
    "github.com/Corphon/daoflow/model"
    "github.com/Corphon/daoflow/system/types"
)

// StateManager 状态管理器
type StateManager struct {
    mu sync.RWMutex

    // 基础配置
    config struct {
        checkInterval    time.Duration  // 状态检查间隔
        stateTimeout     time.Duration  // 状态超时时间
        maxStateHistory  int           // 最大历史记录数
        consistencyLevel string        // 一致性级别
    }

    // 状态存储
    state struct {
        current      *SystemState         // 当前状态
        history      []StateSnapshot      // 状态历史
        transitions  []StateTransition    // 转换历史
        validators   map[string]Validator // 状态验证器
    }

    // 依赖注入
    validator   *StateValidator
    transition  *StateTransition
}

// SystemState 系统状态
type SystemState struct {
    ID          string                 // 状态ID
    Version     int64                  // 版本号
    Components  map[string]*Component  // 组件状态
    Resources   map[string]*Resource   // 资源状态
    Properties  map[string]interface{} // 状态属性
    Timestamp   time.Time             // 状态时间戳
}

// Component 组件状态
type Component struct {
    ID          string                // 组件ID
    Type        string                // 组件类型
    Status      string                // 运行状态
    Health      float64               // 健康度
    Properties  map[string]interface{} // 组件属性
    LastUpdate  time.Time            // 最后更新时间
}

// Resource 资源状态
type Resource struct {
    ID          string                // 资源ID
    Type        string                // 资源类型
    Capacity    float64               // 总容量
    Usage       float64               // 当前使用量
    Allocated   float64               // 已分配量
    Properties  map[string]interface{} // 资源属性
}

// StateSnapshot 状态快照
type StateSnapshot struct {
    ID          string                // 快照ID
    StateID     string                // 状态ID
    Version     int64                 // 版本号
    Data        map[string]interface{} // 快照数据
    Timestamp   time.Time            // 快照时间
}

// Validator 状态验证器接口
type Validator interface {
    Validate(*SystemState) error
    ValidateTransition(*SystemState, *SystemState) error
}

// NewStateManager 创建新的状态管理器
func NewStateManager(
    validator *StateValidator,
    transition *StateTransition) *StateManager {
    
    sm := &StateManager{
        validator:  validator,
        transition: transition,
    }

    // 初始化配置
    sm.config.checkInterval = 1 * time.Second
    sm.config.stateTimeout = 30 * time.Second
    sm.config.maxStateHistory = 1000
    sm.config.consistencyLevel = "strong"

    // 初始化状态
    sm.state.current = &SystemState{
        ID:         generateStateID(),
        Version:    1,
        Components: make(map[string]*Component),
        Resources:  make(map[string]*Resource),
        Properties: make(map[string]interface{}),
        Timestamp:  time.Now(),
    }
    sm.state.history = make([]StateSnapshot, 0)
    sm.state.transitions = make([]StateTransition, 0)
    sm.state.validators = make(map[string]Validator)

    return sm
}

// GetCurrentState 获取当前状态
func (sm *StateManager) GetCurrentState() (*SystemState, error) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()

    if sm.state.current == nil {
        return nil, model.WrapError(nil, model.ErrCodeNotFound, "current state not found")
    }

    return sm.state.current, nil
}

// UpdateState 更新系统状态
func (sm *StateManager) UpdateState(newState *SystemState) error {
    if newState == nil {
        return model.WrapError(nil, model.ErrCodeValidation, "nil state")
    }

    sm.mu.Lock()
    defer sm.mu.Unlock()

    // 验证新状态
    if err := sm.validateState(newState); err != nil {
        return err
    }

    // 验证状态转换
    if err := sm.validateTransition(sm.state.current, newState); err != nil {
        return err
    }

    // 创建快照
    snapshot := sm.createSnapshot(sm.state.current)
    sm.state.history = append(sm.state.history, snapshot)

    // 更新状态
    newState.Version = sm.state.current.Version + 1
    newState.Timestamp = time.Now()
    sm.state.current = newState

    // 记录转换
    sm.recordTransition(snapshot.ID, newState.ID)

    // 清理历史记录
    sm.cleanupHistory()

    return nil
}

// RegisterValidator 注册状态验证器
func (sm *StateManager) RegisterValidator(name string, validator Validator) error {
    if validator == nil {
        return model.WrapError(nil, model.ErrCodeValidation, "nil validator")
    }

    sm.mu.Lock()
    defer sm.mu.Unlock()

    sm.state.validators[name] = validator
    return nil
}

// validateState 验证状态
func (sm *StateManager) validateState(state *SystemState) error {
    // 基础验证
    if state.ID == "" {
        return model.WrapError(nil, model.ErrCodeValidation, "empty state ID")
    }

    // 使用注册的验证器
    for _, validator := range sm.state.validators {
        if err := validator.Validate(state); err != nil {
            return err
        }
    }

    return nil
}

// validateTransition 验证状态转换
func (sm *StateManager) validateTransition(
    current, next *SystemState) error {
    
    // 使用注册的验证器
    for _, validator := range sm.state.validators {
        if err := validator.ValidateTransition(current, next); err != nil {
            return err
        }
    }

    return nil
}

// createSnapshot 创建状态快照
func (sm *StateManager) createSnapshot(state *SystemState) StateSnapshot {
    return StateSnapshot{
        ID:        generateSnapshotID(),
        StateID:   state.ID,
        Version:   state.Version,
        Data:      sm.serializeState(state),
        Timestamp: time.Now(),
    }
}

// cleanupHistory 清理历史记录
func (sm *StateManager) cleanupHistory() {
    if len(sm.state.history) > sm.config.maxStateHistory {
        sm.state.history = sm.state.history[1:]
    }
}

// 辅助函数

func (sm *StateManager) serializeState(
    state *SystemState) map[string]interface{} {
    
    data := make(map[string]interface{})
    // 序列化状态数据
    data["components"] = state.Components
    data["resources"] = state.Resources
    data["properties"] = state.Properties
    return data
}

func generateStateID() string {
    return fmt.Sprintf("state_%d", time.Now().UnixNano())
}

func generateSnapshotID() string {
    return fmt.Sprintf("snap_%d", time.Now().UnixNano())
}
