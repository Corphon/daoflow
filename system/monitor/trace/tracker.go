// system/monitor/trace/tracker.go

package trace

import (
    "context"
    "sync"
    "time"

    "github.com/Corphon/daoflow/model"
    "github.com/Corphon/daoflow/system/types"
)

// Span 表示一个追踪跨度
type Span struct {
    ID          SpanID                  
    TraceID     TraceID                 
    ParentID    SpanID                  
    Name        string                  
    StartTime   time.Time              
    EndTime     time.Time              
    Duration    time.Duration          
    Status      SpanStatus             
    Tags        map[string]string      
    Events      []SpanEvent            
    Metrics     map[string]float64     
    Fields      map[string]interface{} 

    // 新增模型相关字段
    ModelType   model.ModelType        // 关联的模型类型
    ModelState  model.ModelState       // 相关的模型状态
    ModelFlow   model.FlowState        // 流状态
}

// SpanEvent 跨度事件
type SpanEvent struct {
    Time      time.Time              
    Name      string                 
    Type      string                 
    Fields    map[string]interface{} 
    ModelData *model.EventData       // 新增：模型事件数据
}

// Tracker 追踪器
type Tracker struct {
    mu sync.RWMutex

    // 配置使用 types 包的配置
    config types.TraceConfig

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

    // 新增：模型状态管理器
    modelManager *model.StateManager
}

// SpanSubscriber 跨度订阅者接口
type SpanSubscriber interface {
    OnSpan(*Span) error
    OnModelEvent(model.EventData) error // 新增：处理模型事件
}

// NewTracker 创建新的追踪器
func NewTracker(config types.TraceConfig) *Tracker {
    t := &Tracker{
        config:      config,
        activeSpans: make(map[SpanID]*Span),
        spanChan:    make(chan *Span, config.BufferSize),
        modelManager: model.NewStateManager(),
    }

    return t
}

// Start 启动追踪器
func (t *Tracker) Start(ctx context.Context) error {
    t.mu.Lock()
    if t.status.isRunning {
        t.mu.Unlock()
        return model.WrapError(nil, model.ErrCodeOperation, "tracker already running")
    }
    t.status.isRunning = true
    t.mu.Unlock()

    go t.processLoop(ctx)
    return nil
}

// Stop 停止追踪器
func (t *Tracker) Stop() error {
    t.mu.Lock()
    defer t.mu.Unlock()

    if !t.status.isRunning {
        return model.WrapError(nil, model.ErrCodeOperation, "tracker not running")
    }

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

    // 如果设置了模型类型，获取相应的模型状态
    if span.ModelType != model.ModelTypeNone {
        if state, err := t.modelManager.GetState(span.ModelType); err == nil {
            span.ModelState = state
        }
    }

    // 存储活跃跨度
    t.mu.Lock()
    t.activeSpans[span.ID] = span
    t.mu.Unlock()

    return span
}

// EndSpan 结束跨度
func (t *Tracker) EndSpan(span *Span) error {
    if span == nil {
        return model.WrapError(nil, model.ErrCodeValidation, "nil span")
    }

    span.EndTime = time.Now()
    span.Duration = span.EndTime.Sub(span.StartTime)

    // 更新模型状态
    if span.ModelType != model.ModelTypeNone {
        if err := t.updateModelState(span); err != nil {
            return model.WrapError(err, model.ErrCodeOperation, "failed to update model state")
        }
    }

    // 发送跨度
    if err := t.sendSpan(span); err != nil {
        return model.WrapError(err, model.ErrCodeOperation, "failed to send span")
    }

    // 移除活跃跨度
    t.mu.Lock()
    delete(t.activeSpans, span.ID)
    t.mu.Unlock()

    return nil
}

// AddEvent 添加事件到跨度
func (t *Tracker) AddEvent(span *Span, name string, fields map[string]interface{}, modelData *model.EventData) error {
    if span == nil {
        return model.WrapError(nil, model.ErrCodeValidation, "nil span")
    }

    event := SpanEvent{
        Time:      time.Now(),
        Name:      name,
        Fields:    fields,
        ModelData: modelData,
    }

    span.Events = append(span.Events, event)

    // 如果有模型事件数据，通知订阅者
    if modelData != nil {
        t.notifyModelEvent(*modelData)
    }

    return nil
}

// updateModelState 更新模型状态
func (t *Tracker) updateModelState(span *Span) error {
    // 基于跨度信息更新模型状态
    update := model.StateUpdate{
        Type:     span.ModelType,
        State:    span.ModelState,
        Flow:     span.ModelFlow,
        Time:     time.Now(),
        Duration: span.Duration,
    }

    return t.modelManager.UpdateState(update)
}

// notifyModelEvent 通知模型事件
func (t *Tracker) notifyModelEvent(event model.EventData) {
    t.mu.RLock()
    subscribers := t.subscribers
    t.mu.RUnlock()

    for _, sub := range subscribers {
        if err := sub.OnModelEvent(event); err != nil {
            t.recordError(err)
        }
    }
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
func (t *Tracker) sendSpan(span *Span) error {
    // 采样检查
    if !t.shouldSample() {
        return nil
    }

    select {
    case t.spanChan <- span:
        return nil
    default:
        return model.WrapError(nil, model.ErrCodeResource, "span buffer full")
    }
}

// processSpan 处理跨度
func (t *Tracker) processSpan(span *Span) {
    // 通知订阅者
    t.mu.RLock()
    subscribers := t.subscribers
    t.mu.RUnlock()

    for _, subscriber := range subscribers {
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
            if err := t.EndSpan(span); err != nil {
                t.recordError(err)
            }
            delete(t.activeSpans, id)
        }
    }
}

// recordError 记录错误
func (t *Tracker) recordError(err error) {
    t.mu.Lock()
    defer t.mu.Unlock()

    t.status.errors = append(t.status.errors, err)
}

// SpanOption 跨度选项函数类型
type SpanOption func(*Span)

// WithModelType 设置模型类型
func WithModelType(modelType model.ModelType) SpanOption {
    return func(s *Span) {
        s.ModelType = modelType
    }
}

// WithModelState 设置模型状态
func WithModelState(state model.ModelState) SpanOption {
    return func(s *Span) {
        s.ModelState = state
    }
}

// WithModelFlow 设置流状态
func WithModelFlow(flow model.FlowState) SpanOption {
    return func(s *Span) {
        s.ModelFlow = flow
    }
}

// generateID 生成唯一ID
func generateID() string {
    return fmt.Sprintf("span-%d", time.Now().UnixNano())
}
