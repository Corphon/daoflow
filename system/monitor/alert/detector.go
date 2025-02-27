// system/monitor/alert/detector.go

package alert

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/types"
)

// AlertCondition 告警条件
type AlertCondition struct {
	ID        string
	Name      string
	Type      string
	Metric    string
	Operator  string
	Threshold float64
	Duration  time.Duration
	Severity  types.AlertLevel
	Labels    map[string]string
	Actions   []string
	ModelType model.ModelType // 关联的模型类型
}

// AlertState 告警状态
type AlertState struct {
	ConditionID string
	Active      bool
	StartTime   time.Time
	LastUpdate  time.Time
	Value       float64
	Count       int
	ModelState  model.ModelState // 关联的模型状态
}

// Detector 告警检测器
type Detector struct {
	mu sync.RWMutex

	// 配置
	config types.AlertConfig // 使用 types 包的配置类型

	// 条件管理
	conditions map[string]*AlertCondition
	states     map[string]*AlertState

	// 通知通道
	alertChan chan types.AlertData // 使用 types.AlertData

	// 检测状态
	status struct {
		isRunning  bool
		lastCheck  time.Time
		errorCount int
		errors     []error
	}

	// 指标源
	metricsSource MetricsSource

	// 模型状态管理器
	modelStateManager *model.StateManager
}

// MetricsSource 指标数据源接口
type MetricsSource interface {
	GetMetric(name string) (float64, error)
	GetMetrics() (map[string]float64, error)
	GetModelMetrics() (model.ModelMetrics, error) // 添加获取模型指标的方法
}

// --------------------------------------------------------------------------
// NewDetector 创建新的检测器
func NewDetector(config types.AlertConfig, source MetricsSource) *Detector {
	d := &Detector{
		// 基础配置
		config: config,

		// 状态存储
		conditions: make(map[string]*AlertCondition),
		states:     make(map[string]*AlertState),

		// 告警通道
		alertChan: make(chan types.AlertData, config.BufferSize),

		// 指标源
		metricsSource: source,

		// 模型管理器
		modelStateManager: model.NewStateManager(
			model.ModelTypeAlert,  // 告警模型类型
			model.MaxSystemEnergy, // 使用系统最大能量作为容量
		),
	}

	// 初始化状态
	d.status.isRunning = false
	d.status.lastCheck = time.Now()

	return d
}

// Start 启动检测器
func (d *Detector) Start(ctx context.Context) error {
	d.mu.Lock()
	if d.status.isRunning {
		d.mu.Unlock()
		return model.WrapError(nil, model.ErrCodeOperation, "detector already running")
	}
	d.status.isRunning = true
	d.mu.Unlock()

	go d.detectionLoop(ctx)
	return nil
}

// detectionLoop 检测循环
func (d *Detector) detectionLoop(ctx context.Context) {
	ticker := time.NewTicker(d.config.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := d.check(); err != nil {
				d.recordError(err)
			}
		}
	}
}

// Stop 停止检测器
func (d *Detector) Stop() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.status.isRunning {
		return model.WrapError(nil, model.ErrCodeOperation, "detector not running")
	}

	d.status.isRunning = false
	return nil
}

// AddCondition 添加告警条件
func (d *Detector) AddCondition(condition *AlertCondition) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.conditions[condition.ID]; exists {
		return model.WrapError(nil, model.ErrCodeValidation, "condition already exists")
	}

	d.conditions[condition.ID] = condition
	d.states[condition.ID] = &AlertState{
		ConditionID: condition.ID,
		Active:      false,
	}

	return nil
}

// check 执行检查
func (d *Detector) check() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 获取系统指标
	metrics, err := d.metricsSource.GetMetrics()
	if err != nil {
		return model.WrapError(err, model.ErrCodeOperation, "failed to get metrics")
	}

	// 获取模型指标
	modelMetrics, err := d.metricsSource.GetModelMetrics()
	if err != nil {
		return model.WrapError(err, model.ErrCodeOperation, "failed to get model metrics")
	}

	// 检查每个条件
	for id, condition := range d.conditions {
		state := d.states[id]

		// 获取当前值
		value, exists := metrics[condition.Metric]
		if !exists {
			continue
		}

		// 检查条件
		exceeded := d.evaluateCondition(condition, value, modelMetrics)

		if exceeded {
			if !state.Active {
				// 新的告警
				state.Active = true
				state.StartTime = time.Now()
				state.Count = 1
			} else {
				state.Count++
			}

			if time.Since(state.StartTime) >= condition.Duration {
				if err := d.triggerAlert(condition, value, modelMetrics); err != nil {
					d.recordError(err)
				}
			}
		} else {
			if state.Active {
				if err := d.resolveAlert(condition); err != nil {
					d.recordError(err)
				}
			}
			state.Active = false
			state.Count = 0
		}

		state.LastUpdate = time.Now()
		state.Value = value
	}

	d.status.lastCheck = time.Now()
	return nil
}

