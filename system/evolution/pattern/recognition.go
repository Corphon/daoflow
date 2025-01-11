//system/evolution/pattern/recognition.go

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

// PatternRecognizer 模式识别器
type PatternRecognizer struct {
    mu sync.RWMutex

    // 基础配置
    config struct {
        minConfidence   float64         // 最小置信度
        learningRate    float64         // 学习率
        memoryDepth     int             // 记忆深度
        adaptiveRate    bool            // 是否使用自适应学习率
    }

    // 识别状态
    state struct {
        patterns     map[string]*RecognizedPattern  // 已识别模式
        memories     []PatternMemory               // 模式记忆
        statistics   PatternStatistics             // 统计信息
    }

    // 依赖项
    detector   *emergence.PatternDetector
    matcher    *resonance.PatternMatcher
    amplifier  *resonance.ResonanceAmplifier
}

// RecognizedPattern 已识别模式
type RecognizedPattern struct {
    ID           string                 // 模式ID
    Type         string                 // 模式类型
    Signature    PatternSignature       // 模式特征
    Confidence   float64                // 置信度
    Stability    float64                // 稳定性
    FirstSeen    time.Time             // 首次发现时间
    LastSeen     time.Time             // 最后发现时间
    Occurrences  int                    // 出现次数
    Evolution    []PatternState         // 演化历史
}

// PatternSignature 模式特征
type PatternSignature struct {
    Components  []SignatureComponent    // 组成成分
    Structure   map[string]interface{}  // 结构特征
    Dynamics    map[string]float64      // 动态特征
    Context     map[string]string       // 上下文信息
}

// SignatureComponent 特征组件
type SignatureComponent struct {
    Type       string                  // 组件类型
    Properties map[string]float64      // 组件属性
    Weight     float64                 // 权重
    Role       string                  // 角色
}

// PatternMemory 模式记忆
type PatternMemory struct {
    Timestamp    time.Time
    Pattern      *RecognizedPattern
    Context      map[string]interface{}
    Associations []string
}

// PatternStatistics 模式统计
type PatternStatistics struct {
    TotalPatterns    int
    ActivePatterns   int
    Recognition      map[string]float64  // 识别率统计
    Accuracy        map[string]float64  // 准确率统计
    Evolution       []StatPoint         // 演化趋势
}

// StatPoint 统计点
type StatPoint struct {
    Timestamp time.Time
    Metrics   map[string]float64
}

// NewPatternRecognizer 创建新的模式识别器
func NewPatternRecognizer(
    detector *emergence.PatternDetector,
    matcher *resonance.PatternMatcher,
    amplifier *resonance.ResonanceAmplifier) *PatternRecognizer {
    
    pr := &PatternRecognizer{
        detector:  detector,
        matcher:   matcher,
        amplifier: amplifier,
    }

    // 初始化配置
    pr.config.minConfidence = 0.75
    pr.config.learningRate = 0.1
    pr.config.memoryDepth = 100
    pr.config.adaptiveRate = true

    // 初始化状态
    pr.state.patterns = make(map[string]*RecognizedPattern)
    pr.state.memories = make([]PatternMemory, 0)
    pr.state.statistics = PatternStatistics{
        Recognition: make(map[string]float64),
        Accuracy:   make(map[string]float64),
        Evolution:  make([]StatPoint, 0),
    }

    return pr
}

// Recognize 执行模式识别
func (pr *PatternRecognizer) Recognize() error {
    pr.mu.Lock()
    defer pr.mu.Unlock()

    // 获取当前模式
    patterns, err := pr.detector.Detect()
    if err != nil {
        return err
    }

    // 识别新模式
    newPatterns := pr.recognizeNewPatterns(patterns)

    // 更新现有模式
    pr.updateExistingPatterns(patterns)

    // 构建模式记忆
    pr.buildPatternMemory(newPatterns)

    // 更新统计信息
    pr.updateStatistics()

    return nil
}

