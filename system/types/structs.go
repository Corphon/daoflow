// system/types/structs.go

package types

import (
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
	"github.com/Corphon/daoflow/model"
)

// 复用模型包的基础类型
type (
	Vector3D     = model.Vector3D
	ModelState   = model.ModelState
	FieldState   = core.FieldState
	QuantumState = core.QuantumState
)

// MetricPoint 表示一个指标数据点
type MetricPoint struct {
	Timestamp time.Time          `json:"timestamp"` // 时间戳
	Values    map[string]float64 `json:"values"`    // 指标值集合
	Type      string
	Labels    map[string]string
}

// System 系统结构
type System struct {
	mu sync.RWMutex

	// 基础信息
	ID        string    // 系统ID
	Name      string    // 系统名称
	Version   string    // 系统版本
	StartTime time.Time // 启动时间

	// 核心模型 - 使用 model 包的流模型
	models struct {
		yinyang *model.YinYangFlow
		wuxing  *model.WuXingFlow
		bagua   *model.BaGuaFlow
		ganzhi  *model.GanZhiFlow
		unified *model.IntegrateFlow
	}

	// 系统状态
	state struct {
		current  model.SystemState       // 当前状态
		previous model.SystemState       // 前一状态
		changes  []model.StateTransition // 状态变更历史
	}

	// 系统组件
	components struct {
		meta      *MetaSystem      // 元系统组件
		evolution *EvolutionSystem // 演化系统组件
		control   *ControlSystem   // 控制系统组件
		monitor   *MonitorSystem   // 监控系统组件
		resource  *ResourceSystem  // 资源系统组件
	}

	// 系统配置
	config SystemConfig
}

// MetaSystem 元系统组件
type MetaSystem struct {
	// 场状态
	field struct {
		state    core.FieldState   // 场状态
		quantum  core.QuantumState // 量子状态
		coupling [][]float64       // 场耦合矩阵
	}

	// 涌现状态
	emergence struct {
		patterns  []core.EmergentPattern    // 涌现模式
		active    []core.EmergentProperty   // 活跃属性
		potential []core.PotentialEmergence // 潜在涌现
	}

	// 共振状态
	resonance struct {
		state     core.ResonanceState // 共振状态
		coherence float64             // 相干度
		phase     float64             // 相位
	}
}

// EvolutionSystem 演化系统组件
type EvolutionSystem struct {
	// 当前状态
	current struct {
		level     float64        // 演化层级
		direction model.Vector3D // 演化方向
		speed     float64        // 演化速度
		energy    float64        // 演化能量
	}

	// 演化历史
	history struct {
		path    []EvolutionPoint        // 演化路径
		changes []model.StateTransition // 状态变更
		metrics []EvolutionMetrics      // 演化指标
	}
}

// ControlSystem 控制系统组件
type ControlSystem struct {
	// 状态控制
	state struct {
		manager    *model.StateManager   // 状态管理器
		validator  *model.StateValidator // 状态验证器
		transition *model.StateTransitor // 状态转换器
	}

	// 流控制
	flow struct {
		scheduler    *FlowScheduler // 调度器
		balancer     *FlowBalancer  // 平衡器
		backpressure *BackPressure  // 背压控制
	}

	// 同步控制
	sync struct {
		coordinator  *SyncCoordinator   // 同步协调器
		resolver     *ConflictResolver  // 冲突解决器
		synchronizer *StateSynchronizer // 状态同步器
	}
}

// SyncStatus 同步状态
type SyncStatus struct {
	State       string    `json:"state"`        // 同步状态(syncing/completed/failed)
	Progress    float64   `json:"progress"`     // 同步进度(0-1)
	LastSync    time.Time `json:"last_sync"`    // 最后同步时间
	NextSync    time.Time `json:"next_sync"`    // 下次同步时间
	CurrentMode SyncMode  `json:"current_mode"` // 当前同步模式
	Stats       SyncStats `json:"stats"`        // 同步统计
	Errors      []error   `json:"errors"`       // 错误记录
}

// SyncStats 同步统计信息
type SyncStats struct {
	TotalSyncs   int64         `json:"total_syncs"`    // 总同步次数
	SuccessSyncs int64         `json:"success_syncs"`  // 成功同步次数
	FailedSyncs  int64         `json:"failed_syncs"`   // 失败同步次数
	AverageTime  time.Duration `json:"average_time"`   // 平均同步时间
	LastSyncTime time.Duration `json:"last_sync_time"` // 最后同步耗时
}

// MonitorSystem 监控系统组件
type MonitorSystem struct {
	// 指标监控
	metrics struct {
		collector *MetricsCollector // 指标收集器
		storage   *MetricsStorage   // 指标存储
		analyzer  *MetricsAnalyzer  // 指标分析器
	}

	// 告警管理
	alerts struct {
		manager  *AlertManager  // 告警管理器
		handler  *AlertHandler  // 告警处理器
		notifier *AlertNotifier // 告警通知器
	}

	// 健康检查
	health struct {
		checker   *HealthChecker   // 健康检查器
		reporter  *HealthReporter  // 健康报告器
		diagnoser *HealthDiagnoser // 健康诊断器
	}
}

