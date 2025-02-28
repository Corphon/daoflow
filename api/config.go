// api/config.go
package api

import (
	"time"

	"github.com/Corphon/daoflow/system"
)

// Config 使用system的配置
type Config struct {
	// 直接使用system配置
	SystemConfig *system.Config // 系统配置
	LogLevel     string         // 日志级别
	MaxRetries   int            // 最大重试次数
	Timeout      time.Duration  // 超时时间
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		SystemConfig: system.DefaultConfig(),
		LogLevel:     "info",
		MaxRetries:   3,
		Timeout:      time.Minute,
	}
}
