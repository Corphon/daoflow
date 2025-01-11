// system/monitor/alert/handler.go

package alert

import (
    "context"
    "sync"
    "time"

    "github.com/Corphon/daoflow/system/types"
)

// AlertHandler 告警处理器
type AlertHandler struct {
    mu sync.RWMutex

    // 配置
    config struct {
        MaxConcurrent int           // 最大并发处理数
        Timeout       time.Duration // 处理超时
        RetryCount    int          // 重试次数
        QueueSize     int          // 队列大小
    }

    // 处理器注册表
    handlers map[string]HandlerFunc

    // 告警队列
    queue chan types.Alert

    // 处理状态
    status struct {
        isRunning    bool
        activeCount  int
        totalHandled int64
        lastError    error
        errors      []error
    }

    // 处理结果
    results chan HandlerResult
}

// HandlerFunc 告警处理函数类型
type HandlerFunc func(context.Context, types.Alert) error

// HandlerResult 处理结果
type HandlerResult struct {
    AlertID    string
    Handler    string
    Status     string
    Error      error
    StartTime  time.Time
    EndTime    time.Time
    Duration   time.Duration
    RetryCount int
}

// NewAlertHandler 创建新的告警处理器
func NewAlertHandler(config types.AlertConfig) *AlertHandler {
    h := &AlertHandler{
        handlers: make(map[string]HandlerFunc),
        queue:    make(chan types.Alert, config.QueueSize),
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
        return types.NewSystemError(types.ErrRuntime, "handler already running", nil)
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

    h.status.isRunning = false
    return nil
}

// RegisterHandler 注册告警处理器
func (h *AlertHandler) RegisterHandler(name string, handler HandlerFunc) {
    h.mu.Lock()
    defer h.mu.Unlock()

    h.handlers[name] = handler
}

// Handle 处理告警
func (h *AlertHandler) Handle(alert types.Alert) error {
    if !h.status.isRunning {
        return types.NewSystemError(types.ErrRuntime, "handler not running", nil)
    }

    select {
    case h.queue <- alert:
        return nil
    default:
        return types.NewSystemError(types.ErrOverflow, "alert queue full", nil)
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
func (h *AlertHandler) handleAlert(ctx context.Context, alert types.Alert) {
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
            AlertID:   alert.ID,
            Handler:   name,
            StartTime: time.Now(),
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
            result.Error = err
            h.recordError(err)
        } else {
            result.Status = "success"
        }

        h.recordResult(result)
    }
}

// registerDefaultHandlers 注册默认处理器
func (h *AlertHandler) registerDefaultHandlers() {
    // 日志处理器
    h.RegisterHandler("log", func(ctx context.Context, alert types.Alert) error {
        // TODO: 实现日志记录逻辑
        return nil
    })

    // 状态更新处理器
    h.RegisterHandler("status", func(ctx context.Context, alert types.Alert) error {
        // TODO: 实现状态更新逻辑
        return nil
    })

    // 自动恢复处理器
    h.RegisterHandler("recovery", func(ctx context.Context, alert types.Alert) error {
        // TODO: 实现自动恢复逻辑
        return nil
    })

    // 元系统响应处理器
    h.RegisterHandler("meta", func(ctx context.Context, alert types.Alert) error {
        // TODO: 实现元系统响应逻辑
        return nil
    })
}

// recordResult 记录处理结果
func (h *AlertHandler) recordResult(result HandlerResult) {
    select {
    case h.results <- result:
    default:
        h.recordError(types.NewSystemError(types.ErrOverflow, "result buffer full", nil))
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
func (h *AlertHandler) GetStatus() struct {
    IsRunning    bool
    ActiveCount  int
    TotalHandled int64
    LastError    error
    ErrorCount   int
} {
    h.mu.RLock()
    defer h.mu.RUnlock()

    return struct {
        IsRunning    bool
        ActiveCount  int
        TotalHandled int64
        LastError    error
        ErrorCount   int
    }{
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
