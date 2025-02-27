//system/control/ctrlsync/types.go

package ctrlsync

import (
	"time"
)

// Participant 参与者
type Participant struct {
	ID       string                 // 参与者ID
	Type     string                 // 参与者类型
	Role     string                 // 参与角色
	Status   string                 // 当前状态
	JoinTime time.Time              // 加入时间
	LastPing time.Time              // 最后心跳
	Version  string                 // 版本信息
	Metadata map[string]interface{} // 元数据
}

// SyncSession 同步会话
type SyncSession struct {
	ID        string    // 会话ID
	Type      string    // 会话类型
	Status    string    // 会话状态
	StartTime time.Time // 开始时间
	EndTime   time.Time // 结束时间

	// 参与者信息
	Initiator    string            // 发起者
	Participants []string          // 参与者列表
	Roles        map[string]string // 角色映射

	// 同步内容
	Data      map[string]interface{} // 同步数据
	Changes   []StateChange          // 状态变更
	Conflicts []SyncConflict         // 同步冲突

	// 会话控制
	Config     SessionConfig // 会话配置
	Statistics SessionStats  // 会话统计
}

// CoordinationEvent 协调事件
type CoordinationEvent struct {
	ID       string                 // 事件ID
	Type     string                 // 事件类型
	Source   string                 // 事件源
	Target   string                 // 事件目标
	Time     time.Time              // 事件时间
	Data     map[string]interface{} // 事件数据
	Status   string                 // 事件状态
	Priority int                    // 事件优先级
}

// StateChange 状态变更
type StateChange struct {
	Field    string      // 变更字段
	OldValue interface{} // 原值
	NewValue interface{} // 新值
	Time     time.Time   // 变更时间
	Source   string      // 变更源
}

// SyncConflict 同步冲突
type SyncConflict struct {
	ID          string        // 冲突ID
	Type        string        // 冲突类型
	Description string        // 冲突描述
	Fields      []string      // 冲突字段
	Values      []interface{} // 冲突值
	Resolution  string        // 解决方案
	Status      string        // 冲突状态
}

// SessionConfig 会话配置
type SessionConfig struct {
	Timeout        time.Duration // 超时时间
	RetryLimit     int           // 重试限制
	BatchSize      int           // 批次大小
	ValidateRules  []string      // 验证规则
	ConflictPolicy string        // 冲突策略
}

// SessionStats 会话统计
type SessionStats struct {
	TotalChanges  int           // 总变更数
	Conflicts     int           // 冲突数
	ResolvedCount int           // 已解决数
	Duration      time.Duration // 持续时间
	SuccessRate   float64       // 成功率
	SyncLatency   time.Duration // 同步延迟
}

// Constants
const (
	// 参与者状态
	ParticipantStatusActive   = "active"
	ParticipantStatusInactive = "inactive"
	ParticipantStatusPending  = "pending"

	// 会话状态
	SessionStatusInit     = "initialized"
	SessionStatusRunning  = "running"
	SessionStatusComplete = "completed"
	SessionStatusFailed   = "failed"

	// 事件类型
	EventTypeStateChange = "state_change"
	EventTypeConflict    = "conflict"
	EventTypeResolution  = "resolution"

	// 冲突策略
	ConflictPolicyLatestWins = "latest_wins"
	ConflictPolicyMerge      = "merge"
	ConflictPolicyManual     = "manual"
)

// SyncStats 同步统计
type SyncStats struct {
	// 基础计数
	TotalTasks    int64 // 总任务数
	CompletedTask int64 // 已完成任务
	FailedTasks   int64 // 失败任务数

	// 性能指标
	AverageLatency time.Duration // 平均延迟
	MaxLatency     time.Duration // 最大延迟
	ProcessingTime time.Duration // 处理时间
	LastSyncTime   time.Time     // 最后同步时间

	// 同步指标
	SyncRate    float64 // 同步速率
	SuccessRate float64 // 成功率
	ErrorRate   float64 // 错误率

	// 资源统计
	BytesSynced       int64   // 同步数据量
	ResourceUsage     float64 // 资源使用率
	NetworkThroughput float64 // 网络吞吐量
}

//----------------------------------------
