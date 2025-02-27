// core/correlator.go

package core

import (
	"math"
	"sync"
	"time"
)

// CorrelatorConfig 关联器配置
type CorrelatorConfig struct {
	DecayTime      float64 // 关联衰减时间
	MaxCorrelation float64 // 最大关联强度
	UpdateInterval float64 // 更新间隔
}

// Correlator 关联器
type Correlator struct {
	mu sync.RWMutex

	// 状态
	correlations map[string]float64   // 关联强度映射
	lastUpdated  map[string]time.Time // 最后更新时间

	// 配置
	config *CorrelatorConfig
}

// ---------------------------------------------------
// NewCorrelator 创建关联器
func NewCorrelator() *Correlator {
	return &Correlator{
		correlations: make(map[string]float64),
		lastUpdated:  make(map[string]time.Time),
		config: &CorrelatorConfig{
			DecayTime:      1.0,
			MaxCorrelation: 1.0,
			UpdateInterval: 0.1,
		},
	}
}

// Initialize 初始化关联器
func (c *Correlator) Initialize() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.correlations = make(map[string]float64)
	c.lastUpdated = make(map[string]time.Time)

	return nil
}

// SetCorrelation 设置关联强度
func (c *Correlator) SetCorrelation(key string, value float64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if value < 0 || value > c.config.MaxCorrelation {
		return NewCoreError("correlation value out of range")
	}

	c.correlations[key] = value
	c.lastUpdated[key] = time.Now()
	return nil
}

// GetCorrelation 获取关联强度
func (c *Correlator) GetCorrelation(key string) float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.correlations[key]
}

// Update 更新关联状态
func (c *Correlator) Update() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, lastTime := range c.lastUpdated {
		dt := now.Sub(lastTime).Seconds()
		if dt > c.config.UpdateInterval {
			// 应用时间衰减
			c.correlations[key] *= math.Exp(-dt / c.config.DecayTime)
			c.lastUpdated[key] = now
		}
	}
	return nil
}
