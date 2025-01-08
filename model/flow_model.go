// model/flow_model.go

package model

import (
    "sync"
    "time"

    "github.com/Corphon/daoflow/core"
)

// BaseFlowModel 基础流模型
type BaseFlowModel struct {
    mu       sync.RWMutex
    modelType ModelType
    capacity  float64

    // 基础状态
    state    ModelState
    running  bool
    
    // 内部组件
    quantum  *core.QuantumState  // 量子状态
    field    *core.Field        // 统一场
    
    // 控制
    done     chan struct{}
}

// NewBaseFlowModel 创建基础流模型
func NewBaseFlowModel(modelType ModelType, capacity float64) *BaseFlowModel {
    if capacity <= 0 {
        capacity = DefaultCapacity
    }

    return &BaseFlowModel{
        modelType: modelType,
        capacity:  capacity,
        state: ModelState{
            Type:       modelType,
            Energy:     0,
            Properties: make(map[string]interface{}),
            UpdateTime: time.Now(),
        },
        quantum:  core.NewQuantumState(),
        field:    core.NewField(),
        done:     make(chan struct{}),
    }
}

// GetModelType 获取模型类型
func (bm *BaseFlowModel) GetModelType() ModelType {
    return bm.modelType
}

// GetState 获取当前状态
func (bm *BaseFlowModel) GetState() ModelState {
    bm.mu.RLock()
    defer bm.mu.RUnlock()
    return bm.state
}

// Start 启动模型
func (bm *BaseFlowModel) Start() error {
    bm.mu.Lock()
    defer bm.mu.Unlock()

    if bm.running {
        return NewModelError(ErrCodeOperation, "model already started", nil)
    }

    bm.running = true
    bm.done = make(chan struct{})

    // 初始化量子态
    bm.quantum.Initialize()
    
    // 初始化场
    bm.field.Initialize()

    return nil
}

// Stop 停止模型
func (bm *BaseFlowModel) Stop() error {
    bm.mu.Lock()
    defer bm.mu.Unlock()

    if !bm.running {
        return NewModelError(ErrCodeOperation, "model not running", nil)
    }

    bm.running = false
    close(bm.done)
    return nil
}

// Reset 重置模型
func (bm *BaseFlowModel) Reset() error {
    bm.mu.Lock()
    defer bm.mu.Unlock()

    // 停止运行
    if bm.running {
        if err := bm.Stop(); err != nil {
            return err
        }
    }

    // 重置状态
    bm.state = ModelState{
        Type:       bm.modelType,
        Energy:     0,
        Properties: make(map[string]interface{}),
        UpdateTime: time.Now(),
    }

    // 重置量子态
    bm.quantum.Reset()
    
    // 重置场
    bm.field.Reset()

    return nil
}

// GetEnergy 获取能量
func (bm *BaseFlowModel) GetEnergy() float64 {
    bm.mu.RLock()
    defer bm.mu.RUnlock()
    return bm.state.Energy
}

// SetEnergy 设置能量
func (bm *BaseFlowModel) SetEnergy(energy float64) error {
    if !ValidateEnergy(energy) {
        return NewModelError(ErrCodeOperation, "invalid energy value", nil)
    }

    bm.mu.Lock()
    defer bm.mu.Unlock()

    // 更新量子态
    probability := energy / bm.capacity
    bm.quantum.SetProbability(probability)
    
    // 更新场强度
    bm.field.SetStrength(probability)

    // 更新状态
    bm.state.Energy = energy
    bm.state.UpdateTime = time.Now()

    return nil
}

// AdjustEnergy 调整能量
func (bm *BaseFlowModel) AdjustEnergy(delta float64) error {
    bm.mu.Lock()
    defer bm.mu.Unlock()

    newEnergy := bm.state.Energy + delta
    if !ValidateEnergy(newEnergy) {
        return NewModelError(ErrCodeOperation, "energy adjustment out of range", nil)
    }

    // 更新量子态
    probability := newEnergy / bm.capacity
    bm.quantum.SetProbability(probability)
    
    // 更新场强度
    bm.field.SetStrength(probability)

    // 更新状态
    bm.state.Energy = newEnergy
    bm.state.UpdateTime = time.Now()

    return nil
}

// GetPhase 获取相位
func (bm *BaseFlowModel) GetPhase() Phase {
    bm.mu.RLock()
    defer bm.mu.RUnlock()
    return bm.state.Phase
}

// SetPhase 设置相位
func (bm *BaseFlowModel) SetPhase(phase Phase) error {
    if !ValidatePhase(phase) {
        return NewModelError(ErrCodeOperation, "invalid phase value", nil)
    }

    bm.mu.Lock()
    defer bm.mu.Unlock()

    // 更新量子态相位
    bm.quantum.SetPhase(float64(phase))
    
    // 更新场相位
    bm.field.SetPhase(float64(phase))

    // 更新状态
    bm.state.Phase = phase
    bm.state.UpdateTime = time.Now()

    return nil
}

// Transform 基础转换实现
func (bm *BaseFlowModel) Transform(pattern TransformPattern) error {
    bm.mu.Lock()
    defer bm.mu.Unlock()

    if !bm.running {
        return NewModelError(ErrCodeOperation, "model not running", nil)
    }

    // 量子态演化
    bm.quantum.Evolve(pattern.String())
    
    // 场演化
    bm.field.Evolve()

    // 更新状态
    bm.state.UpdateTime = time.Now()

    return nil
}

// updateState 更新状态
func (bm *BaseFlowModel) updateState(updates map[string]interface{}) {
    for k, v := range updates {
        bm.state.Properties[k] = v
    }
    bm.state.UpdateTime = time.Now()
}

// validateState 验证状态
func (bm *BaseFlowModel) validateState() error {
    if !ValidateEnergy(bm.state.Energy) {
        return NewModelError(ErrCodeState, "invalid energy state", nil)
    }

    if !ValidatePhase(bm.state.Phase) {
        return NewModelError(ErrCodeState, "invalid phase state", nil)
    }

    return nil
}
