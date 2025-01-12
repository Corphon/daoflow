//system/system.go

package system

import (
    // 标准库
    "context"
    "errors"
    "fmt"
    "os"
    "runtime"
    "sync"
    "sync/atomic"
    "time"

    // 项目基础包
    "github.com/Corphon/daoflow/system/common"
    "github.com/Corphon/daoflow/system/config"
    "github.com/Corphon/daoflow/system/resource"
    "github.com/Corphon/daoflow/system/types"
    
    // 元系统包
    "github.com/Corphon/daoflow/system/meta/field"
    "github.com/Corphon/daoflow/system/meta/emergence" 
    "github.com/Corphon/daoflow/system/meta/resonance"
    
    // 演化系统包
    "github.com/Corphon/daoflow/system/evolution/pattern"
    "github.com/Corphon/daoflow/system/evolution/mutation"
    "github.com/Corphon/daoflow/system/evolution/adaptation"
    
    // 控制系统包
    "github.com/Corphon/daoflow/system/control/state"
    "github.com/Corphon/daoflow/system/control/flow"
    csync "github.com/Corphon/daoflow/system/control/sync" // 重命名避免冲突
    
    // 监控系统包
    "github.com/Corphon/daoflow/system/monitor/metrics"
    "github.com/Corphon/daoflow/system/monitor/trace"
    "github.com/Corphon/daoflow/system/monitor/alert"
)

// System DaoFlow系统主结构
type System struct {
    mu sync.RWMutex

    // 系统配置
    config struct {
        Name          string                // 系统名称
        Version       string                // 系统版本
        StartTime     time.Time            // 启动时间
        ConfigPath    string                // 配置路径
    }

    // 元系统组件
    meta struct {
        fieldSystem      *field.System      // 场系统
        emergenceSystem  *emergence.System  // 涌现系统
        resonanceSystem  *resonance.System  // 共振系统
    }

    // 演化系统组件
    evolution struct {
        patternSystem    *pattern.System    // 模式系统
        mutationSystem   *mutation.System   // 突变系统
        adaptationSystem *adaptation.System // 适应系统
    }

    // 控制系统组件
    control struct {
        stateManager     *state.Manager     // 状态管理器
        flowController   *flow.Controller   // 流控制器
        syncController   *sync.Controller   // 同步控制器
    }

    // 监控系统组件
    monitor struct {
        metricsSystem    *metrics.System    // 指标系统
        traceSystem      *trace.System      // 追踪系统
        alertSystem      *alert.System      // 告警系统
    }

    // 系统状态
    state struct {
        status        string                // 系统状态
        health        float64               // 健康度
        lastCheck     time.Time            // 最后检查
        metrics       SystemMetrics         // 系统指标
    }
}

// SystemMetrics 系统指标
type SystemMetrics struct {
    Uptime         time.Duration           // 运行时间
    MemoryUsage    uint64                  // 内存使用
    CPUUsage       float64                 // CPU使用率
    GoroutineCount int                     // 协程数量
    ErrorCount     int64                   // 错误计数
    EventCount     int64                   // 事件计数

    // 新增指标
    HeapObjects    uint64        // 堆对象数
    HeapAlloc      uint64        // 堆分配量
    GCPause        time.Duration // GC暂停时间
    ThreadCount    int           // 线程数
    MutexWait      int64        // 互斥锁等待次数
    NetworkIO      NetworkStats  // 网络IO统计
    DiskIO        DiskStats     // 磁盘IO统计
}

// SystemOptions 系统配置选项
type SystemOptions struct {
    Name       string                      // 系统名称
    Version    string                      // 系统版本
    ConfigPath string                      // 配置路径
}

type NetworkStats struct {
    BytesRead    uint64
    BytesWritten uint64
    Connections  int
}

type DiskStats struct {
    BytesRead    uint64
    BytesWritten uint64
    Operations   int64
}

