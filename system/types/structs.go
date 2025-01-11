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

// 其他必要的配置结构...
type (
    MetaConfig struct {
        // 元系统特定配置
    }

    EvoConfig struct {
        // 演化系统特定配置
    }

    ControlConfig struct {
        // 控制系统特定配置
    }

    MonitorConfig struct {
        // 监控系统特定配置
    }

    ResourceConfig struct {
        // 资源系统特定配置
    }
)
