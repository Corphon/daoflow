// model/flow_wuxing.go

package model

import (
    "math"
    "time"

    "github.com/Corphon/daoflow/core"
)

// WuXingConstants 五行常数
const (
    PhaseCount      = 5                // 五行数量
    CycleLength     = 2 * math.Pi      // 周期长度
    GenerateRate    = 0.2              // 相生率
    RestrainRate    = 0.15             // 相克率
    WeakenRate      = 0.1              // 相泄率
    ControlRate     = 0.05             // 相制率
    BalancePoint    = 0.2              // 平衡点(1/5)
)

// WuXingFlow 五行模型
type WuXingFlow struct {
    *BaseFlowModel

    // 五行能量
    energies    map[WuXingPhase]float64
    
    // 量子态组件
    states      map[WuXingPhase]*core.QuantumState
    
    // 场组件
    fields      map[WuXingPhase]*core.Field
    
    // 相位关系
    relationships map[WuXingPhase]map[WuXingPhase]int
}

// NewWuXingFlow 创建五行模型
func NewWuXingFlow() *WuXingFlow {
    base := NewBaseFlowModel(ModelWuXing, 500.0)
    
    wx := &WuXingFlow{
        BaseFlowModel:  base,
        energies:       make(map[WuXingPhase]float64),
        states:         make(map[WuXingPhase]*core.QuantumState),
        fields:         make(map[WuXingPhase]*core.Field),
        relationships:  initializeRelationships(),
    }

    // 初始化五行
    wx.initializePhases()
    
    // 设置初始状态
    wx.state.Phase = PhaseWuXing
    wx.state.Properties["dominant"] = Metal // 默认以金起始

    return wx
}

// initializePhases 初始化五行
func (wx *WuXingFlow) initializePhases() {
    phases := []WuXingPhase{Metal, Wood, Water, Fire, Earth}
    
    for _, phase := range phases {
        // 初始化能量
        wx.energies[phase] = wx.state.Energy / PhaseCount
        
        // 初始化量子态
        wx.states[phase] = core.NewQuantumState()
        wx.states[phase].SetPhase(float64(phase) * CycleLength / PhaseCount)
        
        // 初始化场
        wx.fields[phase] = core.NewField()
        wx.fields[phase].SetStrength(1.0 / PhaseCount)
    }
}

// initializeRelationships 初始化五行关系
func initializeRelationships() map[WuXingPhase]map[WuXingPhase]int {
    relationships := make(map[WuXingPhase]map[WuXingPhase]int)
    
    // 初始化关系映射
    relationships[Metal] = map[WuXingPhase]int{
        Water: RelationGenerate,  // 金生水
        Wood:  RelationRestrain,  // 金克木
        Fire:  RelationWeaken,    // 火克金
        Earth: RelationControl,   // 土生金
    }
    
    relationships[Wood] = map[WuXingPhase]int{
        Fire:  RelationGenerate,  // 木生火
        Earth: RelationRestrain,  // 木克土
        Metal: RelationWeaken,    // 金克木
        Water: RelationControl,   // 水生木
    }
    
    relationships[Water] = map[WuXingPhase]int{
        Wood:  RelationGenerate,  // 水生木
        Fire:  RelationRestrain,  // 水克火
        Earth: RelationWeaken,    // 土克水
        Metal: RelationControl,   // 金生水
    }
    
    relationships[Fire] = map[WuXingPhase]int{
        Earth: RelationGenerate,  // 火生土
        Metal: RelationRestrain,  // 火克金
        Water: RelationWeaken,    // 水克火
        Wood:  RelationControl,   // 木生火
    }
    
    relationships[Earth] = map[WuXingPhase]int{
        Metal: RelationGenerate,  // 土生金
        Water: RelationRestrain,  // 土克水
        Wood:  RelationWeaken,    // 木克土
        Fire:  RelationControl,   // 火生土
    }
    
    return relationships
}

// Transform 五行转换实现
func (wx *WuXingFlow) Transform(pattern TransformPattern) error {
    wx.mu.Lock()
    defer wx.mu.Unlock()

    if !wx.running {
        return NewModelError(ErrCodeOperation, "model not running", nil)
    }

    // 执行转换
    switch pattern {
    case PatternNormal:
        wx.normalTransform()
    case PatternForward:
        wx.generateTransform()
    case PatternReverse:
        wx.restrainTransform()
    case PatternBalance:
        wx.balanceTransform()
    case PatternMutate:
        wx.mutateTransform()
    default:
        return NewModelError(ErrCodeOperation, "invalid transform pattern", nil)
    }

    // 更新量子态
    wx.updateQuantumStates()
    
    // 更新场
    wx.updateFields()
    
    // 更新状态
    wx.updateModelState()

    return nil
}

