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
	models map[string]model.Model

	// System subsystems
	common    *control.Manager   // Common utilities and shared resources
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
	}

	// Event handling
	events struct {
		handlers  map[string][]types.EventHandler // 事件处理器
		queue     chan types.SystemEvent          // 事件队列
		processor *types.EventProcessor           // 事件处理器
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
	ModelConfig     *model.Config
	CommonConfig    *common.Config
	ControlConfig   *control.Config
	EvolutionConfig *evolution.Config
	MetaConfig      *meta.Config
	MonitorConfig   *monitor.Config
}

// New creates a new System instance
func New(cfg *Config) (*System, error) {
	if cfg == nil {
		cfg = defaultConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	sys := &System{
		models: make(map[string]model.Model),
		ctx:    ctx,
		cancel: cancel,
		config: cfg,
	}

	// 初始化事件系统
	sys.events.handlers = make(map[string][]types.EventHandler)
	sys.events.queue = make(chan types.SystemEvent, 1000)
	sys.events.processor = types.NewEventProcessor()

	// 初始化状态
	sys.state.status = "initialized"
	sys.state.startTime = time.Now()
	sys.state.errors = make([]error, 0)
	sys.state.events = make([]types.SystemEvent, 0)
	sys.state.metrics = types.SystemMetrics{
		StartTime: time.Now(),
	}

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
func defaultConfig() *Config {
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

// Start initializes and starts all system components
func (s *System) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return nil
	}

	// 更新系统状态
	s.state.status = "starting"
	startTime := time.Now()

	// 启动核心引擎
	if err := s.core.Start(s.ctx); err != nil {
		s.state.status = "failed"
		return fmt.Errorf("failed to start core engine: %w", err)
	}

	// 启动子系统,按依赖顺序启动
	if err := s.startSubsystems(); err != nil {
		s.core.Stop()
		s.state.status = "failed"
		return fmt.Errorf("failed to start subsystems: %w", err)
	}

	// 启动所有模型
	for name, m := range s.models {
		if err := m.Start(s.ctx); err != nil {
			s.stopSubsystems()
			s.core.Stop()
			s.state.status = "failed"
			return fmt.Errorf("failed to start model %s: %w", name, err)
		}
	}

	// 更新系统状态
	s.isRunning = true
	s.state.status = "running"
	s.state.startTime = startTime

	// 发送系统启动事件
	s.HandleEvent(types.SystemEvent{
		Type:      "system.started",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"startup_time": time.Since(startTime).String(),
		},
	})

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

// Stop gracefully shuts down all system components
func (s *System) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return nil
	}

	// 更新系统状态
	s.state.status = "stopping"

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	// 发送系统关闭事件
	s.HandleEvent(types.SystemEvent{
		Type:      "system.stopping",
		Timestamp: time.Now(),
	})

	// 等待所有事件处理完成
	close(s.events.queue)
	s.events.processor.Wait()

	// 停止所有模型
	for name, m := range s.models {
		if err := m.Stop(); err != nil {
			s.recordError(fmt.Errorf("failed to stop model %s: %w", name, err))
		}
	}

	// 停止所有子系统
	if err := s.stopSubsystems(); err != nil {
		s.recordError(fmt.Errorf("failed to stop subsystems: %w", err))
	}

	// 停止核心引擎
	if err := s.core.Stop(); err != nil {
		s.recordError(fmt.Errorf("failed to stop core engine: %w", err))
	}

	// 取消上下文
	s.cancel()

	// 等待所有组件完全停止或超时
	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			s.recordError(fmt.Errorf("system shutdown timed out"))
		}
	case <-s.waitForComponents():
		// 所有组件已停止
	}

	s.isRunning = false
	s.state.status = "stopped"

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
	s.state.errors = make([]error, 0)
	s.state.events = make([]types.SystemEvent, 0)
	s.state.metrics = types.SystemMetrics{
		StartTime: time.Now(),
	}

	// 重置事件系统
	s.events.handlers = make(map[string][]types.EventHandler)
	s.events.queue = make(chan types.SystemEvent, 1000)
	s.events.processor = types.NewEventProcessor()

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
			if err := m.Start(s.ctx); err != nil {
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
func (s *System) Subscribe(eventType string, handler types.EventHandler) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if handler == nil {
		return types.NewSystemError(types.ErrValidation, "nil handler", nil)
	}

	s.events.handlers[eventType] = append(s.events.handlers[eventType], handler)
	return nil
}

// Unsubscribe 取消事件订阅
func (s *System) Unsubscribe(eventType string, handler types.EventHandler) error {
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
			if err := h.Handle(event); err != nil {
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
	s.state.metrics.Uptime = uptime
	s.state.metrics.LastUpdate = now
	s.state.metrics.Status = s.state.status
	s.state.metrics.ErrorCount = len(s.state.errors)
	s.state.metrics.EventCount = len(s.state.events)

	// 收集子系统指标
	s.state.metrics.Subsystems = map[string]types.SubsystemMetrics{
		"core": {
			Status:     s.core.Status(),
			Metrics:    s.core.GetMetrics(),
			LastUpdate: now,
		},
		"control": {
			Status:     s.control.Status(),
			Metrics:    s.control.GetMetrics(),
			LastUpdate: now,
		},
		"evolution": {
			Status:     s.evolution.Status(),
			Metrics:    s.evolution.GetMetrics(),
			LastUpdate: now,
		},
		"meta": {
			Status:     s.meta.Status(),
			Metrics:    s.meta.GetMetrics(),
			LastUpdate: now,
		},
		"monitor": {
			Status:     s.monitor.Status(),
			Metrics:    s.monitor.GetMetrics(),
			LastUpdate: now,
		},
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
		if health, ok := metrics.Metrics["health"].(float64); ok {
			subsystemScores = append(subsystemScores, health)
		}
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
