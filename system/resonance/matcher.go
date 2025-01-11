//system/resonance/matcher.go

package resonance

import (
    "math"
    "sync"
    "time"

    "github.com/Corphon/daoflow/meta/field"
    "github.com/Corphon/daoflow/meta/emergence"
    "github.com/Corphon/daoflow/model"
    "github.com/Corphon/daoflow/system/types"
)

// PatternMatcher 模式匹配器
type PatternMatcher struct {
    mu sync.RWMutex

    // 基础配置
    config struct {
        matchThreshold float64         // 匹配阈值
        minSimilarity float64         // 最小相似度
        maxDistance   float64         // 最大距离
    }

    // 匹配状态
    state struct {
        matches     map[string]*MatchState    // 当前匹配
        templates   map[string]*MatchTemplate // 匹配模板
        history     []MatchEvent             // 匹配历史
    }

    // 依赖项
    detector  *emergence.PatternDetector
    amplifier *ResonanceAmplifier
}

// MatchState 匹配状态
type MatchState struct {
    ID          string                // 匹配ID
    Template    *MatchTemplate        // 使用的模板
    Pattern     *emergence.EmergentPattern // 匹配的模式
    Similarity  float64               // 相似度
    Confidence  float64               // 置信度
    StartTime   time.Time            // 开始时间
    LastUpdate  time.Time            // 最后更新时间
    Properties  map[string]float64   // 匹配属性
}

// MatchTemplate 匹配模板
type MatchTemplate struct {
    ID          string               // 模板ID
    Type        string               // 模板类型
    Features    []TemplateFeature    // 特征列表
    Weights     map[string]float64   // 特征权重
    Constraints []MatchConstraint    // 匹配约束
    Created     time.Time           // 创建时间
}

// TemplateFeature 模板特征
type TemplateFeature struct {
    Name       string               // 特征名称
    Type       string               // 特征类型
    Value      interface{}          // 特征值
    Tolerance  float64             // 容差
}

// MatchConstraint 匹配约束
type MatchConstraint struct {
    Type      string               // 约束类型
    Target    string               // 约束目标
    Operator  string               // 约束操作符
    Value     interface{}          // 约束值
}

// MatchEvent 匹配事件
type MatchEvent struct {
    Timestamp  time.Time
    MatchID    string
    Type       string
    Template   string
    Pattern    string
    Similarity float64
    Success    bool
}

// NewPatternMatcher 创建新的模式匹配器
func NewPatternMatcher(
    detector *emergence.PatternDetector,
    amplifier *ResonanceAmplifier) *PatternMatcher {
    
    pm := &PatternMatcher{
        detector:  detector,
        amplifier: amplifier,
    }

    // 初始化配置
    pm.config.matchThreshold = 0.75
    pm.config.minSimilarity = 0.6
    pm.config.maxDistance = 0.4

    // 初始化状态
    pm.state.matches = make(map[string]*MatchState)
    pm.state.templates = make(map[string]*MatchTemplate)
    pm.state.history = make([]MatchEvent, 0)

    return pm
}

// Match 执行模式匹配
func (pm *PatternMatcher) Match() error {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    // 获取当前模式
    patterns, err := pm.detector.Detect()
    if err != nil {
        return err
    }

    // 对每个模式进行匹配
    for _, pattern := range patterns {
        matches := pm.matchPattern(pattern)
        
        // 更新匹配状态
        pm.updateMatches(pattern, matches)
    }

    // 清理过期匹配
    pm.cleanupMatches()

    return nil
}

// RegisterTemplate 注册匹配模板
func (pm *PatternMatcher) RegisterTemplate(template *MatchTemplate) error {
    if template == nil {
        return model.WrapError(nil, model.ErrCodeValidation, "nil template")
    }

    pm.mu.Lock()
    defer pm.mu.Unlock()

    // 验证模板
    if err := pm.validateTemplate(template); err != nil {
        return err
    }

    // 存储模板
    pm.state.templates[template.ID] = template

    return nil
}

// matchPattern 匹配单个模式
func (pm *PatternMatcher) matchPattern(
    pattern emergence.EmergentPattern) []*MatchState {
    
    matches := make([]*MatchState, 0)

    // 对每个模板进行匹配
    for _, template := range pm.state.templates {
        if match := pm.matchAgainstTemplate(pattern, template); match != nil {
            matches = append(matches, match)
        }
    }

    return matches
}

