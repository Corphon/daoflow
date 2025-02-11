// system/types/resource_types.go

package types

import (
	"sync"
	"time"
)

// ResourceReq 资源请求
type ResourceReq struct {
	ID          string             // 请求ID
	Type        string             // 资源类型
	Amount      float64            // 请求数量
	Priority    int                // 优先级
	Constraints map[string]float64 // 约束条件
	Deadline    time.Time          // 截止时间
}

// ResourceStats 资源统计
type ResourceStats struct {
	Usage     map[string]float64 // 资源使用情况
	Available map[string]float64 // 可用资源
	Reserved  map[string]float64 // 预留资源
	Total     map[string]float64 // 总资源量
	UpdatedAt time.Time          // 更新时间
}

// ResourceLimits 资源限制
type ResourceLimits struct {
	CPU      float64 // CPU限制
	Memory   int64   // 内存限制
	Storage  int64   // 存储限制
	Network  float64 // 网络限制
	Energy   float64 // 能量限制
	MaxTasks int     // 最大任务数
}

// ResourcePrediction 资源预测
type ResourcePrediction struct {
	Estimates  map[string][]float64 // 预测值
	Confidence float64              // 置信度
	TimeRange  time.Duration        // 预测时间范围
	Factors    map[string]float64   // 影响因子
}

// ResourcePool 资源池
type ResourcePool struct {
	mu sync.RWMutex

	// 资源配置
	config struct {
		Type     string  // 资源类型
		Capacity float64 // 总容量
		MinLimit float64 // 最小限制
		MaxLimit float64 // 最大限制
	}

	// 资源状态
	state struct {
		available  float64   // 可用量
		allocated  float64   // 已分配量
		reserved   float64   // 预留量
		lastUpdate time.Time // 最后更新
	}
}

// ResourceAllocator 资源分配器
type ResourceAllocator struct {
	mu sync.RWMutex

	// 分配配置
	config struct {
		MaxBatch   int           // 最大批量
		Timeout    time.Duration // 分配超时
		RetryLimit int           // 重试限制
	}

	// 分配状态
	state struct {
		requests    map[string]ResourceRequest    // 请求队列
		allocations map[string]ResourceAllocation // 分配记录
		pending     []string                      // 等待处理
	}
}

// ResourceScheduler 资源调度器
type ResourceScheduler struct {
	mu sync.RWMutex

	// 调度配置
	config struct {
		Interval time.Duration  // 调度间隔
		Strategy string         // 调度策略
		Priority map[string]int // 优先级
	}

	// 调度状态
	state struct {
		tasks   map[string]ScheduleTask // 调度任务
		running []string                // 运行中
		waiting []string                // 等待中
	}
}

// ResourceOptimizer 资源优化器
type ResourceOptimizer struct {
	mu sync.RWMutex

	// 优化配置
	config struct {
		Goals       []string           // 优化目标
		Constraints map[string]float64 // 约束条件
		Window      time.Duration      // 优化窗口
	}

	// 优化状态
	state struct {
		plans   []OptimizationPlan   // 优化计划
		metrics []OptimizationMetric // 优化指标
		history []OptimizationResult // 优化历史
	}
}

// ResourceCollector 资源收集器
type ResourceCollector struct {
	mu sync.RWMutex

	// 收集配置
	config struct {
		Interval  time.Duration // 收集间隔
		Metrics   []string      // 收集指标
		BatchSize int           // 批量大小
	}

	// 收集状态
	state struct {
		current     ResourceMetrics   // 当前指标
		history     []ResourceMetrics // 历史指标
		lastCollect time.Time         // 最后收集
	}
}

// ResourceAnalyzer 资源分析器
type ResourceAnalyzer struct {
	mu sync.RWMutex

	// 分析配置
	config struct {
		Rules      []AnalysisRule     // 分析规则
		Thresholds map[string]float64 // 阈值设置
		Window     time.Duration      // 分析窗口
	}

	// 分析状态
	state struct {
		results   []AnalysisResult  // 分析结果
		patterns  []ResourcePattern // 资源模式
		anomalies []ResourceAnomaly // 异常情况
	}
}

// ResourcePredictor 资源预测器
type ResourcePredictor struct {
	mu sync.RWMutex

	// 预测配置
	config struct {
		Horizon    time.Duration // 预测范围
		Models     []string      // 预测模型
		Confidence float64       // 置信度
	}

	// 预测状态
	state struct {
		predictions []ResourcePrediction // 预测结果
		accuracy    map[string]float64   // 预测准确度
		lastUpdate  time.Time            // 最后更新
	}
}

// ResourceRequest 资源请求
type ResourceRequest struct {
	ID        string // 请求ID
	Type      string // 资源类型
	Priority  int    // 优先级
	Resources struct {
		CPU     float64 // CPU需求
		Memory  int64   // 内存需求
		Storage int64   // 存储需求
		Network float64 // 网络需求
		Energy  float64 // 能量需求
	}
	Constraints map[string]float64 // 约束条件
	Deadline    time.Time          // 截止时间
	Status      string             // 请求状态
}

// ResourceAllocation 资源分配
type ResourceAllocation struct {
	RequestID  string // 请求ID
	ResourceID string // 资源ID
	Resources  struct {
		CPU     float64 // CPU分配
		Memory  int64   // 内存分配
		Storage int64   // 存储分配
		Network float64 // 网络分配
		Energy  float64 // 能量分配
	}
	StartTime time.Time // 分配时间
	EndTime   time.Time // 释放时间
	Status    string    // 分配状态
}

