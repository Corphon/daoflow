//system/control/state/validator.go

package state

import (
    "sync"
    "time"

    "github.com/Corphon/daoflow/meta/field"
    "github.com/Corphon/daoflow/model"
    "github.com/Corphon/daoflow/system/types"
)

// StateValidator 状态验证器
type StateValidator struct {
    mu sync.RWMutex

    // 基础配置
    config struct {
        validationTimeout time.Duration // 验证超时时间
        strictMode       bool          // 严格模式
        maxRules        int           // 最大规则数
        cacheEnabled    bool          // 启用缓存
    }

    // 验证状态
    state struct {
        rules      map[string]*ValidationRule   // 验证规则
        cache      map[string]*ValidationCache  // 验证缓存
        results    []ValidationResult           // 验证结果
        metrics    ValidationMetrics            // 验证指标
    }
}

// ValidationRule 验证规则
type ValidationRule struct {
    ID           string                // 规则ID
    Name         string                // 规则名称
    Type         string                // 规则类型
    Target       string                // 验证目标
    Condition    RuleCondition         // 验证条件
    Priority     int                   // 优先级
    Enabled      bool                  // 是否启用
}

// RuleCondition 规则条件
type RuleCondition struct {
    Type         string                // 条件类型
    Expression   string                // 条件表达式
    Parameters   map[string]interface{} // 条件参数
    ErrorMsg     string                // 错误消息
}

// ValidationCache 验证缓存
type ValidationCache struct {
    StateID      string                // 状态ID
    RuleResults  map[string]bool       // 规则结果
    Timestamp    time.Time            // 缓存时间
    TTL          time.Duration        // 生存时间
}

// ValidationResult 验证结果
type ValidationResult struct {
    ID           string                // 结果ID
    StateID      string                // 状态ID
    RuleID       string                // 规则ID
    Valid        bool                  // 是否有效
    Details      map[string]interface{} // 详细信息
    Timestamp    time.Time            // 验证时间
}

// ValidationMetrics 验证指标
type ValidationMetrics struct {
    TotalValidations int              // 总验证次数
    SuccessRate     float64           // 成功率
    AverageLatency  time.Duration     // 平均延迟
    CacheHitRate    float64           // 缓存命中率
    History         []MetricPoint     // 历史指标
}

// MetricPoint 指标点
type MetricPoint struct {
    Timestamp    time.Time
    Values       map[string]float64
}

// NewStateValidator 创建新的状态验证器
func NewStateValidator() *StateValidator {
    sv := &StateValidator{}

    // 初始化配置
    sv.config.validationTimeout = 5 * time.Second
    sv.config.strictMode = true
    sv.config.maxRules = 100
    sv.config.cacheEnabled = true

    // 初始化状态
    sv.state.rules = make(map[string]*ValidationRule)
    sv.state.cache = make(map[string]*ValidationCache)
    sv.state.results = make([]ValidationResult, 0)
    sv.state.metrics = ValidationMetrics{
        History: make([]MetricPoint, 0),
    }

    return sv
}

// RegisterRule 注册验证规则
func (sv *StateValidator) RegisterRule(rule *ValidationRule) error {
    if rule == nil {
        return model.WrapError(nil, model.ErrCodeValidation, "nil rule")
    }

    sv.mu.Lock()
    defer sv.mu.Unlock()

    // 验证规则
    if err := sv.validateRule(rule); err != nil {
        return err
    }

    // 检查规则数量限制
    if len(sv.state.rules) >= sv.config.maxRules {
        return model.WrapError(nil, model.ErrCodeLimit, "max rules reached")
    }

    // 存储规则
    sv.state.rules[rule.ID] = rule

    return nil
}

// ValidateState 验证系统状态
func (sv *StateValidator) ValidateState(state *SystemState) error {
    sv.mu.RLock()
    defer sv.mu.RUnlock()

    startTime := time.Now()

    // 检查缓存
    if sv.config.cacheEnabled {
        if result := sv.checkCache(state.ID); result != nil {
            return sv.processValidationResult(result)
        }
    }

    // 执行验证
    results, err := sv.executeValidation(state)
    if err != nil {
        return err
    }

    // 更新缓存
    if sv.config.cacheEnabled {
        sv.updateCache(state.ID, results)
    }

    // 更新指标
    sv.updateMetrics(startTime)

    // 处理结果
    return sv.processResults(results)
}

// ValidateTransition 验证状态转换
func (sv *StateValidator) ValidateTransition(
    current, next *SystemState) error {
    
    sv.mu.RLock()
    defer sv.mu.RUnlock()

    // 验证转换规则
    for _, rule := range sv.state.rules {
        if !rule.Enabled {
            continue
        }

        if err := sv.validateTransitionRule(rule, current, next); err != nil {
            return err
        }
    }

    return nil
}

// executeValidation 执行验证
func (sv *StateValidator) executeValidation(
    state *SystemState) ([]ValidationResult, error) {
    
    results := make([]ValidationResult, 0)

    // 按优先级排序规则
    rules := sv.sortRulesByPriority()

    // 执行每个规则
    for _, rule := range rules {
        if !rule.Enabled {
            continue
        }

        result := sv.validateRule(rule, state)
        results = append(results, result)

        // 严格模式下，任何失败都立即返回
        if sv.config.strictMode && !result.Valid {
            return results, sv.createValidationError(result)
        }
    }

    return results, nil
}

// 辅助函数

func (sv *StateValidator) validateRule(rule *ValidationRule) error {
    if rule.ID == "" {
        return model.WrapError(nil, model.ErrCodeValidation, "empty rule ID")
    }

    if rule.Condition.Expression == "" {
        return model.WrapError(nil, model.ErrCodeValidation, "empty condition expression")
    }

    return nil
}

func (sv *StateValidator) updateMetrics(startTime time.Time) {
    duration := time.Since(startTime)

    point := MetricPoint{
        Timestamp: time.Now(),
        Values: map[string]float64{
            "latency": float64(duration.Milliseconds()),
            "success_rate": sv.calculateSuccessRate(),
            "cache_hit_rate": sv.calculateCacheHitRate(),
        },
    }

    sv.state.metrics.History = append(sv.state.metrics.History, point)

    // 限制历史记录数量
    if len(sv.state.metrics.History) > maxMetricsHistory {
        sv.state.metrics.History = sv.state.metrics.History[1:]
    }
}

func generateValidationID() string {
    return fmt.Sprintf("val_%d", time.Now().UnixNano())
}

const (
    maxMetricsHistory = 1000
)