// errorCodeToAlertLevel 错误码转告警级别
func errorCodeToAlertLevel(code model.ErrorCode) types.AlertLevel {
	switch code {
	case model.ErrCodeCritical:
		return types.AlertLevelCritical
	case model.ErrCodeError:
		return types.AlertLevelError
	case model.ErrCodeWarning:
		return types.AlertLevelWarning
	default:
		return types.AlertLevelInfo
	}
}

// triggerAlert 触发告警
func (d *Detector) triggerAlert(condition *AlertCondition, value float64, modelMetrics model.ModelMetrics) error {
	alert := types.AlertData{
		ID:      generateAlertID(),
		Type:    condition.Type,
		Level:   condition.Severity,
		Message: generateAlertMessage(condition, value),
		Source:  condition.Name,
		Time:    time.Now(),
		Status:  "firing",
		Labels:  condition.Labels,
		ModelData: &types.ModelAlertData{
			Type:    condition.ModelType,
			Metrics: &modelMetrics,
		},
	}

	select {
	case d.alertChan <- alert:
		return nil
	default:
		return model.WrapError(nil, model.ErrCodeResource, "alert channel full")
	}
}

// resolveAlert 解除告警
func (d *Detector) resolveAlert(condition *AlertCondition) error {
	alert := types.AlertData{
		ID:      generateAlertID(),
		Type:    condition.Type,
		Level:   condition.Severity,
		Message: "Alert resolved: " + condition.Name,
		Source:  condition.Name,
		Time:    time.Now(),
		Status:  "resolved",
		Labels:  condition.Labels,
	}

	select {
	case d.alertChan <- alert:
		return nil
	default:
		return model.WrapError(nil, model.ErrCodeResource, "alert channel full")
	}
}

// evaluateCondition 评估条件
func (d *Detector) evaluateCondition(condition *AlertCondition, value float64, modelMetrics model.ModelMetrics) bool {
	// 基本阈值检查
	thresholdMet := d.checkThreshold(condition.Operator, value, condition.Threshold)

	// 如果配置了模型类型，还需要检查模型指标
	if condition.ModelType != model.ModelTypeNone {
		return thresholdMet && d.checkModelMetrics(condition, modelMetrics)
	}

	return thresholdMet
}

// checkThreshold 检查阈值
func (d *Detector) checkThreshold(operator string, value, threshold float64) bool {
	switch operator {
	case ">":
		return value > threshold
	case ">=":
		return value >= threshold
	case "<":
		return value < threshold
	case "<=":
		return value <= threshold
	case "==":
		return value == threshold
	default:
		return false
	}
}

// checkModelMetrics 检查模型指标
func (d *Detector) checkModelMetrics(condition *AlertCondition, metrics model.ModelMetrics) bool {
	// 根据不同的模型类型检查相应的指标
	switch condition.ModelType {
	case model.ModelYinYang:
		return metrics.YinYang.Balance >= condition.Threshold
	case model.ModelWuXing:
		return metrics.WuXing.Balance >= condition.Threshold
	case model.ModelBaGua:
		return metrics.BaGua.Stability >= condition.Threshold
	case model.ModelGanZhi:
		return metrics.GanZhi.Alignment >= condition.Threshold
	default:
		return true
	}
}

// recordError 记录错误
func (d *Detector) recordError(err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.status.errors = append(d.status.errors, err)
	d.status.errorCount++
}

// GetAlertChannel 获取告警通道
func (d *Detector) GetAlertChannel() <-chan types.AlertData {
	return d.alertChan
}

// generateAlertID 生成告警ID
func generateAlertID() string {
	return fmt.Sprintf("alert-%d", time.Now().UnixNano())
}

// generateAlertMessage 生成告警消息
func generateAlertMessage(condition *AlertCondition, value float64) string {
	return fmt.Sprintf("%s: current value %.2f %s threshold %.2f",
		condition.Name, value, condition.Operator, condition.Threshold)
}
