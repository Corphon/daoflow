// system/system.go

package system

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
	"github.com/Corphon/daoflow/system/evolution"
	"github.com/Corphon/daoflow/system/meta"
	"github.com/Corphon/daoflow/system/monitor"
	"github.com/Corphon/daoflow/system/types"
)

// System represents the main system controller that coordinates all components
type System struct {
	mu sync.RWMutex

	// Core components
	core *core.Engine

	// Model components
	models       map[string]model.Model
	modelManager *model.IntegrateFlow // 集成流模型管理器

	// System subsystems
	common    *common.Manager    // Common utilities and shared resources
	control   *control.Manager   // System control and management
	evolution *evolution.Manager // Evolution and learning capabilities
	meta      *meta.Manager      // Metadata and system information
	monitor   *monitor.Manager   // System monitoring and metrics

	// System state management
	state struct {
		status    string              // 系统状态
		startTime time.Time           // 启动时间
		errors    []error             // 错误记录
		metrics   types.SystemMetrics // 系统指标
		events    []types.SystemEvent // 事件历史
		energy    float64             // 系统能量
	}

	// Event handling
	events struct {
		handlers  map[types.EventType][]types.EventHandler // 事件处理器
		queue     chan types.SystemEvent                   // 事件队列
		processor types.EventProcessor                     // 事件处理器
	}

	// Lifecycle management
	isRunning bool
	ctx       context.Context
	cancel    context.CancelFunc

	// Configuration
	config *Config
}

// Config holds the system configuration
type Config struct {
	CoreConfig      *core.Config
	ModelConfig     *model.ModelConfig
	CommonConfig    *types.CommonConfig
	ControlConfig   *types.ControlConfig
	EvolutionConfig *types.EvoConfig
	MetaConfig      *types.MetaConfig
	MonitorConfig   *types.MonitorConfig
}

// --------------------------------------
// New creates a new System instance
func New(cfg *Config) (*System, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	sys := &System{
		models: make(map[string]model.Model),
		ctx:    ctx,
		cancel: cancel,
		config: cfg,
	}

	// 初始化事件系统
	sys.events.handlers = make(map[types.EventType][]types.EventHandler)
	sys.events.queue = make(chan types.SystemEvent, 1000)
	sys.events.processor = types.NewEventBus()

	// 初始化状态
	sys.state.status = "initialized"
	sys.state.startTime = time.Now()
	sys.state.errors = make([]error, 0)
	sys.state.events = make([]types.SystemEvent, 0)
	sys.state.metrics = types.SystemMetrics{}

	// 初始化模型管理器
	integrateFlow := model.NewIntegrateFlow()
	sys.modelManager = integrateFlow
	sys.models["integrate"] = integrateFlow

	// Initialize core engine
	engine, err := core.NewEngine(cfg.CoreConfig)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize core engine: %w", err)
	}
	sys.core = engine

	// Initialize subsystems
	if err := sys.initializeSubsystems(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize subsystems: %w", err)
	}

	// 启动事件处理
	go sys.processEvents()

	return sys, nil
}

// defaultConfig returns default system configuration
func DefaultConfig() *Config {
	return &Config{
		CoreConfig:      core.DefaultConfig(),
		ModelConfig:     model.DefaultConfig(),
		CommonConfig:    common.DefaultConfig(),
		ControlConfig:   control.DefaultConfig(),
		EvolutionConfig: evolution.DefaultConfig(),
		MetaConfig:      meta.DefaultConfig(),
		MonitorConfig:   monitor.DefaultConfig(),
	}
}

// 初始化Config方法
func (c *Config) DefaultConfig() *Config {
	if c == nil {
		return DefaultConfig()
	}

	// 使用现有配置,未设置的使用默认值
	cfg := DefaultConfig()

	if c.CoreConfig != nil {
		cfg.CoreConfig = c.CoreConfig
	}
	if c.ModelConfig != nil {
		cfg.ModelConfig = c.ModelConfig
	}
	if c.CommonConfig != nil {
		cfg.CommonConfig = c.CommonConfig
	}
	if c.ControlConfig != nil {
		cfg.ControlConfig = c.ControlConfig
	}
	if c.EvolutionConfig != nil {
		cfg.EvolutionConfig = c.EvolutionConfig
	}
	if c.MetaConfig != nil {
		cfg.MetaConfig = c.MetaConfig
	}
	if c.MonitorConfig != nil {
		cfg.MonitorConfig = c.MonitorConfig
	}

	return cfg
}

