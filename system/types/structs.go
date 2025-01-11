// system/types/structs.go

package types

import (
    "sync"
    "time"
    
    "github.com/Corphon/daoflow/model"
)

// 复用模型包的基础类型
type (
    Vector3D = model.Vector3D
    ModelState = model.ModelState
    FieldState = model.FieldState
    QuantumState = model.QuantumState
)

// System 系统结构
type System struct {
    mu sync.RWMutex

    // 基础信息
    ID        string            // 系统ID
    Name      string            // 系统名称
    Version   string            // 系统版本
    StartTime time.Time         // 启动时间

    // 核心模型 - 使用 model 包的流模型
    models struct {
        yinyang *model.YinYangFlow
        wuxing  *model.WuXingFlow
        bagua   *model.BaGuaFlow
        ganzhi  *model.GanZhiFlow
        unified *model.IntegrateFlow
    }

    // 系统状态
    state struct {
        current  model.SystemState    // 当前状态
        previous model.SystemState    // 前一状态
        changes  []StateTransition   // 状态变更历史
    }

    // 系统组件
    components struct {
        meta      *MetaSystem      // 元系统组件
        evolution *EvolutionSystem // 演化系统组件
        control   *ControlSystem   // 控制系统组件
        monitor   *MonitorSystem   // 监控系统组件
        resource  *ResourceSystem  // 资源系统组件
    }

    // 系统配置
    config SystemConfig
}

// MetaSystem 元系统组件
type MetaSystem struct {
    // 场状态
    field struct {
        state    model.FieldState    // 场状态
        quantum  model.QuantumState  // 量子状态
        coupling [][]float64         // 场耦合矩阵
    }

    // 涌现状态
    emergence struct {
        patterns  []EmergentPattern   // 涌现模式
        active    []EmergentProperty  // 活跃属性
        potential []PotentialEmergence // 潜在涌现
    }

    // 共振状态
    resonance struct {
        state     ResonanceState // 共振状态
        coherence float64       // 相干度
        phase     float64       // 相位
    }
}

// EvolutionSystem 演化系统组件
type EvolutionSystem struct {
    // 当前状态
    current struct {
        level     float64        // 演化层级
        direction model.Vector3D // 演化方向
        speed     float64        // 演化速度
        energy    float64        // 演化能量
    }

    // 演化历史
    history struct {
        path    []EvolutionPoint    // 演化路径
        changes []StateTransition   // 状态变更
        metrics []EvolutionMetrics  // 演化指标
    }
}

// ControlSystem 控制系统组件
type ControlSystem struct {
    // 状态控制
    state struct {
        manager    *StateManager   // 状态管理器
        validator  *StateValidator // 状态验证器
        transition *StateTransitor // 状态转换器
    }

    // 流控制
    flow struct {
        scheduler    *FlowScheduler   // 调度器
        balancer     *FlowBalancer    // 平衡器
        backpressure *BackPressure    // 背压控制
    }

    // 同步控制
    sync struct {
        coordinator  *SyncCoordinator // 同步协调器
        resolver     *ConflictResolver // 冲突解决器
        synchronizer *StateSynchronizer // 状态同步器
    }
}

// MonitorSystem 监控系统组件
type MonitorSystem struct {
    // 指标监控
    metrics struct {
        collector *MetricsCollector // 指标收集器
        storage   *MetricsStorage   // 指标存储
        analyzer  *MetricsAnalyzer  // 指标分析器
    }

    // 告警管理
    alerts struct {
        manager   *AlertManager    // 告警管理器
        handler   *AlertHandler    // 告警处理器
        notifier  *AlertNotifier   // 告警通知器
    }

    // 健康检查
    health struct {
        checker   *HealthChecker   // 健康检查器
        reporter  *HealthReporter  // 健康报告器
        diagnoser *HealthDiagnoser // 健康诊断器
    }
}

// ResourceSystem 资源系统组件
type ResourceSystem struct {
    // 资源池
    pool struct {
        cpu    *ResourcePool // CPU资源池
        memory *ResourcePool // 内存资源池
        energy *ResourcePool // 能量资源池
    }

    // 资源管理
    management struct {
        allocator   *ResourceAllocator   // 资源分配器
        scheduler   *ResourceScheduler   // 资源调度器
        optimizer   *ResourceOptimizer   // 资源优化器
    }

    // 资源监控
    monitor struct {
        collector *ResourceCollector     // 资源收集器
        analyzer  *ResourceAnalyzer      // 资源分析器
        predictor *ResourcePredictor     // 资源预测器
    }
}

