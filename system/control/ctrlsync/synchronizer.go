//system/control/ctrlsync/synchronizer.go

package ctrlsync

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/types"
)

// Synchronizer 同步器
type Synchronizer struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		syncInterval time.Duration // 同步间隔
		batchSize    int           // 批处理大小
		retryLimit   int           // 重试限制
		maxQueueSize int           // 最大队列大小
	}

	// 同步状态
	state struct {
		tasks      map[string]*SyncTask      // 同步任务
		queues     map[string]*SyncQueue     // 同步队列
		operations map[string]*SyncOperation // 同步操作
		metrics    SyncMetrics               // 同步指标
	}
}

// SyncTask 同步任务
type SyncTask struct {
	ID       string        // 任务ID
	Type     string        // 任务类型
	Source   *SyncEndpoint // 源端点
	Target   *SyncEndpoint // 目标端点
	State    string        // 任务状态
	Priority int           // 优先级
	Schedule *SyncSchedule // 同步计划
	LastSync time.Time     // 最后同步
}

// SyncEndpoint 同步端点
type SyncEndpoint struct {
	ID         string                 // 端点ID
	Type       string                 // 端点类型
	Location   string                 // 端点位置
	Properties map[string]interface{} // 端点属性
	State      *EndpointState         // 端点状态
}

// EndpointState 端点状态
type EndpointState struct {
	Status     string    // 状态标识
	Version    string    // 数据版本
	LastUpdate time.Time // 最后更新
	Checksum   string    // 数据校验
}

// SyncSchedule 同步计划
type SyncSchedule struct {
	Type       string          // 计划类型
	Interval   time.Duration   // 同步间隔
	TimeWindow *TimeWindow     // 时间窗口
	Conditions []SyncCondition // 同步条件
}

// TimeWindow 时间窗口
type TimeWindow struct {
	Start      time.Time // 开始时间
	End        time.Time // 结束时间
	Recurrence string    // 重复规则
}

// SyncCondition 同步条件
type SyncCondition struct {
	Type       string                 // 条件类型
	Expression string                 // 条件表达式
	Parameters map[string]interface{} // 条件参数
}

// SyncQueue 同步队列
type SyncQueue struct {
	ID       string      // 队列ID
	Type     string      // 队列类型
	Priority int         // 优先级
	Tasks    []*SyncTask // 任务列表
	Stats    QueueStats  // 队列统计
}

// QueueStats 队列统计
type QueueStats struct {
	TotalTasks int // 总任务数
	Pending    int // 等待任务
	Processing int // 处理中任务
	Completed  int // 完成任务
	Failed     int // 失败任务
}

// SyncOperation 同步操作
type SyncOperation struct {
	ID        string           // 操作ID
	TaskID    string           // 任务ID
	Type      string           // 操作类型
	Status    string           // 操作状态
	StartTime time.Time        // 开始时间
	EndTime   time.Time        // 结束时间
	Result    *OperationResult // 操作结果
}

// OperationResult 操作结果
type OperationResult struct {
	Success     bool     // 是否成功
	SyncedItems int      // 同步项数
	Errors      []string // 错误信息
	Checksum    string   // 校验和
}

// SyncMetrics 同步指标
type SyncMetrics struct {
	ActiveTasks    int                 // 活跃任务数
	QueuedTasks    int                 // 队列任务数
	SyncRate       float64             // 同步速率
	SuccessRate    float64             // 成功率
	AverageLatency time.Duration       // 平均延迟
	History        []types.MetricPoint // 历史指标
}

// ------------------------------------
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
		History: make([]types.MetricPoint, 0),
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

