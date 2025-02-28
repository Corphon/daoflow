package types

import (
	"time"

	"github.com/Corphon/daoflow/model"
)

// SystemConfig 系统配置
type SystemConfig struct {
	// 基础配置
	Base struct {
		ID       string // 系统ID
		Name     string // 系统名称
		Version  string // 系统版本
		LogLevel string // 日志级别
	}

	// 子系统配置
	MetaConfig    *MetaConfig    // 元系统配置
	EvoConfig     *EvoConfig     // 演化系统配置
	MonitorConfig *MonitorConfig // 监控系统配置

	// 模型配置 - 使用 model 包的配置
	Model model.ModelConfig

	// 组件配置
	Components struct {
		Meta      MetaConfig     // 元系统配置
		Evolution EvoConfig      // 演化配置
		Control   ControlConfig  // 控制配置
		Monitor   MonitorConfig  // 监控配置
		Resource  ResourceConfig // 资源配置
	}

	// 系统限制
	Limits struct {
		MaxGoroutines int     // 最大协程数
		MaxMemory     int64   // 最大内存使用
		MaxCPU        float64 // 最大CPU使用
	}
}

// 配置结构实现
// MetaConfig 元系统配置
type MetaConfig struct {
	// 场配置
	Field struct {
		InitialStrength float64       `json:"initial_strength"` // 初始场强度
		MinStrength     float64       `json:"min_strength"`     // 最小场强度
		MaxStrength     float64       `json:"max_strength"`     // 最大场强度
		Dimension       int           `json:"dimension"`        // 场维度
		UpdateInterval  time.Duration `json:"update_interval"`  // 场更新间隔

		// 场相互作用配置
		Coupling struct {
			Strength  float64 `json:"strength"`  // 耦合强度
			Range     float64 `json:"range"`     // 作用范围
			Threshold float64 `json:"threshold"` // 耦合阈值
			MaxPairs  int     `json:"max_pairs"` // 最大耦合对数
		} `json:"coupling"`
	} `json:"field"`

	// 量子配置
	Quantum struct {
		InitialState    []complex128  `json:"initial_state"`    // 初始量子态
		Coherence       float64       `json:"coherence"`        // 相干度阈值
		DecoherenceRate float64       `json:"decoherence_rate"` // 退相干率
		MeasureInterval time.Duration `json:"measure_interval"` // 测量间隔

		// 纠缠配置
		Entanglement struct {
			MaxPairs  int           `json:"max_pairs"` // 最大纠缠对数
			Threshold float64       `json:"threshold"` // 纠缠阈值
			Lifetime  time.Duration `json:"lifetime"`  // 纠缠寿命
		} `json:"entanglement"`
	} `json:"quantum"`

	// 涌现配置
	Emergence struct {
		DetectionInterval time.Duration `json:"detection_interval"` // 检测间隔
		MinStrength       float64       `json:"min_strength"`       // 最小强度阈值
		MaxPatterns       int           `json:"max_patterns"`       // 最大模式数

		// 模式配置
		Patterns struct {
			MinLifetime        time.Duration `json:"min_lifetime"`        // 最小生命周期
			StabilityThreshold float64       `json:"stability_threshold"` // 稳定性阈值
			EnergyThreshold    float64       `json:"energy_threshold"`    // 能量阈值
		} `json:"patterns"`
	} `json:"emergence"`

	// 共振配置
	Resonance struct {
		FrequencyRange [2]float64 `json:"frequency_range"` // 频率范围
		MinAmplitude   float64    `json:"min_amplitude"`   // 最小振幅
		PhaseThreshold float64    `json:"phase_threshold"` // 相位阈值

		// 共振条件
		Conditions struct {
			MinCoupling   float64 `json:"min_coupling"`    // 最小耦合强度
			MinCoherence  float64 `json:"min_coherence"`   // 最小相干度
			MaxPhaseShift float64 `json:"max_phase_shift"` // 最大相位偏移
		} `json:"conditions"`
	} `json:"resonance"`
}

