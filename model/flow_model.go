// model/flow_model.go

package model

import (
    "context"
    "sync"
    "time"

    "github.com/Corphon/daoflow/core"
)

// 模型类型
type ModelType uint8

const (
    ModelYinYang ModelType = iota  // 阴阳模型
    ModelWuXing                    // 五行模型
    ModelBaGua                     // 八卦模型
    ModelTianGan                   // 天干模型
    ModelDiZhi                     // 地支模型
)

// FlowModel 流模型接口
type FlowModel interface {
    // 基础流接口
    core.FlowSource
    
    // 获取模型类型
    GetModelType() ModelType
    
    // 相互作用
    Interact(other FlowModel) error
    
    // 能量转换
    Transform(pattern TransformPattern) error
    
    // 获取状态
    GetState() ModelState
}

// ModelState 模型状态
type ModelState struct {
    Energy      float64            // 能量值
    Phase       core.Phase         // 相位
    Nature      Nature             // 阴阳属性
    Properties  map[string]float64 // 属性值映射
}

// TransformPattern 转换模式
type TransformPattern struct {
    SourceType      ModelType  // 源模型类型
    TargetType      ModelType  // 目标模型类型
    TransformRatio  float64    // 转换比例
    EnergyVector    Vector3D   // 能量向量
}

// BaseFlowModel 基础流模型实现
type BaseFlowModel struct {
    mu          sync.RWMutex
    modelType   ModelType
    flow        *core.BaseFlow
    physics     *core.FlowPhysics
    energy      *core.EnergySystem
    field       *core.Field
    
    // 模型状态
    state       ModelState
    
    // 相互作用记录
    interactions map[string]InteractionRecord
    
    // 观察者
    observers   []ModelObserver
    done        chan struct{}
}

// InteractionRecord 相互作用记录
type InteractionRecord struct {
    Timestamp time.Time
    Target    ModelType
    Effect    float64
    Duration  time.Duration
}

// NewBaseFlowModel 创建基础流模型
func NewBaseFlowModel(modelType ModelType, capacity float64) *BaseFlowModel {
    return &BaseFlowModel{
        modelType: modelType,
        flow:      core.NewBaseFlow(&core.FlowConfig{
            MinEnergy:    0,
            MaxEnergy:    capacity,
            FlowInterval: time.Second,
        }),
        physics:  core.NewFlowPhysics(),
        energy:   core.NewEnergySystem(capacity),
        field:    core.NewField(core.ScalarField, 3),
        
        state: ModelState{
            Energy:     capacity * 0.5, // 初始能量50%
            Phase:      core.PhaseWuJi,
            Nature:     NatureBalance,
            Properties: make(map[string]float64),
        },
        
        interactions: make(map[string]InteractionRecord),
        done:        make(chan struct{}),
    }
}

// Interact 实现相互作用
func (bm *BaseFlowModel) Interact(other FlowModel) error {
    bm.mu.Lock()
    defer bm.mu.Unlock()

    // 计算物理相互作用
    pos := core.Vector3D{X: 0, Y: 0, Z: 0} // 可以基于模型位置计算
    fieldStrength := bm.field.CalculateFieldStrength(pos)
    
    // 能量交换
    energyTransfer := bm.calculateEnergyTransfer(other, fieldStrength)
    
    // 更新能量系统
    if err := bm.energy.Convert(
        core.FieldEnergy,
        core.KineticEnergy,
        energyTransfer,
    ); err != nil {
        return err
    }
    
    // 记录相互作用
    bm.interactions[other.GetState().Phase.String()] = InteractionRecord{
        Timestamp: time.Now(),
        Target:    other.GetModelType(),
        Effect:    energyTransfer,
        Duration:  time.Second,
    }
    
    return nil
}

// Transform 实现状态转换
func (bm *BaseFlowModel) Transform(pattern TransformPattern) error {
    bm.mu.Lock()
    defer bm.mu.Unlock()

    // 验证转换模式
    if pattern.SourceType != bm.modelType {
        return ErrInvalidTransform
    }

    // 计算能量变化
    energyChange := bm.calculateTransformEnergy(pattern)
    
    // 应用物理变换
    bm.physics.ApplyYinYangTransformation(pattern.TransformRatio)
    
    // 更新能量状态
    newEnergy := bm.state.Energy + energyChange
    if newEnergy < 0 || newEnergy > bm.energy.GetCapacity() {
        return ErrEnergyOutOfRange
    }
    
    bm.state.Energy = newEnergy
    
    // 更新相位和属性
    bm.updateStateProperties(pattern)
    
    return nil
}

// calculateEnergyTransfer 计算能量传递
func (bm *BaseFlowModel) calculateEnergyTransfer(
    other FlowModel,
    fieldStrength float64,
) float64 {
    // 基于场强和能量差计算传递量
    energyDiff := bm.state.Energy - other.GetState().Energy
    transferRatio := fieldStrength * 0.1 // 10%的场强作为传递系数
    
    return energyDiff * transferRatio
}

// updateStateProperties 更新状态属性
func (bm *BaseFlowModel) updateStateProperties(pattern TransformPattern) {
    // 更新相关属性
    bm.state.Properties["transformRatio"] = pattern.TransformRatio
    bm.state.Properties["fieldStrength"] = bm.field.CalculateFieldStrength(
        core.Vector3D{
            X: pattern.EnergyVector.X,
            Y: pattern.EnergyVector.Y,
            Z: pattern.EnergyVector.Z,
        },
    )
}

// GetModelType 获取模型类型
func (bm *BaseFlowModel) GetModelType() ModelType {
    return bm.modelType
}

// GetState 获取模型状态
func (bm *BaseFlowModel) GetState() ModelState {
    bm.mu.RLock()
    defer bm.mu.RUnlock()
    return bm.state
}
