//system/control/ctrlsync/coordinator.go

package ctrlsync

import (
	"fmt"
	"sync"
	"time"

	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/types"
)

// Coordinator 协调器
type Coordinator struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		coordinationInterval time.Duration // 协调间隔
		consensusTimeout     time.Duration // 共识超时
		maxRetries           int           // 最大重试次数
		quorumSize           int           // 法定人数大小
		maxParticipants      int           // 最大参与者数量
	}

	// 协调状态
	state struct {
		processes  map[string]*Process   // 进程列表
		sessions   map[string]*Session   // 会话列表
		agreements map[string]*Agreement // 共识协议
		metrics    CoordinationMetrics   // 协调指标
	}

	// 依赖项
	resolver     *Resolver
	synchronizer *Synchronizer
}

// Process 进程信息
type Process struct {
	ID           string    // 进程ID
	Type         string    // 进程类型
	State        string    // 进程状态
	Priority     int       // 优先级
	Dependencies []string  // 依赖进程
	LastUpdate   time.Time // 最后更新
}

// Session 协调会话
type Session struct {
	ID           string       // 会话ID
	Type         string       // 会话类型
	Participants []string     // 参与者
	State        SessionState // 会话状态
	StartTime    time.Time    // 开始时间
	Deadline     time.Time    // 截止时间
}

// SessionState 会话状态
type SessionState struct {
	Phase     string     // 当前阶段
	Progress  float64    // 进度
	Decisions []Decision // 决策列表
	Conflicts []Conflict // 冲突列表
}

// Decision 决策信息
type Decision struct {
	ID        string          // 决策ID
	Topic     string          // 决策主题
	Value     interface{}     // 决策值
	Votes     map[string]bool // 投票情况
	Timestamp time.Time       // 决策时间
}

// Agreement 共识协议
type Agreement struct {
	ID         string               // 协议ID
	Type       string               // 协议类型
	Terms      []Term               // 协议条款
	State      string               // 协议状态
	Signatures map[string]Signature // 签名列表
}

// Term 协议条款
type Term struct {
	ID          string       // 条款ID
	Content     string       // 条款内容
	Constraints []Constraint // 约束条件
	Priority    int          // 优先级
}

// Signature 签名信息
type Signature struct {
	ProcessID string    // 进程ID
	Timestamp time.Time // 签名时间
	Valid     bool      // 是否有效
}

// Constraint 约束条件
type Constraint struct {
	Type      string      // 约束类型
	Value     interface{} // 约束值
	Operator  string      // 操作符
	Tolerance float64     // 容差值
}

// CoordinationMetrics 协调指标
type CoordinationMetrics struct {
	ActiveSessions int                 // 活跃会话数
	SuccessRate    float64             // 成功率
	AverageLatency time.Duration       // 平均延迟
	ConflictRate   float64             // 冲突率
	History        []types.MetricPoint // 历史指标
}

// SyncCoordinator 同步协调器
type SyncCoordinator struct {
	mu sync.RWMutex

	// 协调配置
	config struct {
		coordinateInterval time.Duration // 协调间隔
		maxParticipants    int           // 最大参与者
		consensusTimeout   time.Duration // 共识超时
	}

	// 协调状态
	state struct {
		participants map[string]*Participant // 参与者
		sessions     map[string]*SyncSession // 同步会话
		events       []CoordinationEvent     // 协调事件
	}
}

// StateSynchronizer 状态同步器
type StateSynchronizer struct {
	mu sync.RWMutex

	// 同步配置
	config struct {
		syncInterval time.Duration // 同步间隔
		batchSize    int           // 批次大小
		retryLimit   int           // 重试限制
	}

	// 同步状态
	state struct {
		syncs     map[string]*SyncTask // 同步任务
		conflicts []SyncConflict       // 同步冲突
		stats     SyncStats            // 同步统计
	}
}

