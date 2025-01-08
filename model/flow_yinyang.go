// model/flow_yinyang.go

package model

import (
    "math"
    "time"

    "github.com/Corphon/daoflow/core"
)

// YinYangConstants 阴阳常数
const (
    YinYangPeriod     = 2 * math.Pi       // 阴阳周期
    YinYangThreshold  = 0.5               // 阴阳转换阈值
    YinEntropyFactor  = 0.3               // 阴熵因子
    YangEntropyFactor = 0.7               // 阳熵因子
)

// YinYangFlow 阴阳模型
type YinYangFlow struct {
    *BaseFlowModel

    // 阴阳特有状态
    yinEnergy  float64                 // 阴能量
    yangEnergy float64                 // 阳能量
    balance    float64                 // 平衡度
    
    // 量子场组件
    yinField   *core.Field            // 阴场
    yangField  *core.Field            // 阳场
    
    // 量子态组件
    yinState   *core.QuantumState     // 阴量子态
    yangState  *core.QuantumState     // 阳量子态
}

// NewYinYangFlow 创建阴阳模型
func NewYinYangFlow() *YinYangFlow {
    base := NewBaseFlowModel(ModelYinYang, 200.0)
    
    yy := &YinYangFlow{
        BaseFlowModel: base,
        yinField:      core.NewField(),
        yangField:     core.NewField(),
        yinState:      core.NewQuantumState(),
        yangState:     core.NewQuantumState(),
    }

    // 初始化状态
    yy.state.Phase = PhaseYinYang
    yy.state.Nature = NatureYang // 默认从阳开始
    yy.state.Properties["balance"] = 0.5
    
    return yy
}

// Transform 阴阳转换实现
func (yy *YinYangFlow) Transform(pattern TransformPattern) error {
    yy.mu.Lock()
    defer yy.mu.Unlock()

    if !yy.running {
        return NewModelError(ErrCodeOperation, "model not running", nil)
    }

    // 计算阴阳转换
    switch pattern {
    case PatternNormal:
        yy.naturalTransform()
    case PatternForward:
        yy.forwardTransform()
    case PatternReverse:
        yy.reverseTransform()
    case PatternBalance:
        yy.balanceTransform()
    case PatternMutate:
        yy.mutateTransform()
    default:
        return NewModelError(ErrCodeOperation, "invalid transform pattern", nil)
    }

    // 更新量子态
    yy.updateQuantumStates()
    
    // 更新场
    yy.updateFields()
    
    // 更新状态
    yy.updateModelState()

    return nil
}

// naturalTransform 自然转换
func (yy *YinYangFlow) naturalTransform() {
    // 使用量子态演化
    phase := yy.quantum.GetPhase()
    newPhase := math.Mod(phase + YinYangPeriod/360.0, YinYangPeriod)
    yy.quantum.SetPhase(newPhase)
    
    // 计算能量分配
    totalEnergy := yy.state.Energy
    ratio := (math.Sin(newPhase) + 1) / 2 // 转换到 [0,1] 区间
    
    yy.yinEnergy = totalEnergy * (1 - ratio)
    yy.yangEnergy = totalEnergy * ratio
}

// forwardTransform 顺序转换
func (yy *YinYangFlow) forwardTransform() {
    if yy.state.Nature == NatureYin {
        yy.transformToYang()
    } else {
        yy.transformToYin()
    }
}

// reverseTransform 逆序转换
func (yy *YinYangFlow) reverseTransform() {
    if yy.state.Nature == NatureYin {
        yy.transformToYin()
    } else {
        yy.transformToYang()
    }
}

// balanceTransform 平衡转换
func (yy *YinYangFlow) balanceTransform() {
    totalEnergy := yy.state.Energy
    yy.yinEnergy = totalEnergy * 0.5
    yy.yangEnergy = totalEnergy * 0.5
    yy.balance = 1.0
}

