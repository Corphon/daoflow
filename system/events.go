// system/events.go
package system

import (
    "context"
    "sync"
    "time"
)

type PriorityEventQueue struct {
    mu             sync.RWMutex
    highPriority   *DynamicBuffer
    normalPriority *DynamicBuffer
    lowPriority    *DynamicBuffer
    backpressure   *BackpressureHandler
    metrics        *EventMetrics
}

type EventMetrics struct {
    TotalEvents     int64
    ProcessedEvents int64
    DroppedEvents   int64
    AverageLatency  time.Duration
}

func NewPriorityEventQueue(ctx context.Context) *PriorityEventQueue {
    policy := ResizePolicy{
        MinCapacity:    100,
        MaxCapacity:    10000,
        GrowthFactor:   2.0,
        ShrinkFactor:   0.5,
        ResizeInterval: time.Minute,
    }
    
    return &PriorityEventQueue{
        highPriority:   NewDynamicBuffer(200, policy),
        normalPriority: NewDynamicBuffer(500, policy),
        lowPriority:    NewDynamicBuffer(1000, policy),
        backpressure:   NewBackpressureHandler(DropOldest, RetryStrategy{
            MaxRetries:    3,
            RetryInterval: time.Second,
            BackoffFactor: 2.0,
        }),
        metrics: &EventMetrics{},
    }
}