// SystemConfig 系统配置
type SystemConfig struct {
    // 基础配置
    Base struct {
        ID        string                 // 系统ID
        Name      string                 // 系统名称
        Version   string                 // 系统版本
        LogLevel  string                 // 日志级别
    }

    // 模型配置 - 使用 model 包的配置
    Model model.ModelConfig

    // 组件配置
    Components struct {
        Meta      MetaConfig      // 元系统配置
        Evolution EvoConfig       // 演化配置
        Control   ControlConfig   // 控制配置
        Monitor   MonitorConfig   // 监控配置
        Resource  ResourceConfig  // 资源配置
    }

    // 系统限制
    Limits struct {
        MaxGoroutines int     // 最大协程数
        MaxMemory     int64   // 最大内存使用
        MaxCPU       float64  // 最大CPU使用
    }
}

// 配置结构实现
// MetaConfig 元系统配置
type MetaConfig struct {
    // 场配置
    Field struct {
        InitialStrength float64            `json:"initial_strength"` // 初始场强度
        MinStrength    float64            `json:"min_strength"`     // 最小场强度
        MaxStrength    float64            `json:"max_strength"`     // 最大场强度
        Dimension      int                `json:"dimension"`        // 场维度
        UpdateInterval time.Duration      `json:"update_interval"` // 场更新间隔
        
        // 场相互作用配置
        Coupling struct {
            Strength    float64   `json:"strength"`    // 耦合强度
            Range       float64   `json:"range"`       // 作用范围
            Threshold   float64   `json:"threshold"`   // 耦合阈值
            MaxPairs    int       `json:"max_pairs"`   // 最大耦合对数
        } `json:"coupling"`
    } `json:"field"`

    // 量子配置
    Quantum struct {
        InitialState    []complex128 `json:"initial_state"`    // 初始量子态
        Coherence      float64      `json:"coherence"`        // 相干度阈值
        DecoherenceRate float64     `json:"decoherence_rate"` // 退相干率
        MeasureInterval time.Duration `json:"measure_interval"` // 测量间隔
        
        // 纠缠配置
        Entanglement struct {
            MaxPairs     int      `json:"max_pairs"`      // 最大纠缠对数
            Threshold    float64  `json:"threshold"`      // 纠缠阈值
            Lifetime     time.Duration `json:"lifetime"`  // 纠缠寿命
        } `json:"entanglement"`
    } `json:"quantum"`

    // 涌现配置
    Emergence struct {
        DetectionInterval time.Duration `json:"detection_interval"` // 检测间隔
        MinStrength      float64       `json:"min_strength"`       // 最小强度阈值
        MaxPatterns      int           `json:"max_patterns"`       // 最大模式数
        
        // 模式配置
        Patterns struct {
            MinLifetime   time.Duration `json:"min_lifetime"`   // 最小生命周期
            StabilityThreshold float64 `json:"stability_threshold"` // 稳定性阈值
            EnergyThreshold   float64 `json:"energy_threshold"`    // 能量阈值
        } `json:"patterns"`
    } `json:"emergence"`

    // 共振配置
    Resonance struct {
        FrequencyRange [2]float64    `json:"frequency_range"` // 频率范围
        MinAmplitude   float64       `json:"min_amplitude"`   // 最小振幅
        PhaseThreshold float64       `json:"phase_threshold"` // 相位阈值
        
        // 共振条件
        Conditions struct {
            MinCoupling    float64   `json:"min_coupling"`    // 最小耦合强度
            MinCoherence   float64   `json:"min_coherence"`   // 最小相干度
            MaxPhaseShift  float64   `json:"max_phase_shift"` // 最大相位偏移
        } `json:"conditions"`
    } `json:"resonance"`
}

