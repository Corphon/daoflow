// system/types/interface.go

package types

import (
	"context"
	"time"

	"github.com/Corphon/daoflow/core"
	"github.com/Corphon/daoflow/model"
)

// SyncParams 同步参数
type SyncParams struct {
	// 同步目标
	Target      string            // 同步目标标识
	TargetState model.SystemState // 目标状态

	// 同步配置
	Mode       SyncMode      // 同步模式
	Timeout    time.Duration // 同步超时
	RetryCount int           // 重试次数

	// 同步选项
	Options  map[string]interface{} // 自定义选项
	Priority Priority               // 同步优先级

	// 回调函数
	OnProgress func(float64) // 进度回调
	OnComplete func(error)   // 完成回调
}

// OptimizationParams 优化参数
type OptimizationParams struct {
	// 优化目标
	Targets     []string           // 优化目标列表
	Constraints map[string]float64 // 优化约束条件

	// 优化配置
	Strategy      string  // 优化策略
	MaxIterations int     // 最大迭代次数
	Tolerance     float64 // 收敛容差

	// 优化选项
	Options  map[string]interface{} // 自定义选项
	Priority Priority               // 优化优先级

	// 回调函数
	OnProgress func(float64) // 进度回调
	OnComplete func(error)   // 完成回调
}

// SystemEvent 系统事件
type SystemEvent struct {
	// 事件基本信息
	ID        string    // 事件ID
	Type      EventType // 事件类型
	Source    string    // 事件源
	Timestamp time.Time // 事件时间

	// 事件内容
	Message  string            // 事件消息
	Data     interface{}       // 事件数据
	Metadata map[string]string // 事件元数据

	// 事件处理
	Priority Priority // 事件优先级
	Handled  bool     // 是否已处理
	Error    error    // 处理错误
}

// EventHandler 事件处理器
type EventHandler interface {
	// 处理器信息
	GetHandlerID() string
	GetEventTypes() []EventType
	GetPriority() Priority

	// 事件处理
	HandleEvent(event SystemEvent) error
	ShouldHandle(event SystemEvent) bool

	// 生命周期管理
	Initialize() error
	Shutdown() error
}

// 添加事件处理器基础实现
type BaseEventHandler struct {
	ID       string
	Types    []EventType
	Priority Priority
}

func (h *BaseEventHandler) GetHandlerID() string {
	return h.ID
}

func (h *BaseEventHandler) GetEventTypes() []EventType {
	return h.Types
}

func (h *BaseEventHandler) GetPriority() Priority {
	return h.Priority
}

func (h *BaseEventHandler) Initialize() error {
	return nil
}

func (h *BaseEventHandler) Shutdown() error {
	return nil
}

// EventHandlerFunc 事件处理函数类型
type EventHandlerFunc func(SystemEvent) error

// 包装事件处理函数为处理器接口
type funcEventHandler struct {
	BaseEventHandler
	handler EventHandlerFunc
}

// NewEventHandler 创建事件处理器
func NewEventHandler(id string, types []EventType, priority Priority, handler EventHandlerFunc) EventHandler {
	return &funcEventHandler{
		BaseEventHandler: BaseEventHandler{
			ID:       id,
			Types:    types,
			Priority: priority,
		},
		handler: handler,
	}
}

func (h *funcEventHandler) HandleEvent(event SystemEvent) error {
	if h.handler != nil {
		return h.handler(event)
	}
	return nil
}

func (h *funcEventHandler) ShouldHandle(event SystemEvent) bool {
	for _, t := range h.Types {
		if t == event.Type {
			return true
		}
	}
	return false
}

// 添加事件总线
type EventBus interface {
	// 发布事件
	Publish(event SystemEvent) error

	// 订阅管理
	Subscribe(types []EventType, handler EventHandler) error
	Unsubscribe(handler EventHandler) error

	// 处理器管理
	AddHandler(handler EventHandler) error
	RemoveHandler(handlerID string) error

	// 状态查询
	GetHandlers() []EventHandler
	GetSubscriptions(eventType EventType) []EventHandler
}

// SystemInterface 系统核心接口
type SystemInterface interface {
	// 生命周期管理
	Initialize(ctx context.Context) error
	Start() error
	Stop() error
	Reset() error
	Shutdown(ctx context.Context) error

	// 状态访问
	GetState() model.SystemState
	GetModelState() model.ModelState
	GetMetrics() model.ModelMetrics

	// 系统操作
	Transform(pattern model.TransformPattern) error
	Synchronize(params SyncParams) error
	Optimize(params OptimizationParams) error

	// 事件处理
	HandleEvent(event SystemEvent) error
	Subscribe(eventType EventType, handler EventHandler) error
	Unsubscribe(eventType EventType, handler EventHandler) error
}

