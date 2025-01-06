// api/events.go

package api

import (
    "context"
    "sync"
    "time"

    "github.com/Corphon/daoflow/system"
)

// EventType 事件类型
type EventType string

const (
    // 系统事件
    EventSystemStartup    EventType = "system.startup"     // 系统启动
    EventSystemShutdown   EventType = "system.shutdown"    // 系统关闭
    EventSystemError      EventType = "system.error"       // 系统错误
    
    // 生命周期事件
    EventStateChange      EventType = "lifecycle.state"    // 状态变更
    EventHealthCheck      EventType = "lifecycle.health"   // 健康检查
    
    // 模式事件
    EventPatternDetected  EventType = "pattern.detected"   // 模式检测
    EventPatternLearned   EventType = "pattern.learned"    // 模式学习
    EventPatternEvolved   EventType = "pattern.evolved"    // 模式演化
    
    // 能量事件
    EventEnergyLow        EventType = "energy.low"         // 能量不足
    EventEnergyBalance    EventType = "energy.balance"     // 能量平衡
    EventEnergyOverflow   EventType = "energy.overflow"    // 能量溢出
)

// EventPriority 事件优先级
type EventPriority int

const (
    PriorityLow     EventPriority = 0
    PriorityNormal  EventPriority = 1
    PriorityHigh    EventPriority = 2
    PriorityCritical EventPriority = 3
)

// Event 事件结构
type Event struct {
    ID        string                 `json:"id"`         // 事件ID
    Type      EventType              `json:"type"`       // 事件类型
    Priority  EventPriority          `json:"priority"`   // 优先级
    Source    string                 `json:"source"`     // 事件源
    Timestamp time.Time              `json:"timestamp"`  // 发生时间
    Payload   interface{}            `json:"payload"`    // 事件数据
    Metadata  map[string]interface{} `json:"metadata"`   // 元数据
}

// EventFilter 事件过滤器
type EventFilter struct {
    Types     []EventType    `json:"types"`      // 事件类型
    Priority  EventPriority  `json:"priority"`   // 最低优先级
    Source    string         `json:"source"`     // 事件源
    FromTime  time.Time      `json:"from_time"`  // 起始时间
    ToTime    time.Time      `json:"to_time"`    // 结束时间
}

// Subscription 订阅信息
type Subscription struct {
    ID      string       // 订阅ID
    Filter  EventFilter  // 过滤器
    Channel chan Event   // 事件通道
    Active  bool         // 是否活跃
}

// EventsAPI 事件API
type EventsAPI struct {
    mu sync.RWMutex
    system *system.SystemCore

    // 订阅管理
    subscriptions map[string]*Subscription
    
    // 事件缓存
    eventCache []*Event
    cacheSize  int
    
    // 事件统计
    stats *EventStats
    
    ctx    context.Context
    cancel context.CancelFunc

    eventQueue  *system.PriorityEventQueue
}

// EventStats 事件统计
type EventStats struct {
    TotalEvents     int64                  `json:"total_events"`
    EventsByType    map[EventType]int64    `json:"events_by_type"`
    EventsByPriority map[EventPriority]int64 `json:"events_by_priority"`
    LastEventTime   time.Time              `json:"last_event_time"`
}

// NewEventsAPI 创建事件API实例
func NewEventsAPI(sys *system.SystemCore) *EventsAPI {
    ctx, cancel := context.WithCancel(context.Background())
    
    api := &EventsAPI{
        system:        sys,
        subscriptions: make(map[string]*Subscription),
        eventCache:    make([]*Event, 0),
        cacheSize:     1000, // 默认缓存1000条事件
        stats:         &EventStats{
            EventsByType:     make(map[EventType]int64),
            EventsByPriority: make(map[EventPriority]int64),
        },
        ctx:           ctx,
        cancel:        cancel,
    }
    
    return api
}

// 更新 Subscribe 方法
func (e *EventsAPI) Subscribe(filter EventFilter) (*Subscription, error) {
    bufferConfig := system.ResizePolicy{
        MinCapacity:    100,
        MaxCapacity:    10000,
        GrowthFactor:   2.0,
        ShrinkFactor:   0.5,
        ResizeInterval: time.Minute,
    }
    
    buffer := system.NewDynamicBuffer(500, bufferConfig)
    
    sub := &Subscription{
        ID:      generateSubscriptionID(),
        Filter:  filter,
        Channel: buffer.buffer,
        Active:  true,
    }
    
    return sub, nil
}