// recognizeNewPatterns 识别新模式
func (pr *PatternRecognizer) recognizeNewPatterns(
    patterns []emergence.EmergentPattern) []*RecognizedPattern {
    
    newPatterns := make([]*RecognizedPattern, 0)

    for _, pattern := range patterns {
        // 检查是否是新模式
        if pr.isKnownPattern(pattern) {
            continue
        }

        // 提取模式特征
        signature := pr.extractSignature(pattern)

        // 评估模式
        confidence := pr.evaluatePattern(pattern, signature)
        if confidence < pr.config.minConfidence {
            continue
        }

        // 创建新的识别模式
        recognized := &RecognizedPattern{
            ID:          generatePatternID(),
            Type:        determinePatternType(pattern),
            Signature:   signature,
            Confidence:  confidence,
            Stability:   calculateInitialStability(pattern),
            FirstSeen:   time.Now(),
            LastSeen:    time.Now(),
            Occurrences: 1,
            Evolution:   make([]PatternState, 0),
        }

        // 添加到已识别模式
        pr.state.patterns[recognized.ID] = recognized
        newPatterns = append(newPatterns, recognized)
    }

    return newPatterns
}

// updateExistingPatterns 更新现有模式
func (pr *PatternRecognizer) updateExistingPatterns(
    patterns []emergence.EmergentPattern) {
    
    currentTime := time.Now()

    for id, recognized := range pr.state.patterns {
        // 查找匹配的当前模式
        matched := false
        for _, pattern := range patterns {
            if pr.isPatternMatch(recognized, pattern) {
                // 更新模式状态
                pr.updatePatternState(recognized, pattern)
                matched = true
                break
            }
        }

        // 处理未匹配的模式
        if !matched {
            // 检查是否应该保留模式
            if pr.shouldRetainPattern(recognized) {
                // 降低置信度
                recognized.Confidence *= (1 - pr.config.learningRate)
            } else {
                // 移除模式
                delete(pr.state.patterns, id)
            }
        }
    }
}

// buildPatternMemory 构建模式记忆
func (pr *PatternRecognizer) buildPatternMemory(newPatterns []*RecognizedPattern) {
    memory := PatternMemory{
        Timestamp:    time.Now(),
        Pattern:      nil,
        Context:      make(map[string]interface{}),
        Associations: make([]string, 0),
    }

    // 记录新模式
    for _, pattern := range newPatterns {
        memory.Pattern = pattern
        memory.Context = pr.extractContext(pattern)
        memory.Associations = pr.findAssociations(pattern)

        pr.state.memories = append(pr.state.memories, memory)
    }

    // 限制记忆深度
    if len(pr.state.memories) > pr.config.memoryDepth {
        pr.state.memories = pr.state.memories[1:]
    }
}

// 辅助函数

func (pr *PatternRecognizer) isKnownPattern(
    pattern emergence.EmergentPattern) bool {
    
    for _, recognized := range pr.state.patterns {
        if pr.isPatternMatch(recognized, pattern) {
            return true
        }
    }
    return false
}

func (pr *PatternRecognizer) extractSignature(
    pattern emergence.EmergentPattern) PatternSignature {
    
    signature := PatternSignature{
        Components: make([]SignatureComponent, 0),
        Structure:  make(map[string]interface{}),
        Dynamics:   make(map[string]float64),
        Context:    make(map[string]string),
    }

    // 提取组件特征
    for _, comp := range pattern.Components {
        component := SignatureComponent{
            Type:       comp.Type,
            Properties: make(map[string]float64),
            Weight:     comp.Weight,
            Role:       comp.Role,
        }
        
        // 复制属性
        for k, v := range comp.Properties {
            component.Properties[k] = v
        }
        
        signature.Components = append(signature.Components, component)
    }

    // 提取结构特征
    signature.Structure = extractStructuralFeatures(pattern)

    // 提取动态特征
    signature.Dynamics = extractDynamicFeatures(pattern)

    return signature
}

func (pr *PatternRecognizer) evaluatePattern(
    pattern emergence.EmergentPattern,
    signature PatternSignature) float64 {
    
    // 基础置信度
    baseConfidence := pattern.Strength

    // 结构完整性评分
    structureScore := evaluateStructure(signature.Structure)

    // 动态稳定性评分
    stabilityScore := evaluateStability(signature.Dynamics)

    // 组合评分
    confidence := (baseConfidence + structureScore + stabilityScore) / 3.0

    return confidence
}

func generatePatternID() string {
    return fmt.Sprintf("pat_%d", time.Now().UnixNano())
}

const (
    maxHistoryLength = 1000
)

type PatternState struct {
    Timestamp  time.Time
    Confidence float64
    Stability  float64
    Changes    map[string]float64
}
