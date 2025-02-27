// system/monitor/alert/handler.go

package alert

import (
	"context"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/types"
)

// Handler 告警处理器
type Handler struct {
	mu sync.RWMutex

	alerts   map[string]*AlertData  // 活跃告警
	handlers map[string]HandlerFunc // 处理器函数映射
	config   HandlerConfig          // 配置
	status   HandlerStatus          // 状态

	// 告警队列
	queue chan types.AlertData

	// 处理结果
	results chan HandlerResult
}

// AlertHandler 告警处理器
type AlertHandler struct {
	mu sync.RWMutex

	// 配置
	config struct {
		MaxConcurrent int           // 最大并发处理数
		Timeout       time.Duration // 处理超时
		RetryCount    int           // 重试次数
		QueueSize     int           // 队列大小
	}

	// 处理器注册表
	handlers map[string]HandlerFunc

	// 告警队列 - 使用 types.AlertData
	queue chan types.AlertData

	// 处理状态
	status struct {
		isRunning    bool
		activeCount  int
		totalHandled int64
		lastError    error
		errors       []error
	}

	// 处理结果
	results chan HandlerResult

	// 模型状态
	modelState model.ModelState

	// 状态更新回调
	onStatusUpdate func(types.SystemEvent)
}

// HandlerFunc 告警处理函数类型
type HandlerFunc func(context.Context, types.AlertData) error

// HandlerResult 处理结果
type HandlerResult struct {
	AlertData  types.AlertData
	ModelState model.ModelState
	Handler    string
	Status     string
	Error      error
	StartTime  time.Time
	EndTime    time.Time
	Duration   time.Duration
	RetryCount int
}

// AlertData 告警数据
type AlertData struct {
	ID        string                 // 告警ID
	Type      string                 // 告警类型
	Level     types.AlertLevel       // 告警级别
	Message   string                 // 告警消息
	Source    string                 // 告警源
	Time      time.Time              // 发生时间
	Status    string                 // 状态
	Details   map[string]interface{} // 详细信息
	ModelData *types.ModelAlertData  // 模型数据
}

// HandlerConfig 处理器配置
type HandlerConfig struct {
	MaxRetries    int           // 最大重试次数
	RetryDelay    time.Duration // 重试间隔
	MaxQueueSize  int           // 最大队列大小
	MaxConcurrent int           // 最大并发处理数
}

// HandlerStatus 处理器状态
type HandlerStatus struct {
	StartTime    time.Time // 启动时间
	IsRunning    bool      // 是否运行中
	ActiveCount  int       // 活跃处理数
	TotalHandled int64     // 总处理数
	LastError    error     // 最后错误
	ErrorCount   int       // 错误数量
}

// ------------------------------------------------------------
// NewHandler 创建新的告警处理器
func NewHandler() *Handler {
	return &Handler{
		alerts:   make(map[string]*AlertData),
		handlers: make(map[string]HandlerFunc),
		config: HandlerConfig{
			MaxRetries:   3,
			RetryDelay:   time.Second * 5,
			MaxQueueSize: 1000,
		},
		status: HandlerStatus{
			StartTime: time.Now(),
		},
	}
}

// NewAlertHandler 创建新的告警处理器
func NewAlertHandler(config types.AlertConfig) *AlertHandler {
	h := &AlertHandler{
		handlers: make(map[string]HandlerFunc),
		queue:    make(chan types.AlertData, config.QueueSize),
		results:  make(chan HandlerResult, config.QueueSize),
	}

	// 设置配置
	h.config.MaxConcurrent = config.MaxConcurrent
	h.config.Timeout = config.Timeout
	h.config.RetryCount = config.RetryCount
	h.config.QueueSize = config.QueueSize

	// 注册默认处理器
	h.registerDefaultHandlers()

	return h
}

// Start 启动处理器
func (h *AlertHandler) Start(ctx context.Context) error {
	h.mu.Lock()
	if h.status.isRunning {
		h.mu.Unlock()
		return model.WrapError(nil, model.ErrCodeOperation, "handler already running")
	}
	h.status.isRunning = true
	h.mu.Unlock()

	// 启动处理循环
	for i := 0; i < h.config.MaxConcurrent; i++ {
		go h.processLoop(ctx)
	}

	return nil
}

