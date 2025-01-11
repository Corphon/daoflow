// system/monitor/trace/tracker.go

package trace

import (
    "context"
    "sync"
    "time"

    "github.com/Corphon/daoflow/system/types"
)

// SpanID 跨度ID
type SpanID string

// TraceID 追踪ID
type TraceID string

// Span 表示一个追踪跨度
type Span struct {
    ID          SpanID                  // 跨度ID
    TraceID     TraceID                 // 追踪ID
    ParentID    SpanID                  // 父跨度ID
    Name        string                  // 跨度名称
    StartTime   time.Time              // 开始时间
    EndTime     time.Time              // 结束时间
    Duration    time.Duration          // 持续时间
    Status      SpanStatus             // 跨度状态
    Tags        map[string]string      // 标签
    Events      []SpanEvent            // 事件列表
    Metrics     map[string]float64     // 指标
    Fields      map[string]interface{} // 字段数据
}

// SpanEvent 跨度事件
type SpanEvent struct {
    Time     time.Time              // 事件时间
    Name     string                 // 事件名称
    Type     string                 // 事件类型
    Fields   map[string]interface{} // 事件字段
}

// SpanStatus 跨度状态
type SpanStatus uint8

const (
    SpanStatusNone SpanStatus = iota
    SpanStatusOK
    SpanStatusError
)

// Tracker 追踪器
type Tracker struct {
    mu sync.RWMutex

    // 配置
    config struct {
        SampleRate    float64        // 采样率
        MaxSpans      int           // 最大跨度数
        BufferSize    int           // 缓冲区大小
        FlushInterval time.Duration // 刷新间隔
    }

    // 活跃跨度
    activeSpans map[SpanID]*Span

    // 跨度通道
    spanChan chan *Span

    // 订阅者
    subscribers []SpanSubscriber

    // 状态
    status struct {
        isRunning bool
        lastFlush time.Time
        errors    []error
    }
}

// SpanSubscriber 跨度订阅者接口
type SpanSubscriber interface {
    OnSpan(*Span) error
}

// NewTracker 创建新的追踪器
func NewTracker(config types.TraceConfig) *Tracker {
    t := &Tracker{
        activeSpans: make(map[SpanID]*Span),
        spanChan:    make(chan *Span, config.BufferSize),
    }

    // 设置配置
    t.config.SampleRate = config.SampleRate
    t.config.MaxSpans = config.MaxSpans
    t.config.BufferSize = config.BufferSize
    t.config.FlushInterval = config.FlushInterval

    return t
}

// Start 启动追踪器
func (t *Tracker) Start(ctx context.Context) error {
    t.mu.Lock()
    if t.status.isRunning {
        t.mu.Unlock()
        return types.NewSystemError(types.ErrRuntime, "tracker already running", nil)
    }
    t.status.isRunning = true
    t.mu.Unlock()

    // 启动处理循环
    go t.processLoop(ctx)

    return nil
}

// Stop 停止追踪器
func (t *Tracker) Stop() error {
    t.mu.Lock()
    defer t.mu.Unlock()

    t.status.isRunning = false
    return nil
}

// StartSpan 开始一个新的跨度
func (t *Tracker) StartSpan(name string, opts ...SpanOption) *Span {
    span := &Span{
        ID:        SpanID(generateID()),
        TraceID:   TraceID(generateID()),
        Name:      name,
        StartTime: time.Now(),
        Status:    SpanStatusNone,
        Tags:      make(map[string]string),
        Events:    make([]SpanEvent, 0),
        Metrics:   make(map[string]float64),
        Fields:    make(map[string]interface{}),
    }

    // 应用选项
    for _, opt := range opts {
        opt(span)
    }

    // 存储活跃跨度
    t.mu.Lock()
    t.activeSpans[span.ID] = span
    t.mu.Unlock()

    return span
}

// EndSpan 结束跨度
func (t *Tracker) EndSpan(span *Span) {
    if span == nil {
        return
    }

    span.EndTime = time.Now()
    span.Duration = span.EndTime.Sub(span.StartTime)

    // 发送跨度
    t.sendSpan(span)

    // 移除活跃跨度
    t.mu.Lock()
    delete(t.activeSpans, span.ID)
    t.mu.Unlock()
}

// AddEvent 添加事件到跨度
func (t *Tracker) AddEvent(span *Span, name string, fields map[string]interface{}) {
    if span == nil {
        return
    }

    event := SpanEvent{
        Time:   time.Now(),
        Name:   name,
        Fields: fields,
    }

    span.Events = append(span.Events, event)
}

// SetTag 设置跨度标签
func (t *Tracker) SetTag(span *Span, key, value string) {
    if span == nil {
        return
    }

    span.Tags[key] = value
}

// SetMetric 设置跨度指标
func (t *Tracker) SetMetric(span *Span, key string, value float64) {
    if span == nil {
        return
    }

    span.Metrics[key] = value
}

// Subscribe 订阅跨度
func (t *Tracker) Subscribe(subscriber SpanSubscriber) {
    t.mu.Lock()
    defer t.mu.Unlock()

    t.subscribers = append(t.subscribers, subscriber)
}

// processLoop 处理循环
func (t *Tracker) processLoop(ctx context.Context) {
    ticker := time.NewTicker(t.config.FlushInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case span := <-t.spanChan:
            t.processSpan(span)
        case <-ticker.C:
            t.flush()
        }
    }
}

// sendSpan 发送跨度
func (t *Tracker) sendSpan(span *Span) {
    // 采样检查
    if !t.shouldSample() {
        return
    }

    select {
    case t.spanChan <- span:
    default:
        // 缓冲区满时记录错误
        t.recordError(types.NewSystemError(types.ErrOverflow, "span buffer full", nil))
    }
}

// processSpan 处理跨度
func (t *Tracker) processSpan(span *Span) {
    // 通知订阅者
    for _, subscriber := range t.subscribers {
        if err := subscriber.OnSpan(span); err != nil {
            t.recordError(err)
        }
    }
}

// flush 刷新所有活跃跨度
func (t *Tracker) flush() {
    t.mu.Lock()
    defer t.mu.Unlock()

    t.status.lastFlush = time.Now()

    // 结束所有超时的活跃跨度
    for id, span := range t.activeSpans {
        if time.Since(span.StartTime) > t.config.FlushInterval {
            t.EndSpan(span)
            delete(t.activeSpans, id)
        }
    }
}

// shouldSample 是否应该采样
func (t *Tracker) shouldSample() bool {
    return rand.Float64() < t.config.SampleRate
}

// recordError 记录错误
func (t *Tracker) recordError(err error) {
    t.mu.Lock()
    defer t.mu.Unlock()

    t.status.errors = append(t.status.errors, err)
}

// SpanOption 跨度选项函数类型
type SpanOption func(*Span)

// WithParent 设置父跨度
func WithParent(parentID SpanID) SpanOption {
    return func(s *Span) {
        s.ParentID = parentID
    }
}

// WithTags 设置标签
func WithTags(tags map[string]string) SpanOption {
    return func(s *Span) {
        for k, v := range tags {
            s.Tags[k] = v
        }
    }
}

// generateID 生成唯一ID
func generateID() string {
    // 实现唯一ID生成逻辑
    return uuid.New().String()
}
