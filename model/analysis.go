//model/analysis.go

package model

import "time"

// AnalysisResult 分析结果基础接口
type AnalysisResult interface {
	GetID() string
	GetTimestamp() time.Time
	GetType() string
}

// MutationAnalysis 突变分析结果
type MutationAnalysis struct {
	ID           string         // 分析ID
	MutationID   string         // 突变ID
	Causes       []CausalFactor // 因果因素
	Effects      []Effect       // 影响效果
	Correlations []Correlation  // 相关性
	Risk         RiskAssessment // 风险评估
	Created      time.Time      // 创建时间
}

// CausalFactor 因果因素
type CausalFactor struct {
	Type       string   // 因素类型
	Source     string   // 来源
	Weight     float64  // 权重
	Confidence float64  // 置信度
	Evidence   []string // 证据
}

// Effect 影响效果
type Effect struct {
	Target     string        // 影响目标
	Type       string        // 效果类型
	Magnitude  float64       // 影响程度
	Duration   time.Duration // 持续时间
	Reversible bool          // 是否可逆
}

// Correlation 相关性
type Correlation struct {
	SourceID   string        // 源ID
	TargetID   string        // 目标ID
	Type       string        // 相关类型
	Strength   float64       // 相关强度
	Direction  int           // 相关方向
	TimeOffset time.Duration // 时间偏移
}

// RiskAssessment 风险评估
type RiskAssessment struct {
	Level      string       // 风险等级
	Score      float64      // 风险分数
	Factors    []RiskFactor // 风险因素
	Mitigation []string     // 缓解措施
}

// RiskFactor 风险因素
type RiskFactor struct {
	Type        string  // 因素类型
	Impact      float64 // 影响程度
	Probability float64 // 发生概率
	Urgency     int     // 紧急程度
}

//-------------------------------------------------------
// 实现 AnalysisResult 接口
func (ma *MutationAnalysis) GetID() string {
	return ma.ID
}

func (ma *MutationAnalysis) GetTimestamp() time.Time {
	return ma.Created
}

func (ma *MutationAnalysis) GetType() string {
	return "mutation"
}
