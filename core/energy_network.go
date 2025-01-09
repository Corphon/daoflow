// core/energy_network.go

package core

import (
    "math"
    "sync"
)

// EnergyNetworkConstants 能量网络常数
const (
    MinFlowRate     = 0.0   // 最小流动率
    MaxFlowRate     = 1.0   // 最大流动率
    DefaultCapacity = 100.0 // 默认容量
)

// EnergyFlow 能量流动记录
type EnergyFlow struct {
    Source      string  // 源节点
    Target      string  // 目标节点
    Amount      float64 // 流动量
    Timestamp   int64   // 时间戳
}

// NetworkNode 网络节点
type NetworkNode struct {
    ID       string             // 节点标识
    Energy   float64           // 节点能量
    Flows    map[string]float64 // 与其他节点的流动关系
    Capacity float64           // 节点容量
}

// EnergyNetwork 能量网络
type EnergyNetwork struct {
    mu sync.RWMutex

    // 网络结构
    nodes    map[string]*NetworkNode // 节点集合
    flows    []EnergyFlow           // 流动记录
    capacity float64                // 网络总容量

    // 网络状态
    totalEnergy float64 // 总能量
    flowRate    float64 // 流动率
    balance     float64 // 平衡度
}

// NewEnergyNetwork 创建新的能量网络
func NewEnergyNetwork() *EnergyNetwork {
    return &EnergyNetwork{
        nodes:    make(map[string]*NetworkNode),
        flows:    make([]EnergyFlow, 0),
        capacity: DefaultCapacity,
    }
}

// Initialize 初始化网络
func (en *EnergyNetwork) Initialize() error {
    en.mu.Lock()
    defer en.mu.Unlock()

    // 清空节点
    en.nodes = make(map[string]*NetworkNode)
    
    // 清空流动记录
    en.flows = make([]EnergyFlow, 0)

    // 重置状态
    en.totalEnergy = 0
    en.flowRate = MinFlowRate
    en.balance = 1.0

    return nil
}

// AddNode 添加节点
func (en *EnergyNetwork) AddNode(id string, capacity float64) error {
    en.mu.Lock()
    defer en.mu.Unlock()

    if _, exists := en.nodes[id]; exists {
        return NewCoreErrorWithCode(ErrInvalid, "node already exists")
    }

    en.nodes[id] = &NetworkNode{
        ID:       id,
        Energy:   0,
        Flows:    make(map[string]float64),
        Capacity: capacity,
    }

    return nil
}

// UpdateFlow 更新能量流动
func (en *EnergyNetwork) UpdateFlow(from, to string, amount float64) error {
    en.mu.Lock()
    defer en.mu.Unlock()

    // 验证节点
    sourceNode, exists := en.nodes[from]
    if !exists {
        return NewCoreErrorWithCode(ErrInvalid, "source node not found")
    }

    targetNode, exists := en.nodes[to]
    if !exists {
        return NewCoreErrorWithCode(ErrInvalid, "target node not found")
    }

    // 验证能量流动
    if amount < 0 {
        return NewCoreErrorWithCode(ErrRange, "negative energy flow")
    }

    if amount > sourceNode.Energy {
        return NewCoreErrorWithCode(ErrRange, "insufficient energy in source node")
    }

    if targetNode.Energy+amount > targetNode.Capacity {
        return NewCoreErrorWithCode(ErrRange, "target node capacity exceeded")
    }

    // 执行能量转移
    sourceNode.Energy -= amount
    targetNode.Energy += amount

    // 更新流动关系
    sourceNode.Flows[to] = amount
    targetNode.Flows[from] = -amount

    // 记录流动
    en.flows = append(en.flows, EnergyFlow{
        Source:    from,
        Target:    to,
        Amount:    amount,
        Timestamp: GetCurrentTimestamp(),
    })

    // 更新网络状态
    en.updateNetworkState()

    return nil
}

// GetNodeEnergy 获取节点能量
func (en *EnergyNetwork) GetNodeEnergy(id string) (float64, error) {
    en.mu.RLock()
    defer en.mu.RUnlock()

    node, exists := en.nodes[id]
    if !exists {
        return 0, NewCoreErrorWithCode(ErrInvalid, "node not found")
    }

    return node.Energy, nil
}

// GetTotalEnergy 获取总能量
func (en *EnergyNetwork) GetTotalEnergy() float64 {
    en.mu.RLock()
    defer en.mu.RUnlock()
    return en.totalEnergy
}

// GetBalance 获取网络平衡度
func (en *EnergyNetwork) GetBalance() float64 {
    en.mu.RLock()
    defer en.mu.RUnlock()
    return en.balance
}

// updateNetworkState 更新网络状态
func (en *EnergyNetwork) updateNetworkState() {
    // 计算总能量
    totalEnergy := 0.0
    for _, node := range en.nodes {
        totalEnergy += node.Energy
    }
    en.totalEnergy = totalEnergy

    // 计算平均能量
    avgEnergy := totalEnergy / float64(len(en.nodes))

    // 计算能量方差
    variance := 0.0
    for _, node := range en.nodes {
        diff := node.Energy - avgEnergy
        variance += diff * diff
    }
    variance /= float64(len(en.nodes))

    // 计算平衡度 (0-1)
    en.balance = 1 / (1 + math.Sqrt(variance))
}

// GetCurrentTimestamp 获取当前时间戳
func GetCurrentTimestamp() int64 {
    return time.Now().UnixNano()
}
