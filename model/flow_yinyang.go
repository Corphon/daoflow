// model/flow_yinyang.go

package model

import (
	"math"
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
)

// YinYangConstants 阴阳常数
const (
	MaxYinYangEnergy = 100.0 // 最大能量
	TransformRate    = 0.05  // 转换率
	MinPolarity      = -1.0  // 最小极性
	MaxPolarity      = 1.0   // 最大极性
)

// YinYangFlow 阴阳模型
type YinYangFlow struct {
	*BaseFlowModel // 继承基础模型

	// 内部状态 - 对外隐藏实现
	state struct {
		yinEnergy  float64 // 阴能量
		yangEnergy float64 // 阳能量
		polarity   float64 // 极性
		balance    float64 // 平衡度
	}

	// 内部组件 - 使用 core 层功能
	components struct {
		yinField    *core.Field        // 阴场
		yangField   *core.Field        // 阳场
		yinState    *core.QuantumState // 阴量子态
		yangState   *core.QuantumState // 阳量子态
		interaction *core.Interaction  // 相互作用
	}

	mu sync.RWMutex
}

// NewYinYangFlow 创建阴阳模型
func NewYinYangFlow() *YinYangFlow {
	// 创建基础模型
	base := NewBaseFlowModel(ModelYinYang, MaxYinYangEnergy)

	// 创建阴阳模型
	flow := &YinYangFlow{
		BaseFlowModel: base,
	}

	// 初始化内部组件
	flow.initializeComponents()

	return flow
}

// initializeComponents 初始化组件
func (f *YinYangFlow) initializeComponents() {
	// 创建场
	f.components.yinField = core.NewField(core.ScalarField, 3)
	f.components.yangField = core.NewField(core.ScalarField, 3)

	// 创建量子态
	f.components.yinState = core.NewQuantumState()
	f.components.yangState = core.NewQuantumState()

	// 创建相互作用
	f.components.interaction = core.NewInteraction()
}

// Start 启动模型
func (f *YinYangFlow) Start() error {
	if err := f.BaseFlowModel.Start(); err != nil {
		return err
	}

	// 初始化内部状态
	return f.initializeYinYang()
}

// initializeYinYang 初始化阴阳
func (f *YinYangFlow) initializeYinYang() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 初始化能量
	totalEnergy := f.BaseFlowModel.GetState().Energy
	f.state.yinEnergy = totalEnergy / 2
	f.state.yangEnergy = totalEnergy / 2
	f.state.polarity = 0
	f.state.balance = 1.0

	// 初始化场
	if err := f.components.yinField.Initialize(); err != nil {
		return WrapError(err, ErrCodeOperation, "failed to initialize yin field")
	}
	if err := f.components.yangField.Initialize(); err != nil {
		return WrapError(err, ErrCodeOperation, "failed to initialize yang field")
	}

	// 初始化量子态
	if err := f.components.yinState.Initialize(); err != nil {
		return WrapError(err, ErrCodeOperation, "failed to initialize yin state")
	}
	if err := f.components.yangState.Initialize(); err != nil {
		return WrapError(err, ErrCodeOperation, "failed to initialize yang state")
	}

	return nil
}

// Transform 执行阴阳转换
func (f *YinYangFlow) Transform(pattern TransformPattern) error {
	if err := f.BaseFlowModel.Transform(pattern); err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	switch pattern {
	case PatternBalance:
		return f.balanceTransform()
	case PatternForward:
		return f.yinToYangTransform()
	case PatternReverse:
		return f.yangToYinTransform()
	default:
		return f.naturalTransform()
	}
}

// balanceTransform 平衡转换
func (f *YinYangFlow) balanceTransform() error {
	// 计算总能量
	totalEnergy := f.state.yinEnergy + f.state.yangEnergy

	// 均衡分配
	f.state.yinEnergy = totalEnergy / 2
	f.state.yangEnergy = totalEnergy / 2
	f.state.polarity = 0

	// 更新量子态
	if err := f.components.yinState.Reset(); err != nil {
		return err
	}
	if err := f.components.yangState.Reset(); err != nil {
		return err
	}

	// 更新场
	if err := f.components.yinField.Reset(); err != nil {
		return err
	}
	if err := f.components.yangField.Reset(); err != nil {
		return err
	}

	return f.updateState()
}

