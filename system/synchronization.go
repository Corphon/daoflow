// system/synchronization.go

package system

import (
    "math"
    "sync"
    "time"
    "context"

    "github.com/Corphon/daoflow/model"
)

// SyncConstants 同步常数
const (
    MaxPhaseOffset   = math.Pi / 4  // 最大相位偏移
    CouplingStrength = 0.3         // 耦合强度
    SyncThreshold    = 0.8         // 同步阈值
    MinCoherence     = 0.6         // 最小相干度
)

// SynchronizationSystem 同步系统
type SynchronizationSystem struct {
    mu sync.RWMutex

    // 关联系统
    evolution  *EvolutionSystem
    adaptation *AdaptationSystem
    integrate  *model.IntegrateFlow

    // 同步状态
    state struct {
        Coherence    float64            // 相干度
        Phase        map[string]float64 // 各子系统相位
        Frequency    map[string]float64 // 各子系统频率
        Coupling     [][]float64        // 耦合矩阵
    }

    // 同步控制
    control struct {
        reference  OscillatorState  // 参考振荡器
        observers  []Observer       // 观察者列表
        tolerance  float64         // 同步容差
    }

    // 协调器
    coordinator *PhaseCoordinator

    ctx    context.Context
    cancel context.CancelFunc
}

// OscillatorState 振荡器状态
type OscillatorState struct {
    Phase     float64
    Frequency float64
    Amplitude float64
    Energy    float64
}

// PhaseCoordinator 相位协调器
type PhaseCoordinator struct {
    kuramoto *KuramotoNetwork  // Kuramoto网络模型
    phases   map[string]*PhaseState
    lock     sync.RWMutex
}

// PhaseState 相位状态
type PhaseState struct {
    Current   float64
    Target    float64
    Velocity  float64
    Coupling  float64
}

// KuramotoNetwork Kuramoto网络
type KuramotoNetwork struct {
    nodes     []KuramotoNode
    coupling  [][]float64
    natural   []float64
}

// KuramotoNode Kuramoto节点
type KuramotoNode struct {
    phase     float64
    frequency float64
    neighbors []int
}

// NewSynchronizationSystem 创建同步系统
func NewSynchronizationSystem(ctx context.Context, 
    es *EvolutionSystem, 
    as *AdaptationSystem,
    integrate *model.IntegrateFlow) *SynchronizationSystem {
    
    ctx, cancel := context.WithCancel(ctx)

    ss := &SynchronizationSystem{
        evolution:  es,
        adaptation: as,
        integrate:  integrate,
        ctx:       ctx,
        cancel:    cancel,
    }

    // 初始化状态
    ss.initializeState()
    
    // 创建相位协调器
    ss.coordinator = NewPhaseCoordinator()

    go ss.runSynchronization()
    return ss
}

// initializeState 初始化状态
func (ss *SynchronizationSystem) initializeState() {
    ss.state.Phase = make(map[string]float64)
    ss.state.Frequency = make(map[string]float64)
    ss.state.Coupling = make([][]float64, 4) // 4个核心模型
    
    for i := range ss.state.Coupling {
        ss.state.Coupling[i] = make([]float64, 4)
        for j := range ss.state.Coupling[i] {
            if i == j {
                ss.state.Coupling[i][j] = 1.0
            } else {
                ss.state.Coupling[i][j] = CouplingStrength
            }
        }
    }
}

// runSynchronization 运行同步过程
func (ss *SynchronizationSystem) runSynchronization() {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ss.ctx.Done():
            return
        case <-ticker.C:
            ss.synchronize()
        }
    }
}

// synchronize 执行同步
func (ss *SynchronizationSystem) synchronize() {
    ss.mu.Lock()
    defer ss.mu.Unlock()

    // 获取系统状态
    state := ss.integrate.GetSystemState()

    // 更新振荡器状态
    ss.updateOscillators(state)

    // 计算相位差异
    phaseDiffs := ss.calculatePhaseDifferences()

    // 应用Kuramoto模型
    ss.applyKuramotoModel(phaseDiffs)

    // 调整耦合强度
    ss.adjustCoupling()

    // 更新相干度
    ss.updateCoherence()
}

