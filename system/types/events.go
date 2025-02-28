//system/types/events.go

package types

import (
	"sync"
	"time"

	"github.com/Corphon/daoflow/model"
)

// 系统事件类型常量
const (
	// 系统事件
	EventStateChange   EventType = "system.state_change"
	EventHealthCheck   EventType = "system.health_check"
	EventMetricsUpdate EventType = "system.metrics_update"

	// 模型事件
	EventModelChange EventType = "model.change"
	EventModelSync   EventType = "model.sync"
	EventModelError  EventType = "model.error"

	// 流程事件
	EventFlowStart    EventType = "flow.start"
	EventFlowComplete EventType = "flow.complete"
	EventFlowError    EventType = "flow.error"

	// 系统事件
	EventSystemStarted  EventType = "system.started"  // 系统启动
	EventSystemStopping EventType = "system.stopping" // 系统停止中
	EventSystemStopped  EventType = "system.stopped"  // 系统已停止
	EventSystemError    EventType = "system.error"    // 系统错误
	EventSystemWarning  EventType = "system.warning"  // 系统警告

	// 组件事件
	EventComponentStarted EventType = "component.started" // 组件启动
	EventComponentStopped EventType = "component.stopped" // 组件停止
	EventComponentError   EventType = "component.error"   // 组件错误

	// 状态事件
	EventStateChanged    EventType = "state.changed"    // 状态改变
	EventStateTransition EventType = "state.transition" // 状态转换

	// 演化事件
	EventEvolutionStarted      EventType = "evolution.started"       // 演化开始
	EventEvolutionCompleted    EventType = "evolution.completed"     // 演化完成
	EventEvolutionStateChanged EventType = "evolution.state_changed" // 演化状态变更
	EventEvolutionPhaseShift   EventType = "evolution.phase_shift"   // 演化相位转换
	EventEvolutionError        EventType = "evolution.error"         // 演化错误

)

// EventPriority 事件优先级
type EventPriority int

// Event 事件基础结构
type Event struct {
	ID        string                 `json:"id"`        // 事件ID
	Type      EventType              `json:"type"`      // 事件类型
	Priority  EventPriority          `json:"priority"`  // 事件优先级
	Source    string                 `json:"source"`    // 事件源
	Timestamp time.Time              `json:"timestamp"` // 事件时间
	Topic     string                 `json:"topic"`     // 事件主题
	Payload   interface{}            `json:"payload"`   // 事件负载
	Metadata  map[string]interface{} `json:"metadata"`  // 元数据
}

// EventProcessor 事件处理接口
type EventProcessor interface {
	// 事件处理
	ProcessModelEvent(event model.ModelEvent) error
	ProcessSystemEvent(event SystemEvent) error

	// 事件订阅
	Subscribe(eventType EventType, handler EventHandler) error
	Unsubscribe(eventType EventType, handler EventHandler) error
}

// Subscription 事件订阅
type Subscription struct {
	ID        string       `json:"id"`         // 订阅ID
	Topic     string       `json:"topic"`      // 订阅主题
	Handler   EventHandler `json:"-"`          // 事件处理器
	Active    bool         `json:"active"`     // 是否活跃
	CreatedAt time.Time    `json:"created_at"` // 创建时间
}

// EventBusImpl 事件总线实现
type EventBusImpl struct {
	mu sync.RWMutex

	// 事件处理器映射
	handlers map[string]EventHandler               // handlerID -> handler
	topics   map[EventType]map[string]EventHandler // eventType -> handlerID -> handler

	// 配置
	config struct {
		bufferSize int           // 事件缓冲区大小
		timeout    time.Duration // 处理超时时间
	}

	// 状态
	status struct {
		running   bool      // 运行状态
		startTime time.Time // 启动时间
	}
}

// ---------------------------------------
// NewEventBus 创建事件总线
func NewEventBus() *EventBusImpl {
	bus := &EventBusImpl{
		handlers: make(map[string]EventHandler),
		topics:   make(map[EventType]map[string]EventHandler),
	}

	// 初始化配置
	bus.config.bufferSize = 1000
	bus.config.timeout = 30 * time.Second

	// 初始化状态
	bus.status.running = true
	bus.status.startTime = time.Now()

	return bus
}

// AddHandler 添加处理器
func (eb *EventBusImpl) AddHandler(handler EventHandler) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	// 验证处理器
	if handler == nil {
		return NewSystemError(ErrInvalid, "nil handler", nil)
	}

	// 注册处理器
	handlerID := handler.GetHandlerID()
	eb.handlers[handlerID] = handler

	// 为每个事件类型建立映射
	for _, eventType := range handler.GetEventTypes() {
		if _, exists := eb.topics[eventType]; !exists {
			eb.topics[eventType] = make(map[string]EventHandler)
		}
		eb.topics[eventType][handlerID] = handler
	}

	return nil
}

