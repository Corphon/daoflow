// api/config.go

package api

import (
    "context"
    "encoding/json"
    "sync"
    "time"

    "github.com/Corphon/daoflow/system"
)

// ConfigScope 配置作用域
type ConfigScope string

const (
    ScopeGlobal     ConfigScope = "global"     // 全局配置
    ScopeComponent  ConfigScope = "component"  // 组件配置
    ScopePattern    ConfigScope = "pattern"    // 模式配置
    ScopeMetric     ConfigScope = "metric"     // 指标配置
    ScopeEvolution  ConfigScope = "evolution"  // 演化配置
)

// ConfigValue 配置值
type ConfigValue struct {
    Value     interface{}            `json:"value"`      // 配置值
    Type      string                `json:"type"`       // 值类型
    Scope     ConfigScope           `json:"scope"`      // 作用域
    Version   int64                 `json:"version"`    // 版本号
    UpdatedAt time.Time            `json:"updated_at"`  // 更新时间
    UpdatedBy string               `json:"updated_by"`  // 更新者
    Metadata  map[string]interface{} `json:"metadata"`   // 元数据
}

// ConfigHistory 配置历史记录
type ConfigHistory struct {
    Key       string      `json:"key"`        // 配置键
    Value     interface{} `json:"value"`      // 配置值
    Version   int64       `json:"version"`    // 版本号
    Timestamp time.Time   `json:"timestamp"`  // 变更时间
    User      string      `json:"user"`       // 操作用户
    Reason    string      `json:"reason"`     // 变更原因
}

// ConfigValidation 配置验证规则
type ConfigValidation struct {
    Required    bool        `json:"required"`     // 是否必需
    Type        string      `json:"type"`         // 值类型
    Range       []interface{} `json:"range"`      // 取值范围
    Pattern     string      `json:"pattern"`      // 匹配模式
    Constraints []string    `json:"constraints"`  // 约束条件
}

// ConfigAPI 配置管理API
type ConfigAPI struct {
    mu     sync.RWMutex
    system *system.SystemCore
    
    // 配置存储
    configs map[string]*ConfigValue
    
    // 配置验证规则
    validations map[string]*ConfigValidation
    
    // 历史记录
    history []*ConfigHistory
    
    // 变更通知
    events chan ConfigEvent
    
    // 上下文控制
    ctx    context.Context
    cancel context.CancelFunc
}

// ConfigEvent 配置事件
type ConfigEvent struct {
    Type      string      `json:"type"`       // 事件类型
    Key       string      `json:"key"`        // 配置键
    Value     *ConfigValue `json:"value"`      // 配置值
    Version   int64       `json:"version"`    // 版本号
    Timestamp time.Time   `json:"timestamp"`  // 事件时间
}

// NewConfigAPI 创建配置API实例
func NewConfigAPI(sys *system.SystemCore) *ConfigAPI {
    ctx, cancel := context.WithCancel(context.Background())
    
    api := &ConfigAPI{
        system:      sys,
        configs:     make(map[string]*ConfigValue),
        validations: make(map[string]*ConfigValidation),
        history:     make([]*ConfigHistory, 0),
        events:      make(chan ConfigEvent, 100),
        ctx:         ctx,
        cancel:      cancel,
    }
    
    // 初始化默认配置验证规则
    api.initDefaultValidations()
    
    return api
}

// SetConfig 设置配置项
func (c *ConfigAPI) SetConfig(key string, value interface{}, scope ConfigScope, metadata map[string]interface{}) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    // 验证配置
    if err := c.validateConfig(key, value); err != nil {
        return err
    }

    // 获取或创建配置项
    config, exists := c.configs[key]
    if !exists {
        config = &ConfigValue{
            Type:  getValueType(value),
            Scope: scope,
        }
        c.configs[key] = config
    }

    // 更新配置
    prevValue := config.Value
    config.Value = value
    config.Version++
    config.UpdatedAt = time.Now()
    config.UpdatedBy = "Corphon" // 使用当前用户
    config.Metadata = metadata

    // 记录历史
    c.recordHistory(key, prevValue, config, "Manual update")

    // 发送事件
    c.events <- ConfigEvent{
        Type:      "config_updated",
        Key:       key,
        Value:     config,
        Version:   config.Version,
        Timestamp: time.Now(),
    }

    return nil
}

