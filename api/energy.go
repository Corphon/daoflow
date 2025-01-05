// api/energy.go

package api

import (
    "context"
    "sync"
    "time"

    "github.com/Corphon/daoflow/system"
)

// EnergyState 能量状态
type EnergyState string

const (
    StateOptimal    EnergyState = "optimal"    // 最优状态
    StateOverflow   EnergyState = "overflow"   // 能量溢出
    StateUnderflow  EnergyState = "underflow"  // 能量不足
    StateUnstable   EnergyState = "unstable"   // 不稳定
    StateRecharging EnergyState = "recharging" // 充能中
)

// EnergyDistribution 能量分配
type EnergyDistribution struct {
    Pattern    float64 `json:"pattern"`     // 模式识别能量
    Evolution  float64 `json:"evolution"`   // 演化过程能量
    Adaptation float64 `json:"adaptation"`  // 适应性调整能量
    Reserve    float64 `json:"reserve"`     // 能量储备
}

// EnergyConfig 能量配置
type EnergyConfig struct {
    MaxCapacity     float64           `json:"max_capacity"`      // 最大容量
    MinThreshold    float64           `json:"min_threshold"`     // 最小阈值
    Distribution    EnergyDistribution `json:"distribution"`     // 能量分配
    ChargeRate      float64           `json:"charge_rate"`       // 充能速率
    DischargeRate   float64           `json:"discharge_rate"`    // 放电速率
    BalanceInterval time.Duration     `json:"balance_interval"`  // 平衡间隔
}

// EnergyMetrics 能量指标
type EnergyMetrics struct {
    CurrentLevel    float64           `json:"current_level"`    // 当前能量水平
    Efficiency      float64           `json:"efficiency"`       // 能量效率
    Stability      float64           `json:"stability"`        // 能量稳定性
    Distribution   EnergyDistribution `json:"distribution"`    // 实际分配
    State          EnergyState       `json:"state"`           // 能量状态
    LastBalanced   time.Time         `json:"last_balanced"`   // 上次平衡时间
}

// EnergyEvent 能量事件
type EnergyEvent struct {
    Type      string      `json:"type"`       // 事件类型
    Level     float64     `json:"level"`      // 能量水平
    State     EnergyState `json:"state"`      // 能量状态
    Timestamp time.Time   `json:"timestamp"`  // 事件时间
    Details   string      `json:"details"`    // 详细信息
}

// EnergyAPI 能量管理API
type EnergyAPI struct {
    mu     sync.RWMutex
    system *system.SystemCore
    config EnergyConfig
}

// NewEnergyAPI 创建能量API实例
func NewEnergyAPI(sys *system.SystemCore, config EnergyConfig) *EnergyAPI {
    return &EnergyAPI{
        system: sys,
        config: config,
    }
}

// Configure 配置能量系统
func (e *EnergyAPI) Configure(ctx context.Context, config EnergyConfig) error {
    e.mu.Lock()
    defer e.mu.Unlock()

    if err := e.validateConfig(config); err != nil {
        return err
    }

    e.config = config
    return e.system.Energy().Configure(config)
}

// Distribute 调整能量分配
func (e *EnergyAPI) Distribute(ctx context.Context, dist EnergyDistribution) error {
    e.mu.Lock()
    defer e.mu.Unlock()

    if err := e.validateDistribution(dist); err != nil {
        return err
    }

    return e.system.Energy().Distribute(dist)
}

// GetMetrics 获取能量指标
func (e *EnergyAPI) GetMetrics(ctx context.Context) (*EnergyMetrics, error) {
    e.mu.RLock()
    defer e.mu.RUnlock()

    return e.system.Energy().GetMetrics()
}

// Balance 执行能量平衡
func (e *EnergyAPI) Balance(ctx context.Context) error {
    e.mu.Lock()
    defer e.mu.Unlock()

    return e.system.Energy().Balance()
}

// Subscribe 订阅能量事件
func (e *EnergyAPI) Subscribe(ctx context.Context) (<-chan EnergyEvent, error) {
    events := make(chan EnergyEvent, 100)
    
    sysEvents := e.system.Energy().Subscribe()
    
    go func() {
        defer close(events)
        for {
            select {
            case <-ctx.Done():
                return
            case evt := <-sysEvents:
                events <- e.convertEvent(evt)
            }
        }
    }()

    return events, nil
}

// validateConfig 验证能量配置
func (e *EnergyAPI) validateConfig(config EnergyConfig) error {
    if config.MaxCapacity <= 0 {
        return NewError(ErrInvalidConfig, "max capacity must be positive")
    }
    if config.MinThreshold < 0 || config.MinThreshold >= config.MaxCapacity {
        return NewError(ErrInvalidConfig, "invalid min threshold")
    }
    if config.ChargeRate <= 0 || config.DischargeRate <= 0 {
        return NewError(ErrInvalidConfig, "invalid rate settings")
    }
    return e.validateDistribution(config.Distribution)
}

// validateDistribution 验证能量分配
func (e *EnergyAPI) validateDistribution(dist EnergyDistribution) error {
    total := dist.Pattern + dist.Evolution + dist.Adaptation + dist.Reserve
    if total != 1.0 {
        return NewError(ErrInvalidDistribution, "distribution must sum to 1.0")
    }
    return nil
}

// convertEvent 转换系统事件到API事件
func (e *EnergyAPI) convertEvent(sysEvent interface{}) EnergyEvent {
    // 实现事件转换逻辑
    return EnergyEvent{}
}
