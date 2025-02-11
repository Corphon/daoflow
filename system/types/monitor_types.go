// system/types/monitor_types.go

package types

import (
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
	"github.com/Corphon/daoflow/model"
)

// MetricType 指标类型常量
const (
	MetricEnergy MetricType = iota
	MetricField
	MetricQuantum
	MetricEmergence
)

// Alert 告警信息
type Alert struct {
	ID      string                 `json:"id"`      // 告警ID
	Type    string                 `json:"type"`    // 告警类型
	Level   AlertLevel             `json:"level"`   // 告警级别
	Source  string                 `json:"source"`  // 告警源
	Target  string                 `json:"target"`  // 告警目标
	Message string                 `json:"message"` // 告警消息
	Time    time.Time              `json:"time"`    // 发生时间
	Status  string                 `json:"status"`  // 告警状态
	Labels  map[string]string      `json:"labels"`  // 标签
	Details map[string]interface{} `json:"details"` // 详细信息
}

// NotificationChannel 通知渠道类型
type NotificationChannel string

// 通知渠道常量
const (
	ChannelEmail   NotificationChannel = "email"
	ChannelWebhook NotificationChannel = "webhook"
	ChannelMessage NotificationChannel = "message"
	ChannelConsole NotificationChannel = "console"
	ChannelLog     NotificationChannel = "log"
)

// NotificationTarget 通知目标配置
type NotificationTarget struct {
	ID       string              `json:"id"`       // 目标ID
	Name     string              `json:"name"`     // 目标名称
	Channel  NotificationChannel `json:"channel"`  // 通知渠道
	Config   map[string]string   `json:"config"`   // 渠道配置
	Filters  []string            `json:"filters"`  // 告警过滤器
	Template string              `json:"template"` // 消息模板
	Enabled  bool                `json:"enabled"`  // 是否启用
}

// AlertConfig 告警配置
type AlertConfig struct {
	// 基础配置
	Enabled       bool          // 是否启用
	CheckInterval time.Duration // 检查间隔
	MaxAlerts     int           // 最大告警数

	// 处理配置
	MaxConcurrent int           // 最大并发数
	RetryCount    int           // 重试次数
	Timeout       time.Duration // 超时时间
	QueueSize     int           // 队列大小

	// 缓冲配置
	BufferSize int // 缓冲区大小
	BatchSize  int // 批处理大小

	// 阈值配置
	Thresholds map[string]float64 // 阈值配置

	// 通知配置
	RetryInterval time.Duration // 重试间隔
	MaxRetries    int           // 最大重试次数

	// 通知渠道
	Channels  []string          // 通知渠道
	Templates map[string]string // 消息模板
}

// AlertRule 告警规则
type AlertRule struct {
	Name      string        // 规则名称
	Condition string        // 告警条件
	Duration  time.Duration // 持续时间
	Severity  string        // 严重程度
	Actions   []string      // 触发动作
}

// HealthStatus 健康状态
type HealthStatus struct {
	Status    string             // 状态
	Score     float64            // 健康评分
	Details   map[string]float64 // 详细指标
	LastCheck time.Time          // 最后检查时间
	Issues    []Issue            // 存在的问题
}

// Issue 问题信息
type Issue struct {
	ID       string                 // 问题ID
	Type     string                 // 问题类型
	Severity string                 // 严重程度
	Message  string                 // 问题描述
	Source   string                 // 问题源
	Time     time.Time              // 发现时间
	Status   string                 // 问题状态
	Context  map[string]interface{} // 上下文信息
}

// SystemStatus 系统状态
type SystemStatus struct {
	Status     string                     // 运行状态
	Health     float64                    // 健康度
	Components map[string]ComponentStatus // 组件状态
	Resources  map[string]ResourceStatus  // 资源状态
	StartTime  time.Time                  // 启动时间
	UpdateTime time.Time                  // 更新时间
}

// ComponentStatus 组件状态
type ComponentStatus struct {
	Status    string             // 状态
	Health    float64            // 健康度
	Metrics   map[string]float64 // 指标
	LastError string             // 最后错误
}

// ResourceStatus 资源状态
type ResourceStatus struct {
	Usage     float64 // 使用率
	Available float64 // 可用量
	Total     float64 // 总量
	Limit     float64 // 限制
}

// MetricsStorage 指标存储器
type MetricsStorage struct {
	mu sync.RWMutex

	// 存储配置
	config struct {
		RetentionPeriod time.Duration // 保留周期
		MaxCapacity     int           // 最大容量
		BatchSize       int           // 批处理大小
	}

	// 指标数据
	data struct {
		current MetricsData             // 当前指标
		history []MetricsData           // 历史指标
		series  map[string]MetricSeries // 时间序列
	}
}