// NewCoordinator 创建新的协调器
func NewCoordinator(resolver *Resolver, synchronizer *Synchronizer) *Coordinator {
	c := &Coordinator{
		resolver:     resolver,
		synchronizer: synchronizer,
	}

	// 初始化配置
	c.config.coordinationInterval = 1 * time.Second
	c.config.consensusTimeout = 10 * time.Second
	c.config.maxRetries = 3
	c.config.quorumSize = 2
	c.config.maxParticipants = 10

	// 初始化状态
	c.state.processes = make(map[string]*Process)
	c.state.sessions = make(map[string]*Session)
	c.state.agreements = make(map[string]*Agreement)
	c.state.metrics = CoordinationMetrics{
		History: make([]types.MetricPoint, 0),
	}

	return c
}

// Coordinate 执行协调
func (c *Coordinator) Coordinate() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 更新进程状态
	if err := c.updateProcessStates(); err != nil {
		return err
	}

	// 管理会话
	if err := c.manageSessions(); err != nil {
		return err
	}

	// 达成共识
	if err := c.achieveConsensus(); err != nil {
		return err
	}

	// 解决冲突
	if err := c.resolveConflicts(); err != nil {
		return err
	}

	// 更新指标
	c.updateMetrics()

	return nil
}

// RegisterProcess 注册进程
func (c *Coordinator) RegisterProcess(process *Process) error {
	if process == nil {
		return model.WrapError(nil, model.ErrCodeValidation, "nil process")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// 验证进程
	if err := c.validateProcess(process); err != nil {
		return err
	}

	// 存储进程
	c.state.processes[process.ID] = process

	return nil
}

// CreateSession 创建协调会话
func (c *Coordinator) CreateSession(
	participants []string,
	sessionType string) (*Session, error) {

	c.mu.Lock()
	defer c.mu.Unlock()

	// 验证参与者
	if err := c.validateParticipants(participants); err != nil {
		return nil, err
	}

	// 创建会话
	session := &Session{
		ID:           generateSessionID(),
		Type:         sessionType,
		Participants: participants,
		State: SessionState{
			Phase:     "initialized",
			Progress:  0.0,
			Decisions: make([]Decision, 0),
			Conflicts: make([]Conflict, 0),
		},
		StartTime: time.Now(),
		Deadline:  time.Now().Add(c.config.consensusTimeout),
	}

	// 存储会话
	c.state.sessions[session.ID] = session

	return session, nil
}

// validateParticipants 验证参与者列表
func (c *Coordinator) validateParticipants(participants []string) error {
	// 验证参与者列表非空
	if len(participants) == 0 {
		return model.WrapError(nil, model.ErrCodeValidation, "empty participants list")
	}

	// 验证参与者数量不超过限制
	if len(participants) > c.config.maxParticipants {
		return model.WrapError(nil, model.ErrCodeValidation,
			fmt.Sprintf("participants count exceeds limit: %d", c.config.maxParticipants))
	}

	// 验证每个参与者
	seen := make(map[string]bool)
	for _, pid := range participants {
		// 检查参与者ID非空
		if pid == "" {
			return model.WrapError(nil, model.ErrCodeValidation, "empty participant ID")
		}

		// 检查重复参与者
		if seen[pid] {
			return model.WrapError(nil, model.ErrCodeValidation,
				fmt.Sprintf("duplicate participant: %s", pid))
		}
		seen[pid] = true

		// 检查参与者是否存在且活跃
		process, exists := c.state.processes[pid]
		if !exists {
			return model.WrapError(nil, model.ErrCodeNotFound,
				fmt.Sprintf("participant not found: %s", pid))
		}

		if process.State != "active" {
			return model.WrapError(nil, model.ErrCodeValidation,
				fmt.Sprintf("participant not active: %s", pid))
		}
	}

	// 验证参与者数量满足法定人数要求
	if len(participants) < c.config.quorumSize {
		return model.WrapError(nil, model.ErrCodeValidation,
			fmt.Sprintf("insufficient participants for quorum: %d < %d",
				len(participants), c.config.quorumSize))
	}

	return nil
}

// manageSessions 管理会话
func (c *Coordinator) manageSessions() error {
	currentTime := time.Now()

	for id, session := range c.state.sessions {
		// 检查会话是否过期
		if currentTime.After(session.Deadline) {
			if err := c.handleExpiredSession(session); err != nil {
				continue
			}
			delete(c.state.sessions, id)
			continue
		}

		// 更新会话状态
		if err := c.updateSessionState(session); err != nil {
			continue
		}
	}

	return nil
}

// handleExpiredSession 处理过期会话
func (c *Coordinator) handleExpiredSession(session *Session) error {
	// 更新会话状态
	session.State.Phase = "expired"

	// 处理未完成的决策
	for _, decision := range session.State.Decisions {
		if len(decision.Votes) < c.config.quorumSize {
			// 记录未达成共识的决策
			c.state.metrics.SuccessRate = (c.state.metrics.SuccessRate *
				float64(c.state.metrics.ActiveSessions-1)) /
				float64(c.state.metrics.ActiveSessions)
		}
	}

	// 取消未解决的冲突
	for _, conflict := range session.State.Conflicts {
		c.recordConflictResolution(session, &conflict,
			model.WrapError(nil, model.ErrCodeTimeout, "session expired"))
	}

	// 更新指标
	c.state.metrics.ActiveSessions--

	return nil
}

// updateSessionState 更新会话状态
func (c *Coordinator) updateSessionState(session *Session) error {
	// 计算新进度
	newProgress := calculateProgress(session)

	// 检查进度变化
	if newProgress != session.State.Progress {
		session.State.Progress = newProgress

		// 根据进度更新阶段
		switch {
		case newProgress >= 1.0:
			session.State.Phase = "completed"
		case newProgress >= 0.5:
			session.State.Phase = "consensus"
		case newProgress > 0:
			session.State.Phase = "collecting"
		}
	}

	// 检查是否存在新冲突
	if len(session.State.Conflicts) > 0 && session.State.Phase != "resolving" {
		session.State.Phase = "resolving"
	}

	// 更新指标
	if session.State.Phase == "completed" {
		c.state.metrics.SuccessRate = (c.state.metrics.SuccessRate*
			float64(c.state.metrics.ActiveSessions-1) + 1.0) /
			float64(c.state.metrics.ActiveSessions)
	}

	return nil
}

// achieveConsensus 达成共识
func (c *Coordinator) achieveConsensus() error {
	for _, session := range c.state.sessions {
		if session.State.Phase != "consensus" {
			continue
		}

		// 收集决策
		decisions, err := c.collectDecisions(session)
		if err != nil {
			continue
		}

		// 验证决策
		if err := c.validateDecisions(decisions); err != nil {
			continue
		}

		// 创建协议
		agreement := c.createAgreement(session, decisions)

		// 获取签名
		if err := c.collectSignatures(agreement); err != nil {
			continue
		}

		// 完成共识
		c.finalizeConsensus(session, agreement)
	}

	return nil
}

// collectDecisions 收集决策
func (c *Coordinator) collectDecisions(session *Session) ([]Decision, error) {
	decisions := make([]Decision, 0)
	requiredVotes := len(session.Participants)

	// 收集所有待决策项
	for _, decision := range session.State.Decisions {
		if len(decision.Votes) >= requiredVotes {
			decisions = append(decisions, decision)
		}
	}

	if len(decisions) == 0 {
		return nil, model.WrapError(nil, model.ErrCodeValidation, "no valid decisions")
	}

	return decisions, nil
}

// validateDecisions 验证决策
func (c *Coordinator) validateDecisions(decisions []Decision) error {
	for _, decision := range decisions {
		// 验证投票数量
		if len(decision.Votes) < c.config.quorumSize {
			return model.WrapError(nil, model.ErrCodeValidation,
				fmt.Sprintf("insufficient votes for decision %s", decision.ID))
		}

		// 验证一致性
		consensus := true
		var lastVote bool
		for _, vote := range decision.Votes {
			if lastVote != vote && len(decision.Votes) > 1 {
				consensus = false
				break
			}
			lastVote = vote
		}

		if !consensus {
			return model.WrapError(nil, model.ErrCodeConsensus,
				fmt.Sprintf("no consensus reached for decision %s", decision.ID))
		}
	}

	return nil
}

// createAgreement 创建共识协议
func (c *Coordinator) createAgreement(session *Session, decisions []Decision) *Agreement {
	agreement := &Agreement{
		ID:         fmt.Sprintf("agreement_%s", session.ID),
		Type:       session.Type,
		Terms:      make([]Term, 0),
		State:      "pending",
		Signatures: make(map[string]Signature),
	}

	// 将决策转换为协议条款
	for i, decision := range decisions {
		term := Term{
			ID:       fmt.Sprintf("term_%d", i),
			Content:  fmt.Sprintf("%v", decision.Value),
			Priority: i + 1,
			Constraints: []Constraint{
				{
					Type:      "quorum",
					Value:     c.config.quorumSize,
					Operator:  ">=",
					Tolerance: 0,
				},
			},
		}
		agreement.Terms = append(agreement.Terms, term)
	}

	// 存储协议
	c.state.agreements[agreement.ID] = agreement

	return agreement
}

// collectSignatures 收集签名
func (c *Coordinator) collectSignatures(agreement *Agreement) error {
	// 遍历所有参与者收集签名
	for processID, process := range c.state.processes {
		if process.State != "active" {
			continue
		}

		signature := Signature{
			ProcessID: processID,
			Timestamp: time.Now(),
			Valid:     true,
		}

		agreement.Signatures[processID] = signature
	}

	// 检查签名数量是否达到法定人数
	if len(agreement.Signatures) < c.config.quorumSize {
		return model.WrapError(nil, model.ErrCodeConsensus,
			"insufficient signatures for agreement")
	}

	return nil
}

// finalizeConsensus 完成共识
func (c *Coordinator) finalizeConsensus(session *Session, agreement *Agreement) {
	// 更新协议状态
	agreement.State = "completed"

	// 更新会话状态
	session.State.Phase = "completed"
	session.State.Progress = 1.0

	// 更新指标
	c.state.metrics.SuccessRate = (c.state.metrics.SuccessRate*
		float64(c.state.metrics.ActiveSessions-1) + 1.0) /
		float64(c.state.metrics.ActiveSessions)
}

// updateProcessStates 更新进程状态
func (c *Coordinator) updateProcessStates() error {
	// 遍历所有进程
	for id, process := range c.state.processes {
		// 检查进程是否活跃
		if err := c.checkProcessStatus(process); err != nil {
			// 移除不活跃的进程
			delete(c.state.processes, id)
			continue
		}

		// 更新进程依赖状态
		if err := c.updateProcessDependencies(process); err != nil {
			continue
		}

		// 更新进程最后更新时间
		process.LastUpdate = time.Now()
	}

	return nil
}

// resolveConflicts 解决冲突
func (c *Coordinator) resolveConflicts() error {
	for _, session := range c.state.sessions {
		// 检查会话中的冲突
		if len(session.State.Conflicts) == 0 {
			continue
		}

		// 解决每个冲突
		for _, conflict := range session.State.Conflicts {
			// 使用解决器解决冲突
			if err := c.resolver.ResolveConflict(&conflict); err != nil {
				// 记录冲突解决失败
				c.recordConflictResolution(session, &conflict, err)
				continue
			}

			// 更新会话状态
			if err := c.updateSessionAfterResolution(session, &conflict); err != nil {
				continue
			}
		}

		// 清空已解决的冲突
		session.State.Conflicts = make([]Conflict, 0)
	}

	return nil
}

// 辅助函数

func (c *Coordinator) validateProcess(process *Process) error {
	if process.ID == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty process ID")
	}

	if process.Type == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty process type")
	}

	return nil
}

