// system/types/trace_types.go

package types

import (
	"time"

	"github.com/Corphon/daoflow/model"
)

// 复用 model 包中的类型
type (
	TraceStatus = model.ModelState // 使用model.Status作为基础状态类型
)

// TraceID 追踪ID类型
type TraceID string

// SpanID 跨度ID类型
type SpanID string

// TraceEvent 追踪事件
type TraceEvent struct {
	model.ModelEvent // 嵌入模型事件

	TraceID TraceID // 追踪ID
	SpanID  SpanID  // 跨度ID

}

// SpanStatus 跨度状态常量
type SpanStatus string

const (
	SpanStatusNone     SpanStatus = "none"     // 初始状态
	SpanStatusActive   SpanStatus = "active"   // 活动状态
	SpanStatusComplete SpanStatus = "complete" // 完成状态
	SpanStatusError    SpanStatus = "error"    // 错误状态
)

// TraceConfig 追踪配置
type TraceConfig struct {
	// 存储配置
	StoragePath   string // 存储路径
	RetentionDays int    // 保留天数
	BatchSize     int    // 批处理大小
	BufferSize    int    // 缓冲区大小

	// 处理配置
	FlushInterval    time.Duration // 刷新间隔
	AnalysisInterval time.Duration // 分析间隔
	Compression      bool          // 是否启用压缩
	AsyncWrite       bool          // 异步写入

	// 采样配置
	SampleRate   float64 // 采样率
	MaxQueueSize int     // 最大队列大小

	// 追踪选项
	EnableMetrics bool // 启用指标采集
	EnableEvents  bool // 启用事件记录
	IncludeModel  bool // 包含模型信息
}

// TracePattern 追踪模式
type TracePattern struct {
	ID         string                 // 模式ID
	Type       string                 // 模式类型
	Confidence float64                // 置信度
	StartTime  time.Time              // 开始时间
	EndTime    time.Time              // 结束时间
	SpanIDs    []SpanID               // 相关跨度
	Properties map[string]interface{} // 属性
}

// Bottleneck 系统瓶颈
type Bottleneck struct {
	ID         string        // 瓶颈ID
	Type       string        // 瓶颈类型
	Resource   string        // 资源类型
	Severity   float64       // 严重程度
	Impact     float64       // 影响程度
	Duration   time.Duration // 持续时间
	Suggestion string        // 改进建议
}