// normalTransform 常规转换
func (wx *WuXingFlow) normalTransform() {
    // 获取当前主导相位
    dominant := wx.state.Properties["dominant"].(WuXingPhase)
    
    // 计算相生和相克作用
    for target, energy := range wx.energies {
        if relationship, exists := wx.relationships[dominant][target]; exists {
            switch relationship {
            case RelationGenerate:
                wx.transferEnergy(dominant, target, energy*GenerateRate)
            case RelationRestrain:
                wx.transferEnergy(dominant, target, -energy*RestrainRate)
            case RelationWeaken:
                wx.transferEnergy(target, dominant, energy*WeakenRate)
            case RelationControl:
                wx.transferEnergy(target, dominant, -energy*ControlRate)
            }
        }
    }
}

// generateTransform 相生转换
func (wx *WuXingFlow) generateTransform() {
    dominant := wx.state.Properties["dominant"].(WuXingPhase)
    
    // 寻找相生关系
    for target, relationship := range wx.relationships[dominant] {
        if relationship == RelationGenerate {
            wx.transferEnergy(dominant, target, wx.energies[dominant]*GenerateRate)
            wx.state.Properties["dominant"] = target
            break
        }
    }
}

// restrainTransform 相克转换
func (wx *WuXingFlow) restrainTransform() {
    dominant := wx.state.Properties["dominant"].(WuXingPhase)
    
    // 寻找相克关系
    for target, relationship := range wx.relationships[dominant] {
        if relationship == RelationRestrain {
            wx.transferEnergy(dominant, target, -wx.energies[target]*RestrainRate)
            wx.state.Properties["dominant"] = target
            break
        }
    }
}

// balanceTransform 平衡转换
func (wx *WuXingFlow) balanceTransform() {
    totalEnergy := wx.state.Energy
    balanceEnergy := totalEnergy / PhaseCount
    
    for phase := range wx.energies {
        wx.energies[phase] = balanceEnergy
    }
}

// mutateTransform 变异转换
func (wx *WuXingFlow) mutateTransform() {
    // 使用量子涨落
    for phase, state := range wx.states {
        fluctuation := state.GetFluctuation()
        wx.energies[phase] *= (1 + fluctuation)
    }
    
    // 重新归一化
    wx.normalizeEnergies()
}

// transferEnergy 能量转移
func (wx *WuXingFlow) transferEnergy(from, to WuXingPhase, amount float64) {
    if amount > wx.energies[from] {
        amount = wx.energies[from]
    }
    
    wx.energies[from] -= amount
    wx.energies[to] += amount
}

// normalizeEnergies 能量归一化
func (wx *WuXingFlow) normalizeEnergies() {
    totalEnergy := wx.state.Energy
    currentTotal := 0.0
    
    for _, energy := range wx.energies {
        currentTotal += energy
    }
    
    if currentTotal > 0 {
        ratio := totalEnergy / currentTotal
        for phase := range wx.energies {
            wx.energies[phase] *= ratio
        }
    }
}

// updateQuantumStates 更新量子态
func (wx *WuXingFlow) updateQuantumStates() {
    totalEnergy := wx.state.Energy
    
    for phase, state := range wx.states {
        probability := wx.energies[phase] / totalEnergy
        state.SetProbability(probability)
        state.Evolve(phase.String())
    }
    
    // 更新整体量子态
    dominantPhase := wx.state.Properties["dominant"].(WuXingPhase)
    wx.quantum.SetPhase(float64(dominantPhase) * CycleLength / PhaseCount)
    wx.quantum.Evolve("wuxing")
}

// updateFields 更新场
func (wx *WuXingFlow) updateFields() {
    totalEnergy := wx.state.Energy
    
    for phase, field := range wx.fields {
        strength := wx.energies[phase] / totalEnergy
        field.SetStrength(strength)
        field.SetPhase(wx.states[phase].GetPhase())
        field.Evolve()
    }
    
    // 更新统一场
    avgStrength := 0.0
    for _, field := range wx.fields {
        avgStrength += field.GetStrength()
    }
    wx.field.SetStrength(avgStrength / PhaseCount)
    wx.field.Evolve()
}

// updateModelState 更新模型状态
func (wx *WuXingFlow) updateModelState() {
    // 更新状态属性
    for phase, energy := range wx.energies {
        wx.state.Properties[phase.String()] = energy
    }
    
    // 更新主导相位
    maxEnergy := 0.0
    dominant := wx.state.Properties["dominant"].(WuXingPhase)
    
    for phase, energy := range wx.energies {
        if energy > maxEnergy {
            maxEnergy = energy
            dominant = phase
        }
    }
    
    wx.state.Properties["dominant"] = dominant
    wx.state.UpdateTime = time.Now()
}

// GetPhaseEnergy 获取相位能量
func (wx *WuXingFlow) GetPhaseEnergy(phase WuXingPhase) float64 {
    wx.mu.RLock()
    defer wx.mu.RUnlock()
    return wx.energies[phase]
}

// GetDominantPhase 获取主导相位
func (wx *WuXingFlow) GetDominantPhase() WuXingPhase {
    wx.mu.RLock()
    defer wx.mu.RUnlock()
    return wx.state.Properties["dominant"].(WuXingPhase)
}
