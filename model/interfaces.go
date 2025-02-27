// model/interfaces.go

package model

import (
	"math"
	"time"

	"github.com/Corphon/daoflow/core"
)

// Model 模型接口
type Model interface {
	// 基础操作
	Start() error
	Stop() error
	Reset() error
	Close() error

	// 状态管理
	GetState() ModelState
	GetSystemState() SystemState
	Transform(pattern TransformPattern) error

	// 核心操作
	GetCoreState() CoreState
	UpdateCoreState(state CoreState) error
	ValidateCoreState() error

	// 能量操作
	SetEnergy(energy float64) error
	AdjustEnergy(delta float64) error
}

// CoreState core层状态
type CoreState struct {
	QuantumState  *core.QuantumState // 量子态
	FieldState    *core.Field        // 场态
	EnergyState   *core.EnergySystem // 能量态
	InteractState *core.Interaction  // 交互态
	Phase         float64
	Properties    map[string]float64
}

// CoreComponent core层组件
type CoreComponent interface {
	Initialize() error
	Update() error
	Reset() error
	Close() error
}

// ModelType 模型类型
type ModelType uint8

const (
	ModelTypeNone  ModelType = iota
	ModelYinYang             // 阴阳模型
	ModelWuXing              // 五行模型
	ModelBaGua               // 八卦模型
	ModelGanZhi              // 干支模型
	ModelIntegrate           // 集成模型
	ModelTypeAlert           // 告警模型
	ModelTypeMax
)

// Phase 相位
type Phase uint8

const (
	PhaseNone    Phase = iota
	PhaseYin           // 阴相
	PhaseYang          // 阳相
	PhaseYinYang       // 阴阳平衡相位
	PhaseWood          // 木相
	PhaseFire          // 火相
	PhaseEarth         // 土相
	PhaseMetal         // 金相
	PhaseWater         // 水相
	PhaseMax
)

// ProcessPhase 流程相位
type ProcessPhase uint8

const (
	ProcessPhaseNone      ProcessPhase = iota
	ProcessPhaseInitial                // 初始相位
	ProcessPhaseTransform              // 转换相位
	ProcessPhaseStable                 // 稳定相位
	ProcessPhaseComplete               // 完成相位
	ProcessPhaseMax
)

// Phase 相位常量添加
const (
	PhaseTransform Phase = iota + PhaseMax // 转换相位
	Phase_Stable                           // 稳定相位
	Phase_Unstable                         // 不稳定相位
	PhaseNeutral                           // 中性相位
)

// Nature 属性
type Nature uint8

const (
	NatureNone     Nature = iota
	NatureNeutral         // 中性
	NatureStable          // 稳定
	NatureUnstable        // 不稳定
	NatureYin             // 阴性
	NatureYang            // 阳性
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
	GetSystemState() SystemState
	GetType() ModelType
	Close() error

	// Core层操作
	GetCoreState() CoreState
	UpdateCoreState(state CoreState) error
	ValidateCoreState() error

	// 能量操作
	SetEnergy(energy float64) error
	AdjustEnergy(delta float64) error
}

// FlowPattern 流模式
type FlowPattern struct {
	ID         string                 // 模式ID
	Type       string                 // 模式类型
	Flow       FlowModel              // 相关流模型
	Strength   float64                // 模式强度
	Duration   time.Duration          // 持续时间
	State      ModelState             // 关联状态
	Metrics    PatternMetrics         // 模式指标
	Properties map[string]interface{} // 模式属性
	Created    time.Time              // 创建时间
}

// Anomaly 异常
type Anomaly struct {
	ID        string                 // 异常ID
	Type      string                 // 异常类型
	Level     string                 // 严重级别
	Message   string                 // 异常描述
	Source    string                 // 异常来源
	Time      time.Time              // 发生时间
	Data      map[string]interface{} // 异常数据
	Subtype   string                 // 异常子类型
	Severity  float64                // 严重程度
	Value     float64                // 当前值
	Expected  float64                // 期望值
	Threshold float64                // 阈值

}

// SystemState 系统状态
type SystemState struct {
	// 基础属性
	Entropy float64 // 系统熵
	Harmony float64 // 和谐度
	Balance float64 // 平衡度

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

	Phase      Phase                  `json:"phase"`      // 系统相位
	Energy     float64                `json:"energy"`     // 系统能量
	Stability  float64                `json:"stability"`  // 系统稳定性
	Properties map[string]interface{} `json:"properties"` // 系统属性
}

// ModelState 模型状态
type ModelState struct {
	Type       ModelType              // 模型类型
	Energy     float64                // 能量值
	Phase      Phase                  // 相位
	Nature     Nature                 // 属性
	Health     float64                // 健康度
	Properties map[string]interface{} // 扩展属性
	UpdateTime time.Time              // 更新时间

	// 阴阳相关
	YinEnergy  float64 // 阴能量
	YangEnergy float64 // 阳能量
	Harmony    float64 // 和谐度
	Balance    float64 // 平衡度
}

// Vector3D 三维向量
type Vector3D struct {
	X float64
	Y float64
	Z float64
}

// StateValidator 状态验证器
type StateValidator interface {
	ValidateState(state ModelState) error
	ValidateCoreState(state CoreState) error
	ValidateTransition(from, to ModelState) error
}

