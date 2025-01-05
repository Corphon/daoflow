// api/pattern.go

package api

import (
    "context"
    "sync"
    "time"

    "github.com/Corphon/daoflow/system"
)

// PatternType 模式类型
type PatternType string

const (
    TypeBehavior  PatternType = "behavior"  // 行为模式
    TypeResource  PatternType = "resource"  // 资源模式
    TypeEnergy    PatternType = "energy"    // 能量模式
    TypeAnomaly   PatternType = "anomaly"   // 异常模式
    TypeCycle     PatternType = "cycle"     // 周期模式
)

// Pattern 模式定义
type Pattern struct {
    ID          string                 `json:"id"`           // 模式ID
    Type        PatternType            `json:"type"`         // 模式类型
    Name        string                 `json:"name"`         // 模式名称
    Features    map[string]float64     `json:"features"`     // 特征向量
    Weights     map[string]float64     `json:"weights"`      // 特征权重
    Threshold   float64                `json:"threshold"`    // 匹配阈值
    Confidence  float64                `json:"confidence"`   // 置信度
    CreateTime  time.Time              `json:"create_time"`  // 创建时间
    UpdateTime  time.Time              `json:"update_time"`  // 更新时间
    Metadata    map[string]interface{} `json:"metadata"`     // 元数据
}

// PatternMatch 模式匹配结果
type PatternMatch struct {
    PatternID  string    `json:"pattern_id"`   // 模式ID
    Score      float64   `json:"score"`        // 匹配分数
    Timestamp  time.Time `json:"timestamp"`    // 匹配时间
    Details    map[string]interface{} `json:"details"` // 匹配细节
}

// PatternStats 模式统计
type PatternStats struct {
    TotalPatterns    int                `json:"total_patterns"`     // 模式总数
    ActivePatterns   int                `json:"active_patterns"`    // 活跃模式数
    MatchRate        float64            `json:"match_rate"`         // 匹配率
    AverageScore     float64            `json:"average_score"`      // 平均分数
    TypeDistribution map[PatternType]int `json:"type_distribution"` // 类型分布
}

// PatternAPI 模式识别API
type PatternAPI struct {
    mu     sync.RWMutex
    system *system.SystemCore
    
    // 模式库
    patterns map[string]*Pattern
    
    // 学习参数
    learningRate float64
    minConfidence float64
    
    // 事件通道
    events chan PatternEvent
    
    ctx    context.Context
    cancel context.CancelFunc
}

// PatternEvent 模式事件
type PatternEvent struct {
    Type      string      `json:"type"`       // 事件类型
    Pattern   *Pattern    `json:"pattern"`    // 相关模式
    Match     *PatternMatch `json:"match"`    // 匹配结果
    Timestamp time.Time   `json:"timestamp"`  // 事件时间
}

// NewPatternAPI 创建模式API实例
func NewPatternAPI(sys *system.SystemCore, opts *Options) *PatternAPI {
    ctx, cancel := context.WithCancel(context.Background())
    
    api := &PatternAPI{
        system:       sys,
        patterns:     make(map[string]*Pattern),
        learningRate: 0.1,
        minConfidence: 0.6,
        events:       make(chan PatternEvent, 100),
        ctx:          ctx,
        cancel:       cancel,
    }
    
    go api.learn()
    return api
}

// RegisterPattern 注册新模式
func (p *PatternAPI) RegisterPattern(pattern *Pattern) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    if pattern.ID == "" {
        return NewError(ErrInvalidPattern, "pattern ID is required")
    }

    if !isValidPatternType(pattern.Type) {
        return NewError(ErrInvalidPattern, "invalid pattern type")
    }

    pattern.CreateTime = time.Now()
    pattern.UpdateTime = pattern.CreateTime
    p.patterns[pattern.ID] = pattern

    p.events <- PatternEvent{
        Type:      "pattern_registered",
        Pattern:   pattern,
        Timestamp: time.Now(),
    }

    return nil
}

