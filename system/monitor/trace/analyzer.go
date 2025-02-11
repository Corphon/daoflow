// system/monitor/trace/analyzer.go

package trace

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/Corphon/daoflow/core"
	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/types"
)

// 模式分析相关常量
const (
	defaultPatternThreshold = 0.7 // 默认模式偏差阈值
)

// 调用链分析相关常量
const (
	maxChainDepth = 100 // 最大调用链深度
	maxFanOut     = 50  // 最大扇出度
)

// 延迟分析相关常量
const (
	defaultLatencyThreshold = 50 * time.Millisecond  // 默认延迟阈值
	maxLatencyThreshold     = 100 * time.Millisecond // 最大延迟阈值
)

// 资源分析相关常量
const (
	defaultResourceThreshold = 0.8 // 默认资源使用阈值
)

// TraceAnalysis 追踪分析结果
type TraceAnalysis struct {
	ID        string
	Timestamp time.Time
	TraceID   types.TraceID
	Duration  time.Duration
	SpanCount int

	// 系统层面分析
	Patterns    []types.TracePattern
	Bottlenecks []types.Bottleneck
	Metrics     map[string]float64
	Anomalies   []types.Anomaly

	// 模型层面分析
	ModelAnalysis struct {
		State     model.ModelState
		Flow      model.FlowModel
		Patterns  []model.FlowPattern
		Metrics   model.ModelMetrics
		Anomalies []model.Anomaly
	}

	// 量子层面分析
	QuantumAnalysis struct {
		Entanglement float64
		Coherence    float64
		Phase        float64
		States       []*core.QuantumState
	}

	// 场动力学分析
	FieldAnalysis struct {
		Strength  float64
		Coupling  float64
		Resonance float64
		Evolution []*core.FieldState
	}
}

// Analyzer 追踪分析器
type Analyzer struct {
	mu sync.RWMutex

	// 配置
	config types.TraceConfig

	// 数据源
	tracker  *Tracker
	recorder *Recorder

	// 分析缓存
	cache struct {
		traces    map[types.TraceID]*TraceAnalysis
		patterns  []types.TracePattern
		anomalies []types.Anomaly
	}

	// 分析状态
	status struct {
		isRunning    bool
		lastAnalysis time.Time
		errors       []error
	}

	// 模型分析器
	modelAnalyzer *model.Analyzer
}

// QuantumAnalysis 量子分析结果
type QuantumAnalysis struct {
	Entanglement float64              // 量子纠缠度
	Coherence    float64              // 相干性
	Phase        float64              // 相位
	States       []*core.QuantumState // 修改为指针切片类型
}

// NewAnalyzer 创建新的分析器
func NewAnalyzer(tracker *Tracker, recorder *Recorder, config types.TraceConfig) *Analyzer {
	return &Analyzer{
		tracker:       tracker,
		recorder:      recorder,
		config:        config,
		modelAnalyzer: model.NewAnalyzer(),
		cache: struct {
			traces    map[types.TraceID]*TraceAnalysis
			patterns  []types.TracePattern
			anomalies []types.Anomaly
		}{
			traces: make(map[types.TraceID]*TraceAnalysis),
		},
	}
}

// Start 启动分析器
func (a *Analyzer) Start(ctx context.Context) error {
	a.mu.Lock()
	if a.status.isRunning {
		a.mu.Unlock()
		return model.WrapError(nil, model.ErrCodeOperation, "analyzer already running")
	}
	a.status.isRunning = true
	a.mu.Unlock()

	go a.analysisLoop(ctx)
	return nil
}

// analysisLoop 分析循环
func (a *Analyzer) analysisLoop(ctx context.Context) {
	ticker := time.NewTicker(a.config.AnalysisInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := a.analyze(ctx); err != nil {
				// 记录错误但继续运行
				a.mu.Lock()
				a.status.errors = append(a.status.errors, err)
				a.mu.Unlock()
			}
		}
	}
}

// Stop 停止分析器
func (a *Analyzer) Stop() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.status.isRunning {
		return model.WrapError(nil, model.ErrCodeOperation, "analyzer not running")
	}

	a.status.isRunning = false
	return nil
}

// analyze 执行分析
func (a *Analyzer) analyze(ctx context.Context) error {
	// 获取追踪数据
	traces := a.getTracesInWindow()

	for traceID, spans := range traces {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		analysis := &TraceAnalysis{
			ID:        generateAnalysisID(),
			Timestamp: time.Now(),
			TraceID:   traceID,
		}

		// 系统层面分析
		if err := a.analyzeSystemTrace(analysis, spans); err != nil {
			return model.WrapError(err, model.ErrCodeOperation, "system analysis failed")
		}

		// 模型层面分析
		if err := a.analyzeModelTrace(analysis, spans); err != nil {
			return model.WrapError(err, model.ErrCodeOperation, "model analysis failed")
		}

		// 量子层面分析
		if err := a.analyzeQuantumTrace(analysis, spans); err != nil {
			return model.WrapError(err, model.ErrCodeOperation, "quantum analysis failed")
		}

		// 场动力学分析
		if err := a.analyzeFieldTrace(analysis, spans); err != nil {
			return model.WrapError(err, model.ErrCodeOperation, "field analysis failed")
		}

		// 缓存分析结果
		a.cacheAnalysis(analysis)
	}

	return nil
}

