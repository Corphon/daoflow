// system/types/errors.go

package types

import (
    "fmt"
    "runtime"
    "strings"
    "time"
)

// ErrorCode 错误码类型
type ErrorCode uint32

// 系统错误码定义
const (
    // 基础错误 (0-999)
    ErrNone        ErrorCode = 0   // 无错误
    ErrUnknown     ErrorCode = 1   // 未知错误
    ErrInternal    ErrorCode = 2   // 内部错误
    ErrInvalid     ErrorCode = 3   // 无效参数
    ErrNotFound    ErrorCode = 4   // 未找到
    ErrExists      ErrorCode = 5   // 已存在
    
    // 初始化错误 (1000-1999)
    ErrInitialize      ErrorCode = 1000 // 初始化错误
    ErrInitConfig      ErrorCode = 1001 // 配置初始化错误
    ErrInitMeta        ErrorCode = 1002 // 元系统初始化错误
    ErrInitEvolution   ErrorCode = 1003 // 演化系统初始化错误
    ErrInitResource    ErrorCode = 1004 // 资源系统初始化错误
    
    // 运行时错误 (2000-2999)
    ErrRuntime         ErrorCode = 2000 // 运行时错误
    ErrTimeout         ErrorCode = 2001 // 超时错误
    ErrOverflow        ErrorCode = 2002 // 溢出错误
    ErrUnderflow       ErrorCode = 2003 // 下溢错误
    ErrDeadlock        ErrorCode = 2004 // 死锁错误
    ErrRace           ErrorCode = 2005 // 竞争条件
    
    // 状态错误 (3000-3999)
    ErrState          ErrorCode = 3000 // 状态错误
    ErrStateTransition ErrorCode = 3001 // 状态转换错误
    ErrStateValidation ErrorCode = 3002 // 状态验证错误
    ErrStateConflict   ErrorCode = 3003 // 状态冲突
    
    // 资源错误 (4000-4999)
    ErrResource       ErrorCode = 4000 // 资源错误
    ErrResourceAlloc  ErrorCode = 4001 // 资源分配错误
    ErrResourceExhaust ErrorCode = 4002 // 资源耗尽
    ErrResourceLimit   ErrorCode = 4003 // 资源限制
    
    // 演化错误 (5000-5999)
    ErrEvolution      ErrorCode = 5000 // 演化错误
    ErrEvolutionPath  ErrorCode = 5001 // 演化路径错误
    ErrEvolutionStuck ErrorCode = 5002 // 演化停滞
    
    // 适应错误 (6000-6999)
    ErrAdaptation     ErrorCode = 6000 // 适应错误
    ErrAdaptFailed    ErrorCode = 6001 // 适应失败
    ErrAdaptTimeout   ErrorCode = 6002 // 适应超时
    
    // 同步错误 (7000-7999)
    ErrSync          ErrorCode = 7000 // 同步错误
    ErrSyncConflict  ErrorCode = 7001 // 同步冲突
    ErrSyncTimeout   ErrorCode = 7002 // 同步超时
    
    // 量子场错误 (8000-8999)
    ErrQuantum       ErrorCode = 8000 // 量子场错误
    ErrQuantumState  ErrorCode = 8001 // 量子态错误
    ErrQuantumCollapse ErrorCode = 8002 // 量子态崩溃
    
    // 涌现错误 (9000-9999)
    ErrEmergence     ErrorCode = 9000 // 涌现错误
    ErrEmergPattern  ErrorCode = 9001 // 涌现模式错误
    ErrEmergFailure  ErrorCode = 9002 // 涌现失败
)

// SystemError 系统错误结构
type SystemError struct {
    Code      ErrorCode           // 错误码
    Message   string             // 错误消息
    Details   string             // 详细信息
    Cause     error              // 原因错误
    Stack     []string           // 错误堆栈
    Time      time.Time          // 错误时间
    Layer     SystemLayer        // 错误发生层
    Context   map[string]string  // 错误上下文
    Severity  IssueSeverity      // 错误严重度
}

// Error 实现 error 接口
func (e *SystemError) Error() string {
    var b strings.Builder
    
    // 构建错误消息
    b.WriteString(fmt.Sprintf("[%d] %s", e.Code, e.Message))
    
    // 添加严重度
    b.WriteString(fmt.Sprintf(" (Severity: %v)", e.Severity))
    
    // 添加层级信息
    if e.Layer != LayerNone {
        b.WriteString(fmt.Sprintf(" [Layer: %v]", e.Layer))
    }
    
    // 添加详细信息
    if e.Details != "" {
        b.WriteString(fmt.Sprintf("\nDetails: %s", e.Details))
    }
    
    // 添加原因错误
    if e.Cause != nil {
        b.WriteString(fmt.Sprintf("\nCaused by: %v", e.Cause))
    }
    
    // 添加上下文信息
    if len(e.Context) > 0 {
        b.WriteString("\nContext:")
        for k, v := range e.Context {
            b.WriteString(fmt.Sprintf("\n  %s: %s", k, v))
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
func NewSystemError(code ErrorCode, message string, cause error) *SystemError {
    return &SystemError{
        Code:      code,
        Message:   message,
        Cause:     cause,
        Stack:     captureStack(),
        Time:      time.Now(),
        Context:   make(map[string]string),
        Severity:  getSeverityForError(code),
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
    if e.Context == nil {
        e.Context = make(map[string]string)
    }
    e.Context[key] = value
    return e
}

// WithSeverity 设置错误严重度
func (e *SystemError) WithSeverity(severity IssueSeverity) *SystemError {
    e.Severity = severity
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
            Severity: sysErr.Severity,
        }
    }
    
    // 创建新的 SystemError
    return NewSystemError(code, message, err)
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

// getSeverityForError 根据错误码确定严重度
func getSeverityForError(code ErrorCode) IssueSeverity {
    switch {
    case code >= 9000:
        return SeverityCritical
    case code >= 7000:
        return SeverityError
    case code >= 5000:
        return SeverityWarning
    default:
        return SeverityInfo
    }
}

// ErrorHandler 错误处理器接口
type ErrorHandler interface {
    Handle(error) error
}

// 错误处理相关常量
const (
    MaxStackDepth    = 32   // 最大堆栈深度
    MaxErrorHistory  = 100  // 最大错误历史记录数
    MaxRetryAttempts = 3    // 最大重试次数
)
