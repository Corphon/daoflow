// model/flow_integrate.go

package model

import (
	"log"
	"math"
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
)

// IntegrateConstants 集成常数
const (
	IntegrateSyncRate  = 0.2  // 同步率
	IntegrateBalance   = 0.25 // 平衡系数
	ResonanceThreshold = 0.8  // 共振阈值
)

// IntegrateFlow 集成流模型
type IntegrateFlow struct {
	*BaseFlowModel
	mu sync.RWMutex

	// 子模型
	yinyang *YinYangFlow
	wuxing  *WuXingFlow
	bagua   *BaGuaFlow
	ganzhi  *GanZhiFlow

	// 统一场
	unifiedField *core.Field

	// 量子纠缠态
	entangledState *core.QuantumState

	// 系统状态
	systemState SystemState
}

// ---------------------------------------------
// NewIntegrateFlow 创建集成流模型
func NewIntegrateFlow() *IntegrateFlow {
	base := NewBaseFlowModel(ModelIntegrate, 2000.0)

	// 创建子模型
	yinyang := NewYinYangFlow()
	wuxing := NewWuXingFlow()
	bagua := NewBaGuaFlow()
	ganzhi := NewGanZhiFlow()

	return &IntegrateFlow{
		BaseFlowModel: base,
		yinyang:       yinyang,
		wuxing:        wuxing,
		bagua:         bagua,
		ganzhi:        ganzhi,
		// 使用基础模型中的场
		unifiedField: base.components.field,
		// 使用基础模型中的量子态
		entangledState: base.components.quantum,
		systemState: SystemState{
			Energy:    0,
			Entropy:   0,
			Harmony:   1,
			Balance:   1,
			Phase:     PhaseNone,
			Timestamp: time.Now(),
		},
	}
}

// Start 启动集成模型
func (im *IntegrateFlow) Start() error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if im.running {
		return NewModelError(ErrCodeOperation, "model already started", nil)
	}

	// 启动子模型
	if err := im.yinyang.Start(); err != nil {
		return err
	}
	if err := im.wuxing.Start(); err != nil {
		return err
	}
	if err := im.bagua.Start(); err != nil {
		return err
	}
	if err := im.ganzhi.Start(); err != nil {
		return err
	}

	// 初始化统一场
	im.unifiedField.Initialize()

	// 初始化量子纠缠态
	im.entangledState.Initialize()

	im.running = true
	return nil
}

// Stop 停止集成模型
func (im *IntegrateFlow) Stop() error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if !im.running {
		return NewModelError(ErrCodeOperation, "model not running", nil)
	}

	// 停止子模型
	if err := im.yinyang.Stop(); err != nil {
		return err
	}
	if err := im.wuxing.Stop(); err != nil {
		return err
	}
	if err := im.bagua.Stop(); err != nil {
		return err
	}
	if err := im.ganzhi.Stop(); err != nil {
		return err
	}

	im.running = false
	return nil
}

// Transform 集成转换
func (im *IntegrateFlow) Transform(pattern TransformPattern) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if !im.running {
		return NewModelError(ErrCodeOperation, "model not running", nil)
	}

	// 转换子模型
	if err := im.yinyang.Transform(pattern); err != nil {
		return err
	}
	if err := im.wuxing.Transform(pattern); err != nil {
		return err
	}
	if err := im.bagua.Transform(pattern); err != nil {
		return err
	}
	if err := im.ganzhi.Transform(pattern); err != nil {
		return err
	}

	// 同步子模型
	im.synchronizeModels()

	// 更新量子态
	im.updateQuantumStates()

	// 更新场
	im.updateFields()

	// 更新系统状态
	im.updateSystemState()

	return nil
}

// synchronizeModels 同步子模型
func (im *IntegrateFlow) synchronizeModels() {
	// 阴阳与五行同步
	yinYangState := im.yinyang.GetState()
	wuxingState := im.wuxing.GetState()

	syncEnergy := math.Min(yinYangState.Energy, wuxingState.Energy) * IntegrateSyncRate
	im.yinyang.AdjustEnergy(syncEnergy)
	im.wuxing.AdjustEnergy(syncEnergy)

	// 八卦与干支同步
	baguaState := im.bagua.GetState()
	ganzhiState := im.ganzhi.GetState()

	syncEnergy = math.Min(baguaState.Energy, ganzhiState.Energy) * IntegrateSyncRate
	im.bagua.AdjustEnergy(syncEnergy)
	im.ganzhi.AdjustEnergy(syncEnergy)
}

