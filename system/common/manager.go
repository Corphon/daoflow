// system/common/manager.go

package common

import (
	"context"
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
	"github.com/Corphon/daoflow/system/types"
)

// Manager 通用系统管理器
type Manager struct {
	mu sync.RWMutex

	// 基础配置
	config *types.CommonConfig

	// 共享资源
	resources struct {
		fields   map[string]*core.Field        // 共享场
		states   map[string]*core.QuantumState // 共享量子态
		patterns map[string]*types.Pattern     // 共享模式
		energies map[string]float64            // 共享能量
	}

	// 状态信息
	status struct {
		isRunning  bool
		startTime  time.Time
		lastUpdate time.Time
		errors     []error
	}

	// 上下文控制
	ctx    context.Context
	cancel context.CancelFunc
}

//----------------------------------------------------------
// NewManager 创建新的管理器实例
func NewManager(cfg *types.CommonConfig) (*Manager, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	m := &Manager{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
	}

	// 初始化共享资源
	m.resources.fields = make(map[string]*core.Field)
	m.resources.states = make(map[string]*core.QuantumState)
	m.resources.patterns = make(map[string]*types.Pattern)
	m.resources.energies = make(map[string]float64)

	return m, nil
}

// DefaultConfig 返回默认配置
func DefaultConfig() *types.CommonConfig {
	return &types.CommonConfig{
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

// Start 启动管理器
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.status.isRunning {
		return nil
	}

	m.status.isRunning = true
	m.status.startTime = time.Now()
	return nil
}

// Stop 停止管理器
func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.status.isRunning {
		return nil
	}

	m.cancel()
	m.status.isRunning = false
	return nil
}

// Status 获取管理器状态
func (m *Manager) Status() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.status.isRunning {
		return "running"
	}
	return "stopped"
}

// Wait 等待管理器停止
func (m *Manager) Wait() {
	<-m.ctx.Done()
}

// GetMetrics 获取管理器指标
func (m *Manager) GetMetrics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"fields_count":   len(m.resources.fields),
		"states_count":   len(m.resources.states),
		"patterns_count": len(m.resources.patterns),
		"energy_total":   m.getTotalEnergy(),
		"uptime":         time.Since(m.status.startTime).String(),
		"error_count":    len(m.status.errors),
	}
}

// 辅助方法

// getTotalEnergy 获取总能量
func (m *Manager) getTotalEnergy() float64 {
	total := 0.0
	for _, energy := range m.resources.energies {
		total += energy
	}
	return total
}

// Restore 恢复系统
func (m *Manager) Restore(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 重置状态
	m.resources.fields = make(map[string]*core.Field)
	m.resources.states = make(map[string]*core.QuantumState)
	m.resources.patterns = make(map[string]*types.Pattern)
	m.resources.energies = make(map[string]float64)

	m.status.lastUpdate = time.Now()
	m.status.errors = make([]error, 0)

	return nil
}