func (c *Coordinator) updateMetrics() {
	point := types.MetricPoint{ // 使用 types.MetricPoint 而不是 MetricPoint
		Timestamp: time.Now(),
		Values: map[string]float64{
			"active_sessions": float64(len(c.state.sessions)),
			"success_rate":    c.calculateSuccessRate(),
			"conflict_rate":   c.calculateConflictRate(),
		},
	}

	c.state.metrics.History = append(c.state.metrics.History, point)

	// 限制历史记录数量
	if len(c.state.metrics.History) > types.MaxMetricsHistory {
		c.state.metrics.History = c.state.metrics.History[1:]
	}
}

// 辅助函数 - 添加到文件末尾现有辅助函数部分

// calculateSuccessRate 计算成功率
func (c *Coordinator) calculateSuccessRate() float64 {
	if c.state.metrics.ActiveSessions == 0 {
		return 1.0 // 无活跃会话时返回完美分数
	}

	successCount := 0
	totalSessions := 0

	// 统计会话完成情况
	for _, session := range c.state.sessions {
		if session.State.Phase == "completed" {
			successCount++
		}
		totalSessions++
	}

	// 从历史指标中获取之前的成功会话数
	for _, point := range c.state.metrics.History {
		if successRate, exists := point.Values["success_rate"]; exists {
			// 使用历史数据加权计算
			return (successRate*0.7 + float64(successCount)/float64(totalSessions)*0.3)
		}
	}

	return float64(successCount) / float64(totalSessions)
}

