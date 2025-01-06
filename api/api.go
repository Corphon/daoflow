// api/api.go

package api

import (
    "context"
    "sync"

    "github.com/Corphon/daoflow/system"
)

// APIVersion API 版本
const APIVersion = "v2.0.0"

// DaoFlowAPI 道流API客户端
type DaoFlowAPI struct {
    mu sync.RWMutex

    // 系统核心
    system *system.SystemCore

    // API组件
    lifecycle   *LifecycleAPI    // 生命周期管理
    evolution   *EvolutionAPI    // 演化控制
    energy      *EnergyAPI       // 能量管理
    pattern     *PatternAPI      // 模式识别
    metrics     *MetricsAPI      // 指标监控
    config      *ConfigAPI       // 配置管理
    events      *EventsAPI       // 事件订阅
    health      *HealthAPI       // 健康检查

    // 上下文控制
    ctx    context.Context
    cancel context.CancelFunc

    adapter    *SystemAdapter        // 添加系统适配器
    bufferPool map[string]*DynamicBuffer // 添加缓冲池
}

// Options API选项配置
type Options struct {
    SystemConfig  *system.SystemConfig  // 系统配置
    Timeout       int                   // 超时设置(秒)
    MaxRetries    int                   // 最大重试次数
    AsyncEnabled  bool                  // 是否启用异步
    Debug         bool                  // 是否开启调试
}

// DefaultOptions 默认选项
var DefaultOptions = &Options{
    Timeout:      30,
    MaxRetries:   3,
    AsyncEnabled: true,
    Debug:        false,
}

// NewDaoFlowAPI 创建新的API客户端
func NewDaoFlowAPI(opts *Options) (*DaoFlowAPI, error) {
    if opts == nil {
        opts = DefaultOptions
    }

    ctx, cancel := context.WithCancel(context.Background())

    // 创建系统核心
    systemCore, err := system.NewSystemCore(ctx, opts.SystemConfig)
    if err != nil {
        cancel()
        return nil, err
    }

    api := &DaoFlowAPI{
        system: systemCore,
        ctx:    ctx,
        cancel: cancel,
    }

    // 初始化API组件
    if err := api.initializeComponents(opts); err != nil {
        api.Close()
        return nil, err
    }

    return api, nil
}

// initializeComponents 初始化API组件
func (api *DaoFlowAPI) initializeComponents(opts *Options) error {
    // 初始化生命周期API
    api.lifecycle = NewLifecycleAPI(api.system, opts)

    // 初始化演化API
    api.evolution = NewEvolutionAPI(api.system, opts)

    // 初始化能量API
    api.energy = NewEnergyAPI(api.system, opts)

    // 初始化模式API
    api.pattern = NewPatternAPI(api.system, opts)

    // 初始化指标API
    api.metrics = NewMetricsAPI(api.system, opts)

    // 初始化配置API
    api.config = NewConfigAPI(api.system, opts)

    // 初始化事件API
    api.events = NewEventsAPI(api.system, opts)

    // 初始化健康API
    api.health = NewHealthAPI(api.system, opts)

    return nil
}

// Lifecycle 获取生命周期API
func (api *DaoFlowAPI) Lifecycle() *LifecycleAPI {
    return api.lifecycle
}

// Evolution 获取演化API
func (api *DaoFlowAPI) Evolution() *EvolutionAPI {
    return api.evolution
}

// Energy 获取能量API
func (api *DaoFlowAPI) Energy() *EnergyAPI {
    return api.energy
}

// Pattern 获取模式API
func (api *DaoFlowAPI) Pattern() *PatternAPI {
    return api.pattern
}

// Metrics 获取指标API
func (api *DaoFlowAPI) Metrics() *MetricsAPI {
    return api.metrics
}

// Config 获取配置API
func (api *DaoFlowAPI) Config() *ConfigAPI {
    return api.config
}

// Events 获取事件API
func (api *DaoFlowAPI) Events() *EventsAPI {
    return api.events
}

// Health 获取健康API
func (api *DaoFlowAPI) Health() *HealthAPI {
    return api.health
}

// Version 获取API版本
func (api *DaoFlowAPI) Version() string {
    return APIVersion
}

// Close 关闭API客户端
func (api *DaoFlowAPI) Close() error {
    api.cancel()
    return api.system.Close()
}

// 使用示例
/*
func Example() {
    // 创建API客户端
    opts := &Options{
        SystemConfig: &system.SystemConfig{
            Capacity: 2000.0,
            Threshold: 0.7,
        },
        Debug: true,
    }
    
    client, err := NewDaoFlowAPI(opts)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // 初始化系统
    if err := client.Lifecycle().Initialize(); err != nil {
        log.Fatal(err)
    }

    // 启动系统
    if err := client.Lifecycle().Start(); err != nil {
        log.Fatal(err)
    }

    // 监控演化事件
    eventChan, err := client.Events().Subscribe([]string{"evolution"})
    if err != nil {
        log.Fatal(err)
    }

    go func() {
        for event := range eventChan {
            log.Printf("Evolution event: %+v", event)
        }
    }()

    // 触发演化
    if err := client.Evolution().TriggerEvolution("optimize"); err != nil {
        log.Fatal(err)
    }

    // 获取系统状态
    status, err := client.Lifecycle().GetStatus()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("System status: %+v", status)
}
*/