// initializeSubsystems initializes all system subsystems
func (s *System) initializeSubsystems() error {
	var err error

	// Initialize common manager
	s.common, err = common.NewManager(s.config.CommonConfig)
	if err != nil {
		return err
	}

	// Initialize control manager
	s.control, err = control.NewManager(s.config.ControlConfig)
	if err != nil {
		return err
	}

	// Initialize evolution manager
	s.evolution, err = evolution.NewManager(s.config.EvolutionConfig)
	if err != nil {
		return err
	}

	// Initialize meta manager
	s.meta, err = meta.NewManager(s.config.MetaConfig)
	if err != nil {
		return err
	}

	// Initialize monitor manager
	s.monitor, err = monitor.NewManager(s.config.MonitorConfig)
	if err != nil {
		return err
	}

	return nil
}

// Initialize 初始化系统
func (s *System) Initialize(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return types.ErrAlreadyRunning
	}

	// 初始化上下文
	s.ctx = ctx
	s.state.startTime = time.Now()
	s.state.status = "initializing"

	// 初始化各组件
	if err := s.initializeSubsystems(); err != nil {
		return fmt.Errorf("failed to initialize subsystems: %w", err)
	}

	// 注入依赖关系
	if err := s.injectDependencies(); err != nil {
		return fmt.Errorf("failed to inject dependencies: %w", err)
	}

	s.state.status = "initialized"
	return nil
}

// Start 启动系统
func (s *System) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return types.ErrAlreadyRunning
	}

	s.state.status = "starting"

	// 启动所有组件
	if err := s.startComponents(); err != nil {
		s.state.status = "failed"
		return fmt.Errorf("failed to start components: %w", err)
	}

	// 更新系统状态
	s.isRunning = true
	s.state.status = "running"

	// 发送系统启动事件
	s.HandleEvent(types.SystemEvent{
		Type:      types.EventSystemStarted,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"startup_time": time.Since(s.state.startTime).String(),
		},
	})

	return nil
}

// startComponents 启动所有组件
func (s *System) startComponents() error {
	// 1. 启动核心引擎
	if err := s.core.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize core engine: %w", err)
	}

	// 2. 启动子系统,按依赖顺序启动
	if err := s.startSubsystems(); err != nil {
		return fmt.Errorf("failed to start subsystems: %w", err)
	}

	// 3. 启动所有模型
	for name, m := range s.models {
		if err := m.Start(); err != nil {
			s.stopSubsystems()
			return fmt.Errorf("failed to start model %s: %w", name, err)
		}
	}

	return nil
}

// startSubsystems starts all subsystems in dependency order
func (s *System) startSubsystems() error {
	// 1. 启动公共子系统
	if err := s.common.Start(s.ctx); err != nil {
		return fmt.Errorf("failed to start common subsystem: %w", err)
	}

	// 2. 启动控制子系统
	if err := s.control.Start(s.ctx); err != nil {
		s.common.Stop()
		return fmt.Errorf("failed to start control subsystem: %w", err)
	}

	// 3. 启动演化子系统
	if err := s.evolution.Start(s.ctx); err != nil {
		s.control.Stop()
		s.common.Stop()
		return fmt.Errorf("failed to start evolution subsystem: %w", err)
	}

	// 4. 启动元数据子系统
	if err := s.meta.Start(s.ctx); err != nil {
		s.evolution.Stop()
		s.control.Stop()
		s.common.Stop()
		return fmt.Errorf("failed to start meta subsystem: %w", err)
	}

	// 5. 启动监控子系统
	if err := s.monitor.Start(s.ctx); err != nil {
		s.meta.Stop()
		s.evolution.Stop()
		s.control.Stop()
		s.common.Stop()
		return fmt.Errorf("failed to start monitor subsystem: %w", err)
	}

	return nil
}

