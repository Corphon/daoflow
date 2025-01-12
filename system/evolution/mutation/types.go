// system/evolution/mutation/types.go

package mutation

import (
    "github.com/Corphon/daoflow/system/common"
)

// Mutation 突变类型
type Mutation struct {
    common.BasePattern             // 嵌入基础模式结构
    Source      *MutationSource   // 突变源
    Changes     []MutationChange  // 变化列表
    Probability float64           // 发生概率
    Status      string           // 当前状态
}

// 确保实现了 SharedPattern 接口
var _ common.SharedPattern = (*Mutation)(nil)
