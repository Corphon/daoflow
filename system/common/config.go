// system/common/config.go
package common

import "time"

// Config 公共系统配置
type Config struct {
	// 基础配置
	Base struct {
		UpdateInterval time.Duration `json:"update_interval"` // 更新间隔
		MaxRetries     int           `json:"max_retries"`     // 最大重试次数
		Timeout        time.Duration `json:"timeout"`         // 超时时间
	} `json:"base"`

	// 资源限制
	Resources struct {
		MaxFields    int     `json:"max_fields"`    // 最大场数量
		MaxStates    int     `json:"max_states"`    // 最大状态数
		MaxPatterns  int     `json:"max_patterns"`  // 最大模式数
		MaxEnergy    float64 `json:"max_energy"`    // 最大能量值
		ReserveRatio float64 `json:"reserve_ratio"` // 预留比例
	} `json:"resources"`

	// 共享配置
	Sharing struct {
		EnableCache  bool          `json:"enable_cache"`  // 启用缓存
		CacheTTL     time.Duration `json:"cache_ttl"`     // 缓存过期时间
		SyncInterval time.Duration `json:"sync_interval"` // 同步间隔
	} `json:"sharing"`

	// 监控配置
	Monitor struct {
		EnableMetrics bool    `json:"enable_metrics"` // 启用指标
		SampleRate    float64 `json:"sample_rate"`    // 采样率
		ErrorRate     float64 `json:"error_rate"`     // 错误率阈值
	} `json:"monitor"`
}

// DefaultConfig 返回默认配置
func NewDefaultConfig() *Config {
	return &Config{
		Base: struct {
			UpdateInterval time.Duration `json:"update_interval"`
			MaxRetries     int           `json:"max_retries"`
			Timeout        time.Duration `json:"timeout"`
		}{
			UpdateInterval: time.Second,
			MaxRetries:     3,
			Timeout:        time.Minute,
		},

		Resources: struct {
			MaxFields    int     `json:"max_fields"`
			MaxStates    int     `json:"max_states"`
			MaxPatterns  int     `json:"max_patterns"`
			MaxEnergy    float64 `json:"max_energy"`
			ReserveRatio float64 `json:"reserve_ratio"`
		}{
			MaxFields:    100,
			MaxStates:    100,
			MaxPatterns:  1000,
			MaxEnergy:    1000.0,
			ReserveRatio: 0.2,
		},
	}
}
