//system/control/sync/resolver.go

package sync

import (
    "sync"
    "time"

    "github.com/Corphon/daoflow/meta/field"
    "github.com/Corphon/daoflow/model"
    "github.com/Corphon/daoflow/system/types"
)

// Resolver 解决器
type Resolver struct {
    mu sync.RWMutex

    // 基础配置
    config struct {
        resolutionTimeout time.Duration // 解决超时
        maxAttempts      int           // 最大尝试次数
        minConfidence    float64       // 最小置信度
        autoResolve      bool          // 自动解决
    }

    // 解决状态
    state struct {
        conflicts   map[string]*Conflict     // 冲突列表
        strategies  map[string]*Strategy     // 策略列表
        resolutions map[string]*Resolution   // 解决方案
        metrics     ResolutionMetrics       // 解决指标
    }
}

// Conflict 冲突信息
type Conflict struct {
    ID           string                // 冲突ID
    Type         string                // 冲突类型
    Status       string                // 冲突状态
    Priority     int                   // 优先级
    Parties      []Party               // 冲突方
    Resources    []Resource            // 相关资源
    Created      time.Time            // 创建时间
    Updated      time.Time            // 更新时间
}

// Party 冲突方
type Party struct {
    ID           string                // 参与方ID
    Type         string                // 参与方类型
    Role         string                // 参与角色
    Position     interface{}           // 立场信息
    Requirements []Requirement         // 需求列表
}

// Resource 相关资源
type Resource struct {
    ID           string                // 资源ID
    Type         string                // 资源类型
    State        string                // 资源状态
    Constraints  []Constraint          // 资源约束
    Dependencies []string              // 依赖资源
}

// Requirement 需求信息
type Requirement struct {
    ID           string                // 需求ID
    Type         string                // 需求类型
    Priority     int                   // 优先级
    Constraints  []Constraint          // 需求约束
    Flexibility  float64               // 灵活度
}

// Strategy 解决策略
type Strategy struct {
    ID           string                // 策略ID
    Type         string                // 策略类型
    Priority     int                   // 优先级
    Conditions   []Condition           // 应用条件
    Actions      []Action              // 策略动作
    Success      float64               // 成功率
}

// Condition 应用条件
type Condition struct {
    Type         string                // 条件类型
    Value        interface{}           // 条件值
    Operator     string                // 操作符
    Weight       float64               // 权重
}

// Action 策略动作
type Action struct {
    Type         string                // 动作类型
    Target       string                // 目标对象
    Operation    string                // 操作类型
    Parameters   map[string]interface{} // 操作参数
}

// Resolution 解决方案
type Resolution struct {
    ID           string                // 方案ID
    ConflictID   string                // 冲突ID
    Type         string                // 方案类型
    Status       string                // 方案状态
    Steps        []ResolutionStep      // 解决步骤
    Results      ResolutionResult      // 解决结果
    Created      time.Time            // 创建时间
}

// ResolutionStep 解决步骤
type ResolutionStep struct {
    ID           string                // 步骤ID
    Type         string                // 步骤类型
    Action       Action                // 执行动作
    Status       string                // 步骤状态
    StartTime    time.Time            // 开始时间
    EndTime      time.Time            // 结束时间
}

// ResolutionResult 解决结果
type ResolutionResult struct {
    Success      bool                  // 是否成功
    Confidence   float64               // 置信度
    Impact       map[string]float64    // 影响评估
    Feedback     []string              // 反馈信息
}

// ResolutionMetrics 解决指标
type ResolutionMetrics struct {
    TotalConflicts  int               // 总冲突数
    ResolvedCount   int               // 已解决数
    SuccessRate     float64           // 成功率
    AverageTime     time.Duration     // 平均耗时
    History         []MetricPoint     // 历史指标
}

