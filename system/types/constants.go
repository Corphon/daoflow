// system/types/constants.go

package types

import (
    "github.com/Corphon/daoflow/model"
)

// 系统级阈值常量
const (
    // 历史记录限制
    MaxHistoryLength    = 1000  // 通用历史记录最大长度
    MaxMetricsHistory  = 1000  // 指标历史最大长度
    MaxObservations   = 1000  // 观察记录最大长度
    
    // 时间相关常量
    DefaultTimeWindow  = 24 * time.Hour  // 默认时间窗口
    MinTimeWindow     = time.Minute     // 最小时间窗口
    MaxTimeWindow     = 7 * 24 * time.Hour // 最大时间窗口
    
    // 阈值常量
    DefaultThreshold  = 0.75  // 默认阈值
    MinThreshold     = 0.1   // 最小阈值
    MaxThreshold     = 0.99  // 最大阈值
    
    // 容量相关常量
    DefaultBatchSize  = 100   // 默认批处理大小
    MaxBatchSize     = 1000  // 最大批处理大小
    DefaultQueueSize = 1000  // 默认队列大小
    MaxQueueSize    = 10000 // 最大队列大小
)

// 复用 model 包中的类型
type (
    SystemState = model.SystemState
    ModelType = model.ModelType
    Phase = model.Phase
    Nature = model.Nature
    Element = model.Element
    TransformPattern = model.TransformPattern
)

// SystemLayer 系统层级常量
type SystemLayer uint8

const (
    LayerNone SystemLayer = iota
    LayerMeta             // 元系统层
    LayerEvolution        // 演化系统层
    LayerControl         // 控制系统层
    LayerResource        // 资源系统层
    LayerMonitor         // 监控系统层
    LayerMax
)

// ComponentState 组件状态常量
type ComponentState uint8

const (
    CompStateActive ComponentState = iota
    CompStateInactive
    CompStateError
)

// MetricType 指标类型常量 - 扩展模型层的指标类型
type MetricType uint8

const (
    MetricNone MetricType = iota
    MetricSystem         // 系统指标
    MetricProcess        // 处理指标
    MetricResource       // 资源指标
    MetricPerformance    // 性能指标
    MetricSecurity       // 安全指标
)

// EventType 事件类型常量
type EventType uint8

const (
    EventStateChange EventType = iota
    EventResource              // 资源事件
    EventMetric               // 指标事件
    EventAlert                // 告警事件
    EventAudit                // 审计事件
    EventSystem               // 系统事件
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

// SyncMode 同步模式常量
type SyncMode uint8

const (
    SyncImmediate SyncMode = iota  // 立即同步
    SyncBatch                      // 批量同步
    SyncPeriodic                   // 周期同步
    SyncAdaptive                   // 自适应同步
)

// 系统级阈值常量
const (
    // 系统容量
    MinSystemCapacity = 100
    MaxSystemCapacity = 10000
    DefaultCapacity   = 1000

    // 处理限制
    MaxConcurrent    = 100    // 最大并发数
    MaxQueueSize     = 10000  // 最大队列大小
    MaxBatchSize     = 1000   // 最大批处理大小
    
    // 时间限制
    MinInterval      = 100    // 最小间隔(ms)
    MaxInterval      = 60000  // 最大间隔(ms)
    DefaultInterval  = 1000   // 默认间隔(ms)
    
    // 资源限制
    MinWorkers       = 1      // 最小工作协程数
    MaxWorkers       = 100    // 最大工作协程数
    DefaultWorkers   = 10     // 默认工作协程数
)

// 监控常量
const (
    // 指标采集
    MetricsInterval   = 15    // 指标采集间隔(s)
    MetricsRetention  = 7     // 指标保留天数
    
    // 告警阈值
    AlertCriticalThreshold = 0.9  // 严重告警阈值
    AlertWarningThreshold  = 0.7  // 警告告警阈值
    AlertInfoThreshold     = 0.5  // 信息告警阈值
)

// 配置常量
const (
    // 路径配置
    DefaultConfigPath = "configs/"
    DefaultLogPath    = "logs/"
    DefaultDataPath   = "data/"
    
    // 文件配置
    MaxFileSize      = 100 << 20  // 最大文件大小(100MB)
    MaxFileAge       = 30         // 最大文件保留天数
    
    // 验证配置
    MinNameLen       = 2
    MaxNameLen       = 64
    MinDescLen       = 0
    MaxDescLen       = 256
)
