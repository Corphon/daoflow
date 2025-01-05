// api/evolution.go

package api

import (
    "context"
    "sync"
    "time"

    "github.com/Corphon/daoflow/system"
)

// EvolutionConfig 演化配置
type EvolutionConfig struct {
    Mode          EvolutionMode       // 演化模式
    Patterns      []PatternConfig     // 模式配置
    Strategies    []StrategyConfig    // 策略配置
    EnergyControl EnergyConfig        // 能量控制
    Constraints   SystemConstraints   // 系统约束
}

// PatternConfig 模式配置
type PatternConfig struct {
    Type      string    // 模式类型
    Weight    float64   // 权重系数
    Threshold float64   // 识别阈值
    Features  []string  // 特征列表
}

// StrategyConfig 策略配置
type StrategyConfig struct {
    Type       string             // 策略类型
    Priority   int                // 优先级
    Conditions map[string]float64 // 触发条件
    Actions    []string           // 动作列表
}

// EnergyConfig 能量控制配置
type EnergyConfig struct {
    BalanceRatio float64 // 平衡比例
    MinLevel     float64 // 最低水平
    MaxLevel     float64 // 最高水平
    Distribution map[string]float64 // 能量分配
}

// SystemConstraints 系统约束
type SystemConstraints struct {
    StabilityThreshold float64 // 稳定性阈值
    AdaptationRate     float64 // 适应速率
    EnergyEfficiency   float64 // 能量效率
    TimeWindow         time.Duration // 时间窗口
}

// EvolutionAPI 演化控制API
type EvolutionAPI struct {
    mu     sync.RWMutex
    system *system.SystemCore
    config EvolutionConfig
}

// NewEvolution 创建演化控制实例
func NewEvolution(sys *system.SystemCore, config EvolutionConfig) *EvolutionAPI {
    return &EvolutionAPI{
        system: sys,
        config: config,
    }
}

// ApplyStrategy 应用演化策略
func (e *EvolutionAPI) ApplyStrategy(ctx context.Context, strategy StrategyConfig) error {
    e.mu.Lock()
    defer e.mu.Unlock()
    
    return e.system.Evolution().ApplyStrategy(strategy)
}

// UpdatePatterns 更新模式识别配置
func (e *EvolutionAPI) UpdatePatterns(patterns []PatternConfig) error {
    e.mu.Lock()
    defer e.mu.Unlock()
    
    return e.system.Evolution().UpdatePatterns(patterns)
}

// AdjustEnergy 调整能量分配
func (e *EvolutionAPI) AdjustEnergy(config EnergyConfig) error {
    e.mu.Lock()
    defer e.mu.Unlock()
    
    return e.system.Evolution().AdjustEnergy(config)
}

// GetEvolutionStatus 获取演化状态
func (e *EvolutionAPI) GetEvolutionStatus() (*EvolutionStatus, error) {
    e.mu.RLock()
    defer e.mu.RUnlock()
    
    return e.system.Evolution().GetStatus()
}

// SubscribeEvents 订阅演化事件
func (e *EvolutionAPI) SubscribeEvents() (<-chan EvolutionEvent, error) {
    return e.system.Evolution().Subscribe()
}
