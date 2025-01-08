// model/errors.go

package model

import (
    "errors"
    "fmt"
    "runtime"
    "strings"
)

// 保留原有的预定义错误常量
var (
    // 基础错误
    ErrModelNotInitialized = errors.New("model not initialized")
    // ... (保留原有的所有错误常量)
)

// ErrorCode 错误代码类型
type ErrorCode int

const (
    // 保留原有的错误码定义
    ErrCodeNone ErrorCode = iota
    // ... (保留原有的所有错误码)
)

// ModelError 模型错误类型
type ModelError struct {
    Code    ErrorCode   // 错误代码
    Message string      // 错误消息
    Cause   error      // 原因错误（改名以避免与标准库冲突）
    Stack   string     // 新增：堆栈信息
}

// NewModelError 创建新的模型错误
func NewModelError(code ErrorCode, message string, cause error) *ModelError {
    var stack strings.Builder
    
    // 获取堆栈信息
    for i := 1; i < 5; i++ {
        pc, file, line, ok := runtime.Caller(i)
        if !ok {
            break
        }
        fn := runtime.FuncForPC(pc)
        if fn == nil {
            continue
        }
        parts := strings.Split(file, "/")
        if len(parts) > 2 {
            file = strings.Join(parts[len(parts)-2:], "/")
        }
        stack.WriteString(fmt.Sprintf("%s:%d %s\n", file, line, fn.Name()))
    }

    return &ModelError{
        Code:    code,
        Message: message,
        Cause:   cause,
        Stack:   stack.String(),
    }
}

// Error 实现 error 接口
func (e *ModelError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
    }
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 实现 errors.Unwrap 接口
func (e *ModelError) Unwrap() error {
    return e.Cause
}

// String 实现 Stringer 接口
func (e ErrorCode) String() string {
    // ... (保留原有的实现)
}

// 保留并增强原有的工具函数
func IsModelError(err error) bool {
    var modelErr *ModelError
    return errors.As(err, &modelErr)
}

func GetErrorCode(err error) ErrorCode {
    var modelErr *ModelError
    if errors.As(err, &modelErr) {
        return modelErr.Code
    }
    return ErrCodeNone
}

// 新增实用工具函数
func ErrorStack(err error) string {
    var modelErr *ModelError
    if errors.As(err, &modelErr) {
        return modelErr.Stack
    }
    return ""
}

func FormatError(err error) string {
    if err == nil {
        return ""
    }

    var modelErr *ModelError
    if errors.As(err, &modelErr) {
        return fmt.Sprintf("Error: %s\nStack:\n%s", modelErr.Error(), modelErr.Stack)
    }
    return err.Error()
}

// 增强的错误包装函数
func WrapError(err error, code ErrorCode, message string) error {
    if err == nil {
        return nil
    }
    return NewModelError(code, message, err)
}