// updateQuantumStates 更新量子态
func (im *IntegrateFlow) updateQuantumStates() {
	// 更新纠缠态
	yinYangProb := im.yinyang.GetState().Energy / im.capacity
	wuxingProb := im.wuxing.GetState().Energy / im.capacity
	baguaProb := im.bagua.GetState().Energy / im.capacity
	ganzhiProb := im.ganzhi.GetState().Energy / im.capacity

	avgProb := (yinYangProb + wuxingProb + baguaProb + ganzhiProb) / 4
	im.entangledState.SetProbability(avgProb)
	im.entangledState.Evolve("integrate")
}

// updateFields 更新场
func (im *IntegrateFlow) updateFields() {
	// 使用状态计算场强度
	yinYangState := im.yinyang.GetState()
	wuxingState := im.wuxing.GetState()
	baguaState := im.bagua.GetState()
	ganzhiState := im.ganzhi.GetState()

	// 计算总场强度
	totalEnergy := yinYangState.Energy + wuxingState.Energy +
		baguaState.Energy + ganzhiState.Energy
	averageStrength := totalEnergy / 4.0

	// 更新统一场
	im.unifiedField.SetStrength(averageStrength)
	im.unifiedField.SetPhase(im.entangledState.GetPhase())
	im.unifiedField.Evolve()
}

// updateSystemState 更新系统状态
func (im *IntegrateFlow) updateSystemState() {
	// 计算总能量
	im.systemState.Energy = im.yinyang.GetState().Energy +
		im.wuxing.GetState().Energy +
		im.bagua.GetState().Energy +
		im.ganzhi.GetState().Energy

	// 计算熵
	im.systemState.Entropy = im.calculateSystemEntropy()

	// 计算和谐度
	im.systemState.Harmony = im.calculateSystemHarmony()

	// 计算平衡度
	im.systemState.Balance = im.calculateSystemBalance()

	// 更新时间戳
	im.systemState.Timestamp = time.Now()

	// 通过状态管理器更新状态
	if im.stateManager != nil {
		// 使用 GetModelState 获取当前状态
		currentState := im.stateManager.GetModelState()

		// 更新状态值
		currentState.Energy = im.systemState.Energy
		currentState.Properties["entropy"] = im.systemState.Entropy
		currentState.Properties["harmony"] = im.systemState.Harmony
		currentState.Properties["balance"] = im.systemState.Balance
		currentState.UpdateTime = im.systemState.Timestamp

		// 使用 UpdateState 保存更新后的状态
		if err := im.stateManager.UpdateState(); err != nil {
			// 这里可以处理错误，或者记录日志
			// 为了保持与原代码一致，我们暂时不返回错误
			log.Printf("Failed to update state: %v", err)
		}
	}
}

// calculateSystemEntropy 计算系统熵
func (im *IntegrateFlow) calculateSystemEntropy() float64 {
	if im.systemState.Energy <= 0 {
		return 0
	}

	// 使用量子态计算熵
	return -im.entangledState.GetProbability() * math.Log(im.entangledState.GetProbability())
}

// calculateSystemHarmony 计算系统和谐度
func (im *IntegrateFlow) calculateSystemHarmony() float64 {
	// 基于场强度计算和谐度
	fieldStrength := im.unifiedField.GetStrength()
	return math.Min(1.0, fieldStrength/ResonanceThreshold)
}

// calculateSystemBalance 计算系统平衡度
func (im *IntegrateFlow) calculateSystemBalance() float64 {
	if im.systemState.Energy <= 0 {
		return 1
	}

	// 计算各子系统能量比例的方差
	totalEnergy := im.systemState.Energy
	energyRatios := []float64{
		im.yinyang.GetState().Energy / totalEnergy,
		im.wuxing.GetState().Energy / totalEnergy,
		im.bagua.GetState().Energy / totalEnergy,
		im.ganzhi.GetState().Energy / totalEnergy,
	}

	variance := 0.0
	meanRatio := 0.25 // 理想平均比例
	for _, ratio := range energyRatios {
		diff := ratio - meanRatio
		variance += diff * diff
	}
	variance /= 4

	// 转换为平衡度（0-1）
	return 1 - math.Min(1, variance/IntegrateBalance)
}

// GetSystemState 获取系统状态
func (im *IntegrateFlow) GetSystemState() SystemState {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return im.systemState
}