// calculateConflictRate 计算冲突率
func (c *Coordinator) calculateConflictRate() float64 {
	if c.state.metrics.ActiveSessions == 0 {
		return 0.0 // 无活跃会话时返回零冲突
	}

	conflictCount := 0
	totalDecisions := 0

	// 统计所有会话中的冲突
	for _, session := range c.state.sessions {
		conflictCount += len(session.State.Conflicts)
		totalDecisions += len(session.State.Decisions)
	}

	if totalDecisions == 0 {
		return 0.0
	}

	// 计算当前冲突率
	currentRate := float64(conflictCount) / float64(totalDecisions)

	// 从历史指标中获取之前的冲突率
	for _, point := range c.state.metrics.History {
		if conflictRate, exists := point.Values["conflict_rate"]; exists {
			// 使用历史数据平滑计算
			return (conflictRate*0.8 + currentRate*0.2)
		}
	}

	return currentRate
}
func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}

// checkProcessStatus 检查进程状态
func (c *Coordinator) checkProcessStatus(process *Process) error {
	// 检查最后更新时间
	if time.Since(process.LastUpdate) > c.config.consensusTimeout {
		return model.WrapError(nil, model.ErrCodeTimeout, "process timeout")
	}

	return nil
}

// updateProcessDependencies 更新进程依赖状态
func (c *Coordinator) updateProcessDependencies(process *Process) error {
	for _, depID := range process.Dependencies {
		if _, exists := c.state.processes[depID]; !exists {
			return model.WrapError(nil, model.ErrCodeDependency, "dependency not found")
		}
	}
	return nil
}