// EvoConfig 演化系统配置
type EvoConfig struct {
	// 基本演化参数
	Base struct {
		InitialLevel   float64       `json:"initial_level"`   // 初始演化级别
		MinLevel       float64       `json:"min_level"`       // 最小演化级别
		MaxLevel       float64       `json:"max_level"`       // 最大演化级别
		UpdateInterval time.Duration `json:"update_interval"` // 更新间隔
	} `json:"base"`

	// 能量配置
	Energy struct {
		InitialEnergy   float64 `json:"initial_energy"`   // 初始能量
		MinEnergy       float64 `json:"min_energy"`       // 最小能量
		MaxEnergy       float64 `json:"max_energy"`       // 最大能量
		DissipationRate float64 `json:"dissipation_rate"` // 能量耗散率
	} `json:"energy"`

	// 模式生成配置
	Pattern *PatternConfig `json:"pattern"`

	// 识别配置
	Recognition *RecognitionConfig `json:"recognition"`

	// 演化配置
	Evolution *EvolutionConfig `json:"evolution"`

	// 突变配置
	Mutation *MutationConfig `json:"mutation"`

	// 适应配置
	Adaptation *AdaptationConfig `json:"adaptation"`

	//适应策略配置
	Strategy *StrategyConfig `json:"strategy"`

	// 历史记录配置
	MaxHistorySize int `json:"max_history_size"` // 最大历史记录大小

	// 目标配置
	Target *struct {
		Properties map[string]interface{} `json:"properties"` // 目标属性
	} `json:"target"`
}

// StrategyConfig 适应策略配置
type StrategyConfig struct {
	// 基础配置
	Base struct {
		UpdateInterval    time.Duration `json:"update_interval"`    // 更新间隔
		MaxStrategies     int           `json:"max_strategies"`     // 最大策略数
		MinEffectiveness  float64       `json:"min_effectiveness"`  // 最小有效性
		AdaptiveThreshold float64       `json:"adaptive_threshold"` // 自适应阈值
	} `json:"base"`

	// 规则配置
	Rules struct {
		MaxRules      int           `json:"max_rules"`      // 最大规则数
		MinConfidence float64       `json:"min_confidence"` // 最小置信度
		UpdateRate    float64       `json:"update_rate"`    // 更新频率
		PruneInterval time.Duration `json:"prune_interval"` // 清理间隔
	} `json:"rules"`

	// 执行配置
	Execution struct {
		MaxRetries    int           `json:"max_retries"`    // 最大重试次数
		Timeout       time.Duration `json:"timeout"`        // 执行超时
		BatchSize     int           `json:"batch_size"`     // 批处理大小
		RetryInterval time.Duration `json:"retry_interval"` // 重试间隔
	} `json:"execution"`

	// 评估配置
	Evaluation struct {
		SuccessThreshold float64 `json:"success_threshold"` // 成功阈值
		WeightDecay      float64 `json:"weight_decay"`      // 权重衰减
		HistorySize      int     `json:"history_size"`      // 历史大小
		MinSamples       int     `json:"min_samples"`       // 最小样本数
	} `json:"evaluation"`
}

