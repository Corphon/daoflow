// model/types.go

package model

import (
    "sync"
    "time"
)

// ModelType 模型类型
type ModelType uint8

const (
    ModelYinYang ModelType = iota  // 阴阳模型
    ModelWuXing                    // 五行模型
    ModelBaGua                     // 八卦模型
    ModelGanZhi                    // 干支模型
    ModelIntegrate                 // 集成模型
)

// Nature 阴阳性质
type Nature uint8

const (
    NatureYin     Nature = iota  // 阴性
    NatureYang                   // 阳性
    NatureBalance               // 平衡
)

// Phase 相位状态
type Phase uint8

const (
    PhaseWuJi    Phase = iota  // 无极
    PhaseTaiJi                 // 太极
    PhaseYinYang               // 两仪
    PhaseWuXing                // 五行
    PhaseBaGua                 // 八卦
)

// Vector3D 三维向量
type Vector3D struct {
    X, Y, Z float64
}

// ModelState 模型状态
type ModelState struct {
    Energy     float64             // 能量值
    Phase      Phase               // 相位
    Nature     Nature              // 阴阳属性
    Properties map[string]float64  // 属性映射
}

// StateTransition 状态转换
type StateTransition struct {
    FromState   string                // 源状态
    ToState     string                // 目标状态
    Timestamp   time.Time             // 转换时间
    Properties  map[string]float64    // 转换属性
    Reason      string                // 转换原因
}

// TransformPattern 转换模式
type TransformPattern struct {
    SourceType     ModelType  // 源模型类型
    TargetType     ModelType  // 目标模型类型
    TransformRatio float64    // 转换比例
    EnergyVector   Vector3D   // 能量向量
}

// InteractionRecord 交互记录
type InteractionRecord struct {
    Timestamp  time.Time  // 交互时间
    TargetType ModelType  // 目标类型
    Effect     float64    // 效果值
    Duration   time.Duration  // 持续时间
}

// ModelObserver 模型观察者接口
type ModelObserver interface {
    OnStateChange(state ModelState)
    OnTransform(transition StateTransition)
}

// SafeCounter 线程安全计数器
type SafeCounter struct {
    mu    sync.RWMutex
    value int64
}

func (sc *SafeCounter) Increment() int64 {
    sc.mu.Lock()
    defer sc.mu.Unlock()
    sc.value++
    return sc.value
}

func (sc *SafeCounter) GetValue() int64 {
    sc.mu.RLock()
    defer sc.mu.RUnlock()
    return sc.value
}