// Stop 停止系统
func (s *System) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return nil
	}

	s.state.status = "stopping"

	// 发送系统停止事件
	s.HandleEvent(types.SystemEvent{
		Type:      types.EventSystemStopping,
		Timestamp: time.Now(),
	})

	// 关闭所有组件
	if err := s.stopComponents(); err != nil {
		s.recordError(fmt.Errorf("failed to stop components: %w", err))
	}

	s.isRunning = false
	s.state.status = "stopped"

	return nil
}

// stopComponents 停止所有组件
func (s *System) stopComponents() error {
	// 1. 停止所有模型
	for name, m := range s.models {
		if err := m.Stop(); err != nil {
			s.recordError(fmt.Errorf("failed to stop model %s: %w", name, err))
		}
	}

	// 2. 停止所有子系统
	if err := s.stopSubsystems(); err != nil {
		s.recordError(fmt.Errorf("failed to stop subsystems: %w", err))
	}

	// 3. 关闭核心引擎
	if err := s.core.Shutdown(); err != nil {
		s.recordError(fmt.Errorf("failed to stop core engine: %w", err))
	}

	return nil
}

// Shutdown 关闭系统
func (s *System) Shutdown(ctx context.Context) error {
	// 设置关闭超时
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 停止系统
	if err := s.Stop(); err != nil {
		return err
	}

	// 等待所有组件完全停止或超时
	select {
	case <-shutdownCtx.Done():
		if shutdownCtx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("system shutdown timed out")
		}
	case <-s.waitForComponents():
		// 所有组件已停止
	}

	return nil
}

// waitForComponents waits for all components to stop
func (s *System) waitForComponents() chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)

		// 等待监控系统停止
		s.monitor.Wait()

		// 等待元系统停止
		s.meta.Wait()

		// 等待演化系统停止
		s.evolution.Wait()

		// 等待控制系统停止
		s.control.Wait()

		// 等待公共系统停止
		s.common.Wait()
	}()
	return done
}

// stopSubsystems stops all subsystems in reverse order
func (s *System) stopSubsystems() error {
	if err := s.monitor.Stop(); err != nil {
		return err
	}

	if err := s.meta.Stop(); err != nil {
		return err
	}

	if err := s.evolution.Stop(); err != nil {
		return err
	}

	if err := s.control.Stop(); err != nil {
		return err
	}

	if err := s.common.Stop(); err != nil {
		return err
	}

	return nil
}

// Reset resets the system to its initial state
func (s *System) Reset() error {
	if s.isRunning {
		if err := s.Stop(); err != nil {
			return fmt.Errorf("failed to stop system: %w", err)
		}
	}

	// 重置所有状态
	s.state.status = "resetting"
	s.state.startTime = time.Now()
	s.state.errors = make([]error, 0)
	s.state.events = make([]types.SystemEvent, 0)
	s.state.metrics = types.SystemMetrics{}

	// 重置事件系统
	s.events.handlers = make(map[types.EventType][]types.EventHandler)
	s.events.queue = make(chan types.SystemEvent, 1000)
	s.events.processor = types.NewEventBus()

	// 重置上下文
	s.ctx, s.cancel = context.WithCancel(context.Background())

	// 重新初始化所有组件
	if err := s.initializeSubsystems(); err != nil {
		return fmt.Errorf("failed to reinitialize subsystems: %w", err)
	}

	s.state.status = "reset"
	return nil
}

// Subsystem access methods

// Common returns the common utilities manager
func (s *System) Common() *common.Manager {
	return s.common
}

// Control returns the system control manager
func (s *System) Control() *control.Manager {
	return s.control
}

// Evolution returns the evolution manager
func (s *System) Evolution() *evolution.Manager {
	return s.evolution
}

// Meta returns the metadata manager
func (s *System) Meta() *meta.Manager {
	return s.meta
}

// Monitor returns the system monitor manager
func (s *System) Monitor() *monitor.Manager {
	return s.monitor
}

// Model management methods

