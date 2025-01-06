// api/adapter.go
package api

// SystemAdapter 系统适配器
type SystemAdapter struct {
    core          *system.SystemCore
    legacySupport bool
    eventMapper   *EventMapper
    bufferMapper  *BufferMapper
}

// EventMapper 事件映射器
type EventMapper struct {
    oldToNew map[string]string
    newToOld map[string]string
}

// BufferMapper 缓冲区映射器
type BufferMapper struct {
    oldToNew map[string]string
    newToOld map[string]string
}