// getTracesInWindow 获取时间窗口内的追踪数据
func (a *Analyzer) getTracesInWindow() map[types.TraceID][]*Span {
	a.mu.RLock()
	defer a.mu.RUnlock()

	traces := make(map[types.TraceID][]*Span)
	cutoff := time.Now().Add(-a.config.AnalysisInterval)

	// 从recorder获取原始数据
	records := a.recorder.GetRecords()

	// 按TraceID分组并过滤时间窗口
	for _, record := range records {
		if record.Timestamp.After(cutoff) {
			traces[record.TraceID] = append(traces[record.TraceID], record.Data.(*Span))
		}
	}

	return traces
}

// generateAnalysisID 生成分析ID
func generateAnalysisID() string {
	return fmt.Sprintf("analysis-%d", time.Now().UnixNano())
}

// analyzeSystemTrace 分析系统层面的追踪
func (a *Analyzer) analyzeSystemTrace(analysis *TraceAnalysis, spans []*Span) error {
	// 检测系统模式
	patterns := a.detectSystemPatterns(spans)
	analysis.Patterns = patterns

	// 检测瓶颈
	bottlenecks := a.detectBottlenecks(spans)
	analysis.Bottlenecks = bottlenecks

	// 计算指标
	metrics := a.calculateSystemMetrics(spans)
	analysis.Metrics = metrics

	// 检测异常
	anomalies := a.detectSystemAnomalies(spans, patterns)
	analysis.Anomalies = anomalies

	return nil
}

// detectSystemPatterns 检测系统模式
func (a *Analyzer) detectSystemPatterns(spans []*Span) []types.TracePattern {
	patterns := make([]types.TracePattern, 0)

	// 基于时间窗口分组
	groups := groupSpansByTime(spans, a.config.AnalysisInterval)

	// 对每个时间窗口进行模式检测
	for _, group := range groups {
		// 检测执行路径模式
		if pattern := detectPathPattern(group); pattern != nil {
			patterns = append(patterns, *pattern)
		}

		// 检测调用链模式
		if pattern := detectChainPattern(group); pattern != nil {
			patterns = append(patterns, *pattern)
		}
	}

	return patterns
}

// groupSpansByTime 按时间窗口对跨度分组
func groupSpansByTime(spans []*Span, window time.Duration) [][]*Span {
	groups := make([][]*Span, 0)
	if len(spans) == 0 {
		return groups
	}

	// 按开始时间排序
	sort.Slice(spans, func(i, j int) bool {
		return spans[i].StartTime.Before(spans[j].StartTime)
	})

	currentGroup := []*Span{spans[0]}
	groupStart := spans[0].StartTime

	for i := 1; i < len(spans); i++ {
		if spans[i].StartTime.Sub(groupStart) > window {
			groups = append(groups, currentGroup)
			currentGroup = []*Span{spans[i]}
			groupStart = spans[i].StartTime
		} else {
			currentGroup = append(currentGroup, spans[i])
		}
	}

	groups = append(groups, currentGroup)
	return groups
}

// detectPathPattern 检测执行路径模式
func detectPathPattern(spans []*Span) *types.TracePattern {
	if len(spans) < 2 {
		return nil
	}

	// 构建路径图
	graph := buildPathGraph(spans)

	// 分析路径特征
	if pattern := analyzePathPattern(graph); pattern != nil {
		pattern.Type = "execution_path"
		pattern.StartTime = spans[0].StartTime
		pattern.EndTime = spans[len(spans)-1].EndTime
		return pattern
	}

	return nil
}

// PathGraph 路径图结构
type PathGraph struct {
	Nodes map[string]*Span    // 节点集合
	Edges map[string][]string // 边集合
	Entry string              // 入口节点
	Exit  string              // 出口节点
}

// buildPathGraph 构建路径图
func buildPathGraph(spans []*Span) *PathGraph {
	graph := &PathGraph{
		Nodes: make(map[string]*Span),
		Edges: make(map[string][]string),
	}

	for _, span := range spans {
		// 将SpanID转换为string
		spanID := string(span.ID)
		parentID := string(span.ParentID)

		graph.Nodes[spanID] = span

		// 构建边关系
		if parentID != "" {
			graph.Edges[parentID] = append(
				graph.Edges[parentID],
				spanID,
			)
		}
	}

	// 识别入口和出口
	graph.Entry = string(spans[0].ID)
	graph.Exit = string(spans[len(spans)-1].ID)

	return graph
}

// analyzePathPattern 分析路径特征
func analyzePathPattern(graph *PathGraph) *types.TracePattern {
	if graph == nil {
		return nil
	}

	// 提取路径特征
	pattern := &types.TracePattern{
		ID:         generateAnalysisID(),
		Type:       "execution_path",
		Properties: make(map[string]interface{}),
	}

	// 分析路径特征
	pattern.Properties["path_length"] = len(graph.Nodes)
	pattern.Properties["branch_count"] = countBranches(graph)
	pattern.Properties["max_depth"] = calculatePathDepth(graph)

	// 计算置信度
	pattern.Confidence = calculatePathConfidence(graph)

	return pattern
}