// ControlConfig 控制系统配置
type ControlConfig struct {
	// 基本控制参数
	Base struct {
		UpdateRate int           `json:"update_rate"` // 更新频率
		MaxLatency time.Duration `json:"max_latency"` // 最大延迟
		BufferSize int           `json:"buffer_size"` // 缓冲区大小
		Timeout    time.Duration `json:"timeout"`     // 超时时间
	} `json:"base"`

	// 反馈控制
	Feedback struct {
		Enabled      bool          `json:"enabled"`       // 是否启用
		Sensitivity  float64       `json:"sensitivity"`   // 灵敏度
		ResponseTime time.Duration `json:"response_time"` // 响应时间

		// PID控制器参数
		PID struct {
			Proportional float64 `json:"proportional"` // 比例系数
			Integral     float64 `json:"integral"`     // 积分系数
			Derivative   float64 `json:"derivative"`   // 微分系数
			WindupGuard  float64 `json:"windup_guard"` // 积分限幅
		} `json:"pid"`
	} `json:"feedback"`

	// 稳定性控制
	Stability struct {
		CheckInterval time.Duration `json:"check_interval"` // 检查间隔
		MinThreshold  float64       `json:"min_threshold"`  // 最小阈值
		MaxDeviation  float64       `json:"max_deviation"`  // 最大偏差

		// 修正参数
		Correction struct {
			Strength    float64       `json:"strength"`     // 修正强度
			MaxAttempts int           `json:"max_attempts"` // 最大尝试次数
			CoolDown    time.Duration `json:"cool_down"`    // 冷却时间
		} `json:"correction"`
	} `json:"stability"`

	// 优化控制
	Optimization struct {
		Enabled  bool          `json:"enabled"`  // 是否启用
		Strategy string        `json:"strategy"` // 优化策略
		Interval time.Duration `json:"interval"` // 优化间隔

		// 目标参数
		Objectives struct {
			Energy      float64            `json:"energy"`      // 能量目标
			Performance float64            `json:"performance"` // 性能目标
			Stability   float64            `json:"stability"`   // 稳定性目标
			Weights     map[string]float64 `json:"weights"`     // 权重
		} `json:"objectives"`
	} `json:"optimization"`

	// 基础任务配置
	Tasks struct {
		MaxTasks     int           `json:"max_tasks"`     // 最大任务数
		MaxWorkflows int           `json:"max_workflows"` // 最大工作流数
		TaskTimeout  time.Duration `json:"task_timeout"`  // 任务超时时间
		RetryLimit   int           `json:"retry_limit"`   // 重试次数限制
		BufferSize   int           `json:"buffer_size"`   // 缓冲区大小
		UpdateRate   time.Duration `json:"update_rate"`   // 更新频率
	} `json:"tasks"`

	// 流控制配置
	Flow struct {
		MaxConcurrent int           `json:"max_concurrent"` // 最大并发数
		QueueSize     int           `json:"queue_size"`     // 队列大小
		BatchSize     int           `json:"batch_size"`     // 批处理大小
		FlowTimeout   time.Duration `json:"flow_timeout"`   // 流超时时间
	} `json:"flow"`

	// 资源控制配置
	Resource struct {
		MaxAlloc     float64 `json:"max_alloc"`     // 最大分配量
		MinAlloc     float64 `json:"min_alloc"`     // 最小分配量
		ReserveRatio float64 `json:"reserve_ratio"` // 预留比例
	} `json:"resource"`
}

// MonitorConfig 监控系统配置
type MonitorConfig struct {
	// 基本监控参数
	Base struct {
		SampleInterval time.Duration `json:"sample_interval"` // 采样间隔
		BatchSize      int           `json:"batch_size"`      // 批处理大小
		BufferSize     int           `json:"buffer_size"`     // 缓冲区大小
		RetentionTime  time.Duration `json:"retention_time"`  // 保留时间
		MaxHistory     int           `json:"max_history"`     // 最大历史记录字段
	} `json:"base"`

	// 指标配置
	Metrics struct {
		// 基础设置
		Enabled     bool          `json:"enabled"`      // 是否启用
		Interval    time.Duration `json:"interval"`     // 监控间隔
		HistorySize int           `json:"history_size"` // 历史大小

		// 指标类型
		EnabledTypes  []string `json:"enabled_types"`  // 启用的指标类型
		CustomMetrics []string `json:"custom_metrics"` // 自定义指标

		// 聚合配置
		Aggregation struct {
			Interval   time.Duration `json:"interval"`    // 聚合间隔
			Functions  []string      `json:"functions"`   // 聚合函数
			WindowSize int           `json:"window_size"` // 窗口大小
		} `json:"aggregation"`
	} `json:"metrics"`

	// 添加报告配置
	Report struct {
		ReportInterval time.Duration      `json:"report_interval"` // 报告间隔
		ReportFormat   string             `json:"report_format"`   // 报告格式
		OutputPath     string             `json:"output_path"`     // 输出路径
		Thresholds     map[string]float64 `json:"thresholds"`      // 报告阈值
		Filters        []string           `json:"filters"`         // 指标过滤器
	} `json:"report"`

	// 告警配置
	Alert AlertConfig `json:"alert"`

	// 健康检查配置
	Health struct {
		CheckInterval time.Duration `json:"check_interval"` // 检查间隔
		Timeout       time.Duration `json:"timeout"`        // 超时时间
		RetryCount    int           `json:"retry_count"`    // 重试次数

		// 检查项配置
		Checks struct {
			Required   []string           `json:"required"`   // 必需检查项
			Optional   []string           `json:"optional"`   // 可选检查项
			Thresholds map[string]float64 `json:"thresholds"` // 阈值
		} `json:"checks"`
	} `json:"health"`

	// 追踪配置
	Trace struct {
		Enabled       bool          `json:"enabled"`        // 是否启用
		SampleRate    float64       `json:"sample_rate"`    // 采样率
		BufferSize    int           `json:"buffer_size"`    // 缓冲区大小
		MaxSpans      int           `json:"max_spans"`      // 最大跨度数
		FlushInterval time.Duration `json:"flush_interval"` // 刷新间隔
		StoragePath   string        `json:"storage_path"`   // 存储路径

		// 过滤器配置
		Filters struct {
			MinDuration time.Duration `json:"min_duration"` // 最小持续时间
			MaxDuration time.Duration `json:"max_duration"` // 最大持续时间
			Types       []string      `json:"types"`        // 跟踪类型
			Tags        []string      `json:"tags"`         // 标签过滤
		} `json:"filters"`
	} `json:"trace"`
}

