// system/evolution/manager.go

package evolution

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/common"
	"github.com/Corphon/daoflow/system/control"
	"github.com/Corphon/daoflow/system/evolution/adaptation"
	"github.com/Corphon/daoflow/system/evolution/mutation"
	"github.com/Corphon/daoflow/system/evolution/pattern"
	"github.com/Corphon/daoflow/system/types"
)

// Manager 演化系统管理器
type Manager struct {
	mu sync.RWMutex

	// 基础配置
	config *types.EvoConfig

	// 演化组件
	components struct {
		patternGen  *pattern.PatternGenerator        // 模式生成器
		patternRec  *pattern.PatternRecognizer       // 模式识别器
		evoMatcher  *pattern.EvolutionMatcher        // 演化匹配器
		mutDetector *mutation.MutationDetector       // 突变检测器
		mutHandler  *mutation.MutationHandler        // 突变处理器
		adapLearn   *adaptation.AdaptiveLearning     // 适应性学习
		adapStrat   *adaptation.AdaptationStrategy   // 适应策略
		optimizer   *adaptation.AdaptiveOptimization // 优化器
	}

	// 演化状态
	state struct {
		status    string                 // 运行状态
		startTime time.Time              // 启动时间
		evolution types.EvolutionStatus  // 演化状态
		metrics   map[string]float64     // 演化指标
		history   []types.EvolutionPoint // 演化历史
	}

	// 核心依赖
	core    *core.Engine
	common  *common.Manager
	control *control.Manager

	// 上下文控制
	ctx    context.Context
	cancel context.CancelFunc

	// 观察者列表
	observers []types.StateObserver
}

// NewManager 创建新的管理器实例
func NewManager(cfg *types.EvoConfig) (*Manager, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	m := &Manager{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
	}

	// 初始化状态
	m.state.status = "initialized"
	m.state.startTime = time.Now()
	m.state.metrics = make(map[string]float64)
	m.state.history = make([]types.EvolutionPoint, 0)

	return m, nil
}

