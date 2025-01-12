// system/common/patterns.go

package common

import (
    "time"
    "github.com/Corphon/daoflow/system/types"
)

// SharedPattern 共享模式接口
type SharedPattern interface {
    GetID() string
    GetType() string
    GetStrength() float64
    GetStability() float64
    GetTimestamp() time.Time
}

// BasePattern 基础模式结构
type BasePattern struct {
    ID        string
    Type      string
    Strength  float64
    Stability float64
    Created   time.Time
}

// 实现 SharedPattern 接口
func (bp *BasePattern) GetID() string         { return bp.ID }
func (bp *BasePattern) GetType() string       { return bp.Type }
func (bp *BasePattern) GetStrength() float64  { return bp.Strength }
func (bp *BasePattern) GetStability() float64 { return bp.Stability }
func (bp *BasePattern) GetTimestamp() time.Time { return bp.Created }

// PatternAnalyzer 模式分析器接口
type PatternAnalyzer interface {
    AnalyzePattern(SharedPattern) (float64, error)
    ComparePatterns(p1, p2 SharedPattern) (float64, error)
}

// PatternMatcher 模式匹配器接口
type PatternMatcher interface {
    MatchPatterns(patterns []SharedPattern) ([]PatternMatch, error)
}

// PatternMatch 模式匹配结果
type PatternMatch struct {
    Source      SharedPattern
    Target      SharedPattern
    Similarity  float64
    Confidence  float64
}

// PatternEventEmitter 模式事件发射器接口
type PatternEventEmitter interface {
    EmitPatternEvent(event PatternEvent)
    AddPatternListener(listener PatternEventListener)
    RemovePatternListener(listener PatternEventListener)
}

// PatternEventListener 模式事件监听器接口
type PatternEventListener interface {
    OnPatternEvent(event PatternEvent)
}

// PatternEvent 模式事件
type PatternEvent struct {
    Type      string
    Pattern   SharedPattern
    Timestamp time.Time
    Data      map[string]interface{}
}

// PatternProcessor 模式处理器接口
type PatternProcessor interface {
    ProcessPattern(pattern SharedPattern) error
    GetProcessingResult() ProcessingResult
}

// ProcessingResult 处理结果
type ProcessingResult struct {
    Success    bool
    Score      float64
    Changes    []PatternChange
    Timestamp  time.Time
}

// PatternChange 模式变化
type PatternChange struct {
    Field     string
    OldValue  interface{}
    NewValue  interface{}
    Delta     float64
}