// ResourceConfig 资源系统配置
/*type ResourceConfig struct {
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
			InitialLevel float64 `json:"initial_level"` // 初始能级
			MinLevel     float64 `json:"min_level"`     // 最小能级
			MaxLevel     float64 `json:"max_level"`     // 最大能级
			FlowRate     float64 `json:"flow_rate"`     // 流动率
		} `json:"energy"`
	} `json:"pool"`

	// 分配策略
	Allocation struct {
		Strategy  string        `json:"strategy"`   // 分配策略
		QueueSize int           `json:"queue_size"` // 队列大小
		Timeout   time.Duration `json:"timeout"`    // 分配超时

		// 优先级配置
		Priority struct {
			Levels       int     `json:"levels"`        // 优先级级别
			DefaultLevel int     `json:"default_level"` // 默认级别
			BoostFactor  float64 `json:"boost_factor"`  // 提升因子
		} `json:"priority"`
	} `json:"allocation"`

	// 负载均衡
	Balance struct {
		Enabled       bool          `json:"enabled"`        // 是否启用
		CheckInterval time.Duration `json:"check_interval"` // 检查间隔
		Threshold     float64       `json:"threshold"`      // 均衡阈值

		// 迁移配置
		Migration struct {
			BatchSize   int           `json:"batch_size"`   // 批处理大小
			CoolDown    time.Duration `json:"cool_down"`    // 冷却时间
			MaxAttempts int           `json:"max_attempts"` // 最大尝试数
		} `json:"migration"`
	} `json:"balance"`

	// 资源监控
	Monitoring struct {
		Enabled     bool          `json:"enabled"`      // 是否启用
		Interval    time.Duration `json:"interval"`     // 监控间隔
		HistorySize int           `json:"history_size"` // 历史大小

		// 阈值配置
		Thresholds struct {
			CPU         float64 `json:"cpu"`         // CPU阈值
			Memory      float64 `json:"memory"`      // 内存阈值
			Energy      float64 `json:"energy"`      // 能量阈值
			Utilization float64 `json:"utilization"` // 使用率阈值
		} `json:"thresholds"`
	} `json:"monitoring"`
}*/

// CommonConfig 公共系统配置
type CommonConfig struct {
	// 基础配置
	Base struct {
		UpdateInterval time.Duration `json:"update_interval"` // 更新间隔
		MaxRetries     int           `json:"max_retries"`     // 最大重试次数
		Timeout        time.Duration `json:"timeout"`         // 超时时间
	} `json:"base"`

	// 资源限制
	Resources struct {
		MaxFields    int     `json:"max_fields"`    // 最大场数量
		MaxStates    int     `json:"max_states"`    // 最大状态数
		MaxPatterns  int     `json:"max_patterns"`  // 最大模式数
		MaxEnergy    float64 `json:"max_energy"`    // 最大能量值
		ReserveRatio float64 `json:"reserve_ratio"` // 预留比例
	} `json:"resources"`

	// 共享配置
	Sharing struct {
		EnableCache  bool          `json:"enable_cache"`  // 启用缓存
		CacheTTL     time.Duration `json:"cache_ttl"`     // 缓存过期时间
		SyncInterval time.Duration `json:"sync_interval"` // 同步间隔
	} `json:"sharing"`

	// 监控配置
	Monitor struct {
		EnableMetrics bool    `json:"enable_metrics"` // 启用指标
		SampleRate    float64 `json:"sample_rate"`    // 采样率
		ErrorRate     float64 `json:"error_rate"`     // 错误率阈值
	} `json:"monitor"`
}