// DefaultConfig 返回默认配置
func DefaultConfig() *types.EvoConfig {
	return &types.EvoConfig{
		Pattern: &types.PatternConfig{
			Base: struct {
				GenerationRate float64       `json:"generation_rate"`
				MutationRate   float64       `json:"mutation_rate"`
				ComplexityBias float64       `json:"complexity_bias"`
				EnergyBalance  float64       `json:"energy_balance"`
				UpdateInterval time.Duration `json:"update_interval"`
			}{
				GenerationRate: 0.1,
				MutationRate:   0.01,
				ComplexityBias: 0.5,
				EnergyBalance:  0.7,
				UpdateInterval: time.Second,
			},
			Template: struct {
				MaxTemplates  int     `json:"max_templates"`
				MinSuccess    float64 `json:"min_success"`
				MaxComponents int     `json:"max_components"`
				MaxRelations  int     `json:"max_relations"`
			}{
				MaxTemplates:  100,
				MinSuccess:    0.6,
				MaxComponents: 50,
				MaxRelations:  100,
			},
		},
		Recognition: &types.RecognitionConfig{
			Base: struct {
				MinConfidence  float64       `json:"min_confidence"`
				LearningRate   float64       `json:"learning_rate"`
				MemoryDepth    int           `json:"memory_depth"`
				AdaptiveRate   bool          `json:"adaptive_rate"`
				UpdateInterval time.Duration `json:"update_interval"`
			}{
				MinConfidence:  0.7,
				LearningRate:   0.1,
				MemoryDepth:    100,
				AdaptiveRate:   true,
				UpdateInterval: time.Second,
			},
			Evaluation: struct {
				StructureWeight float64 `json:"structure_weight"`
				DynamicsWeight  float64 `json:"dynamics_weight"`
				ContextWeight   float64 `json:"context_weight"`
				StabilityFactor float64 `json:"stability_factor"`
				TimeDecayFactor float64 `json:"time_decay_factor"`
			}{
				StructureWeight: 0.4,
				DynamicsWeight:  0.3,
				ContextWeight:   0.3,
				StabilityFactor: 0.7,
				TimeDecayFactor: 0.1,
			},
			Memory: struct {
				MaxSize       int           `json:"max_size"`
				RetentionTime time.Duration `json:"retention_time"`
				PruneInterval time.Duration `json:"prune_interval"`
				MinRelevance  float64       `json:"min_relevance"`
			}{
				MaxSize:       1000,
				RetentionTime: time.Hour * 24,
				PruneInterval: time.Hour,
				MinRelevance:  0.3,
			},
		},
		Evolution: &types.EvolutionConfig{
			MatchThreshold: 0.7,
			EvolutionDepth: 5,
			AdaptiveBias:   0.3,
			ContextWeight:  0.5,
			Rules: struct {
				MinConfidence float64 `json:"min_confidence"`
				MaxRules      int     `json:"max_rules"`
				UpdateRate    float64 `json:"update_rate"`
			}{
				MinConfidence: 0.6,
				MaxRules:      100,
				UpdateRate:    0.1,
			},
			Trajectory: struct {
				MaxLength      int           `json:"max_length"`
				MaxAge         time.Duration `json:"max_age"`
				PruneRate      float64       `json:"prune_rate"`
				MinProbability float64       `json:"min_probability"`
			}{
				MaxLength:      1000,
				MaxAge:         time.Hour * 24,
				PruneRate:      0.1,
				MinProbability: 0.3,
			},
		},
		Mutation: &types.MutationConfig{
			Detection: struct {
				Threshold       float64       `json:"threshold"`
				TimeWindow      time.Duration `json:"time_window"`
				Sensitivity     float64       `json:"sensitivity"`
				StabilityFactor float64       `json:"stability_factor"`
				MinSamples      int           `json:"min_samples"`
				MaxSamples      int           `json:"max_samples"`
			}{
				Threshold:       0.7,
				TimeWindow:      time.Minute * 5,
				Sensitivity:     0.8,
				StabilityFactor: 0.6,
				MinSamples:      10,
				MaxSamples:      1000,
			},
			Handler: struct {
				ResponseThreshold float64       `json:"response_threshold"`
				MaxRetries        int           `json:"max_retries"`
				StabilityTarget   float64       `json:"stability_target"`
				AdaptiveResponse  bool          `json:"adaptive_response"`
				ActionTimeout     time.Duration `json:"action_timeout"`
			}{
				ResponseThreshold: 0.8,
				MaxRetries:        3,
				StabilityTarget:   0.7,
				AdaptiveResponse:  true,
				ActionTimeout:     time.Second * 30,
			},
		},
		Adaptation: &types.AdaptationConfig{
			Learning: struct {
				LearningRate    float64       `json:"learning_rate"`
				MemoryCapacity  int           `json:"memory_capacity"`
				ExplorationRate float64       `json:"exploration_rate"`
				DecayFactor     float64       `json:"decay_factor"`
				UpdateInterval  time.Duration `json:"update_interval"`
			}{
				LearningRate:    0.1,
				MemoryCapacity:  1000,
				ExplorationRate: 0.2,
				DecayFactor:     0.95,
				UpdateInterval:  time.Second * 5,
			},
			Pattern: struct {
				MinConfidence float64       `json:"min_confidence"`
				MaxPatterns   int           `json:"max_patterns"`
				PruneInterval time.Duration `json:"prune_interval"`
				RetentionTime time.Duration `json:"retention_time"`
			}{
				MinConfidence: 0.6,
				MaxPatterns:   100,
				PruneInterval: time.Hour,
				RetentionTime: time.Hour * 24,
			},
		},
		Strategy: &types.StrategyConfig{
			Base: struct {
				UpdateInterval    time.Duration `json:"update_interval"`
				MaxStrategies     int           `json:"max_strategies"`
				MinEffectiveness  float64       `json:"min_effectiveness"`
				AdaptiveThreshold float64       `json:"adaptive_threshold"`
			}{
				UpdateInterval:    time.Hour,
				MaxStrategies:     100,
				MinEffectiveness:  0.5,
				AdaptiveThreshold: 0.7,
			},
			Rules: struct {
				MaxRules      int           `json:"max_rules"`
				MinConfidence float64       `json:"min_confidence"`
				UpdateRate    float64       `json:"update_rate"`
				PruneInterval time.Duration `json:"prune_interval"`
			}{
				MaxRules:      100,
				MinConfidence: 0.6,
				UpdateRate:    0.1,
				PruneInterval: time.Hour,
			},
			Execution: struct {
				MaxRetries    int           `json:"max_retries"`
				Timeout       time.Duration `json:"timeout"`
				BatchSize     int           `json:"batch_size"`
				RetryInterval time.Duration `json:"retry_interval"`
			}{
				MaxRetries:    3,
				Timeout:       time.Minute * 5,
				BatchSize:     10,
				RetryInterval: time.Second * 30,
			},
			Evaluation: struct {
				SuccessThreshold float64 `json:"success_threshold"`
				WeightDecay      float64 `json:"weight_decay"`
				HistorySize      int     `json:"history_size"`
				MinSamples       int     `json:"min_samples"`
			}{
				SuccessThreshold: 0.7,
				WeightDecay:      0.1,
				HistorySize:      1000,
				MinSamples:       10,
			},
		},
	}
}

