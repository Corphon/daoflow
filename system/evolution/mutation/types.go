// system/evolution/mutation/types.go

package mutation

import (
	"time"

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

// Mutation 实现 Mutation 接口
type Mutation struct {
	ID          string           // 突变ID
	Type        string           // 突变类型
	Source      *MutationSource  // 突变源
	Changes     []MutationChange // 变化列表
	Severity    float64          // 严重程度
	Probability float64          // 发生概率
	DetectedAt  time.Time        // 检测时间
	LastUpdate  time.Time        // 最后更新时间
	Status      string           // 当前状态
}

// 实现 SharedPattern 接口方法
func (m *Mutation) GetID() string {
	return m.ID
}

func (m *Mutation) GetType() string {
	return m.Type
}

func (m *Mutation) GetStrength() float64 {
	return m.Severity // 使用严重度作为强度
}

func (m *Mutation) GetStability() float64 {
	return 1 - m.Probability // 概率越高稳定性越低
}

func (m *Mutation) GetTimestamp() time.Time {
	return m.DetectedAt
}

// 实现 Mutation 接口方法
func (m *MutationType) GetSource() string                   { return m.Source }
func (m *MutationType) GetTarget() string                   { return m.Target }
func (m *MutationType) GetProbability() float64             { return m.Probability }
func (m *MutationType) GetChanges() []common.MutationChange { return m.Changes }

// 确保实现了 SharedPattern 接口
var _ common.SharedPattern = (*Mutation)(nil)
