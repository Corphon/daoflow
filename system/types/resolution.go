// system/types/resolution.go

package types

import (
	"time"
)

// ResolutionRecord 冲突解决记录
type ResolutionRecord struct {
	ConflictID   string                 `json:"conflict_id"`   // 冲突ID
	SessionID    string                 `json:"session_id"`    // 会话ID
	ResolvedAt   time.Time              `json:"resolved_at"`   // 解决时间
	Status       string                 `json:"status"`        // 解决状态
	Error        error                  `json:"error"`         // 错误信息
	ErrorDetails *ErrorDetails          `json:"error_details"` // 错误详情
	Participants []string               `json:"participants"`  // 参与者列表
	Details      map[string]interface{} `json:"details"`       // 详细信息
}

// ErrorDetails 错误详细信息
type ErrorDetails struct {
	Code    string                 `json:"code"`    // 错误代码
	Message string                 `json:"message"` // 错误消息
	Time    time.Time              `json:"time"`    // 错误时间
	Stack   string                 `json:"stack"`   // 堆栈信息（可选）
	Context map[string]interface{} `json:"context"` // 错误上下文
}

// ResolutionResult 解决结果
type ResolutionResult struct {
	Success    bool               `json:"success"`    // 是否成功
	Confidence float64            `json:"confidence"` // 置信度
	Impact     map[string]float64 `json:"impact"`     // 影响评估
	Metrics    ResolutionMetrics  `json:"metrics"`    // 解决指标
	Feedback   []string           `json:"feedback"`   // 反馈信息
}

// ResolutionMetrics 解决指标
type ResolutionMetrics struct {
	Duration    time.Duration `json:"duration"`     // 解决耗时
	StepCount   int           `json:"step_count"`   // 步骤数量
	RetryCount  int           `json:"retry_count"`  // 重试次数
	SuccessRate float64       `json:"success_rate"` // 成功率
}

//-------------------------------------------------------------
