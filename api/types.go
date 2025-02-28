// api/types.go
package api

import (
	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/types"
)

// 类型别名引用
type (
	Pattern      = model.FlowPattern
	EventType    = types.EventType
	EventHandler = types.EventHandler
	SystemEvent  = types.SystemEvent
	SystemState  = types.SystemState
)