// recordConflictResolution 记录冲突解决结果
func (c *Coordinator) recordConflictResolution(session *Session, conflict *Conflict, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 更新会话统计
	session.State.Phase = "resolving"

	// 创建解决记录
	resolution := &types.ResolutionRecord{
		ConflictID:   conflict.ID,
		SessionID:    session.ID,
		ResolvedAt:   time.Now(),
		Status:       "failed",
		Error:        err,
		Participants: session.Participants,
		Details: map[string]interface{}{
			"conflict_type": conflict.Type,
			"priority":      conflict.Priority,
			"phase":         session.State.Phase,
			"progress":      session.State.Progress,
		},
	}

	// 根据错误更新状态
	if err == nil {
		resolution.Status = "resolved"
		// 更新指标
		c.state.metrics.ConflictRate = (c.state.metrics.ConflictRate *
			float64(c.state.metrics.ActiveSessions-1)) /
			float64(c.state.metrics.ActiveSessions)
	} else {
		// 记录错误
		resolution.ErrorDetails = &types.ErrorDetails{
			Code:    types.ErrorCodeToString(model.GetErrorCode(err)), // Use conversion
			Message: err.Error(),
			Time:    time.Now(),
		}
		// 更新失败指标
		c.state.metrics.ConflictRate = (c.state.metrics.ConflictRate*
			float64(c.state.metrics.ActiveSessions-1) + 1.0) /
			float64(c.state.metrics.ActiveSessions)
	}

	// 添加到历史记录
	point := types.MetricPoint{
		Timestamp: time.Now(),
		Values: map[string]float64{
			"conflict_resolved": types.BoolToFloat64(resolution.Status == "resolved"), // Use conversion
			"resolution_time":   float64(resolution.ResolvedAt.Sub(session.StartTime).Seconds()),
		},
	}
	c.state.metrics.History = append(c.state.metrics.History, point)

	// 限制历史记录数量
	if len(c.state.metrics.History) > types.MaxMetricsHistory {
		c.state.metrics.History = c.state.metrics.History[1:]
	}
}