// CoreComponentManager Core组件管理器
type CoreComponentManager interface {
	// 组件管理
	RegisterComponent(name string, component CoreComponent) error
	GetComponent(name string) (CoreComponent, error)
	RemoveComponent(name string) error

	// 状态同步
	SyncComponents() error
	ValidateComponents() error

	// 生命周期
	Initialize() error
	Close() error
}

// Core相关常量
const (
	// 量子态阈值
	MinQuantumCoherence = 0.1
	MaxQuantumCoherence = 1.0

	// 场态阈值
	MinFieldStrength = 0.0
	MaxFieldStrength = 100.0

	// 能量阈值
	MinSystemEnergy = 0.0
	MaxSystemEnergy = 1000.0

	// 交互阈值
	MinInteraction = 0.0
	MaxInteraction = 1.0
)

// StateTransition 状态转换记录
type StateTransition struct {
	ID        string                 // 转换标识
	From      SystemState            // 起始状态
	To        SystemState            // 目标状态
	Type      TransformPattern       // 转换类型
	Timestamp time.Time              // 转换时间
	Duration  time.Duration          // 转换耗时
	Success   bool                   // 是否成功
	Error     error                  // 错误信息
	Metadata  map[string]interface{} // 额外元数据
}

// StateTransitionResult 状态转换结果
type StateTransitionResult struct {
	Success     bool                   // 是否成功
	Error       error                  // 错误信息
	Timestamp   time.Time              // 完成时间
	Duration    time.Duration          // 转换耗时
	ResultState SystemState            // 结果状态
	Details     map[string]interface{} // 详细信息
}

// Field 场接口
type Field interface {
	GetState() (*core.FieldState, error)
}

// FieldState 场状态
type FieldState struct {
	Energy     float64            // 场能量
	Elements   []*WuXingElement   // 元素列表
	Properties map[string]float64 // 属性集
	Timestamp  time.Time          // 时间戳
	Quantum    *core.QuantumState // 量子态
}

// Element 元素结构体定义
type Element struct {
	Type       string             // 元素类型
	Energy     float64            // 能量值
	Properties map[string]float64 // 属性集
}

// 添加相关的访问方法-------------------------------------
func (e *Element) GetType() string {
	return e.Type
}

func (e *Element) GetEnergy() float64 {
	return e.Energy
}

func (e *Element) GetProperties() map[string]float64 {
	return e.Properties
}

func (e *Element) SetEnergy(energy float64) {
	e.Energy = energy
}

// ToFloat64 将Phase转换为float64
func (p Phase) ToFloat64() float64 {
	return float64(p)
}

// FromFloat64 从float64创建Phase
func FromFloat64(f float64) Phase {
	// 标准化到[0, PhaseMax)范围
	normalized := math.Mod(f, float64(PhaseMax))
	if normalized < 0 {
		normalized += float64(PhaseMax)
	}
	return Phase(normalized)
}

// GetElements 获取元素列表
func (fs *FieldState) GetElements() []*WuXingElement {
	if fs == nil {
		return nil
	}
	return fs.Elements
}

// GetEnergyLevel 获取能量等级
func (fs *FieldState) GetEnergyLevel() float64 {
	if fs == nil {
		return 0
	}
	return fs.Energy
}

// GetEnergyFlow 获取能量流动
func (fs *FieldState) GetEnergyFlow() float64 {
	if fs == nil {
		return 0
	}
	if value, exists := fs.Properties["energy_flow"]; exists {
		return value
	}
	return 0
}

// GetEnergyDistribution 获取能量分布
func (fs *FieldState) GetEnergyDistribution() map[core.Point]float64 {
	distribution := make(map[core.Point]float64)

	// 将元素能量映射到空间点上
	for i, elem := range fs.Elements {
		if energy := elem.GetEnergy(); energy > 0 {
			// 基于元素序号创建二维网格坐标
			x := i % 10 // 10x10网格
			y := i / 10
			point := core.Point{X: x, Y: y}
			distribution[point] = energy
		}
	}

	return distribution
}

// GetQuantumState 获取量子态
func (fs *FieldState) GetQuantumState() *core.QuantumState {
	if fs == nil {
		return nil
	}
	return fs.Quantum
}

// HasElement 判断是否存在指定类型的元素
func (fs *FieldState) HasElement(elementType string) bool {
	if fs == nil || len(fs.Elements) == 0 {
		return false
	}
	for _, elem := range fs.Elements {
		if elem.GetType() == elementType {
			return true
		}
	}
	return false
}

// HasEnergyLevel 判断能量是否达到指定水平
func (fs *FieldState) HasEnergyLevel(level float64) bool {
	if fs == nil {
		return false
	}
	return fs.Energy >= level
}

// HasQuantumState 判断是否存在特定量子态
func (fs *FieldState) HasQuantumState(stateValue float64) bool {
	if fs == nil || fs.Quantum == nil {
		return false
	}
	// 检查量子态概率
	prob := fs.Quantum.GetProbability()
	return math.Abs(prob-stateValue) < 0.1 // 允许0.1的误差范围
}

// GetFieldStrength 获取场强度
func (fs *FieldState) GetFieldStrength() float64 {
	// 检查nil
	if fs == nil {
		return 0
	}

	// 从Properties获取场强度
	if strength, exists := fs.Properties["field_strength"]; exists {
		return strength
	}

	// 从Energy估算场强度
	return fs.Energy * 0.8 // 默认场强度为能量的80%
}
