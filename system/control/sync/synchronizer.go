//system/control/sync/synchronizer.go

package sync

import (
    "sync"
    "time"

    "github.com/Corphon/daoflow/meta/field"
    "github.com/Corphon/daoflow/model"
    "github.com/Corphon/daoflow/system/types"
)

// Synchronizer 同步器
type Synchronizer struct {
    mu sync.RWMutex

    // 基础配置
    config struct {
        syncInterval    time.Duration // 同步间隔
        batchSize      int           // 批处理大小
        retryLimit     int           // 重试限制
        maxQueueSize   int           // 最大队列大小
    }

    // 同步状态
    state struct {
        tasks       map[string]*SyncTask      // 同步任务
        queues      map[string]*SyncQueue     // 同步队列
        operations  map[string]*SyncOperation // 同步操作
        metrics     SyncMetrics              // 同步指标
    }
}

// SyncTask 同步任务
type SyncTask struct {
    ID           string                // 任务ID
    Type         string                // 任务类型
    Source       *SyncEndpoint        // 源端点
    Target       *SyncEndpoint        // 目标端点
    State        string                // 任务状态
    Priority     int                   // 优先级
    Schedule     *SyncSchedule        // 同步计划
    LastSync     time.Time            // 最后同步
}

// SyncEndpoint 同步端点
type SyncEndpoint struct {
    ID           string                // 端点ID
    Type         string                // 端点类型
    Location     string                // 端点位置
    Properties   map[string]interface{} // 端点属性
    State        *EndpointState        // 端点状态
}

// EndpointState 端点状态
type EndpointState struct {
    Status       string                // 状态标识
    Version      string                // 数据版本
    LastUpdate   time.Time            // 最后更新
    Checksum     string                // 数据校验
}

// SyncSchedule 同步计划
type SyncSchedule struct {
    Type         string                // 计划类型
    Interval     time.Duration         // 同步间隔
    TimeWindow   *TimeWindow           // 时间窗口
    Conditions   []SyncCondition       // 同步条件
}

// TimeWindow 时间窗口
type TimeWindow struct {
    Start        time.Time            // 开始时间
    End          time.Time            // 结束时间
    Recurrence   string               // 重复规则
}

// SyncCondition 同步条件
type SyncCondition struct {
    Type         string                // 条件类型
    Expression   string                // 条件表达式
    Parameters   map[string]interface{} // 条件参数
}

// SyncQueue 同步队列
type SyncQueue struct {
    ID           string                // 队列ID
    Type         string                // 队列类型
    Priority     int                   // 优先级
    Tasks        []*SyncTask           // 任务列表
    Stats        QueueStats            // 队列统计
}

// QueueStats 队列统计
type QueueStats struct {
    TotalTasks    int                 // 总任务数
    Pending       int                 // 等待任务
    Processing    int                 // 处理中任务
    Completed     int                 // 完成任务
    Failed        int                 // 失败任务
}

// SyncOperation 同步操作
type SyncOperation struct {
    ID           string                // 操作ID
    TaskID       string                // 任务ID
    Type         string                // 操作类型
    Status       string                // 操作状态
    StartTime    time.Time            // 开始时间
    EndTime      time.Time            // 结束时间
    Result       *OperationResult      // 操作结果
}

// OperationResult 操作结果
type OperationResult struct {
    Success      bool                  // 是否成功
    SyncedItems  int                  // 同步项数
    Errors       []string              // 错误信息
    Checksum     string                // 校验和
}

// SyncMetrics 同步指标
type SyncMetrics struct {
    ActiveTasks    int                // 活跃任务数
    QueuedTasks    int                // 队列任务数
    SyncRate       float64            // 同步速率
    SuccessRate    float64            // 成功率
    AverageLatency time.Duration      // 平均延迟
    History        []MetricPoint      // 历史指标
}

// NewSynchronizer 创建新的同步器
func NewSynchronizer() *Synchronizer {
    s := &Synchronizer{}

    // 初始化配置
    s.config.syncInterval = 5 * time.Second
    s.config.batchSize = 100
    s.config.retryLimit = 3
    s.config.maxQueueSize = 1000

    // 初始化状态
    s.state.tasks = make(map[string]*SyncTask)
    s.state.queues = make(map[string]*SyncQueue)
    s.state.operations = make(map[string]*SyncOperation)
    s.state.metrics = SyncMetrics{
        History: make([]MetricPoint, 0),
    }

    return s
}

// RegisterTask 注册同步任务
func (s *Synchronizer) RegisterTask(task *SyncTask) error {
    if task == nil {
        return model.WrapError(nil, model.ErrCodeValidation, "nil task")
    }

    s.mu.Lock()
    defer s.mu.Unlock()

    // 验证任务
    if err := s.validateTask(task); err != nil {
        return err
    }

    // 存储任务
    s.state.tasks[task.ID] = task

    // 添加到相应队列
    if err := s.enqueueTask(task); err != nil {
        delete(s.state.tasks, task.ID)
        return err
    }

    return nil
}

// Synchronize 执行同步
func (s *Synchronizer) Synchronize() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    // 处理队列中的任务
    for _, queue := range s.state.queues {
        if err := s.processQueue(queue); err != nil {
            continue
        }
    }

    // 更新指标
    s.updateMetrics()

    return nil
}

// processQueue 处理同步队列
func (s *Synchronizer) processQueue(queue *SyncQueue) error {
    // 检查队列状态
    if len(queue.Tasks) == 0 {
        return nil
    }

    // 按批次处理任务
    for i := 0; i < len(queue.Tasks); i += s.config.batchSize {
        end := i + s.config.batchSize
        if end > len(queue.Tasks) {
            end = len(queue.Tasks)
        }

        batch := queue.Tasks[i:end]
        if err := s.processBatch(batch); err != nil {
            continue
        }
    }

    return nil
}

// processBatch 处理任务批次
func (s *Synchronizer) processBatch(tasks []*SyncTask) error {
    for _, task := range tasks {
        // 创建同步操作
        operation := &SyncOperation{
            ID:        generateOperationID(),
            TaskID:    task.ID,
            Type:      task.Type,
            Status:    "started",
            StartTime: time.Now(),
        }

        // 执行同步操作
        if err := s.executeOperation(operation); err != nil {
            s.handleOperationError(operation, err)
            continue
        }

        // 更新任务状态
        s.updateTaskState(task, operation)
    }

    return nil
}

// 辅助函数

func (s *Synchronizer) validateTask(task *SyncTask) error {
    if task.ID == "" {
        return model.WrapError(nil, model.ErrCodeValidation, "empty task ID")
    }

    if task.Source == nil || task.Target == nil {
        return model.WrapError(nil, model.ErrCodeValidation, "invalid endpoints")
    }

    return nil
}

func (s *Synchronizer) updateMetrics() {
    point := MetricPoint{
        Timestamp: time.Now(),
        Values: map[string]float64{
            "active_tasks":    float64(len(s.state.tasks)),
            "queued_tasks":    float64(s.state.metrics.QueuedTasks),
            "success_rate":    s.state.metrics.SuccessRate,
            "sync_rate":       s.state.metrics.SyncRate,
        },
    }

    s.state.metrics.History = append(s.state.metrics.History, point)

    // 限制历史记录数量
    if len(s.state.metrics.History) > maxMetricsHistory {
        s.state.metrics.History = s.state.metrics.History[1:]
    }
}

func generateOperationID() string {
    return fmt.Sprintf("op_%d", time.Now().UnixNano())
}

const (
    maxMetricsHistory = 1000
)