// ResourceSystem 资源系统组件
type ResourceSystem struct {
	// 资源池
	pool struct {
		cpu    *ResourcePool // CPU资源池
		memory *ResourcePool // 内存资源池
		energy *ResourcePool // 能量资源池
	}

	// 资源管理
	management struct {
		allocator *ResourceAllocator // 资源分配器
		scheduler *ResourceScheduler // 资源调度器
		optimizer *ResourceOptimizer // 资源优化器
	}

	// 资源监控
	monitor struct {
		collector *ResourceCollector // 资源收集器
		analyzer  *ResourceAnalyzer  // 资源分析器
		predictor *ResourcePredictor // 资源预测器
	}
}

// OptimizationStatus 优化状态
type OptimizationStatus struct {
	State      string             `json:"state"`      // 当前状态(optimizing/completed/failed)
	Progress   float64            `json:"progress"`   // 优化进度(0-1)
	Goals      OptimizationGoals  `json:"goals"`      // 优化目标
	Results    map[string]float64 `json:"results"`    // 优化结果
	StartTime  time.Time          `json:"start_time"` // 开始时间
	EndTime    time.Time          `json:"end_time"`   // 结束时间
	Iterations int                `json:"iterations"` // 迭代次数
	Error      string             `json:"error"`      // 错误信息
}

// OptimizationGoals 优化目标
type OptimizationGoals struct {
	Targets     map[string]float64    `json:"targets"`     // 目标值
	Weights     map[string]float64    `json:"weights"`     // 目标权重
	Constraints map[string]Constraint `json:"constraints"` // 约束条件
	TimeLimit   time.Duration         `json:"time_limit"`  // 时间限制
	MinGain     float64               `json:"min_gain"`    // 最小增益
}

// Constraint 约束条件
type Constraint struct {
	Min      float64  `json:"min"`      // 最小值
	Max      float64  `json:"max"`      // 最大值
	Equals   *float64 `json:"equals"`   // 相等值(可选)
	Tolerant float64  `json:"tolerant"` // 容差
}

// ModelEvent 模型事件
type ModelEvent struct {
	ID        string                 `json:"id"`         // 事件ID
	Type      string                 `json:"type"`       // 事件类型
	ModelType model.ModelType        `json:"model_type"` // 模型类型
	State     model.ModelState       `json:"state"`      // 模型状态
	Changes   []StateChange          `json:"changes"`    // 状态变更
	Timestamp time.Time              `json:"timestamp"`  // 发生时间
	Details   map[string]interface{} `json:"details"`    // 详细信息
}

// MonitorMetrics 监控系统指标
type MonitorMetrics struct {
	// 基础指标
	Basic struct {
		ActiveHandlers int     // 活跃处理器数
		QueueLength    int     // 队列长度
		ErrorRate      float64 // 错误率
		Uptime         float64 // 运行时间(秒)
	}

	// 性能指标
	Performance struct {
		AverageLatency   float64 // 平均延迟
		ProcessingRate   float64 // 处理速率
		ResourceUsage    float64 // 资源使用率
		ConcurrencyLevel int     // 并发级别
	}

	// 状态指标
	Status struct {
		TotalAlerts    int64  // 总告警数
		HandledAlerts  int64  // 已处理告警数
		PendingAlerts  int    // 待处理告警数
		LastUpdateTime string // 最后更新时间
	}

	// 历史记录
	History []MetricPoint // 历史指标点
}

// ----------------------------------------------------------
// ToModelSystemState converts to model.SystemState
func ToModelSystemState(s *model.SystemState) *model.SystemState {
	if s == nil {
		return nil
	}
	return &model.SystemState{
		Energy:     s.Energy,
		Entropy:    s.Properties["entropy"].(float64),
		Harmony:    s.Properties["harmony"].(float64),
		Balance:    s.Properties["balance"].(float64),
		Phase:      model.Phase(s.Properties["phase"].(int)),
		Timestamp:  s.Timestamp,
		Properties: s.Properties,
	}
}

// FromModelSystemState creates SystemState from model.SystemState
func FromModelSystemState(s *model.SystemState) *SystemState {
	if s == nil {
		return nil
	}

	state := &SystemState{
		Energy:     s.Energy,
		Properties: s.Properties,
		Timestamp:  s.Timestamp,
	}

	state.Properties["entropy"] = s.Entropy
	state.Properties["harmony"] = s.Harmony
	state.Properties["balance"] = s.Balance
	state.Properties["phase"] = int(s.Phase)

	return state
}
