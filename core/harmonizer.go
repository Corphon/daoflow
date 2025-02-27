// core/harmonizer.go

package core

import (
	"math"
	"sync"
)

// HarmonyState 和谐状态
type HarmonyState struct {
	Value      float64            // 和谐值
	Components map[string]float64 // 分项和谐度
}

// Harmonizer 和谐器
type Harmonizer struct {
	mu sync.RWMutex

	// 状态属性
	harmony    float64            // 总体和谐度
	components map[string]float64 // 分项和谐度
	weights    map[string]float64 // 权重配置
}

// -----------------------------------------------
// NewHarmonizer 创建新的和谐器
func NewHarmonizer() *Harmonizer {
	return &Harmonizer{
		harmony:    1.0,
		components: make(map[string]float64),
		weights:    make(map[string]float64),
	}
}

// Initialize 初始化和谐器
func (h *Harmonizer) Initialize() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.harmony = 1.0
	h.components = make(map[string]float64)
	h.weights = make(map[string]float64)

	return nil
}

// SetWeight 设置分量权重
func (h *Harmonizer) SetWeight(component string, weight float64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.weights[component] = weight
}

// UpdateComponent 更新分量和谐度
func (h *Harmonizer) UpdateComponent(component string, value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 确保值在[0,1]范围内
	value = math.Max(0, math.Min(1, value))
	h.components[component] = value

	// 重新计算总体和谐度
	h.calculateHarmony()
}

// GetHarmony 获取总体和谐度
func (h *Harmonizer) GetHarmony() float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.harmony
}

// GetState 获取和谐状态
func (h *Harmonizer) GetState() HarmonyState {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// 复制组件状态
	components := make(map[string]float64)
	for k, v := range h.components {
		components[k] = v
	}

	return HarmonyState{
		Value:      h.harmony,
		Components: components,
	}
}

// calculateHarmony 计算总体和谐度
func (h *Harmonizer) calculateHarmony() {
	if len(h.components) == 0 {
		h.harmony = 1.0
		return
	}

	totalWeight := 0.0
	weightedSum := 0.0

	for component, value := range h.components {
		weight := h.weights[component]
		if weight == 0 {
			weight = 1.0 // 默认权重
		}
		totalWeight += weight
		weightedSum += value * weight
	}

	if totalWeight > 0 {
		h.harmony = weightedSum / totalWeight
	} else {
		h.harmony = 0
	}
}

// Close 关闭和谐器
func (h *Harmonizer) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.components = nil
	h.weights = nil
	return nil
}
