// model/config.go

package model

import (
    "encoding/json"
    "sync"
    "time"
)

// ModelConfig 模型配置
type ModelConfig struct {
    // 基础配置
    Base struct {
        MaxEnergy     float64 `json:"maxEnergy"`     // 最大能量
        UpdateRate    float64 `json:"updateRate"`    // 更新速率
        SyncInterval  int64   `json:"syncInterval"`  // 同步间隔(ms)
        EnableLogging bool    `json:"enableLogging"` // 启用日志
    }

    // 阴阳模型配置
    YinYang struct {
        BalanceThreshold float64 `json:"balanceThreshold"` // 平衡阈值
        TransformRate    float64 `json:"transformRate"`    // 转换率
        ResonanceRate    float64 `json:"resonanceRate"`    // 共振率
    }

    // 五行模型配置
    WuXing struct {
        CycleThreshold  float64 `json:"cycleThreshold"`  // 循环阈值
        ElementCapacity float64 `json:"elementCapacity"` // 元素容量
        FlowRate        float64 `json:"flowRate"`        // 流动率
    }

    // 八卦模型配置
    BaGua struct {
        ChangeThreshold float64 `json:"changeThreshold"` // 变化阈值
        ResonanceRate   float64 `json:"resonanceRate"`   // 共振率
        HarmonyLevel    float64 `json:"harmonyLevel"`    // 和谐度
    }

    // 干支模型配置
    GanZhi struct {
        CycleLength     int     `json:"cycleLength"`     // 周期长度
        HarmonyThreshold float64 `json:"harmonyThreshold"` // 和谐阈值
        ElementBalance   float64 `json:"elementBalance"`   // 五行平衡度
    }

    // 资源限制
    Resources struct {
        MaxMemory      int64 `json:"maxMemory"`      // 最大内存使用(MB)
        MaxGoroutines  int   `json:"maxGoroutines"`  // 最大协程数
        MaxOpenFiles   int   `json:"maxOpenFiles"`   // 最大打开文件数
    }

    // 监控配置
    Monitoring struct {
        EnableMetrics  bool   `json:"enableMetrics"`  // 启用指标
        MetricsPort    int    `json:"metricsPort"`    // 指标端口
        CollectInterval int64 `json:"collectInterval"` // 采集间隔(ms)
    }

    lastUpdate time.Time // 最后更新时间
    mu         sync.RWMutex
}

// ConfigManager 配置管理器
type ConfigManager struct {
    config     *ModelConfig
    validators []ConfigValidator
    observers  []ConfigObserver
    mu         sync.RWMutex
}

// ConfigValidator 配置验证器接口
type ConfigValidator interface {
    Validate(*ModelConfig) error
}

// ConfigObserver 配置观察者接口
type ConfigObserver interface {
    OnConfigUpdate(*ModelConfig)
}

// NewConfigManager 创建配置管理器
func NewConfigManager() *ConfigManager {
    return &ConfigManager{
        config:     newDefaultConfig(),
        validators: make([]ConfigValidator, 0),
        observers:  make([]ConfigObserver, 0),
    }
}

// newDefaultConfig 创建默认配置
func newDefaultConfig() *ModelConfig {
    config := &ModelConfig{}

    // 设置基础配置默认值
    config.Base.MaxEnergy = 1000.0
    config.Base.UpdateRate = 0.1
    config.Base.SyncInterval = 1000
    config.Base.EnableLogging = true

    // 设置阴阳模型默认值
    config.YinYang.BalanceThreshold = 0.1
    config.YinYang.TransformRate = 0.05
    config.YinYang.ResonanceRate = 0.08

    // 设置五行模型默认值
    config.WuXing.CycleThreshold = 0.3
    config.WuXing.ElementCapacity = 20.0
    config.WuXing.FlowRate = 0.05

    // 设置八卦模型默认值
    config.BaGua.ChangeThreshold = 0.2
    config.BaGua.ResonanceRate = 0.08
    config.BaGua.HarmonyLevel = 0.7

    // 设置干支模型默认值
    config.GanZhi.CycleLength = 60
    config.GanZhi.HarmonyThreshold = 0.7
    config.GanZhi.ElementBalance = 0.8

    // 设置资源限制默认值
    config.Resources.MaxMemory = 1024
    config.Resources.MaxGoroutines = 100
    config.Resources.MaxOpenFiles = 1000

    // 设置监控配置默认值
    config.Monitoring.EnableMetrics = true
    config.Monitoring.MetricsPort = 9090
    config.Monitoring.CollectInterval = 5000

    config.lastUpdate = time.Now()
    return config
}

// LoadConfig 加载配置
func (cm *ConfigManager) LoadConfig(data []byte) error {
    cm.mu.Lock()
    defer cm.mu.Unlock()

    // 解析配置
    newConfig := new(ModelConfig)
    if err := json.Unmarshal(data, newConfig); err != nil {
        return WrapError(err, ErrCodeValidation, "failed to parse config")
    }

    // 验证配置
    if err := cm.validateConfig(newConfig); err != nil {
        return err
    }

    // 更新配置
    cm.config = newConfig
    cm.config.lastUpdate = time.Now()

    // 通知观察者
    cm.notifyObservers()

    return nil
}

// GetConfig 获取配置
func (cm *ConfigManager) GetConfig() *ModelConfig {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    return cm.config
}

// RegisterValidator 注册验证器
func (cm *ConfigManager) RegisterValidator(validator ConfigValidator) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    cm.validators = append(cm.validators, validator)
}

// RegisterObserver 注册观察者
func (cm *ConfigManager) RegisterObserver(observer ConfigObserver) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    cm.observers = append(cm.observers, observer)
}

// validateConfig 验证配置
func (cm *ConfigManager) validateConfig(config *ModelConfig) error {
    for _, validator := range cm.validators {
        if err := validator.Validate(config); err != nil {
            return WrapError(err, ErrCodeValidation, "config validation failed")
        }
    }
    return nil
}

// notifyObservers 通知观察者
func (cm *ConfigManager) notifyObservers() {
    for _, observer := range cm.observers {
        observer.OnConfigUpdate(cm.config)
    }
}

// BaseConfigValidator 基础配置验证器
type BaseConfigValidator struct{}

// Validate 验证基础配置
func (v *BaseConfigValidator) Validate(config *ModelConfig) error {
    // 验证能量限制
    if config.Base.MaxEnergy <= 0 {
        return NewModelError(ErrCodeValidation, "max energy must be positive", nil)
    }

    // 验证更新速率
    if config.Base.UpdateRate <= 0 || config.Base.UpdateRate > 1 {
        return NewModelError(ErrCodeValidation, "update rate must be between 0 and 1", nil)
    }

    // 验证同步间隔
    if config.Base.SyncInterval < 100 {
        return NewModelError(ErrCodeValidation, "sync interval must be at least 100ms", nil)
    }

    return nil
}

// ResourceConfigValidator 资源配置验证器
type ResourceConfigValidator struct{}

// Validate 验证资源配置
func (v *ResourceConfigValidator) Validate(config *ModelConfig) error {
    // 验证内存限制
    if config.Resources.MaxMemory < 128 {
        return NewModelError(ErrCodeValidation, "max memory must be at least 128MB", nil)
    }

    // 验证协程限制
    if config.Resources.MaxGoroutines < 10 {
        return NewModelError(ErrCodeValidation, "max goroutines must be at least 10", nil)
    }

    // 验证文件限制
    if config.Resources.MaxOpenFiles < 100 {
        return NewModelError(ErrCodeValidation, "max open files must be at least 100", nil)
    }

    return nil
}
