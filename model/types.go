// model/types.go

package model

import (
    "time"
)

// ModelType 模型类型
type ModelType uint8

const (
    ModelNone ModelType = iota
    ModelYinYang        // 阴阳模型
    ModelWuXing        // 五行模型
    ModelBaGua         // 八卦模型
    ModelGanZhi        // 天干地支模型
    ModelIntegrate     // 集成模型
)

// Phase 相位类型
type Phase uint8

const (
    PhaseNone Phase = iota
    PhaseYinYang     // 阴阳相位
    PhaseWuXing      // 五行相位
    PhaseBaGua       // 八卦相位
    PhaseGanZhi      // 天干地支相位
)

// Nature 阴阳性质
type Nature uint8

const (
    NatureNone Nature = iota
    NatureYin         // 阴性
    NatureYang        // 阳性
)

// WuXingPhase 五行相位
type WuXingPhase uint8

const (
    WuXingNone WuXingPhase = iota
    Metal                   // 金
    Wood                    // 木
    Water                   // 水
    Fire                    // 火
    Earth                   // 土
)

// TransformPattern 转换模式
type TransformPattern uint8

const (
    PatternNone TransformPattern = iota
    PatternNormal               // 常规转换
    PatternForward              // 顺序转换
    PatternReverse              // 逆序转换
    PatternBalance              // 平衡转换
    PatternMutate              // 变异转换
)

// ModelState 模型状态
type ModelState struct {
    Type       ModelType             // 模型类型
    Phase      Phase                 // 当前相位
    Energy     float64               // 当前能量
    Nature     Nature                // 当前性质
    Properties map[string]interface{} // 属性集合
    UpdateTime time.Time             // 更新时间
}

// FlowModel 流模型接口
type FlowModel interface {
    // 基本信息
    GetModelType() ModelType
    GetState() ModelState
    
    // 状态控制
    Start() error
    Stop() error
    Reset() error
    
    // 能量控制
    GetEnergy() float64
    SetEnergy(energy float64) error
    AdjustEnergy(delta float64) error
    
    // 转换控制
    Transform(pattern TransformPattern) error
    
    // 相位控制
    GetPhase() Phase
    SetPhase(phase Phase) error
}

// StateTransition 状态转换
type StateTransition struct {
    From      ModelState
    To        ModelState
    Pattern   TransformPattern
    Timestamp time.Time
    Changes   map[string]interface{}
}

// Vector3D 三维向量
type Vector3D struct {
    X float64
    Y float64
    Z float64
}

// Point2D 二维点
type Point2D struct {
    X float64
    Y float64
}

// 五行关系常量
const (
    RelationNone = iota
    RelationGenerate  // 相生
    RelationRestrain  // 相克
    RelationWeaken    // 相泄
    RelationControl   // 相制
)

// 系统常量
const (
    DefaultCapacity    = 1000.0  // 默认容量
    DefaultEnergy     = 100.0   // 默认能量
    MinEnergy        = 0.0     // 最小能量
    MaxEnergy        = 10000.0 // 最大能量
    EnergyThreshold  = 0.1     // 能量阈值
    PhaseThreshold   = 0.2     // 相位阈值
    BalanceThreshold = 0.3     // 平衡阈值
)

// String 方法实现

func (mt ModelType) String() string {
    switch mt {
    case ModelYinYang:
        return "YinYang"
    case ModelWuXing:
        return "WuXing"
    case ModelBaGua:
        return "BaGua"
    case ModelGanZhi:
        return "GanZhi"
    case ModelIntegrate:
        return "Integrate"
    default:
        return "Unknown"
    }
}

func (p Phase) String() string {
    switch p {
    case PhaseYinYang:
        return "YinYang"
    case PhaseWuXing:
        return "WuXing"
    case PhaseBaGua:
        return "BaGua"
    case PhaseGanZhi:
        return "GanZhi"
    default:
        return "Unknown"
    }
}

func (n Nature) String() string {
    switch n {
    case NatureYin:
        return "Yin"
    case NatureYang:
        return "Yang"
    default:
        return "Unknown"
    }
}

func (w WuXingPhase) String() string {
    switch w {
    case Metal:
        return "Metal"
    case Wood:
        return "Wood"
    case Water:
        return "Water"
    case Fire:
        return "Fire"
    case Earth:
        return "Earth"
    default:
        return "Unknown"
    }
}

// ValidateEnergy 验证能量值
func ValidateEnergy(energy float64) bool {
    return energy >= MinEnergy && energy <= MaxEnergy
}

// ValidatePhase 验证相位
func ValidatePhase(phase Phase) bool {
    return phase >= PhaseYinYang && phase <= PhaseGanZhi
}

// ValidateNature 验证阴阳性质
func ValidateNature(nature Nature) bool {
    return nature == NatureYin || nature == NatureYang
}
