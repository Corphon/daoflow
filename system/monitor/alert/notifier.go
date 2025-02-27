// system/monitor/alert/notifier.go

package alert

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Corphon/daoflow/system/types"
)

// NotificationChannel 通知渠道类型
type NotificationChannel string

const (
	ChannelEmail   NotificationChannel = "email"
	ChannelWebhook NotificationChannel = "webhook"
	ChannelMessage NotificationChannel = "message"
	ChannelConsole NotificationChannel = "console"
	ChannelLog     NotificationChannel = "log"
)

// NotificationTarget 通知目标配置
type NotificationTarget struct {
	ID       string              // 目标ID
	Name     string              // 目标名称
	Channel  NotificationChannel // 通知渠道
	Config   map[string]string   // 渠道配置
	Filters  []string            // 告警过滤器
	Template string              // 消息模板
	Enabled  bool                // 是否启用
}

// NotificationResult 通知结果
type NotificationResult struct {
	ID        string    // 结果ID
	TargetID  string    // 目标ID
	AlertID   string    // 告警ID
	Status    string    // 状态
	Error     error     // 错误信息
	Timestamp time.Time // 时间戳
}

// Notifier 告警通知器
type Notifier struct {
	mu sync.RWMutex

	// 配置
	config struct {
		RetryInterval time.Duration
		MaxRetries    int
		BatchSize     int
		QueueSize     int
	}

	// 通知目标管理
	targets map[string]*NotificationTarget

	// 通知渠道处理器
	channels map[NotificationChannel]NotificationHandler

	// 通知队列
	queue   chan types.Alert
	results chan NotificationResult

	// 状态
	status struct {
		isRunning  bool
		lastNotify time.Time
		errorCount int
		errors     []error
	}
}

// NotificationHandler 通知处理器接口
type NotificationHandler interface {
	Send(ctx context.Context, target *NotificationTarget, alert types.Alert) error
	Validate(config map[string]string) error
}

// ---------------------------------------------------
// NewNotifier 创建新的通知器
func NewNotifier(config types.AlertConfig) *Notifier {
	n := &Notifier{
		targets:  make(map[string]*NotificationTarget),
		channels: make(map[NotificationChannel]NotificationHandler),
		queue:    make(chan types.Alert, config.QueueSize),
		results:  make(chan NotificationResult, config.QueueSize),
	}

	// 设置配置
	n.config.RetryInterval = config.RetryInterval
	n.config.MaxRetries = config.MaxRetries
	n.config.BatchSize = config.BatchSize
	n.config.QueueSize = config.QueueSize

	// 注册默认通知渠道
	n.registerDefaultChannels()

	return n
}

// Start 启动通知器
func (n *Notifier) Start(ctx context.Context) error {
	n.mu.Lock()
	if n.status.isRunning {
		n.mu.Unlock()
		return types.NewSystemError(types.ErrRuntime, "notifier already running", nil)
	}
	n.status.isRunning = true
	n.mu.Unlock()

	// 启动通知处理循环
	go n.processLoop(ctx)

	return nil
}

// Stop 停止通知器
func (n *Notifier) Stop() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.status.isRunning = false
	return nil
}

// AddTarget 添加通知目标
func (n *Notifier) AddTarget(target *NotificationTarget) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	// 验证目标配置
	if handler, exists := n.channels[target.Channel]; exists {
		if err := handler.Validate(target.Config); err != nil {
			return types.WrapError(err, types.ErrInvalid, "invalid target configuration")
		}
	} else {
		return types.NewSystemError(types.ErrNotFound, "notification channel not found", nil)
	}

	n.targets[target.ID] = target
	return nil
}

// RemoveTarget 移除通知目标
func (n *Notifier) RemoveTarget(id string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if _, exists := n.targets[id]; !exists {
		return types.NewSystemError(types.ErrNotFound, "target not found", nil)
	}

	delete(n.targets, id)
	return nil
}

