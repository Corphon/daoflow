// core/errors.go

package core

import (
	"fmt"
	"runtime"
	"strings"
)

// CoreError 核心错误类型
type CoreError struct {
	Message string    // 错误消息
	Code    ErrorCode // 错误码
	Stack   []string  // 错误堆栈
	cause   error     // 原因错误
}

// ErrorCode 错误码类型
type ErrorCode string

const (
	// 基础错误
	ErrInvalid    ErrorCode = "INVALID"    // 无效值
	ErrRange      ErrorCode = "RANGE"      // 超出范围
	ErrState      ErrorCode = "STATE"      // 状态错误
	ErrInitialize ErrorCode = "INITIALIZE" // 初始化错误

	// 量子相关错误
	ErrQuantum   ErrorCode = "QUANTUM"   // 量子态错误
	ErrSuperpose ErrorCode = "SUPERPOSE" // 叠加态错误
	ErrEntangle  ErrorCode = "ENTANGLE"  // 纠缠态错误

	// 场相关错误
	ErrField       ErrorCode = "FIELD"       // 场错误
	ErrPotential   ErrorCode = "POTENTIAL"   // 势场错误
	ErrInteraction ErrorCode = "INTERACTION" // 相互作用错误

	// 能量相关错误
	ErrEnergy       ErrorCode = "ENERGY"       // 能量错误
	ErrTransform    ErrorCode = "TRANSFORM"    // 转换错误
	ErrConservation ErrorCode = "CONSERVATION" // 守恒错误
)

// NewCoreError 创建新的核心错误
func NewCoreError(message string) error {
	return &CoreError{
		Message: message,
		Code:    ErrInvalid,
		Stack:   captureStack(),
	}
}

// NewCoreErrorWithCode 创建带错误码的核心错误
func NewCoreErrorWithCode(code ErrorCode, message string) error {
	return &CoreError{
		Message: message,
		Code:    code,
		Stack:   captureStack(),
	}
}

// WrapCoreError 包装错误
func WrapCoreError(err error, code ErrorCode, message string) error {
	if err == nil {
		return nil
	}

	return &CoreError{
		Message: message,
		Code:    code,
		Stack:   captureStack(),
		cause:   err,
	}
}

// Error 实现 error 接口
func (e *CoreError) Error() string {
	var b strings.Builder

	// 构建错误消息
	b.WriteString(fmt.Sprintf("[%s] %s", e.Code, e.Message))

	// 添加原因错误
	if e.cause != nil {
		b.WriteString(fmt.Sprintf("\nCaused by: %v", e.cause))
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

// Unwrap 实现 errors.Unwrap 接口
func (e *CoreError) Unwrap() error {
	return e.cause
}

// GetCode 获取错误码
func (e *CoreError) GetCode() ErrorCode {
	return e.Code
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
