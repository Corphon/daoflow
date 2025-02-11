//system/control/flow/scheduler.go

package flow

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/Corphon/daoflow/model"
)

// Scheduler 调度器
type Scheduler struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		maxTasks       int           // 最大任务数
		taskTimeout    time.Duration // 任务超时时间
		retryLimit     int           // 重试限制
		priorityLevels int           // 优先级层级
	}

	// 调度状态
	state struct {
		tasks     map[string]*Task         // 任务列表
		queues    map[int]*TaskQueue       // 优先级队列
		executors map[string]*TaskExecutor // 任务执行器
		history   []TaskHistory            // 任务历史
		metrics   SchedulerMetrics         // 调度指标
	}

	// 依赖项
	backpressure *BackpressureManager
}

// Task 任务
type Task struct {
	ID           string                 // 任务ID
	Type         string                 // 任务类型
	Priority     int                    // 优先级
	Status       string                 // 任务状态
	Parameters   map[string]interface{} // 任务参数
	Dependencies []string               // 依赖任务
	Created      time.Time              // 创建时间
	StartTime    time.Time              // 开始时间
	EndTime      time.Time              // 结束时间
	Deadline     time.Time              // 截止时间
	Retries      int                    // 重试次数
}

// TaskQueue 任务队列
type TaskQueue struct {
	Priority int        // 优先级
	Tasks    []*Task    // 任务列表
	Capacity int        // 队列容量
	Stats    QueueStats // 队列统计
}

// QueueStats 队列统计
type QueueStats struct {
	TotalTasks     int           // 总任务数
	CompletedTasks int           // 完成任务数
	FailedTasks    int           // 失败任务数
	AverageWait    time.Duration // 平均等待时间
}

// TaskExecutor 任务执行器
type TaskExecutor struct {
	ID          string              // 执行器ID
	Type        string              // 执行器类型
	Status      string              // 执行器状态
	Capacity    int                 // 处理容量
	CurrentLoad int                 // 当前负载
	Performance ExecutorPerformance // 性能指标
}

// ExecutorPerformance 执行器性能
type ExecutorPerformance struct {
	SuccessRate float64       // 成功率
	Throughput  float64       // 吞吐量
	Latency     time.Duration // 处理延迟
	ErrorRate   float64       // 错误率
	LastUpdate  time.Time     // 最后更新时间
}

// TaskHistory 任务历史
type TaskHistory struct {
	TaskID    string                 // 任务ID
	Type      string                 // 记录类型
	Status    string                 // 任务状态
	Details   map[string]interface{} // 详细信息
	Timestamp time.Time              // 记录时间
}

// SchedulerMetrics 调度指标
type SchedulerMetrics struct {
	ActiveTasks  int               // 活跃任务数
	QueueLength  map[int]int       // 队列长度
	Throughput   float64           // 系统吞吐量
	LatencyStats LatencyStatistics // 延迟统计
	History      []MetricPoint     // 历史指标
}

// LatencyStatistics 延迟统计
type LatencyStatistics struct {
	Average time.Duration // 平均延迟
	P95     time.Duration // 95分位延迟
	P99     time.Duration // 99分位延迟
	Max     time.Duration // 最大延迟
}

// FlowEvent 流事件记录
type FlowEvent struct {
	FlowID    string                 // 流ID
	Type      string                 // 事件类型
	Status    string                 // 流状态
	Details   map[string]interface{} // 详细信息
	Timestamp time.Time              // 记录时间
}

const (
	maxHistoryLength = 1000
)

// FlowScheduler 流调度器
type FlowScheduler struct {
	mu sync.RWMutex

	// 调度配置
	config struct {
		maxConcurrent    int           // 最大并发数
		queueCapacity    int           // 队列容量
		scheduleInterval time.Duration // 调度间隔
		priorityLevels   int           // 优先级级别
	}

	// 调度状态
	state struct {
		activeFlows    map[string]*FlowInfo // 活动流
		pendingQueue   *FlowQueue           // 等待队列
		completedFlows []string             // 已完成流
		history        []FlowEvent          // 事件历史记录
	}
}

