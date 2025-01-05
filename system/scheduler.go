// system/scheduler.go

package system

import (
    "context"
    "math"
    "sync"
    "time"
    "container/heap"

    "github.com/Corphon/daoflow/model"
)

// SchedulerConstants 调度常数
const (
    MaxTaskQueue     = 1000    // 最大任务队列长度
    MinPriority     = 0       // 最小优先级
    MaxPriority     = 10      // 最大优先级
    ScheduleInterval = time.Millisecond * 100 // 调度间隔
    ResourceBuffer   = 0.2     // 资源缓冲比例
)

// SchedulerSystem 调度系统
type SchedulerSystem struct {
    mu sync.RWMutex

    // 关联系统
    systemCore *SystemCore

    // 调度状态
    state struct {
        Tasks      *TaskQueue         // 任务队列
        Resources  *ResourceManager   // 资源管理器
        Workers    *WorkerPool        // 工作池
        Balancer   *LoadBalancer     // 负载均衡器
    }

    // 策略控制
    policy struct {
        Priority   *PriorityPolicy   // 优先级策略
        Resource   *ResourcePolicy   // 资源策略
        Balance    *BalancePolicy    // 均衡策略
    }

    metrics *SchedulerMetrics
    ctx     context.Context
    cancel  context.CancelFunc
}

// Task 任务
type Task struct {
    ID          string
    Type        TaskType
    Priority    int
    Resources   ResourceRequirement
    State       TaskState
    Progress    float64
    StartTime   time.Time
    Deadline    time.Time
}

// TaskType 任务类型
type TaskType uint8

const (
    TaskCompute TaskType = iota // 计算任务
    TaskIO                     // IO任务
    TaskSync                   // 同步任务
    TaskControl               // 控制任务
)

// TaskState 任务状态
type TaskState uint8

const (
    TaskPending TaskState = iota
    TaskRunning
    TaskCompleted
    TaskFailed
)

// ResourceRequirement 资源需求
type ResourceRequirement struct {
    CPU      float64 // CPU需求
    Memory   float64 // 内存需求
    Energy   float64 // 能量需求
}

// ResourceManager 资源管理器
type ResourceManager struct {
    total     ResourcePool    // 总资源
    available ResourcePool    // 可用资源
    allocated map[string]ResourcePool // 已分配资源
}

// WorkerPool 工作池
type WorkerPool struct {
    workers    map[string]*Worker
    capacity   int
    active     int
}

// Worker 工作单元
type Worker struct {
    ID        string
    State     WorkerState
    Task      *Task
    Resources ResourcePool
    Stats     WorkerStats
}

// LoadBalancer 负载均衡器
type LoadBalancer struct {
    strategy    BalanceStrategy
    metrics     map[string]float64
    history     []BalanceRecord
}

// NewSchedulerSystem 创建调度系统
func NewSchedulerSystem(ctx context.Context, sc *SystemCore) *SchedulerSystem {
    ctx, cancel := context.WithCancel(ctx)

    ss := &SchedulerSystem{
        systemCore: sc,
        ctx:       ctx,
        cancel:    cancel,
    }

    // 初始化状态
    ss.initializeState()

    // 初始化策略
    ss.initializePolicy()

    go ss.run()
    return ss
}

// initializeState 初始化状态
func (ss *SchedulerSystem) initializeState() {
    // 初始化任务队列
    ss.state.Tasks = NewTaskQueue()

    // 初始化资源管理器
    ss.state.Resources = NewResourceManager()

    // 初始化工作池
    ss.state.Workers = NewWorkerPool(runtime.NumCPU())

    // 初始化负载均衡器
    ss.state.Balancer = NewLoadBalancer()
}

// run 运行调度器
func (ss *SchedulerSystem) run() {
    ticker := time.NewTicker(ScheduleInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ss.ctx.Done():
            return
        case <-ticker.C:
            ss.schedule()
        }
    }
}

// schedule 执行调度
func (ss *SchedulerSystem) schedule() {
    ss.mu.Lock()
    defer ss.mu.Unlock()

    // 更新资源状态
    ss.updateResourceState()

    // 任务调度
    ss.scheduleTasks()

    // 负载均衡
    ss.balance()

    // 更新指标
    ss.updateMetrics()
}

// scheduleTasks 调度任务
func (ss *SchedulerSystem) scheduleTasks() {
    for !ss.state.Tasks.Empty() {
        // 获取最高优先级任务
        task := ss.state.Tasks.Peek()

        // 检查资源是否满足
        if !ss.checkResourceAvailability(task) {
            break
        }

        // 分配任务
        if worker := ss.allocateTask(task); worker != nil {
            ss.state.Tasks.Pop()
            worker.AssignTask(task)
        } else {
            break
        }
    }
}

// checkResourceAvailability 检查资源可用性
func (ss *SchedulerSystem) checkResourceAvailability(task *Task) bool {
    required := task.Resources
    available := ss.state.Resources.available

    return available.CPU >= required.CPU &&
           available.Memory >= required.Memory &&
           available.Energy >= required.Energy
}

// allocateTask 分配任务
func (ss *SchedulerSystem) allocateTask(task *Task) *Worker {
    // 使用最适合策略选择worker
    worker := ss.selectBestWorker(task)
    if worker == nil {
        return nil
    }

    // 分配资源
    if err := ss.allocateResources(worker, task.Resources); err != nil {
        return nil
    }

    return worker
}

// selectBestWorker 选择最佳worker
func (ss *SchedulerSystem) selectBestWorker(task *Task) *Worker {
    var bestWorker *Worker
    var minScore float64 = math.MaxFloat64

    for _, worker := range ss.state.Workers.workers {
        if worker.State != WorkerIdle {
            continue
        }

        score := ss.calculateWorkerScore(worker, task)
        if score < minScore {
            minScore = score
            bestWorker = worker
        }
    }

    return bestWorker
}

// balance 执行负载均衡
func (ss *SchedulerSystem) balance() {
    // 计算当前负载分布
    loadDistribution := ss.calculateLoadDistribution()

    // 检查是否需要均衡
    if ss.needRebalance(loadDistribution) {
        // 执行负载重分配
        ss.rebalanceLoad(loadDistribution)
    }
}

// SubmitTask 提交任务
func (ss *SchedulerSystem) SubmitTask(task *Task) error {
    ss.mu.Lock()
    defer ss.mu.Unlock()

    // 验证任务
    if err := ss.validateTask(task); err != nil {
        return err
    }

    // 应用优先级策略
    ss.applyPriorityPolicy(task)

    // 入队
    ss.state.Tasks.Push(task)

    return nil
}

// GetSchedulerStatus 获取调度器状态
func (ss *SchedulerSystem) GetSchedulerStatus() map[string]interface{} {
    ss.mu.RLock()
    defer ss.mu.RUnlock()

    return map[string]interface{}{
        "tasks":     ss.state.Tasks.Status(),
        "resources": ss.state.Resources.Status(),
        "workers":   ss.state.Workers.Status(),
        "balance":   ss.state.Balancer.Status(),
        "metrics":   ss.metrics,
    }
}

// Close 关闭调度系统
func (ss *SchedulerSystem) Close() error {
    ss.cancel()
    return nil
}