// countBranches 统计分支数量
func countBranches(graph *PathGraph) int {
	branchCount := 0
	for _, edges := range graph.Edges {
		if len(edges) > 1 {
			branchCount += len(edges) - 1
		}
	}
	return branchCount
}

// calculatePathDepth 计算路径最大深度
func calculatePathDepth(graph *PathGraph) int {
	visited := make(map[string]int)
	return dfsDepth(graph, graph.Entry, visited)
}

// dfsDepth DFS辅助函数
func dfsDepth(graph *PathGraph, node string, visited map[string]int) int {
	if depth, ok := visited[node]; ok {
		return depth
	}

	maxDepth := 0
	for _, next := range graph.Edges[node] {
		depth := dfsDepth(graph, next, visited)
		if depth > maxDepth {
			maxDepth = depth
		}
	}

	visited[node] = maxDepth + 1
	return maxDepth + 1
}

// calculatePathConfidence 计算路径置信度
func calculatePathConfidence(graph *PathGraph) float64 {
	// 基于路径特征计算置信度
	branchFactor := float64(countBranches(graph)) / float64(len(graph.Nodes))
	depthFactor := float64(calculatePathDepth(graph)) / float64(len(graph.Nodes))

	// 综合评估置信度
	return (1.0 - branchFactor*0.3) * (1.0 - depthFactor*0.2)
}

// detectChainPattern 检测调用链模式
func detectChainPattern(spans []*Span) *types.TracePattern {
	if len(spans) < 2 {
		return nil
	}

	// 构建调用链
	chain := buildCallChain(spans)

	// 分析链路特征
	if pattern := analyzeChainPattern(chain); pattern != nil {
		pattern.Type = "call_chain"
		pattern.StartTime = spans[0].StartTime
		pattern.EndTime = spans[len(spans)-1].EndTime
		return pattern
	}

	return nil
}

// CallChain 调用链结构
type CallChain struct {
	Root     *Span               // 根节点
	Nodes    map[string]*Span    // 所有节点
	Children map[string][]string // 子节点关系
	Depth    int                 // 调用深度
}

// buildCallChain 构建调用链
func buildCallChain(spans []*Span) *CallChain {
	chain := &CallChain{
		Nodes:    make(map[string]*Span),
		Children: make(map[string][]string),
	}

	// 构建节点映射
	for _, span := range spans {
		spanID := string(span.ID)
		chain.Nodes[spanID] = span

		// 处理父子关系
		if span.ParentID != "" {
			parentID := string(span.ParentID)
			chain.Children[parentID] = append(chain.Children[parentID], spanID)
		} else {
			chain.Root = span
		}
	}

	// 计算调用深度
	chain.Depth = calculateChainDepth(chain)

	return chain
}

// analyzeChainPattern 分析调用链特征
func analyzeChainPattern(chain *CallChain) *types.TracePattern {
	if chain == nil || chain.Root == nil {
		return nil
	}

	pattern := &types.TracePattern{
		ID:         generateAnalysisID(),
		Type:       "call_chain",
		Properties: make(map[string]interface{}),
	}

	// 分析链路特征
	pattern.Properties["chain_depth"] = chain.Depth
	pattern.Properties["node_count"] = len(chain.Nodes)
	pattern.Properties["fan_out"] = calculateFanOut(chain)

	// 计算置信度
	pattern.Confidence = calculateChainConfidence(chain)

	return pattern
}

// calculateFanOut 计算调用链的扇出度
func calculateFanOut(chain *CallChain) float64 {
	if len(chain.Children) == 0 {
		return 0
	}

	// 计算平均子节点数
	totalChildren := 0
	for _, children := range chain.Children {
		totalChildren += len(children)
	}

	return float64(totalChildren) / float64(len(chain.Children))
}

// calculateChainConfidence 计算调用链置信度
func calculateChainConfidence(chain *CallChain) float64 {
	if chain == nil || chain.Root == nil {
		return 0
	}

	// 基于深度和扇出度计算置信度
	depth := float64(chain.Depth)
	fanOut := calculateFanOut(chain)

	// 深度和扇出度的权重
	depthWeight := 0.6
	fanOutWeight := 0.4

	// 计算归一化的置信度
	confidence := (depth*depthWeight + fanOut*fanOutWeight) /
		(maxChainDepth*depthWeight + maxFanOut*fanOutWeight)

	return math.Max(0, math.Min(1, confidence))
}

// calculateChainDepth 计算调用链深度
func calculateChainDepth(chain *CallChain) int {
	if chain == nil || chain.Root == nil {
		return 0
	}

	// 使用DFS计算最大深度
	depth := 0
	visited := make(map[string]int)

	// 从根节点开始DFS遍历
	rootID := string(chain.Root.ID)
	depth = dfsChainDepth(chain, rootID, visited)

	return depth
}

