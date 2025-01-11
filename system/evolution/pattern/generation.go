//system/evolution/pattern/generation.go

package pattern

import (
    "sync"
    "time"

    "github.com/Corphon/daoflow/meta/field"
    "github.com/Corphon/daoflow/meta/emergence"
    "github.com/Corphon/daoflow/meta/resonance"
    "github.com/Corphon/daoflow/model"
    "github.com/Corphon/daoflow/system/types"
)

// PatternGenerator 模式生成器
type PatternGenerator struct {
    mu sync.RWMutex

    // 基础配置
    config struct {
        generationRate  float64         // 生成率
        mutationRate   float64         // 变异率
        complexityBias float64         // 复杂度偏好
        energyBalance  float64         // 能量平衡因子
    }

    // 生成状态
    state struct {
        templates    map[string]*GenerationTemplate  // 生成模板
        candidates   []*PatternCandidate            // 候选模式
        history      []GenerationEvent              // 生成历史
        metrics      GenerationMetrics              // 生成指标
    }

    // 依赖项
    recognizer *PatternRecognizer
    matcher    *EvolutionMatcher
}

// GenerationTemplate 生成模板
type GenerationTemplate struct {
    ID           string                // 模板ID
    Type         string                // 模板类型
    Structure    TemplateStructure     // 结构定义
    Constraints  []GenerationConstraint // 生成约束
    Properties   map[string]Range      // 属性范围
    Success      float64               // 成功率
    UsageCount   int                   // 使用次数
}

// TemplateStructure 模板结构
type TemplateStructure struct {
    Components   []ComponentSpec       // 组件规格
    Relations    []RelationSpec       // 关系规格
    Dynamics     DynamicsSpec         // 动态规格
}

// ComponentSpec 组件规格
type ComponentSpec struct {
    Type         string               // 组件类型
    Required     bool                 // 是否必需
    Properties   map[string]Range     // 属性范围
    Quantity     Range                // 数量范围
}

// RelationSpec 关系规格
type RelationSpec struct {
    Source       string               // 源组件类型
    Target       string               // 目标组件类型
    Type         string               // 关系类型
    Strength     Range                // 强度范围
}

// DynamicsSpec 动态规格
type DynamicsSpec struct {
    TimeScale    Range                // 时间尺度
    Evolution    []EvolutionRule      // 演化规则
    Stability    Range                // 稳定性范围
}

// Range 数值范围
type Range struct {
    Min          float64
    Max          float64
    Preferred    float64
}

// GenerationConstraint 生成约束
type GenerationConstraint struct {
    Type         string               // 约束类型
    Target       string               // 约束目标
    Condition    string               // 约束条件
    Value        interface{}          // 约束值
}

// PatternCandidate 候选模式
type PatternCandidate struct {
    ID           string               // 候选ID
    Template     string               // 使用的模板
    Pattern      *emergence.EmergentPattern // 生成的模式
    Score        float64              // 评分
    Generation   int                  // 生成代数
    Created      time.Time            // 创建时间
}

// GenerationEvent 生成事件
type GenerationEvent struct {
    Timestamp    time.Time
    TemplateID   string
    PatternID    string
    Success      bool
    Score        float64
    Changes      map[string]interface{}
}

// GenerationMetrics 生成指标
type GenerationMetrics struct {
    TotalGenerated  int
    SuccessRate     float64
    AverageScore    float64
    ComplexityDist  map[float64]int
    Evolution       []MetricPoint
}

// MetricPoint 指标点
type MetricPoint struct {
    Timestamp    time.Time
    Metrics      map[string]float64
}

// NewPatternGenerator 创建新的模式生成器
func NewPatternGenerator(
    recognizer *PatternRecognizer,
    matcher *EvolutionMatcher) *PatternGenerator {
    
    pg := &PatternGenerator{
        recognizer: recognizer,
        matcher:    matcher,
    }

    // 初始化配置
    pg.config.generationRate = 0.3
    pg.config.mutationRate = 0.1
    pg.config.complexityBias = 0.4
    pg.config.energyBalance = 0.7

    // 初始化状态
    pg.state.templates = make(map[string]*GenerationTemplate)
    pg.state.candidates = make([]*PatternCandidate, 0)
    pg.state.history = make([]GenerationEvent, 0)
    pg.state.metrics = GenerationMetrics{
        ComplexityDist: make(map[float64]int),
        Evolution:      make([]MetricPoint, 0),
    }

    return pg
}

