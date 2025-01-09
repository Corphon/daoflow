// model/interfaces.go

package model

import (
	"time"
)

// ModelType 模型类型
type ModelType uint8

const (
	ModelTypeNone  ModelType = iota
	ModelYinYang             // 阴阳模型
	ModelWuXing              // 五行模型
	ModelBaGua               // 八卦模型
	ModelGanZhi              // 干支模型
	ModelIntegrate           // 集成模型
	ModelTypeMax
)

// Phase 相位
type Phase uint8

const (
	PhaseNone  Phase = iota
	PhaseYin         // 阴相
	PhaseYang        // 阳相
	PhaseWood        // 木相
	PhaseFire        // 火相
	PhaseEarth       // 土相
	PhaseMetal       // 金相
	PhaseWater       // 水相
	PhaseMax
)

// Nature 属性
type Nature uint8

const (
	NatureNeutral Nature = iota
	NatureYin            // 阴性
	NatureYang           // 阳性
	NatureMax
)

// TransformPattern 转换模式
type TransformPattern uint8

const (
	PatternNone    TransformPattern = iota
	PatternNormal                   // 常规转换
	PatternForward                  // 顺序转换
	PatternReverse                  // 逆序转换
	PatternBalance                  // 平衡转换
	PatternMutate                   // 变异转换
	PatternMax
)

// FlowModel 流模型接口
type FlowModel interface {
	// 基础操作
	Start() error
	Stop() error
	Transform(pattern TransformPattern) error
	GetState() ModelState
	Close() error
}

// SystemState 系统状态
type SystemState struct {
	// 基础属性
	Energy    float64   // 系统能量
	Entropy   float64   // 系统熵
	Harmony   float64   // 和谐度
	Balance   float64   // 平衡度
	Phase     Phase     // 系统相位
	Timestamp time.Time // 时间戳

	// 子系统能量
	YinYang      float64 // 阴阳能量
	WuXingEnergy float64 // 五行能量
	BaGuaEnergy  float64 // 八卦能量
	GanZhiEnergy float64 // 干支能量

	// 系统详情
	System struct {
		Energy       float64 // 总能量
		Entropy      float64 // 系统熵
		WuXingEnergy float64 // 五行能量
		BaGuaEnergy  float64 // 八卦能量
		GanZhiEnergy float64 // 干支能量
	}

	// 扩展属性
	Properties map[string]interface{}
}

// ModelState 模型状态
type ModelState struct {
	Type       ModelType              // 模型类型
	Energy     float64                // 能量值
	Phase      Phase                  // 相位
	Nature     Nature                 // 属性
	Properties map[string]interface{} // 扩展属性
	UpdateTime time.Time              // 更新时间
}

// Vector3D 三维向量
type Vector3D struct {
	X float64
	Y float64
	Z float64
}
