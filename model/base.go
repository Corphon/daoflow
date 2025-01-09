// model/base.go

package model

import (
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
)

// BaseFlowModel 基础流模型
type BaseFlowModel struct {
	mu sync.RWMutex

	// 模型标识
	modelType ModelType
	capacity  float64

	// 基础状态
	state   ModelState
	running bool

	// 状态管理
	stateManager *StateManager

	// 内部组件 - 对外隐藏核心实现
	components struct {
		quantum *core.QuantumState
		field   *core.Field
		energy  *core.EnergySystem
	}

	// 控制
	done chan struct{}
}

// NewBaseFlowModel 创建基础流模型
func NewBaseFlowModel(modelType ModelType, capacity float64) *BaseFlowModel {
	base := &BaseFlowModel{
		modelType: modelType,
		capacity:  capacity,
		state: ModelState{
			Type:       modelType,
			Energy:     0,
			Properties: make(map[string]interface{}),
			UpdateTime: time.Now(),
		},
		done: make(chan struct{}),
	}

	// 初始化状态管理器
	base.stateManager = NewStateManager(modelType, capacity)

	// 初始化内部组件 - 通过专门的方法初始化量子态和场
	base.initializeComponents()

	return base
}

// initializeComponents 初始化组件
func (b *BaseFlowModel) initializeComponents() {
	b.components.quantum = core.NewQuantumState()
	b.components.field = core.NewField(core.ScalarField, 3)
	b.components.energy = core.NewEnergySystem(b.capacity)
}

// Start 启动模型
func (b *BaseFlowModel) Start() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.running {
		return NewModelError(ErrCodeOperation, "model already started", nil)
	}

	// 初始化内部状态
	if err := b.initializeState(); err != nil {
		return WrapError(err, ErrCodeOperation, "failed to initialize state")
	}

	b.running = true
	return nil
}

// Stop 停止模型
func (b *BaseFlowModel) Stop() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.running {
		return NewModelError(ErrCodeOperation, "model not running", nil)
	}

	// 保存最终状态
	if err := b.stateManager.UpdateState(); err != nil {
		return WrapError(err, ErrCodeOperation, "failed to save final state")
	}

	b.running = false
	return nil
}

// Transform 执行状态转换
func (b *BaseFlowModel) Transform(pattern TransformPattern) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.running {
		return NewModelError(ErrCodeOperation, "model not running", nil)
	}

	// 获取当前状态
	state := b.stateManager.GetModelState()

	// 执行转换
	if err := b.stateManager.transformer.ApplyTransform(
		pattern,
		state,
		b.components.quantum,
		b.components.field,
		b.components.energy,
	); err != nil {
		return WrapError(err, ErrCodeTransform, "transform failed")
	}

	// 更新状态
	return b.stateManager.UpdateState()
}

// GetState 获取模型状态
func (b *BaseFlowModel) GetState() ModelState {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.stateManager.GetModelState()
}

// GetSystemState 获取系统状态
func (b *BaseFlowModel) GetSystemState() SystemState {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.stateManager.GetSystemState()
}

// initializeState 初始化状态
func (b *BaseFlowModel) initializeState() error {
	// 初始化量子态
	b.components.quantum.Initialize()

	// 初始化场
	b.components.field.Initialize()

	// 初始化能量系统
	// 为 EnergySystem 添加初始化方法
	initialState := map[core.EnergyType]float64{
		core.PotentialEnergy: 0,
		core.KineticEnergy:   0,
		core.ThermalEnergy:   0,
		core.FieldEnergy:     0,
	}
	if err := b.components.energy.TransformEnergy(initialState); err != nil {
		return err
	}

	// 更新状态
	return b.stateManager.UpdateState()
}

// Close 关闭模型
func (b *BaseFlowModel) Close() error {
	if err := b.Stop(); err != nil {
		return err
	}

	// 清理资源
	b.components.quantum = nil
	b.components.field = nil
	b.components.energy = nil

	return nil
}

// 以下是内部辅助方法

// validateState 验证状态
func (b *BaseFlowModel) validateState() error {
	if b.components.energy.GetTotalEnergy() > b.capacity {
		return NewModelError(ErrCodeState, "energy exceeds capacity", nil)
	}
	return nil
}

// checkRunning 检查运行状态
func (b *BaseFlowModel) checkRunning() error {
	if !b.running {
		return NewModelError(ErrCodeOperation, "model not running", nil)
	}
	return nil
}

// getInternalState 获取内部状态 - 仅供model层使用
func (b *BaseFlowModel) getInternalState() (*core.QuantumState, *core.Field, *core.EnergySystem) {
	return b.components.quantum, b.components.field, b.components.energy
}
