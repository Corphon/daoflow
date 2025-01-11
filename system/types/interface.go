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
    
    // 状态管理
    GetState() SystemState
    GetStatus() SystemStatus
    GetMetrics() MetricsData
    GetHealth() HealthStatus
    
    // 演化控制
    Evolve(params EvolutionParams) error
    Adapt(params AdaptationParams) error
    Synchronize(params SyncParams) error
    Optimize(params OptimizationParams) error
    
    // 资源管理
    AllocateResources(req ResourceReq) error
    ReleaseResources(id string) error
    GetResourceStatus() ResourceStats
    
    // 事件处理
    HandleEvent(event SystemEvent) error
    Subscribe(eventType EventType, handler EventHandler) error
    Unsubscribe(eventType EventType, handler EventHandler) error
}

// MetaSystemInterface 元系统接口
type MetaSystemInterface interface {
    // 场操作
    InitializeField(params FieldParams) error
    UpdateField(params FieldParams) error
    GetFieldState() FieldState
    
    // 量子场操作
    InitializeQuantumField() error
    UpdateQuantumState(state QuantumState) error
    GetQuantumState() QuantumState
    MeasureQuantumState() (float64, error)
    
    // 涌现操作
    DetectEmergence() []EmergentProperty
    PredictEmergence() []PotentialEmergence
    HandleEmergence(property EmergentProperty) error
    AnalyzePattern(pattern EmergentPattern) error
    
    // 共振操作
    InitializeResonance(params ResonanceParams) error
    UpdateResonance(params ResonanceParams) error
    GetResonanceState() ResonanceState
    MaintainCoherence(threshold float64) error
}

// EvolutionInterface 演化接口
type EvolutionInterface interface {
    // 演化控制
    SetEvolutionParams(params EvolutionParams) error
    GetEvolutionStatus() EvolutionMetrics
    AdjustEvolution(adjustment float64) error
    
    // 路径管理
    PlanEvolutionPath(target SystemState) ([]EvolutionPoint, error)
    ValidateEvolutionPath(path []EvolutionPoint) error
    ExecuteEvolutionStep() error
    RollbackEvolution() error
    
    // 能量管理
    CalculateEvolutionEnergy() float64
    OptimizeEnergyUsage() error
    PredictEnergyNeeds() float64
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
    // 指标管理
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

// EventHandler 事件处理器接口
type EventHandler interface {
    HandleEvent(event SystemEvent) error
    GetHandlerID() string
    GetEventTypes() []EventType
}

// ConfigManager 配置管理接口
type ConfigManager interface {
    // 配置操作
    LoadConfig(path string) error
    SaveConfig(path string) error
    ValidateConfig() error
    
    // 配置访问
    GetValue(key string) interface{}
    SetValue(key string, value interface{}) error
    GetSection(section string) interface{}
    
    // 配置监控
    WatchConfig(handler ConfigHandler) error
    StopWatch(handler ConfigHandler) error
}

// ConfigHandler 配置处理器接口
type ConfigHandler interface {
    OnConfigChange(old, new interface{}) error
    GetHandlerID() string
}

// StateManager 状态管理接口
type StateManager interface {
    // 状态操作
    GetCurrentState() SystemState
    SetState(state SystemState) error
    ValidateState(state SystemState) error
    
    // 状态转换
    TransitionTo(target SystemState) error
    RollbackState() error
    GetStateHistory() []StateTransition
    
    // 状态监控
    WatchState(handler StateHandler) error
    StopWatch(handler StateHandler) error
}

// StateHandler 状态处理器接口
type StateHandler interface {
    OnStateChange(old, new SystemState) error
    GetHandlerID() string
}

// MetricsCollector 指标收集器接口
type MetricsCollector interface {
    // 指标收集
    Collect() MetricsData
    Process(data MetricsData) error
    Store(data MetricsData) error
    
    // 指标查询
    Query(filter MetricsFilter) []MetricsData
    GetLatest() MetricsData
    GetHistory(duration time.Duration) []MetricsData
    
    // 配置管理
    Configure(config MetricsConfig) error
    GetConfig() MetricsConfig
}

// AlertManager 告警管理器接口
type AlertManager interface {
    // 告警处理
    CheckAlerts() []Alert
    HandleAlert(alert Alert) error
    ClearAlert(id string) error
    
    // 告警配置
    ConfigureAlerts(config AlertConfig) error
    GetAlertConfig() AlertConfig
    
    // 告警查询
    GetActiveAlerts() []Alert
    GetAlertHistory(filter AlertFilter) []Alert
}

// HealthChecker 健康检查接口
type HealthChecker interface {
    // 健康检查
    Check() HealthStatus
    DiagnoseIssue(issue Issue) error
    GetHealthMetrics() HealthMetrics
    
    // 配置
    SetHealthChecks(checks []HealthCheck) error
    ConfigureThresholds(thresholds map[string]float64) error
}

// ResourcePredictor 资源预测接口
type ResourcePredictor interface {
    // 预测
    PredictUsage(duration time.Duration) ResourcePrediction
    AnalyzeTrends() []ResourceTrend
    OptimizePrediction() error
    
    // 配置
    ConfigurePredictor(config PredictorConfig) error
    GetAccuracy() float64
}