// Start 启动管理器
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.state.status == "running" {
		return nil
	}

	// 初始化并启动所有组件
	if err := m.initComponents(); err != nil {
		return err
	}

	m.state.status = "running"
	m.state.startTime = time.Now()
	return nil
}

// Stop 停止管理器
func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.state.status != "running" {
		return nil
	}

	m.cancel()
	m.state.status = "stopped"
	return nil
}

// Status 获取管理器状态
func (m *Manager) Status() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state.status
}

// Wait 等待管理器停止
func (m *Manager) Wait() {
	<-m.ctx.Done()
}

// GetMetrics 获取管理器指标
func (m *Manager) GetMetrics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"status":       m.state.status,
		"uptime":       time.Since(m.state.startTime).String(),
		"evolution":    m.state.evolution,
		"metrics":      m.state.metrics,
		"history_size": len(m.state.history),
	}
}

// InjectCore 注入核心引擎
func (m *Manager) InjectCore(core *core.Engine) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.core = core
}

// 私有方法

// initComponents 初始化组件
func (m *Manager) initComponents() error {
	// 创建模式生成器
	patternGen, err := pattern.NewPatternGenerator(m.config.Pattern)
	if err != nil {
		return fmt.Errorf("failed to create pattern generator: %w", err)
	}
	m.components.patternGen = patternGen

	// 创建模式识别器
	patternRec, err := pattern.NewPatternRecognizer(m.config.Recognition)
	if err != nil {
		return fmt.Errorf("failed to create pattern recognizer: %w", err)
	}
	m.components.patternRec = patternRec

	// 创建演化匹配器
	evoMatcher, err := pattern.NewEvolutionMatcher(patternRec, m.config.Evolution)
	if err != nil {
		return fmt.Errorf("failed to create evolution matcher: %w", err)
	}
	m.components.evoMatcher = evoMatcher

	// 创建突变检测器
	mutDetector, err := mutation.NewMutationDetector(m.config.Mutation)
	if err != nil {
		return fmt.Errorf("failed to create mutation detector: %w", err)
	}
	m.components.mutDetector = mutDetector

	// 创建突变处理器
	mutHandler, err := mutation.NewMutationHandler(mutDetector, m.config.Mutation)
	if err != nil {
		return fmt.Errorf("failed to create mutation handler: %w", err)
	}
	m.components.mutHandler = mutHandler

	// 创建适应性学习组件
	adapLearn, err := adaptation.NewAdaptiveLearning(evoMatcher, m.config.Adaptation)
	if err != nil {
		return fmt.Errorf("failed to create adaptive learning: %w", err)
	}
	m.components.adapLearn = adapLearn

	// 创建适应策略组件
	adapStrat, err := adaptation.NewAdaptationStrategy(evoMatcher, mutHandler)
	if err != nil {
		return fmt.Errorf("failed to create adaptation strategy: %w", err)
	}
	m.components.adapStrat = adapStrat

	// 创建优化器
	optimizer := adaptation.NewAdaptiveOptimization(adapStrat, adapLearn)
	if optimizer == nil {
		return fmt.Errorf("failed to create optimizer")
	}
	m.components.optimizer = optimizer

	return nil
}