// Stop 停止处理器
func (h *AlertHandler) Stop() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.status.isRunning {
		return model.WrapError(nil, model.ErrCodeOperation, "handler not running")
	}

	h.status.isRunning = false
	return nil
}

// Handle 处理告警
func (h *AlertHandler) Handle(alert types.AlertData) error {
	if !h.status.isRunning {
		return model.WrapError(nil, model.ErrCodeOperation, "handler not running")
	}

	select {
	case h.queue <- alert:
		return nil
	default:
		return model.WrapError(nil, model.ErrCodeResource, "alert queue full")
	}
}

// processLoop 处理循环
func (h *AlertHandler) processLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case alert := <-h.queue:
			h.handleAlert(ctx, alert)
		}
	}
}

// handleAlert 处理单个告警
func (h *AlertHandler) handleAlert(ctx context.Context, alert types.AlertData) {
	h.mu.Lock()
	h.status.activeCount++
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		h.status.activeCount--
		h.status.totalHandled++
		h.mu.Unlock()
	}()

	// 创建处理上下文
	ctx, cancel := context.WithTimeout(ctx, h.config.Timeout)
	defer cancel()

	// 执行所有注册的处理器
	for name, handler := range h.handlers {
		result := HandlerResult{
			AlertData:  alert,
			ModelState: h.modelState,
			Handler:    name,
			StartTime:  time.Now(),
		}

		// 重试机制
		var err error
		for retry := 0; retry <= h.config.RetryCount; retry++ {
			result.RetryCount = retry

			if err = handler(ctx, alert); err == nil {
				break
			}

			// 检查上下文是否已取消
			if ctx.Err() != nil {
				err = ctx.Err()
				break
			}

			// 最后一次重试
			if retry == h.config.RetryCount {
				break
			}

			// 等待后重试
			time.Sleep(time.Second * time.Duration(retry+1))
		}

		// 记录结果
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		if err != nil {
			result.Status = "failed"
			result.Error = model.WrapError(err, model.ErrCodeOperation, "handler execution failed")
			h.recordError(result.Error)
		} else {
			result.Status = "success"
		}

		h.recordResult(result)
	}
}

// RegisterHandler 注册告警处理器
func (h *AlertHandler) RegisterHandler(name string, handler HandlerFunc) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 参数验证
	if name == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty handler name")
	}
	if handler == nil {
		return model.WrapError(nil, model.ErrCodeValidation, "nil handler function")
	}

	// 检查是否已存在
	if _, exists := h.handlers[name]; exists {
		return model.WrapError(nil, model.ErrCodeValidation, "handler already registered")
	}

	// 注册处理器
	h.handlers[name] = handler
	return nil
}