// enqueueTask 将任务加入队列
func (s *Synchronizer) enqueueTask(task *SyncTask) error {
	// 根据任务类型获取或创建队列
	queue, exists := s.state.queues[task.Type]
	if !exists {
		queue = &SyncQueue{
			ID:       fmt.Sprintf("queue_%s", task.Type),
			Type:     task.Type,
			Priority: task.Priority,
			Tasks:    make([]*SyncTask, 0),
			Stats: QueueStats{
				TotalTasks: 0,
				Pending:    0,
				Processing: 0,
				Completed:  0,
				Failed:     0,
			},
		}
		s.state.queues[task.Type] = queue
	}

	// 检查队列容量
	if len(queue.Tasks) >= s.config.maxQueueSize {
		return model.WrapError(nil, model.ErrCodeResource,
			fmt.Sprintf("queue %s is full", task.Type))
	}

	// 添加到队列
	queue.Tasks = append(queue.Tasks, task)
	queue.Stats.TotalTasks++
	queue.Stats.Pending++

	// 更新指标
	s.state.metrics.QueuedTasks++

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

// executeOperation 执行同步操作
func (s *Synchronizer) executeOperation(operation *SyncOperation) error {
	// 存储操作
	s.state.operations[operation.ID] = operation

	// 获取相关任务
	task, exists := s.state.tasks[operation.TaskID]
	if !exists {
		return model.WrapError(nil, model.ErrCodeNotFound, "task not found")
	}

	// 执行同步
	result := &OperationResult{
		Success:     true,
		SyncedItems: 0,
		Errors:      make([]string, 0),
		Checksum:    "",
	}

	// 验证端点状态
	if err := s.validateEndpoints(task.Source, task.Target); err != nil {
		result.Success = false
		result.Errors = append(result.Errors, err.Error())
		operation.Result = result
		return err
	}

	// 执行数据同步
	syncedItems, checksum, err := s.syncData(task.Source, task.Target)
	if err != nil {
		result.Success = false
		result.Errors = append(result.Errors, err.Error())
	} else {
		result.SyncedItems = syncedItems
		result.Checksum = checksum
	}

	// 更新操作结果
	operation.Status = "completed"
	operation.EndTime = time.Now()
	operation.Result = result

	return nil
}

// handleOperationError 处理操作错误
func (s *Synchronizer) handleOperationError(operation *SyncOperation, err error) {
	operation.Status = "failed"
	operation.EndTime = time.Now()
	operation.Result = &OperationResult{
		Success: false,
		Errors:  []string{err.Error()},
	}

	// 更新队列统计并计算成功率
	var failedCount int
	if queue, exists := s.state.queues[operation.Type]; exists {
		queue.Stats.Failed++
		queue.Stats.Processing--
		failedCount = queue.Stats.Failed
	}

	// 更新全局指标
	if s.state.metrics.QueuedTasks > 0 {
		s.state.metrics.SuccessRate = float64(s.state.metrics.QueuedTasks-failedCount) /
			float64(s.state.metrics.QueuedTasks)
	} else {
		s.state.metrics.SuccessRate = 0
	}
}

// updateTaskState 更新任务状态
func (s *Synchronizer) updateTaskState(task *SyncTask, operation *SyncOperation) {
	// 更新任务状态
	if operation.Result != nil && operation.Result.Success {
		task.State = "completed"
		task.LastSync = operation.EndTime

		// 更新端点状态
		if task.Source.State != nil {
			task.Source.State.LastUpdate = operation.EndTime
			task.Source.State.Checksum = operation.Result.Checksum
		}
		if task.Target.State != nil {
			task.Target.State.LastUpdate = operation.EndTime
			task.Target.State.Checksum = operation.Result.Checksum
		}
	} else {
		task.State = "failed"
	}

	// 更新队列统计
	if queue, exists := s.state.queues[task.Type]; exists {
		if task.State == "completed" {
			queue.Stats.Completed++
		}
		queue.Stats.Processing--
	}
}

// validateEndpoints 验证端点状态
func (s *Synchronizer) validateEndpoints(source, target *SyncEndpoint) error {
	if source.State == nil || target.State == nil {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid endpoint state")
	}

	if source.State.Status != "active" || target.State.Status != "active" {
		return model.WrapError(nil, model.ErrCodeValidation, "endpoints not active")
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
	point := types.MetricPoint{
		Timestamp: time.Now(),
		Values: map[string]float64{
			"active_tasks": float64(len(s.state.tasks)),
			"queued_tasks": float64(s.state.metrics.QueuedTasks),
			"success_rate": s.state.metrics.SuccessRate,
			"sync_rate":    s.state.metrics.SyncRate,
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

// syncData 执行数据同步
func (s *Synchronizer) syncData(source, target *SyncEndpoint) (int, string, error) {
	// 检查端点配置
	if source.Location == "" || target.Location == "" {
		return 0, "", model.WrapError(nil, model.ErrCodeValidation,
			"invalid endpoint location")
	}

	// 计算同步开始时间
	startTime := time.Now()

	// 获取源数据
	sourceData, err := s.fetchEndpointData(source)
	if err != nil {
		return 0, "", model.WrapError(err, model.ErrCodeOperation,
			"failed to fetch source data")
	}

	// 计算数据校验和
	checksum := calculateDataChecksum(sourceData)

	// 检查目标端是否需要同步
	if target.State != nil && target.State.Checksum == checksum {
		return 0, checksum, nil // 数据一致，无需同步
	}

	// 执行数据同步
	syncedItems, err := s.transferData(sourceData, target)
	if err != nil {
		return 0, "", model.WrapError(err, model.ErrCodeOperation,
			"failed to transfer data")
	}

	// 更新同步统计
	s.updateSyncStats(startTime, syncedItems)

	return syncedItems, checksum, nil
}

// fetchEndpointData 获取端点数据
func (s *Synchronizer) fetchEndpointData(endpoint *SyncEndpoint) (interface{}, error) {
	// 根据端点类型获取数据
	switch endpoint.Type {
	case "file":
		return s.readFileData(endpoint.Location)
	case "database":
		return s.readDatabaseData(endpoint.Properties)
	case "api":
		return s.fetchAPIData(endpoint.Location, endpoint.Properties)
	default:
		return nil, model.WrapError(nil, model.ErrCodeValidation,
			fmt.Sprintf("unsupported endpoint type: %s", endpoint.Type))
	}
}

// transferData 传输数据到目标端
func (s *Synchronizer) transferData(data interface{}, target *SyncEndpoint) (int, error) {
	// 根据目标端点类型写入数据
	switch target.Type {
	case "file":
		return s.writeFileData(data, target.Location)
	case "database":
		return s.writeDatabaseData(data, target.Properties)
	case "api":
		return s.sendAPIData(data, target.Location, target.Properties)
	default:
		return 0, model.WrapError(nil, model.ErrCodeValidation,
			fmt.Sprintf("unsupported target type: %s", target.Type))
	}
}

// updateSyncStats 更新同步统计信息
func (s *Synchronizer) updateSyncStats(startTime time.Time, syncedItems int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 计算同步延迟
	latency := time.Since(startTime)

	// 更新性能指标
	if s.state.metrics.AverageLatency == 0 {
		s.state.metrics.AverageLatency = latency
	} else {
		s.state.metrics.AverageLatency = (s.state.metrics.AverageLatency + latency) / 2
	}

	// 更新同步速率
	if syncedItems > 0 {
		s.state.metrics.SyncRate = float64(syncedItems) / latency.Seconds()
	}
}

// calculateDataChecksum 计算数据校验和
func calculateDataChecksum(data interface{}) string {
	// 使用 SHA-256 计算校验和
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", data)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// readFileData 读取文件数据
func (s *Synchronizer) readFileData(location string) (interface{}, error) {
	// 验证文件路径
	if location == "" {
		return nil, model.WrapError(nil, model.ErrCodeValidation,
			"empty file location")
	}

	// 读取文件数据
	data, err := os.ReadFile(location)
	if err != nil {
		return nil, model.WrapError(err, model.ErrCodeIO,
			fmt.Sprintf("failed to read file: %s", location))
	}

	// 解析数据（根据文件扩展名）
	ext := strings.ToLower(filepath.Ext(location))
	switch ext {
	case ".json":
		var result interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, model.WrapError(err, model.ErrCodeTransform,
				"failed to parse JSON data")
		}
		return result, nil
	// 可以添加其他格式的支持
	default:
		return data, nil
	}
}

// readDatabaseData 读取数据库数据
func (s *Synchronizer) readDatabaseData(props map[string]interface{}) (interface{}, error) {
	// 验证数据库配置
	dsn, ok := props["dsn"].(string)
	if !ok {
		return nil, model.WrapError(nil, model.ErrCodeValidation,
			"missing database connection string")
	}

	query, ok := props["query"].(string)
	if !ok {
		return nil, model.WrapError(nil, model.ErrCodeValidation,
			"missing query statement")
	}

	// 获取数据库类型
	dbType, _ := props["type"].(string)
	if dbType == "" {
		dbType = "default" // 默认数据库类型
	}

	// 执行查询
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := s.executeQuery(ctx, dbType, dsn, query)
	if err != nil {
		return nil, model.WrapError(err, model.ErrCodeOperation,
			"failed to execute database query")
	}

	return result, nil
}

// fetchAPIData 获取API数据
func (s *Synchronizer) fetchAPIData(url string, props map[string]interface{}) (interface{}, error) {
	// 验证API配置
	if url == "" {
		return nil, model.WrapError(nil, model.ErrCodeValidation,
			"empty API URL")
	}

	// 创建请求上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 准备请求
	method, _ := props["method"].(string)
	if method == "" {
		method = "GET"
	}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, model.WrapError(err, model.ErrCodeOperation,
			"failed to create API request")
	}

	// 添加请求头
	if headers, ok := props["headers"].(map[string]string); ok {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, model.WrapError(err, model.ErrCodeOperation,
			"failed to execute API request")
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, model.WrapError(err, model.ErrCodeIO,
			"failed to read API response")
	}

	// 检查响应状态
	if resp.StatusCode >= 400 {
		return nil, model.WrapError(nil, model.ErrCodeOperation,
			fmt.Sprintf("API request failed with status: %d", resp.StatusCode))
	}

	// 解析响应数据
	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, model.WrapError(err, model.ErrCodeTransform,
			"failed to parse API response")
	}

	return result, nil
}

// executeQuery 执行数据库查询
func (s *Synchronizer) executeQuery(ctx context.Context, dbType, dsn, query string) (interface{}, error) {
	// 验证输入
	if query == "" {
		return nil, model.WrapError(nil, model.ErrCodeValidation, "empty query")
	}

	// 只支持 SQLite
	if dbType != "sqlite" && dbType != "default" {
		return nil, model.WrapError(nil, model.ErrCodeValidation,
			"only sqlite database is supported")
	}

	return s.querySQLite(ctx, dsn, query)
}

// writeFileData 写入文件数据
func (s *Synchronizer) writeFileData(data interface{}, location string) (int, error) {
	// 验证文件路径
	if location == "" {
		return 0, model.WrapError(nil, model.ErrCodeValidation,
			"empty file location")
	}

	// 创建目录
	dir := filepath.Dir(location)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return 0, model.WrapError(err, model.ErrCodeIO,
			fmt.Sprintf("failed to create directory: %s", dir))
	}

	// 序列化数据
	var bytes []byte
	var err error

	// 根据文件扩展名选择序列化方式
	ext := strings.ToLower(filepath.Ext(location))
	switch ext {
	case ".json":
		bytes, err = json.Marshal(data)
		if err != nil {
			return 0, model.WrapError(err, model.ErrCodeTransform,
				"failed to marshal JSON data")
		}
	default:
		// 对于其他类型，尝试直接转换为字节
		if b, ok := data.([]byte); ok {
			bytes = b
		} else {
			bytes = []byte(fmt.Sprintf("%v", data))
		}
	}

	// 写入文件
	if err := os.WriteFile(location, bytes, 0644); err != nil {
		return 0, model.WrapError(err, model.ErrCodeIO,
			fmt.Sprintf("failed to write file: %s", location))
	}

	return len(bytes), nil
}

// writeDatabaseData 写入数据库数据
func (s *Synchronizer) writeDatabaseData(data interface{}, props map[string]interface{}) (int, error) {
	// 验证数据库配置
	dsn, ok := props["dsn"].(string)
	if !ok {
		return 0, model.WrapError(nil, model.ErrCodeValidation,
			"missing database connection string")
	}

	table, ok := props["table"].(string)
	if !ok {
		return 0, model.WrapError(nil, model.ErrCodeValidation,
			"missing table name")
	}

	// 获取数据库类型
	dbType, _ := props["type"].(string)
	if dbType == "" {
		dbType = "default"
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 转换数据为可写入格式
	records, err := s.convertToRecords(data)
	if err != nil {
		return 0, model.WrapError(err, model.ErrCodeTransform,
			"failed to convert data to records")
	}

	// 执行写入
	count, err := s.executeWrite(ctx, dbType, dsn, table, records)
	if err != nil {
		return 0, model.WrapError(err, model.ErrCodeOperation,
			"failed to write database data")
	}

	return count, nil
}

// sendAPIData 发送API数据
func (s *Synchronizer) sendAPIData(data interface{}, url string, props map[string]interface{}) (int, error) {
	// 验证API配置
	if url == "" {
		return 0, model.WrapError(nil, model.ErrCodeValidation,
			"empty API URL")
	}

	// 准备请求数据
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, model.WrapError(err, model.ErrCodeTransform,
			"failed to marshal API data")
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 准备请求
	method, _ := props["method"].(string)
	if method == "" {
		method = "POST"
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, model.WrapError(err, model.ErrCodeOperation,
			"failed to create API request")
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	if headers, ok := props["headers"].(map[string]string); ok {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, model.WrapError(err, model.ErrCodeOperation,
			"failed to send API request")
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode >= 400 {
		return 0, model.WrapError(nil, model.ErrCodeOperation,
			fmt.Sprintf("API request failed with status: %d", resp.StatusCode))
	}

	return len(jsonData), nil
}

// convertToRecords 转换数据为记录格式
func (s *Synchronizer) convertToRecords(data interface{}) ([]map[string]interface{}, error) {
	var records []map[string]interface{}

	switch v := data.(type) {
	case []map[string]interface{}:
		records = v
	case map[string]interface{}:
		records = []map[string]interface{}{v}
	default:
		// 尝试将其他类型转换为JSON然后解析
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(jsonData, &records); err != nil {
			return nil, err
		}
	}

	return records, nil
}

// executeWrite 执行数据库写入操作
func (s *Synchronizer) executeWrite(ctx context.Context, dbType, dsn, table string, records []map[string]interface{}) (int, error) {
	// 验证参数
	if table == "" {
		return 0, model.WrapError(nil, model.ErrCodeValidation, "empty table name")
	}
	if len(records) == 0 {
		return 0, nil
	}

	// 只支持 SQLite
	if dbType != "sqlite" && dbType != "default" {
		return 0, model.WrapError(nil, model.ErrCodeValidation,
			"only sqlite database is supported")
	}

	return s.writeToSQLite(ctx, dsn, table, records)
}

// writeToSQLite 写入SQLite数据库
func (s *Synchronizer) writeToSQLite(ctx context.Context, dsn, table string, records []map[string]interface{}) (int, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return 0, model.WrapError(err, model.ErrCodeOperation,
			"failed to open sqlite database")
	}
	defer db.Close()

	// 开启事务
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, model.WrapError(err, model.ErrCodeOperation,
			"failed to start transaction")
	}
	defer tx.Rollback()

	// 获取列名
	if len(records) == 0 {
		return 0, nil
	}
	columns := make([]string, 0, len(records[0]))
	for col := range records[0] {
		columns = append(columns, col)
	}

	// 构建插入语句
	placeholders := make([]string, len(columns))
	for i := range columns {
		placeholders[i] = "?"
	}
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	// 执行插入
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return 0, model.WrapError(err, model.ErrCodeOperation,
			"failed to prepare statement")
	}
	defer stmt.Close()

	count := 0
	for _, record := range records {
		values := make([]interface{}, len(columns))
		for i, col := range columns {
			values[i] = record[col]
		}

		if _, err := stmt.ExecContext(ctx, values...); err != nil {
			return count, model.WrapError(err, model.ErrCodeOperation,
				"failed to insert record")
		}
		count++
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return count, model.WrapError(err, model.ErrCodeOperation,
			"failed to commit transaction")
	}

	return count, nil
}

// querySQLite SQLite查询实现
func (s *Synchronizer) querySQLite(ctx context.Context, dsn, query string) (interface{}, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, model.WrapError(err, model.ErrCodeOperation,
			"failed to open sqlite database")
	}
	defer db.Close()

	// 执行查询
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, model.WrapError(err, model.ErrCodeOperation,
			"failed to execute query")
	}
	defer rows.Close()

	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		return nil, model.WrapError(err, model.ErrCodeOperation,
			"failed to get columns")
	}

	// 构造结果集
	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}

		if err := rows.Scan(values...); err != nil {
			return nil, model.WrapError(err, model.ErrCodeOperation,
				"failed to scan row")
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}

	return results, nil
}