// Generate 生成新模式
func (pg *PatternGenerator) Generate() error {
    pg.mu.Lock()
    defer pg.mu.Unlock()

    // 选择模板
    template := pg.selectTemplate()
    if template == nil {
        return model.WrapError(nil, model.ErrCodeOperation, "no suitable template")
    }

    // 生成候选模式
    candidates := pg.generateCandidates(template)

    // 评估候选模式
    evaluated := pg.evaluateCandidates(candidates)

    // 选择最佳候选
    selected := pg.selectBestCandidates(evaluated)

    // 优化选中的模式
    optimized := pg.optimizePatterns(selected)

    // 更新生成指标
    pg.updateMetrics(optimized)

    return nil
}

// RegisterTemplate 注册生成模板
func (pg *PatternGenerator) RegisterTemplate(template *GenerationTemplate) error {
    if template == nil {
        return model.WrapError(nil, model.ErrCodeValidation, "nil template")
    }

    pg.mu.Lock()
    defer pg.mu.Unlock()

    // 验证模板
    if err := pg.validateTemplate(template); err != nil {
        return err
    }

    // 存储模板
    pg.state.templates[template.ID] = template

    return nil
}

// generateCandidates 生成候选模式
func (pg *PatternGenerator) generateCandidates(
    template *GenerationTemplate) []*PatternCandidate {
    
    candidates := make([]*PatternCandidate, 0)

    // 生成多个候选
    for i := 0; i < maxCandidates; i++ {
        // 构建基础模式
        pattern := pg.buildBasePattern(template)
        
        // 应用变异
        if rand.Float64() < pg.config.mutationRate {
            pattern = pg.mutatePattern(pattern)
        }

        // 创建候选
        candidate := &PatternCandidate{
            ID:         generateCandidateID(),
            Template:   template.ID,
            Pattern:    pattern,
            Generation: 0,
            Created:    time.Now(),
        }

        candidates = append(candidates, candidate)
    }

    return candidates
}

// evaluateCandidates 评估候选模式
func (pg *PatternGenerator) evaluateCandidates(
    candidates []*PatternCandidate) []*PatternCandidate {
    
    for _, candidate := range candidates {
        // 计算基础分数
        baseScore := pg.calculateBaseScore(candidate.Pattern)
        
        // 评估复杂度
        complexityScore := pg.evaluateComplexity(candidate.Pattern)
        
        // 检查能量平衡
        energyScore := pg.checkEnergyBalance(candidate.Pattern)
        
        // 组合得分
        candidate.Score = pg.combineScores(baseScore, complexityScore, energyScore)
    }

    return candidates
}

// selectBestCandidates 选择最佳候选
func (pg *PatternGenerator) selectBestCandidates(
    candidates []*PatternCandidate) []*PatternCandidate {
    
    // 按分数排序
    sort.Slice(candidates, func(i, j int) bool {
        return candidates[i].Score > candidates[j].Score
    })

    // 选择前N个
    selected := candidates
    if len(selected) > maxSelected {
        selected = selected[:maxSelected]
    }

    return selected
}

// optimizePatterns 优化模式
func (pg *PatternGenerator) optimizePatterns(
    patterns []*PatternCandidate) []*PatternCandidate {
    
    optimized := make([]*PatternCandidate, 0)

    for _, pattern := range patterns {
        // 应用优化规则
        improved := pg.optimizePattern(pattern)
        
        // 检查优化效果
        if improved.Score > pattern.Score {
            optimized = append(optimized, improved)
        } else {
            optimized = append(optimized, pattern)
        }
    }

    return optimized
}

// updateMetrics 更新指标
func (pg *PatternGenerator) updateMetrics(patterns []*PatternCandidate) {
    metrics := pg.state.metrics
    
    // 更新总数
    metrics.TotalGenerated += len(patterns)
    
    // 计算成功率
    successCount := 0
    totalScore := 0.0
    
    for _, pattern := range patterns {
        if pattern.Score >= successThreshold {
            successCount++
        }
        totalScore += pattern.Score
    }
    
    metrics.SuccessRate = float64(successCount) / float64(len(patterns))
    metrics.AverageScore = totalScore / float64(len(patterns))
    
    // 记录演化点
    point := MetricPoint{
        Timestamp: time.Now(),
        Metrics: map[string]float64{
            "success_rate":   metrics.SuccessRate,
            "average_score":  metrics.AverageScore,
        },
    }
    
    metrics.Evolution = append(metrics.Evolution, point)
}

// 辅助函数

func (pg *PatternGenerator) validateTemplate(template *GenerationTemplate) error {
    if template.ID == "" {
        return model.WrapError(nil, model.ErrCodeValidation, "empty template ID")
    }
    
    // 验证结构
    if err := pg.validateStructure(template.Structure); err != nil {
        return err
    }
    
    // 验证约束
    if err := pg.validateConstraints(template.Constraints); err != nil {
        return err
    }
    
    return nil
}

func generateCandidateID() string {
    return fmt.Sprintf("cand_%d", time.Now().UnixNano())
}

const (
    maxCandidates = 100
    maxSelected = 10
    successThreshold = 0.7
)