// GetConfig 获取配置项
func (c *ConfigAPI) GetConfig(key string) (*ConfigValue, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    config, exists := c.configs[key]
    if !exists {
        return nil, NewError(ErrConfigNotFound, "config not found")
    }

    return config, nil
}

// DeleteConfig 删除配置项
func (c *ConfigAPI) DeleteConfig(key string) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    config, exists := c.configs[key]
    if !exists {
        return NewError(ErrConfigNotFound, "config not found")
    }

    // 记录删除历史
    c.recordHistory(key, config.Value, nil, "Manual deletion")

    delete(c.configs, key)

    // 发送事件
    c.events <- ConfigEvent{
        Type:      "config_deleted",
        Key:       key,
        Timestamp: time.Now(),
    }

    return nil
}

// GetConfigsByScope 获取指定作用域的配置
func (c *ConfigAPI) GetConfigsByScope(scope ConfigScope) map[string]*ConfigValue {
    c.mu.RLock()
    defer c.mu.RUnlock()

    result := make(map[string]*ConfigValue)
    for key, config := range c.configs {
        if config.Scope == scope {
            result[key] = config
        }
    }
    return result
}

// SetValidation 设置配置验证规则
func (c *ConfigAPI) SetValidation(key string, validation *ConfigValidation) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.validations[key] = validation
    return nil
}

// GetHistory 获取配置历史记录
func (c *ConfigAPI) GetHistory(key string) ([]*ConfigHistory, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    var records []*ConfigHistory
    for _, record := range c.history {
        if record.Key == key {
            records = append(records, record)
        }
    }
    return records, nil
}

// Subscribe 订阅配置事件
func (c *ConfigAPI) Subscribe() (<-chan ConfigEvent, error) {
    return c.events, nil
}

// Export 导出配置
func (c *ConfigAPI) Export() ([]byte, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    return json.Marshal(c.configs)
}

// Import 导入配置
func (c *ConfigAPI) Import(data []byte) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    configs := make(map[string]*ConfigValue)
    if err := json.Unmarshal(data, &configs); err != nil {
        return err
    }

    // 验证所有配置
    for key, config := range configs {
        if err := c.validateConfig(key, config.Value); err != nil {
            return err
        }
    }

    // 更新配置
    c.configs = configs
    return nil
}

// validateConfig 验证配置值
func (c *ConfigAPI) validateConfig(key string, value interface{}) error {
    validation, exists := c.validations[key]
    if !exists {
        return nil // 无验证规则时默认通过
    }

    // 实现验证逻辑
    return nil
}

// recordHistory 记录配置历史
func (c *ConfigAPI) recordHistory(key string, oldValue interface{}, newValue *ConfigValue, reason string) {
    record := &ConfigHistory{
        Key:       key,
        Value:     newValue.Value,
        Version:   newValue.Version,
        Timestamp: time.Now(),
        User:      newValue.UpdatedBy,
        Reason:    reason,
    }
    c.history = append(c.history, record)
}

// initDefaultValidations 初始化默认验证规则
func (c *ConfigAPI) initDefaultValidations() {
    // 添加默认验证规则
}

// getValueType 获取值类型
func getValueType(value interface{}) string {
    switch value.(type) {
    case string:
        return "string"
    case int, int32, int64:
        return "integer"
    case float32, float64:
        return "float"
    case bool:
        return "boolean"
    case []interface{}:
        return "array"
    case map[string]interface{}:
        return "object"
    default:
        return "unknown"
    }
}

// Close 关闭API
func (c *ConfigAPI) Close() error {
    c.cancel()
    close(c.events)
    return nil
}