// yinToYangTransform 阴转阳
func (f *YinYangFlow) yinToYangTransform() error {
	// 计算转换量
	transferAmount := f.state.yinEnergy * TransformRate

	// 执行转换
	f.state.yinEnergy -= transferAmount
	f.state.yangEnergy += transferAmount

	// 更新极性
	f.state.polarity = math.Min(f.state.polarity+TransformRate, MaxPolarity)

	// 更新量子态
	return f.updateQuantumStates()
}

// yangToYinTransform 阳转阴
func (f *YinYangFlow) yangToYinTransform() error {
	// 计算转换量
	transferAmount := f.state.yangEnergy * TransformRate

	// 执行转换
	f.state.yangEnergy -= transferAmount
	f.state.yinEnergy += transferAmount

	// 更新极性
	f.state.polarity = math.Max(f.state.polarity-TransformRate, MinPolarity)

	// 更新量子态
	return f.updateQuantumStates()
}

// naturalTransform 自然转换
func (f *YinYangFlow) naturalTransform() error {
	// 计算能量差异
	energyDiff := math.Abs(f.state.yinEnergy - f.state.yangEnergy)

	// 如果差异小于阈值，保持平衡
	if energyDiff < BalanceThreshold {
		return nil
	}

	// 根据极性决定转换方向
	if f.state.polarity > 0 {
		return f.yinToYangTransform()
	}
	return f.yangToYinTransform()
}

// updateQuantumStates 更新量子态
func (f *YinYangFlow) updateQuantumStates() error {
	// 更新阴量子态
	if err := f.components.yinState.SetEnergy(f.state.yinEnergy); err != nil {
		return err
	}

	// 更新阳量子态
	if err := f.components.yangState.SetEnergy(f.state.yangEnergy); err != nil {
		return err
	}

	// 更新相互作用
	return f.components.interaction.Update(
		f.components.yinState,
		f.components.yangState,
	)
}

// updateState 更新状态
func (f *YinYangFlow) updateState() error {
	// 计算平衡度
	f.state.balance = 1 - math.Abs(f.state.polarity)

	// 更新基础状态
	modelState := f.GetState()
	modelState.Energy = f.state.yinEnergy + f.state.yangEnergy
	modelState.Phase = f.determinePhase()
	modelState.Nature = f.determineNature()
	modelState.UpdateTime = time.Now()

	return nil
}

// determinePhase 确定相位
func (f *YinYangFlow) determinePhase() Phase {
	if f.state.polarity > 0 {
		return PhaseYang
	}
	return PhaseYin
}

// determineNature 确定属性
func (f *YinYangFlow) determineNature() Nature {
	if math.Abs(f.state.polarity) < BalanceThreshold {
		return NatureNeutral
	}
	if f.state.polarity > 0 {
		return NatureYang
	}
	return NatureYin
}

// Close 关闭模型
func (f *YinYangFlow) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 清理内部组件
	f.components.yinField = nil
	f.components.yangField = nil
	f.components.yinState = nil
	f.components.yangState = nil
	f.components.interaction = nil

	return f.BaseFlowModel.Close()
}

// AdjustEnergy 调整阴阳能量
func (f *YinYangFlow) AdjustEnergy(delta float64) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 首先调用基类方法
	if err := f.BaseFlowModel.AdjustEnergy(delta); err != nil {
		return err
	}

	// 按照阴阳比例分配能量
	totalEnergy := f.state.yinEnergy + f.state.yangEnergy
	deltaYin := delta * (f.state.yinEnergy / totalEnergy)
	deltaYang := delta * (f.state.yangEnergy / totalEnergy)

	f.state.yinEnergy += deltaYin
	f.state.yangEnergy += deltaYang

	// 更新内部组件
	if err := f.updateQuantumStates(); err != nil {
		return err
	}

	return nil
}