// EvoConfig 演化系统配置
type EvoConfig struct {
    // 基本演化参数
    Base struct {
        InitialLevel    float64       `json:"initial_level"`    // 初始演化级别
        MinLevel        float64       `json:"min_level"`        // 最小演化级别
        MaxLevel        float64       `json:"max_level"`        // 最大演化级别
        UpdateInterval  time.Duration `json:"update_interval"`  // 更新间隔
    } `json:"base"`

    // 能量配置
    Energy struct {
        InitialEnergy   float64   `json:"initial_energy"`   // 初始能量
        MinEnergy       float64   `json:"min_energy"`       // 最小能量
        MaxEnergy       float64   `json:"max_energy"`       // 最大能量
        DissipationRate float64   `json:"dissipation_rate"` // 能量耗散率
    } `json:"energy"`

    // 路径配置
    Path struct {
        MaxPoints       int       `json:"max_points"`       // 最大路径点数
        MinDistance     float64   `json:"min_distance"`     // 最小点距离
        MaxDeviation    float64   `json:"max_deviation"`    // 最大偏差
        OptimizeInterval time.Duration `json:"optimize_interval"` // 优化间隔
    } `json:"path"`

    // 状态转换配置
    Transition struct {
        MinEnergy      float64       `json:"min_energy"`      // 最小转换能量
        MaxRetries     int           `json:"max_retries"`     // 最大重试次数
        CooldownPeriod time.Duration `json:"cooldown_period"` // 冷却期
        
        // 转换规则
        Rules struct {
            AllowedStates []string   `json:"allowed_states"`  // 允许的状态
            Priorities    map[string]int `json:"priorities"`   // 状态优先级
            Constraints   map[string]float64 `json:"constraints"` // 约束条件
        } `json:"rules"`
    } `json:"transition"`

    // 适应性配置
    Adaptation struct {
        LearningRate   float64   `json:"learning_rate"`   // 学习率
        MemorySize     int       `json:"memory_size"`     // 记忆大小
        UpdateThreshold float64  `json:"update_threshold"` // 更新阈值
        
        // 策略配置
        Strategy struct {
            Type           string    `json:"type"`            // 策略类型
            Parameters     map[string]float64 `json:"parameters"` // 策略参数
            UpdateInterval time.Duration `json:"update_interval"` // 更新间隔
        } `json:"strategy"`
    } `json:"adaptation"`
}

// ControlConfig 控制系统配置
type ControlConfig struct {
    // 基本控制参数
    Base struct {
        UpdateRate     int           `json:"update_rate"`     // 更新频率
        MaxLatency     time.Duration `json:"max_latency"`     // 最大延迟
        BufferSize     int           `json:"buffer_size"`     // 缓冲区大小
        Timeout        time.Duration `json:"timeout"`         // 超时时间
    } `json:"base"`

    // 反馈控制
    Feedback struct {
        Enabled        bool      `json:"enabled"`         // 是否启用
        Sensitivity    float64   `json:"sensitivity"`     // 灵敏度
        ResponseTime   time.Duration `json:"response_time"` // 响应时间
        
        // PID控制器参数
        PID struct {
            Proportional float64 `json:"proportional"` // 比例系数
            Integral     float64 `json:"integral"`     // 积分系数
            Derivative   float64 `json:"derivative"`   // 微分系数
            WindupGuard float64 `json:"windup_guard"` // 积分限幅
        } `json:"pid"`
    } `json:"feedback"`

    // 稳定性控制
    Stability struct {
        CheckInterval  time.Duration `json:"check_interval"`  // 检查间隔
        MinThreshold   float64      `json:"min_threshold"`   // 最小阈值
        MaxDeviation   float64      `json:"max_deviation"`   // 最大偏差
        
        // 修正参数
        Correction struct {
            Strength     float64   `json:"strength"`     // 修正强度
            MaxAttempts  int       `json:"max_attempts"` // 最大尝试次数
            CoolDown     time.Duration `json:"cool_down"`// 冷却时间
        } `json:"correction"`
    } `json:"stability"`

    // 优化控制
    Optimization struct {
        Enabled        bool      `json:"enabled"`        // 是否启用
        Strategy       string    `json:"strategy"`       // 优化策略
        Interval       time.Duration `json:"interval"`   // 优化间隔
        
        // 目标参数
        Objectives struct {
            Energy      float64 `json:"energy"`      // 能量目标
            Performance float64 `json:"performance"` // 性能目标
            Stability   float64 `json:"stability"`   // 稳定性目标
            Weights     map[string]float64 `json:"weights"` // 权重
        } `json:"objectives"`
    } `json:"optimization"`
}

