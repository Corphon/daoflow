// system/common/resonance.go

package common

import (
    "time"
)

// ResonanceDetector 共振检测器接口
type ResonanceDetector interface {
    DetectResonance(p1, p2 SharedPattern) (Resonance, error)
    GetDetectionState() ResonanceState
}

// Resonance 共振接口
type Resonance interface {
    GetStrength() float64
    GetFrequency() float64
    GetPhase() float64
    GetCoherence() float64
    GetDuration() time.Duration
}

// ResonanceState 共振状态
type ResonanceState struct {
    ActiveResonances  int
    TotalEnergy       float64
    Stability        float64
    LastUpdate       time.Time
}
