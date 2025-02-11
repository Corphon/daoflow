// model/state_transitor.go

package model

import (
	"sync"
	"time"
)

// StateTransitor 状态转换器
type StateTransitor struct {
	mu sync.RWMutex

	// 转换配置
	config struct {
		maxAttempts   int           // 最大尝试次数
		timeout       time.Duration // 转换超时
		validateSteps bool          // 是否验证步骤
		autoRollback  bool          // 自动回滚
	}

	// 转换状态
	state struct {
		current   *TransitionState  // 当前转换状态
		history   []StateTransition // 历史记录
		rollbacks []RollbackInfo    // 回滚信息
	}
}

// TransitionState 转换状态
type TransitionState struct {
	ID        string
	From      ModelState
	To        ModelState
	Steps     []TransitionStep
	StartTime time.Time
	Status    string
}

// TransitionStep 转换步骤
type TransitionStep struct {
	ID         string                 // 步骤ID
	Type       string                 // 步骤类型
	Action     string                 // 执行动作
	Parameters map[string]interface{} // 步骤参数
	Validation []ValidationRule       // 验证规则
	Status     string                 // 步骤状态
	StartTime  time.Time              // 开始时间
	EndTime    time.Time              // 结束时间
	Error      error                  // 执行错误
}

// ValidationRule 验证规则
type ValidationRule struct {
	Type       string                 // 规则类型
	Target     string                 // 验证目标
	Condition  string                 // 验证条件
	Parameters map[string]interface{} // 规则参数
	ErrorMsg   string                 // 错误消息
}

// RollbackInfo 回滚信息
type RollbackInfo struct {
	TransitionID string
	Steps        []RollbackStep
	Timestamp    time.Time
	Success      bool
}

// RollbackStep 回滚步骤
type RollbackStep struct {
	ID         string                 // 步骤ID
	TargetStep string                 // 对应转换步骤ID
	Action     string                 // 回滚动作
	Parameters map[string]interface{} // 回滚参数
	Status     string                 // 执行状态
	StartTime  time.Time              // 开始时间
	EndTime    time.Time              // 结束时间
	Error      error                  // 执行错误
}

// NewStateTransitor 创建状态转换器
func NewStateTransitor() *StateTransitor {
	st := &StateTransitor{}

	// 初始化配置
	st.config.maxAttempts = 3
	st.config.timeout = 30 * time.Second
	st.config.validateSteps = true
	st.config.autoRollback = true

	// 初始化状态
	st.state.history = make([]StateTransition, 0)
	st.state.rollbacks = make([]RollbackInfo, 0)

	return st
}
