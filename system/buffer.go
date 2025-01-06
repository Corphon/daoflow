// system/buffer.go
package system

import (
    "sync"
    "time"
)

// BufferMetrics 缓冲区指标
type BufferMetrics struct {
    Utilization    float64
    DropRate       float64
    Latency        time.Duration
    LastResize     time.Time
}

// DynamicBuffer 动态缓冲区
type DynamicBuffer struct {
    mu            sync.RWMutex
    capacity      int
    threshold     float64
    metrics       *BufferMetrics
    buffer        chan interface{}
    resizePolicy  ResizePolicy
}

// ResizePolicy 缓冲区调整策略
type ResizePolicy struct {
    MinCapacity     int
    MaxCapacity     int
    GrowthFactor    float64
    ShrinkFactor    float64
    ResizeInterval  time.Duration
}

func NewDynamicBuffer(initialCap int, policy ResizePolicy) *DynamicBuffer {
    return &DynamicBuffer{
        capacity:     initialCap,
        threshold:    0.75, // 75% 作为默认阈值
        metrics:      &BufferMetrics{},
        buffer:      make(chan interface{}, initialCap),
        resizePolicy: policy,
    }
}

func (db *DynamicBuffer) adjustCapacity() {
    db.mu.Lock()
    defer db.mu.Unlock()
    
    utilization := float64(len(db.buffer)) / float64(db.capacity)
    
    switch {
    case utilization > db.threshold && db.capacity < db.resizePolicy.MaxCapacity:
        newCap := int(float64(db.capacity) * db.resizePolicy.GrowthFactor)
        if newCap > db.resizePolicy.MaxCapacity {
            newCap = db.resizePolicy.MaxCapacity
        }
        db.resize(newCap)
        
    case utilization < db.threshold/2 && db.capacity > db.resizePolicy.MinCapacity:
        newCap := int(float64(db.capacity) * db.resizePolicy.ShrinkFactor)
        if newCap < db.resizePolicy.MinCapacity {
            newCap = db.resizePolicy.MinCapacity
        }
        db.resize(newCap)
    }
}

func (db *DynamicBuffer) resize(newCapacity int) {
    newBuffer := make(chan interface{}, newCapacity)
    close(db.buffer)
    db.buffer = newBuffer
    db.capacity = newCapacity
    db.metrics.LastResize = time.Now()
}
