// system/evolution/pattern/types.go

package pattern

import (
    "github.com/Corphon/daoflow/system/common"
)

// RecognizedPattern 识别的模式
type RecognizedPattern struct {
    common.BasePattern              // 嵌入基础模式结构
    Signature    PatternSignature   // 模式特征
    Evolution    []PatternState     // 演化历史
    Properties   map[string]float64 // 附加属性
}

// 确保实现了 SharedPattern 接口
var _ common.SharedPattern = (*RecognizedPattern)(nil)
