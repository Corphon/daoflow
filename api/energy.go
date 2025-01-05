// api/energy.go

package api

import (
    "context"
    "sync"

    "github.com/Corphon/daoflow/system"
)

// EnergyAPI 能量管理API
type EnergyAPI struct {
    mu     sync.RWMutex
    system *system.SystemCore
    opts   *Options
}

// NewEnergyAPI 创建能量API实例
func NewEnergyAPI(sys *system.SystemCore, opts *Options) *EnergyAPI {
    return &EnergyAPI{
        system: sys,
        opts:   opts,
    }
}

// GetEnergy 获取当前能量值
func (e *EnergyAPI) GetEnergy() (float64, error) {
    e.mu.RLock()
    defer e.mu.RUnlock()
    
    return e.system.GetEnergy(), nil
}

// SetEnergy 设置系统能量值
func (e *EnergyAPI) SetEnergy(ctx context.Context, value float64) error {
    e.mu.Lock()
    defer e.mu.Unlock()

    return e.system.SetEnergy(value)
}

// GetCapacity 获取能量容量
func (e *EnergyAPI) GetCapacity() (float64, error) {
    e.mu.RLock()
    defer e.mu.RUnlock()

    return e.system.GetEnergyCapacity(), nil
}

// GetUsage 获取能量使用情况
func (e *EnergyAPI) GetUsage() (map[string]float64, error) {
    e.mu.RLock()
    defer e.mu.RUnlock()

    return e.system.GetEnergyUsage(), nil
}

// Subscribe 订阅能量变化事件
func (e *EnergyAPI) Subscribe() (<-chan float64, error) {
    return e.system.SubscribeEnergy(), nil
}
