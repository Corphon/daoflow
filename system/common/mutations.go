// system/common/mutations.go

package common

import (
    "time"
)

// MutationDetector 突变检测器接口
type MutationDetector interface {
    DetectMutations(pattern SharedPattern) ([]Mutation, error)
    GetDetectionMetrics() DetectionMetrics
}

// Mutation 突变接口
type Mutation interface {
    SharedPattern
    GetSource() string
    GetTarget() string
    GetProbability() float64
    GetChanges() []MutationChange
}

// MutationChange 突变变化
type MutationChange struct {
    Property    string
    OldValue    interface{}
    NewValue    interface{}
    Delta       float64
    Timestamp   time.Time
}

// DetectionMetrics 检测指标
type DetectionMetrics struct {
    TotalDetected  int
    SuccessRate    float64
    AverageLatency time.Duration
    LastDetection  time.Time
}

// MutationHandler 突变处理器接口
type MutationHandler interface {
    HandleMutation(mutation Mutation) error
    GetHandlingMetrics() HandlingMetrics
}

// HandlingMetrics 处理指标
type HandlingMetrics struct {
    TotalHandled   int
    SuccessRate    float64
    AverageLatency time.Duration
    LastHandled    time.Time
}
