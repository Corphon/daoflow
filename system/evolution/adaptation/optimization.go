//system/evolution/adaptation/optimization.go

package adaptation

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system/types"
)

const (
	maxConvergenceHistory = 1000
)

// AdaptiveOptimization 适应性优化系统
type AdaptiveOptimization struct {
	mu sync.RWMutex

	// 基础配置
	config struct {
		optimizationInterval time.Duration // 优化间隔
		improvementThreshold float64       // 改进阈值
		maxIterations        int           // 最大迭代次数
		convergenceRate      float64       // 收敛率
	}

	// 优化状态
	state struct {
		optimizations map[string]*Optimization          // 当前优化
		objectives    map[string]*OptimizationObjective // 优化目标
		constraints   []OptimizationConstraint          // 优化约束
		metrics       OptimizationMetrics               // 优化指标
	}

	// 依赖项
	strategy *AdaptationStrategy
	learning *AdaptiveLearning
}

// Optimization 优化实例
type Optimization struct {
	ID         string                  // 优化ID
	Type       string                  // 优化类型
	Target     string                  // 优化目标
	Parameters []OptimizationParameter // 优化参数
	Progress   OptimizationProgress    // 优化进度
	Result     OptimizationResult      // 优化结果
	StartTime  time.Time               // 开始时间
	LastUpdate time.Time               // 最后更新
}

// OptimizationParameter 优化参数
type OptimizationParameter struct {
	Name         string      // 参数名称
	CurrentValue interface{} // 当前值
	Range        [2]float64  // 取值范围
	Step         float64     // 调整步长
	Weight       float64     // 权重
}

// OptimizationProgress 优化进度
type OptimizationProgress struct {
	Iteration    int     // 当前迭代
	BestValue    float64 // 最佳值
	CurrentValue float64 // 当前值
	Improvement  float64 // 改进幅度
	Status       string  // 优化状态
}

// OptimizationResult 优化结果
type OptimizationResult struct {
	Success     bool                   // 是否成功
	FinalValue  float64                // 最终值
	Parameters  map[string]interface{} // 最优参数
	Performance map[string]float64     // 性能指标
	Duration    time.Duration          // 优化时长
}

// OptimizationObjective 优化目标
type OptimizationObjective struct {
	ID        string             // 目标ID
	Name      string             // 目标名称
	Type      string             // 目标类型
	TargetID  string             // 目标标识符
	Target    float64            // 目标值
	Weight    float64            // 权重
	Evaluator ObjectiveEvaluator // 评估函数
}

// OptimizationConstraint 优化约束
type OptimizationConstraint struct {
	Type      string      // 约束类型
	Target    string      // 约束目标
	Condition string      // 约束条件
	Value     interface{} // 约束值
	Priority  int         // 优先级
}

// ObjectiveEvaluator 目标评估器
type ObjectiveEvaluator struct {
	Function   func(interface{}) float64 // 评估函数
	Parameters map[string]interface{}    // 评估参数
	Threshold  float64                   // 评估阈值
}

// OptimizationMetrics 优化指标
type OptimizationMetrics struct {
	TotalOptimizations int                // 总优化次数
	SuccessRate        float64            // 成功率
	AverageImprovement float64            // 平均改进
	Convergence        []ConvergencePoint // 收敛曲线
}

// ConvergencePoint 收敛点
type ConvergencePoint struct {
	Iteration   int
	Value       float64
	Improvement float64
}

// -----------------------------------------------------
// NewAdaptiveOptimization 创建新的适应性优化系统
func NewAdaptiveOptimization(
	strategy *AdaptationStrategy,
	learning *AdaptiveLearning) *AdaptiveOptimization {

	ao := &AdaptiveOptimization{
		strategy: strategy,
		learning: learning,
	}

	// 初始化配置
	ao.config.optimizationInterval = 1 * time.Hour
	ao.config.improvementThreshold = 0.01
	ao.config.maxIterations = 100
	ao.config.convergenceRate = 0.001

	// 初始化状态
	ao.state.optimizations = make(map[string]*Optimization)
	ao.state.objectives = make(map[string]*OptimizationObjective)
	ao.state.constraints = make([]OptimizationConstraint, 0)
	ao.state.metrics = OptimizationMetrics{
		Convergence: make([]ConvergencePoint, 0),
	}

	return ao
}