// dfsChainDepth DFS计算深度辅助函数
func dfsChainDepth(chain *CallChain, nodeID string, visited map[string]int) int {
	// 已访问过的节点直接返回其深度
	if depth, ok := visited[nodeID]; ok {
		return depth
	}

	maxChildDepth := 0
	// 遍历所有子节点
	for _, childID := range chain.Children[nodeID] {
		childDepth := dfsChainDepth(chain, childID, visited)
		if childDepth > maxChildDepth {
			maxChildDepth = childDepth
		}
	}

	// 当前节点深度为最大子节点深度+1
	depth := maxChildDepth + 1
	visited[nodeID] = depth
	return depth
}

// detectBottlenecks 检测系统瓶颈
func (a *Analyzer) detectBottlenecks(spans []*Span) []types.Bottleneck {
	bottlenecks := make([]types.Bottleneck, 0)

	// 检测延迟瓶颈
	if b := detectLatencyBottleneck(spans); b != nil {
		bottlenecks = append(bottlenecks, *b)
	}

	// 检测资源瓶颈
	if b := detectResourceBottleneck(spans); b != nil {
		bottlenecks = append(bottlenecks, *b)
	}

	return bottlenecks
}

// detectLatencyBottleneck 检测延迟瓶颈
func detectLatencyBottleneck(spans []*Span) *types.Bottleneck {
	if len(spans) == 0 {
		return nil
	}

	// 计算平均延迟和标准差
	var totalLatency time.Duration
	for _, span := range spans {
		totalLatency += span.Duration
	}
	avgLatency := totalLatency / time.Duration(len(spans))

	// 如果平均延迟超过阈值则判定为瓶颈
	if avgLatency > defaultLatencyThreshold {
		return &types.Bottleneck{
			Type:     "latency",
			Resource: "system",
			Severity: calculateLatencySeverity(avgLatency),
			Duration: avgLatency,
		}
	}
	return nil
}

// calculateLatencySeverity 计算延迟严重程度
func calculateLatencySeverity(latency time.Duration) float64 {
	// 根据延迟时间计算严重程度 0-1
	normalized := float64(latency) / float64(maxLatencyThreshold)
	return math.Max(0, math.Min(1, normalized))
}

// detectResourceBottleneck 检测资源瓶颈
func detectResourceBottleneck(spans []*Span) *types.Bottleneck {
	// 统计资源使用
	resourceUsage := calculateResourceUsage(spans)

	// 检查是否超过阈值
	for resource, usage := range resourceUsage {
		if usage > defaultResourceThreshold {
			return &types.Bottleneck{
				Type:     "resource",
				Resource: resource,
				Severity: calculateResourceSeverity(usage),
				Impact:   usage,
			}
		}
	}
	return nil
}

// calculateResourceUsage 计算资源使用情况
func calculateResourceUsage(spans []*Span) map[string]float64 {
	usage := make(map[string]float64)
	if len(spans) == 0 {
		return usage
	}

	// 统计资源使用
	for _, span := range spans {
		if cpu, ok := span.Metrics["cpu_usage"]; ok {
			usage["cpu"] += cpu
		}
		if mem, ok := span.Metrics["memory_usage"]; ok {
			usage["memory"] += mem
		}
	}

	// 计算平均使用率
	count := float64(len(spans))
	for resource := range usage {
		usage[resource] /= count
	}

	return usage
}

// calculateResourceSeverity 计算资源瓶颈严重程度
func calculateResourceSeverity(usage float64) float64 {
	// 基于使用率计算严重程度 0-1
	return math.Max(0, math.Min(1, (usage-defaultResourceThreshold)/(1-defaultResourceThreshold)))
}

// calculateSystemMetrics 计算系统指标
func (a *Analyzer) calculateSystemMetrics(spans []*Span) map[string]float64 {
	metrics := make(map[string]float64)

	// 计算基础指标
	metrics["request_count"] = float64(len(spans))
	metrics["error_rate"] = calculateErrorRate(spans)
	metrics["avg_latency"] = calculateAvgLatency(spans)

	// 计算资源指标
	metrics["cpu_usage"] = calculateCPUUsage(spans)
	metrics["memory_usage"] = calculateMemoryUsage(spans)

	return metrics
}

// calculateErrorRate 计算错误率
func calculateErrorRate(spans []*Span) float64 {
	if len(spans) == 0 {
		return 0.0
	}

	errorCount := 0
	for _, span := range spans {
		if span.Status == types.SpanStatusError {
			errorCount++
		}
	}

	return float64(errorCount) / float64(len(spans))
}

// calculateAvgLatency 计算平均延迟
func calculateAvgLatency(spans []*Span) float64 {
	if len(spans) == 0 {
		return 0.0
	}

	var totalLatency time.Duration
	for _, span := range spans {
		totalLatency += span.Duration
	}

	return float64(totalLatency.Milliseconds()) / float64(len(spans))
}

