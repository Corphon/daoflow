// system/types/constants.go

package types

// SystemLayer 系统层级常量
type SystemLayer uint8

const (
    LayerNone SystemLayer = iota
    LayerMeta             // 元系统层
    LayerEvolution        // 演化系统层
    LayerStructure        // 结构系统层
    LayerSync            // 同步系统层
    LayerMax
)

// SystemState 系统状态常量
type SystemState uint8

const (
    StateNone SystemState = iota
    StateInitializing    // 初始化中
    StateRunning        // 运行中
    StatePaused         // 暂停
    StateStopped        // 停止
    StateError         // 错误
    StateMax
)

// ComponentState 组件状态常量
type ComponentState uint8

const (
    CompStateActive ComponentState = iota
    CompStateInactive
    CompStateError
)

// MetricType 指标类型常量
type MetricType uint8

const (
    MetricNone MetricType = iota
    MetricEnergy          // 能量指标
    MetricField           // 场指标
    MetricQuantum         // 量子指标
    MetricEmergence       // 涌现指标
    MetricMax
)

// EventType 事件类型常量
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

// Priority 优先级常量
type Priority uint8

const (
    PriorityLowest Priority = iota
    PriorityLow
    PriorityNormal
    PriorityHigh
    PriorityHighest
)

// PatternType 模式类型常量
type PatternType uint8

const (
    PatternNone PatternType = iota
    PatternLinear          // 线性模式
    PatternCyclic          // 循环模式
    PatternSpiral          // 螺旋模式
    PatternEmergent        // 涌现模式
    PatternMax
)

// AdaptMode 适应模式常量
type AdaptMode uint8

const (
    AdaptReactive AdaptMode = iota
    AdaptProactive
    AdaptHybrid
)

// SyncMode 同步模式常量
type SyncMode uint8

const (
    SyncPhase SyncMode = iota
    SyncFrequency
    SyncEnergy
)

// OptimMethod 优化方法常量
type OptimMethod uint8

const (
    OptimGradient OptimMethod = iota
    OptimGenetic
    OptimParticleSwarm
)

// IssueType 问题类型常量
type IssueType uint8

const (
    IssueResource IssueType = iota
    IssuePerformance
    IssueSecurity
    IssueStability
)

// IssueSeverity 问题严重度常量
type IssueSeverity uint8

const (
    SeverityInfo IssueSeverity = iota
    SeverityWarning
    SeverityError
    SeverityCritical
)

// 系统阈值常量
const (
    // 能量阈值
    MinEnergy     = 0.0
    MaxEnergy     = 100.0
    DefaultEnergy = 50.0

    // 场强度阈值
    MinFieldStrength = 0.0
    MaxFieldStrength = 1.0
    DefaultStrength  = 0.5

    // 相位阈值
    MinPhase     = 0.0
    MaxPhase     = 2 * 3.14159265359 // 2π
    DefaultPhase = 0.0

    // 时间常量
    DefaultTimeout = 30 // 默认超时时间(秒)
    MinInterval   = 1  // 最小间隔时间(秒)
    MaxInterval   = 60 // 最大间隔时间(秒)
)

// 容量常量
const (
    MinCapacity     = 1
    MaxCapacity     = 1000
    DefaultCapacity = 100

    // 缓冲区大小
    MinBufferSize     = 16
    DefaultBufferSize = 256
    MaxBufferSize     = 4096
)

// 维度常量
const (
    // 空间维度
    SpaceDimension2D = 2
    SpaceDimension3D = 3
    SpaceDimension4D = 4

    // 场维度
    FieldDimensionScalar = 1
    FieldDimensionVector = 3
    FieldDimensionTensor = 9
)

// 配置默认值
const (
    DefaultConfigPath = "config/system.yaml"
    DefaultLogPath    = "log/system.log"
    DefaultDataPath   = "data/system"

    // 监控配置
    DefaultMetricsInterval = 10 // 秒
    DefaultReportInterval = 60  // 秒
    DefaultRetentionDays  = 30  // 天

    // 系统限制
    MaxGoroutines  = 10000 // 最大协程数
    MaxConnections = 1000  // 最大连接数
    MaxQueueSize   = 10000 // 最大队列大小
    MaxRetries     = 3     // 最大重试次数

    // 验证限制
    MinNameLength = 3
    MaxNameLength = 64
    MinDescLength = 0
    MaxDescLength = 256
)