// MetricsAnalyzer 指标分析器
type MetricsAnalyzer struct {
	mu sync.RWMutex

	// 分析配置
	config struct {
		Thresholds map[string]float64       // 阈值配置
		Intervals  map[string]time.Duration // 分析间隔
		Patterns   []string                 // 模式识别
	}

	// 分析缓存
	cache struct {
		lastAnalysis time.Time
		results      []AnalysisResult
		predictions  []Prediction
	}
}

// AlertManager 告警管理器
type AlertManager struct {
	mu sync.RWMutex

	// 告警配置
	config AlertConfig

	// 告警状态
	state struct {
		activeAlerts map[string]Alert
		alertHistory []Alert
		mutedAlerts  map[string]time.Time
	}

	// 处理器和通知器
	handler  *AlertHandler
	notifier *AlertNotifier
}

// AlertHandler 告警处理器
type AlertHandler struct {
	mu sync.RWMutex

	// 处理配置
	config struct {
		MaxRetries    int
		RetryInterval time.Duration
		Timeout       time.Duration
	}

	// 处理状态
	state struct {
		processing map[string]Alert
		processed  []Alert
		failures   map[string]error
	}
}

// AlertNotifier 告警通知器
type AlertNotifier struct {
	Mu sync.RWMutex

	// 配置
	Config struct {
		Channels   []string          // 通知渠道
		Templates  map[string]string // 消息模板
		Throttling time.Duration     // 节流间隔
	}

	// 状态
	State struct {
		PendingNotifications []Alert   // 待处理通知
		SentNotifications    []Alert   // 已发送通知
		LastNotification     time.Time // 最后通知时间
	}
}

// HealthChecker 健康检查器
type HealthChecker struct {
	mu sync.RWMutex

	// 检查配置
	config struct {
		Intervals  map[string]time.Duration
		Thresholds map[string]float64
		Required   []string
	}

	// 检查状态
	state struct {
		checks    map[string]HealthCheck
		results   map[string]HealthResult
		lastCheck time.Time
	}
}

// HealthReporter 健康报告器
type HealthReporter struct {
	mu sync.RWMutex

	// 报告配置
	config struct {
		Interval time.Duration
		Format   string
		Output   string
	}

	// 报告状态
	state struct {
		reports    []HealthReport
		lastReport time.Time
		statistics HealthStats
	}
}

// HealthDiagnoser 健康诊断器
type HealthDiagnoser struct {
	mu sync.RWMutex

	// 诊断配置
	config struct {
		Rules    []DiagnosisRule
		Severity map[string]int
		Actions  map[string]string
	}

	// 诊断状态
	state struct {
		diagnoses []Diagnosis
		issues    map[string]Issue
		solutions map[string]Solution
	}
}

// MetricStatus 指标状态类型
type MetricStatus string

const (
	StatusOK      MetricStatus = "OK"      // 正常
	StatusError   MetricStatus = "ERROR"   // 错误
	StatusWarn    MetricStatus = "WARN"    // 警告
	StatusUnknown MetricStatus = "UNKNOWN" // 未知
)

// MetricsData 指标数据
type MetricsData struct {
	// 基础信息
	ID        string       `json:"id"`        // 指标ID
	Timestamp time.Time    `json:"timestamp"` // 时间戳
	Status    MetricStatus `json:"status"`    // 指标状态

	// 系统指标
	System struct {
		Energy    float64            `json:"energy"`
		Field     *core.FieldState   `json:"field"`     // 改为指针
		Quantum   *core.QuantumState `json:"quantum"`   // 改为指针
		Emergence *EmergentProperty  `json:"emergence"` // 改为指针
	} `json:"system"`

	// 模型指标
	Model struct {
		Integration float64 `json:"integration"` // 整体集成度
		Coherence   float64 `json:"coherence"`   // 整体相干性
		Emergence   float64 `json:"emergence"`   // 涌现程度
	} `json:"model"`

	// 自定义指标
	Custom map[string]interface{} `json:"custom"` // 自定义指标
}

// MetricSeries 指标序列
type MetricSeries struct {
	Name        string            `json:"name"`        // 指标名称
	Type        string            `json:"type"`        // 指标类型
	Description string            `json:"description"` // 描述
	Unit        string            `json:"unit"`        // 单位
	Labels      map[string]string `json:"labels"`      // 标签
	Values      []MetricValue     `json:"values"`      // 历史值
	Statistics  MetricStats       `json:"statistics"`  // 统计信息
}