// MatchPattern 执行模式匹配
func (p *PatternAPI) MatchPattern(ctx context.Context, features map[string]float64) ([]*PatternMatch, error) {
    p.mu.RLock()
    defer p.mu.RUnlock()

    var matches []*PatternMatch
    
    for _, pattern := range p.patterns {
        score := p.calculateMatchScore(pattern, features)
        
        if score >= pattern.Threshold {
            match := &PatternMatch{
                PatternID: pattern.ID,
                Score:    score,
                Timestamp: time.Now(),
                Details: map[string]interface{}{
                    "features": features,
                    "threshold": pattern.Threshold,
                },
            }
            matches = append(matches, match)
            
            p.events <- PatternEvent{
                Type:   "pattern_matched",
                Pattern: pattern,
                Match:  match,
                Timestamp: time.Now(),
            }
        }
    }

    return matches, nil
}

// UpdatePattern 更新模式
func (p *PatternAPI) UpdatePattern(id string, updates map[string]interface{}) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    pattern, exists := p.patterns[id]
    if !exists {
        return NewError(ErrPatternNotFound, "pattern not found")
    }

    // 应用更新
    if features, ok := updates["features"].(map[string]float64); ok {
        pattern.Features = features
    }
    if weights, ok := updates["weights"].(map[string]float64); ok {
        pattern.Weights = weights
    }
    if threshold, ok := updates["threshold"].(float64); ok {
        pattern.Threshold = threshold
    }

    pattern.UpdateTime = time.Now()

    p.events <- PatternEvent{
        Type:      "pattern_updated",
        Pattern:   pattern,
        Timestamp: time.Now(),
    }

    return nil
}

// GetPattern 获取模式信息
func (p *PatternAPI) GetPattern(id string) (*Pattern, error) {
    p.mu.RLock()
    defer p.mu.RUnlock()

    pattern, exists := p.patterns[id]
    if !exists {
        return nil, NewError(ErrPatternNotFound, "pattern not found")
    }

    return pattern, nil
}

// GetStats 获取模式统计信息
func (p *PatternAPI) GetStats() (*PatternStats, error) {
    p.mu.RLock()
    defer p.mu.RUnlock()

    stats := &PatternStats{
        TotalPatterns: len(p.patterns),
        TypeDistribution: make(map[PatternType]int),
    }

    var totalScore float64
    activeCount := 0

    for _, pattern := range p.patterns {
        stats.TypeDistribution[pattern.Type]++
        if pattern.Confidence >= p.minConfidence {
            activeCount++
        }
        totalScore += pattern.Confidence
    }

    stats.ActivePatterns = activeCount
    if len(p.patterns) > 0 {
        stats.AverageScore = totalScore / float64(len(p.patterns))
    }

    return stats, nil
}

// Subscribe 订阅模式事件
func (p *PatternAPI) Subscribe() (<-chan PatternEvent, error) {
    return p.events, nil
}

// learn 持续学习协程
func (p *PatternAPI) learn() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-p.ctx.Done():
            return
        case <-ticker.C:
            p.updatePatternConfidence()
        }
    }
}

// calculateMatchScore 计算匹配分数
func (p *PatternAPI) calculateMatchScore(pattern *Pattern, features map[string]float64) float64 {
    var score float64
    var totalWeight float64

    for feature, weight := range pattern.Weights {
        if value, exists := features[feature]; exists {
            score += weight * value
        }
        totalWeight += weight
    }

    if totalWeight > 0 {
        return score / totalWeight
    }
    return 0
}

// updatePatternConfidence 更新模式置信度
func (p *PatternAPI) updatePatternConfidence() {
    p.mu.Lock()
    defer p.mu.Unlock()

    // 实现模式置信度更新逻辑
}

// Close 关闭API
func (p *PatternAPI) Close() error {
    p.cancel()
    close(p.events)
    return nil
}

// isValidPatternType 验证模式类型
func isValidPatternType(pt PatternType) bool {
    switch pt {
    case TypeBehavior, TypeResource, TypeEnergy, TypeAnomaly, TypeCycle:
        return true
    default:
        return false
    }
}
