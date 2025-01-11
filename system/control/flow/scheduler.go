//system/control/flow/scheduler.go

package flow

import (
    "sync"
    "time"

    "github.com/Corphon/daoflow/meta/field"
    "github.com/Corphon/daoflow/model"
    "github.com/Corphon/daoflow/system/types"
)

// Scheduler 调度器
type Scheduler struct {
    mu sync.RWMutex

    // 基础配置
    config struct {
        maxTasks         int           // 最大任务数
        taskTimeout      time.Duration // 任务超时时间
        retryLimit       int           // 重试限制
        priorityLevels   int           // 优先级层级
    }

    // 调度状态
    state struct {
        tasks       map[string]*Task          // 任务列表
        queues      map[int]*TaskQueue       // 优先级队列
        executors   map[string]*TaskExecutor  // 任务执行器
        history     []TaskHistory            // 任务历史
        metrics     SchedulerMetrics         // 调度指标
    }

    // 依赖项
    backpressure *BackpressureManager
}

// Task 任务
type Task struct {
    ID           string                // 任务ID
    Type         string                // 任务类型
    Priority     int                   // 优先级
    Status       string                // 任务状态
    Parameters   map[string]interface{} // 任务参数
    Dependencies []string              // 依赖任务
    Created      time.Time            // 创建时间
    Deadline     time.Time            // 截止时间
    Retries      int                  // 重试次数
}

// TaskQueue 任务队列
type TaskQueue struct {
    Priority     int                   // 优先级
    Tasks        []*Task               // 任务列表
    Capacity     int                   // 队列容量
    Stats        QueueStats            // 队列统计
}

// QueueStats 队列统计
type QueueStats struct {
    TotalTasks    int                 // 总任务数
    CompletedTasks int                // 完成任务数
    FailedTasks   int                 // 失败任务数
    AverageWait   time.Duration       // 平均等待时间
}

// TaskExecutor 任务执行器
type TaskExecutor struct {
    ID           string                // 执行器ID
    Type         string                // 执行器类型
    Status       string                // 执行器状态
    Capacity     int                   // 处理容量
    CurrentLoad  int                   // 当前负载
    Performance  ExecutorPerformance   // 性能指标
}

// ExecutorPerformance 执行器性能
type ExecutorPerformance struct {
    SuccessRate  float64              // 成功率
    Throughput   float64              // 吞吐量
    Latency      time.Duration        // 处理延迟
    ErrorRate    float64              // 错误率
}

// TaskHistory 任务历史
type TaskHistory struct {
    TaskID       string                // 任务ID
    Type         string                // 记录类型
    Status       string                // 任务状态
    Details      map[string]interface{} // 详细信息
    Timestamp    time.Time            // 记录时间
}

// SchedulerMetrics 调度指标
type SchedulerMetrics struct {
    ActiveTasks   int                 // 活跃任务数
    QueueLength   map[int]int         // 队列长度
    Throughput    float64             // 系统吞吐量
    LatencyStats  LatencyStatistics   // 延迟统计
    History       []MetricPoint       // 历史指标
}

// LatencyStatistics 延迟统计
type LatencyStatistics struct {
    Average      time.Duration        // 平均延迟
    P95          time.Duration        // 95分位延迟
    P99          time.Duration        // 99分位延迟
    Max          time.Duration        // 最大延迟
}

// NewScheduler 创建新的调度器
func NewScheduler(backpressure *BackpressureManager) *Scheduler {
    s := &Scheduler{
        backpressure: backpressure,
    }

    // 初始化配置
    s.config.maxTasks = 1000
    s.config.taskTimeout = 30 * time.Second
    s.config.retryLimit = 3
    s.config.priorityLevels = 5

    // 初始化状态
    s.state.tasks = make(map[string]*Task)
    s.state.queues = make(map[int]*TaskQueue)
    s.state.executors = make(map[string]*TaskExecutor)
    s.state.history = make([]TaskHistory, 0)
    s.state.metrics = SchedulerMetrics{
        QueueLength: make(map[int]int),
        History:    make([]MetricPoint, 0),
    }

    // 初始化优先级队列
    for i := 0; i < s.config.priorityLevels; i++ {
        s.state.queues[i] = &TaskQueue{
            Priority: i,
            Tasks:   make([]*Task, 0),
            Capacity: s.config.maxTasks / s.config.priorityLevels,
        }
    }

    return s
}

// Schedule 调度任务
func (s *Scheduler) Schedule(task *Task) error {
    if task == nil {
        return model.WrapError(nil, model.ErrCodeValidation, "nil task")
    }

    s.mu.Lock()
    defer s.mu.Unlock()

    // 验证任务
    if err := s.validateTask(task); err != nil {
        return err
    }

    // 检查系统负载
    if err := s.checkSystemLoad(); err != nil {
        return err
    }

    // 分配到相应队列
    if err := s.enqueueTask(task); err != nil {
        return err
    }

    // 记录任务
    s.recordTask(task, "scheduled", nil)

    return nil
}

// Execute 执行调度
func (s *Scheduler) Execute() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    // 检查执行器状态
    if err := s.checkExecutors(); err != nil {
        return err
    }

    // 处理每个优先级队列
    for priority := 0; priority < s.config.priorityLevels; priority++ {
        if err := s.processQueue(priority); err != nil {
            continue
        }
    }

    // 更新指标
    s.updateMetrics()

    return nil
}

// processQueue 处理任务队列
func (s *Scheduler) processQueue(priority int) error {
    queue := s.state.queues[priority]
    if queue == nil || len(queue.Tasks) == 0 {
        return nil
    }

    // 获取可用执行器
    executor := s.findAvailableExecutor()
    if executor == nil {
        return model.WrapError(nil, model.ErrCodeResource, "no available executor")
    }

    // 处理队列中的任务
    for i := 0; i < len(queue.Tasks); i++ {
        task := queue.Tasks[i]

        // 检查任务状态
        if !s.isTaskExecutable(task) {
            continue
        }

        // 执行任务
        if err := s.executeTask(executor, task); err != nil {
            s.handleTaskError(task, err)
            continue
        }

        // 更新队列
        queue.Tasks = append(queue.Tasks[:i], queue.Tasks[i+1:]...)
        i--
    }

    return nil
}

// 辅助函数

func (s *Scheduler) validateTask(task *Task) error {
    if task.ID == "" {
        return model.WrapError(nil, model.ErrCodeValidation, "empty task ID")
    }

    if task.Priority >= s.config.priorityLevels {
        return model.WrapError(nil, model.ErrCodeValidation, "invalid priority level")
    }

    return nil
}

func (s *Scheduler) recordTask(
    task *Task,
    recordType string,
    details map[string]interface{}) {
    
    record := TaskHistory{
        TaskID:    task.ID,
        Type:      recordType,
        Status:    task.Status,
        Details:   details,
        Timestamp: time.Now(),
    }

    s.state.history = append(s.state.history, record)

    // 限制历史记录数量
    if len(s.state.history) > maxHistoryLength {
        s.state.history = s.state.history[1:]
    }
}

const (
    maxHistoryLength = 1000
)
