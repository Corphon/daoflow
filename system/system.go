//system/system.go

package system

import (
    "fmt"
    "sync"
    "time"

    // 内部包导入
    "github.com/Corphon/daoflow/system/meta/field"
    "github.com/Corphon/daoflow/system/meta/emergence"
    "github.com/Corphon/daoflow/system/meta/resonance"
    "github.com/Corphon/daoflow/system/evolution/pattern"
    "github.com/Corphon/daoflow/system/evolution/mutation"
    "github.com/Corphon/daoflow/system/evolution/adaptation"
    "github.com/Corphon/daoflow/system/control/state"
    "github.com/Corphon/daoflow/system/control/flow"
    "github.com/Corphon/daoflow/system/control/sync"
    "github.com/Corphon/daoflow/system/monitor/metrics"
    "github.com/Corphon/daoflow/system/monitor/trace"
    "github.com/Corphon/daoflow/system/monitor/alert"
    "github.com/Corphon/daoflow/system/types"
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
}

// SystemOptions 系统配置选项
type SystemOptions struct {
    Name       string                      // 系统名称
    Version    string                      // 系统版本
    ConfigPath string                      // 配置路径
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
        return fmt.Errorf("system in invalid state for starting: %s", s.state.status)
    }

    // 启动各子系统
    if err := s.startSubSystems(); err != nil {
        return fmt.Errorf("failed to start subsystems: %v", err)
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