// NewResolver 创建新的解决器
func NewResolver() *Resolver {
    r := &Resolver{}

    // 初始化配置
    r.config.resolutionTimeout = 30 * time.Second
    r.config.maxAttempts = 3
    r.config.minConfidence = 0.75
    r.config.autoResolve = true

    // 初始化状态
    r.state.conflicts = make(map[string]*Conflict)
    r.state.strategies = make(map[string]*Strategy)
    r.state.resolutions = make(map[string]*Resolution)
    r.state.metrics = ResolutionMetrics{
        History: make([]MetricPoint, 0),
    }

    return r
}

// RegisterConflict 注册冲突
func (r *Resolver) RegisterConflict(conflict *Conflict) error {
    if conflict == nil {
        return model.WrapError(nil, model.ErrCodeValidation, "nil conflict")
    }

    r.mu.Lock()
    defer r.mu.Unlock()

    // 验证冲突
    if err := r.validateConflict(conflict); err != nil {
        return err
    }

    // 存储冲突
    r.state.conflicts[conflict.ID] = conflict

    // 如果开启自动解决，尝试解决冲突
    if r.config.autoResolve {
        go r.resolveConflict(conflict)
    }

    return nil
}

// RegisterStrategy 注册策略
func (r *Resolver) RegisterStrategy(strategy *Strategy) error {
    if strategy == nil {
        return model.WrapError(nil, model.ErrCodeValidation, "nil strategy")
    }

    r.mu.Lock()
    defer r.mu.Unlock()

    // 验证策略
    if err := r.validateStrategy(strategy); err != nil {
        return err
    }

    // 存储策略
    r.state.strategies[strategy.ID] = strategy

    return nil
}

// Resolve 解决冲突
func (r *Resolver) Resolve(conflictID string) (*Resolution, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    // 获取冲突
    conflict, exists := r.state.conflicts[conflictID]
    if !exists {
        return nil, model.WrapError(nil, model.ErrCodeNotFound, "conflict not found")
    }

    // 创建解决方案
    resolution := &Resolution{
        ID:         generateResolutionID(),
        ConflictID: conflictID,
        Type:       determineResolutionType(conflict),
        Status:     "initiated",
        Steps:      make([]ResolutionStep, 0),
        Created:    time.Now(),
    }

    // 执行解决方案
    if err := r.executeResolution(resolution); err != nil {
        return nil, err
    }

    // 存储解决方案
    r.state.resolutions[resolution.ID] = resolution

    // 更新指标
    r.updateMetrics()

    return resolution, nil
}

// 辅助函数

func (r *Resolver) validateConflict(conflict *Conflict) error {
    if conflict.ID == "" {
        return model.WrapError(nil, model.ErrCodeValidation, "empty conflict ID")
    }

    if len(conflict.Parties) == 0 {
        return model.WrapError(nil, model.ErrCodeValidation, "no conflict parties")
    }

    return nil
}

func (r *Resolver) validateStrategy(strategy *Strategy) error {
    if strategy.ID == "" {
        return model.WrapError(nil, model.ErrCodeValidation, "empty strategy ID")
    }

    if len(strategy.Actions) == 0 {
        return model.WrapError(nil, model.ErrCodeValidation, "no strategy actions")
    }

    return nil
}

func (r *Resolver) updateMetrics() {
    point := MetricPoint{
        Timestamp: time.Now(),
        Values: map[string]float64{
            "total_conflicts":  float64(len(r.state.conflicts)),
            "resolved_count":   float64(r.state.metrics.ResolvedCount),
            "success_rate":     r.state.metrics.SuccessRate,
        },
    }

    r.state.metrics.History = append(r.state.metrics.History, point)

    // 限制历史记录数量
    if len(r.state.metrics.History) > maxMetricsHistory {
        r.state.metrics.History = r.state.metrics.History[1:]
    }
}

func generateResolutionID() string {
    return fmt.Sprintf("res_%d", time.Now().UnixNano())
}

const (
    maxMetricsHistory = 1000
)
