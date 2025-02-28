// api/errors.go
package api

import (
	"github.com/Corphon/daoflow/system/types"
)

var (
	// 系统状态错误
	ErrAlreadyRunning = types.ErrAlreadyRunning
	ErrNotRunning     = types.ErrNotRunning
	ErrInitialized    = types.ErrInitialized
	ErrNotInitialized = types.ErrNotInitialized

	// 模型相关错误
	ErrModelNotFound      = types.ErrModelNotFound
	ErrModelAlreadyExists = types.ErrModelAlreadyExists
	ErrModelInitFailed    = types.ErrModelInitFailed
	ErrModelStartFailed   = types.ErrModelStartFailed
	ErrModelStopFailed    = types.ErrModelStopFailed
)

// 使用type别名引用错误码类型
type ErrorCode = types.ErrorCode