// RegisterModel adds a new model to the system
func (s *System) RegisterModel(name string, m model.Model) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.models[name]; exists {
		return types.ErrModelAlreadyExists
	}

	s.models[name] = m
	return nil
}

// RegisterModels registers multiple models at once
func (s *System) RegisterModels(models map[string]model.Model) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 预检查
	for name := range models {
		if _, exists := s.models[name]; exists {
			return fmt.Errorf("model %s already exists", name)
		}
	}

	// 批量注册
	for name, m := range models {
		s.models[name] = m
	}

	// 如果系统已运行,启动新注册的模型
	if s.isRunning {
		for name, m := range models {
			if err := m.Start(); err != nil {
				return fmt.Errorf("failed to start model %s: %w", name, err)
			}
		}
	}

	return nil
}

// UnregisterModel safely removes a model
func (s *System) UnregisterModel(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	model, exists := s.models[name]
	if !exists {
		return fmt.Errorf("model %s not found", name)
	}

	// 如果模型正在运行,先停止它
	if s.isRunning {
		if err := model.Stop(); err != nil {
			return fmt.Errorf("failed to stop model %s: %w", name, err)
		}
	}

	// 移除模型
	delete(s.models, name)

	return nil
}

// GetModel retrieves a registered model by name
func (s *System) GetModel(name string) (model.Model, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	m, exists := s.models[name]
	if !exists {
		return nil, types.ErrModelNotFound
	}

	return m, nil
}

// Core returns the core engine instance
func (s *System) Core() *core.Engine {
	return s.core
}

// IsRunning returns the current system state
func (s *System) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isRunning
}

// Context returns the system context
func (s *System) Context() context.Context {
	return s.ctx
}

// HandleEvent 处理系统事件
func (s *System) HandleEvent(event types.SystemEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查系统状态
	if !s.isRunning {
		return types.NewSystemError(types.ErrState, "system not running", nil)
	}

	// 添加到事件队列
	select {
	case s.events.queue <- event:
		// 成功添加到队列
	default:
		return types.NewSystemError(types.ErrQueue, "event queue full", nil)
	}

	// 记录事件
	s.state.events = append(s.state.events, event)
	if len(s.state.events) > types.MaxEventHistory {
		s.state.events = s.state.events[1:]
	}

	return nil
}

// Subscribe 订阅事件
func (s *System) Subscribe(eventType types.EventType, handler types.EventHandler) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if handler == nil {
		return types.NewSystemError(types.ErrValidation, "nil handler", nil)
	}

	s.events.handlers[eventType] = append(s.events.handlers[eventType], handler)
	return nil
}

// Unsubscribe 取消事件订阅
func (s *System) Unsubscribe(eventType types.EventType, handler types.EventHandler) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	handlers := s.events.handlers[eventType]
	for i, h := range handlers {
		if h == handler {
			s.events.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			return nil
		}
	}

	return types.NewSystemError(types.ErrNotFound, "handler not found", nil)
}

// processEvents 处理事件队列
func (s *System) processEvents() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case event := <-s.events.queue:
			s.dispatchEvent(event)
		}
	}
}

// dispatchEvent 分发事件到处理器
func (s *System) dispatchEvent(event types.SystemEvent) {
	s.mu.RLock()
	handlers := s.events.handlers[event.Type]
	s.mu.RUnlock()

	for _, handler := range handlers {
		go func(h types.EventHandler) {
			if err := h.HandleEvent(event); err != nil {
				s.recordError(err)
			}
		}(handler)
	}
}

// recordError records a system error
func (s *System) recordError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.state.errors = append(s.state.errors, err)
	if len(s.state.errors) > types.MaxErrorHistory {
		s.state.errors = s.state.errors[1:]
	}

	// 触发错误事件
	s.HandleEvent(types.SystemEvent{
		Type:      "system.error",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"error": err.Error(),
		},
	})
}

