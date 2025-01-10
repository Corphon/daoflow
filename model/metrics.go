// model/metrics.go

package model

// ModelMetrics 模型指标
type ModelMetrics struct {
	YinYang struct {
		Balance   float64 `json:"balance"`   // 阴阳平衡度
		Harmony   float64 `json:"harmony"`   // 阴阳和谐度
		Transform float64 `json:"transform"` // 转化率
	} `json:"yin_yang"`

	WuXing struct {
		Cycles  float64 `json:"cycles"`  // 五行循环强度
		Energy  float64 `json:"energy"`  // 五行能量水平
		Balance float64 `json:"balance"` // 五行平衡度
	} `json:"wu_xing"`

	BaGua struct {
		Patterns  float64 `json:"patterns"`  // 八卦模式强度
		Changes   float64 `json:"changes"`   // 变化频率
		Stability float64 `json:"stability"` // 稳定性
	} `json:"ba_gua"`

	GanZhi struct {
		Alignment float64 `json:"alignment"` // 天干地支对齐度
		Cycle     float64 `json:"cycle"`     // 周期完整度
		Strength  float64 `json:"strength"`  // 作用强度
	} `json:"gan_zhi"`

	Integration float64 `json:"integration"` // 整体集成度
	Coherence   float64 `json:"coherence"`   // 整体相干性
	Emergence   float64 `json:"emergence"`   // 涌现程度
}
