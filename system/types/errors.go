// system/types/errors.go

package types

import (
    "fmt"
    "runtime"
    "strings"
    "time"
)

// ErrorCode 错误码
type ErrorCode string

// 系统层错误码
const (
    // 基础错误
    ErrNone          ErrorCode = ""
    ErrInternal      ErrorCode = "internal"      // 内部错误
    ErrInvalid       ErrorCode = "invalid"       // 无效参数
    ErrNotFound      ErrorCode = "not_found"     // 未找到
    ErrExists        ErrorCode = "exists"        // 已存在
    
    // 初始化错误
    ErrInitSystem    ErrorCode = "init_system"   // 系统初始化错误
    ErrInitMeta      ErrorCode = "init_meta"     // 元系统初始化错误
    ErrInitEvolution ErrorCode = "init_evolution" // 演化系统初始化错误
    
    // 运行时错误
    ErrRuntime       ErrorCode = "runtime"       // 运行时错误
    ErrTimeout       ErrorCode = "timeout"       // 超时错误
    ErrOverflow      ErrorCode = "overflow"      // 溢出错误
    ErrUnderflow     ErrorCode = "underflow"     // 下溢错误
    
    // 状态错误
    ErrState        ErrorCode = "state"         // 状态错误
    ErrTransition   ErrorCode = "transition"    // 转换错误
    ErrValidation   ErrorCode = "validation"    // 验证错误
    
    // 资源错误
    ErrResource     ErrorCode = "resource"      // 资源错误
    ErrCapacity     ErrorCode = "capacity"      // 容量错误
    ErrExhausted    ErrorCode = "exhausted"     // 资源耗尽
    
    // 同步错误
    ErrSync         ErrorCode = "sync"          // 同步错误
    ErrDeadlock     ErrorCode = "deadlock"      // 死锁错误
    ErrRace         ErrorCode = "race"          // 竞争条件
    
    // 配置错误
    ErrConfig       ErrorCode = "config"        // 配置错误
    ErrParse        ErrorCode = "parse"         // 解析错误
    ErrValidate     ErrorCode = "validate"      // 校验错误
)

// SystemError 系统错误
type SystemError struct {
    Code      ErrorCode           // 错误码
    Message   string             // 错误消息
    Details   string             // 详细信息
    Cause     error              // 原因错误
    Stack     []string           // 错误堆栈
    Time      time.Time          // 错误时间
    Layer     SystemLayer        // 错误发生层
    Context   map[string]string  // 错误上下文
}

// Error 实现 error 接口
func (e *SystemError) Error() string {
    var b strings.Builder
    
    // 构建错误消息
    b.WriteString(fmt.Sprintf("[%s] %s", e.Code, e.Message))
    
    // 添加层级信息
    if e.Layer != LayerNone {
        b.WriteString(fmt.Sprintf(" (Layer: %v)", e.Layer))
    }
    
    // 添加详细信息
    if e.Details != "" {
        b.WriteString(fmt.Sprintf("\nDetails: %s", e.Details))
    }
    
    // 添加原因错误
    if e.Cause != nil {
        b.WriteString(fmt.Sprintf("\nCaused by: %v", e.Cause))
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

// NewSystemError 创建系统错误
func NewSystemError(code ErrorCode, message string, cause error) *SystemError {
    return &SystemError{
        Code:     code,
        Message:  message,
        Cause:    cause,
        Stack:    captureStack(),
        Time:     time.Now(),
        Context:  make(map[string]string),
    }
}

// WithDetails 添加详细信息
func (e *SystemError) WithDetails(details string) *SystemError {
    e.Details = details
    return e
}

// WithLayer 设置错误层级
func (e *SystemError) WithLayer(layer SystemLayer) *SystemError {
    e.Layer = layer
    return e
}

// WithContext 添加上下文信息
func (e *SystemError) WithContext(key, value string) *SystemError {
    e.Context[key] = value
    return e
}

// WrapError 包装错误
func WrapError(err error, code ErrorCode, message string) *SystemError {
    if err == nil {
        return nil
    }
    
    // 如果已经是 SystemError，则添加到堆栈
    if sysErr, ok := err.(*SystemError); ok {
        return &SystemError{
            Code:     code,
            Message:  message,
            Cause:    sysErr,
            Stack:    append(captureStack(), sysErr.Stack...),
            Time:     time.Now(),
            Context:  make(map[string]string),
        }
    }
    
    // 创建新的 SystemError
    return NewSystemError(code, message, err)
}

// IsSystemError 检查是否为系统错误
func IsSystemError(err error) bool {
    _, ok := err.(*SystemError)
    return ok
}

// GetErrorCode 获取错误码
func GetErrorCode(err error) ErrorCode {
    if sysErr, ok := err.(*SystemError); ok {
        return sysErr.Code
    }
    return ErrNone
}

// GetErrorLayer 获取错误层级
func GetErrorLayer(err error) SystemLayer {
    if sysErr, ok := err.(*SystemError); ok {
        return sysErr.Layer
    }
    return LayerNone
}

// captureStack 捕获堆栈信息
func captureStack() []string {
    const maxDepth = 32
    stack := make([]string, 0, maxDepth)
    
    // 跳过错误处理相关的帧
    skip := 2
    
    // 收集堆栈信息
    for i := skip; i < maxDepth+skip; i++ {
        pc, file, line, ok := runtime.Caller(i)
        if !ok {
            break
        }
        
        // 获取函数名
        fn := runtime.FuncForPC(pc)
        if fn == nil {
            break
        }
        
        // 构建堆栈帧信息
        frame := fmt.Sprintf("%s:%d %s", file, line, fn.Name())
        stack = append(stack, frame)
    }
    
    return stack
}

// ErrorHandler 错误处理器接口
type ErrorHandler interface {
    Handle(error) error
}

// 错误处理相关的常量
const (
    MaxStackDepth    = 32    // 最大堆栈深度
    MaxErrorHistory  = 100   // 最大错误历史记录数
    MaxRetryAttempts = 3     // 最大重试次数
)
