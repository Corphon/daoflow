// model/flow_bagua.go

package model

import (
    "sync"
    "time"

    "github.com/Corphon/daoflow/core"
)

// BaGuaConstants 八卦常数
const (
    MaxTrigramEnergy = 12.5   // 每个卦象最大能量
    ChangeThreshold  = 0.2    // 变化阈值
    ResonanceRate   = 0.08   // 共振率
)

// Trigram 卦象
type Trigram uint8

const (
    Qian Trigram = iota // 乾 ☰
    Dui                 // 兑 ☱
    Li                  // 离 ☲
    Zhen               // 震 ☳
    Xun                // 巽 ☴
    Kan                // 坎 ☵
    Gen                // 艮 ☶
    Kun                // 坤 ☷
)

// BaGuaFlow 八卦模型
type BaGuaFlow struct {
    *BaseFlowModel // 继承基础模型

    // 八卦状态 - 内部使用
    state struct {
        trigrams    map[Trigram]*TrigramState
        resonance   float64
        harmony     float64
        changes     []Change
    }

    // 内部组件 - 使用 core 层功能
    components struct {
        fields      map[Trigram]*core.Field        // 卦象场
        states      map[Trigram]*core.QuantumState // 量子态
        resonator   *core.Resonator               // 共振器
        correlator  *core.Correlator              // 关联器
    }

    mu sync.RWMutex
}

// TrigramState 卦象状态
type TrigramState struct {
    Energy     float64
    Lines      [3]bool    // 三爻状态
    Resonance  float64
    Relations  map[Trigram]float64
    LastChange time.Time
}

// Change 变化记录
type Change struct {
    From      Trigram
    To        Trigram
    Type      ChangeType
    Timestamp time.Time
}

// ChangeType 变化类型
type ChangeType uint8

const (
    NoChange ChangeType = iota
    NaturalChange      // 自然变化
    ForcedChange       // 强制变化
    ResonantChange     // 共振变化
)

// NewBaGuaFlow 创建八卦模型
func NewBaGuaFlow() *BaGuaFlow {
    // 创建基础模型
    base := NewBaseFlowModel(ModelBaGua, MaxTrigramEnergy*8)

    // 创建八卦模型
    flow := &BaGuaFlow{
        BaseFlowModel: base,
    }

    // 初始化状态
    flow.state.trigrams = make(map[Trigram]*TrigramState)
    flow.initializeTrigrams()

    // 初始化组件
    flow.initializeComponents()

    return flow
}

// initializeTrigrams 初始化卦象
func (f *BaGuaFlow) initializeTrigrams() {
    trigrams := []Trigram{Qian, Dui, Li, Zhen, Xun, Kan, Gen, Kun}
    for _, tri := range trigrams {
        f.state.trigrams[tri] = &TrigramState{
            Energy:    MaxTrigramEnergy / 8,
            Lines:     [3]bool{},
            Resonance: 0,
            Relations: make(map[Trigram]float64),
            LastChange: time.Now(),
        }
    }
}

// initializeComponents 初始化组件
func (f *BaGuaFlow) initializeComponents() {
    // 初始化场和量子态
    f.components.fields = make(map[Trigram]*core.Field)
    f.components.states = make(map[Trigram]*core.QuantumState)

    for tri := range f.state.trigrams {
        f.components.fields[tri] = core.NewField()
        f.components.states[tri] = core.NewQuantumState()
    }

    // 初始化共振器和关联器
    f.components.resonator = core.NewResonator()
    f.components.correlator = core.NewCorrelator()
}

// Start 启动模型
func (f *BaGuaFlow) Start() error {
    if err := f.BaseFlowModel.Start(); err != nil {
        return err
    }

    return f.initializeBaGua()
}

// initializeBaGua 初始化八卦
func (f *BaGuaFlow) initializeBaGua() error {
    f.mu.Lock()
    defer f.mu.Unlock()

    // 初始化场
    for _, field := range f.components.fields {
        if err := field.Initialize(); err != nil {
            return WrapError(err, ErrCodeOperation, "failed to initialize field")
        }
    }

    // 初始化量子态
    for _, state := range f.components.states {
        if err := state.Initialize(); err != nil {
            return WrapError(err, ErrCodeOperation, "failed to initialize quantum state")
        }
    }

    // 初始化共振器
    return f.components.resonator.Initialize()
}

// Transform 执行八卦转换
func (f *BaGuaFlow) Transform(pattern TransformPattern) error {
    if err := f.BaseFlowModel.Transform(pattern); err != nil {
        return err
    }

    f.mu.Lock()
    defer f.mu.Unlock()

    switch pattern {
    case PatternBalance:
        return f.balanceTrigrams()
    case PatternForward:
        return f.naturalTransform()
    case PatternReverse:
        return f.resonantTransform()
    default:
        return f.adaptiveTransform()
    }
}

// balanceTrigrams 平衡卦象
func (f *BaGuaFlow) balanceTrigrams() error {
    // 计算总能量
    totalEnergy := 0.0
    for _, state := range f.state.trigrams {
        totalEnergy += state.Energy
    }

    // 平均分配能量
    averageEnergy := totalEnergy / 8
    for tri, state := range f.state.trigrams {
        state.Energy = averageEnergy
        if err := f.components.states[tri].SetEnergy(averageEnergy); err != nil {
            return err
        }
    }

    return f.updateTrigramStates()
}

// naturalTransform 自然变化
func (f *BaGuaFlow) naturalTransform() error {
    for tri, state := range f.state.trigrams {
        // 检查能量阈值
        if state.Energy > ChangeThreshold*MaxTrigramEnergy {
            // 寻找相应变化卦
            targetTri := f.findNaturalChange(tri)
            if targetTri != tri {
                if err := f.changeTrigram(tri, targetTri, NaturalChange); err != nil {
                    return err
                }
            }
        }
    }
    return f.updateTrigramStates()
}

// resonantTransform 共振变化
func (f *BaGuaFlow) resonantTransform() error {
    // 更新共振器
    if err := f.components.resonator.Update(); err != nil {
        return err
    }

    // 检查共振条件
    resonatingPairs := f.findResonatingPairs()
    for _, pair := range resonatingPairs {
        if err := f.applyResonance(pair[0], pair[1]); err != nil {
            return err
        }
    }

    return f.updateTrigramStates()
}

// adaptiveTransform 适应性变化
func (f *BaGuaFlow) adaptiveTransform() error {
    // 计算系统熵
    entropy := f.calculateSystemEntropy()

    // 根据熵值调整变化
    if entropy > ChangeThreshold {
        return f.naturalTransform()
    }
    return f.resonantTransform()
}

// updateTrigramStates 更新卦象状态
func (f *BaGuaFlow) updateTrigramStates() error {
    // 更新量子态
    for tri, state := range f.components.states {
        if err := state.Update(); err != nil {
            return err
        }
        // 更新卦象能量
        f.state.trigrams[tri].Energy = state.GetEnergy()
    }

    // 更新共振度
    f.state.resonance = f.components.resonator.GetResonance()

    // 更新和谐度
    f.state.harmony = f.calculateHarmony()

    return nil
}

// Close 关闭模型
func (f *BaGuaFlow) Close() error {
    f.mu.Lock()
    defer f.mu.Unlock()

    // 清理组件
    for tri := range f.components.fields {
        f.components.fields[tri] = nil
        f.components.states[tri] = nil
    }
    f.components.resonator = nil
    f.components.correlator = nil

    return f.BaseFlowModel.Close()
}