// Optimize 执行优化
func (ao *AdaptiveOptimization) Optimize() error {
	ao.mu.Lock()
	defer ao.mu.Unlock()

	// 更新优化目标
	if err := ao.updateObjectives(); err != nil {
		return err
	}

	// 执行优化迭代
	if err := ao.runOptimization(); err != nil {
		return err
	}

	// 应用优化结果
	if err := ao.applyOptimization(); err != nil {
		return err
	}

	// 更新指标
	ao.updateMetrics()

	return nil
}

// RegisterObjective 注册优化目标
func (ao *AdaptiveOptimization) RegisterObjective(
	objective *OptimizationObjective) error {

	if objective == nil {
		return model.WrapError(nil, model.ErrCodeValidation, "nil objective")
	}

	ao.mu.Lock()
	defer ao.mu.Unlock()

	// 验证目标
	if err := ao.validateObjective(objective); err != nil {
		return err
	}

	// 存储目标
	ao.state.objectives[objective.ID] = objective

	return nil
}

// updateObjectives 更新优化目标
func (ao *AdaptiveOptimization) updateObjectives() error {
	// 获取当前系统状态
	state, err := ao.getCurrentState()
	if err != nil {
		return err
	}

	// 更新每个目标
	for _, objective := range ao.state.objectives {
		// 评估目标当前值
		currentValue := objective.Evaluator.Function(state)

		// 检查是否需要优化
		if ao.needsOptimization(objective, currentValue) {
			// 创建新的优化实例
			optimization := ao.createOptimization(objective, currentValue)
			ao.state.optimizations[optimization.ID] = optimization
		}
	}

	return nil
}

// getCurrentState 获取当前系统状态
func (ao *AdaptiveOptimization) getCurrentState() (*types.SystemState, error) {
	// 通过策略管理器获取系统状态
	modelState, err := ao.strategy.getCurrentState()
	if err != nil {
		return nil, err
	}

	// 将 model.SystemState 转换为 types.SystemState
	return types.FromModelSystemState(modelState), nil
}

// needsOptimization 检查目标是否需要优化
func (ao *AdaptiveOptimization) needsOptimization(
	objective *OptimizationObjective,
	currentValue float64) bool {

	// 检查目标值与当前值的差距
	if math.Abs(objective.Target-currentValue) < ao.config.improvementThreshold {
		return false
	}

	// 检查是否已有正在进行的优化
	for _, opt := range ao.state.optimizations {
		if opt.Target == objective.ID && opt.Progress.Status == "active" {
			return false
		}
	}

	return true
}

// createOptimization 创建新的优化实例
func (ao *AdaptiveOptimization) createOptimization(
	objective *OptimizationObjective,
	currentValue float64) *Optimization {

	return &Optimization{
		ID:         generateOptimizationID(),
		Type:       objective.Type,
		Target:     objective.ID,
		StartTime:  time.Now(),
		LastUpdate: time.Now(),
		Progress: OptimizationProgress{
			Iteration:    0,
			BestValue:    currentValue,
			CurrentValue: currentValue,
			Status:       "active",
		},
		Parameters: []OptimizationParameter{
			{
				Name:         "target",
				CurrentValue: currentValue,
				Range:        [2]float64{currentValue * 0.5, currentValue * 1.5},
				Step:         ao.config.improvementThreshold,
				Weight:       1.0,
			},
		},
	}
}

// runOptimization 执行优化迭代
func (ao *AdaptiveOptimization) runOptimization() error {
	for _, opt := range ao.state.optimizations {
		// 检查优化状态
		if opt.Progress.Status != "active" {
			continue
		}

		// 执行优化迭代
		for i := opt.Progress.Iteration; i < ao.config.maxIterations; i++ {
			// 生成新参数
			newParams := ao.generateNewParameters(opt)

			// 评估新参数
			value := ao.evaluateParameters(opt, newParams)

			// 更新最优解
			if value > opt.Progress.BestValue {
				ao.updateOptimization(opt, newParams, value)
			}

			// 检查收敛
			if ao.hasConverged(opt) {
				opt.Progress.Status = "converged"
				break
			}

			opt.Progress.Iteration = i + 1
		}

		// 检查完成条件
		if ao.isOptimizationComplete(opt) {
			ao.finalizeOptimization(opt)
		}
	}

	return nil
}