// 添加系统状态常量
const (
    SystemStateInitialized = "initialized"
    SystemStateRunning     = "running"
    SystemStateStopping    = "stopping"
    SystemStateStopped     = "stopped"
    SystemStateError       = "error"
)

// 添加资源管理
type SystemResources struct {
    CPU     *CPUResource
    Memory  *MemoryResource
    Network *NetworkResource
    Disk    *DiskResource
}

func (s *System) initializeResources() error {
    s.resources = &SystemResources{
        CPU:     NewCPUResource(s.config.CPULimit),
        Memory:  NewMemoryResource(s.config.MemoryLimit),
        Network: NewNetworkResource(s.config.NetworkLimit),
        Disk:    NewDiskResource(s.config.DiskLimit),
    }
    return s.resources.Initialize()
}

func (s *System) monitorResources() {
    ticker := time.NewTicker(time.Second)
    for range ticker.C {
        s.updateResourceMetrics()
    }
}

// 添加健康检查
type HealthCheck struct {
    Name       string
    Check      func() error
    Interval   time.Duration
    Timeout    time.Duration
    Required   bool
}

func (s *System) initializeHealthChecks() {
    s.healthChecks = []HealthCheck{
        {
            Name:     "memory",
            Check:    s.checkMemoryHealth,
            Interval: time.Second * 30,
            Required: true,
        },
        {
            Name:     "subsystems",
            Check:    s.checkSubsystemsHealth,
            Interval: time.Minute,
            Required: true,
        },
    }
}

// 添加动态配置支持
type SystemConfig struct {
    Basic    BasicConfig
    Meta     MetaConfig
    Evolution EvolutionConfig
    Control   ControlConfig
    Monitor   MonitorConfig
    Resources ResourceConfig
}

func (s *System) UpdateConfig(newConfig SystemConfig) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if err := s.validateConfig(newConfig); err != nil {
        return err
    }

    if err := s.applyConfig(newConfig); err != nil {
        return err
    }

    return s.notifyConfigUpdate()
}

// 添加事件系统
type SystemEvent struct {
    Type      string
    Source    string
    Timestamp time.Time
    Data      map[string]interface{}
}

func (s *System) EmitEvent(evt SystemEvent) {
    s.eventChan <- evt
    s.incrementEventCount()
}

func (s *System) processEvents() {
    for evt := range s.eventChan {
        s.handleEvent(evt)
    }
}

// 添加状态转换验证
func (s *System) validateStateTransition(targetState string) error {
    validTransitions := map[string][]string{
        SystemStateInitialized: {SystemStateRunning},
        SystemStateRunning:     {SystemStateStopping, SystemStateError},
        SystemStateStopping:    {SystemStateStopped, SystemStateError},
        SystemStateStopped:     {SystemStateInitialized},
        SystemStateError:       {SystemStateInitialized},
    }

    if valid, ok := validTransitions[s.state.status]; ok {
        for _, v := range valid {
            if v == targetState {
                return nil
            }
        }
    }

    return types.NewSystemError(
        types.ErrState,
        "invalid state transition",
        nil,
    ).WithContext("from", s.state.status).
      WithContext("to", targetState)
}

// NewSystem 创建新的DaoFlow系统
func NewSystem(opts SystemOptions) (*System, error) {
    s := &System{}

    // 初始化配置
    s.config.Name = opts.Name
    s.config.Version = opts.Version
    s.config.StartTime = time.Now()
    s.config.ConfigPath = opts.ConfigPath

    // 初始化元系统
    if err := s.initMetaSystems(); err != nil {
        return nil, fmt.Errorf("failed to initialize meta systems: %v", err)
    }

    // 初始化演化系统
    if err := s.initEvolutionSystems(); err != nil {
        return nil, fmt.Errorf("failed to initialize evolution systems: %v", err)
    }

    // 初始化控制系统
    if err := s.initControlSystems(); err != nil {
        return nil, fmt.Errorf("failed to initialize control systems: %v", err)
    }

    // 初始化监控系统
    if err := s.initMonitorSystems(); err != nil {
        return nil, fmt.Errorf("failed to initialize monitor systems: %v", err)
    }

    // 初始化系统状态
    s.state.status = "initialized"
    s.state.health = 1.0
    s.state.lastCheck = time.Now()

    return s, nil
}