// PatternConfig 模式生成配置
type PatternConfig struct {
	// 基础生成参数
	Base struct {
		GenerationRate float64       `json:"generation_rate"` // 生成率
		MutationRate   float64       `json:"mutation_rate"`   // 变异率
		ComplexityBias float64       `json:"complexity_bias"` // 复杂度偏好
		EnergyBalance  float64       `json:"energy_balance"`  // 能量平衡因子
		UpdateInterval time.Duration `json:"update_interval"` // 更新间隔
	} `json:"base"`

	// 模板配置
	Template struct {
		MaxTemplates  int     `json:"max_templates"`  // 最大模板数
		MinSuccess    float64 `json:"min_success"`    // 最小成功率
		MaxComponents int     `json:"max_components"` // 最大组件数
		MaxRelations  int     `json:"max_relations"`  // 最大关系数
	} `json:"template"`

	// 生成约束
	Constraints struct {
		MinEnergy     float64 `json:"min_energy"`     // 最小能量
		MaxEnergy     float64 `json:"max_energy"`     // 最大能量
		MinComplexity float64 `json:"min_complexity"` // 最小复杂度
		MaxComplexity float64 `json:"max_complexity"` // 最大复杂度
	} `json:"constraints"`

	// 评估配置
	Evaluation struct {
		BaseWeight       float64 `json:"base_weight"`       // 基础权重
		ComplexityWeight float64 `json:"complexity_weight"` // 复杂度权重
		EnergyWeight     float64 `json:"energy_weight"`     // 能量权重
		SuccessThreshold float64 `json:"success_threshold"` // 成功阈值
	} `json:"evaluation"`

	// 优化配置
	Optimization struct {
		Enabled       bool    `json:"enabled"`        // 是否启用优化
		MaxIterations int     `json:"max_iterations"` // 最大迭代次数
		StopThreshold float64 `json:"stop_threshold"` // 停止阈值
		LearningRate  float64 `json:"learning_rate"`  // 学习率
	} `json:"optimization"`
}

// EvolutionConfig 演化配置
type EvolutionConfig struct {
	// 匹配配置
	MatchThreshold float64 `json:"match_threshold"` // 匹配阈值
	EvolutionDepth int     `json:"evolution_depth"` // 演化深度
	AdaptiveBias   float64 `json:"adaptive_bias"`   // 自适应偏差
	ContextWeight  float64 `json:"context_weight"`  // 上下文权重

	// 演化规则
	Rules struct {
		MinConfidence float64 `json:"min_confidence"` // 最小置信度
		MaxRules      int     `json:"max_rules"`      // 最大规则数
		UpdateRate    float64 `json:"update_rate"`    // 规则更新率
	} `json:"rules"`

	// 轨迹配置
	Trajectory struct {
		MaxLength      int           `json:"max_length"`      // 最大长度
		MaxAge         time.Duration `json:"max_age"`         // 最大保留时间
		PruneRate      float64       `json:"prune_rate"`      // 裁剪率
		MinProbability float64       `json:"min_probability"` // 最小概率
	} `json:"trajectory"`

	// 预测配置
	Prediction struct {
		Horizon     time.Duration `json:"horizon"`      // 预测周期
		StepSize    time.Duration `json:"step_size"`    // 步长
		Confidence  float64       `json:"confidence"`   // 置信度阈值
		MaxBranches int           `json:"max_branches"` // 最大分支数
	} `json:"prediction"`

	// 上下文配置
	Context struct {
		UpdateInterval time.Duration `json:"update_interval"` // 更新间隔
		HistoryLength  int           `json:"history_length"`  // 历史长度
		DecayRate      float64       `json:"decay_rate"`      // 衰减率
		MinInfluence   float64       `json:"min_influence"`   // 最小影响度
	} `json:"context"`
}