// calculateCPUUsage 计算CPU使用率
func calculateCPUUsage(spans []*Span) float64 {
	if len(spans) == 0 {
		return 0.0
	}

	var totalCPU float64
	for _, span := range spans {
		if cpu, ok := span.Metrics["cpu_usage"]; ok {
			totalCPU += cpu
		}
	}

	return totalCPU / float64(len(spans))
}

// calculateMemoryUsage 计算内存使用率
func calculateMemoryUsage(spans []*Span) float64 {
	if len(spans) == 0 {
		return 0.0
	}

	var totalMemory float64
	for _, span := range spans {
		if mem, ok := span.Metrics["memory_usage"]; ok {
			totalMemory += mem
		}
	}

	return totalMemory / float64(len(spans))
}

// detectSystemAnomalies 检测系统异常
func (a *Analyzer) detectSystemAnomalies(spans []*Span, patterns []types.TracePattern) []types.Anomaly {
	anomalies := make([]types.Anomaly, 0)

	// 检测性能异常
	if anomaly := detectPerformanceAnomaly(spans); anomaly != nil {
		anomalies = append(anomalies, *anomaly)
	}

	// 检测模式异常 - 移除spans参数
	if anomaly := detectPatternAnomaly(patterns); anomaly != nil {
		anomalies = append(anomalies, *anomaly)
	}

	return anomalies
}

// detectPerformanceAnomaly 检测性能异常
func detectPerformanceAnomaly(spans []*Span) *types.Anomaly {
	if len(spans) == 0 {
		return nil
	}

	// 计算平均延迟
	avgLatency := calculateAvgLatency(spans)
	if avgLatency > float64(defaultLatencyThreshold) {
		return &types.Anomaly{
			Type:       "performance",
			Severity:   calculateLatencySeverity(time.Duration(avgLatency) * time.Millisecond),
			Metric:     "latency",
			Threshold:  float64(defaultLatencyThreshold),
			Value:      avgLatency,
			DetectedAt: time.Now(),
		}
	}

	return nil
}

// detectPatternAnomaly 检测模式异常
func detectPatternAnomaly(patterns []types.TracePattern) *types.Anomaly {
	if len(patterns) == 0 {
		return nil
	}

	// 分析模式偏差
	deviation := calculatePatternDeviation(patterns)
	if deviation > defaultPatternThreshold {
		return &types.Anomaly{
			Type:       "pattern",
			Severity:   deviation,
			Metric:     "pattern_deviation",
			Threshold:  defaultPatternThreshold,
			Value:      deviation,
			DetectedAt: time.Now(),
		}
	}

	return nil
}

// calculatePatternDeviation 计算模式偏差
func calculatePatternDeviation(patterns []types.TracePattern) float64 {
	if len(patterns) < 2 {
		return 0
	}

	// 计算基准模式
	baseline := calculateBaselinePattern(patterns)

	// 计算偏差
	totalDeviation := 0.0
	for _, pattern := range patterns {
		deviation := calculateSinglePatternDeviation(pattern, baseline)
		totalDeviation += deviation
	}

	return totalDeviation / float64(len(patterns))
}

// calculateBaselinePattern 计算基准模式
func calculateBaselinePattern(patterns []types.TracePattern) map[string]float64 {
	baseline := make(map[string]float64)

	// 计算关键指标的平均值
	for _, pattern := range patterns {
		for key, value := range pattern.Properties {
			if v, ok := value.(float64); ok {
				baseline[key] += v
			}
		}
	}

	// 归一化
	for key := range baseline {
		baseline[key] /= float64(len(patterns))
	}

	return baseline
}

// calculateSinglePatternDeviation 计算单个模式的偏差
func calculateSinglePatternDeviation(pattern types.TracePattern, baseline map[string]float64) float64 {
	deviation := 0.0
	count := 0.0

	// 计算各指标偏差
	for key, value := range pattern.Properties {
		if v, ok := value.(float64); ok {
			if baseValue, exists := baseline[key]; exists {
				deviation += math.Abs(v - baseValue)
				count++
			}
		}
	}

	if count == 0 {
		return 0
	}

	return deviation / count
}

// analyzeModelTrace 分析模型层面的追踪
func (a *Analyzer) analyzeModelTrace(analysis *TraceAnalysis, spans []*Span) error {
	modelSpans := a.filterModelSpans(spans)
	if len(modelSpans) == 0 {
		return nil
	}

	// 分析模型状态
	state := a.analyzeModelState(modelSpans)
	analysis.ModelAnalysis.State = state

	// 分析流状态
	flow := a.analyzeModelFlow(modelSpans)
	analysis.ModelAnalysis.Flow = flow

	// 检测模型模式
	patterns := a.modelAnalyzer.DetectPatterns(modelSpans)
	analysis.ModelAnalysis.Patterns = patterns

	// 计算模型指标
	metrics := a.modelAnalyzer.CalculateMetrics(modelSpans)
	analysis.ModelAnalysis.Metrics = metrics

	// 检测模型异常
	anomalies := a.modelAnalyzer.DetectAnomalies(modelSpans)
	analysis.ModelAnalysis.Anomalies = anomalies

	return nil
}