// registerDefaultHandlers 注册默认处理器
func (h *AlertHandler) registerDefaultHandlers() {
	// 日志处理器
	h.RegisterHandler("log", func(ctx context.Context, alert types.AlertData) error {
		// 创建日志记录
		logEntry := struct {
			Time    time.Time
			Level   string
			Type    string
			Message string
			Details map[string]interface{}
			ModelID string
			Metrics *types.AlertMetrics
		}{
			Time:    time.Now(),
			Level:   string(alert.Level),
			Type:    alert.Type,
			Message: alert.Message,
			Details: alert.Details,
			ModelID: alert.ModelID,
			Metrics: alert.Metrics,
		}

		// 序列化日志
		data, err := json.Marshal(logEntry)
		if err != nil {
			return model.WrapError(err, model.ErrCodeOperation, "failed to marshal log entry")
		}

		// 写入日志文件
		logPath := filepath.Join("logs", "alerts", time.Now().Format("2006-01-02")+".log")
		if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
			return model.WrapError(err, model.ErrCodeOperation, "failed to create log directory")
		}

		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return model.WrapError(err, model.ErrCodeOperation, "failed to open log file")
		}
		defer f.Close()

		if _, err := f.Write(append(data, '\n')); err != nil {
			return model.WrapError(err, model.ErrCodeOperation, "failed to write log")
		}

		return nil
	})

	// 状态更新处理器
	h.RegisterHandler("status", func(ctx context.Context, alert types.AlertData) error {
		h.mu.Lock()
		defer h.mu.Unlock()

		// 更新系统状态指标
		metrics := types.SystemMetrics{
			AlertCount:    h.status.totalHandled + 1,
			LastAlertTime: time.Now(),
			AlertLevels: map[types.AlertLevel]int{
				alert.Level: 1,
			},
		}

		// 基于告警级别更新健康度
		switch alert.Level {
		case types.AlertLevelCritical:
			metrics.Health = 0.2
		case types.AlertLevelError:
			metrics.Health = 0.4
		case types.AlertLevelWarning:
			metrics.Health = 0.6
		case types.AlertLevelInfo:
			metrics.Health = 0.8
		default:
			metrics.Health = 1.0
		}

		// 触发状态更新事件
		event := types.SystemEvent{
			Type:      "alert.status_update",
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"alert_type":  alert.Type,
				"alert_level": alert.Level,
				"metrics":     metrics,
			},
		}

		// 通知系统状态更新
		if h.onStatusUpdate != nil {
			h.onStatusUpdate(event)
		}

		return nil
	})

	// 元系统响应处理器
	h.RegisterHandler("meta", func(ctx context.Context, alert types.AlertData) error {
		// 更新模型状态
		if err := h.updateModelState(alert); err != nil {
			return model.WrapError(err, model.ErrCodeOperation, "failed to update model state")
		}

		// 根据告警类型执行相应的元系统响应
		switch alert.Type {
		case "energy_anomaly":
			if err := h.handleEnergyAnomaly(alert); err != nil {
				return model.WrapError(err, model.ErrCodeOperation, "failed to handle energy anomaly")
			}
		case "coherence_violation":
			if err := h.handleCoherenceViolation(alert); err != nil {
				return model.WrapError(err, model.ErrCodeOperation, "failed to handle coherence violation")
			}
		case "quantum_fluctuation":
			if err := h.handleQuantumFluctuation(alert); err != nil {
				return model.WrapError(err, model.ErrCodeOperation, "failed to handle quantum fluctuation")
			}
		}

		return nil
	})
}

// handleEnergyAnomaly 处理能量异常
func (h *AlertHandler) handleEnergyAnomaly(alert types.AlertData) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if alert.ModelData == nil || alert.ModelData.Metrics == nil {
		return model.WrapError(nil, model.ErrCodeValidation, "missing model metrics")
	}

	// 获取能量指标
	metrics := alert.ModelData.Metrics
	if metrics.Energy == (model.Energy{}) {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid energy metrics")
	}

	energy := metrics.Energy.Total
	threshold := metrics.Energy.Average + metrics.Energy.Variance*2

	// 根据能量异常类型采取相应措施
	if energy > threshold {
		// 能量过高 - 执行能量耗散
		newEnergy := metrics.Energy.Average * 0.8
		h.modelState.Energy = newEnergy
		h.modelState.Phase = model.PhaseTransform
		h.modelState.Nature = model.NatureUnstable
	} else if energy < metrics.Energy.Average/2 {
		// 能量过低 - 执行能量补充
		newEnergy := metrics.Energy.Average * 1.2
		h.modelState.Energy = newEnergy
		h.modelState.Phase = model.PhaseYang
		h.modelState.Nature = model.NatureNeutral
	}

	// 触发状态更新事件
	h.notifyStateChange("energy_adjustment", map[string]interface{}{
		"old_energy": energy,
		"new_energy": h.modelState.Energy,
		"threshold":  threshold,
	})

	return nil
}

// handleCoherenceViolation 处理相干性违规
func (h *AlertHandler) handleCoherenceViolation(alert types.AlertData) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if alert.ModelData == nil || alert.ModelData.Quantum == nil {
		return model.WrapError(nil, model.ErrCodeValidation, "missing quantum state")
	}

	// 获取量子状态
	quantum := alert.ModelData.Quantum
	coherence := quantum.GetCoherence()
	threshold := 0.6 // 相干性阈值

	// 根据相干性违规类型采取措施
	if coherence < threshold {
		// 相干性过低 - 执行量子态重整
		h.modelState.Phase = model.PhaseTransform
		h.modelState.Nature = model.NatureUnstable
		h.modelState.Energy *= 0.9 // 降低能量以稳定系统

		// 应用量子态调整
		adjustments := map[string]float64{
			"entanglement": math.Min(1.0, quantum.GetEntanglement()*1.2),
			"phase":        normalizePhase(quantum.GetPhase() + math.Pi/4),
			"stability":    quantum.GetStability() * 1.1,
		}

		// 转换为 interface{} map
		interfaceMap := make(map[string]interface{}, len(adjustments))
		for k, v := range adjustments {
			interfaceMap[k] = v
		}

		// 触发量子态更新事件
		h.notifyStateChange("quantum_adjustment", interfaceMap)
	}

	return nil
}

