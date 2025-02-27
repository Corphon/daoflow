// model/events.go

package model

import (
	"fmt"
	"time"
)

// ModelEventType 模型事件类型
type ModelEventType string

const (
	EventStateChange   ModelEventType = "state_change"   // 状态变更
	EventTransform     ModelEventType = "transform"      // 模型转换
	EventEnergyChange  ModelEventType = "energy_change"  // 能量变化
	EventPhaseShift    ModelEventType = "phase_shift"    // 相位转移
	EventFieldChange   ModelEventType = "field_change"   // 场变化
	EventQuantumChange ModelEventType = "quantum_change" // 量子态变化
	EventEmergence     ModelEventType = "emergence"      // 涌现现象
	EventResonance     ModelEventType = "resonance"      // 共振现象
)

// ModelEvent 模型事件
type ModelEvent struct {
	ID        string                 `json:"id"`         // 事件ID
	Type      ModelEventType         `json:"type"`       // 事件类型
	ModelType ModelType              `json:"model_type"` // 模型类型
	State     ModelState             `json:"state"`      // 模型状态
	Changes   []ModelStateChange     `json:"changes"`    // 状态变更
	Timestamp time.Time              `json:"timestamp"`  // 发生时间
	Source    string                 `json:"source"`     // 事件源
	Target    string                 `json:"target"`     // 事件目标
	Details   map[string]interface{} `json:"details"`    // 详细信息
	// 新增字段支持模型特性
	Energy float64 `json:"energy"` // 能量值
	Phase  Phase   `json:"phase"`  // 相位
	Nature Nature  `json:"nature"` // 属性
}

// ModelStateChange 模型状态变更
type ModelStateChange struct {
	Field     string      `json:"field"`     // 变更字段
	OldValue  interface{} `json:"old_value"` // 原值
	NewValue  interface{} `json:"new_value"` // 新值
	Timestamp time.Time   `json:"timestamp"` // 变更时间
}

// ModelEventHandler 模型事件处理器接口
type ModelEventHandler interface {
	HandleModelEvent(event ModelEvent) error
	GetEventTypes() []ModelEventType
}

// ModelEventEmitter 模型事件发射器接口
type ModelEventEmitter interface {
	EmitEvent(event ModelEvent) error
	AddHandler(handler ModelEventHandler) error
	RemoveHandler(handler ModelEventHandler) error
}

// BaseModelEventHandler 基础模型事件处理器
type BaseModelEventHandler struct {
	types []ModelEventType
}

// ---------------------------------------------
func (h *BaseModelEventHandler) GetEventTypes() []ModelEventType {
	return h.types
}

// NewModelEvent 创建新的模型事件
func NewModelEvent(eventType ModelEventType, modelType ModelType, state ModelState) ModelEvent {
	return ModelEvent{
		ID:        generateEventID(),
		Type:      eventType,
		ModelType: modelType,
		State:     state,
		Timestamp: time.Now(),
		Changes:   make([]ModelStateChange, 0),
		Details:   make(map[string]interface{}),
	}
}

// AddStateChange 添加状态变更
func (e *ModelEvent) AddStateChange(field string, oldValue, newValue interface{}) {
	change := ModelStateChange{
		Field:     field,
		OldValue:  oldValue,
		NewValue:  newValue,
		Timestamp: time.Now(),
	}
	e.Changes = append(e.Changes, change)
}

// generateEventID 生成事件ID
func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}
