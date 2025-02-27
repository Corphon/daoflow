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

	// 基本信息
	ID        string    `json:"id"`         // 共振标识
	Type      string    `json:"type"`       // 共振类型
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间

	// 共振参数
	Amplitude float64 `json:"amplitude"` // 振幅
	Frequency float64 `json:"frequency"` // 频率
	Phase     float64 `json:"phase"`     // 相位
	Energy    float64 `json:"energy"`    // 能量
	Coherence float64 `json:"coherence"` // 相干度
	Resonance float64 `json:"resonance"` // 共振强度
	Stability float64 `json:"stability"` // 稳定性

	// 关联信息
	Source struct {
		ID         string                 `json:"id"`
		Type       string                 `json:"type"`
		Properties map[string]interface{} `json:"properties"`
	} `json:"source"`

	Target struct {
		ID         string                 `json:"id"`
		Type       string                 `json:"type"`
		Properties map[string]interface{} `json:"properties"`
	} `json:"target"`

	// 元数据
	Metadata map[string]interface{} `json:"metadata"` // 附加信息
}

// ResonanceMetrics 共振指标
type ResonanceMetrics struct {
	ActiveCount     int     `json:"active_count"`     // 活跃共振数量
	TotalEnergy     float64 `json:"total_energy"`     // 总能量
	AvgCoherence    float64 `json:"avg_coherence"`    // 平均相干度
	AvgResonance    float64 `json:"avg_resonance"`    // 平均共振强度
	SystemStability float64 `json:"system_stability"` // 系统稳定性
}

// Constants for resonance
const (
	MaxResonanceAge = 24 * time.Hour // 共振状态最大存活时间
)
