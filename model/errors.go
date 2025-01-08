// model/errors.go

package model

import "errors"

var (
    // 基础错误
    ErrModelNotInitialized = errors.New("model not initialized")
    ErrModelAlreadyStarted = errors.New("model already started")
    ErrModelNotStarted     = errors.New("model not started")
    ErrInvalidModelType    = errors.New("invalid model type")

    // 状态错误
    ErrInvalidState       = errors.New("invalid state")
    ErrStateTransition    = errors.New("invalid state transition")
    ErrStateNotReady      = errors.New("state not ready for operation")
    ErrStateLocked        = errors.New("state is locked")

    // 能量错误
    ErrEnergyOutOfRange   = errors.New("energy value out of valid range")
    ErrEnergyOverflow     = errors.New("energy overflow")
    ErrEnergyUnderflow    = errors.New("energy underflow")
    ErrEnergyImbalance    = errors.New("system energy imbalance")

    // 相互作用错误
    ErrInvalidInteraction = errors.New("invalid model interaction")
    ErrInteractionFailed  = errors.New("model interaction failed")
    ErrInvalidTransform   = errors.New("invalid transformation pattern")
    
    // 量子态错误
    ErrQuantumDecoherence = errors.New("quantum state decoherence")
    ErrQuantumCollapse    = errors.New("quantum state collapse")
    ErrQuantumEntangle    = errors.New("quantum entanglement failed")

    // 场论错误
    ErrFieldOverload      = errors.New("field strength overload")
    ErrFieldInterference  = errors.New("destructive field interference")
    ErrFieldResonance     = errors.New("unstable field resonance")

    // 系统集成错误
    ErrSystemUnbalanced   = errors.New("system is unbalanced")
    ErrSystemOverload     = errors.New("system overload")
    ErrSystemInstability  = errors.New("system instability detected")
    ErrSystemAsynchrony   = errors.New("system components out of sync")

    // 转换错误
    ErrInvalidPhaseTransition = errors.New("invalid phase transition")
    ErrInvalidNatureChange    = errors.New("invalid nature change")
    ErrInvalidElementChange   = errors.New("invalid element change")

    // 观察者错误
    ErrObserverNotFound      = errors.New("observer not found")
    ErrObserverAlreadyExists = errors.New("observer already exists")
    ErrObserverFailed        = errors.New("observer notification failed")

    // 配置错误
    ErrInvalidConfiguration  = errors.New("invalid model configuration")
    ErrConfigurationMismatch = errors.New("configuration mismatch")
    ErrInvalidParameter      = errors.New("invalid parameter value")

    // 资源错误
    ErrResourceExhausted    = errors.New("system resources exhausted")
    ErrResourceUnavailable  = errors.New("required resource unavailable")
    ErrResourceLocked       = errors.New("resource is locked")

    // 周期错误
    ErrCycleInterrupted     = errors.New("cycle interrupted")
    ErrCycleOutOfSync       = errors.New("cycle out of synchronization")
    ErrCycleOverflow        = errors.New("cycle counter overflow")

    // 核心集成错误
    ErrCoreIntegration      = errors.New("core integration failed")
    ErrCoreStateInvalid     = errors.New("core state invalid")
    ErrCoreOperationFailed  = errors.New("core operation failed")
)

// ErrorCode 错误代码类型
type ErrorCode int

const (
    // 基础错误码
    ErrCodeNone ErrorCode = iota
    ErrCodeInitialization
    ErrCodeOperation
    ErrCodeState
    ErrCodeEnergy
    ErrCodeInteraction
    ErrCodeQuantum
    ErrCodeField
    ErrCodeSystem
    ErrCodeTransformation
    ErrCodeObserver
    ErrCodeConfiguration
    ErrCodeResource
    ErrCodeCycle
    ErrCodeCore
)

// ModelError 模型错误类型
type ModelError struct {
    Code    ErrorCode
    Message string
    Err     error
}

// Error 实现error接口
func (e *ModelError) Error() string {
    if e.Err != nil {
        return e.Message + ": " + e.Err.Error()
    }
    return e.Message
}

// Unwrap 获取底层错误
func (e *ModelError) Unwrap() error {
    return e.Err
}

// NewModelError 创建新的模型错误
func NewModelError(code ErrorCode, message string, err error) *ModelError {
    return &ModelError{
        Code:    code,
        Message: message,
        Err:     err,
    }
}

// IsModelError 检查错误类型
func IsModelError(err error) bool {
    var modelErr *ModelError
    return errors.As(err, &modelErr)
}

// GetErrorCode 获取错误代码
func GetErrorCode(err error) ErrorCode {
    var modelErr *ModelError
    if errors.As(err, &modelErr) {
        return modelErr.Code
    }
    return ErrCodeNone
}

// WrapError 包装错误
func WrapError(err error, message string) error {
    if err == nil {
        return nil
    }
    var code ErrorCode
    if modelErr, ok := err.(*ModelError); ok {
        code = modelErr.Code
    }
    return NewModelError(code, message, err)
}