// Close 关闭集成模型
func (im *IntegrateFlow) Close() error {
	if err := im.Stop(); err != nil {
		return err
	}

	// 关闭子模型
	if err := im.yinyang.Close(); err != nil {
		return err
	}
	if err := im.wuxing.Close(); err != nil {
		return err
	}
	if err := im.bagua.Close(); err != nil {
		return err
	}
	if err := im.ganzhi.Close(); err != nil {
		return err
	}

	return im.BaseFlowModel.Close()
}

// GetCoreState 获取核心状态
func (im *IntegrateFlow) GetCoreState() CoreState {
	im.mu.RLock()
	defer im.mu.RUnlock()

	// 创建CoreState
	coreState := CoreState{
		QuantumState: im.entangledState,
		FieldState:   im.unifiedField,
		EnergyState:  im.components.energy,
		Phase:        float64(im.systemState.Phase),
		Properties: map[string]float64{
			"harmony": im.systemState.Harmony,
			"balance": im.systemState.Balance,
			"entropy": im.systemState.Entropy,
			"yinyang": im.yinyang.GetState().Energy,
			"wuxing":  im.wuxing.GetState().Energy,
			"bagua":   im.bagua.GetState().Energy,
			"ganzhi":  im.ganzhi.GetState().Energy,
		},
	}

	return coreState
}

// UpdateCoreState 更新核心状态
func (im *IntegrateFlow) UpdateCoreState(state CoreState) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	// 更新量子态
	if state.QuantumState != nil {
		*im.entangledState = *state.QuantumState
	}

	// 更新场
	if state.FieldState != nil {
		*im.unifiedField = *state.FieldState
	}

	// 更新能量系统
	if im.components.energy != nil && state.EnergyState != nil {
		*im.components.energy = *state.EnergyState
	}

	// 更新相位
	im.systemState.Phase = Phase(state.Phase)

	// 更新属性
	if state.Properties != nil {
		if harmony, ok := state.Properties["harmony"]; ok {
			im.systemState.Harmony = harmony
		}
		if balance, ok := state.Properties["balance"]; ok {
			im.systemState.Balance = balance
		}
		if entropy, ok := state.Properties["entropy"]; ok {
			im.systemState.Entropy = entropy
		}
	}

	return nil
}

// ValidateCoreState 验证核心状态
func (im *IntegrateFlow) ValidateCoreState() error {
	im.mu.RLock()
	defer im.mu.RUnlock()

	// 验证量子态
	if im.entangledState == nil {
		return NewModelError(ErrCodeValidation, "nil quantum state", nil)
	}

	// 验证场
	if im.unifiedField == nil {
		return NewModelError(ErrCodeValidation, "nil field", nil)
	}

	// 验证能量系统
	if im.components.energy == nil {
		return NewModelError(ErrCodeValidation, "nil energy system", nil)
	}

	// 验证系统状态
	if im.systemState.Energy < 0 || im.systemState.Energy > MaxSystemEnergy {
		return NewModelError(ErrCodeValidation, "invalid energy value", nil)
	}

	if im.systemState.Harmony < 0 || im.systemState.Harmony > 1 {
		return NewModelError(ErrCodeValidation, "invalid harmony value", nil)
	}

	if im.systemState.Balance < 0 || im.systemState.Balance > 1 {
		return NewModelError(ErrCodeValidation, "invalid balance value", nil)
	}

	if im.systemState.Entropy < 0 {
		return NewModelError(ErrCodeValidation, "invalid entropy value", nil)
	}

	return nil
}

// model/flow_integrate.go 文件

// GetYinYangFlow 获取阴阳流模型
func (im *IntegrateFlow) GetYinYangFlow() *YinYangFlow {
	im.mu.RLock()
	defer im.mu.RUnlock()

	return im.yinyang
}

// GetBaGuaFlow 获取八卦流模型
func (im *IntegrateFlow) GetBaGuaFlow() *BaGuaFlow {
	im.mu.RLock()
	defer im.mu.RUnlock()

	return im.bagua
}

// GetGanZhiFlow 获取干支流模型
func (im *IntegrateFlow) GetGanZhiFlow() *GanZhiFlow {
	im.mu.RLock()
	defer im.mu.RUnlock()

	return im.ganzhi
}

// GetWuXingFlow 获取五行流模型
func (im *IntegrateFlow) GetWuXingFlow() *WuXingFlow {
	im.mu.RLock()
	defer im.mu.RUnlock()

	return im.wuxing
}
