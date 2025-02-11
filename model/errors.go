// model/errors.go

package model

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
)

// ErrorCode 错误码
type ErrorCode string

const (
	// 操作相关错误
	ErrCodeOperation  ErrorCode = "OPERATION"  // 操作错误
	ErrCodeState      ErrorCode = "STATE"      // 状态错误
	ErrCodeTransform  ErrorCode = "TRANSFORM"  // 转换错误
	ErrCodeSync       ErrorCode = "SYNC"       // 同步错误
	ErrCodeValidation ErrorCode = "VALIDATION" // 验证错误
	ErrCodeInit       ErrorCode = "INIT"       // 初始化错误

	// 模型相关错误
	ErrCodeModel   ErrorCode = "MODEL"   // 模型错误
	ErrCodeYinYang ErrorCode = "YINYANG" // 阴阳模型错误
	ErrCodeWuXing  ErrorCode = "WUXING"  // 五行模型错误
	ErrCodeBaGua   ErrorCode = "BAGUA"   // 八卦模型错误
	ErrCodeGanZhi  ErrorCode = "GANZHI"  // 干支模型错误

	// 资源相关错误
	ErrCodeResource  ErrorCode = "RESOURCE"  // 资源错误
	ErrCodeEnergy    ErrorCode = "ENERGY"    // 能量错误
	ErrCodeField     ErrorCode = "FIELD"     // 场错误
	ErrCodeQuantum   ErrorCode = "QUANTUM"   // 量子态错误
	ErrCodeNotFound  ErrorCode = "NOTFOUND"  // 未找到错误
	ErrCodeDuplicate ErrorCode = "DUPLICATE" // 重复错误
	ErrCodeLimit     ErrorCode = "LIMIT"     // 限制错误

	//new
	ErrCodeInvalid  ErrorCode = "invalid"  // 无效参数
	ErrCodeRange    ErrorCode = "range"    // 超出范围
	ErrCodeNone     ErrorCode = ""         // 无错误
	ErrCodeInternal ErrorCode = "internal" // 内部错误
	ErrCodeIO       ErrorCode = "io"       // IO错误

	// 严重级别错误码
	ErrCodeCritical ErrorCode = "CRITICAL" // 严重错误
	ErrCodeError    ErrorCode = "ERROR"    // 一般错误
	ErrCodeWarning  ErrorCode = "WARNING"  // 警告
	ErrCodeInfo     ErrorCode = "INFO"     // 信息

	// 依赖相关错误
	ErrCodeDependency ErrorCode = "DEPENDENCY" // 依赖错误

	// 时间相关错误
	ErrCodeTimeout  ErrorCode = "TIMEOUT"  // 超时错误
	ErrCodeDeadline ErrorCode = "DEADLINE" // 截止时间错误
	ErrCodeInterval ErrorCode = "INTERVAL" // 间隔错误

	// 共识相关错误
	ErrCodeConsensus ErrorCode = "CONSENSUS" // 共识错误
	ErrCodeQuorum    ErrorCode = "QUORUM"    // 法定人数错误
	ErrCodeVote      ErrorCode = "VOTE"      // 投票错误
	ErrCodeAgreement ErrorCode = "AGREEMENT" // 协议错误
)

// ModelError 模型错误
type ModelError struct {
	Code    ErrorCode // 错误码
	Message string    // 错误消息
	Cause   error     // 原因错误
	Stack   []string  // 错误堆栈
}

// Error 实现 error 接口
func (e *ModelError) Error() string {
	var b strings.Builder

	// 构建错误消息
	b.WriteString(fmt.Sprintf("[%s] %s", e.Code, e.Message))

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

// NewModelError 创建新的模型错误
func NewModelError(code ErrorCode, message string, cause error) *ModelError {
	return &ModelError{
		Code:    code,
		Message: message,
		Cause:   cause,
		Stack:   captureStack(),
	}
}

// WrapError 包装错误
func WrapError(err error, code ErrorCode, message string) *ModelError {
	if err == nil {
		return nil
	}

	// 如果已经是 ModelError，则添加到堆栈
	if modelErr, ok := err.(*ModelError); ok {
		return &ModelError{
			Code:    code,
			Message: message,
			Cause:   modelErr,
			Stack:   append(captureStack(), modelErr.Stack...),
		}
	}

	// 创建新的 ModelError
	return NewModelError(code, message, err)
}

// IsModelError 检查是否为模型错误
func IsModelError(err error) bool {
	_, ok := err.(*ModelError)
	return ok
}

// GetErrorCode 获取错误码
func GetErrorCode(err error) ErrorCode {
	if modelErr, ok := err.(*ModelError); ok {
		return modelErr.Code
	}
	return ""
}

// GetRootCause 获取根本原因
func GetRootCause(err error) error {
	if modelErr, ok := err.(*ModelError); ok {
		if modelErr.Cause == nil {
			return modelErr
		}
		return GetRootCause(modelErr.Cause)
	}
	return err
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

// DefaultErrorHandler 默认错误处理器
type DefaultErrorHandler struct {
	// 错误处理配置
	config struct {
		MaxRetries int
		LogErrors  bool
	}
}

// NewDefaultErrorHandler 创建默认错误处理器
func NewDefaultErrorHandler() *DefaultErrorHandler {
	handler := &DefaultErrorHandler{}
	handler.config.MaxRetries = 3
	handler.config.LogErrors = true
	return handler
}

// Handle 处理错误
func (h *DefaultErrorHandler) Handle(err error) error {
	if err == nil {
		return nil
	}

	// 记录错误
	if h.config.LogErrors {
		logError(err)
	}

	// 获取错误码
	code := GetErrorCode(err)

	// 根据错误码处理
	switch code {
	case ErrCodeOperation:
		return h.handleOperationError(err)
	case ErrCodeState:
		return h.handleStateError(err)
	case ErrCodeTransform:
		return h.handleTransformError(err)
	case ErrCodeSync:
		return h.handleSyncError(err)
	default:
		return err
	}
}

// handleOperationError 处理操作错误
func (h *DefaultErrorHandler) handleOperationError(err error) error {
	// 实现操作错误处理逻辑
	return err
}

// handleStateError 处理状态错误
func (h *DefaultErrorHandler) handleStateError(err error) error {
	// 实现状态错误处理逻辑
	return err
}

// handleTransformError 处理转换错误
func (h *DefaultErrorHandler) handleTransformError(err error) error {
	// 实现转换错误处理逻辑
	return err
}

// handleSyncError 处理同步错误
func (h *DefaultErrorHandler) handleSyncError(err error) error {
	// 实现同步错误处理逻辑
	return err
}

// logError 记录错误
func logError(err error) {
	// 实现错误日志记录
	fmt.Printf("Error occurred: %v\n", err)
}

// 全局错误处理器
var (
	globalHandler ErrorHandler
	handlerOnce   sync.Once
)

// GetErrorHandler 获取全局错误处理器
func GetErrorHandler() ErrorHandler {
	handlerOnce.Do(func() {
		if globalHandler == nil {
			globalHandler = NewDefaultErrorHandler()
		}
	})
	return globalHandler
}

// SetErrorHandler 设置全局错误处理器
func SetErrorHandler(handler ErrorHandler) {
	if handler != nil {
		globalHandler = handler
	}
}
