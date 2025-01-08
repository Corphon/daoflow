// model/flow_model.go

package model

import (
    "context"
    "sync"
    "time"

    "github.com/Corphon/daoflow/core"
)

// FlowModel 流模型接口
type FlowModel interface {
    // 基础流控制
    Initialize() error
    Start() error
    Stop() error
    
    // 状态管理
    GetModelType() ModelType
    GetState() ModelState
    
    // 模型交互
    Interact(other FlowModel) error
    Transform(pattern TransformPattern) error
    
    // 能量管理
    AdjustEnergy(delta float64) error
}

// BaseFlowModel 基础流模型实现
type BaseFlowModel struct {
    mu          sync.RWMutex
    modelType   ModelType
    capacity    float64
    
    // 状态管理
    state       ModelState
    
    // 子系统组件
    coreFlow    *core.Flow        // 核心流体系统
    corePhysics *core.FlowPhysics // 物理特性系统
    
    // 状态追踪
    interactions map[string]InteractionRecord
    
    // 观察者
    observers   []ModelObserver
    done        chan struct{}
}

// NewBaseFlowModel 创建基础流模型
func NewBaseFlowModel(modelType ModelType, capacity float64) *BaseFlowModel {
    return &BaseFlowModel{
        modelType: modelType,
        capacity:  capacity,
        state: ModelState{
            Energy:     0,
            Phase:      PhaseWuJi,
            Nature:     NatureBalance,
            Properties: make(map[string]float64),
        },
        coreFlow:    core.NewFlow(&core.FlowConfig{
            MinEnergy:    0,
            MaxEnergy:    capacity,
            FlowInterval: time.Second,
        }),
        corePhysics: core.NewFlowPhysics(),
        interactions: make(map[string]InteractionRecord),
        observers:   make([]ModelObserver, 0),
        done:        make(chan struct{}),
    }
}

// Initialize 初始化模型
func (bm *BaseFlowModel) Initialize() error {
    bm.mu.Lock()
    defer bm.mu.Unlock()
    
    // 初始化核心流体系统
    if err := bm.coreFlow.Initialize(); err != nil {
        return err
    }
    
    // 更新初始状态
    bm.updateStateFromCore()
    return nil
}

// Start 启动模型
func (bm *BaseFlowModel) Start() error {
    bm.mu.Lock()
    defer bm.mu.Unlock()
    
    // 启动核心流体系统
    return bm.coreFlow.Start()
}

// Stop 停止模型
func (bm *BaseFlowModel) Stop() error {
    bm.mu.Lock()
    defer bm.mu.Unlock()
    
    close(bm.done)
    return bm.coreFlow.Stop()
}

// GetModelType 获取模型类型
func (bm *BaseFlowModel) GetModelType() ModelType {
    return bm.modelType
}

// GetState 获取当前状态
func (bm *BaseFlowModel) GetState() ModelState {
    bm.mu.RLock()
    defer bm.mu.RUnlock()
    
    bm.updateStateFromCore()
    return bm.state
}

// Interact 模型交互
func (bm *BaseFlowModel) Interact(other FlowModel) error {
    bm.mu.Lock()
    defer bm.mu.Unlock()
    
    // 计算交互效果
    effect := bm.calculateInteractionEffect(other)
    
    // 记录交互
    bm.recordInteraction(other.GetModelType().String(), effect)
    
    // 应用交互效果
    return bm.AdjustEnergy(effect)
}

// Transform 状态转换
func (bm *BaseFlowModel) Transform(pattern TransformPattern) error {
    bm.mu.Lock()
    defer bm.mu.Unlock()
    
    if pattern.SourceType != bm.modelType {
        return ErrInvalidTransform
    }
    
    // 计算能量变化
    energyChange := bm.calculateTransformEnergy(pattern)
    
    // 应用物理变换
    if err := bm.corePhysics.ApplyYinYangTransformation(pattern.TransformRatio); err != nil {
        return err
    }
    
    // 更新能量
    if err := bm.AdjustEnergy(energyChange); err != nil {
        return err
    }
    
    // 更新状态
    bm.updateStateFromCore()
    return nil
}

// AdjustEnergy 调整能量
func (bm *BaseFlowModel) AdjustEnergy(delta float64) error {
    newEnergy := bm.state.Energy + delta
    if newEnergy < 0 || newEnergy > bm.capacity {
        return ErrEnergyOutOfRange
    }
    
    // 更新核心流体系统的能量
    if err := bm.coreFlow.SetEnergy(newEnergy); err != nil {
        return err
    }
    
    // 更新状态
    bm.state.Energy = newEnergy
    bm.notifyObservers()
    return nil
}

// 内部辅助方法

// updateStateFromCore 从核心系统更新状态
func (bm *BaseFlowModel) updateStateFromCore() {
    coreState := bm.coreFlow.GetState()
    physicsState := bm.corePhysics.GetState()
    
    bm.state.Energy = coreState.Energy
    bm.state.Properties["density"] = physicsState.Density
    bm.state.Properties["temperature"] = physicsState.Temperature
    bm.state.Properties["pressure"] = physicsState.Pressure
    bm.state.Properties["entropy"] = physicsState.Entropy
}

// calculateInteractionEffect 计算交互效果
func (bm *BaseFlowModel) calculateInteractionEffect(other FlowModel) float64 {
    otherState := other.GetState()
    energyDiff := bm.state.Energy - otherState.Energy
    return energyDiff * 0.1 // 简单的能量平衡效应
}

// calculateTransformEnergy 计算转换能量
func (bm *BaseFlowModel) calculateTransformEnergy(pattern TransformPattern) float64 {
    return bm.state.Energy * pattern.TransformRatio
}

// recordInteraction 记录交互
func (bm *BaseFlowModel) recordInteraction(targetType string, effect float64) {
    bm.interactions[targetType] = InteractionRecord{
        Timestamp: time.Now(),
        Effect:    effect,
        Duration:  time.Second,
    }
}

// notifyObservers 通知观察者
func (bm *BaseFlowModel) notifyObservers() {
    for _, observer := range bm.observers {
        observer.OnStateChange(bm.state)
    }
}

// AddObserver 添加观察者
func (bm *BaseFlowModel) AddObserver(observer ModelObserver) {
    bm.mu.Lock()
    defer bm.mu.Unlock()
    bm.observers = append(bm.observers, observer)
}

// RemoveObserver 移除观察者
func (bm *BaseFlowModel) RemoveObserver(observer ModelObserver) {
    bm.mu.Lock()
    defer bm.mu.Unlock()
    for i, obs := range bm.observers {
        if obs == observer {
            bm.observers = append(bm.observers[:i], bm.observers[i+1:]...)
            break
        }
    }
}
