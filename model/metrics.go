// model/metrics.go

package model

import (
	"time"

	"github.com/Corphon/daoflow/core"
)

// Performance 性能指标
type Performance struct {
	Throughput float64 `json:"throughput"` // 吞吐量
	QPS        float64 `json:"qps"`        // 每秒查询数
	ErrorRate  float64 `json:"error_rate"` // 错误率
	Latency    float64 `json:"latency"`    // 延迟
}

// State 状态指标
type State struct {
	Stability   float64 `json:"stability"`   // 稳定性指标
	Transitions int     `json:"transitions"` // 状态转换次数
	Uptime      float64 `json:"uptime"`      // 运行时间
}

// Energy 能量指标
type Energy struct {
	Total    float64 `json:"total"`    // 总能量
	Average  float64 `json:"average"`  // 平均能量
	Variance float64 `json:"variance"` // 能量方差
}

// ModelMetrics 模型指标
type ModelMetrics struct {
	Basic struct {
		TotalSpans int     `json:"total_spans"`
		ErrorRate  float64 `json:"error_rate"`
		Latency    float64 `json:"latency"`
	} `json:"basic"`

	Energy      Energy      `json:"energy"`
	State       State       `json:"state"`
	Performance Performance `json:"performance"`

	Quantum *core.QuantumState `json:"quantum"` // 量子态
	Field   *core.FieldState   `json:"field"`   // 场态

	YinYang struct {
		Balance   float64 `json:"balance"`   // 阴阳平衡度
		Harmony   float64 `json:"harmony"`   // 阴阳和谐度
		Transform float64 `json:"transform"` // 转化率
	} `json:"yin_yang"`

	WuXing struct {
		Cycles  float64 `json:"cycles"`  // 五行循环强度
		Energy  float64 `json:"energy"`  // 五行能量水平
		Balance float64 `json:"balance"` // 五行平衡度
	} `json:"wu_xing"`

	BaGua struct {
		Patterns  float64 `json:"patterns"`  // 八卦模式强度
		Changes   float64 `json:"changes"`   // 变化频率
		Stability float64 `json:"stability"` // 稳定性
	} `json:"ba_gua"`

	GanZhi struct {
		Alignment float64 `json:"alignment"` // 天干地支对齐度
		Cycle     float64 `json:"cycle"`     // 周期完整度
		Strength  float64 `json:"strength"`  // 作用强度
	} `json:"gan_zhi"`

	Integration float64 `json:"integration"` // 整体集成度
	Coherence   float64 `json:"coherence"`   // 整体相干性
	Emergence   float64 `json:"emergence"`   // 涌现程度
}

// MetricsData 指标数据
type MetricsData struct {
	// 基础信息
	ID        string    `json:"id"`        // 指标ID
	Timestamp time.Time `json:"timestamp"` // 时间戳

	// 系统指标
	System struct {
		Energy    float64               `json:"energy"`    // 系统能量
		Field     core.FieldState       `json:"field"`     // 场状态
		Quantum   core.QuantumState     `json:"quantum"`   // 量子状态
		Emergence core.EmergentProperty `json:"emergence"` // 涌现属性
	} `json:"system"`

	// 模型指标
	Model ModelMetrics `json:"model"` // 模型指标

	// 扩展指标
	Custom map[string]interface{} `json:"custom"` // 自定义指标
}

// MetricsFilter 指标过滤器
type MetricsFilter struct {
	// 时间范围
	StartTime time.Time // 开始时间
	EndTime   time.Time // 结束时间

	// 指标选择
	Types   []string // 指标类型
	Sources []string // 指标来源
	Labels  []string // 指标标签

	// 值过滤
	MinValue float64 // 最小值
	MaxValue float64 // 最大值

	// 聚合选项
	GroupBy []string // 分组字段
	SortBy  string   // 排序字段
	Limit   int      // 限制数量

	// 自定义过滤条件
	Conditions map[string]interface{}
}

// EnergyMetrics 能量指标
type EnergyMetrics struct {
	Total    float64 `json:"total"`    // 总能量
	Average  float64 `json:"average"`  // 平均能量
	Variance float64 `json:"variance"` // 能量方差
}

// StateMetrics 状态指标
type StateMetrics struct {
	Transitions int     `json:"transitions"` // 状态转换次数
	Stability   float64 `json:"stability"`   // 稳定性
	Phase       Phase   `json:"phase"`       // 当前相位
}

// PerformanceMetrics 性能指标
type PerformanceMetrics struct {
	ThroughPut float64 `json:"throughput"` // 吞吐量
	Latency    float64 `json:"latency"`    // 延迟
	ErrorRate  float64 `json:"error_rate"` // 错误率
}

// ----------------------------------------------------
// GetTotalEnergy 获取总能量
func (m *ModelMetrics) GetTotalEnergy() float64 {
	return m.Energy.Total
}

// GetPreviousEnergy 获取前一个能量值
func (m *ModelMetrics) GetPreviousEnergy() float64 {
	if m.Energy.Average == 0 {
		return m.Energy.Total
	}
	return m.Energy.Total - m.Energy.Variance
}
