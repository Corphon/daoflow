// system/interface.go

package system

import (
	"context"
	"time"

	"github.com/Corphon/daoflow/model"
)

// SystemInterface 系统接口
type SystemInterface interface {
	// 核心功能
	Initialize(ctx context.Context) error
	Start() error
	Stop() error
	Reset() error

	// 状态管理
	GetStatus() SystemStatus
	GetMetrics() Metrics
	GetHealth() HealthStatus

	// 演化控制
	Evolve(params EvolutionParams) error
	Adapt(params AdaptationParams) error
	Synchronize(params SyncParams) error
	Optimize(params OptimizationParams) error

	// 资源管理
	AllocateResources(req ResourceReq) error
	ReleaseResources(id string) error
	GetResourceStatus() ResourceStat
}

// ResourceReq 资源请求结构
type ResourceReq struct {
	ID       string    `json:"id"`       // 资源请求ID
	Type     string    `json:"type"`     // 资源类型
	Amount   float64   `json:"amount"`   // 请求数量
	Priority int       `json:"priority"` // 优先级
	Deadline time.Time `json:"deadline"` // 截止时间

	// 具体资源需求
	CPU    float64 `json:"cpu"`    // CPU需求 (核心数)
	Memory float64 `json:"memory"` // 内存需求 (MB)
	Energy float64 `json:"energy"` // 能量需求 (单位能量)

	// 可选参数
	Labels   map[string]string `json:"labels"`   // 资源标签
	Metadata map[string]string `json:"metadata"` // 元数据
}

// ResourceStat 资源状态结构
type ResourceStat struct {
	// 基础信息
	Timestamp time.Time `json:"timestamp"` // 状态更新时间

	// 容量信息
	Total struct {
		CPU    float64 `json:"cpu"`    // 总CPU核心数
		Memory float64 `json:"memory"` // 总内存(MB)
		Energy float64 `json:"energy"` // 总能量单位
	} `json:"total"`

	// 已用资源
	Used struct {
		CPU    float64 `json:"cpu"`    // 已用CPU
		Memory float64 `json:"memory"` // 已用内存
		Energy float64 `json:"energy"` // 已用能量
	} `json:"used"`

	// 可用资源
	Available struct {
		CPU    float64 `json:"cpu"`    // 可用CPU
		Memory float64 `json:"memory"` // 可用内存
		Energy float64 `json:"energy"` // 可用能量
	} `json:"available"`

	// 使用率
	Utilization struct {
		CPU    float64 `json:"cpu"`    // CPU使用率
		Memory float64 `json:"memory"` // 内存使用率
		Energy float64 `json:"energy"` // 能量使用率
	} `json:"utilization"`

	// 资源分配情况
	Allocations []ResourceAllocation `json:"allocations"` // 当前资源分配情况

	// 健康状态
	Health float64 `json:"health"` // 资源健康度
}

// ResourceAllocation 资源分配记录
type ResourceAllocation struct {
	ID        string    `json:"id"`         // 分配ID
	ReqID     string    `json:"req_id"`     // 请求ID
	Type      string    `json:"type"`       // 资源类型
	Amount    float64   `json:"amount"`     // 分配数量
	StartTime time.Time `json:"start_time"` // 分配开始时间
	EndTime   time.Time `json:"end_time"`   // 预计释放时间
	Status    string    `json:"status"`     // 分配状态
}

// Metrics 系统指标
type Metrics struct {
	Performance PerformanceMetrics `json:"performance"`
	Resources   ResourceMetrics    `json:"resources"`
	System      SystemMetrics      `json:"system"`
	Models      model.ModelMetrics `json:"models"`
	Timestamp   time.Time          `json:"timestamp"`
}

// PerformanceMetrics 性能指标
type PerformanceMetrics struct {
	CPU struct {
		Usage       float64 `json:"usage"`       // CPU使用率
		Load        float64 `json:"load"`        // CPU负载
		Temperature float64 `json:"temperature"` // CPU温度
	} `json:"cpu"`

	Memory struct {
		Usage     float64 `json:"usage"`     // 内存使用率
		Available float64 `json:"available"` // 可用内存
		Cached    float64 `json:"cached"`    // 缓存使用
	} `json:"memory"`

	Network struct {
		Throughput float64 `json:"throughput"` // 网络吞吐量
		Latency    float64 `json:"latency"`    // 网络延迟
		ErrorRate  float64 `json:"error_rate"` // 错误率
	} `json:"network"`

	QPS          float64 `json:"qps"`           // 每秒查询数
	Concurrency  float64 `json:"concurrency"`   // 并发度
	ResponseTime float64 `json:"response_time"` // 响应时间
}

