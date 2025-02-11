//system/control/flow/balancer.go

package flow

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/Corphon/daoflow/model"
)

const (
	maxMetricsHistory = 1000
)

// LoadBalancer 负载均衡器
type LoadBalancer struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		balanceInterval time.Duration // 均衡间隔
		maxLoad         float64       // 最大负载
		targetLoad      float64       // 目标负载
		smoothingFactor float64       // 平滑因子
	}

	// 均衡状态
	state struct {
		nodes       map[string]*Node       // 节点列表
		workloads   map[string]*Workload   // 工作负载
		allocations map[string]*Allocation // 资源分配
		metrics     BalancerMetrics        // 均衡指标
	}

	// 依赖项
	scheduler    *Scheduler
	backpressure *BackpressureManager
}

// Node 节点信息
type Node struct {
	ID            string           // 节点ID
	Type          string           // 节点类型
	Capacity      ResourceCapacity // 资源容量
	Status        string           // 节点状态
	Health        float64          // 健康度
	LastHeartbeat time.Time        // 最后心跳
}

// ResourceCapacity 资源容量
type ResourceCapacity struct {
	CPU     float64 // CPU容量
	Memory  int64   // 内存容量
	Storage int64   // 存储容量
	Network float64 // 网络带宽
}

// Workload 工作负载
type Workload struct {
	ID           string                 // 负载ID
	Type         string                 // 负载类型
	Resources    ResourceRequirement    // 资源需求
	Priority     int                    // 优先级
	Distribution []WorkloadDistribution // 负载分布
}

// ResourceRequirement 资源需求
type ResourceRequirement struct {
	MinCPU     float64 // 最小CPU
	MinMemory  int64   // 最小内存
	MinStorage int64   // 最小存储
	MinNetwork float64 // 最小带宽
}

// WorkloadDistribution 负载分布
type WorkloadDistribution struct {
	NodeID      string              // 节点ID
	Percentage  float64             // 分配百分比
	Performance DistributionMetrics // 分布指标
}

// Allocation 资源分配
type Allocation struct {
	ID         string             // 分配ID
	WorkloadID string             // 负载ID
	NodeID     string             // 节点ID
	Resources  ResourceAllocation // 资源分配
	Status     string             // 分配状态
	StartTime  time.Time          // 开始时间
}

// ResourceAllocation 资源分配
type ResourceAllocation struct {
	CPU     float64 // CPU分配
	Memory  int64   // 内存分配
	Storage int64   // 存储分配
	Network float64 // 带宽分配
}

// DistributionMetrics 分布指标
type DistributionMetrics struct {
	ResponseTime time.Duration // 响应时间
	Throughput   float64       // 吞吐量
	ErrorRate    float64       // 错误率
	Utilization  float64       // 资源利用率
}

