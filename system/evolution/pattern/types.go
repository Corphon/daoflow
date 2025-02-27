// system/evolution/pattern/types.go

package pattern

import (
	"time"

	"github.com/Corphon/daoflow/system/common"
	"github.com/Corphon/daoflow/system/meta/emergence"
)

// 确保实现了 SharedPattern 接口
var _ common.SharedPattern = (*RecognizedPattern)(nil)

// RecognizedPattern 识别的模式
type RecognizedPattern struct {
	common.BasePattern                            // 嵌入基础模式结构
	Pattern            *emergence.EmergentPattern // 原始模式
	Signature          PatternSignature           // 模式特征
	Evolution          []PatternState             // 演化历史
	Properties         map[string]float64         // 附加属性
	Context            map[string]float64         // 上下文环境因素

	ID       string             // 模式ID
	Type     string             // 模式类型
	Features map[string]float64 // 特征向量
	Created  time.Time          // 创建时间

	Active     bool      // 是否活跃
	Formation  time.Time // 形成时间
	LastUpdate time.Time // 最后更新时间

	Confidence  float64   // 置信度
	Stability   float64   // 稳定性
	FirstSeen   time.Time // 首次发现时间
	LastSeen    time.Time // 最后发现时间
	Occurrences int       // 出现次数
	Strength    float64
}

// PatternState 模式状态
type PatternState struct {
	Pattern    *emergence.EmergentPattern // 模式
	Active     bool                       // 是否活跃
	Duration   time.Duration              // 持续时间
	LastUpdate time.Time                  // 最后更新时间
	Properties map[string]float64         // 状态属性
}

// EvolutionState 演化状态结构
type EvolutionState struct {
	// 基础数据
	Matches      map[string]*EvolutionMatch // 当前匹配
	Trajectories map[string]*EvolutionPath  // 演化轨迹
	Context      *MatchingContext           // 匹配上下文

	// 添加模式集合
	Patterns map[string]*RecognizedPattern // 当前识别的模式

	// 统计指标
	Metrics struct {
		ActivityLevel float64 // 活动水平
		EnergyLevel   float64 // 能量水平
		Stability     float64 // 稳定性
		ChangeRate    float64 // 变化率
	}

	// 时间信息
	LastUpdate time.Time // 最后更新时间
}

// -------------------------------------------------------------------------
// convertToPatternState 转换EmergentPattern.PatternState到本地PatternState
func convertToPatternState(state *emergence.EmergentPattern) PatternState {
	return PatternState{
		Pattern:    state,
		Properties: state.Properties,
		Active:     true,
		Duration:   time.Since(state.Formation),
		LastUpdate: state.LastUpdate,
	}
}

// convertPatternState 转换emergence.PatternState到本地PatternState
func convertPatternState(state emergence.PatternState) PatternState {
	return PatternState{
		Pattern:    state.Pattern,
		Active:     state.Active,
		Duration:   state.Duration,
		LastUpdate: state.LastUpdate,
		Properties: state.Properties,
	}
}

// convertLocalPatternState 转换本地PatternState到emergence.PatternState
func convertLocalPatternState(state PatternState) emergence.PatternState {
	return emergence.PatternState{
		Pattern:    state.Pattern,
		Active:     state.Active,
		Duration:   state.Duration,
		LastUpdate: state.LastUpdate,
		Properties: state.Properties,
	}
}