// MetricValue 指标值
type MetricValue struct {
	Value     float64           `json:"value"`     // 值
	Timestamp time.Time         `json:"timestamp"` // 时间戳
	Labels    map[string]string `json:"labels"`    // 标签
}

// MetricStats 指标统计
type MetricStats struct {
	Min        float64   `json:"min"`         // 最小值
	Max        float64   `json:"max"`         // 最大值
	Sum        float64   `json:"sum"`         // 总和
	Count      int64     `json:"count"`       // 计数
	Average    float64   `json:"average"`     // 平均值
	Variance   float64   `json:"variance"`    // 方差
	LastUpdate time.Time `json:"last_update"` // 最后更新
}

// AnalysisResult 分析结果
type AnalysisResult struct {
	ID         string             // 结果ID
	Type       string             // 分析类型
	Timestamp  time.Time          // 分析时间
	Metrics    []string           // 相关指标
	Values     map[string]float64 // 分析值
	Patterns   []Pattern          // 识别模式
	Anomalies  []Anomaly          // 异常项
	Confidence float64            // 置信度
}

// Prediction 预测结果
type Prediction struct {
	ID         string           // 预测ID
	Type       string           // 预测类型
	Target     string           // 预测目标
	Values     []PredictedValue // 预测值
	Confidence float64          // 置信度
	Window     time.Duration    // 预测窗口
	UpdatedAt  time.Time        // 更新时间
}

// 辅助类型定义
type Pattern struct {
	Type       string                 // 模式类型
	Score      float64                // 匹配分数
	Components []string               // 组成部分
	Properties map[string]interface{} // 属性
}

type Anomaly struct {
	Type       string    // 异常类型
	Severity   float64   // 严重程度
	Metric     string    // 相关指标
	Threshold  float64   // 触发阈值
	Value      float64   // 实际值
	DetectedAt time.Time // 检测时间
}

type PredictedValue struct {
	Value       float64    // 预测值
	Timestamp   time.Time  // 时间点
	Range       [2]float64 // 置信区间
	Probability float64    // 概率
}

// HealthCheck 健康检查定义
type HealthCheck struct {
	ID        string        // 检查ID
	Type      string        // 检查类型
	Name      string        // 检查名称
	Target    string        // 检查目标
	Interval  time.Duration // 检查间隔
	Timeout   time.Duration // 超时时间
	Validator func() error  // 验证函数
	Required  bool          // 是否必需
}

// HealthResult 健康检查结果
type HealthResult struct {
	CheckID   string                 // 检查ID
	Status    string                 // 状态
	Score     float64                // 健康评分
	Message   string                 // 消息
	Details   map[string]interface{} // 详情
	StartTime time.Time              // 开始时间
	EndTime   time.Time              // 结束时间
	Duration  time.Duration          // 耗时
}

// HealthReport 健康报告
type HealthReport struct {
	ID        string        // 报告ID
	Timestamp time.Time     // 生成时间
	Duration  time.Duration // 报告周期

	// 系统概况
	Summary struct {
		Status string        // 系统状态
		Health float64       // 健康度
		Uptime time.Duration // 运行时间
		Issues int           // 问题数量
	}

	// 组件状态
	Components map[string]ComponentStatus

	// 检查结果
	Checks []HealthResult

	// 资源状态
	Resources map[string]ResourceStatus

	// 指标统计
	Metrics MetricsData

	// 问题和建议
	Issues          []Issue
	Recommendations []string
}

// HealthStats 健康统计
type HealthStats struct {
	// 基础统计
	TotalChecks  int64 // 总检查次数
	PassedChecks int64 // 通过检查数
	FailedChecks int64 // 失败检查数

	// 健康度统计
	AverageHealth float64 // 平均健康度
	MinHealth     float64 // 最低健康度
	MaxHealth     float64 // 最高健康度

	// 性能统计
	AverageLatency time.Duration // 平均延迟
	MaxLatency     time.Duration // 最大延迟

	// 故障统计
	TotalIssues    int64 // 总问题数
	CriticalIssues int64 // 严重问题数
	WarningIssues  int64 // 警告问题数

	// 时间统计
	LastUpdate   time.Time     // 最后更新
	ReportPeriod time.Duration // 统计周期
}