// generateNewParameters 生成新的优化参数
func (ao *AdaptiveOptimization) generateNewParameters(opt *Optimization) []OptimizationParameter {
	newParams := make([]OptimizationParameter, len(opt.Parameters))
	copy(newParams, opt.Parameters)

	// 随机调整参数
	for i := range newParams {
		param := &newParams[i]
		// 在参数范围内随机调整
		delta := (param.Range[1] - param.Range[0]) * (rand.Float64()*2 - 1) * param.Step
		if current, ok := param.CurrentValue.(float64); ok {
			// 确保新值在范围内
			newValue := math.Max(param.Range[0],
				math.Min(param.Range[1], current+delta))
			param.CurrentValue = newValue
		}
	}

	return newParams
}

// evaluateParameters 评估参数效果
func (ao *AdaptiveOptimization) evaluateParameters(opt *Optimization, params []OptimizationParameter) float64 {
	// 获取目标对象
	objective := ao.state.objectives[opt.Target]
	if objective == nil {
		return 0
	}

	// 构造评估输入
	input := make(map[string]interface{})
	for _, param := range params {
		input[param.Name] = param.CurrentValue
	}

	// 执行评估
	return objective.Evaluator.Function(input)
}

// updateOptimization 更新优化状态
func (ao *AdaptiveOptimization) updateOptimization(opt *Optimization, params []OptimizationParameter, value float64) {
	// 计算改进幅度
	improvement := (value - opt.Progress.BestValue) / math.Abs(opt.Progress.BestValue)

	// 更新最优解
	opt.Parameters = params
	opt.Progress.BestValue = value
	opt.Progress.CurrentValue = value
	opt.Progress.Improvement = improvement
	opt.LastUpdate = time.Now()

	// 记录收敛点
	ao.state.metrics.Convergence = append(ao.state.metrics.Convergence, ConvergencePoint{
		Iteration:   opt.Progress.Iteration,
		Value:       value,
		Improvement: improvement,
	})

	// 限制收敛历史长度
	if len(ao.state.metrics.Convergence) > maxConvergenceHistory {
		ao.state.metrics.Convergence = ao.state.metrics.Convergence[1:]
	}
}

// hasConverged 检查是否收敛
func (ao *AdaptiveOptimization) hasConverged(opt *Optimization) bool {
	if opt == nil {
		return false
	}

	// 检查迭代次数是否足够
	if opt.Progress.Iteration < 3 {
		return false
	}

	// 获取优化实例的最近几次改进记录
	recentImprovements := make([]float64, 0)
	for _, point := range ao.state.metrics.Convergence {
		if point.Iteration > opt.Progress.Iteration-3 {
			recentImprovements = append(recentImprovements, point.Improvement)
		}
		if len(recentImprovements) >= 3 {
			break
		}
	}

	// 至少需要3个数据点才能判断收敛
	if len(recentImprovements) < 3 {
		return false
	}

	// 检查最近的改进幅度
	for _, improvement := range recentImprovements {
		if math.Abs(improvement) > ao.config.convergenceRate {
			return false
		}
	}

	// 检查是否达到目标值
	if objective := ao.state.objectives[opt.Target]; objective != nil {
		currentError := math.Abs(opt.Progress.CurrentValue - objective.Target)
		if currentError > ao.config.improvementThreshold {
			return false
		}
	}

	return true
}

// isOptimizationComplete 检查优化是否完成
func (ao *AdaptiveOptimization) isOptimizationComplete(opt *Optimization) bool {
	// 检查状态
	if opt.Progress.Status == "converged" {
		return true
	}

	// 检查迭代次数
	if opt.Progress.Iteration >= ao.config.maxIterations {
		return true
	}

	// 检查目标值是否达到
	if objective := ao.state.objectives[opt.Target]; objective != nil {
		if math.Abs(opt.Progress.BestValue-objective.Target) < ao.config.improvementThreshold {
			return true
		}
	}

	return false
}

// finalizeOptimization 完成优化处理
func (ao *AdaptiveOptimization) finalizeOptimization(opt *Optimization) {
	// 设置结果
	opt.Result = OptimizationResult{
		Success:    opt.Progress.Status == "converged",
		FinalValue: opt.Progress.BestValue,
		Parameters: make(map[string]interface{}),
		Performance: map[string]float64{
			"iterations":  float64(opt.Progress.Iteration),
			"improvement": opt.Progress.Improvement,
		},
		Duration: time.Since(opt.StartTime),
	}

	// 复制最优参数
	for _, param := range opt.Parameters {
		opt.Result.Parameters[param.Name] = param.CurrentValue
	}

	// 更新状态
	opt.Progress.Status = "completed"
}