// updateEvolutionStatus 更新演化状态
func (m *Manager) updateEvolutionStatus() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 获取当前演化阶段
	phase := m.calculateEvolutionPhase()

	// 计算演化方向
	direction := m.calculateEvolutionDirection()

	// 计算系统状态
	stability := m.calculateSystemStability()
	energy := m.calculateSystemEnergy()
	progress := m.calculateEvolutionProgress()

	// 更新演化状态
	m.state.evolution = types.EvolutionStatus{
		Phase:     phase,
		Direction: direction,
		Progress:  progress,
		Stability: stability,
		Energy:    energy,
		UpdatedAt: time.Now(),
	}

	// 记录演化点
	point := types.EvolutionPoint{
		State: model.SystemState{
			Phase:      phase,
			Energy:     energy,
			Stability:  stability,
			Properties: m.collectSystemProperties(),
		},
		Energy:    energy,
		Timestamp: time.Now(),
		Meta: map[string]interface{}{
			"direction": direction,
			"progress":  progress,
		},
	}

	// 更新历史记录
	m.state.history = append(m.state.history, point)
	if len(m.state.history) > m.config.MaxHistorySize {
		m.state.history = m.state.history[1:]
	}

	// 更新指标
	m.state.metrics["stability"] = stability
	m.state.metrics["energy"] = energy
	m.state.metrics["progress"] = progress
	m.state.metrics["learning_rate"] = m.components.adapLearn.GetLearningRate()
	m.state.metrics["mutation_rate"] = m.components.mutDetector.GetMutationRate()
}

// 辅助方法
func (m *Manager) calculateEvolutionPhase() model.Phase {
	// 根据系统状态确定演化阶段
	stability := m.calculateSystemStability()
	energy := m.calculateSystemEnergy()

	if stability < 0.3 {
		return model.Phase_Unstable
	} else if stability > 0.8 {
		return model.Phase_Stable
	} else if energy > 0.7 {
		return model.PhaseTransform
	}
	return model.PhaseNeutral
}

// 通过分析历史演化点计算方向
func (m *Manager) calculateEvolutionDirection() model.Vector3D {

	if len(m.state.history) < 2 {
		return model.Vector3D{X: 0, Y: 0, Z: 0}
	}

	latest := m.state.history[len(m.state.history)-1]
	previous := m.state.history[len(m.state.history)-2]

	return model.Vector3D{
		X: latest.Energy - previous.Energy,
		Y: latest.State.Stability - previous.State.Stability,
		Z: calculateSystemComplexity(latest) - calculateSystemComplexity(previous),
	}
}

// calculateSystemComplexity 计算系统复杂度
func calculateSystemComplexity(point types.EvolutionPoint) float64 {
	// 基础复杂度（根据属性数量）
	baseComplexity := float64(len(point.State.Properties)) * 0.1

	// 结构复杂度
	structuralComplexity := 0.0
	if coreState, ok := point.State.Properties["core_state"].(map[string]interface{}); ok {
		structuralComplexity = float64(len(coreState)) * 0.2
	}

	// 模式复杂度
	patternComplexity := 0.0
	if patternCount, ok := point.State.Properties["pattern_count"].(int); ok {
		patternComplexity = float64(patternCount) * 0.3
	}

	// 突变复杂度
	mutationComplexity := 0.0
	if mutationCount, ok := point.State.Properties["mutation_count"].(int); ok {
		mutationComplexity = float64(mutationCount) * 0.2
	}

	// 能量复杂度（基于能量水平的非线性函数）
	energyComplexity := math.Log1p(point.Energy) * 0.2

	// 综合复杂度（归一化到0-1范围）
	totalComplexity := (baseComplexity +
		structuralComplexity +
		patternComplexity +
		mutationComplexity +
		energyComplexity) / 2.0

	return math.Max(0.0, math.Min(1.0, totalComplexity))
}

func (m *Manager) calculateSystemStability() float64 {
	patterns := m.components.patternRec.GetPatterns()

	return calculateAverageStability(patterns)
}