// DiagnosisRule 诊断规则
type DiagnosisRule struct {
	ID         string               // 规则ID
	Name       string               // 规则名称
	Type       string               // 规则类型
	Conditions []DiagnosisCondition // 诊断条件
	Actions    []DiagnosisAction    // 诊断动作
	Severity   string               // 严重程度
	Priority   int                  // 优先级
	Enabled    bool                 // 是否启用
}

// DiagnosisCondition 诊断条件
type DiagnosisCondition struct {
	Type     string        // 条件类型
	Target   string        // 检查目标
	Operator string        // 操作符
	Value    interface{}   // 比较值
	Duration time.Duration // 持续时间
}

// DiagnosisAction 诊断动作
type DiagnosisAction struct {
	Type       string                 // 动作类型
	Target     string                 // 目标对象
	Operation  string                 // 操作类型
	Parameters map[string]interface{} // 操作参数
}

// Diagnosis 诊断结果
type Diagnosis struct {
	ID        string        // 诊断ID
	RuleID    string        // 规则ID
	Type      string        // 诊断类型
	Level     string        // 严重等级
	Message   string        // 诊断消息
	Issues    []Issue       // 发现的问题
	Solutions []Solution    // 建议的解决方案
	StartTime time.Time     // 开始时间
	EndTime   time.Time     // 结束时间
	Duration  time.Duration // 诊断耗时
}

// Solution 解决方案
type Solution struct {
	ID        string         // 方案ID
	Type      string         // 方案类型
	IssueID   string         // 问题ID
	Steps     []SolutionStep // 解决步骤
	Status    string         // 方案状态
	Priority  int            // 优先级
	Success   bool           // 是否成功
	StartTime time.Time      // 开始时间
	EndTime   time.Time      // 结束时间
}

// SolutionStep 解决步骤
type SolutionStep struct {
	ID        string          // 步骤ID
	Type      string          // 步骤类型
	Action    DiagnosisAction // 执行动作
	Status    string          // 步骤状态
	Result    string          // 执行结果
	StartTime time.Time       // 开始时间
	EndTime   time.Time       // 结束时间
}

// EmergentProperty 涌现属性
type EmergentProperty struct {
	ID         string             `json:"id"`         // 属性ID
	Type       string             `json:"type"`       // 属性类型
	Pattern    *EmergentPattern   `json:"pattern"`    // 关联模式
	Properties map[string]float64 `json:"properties"` // 属性值
	Value      float64            `json:"value"`      // 当前值
	Stability  float64            `json:"stability"`  // 稳定性
	Energy     float64            `json:"energy"`     // 能量值
	Created    time.Time          `json:"created"`    // 创建时间
	Updated    time.Time          `json:"updated"`    // 更新时间
}

// EmergentPattern 涌现模式
type EmergentPattern struct {
	ID         string             `json:"id"`         // 模式ID
	Type       string             `json:"type"`       // 模式类型
	Components []PatternComponent `json:"components"` // 组成成分
	Properties map[string]float64 `json:"properties"` // 模式属性
	Strength   float64            `json:"strength"`   // 模式强度
	Energy     float64            `json:"energy"`     // 模式能量
	Created    time.Time          `json:"created"`    // 创建时间
	// 涌现特征
	Complexity float64 `json:"complexity"` // 复杂度
	Stability  float64 `json:"stability"`  // 稳定性
	Coupling   float64 `json:"coupling"`   // 耦合度
}

// PatternComponent 模式组件
type PatternComponent struct {
	Type   string             `json:"type"`   // 组件类型
	Weight float64            `json:"weight"` // 权重
	Role   string             `json:"role"`   // 角色
	State  map[string]float64 `json:"state"`  // 状态
}

// AlertLevel 告警级别
type AlertLevel string

const (
	AlertLevelInfo     AlertLevel = "INFO"     // 信息
	AlertLevelWarning  AlertLevel = "WARNING"  // 警告
	AlertLevelError    AlertLevel = "ERROR"    // 错误
	AlertLevelCritical AlertLevel = "CRITICAL" // 严重
)

// AlertData 告警数据
type AlertData struct {
	// 基础信息
	ID      string            `json:"id"`      // 告警ID
	Type    string            `json:"type"`    // 告警类型
	Level   AlertLevel        `json:"level"`   // 告警级别
	Message string            `json:"message"` // 告警消息
	Source  string            `json:"source"`  // 告警源
	Time    time.Time         `json:"time"`    // 发生时间
	Status  string            `json:"status"`  // 告警状态
	Labels  map[string]string `json:"labels"`  // 标签

	// 指标数据
	MetricName  string  // 指标名称
	MetricValue float64 // 指标值
	Threshold   float64 // 触发阈值

	// 模型数据
	ModelData *ModelAlertData `json:"model_data,omitempty"` // 模型相关数据

	// 扩展信息
	Details map[string]interface{} `json:"details,omitempty"` // 详细信息
}