// analyzeModelState 分析模型状态
func (a *Analyzer) analyzeModelState(spans []*Span) model.ModelState {
	state := model.ModelState{}

	// 获取最后一个模型状态
	for _, span := range spans {
		if span.ModelState != nil {
			state = *span.ModelState
			break
		}
	}
	return state
}

// analyzeModelFlow 分析流状态
func (a *Analyzer) analyzeModelFlow(spans []*Span) model.FlowModel {
	var flow model.FlowModel

	// 获取最后一个流状态
	for _, span := range spans {
		if span.ModelType != model.ModelTypeNone {
			flow = span.ModelFlow
			break
		}
	}
	return flow
}

// analyzeQuantumTrace 分析量子层面的追踪
func (a *Analyzer) analyzeQuantumTrace(analysis *TraceAnalysis, spans []*Span) error {
	// 提取量子态相关的跨度
	quantumSpans := a.filterQuantumSpans(spans)
	if len(quantumSpans) == 0 {
		return nil
	}

	// 分析量子纠缠
	entanglement := a.calculateEntanglement(quantumSpans)
	analysis.QuantumAnalysis.Entanglement = entanglement

	// 分析相干性
	coherence := a.calculateCoherence(quantumSpans)
	analysis.QuantumAnalysis.Coherence = coherence

	// 分析相位
	phase := a.calculatePhase(quantumSpans)
	analysis.QuantumAnalysis.Phase = phase

	// 提取量子态序列
	states := a.extractQuantumStates(quantumSpans)
	analysis.QuantumAnalysis.States = states

	return nil
}

// extractQuantumStates 从跨度中提取量子态序列
func (a *Analyzer) extractQuantumStates(spans []*Span) []*core.QuantumState {
	// 改为指针切片
	states := make([]*core.QuantumState, 0)

	for _, span := range spans {
		state, ok := span.Fields["quantum_state"].(*core.QuantumState)
		if !ok {
			continue
		}
		// 直接append指针
		states = append(states, state)
	}

	return states
}

// extractFieldEvolution 从追踪跨度中提取场态演化序列
func (a *Analyzer) extractFieldEvolution(spans []*Span) []*core.FieldState {
	// 改为指针切片
	states := make([]*core.FieldState, 0)

	for _, span := range spans {
		// 获取场态
		state, ok := span.Fields["field_state"].(*core.FieldState)
		if !ok {
			continue
		}

		// 直接append指针
		states = append(states, state)
	}

	// 按时间排序
	sort.Slice(states, func(i, j int) bool {
		return states[i].Timestamp.Before(states[j].Timestamp)
	})

	return states
}

// analyzeFieldTrace 分析场动力学追踪
func (a *Analyzer) analyzeFieldTrace(analysis *TraceAnalysis, spans []*Span) error {
	// 提取场相关的跨度
	fieldSpans := a.filterFieldSpans(spans)
	if len(fieldSpans) == 0 {
		return nil
	}

	// 分析场强度
	strength := a.calculateFieldStrength(fieldSpans)
	analysis.FieldAnalysis.Strength = strength

	// 分析场耦合
	coupling := a.calculateFieldCoupling(fieldSpans)
	analysis.FieldAnalysis.Coupling = coupling

	// 分析共振
	resonance := a.calculateResonance(fieldSpans)
	analysis.FieldAnalysis.Resonance = resonance

	// 提取场态演化序列
	evolution := a.extractFieldEvolution(fieldSpans)
	analysis.FieldAnalysis.Evolution = evolution

	return nil
}

// 过滤器方法
func (a *Analyzer) filterModelSpans(spans []*Span) []*Span {
	modelSpans := make([]*Span, 0)
	for _, span := range spans {
		if span.ModelType != model.ModelTypeNone {
			modelSpans = append(modelSpans, span)
		}
	}
	return modelSpans
}

func (a *Analyzer) filterQuantumSpans(spans []*Span) []*Span {
	quantumSpans := make([]*Span, 0)
	for _, span := range spans {
		if _, ok := span.Fields["quantum_state"]; ok {
			quantumSpans = append(quantumSpans, span)
		}
	}
	return quantumSpans
}

func (a *Analyzer) filterFieldSpans(spans []*Span) []*Span {
	fieldSpans := make([]*Span, 0)
	for _, span := range spans {
		if _, ok := span.Fields["field_state"]; ok {
			fieldSpans = append(fieldSpans, span)
		}
	}
	return fieldSpans
}

// 缓存方法
func (a *Analyzer) cacheAnalysis(analysis *TraceAnalysis) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.cache.traces[analysis.TraceID] = analysis
	a.status.lastAnalysis = analysis.Timestamp
}

// 辅助方法
func (a *Analyzer) calculateEntanglement(spans []*Span) float64 {
	if len(spans) < 2 {
		return 0.0
	}

	var totalEntanglement float64
	pairCount := 0

	// 计算所有量子态对之间的纠缠度
	for i := 0; i < len(spans)-1; i++ {
		for j := i + 1; j < len(spans); j++ {
			state1, ok1 := spans[i].Fields["quantum_state"].(*core.QuantumState)
			state2, ok2 := spans[j].Fields["quantum_state"].(*core.QuantumState)

			if !ok1 || !ok2 {
				continue
			}

			// 计算两个量子态之间的纠缠度
			entanglement := calculatePairEntanglement(state1, state2)
			totalEntanglement += entanglement
			pairCount++
		}
	}

	if pairCount == 0 {
		return 0.0
	}

	return totalEntanglement / float64(pairCount)
}

