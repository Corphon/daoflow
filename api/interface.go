// api/interface.go

package api

import (
	"context"

	"github.com/Corphon/daoflow/core"
	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/types"
)

// DaoFlow 主接口
type DaoFlow interface {
	// 基础接口
	types.SystemInterface

	// 核心系统访问
	GetEnergySystem() *core.EnergySystem
	GetFieldSystem() *core.FieldSystem
	GetQuantumSystem() *core.QuantumSystem

	// 模型访问
	GetYinYangFlow() *model.YinYangFlow
	GetWuXingFlow() *model.WuXingFlow
	GetBaGuaFlow() *model.BaGuaFlow
	GetGanZhiFlow() *model.GanZhiFlow

	TransformModel(ctx context.Context, pattern model.TransformPattern) error
}

// AI能力接口
/*type AIAPI interface {
	// 量子预测 - 基于core的量子计算
	PredictQuantumState(state *core.QuantumState) (*Prediction, error)

	// 模式识别 - 基于model的阴阳五行
	RecognizePattern(data *PatternData) (*Pattern, error)
	DetectEmergence(timeWindow time.Duration) ([]*EmergentPattern, error)

	// 优化能力 - 基于system的演化
	OptimizeEnergy(target *EnergyTarget) (*Solution, error)
	BalanceElements(elements []Element) (*BalanceResult, error)

	// 学习能力
	LearnFromHistory(history []StateTransition) error
	AdaptStrategy(feedback *Feedback) error
}
*/