// ResourceMetrics 资源指标
type ResourceMetrics struct {
	Energy struct {
		Total      float64 `json:"total"`      // 总能量
		Used       float64 `json:"used"`       // 已用能量
		Available  float64 `json:"available"`  // 可用能量
		Efficiency float64 `json:"efficiency"` // 能量效率
	} `json:"energy"`

	Computation struct {
		Capacity    float64 `json:"capacity"`    // 计算容量
		Utilization float64 `json:"utilization"` // 利用率
		Queue       float64 `json:"queue"`       // 队列长度
	} `json:"computation"`

	Storage struct {
		Total     float64 `json:"total"`     // 总存储
		Used      float64 `json:"used"`      // 已用存储
		Available float64 `json:"available"` // 可用存储
	} `json:"storage"`

	Allocation   float64 `json:"allocation"`   // 资源分配率
	Distribution float64 `json:"distribution"` // 资源分布均衡度
	Utilization  float64 `json:"utilization"`  // 总体利用率
}

// SystemMetrics 系统指标
type SystemMetrics struct {
	Evolution struct {
		Level     float64 `json:"level"`     // 演化等级
		Speed     float64 `json:"speed"`     // 演化速度
		Direction float64 `json:"direction"` // 演化方向
	} `json:"evolution"`

	Adaptation struct {
		Rate      float64 `json:"rate"`      // 适应率
		Accuracy  float64 `json:"accuracy"`  // 适应准确度
		Stability float64 `json:"stability"` // 适应稳定性
	} `json:"adaptation"`

	Synchronization struct {
		Degree    float64 `json:"degree"`    // 同步度
		Phase     float64 `json:"phase"`     // 相位
		Coherence float64 `json:"coherence"` // 相干性
	} `json:"synchronization"`

	Health    float64 `json:"health"`    // 系统健康度
	Stability float64 `json:"stability"` // 系统稳定性
	Balance   float64 `json:"balance"`   // 系统平衡度
}

// SystemStatus 系统状态
type SystemStatus struct {
	State      model.SystemState `json:"state"`       // 使用 model.SystemState
	ModelState model.ModelState  `json:"model_state"` // 使用 model.ModelState
	Time       time.Time         `json:"time"`
	Components []Component       `json:"components"`
	Resources  ResourceStat      `json:"resources"`
}

// SystemState 系统状态
type SystemState struct {
	Evolution struct {
		Level              float64   // 当前演化等级
		EnergyDistribution []float64 // 能量分布
		Direction          Vector3D  // 演化方向
		Speed              float64   // 演化速度
		Phase              float64   // 演化相位
	}

	Adaptation struct {
		Fitness      float64   // 适应度
		LearningRate float64   // 学习率
		Memory       []float64 // 记忆状态
		Response     []float64 // 响应历史
		Plasticity   float64   // 可塑性
	}

	Synchronization struct {
		Phase     float64 // 同步相位
		Coupling  float64 // 耦合强度
		Coherence float64 // 相干度
		Stability float64 // 稳定性
	}

	Optimization struct {
		Objective  float64 // 目标值
		Progress   float64 // 优化进度
		Efficiency float64 // 优化效率
		Quality    float64 // 解的质量
	}

	Emergence struct {
		Complexity  float64 // 系统复杂度
		Novelty     float64 // 新颖度
		Integration float64 // 集成度
		Potential   float64 // 涌现潜力
	}

	Energy float64   // 系统总能量
	Time   time.Time // 状态时间戳
}

// Component 组件信息
type Component struct {
	ID         string         `json:"id"`
	Type       ComponentType  `json:"type"`
	State      ComponentState `json:"state"`
	Health     float64        `json:"health"`
	LastUpdate time.Time      `json:"last_update"`
}

// ComponentType 组件类型
type ComponentType string

const (
	CompEvolution       ComponentType = "evolution"
	CompAdaptation      ComponentType = "adaptation"
	CompSynchronization ComponentType = "synchronization"
	CompOptimization    ComponentType = "optimization"
	CompEmergence       ComponentType = "emergence"
)

// ComponentState 组件状态
type ComponentState uint8

const (
	CompStateActive ComponentState = iota
	CompStateInactive
	CompStateError
)

// EvolutionParams 演化参数
type EvolutionParams struct {
	Target      EvolutionTarget `json:"target"`
	Constraints []Constraint    `json:"constraints"`
	Speed       float64         `json:"speed"`
	Threshold   float64         `json:"threshold"`
}

// EvolutionTarget 演化目标
type EvolutionTarget struct {
	Type     TargetType  `json:"type"`
	Value    interface{} `json:"value"`
	Priority int         `json:"priority"`
}

// TargetType 目标类型
type TargetType uint8

const (
	TargetStability TargetType = iota
	TargetEfficiency
	TargetResilience
	TargetHarmony
)

// AdaptationParams 适应参数
type AdaptationParams struct {
	Mode      AdaptMode    `json:"mode"`
	Patterns  []Pattern    `json:"patterns"`
	LearnRate float64      `json:"learn_rate"`
	Memory    MemoryConfig `json:"memory"`
}

// AdaptMode 适应模式
type AdaptMode uint8

const (
	AdaptReactive AdaptMode = iota
	AdaptProactive
	AdaptHybrid
)

// SyncParams 同步参数
type SyncParams struct {
	Target   SyncTarget             `json:"target"`
	Pattern  model.TransformPattern `json:"pattern"` // 使用 model.TransformPattern
	Coupling [][]float64            `json:"coupling"`
	Phase    float64                `json:"phase"`
}

