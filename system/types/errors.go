// system/types/errors.go

package types

import (
    "fmt"
    "runtime"
    "strings"
    "time"

    "github.com/Corphon/daoflow/model"
)

// 复用 model 包的错误码
type ErrorCode = model.ErrorCode

// 扩展系统层错误码
const (
    // 系统层错误码从 10000 开始，避免与 model 层冲突
    ErrSystem        ErrorCode = 10000 + iota  // 系统错误
    ErrComponent                               // 组件错误
    ErrMonitor                                 // 监控错误
    ErrResource                                // 资源错误
    ErrConfig                                  // 配置错误
    ErrValidation                              // 验证错误
    ErrPermission                              // 权限错误
    ErrSecurity                                // 安全错误
    ErrNetwork                                 // 网络错误
    ErrStorage                                 // 存储错误
)

// SystemError 系统错误结构
type SystemError struct {
    ModelErr *model.ModelError  // 包含模型层错误
    Code     ErrorCode         // 系统层错误码
    Layer    SystemLayer       // 错误发生层
    Message  string           // 错误消息
    Details  string           // 详细信息
    Time     time.Time        // 错误发生时间
    Stack    []string         // 错误堆栈
    Context  map[string]any   // 错误上下文
}

// Error 实现 error 接口
func (e *SystemError) Error() string {
    var b strings.Builder

    // 构建错误消息
    b.WriteString(fmt.Sprintf("[System Error %d] ", e.Code))
    if e.Layer != LayerNone {
        b.WriteString(fmt.Sprintf("[Layer: %v] ", e.Layer))
    }
    b.WriteString(e.Message)

    // 添加模型层错误信息
    if e.ModelErr != nil {
        b.WriteString(fmt.Sprintf("\nModel Error: %v", e.ModelErr))
    }

    // 添加详细信息
    if e.Details != "" {
        b.WriteString(fmt.Sprintf("\nDetails: %s", e.Details))
    }

    // 添加上下文信息
    if len(e.Context) > 0 {
        b.WriteString("\nContext:")
        for k, v := range e.Context {
            b.WriteString(fmt.Sprintf("\n  %s: %v", k, v))
        }
    }

    // 添加堆栈信息
    if len(e.Stack) > 0 {
        b.WriteString("\nStack trace:")
        for i, frame := range e.Stack {
            b.WriteString(fmt.Sprintf("\n  %d: %s", i+1, frame))
        }
    }

    return b.String()
}

// NewSystemError 创建新的系统错误
func NewSystemError(code ErrorCode, message string, modelErr error) *SystemError {
    var mErr *model.ModelError
    if modelErr != nil {
        if me, ok := modelErr.(*model.ModelError); ok {
            mErr = me
        }
    }

    return &SystemError{
        ModelErr: mErr,
        Code:    code,
        Message: message,
        Time:    time.Now(),
        Stack:   captureStack(),
        Context: make(map[string]any),
    }
}

// WithLayer 设置错误层级
func (e *SystemError) WithLayer(layer SystemLayer) *SystemError {
    e.Layer = layer
    return e
}

// WithDetails 添加详细信息
func (e *SystemError) WithDetails(details string) *SystemError {
    e.Details = details
    return e
}

// WithContext 添加上下文信息
func (e *SystemError) WithContext(key string, value any) *SystemError {
    e.Context[key] = value
    return e
}

// WrapError 包装错误
func WrapError(err error, code ErrorCode, message string) *SystemError {
    if err == nil {
        return nil
    }

    sysErr := NewSystemError(code, message, err)

    // 如果是系统错误，继承其上下文
    if se, ok := err.(*SystemError); ok {
        for k, v := range se.Context {
            sysErr.Context[k] = v
        }
        sysErr.Layer = se.Layer
    }

    return sysErr
}

// captureStack 捕获堆栈信息
func captureStack() []string {
    const maxDepth = 32
    stack := make([]string, 0, maxDepth)
    
    // 跳过错误处理相关的帧
    skip := 2
    
    for i := skip; i < maxDepth+skip; i++ {
        pc, file, line, ok := runtime.Caller(i)
        if !ok {
            break
        }
        
        fn := runtime.FuncForPC(pc)
        if fn == nil {
            break
        }
        
        frame := fmt.Sprintf("%s:%d %s", file, line, fn.Name())
        stack = append(stack, frame)
    }
    
    return stack
}

// IsSystemError 检查是否为系统错误
func IsSystemError(err error) bool {
    _, ok := err.(*SystemError)
    return ok
}

// IsModelError 检查是否为模型错误
func IsModelError(err error) bool {
    if se, ok := err.(*SystemError); ok {
        return se.ModelErr != nil
    }
    return model.IsModelError(err)
}

// GetErrorCode 获取错误码
func GetErrorCode(err error) ErrorCode {
    if se, ok := err.(*SystemError); ok {
        return se.Code
    }
    return model.GetErrorCode(err)
}

// GetErrorLayer 获取错误层级
func GetErrorLayer(err error) SystemLayer {
    if se, ok := err.(*SystemError); ok {
        return se.Layer
    }
    return LayerNone
}

// ErrorHandler 系统错误处理器
type ErrorHandler interface {
    Handle(error) error
    GetPriority() Priority
}

// DefaultErrorHandler 默认错误处理器
type DefaultErrorHandler struct {
    priority Priority
    retries  int
}

// NewDefaultErrorHandler 创建默认错误处理器
func NewDefaultErrorHandler(priority Priority) *DefaultErrorHandler {
    return &DefaultErrorHandler{
        priority: priority,
        retries:  3,
    }
}

// Handle 处理错误
func (h *DefaultErrorHandler) Handle(err error) error {
    if err == nil {
        return nil
    }

    // 如果是模型错误，交给模型层处理
    if IsModelError(err) {
        return model.GetErrorHandler().Handle(err)
    }

    // 处理系统错误
    if se, ok := err.(*SystemError); ok {
        // 根据错误层级选择处理策略
        switch se.Layer {
        case LayerMeta:
            return h.handleMetaError(se)
        case LayerEvolution:
            return h.handleEvolutionError(se)
        case LayerControl:
            return h.handleControlError(se)
        case LayerResource:
            return h.handleResourceError(se)
        case LayerMonitor:
            return h.handleMonitorError(se)
        default:
            return h.handleGenericError(se)
        }
    }

    // 包装未知错误
    return NewSystemError(ErrSystem, "unknown error occurred", err)
}

// GetPriority 获取处理器优先级
func (h *DefaultErrorHandler) GetPriority() Priority {
    return h.priority
}

// 内部错误处理方法
func (h *DefaultErrorHandler) handleMetaError(err *SystemError) error {
    // TODO: 实现元系统错误处理逻辑
    return err
}

func (h *DefaultErrorHandler) handleEvolutionError(err *SystemError) error {
    // TODO: 实现演化系统错误处理逻辑
    return err
}

func (h *DefaultErrorHandler) handleControlError(err *SystemError) error {
    // TODO: 实现控制系统错误处理逻辑
    return err
}

func (h *DefaultErrorHandler) handleResourceError(err *SystemError) error {
    // TODO: 实现资源系统错误处理逻辑
    return err
}

func (h *DefaultErrorHandler) handleMonitorError(err *SystemError) error {
    // TODO: 实现监控系统错误处理逻辑
    return err
}

func (h *DefaultErrorHandler) handleGenericError(err *SystemError) error {
    // TODO: 实现通用错误处理逻辑
    return err
}
