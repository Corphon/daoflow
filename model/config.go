// model/config.go

package model

import (
    "time"
)

// ModelConfig 模型配置
type ModelConfig struct {
    // 基础配置
    ModelType       ModelType        // 模型类型
    Capacity        float64         // 能量容量 
    InitialEnergy   float64         // 初始能量
    UpdateInterval  time.Duration   // 更新间隔
    
    // 系统限制
    EnergyMinimum   float64         // 最小能量
    EnergyMaximum   float64         // 最大能量
    StateTimeout    time.Duration   // 状态超时
    
    // 物理参数
    DampingRatio    float64         // 阻尼系数
    Frequency       float64         // 振动频率
    PhaseOffset     float64         // 相位偏移
    
    // 转换参数
    TransformThreshold  float64     // 转换阈值
    TransformCooldown   time.Duration // 转换冷却
    
    // 同步设置
    SyncEnabled     bool            // 启用同步
    SyncInterval    time.Duration   // 同步间隔
}

// DefaultConfig 默认配置
var DefaultConfig = ModelConfig{
    // 基础配置
    Capacity:       1000.0,
    InitialEnergy:  100.0,
    UpdateInterval: time.Second,
    
    // 系统限制
    EnergyMinimum:  0.0,
    EnergyMaximum:  1000.0,
    StateTimeout:   time.Minute * 5,
    
    // 物理参数
    DampingRatio:   0.05,
    Frequency:      2 * math.Pi / (24 * 3600), // 一天一个周期
    PhaseOffset:    0.0,
    
    // 转换参数
    TransformThreshold: 0.7,
    TransformCooldown:  time.Second * 10,
    
    // 同步设置
    SyncEnabled:    true,
    SyncInterval:   time.Second * 5,
}

// 模型特定配置
var (
    // 阴阳模型配置
    YinYangConfig = ModelConfig{
        ModelType:      ModelYinYang,
        Capacity:       200.0,
        InitialEnergy:  100.0,
        DampingRatio:   0.1,
        Frequency:      2 * math.Pi / (12 * 3600), // 12小时周期
    }
    
    // 五行模型配置
    WuXingConfig = ModelConfig{
        ModelType:      ModelWuXing,
        Capacity:       500.0,
        InitialEnergy:  300.0,
        DampingRatio:   0.08,
        Frequency:      2 * math.Pi / (24 * 3600), // 24小时周期
    }
    
    // 八卦模型配置
    BaGuaConfig = ModelConfig{
        ModelType:      ModelBaGua,
        Capacity:       800.0,
        InitialEnergy:  400.0,
        DampingRatio:   0.06,
        Frequency:      2 * math.Pi / (3 * 24 * 3600), // 3天周期
    }
    
    // 天干地支模型配置
    GanZhiConfig = ModelConfig{
        ModelType:      ModelGanZhi,
        Capacity:       2200.0,
        InitialEnergy:  1000.0,
        DampingRatio:   0.03,
        Frequency:      2 * math.Pi / (60 * 24 * 3600), // 60天周期
    }
)

// ConfigOption 配置选项函数类型
type ConfigOption func(*ModelConfig)

// WithCapacity 设置容量
func WithCapacity(capacity float64) ConfigOption {
    return func(c *ModelConfig) {
        c.Capacity = capacity
    }
}

// WithInitialEnergy 设置初始能量
func WithInitialEnergy(energy float64) ConfigOption {
    return func(c *ModelConfig) {
        c.InitialEnergy = energy
    }
}

// WithUpdateInterval 设置更新间隔
func WithUpdateInterval(interval time.Duration) ConfigOption {
    return func(c *ModelConfig) {
        c.UpdateInterval = interval
    }
}

// WithPhysicsParams 设置物理参数
func WithPhysicsParams(damping, frequency, phase float64) ConfigOption {
    return func(c *ModelConfig) {
        c.DampingRatio = damping
        c.Frequency = frequency
        c.PhaseOffset = phase
    }
}

// WithTransformParams 设置转换参数
func WithTransformParams(threshold float64, cooldown time.Duration) ConfigOption {
    return func(c *ModelConfig) {
        c.TransformThreshold = threshold
        c.TransformCooldown = cooldown
    }
}

// WithSyncParams 设置同步参数
func WithSyncParams(enabled bool, interval time.Duration) ConfigOption {
    return func(c *ModelConfig) {
        c.SyncEnabled = enabled
        c.SyncInterval = interval
    }
}

// NewConfig 创建新配置
func NewConfig(modelType ModelType, opts ...ConfigOption) ModelConfig {
    // 基于模型类型选择基础配置
    var config ModelConfig
    switch modelType {
    case ModelYinYang:
        config = YinYangConfig
    case ModelWuXing:
        config = WuXingConfig
    case ModelBaGua:
        config = BaGuaConfig
    case ModelGanZhi:
        config = GanZhiConfig
    default:
        config = DefaultConfig
    }
    
    // 应用配置选项
    for _, opt := range opts {
        opt(&config)
    }
    
    return config
}

// Validate 验证配置
func (c ModelConfig) Validate() error {
    if c.Capacity <= 0 {
        return NewModelError(ErrCodeConfiguration, "invalid capacity", nil)
    }
    
    if c.InitialEnergy < c.EnergyMinimum || c.InitialEnergy > c.EnergyMaximum {
        return NewModelError(ErrCodeConfiguration, "invalid initial energy", nil)
    }
    
    if c.UpdateInterval < time.Millisecond {
        return NewModelError(ErrCodeConfiguration, "update interval too small", nil)
    }
    
    if c.DampingRatio < 0 || c.DampingRatio > 1 {
        return NewModelError(ErrCodeConfiguration, "invalid damping ratio", nil)
    }
    
    return nil
}
