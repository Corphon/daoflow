// api/evolution.go

package api

import (
    "context"
    "sync"

    "github.com/Corphon/daoflow/system"
)

// EvolutionMode 演化模式类型
type EvolutionMode string

const (
    ModeOptimize  EvolutionMode = "optimize"  // 优化模式
    ModeAdaptive  EvolutionMode = "adaptive"  // 适应模式
    ModeExplore   EvolutionMode = "explore"   // 探索模式
    ModeStabilize EvolutionMode = "stabilize" // 稳定模式
)

// EvolutionAPI 演化控制API
type EvolutionAPI struct {
    mu     sync.RWMutex
    system *system.SystemCore
    opts   *Options
}

// NewEvolutionAPI 创建演化API实例
func NewEvolutionAPI(sys *system.SystemCore, opts *Options) *EvolutionAPI {
    return &EvolutionAPI{
        system: sys,
        opts:   opts,
    }
}

// Evolve 触发系统演化
func (e *EvolutionAPI) Evolve(ctx context.Context, mode EvolutionMode) error {
    e.mu.Lock()
    defer e.mu.Unlock()

    return e.system.Evolution().Evolve(mode)
}

// GetMode 获取当前演化模式
func (e *EvolutionAPI) GetMode() (EvolutionMode, error) {
    e.mu.RLock()
    defer e.mu.RUnlock()

    return e.system.Evolution().GetMode(), nil
}

// GetStatus 获取演化状态
func (e *EvolutionAPI) GetStatus() (map[string]interface{}, error) {
    e.mu.RLock()
    defer e.mu.RUnlock()

    return e.system.Evolution().GetStatus(), nil
}

// GetMetrics 获取演化指标
func (e *EvolutionAPI) GetMetrics() (map[string]float64, error) {
    e.mu.RLock()
    defer e.mu.RUnlock()

    return e.system.Evolution().GetMetrics(), nil
}

// Subscribe 订阅演化事件
func (e *EvolutionAPI) Subscribe() (<-chan map[string]interface{}, error) {
    return e.system.Evolution().Subscribe(), nil
}
