// model/state.go

package model

import (
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
)

// StateManager 状态管理器
type StateManager struct {
	mu sync.RWMutex

	// 公开状态 - system 层可访问
	modelState  ModelState  // 模型状态
	systemState SystemState // 系统状态

	// 内部状态 - 仅 model 层使用
	internal struct {
		quantum *core.QuantumState // 量子态
		field   *core.Field        // 场
		energy  *core.EnergySystem // 能量
	}

	// 状态转换器
	transformer *StateTransformer
}

// NewStateManager 创建状态管理器
func NewStateManager(modelType ModelType, capacity float64) *StateManager {
	sm := &StateManager{
		modelState: ModelState{
			Type:       modelType,
			Energy:     0,
			Phase:      PhaseNone,
			Nature:     NatureNeutral,
			Properties: make(map[string]interface{}),
			UpdateTime: time.Now(),
		},
		systemState: SystemState{
			Energy:     0,
			Entropy:    0,
			Harmony:    1,
			Balance:    1,
			Phase:      PhaseNone,
			Timestamp:  time.Now(),
			Properties: make(map[string]interface{}),
		},
	}

	// 初始化内部状态
	sm.internal.quantum = core.NewQuantumState()
	sm.internal.field = core.NewField(core.ScalarField, 3) // 使用标量场和三维空间
	sm.internal.energy = core.NewEnergySystem(capacity)

	// 创建状态转换器
	sm.transformer = NewStateTransformer()

	return sm
}

// GetModelState 获取模型状态
func (sm *StateManager) GetModelState() ModelState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.modelState
}

// GetSystemState 获取系统状态
func (sm *StateManager) GetSystemState() SystemState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.systemState
}

// UpdateState 更新状态
func (sm *StateManager) UpdateState() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// 从内部状态更新公开状态
	if err := sm.updateFromInternal(); err != nil {
		return WrapError(err, ErrCodeState, "failed to update from internal state")
	}

	// 更新时间戳
	now := time.Now()
	sm.modelState.UpdateTime = now
	sm.systemState.Timestamp = now

	return nil
}

// updateFromInternal 从内部状态更新
func (sm *StateManager) updateFromInternal() error {
	// 更新能量
	energy := sm.internal.energy.GetTotalEnergy()
	sm.modelState.Energy = energy
	sm.systemState.Energy = energy

	// 更新熵
	entropy := sm.internal.quantum.GetEntropy()
	sm.modelState.Properties["entropy"] = entropy
	sm.systemState.Entropy = entropy

	// 更新场强度相关指标
	fieldStrength := sm.internal.field.GetStrength()
	sm.systemState.Harmony = sm.calculateHarmony(fieldStrength)
	sm.systemState.Balance = sm.calculateBalance()

	// 更新相位
	phase := sm.internal.quantum.GetPhase()
	sm.modelState.Phase = Phase(phase)
	sm.systemState.Phase = sm.modelState.Phase

	// 更新时间戳
	now := time.Now()
	sm.modelState.UpdateTime = now
	sm.systemState.Timestamp = now

	return nil
}

// calculateHarmony 计算和谐度
func (sm *StateManager) calculateHarmony(fieldStrength float64) float64 {
	// 和谐度与场强度和量子相干性相关
	coherence := sm.internal.quantum.GetCoherence()
	return (fieldStrength + coherence) / 2
}

// calculateBalance 计算平衡度
func (sm *StateManager) calculateBalance() float64 {
	// 平衡度由能量分布和场均匀性决定
	energyBalance := sm.internal.energy.GetBalance()
	fieldUniformity := sm.internal.field.GetUniformity()
	return (energyBalance + fieldUniformity) / 2
}

// StateTransformer 状态转换器
type StateTransformer struct {
	rules map[TransformPattern]TransformRule
}

// TransformRule 转换规则
type TransformRule struct {
	Validate  func(ModelState) bool
	Transform func(*core.QuantumState, *core.Field, *core.EnergySystem) error
}

// ---------------------------------------------
// NewStateTransformer 创建状态转换器
func NewStateTransformer() *StateTransformer {
	st := &StateTransformer{
		rules: make(map[TransformPattern]TransformRule),
	}
	st.initializeRules()
	return st
}

// initializeRules 初始化转换规则
func (st *StateTransformer) initializeRules() {
	// 常规转换
	st.rules[PatternNormal] = TransformRule{
		Validate: func(state ModelState) bool {
			return state.Energy > 0
		},
		Transform: func(q *core.QuantumState, f *core.Field, e *core.EnergySystem) error {
			// 执行常规量子演化
			q.Evolve(core.PatternIntegrate) // 使用集成模式演化

			// 更新场
			if err := f.Update(q.GetPhase()); err != nil {
				return err
			}

			return nil
		},
	}

	// 平衡转换
	st.rules[PatternBalance] = TransformRule{
		Validate: func(state ModelState) bool {
			return true
		},
		Transform: func(q *core.QuantumState, f *core.Field, e *core.EnergySystem) error {
			// 重新分配能量
			equalShare := e.GetTotalEnergy() / 4
			energyMap := map[core.EnergyType]float64{
				core.PotentialEnergy: equalShare,
				core.KineticEnergy:   equalShare,
				core.ThermalEnergy:   equalShare,
				core.FieldEnergy:     equalShare,
			}
			if err := e.TransformEnergy(energyMap); err != nil {
				return err
			}

			// 重置量子态
			q.Reset()

			// 重置场
			if err := f.Reset(); err != nil {
				return err
			}

			return nil
		},
	}
}

// ApplyTransform 应用转换
func (st *StateTransformer) ApplyTransform(
	pattern TransformPattern,
	state ModelState,
	quantum *core.QuantumState,
	field *core.Field,
	energy *core.EnergySystem) error {

	rule, exists := st.rules[pattern]
	if !exists {
		return NewModelError(ErrCodeTransform, "unknown transform pattern", nil)
	}

	if !rule.Validate(state) {
		return NewModelError(ErrCodeState, "invalid state for transformation", nil)
	}

	return rule.Transform(quantum, field, energy)
}
