//system/control/flow/balancer.go

package flow

import (
    "sync"
    "time"

    "github.com/Corphon/daoflow/meta/field"
    "github.com/Corphon/daoflow/model"
    "github.com/Corphon/daoflow/system/types"
)

// LoadBalancer 负载均衡器
type LoadBalancer struct {
    mu sync.RWMutex

    // 基础配置
    config struct {
        balanceInterval time.Duration  // 均衡间隔
        maxLoad        float64        // 最大负载
        targetLoad     float64        // 目标负载
        smoothingFactor float64       // 平滑因子
    }

    // 均衡状态
    state struct {
        nodes       map[string]*Node         // 节点列表
        workloads   map[string]*Workload     // 工作负载
        allocations map[string]*Allocation   // 资源分配
        metrics     BalancerMetrics         // 均衡指标
    }

    // 依赖项
    scheduler    *Scheduler
    backpressure *BackpressureManager
}

// Node 节点信息
type Node struct {
    ID           string                // 节点ID
    Type         string                // 节点类型
    Capacity     ResourceCapacity      // 资源容量
    Status       string                // 节点状态
    Health       float64               // 健康度
    LastHeartbeat time.Time           // 最后心跳
}

// ResourceCapacity 资源容量
type ResourceCapacity struct {
    CPU          float64              // CPU容量
    Memory       int64                // 内存容量
    Storage      int64                // 存储容量
    Network      float64              // 网络带宽
}

// Workload 工作负载
type Workload struct {
    ID           string                // 负载ID
    Type         string                // 负载类型
    Resources    ResourceRequirement   // 资源需求
    Priority     int                   // 优先级
    Distribution []WorkloadDistribution // 负载分布
}

// ResourceRequirement 资源需求
type ResourceRequirement struct {
    MinCPU       float64              // 最小CPU
    MinMemory    int64                // 最小内存
    MinStorage   int64                // 最小存储
    MinNetwork   float64              // 最小带宽
}

// WorkloadDistribution 负载分布
type WorkloadDistribution struct {
    NodeID       string                // 节点ID
    Percentage   float64               // 分配百分比
    Performance  DistributionMetrics   // 分布指标
}

// Allocation 资源分配
type Allocation struct {
    ID           string                // 分配ID
    WorkloadID   string                // 负载ID
    NodeID       string                // 节点ID
    Resources    ResourceAllocation    // 资源分配
    Status       string                // 分配状态
    StartTime    time.Time            // 开始时间
}

// ResourceAllocation 资源分配
type ResourceAllocation struct {
    CPU          float64              // CPU分配
    Memory       int64                // 内存分配
    Storage      int64                // 存储分配
    Network      float64              // 带宽分配
}

// DistributionMetrics 分布指标
type DistributionMetrics struct {
    ResponseTime time.Duration        // 响应时间
    Throughput   float64              // 吞吐量
    ErrorRate    float64              // 错误率
    Utilization  float64              // 资源利用率
}

// BalancerMetrics 均衡器指标
type BalancerMetrics struct {
    TotalNodes    int                 // 总节点数
    ActiveNodes   int                 // 活跃节点数
    LoadVariance  float64             // 负载方差
    Imbalance     float64             // 不平衡度
    History       []MetricPoint       // 历史指标
}

// NewLoadBalancer 创建新的负载均衡器
func NewLoadBalancer(
    scheduler *Scheduler,
    backpressure *BackpressureManager) *LoadBalancer {
    
    lb := &LoadBalancer{
        scheduler:    scheduler,
        backpressure: backpressure,
    }

    // 初始化配置
    lb.config.balanceInterval = 5 * time.Second
    lb.config.maxLoad = 0.9
    lb.config.targetLoad = 0.7
    lb.config.smoothingFactor = 0.3

    // 初始化状态
    lb.state.nodes = make(map[string]*Node)
    lb.state.workloads = make(map[string]*Workload)
    lb.state.allocations = make(map[string]*Allocation)
    lb.state.metrics = BalancerMetrics{
        History: make([]MetricPoint, 0),
    }

    return lb
}

// Balance 执行负载均衡
func (lb *LoadBalancer) Balance() error {
    lb.mu.Lock()
    defer lb.mu.Unlock()

    // 更新节点状态
    if err := lb.updateNodeStatus(); err != nil {
        return err
    }

    // 分析工作负载
    if err := lb.analyzeWorkloads(); err != nil {
        return err
    }

    // 重新分配资源
    if err := lb.rebalanceResources(); err != nil {
        return err
    }

    // 应用新的分配
    if err := lb.applyAllocations(); err != nil {
        return err
    }

    // 更新指标
    lb.updateMetrics()

    return nil
}

// RegisterNode 注册新节点
func (lb *LoadBalancer) RegisterNode(node *Node) error {
    if node == nil {
        return model.WrapError(nil, model.ErrCodeValidation, "nil node")
    }

    lb.mu.Lock()
    defer lb.mu.Unlock()

    // 验证节点
    if err := lb.validateNode(node); err != nil {
        return err
    }

    // 存储节点
    lb.state.nodes[node.ID] = node

    return nil
}

// analyzeWorkloads 分析工作负载
func (lb *LoadBalancer) analyzeWorkloads() error {
    for _, workload := range lb.state.workloads {
        // 检查资源需求
        if err := lb.checkResourceRequirements(workload); err != nil {
            continue
        }

        // 评估当前分布
        lb.evaluateDistribution(workload)

        // 计算最优分布
        if err := lb.calculateOptimalDistribution(workload); err != nil {
            continue
        }
    }

    return nil
}

// rebalanceResources 重新分配资源
func (lb *LoadBalancer) rebalanceResources() error {
    // 根据优先级排序工作负载
    workloads := lb.sortWorkloadsByPriority()

    // 为每个工作负载分配资源
    for _, workload := range workloads {
        if err := lb.allocateResources(workload); err != nil {
            continue
        }
    }

    return nil
}

// 辅助函数

func (lb *LoadBalancer) validateNode(node *Node) error {
    if node.ID == "" {
        return model.WrapError(nil, model.ErrCodeValidation, "empty node ID")
    }

    if node.Capacity.CPU <= 0 {
        return model.WrapError(nil, model.ErrCodeValidation, "invalid CPU capacity")
    }

    if node.Capacity.Memory <= 0 {
        return model.WrapError(nil, model.ErrCodeValidation, "invalid memory capacity")
    }

    return nil
}

func (lb *LoadBalancer) updateMetrics() {
    point := MetricPoint{
        Timestamp: time.Now(),
        Values: map[string]float64{
            "load_variance": lb.calculateLoadVariance(),
            "imbalance":    lb.calculateImbalance(),
        },
    }

    lb.state.metrics.History = append(lb.state.metrics.History, point)

    // 限制历史记录数量
    if len(lb.state.metrics.History) > maxMetricsHistory {
        lb.state.metrics.History = lb.state.metrics.History[1:]
    }
}

const (
    maxMetricsHistory = 1000
)