// ScheduleTask 调度任务
type ScheduleTask struct {
	ID           string              // 任务ID
	Type         string              // 任务类型
	Priority     int                 // 优先级
	Resources    ResourceRequirement // 资源需求
	Status       string              // 任务状态
	Progress     float64             // 进度
	StartTime    time.Time           // 开始时间
	EndTime      time.Time           // 结束时间
	Dependencies []string            // 依赖任务
	Metrics      TaskMetrics         // 任务指标
}

// ResourceRequirement 资源需求
type ResourceRequirement struct {
	CPU     float64 // CPU需求
	Memory  int64   // 内存需求
	Storage int64   // 存储需求
	Network float64 // 网络需求
	Energy  float64 // 能量需求
}

// TaskMetrics 任务指标
type TaskMetrics struct {
	ProcessingTime time.Duration      // 处理时间
	WaitingTime    time.Duration      // 等待时间
	ResourceUsage  map[string]float64 // 资源使用
	ErrorCount     int                // 错误次数
}

// OptimizationPlan 优化计划
type OptimizationPlan struct {
	ID        string             // 计划ID
	Goals     []OptimizationGoal // 优化目标
	Steps     []OptimizationStep // 优化步骤
	Resources []string           // 涉及资源
	StartTime time.Time          // 开始时间
	EndTime   time.Time          // 结束时间
	Status    string             // 计划状态
}

// OptimizationGoal 优化目标
type OptimizationGoal struct {
	Type      string  // 目标类型
	Target    float64 // 目标值
	Weight    float64 // 权重
	Threshold float64 // 阈值
}

// OptimizationStep 优化步骤
type OptimizationStep struct {
	ID         string             // 步骤ID
	Action     string             // 优化动作
	Resource   string             // 目标资源
	Parameters map[string]float64 // 参数
	Status     string             // 执行状态
}

// OptimizationMetric 优化指标
type OptimizationMetric struct {
	Type      string    // 指标类型
	Value     float64   // 当前值
	Target    float64   // 目标值
	Trend     float64   // 变化趋势
	UpdatedAt time.Time // 更新时间
}

// OptimizationResult 优化结果
type OptimizationResult struct {
	PlanID    string               // 计划ID
	Success   bool                 // 是否成功
	Metrics   []OptimizationMetric // 优化指标
	Changes   []ResourceChange     // 资源变更
	StartTime time.Time            // 开始时间
	EndTime   time.Time            // 结束时间
	Duration  time.Duration        // 执行时长
}

// ResourceChange 资源变更
type ResourceChange struct {
	ResourceID string    // 资源ID
	Type       string    // 变更类型
	OldValue   float64   // 原始值
	NewValue   float64   // 新值
	Timestamp  time.Time // 变更时间
}

// ResourceMetrics 资源指标
type ResourceMetrics struct {
	// 基础指标
	ID        string    // 指标ID
	Timestamp time.Time // 采集时间

	// CPU指标
	CPU struct {
		Usage       float64 // 使用率
		Load        float64 // 负载
		Temperature float64 // 温度
		Frequency   float64 // 频率
	}

	// 内存指标
	Memory struct {
		Total    int64 // 总量
		Used     int64 // 已用
		Free     int64 // 空闲
		Cached   int64 // 缓存
		SwapUsed int64 // 交换使用
	}

	// 存储指标
	Storage struct {
		Total     int64         // 总量
		Used      int64         // 已用
		Free      int64         // 空闲
		IOps      int64         // IO操作数
		IOLatency time.Duration // IO延迟
	}

	// 网络指标
	Network struct {
		BytesIn  int64         // 入流量
		BytesOut int64         // 出流量
		Packets  int64         // 包数量
		Errors   int64         // 错误数
		Latency  time.Duration // 网络延迟
	}

	// 资源池指标
	Pools map[string]struct {
		Capacity    float64 // 容量
		Available   float64 // 可用
		Reserved    float64 // 预留
		Utilization float64 // 利用率
	}

	// 统计信息
	Statistics struct {
		SampleCount int64   // 采样数
		ErrorCount  int64   // 错误数
		MaxUsage    float64 // 最大使用
		MinUsage    float64 // 最小使用
		AvgUsage    float64 // 平均使用
	}
}

// AnalysisRule 分析规则
type AnalysisRule struct {
	ID         string                 // 规则ID
	Name       string                 // 规则名称
	Target     string                 // 分析目标
	Type       string                 // 规则类型
	Condition  string                 // 条件表达式
	Parameters map[string]interface{} // 参数
	Priority   int                    // 优先级
}

// ResourcePattern 资源模式
type ResourcePattern struct {
	ID         string    // 模式ID
	Name       string    // 模式名称
	Type       string    // 模式类型
	Metrics    []string  // 相关指标
	Values     []float64 // 特征值
	Confidence float64   // 置信度
	DetectedAt time.Time // 检测时间
}

// ResourceAnomaly 资源异常
type ResourceAnomaly struct {
	ID          string    // 异常ID
	Type        string    // 异常类型
	Source      string    // 异常源
	Level       string    // 异常级别
	Message     string    // 异常描述
	MetricName  string    // 指标名
	MetricValue float64   // 指标值
	Threshold   float64   // 阈值
	DetectedAt  time.Time // 检测时间
	ResolvedAt  time.Time // 解决时间
}
