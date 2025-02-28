// system/types/evolution_params.go

package types

import (
	"time"

	"github.com/Corphon/daoflow/model"
)

// EvolutionParams 演化参数
type EvolutionParams struct {
	Mode        string             // 演化模式
	Rate        float64            // 演化速率
	Duration    time.Duration      // 演化周期
	Target      model.SystemState  // 目标状态
	Constraints map[string]float64 // 约束条件
	StepSize    float64            // 步长
	MaxSteps    int                // 最大步数
}

// EvolutionPoint 演化点
type EvolutionPoint struct {
	State     model.SystemState      // 系统状态
	Energy    float64                // 能量值
	Timestamp time.Time              // 时间戳
	Meta      map[string]interface{} // 元数据
}

// EvolutionMetrics 演化指标
type EvolutionMetrics struct {
	Progress   float64       // 进度
	Stability  float64       // 稳定性
	Efficiency float64       // 效率
	Success    bool          // 是否成功
	Duration   time.Duration // 持续时间
}

// EvolutionStatus 演化状态
type EvolutionStatus struct {
	Phase     model.Phase    // 当前阶段
	Direction model.Vector3D // 演化方向
	Progress  float64        // 完成进度
	Stability float64        // 系统稳定性
	Energy    float64        // 系统能量
	UpdatedAt time.Time      // 更新时间
}

// EvolutionPath 演化路径
type EvolutionPath struct {
	Points    []EvolutionPoint // 演化点列表
	Metrics   EvolutionMetrics // 路径指标
	Valid     bool             // 是否有效
	CreatedAt time.Time        // 创建时间
}