// RecognitionConfig 模式识别配置
type RecognitionConfig struct {
	// 基础识别参数
	Base struct {
		MinConfidence  float64       `json:"min_confidence"`  // 最小置信度
		LearningRate   float64       `json:"learning_rate"`   // 学习率
		MemoryDepth    int           `json:"memory_depth"`    // 记忆深度
		AdaptiveRate   bool          `json:"adaptive_rate"`   // 是否自适应学习率
		UpdateInterval time.Duration `json:"update_interval"` // 更新间隔
	} `json:"base"`

	// 模式评估配置
	Evaluation struct {
		StructureWeight float64 `json:"structure_weight"`  // 结构权重
		DynamicsWeight  float64 `json:"dynamics_weight"`   // 动态权重
		ContextWeight   float64 `json:"context_weight"`    // 上下文权重
		StabilityFactor float64 `json:"stability_factor"`  // 稳定性因子
		TimeDecayFactor float64 `json:"time_decay_factor"` // 时间衰减因子
	} `json:"evaluation"`

	// 记忆管理配置
	Memory struct {
		MaxSize       int           `json:"max_size"`       // 最大记忆数量
		RetentionTime time.Duration `json:"retention_time"` // 保留时间
		PruneInterval time.Duration `json:"prune_interval"` // 清理间隔
		MinRelevance  float64       `json:"min_relevance"`  // 最小关联度
	} `json:"memory"`

	// 关联分析配置
	Association struct {
		MaxAssociations int           `json:"max_associations"` // 最大关联数
		MinSimilarity   float64       `json:"min_similarity"`   // 最小相似度
		DepthLimit      int           `json:"depth_limit"`      // 深度限制
		TimeWindow      time.Duration `json:"time_window"`      // 时间窗口
	} `json:"association"`

	// 统计配置
	Statistics struct {
		EnableMetrics  bool    `json:"enable_metrics"`  // 启用指标
		SampleInterval int     `json:"sample_interval"` // 采样间隔
		MaxHistory     int     `json:"max_history"`     // 最大历史记录
		MinAccuracy    float64 `json:"min_accuracy"`    // 最小准确度
	} `json:"statistics"`
}

// AdaptationConfig 适应性配置
type AdaptationConfig struct {
	// 学习配置
	Learning struct {
		LearningRate    float64       `json:"learning_rate"`    // 学习率
		MemoryCapacity  int           `json:"memory_capacity"`  // 记忆容量
		ExplorationRate float64       `json:"exploration_rate"` // 探索率
		DecayFactor     float64       `json:"decay_factor"`     // 衰减因子
		UpdateInterval  time.Duration `json:"update_interval"`  // 更新间隔
	} `json:"learning"`

	// 模式配置
	Pattern struct {
		MinConfidence float64       `json:"min_confidence"` // 最小置信度
		MaxPatterns   int           `json:"max_patterns"`   // 最大模式数
		PruneInterval time.Duration `json:"prune_interval"` // 裁剪间隔
		RetentionTime time.Duration `json:"retention_time"` // 保留时间
	} `json:"pattern"`

	// 模型配置
	Model struct {
		BatchSize       int     `json:"batch_size"`       // 批次大小
		MaxIterations   int     `json:"max_iterations"`   // 最大迭代次数
		MinAccuracy     float64 `json:"min_accuracy"`     // 最小准确率
		ValidationRatio float64 `json:"validation_ratio"` // 验证比例
	} `json:"model"`

	// 知识配置
	Knowledge struct {
		MaxSize       int           `json:"max_size"`       // 最大知识数
		MinConfidence float64       `json:"min_confidence"` // 最小置信度
		UpdateRate    float64       `json:"update_rate"`    // 更新频率
		ExpireTime    time.Duration `json:"expire_time"`    // 过期时间
	} `json:"knowledge"`

	// 策略配置
	Strategy struct {
		MaxRules          int           `json:"max_rules"`          // 最大规则数
		MinEffectiveness  float64       `json:"min_effectiveness"`  // 最小有效性
		OptimizeInterval  time.Duration `json:"optimize_interval"`  // 优化间隔
		AdaptiveThreshold float64       `json:"adaptive_threshold"` // 适应阈值
	} `json:"strategy"`
}