// SyncTarget 同步目标
type SyncTarget struct {
	Systems   []string `json:"systems"`
	Mode      SyncMode `json:"mode"`
	Threshold float64  `json:"threshold"`
}

// SyncMode 同步模式
type SyncMode uint8

const (
	SyncPhase SyncMode = iota
	SyncFrequency
	SyncEnergy
)

// OptimizationParams 优化参数
type OptimizationParams struct {
	Objective Objective       `json:"objective"`
	ModelType model.ModelType `json:"model_type"` // 使用 model.ModelType
	Method    OptimMethod     `json:"method"`
	Config    OptimConfig     `json:"config"`
}

// MemoryConfig 记忆配置
type MemoryConfig struct {
	Capacity  int     `json:"capacity"`   // 记忆容量
	Duration  int64   `json:"duration"`   // 记忆持续时间(秒)
	DecayRate float64 `json:"decay_rate"` // 衰减率
	Threshold float64 `json:"threshold"`  // 记忆阈值

	// 记忆优先级配置
	Priority struct {
		Recent    float64 `json:"recent"`    // 最近事件权重
		Frequent  float64 `json:"frequent"`  // 频繁事件权重
		Important float64 `json:"important"` // 重要事件权重
	} `json:"priority"`

	// 记忆类型配置
	Types struct {
		ShortTerm bool `json:"short_term"` // 启用短期记忆
		LongTerm  bool `json:"long_term"`  // 启用长期记忆
		Working   bool `json:"working"`    // 启用工作记忆
	} `json:"types"`
}

// OptimConfig 优化配置
type OptimConfig struct {
	// 基本优化参数
	MaxIterations int     `json:"max_iterations"` // 最大迭代次数
	Tolerance     float64 `json:"tolerance"`      // 收敛容差
	LearningRate  float64 `json:"learning_rate"`  // 学习率
	Momentum      float64 `json:"momentum"`       // 动量因子

	// 约束配置
	Constraints struct {
		UseConstraints bool        `json:"use_constraints"` // 是否使用约束
		Penalty        float64     `json:"penalty"`         // 违反约束的惩罚因子
		Bounds         [][]float64 `json:"bounds"`          // 参数边界
	} `json:"constraints"`

	// 高级优化选项
	Advanced struct {
		ParallelDegree int     `json:"parallel_degree"` // 并行度
		BatchSize      int     `json:"batch_size"`      // 批处理大小
		Regularization float64 `json:"regularization"`  // 正则化系数
	} `json:"advanced"`

	// 早停策略
	EarlyStop struct {
		Enabled  bool    `json:"enabled"`   // 是否启用早停
		Patience int     `json:"patience"`  // 容忍次数
		MinDelta float64 `json:"min_delta"` // 最小改善量
	} `json:"early_stop"`
}

// Objective 优化目标
type Objective struct {
	Function string      `json:"function"`
	Weights  []float64   `json:"weights"`
	Bounds   [][]float64 `json:"bounds"`
}

// OptimMethod 优化方法
type OptimMethod uint8

const (
	OptimGradient OptimMethod = iota
	OptimGenetic
	OptimParticleSwarm
)

// HealthStatus 健康状态
type HealthStatus struct {
	Score     float64   `json:"score"`
	Issues    []Issue   `json:"issues"`
	LastCheck time.Time `json:"last_check"`
}

// Issue 问题记录
type Issue struct {
	Type     IssueType     `json:"type"`
	Severity IssueSeverity `json:"severity"`
	Message  string        `json:"message"`
	Time     time.Time     `json:"time"`
}

// IssueType 问题类型
type IssueType uint8

const (
	IssueResource IssueType = iota
	IssuePerformance
	IssueSecurity
	IssueStability
)

// IssueSeverity 问题严重度
type IssueSeverity uint8

const (
	SeverityInfo IssueSeverity = iota
	SeverityWarning
	SeverityError
	SeverityCritical
)

// SystemError 系统错误
type SystemError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details"`
	Time    time.Time `json:"time"`
}

// ErrorCode 错误代码
type ErrorCode uint32

const (
	ErrNone            ErrorCode = 0
	ErrInitialize      ErrorCode = 1000
	ErrEvolution       ErrorCode = 2000
	ErrAdaptation      ErrorCode = 3000
	ErrSynchronization ErrorCode = 4000
	ErrOptimization    ErrorCode = 5000
	ErrResource        ErrorCode = 6000
)

// EventType 事件类型
type EventType uint8

const (
	EventStateChange EventType = iota
	EventEvolution
	EventAdaptation
	EventSync
	EventOptimization
	EventResource
	EventError
)

// Event 系统事件
type Event struct {
	Type   EventType   `json:"type"`
	Source string      `json:"source"`
	Target string      `json:"target"`
	Data   interface{} `json:"data"`
	Time   time.Time   `json:"time"`
}

// Config 配置接口
type Config interface {
	Validate() error
	Apply(interface{}) error
	GetValue(key string) interface{}
	SetValue(key string, value interface{}) error
}