// RegisterChannel 注册通知渠道
func (n *Notifier) RegisterChannel(channel NotificationChannel, handler NotificationHandler) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.channels[channel] = handler
}

// Notify 发送告警通知
func (n *Notifier) Notify(alert types.Alert) error {
	if !n.status.isRunning {
		return types.NewSystemError(types.ErrRuntime, "notifier not running", nil)
	}

	select {
	case n.queue <- alert:
		return nil
	default:
		return types.NewSystemError(types.ErrOverflow, "notification queue full", nil)
	}
}

// processLoop 通知处理循环
func (n *Notifier) processLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case alert := <-n.queue:
			n.processAlert(ctx, alert)
		}
	}
}

// processAlert 处理单个告警
func (n *Notifier) processAlert(ctx context.Context, alert types.Alert) {
	n.mu.RLock()
	targets := make([]*NotificationTarget, 0)
	for _, target := range n.targets {
		if target.Enabled && n.shouldNotify(target, alert) {
			targets = append(targets, target)
		}
	}
	n.mu.RUnlock()

	// 并发发送通知
	var wg sync.WaitGroup
	for _, target := range targets {
		wg.Add(1)
		go func(t *NotificationTarget) {
			defer wg.Done()
			n.sendNotification(ctx, t, alert)
		}(target)
	}
	wg.Wait()
}

// shouldNotify 检查是否应该通知目标
func (n *Notifier) shouldNotify(target *NotificationTarget, alert types.Alert) bool {
	// 检查过滤器
	for _, filter := range target.Filters {
		if !matchFilter(filter, alert) {
			return false
		}
	}
	return true
}

// sendNotification 发送通知
func (n *Notifier) sendNotification(ctx context.Context, target *NotificationTarget, alert types.Alert) {
	handler, exists := n.channels[target.Channel]
	if !exists {
		n.recordError(types.NewSystemError(types.ErrNotFound, "channel handler not found", nil))
		return
	}

	// 重试机制
	for retry := 0; retry < n.config.MaxRetries; retry++ {
		err := handler.Send(ctx, target, alert)
		if err == nil {
			n.recordResult(NotificationResult{
				ID:        generateResultID(),
				TargetID:  target.ID,
				AlertID:   alert.ID,
				Status:    "success",
				Timestamp: time.Now(),
			})
			return
		}

		// 记录错误但继续重试
		n.recordError(err)
		time.Sleep(n.config.RetryInterval)
	}

	// 最终失败
	n.recordResult(NotificationResult{
		ID:        generateResultID(),
		TargetID:  target.ID,
		AlertID:   alert.ID,
		Status:    "failed",
		Error:     types.NewSystemError(types.ErrRuntime, "max retries exceeded", nil),
		Timestamp: time.Now(),
	})
}

// ConsoleNotifier 控制台通知处理器
type ConsoleNotifier struct{}

// Send 发送通知到控制台
func (n *ConsoleNotifier) Send(ctx context.Context, target *NotificationTarget, alert types.Alert) error {
	fmt.Printf("[%s] %s: %s\n", alert.Level, target.Name, alert.Message)
	return nil
}

// Validate 验证配置
func (n *ConsoleNotifier) Validate(config map[string]string) error {
	return nil // 控制台通知不需要配置
}

// LogNotifier 日志通知处理器
type LogNotifier struct{}

// Send 发送通知到日志
func (n *LogNotifier) Send(ctx context.Context, target *NotificationTarget, alert types.Alert) error {
	log.Printf("[%s] Target: %s, Alert: %s\n", alert.Level, target.Name, alert.Message)
	return nil
}

// Validate 验证配置
func (n *LogNotifier) Validate(config map[string]string) error {
	return nil // 默认日志通知不需要配置
}

// registerDefaultChannels 注册默认通知渠道
func (n *Notifier) registerDefaultChannels() {
	// 注册控制台通知处理器
	n.RegisterChannel(ChannelConsole, &ConsoleNotifier{})

	// 注册日志通知处理器
	n.RegisterChannel(ChannelLog, &LogNotifier{})
}