// calculatePairEntanglement 计算两个量子态之间的纠缠度
func calculatePairEntanglement(state1, state2 *core.QuantumState) float64 {
	// 计算态矢量的内积
	overlap, err := state1.DotProduct(state2)
	if err != nil {
		return 0.0 // 如果计算内积失败，返回0纠缠度
	}

	// 计算纠缠度：使用冯诺依曼熵的简化形式
	// E = -Tr(ρ log₂ ρ)，其中 ρ 是约化密度矩阵
	p := math.Abs(real(overlap))
	if p == 0 || p == 1 {
		return 0.0
	}

	entropy := -p*math.Log2(p) - (1-p)*math.Log2(1-p)
	return entropy
}

func (a *Analyzer) calculateCoherence(spans []*Span) float64 {
	if len(spans) == 0 {
		return 0.0
	}

	var totalCoherence float64
	validSpans := 0

	for _, span := range spans {
		state, ok := span.Fields["quantum_state"].(*core.QuantumState)
		if !ok {
			continue
		}

		// 计算相干性：使用态的纯度作为相干性度量
		// C = |⟨ψ|ψ⟩|²
		purity := state.CalculatePurity()
		phaseStability := calculatePhaseStability(state)

		// 综合考虑纯度和相位稳定性
		coherence := math.Sqrt(purity * phaseStability)
		totalCoherence += coherence
		validSpans++
	}

	if validSpans == 0 {
		return 0.0
	}

	return totalCoherence / float64(validSpans)
}

func (a *Analyzer) calculatePhase(spans []*Span) float64 {
	if len(spans) == 0 {
		return 0.0
	}

	var phases []float64
	for _, span := range spans {
		state, ok := span.Fields["quantum_state"].(*core.QuantumState)
		if !ok {
			continue
		}

		// 获取量子态的相位
		phase := state.GetPhase()
		phases = append(phases, phase)
	}

	if len(phases) == 0 {
		return 0.0
	}

	// 计算平均相位和相位相干性
	avgPhase := calculateAveragePhase(phases)
	phaseCoherence := calculatePhaseCoherence(phases, avgPhase)

	return normalizePhase(avgPhase * phaseCoherence)
}

func (a *Analyzer) calculateFieldStrength(spans []*Span) float64 {
	if len(spans) == 0 {
		return 0.0
	}

	var totalStrength float64
	var weights float64

	for _, span := range spans {
		state, ok := span.Fields["field_state"].(*core.FieldState)
		if !ok {
			continue
		}

		// 获取基本场强
		baseStrength := state.GetStrength()

		// 考虑时间衰减
		timeFactor := calculateTimeFactor(span.StartTime, span.EndTime)

		// 考虑空间分布
		spatialFactor := calculateSpatialFactor(state.GetDistribution())

		// 计算加权场强
		weightedStrength := baseStrength * timeFactor * spatialFactor

		// 累加加权和
		weight := span.Duration.Seconds()
		totalStrength += weightedStrength * weight
		weights += weight
	}

	if weights == 0 {
		return 0.0
	}

	return totalStrength / weights
}

func (a *Analyzer) calculateFieldCoupling(spans []*Span) float64 {
	if len(spans) < 2 {
		return 0.0
	}

	var totalCoupling float64
	var couplingCount int

	// 计算场之间的耦合强度
	for i := 0; i < len(spans)-1; i++ {
		for j := i + 1; j < len(spans); j++ {
			field1, ok1 := spans[i].Fields["field_state"].(*core.FieldState)
			field2, ok2 := spans[j].Fields["field_state"].(*core.FieldState)

			if !ok1 || !ok2 {
				continue
			}

			// 计算场间耦合
			coupling := calculateFieldInteraction(field1, field2)

			// 考虑时空关联
			spacetimeFactor := calculateSpacetimeCorrelation(spans[i], spans[j])

			totalCoupling += coupling * spacetimeFactor
			couplingCount++
		}
	}

	if couplingCount == 0 {
		return 0.0
	}

	return totalCoupling / float64(couplingCount)
}

func (a *Analyzer) calculateResonance(spans []*Span) float64 {
	if len(spans) == 0 {
		return 0.0
	}

	// 提取场的频率特征
	frequencies := make([]float64, 0)
	amplitudes := make([]float64, 0)

	for _, span := range spans {
		field, ok := span.Fields["field_state"].(*core.FieldState)
		if !ok {
			continue
		}

		freq := field.GetFrequency()
		amp := field.GetAmplitude()

		frequencies = append(frequencies, freq)
		amplitudes = append(amplitudes, amp)
	}

	if len(frequencies) == 0 {
		return 0.0
	}

	// 计算共振强度
	resonanceStrength := calculateResonanceStrength(frequencies, amplitudes)

	// 计算共振稳定性
	resonanceStability := calculateResonanceStability(frequencies)

	// 综合评估
	return math.Sqrt(resonanceStrength * resonanceStability)
}