// Publish 发布事件
func (e *EventsAPI) Publish(evt Event) error {
    e.mu.Lock()
    defer e.mu.Unlock()

    // 更新统计信息
    e.stats.TotalEvents++
    e.stats.EventsByType[evt.Type]++
    e.stats.EventsByPriority[evt.Priority]++
    e.stats.LastEventTime = evt.Timestamp

    // 缓存事件
    e.cacheEvent(&evt)

    // 分发事件给订阅者
    for _, sub := range e.subscriptions {
        if sub.Active && e.matchFilter(evt, sub.Filter) {
            select {
            case sub.Channel <- evt:
                // 事件已发送
            default:
                // 通道已满，跳过
            }
        }
    }

    return nil
}

// Subscribe 订阅事件
func (e *EventsAPI) Subscribe(filter EventFilter) (*Subscription, error) {
    e.mu.Lock()
    defer e.mu.Unlock()

    sub := &Subscription{
        ID:      generateSubscriptionID(),
        Filter:  filter,
        Channel: make(chan Event, 100),
        Active:  true,
    }

    e.subscriptions[sub.ID] = sub
    return sub, nil
}

// Unsubscribe 取消订阅
func (e *EventsAPI) Unsubscribe(subID string) error {
    e.mu.Lock()
    defer e.mu.Unlock()

    sub, exists := e.subscriptions[subID]
    if !exists {
        return NewError(ErrSubscriptionNotFound, "subscription not found")
    }

    sub.Active = false
    close(sub.Channel)
    delete(e.subscriptions, subID)

    return nil
}

// GetEvents 获取历史事件
func (e *EventsAPI) GetEvents(filter EventFilter) ([]*Event, error) {
    e.mu.RLock()
    defer e.mu.RUnlock()

    var events []*Event
    for _, evt := range e.eventCache {
        if e.matchFilter(*evt, filter) {
            events = append(events, evt)
        }
    }

    return events, nil
}

// GetStats 获取事件统计信息
func (e *EventsAPI) GetStats() *EventStats {
    e.mu.RLock()
    defer e.mu.RUnlock()
    
    return e.stats
}

// matchFilter 检查事件是否匹配过滤器
func (e *EventsAPI) matchFilter(evt Event, filter EventFilter) bool {
    // 检查事件类型
    if len(filter.Types) > 0 {
        typeMatch := false
        for _, t := range filter.Types {
            if evt.Type == t {
                typeMatch = true
                break
            }
        }
        if !typeMatch {
            return false
        }
    }

    // 检查优先级
    if evt.Priority < filter.Priority {
        return false
    }

    // 检查事件源
    if filter.Source != "" && evt.Source != filter.Source {
        return false
    }

    // 检查时间范围
    if !filter.FromTime.IsZero() && evt.Timestamp.Before(filter.FromTime) {
        return false
    }
    if !filter.ToTime.IsZero() && evt.Timestamp.After(filter.ToTime) {
        return false
    }

    return true
}

// cacheEvent 缓存事件
func (e *EventsAPI) cacheEvent(evt *Event) {
    e.eventCache = append(e.eventCache, evt)
    if len(e.eventCache) > e.cacheSize {
        e.eventCache = e.eventCache[1:]
    }
}

// generateSubscriptionID 生成订阅ID
func generateSubscriptionID() string {
    // 实现ID生成逻辑
    return time.Now().Format("20060102150405") + randomString(6)
}

// randomString 生成随机字符串
func randomString(n int) string {
    // 实现随机字符串生成逻辑
    return "random"
}

// Close 关闭API
func (e *EventsAPI) Close() error {
    e.mu.Lock()
    defer e.mu.Unlock()

    e.cancel()
    
    // 关闭所有订阅
    for _, sub := range e.subscriptions {
        sub.Active = false
        close(sub.Channel)
    }
    
    return nil
}
