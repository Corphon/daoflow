// system/meta/manager.go

package meta

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/common"
	"github.com/Corphon/daoflow/system/control"
	"github.com/Corphon/daoflow/system/meta/emergence"
	"github.com/Corphon/daoflow/system/meta/field"
	"github.com/Corphon/daoflow/system/meta/resonance"
	"github.com/Corphon/daoflow/system/types"
)

// Manager 元系统管理器
type Manager struct {
	mu sync.RWMutex

	// 基础配置
	config *types.MetaConfig

	// 元系统组件
	components struct {
		field     *field.UnifiedField           // 统一场
		detector  *emergence.PatternDetector    // 模式检测器
		matcher   *resonance.PatternMatcher     // 模式匹配器
		amplifier *resonance.ResonanceAmplifier // 共振放大器
	}

	// 元系统状态
	state struct {
		status    string                  // 运行状态
		startTime time.Time               // 启动时间
		emergence []types.EmergentPattern // 涌现模式
		resonance []common.ResonanceState // 共振状态
		energy    float64                 // 系统能量
		metrics   map[string]float64      // 系统指标
	}

	// 核心依赖
	core    *core.Engine
	common  *common.Manager
	control *control.Manager

	// 上下文控制
	ctx    context.Context
	cancel context.CancelFunc
}

// NewManager 创建新的管理器实例
func NewManager(cfg *types.MetaConfig) (*Manager, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	m := &Manager{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
	}

	// 初始化组件
	if err := m.initComponents(); err != nil {
		cancel()
		return nil, err
	}

	// 初始化状态
	m.state.status = "initialized"
	m.state.startTime = time.Now()
	m.state.metrics = make(map[string]float64)

	return m, nil
}

// DefaultConfig 返回默认配置
func DefaultConfig() *types.MetaConfig {
	return &types.MetaConfig{
		Field: struct {
			InitialStrength float64       `json:"initial_strength"`
			MinStrength     float64       `json:"min_strength"`
			MaxStrength     float64       `json:"max_strength"`
			Dimension       int           `json:"dimension"`
			UpdateInterval  time.Duration `json:"update_interval"`
			Coupling        struct {
				Strength  float64 `json:"strength"`
				Range     float64 `json:"range"`
				Threshold float64 `json:"threshold"`
				MaxPairs  int     `json:"max_pairs"`
			} `json:"coupling"`
		}{
			InitialStrength: 1.0,
			MinStrength:     0.1,
			MaxStrength:     10.0,
			Dimension:       3,
			UpdateInterval:  time.Second,
			Coupling: struct {
				Strength  float64 `json:"strength"`
				Range     float64 `json:"range"`
				Threshold float64 `json:"threshold"`
				MaxPairs  int     `json:"max_pairs"`
			}{
				Strength:  0.7,
				Range:     1.0,
				Threshold: 0.5,
				MaxPairs:  100,
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
		"status":    m.state.status,
		"uptime":    time.Since(m.state.startTime).String(),
		"energy":    m.state.energy,
		"emergence": len(m.state.emergence),
		"resonance": len(m.state.resonance),
		"field":     m.components.field.GetMetrics(),
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
	// 1. 初始化统一场
	field, err := field.NewUnifiedField(m.config.Field.InitialStrength)
	if err != nil {
		return err
	}
	m.components.field = field

	// 2. 初始化模式检测器
	detector := emergence.NewPatternDetector(field)
	if detector == nil {
		return fmt.Errorf("failed to create pattern detector")
	}
	m.components.detector = detector

	// 3. 初始化属性生成器
	propertyGenerator := emergence.NewPropertyGenerator(detector, field)
	if propertyGenerator == nil {
		return fmt.Errorf("failed to create property generator")
	}

	// 4. 初始化模式匹配器
	matcher := resonance.NewPatternMatcher(detector, nil) // amplifier will be set later
	if matcher == nil {
		return fmt.Errorf("failed to create pattern matcher")
	}
	m.components.matcher = matcher

	// 5. 初始化共振放大器
	amplifier := resonance.NewResonanceAmplifier(field, detector, propertyGenerator)
	if amplifier == nil {
		return fmt.Errorf("failed to create resonance amplifier")
	}
	m.components.amplifier = amplifier

	// 6. 设置匹配器的放大器引用
	matcher.SetAmplifier(amplifier)

	return nil
}

// startComponents 启动组件
func (m *Manager) startComponents() error {
	// 1. 启动统一场
	if err := m.components.field.Start(m.ctx); err != nil {
		return fmt.Errorf("failed to start field: %w", err)
	}

	// 2. 启动模式检测器
	if err := m.components.detector.Start(m.ctx); err != nil {
		m.components.field.Stop()
		return fmt.Errorf("failed to start detector: %w", err)
	}

	// 3. 启动模式匹配器
	if err := m.components.matcher.Start(m.ctx); err != nil {
		m.components.detector.Stop()
		m.components.field.Stop()
		return fmt.Errorf("failed to start matcher: %w", err)
	}

	// 4. 启动共振放大器
	if err := m.components.amplifier.Start(m.ctx); err != nil {
		m.components.matcher.Stop()
		m.components.detector.Stop()
		m.components.field.Stop()
		return fmt.Errorf("failed to start amplifier: %w", err)
	}

	return nil
}

// stopComponents 停止组件
func (m *Manager) stopComponents() error {
	var errs []error

	// 按依赖关系的反序停止组件
	// 1. 停止共振放大器
	if err := m.components.amplifier.Stop(); err != nil {
		errs = append(errs, fmt.Errorf("failed to stop amplifier: %w", err))
	}

	// 2. 停止模式匹配器
	if err := m.components.matcher.Stop(); err != nil {
		errs = append(errs, fmt.Errorf("failed to stop matcher: %w", err))
	}

	// 3. 停止模式检测器
	if err := m.components.detector.Stop(); err != nil {
		errs = append(errs, fmt.Errorf("failed to stop detector: %w", err))
	}

	// 4. 停止统一场
	if err := m.components.field.Stop(); err != nil {
		errs = append(errs, fmt.Errorf("failed to stop field: %w", err))
	}

	// 如果有错误，返回组合后的错误信息
	if len(errs) > 0 {
		var combined error
		for i, err := range errs {
			if i == 0 {
				combined = fmt.Errorf("failed to stop components: %w", err)
			} else {
				combined = fmt.Errorf("%v; %w", combined, err)
			}
		}
		return combined
	}

	return nil
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
	m.state.emergence = make([]types.EmergentPattern, 0)
	m.state.resonance = make([]common.ResonanceState, 0)
	m.state.energy = 0
	m.state.metrics = make(map[string]float64)

	// 重置组件
	return m.initComponents()
}
