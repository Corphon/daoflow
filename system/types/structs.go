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
    
    // 系统基础信息
    ID        string
    Name      string
    Version   string
    StartTime time.Time
    
    // 系统状态
    state struct {
        Current   SystemState
        Previous  SystemState
        LastChange time.Time
    }
    
    // 元系统组件
    meta struct {
        Field struct {
            State     FieldState
            Strength  float64
            Vector    []float64
            Tensor    [][]float64
        }
        
        Quantum struct {
            Wave      []complex128
            Phase     float64
            Entangled bool
            Coherence float64
        }
        
        Emergence struct {
            Patterns  []EmergentPattern
            Active    []EmergentProperty
            Potential []PotentialEmergence
        }
    }
    
    // 演化系统组件
    evolution struct {
        Level     float64
        Direction Vector3D
        Speed     float64
        Path      []EvolutionPoint
        History   []StateTransition
    }
    
    // 资源管理
    resources struct {
        Pool      ResourcePool
        Queue     []ResourceReq
        Active    []ResourceAllocation
        History   []ResourceEvent
    }
    
    // 系统配置
    config SystemConfig
}

// FieldState 场状态
type FieldState struct {
    // 基本场属性
    Strength   float64            // 场强度
    Phase      float64            // 场相位
    Polarity   float64            // 场极性
    Dimension  int               // 场维度
    
    // 场向量和张量
    Vector     []float64         // 场向量
    Tensor     [][]float64       // 场张量
    Gradient   []float64         // 场梯度
    
    // 场特性
    Properties map[string]float64 // 场属性
    Coupling   [][]float64       // 耦合矩阵
    Symmetry   string            // 对称性
}

// QuantumState 量子态
type QuantumState struct {
    // 波函数
    Wave        []complex128     // 波函数
    Phase       float64         // 相位
    Amplitude   float64         // 振幅
    
    // 量子特性
    Entangled   bool           // 纠缠状态
    Coherence   float64        // 相干度
    Superposed  bool           // 叠加态
    
    // 测量结果
    Measurement struct {
        Value      float64
        Certainty  float64
        Timestamp  time.Time
    }
}

// EmergentPattern 涌现模式
type EmergentPattern struct {
    ID          string
    Type        string
    Properties  map[string]float64
    Components  []string
    Strength    float64
    Stability   float64
    Formation   time.Time
}

// Vector3D 三维向量
type Vector3D struct {
    X float64
    Y float64
    Z float64
}

// EvolutionPoint 演化点
type EvolutionPoint struct {
    State     SystemState
    Energy    float64
    Time      time.Time
    Position  Vector3D
}

// ResourcePool 资源池
type ResourcePool struct {
    // 计算资源
    CPU struct {
        Total     float64
        Used      float64
        Reserved  float64
    }
    
    // 内存资源
    Memory struct {
        Total     float64
        Used      float64
        Cached    float64
    }
    
    // 能量资源
    Energy struct {
        Current   float64
        Max       float64
        Min       float64
        Flow      float64
    }
}

// ResourceEvent 资源事件
type ResourceEvent struct {
    Type      string
    Resource  string
    Amount    float64
    Time      time.Time
    Status    string
}

// SystemConfig 系统配置
type SystemConfig struct {
    // 基础配置
    Base struct {
        Name         string
        Description  string
        Version     string
        LogLevel    string
    }
    
    // 性能配置
    Performance struct {
        MaxGoroutines  int
        BufferSize     int
        Timeout        time.Duration
    }
    
    // 资源限制
    Limits struct {
        MaxCPU     float64
        MaxMemory  float64
        MaxEnergy  float64
    }
    
    // 演化参数
    Evolution struct {
        InitialLevel    float64
        MinSpeed        float64
        MaxSpeed        float64
        AdaptRate      float64
    }
    
    // 场配置
    Field struct {
        InitStrength   float64
        MinStrength    float64
        MaxStrength    float64
        Dimension     int
    }
    
    // 量子配置
    Quantum struct {
        WaveSize      int
        CoherenceThreshold float64
        EntanglementLimit  int
    }
}

// StateTransition 状态转换记录
type StateTransition struct {
    From      SystemState
    To        SystemState
    Reason    string
    Energy    float64
    Time      time.Time
}

// MetricsData 指标数据
type MetricsData struct {
    // 系统指标
    System struct {
        State     SystemState
        Health    float64
        Load      float64
        Uptime    time.Duration
    }
    
    // 性能指标
    Performance struct {
        CPU       float64
        Memory    float64
        Latency   float64
        QPS       float64
    }
    
    // 资源指标
    Resources struct {
        Usage     float64
        Available float64
        Efficiency float64
    }
    
    // 演化指标
    Evolution struct {
        Progress  float64
        Speed     float64
        Quality   float64
    }
    
    // 时间戳
    Timestamp time.Time
}

// AnalysisResult 分析结果
type AnalysisResult struct {
    // 基本信息
    ID        string
    Type      string
    Time      time.Time
    
    // 分析数据
    Data struct {
        Values    map[string]float64
        Patterns  []string
        Trends    []float64
    }
    
    // 评估结果
    Evaluation struct {
        Score     float64
        Confidence float64
        Risk      float64
    }
    
    // 建议actions
    Recommendations []string
}