// updateSessionAfterResolution 更新会话状态
func (c *Coordinator) updateSessionAfterResolution(session *Session, conflict *Conflict) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 更新冲突状态
	for i, existingConflict := range session.State.Conflicts {
		if existingConflict.ID == conflict.ID {
			session.State.Conflicts[i].Status = "resolved"
			session.State.Conflicts[i].Resolution = "auto_resolved"
			break
		}
	}

	// 更新会话进度
	newProgress := calculateProgress(session)
	if newProgress != session.State.Progress {
		session.State.Progress = newProgress

		// 根据进度更新阶段
		switch {
		case newProgress >= 1.0:
			session.State.Phase = "completed"
		case len(session.State.Conflicts) > 0:
			session.State.Phase = "resolving"
		case newProgress >= 0.5:
			session.State.Phase = "consensus"
		default:
			session.State.Phase = "collecting"
		}
	}

	// 检查是否所有冲突都已解决
	allResolved := true
	for _, c := range session.State.Conflicts {
		if c.Status != "resolved" {
			allResolved = false
			break
		}
	}

	// 如果所有冲突已解决且进度足够，进入共识阶段
	if allResolved && session.State.Progress >= 0.5 {
		session.State.Phase = "consensus"
	}

	// 更新会话指标
	c.updateSessionMetrics(session)

	return nil
}

// updateSessionMetrics 更新会话相关指标
func (c *Coordinator) updateSessionMetrics(session *Session) {
	totalConflicts := len(session.State.Conflicts)
	resolvedConflicts := 0
	for _, conflict := range session.State.Conflicts {
		if conflict.Status == "resolved" {
			resolvedConflicts++
		}
	}

	// 更新解决率
	if totalConflicts > 0 {
		resolutionRate := float64(resolvedConflicts) / float64(totalConflicts)
		c.state.metrics.SuccessRate = (c.state.metrics.SuccessRate*
			float64(c.state.metrics.ActiveSessions-1) + resolutionRate) /
			float64(c.state.metrics.ActiveSessions)
	}
}

// calculateProgress 计算会话进度
func calculateProgress(session *Session) float64 {
	totalDecisions := len(session.State.Decisions)
	if totalDecisions == 0 {
		return 0.0
	}

	completedDecisions := 0
	for _, decision := range session.State.Decisions {
		if len(decision.Votes) >= 2 { // 假设需要至少2个投票
			completedDecisions++
		}
	}

	return float64(completedDecisions) / float64(totalDecisions)
}