// MonitorConfig 监控系统配置
type MonitorConfig struct {
    // 基本监控参数
    Base struct {
        SampleInterval time.Duration `json:"sample_interval"` // 采样间隔
        BatchSize      int          `json:"batch_size"`      // 批处理大小
        BufferSize     int          `json:"buffer_size"`     // 缓冲区大小
        RetentionTime  time.Duration `json:"retention_time"` // 保留时间
    } `json:"base"`

    // 指标配置
    Metrics struct {
        EnabledTypes   []string   `json:"enabled_types"`   // 启用的指标类型
        CustomMetrics  []string   `json:"custom_metrics"`  // 自定义指标
        
        // 聚合配置
        Aggregation struct {
            Interval    time.Duration `json:"interval"`    // 聚合间隔
            Functions   []string      `json:"functions"`   // 聚合函数
            WindowSize  int          `json:"window_size"` // 窗口大小
        } `json:"aggregation"`
    } `json:"metrics"`

    // 告警配置
    Alerts struct {
        Enabled        bool      `json:"enabled"`        // 是否启用
        CheckInterval  time.Duration `json:"check_interval"` // 检查间隔
        MaxAlerts      int       `json:"max_alerts"`     // 最大告警数
        
        // 通知配置
        Notification struct {
            Channels     []string   `json:"channels"`     // 通知渠道
            MinInterval  time.Duration `json:"min_interval"` // 最小间隔
            MaxRetries   int          `json:"max_retries"`  // 最大重试次数
        } `json:"notification"`
    } `json:"alerts"`

    // 健康检查配置
    Health struct {
        CheckInterval  time.Duration `json:"check_interval"` // 检查间隔
        Timeout        time.Duration `json:"timeout"`        // 超时时间
        RetryCount     int          `json:"retry_count"`    // 重试次数
        
        // 检查项配置
        Checks struct {
            Required     []string   `json:"required"`     // 必需检查项
            Optional     []string   `json:"optional"`     // 可选检查项
            Thresholds   map[string]float64 `json:"thresholds"` // 阈值
        } `json:"checks"`
    } `json:"health"`
}

// ResourceConfig 资源系统配置
type ResourceConfig struct {
    // 资源池配置
    Pool struct {
        CPU struct {
            MinAllocation float64 `json:"min_allocation"` // 最小分配
            MaxAllocation float64 `json:"max_allocation"` // 最大分配
            ReserveRatio  float64 `json:"reserve_ratio"`  // 预留比例
        } `json:"cpu"`
        
        Memory struct {
            MinAllocation float64 `json:"min_allocation"` // 最小分配
            MaxAllocation float64 `json:"max_allocation"` // 最大分配
            CacheRatio    float64 `json:"cache_ratio"`    // 缓存比例
        } `json:"memory"`
        
        Energy struct {
            InitialLevel  float64 `json:"initial_level"`  // 初始能级
            MinLevel      float64 `json:"min_level"`      // 最小能级
            MaxLevel      float64 `json:"max_level"`      // 最大能级
            FlowRate      float64 `json:"flow_rate"`      // 流动率
        } `json:"energy"`
    } `json:"pool"`

    // 分配策略
    Allocation struct {
        Strategy       string    `json:"strategy"`       // 分配策略
        QueueSize      int       `json:"queue_size"`     // 队列大小
        Timeout        time.Duration `json:"timeout"`    // 分配超时
        
        // 优先级配置
        Priority struct {
            Levels        int    `json:"levels"`        // 优先级级别
            DefaultLevel  int    `json:"default_level"` // 默认级别
            BoostFactor   float64 `json:"boost_factor"` // 提升因子
        } `json:"priority"`
    } `json:"allocation"`

    // 负载均衡
    Balance struct {
        Enabled        bool      `json:"enabled"`        // 是否启用
        CheckInterval  time.Duration `json:"check_interval"` // 检查间隔
        Threshold      float64   `json:"threshold"`      // 均衡阈值
        
        // 迁移配置
        Migration struct {
            BatchSize     int    `json:"batch_size"`     // 批处理大小
            CoolDown      time.Duration `json:"cool_down"`// 冷却时间
            MaxAttempts   int    `json:"max_attempts"`   // 最大尝试数
        } `json:"migration"`
    } `json:"balance"`

    // 资源监控
    Monitoring struct {
        Enabled        bool      `json:"enabled"`        // 是否启用
        Interval       time.Duration `json:"interval"`   // 监控间隔
        HistorySize    int       `json:"history_size"`  // 历史大小
        
        // 阈值配置
        Thresholds struct {
            CPU         float64 `json:"cpu"`          // CPU阈值
            Memory      float64 `json:"memory"`       // 内存阈值
            Energy      float64 `json:"energy"`       // 能量阈值
            Utilization float64 `json:"utilization"`  // 使用率阈值
        } `json:"thresholds"`
    } `json:"monitoring"`
}