// handleQuantumFluctuation 处理量子波动
func (h *AlertHandler) handleQuantumFluctuation(alert types.AlertData) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if alert.ModelData == nil || alert.ModelData.Field == nil {
		return model.WrapError(nil, model.ErrCodeValidation, "missing field state")
	}

	// 获取场状态
	field := alert.ModelData.Field
	fluctuation := field.GetEnergyFlow()
	threshold := 0.3 // 波动阈值

	// 处理量子波动
	if math.Abs(fluctuation) > threshold {
		// 计算平衡因子
		balanceFactor := 1.0 / (1.0 + math.Abs(fluctuation))

		// 调整模型状态
		h.modelState.Energy *= balanceFactor
		h.modelState.Phase = model.PhaseYinYang // 使用 PhaseYinYang 常量
		h.modelState.Nature = model.NatureNeutral

		// 应用场调整
		adjustments := map[string]float64{
			"energy_flow":    fluctuation * balanceFactor,
			"field_strength": field.GetFieldStrength() * balanceFactor,
			"balance":        calculateFieldBalance(field),
		}

		// 转换为 interface{} map
		interfaceMap := make(map[string]interface{}, len(adjustments))
		for k, v := range adjustments {
			interfaceMap[k] = v
		}

		// 触发场态更新事件
		h.notifyStateChange("field_adjustment", interfaceMap)
	}

	return nil
}

// 辅助函数

// notifyStateChange 通知状态变化
func (h *AlertHandler) notifyStateChange(eventType string, data map[string]interface{}) {
	if h.onStatusUpdate == nil {
		return
	}

	// 转换事件类型
	var sysEventType types.EventType
	switch eventType {
	case "energy_adjustment":
		sysEventType = types.EventAlertEnergyAdjustment
	case "quantum_adjustment":
		sysEventType = types.EventAlertQuantumAdjustment
	case "field_adjustment":
		sysEventType = types.EventAlertFieldAdjustment
	default:
		sysEventType = types.EventAlertStatusUpdate
	}

	event := types.SystemEvent{
		Type:      sysEventType,
		Timestamp: time.Now(),
		Data:      data,
	}

	h.onStatusUpdate(event)
}

// normalizePhase 标准化相位到 [0, 2π)
func normalizePhase(phase float64) float64 {
	return math.Mod(phase+2*math.Pi, 2*math.Pi)
}

// calculateFieldBalance 计算场平衡度
func calculateFieldBalance(field *model.FieldState) float64 {
	if field == nil {
		return 0
	}
	return 1.0 / (1.0 + math.Abs(field.GetEnergyFlow()))
}

// updateModelState 更新模型状态
func (h *AlertHandler) updateModelState(alert types.AlertData) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 检查模型数据
	if alert.ModelData == nil {
		return nil // 无模型数据,不需要更新
	}

	// 准备新状态
	newState := h.modelState

	// 基于告警级别更新状态
	switch alert.Level {
	case types.AlertLevelCritical:
		newState.Health = 0.2
		newState.Nature = model.NatureUnstable
	case types.AlertLevelError:
		newState.Health = 0.4
		newState.Nature = model.NatureUnstable
	case types.AlertLevelWarning:
		newState.Health = 0.6
		newState.Nature = model.NatureNeutral
	default:
		newState.Health = 0.8
		newState.Nature = model.NatureStable
	}

	// 根据告警类型更新特定字段
	switch alert.Type {
	case "energy_low":
		newState.Energy = alert.ModelData.Metrics.Energy.Total * 0.8
	case "field_high":
		newState.Phase = model.PhaseTransform
	case "coherence_low":
		newState.Nature = model.NatureUnstable
	case "quantum_anomaly":
		newState.Phase = model.PhaseYin
	}

	// 更新时间戳
	newState.UpdateTime = time.Now()

	// 应用新状态
	h.modelState = newState

	return nil
}