// matchAgainstTemplate 与模板进行匹配
func (pm *PatternMatcher) matchAgainstTemplate(
    pattern emergence.EmergentPattern,
    template *MatchTemplate) *MatchState {
    
    // 计算特征相似度
    similarity := pm.calculateSimilarity(pattern, template)
    if similarity < pm.config.minSimilarity {
        return nil
    }

    // 检查约束条件
    if !pm.checkConstraints(pattern, template) {
        return nil
    }

    // 创建匹配状态
    match := &MatchState{
        ID:         generateMatchID(),
        Template:   template,
        Pattern:    &pattern,
        Similarity: similarity,
        Confidence: pm.calculateConfidence(similarity, pattern, template),
        StartTime:  time.Now(),
        LastUpdate: time.Now(),
        Properties: make(map[string]float64),
    }

    // 提取匹配属性
    match.Properties = pm.extractMatchProperties(pattern, template)

    return match
}

// calculateSimilarity 计算相似度
func (pm *PatternMatcher) calculateSimilarity(
    pattern emergence.EmergentPattern,
    template *MatchTemplate) float64 {
    
    totalWeight := 0.0
    weightedSimilarity := 0.0

    // 计算每个特征的相似度
    for _, feature := range template.Features {
        weight := template.Weights[feature.Name]
        similarity := pm.calculateFeatureSimilarity(pattern, feature)
        
        weightedSimilarity += weight * similarity
        totalWeight += weight
    }

    if totalWeight == 0 {
        return 0
    }

    return weightedSimilarity / totalWeight
}

// calculateFeatureSimilarity 计算特征相似度
func (pm *PatternMatcher) calculateFeatureSimilarity(
    pattern emergence.EmergentPattern,
    feature TemplateFeature) float64 {
    
    switch feature.Type {
    case "numeric":
        return pm.calculateNumericSimilarity(pattern, feature)
    case "categorical":
        return pm.calculateCategoricalSimilarity(pattern, feature)
    case "structural":
        return pm.calculateStructuralSimilarity(pattern, feature)
    default:
        return 0
    }
}

// checkConstraints 检查约束条件
func (pm *PatternMatcher) checkConstraints(
    pattern emergence.EmergentPattern,
    template *MatchTemplate) bool {
    
    for _, constraint := range template.Constraints {
        if !pm.evaluateConstraint(pattern, constraint) {
            return false
        }
    }

    return true
}

// updateMatches 更新匹配状态
func (pm *PatternMatcher) updateMatches(
    pattern emergence.EmergentPattern,
    newMatches []*MatchState) {
    
    // 记录匹配事件
    for _, match := range newMatches {
        event := MatchEvent{
            Timestamp: time.Now(),
            MatchID:   match.ID,
            Type:     "new_match",
            Template: match.Template.ID,
            Pattern:  pattern.ID,
            Similarity: match.Similarity,
            Success:   true,
        }
        pm.recordMatchEvent(event)

        // 更新或添加匹配状态
        pm.state.matches[match.ID] = match
    }
}

// cleanupMatches 清理过期匹配
func (pm *PatternMatcher) cleanupMatches() {
    threshold := time.Now().Add(-matchTimeout)

    for id, match := range pm.state.matches {
        if match.LastUpdate.Before(threshold) {
            delete(pm.state.matches, id)
            
            // 记录过期事件
            event := MatchEvent{
                Timestamp: time.Now(),
                MatchID:   id,
                Type:     "expired",
                Template: match.Template.ID,
                Pattern:  match.Pattern.ID,
                Success:  false,
            }
            pm.recordMatchEvent(event)
        }
    }
}

// 辅助函数

func (pm *PatternMatcher) validateTemplate(template *MatchTemplate) error {
    if len(template.Features) == 0 {
        return model.WrapError(nil, model.ErrCodeValidation, "template has no features")
    }

    // 验证权重
    totalWeight := 0.0
    for _, weight := range template.Weights {
        if weight < 0 {
            return model.WrapError(nil, model.ErrCodeValidation, "negative weight")
        }
        totalWeight += weight
    }

    if math.Abs(totalWeight-1.0) > 1e-6 {
        return model.WrapError(nil, model.ErrCodeValidation, "weights must sum to 1")
    }

    return nil
}

func generateMatchID() string {
    return fmt.Sprintf("match_%d", time.Now().UnixNano())
}

func (pm *PatternMatcher) recordMatchEvent(event MatchEvent) {
    pm.state.history = append(pm.state.history, event)

    // 限制历史记录长度
    if len(pm.state.history) > maxHistoryLength {
        pm.state.history = pm.state.history[1:]
    }
}

const (
    maxHistoryLength = 1000
    matchTimeout    = 1 * time.Hour
)