// calculateAverageStability 计算模式平均稳定性
func calculateAverageStability(patterns []*pattern.RecognizedPattern) float64 {
	if len(patterns) == 0 {
		return 1.0 // 没有模式时认为系统稳定
	}

	totalStability := 0.0
	activeCount := 0

	for _, p := range patterns {
		if p.Active {
			totalStability += p.Stability
			activeCount++
		}
	}

	if activeCount == 0 {
		return 1.0
	}

	// 计算平均稳定性并确保在 [0,1] 范围内
	return math.Max(0.0, math.Min(1.0, totalStability/float64(activeCount)))
}

func (m *Manager) calculateSystemEnergy() float64 {
	if m.core == nil {
		return 0
	}
	return m.core.GetTotalEnergy()
}

func (m *Manager) calculateEvolutionProgress() float64 {
	if m.config.Target == nil {
		return 0
	}
	current := m.collectSystemProperties()
	return calculateStateDistance(current, m.config.Target.Properties)
}

// calculateStateDistance 计算系统状态之间的距离
func calculateStateDistance(current, target map[string]interface{}) float64 {
	if len(target) == 0 {
		return 1.0 // 没有目标时返回最大距离
	}

	totalDiff := 0.0
	matchCount := 0

	// 为不同类型的属性定义权重
	weights := map[string]float64{
		"core_state":     0.4, // 核心状态权重
		"pattern_count":  0.3, // 模式数量权重
		"mutation_count": 0.3, // 突变数量权重
	}

	for key, targetValue := range target {
		currentValue, exists := current[key]
		if !exists {
			continue
		}

		// 根据不同类型计算差异
		diff := 0.0
		switch v := targetValue.(type) {
		case float64:
			if cv, ok := currentValue.(float64); ok {
				diff = math.Abs(cv - v)
			}
		case int:
			if cv, ok := currentValue.(int); ok {
				diff = math.Abs(float64(cv - v))
			}
		case map[string]interface{}:
			if cv, ok := currentValue.(map[string]interface{}); ok {
				// 递归计算嵌套属性的差异
				diff = calculateMapDistance(cv, v)
			}
		}

		// 应用权重
		if weight, ok := weights[key]; ok {
			diff *= weight
		}

		totalDiff += diff
		matchCount++
	}

	if matchCount == 0 {
		return 1.0
	}

	// 归一化到 [0,1] 范围
	distance := totalDiff / float64(matchCount)
	return math.Max(0.0, math.Min(1.0, distance))
}

// calculateMapDistance 计算两个map之间的距离
func calculateMapDistance(map1, map2 map[string]interface{}) float64 {
	if len(map1) == 0 || len(map2) == 0 {
		return 1.0
	}

	totalDiff := 0.0
	matchCount := 0

	for key, value1 := range map1 {
		if value2, exists := map2[key]; exists {
			switch v1 := value1.(type) {
			case float64:
				if v2, ok := value2.(float64); ok {
					totalDiff += math.Abs(v1 - v2)
					matchCount++
				}
			case int:
				if v2, ok := value2.(int); ok {
					totalDiff += math.Abs(float64(v1 - v2))
					matchCount++
				}
			}
		}
	}

	if matchCount == 0 {
		return 1.0
	}

	return totalDiff / float64(matchCount)
}

func (m *Manager) collectSystemProperties() map[string]interface{} {
	props := make(map[string]interface{})
	if m.core != nil {
		props["core_state"] = m.core.GetState()
	}
	props["pattern_count"] = len(m.components.patternRec.GetPatterns())
	props["mutation_count"] = m.components.mutDetector.GetMutationCount()
	return props
}

// InjectDependencies 注入组件依赖
func (m *Manager) InjectDependencies(core *core.Engine, common *common.Manager, control *control.Manager) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 注入核心引擎
	if core == nil {
		return model.WrapError(nil, model.ErrCodeDependency, "core engine is nil")
	}
	m.core = core

	// 注入通用管理器
	if common == nil {
		return model.WrapError(nil, model.ErrCodeDependency, "common manager is nil")
	}
	m.common = common

	// 注入控制管理器
	if control == nil {
		return model.WrapError(nil, model.ErrCodeDependency, "control manager is nil")
	}
	m.control = control

	return nil
}

