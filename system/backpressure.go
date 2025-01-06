// system/backpressure.go
package system

import (
    "time"
)

type BackpressureHandler struct {
    dropPolicy    DropPolicy
    retryStrategy RetryStrategy
    metrics       *QueueMetrics
}

type DropPolicy int

const (
    DropOldest DropPolicy = iota
    DropNewest
    DropRandom
)

type RetryStrategy struct {
    MaxRetries    int
    RetryInterval time.Duration
    BackoffFactor float64
}

type QueueMetrics struct {
    DroppedCount   int64
    RetryCount     int64
    SuccessRate    float64
    AverageLatency time.Duration
}

func NewBackpressureHandler(policy DropPolicy, strategy RetryStrategy) *BackpressureHandler {
    return &BackpressureHandler{
        dropPolicy:    policy,
        retryStrategy: strategy,
        metrics:      &QueueMetrics{},
    }
}

func (bh *BackpressureHandler) handleOverflow(event interface{}) error {
    switch bh.dropPolicy {
    case DropOldest:
        return bh.dropOldestAndRetry(event)
    case DropNewest:
        return bh.dropNewestAndRetry(event)
    default:
        return bh.dropRandomAndRetry(event)
    }
}