// recordResult 记录处理结果
func (h *AlertHandler) recordResult(result HandlerResult) {
	select {
	case h.results <- result:
	default:
		h.recordError(model.WrapError(nil, model.ErrCodeResource, "result buffer full"))
	}
}

// recordError 记录错误
func (h *AlertHandler) recordError(err error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.status.lastError = err
	h.status.errors = append(h.status.errors, err)
}

// GetStatus 获取处理器状态
func (h *AlertHandler) GetStatus() types.HandlerStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return types.HandlerStatus{
		IsRunning:    h.status.isRunning,
		ActiveCount:  h.status.activeCount,
		TotalHandled: h.status.totalHandled,
		LastError:    h.status.lastError,
		ErrorCount:   len(h.status.errors),
	}
}

// GetResults 获取处理结果通道
func (h *AlertHandler) GetResults() <-chan HandlerResult {
	return h.results
}

// SetStatusUpdateCallback 设置状态更新回调函数
func (h *AlertHandler) SetStatusUpdateCallback(callback func(types.SystemEvent)) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.onStatusUpdate = callback
}

// handleAlert 处理单个告警
func (h *Handler) handleAlert(ctx context.Context, alert types.AlertData) {
	h.mu.Lock()
	h.status.ActiveCount++
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		h.status.ActiveCount--
		h.status.TotalHandled++
		h.mu.Unlock()
	}()

	// 创建处理结果
	result := HandlerResult{
		AlertData: alert,
		StartTime: time.Now(),
	}

	// 执行处理器
	handler, exists := h.handlers[alert.Type]
	if !exists {
		result.Status = "failed"
		result.Error = model.WrapError(nil, model.ErrCodeValidation, "no handler for alert type")
		h.recordResult(result)
		return
	}

	// 带重试的处理
	var err error
	for retry := 0; retry < h.config.MaxRetries; retry++ {
		// 检查上下文是否取消
		select {
		case <-ctx.Done():
			result.Status = "cancelled"
			result.Error = ctx.Err()
			h.recordResult(result)
			return
		default:
		}

		err = handler(ctx, alert)
		if err == nil {
			break
		}

		if retry < h.config.MaxRetries-1 {
			time.Sleep(h.config.RetryDelay)
		}
	}

	// 更新结果
	result.EndTime = time.Now()
	if err != nil {
		result.Status = "failed"
		result.Error = err
		h.recordError(err)
	} else {
		result.Status = "success"
	}

	h.recordResult(result)
}

// recordResult 记录处理结果
func (h *Handler) recordResult(result HandlerResult) {
	select {
	case h.results <- result:
	default:
		h.recordError(model.WrapError(nil, model.ErrCodeResource, "result buffer full"))
	}
}

// recordError 记录错误
func (h *Handler) recordError(err error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.status.LastError = err
	h.status.ErrorCount++
}

// processLoop 处理循环
func (h *Handler) processLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case alert := <-h.queue:
			h.handleAlert(ctx, alert)
		}
	}
}

// Start 启动处理器
func (h *Handler) Start(ctx context.Context) error {
	h.mu.Lock()
	if h.status.IsRunning {
		h.mu.Unlock()
		return model.WrapError(nil, model.ErrCodeOperation, "handler already running")
	}
	h.status.IsRunning = true
	h.mu.Unlock()

	// 创建处理队列
	h.queue = make(chan types.AlertData, h.config.MaxQueueSize)
	h.results = make(chan HandlerResult, h.config.MaxQueueSize)

	// 启动处理循环
	for i := 0; i < h.config.MaxConcurrent; i++ {
		go h.processLoop(ctx)
	}

	return nil
}

// Stop 停止处理器
func (h *Handler) Stop() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.status.IsRunning {
		return model.WrapError(nil, model.ErrCodeOperation, "handler not running")
	}

	// 关闭通道
	close(h.queue)
	close(h.results)

	// 更新状态
	h.status.IsRunning = false

	return nil
}
