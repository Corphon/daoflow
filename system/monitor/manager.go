// system/monitor/manager.go

package monitor

import (
	"context"
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/common"
	"github.com/Corphon/daoflow/system/monitor/alert"
	"github.com/Corphon/daoflow/system/monitor/metrics"
	"github.com/Corphon/daoflow/system/monitor/trace"
	"github.com/Corphon/daoflow/system/types"
)

// Manager 监控系统管理器
type Manager struct {
	mu sync.RWMutex

	// 基础配置
	config *types.MonitorConfig

	// 监控组件
	components struct {
		collector *metrics.Collector // 指标收集器
		analyzer  *metrics.Analyzer  // 指标分析器
		reporter  *metrics.Reporter  // 指标报告器
		detector  *alert.Detector    // 告警检测器
		handler   *alert.Handler     // 告警处理器
		notifier  *alert.Notifier    // 告警通知器
		tracker   *trace.Tracker     // 追踪器
		recorder  *trace.Recorder    // 记录器
		analyzer2 *trace.Analyzer    // 追踪分析器
	}

	// 监控状态
	state struct {
		status     string               // 运行状态
		startTime  time.Time            // 启动时间
		lastUpdate time.Time            // 最后更新
		metrics    types.MonitorMetrics // 监控指标
		errors     []error              // 错误记录
	}

	// 核心依赖
	core   *core.Engine
	common *common.Manager

	// 上下文控制
	ctx    context.Context
	cancel context.CancelFunc
}

// ----------------------------------------------------------
// NewManager 创建新的管理器实例
func NewManager(cfg *types.MonitorConfig) (*Manager, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	m := &Manager{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
	}

	// 初始化各组件
	if err := m.initComponents(); err != nil {
		cancel()
		return nil, err
	}

	// 初始化状态
	m.state.status = "initialized"
	m.state.startTime = time.Now()

	return m, nil
}

