// system/types/structs.go

package types

import (
    "sync"
    "time"
    
    "github.com/Corphon/daoflow/model"
)

// System 系统主结构
type System struct {
    mu sync.RWMutex
    
    // 基础信息
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Version   string    `json:"version"`
    StartTime time.Time `json:"start_time"`
    
    // 系统状态
    state struct {
        Current    SystemState    `json:"current"`
        Previous   SystemState    `json:"previous"`
        LastChange time.Time      `json:"last_change"`
        History    []StateTransition `json:"history"`
    }
    
    // 元系统组件
    meta struct {
        Field struct {
            State     FieldState  `json:"state"`
            Quantum   QuantumState `json:"quantum"`
            Coupling  [][]float64  `json:"coupling"`
        }
        
        Emergence struct {
            Patterns  []EmergentPattern   `json:"patterns"`
            Active    []EmergentProperty  `json:"active"`
            Potential []PotentialEmergence `json:"potential"`
        }
        
        Resonance struct {
            State     ResonanceState `json:"state"`
            Coherence float64       `json:"coherence"`
            Phase     float64       `json:"phase"`
        }
    }
    
    // 演化组件
    evolution struct {
        Current struct {
            Level     float64     `json:"level"`
            Direction Vector3D    `json:"direction"`
            Speed     float64     `json:"speed"`
            Energy    float64     `json:"energy"`
        }
        
        History struct {
            Path      []EvolutionPoint   `json:"path"`
            Changes   []StateTransition  `json:"changes"`
            Metrics   []EvolutionMetrics `json:"metrics"`
        }
    }
    
    // 资源管理
    resources ResourceManager
    
    // 监控组件
    monitor struct {
        Metrics    MetricsCollector `json:"metrics"`
        Alerts     []Alert          `json:"alerts"`
        Status     HealthStatus     `json:"status"`
    }
    
    // 配置
    config SystemConfig
}

// FieldState 场状态
type FieldState struct {
    // 基本属性
    Strength   float64            `json:"strength"`
    Phase      float64            `json:"phase"`
    Polarity   float64            `json:"polarity"`
    Dimension  int               `json:"dimension"`
    
    // 场向量和张量
    Vector     []float64          `json:"vector"`
    Tensor     [][]float64        `json:"tensor"`
    Gradient   []float64          `json:"gradient"`
    
    // 特性
    Properties map[string]float64  `json:"properties"`
    Symmetry   string             `json:"symmetry"`
}

// QuantumState 量子态
type QuantumState struct {
    Wave       []complex128  `json:"wave"`
    Phase      float64      `json:"phase"`
    Amplitude  float64      `json:"amplitude"`
    Entangled  bool         `json:"entangled"`
    Coherence  float64      `json:"coherence"`
    Superposed bool         `json:"superposed"`
    
    Measurement struct {
        Value     float64   `json:"value"`
        Certainty float64   `json:"certainty"`
        Time      time.Time `json:"time"`
    } `json:"measurement"`
}

// ResonanceState 共振状态
type ResonanceState struct {
    Active    bool      `json:"active"`
    Frequency float64   `json:"frequency"`
    Amplitude float64   `json:"amplitude"`
    Phase     float64   `json:"phase"`
    Coupling  float64   `json:"coupling"`
}

// EmergentPattern 涌现模式
type EmergentPattern struct {
    ID         string              `json:"id"`
    Type       string              `json:"type"`
    Properties map[string]float64  `json:"properties"`
    Components []string            `json:"components"`
    Strength   float64            `json:"strength"`
    Stability  float64            `json:"stability"`
    Formation  time.Time          `json:"formation"`
}

// EmergentProperty 涌现属性
type EmergentProperty struct {
    ID         string             `json:"id"`
    Pattern    EmergentPattern    `json:"pattern"`
    State      map[string]float64 `json:"state"`
    Energy     float64           `json:"energy"`
    Time       time.Time         `json:"time"`
}

// PotentialEmergence 潜在涌现
type PotentialEmergence struct {
    Pattern     EmergentPattern `json:"pattern"`
    Probability float64        `json:"probability"`
    TimeFrame   time.Duration  `json:"time_frame"`
    Impact      float64        `json:"impact"`
}

// Vector3D 三维向量
type Vector3D struct {
    X float64 `json:"x"`
    Y float64 `json:"y"`
    Z float64 `json:"z"`
}

// EvolutionPoint 演化点
type EvolutionPoint struct {
    State     SystemState `json:"state"`
    Energy    float64    `json:"energy"`
    Time      time.Time  `json:"time"`
    Position  Vector3D   `json:"position"`
}

// EvolutionMetrics 演化指标
type EvolutionMetrics struct {
    Level     float64   `json:"level"`
    Speed     float64   `json:"speed"`
    Direction Vector3D  `json:"direction"`
    Energy    float64   `json:"energy"`
    Time      time.Time `json:"time"`
}

// ResourceManager 资源管理器
type ResourceManager struct {
    Pool      ResourcePool         `json:"pool"`
    Queue     []ResourceReq        `json:"queue"`
    Active    []ResourceAllocation `json:"active"`
    History   []ResourceEvent      `json:"history"`
    Stats     ResourceStats        `json:"stats"`
}

// ResourcePool 资源池
type ResourcePool struct {
    CPU struct {
        Total     float64 `json:"total"`
        Used      float64 `json:"used"`
        Reserved  float64 `json:"reserved"`
    } `json:"cpu"`
    
    Memory struct {
        Total     float64 `json:"total"`
        Used      float64 `json:"used"`
        Cached    float64 `json:"cached"`
    } `json:"memory"`
    
    Energy struct {
        Current   float64 `json:"current"`
        Max       float64 `json:"max"`
        Min       float64 `json:"min"`
        Flow      float64 `json:"flow"`
    } `json:"energy"`
}

// ResourceStats 资源统计
type ResourceStats struct {
    Utilization  float64   `json:"utilization"`
    Efficiency   float64   `json:"efficiency"`
    Balance      float64   `json:"balance"`
    LastUpdate   time.Time `json:"last_update"`
}

// Alert 告警信息
type Alert struct {
    ID        string        `json:"id"`
    Type      string        `json:"type"`
    Level     IssueSeverity `json:"level"`
    Message   string        `json:"message"`
    Source    string        `json:"source"`
    Time      time.Time     `json:"time"`
    Status    string        `json:"status"`
}

// MetricsCollector 指标收集器
type MetricsCollector struct {
    Current  MetricsData    `json:"current"`
    History  []MetricsData  `json:"history"`
    Config   MetricsConfig  `json:"config"`
}

// MetricsConfig 指标配置
type MetricsConfig struct {
    Interval    time.Duration          `json:"interval"`
    Retention   time.Duration          `json:"retention"`
    Thresholds  map[string]float64     `json:"thresholds"`
    Filters     []string               `json:"filters"`
}

// StateTransition 状态转换
type StateTransition struct {
    From      SystemState `json:"from"`
    To        SystemState `json:"to"`
    Reason    string     `json:"reason"`
    Energy    float64    `json:"energy"`
    Time      time.Time  `json:"time"`
}

// SystemEvent 系统事件
type SystemEvent struct {
    ID        string      `json:"id"`
    Type      EventType   `json:"type"`
    Source    string      `json:"source"`
    Target    string      `json:"target"`
    Data      interface{} `json:"data"`
    Time      time.Time   `json:"time"`
    Status    string      `json:"status"`
}
