// core/types.go

package core

import (
	"sync"
	"time"
)

// FieldState 场状态
type FieldState struct {
	mu        sync.RWMutex // 添加读写锁
	Strength  [][]float64  // 场强度分布
	Potential [][]float64  // 势能分布
	Gradient  [][]Vector3D // 梯度分布
	Phase     float64      // 场相位
	Energy    float64      // 场能量
	Frequency float64      // 场频率
	Amplitude float64      // 场振幅
	Timestamp time.Time    // 状态时间戳
	Flow      float64      // 能量流
	Dimension int          // 维度
}

// FieldParams 场参数
type FieldParams struct {
	Type      FieldType              // 场类型
	Dimension int                    // 场维度
	GridSize  int                    // 网格大小
	Boundary  []float64              // 边界条件
	Config    map[string]interface{} // 配置参数
}

// EmergentPattern 涌现模式
type EmergentPattern struct {
	ID         string // 模式ID
	Type       string // 模式类型
	Components []struct {
		Type  string
		Value float64
	}
	Properties map[string]float64 // 模式属性
	Strength   float64            // 模式强度
	Energy     float64            // 模式能量
	Created    time.Time          // 创建时间
}

// PotentialPattern 潜在模式
type PotentialPattern struct {
	Type        string        // 模式类型
	Probability float64       // 出现概率
	Energy      float64       // 所需能量
	Components  []string      // 组件要求
	TimeFrame   time.Duration // 预计时间框架
}

// PotentialEmergence 潜在涌现
type PotentialEmergence struct {
	ID          string             // 涌现ID
	Type        string             // 涌现类型
	Probability float64            // 出现概率
	Patterns    []PotentialPattern // 相关模式
	Energy      float64            // 预计能量需求
	Properties  map[string]float64 // 预测属性
	TimeWindow  time.Duration      // 时间窗口
	Conditions  map[string]float64 // 触发条件
}

// EmergentProperty 涌现属性
type EmergentProperty struct {
	ID         string             // 属性ID
	Type       string             // 属性类型
	Pattern    *EmergentPattern   // 关联模式
	Properties map[string]float64 // 属性值
	Energy     float64            // 能量值
	Created    time.Time          // 创建时间
}

// Point 二维点结构
type Point struct {
	X int
	Y int
}

// -----------------------------------------------
// GetFrequency 获取场频率
func (fs *FieldState) GetFrequency() float64 {
	return fs.Frequency
}

// GetAmplitude 获取场振幅
func (fs *FieldState) GetAmplitude() float64 {
	return fs.Amplitude
}
