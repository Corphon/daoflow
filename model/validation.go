// model/validation.go

package model

import "math"

// 系统常量
const (
	// 能量相关常量
	MinEnergy     = 0.0    // 最小能量值
	MaxEnergy     = 1000.0 // 最大能量值
	DefaultEnergy = 0.0    // 默认能量值

	// 容量相关常量
	MinCapacity     = 100.0  // 最小容量
	DefaultCapacity = 1000.0 // 默认容量

	// 阈值常量
	BalanceThreshold = 0.1 // 平衡阈值
	HarmonyThreshold = 0.7 // 和谐阈值
)

// ValidateEnergy 验证能量值
func ValidateEnergy(energy float64) bool {
	// 检查是否是有限数
	if math.IsNaN(energy) || math.IsInf(energy, 0) {
		return false
	}

	// 检查范围
	return energy >= MinEnergy && energy <= MaxEnergy
}

// ValidatePhase 验证相位
func ValidatePhase(phase Phase) bool {
	return phase >= PhaseNone && phase < PhaseMax
}

// ValidateNature 验证属性
func ValidateNature(nature Nature) bool {
	return nature >= NatureNeutral && nature < NatureMax
}

// ValidateModelType 验证模型类型
func ValidateModelType(modelType ModelType) bool {
	return modelType > ModelTypeNone && modelType < ModelTypeMax
}

// ValidateTransformPattern 验证转换模式
func ValidateTransformPattern(pattern TransformPattern) bool {
	return pattern > PatternNone && pattern < PatternMax
}

// ValidateCapacity 验证容量
func ValidateCapacity(capacity float64) bool {
	// 检查是否是有限正数
	if math.IsNaN(capacity) || math.IsInf(capacity, 0) || capacity <= 0 {
		return false
	}

	return capacity >= MinCapacity
}

// ValidateState 验证模型状态
func ValidateState(state *ModelState) error {
	if state == nil {
		return NewModelError(ErrCodeValidation, "state is nil", nil)
	}

	if !ValidateModelType(state.Type) {
		return NewModelError(ErrCodeValidation, "invalid model type", nil)
	}

	if !ValidateEnergy(state.Energy) {
		return NewModelError(ErrCodeValidation, "invalid energy value", nil)
	}

	if !ValidatePhase(state.Phase) {
		return NewModelError(ErrCodeValidation, "invalid phase value", nil)
	}

	if !ValidateNature(state.Nature) {
		return NewModelError(ErrCodeValidation, "invalid nature value", nil)
	}

	return nil
}

// ValidateSystemState 验证系统状态
func ValidateSystemState(state *SystemState) error {
	if state == nil {
		return NewModelError(ErrCodeValidation, "system state is nil", nil)
	}

	if !ValidateEnergy(state.Energy) {
		return NewModelError(ErrCodeValidation, "invalid system energy", nil)
	}

	if state.Entropy < 0 {
		return NewModelError(ErrCodeValidation, "invalid entropy value", nil)
	}

	if state.Harmony < 0 || state.Harmony > 1 {
		return NewModelError(ErrCodeValidation, "invalid harmony value", nil)
	}

	if state.Balance < 0 || state.Balance > 1 {
		return NewModelError(ErrCodeValidation, "invalid balance value", nil)
	}

	return nil
}
