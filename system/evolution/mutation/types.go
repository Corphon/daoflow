// system/evolution/mutation/types.go

package mutation

import (
    "github.com/Corphon/daoflow/system/common"
)

// MutationType 实现 Mutation 接口
type MutationType struct {
    common.BasePattern
    Source      string
    Target      string
    Probability float64
    Changes     []common.MutationChange
}

// 实现 Mutation 接口方法
func (m *MutationType) GetSource() string { return m.Source }
func (m *MutationType) GetTarget() string { return m.Target }
func (m *MutationType) GetProbability() float64 { return m.Probability }
func (m *MutationType) GetChanges() []common.MutationChange { return m.Changes }

// 确保实现了 SharedPattern 接口
var _ common.SharedPattern = (*Mutation)(nil)