// updateMetrics 更新系统指标
func (s *System) updateMetrics() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	uptime := now.Sub(s.state.startTime)

	// 更新基本指标
	s.state.metrics.Status = s.state.status
	s.state.metrics.Health = s.calculateSystemHealth()
	s.state.metrics.AlertCount = 0              // TODO: 实现告警计数
	s.state.metrics.LastAlertTime = time.Time{} // TODO: 实现最后告警时间
	s.state.metrics.AlertLevels = make(map[types.AlertLevel]int)

	// 更新时序信息
	s.state.metrics.Timestamp = now
	s.state.metrics.Period = uptime.String()

	// 更新运行指标
	s.state.metrics.Uptime = uptime
	s.state.metrics.ErrorCount = len(s.state.errors)
	s.state.metrics.EventCount = len(s.state.events)

	// 更新统计信息
	s.state.metrics.Stats.LastUpdateTime = now
	s.state.metrics.Stats.TotalRequests = 0 // TODO: 实现请求统计
	s.state.metrics.Stats.SuccessCount = 0  // TODO: 实现成功计数
	s.state.metrics.Stats.FailureCount = 0  // TODO: 实现失败计数

	// 更新资源指标
	s.state.metrics.CPU = 0        // TODO: 实现CPU使用率
	s.state.metrics.Memory = 0     // TODO: 实现内存使用率
	s.state.metrics.Goroutines = 0 // TODO: 实现协程数

	// 收集子系统指标
	s.state.metrics.Subsystems = make(map[string]types.SubsystemMetrics)

	// 核心引擎指标
	s.state.metrics.Subsystems["core"] = types.SubsystemMetrics{
		Status:     s.state.status,
		Health:     1.0, // 基础健康度
		LastUpdate: now,
		Metrics:    make(map[string]float64),
	}

	// 控制子系统指标
	s.state.metrics.Subsystems["control"] = types.SubsystemMetrics{
		Status:     s.state.status,
		Health:     1.0,
		LastUpdate: now,
		Metrics:    make(map[string]float64),
	}

	// 演化子系统指标
	s.state.metrics.Subsystems["evolution"] = types.SubsystemMetrics{
		Status:     s.state.status,
		Health:     1.0,
		LastUpdate: now,
		Metrics:    make(map[string]float64),
	}

	// 元数据子系统指标
	s.state.metrics.Subsystems["meta"] = types.SubsystemMetrics{
		Status:     s.state.status,
		Health:     1.0,
		LastUpdate: now,
		Metrics:    make(map[string]float64),
	}

	// 监控子系统指标
	s.state.metrics.Subsystems["monitor"] = types.SubsystemMetrics{
		Status:     s.state.status,
		Health:     1.0,
		LastUpdate: now,
		Metrics:    make(map[string]float64),
	}

	// 计算系统健康度
	s.state.metrics.Health = s.calculateSystemHealth()
}

// calculateSystemHealth 计算系统整体健康度
func (s *System) calculateSystemHealth() float64 {
	// 基础分值
	baseScore := 1.0

	// 根据错误数量扣分
	errorPenalty := math.Min(float64(len(s.state.errors))*0.1, 0.5)
	baseScore -= errorPenalty

	// 检查子系统状态
	subsystemScores := make([]float64, 0)
	for _, metrics := range s.state.metrics.Subsystems {
		subsystemScores = append(subsystemScores, metrics.Health)
	}

	// 计算子系统平均健康度
	avgSubsystemHealth := 0.0
	if len(subsystemScores) > 0 {
		total := 0.0
		for _, score := range subsystemScores {
			total += score
		}
		avgSubsystemHealth = total / float64(len(subsystemScores))
	}

	// 综合评分
	return math.Max(0, math.Min(1, baseScore*0.4+avgSubsystemHealth*0.6))
}

// GetMetrics 获取系统指标
func (s *System) GetMetrics() types.SystemMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 更新指标
	s.updateMetrics()

	// 返回指标副本
	metrics := s.state.metrics
	return metrics
}

// GetStatus 获取系统状态
func (s *System) GetStatus() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state.status
}

// GetErrors 获取系统错误
func (s *System) GetErrors() []error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 返回错误副本
	errors := make([]error, len(s.state.errors))
	copy(errors, s.state.errors)
	return errors
}