// 辅助函数
func calculateAveragePhase(phases []float64) float64 {
	if len(phases) == 0 {
		return 0.0
	}

	sumSin := 0.0
	sumCos := 0.0
	for _, phase := range phases {
		sumSin += math.Sin(phase)
		sumCos += math.Cos(phase)
	}

	return math.Atan2(sumSin, sumCos)
}

// 修改analyzer.go中的calculatePhaseStability:
func calculatePhaseStability(state *core.QuantumState) float64 {
	phaseVariation := state.GetPhaseVariation()
	return math.Exp(-phaseVariation * phaseVariation)
}
func calculatePhaseCoherence(phases []float64, avgPhase float64) float64 {
	if len(phases) == 0 {
		return 0.0
	}

	var sumDeviation float64
	for _, phase := range phases {
		deviation := math.Abs(normalizePhase(phase - avgPhase))
		sumDeviation += deviation * deviation
	}

	return math.Exp(-sumDeviation / float64(len(phases)))
}

func normalizePhase(phase float64) float64 {
	// 将相位标准化到 [-π, π] 区间
	normalized := math.Mod(phase, 2*math.Pi)
	if normalized > math.Pi {
		normalized -= 2 * math.Pi
	} else if normalized < -math.Pi {
		normalized += 2 * math.Pi
	}
	return normalized
}

func calculateTimeFactor(start, end time.Time) float64 {
	duration := end.Sub(start).Seconds()
	return math.Exp(-duration / 3600.0) // 1小时特征时间
}

func calculateSpatialFactor(distribution []float64) float64 {
	if len(distribution) == 0 {
		return 1.0
	}

	// 计算空间分布的均匀性
	var sum, sumSq float64
	for _, value := range distribution {
		sum += value
		sumSq += value * value
	}

	mean := sum / float64(len(distribution))
	variance := sumSq/float64(len(distribution)) - mean*mean

	return 1.0 / (1.0 + variance)
}

func calculateFieldInteraction(field1, field2 *core.FieldState) float64 {
	// 计算场之间的相互作用强度
	overlap := field1.CalculateOverlap(field2)
	strength := math.Sqrt(field1.GetStrength() * field2.GetStrength())
	return overlap * strength
}

func calculateSpacetimeCorrelation(span1, span2 *Span) float64 {
	// 计算时空关联度
	timeCorr := calculateTimeCorrelation(span1.StartTime, span1.EndTime,
		span2.StartTime, span2.EndTime)
	spaceCorr := calculateSpaceCorrelation(span1, span2)
	return math.Sqrt(timeCorr * spaceCorr)
}

// calculateTimeCorrelation 计算时间相关性
func calculateTimeCorrelation(start1, end1, start2, end2 time.Time) float64 {
	// 计算时间重叠度
	overlap := math.Min(end1.Sub(start1).Seconds(), end2.Sub(start2).Seconds())
	if overlap <= 0 {
		return 0
	}
	return math.Exp(-overlap / 3600.0) // 1小时特征时间
}

// calculateSpaceCorrelation 计算空间相关性
func calculateSpaceCorrelation(span1, span2 *Span) float64 {
	// 通过场状态分布计算空间相关性
	if field1, ok1 := span1.Fields["field_state"].(*core.FieldState); ok1 {
		if field2, ok2 := span2.Fields["field_state"].(*core.FieldState); ok2 {
			return field1.CalculateOverlap(field2)
		}
	}
	return 0
}
func calculateResonanceStrength(frequencies, amplitudes []float64) float64 {
	if len(frequencies) != len(amplitudes) || len(frequencies) == 0 {
		return 0.0
	}

	// 计算频率匹配度和振幅增强
	var resonanceSum float64
	for i := 0; i < len(frequencies)-1; i++ {
		for j := i + 1; j < len(frequencies); j++ {
			freqMatch := calculateFrequencyMatch(frequencies[i], frequencies[j])
			ampProduct := amplitudes[i] * amplitudes[j]
			resonanceSum += freqMatch * ampProduct
		}
	}

	return resonanceSum / float64(len(frequencies)*(len(frequencies)-1)/2)
}

func calculateFrequencyMatch(f1, f2 float64) float64 {
	// 计算频率匹配度，使用高斯函数
	diff := math.Abs(f1 - f2)
	return math.Exp(-diff * diff / (2.0 * 0.1)) // 0.1是带宽参数
}

func calculateResonanceStability(frequencies []float64) float64 {
	if len(frequencies) < 2 {
		return 1.0
	}

	// 计算频率的稳定性
	var sum, sumSq float64
	for _, f := range frequencies {
		sum += f
		sumSq += f * f
	}

	mean := sum / float64(len(frequencies))
	variance := sumSq/float64(len(frequencies)) - mean*mean

	return 1.0 / (1.0 + variance)
}
