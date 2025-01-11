// system/types/interface.go

package types

import (
    "context"
    "time"
    
    "github.com/Corphon/daoflow/model"
)

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
    GetMetrics() MetricsData
    
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
    InitializeField(params model.FieldParams) error
    UpdateField(state model.FieldState) error
    GetFieldState() model.FieldState
    
    // 量子操作 - 使用 model 的量子状态
    InitializeQuantumField() error
    UpdateQuantumState(state model.QuantumState) error
    GetQuantumState() model.QuantumState
    MeasureQuantumState() (float64, error)
    
    // 涌现操作
    DetectEmergence() []EmergentProperty
    PredictEmergence() []PotentialEmergence
    HandleEmergence(property EmergentProperty) error
    
    // 共振操作
    InitializeResonance(params ResonanceParams) error
    UpdateResonance(params ResonanceParams) error
    GetResonanceState() ResonanceState
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
    CollectMetrics() MetricsData
    ProcessMetrics(data MetricsData) error
    StoreMetrics(data MetricsData) error
    
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
    CollectSystemMetrics() MetricsData
    
    // 指标处理
    ProcessMetrics(modelMetrics model.ModelMetrics, sysMetrics MetricsData) error
    StoreMetrics(data MetricsData) error
    
    // 指标查询
    Query(filter MetricsFilter) []MetricsData
    GetHistory(duration time.Duration) []MetricsData
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
