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
    GetMetrics() SystemMetrics
    GetHealth() HealthStatus

    // 演化控制
    Evolve(params EvolutionParams) error
    Adapt(params AdaptationParams) error
    Synchronize(params SyncParams) error
    Optimize(params OptimizationParams) error

    // 资源管理
    AllocateResources(req ResourceRequest) error
    ReleaseResources(id string) error
    GetResourceStatus() ResourceStatus
}

// SystemStatus 系统状态
type SystemStatus struct {
    State       SystemState  `json:"state"`
    Time        time.Time    `json:"time"`
    Components  []Component  `json:"components"`
    Resources   ResourcePool `json:"resources"`
    Metrics     Metrics     `json:"metrics"`
}

// SystemState 系统状态类型
type SystemState uint8

const (
    StateInitializing SystemState = iota
    StateRunning
    StatePaused
    StateError
    StateShutdown
)

// Component 组件信息
type Component struct {
    ID          string          `json:"id"`
    Type        ComponentType   `json:"type"`
    State       ComponentState  `json:"state"`
    Health      float64         `json:"health"`
    LastUpdate  time.Time       `json:"last_update"`
}

// ComponentType 组件类型
type ComponentType string

const (
    CompEvolution      ComponentType = "evolution"
    CompAdaptation    ComponentType = "adaptation"
    CompSynchronization ComponentType = "synchronization"
    CompOptimization   ComponentType = "optimization"
    CompEmergence      ComponentType = "emergence"
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
    Target     EvolutionTarget `json:"target"`
    Constraints []Constraint   `json:"constraints"`
    Speed      float64        `json:"speed"`
    Threshold  float64        `json:"threshold"`
}

// EvolutionTarget 演化目标
type EvolutionTarget struct {
    Type      TargetType     `json:"type"`
    Value     interface{}    `json:"value"`
    Priority  int           `json:"priority"`
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
    Mode        AdaptMode      `json:"mode"`
    Patterns    []Pattern      `json:"patterns"`
    LearnRate   float64        `json:"learn_rate"`
    Memory      MemoryConfig   `json:"memory"`
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
    Target      SyncTarget     `json:"target"`
    Coupling    [][]float64    `json:"coupling"`
    Frequency   float64        `json:"frequency"`
    Phase       float64        `json:"phase"`
}

// SyncTarget 同步目标
type SyncTarget struct {
    Systems    []string       `json:"systems"`
    Mode       SyncMode      `json:"mode"`
    Threshold  float64       `json:"threshold"`
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
    Objective    Objective      `json:"objective"`
    Constraints  []Constraint   `json:"constraints"`
    Method       OptimMethod    `json:"method"`
    Config       OptimConfig    `json:"config"`
}

// Objective 优化目标
type Objective struct {
    Function    string         `json:"function"`
    Weights     []float64      `json:"weights"`
    Bounds      [][]float64    `json:"bounds"`
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
    Score       float64        `json:"score"`
    Issues      []Issue        `json:"issues"`
    LastCheck   time.Time      `json:"last_check"`
}

// Issue 问题记录
type Issue struct {
    Type        IssueType      `json:"type"`
    Severity    IssueSeverity  `json:"severity"`
    Message     string         `json:"message"`
    Time        time.Time      `json:"time"`
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
    Code    ErrorCode   `json:"code"`
    Message string     `json:"message"`
    Details string     `json:"details"`
    Time    time.Time  `json:"time"`
}

// ErrorCode 错误代码
type ErrorCode uint32

const (
    ErrNone           ErrorCode = 0
    ErrInitialize     ErrorCode = 1000
    ErrEvolution      ErrorCode = 2000
    ErrAdaptation     ErrorCode = 3000
    ErrSynchronization ErrorCode = 4000
    ErrOptimization   ErrorCode = 5000
    ErrResource       ErrorCode = 6000
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
    Type      EventType    `json:"type"`
    Source    string      `json:"source"`
    Target    string      `json:"target"`
    Data      interface{} `json:"data"`
    Time      time.Time   `json:"time"`
}

// Config 配置接口
type Config interface {
    Validate() error
    Apply(interface{}) error
    GetValue(key string) interface{}
    SetValue(key string, value interface{}) error
}