// MetaSystemInterface 元系统接口
type MetaSystemInterface interface {
	// 场操作 - 使用 model 的场状态
	InitializeField(params core.FieldParams) error
	UpdateField(state core.FieldState) error
	GetFieldState() core.FieldState

	// 量子操作 - 使用 model 的量子状态
	InitializeQuantumField() error
	UpdateQuantumState(state core.QuantumState) error
	GetQuantumState() core.QuantumState
	MeasureQuantumState() (float64, error)

	// 涌现操作
	DetectEmergence() []core.EmergentPattern
	PredictEmergence() []core.PotentialPattern
	HandleEmergence(property core.EmergentProperty) error

	// 共振操作
	InitializeResonance(params core.ResonanceParams) error
	UpdateResonance(params core.ResonanceParams) error
	GetResonanceState() core.ResonanceState
	MaintainCoherence(threshold float64) error
}

// EvolutionInterface 演化接口
type EvolutionInterface interface {
	// 演化控制 - 与 model 的演化对齐
	SetEvolutionParams(params EvolutionParams) error
	GetEvolutionMetrics() model.ModelMetrics
	AdjustEvolution(delta float64) error

	// 路径管理
	PlanEvolutionPath(target model.SystemState) ([]EvolutionPoint, error)
	ValidateEvolutionPath(path []EvolutionPoint) error
	ExecuteEvolutionStep() error
	RollbackEvolution() error
}

// ResourceInterface 资源管理接口
type ResourceInterface interface {
	// 资源分配
	AllocateResource(req ResourceReq) error
	ReleaseResource(id string) error
	ReserveResources(reqs []ResourceReq) error

	// 资源监控
	GetResourceStats() ResourceStats
	MonitorResource(id string) error
	SetResourceLimits(limits ResourceLimits) error

	// 资源优化
	OptimizeResourceUsage() error
	BalanceResources() error
	PredictResourceNeeds() ResourcePrediction
}

// MonitorInterface 监控接口
type MonitorInterface interface {
	// 指标管理 - 包含 model 的指标
	CollectMetrics() model.MetricsData
	ProcessMetrics(data model.MetricsData) error
	StoreMetrics(data model.MetricsData) error

	// 告警管理
	CheckAlerts() []Alert
	HandleAlert(alert Alert) error
	ConfigureAlerts(config AlertConfig) error

	// 健康检查
	HealthCheck() HealthStatus
	DiagnoseIssue(issue Issue) error
	ReportStatus() SystemStatus
}

// ConfigInterface 配置接口
type ConfigInterface interface {
	// 配置管理
	LoadConfig(path string) error
	SaveConfig(path string) error
	ValidateConfig() error

	// 模型配置 - 与 model 配置集成
	GetModelConfig() model.ModelConfig
	UpdateModelConfig(config model.ModelConfig) error

	// 系统配置
	GetSystemConfig() SystemConfig
	UpdateSystemConfig(config SystemConfig) error
}

// StateHandler 状态处理接口
type StateHandler interface {
	// 状态变更处理
	OnStateChange(old, new model.SystemState) error
	OnModelStateChange(old, new model.ModelState) error

	// 处理器标识
	GetHandlerID() string
	GetPriority() Priority
}

// MetricsCollector 指标收集接口
type MetricsCollector interface {
	// 指标收集 - 包含 model 指标
	CollectModelMetrics() model.ModelMetrics
	CollectSystemMetrics() model.MetricsData

	// 指标处理
	ProcessMetrics(modelMetrics model.ModelMetrics, sysMetrics model.MetricsData) error
	StoreMetrics(data model.MetricsData) error

	// 指标查询
	Query(filter model.MetricsFilter) []model.MetricsData
	GetHistory(duration time.Duration) []model.MetricsData
}

// SyncController 同步控制接口
type SyncController interface {
	// 同步控制
	SyncModelState(state model.ModelState) error
	SyncSystemState(state model.SystemState) error

	// 同步配置
	SetSyncMode(mode SyncMode) error
	GetSyncStatus() SyncStatus
}

// OptimizationController 优化控制接口
type OptimizationController interface {
	// 优化控制
	OptimizeModelParams(params model.ModelConfig) error
	OptimizeSystemParams(params SystemConfig) error

	// 优化状态
	GetOptimizationStatus() OptimizationStatus
	SetOptimizationGoals(goals OptimizationGoals) error
}

// EventProcessor 事件处理接口
type EventProcessor interface {
	// 事件处理
	ProcessModelEvent(event model.ModelEvent) error
	ProcessSystemEvent(event SystemEvent) error

	// 事件订阅
	Subscribe(eventType EventType, handler EventHandler) error
	Unsubscribe(eventType EventType, handler EventHandler) error
}