// BalancerMetrics 均衡器指标
type BalancerMetrics struct {
	TotalNodes   int           // 总节点数
	ActiveNodes  int           // 活跃节点数
	LoadVariance float64       // 负载方差
	Imbalance    float64       // 不平衡度
	History      []MetricPoint // 历史指标
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

// updateNodeStatus 更新节点状态
func (lb *LoadBalancer) updateNodeStatus() error {
	// 更新节点状态和健康度
	for _, node := range lb.state.nodes {
		// 检查心跳超时
		if time.Since(node.LastHeartbeat) > lb.config.balanceInterval*2 {
			node.Status = "offline"
			node.Health = 0
			continue
		}

		// 计算当前负载
		currentLoad := lb.calculateNodeLoad(node.ID)

		// 更新节点状态
		if currentLoad >= lb.config.maxLoad {
			node.Status = "overloaded"
			node.Health = math.Max(0, 1-(currentLoad-lb.config.maxLoad)/0.1)
		} else if currentLoad >= lb.config.targetLoad {
			node.Status = "busy"
			node.Health = 1 - (currentLoad-lb.config.targetLoad)/(lb.config.maxLoad-lb.config.targetLoad)
		} else {
			node.Status = "healthy"
			node.Health = 1
		}
	}

	return nil
}

// applyAllocations 应用资源分配
func (lb *LoadBalancer) applyAllocations() error {
	for _, allocation := range lb.state.allocations {
		// 跳过已完成的分配
		if allocation.Status == "completed" {
			continue
		}

		// 获取节点和工作负载
		node, exists := lb.state.nodes[allocation.NodeID]
		if !exists {
			continue
		}

		workload, exists := lb.state.workloads[allocation.WorkloadID]
		if !exists {
			continue
		}

		// 验证分配是否可行
		if !lb.validateAllocation(allocation, node) {
			allocation.Status = "failed"
			continue
		}

		// 应用分配
		allocation.Status = "active"
		allocation.StartTime = time.Now()

		// 更新工作负载分布
		lb.updateWorkloadDistribution(workload, allocation)
	}

	return nil
}

// validateAllocation 验证资源分配是否可行
func (lb *LoadBalancer) validateAllocation(allocation *Allocation, node *Node) bool {
	// 检查CPU分配
	if allocation.Resources.CPU > node.Capacity.CPU {
		return false
	}

	// 检查内存分配
	if allocation.Resources.Memory > node.Capacity.Memory {
		return false
	}

	// 检查存储分配
	if allocation.Resources.Storage > node.Capacity.Storage {
		return false
	}

	// 检查网络分配
	if allocation.Resources.Network > node.Capacity.Network {
		return false
	}

	return true
}

// calculateNodeLoad 计算节点当前负载
func (lb *LoadBalancer) calculateNodeLoad(nodeID string) float64 {
	var totalLoad float64

	for _, allocation := range lb.state.allocations {
		if allocation.NodeID == nodeID && allocation.Status == "active" {
			// 计算CPU负载贡献
			cpuLoad := allocation.Resources.CPU / lb.state.nodes[nodeID].Capacity.CPU
			// 计算内存负载贡献
			memLoad := float64(allocation.Resources.Memory) / float64(lb.state.nodes[nodeID].Capacity.Memory)
			// 取最大值作为该分配的负载贡献
			totalLoad += math.Max(cpuLoad, memLoad)
		}
	}

	return totalLoad
}

// updateWorkloadDistribution 更新工作负载分布
func (lb *LoadBalancer) updateWorkloadDistribution(workload *Workload, allocation *Allocation) {
	// 查找现有分布
	var distribution *WorkloadDistribution
	for i := range workload.Distribution {
		if workload.Distribution[i].NodeID == allocation.NodeID {
			distribution = &workload.Distribution[i]
			break
		}
	}

	// 如果没有现有分布，创建新的
	if distribution == nil {
		workload.Distribution = append(workload.Distribution, WorkloadDistribution{
			NodeID:     allocation.NodeID,
			Percentage: 0,
			Performance: DistributionMetrics{
				ResponseTime: 0,
				Throughput:   0,
				ErrorRate:    0,
				Utilization:  0,
			},
		})
		distribution = &workload.Distribution[len(workload.Distribution)-1]
	}

	// 更新分布指标
	distribution.Percentage = lb.calculateDistributionPercentage(workload, allocation)
}

// calculateDistributionPercentage 计算负载分布百分比
func (lb *LoadBalancer) calculateDistributionPercentage(workload *Workload, allocation *Allocation) float64 {
	// 计算当前节点分配的资源比例
	var totalResources, nodeResources float64

	// 汇总所有分配的资源
	for _, alloc := range lb.state.allocations {
		if alloc.WorkloadID == workload.ID && alloc.Status == "active" {
			// 使用CPU和内存的最大相对值作为资源度量
			cpuRatio := alloc.Resources.CPU / workload.Resources.MinCPU
			memRatio := float64(alloc.Resources.Memory) / float64(workload.Resources.MinMemory)
			resources := math.Max(cpuRatio, memRatio)

			totalResources += resources
			if alloc.ID == allocation.ID {
				nodeResources = resources
			}
		}
	}

	// 如果没有资源分配，返回0
	if totalResources == 0 {
		return 0
	}

	// 计算百分比并使用平滑因子
	currentPercentage := nodeResources / totalResources

	// 查找当前分布的旧百分比
	var oldPercentage float64
	for _, dist := range workload.Distribution {
		if dist.NodeID == allocation.NodeID {
			oldPercentage = dist.Percentage
			break
		}
	}

	// 应用平滑因子计算最终百分比
	smoothedPercentage := oldPercentage*(1-lb.config.smoothingFactor) +
		currentPercentage*lb.config.smoothingFactor

	// 确保结果在[0,1]范围内
	return math.Max(0, math.Min(1, smoothedPercentage))
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

// checkResourceRequirements 检查资源需求
func (lb *LoadBalancer) checkResourceRequirements(workload *Workload) error {
	// 验证最小资源需求
	if workload.Resources.MinCPU <= 0 {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid minimum CPU requirement")
	}
	if workload.Resources.MinMemory <= 0 {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid minimum memory requirement")
	}
	if workload.Resources.MinStorage <= 0 {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid minimum storage requirement")
	}
	if workload.Resources.MinNetwork <= 0 {
		return model.WrapError(nil, model.ErrCodeValidation, "invalid minimum network requirement")
	}

	// 检查是否有足够的资源满足需求
	var totalCPU float64
	var totalMemory, totalStorage int64
	var totalNetwork float64

	for _, node := range lb.state.nodes {
		if node.Status == "healthy" {
			totalCPU += node.Capacity.CPU
			totalMemory += node.Capacity.Memory
			totalStorage += node.Capacity.Storage
			totalNetwork += node.Capacity.Network
		}
	}

	if totalCPU < workload.Resources.MinCPU {
		return model.WrapError(nil, model.ErrCodeResource, "insufficient CPU resources")
	}
	if totalMemory < workload.Resources.MinMemory {
		return model.WrapError(nil, model.ErrCodeResource, "insufficient memory resources")
	}

	return nil
}

// evaluateDistribution 评估当前分布
func (lb *LoadBalancer) evaluateDistribution(workload *Workload) {
	for i := range workload.Distribution {
		dist := &workload.Distribution[i]
		node := lb.state.nodes[dist.NodeID]
		if node == nil {
			continue
		}

		// 更新性能指标
		currentLoad := lb.calculateNodeLoad(node.ID)
		dist.Performance.Utilization = currentLoad

		// 计算响应时间（基于负载的简单模型）
		baseResponseTime := time.Millisecond * 100
		loadFactor := 1 + currentLoad
		dist.Performance.ResponseTime = time.Duration(float64(baseResponseTime) * loadFactor)

		// 计算吞吐量（反比于负载）
		dist.Performance.Throughput = (1 - currentLoad) * 100

		// 计算错误率（与负载正相关）
		dist.Performance.ErrorRate = math.Max(0, (currentLoad-0.8)*100)
	}
}

// calculateOptimalDistribution 计算最优分布
func (lb *LoadBalancer) calculateOptimalDistribution(workload *Workload) error {
	type nodeScore struct {
		nodeID string
		score  float64
	}

	// 计算每个健康节点的得分
	var scores []nodeScore
	for id, node := range lb.state.nodes {
		if node.Status != "healthy" {
			continue
		}

		// 基础得分从健康度开始
		score := node.Health

		// 考虑当前负载
		currentLoad := lb.calculateNodeLoad(id)
		loadScore := 1 - currentLoad

		// 考虑资源匹配度
		cpuMatch := node.Capacity.CPU / workload.Resources.MinCPU
		memMatch := float64(node.Capacity.Memory) / float64(workload.Resources.MinMemory)
		resourceScore := math.Min(cpuMatch, memMatch)

		// 综合评分
		finalScore := (score + loadScore + resourceScore) / 3

		scores = append(scores, nodeScore{id, finalScore})
	}

	if len(scores) == 0 {
		return model.WrapError(nil, model.ErrCodeResource, "no healthy nodes available")
	}

	// 根据得分排序
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// 更新工作负载的理想分布
	idealDistribution := make([]WorkloadDistribution, len(scores))
	totalScore := 0.0
	for _, s := range scores {
		totalScore += s.score
	}

	for i, s := range scores {
		percentage := s.score / totalScore
		idealDistribution[i] = WorkloadDistribution{
			NodeID:     s.nodeID,
			Percentage: percentage,
		}
	}

	// 存储计算结果供后续重平衡使用
	workload.Distribution = idealDistribution

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

// sortWorkloadsByPriority 根据工作负载优先级排序
func (lb *LoadBalancer) sortWorkloadsByPriority() []*Workload {
	workloads := make([]*Workload, 0, len(lb.state.workloads))
	for _, w := range lb.state.workloads {
		workloads = append(workloads, w)
	}

	// 按优先级降序排序
	sort.Slice(workloads, func(i, j int) bool {
		return workloads[i].Priority > workloads[j].Priority
	})

	return workloads
}

// allocateResources 为工作负载分配资源
func (lb *LoadBalancer) allocateResources(workload *Workload) error {
	// 获取理想分布
	if len(workload.Distribution) == 0 {
		return model.WrapError(nil, model.ErrCodeResource, "no distribution plan")
	}

	// 为每个分布创建或更新分配
	for _, dist := range workload.Distribution {
		node := lb.state.nodes[dist.NodeID]
		if node == nil || node.Status != "healthy" {
			continue
		}

		// 计算应分配的资源量
		resourceAlloc := ResourceAllocation{
			CPU:     workload.Resources.MinCPU * dist.Percentage,
			Memory:  int64(float64(workload.Resources.MinMemory) * dist.Percentage),
			Storage: int64(float64(workload.Resources.MinStorage) * dist.Percentage),
			Network: workload.Resources.MinNetwork * dist.Percentage,
		}

		// 创建或更新分配
		allocation := &Allocation{
			ID:         fmt.Sprintf("%s-%s", workload.ID, node.ID),
			WorkloadID: workload.ID,
			NodeID:     node.ID,
			Resources:  resourceAlloc,
			Status:     "pending",
		}

		// 验证分配是否可行
		if !lb.validateAllocation(allocation, node) {
			continue
		}

		// 存储分配
		lb.state.allocations[allocation.ID] = allocation
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
			"imbalance":     lb.calculateImbalance(),
		},
	}

	lb.state.metrics.History = append(lb.state.metrics.History, point)

	// 限制历史记录数量
	if len(lb.state.metrics.History) > maxMetricsHistory {
		lb.state.metrics.History = lb.state.metrics.History[1:]
	}
}

// calculateLoadVariance 计算负载方差
func (lb *LoadBalancer) calculateLoadVariance() float64 {
	var totalLoad, count float64
	loads := make([]float64, 0)

	// 收集所有活跃节点的负载
	for id, node := range lb.state.nodes {
		if node.Status != "offline" {
			load := lb.calculateNodeLoad(id)
			loads = append(loads, load)
			totalLoad += load
			count++
		}
	}

	if count == 0 {
		return 0
	}

	// 计算平均负载
	meanLoad := totalLoad / count

	// 计算方差
	variance := 0.0
	for _, load := range loads {
		diff := load - meanLoad
		variance += diff * diff
	}
	variance /= count

	return variance
}

// calculateImbalance 计算不平衡度
func (lb *LoadBalancer) calculateImbalance() float64 {
	var minLoad, maxLoad float64
	first := true

	// 查找最大和最小负载
	for id, node := range lb.state.nodes {
		if node.Status == "offline" {
			continue
		}

		load := lb.calculateNodeLoad(id)
		if first {
			minLoad = load
			maxLoad = load
			first = false
		} else {
			minLoad = math.Min(minLoad, load)
			maxLoad = math.Max(maxLoad, load)
		}
	}

	if first { // 没有活跃节点
		return 0
	}

	// 计算不平衡度 (最大负载和最小负载之间的差异)
	// 归一化到 [0,1] 范围
	if maxLoad == 0 {
		return 0
	}
	return (maxLoad - minLoad) / maxLoad
}