// LoadAdjustment 负载调整记录
type LoadAdjustment struct {
	ResourceID string    // 资源ID
	Factor     float64   // 调整因子
	Timestamp  time.Time // 调整时间
}

// FlowBalancer 流平衡器
type FlowBalancer struct {
	mu sync.RWMutex

	// 平衡配置
	config struct {
		balanceInterval time.Duration // 平衡间隔
		threshold       float64       // 平衡阈值
		maxAdjustment   float64       // 最大调整量
		maxLoad         float64       // 最大负载
	}

	// 平衡状态
	state struct {
		loads       map[string]float64 // 负载情况
		adjustments []LoadAdjustment   // 调整记录
	}
}

// PressureRecord 压力记录
type PressureRecord struct {
	Timestamp time.Time          // 记录时间
	Pressure  float64            // 压力值
	Source    string             // 压力来源
	Status    string             // 当前状态
	Details   map[string]float64 // 详细信息
	Action    string             // 采取的动作
}

// BackPressure 背压控制器
type BackPressure struct {
	mu sync.RWMutex

	// 背压配置
	config struct {
		pressureThreshold float64       // 压力阈值
		releaseRate       float64       // 释放率
		checkInterval     time.Duration // 检查间隔
	}

	// 背压状态
	state struct {
		pressure    float64            // 当前压力
		constraints map[string]float64 // 约束条件
		history     []PressureRecord   // 压力记录
	}
}

// FlowInfo 流信息
type FlowInfo struct {
	ID           string             // 流ID
	Type         string             // 流类型
	Priority     int                // 优先级
	Status       string             // 流状态
	StartTime    time.Time          // 开始时间
	EndTime      time.Time          // 结束时间
	Resources    map[string]float64 // 资源使用
	Metrics      FlowMetrics        // 流指标
	Dependencies []string           // 依赖流
}

// FlowQueue 流队列
type FlowQueue struct {
	Items    []*FlowInfo // 队列项
	Capacity int         // 队列容量
	Head     int         // 队列头
	Tail     int         // 队列尾
	Size     int         // 当前大小
	Stats    QueueStats  // 队列统计
}

// FlowMetrics 流指标
type FlowMetrics struct {
	ProcessingTime time.Duration      // 处理时间
	WaitingTime    time.Duration      // 等待时间
	ThroughPut     float64            // 吞吐量
	ErrorRate      float64            // 错误率
	ResourceUsage  map[string]float64 // 资源使用率
}

// NewFlowBalancer 创建新的流平衡器
func NewFlowBalancer() *FlowBalancer {
	lb := &FlowBalancer{}

	// 初始化配置
	lb.config.balanceInterval = 5 * time.Second
	lb.config.threshold = 0.8
	lb.config.maxAdjustment = 0.5
	lb.config.maxLoad = 0.9

	// 初始化状态
	lb.state.loads = make(map[string]float64)
	lb.state.adjustments = make([]LoadAdjustment, 0)

	return lb
}

// GetFlowPressure 获取流的压力值
func (s *FlowScheduler) GetFlowPressure(flowID string) (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 获取流信息
	flow, exists := s.state.activeFlows[flowID]
	if !exists {
		return 0, fmt.Errorf("flow not found: %s", flowID)
	}

	// 计算流压力
	// 基于处理时间、等待时间和资源使用率计算
	pressure := 0.0

	// 处理时间贡献
	if flow.Metrics.ProcessingTime > 0 {
		pressure += float64(flow.Metrics.ProcessingTime.Milliseconds()) / 1000.0 * 0.4 // 40% 权重
	}

	// 等待时间贡献
	if flow.Metrics.WaitingTime > 0 {
		pressure += float64(flow.Metrics.WaitingTime.Milliseconds()) / 1000.0 * 0.3 // 30% 权重
	}

	// 资源使用贡献
	totalResourceUsage := 0.0
	for _, usage := range flow.Metrics.ResourceUsage {
		totalResourceUsage += usage
	}
	if len(flow.Metrics.ResourceUsage) > 0 {
		pressure += (totalResourceUsage / float64(len(flow.Metrics.ResourceUsage))) * 0.3 // 30% 权重
	}

	// 归一化到 [0,1] 区间
	return math.Min(1.0, math.Max(0.0, pressure)), nil
}