// recordResult 记录通知结果
func (n *Notifier) recordResult(result NotificationResult) {
	select {
	case n.results <- result:
	default:
		n.recordError(types.NewSystemError(types.ErrOverflow, "result buffer full", nil))
	}
}

// recordError 记录错误
func (n *Notifier) recordError(err error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.status.errors = append(n.status.errors, err)
	n.status.errorCount++
}

// generateResultID 生成结果ID
func generateResultID() string {
	return fmt.Sprintf("result-%d", time.Now().UnixNano())
}

// matchFilter 匹配过滤器
func matchFilter(filter string, alert types.Alert) bool {
	// 解析过滤器表达式
	parts := strings.Split(filter, ":")
	if len(parts) != 2 {
		return false
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	// 基于不同的过滤键进行匹配
	switch key {
	case "level":
		// 级别匹配
		return string(alert.Level) == value

	case "type":
		// 类型匹配
		return alert.Type == value

	case "source":
		// 来源匹配
		return alert.Source == value

	case "label":
		// 标签匹配
		if labelParts := strings.Split(value, "="); len(labelParts) == 2 {
			labelKey := strings.TrimSpace(labelParts[0])
			labelValue := strings.TrimSpace(labelParts[1])
			actualValue, exists := alert.Labels[labelKey]
			return exists && actualValue == labelValue
		}

	case "status":
		// 状态匹配
		return alert.Status == value
	}

	return false
}

// WebhookNotifier Webhook通知处理器
type WebhookNotifier struct{}

// Send 发送通知到Webhook
func (n *WebhookNotifier) Send(ctx context.Context, target *NotificationTarget, alert types.Alert) error {
	// TODO: 实现webhook通知逻辑
	return nil
}

// Validate 验证配置
func (n *WebhookNotifier) Validate(config map[string]string) error {
	if _, exists := config["url"]; !exists {
		return types.NewSystemError(types.ErrInvalid, "webhook URL not configured", nil)
	}
	return nil
}

// MessageNotifier 消息通知处理器
type MessageNotifier struct{}

// Send 发送通知消息
func (n *MessageNotifier) Send(ctx context.Context, target *NotificationTarget, alert types.Alert) error {
	// TODO: 实现消息通知逻辑
	return nil
}

// Validate 验证配置
func (n *MessageNotifier) Validate(config map[string]string) error {
	if _, exists := config["endpoint"]; !exists {
		return types.NewSystemError(types.ErrInvalid, "message endpoint not configured", nil)
	}
	return nil
}

// NewAlertNotifier 创建新的告警通知器
func NewAlertNotifier(config types.AlertConfig) *types.AlertNotifier {
	notifier := &types.AlertNotifier{}

	// 初始化配置
	notifier.Config.Channels = config.Channels
	notifier.Config.Templates = config.Templates
	notifier.Config.Throttling = config.RetryInterval

	// 初始化状态
	notifier.State.PendingNotifications = make([]types.Alert, 0)
	notifier.State.SentNotifications = make([]types.Alert, 0)
	notifier.State.LastNotification = time.Time{}

	// 创建并启动通知处理器
	baseNotifier := NewNotifier(config)

	// 注册默认通知渠道
	for _, channel := range config.Channels {
		switch channel {
		case string(ChannelConsole):
			baseNotifier.RegisterChannel(ChannelConsole, &ConsoleNotifier{})
		case string(ChannelLog):
			baseNotifier.RegisterChannel(ChannelLog, &LogNotifier{})
		case string(ChannelWebhook):
			baseNotifier.RegisterChannel(ChannelWebhook, &WebhookNotifier{})
		case string(ChannelMessage):
			baseNotifier.RegisterChannel(ChannelMessage, &MessageNotifier{})
		}
	}

	// 启动通知处理
	go func() {
		ctx := context.Background()
		if err := baseNotifier.Start(ctx); err != nil {
			log.Printf("Failed to start notifier: %v", err)
		}
	}()

	return notifier
}