// Restore 恢复系统
func (m *Manager) Restore(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 重置状态
	m.state.evolution = types.EvolutionStatus{}
	m.state.metrics = make(map[string]float64)
	m.state.history = make([]types.EvolutionPoint, 0)

	// 重置组件
	return m.initComponents()
}

// UpdateState 更新演化系统状态
func (m *Manager) UpdateState() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 更新演化状态
	m.updateEvolutionStatus()

	// 触发状态变更事件
	m.notifyStateChange()

	return nil
}

// 添加观察者管理方法
func (m *Manager) RegisterObserver(observer types.StateObserver) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.observers = append(m.observers, observer)
}

func (m *Manager) UnregisterObserver(observer types.StateObserver) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, obs := range m.observers {
		if obs == observer {
			m.observers = append(m.observers[:i], m.observers[i+1:]...)
			break
		}
	}
}

// 添加状态变更通知方法
func (m *Manager) notifyStateChange() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 构建事件
	event := types.SystemEvent{
		Type:      types.EventEvolutionStateChanged,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"phase":     m.state.evolution.Phase,
			"energy":    m.state.evolution.Energy,
			"stability": m.state.evolution.Stability,
			"progress":  m.state.evolution.Progress,
			"metrics":   m.state.metrics,
		},
	}

	// 通知所有观察者
	for _, observer := range m.observers {
		observer.OnStateChange(event)
	}
}

// DetectPattern 检测数据中的模式
func (m *Manager) DetectPattern(data interface{}) (*model.FlowPattern, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.components.patternRec.DetectPattern(data)
}

// AnalyzePattern 分析模式
func (m *Manager) AnalyzePattern(pattern *model.FlowPattern) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.components.patternRec.AnalyzePattern(pattern)
}

// Optimize 执行系统优化
func (m *Manager) Optimize(params types.OptimizationParams) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 验证参数
	if err := m.validateOptimizationParams(params); err != nil {
		return err
	}

	// 创建优化目标
	objectives := m.createOptimizationObjectives(params)

	// 执行优化过程 - 使用Goals.Constraints而不是params.Constraints
	return m.components.optimizer.RunOptimization(objectives, params.Goals.Constraints)
}

// validateOptimizationParams 验证优化参数
func (m *Manager) validateOptimizationParams(params types.OptimizationParams) error {
	// 验证基本参数
	if params.MaxIterations <= 0 {
		params.MaxIterations = 100 // 设置默认值
	}

	// 验证目标
	if len(params.Goals.Targets) == 0 {
		return fmt.Errorf("no optimization targets specified")
	}

	// 验证约束
	for name, constraint := range params.Goals.Constraints {
		if constraint.Max < constraint.Min {
			return fmt.Errorf("invalid constraint for %s: max < min", name)
		}
	}

	return nil
}

// createOptimizationObjectives 创建优化目标
func (m *Manager) createOptimizationObjectives(params types.OptimizationParams) []*adaptation.OptimizationObjective {
	objectives := make([]*adaptation.OptimizationObjective, 0, len(params.Goals.Targets))

	for target, value := range params.Goals.Targets {
		weight := 1.0
		if w, exists := params.Goals.Weights[target]; exists {
			weight = w
		}

		objective := &adaptation.OptimizationObjective{
			ID:       fmt.Sprintf("obj_%s", target),
			Name:     target,
			Type:     "system",
			TargetID: target,
			Target:   value,
			Weight:   weight,
			Evaluator: adaptation.ObjectiveEvaluator{
				Function:  m.createEvaluatorFunction(target),
				Threshold: 0.01,
			},
		}

		objectives = append(objectives, objective)
	}

	return objectives
}

// createEvaluatorFunction 创建评估函数
func (m *Manager) createEvaluatorFunction(target string) func(interface{}) float64 {
	return func(input interface{}) float64 {
		state, ok := input.(*types.SystemState)
		if !ok {
			return 0
		}

		switch target {
		case "performance":
			return state.Properties["performance"].(float64)
		case "stability":
			return state.Stability
		case "energy":
			return state.Energy
		case "harmony":
			return state.Harmony
		default:
			if value, exists := state.Properties[target].(float64); exists {
				return value
			}
			return 0
		}
	}
}