// Start 启动系统
func (s *System) Start() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if s.state.status != "initialized" && s.state.status != "stopped" {
        return types.NewSystemError(
            types.ErrState,
            "invalid system state for starting",
            nil,
        ).WithContext("current_state", s.state.status)
    }

    if err := s.startSubSystems(); err != nil {
        return types.WrapSystemError(
            err,
            types.ErrSystem,
            "failed to start subsystems",
        )
    }

    s.state.status = "running"
    return nil
}

// Stop 停止系统
func (s *System) Stop() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if s.state.status != "running" {
        return fmt.Errorf("system not in running state: %s", s.state.status)
    }

    // 停止各子系统
    if err := s.stopSubSystems(); err != nil {
        return fmt.Errorf("failed to stop subsystems: %v", err)
    }

    s.state.status = "stopped"
    return nil
}

// GetStatus 获取系统状态
func (s *System) GetStatus() types.SystemStatus {
    s.mu.RLock()
    defer s.mu.RUnlock()

    return types.SystemStatus{
        Status:    s.state.status,
        Health:    s.state.health,
        LastCheck: s.state.lastCheck,
        Metrics:   s.state.metrics,
    }
}

// UpdateMetrics 更新系统指标
func (s *System) UpdateMetrics() {
    s.mu.Lock()
    defer s.mu.Unlock()

    s.state.metrics = SystemMetrics{
        Uptime:         time.Since(s.config.StartTime),
        MemoryUsage:    s.calculateMemoryUsage(),
        CPUUsage:       s.calculateCPUUsage(),
        GoroutineCount: s.getGoroutineCount(),
        ErrorCount:     s.getErrorCount(),
        EventCount:     s.getEventCount(),
    }

    s.state.lastCheck = time.Now()
}

// 添加系统生命周期方法
func (s *System) Initialize(ctx context.Context) error {
    if err := s.loadConfig(); err != nil {
        return err
    }
    if err := s.validateConfig(); err != nil {
        return err
    }
    return s.initializeComponents()
}

func (s *System) Shutdown(ctx context.Context) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if err := s.validateStateTransition(SystemStateStopping); err != nil {
        return err
    }

    s.state.status = SystemStateStopping
    
    // 优雅关闭
    if err := s.gracefulShutdown(ctx); err != nil {
        s.state.status = SystemStateError
        return err
    }

    s.state.status = SystemStateStopped
    return nil
}

// 内部辅助方法
func (s *System) initMetaSystems() error {
    // 初始化场系统
    fieldSys, err := field.NewSystem()
    if err != nil {
        return err
    }
    s.meta.fieldSystem = fieldSys

    // 初始化涌现系统
    emergenceSys, err := emergence.NewSystem()
    if err != nil {
        return err
    }
    s.meta.emergenceSystem = emergenceSys

    // 初始化共振系统
    resonanceSys, err := resonance.NewSystem()
    if err != nil {
        return err
    }
    s.meta.resonanceSystem = resonanceSys

    return nil
}

func (s *System) initEvolutionSystems() error {
    // 初始化模式系统
    patternSys, err := pattern.NewSystem()
    if err != nil {
        return err
    }
    s.evolution.patternSystem = patternSys

    // 初始化突变系统
    mutationSys, err := mutation.NewSystem()
    if err != nil {
        return err
    }
    s.evolution.mutationSystem = mutationSys

    // 初始化适应系统
    adaptationSys, err := adaptation.NewSystem()
    if err != nil {
        return err
    }
    s.evolution.adaptationSystem = adaptationSys

    return nil
}