// MutationConfig 突变系统配置
type MutationConfig struct {
	// 检测配置
	Detection struct {
		Threshold       float64       `json:"threshold"`        // 检测阈值
		TimeWindow      time.Duration `json:"time_window"`      // 时间窗口
		Sensitivity     float64       `json:"sensitivity"`      // 灵敏度
		StabilityFactor float64       `json:"stability_factor"` // 稳定性因子
		MinSamples      int           `json:"min_samples"`      // 最小样本数
		MaxSamples      int           `json:"max_samples"`      // 最大样本数
	} `json:"detection"`

	// 处理器配置
	Handler struct {
		ResponseThreshold float64       `json:"response_threshold"` // 响应阈值
		MaxRetries        int           `json:"max_retries"`        // 最大重试次数
		StabilityTarget   float64       `json:"stability_target"`   // 稳定性目标
		AdaptiveResponse  bool          `json:"adaptive_response"`  // 自适应响应
		ActionTimeout     time.Duration `json:"action_timeout"`     // 动作超时
	} `json:"handler"`

	// 基准线配置
	Baseline struct {
		UpdateInterval time.Duration `json:"update_interval"` // 更新间隔
		MinConfidence  float64       `json:"min_confidence"`  // 最小置信度
		MaxHistory     int           `json:"max_history"`     // 最大历史记录
		PruneInterval  time.Duration `json:"prune_interval"`  // 清理间隔
	} `json:"baseline"`

	// 性能配置
	Performance struct {
		MaxConcurrent  int           `json:"max_concurrent"`   // 最大并发数
		BufferSize     int           `json:"buffer_size"`      // 缓冲区大小
		MaxQueueLength int           `json:"max_queue_length"` // 最大队列长度
		BatchTimeout   time.Duration `json:"batch_timeout"`    // 批处理超时
	} `json:"performance"`

	// 指标配置
	Metrics struct {
		EnableHistory bool          `json:"enable_history"` // 启用历史记录
		SampleRate    float64       `json:"sample_rate"`    // 采样率
		WindowSize    time.Duration `json:"window_size"`    // 窗口大小
		FlushInterval time.Duration `json:"flush_interval"` // 刷新间隔
	} `json:"metrics"`
}

// CPUConfig CPU资源配置
type CPUConfig struct {
	MinAllocation float64 `json:"min_allocation"`
	MaxAllocation float64 `json:"max_allocation"`
	ReserveRatio  float64 `json:"reserve_ratio"`
}

// MemoryConfig 内存资源配置
type MemoryConfig struct {
	MinAllocation float64 `json:"min_allocation"`
	MaxAllocation float64 `json:"max_allocation"`
	CacheRatio    float64 `json:"cache_ratio"`
}

// EnergyConfig 能量资源配置
type EnergyConfig struct {
	InitialLevel float64 `json:"initial_level"`
	MinLevel     float64 `json:"min_level"`
	MaxLevel     float64 `json:"max_level"`
	FlowRate     float64 `json:"flow_rate"`
}

// ResourcePoolConfig 资源池配置
type ResourcePoolConfig struct {
	CPU    CPUConfig    `json:"cpu"`
	Memory MemoryConfig `json:"memory"`
	Energy EnergyConfig `json:"energy"`
}

// ResourceConfig 资源配置
type ResourceConfig struct {
	Pool ResourcePoolConfig `json:"pool"`
}

// -------------------------------------------
// DefaultResourceConfig 返回默认资源配置
func DefaultResourceConfig() *ResourceConfig {
	return &ResourceConfig{
		Pool: ResourcePoolConfig{
			CPU: CPUConfig{
				MinAllocation: 0.1, // 10%
				MaxAllocation: 0.9, // 90%
				ReserveRatio:  0.2, // 20%
			},
			Memory: MemoryConfig{
				MinAllocation: 0.1, // 10%
				MaxAllocation: 0.8, // 80%
				CacheRatio:    0.3, // 30%
			},
			Energy: EnergyConfig{
				InitialLevel: 1000.0,
				MinLevel:     100.0,
				MaxLevel:     10000.0,
				FlowRate:     0.1,
			},
		},
	}
}