// updateOscillators 更新振荡器
func (ss *SynchronizationSystem) updateOscillators(state model.SystemState) {
    // 更新各子系统的相位和频率
    ss.updateSubsystemOscillator("yinyang", state.YinYang)
    ss.updateSubsystemOscillator("wuxing", state.WuXing...)
    ss.updateSubsystemOscillator("bagua", state.BaGua...)
    ss.updateSubsystemOscillator("ganzhi", state.GanZhi...)
}

// updateSubsystemOscillator 更新子系统振荡器
func (ss *SynchronizationSystem) updateSubsystemOscillator(
    name string, 
    values ...float64) {
    
    // 计算平均能量
    avgEnergy := 0.0
    for _, v := range values {
        avgEnergy += v
    }
    avgEnergy /= float64(len(values))

    // 更新相位
    currentPhase := ss.state.Phase[name]
    naturalFreq := 2 * math.Pi / float64(len(values))
    
    // 计算新相位
    newPhase := currentPhase + naturalFreq*avgEnergy/100.0
    ss.state.Phase[name] = math.Mod(newPhase, 2*math.Pi)
    
    // 更新频率
    ss.state.Frequency[name] = naturalFreq * avgEnergy/100.0
}

// calculatePhaseDifferences 计算相位差异
func (ss *SynchronizationSystem) calculatePhaseDifferences() [][]float64 {
    n := len(ss.state.Phase)
    diffs := make([][]float64, n)
    for i := range diffs {
        diffs[i] = make([]float64, n)
    }

    systems := []string{"yinyang", "wuxing", "bagua", "ganzhi"}
    for i, si := range systems {
        for j, sj := range systems {
            if i != j {
                diff := math.Abs(ss.state.Phase[si] - ss.state.Phase[sj])
                if diff > math.Pi {
                    diff = 2*math.Pi - diff
                }
                diffs[i][j] = diff
            }
        }
    }

    return diffs
}

// applyKuramotoModel 应用Kuramoto模型
func (ss *SynchronizationSystem) applyKuramotoModel(phaseDiffs [][]float64) {
    for i := range ss.state.Phase {
        // 计算相位更新
        deltaPhase := 0.0
        for j := range ss.state.Phase {
            if i != j {
                deltaPhase += ss.state.Coupling[i][j] * 
                    math.Sin(phaseDiffs[i][j])
            }
        }
        
        // 更新相位
        systems := []string{"yinyang", "wuxing", "bagua", "ganzhi"}
        ss.state.Phase[systems[i]] += deltaPhase * ss.control.tolerance
    }
}

// adjustCoupling 调整耦合强度
func (ss *SynchronizationSystem) adjustCoupling() {
    for i := range ss.state.Coupling {
        for j := range ss.state.Coupling[i] {
            if i != j {
                // 基于相干度调整耦合强度
                coherenceFactor := ss.state.Coherence
                if coherenceFactor < MinCoherence {
                    // 增强耦合
                    ss.state.Coupling[i][j] *= 1.1
                } else {
                    // 维持或减弱耦合
                    ss.state.Coupling[i][j] *= 0.99
                }
                
                // 限制耦合强度范围
                ss.state.Coupling[i][j] = math.Max(0.1,
                    math.Min(1.0, ss.state.Coupling[i][j]))
            }
        }
    }
}

// updateCoherence 更新相干度
func (ss *SynchronizationSystem) updateCoherence() {
    // 使用序参数计算相干度
    var sumSin, sumCos float64
    n := float64(len(ss.state.Phase))

    for _, phase := range ss.state.Phase {
        sumSin += math.Sin(phase)
        sumCos += math.Cos(phase)
    }

    // r = |Σ exp(iθ)|/N
    ss.state.Coherence = math.Sqrt(sumSin*sumSin+sumCos*sumCos) / n
}

// GetSynchronizationStatus 获取同步状态
func (ss *SynchronizationSystem) GetSynchronizationStatus() map[string]interface{} {
    ss.mu.RLock()
    defer ss.mu.RUnlock()

    return map[string]interface{}{
        "coherence": ss.state.Coherence,
        "phases":    ss.state.Phase,
        "coupling":  ss.state.Coupling,
    }
}

// Close 关闭同步系统
func (ss *SynchronizationSystem) Close() error {
    ss.cancel()
    return nil
}