func (s *System) initControlSystems() error {
    // 初始化状态管理器
    stateManager, err := state.NewManager()
    if err != nil {
        return err
    }
    s.control.stateManager = stateManager

    // 初始化流控制器
    flowController, err := flow.NewController()
    if err != nil {
        return err
    }
    s.control.flowController = flowController

    // 初始化同步控制器
    syncController, err := sync.NewController()
    if err != nil {
        return err
    }
    s.control.syncController = syncController

    return nil
}

func (s *System) initMonitorSystems() error {
    // 初始化指标系统
    metricsSys, err := metrics.NewSystem()
    if err != nil {
        return err
    }
    s.monitor.metricsSystem = metricsSys

    // 初始化追踪系统
    traceSys, err := trace.NewSystem()
    if err != nil {
        return err
    }
    s.monitor.traceSystem = traceSys

    // 初始化告警系统
    alertSys, err := alert.NewSystem()
    if err != nil {
        return err
    }
    s.monitor.alertSystem = alertSys

    return nil
}

func (s *System) startSubSystems() error {
    // 按顺序启动各子系统
    if err := s.startMetaSystems(); err != nil {
        return err
    }
    if err := s.startEvolutionSystems(); err != nil {
        return err
    }
    if err := s.startControlSystems(); err != nil {
        return err
    }
    if err := s.startMonitorSystems(); err != nil {
        return err
    }
    return nil
}

func (s *System) stopSubSystems() error {
    // 按顺序停止各子系统
    if err := s.stopMonitorSystems(); err != nil {
        return err
    }
    if err := s.stopControlSystems(); err != nil {
        return err
    }
    if err := s.stopEvolutionSystems(); err != nil {
        return err
    }
    if err := s.stopMetaSystems(); err != nil {
        return err
    }
    return nil
}

// 计算系统指标的辅助方法
func (s *System) calculateMemoryUsage() uint64 {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    // 返回当前程序使用的内存量(字节)
    // Alloc 表示已分配的内存
    // Sys 表示从系统获取的内存
    return m.Alloc
}

func (s *System) calculateCPUUsage() float64 {
    // 使用简单的CPU使用率计算方法
    startTime := time.Now()
    startCPU := getCPUTime()
    
    // 等待一小段时间来计算CPU使用率
    time.Sleep(100 * time.Millisecond)
    
    endTime := time.Now()
    endCPU := getCPUTime()
    
    // 计算CPU使用率
    cpuTime := endCPU - startCPU
    realTime := endTime.Sub(startTime).Seconds()
    
    if realTime > 0 {
        return cpuTime / realTime * 100.0
    }
    return 0.0
}

// getCPUTime 获取CPU时间
func getCPUTime() float64 {
    // 使用runtime.ReadMemStats来间接评估CPU时间
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    // 使用GC时间和系统CPU时间的组合来估算
    return float64(m.GCSys) / 1e9
}

func (s *System) getGoroutineCount() int {
    return runtime.NumGoroutine()
}

// 错误计数器
var errorCounter int64 = 0

func (s *System) getErrorCount() int64 {
    return atomic.LoadInt64(&errorCounter)
}

func (s *System) incrementErrorCount() {
    atomic.AddInt64(&errorCounter, 1)
}

// 事件计数器
var eventCounter int64 = 0

func (s *System) getEventCount() int64 {
    return atomic.LoadInt64(&eventCounter)
}

func (s *System) incrementEventCount() {
    atomic.AddInt64(&eventCounter, 1)
}

// 启动各个子系统组的方法
func (s *System) startMetaSystems() error {
    if s.meta.fieldSystem != nil {
        if err := s.meta.fieldSystem.Start(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to start field system: %v", err)
        }
    }

    if s.meta.emergenceSystem != nil {
        if err := s.meta.emergenceSystem.Start(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to start emergence system: %v", err)
        }
    }

    if s.meta.resonanceSystem != nil {
        if err := s.meta.resonanceSystem.Start(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to start resonance system: %v", err)
        }
    }

    s.incrementEventCount()
    return nil
}