// GetEvents 获取系统事件
func (s *System) GetEvents() []types.SystemEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 返回事件副本
	events := make([]types.SystemEvent, len(s.state.events))
	copy(events, s.state.events)
	return events
}

// GetSubsystemStatus 获取子系统状态
func (s *System) GetSubsystemStatus() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]string{
		"core":      s.core.Status(),
		"common":    s.common.Status(),
		"control":   s.control.Status(),
		"evolution": s.evolution.Status(),
		"meta":      s.meta.Status(),
		"monitor":   s.monitor.Status(),
	}
}

// GetDependencies 获取系统依赖关系
func (s *System) GetDependencies() map[string][]string {
	return map[string][]string{
		"core":      {}, // core 是基础层,无依赖
		"common":    {"core"},
		"control":   {"core", "common"},
		"evolution": {"core", "common", "control"},
		"meta":      {"core", "common", "control"},
		"monitor":   {"core", "common"},
	}
}

// ValidateDependencies 验证依赖关系
func (s *System) ValidateDependencies() error {
	deps := s.GetDependencies()

	// 验证每个组件的依赖
	for component, dependencies := range deps {
		for _, dep := range dependencies {
			if !s.isComponentRunning(dep) {
				return fmt.Errorf("dependency %s not running for component %s",
					dep, component)
			}
		}
	}

	return nil
}

// isComponentRunning 检查组件是否运行中
func (s *System) isComponentRunning(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	switch name {
	case "core":
		return s.core != nil && s.core.Status() == "running"
	case "common":
		return s.common != nil && s.common.Status() == "running"
	case "control":
		return s.control != nil && s.control.Status() == "running"
	case "evolution":
		return s.evolution != nil && s.evolution.Status() == "running"
	case "meta":
		return s.meta != nil && s.meta.Status() == "running"
	case "monitor":
		return s.monitor != nil && s.monitor.Status() == "running"
	default:
		return false
	}
}

// injectDependencies 注入组件依赖
func (s *System) injectDependencies() error {
	// 注入 Control 依赖
	if err := s.control.InjectDependencies(
		s.core,
		s.common,
	); err != nil {
		return fmt.Errorf("failed to inject control dependencies: %w", err)
	}

	// 注入 Evolution 依赖
	if err := s.evolution.InjectDependencies(
		s.core,
		s.common,
		s.control,
	); err != nil {
		return fmt.Errorf("failed to inject evolution dependencies: %w", err)
	}

	// 注入 Meta 依赖
	if err := s.meta.InjectDependencies(
		s.core,
		s.common,
		s.control,
	); err != nil {
		return fmt.Errorf("failed to inject meta dependencies: %w", err)
	}

	// 注入 Monitor 依赖
	if err := s.monitor.InjectDependencies(
		s.core,
		s.common,
	); err != nil {
		return fmt.Errorf("failed to inject monitor dependencies: %w", err)
	}

	return nil
}

// Coordinate 协调系统状态
func (s *System) Coordinate() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. 验证依赖关系
	if err := s.ValidateDependencies(); err != nil {
		return fmt.Errorf("dependency validation failed: %w", err)
	}

	// 2. 检查系统健康状态
	health := s.calculateSystemHealth()
	if health < 0.5 {
		return fmt.Errorf("system health too low: %f", health)
	}

	// 3. 更新系统指标
	s.updateMetrics()

	// 4. 协调子系统状态
	for name, status := range s.GetSubsystemStatus() {
		if status != "running" {
			s.HandleEvent(types.SystemEvent{
				Type:      "system.coordination",
				Timestamp: time.Now(),
				Data: map[string]interface{}{
					"component": name,
					"status":    status,
					"action":    "recovery_needed",
				},
			})
		}
	}

	return nil
}

// RestoreSubsystem 恢复子系统
func (s *System) RestoreSubsystem(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch name {
	case "common":
		return s.common.Restore(s.ctx)
	case "control":
		return s.control.Restore(s.ctx)
	case "evolution":
		return s.evolution.Restore(s.ctx)
	case "meta":
		return s.meta.Restore(s.ctx)
	case "monitor":
		return s.monitor.Restore(s.ctx)
	default:
		return fmt.Errorf("unknown subsystem: %s", name)
	}
}