// ModelAlertData 模型告警数据
type ModelAlertData struct {
	Type    model.ModelType    // 模型类型
	Metrics model.ModelMetrics // 模型指标
	State   model.ModelState   // 模型状态
}

// HandlerStatus 处理器状态
type HandlerStatus struct {
	IsRunning    bool      // 是否运行中
	ActiveCount  int       // 活跃处理数
	TotalHandled int64     // 总处理数
	LastError    error     // 最后错误
	ErrorCount   int       // 错误数量
	StartTime    time.Time // 启动时间
	LastUpdate   time.Time // 最后更新
}

// MetricsConfig 指标配置结构
type MetricsConfig struct {
	// 基础配置
	Base struct {
		Enabled        bool               `json:"enabled"`          // 是否启用
		Interval       time.Duration      `json:"interval"`         // 采集间隔
		SampleInterval time.Duration      `json:"sample_interval"`  // 采样间隔
		BufferSize     int                `json:"buffer_size"`      // 缓冲区大小
		MaxHistory     int                `json:"max_history"`      // 最大历史记录
		MaxHistorySize int                `json:"max_history_size"` // 最大历史大小
		HistorySize    int                `json:"history_size"`     // 历史大小容量
		Thresholds     map[string]float64 `json:"thresholds"`       // 阈值配置
	} `json:"base"`

	// 报告配置
	Report struct {
		ReportInterval time.Duration      `json:"report_interval"` // 报告间隔
		ReportFormat   string             `json:"report_format"`   // 报告格式
		OutputPath     string             `json:"output_path"`     // 输出路径
		Thresholds     map[string]float64 `json:"thresholds"`      // 报告阈值
		Filters        []string           `json:"filters"`         // 指标过滤器
	} `json:"report"`

	// 指标配置
	Metrics struct {
		System      bool `json:"system"`      // 系统指标
		Process     bool `json:"process"`     // 进程指标
		Resource    bool `json:"resource"`    // 资源指标
		Performance bool `json:"performance"` // 性能指标

	} `json:"metrics"`

	// 告警配置
	Alerts struct {
		Enabled      bool    `json:"enabled"`        // 启用告警
		MinSeverity  float64 `json:"min_severity"`   // 最小严重度
		MaxQueueSize int     `json:"max_queue_size"` // 最大队列大小
	} `json:"alerts"`
}

// Insight 分析洞察
type Insight struct {
	ID             string                 `json:"id"`              // 洞察ID
	Type           string                 `json:"type"`            // 洞察类型
	Level          AlertSeverity          `json:"level"`           // 严重程度
	Message        string                 `json:"message"`         // 洞察消息
	Source         string                 `json:"source"`          // 洞察来源
	Recommendation string                 `json:"recommendation"`  // 建议措施
	RelatedMetrics []string               `json:"related_metrics"` // 相关指标
	Details        map[string]interface{} `json:"details"`         // 详细信息
	Created        time.Time              `json:"created"`         // 创建时间
}

// AlertSeverity 转 AlertLevel 的转换方法
func AlertSeverityToLevel(severity AlertSeverity) AlertLevel {
	switch severity {
	case SeverityInfo:
		return AlertLevelInfo
	case SeverityWarning:
		return AlertLevelWarning
	case SeverityError:
		return AlertLevelError
	case SeverityCritical:
		return AlertLevelCritical
	default:
		return AlertLevelInfo
	}
}

// Report 监控报告
type Report struct {
	// 基本信息
	ID        string    `json:"id"`        // 报告ID
	Timestamp time.Time `json:"timestamp"` // 生成时间
	Period    string    `json:"period"`    // 报告周期

	// 系统概览
	Summary struct {
		Status string  `json:"status"` // 系统状态
		Health float64 `json:"health"` // 健康度
		Issues int     `json:"issues"` // 问题数量
	} `json:"summary"`

	// 详细指标
	Metrics MetricsData `json:"metrics"` // 当前指标数据

	// 趋势分析
	Trends struct {
		Energy    []float64 `json:"energy"`    // 能量趋势
		Field     []float64 `json:"field"`     // 场趋势
		Coherence []float64 `json:"coherence"` // 相干性趋势
	} `json:"trends"`

	// 建议措施
	Recommendations []string `json:"recommendations"` // 优化建议
}