func (s *System) startEvolutionSystems() error {
    if s.evolution.patternSystem != nil {
        if err := s.evolution.patternSystem.Start(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to start pattern system: %v", err)
        }
    }

    if s.evolution.mutationSystem != nil {
        if err := s.evolution.mutationSystem.Start(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to start mutation system: %v", err)
        }
    }

    if s.evolution.adaptationSystem != nil {
        if err := s.evolution.adaptationSystem.Start(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to start adaptation system: %v", err)
        }
    }

    s.incrementEventCount()
    return nil
}

func (s *System) startControlSystems() error {
    if s.control.stateManager != nil {
        if err := s.control.stateManager.Start(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to start state manager: %v", err)
        }
    }

    if s.control.flowController != nil {
        if err := s.control.flowController.Start(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to start flow controller: %v", err)
        }
    }

    if s.control.syncController != nil {
        if err := s.control.syncController.Start(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to start sync controller: %v", err)
        }
    }

    s.incrementEventCount()
    return nil
}

func (s *System) startMonitorSystems() error {
    if s.monitor.metricsSystem != nil {
        if err := s.monitor.metricsSystem.Start(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to start metrics system: %v", err)
        }
    }

    if s.monitor.traceSystem != nil {
        if err := s.monitor.traceSystem.Start(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to start trace system: %v", err)
        }
    }

    if s.monitor.alertSystem != nil {
        if err := s.monitor.alertSystem.Start(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to start alert system: %v", err)
        }
    }

    s.incrementEventCount()
    return nil
}

// 停止各个子系统组的方法
func (s *System) stopMetaSystems() error {
    if s.meta.resonanceSystem != nil {
        if err := s.meta.resonanceSystem.Stop(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to stop resonance system: %v", err)
        }
    }

    if s.meta.emergenceSystem != nil {
        if err := s.meta.emergenceSystem.Stop(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to stop emergence system: %v", err)
        }
    }

    if s.meta.fieldSystem != nil {
        if err := s.meta.fieldSystem.Stop(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to stop field system: %v", err)
        }
    }

    s.incrementEventCount()
    return nil
}

func (s *System) stopEvolutionSystems() error {
    if s.evolution.adaptationSystem != nil {
        if err := s.evolution.adaptationSystem.Stop(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to stop adaptation system: %v", err)
        }
    }

    if s.evolution.mutationSystem != nil {
        if err := s.evolution.mutationSystem.Stop(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to stop mutation system: %v", err)
        }
    }

    if s.evolution.patternSystem != nil {
        if err := s.evolution.patternSystem.Stop(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to stop pattern system: %v", err)
        }
    }

    s.incrementEventCount()
    return nil
}

func (s *System) stopControlSystems() error {
    if s.control.syncController != nil {
        if err := s.control.syncController.Stop(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to stop sync controller: %v", err)
        }
    }

    if s.control.flowController != nil {
        if err := s.control.flowController.Stop(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to stop flow controller: %v", err)
        }
    }

    if s.control.stateManager != nil {
        if err := s.control.stateManager.Stop(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to stop state manager: %v", err)
        }
    }

    s.incrementEventCount()
    return nil
}

func (s *System) stopMonitorSystems() error {
    if s.monitor.alertSystem != nil {
        if err := s.monitor.alertSystem.Stop(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to stop alert system: %v", err)
        }
    }

    if s.monitor.traceSystem != nil {
        if err := s.monitor.traceSystem.Stop(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to stop trace system: %v", err)
        }
    }

    if s.monitor.metricsSystem != nil {
        if err := s.monitor.metricsSystem.Stop(); err != nil {
            s.incrementErrorCount()
            return fmt.Errorf("failed to stop metrics system: %v", err)
        }
    }

    s.incrementEventCount()
    return nil
}
