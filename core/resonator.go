// core/resonator.go

package core

import (
    "math"
    "math/cmplx"
    "sync"
    "time"
)

// ResonatorConfig 共振器配置
type ResonatorConfig struct {
    BaseFrequency    float64  // 基础频率
    DecayRate        float64  // 衰减率
    CoherenceLength  int      // 相干长度
    MaxHistorySize   int      // 历史记录最大长度
}

// Resonator 共振器
type Resonator struct {
    mu sync.RWMutex

    // 状态参数
    amplitude   float64             // 振幅
    frequency   float64             // 频率
    phase      float64             // 相位
    energy     float64             // 能量
    coherence  float64             // 相干度
    resonance  float64             // 共振强度
    
    // 时间相关
    startTime   time.Time          // 开始时间
    lastUpdate  time.Time          // 最后更新时间
    
    // 配置
    config     *ResonatorConfig
    
    // 历史记录
    history    []ResonanceState
}

// ResonanceState 共振状态
type ResonanceState struct {
    Amplitude  float64
    Frequency  float64
    Phase      float64
    Energy     float64
    Timestamp  time.Time
}

// NewResonator 创建共振器
func NewResonator() *Resonator {
    return &Resonator{
        config: &ResonatorConfig{
            BaseFrequency:   1.0,
            DecayRate:       0.01,
            CoherenceLength: 100,
            MaxHistorySize:  1000,
        },
        history: make([]ResonanceState, 0),
        startTime: time.Now(),
        lastUpdate: time.Now(),
    }
}

// Initialize 初始化共振器
func (r *Resonator) Initialize() error {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.amplitude = 0.0
    r.frequency = r.config.BaseFrequency
    r.phase = 0.0
    r.energy = 0.0
    r.coherence = 1.0
    r.resonance = 0.0
    
    r.startTime = time.Now()
    r.lastUpdate = time.Now()
    r.history = make([]ResonanceState, 0)
    
    return nil
}

// Update 更新共振器状态
func (r *Resonator) Update() error {
    r.mu.Lock()
    defer r.mu.Unlock()

    now := time.Now()
    dt := now.Sub(r.lastUpdate).Seconds()

    // 更新相位
    r.phase += 2 * math.Pi * r.frequency * dt
    r.phase = math.Mod(r.phase, 2*math.Pi)

    // 更新振幅（考虑衰减）
    r.amplitude *= math.Exp(-r.config.DecayRate * dt)

    // 更新能量
    r.energy = 0.5 * r.amplitude * r.amplitude

    // 更新相干度
    r.updateCoherence()

    // 更新共振强度
    r.updateResonance()

    // 记录状态
    r.recordState()

    r.lastUpdate = now
    return nil
}

// GetResonance 获取共振强度
func (r *Resonator) GetResonance() float64 {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return r.resonance
}

// 内部方法
func (r *Resonator) updateCoherence() {
    if len(r.history) < 2 {
        return
    }

    var sumPhase complex128
    for _, state := range r.history {
        sumPhase += cmplx.Rect(1.0, state.Phase)
    }

    r.coherence = cmplx.Abs(sumPhase) / float64(len(r.history))
}

func (r *Resonator) updateResonance() {
    r.resonance = r.amplitude * r.coherence
}

func (r *Resonator) recordState() {
    state := ResonanceState{
        Amplitude: r.amplitude,
        Frequency: r.frequency,
        Phase:     r.phase,
        Energy:    r.energy,
        Timestamp: time.Now(),
    }

    r.history = append(r.history, state)
    if len(r.history) > r.config.MaxHistorySize {
        r.history = r.history[1:]
    }
}