// ListModels 获取所有注册的模型名称列表
func (s *System) ListModels() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	names := make([]string, 0, len(s.models))
	for name := range s.models {
		names = append(names, name)
	}
	return names
}

// TransformModel 执行模型转换
func (s *System) TransformModel(ctx context.Context, pattern model.TransformPattern) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return types.ErrNotRunning
	}

	// 获取并验证当前状态
	state := s.getCurrentState()
	if err := model.ValidateSystemState(state); err != nil {
		return err
	}

	// 执行转换
	for name, m := range s.models {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := m.Transform(pattern); err != nil {
				return fmt.Errorf("failed to transform model %s: %w", name, err)
			}
		}
	}

	return s.evolution.UpdateState()
}

// getCurrentState 获取当前系统状态
func (s *System) getCurrentState() *model.SystemState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 转换status为Phase
	var phase model.Phase
	switch s.state.status {
	case "running":
		phase = model.PhaseTransform
	case "stable":
		phase = model.Phase_Stable
	case "unstable":
		phase = model.Phase_Unstable
	default:
		phase = model.PhaseNeutral
	}

	return &model.SystemState{
		Energy:    s.core.GetTotalEnergy(),
		Phase:     phase,
		Timestamp: time.Now(),
		Properties: map[string]interface{}{
			"status":  s.state.status,
			"health":  s.calculateSystemHealth(),
			"metrics": s.state.metrics,
		},
	}
}

// GetEnergy 获取系统总能量
func (s *System) GetEnergy() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.core.GetTotalEnergy()
}

// AdjustEnergy 调整系统总能量
func (s *System) AdjustEnergy(delta float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 验证参数
	if delta < -1.0 || delta > 1.0 {
		return types.ErrInvalidParameter
	}

	// 调整core层能量
	currentEnergy := s.core.GetTotalEnergy()
	newEnergy := currentEnergy + delta

	// 确保能量在[0,1]范围内
	if newEnergy < 0 || newEnergy > 1.0 {
		return types.ErrEnergyOutOfRange
	}

	// 更新系统状态
	s.state.energy = newEnergy

	return nil
}

// GetEnergySystem 获取能量系统
func (s *System) GetEnergySystem() *core.EnergySystem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.core.GetEnergySystem()
}

// PublishEvent 发布系统事件
func (s *System) PublishEvent(event types.Event) error {
	// 创建SystemEvent
	sysEvent := types.SystemEvent{
		ID:        event.ID,
		Type:      event.Type,
		Source:    event.Source,
		Timestamp: event.Timestamp,
		Data:      event.Payload,
		Metadata:  map[string]string{},
	}

	// 转换元数据
	for k, v := range event.Metadata {
		if strVal, ok := v.(string); ok {
			sysEvent.Metadata[k] = strVal
		}
	}

	// 使用HandleEvent处理
	return s.HandleEvent(sysEvent)
}

// GetYinYangFlow 获取阴阳流模型
func (s *System) GetYinYangFlow() *model.YinYangFlow {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 从模型管理器获取阴阳流模型实例
	return s.modelManager.GetYinYangFlow()
}

// TransformYinYang 执行阴阳模型转换
func (s *System) TransformYinYang(pattern model.TransformPattern) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return types.ErrNotRunning
	}

	// 获取阴阳流模型
	yinyang := s.modelManager.GetYinYangFlow()
	if yinyang == nil {
		return types.NewSystemError(types.ErrCodeModel, "yinyang model not found", nil)
	}

	// 执行转换
	return yinyang.Transform(pattern)
}

// GetBaGuaFlow 获取八卦模型
func (s *System) GetBaGuaFlow() *model.BaGuaFlow {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 从模型管理器获取八卦流模型实例
	return s.modelManager.GetBaGuaFlow()
}

// GetFieldSystem 获取场系统
func (s *System) GetFieldSystem() *core.FieldSystem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.core.GetFieldSystem()
}