// GetResourcePressure 获取资源的压力值
func (lb *FlowBalancer) GetResourcePressure(resourceID string) (float64, error) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	// 获取资源负载
	load, exists := lb.state.loads[resourceID]
	if !exists {
		return 0, fmt.Errorf("resource not found: %s", resourceID)
	}

	// 负载直接作为压力值返回（假设负载已归一化到 [0,1] 区间）
	return load, nil
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
		History:     make([]MetricPoint, 0),
	}

	// 初始化优先级队列
	for i := 0; i < s.config.priorityLevels; i++ {
		s.state.queues[i] = &TaskQueue{
			Priority: i,
			Tasks:    make([]*Task, 0),
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

// checkSystemLoad 检查系统负载
func (s *Scheduler) checkSystemLoad() error {
	totalTasks := 0
	for _, queue := range s.state.queues {
		totalTasks += len(queue.Tasks)
	}

	// 检查总任务数是否超过限制
	if totalTasks >= s.config.maxTasks {
		return model.WrapError(nil, model.ErrCodeResource,
			"system overloaded: max tasks limit reached")
	}

	// 检查执行器负载
	activeExecutors := 0
	totalLoad := 0
	for _, executor := range s.state.executors {
		if executor.Status == "active" {
			activeExecutors++
			totalLoad += executor.CurrentLoad
		}
	}

	if activeExecutors > 0 {
		avgLoad := float64(totalLoad) / float64(activeExecutors)
		if avgLoad > 0.9 { // 90% 负载阈值
			return model.WrapError(nil, model.ErrCodeResource,
				"system overloaded: high executor load")
		}
	}

	return nil
}

// enqueueTask 将任务加入队列
func (s *Scheduler) enqueueTask(task *Task) error {
	// 获取对应优先级的队列
	queue, exists := s.state.queues[task.Priority]
	if !exists {
		return model.WrapError(nil, model.ErrCodeValidation,
			fmt.Sprintf("invalid priority level: %d", task.Priority))
	}

	// 检查队列容量
	if len(queue.Tasks) >= queue.Capacity {
		return model.WrapError(nil, model.ErrCodeResource,
			fmt.Sprintf("queue %d is full", task.Priority))
	}

	// 存储任务
	s.state.tasks[task.ID] = task

	// 添加到队列
	queue.Tasks = append(queue.Tasks, task)
	task.Status = "queued"

	// 更新队列统计
	queue.Stats.TotalTasks++

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

// checkExecutors 检查执行器状态
func (s *Scheduler) checkExecutors() error {
	activeCount := 0

	// 检查每个执行器的状态
	for _, executor := range s.state.executors {
		// 更新执行器状态
		if time.Since(executor.Performance.LastUpdate) > s.config.taskTimeout {
			executor.Status = "inactive"
			continue
		}

		// 检查执行器性能
		if executor.Status == "active" {
			activeCount++
			// 检查性能指标
			if executor.Performance.ErrorRate > 0.5 ||
				executor.Performance.SuccessRate < 0.3 {
				executor.Status = "degraded"
			}
		}
	}

	// 确保至少有一个活跃执行器
	if activeCount == 0 {
		return model.WrapError(nil, model.ErrCodeResource,
			"no active executors available")
	}

	return nil
}

// updateMetrics 更新调度指标
func (s *Scheduler) updateMetrics() {
	point := MetricPoint{
		Timestamp: time.Now(),
		Values:    make(map[string]float64),
	}

	// 计算活跃任务数
	activeTasks := 0
	for _, task := range s.state.tasks {
		if task.Status == "executing" {
			activeTasks++
		}
	}
	point.Values["active_tasks"] = float64(activeTasks)

	// 计算队列长度
	for priority, queue := range s.state.queues {
		s.state.metrics.QueueLength[priority] = len(queue.Tasks)
		point.Values[fmt.Sprintf("queue_%d_length", priority)] = float64(len(queue.Tasks))
	}

	// 计算系统吞吐量
	totalCompleted := 0
	for _, queue := range s.state.queues {
		totalCompleted += queue.Stats.CompletedTasks
	}
	windowSize := 5 * time.Minute // 5分钟窗口
	s.state.metrics.Throughput = float64(totalCompleted) / windowSize.Seconds()
	point.Values["throughput"] = s.state.metrics.Throughput

	// 更新延迟统计
	s.updateLatencyStats()
	point.Values["avg_latency"] = float64(s.state.metrics.LatencyStats.Average.Milliseconds())
	point.Values["p95_latency"] = float64(s.state.metrics.LatencyStats.P95.Milliseconds())
	point.Values["p99_latency"] = float64(s.state.metrics.LatencyStats.P99.Milliseconds())

	// 添加指标点
	s.state.metrics.History = append(s.state.metrics.History, point)

	// 限制历史记录长度
	if len(s.state.metrics.History) > maxHistoryLength {
		s.state.metrics.History = s.state.metrics.History[1:]
	}
}

// updateLatencyStats 更新延迟统计信息
func (s *Scheduler) updateLatencyStats() {
	var latencies []time.Duration

	// 收集最近任务的延迟数据
	for _, task := range s.state.tasks {
		if task.Status == "completed" {
			latency := task.EndTime.Sub(task.Created)
			latencies = append(latencies, latency)
		}
	}

	if len(latencies) == 0 {
		return
	}

	// 计算平均延迟
	var totalLatency time.Duration
	for _, latency := range latencies {
		totalLatency += latency
	}
	s.state.metrics.LatencyStats.Average = totalLatency / time.Duration(len(latencies))

	// 计算百分位数
	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	p95Index := int(float64(len(latencies)) * 0.95)
	p99Index := int(float64(len(latencies)) * 0.99)

	s.state.metrics.LatencyStats.P95 = latencies[p95Index]
	s.state.metrics.LatencyStats.P99 = latencies[p99Index]
	s.state.metrics.LatencyStats.Max = latencies[len(latencies)-1]
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

// findAvailableExecutor 获取可用执行器
func (s *Scheduler) findAvailableExecutor() *TaskExecutor {
	var selected *TaskExecutor
	minLoad := math.MaxInt

	// 寻找负载最小的活跃执行器
	for _, executor := range s.state.executors {
		if executor.Status != "active" {
			continue
		}
		if executor.CurrentLoad < minLoad &&
			executor.CurrentLoad < executor.Capacity {
			selected = executor
			minLoad = executor.CurrentLoad
		}
	}

	return selected
}

// isTaskExecutable 检查任务是否可执行
func (s *Scheduler) isTaskExecutable(task *Task) bool {
	// 检查任务状态
	if task.Status != "queued" && task.Status != "retry" {
		return false
	}

	// 检查依赖任务
	for _, depID := range task.Dependencies {
		dep, exists := s.state.tasks[depID]
		if !exists || dep.Status != "completed" {
			return false
		}
	}

	// 检查截止时间
	if !task.Deadline.IsZero() && time.Now().After(task.Deadline) {
		task.Status = "timeout"
		return false
	}

	return true
}

// executeTask 执行任务
func (s *Scheduler) executeTask(executor *TaskExecutor, task *Task) error {
	// 更新任务状态
	task.Status = "executing"
	task.StartTime = time.Now()

	// 更新执行器负载
	executor.CurrentLoad++

	// 记录执行开始
	s.recordTask(task, "execution_start", map[string]interface{}{
		"executor_id": executor.ID,
		"start_time":  task.StartTime,
	})

	// 模拟任务执行
	// 实际环境中这里应该调用实际的执行逻辑
	time.Sleep(100 * time.Millisecond)

	// 更新任务完成状态
	task.Status = "completed"
	task.EndTime = time.Now()

	// 更新执行器状态
	executor.CurrentLoad--
	executor.Performance.LastUpdate = time.Now()

	// 更新统计信息
	queue := s.state.queues[task.Priority]
	queue.Stats.CompletedTasks++

	// 记录执行完成
	s.recordTask(task, "execution_complete", map[string]interface{}{
		"executor_id": executor.ID,
		"end_time":    task.EndTime,
		"duration":    task.EndTime.Sub(task.StartTime),
	})

	return nil
}

// handleTaskError 处理任务错误
func (s *Scheduler) handleTaskError(task *Task, err error) {
	// 增加重试次数
	task.Retries++

	// 检查是否超过重试限制
	if task.Retries >= s.config.retryLimit {
		task.Status = "failed"
		// 更新统计
		queue := s.state.queues[task.Priority]
		queue.Stats.FailedTasks++
	} else {
		task.Status = "retry"
	}

	// 记录错误
	s.recordTask(task, "execution_error", map[string]interface{}{
		"error":       err.Error(),
		"retry_count": task.Retries,
	})
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

// ThrottleFlow 节流流处理
func (s *FlowScheduler) ThrottleFlow(flowID string, rate float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 获取流信息
	flow, exists := s.state.activeFlows[flowID]
	if !exists {
		return fmt.Errorf("flow not found: %s", flowID)
	}

	// 验证速率范围
	if rate < 0 || rate > 1 {
		return fmt.Errorf("invalid throttle rate: %f, must be between 0 and 1", rate)
	}

	// 应用节流
	flow.Metrics.ThroughPut *= rate
	flow.Status = "throttled"

	// 记录节流操作
	s.recordFlowEvent(flow, "throttle", map[string]interface{}{
		"rate":      rate,
		"timestamp": time.Now(),
	})

	return nil
}

// recordFlowEvent 记录流事件
func (s *FlowScheduler) recordFlowEvent(
	flow *FlowInfo,
	eventType string,
	details map[string]interface{}) {

	event := FlowEvent{
		FlowID:    flow.ID,
		Type:      eventType,
		Status:    flow.Status,
		Details:   details,
		Timestamp: time.Now(),
	}

	// 初始化历史记录切片（如果需要）
	if s.state.history == nil {
		s.state.history = make([]FlowEvent, 0)
	}

	// 添加事件记录
	s.state.history = append(s.state.history, event)

	// 限制历史记录长度
	if len(s.state.history) > maxHistoryLength {
		s.state.history = s.state.history[1:]
	}
}

// ScaleResource 调整资源规模
func (lb *FlowBalancer) ScaleResource(resourceID string, factor float64) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// 验证缩放因子
	if factor <= 0 {
		return fmt.Errorf("invalid scale factor: %f, must be positive", factor)
	}

	// 获取当前负载
	currentLoad, exists := lb.state.loads[resourceID]
	if !exists {
		return fmt.Errorf("resource not found: %s", resourceID)
	}

	// 计算新的负载容量
	newLoad := currentLoad / factor

	// 确保在合理范围内
	if newLoad > lb.config.maxLoad {
		newLoad = lb.config.maxLoad
	}

	// 更新负载
	lb.state.loads[resourceID] = newLoad

	// 记录调整
	lb.state.adjustments = append(lb.state.adjustments, LoadAdjustment{
		ResourceID: resourceID,
		Factor:     factor,
		Timestamp:  time.Now(),
	})

	return nil
}
