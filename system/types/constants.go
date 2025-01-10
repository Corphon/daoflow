// system/types/constants.go

package types

// SystemLayer 系统层级常量
type SystemLayer uint8

const (
	LayerNone      SystemLayer = iota
	LayerMeta                  // 元系统层
	LayerEvolution             // 演化系统层
	LayerStructure             // 结构系统层
	LayerSync                  // 同步系统层
	LayerMax
)

// SystemState 系统状态常量
type SystemState uint8

const (
	StateNone         SystemState = iota
	StateInitializing             // 初始化中
	StateRunning                  // 运行中
	StatePaused                   // 暂停
	StateStopped                  // 停止
	StateError                    // 错误
	StateMax
)

// MetricType 指标类型常量
type MetricType uint8

const (
	MetricNone      MetricType = iota
	MetricEnergy               // 能量指标
	MetricField                // 场指标
	MetricQuantum              // 量子指标
	MetricEmergence            // 涌现指标
	MetricMax
)

// ThresholdConstants 阈值常量
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
	MinInterval    = 1  // 最小间隔时间(秒)
	MaxInterval    = 60 // 最大间隔时间(秒)
)

// CapacityConstants 容量常量
const (
	MinCapacity     = 1
	MaxCapacity     = 1000
	DefaultCapacity = 100

	// 缓冲区大小
	MinBufferSize     = 16
	DefaultBufferSize = 256
	MaxBufferSize     = 4096
)

// ScalingFactors 缩放因子常量
const (
	MinScaleFactor     = 0.1
	MaxScaleFactor     = 10.0
	DefaultScaleFactor = 1.0

	// 衰减因子
	DecayFactor = 0.95
	// 增长因子
	GrowthFactor = 1.05
)

// Priorities 优先级常量
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
	PatternNone     PatternType = iota
	PatternLinear               // 线性模式
	PatternCyclic               // 循环模式
	PatternSpiral               // 螺旋模式
	PatternEmergent             // 涌现模式
	PatternMax
)

// DimensionConstants 维度常量
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

// ConfigurationDefaults 配置默认值
const (
	DefaultConfigPath = "config/system.yaml"
	DefaultLogPath    = "log/system.log"
	DefaultDataPath   = "data/system"

	// 监控配置
	DefaultMetricsInterval = 10 // 秒
	DefaultReportInterval  = 60 // 秒
	DefaultRetentionDays   = 30 // 天
)

// ValidationConstants 验证常量
const (
	MinNameLength = 3
	MaxNameLength = 64

	MinDescLength = 0
	MaxDescLength = 256
)

// StatusCodes 状态码常量
type StatusCode uint16

const (
	StatusSuccess StatusCode = iota
	StatusWarning
	StatusError
	StatusFatal
)

// FeatureFlags 功能标志常量
const (
	FeatureMetrics    = "metrics"
	FeatureTracing    = "tracing"
	FeatureDebugging  = "debugging"
	FeatureMonitoring = "monitoring"
)

// 系统限制常量
const (
	MaxGoroutines  = 10000 // 最大协程数
	MaxConnections = 1000  // 最大连接数
	MaxQueueSize   = 10000 // 最大队列大小
	MaxRetries     = 3     // 最大重试次数
)