// GetGanZhiFlow 获取干支模型
func (s *System) GetGanZhiFlow() *model.GanZhiFlow {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 从模型管理器获取干支流模型实例
	return s.modelManager.GetGanZhiFlow()
}

// GetModelMetrics 获取模型指标
func (s *System) GetModelMetrics() model.ModelMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 更新系统指标
	s.updateMetrics()

	// 初始化ModelMetrics
	metrics := model.ModelMetrics{}

	// 转换能量指标
	metrics.Energy.Total = s.state.metrics.System.Energy
	metrics.Energy.Average = s.calculateAverageEnergy()
	metrics.Energy.Variance = s.calculateEnergyVariance()

	// 转换场和量子状态
	metrics.Quantum = s.state.metrics.System.Quantum
	metrics.Field = s.state.metrics.System.Field

	// 设置性能指标
	metrics.Performance.Throughput = float64(s.state.metrics.Stats.SuccessCount) / 60.0 // 每分钟
	metrics.Performance.QPS = float64(s.state.metrics.Stats.TotalRequests) / 60.0
	metrics.Performance.ErrorRate = float64(s.state.metrics.ErrorCount) / math.Max(1.0, float64(s.state.metrics.Stats.TotalRequests))

	return metrics
}

// calculateAverageEnergy 计算平均能量
func (s *System) calculateAverageEnergy() float64 {
	// 从历史指标中获取能量数据
	total := 0.0
	count := 0

	// 使用系统历史指标
	for _, event := range s.state.events {
		if metrics, ok := event.Data.(map[string]interface{}); ok {
			if energy, exists := metrics["energy"].(float64); exists {
				total += energy
				count++
			}
		}
	}

	// 防止除零错误
	if count == 0 {
		return s.state.metrics.System.Energy // 返回当前能量作为默认值
	}

	return total / float64(count)
}

// calculateEnergyVariance 计算能量方差
func (s *System) calculateEnergyVariance() float64 {
	avg := s.calculateAverageEnergy()
	sumSquares := 0.0
	count := 0

	// 计算平方差之和
	for _, event := range s.state.events {
		if metrics, ok := event.Data.(map[string]interface{}); ok {
			if energy, exists := metrics["energy"].(float64); exists {
				diff := energy - avg
				sumSquares += diff * diff
				count++
			}
		}
	}

	// 防止除零错误
	if count == 0 {
		return 0.01 // 返回一个小的默认方差
	}

	return sumSquares / float64(count)
}

// GetModelState 获取模型状态
func (s *System) GetModelState() model.ModelState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 使用集成模型管理器获取模型状态
	return s.modelManager.GetState()
}

// GetQuantumSystem 获取量子系统
func (s *System) GetQuantumSystem() *core.QuantumSystem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.core.GetQuantumSystem()
}

// GetState 获取系统状态
func (s *System) GetState() model.SystemState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 使用现有的getCurrentState方法获取完整状态
	state := s.getCurrentState()

	// 检查是否为nil以防止空指针
	if state == nil {
		return model.SystemState{
			Energy:     0,
			Entropy:    0,
			Harmony:    0,
			Balance:    0,
			Phase:      model.PhaseNone,
			Timestamp:  time.Now(),
			Properties: make(map[string]interface{}),
		}
	}

	return *state
}

// GetWuXingFlow 获取五行模型
func (s *System) GetWuXingFlow() *model.WuXingFlow {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 从模型管理器获取五行流模型实例
	return s.modelManager.GetWuXingFlow()
}

// Optimize 执行系统优化
func (s *System) Optimize(params types.OptimizationParams) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return types.ErrNotRunning
	}

	// 委托给evolution管理器处理优化
	return s.evolution.Optimize(params)
}

// Synchronize 同步系统状态
func (s *System) Synchronize(params types.SyncParams) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return types.ErrNotRunning
	}

	// 委托给控制子系统处理同步
	return s.control.Synchronize(params)
}

// Transform 执行系统转换
func (s *System) Transform(pattern model.TransformPattern) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return types.ErrNotRunning
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 委托给TransformModel处理
	return s.TransformModel(ctx, pattern)
}