// applyOptimization 应用优化结果
func (ao *AdaptiveOptimization) applyOptimization() error {
	for _, opt := range ao.state.optimizations {
		if opt.Result.Success {
			// 应用优化参数
			if err := ao.applyOptimizedParameters(opt); err != nil {
				continue
			}

			// 更新学习系统
			ao.updateLearningSystem(opt)
		}
	}

	return nil
}

// applyOptimizedParameters 应用优化后的参数
func (ao *AdaptiveOptimization) applyOptimizedParameters(opt *Optimization) error {
	// 获取目标对象
	objective := ao.state.objectives[opt.Target]
	if objective == nil {
		return fmt.Errorf("objective %s not found", opt.Target)
	}

	// 构造应用参数
	params := make(map[string]interface{})
	for _, param := range opt.Parameters {
		params[param.Name] = param.CurrentValue
	}

	// 通过策略应用参数
	switch objective.Type {
	case "system":
		// 系统级参数调整
		return ao.strategy.UpdateParameters(objective.TargetID, params) // 使用TargetID
	case "component":
		// 组件级参数调整
		return ao.strategy.mutationHandler.AdjustParameter(objective.TargetID, params) // 使用TargetID
	default:
		return fmt.Errorf("unknown optimization type: %s", objective.Type)
	}
}

// updateLearningSystem 更新学习系统
func (ao *AdaptiveOptimization) updateLearningSystem(opt *Optimization) {
	// 创建优化经验
	experience := LearningExperience{
		ID:        fmt.Sprintf("opt_%s", opt.ID),
		Type:      "optimization",
		Timestamp: time.Now(),
		Context: map[string]interface{}{
			"objective_type":   opt.Type,
			"objective_target": opt.Target,
			"parameters":       opt.Result.Parameters,
		},
		Result: LearningResult{
			Status: "success",
			Metrics: map[string]float64{
				"improvement": opt.Progress.Improvement,
				"iterations":  float64(opt.Progress.Iteration),
				"final_value": opt.Result.FinalValue,
			},
			Duration: opt.Result.Duration,
		},
	}

	// 添加到学习系统
	ao.learning.addExperience(experience)
}

// 辅助函数

func (ao *AdaptiveOptimization) validateObjective(
	objective *OptimizationObjective) error {

	if objective.ID == "" {
		return model.WrapError(nil, model.ErrCodeValidation, "empty objective ID")
	}

	if objective.Evaluator.Function == nil {
		return model.WrapError(nil, model.ErrCodeValidation, "missing evaluator function")
	}

	return nil
}

func (ao *AdaptiveOptimization) updateMetrics() {
	total := len(ao.state.optimizations)
	success := 0
	totalImprovement := 0.0

	for _, opt := range ao.state.optimizations {
		if opt.Result.Success {
			success++
			totalImprovement += opt.Progress.Improvement
		}
	}

	// 更新统计
	ao.state.metrics.TotalOptimizations = total
	ao.state.metrics.SuccessRate = float64(success) / float64(total)
	ao.state.metrics.AverageImprovement = totalImprovement / float64(success)
}

func generateOptimizationID() string {
	return fmt.Sprintf("opt_%d", time.Now().UnixNano())
}

// RunOptimization 执行指定优化目标的优化
func (ao *AdaptiveOptimization) RunOptimization(objectives []*OptimizationObjective, constraints map[string]types.Constraint) error {
	// 锁定状态
	ao.mu.Lock()
	defer ao.mu.Unlock()

	// 注册优化目标
	for _, objective := range objectives {
		if err := ao.RegisterObjective(objective); err != nil {
			return err
		}
	}

	// 注册约束
	for name, constraint := range constraints {
		ao.state.constraints = append(ao.state.constraints, OptimizationConstraint{
			Type:      "system",
			Target:    name,
			Condition: "range",
			Value: map[string]float64{
				"min": constraint.Min,
				"max": constraint.Max,
			},
			Priority: 1,
		})
	}

	// 执行优化
	if err := ao.runOptimization(); err != nil {
		return err
	}

	// 应用优化结果
	if err := ao.applyOptimization(); err != nil {
		return err
	}

	// 更新指标
	ao.updateMetrics()

	return nil
}