// DefaultConfig 返回默认配置
func DefaultConfig() *types.MonitorConfig {
	return &types.MonitorConfig{
		Base: struct {
			SampleInterval time.Duration `json:"sample_interval"` // 采样间隔
			BatchSize      int           `json:"batch_size"`      // 批处理大小
			BufferSize     int           `json:"buffer_size"`     // 缓冲区大小
			RetentionTime  time.Duration `json:"retention_time"`  // 保留时间
			MaxHistory     int           `json:"max_history"`     // 最大历史记录字段
		}{
			SampleInterval: time.Second,
			BatchSize:      100,
			BufferSize:     1000,
			RetentionTime:  24 * time.Hour,
		},
		Metrics: struct {
			Enabled       bool          `json:"enabled"`
			Interval      time.Duration `json:"interval"`
			HistorySize   int           `json:"history_size"`
			EnabledTypes  []string      `json:"enabled_types"`
			CustomMetrics []string      `json:"custom_metrics"`
			Aggregation   struct {
				Interval   time.Duration `json:"interval"`
				Functions  []string      `json:"functions"`
				WindowSize int           `json:"window_size"`
			} `json:"aggregation"`
		}{
			Enabled:       true,
			Interval:      5 * time.Second,
			HistorySize:   1000,
			EnabledTypes:  []string{"system", "process", "resource", "performance"},
			CustomMetrics: []string{},
			Aggregation: struct {
				Interval   time.Duration `json:"interval"`
				Functions  []string      `json:"functions"`
				WindowSize int           `json:"window_size"`
			}{
				Interval:   time.Minute,
				Functions:  []string{"avg", "max", "min"},
				WindowSize: 60,
			},
		},
		Alert: types.AlertConfig{
			Enabled:       true,
			CheckInterval: time.Minute,
			MaxAlerts:     1000,
			MaxConcurrent: 10,
			RetryCount:    3,
			Timeout:       time.Minute,
			QueueSize:     1000,
			BufferSize:    1000,
			BatchSize:     100,
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

	// 启动各组件
	if err := m.startComponents(); err != nil {
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

	// 停止各组件
	if err := m.stopComponents(); err != nil {
		return err
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
		"status":      m.state.status,
		"uptime":      time.Since(m.state.startTime).String(),
		"collector":   m.components.collector.GetMetricsData(),
		"alerts":      len(m.components.detector.GetAlertChannel()),
		"traces":      m.components.tracker.GetMetrics(),
		"error_count": len(m.state.errors),
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
	// 转换配置类型
	metricsConfig := types.MetricsConfig{
		Base: struct {
			Enabled        bool               `json:"enabled"`
			Interval       time.Duration      `json:"interval"`
			SampleInterval time.Duration      `json:"sample_interval"`
			BufferSize     int                `json:"buffer_size"`
			MaxHistory     int                `json:"max_history"`
			MaxHistorySize int                `json:"max_history_size"`
			HistorySize    int                `json:"history_size"`
			Thresholds     map[string]float64 `json:"thresholds"`
		}{
			Enabled:        m.config.Metrics.Enabled,
			Interval:       m.config.Metrics.Interval,
			SampleInterval: m.config.Base.SampleInterval,
			BufferSize:     m.config.Base.BufferSize,
			MaxHistory:     m.config.Base.MaxHistory,
			MaxHistorySize: m.config.Metrics.HistorySize,
			HistorySize:    m.config.Metrics.HistorySize,
			Thresholds:     make(map[string]float64),
		},
		Report: m.config.Report,
	}

	// 创建指标收集器
	collector := metrics.NewCollector(metricsConfig)
	m.components.collector = collector

	// 创建指标分析器
	analyzer := metrics.NewAnalyzer(collector, metricsConfig)
	m.components.analyzer = analyzer

	// 创建指标报告器
	reporter := metrics.NewReporter(collector, metricsConfig)
	m.components.reporter = reporter

	// 创建告警检测器
	detector := alert.NewDetector(m.config.Alert, collector)
	m.components.detector = detector

	// 创建告警处理器
	handler := alert.NewHandler()
	m.components.handler = handler

	// 创建告警通知器
	notifier := alert.NewNotifier(m.config.Alert)
	m.components.notifier = notifier

	// 配置转换
	traceConfig := types.TraceConfig{
		StoragePath:   m.config.Trace.StoragePath,
		RetentionDays: m.config.Base.RetentionTime,
		BatchSize:     m.config.Base.BatchSize,
		BufferSize:    m.config.Trace.BufferSize,
		FlushInterval: m.config.Trace.FlushInterval,
		AsyncWrite:    true,
		SampleRate:    m.config.Trace.SampleRate,
		MaxQueueSize:  m.config.Trace.MaxSpans,
		EnableMetrics: true,
		EnableEvents:  true,
		IncludeModel:  true,
	}

	// 创建追踪器
	tracker := trace.NewTracker(traceConfig)
	m.components.tracker = tracker

	// 创建记录器
	recorder := trace.NewRecorder(traceConfig)
	m.components.recorder = recorder

	// 创建追踪分析器
	analyzer2 := trace.NewAnalyzer(tracker, recorder, traceConfig)
	m.components.analyzer2 = analyzer2

	return nil
}

// startComponents 启动组件
func (m *Manager) startComponents() error {
	// 按依赖顺序启动
	if err := m.components.collector.Start(m.ctx); err != nil {
		return err
	}
	if err := m.components.analyzer.Start(m.ctx); err != nil {
		return err
	}
	if err := m.components.reporter.Start(m.ctx); err != nil {
		return err
	}
	if err := m.components.detector.Start(m.ctx); err != nil {
		return err
	}
	if err := m.components.handler.Start(m.ctx); err != nil {
		return err
	}
	if err := m.components.notifier.Start(m.ctx); err != nil {
		return err
	}
	if err := m.components.tracker.Start(m.ctx); err != nil {
		return err
	}
	if err := m.components.recorder.Start(m.ctx); err != nil {
		return err
	}
	if err := m.components.analyzer2.Start(m.ctx); err != nil {
		return err
	}
	return nil
}

// stopComponents 停止组件
func (m *Manager) stopComponents() error {
	// 按依赖反序停止
	if err := m.components.analyzer2.Stop(); err != nil {
		return err
	}
	if err := m.components.recorder.Stop(); err != nil {
		return err
	}
	if err := m.components.tracker.Stop(); err != nil {
		return err
	}
	if err := m.components.notifier.Stop(); err != nil {
		return err
	}
	if err := m.components.handler.Stop(); err != nil {
		return err
	}
	if err := m.components.detector.Stop(); err != nil {
		return err
	}
	if err := m.components.reporter.Stop(); err != nil {
		return err
	}
	if err := m.components.analyzer.Stop(); err != nil {
		return err
	}
	if err := m.components.collector.Stop(); err != nil {
		return err
	}
	return nil
}

// InjectDependencies 注入组件依赖
func (m *Manager) InjectDependencies(core *core.Engine, common *common.Manager) error {
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

	return nil
}

// Restore 恢复系统
func (m *Manager) Restore(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 重置状态
	m.state.metrics = types.MonitorMetrics{}
	m.state.errors = make([]error, 0)
	m.state.lastUpdate = time.Now()

	// 重置组件
	return m.initComponents()
}