// RemoveHandler 移除处理器
func (eb *EventBusImpl) RemoveHandler(handlerID string) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	// 验证处理器存在
	handler, exists := eb.handlers[handlerID]
	if !exists {
		return NewSystemError(ErrNotFound, "handler not found", nil)
	}

	// 从所有主题中移除
	for _, eventType := range handler.GetEventTypes() {
		if handlers, exists := eb.topics[eventType]; exists {
			delete(handlers, handlerID)
		}
	}

	// 从处理器映射中移除
	delete(eb.handlers, handlerID)

	return nil
}

// Unsubscribe 取消订阅
func (eb *EventBusImpl) Unsubscribe(eventType EventType, handler EventHandler) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if handler == nil {
		return NewSystemError(ErrInvalid, "nil handler", nil)
	}

	handlerID := handler.GetHandlerID()

	// 验证处理器存在于指定事件类型
	handlers, exists := eb.topics[eventType]
	if !exists {
		return NewSystemError(ErrNotFound, "event type not found", nil)
	}

	if _, exists := handlers[handlerID]; !exists {
		return NewSystemError(ErrNotFound, "handler not found for event type", nil)
	}

	// 从指定事件类型中移除处理器
	delete(eb.topics[eventType], handlerID)

	// 检查处理器是否还订阅了其他事件类型
	stillSubscribed := false
	for t, handlers := range eb.topics {
		if t != eventType {
			if _, exists := handlers[handlerID]; exists {
				stillSubscribed = true
				break
			}
		}
	}

	// 如果没有其他订阅,从处理器映射中移除
	if !stillSubscribed {
		delete(eb.handlers, handlerID)
	}

	return nil
}

// UnsubscribeMulti 批量取消订阅
func (eb *EventBusImpl) UnsubscribeMulti(types []EventType, handler EventHandler) error {
	for _, eventType := range types {
		if err := eb.Unsubscribe(eventType, handler); err != nil {
			return err
		}
	}
	return nil
}

// GetHandlers 获取所有处理器
func (eb *EventBusImpl) GetHandlers() []EventHandler {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	handlers := make([]EventHandler, 0, len(eb.handlers))
	for _, handler := range eb.handlers {
		handlers = append(handlers, handler)
	}
	return handlers
}

// GetSubscriptions 获取指定事件类型的订阅处理器
func (eb *EventBusImpl) GetSubscriptions(eventType EventType) []EventHandler {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	handlers := make([]EventHandler, 0)
	if topicHandlers, exists := eb.topics[eventType]; exists {
		for _, handler := range topicHandlers {
			handlers = append(handlers, handler)
		}
	}
	return handlers
}

// Publish 发布事件
func (eb *EventBusImpl) Publish(event SystemEvent) error {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if !eb.status.running {
		return NewSystemError(ErrInvalid, "nil handler", nil)
	}

	// 获取该事件类型的所有处理器
	handlers, exists := eb.topics[event.Type]
	if !exists {
		return nil // 没有处理器，直接返回
	}

	// 调用所有相关处理器
	for _, handler := range handlers {
		if handler.ShouldHandle(event) {
			if err := handler.HandleEvent(event); err != nil {
				return err
			}
		}
	}

	return nil
}

// Subscribe implements EventProcessor interface
func (eb *EventBusImpl) Subscribe(eventType EventType, handler EventHandler) error {
	return eb.doSubscribe([]EventType{eventType}, handler)
}

// SubscribeMulti implements EventBus interface for multiple event types
func (eb *EventBusImpl) SubscribeMulti(types []EventType, handler EventHandler) error {
	return eb.doSubscribe(types, handler)
}

// doSubscribe internal implementation of subscription
func (eb *EventBusImpl) doSubscribe(types []EventType, handler EventHandler) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	// 验证处理器
	if handler == nil {
		return NewSystemError(ErrInvalid, "nil handler", nil)
	}

	// 注册处理器
	handlerID := handler.GetHandlerID()
	if _, exists := eb.handlers[handlerID]; !exists {
		eb.handlers[handlerID] = handler
	}

	// 建立事件类型到处理器的映射
	for _, eventType := range types {
		if _, exists := eb.topics[eventType]; !exists {
			eb.topics[eventType] = make(map[string]EventHandler)
		}
		eb.topics[eventType][handlerID] = handler
	}

	return nil
}

// EventBusImpl需要实现ProcessModelEvent方法和ProcessSystemEvent方法
func (eb *EventBusImpl) ProcessModelEvent(event model.ModelEvent) error {
	// 转换为系统事件
	sysEvent := SystemEvent{
		ID:        event.ID,
		Type:      EventModelChange,
		Source:    "model",
		Timestamp: event.Timestamp,
		Data:      event,
	}

	return eb.Publish(sysEvent)
}

func (eb *EventBusImpl) ProcessSystemEvent(event SystemEvent) error {
	return eb.Publish(event)
}