// mutateTransform 变异转换
func (yy *YinYangFlow) mutateTransform() {
    // 使用量子涨落
    fluctuation := yy.quantum.GetFluctuation()
    
    // 计算新的能量分配
    totalEnergy := yy.state.Energy
    ratio := 0.5 + fluctuation
    
    yy.yinEnergy = totalEnergy * ratio
    yy.yangEnergy = totalEnergy * (1 - ratio)
}

// transformToYin 转换到阴
func (yy *YinYangFlow) transformToYin() {
    totalEnergy := yy.state.Energy
    transferEnergy := totalEnergy * 0.2 // 每次转换20%
    
    yy.yinEnergy += transferEnergy
    yy.yangEnergy -= transferEnergy
    
    if yy.yinEnergy > yy.yangEnergy {
        yy.state.Nature = NatureYin
    }
}

// transformToYang 转换到阳
func (yy *YinYangFlow) transformToYang() {
    totalEnergy := yy.state.Energy
    transferEnergy := totalEnergy * 0.2 // 每次转换20%
    
    yy.yangEnergy += transferEnergy
    yy.yinEnergy -= transferEnergy
    
    if yy.yangEnergy > yy.yinEnergy {
        yy.state.Nature = NatureYang
    }
}

// updateQuantumStates 更新量子态
func (yy *YinYangFlow) updateQuantumStates() {
    totalEnergy := yy.state.Energy
    
    // 更新阴态
    yinProb := yy.yinEnergy / totalEnergy
    yy.yinState.SetProbability(yinProb)
    yy.yinState.Evolve("yin")
    
    // 更新阳态
    yangProb := yy.yangEnergy / totalEnergy
    yy.yangState.SetProbability(yangProb)
    yy.yangState.Evolve("yang")
    
    // 更新整体量子态
    yy.quantum.SetProbability((yinProb + yangProb) / 2)
    yy.quantum.Evolve("yinyang")
}

// updateFields 更新场
func (yy *YinYangFlow) updateFields() {
    totalEnergy := yy.state.Energy
    
    // 更新阴场
    yy.yinField.SetStrength(yy.yinEnergy / totalEnergy)
    yy.yinField.SetPhase(yy.quantum.GetPhase())
    yy.yinField.Evolve()
    
    // 更新阳场
    yy.yangField.SetStrength(yy.yangEnergy / totalEnergy)
    yy.yangField.SetPhase(yy.quantum.GetPhase() + math.Pi) // 反相位
    yy.yangField.Evolve()
    
    // 更新统一场
    fieldStrength := (yy.yinField.GetStrength() + yy.yangField.GetStrength()) / 2
    yy.field.SetStrength(fieldStrength)
    yy.field.Evolve()
}

// updateModelState 更新模型状态
func (yy *YinYangFlow) updateModelState() {
    // 计算平衡度
    totalEnergy := yy.state.Energy
    if totalEnergy > 0 {
        yinRatio := yy.yinEnergy / totalEnergy
        yangRatio := yy.yangEnergy / totalEnergy
        yy.balance = 1 - math.Abs(yinRatio - yangRatio)
    }

    // 更新状态属性
    yy.state.Properties["yinEnergy"] = yy.yinEnergy
    yy.state.Properties["yangEnergy"] = yy.yangEnergy
    yy.state.Properties["balance"] = yy.balance
    yy.state.Properties["phase"] = yy.quantum.GetPhase()
    yy.state.UpdateTime = time.Now()
}

// GetYinYangRatio 获取阴阳比例
func (yy *YinYangFlow) GetYinYangRatio() (float64, float64) {
    yy.mu.RLock()
    defer yy.mu.RUnlock()
    
    totalEnergy := yy.state.Energy
    if totalEnergy <= 0 {
        return 0.5, 0.5
    }
    
    return yy.yinEnergy/totalEnergy, yy.yangEnergy/totalEnergy
}

// GetBalance 获取平衡度
func (yy *YinYangFlow) GetBalance() float64 {
    yy.mu.RLock()
    defer yy.mu.RUnlock()
    return yy.balance
}
