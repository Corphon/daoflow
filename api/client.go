// api/client.go

package api

import (
	"context"
	"fmt"
	"sync"

	"github.com/Corphon/daoflow/core"
	"github.com/Corphon/daoflow/model"
	"github.com/Corphon/daoflow/system"
	"github.com/Corphon/daoflow/system/types"
)

// Client 实现DaoFlow接口
type Client struct {
	mu        sync.RWMutex
	sys       *system.System
	isRunning bool
}

var _ DaoFlow = (*Client)(nil)

// NewClient 创建客户端-------------------------------
func NewClient(cfg *system.Config) (*Client, error) {
	// 创建并初始化DaoFlow客户端实例，该实例是与DaoFlow系统交互的主要接口。
	// 客户端封装了系统的所有功能，提供易用的API访问道势系统的核心能力。
	//
	// 参数:
	//   - cfg *system.Config: 系统配置,可包含以下设置:
	//     - CoreConfig: 核心引擎配置(如量子态、场、能量系统参数)
	//     - ModelConfig: 模型配置(阴阳流、五行、八卦等参数)
	//     - EvolutionConfig: 演化系统配置
	//     - ControlConfig: 控制系统配置
	//     - MetadataConfig: 元数据系统配置
	//     - LogLevel: 日志级别
	//   - 如果传入nil,将使用系统默认配置
	//
	// 返回值:
	//   - *Client: DaoFlow客户端实例
	//   - error: 创建失败时返回错误
	//
	// 示例:
	//   // 使用默认配置创建客户端
	//   client, err := api.NewClient(nil)
	//   if err != nil {
	//       log.Fatalf("创建客户端失败: %v", err)
	//   }
	//
	//   // 使用自定义配置
	//   cfg := &system.Config{
	//       CoreConfig: &core.Config{
	//           InitialEnergy: 0.8,
	//           EnergyDecayRate: 0.01,
	//       },
	//       LogLevel: "info",
	//   }
	//   client, err := api.NewClient(cfg)
	if cfg == nil {
		cfg = system.DefaultConfig()
	}

	sys, err := system.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create system: %w", err)
	}

	return &Client{
		sys: sys,
	}, nil
}

// Initialize 初始化系统
func (c *Client) Initialize(ctx context.Context) error {
	// 在Start方法调用前必须先调用此方法完成系统初始化，包括:
	// - 初始化各个子系统(控制层、演化层、元数据层、监控层)
	// - 建立组件间依赖关系
	// - 准备事件处理系统
	// - 加载配置参数
	//
	// 参数:
	//   - ctx context.Context: 上下文对象，可用于:
	//     - 设置超时限制: context.WithTimeout(parentCtx, 30*time.Second)
	//     - 传递取消信号: context.WithCancel(parentCtx)
	//     - 跨组件传递数据
	//
	// 返回值:
	//   - error: 初始化过程中的错误
	//
	// 示例:
	//   // 创建带30秒超时的上下文
	//   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//   defer cancel()
	//
	//   // 初始化系统
	//   if err := client.Initialize(ctx); err != nil {
	//       log.Fatalf("系统初始化失败: %v", err)
	//   }
	//
	//   // 初始化后启动系统
	//   if err := client.Start(); err != nil {
	//       log.Fatalf("系统启动失败: %v", err)
	//   }
	return c.sys.Initialize(ctx)
}

// Start 启动系统
func (c *Client) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 调用此方法启动DaoFlow系统及其所有子系统，包括：
	// - 核心引擎
	// - 控制层
	// - 演化层
	// - 元数据层
	// - 监控层
	// - 已注册的所有模型(如阴阳流、五行流等)
	//
	// 必须先调用 Initialize方法 完成初始化后才能启动系统。
	// 启动成功后，系统开始处理事件、执行模型转换、监控状态等。
	//
	// 返回值:
	//   - error: 启动失败时返回错误
	//     - types.ErrAlreadyRunning: 如果系统已经在运行
	//
	// 示例:
	//   // 初始化系统
	//   ctx := context.Background()
	//   if err := client.Initialize(ctx); err != nil {
	//       log.Printf("初始化失败: %v", err)
	//       return err
	//   }
	//
	//   // 启动系统
	//   if err := client.Start(); err != nil {
	//       log.Printf("启动失败: %v", err)
	//       return err
	//   }
	//
	//   fmt.Println("系统已成功启动")
	if c.isRunning {
		return types.ErrAlreadyRunning
	}

	if err := c.sys.Start(); err != nil {
		return err
	}

	c.isRunning = true
	return nil
}

// Stop 停止系统
func (c *Client) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 停止DaoFlow系统及其所有子系统的运行，包括:
	// - 核心引擎
	// - 控制层
	// - 演化层
	// - 元数据层
	// - 监控层
	// - 所有运行中的模型
	//
	// 停止后系统将释放资源并进入停止状态。已注册的模型会保留，但会停止运行。
	// 停止后可以再次通过 Start方法 启动系统。
	//
	// 返回值:
	//   - error: 停止过程中的任何错误
	//     - nil: 停止成功或系统本来就未运行
	//
	// 示例:
	//   // 停止系统
	//   if err := client.Stop(); err != nil {
	//       log.Printf("停止系统失败: %v", err)
	//       return err
	//   }
	//
	//   fmt.Println("系统已停止")
	//
	//   // 系统可以再次启动
	//   if err := client.Start(); err != nil {
	//       log.Printf("重启系统失败: %v", err)
	//   }
	if !c.isRunning {
		return nil
	}

	if err := c.sys.Stop(); err != nil {
		return err
	}

	c.isRunning = false
	return nil
}

// Close 关闭客户端
func (c *Client) Close() error {
	// 关闭客户端并释放所有相关资源。此方法会停止系统的所有组件，
	// 但会保留所有已注册的模型，以便后续可以重新启动。
	// 实际上此方法是调用Stop()方法的便捷方式，用于实现io.Closer接口。
	//
	// 返回值:
	//   - error: 关闭过程中发生的错误
	//     - nil: 关闭成功或系统已经处于停止状态
	//
	// 示例:
	//   // 创建客户端
	//   client, err := api.NewClient(nil)
	//   if err != nil {
	//       log.Fatalf("创建客户端失败: %v", err)
	//   }
	//
	//   // 使用defer关闭客户端
	//   defer client.Close()
	//
	//   // 或在完成使用后显式关闭
	//   if err := client.Close(); err != nil {
	//       log.Printf("关闭客户端失败: %v", err)
	//   }
	return c.Stop()
}

// Shutdown 关闭系统
func (c *Client) Shutdown(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 关闭整个系统并释放所有资源，包括停止所有子系统、组件和服务。
	// 此方法执行有序关闭，确保数据完整性和资源正确释放。
	// 会等待所有组件完成关闭操作或超时(默认30秒)。
	//
	// 参数:
	//   - ctx context.Context: 上下文对象，用于控制关闭操作的超时和取消
	//     - 使用context.WithTimeout可设置自定义超时
	//     - 使用context.WithCancel可手动取消关闭过程
	//
	// 返回值:
	//   - error: 关闭过程中发生的错误
	//     - nil: 成功关闭所有系统组件
	//     - "system shutdown timed out": 关闭超时
	//
	// 示例:
	//   // 使用默认超时关闭
	//   err := client.Shutdown(context.Background())
	//
	//   // 指定10秒超时关闭
	//   ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//   defer cancel()
	//   err := client.Shutdown(ctx)
	//   if err != nil {
	//       log.Printf("系统关闭失败: %v", err)
	//   }
	// 重置运行状态
	c.isRunning = false

	// 委托给system处理
	return c.sys.Shutdown(ctx)
}

//------------------------------------------

// GetSystem 获取system实例
func (c *Client) GetSystem() *system.System {
	// 此方法返回客户端所管理的系统实例，允许直接访问系统层的功能。
	// 注意：除非有特殊需求，应优先使用Client提供的高级API而不是直接操作系统实例。
	//
	// 返回值:
	//   - *system.System: 系统实例，提供以下功能:
	//     1. 底层组件访问
	//     2. 系统核心操作
	//     3. 子系统管理
	//     4. 事件处理
	//     5. 资源协调
	//
	// 使用示例:
	//   // 获取系统实例
	//   sys := client.GetSystem()
	//
	//   // 访问系统组件
	//   status := sys.GetStatus()
	//   metrics := sys.GetMetrics()
	//
	//   // 检查系统健康状态
	//   health := sys.calculateSystemHealth()
	//   fmt.Printf("系统健康度: %.2f\n", health)
	return c.sys
}

// GetSystemStatus 获取系统当前运行状态
func (c *Client) GetSystemStatus() string {
	// 返回值:
	//   - string: 系统状态，可能的值包括:
	//     - "running": 系统正在运行
	//     - "stopped": 系统已停止
	//     - "initialized": 系统已初始化但未运行
	//     - "error": 系统处于错误状态
	//     - "restarting": 系统正在重启
	//
	// 使用示例:
	//   status := client.GetSystemStatus()
	//   if status == "running" {
	//       fmt.Println("系统正常运行中")
	//   } else if status == "error" {
	//       fmt.Println("系统出现错误，需要检查")
	//   } else {
	//       fmt.Printf("系统当前状态: %s\n", status)
	//   }
	return c.sys.GetStatus()
}

// GetSystemMetrics 获取系统指标
func (c *Client) GetSystemMetrics() map[string]interface{} {
	// 该方法返回当前系统的所有指标数据，包括系统状态、健康度、能量、资源利用率等关键信息。
	// 返回格式为map[string]interface{}，方便直接使用或序列化。
	//
	// 返回值:
	//   - map[string]interface{}: 系统指标映射，包含以下主要类别:
	//     - "status": 系统运行状态(如"running", "error"等)
	//     - "health": 系统健康度(0-1范围)
	//     - "alert_count": 当前告警数量
	//     - "system": 系统核心指标(包含能量、场、量子态等)
	//     - "performance": 性能指标(如QPS、延迟、错误率等)
	//     - "resources": 资源利用率指标(如CPU、内存使用率)
	//     - "subsystems": 各子系统指标(控制层、演化层、元数据层等)
	//
	// 使用示例:
	//   metrics := client.GetSystemMetrics()
	//
	//   // 获取系统状态和健康度
	//   fmt.Printf("系统状态: %s, 健康度: %.2f\n",
	//       metrics["status"], metrics["health"])
	//
	//   // 获取系统能量
	//   systemMetrics := metrics["system"].(map[string]interface{})
	//   energy := systemMetrics["energy"].(float64)
	//   fmt.Printf("系统能量: %.2f\n", energy)
	//
	//   // 检查性能指标
	//   perfMetrics := metrics["performance"].(map[string]interface{})
	//   qps := perfMetrics["qps"].(float64)
	//   latency := perfMetrics["latency"].(string)
	//   fmt.Printf("QPS: %.2f, 延迟: %s\n", qps, latency)
	metrics := c.sys.GetMetrics()
	return metrics.ToMap()
}

// RegisterModel 注册模型
func (c *Client) RegisterModel(name string, m model.Model) error {
	// 将实现了model.Model接口的模型实例注册到系统中，使其受到系统管理并能参与系统运行。
	// 注册后的模型可以通过名称进行引用，并会根据系统状态自动启动或停止。
	// 如果系统已经运行，新注册的模型会自动启动。
	//
	// 参数:
	//   - name: 模型唯一标识符，用于后续获取或操作此模型
	//   - m: 实现了model.Model接口的模型实例
	//
	// 返回值:
	//   - error: 注册失败时返回错误，可能的错误包括:
	//     - 同名模型已存在
	//     - 模型初始化失败
	//     - 系统已运行但模型启动失败
	//
	// 示例:
	//   // 创建并注册阴阳流模型
	//   yinYangFlow := model.NewYinYangFlow()
	//   if err := client.RegisterModel("yinyang", yinYangFlow); err != nil {
	//       log.Printf("模型注册失败: %v", err)
	//       return err
	//   }
	//
	//   // 创建并注册五行流模型
	//   wuXingFlow := model.NewWuXingFlow()
	//   if err := client.RegisterModel("wuxing", wuXingFlow); err != nil {
	//       log.Printf("模型注册失败: %v", err)
	//       return err
	//   }
	return c.sys.RegisterModel(name, m)
}

// GetModel 获取已注册的模型
func (c *Client) GetModel(name string) (model.Model, error) {
	// 根据名称获取系统中已注册的模型实例，可用于查询模型状态或执行模型特定操作。
	//
	// 参数:
	//   - name: 模型名称，必须是通过RegisterModel注册的唯一标识符
	//
	// 返回值:
	//   - model.Model: 模型实例，提供模型特定功能
	//   - error: 获取模型时遇到的错误，常见错误包括:
	//     - types.ErrModelNotFound: 指定名称的模型未注册
	//     - types.ErrSystem: 系统内部错误
	//
	// 示例:
	//   // 获取阴阳流模型
	//   yinyang, err := client.GetModel("yinyang")
	//   if err != nil {
	//       log.Printf("获取模型失败: %v", err)
	//       return err
	//   }
	//
	//   // 使用获取的模型
	//   state := yinyang.GetState()
	//   fmt.Printf("模型能量: %.2f\n", state.Energy)
	//
	//   // 执行模型转换
	//   err = yinyang.Transform(model.PatternBalance)
	return c.sys.GetModel(name)
}

// ListModels 获取系统中已注册的所有模型列表
func (c *Client) ListModels() []string {
	// 该方法返回当前系统中注册的所有模型的名称列表。
	// 这些模型名称可用于后续通过GetModel方法获取具体模型实例。
	//
	// 返回值:
	//   - []string: 模型名称列表，包含所有通过RegisterModel注册的模型名称
	//
	// 使用示例:
	//   // 获取所有已注册模型
	//   modelNames := client.ListModels()
	//
	//   // 遍历模型列表
	//   for _, name := range modelNames {
	//       fmt.Printf("发现模型: %s\n", name)
	//
	//       // 获取特定模型
	//       if name == "yinyang" {
	//           model, err := client.GetModel(name)
	//           if err == nil {
	//               state := model.GetState()
	//               fmt.Printf("阴阳模型能量: %.2f\n", state.Energy)
	//           }
	//       }
	//   }
	return c.sys.ListModels()
}

// UnregisterModel 注销模型
func (c *Client) UnregisterModel(name string) error {
	// 将指定名称的模型从系统中注销并释放相关资源。
	// 如果模型正在运行，会先停止模型然后再移除。
	// 注销后的模型将不再接收系统事件，也不会参与系统运行。
	//
	// 参数:
	//   - name: 模型名称，必须是通过RegisterModel注册的唯一标识符
	//
	// 返回值:
	//   - error: 注销失败时返回错误，可能的错误包括:
	//     - "model not found": 指定名称的模型未注册
	//     - "failed to stop model": 停止模型时发生错误
	//
	// 示例:
	//   // 注销阴阳流模型
	//   if err := client.UnregisterModel("yinyang"); err != nil {
	//       log.Printf("模型注销失败: %v", err)
	//       return err
	//   }
	//
	//   // 检查模型是否已注销
	//   models := client.ListModels()
	//   for _, m := range models {
	//       if m == "yinyang" {
	//           log.Printf("模型注销失败，模型仍存在")
	//       }
	//   }
	return c.sys.UnregisterModel(name)
}

// ModelAPI实现
// TransformModel 执行模型转换操作
func (c *Client) TransformModel(ctx context.Context, pattern model.TransformPattern) error {
	// 将指定的转换模式应用于系统中已注册的所有模型，触发相应的状态转换。
	// 此方法是高级模型操作API，可用于同时转换多个模型或执行系统级转换。
	//
	// 参数:
	//   - ctx: 上下文对象，用于控制操作超时或取消
	//   - pattern: 转换模式，可选值:
	//     - model.PatternBalance: 平衡模式，使系统达到平衡状态
	//     - model.PatternForward: 前进模式，增加能量和活性
	//     - model.PatternReverse: 逆向模式，减少能量和活性
	//     - model.PatternNormal: 常规模式，执行标准转换
	//
	// 返回值:
	//   - error: 转换执行过程中的错误
	//
	// 示例:
	//   // 执行平衡转换，使所有模型达到和谐状态
	//   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//   defer cancel()
	//   err := client.TransformModel(ctx, model.PatternBalance)
	//
	//   // 执行前进转换，促进系统演化
	//   err = client.TransformModel(context.Background(), model.PatternForward)
	return c.sys.TransformModel(ctx, pattern)
}

// GetEnergy 返回系统当前的总能量值。
func (c *Client) GetEnergy() float64 {
	// 能量值范围为0-1之间的归一化值,表示系统的整体能量水平。
	// - 0表示系统能量耗尽
	// - 1表示系统能量充盈
	// 该方法用于:
	// 1. 监控系统能量状态
	// 2. 判断系统是否需要补充能量
	// 3. 评估系统运行状态
	// 4. 作为其他模块的决策依据
	return c.sys.GetEnergy()
}

// AdjustEnergy 调整系统能量水平
func (c *Client) AdjustEnergy(delta float64) error {
	// 参数:
	//   - delta: 能量调整量,范围[-1.0, 1.0]
	//     正值表示增加能量
	//     负值表示减少能量
	//
	// 返回值:
	//   - error: 调整失败时返回错误
	//     - ErrInvalidParameter: delta超出[-1.0, 1.0]范围
	//     - ErrEnergyOutOfRange: 调整后能量超出[0.0, 1.0]范围
	//
	// 示例:
	//   调整能量:
	//     client.AdjustEnergy(0.1)  // 增加10%能量
	//     client.AdjustEnergy(-0.2) // 减少20%能量
	return c.sys.AdjustEnergy(delta)
}

// GetEnergySystem 获取系统的能量系统实例
func (c *Client) GetEnergySystem() *core.EnergySystem {
	// 返回值:
	//   - *core.EnergySystem: 能量系统实例,提供以下功能:
	//     1. 获取系统总能量: GetTotalEnergy()
	//     2. 调整能量水平: TransformEnergy()
	//     3. 能量形态转换: Convert()
	//     4. 获取能量状态: GetEnergyState()
	//     5. 获取系统平衡度: GetBalance()
	//
	// 使用示例:
	//   energySystem := client.GetEnergySystem()
	//   totalEnergy := energySystem.GetTotalEnergy()
	//   balance := energySystem.GetBalance()
	return c.sys.GetEnergySystem()
}

// PatternAPI实现
// DetectPattern 检测输入数据中的模式
func (c *Client) DetectPattern(data interface{}) (*model.FlowPattern, error) {
	// 参数:
	//   - data: 待分析数据,支持以下类型:
	//     1. map[string]interface{}: 结构化数据
	//     2. []float64: 时序数据
	//     3. string: JSON格式数据
	//
	// 返回值:
	//   - *model.FlowPattern: 检测到的流模式
	//   - error: 检测过程中的错误
	//
	// 示例:
	//   data := map[string]interface{}{
	//       "energy": 0.8,
	//       "phase": 0.5,
	//       "timestamp": time.Now(),
	//   }
	//   pattern, err := client.DetectPattern(data)
	return c.sys.Evolution().DetectPattern(data)
}

// AnalyzePattern 分析已检测到的模式
func (c *Client) AnalyzePattern(pattern *model.FlowPattern) error {
	// 参数:
	//   - pattern: 待分析的流模式，包含以下信息:
	//     1. 模式类型 (Type): 如"energy_wave"、"phase_shift"等
	//     2. 模式强度 (Strength): 0-1之间的归一化值
	//     3. 模式特征 (Features): 关键特征向量
	//     4. 时间戳 (Timestamp): 模式检测时间
	//
	// 返回值:
	//   - error: 分析过程中的错误
	//     - ErrInvalidPattern: 无效的模式
	//     - ErrAnalysisFailed: 分析失败
	//
	// 示例:
	//   // 1. 检测并分析模式
	//   data := getSystemState()
	//   pattern, _ := client.DetectPattern(data)
	//   if err := client.AnalyzePattern(pattern); err != nil {
	//       log.Printf("模式分析失败: %v", err)
	//       return err
	//   }
	//   // 2. 使用分析结果
	//   fmt.Printf("模式类型: %s\n", pattern.Type)
	//   fmt.Printf("模式强度: %.2f\n", pattern.Strength)
	//   fmt.Printf("特征向量: %v\n", pattern.Features)
	return c.sys.Evolution().AnalyzePattern(pattern)
}

// Subscribe 订阅系统事件
func (c *Client) Subscribe(eventType types.EventType, handler types.EventHandler) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 参数:
	//   - eventType: 事件类型,可选值包括:
	//     - types.EventStateChange: 状态变更事件
	//     - types.EventSystemStarted: 系统启动事件
	//     - types.EventSystemStopped: 系统停止事件
	//     - types.EventModelChange: 模型变更事件
	//   - handler: 事件处理器,需实现types.EventHandler接口
	//
	// 返回值:
	//   - error: 订阅失败时返回错误
	//
	// 示例:
	//   handler := NewEventHandler("myHandler", func(event SystemEvent) error {
	//       log.Printf("收到事件: %v", event)
	//       return nil
	//   })
	//   err := client.Subscribe(types.EventStateChange, handler)
	return c.sys.Subscribe(eventType, handler)
}

// Unsubscribe 取消订阅系统事件
func (c *Client) Unsubscribe(eventType types.EventType, handler types.EventHandler) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 参数:
	//   - eventType: 事件类型,可选值包括:
	//     - types.EventStateChange: 状态变更事件
	//     - types.EventSystemStarted: 系统启动事件
	//     - types.EventSystemStopped: 系统停止事件
	//     - types.EventModelChange: 模型变更事件
	//   - handler: 要取消订阅的事件处理器
	//
	// 返回值:
	//   - error: 取消订阅失败时返回错误
	//     - ErrInvalidParameter: 无效的处理器
	//     - ErrNotFound: 未找到对应的订阅
	//
	// 示例:
	//   // 取消订阅事件
	//   handler := GetEventHandler("myHandler")
	//   err := client.Unsubscribe(types.EventStateChange, handler)
	//   if err != nil {
	//       log.Printf("取消订阅失败: %v", err)
	//       return err
	//   }
	return c.sys.Unsubscribe(eventType, handler)
}

// PublishEvent 发布系统事件
func (c *Client) PublishEvent(event types.Event) error {
	// 参数:
	//   - event: 要发布的事件,包含以下信息:
	//     - Type: 事件类型(必填)
	//     - Payload: 事件负载数据
	//     - Metadata: 事件元数据
	//
	// 返回值:
	//   - error: 发布失败时返回错误
	//     - ErrState: 系统未运行
	//     - ErrQueue: 事件队列已满
	//
	// 示例:
	//   event := types.Event{
	//       Type:      types.EventStateChange,
	//       Timestamp: time.Now(),
	//       Source:    "client",
	//       Payload:   map[string]interface{}{"energy": 0.8},
	//       Metadata:  map[string]interface{}{"priority": "high"},
	//   }
	//   if err := client.PublishEvent(event); err != nil {
	//       log.Printf("发布事件失败: %v", err)
	//   }
	return c.sys.PublishEvent(event)
}

// DaoAPI实现
// GetYinYangFlow 获取系统的阴阳流模型实例
func (c *Client) GetYinYangFlow() *model.YinYangFlow {
	// 返回值:
	//   - *model.YinYangFlow: 阴阳流模型,提供以下功能:
	//     1. 阴阳能量转换: Transform()
	//     2. 阴阳平衡调节: balanceTransform()
	//     3. 获取阴阳状态: GetState()
	//     4. 调整阴阳能量: AdjustEnergy()
	//
	// 使用示例:
	//   yinyang := client.GetYinYangFlow()
	//   state := yinyang.GetState()
	//   fmt.Printf("阴能量: %.2f, 阳能量: %.2f, 平衡度: %.2f\n",
	//              state.YinEnergy, state.YangEnergy, state.Balance)
	//
	//   // 执行阴阳平衡
	//   yinyang.Transform(model.PatternBalance)
	return c.sys.GetYinYangFlow()
}

// TransformYinYang 执行阴阳模型转换
func (c *Client) TransformYinYang(pattern model.TransformPattern) error {
	// 参数:
	//   - pattern: 转换模式,可选值:
	//     - model.PatternBalance: 阴阳平衡模式,使系统能量均衡分配
	//     - model.PatternForward: 阴转阳模式,增加阳能量
	//     - model.PatternReverse: 阳转阴模式,增加阴能量
	//     - model.PatternNatural: 自然模式,根据当前状态自动选择转换方向
	//
	// 返回值:
	//   - error: 转换失败时返回错误
	//
	// 示例:
	//   // 执行阴阳平衡
	//   err := client.TransformYinYang(model.PatternBalance)
	//
	//   // 执行阴转阳
	//   err := client.TransformYinYang(model.PatternForward)
	return c.sys.TransformYinYang(pattern)
}

// GetBaGuaFlow 获取系统的八卦流模型实例
func (c *Client) GetBaGuaFlow() *model.BaGuaFlow {
	// 返回值:
	//   - *model.BaGuaFlow: 八卦流模型,提供以下功能:
	//     1. 八卦转换: Transform()
	//     2. 卦象变化: naturalTransform(), resonantTransform()
	//     3. 获取八卦状态: GetState()
	//     4. 调整八卦能量: AdjustEnergy()
	//
	// 使用示例:
	//   bagua := client.GetBaGuaFlow()
	//   state := bagua.GetState()
	//
	//   // 执行八卦平衡
	//   bagua.Transform(model.PatternBalance)
	//
	//   // 查看卦象状态
	//   harmony := state.Properties["harmony"]
	//   fmt.Printf("八卦和谐度: %.2f\n", harmony)
	return c.sys.GetBaGuaFlow()
}

// GetFieldSystem 获取系统的场系统实例
func (c *Client) GetFieldSystem() *core.FieldSystem {
	// 返回值:
	//   - *core.FieldSystem: 场系统实例,提供以下功能:
	//     1. 场演化: Evolve()
	//     2. 场叠加: Superposition()
	//     3. 场干涉: CalculateInterference()
	//     4. 场梯度计算: CalculateGradient()
	//     5. 获取场状态: GetFieldState()
	//
	// 使用示例:
	//   fieldSystem := client.GetFieldSystem()
	//
	//   // 获取场强度
	//   strength := fieldSystem.GetStrength()
	//
	//   // 计算场干涉
	//   interference := fieldSystem.CalculateInterference(pos)
	//
	//   // 演化场状态
	//   fieldSystem.Evolve()
	return c.sys.GetFieldSystem()
}

// GetGanZhiFlow 获取系统的干支流模型实例
func (c *Client) GetGanZhiFlow() *model.GanZhiFlow {
	// 返回值:
	//   - *model.GanZhiFlow: 干支流模型,提供以下功能:
	//     1. 干支转换: Transform()
	//     2. 天干地支周期变换: cyclicTransform()
	//     3. 获取干支状态: GetState()
	//     4. 调整干支能量: AdjustEnergy()
	//
	// 使用示例:
	//   ganzhi := client.GetGanZhiFlow()
	//   state := ganzhi.GetState()
	//
	//   // 执行干支周期转换
	//   ganzhi.Transform(model.PatternForward)
	//
	//   // 获取当前干支组合
	//   stem, branch := ganzhi.getCurrentGanZhi()
	//   fmt.Printf("当前干支: %v %v\n", stem, branch)
	return c.sys.GetGanZhiFlow()
}

// GetMetrics 获取模型指标
func (c *Client) GetMetrics() model.ModelMetrics {
	// 返回值:
	//   - model.ModelMetrics: 模型指标，包含以下内容:
	//     1. Energy: 能量指标 (总量、平均值、方差)
	//     2. State: 状态指标 (稳定性、转换次数、运行时间)
	//     3. Performance: 性能指标 (吞吐量、QPS、错误率)
	//     4. Quantum: 量子态指标
	//     5. Field: 场态指标
	//
	// 使用示例:
	//   metrics := client.GetMetrics()
	//
	//   // 查看能量指标
	//   fmt.Printf("总能量: %.2f, 平均能量: %.2f\n",
	//              metrics.Energy.Total, metrics.Energy.Average)
	//
	//   // 查看场强度
	//   fieldStrength := metrics.Field.GetStrength()
	return c.sys.GetModelMetrics()
}

// GetModelState 获取当前系统的模型状态
func (c *Client) GetModelState() model.ModelState {
	// 返回值:
	//   - model.ModelState: 当前模型状态,包含以下信息:
	//     1. 能量 (Energy): 当前能量水平
	//     2. 相位 (Phase): 当前系统相位
	//     3. 性质 (Nature): 当前系统属性
	//     4. 健康度 (Health): 系统健康状态
	//     5. 各模型特定属性 (Properties)
	//
	// 使用示例:
	//   state := client.GetModelState()
	//
	//   // 查看当前能量水平
	//   fmt.Printf("当前能量: %.2f\n", state.Energy)
	//
	//   // 检查系统健康度
	//   if state.Health < 0.5 {
	//       fmt.Println("系统状态不佳,需要维护")
	//   }
	//
	//   // 获取特定属性
	//   stability := state.Properties["stability"].(float64)
	return c.sys.GetModelState()
}

// GetQuantumSystem 获取系统的量子系统实例
func (c *Client) GetQuantumSystem() *core.QuantumSystem {
	// 返回值:
	//   - *core.QuantumSystem: 量子系统实例，提供以下功能:
	//     1. 获取量子态集合: GetStates()
	//     2. 获取相干性: GetCoherence()
	//     3. 获取纠缠度: GetEntanglement()
	//     4. 执行量子演化: Evolve(pattern)
	//     5. 量子测量: Measure()
	//
	// 使用示例:
	//   quantum := client.GetQuantumSystem()
	//
	//   // 获取量子相干性
	//   coherence := quantum.GetCoherence()
	//
	//   // 获取所有量子态
	//   states := quantum.GetStates()
	//   for i, state := range states {
	//       fmt.Printf("状态 %d: 相位=%.2f, 能量=%.2f\n",
	//                  i, state.GetPhase(), state.GetEnergy())
	//   }
	return c.sys.GetQuantumSystem()
}

// GetState 获取当前系统状态
func (c *Client) GetState() model.SystemState {
	// 返回值:
	//   - model.SystemState: 当前系统状态，包含以下信息:
	//     1. Energy: 系统能量
	//     2. Entropy: 系统熵
	//     3. Harmony: 系统和谐度
	//     4. Balance: 系统平衡度
	//     5. Phase: 当前相位
	//     6. Properties: 系统属性集合
	//     7. Timestamp: 状态时间戳
	//
	// 使用示例:
	//   state := client.GetState()
	//
	//   // 检查系统能量
	//   fmt.Printf("系统能量: %.2f\n", state.Energy)
	//
	//   // 查看系统平衡度
	//   if state.Balance < 0.3 {
	//       fmt.Println("系统失衡，需要调整")
	//   }
	//
	//   // 获取自定义属性
	//   if cycleType, ok := state.Properties["cycle_type"].(string); ok {
	//       fmt.Printf("当前周期类型: %s\n", cycleType)
	//   }
	return c.sys.GetState()
}

// GetWuXingFlow 获取系统的五行流模型实例
func (c *Client) GetWuXingFlow() *model.WuXingFlow {
	// 返回值:
	//   - *model.WuXingFlow: 五行流模型,提供以下功能:
	//     1. 五行转换: Transform()
	//     2. 五行相生相克: generateTransform(), constrainTransform()
	//     3. 获取五行状态: GetState()
	//     4. 调整五行能量: AdjustEnergy()
	//
	// 使用示例:
	//   wuxing := client.GetWuXingFlow()
	//   state := wuxing.GetState()
	//
	//   // 执行五行相生转换
	//   wuxing.Transform(model.PatternForward)
	//
	//   // 获取木元素能量
	//   woodEnergy := wuxing.GetWuXingElementEnergy("Wood")
	//   fmt.Printf("木元素能量: %.2f\n", woodEnergy)
	return c.sys.GetWuXingFlow()
}

// HandleEvent 处理系统事件
func (c *Client) HandleEvent(event types.SystemEvent) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 该方法将事件提交到系统的事件处理系统，用于组件间通信、状态变更通知等场景。
	// 系统会根据事件类型将其分发给已注册的相应处理器。
	//
	// 参数:
	//   - event: 系统事件对象，需包含以下信息:
	//     - Type: 事件类型，如types.EventStateChange、types.EventSystemStarted等
	//     - Data: 事件数据(可选)
	//     - Source: 事件源(可选)
	//     - Timestamp: 时间戳(可选，默认为当前时间)
	//
	// 返回值:
	//   - error: 事件提交失败时返回错误
	//     - types.ErrState: 系统未运行
	//     - types.ErrQueue: 事件队列已满
	//
	// 示例:
	//   // 创建并处理系统状态变更事件
	//   event := types.SystemEvent{
	//       Type:      types.EventStateChange,
	//       Timestamp: time.Now(),
	//       Source:    "client",
	//       Data: map[string]interface{}{
	//           "oldState": "initializing",
	//           "newState": "running",
	//           "component": "quantum",
	//       },
	//   }
	//
	//   if err := client.HandleEvent(event); err != nil {
	//       log.Printf("事件处理失败: %v", err)
	//   }
	return c.sys.HandleEvent(event)
}

// Optimize 执行系统优化
func (c *Client) Optimize(params types.OptimizationParams) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 该方法用于触发系统优化过程，根据指定的优化参数调整系统配置以达到最佳性能。
	// 优化过程将分析当前系统状态，并根据指定的目标和约束自动调整系统参数。
	//
	// 参数:
	//   - params: 优化参数，包含:
	//     - Goals: 优化目标(如性能、能源效率、稳定性等)
	//     - Targets: 目标值映射
	//     - Constraints: 优化约束条件
	//     - MaxIterations: 最大迭代次数
	//     - TimeLimit: 优化时间限制
	//
	// 返回值:
	//   - error: 优化过程中的错误
	//
	// 示例:
	//   params := types.OptimizationParams{
	//       Goals: types.OptimizationGoals{
	//           Targets: map[string]float64{
	//               "performance": 0.9,
	//               "stability": 0.8,
	//           },
	//           Weights: map[string]float64{
	//               "performance": 0.7,
	//               "stability": 0.3,
	//           },
	//           Constraints: map[string]types.Constraint{
	//               "memory": {Max: 0.8},
	//               "cpu": {Max: 0.7},
	//           },
	//           TimeLimit: 5 * time.Minute,
	//       },
	//       MaxIterations: 100,
	//   }
	//
	//   if err := client.Optimize(params); err != nil {
	//       log.Printf("优化失败: %v", err)
	//   }
	if !c.isRunning {
		return types.ErrNotRunning
	}

	return c.sys.Optimize(params)
}

// Reset 重置系统状态
func (c *Client) Reset() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 重置系统到初始化状态，停止所有运行中的组件，清除所有状态数据并重新初始化各子系统。
	// 此操作会丢失所有未保存的数据和运行状态，包括:
	//   - 所有模型的当前状态
	//   - 未处理的事件队列
	//   - 系统指标和统计数据
	//   - 错误和运行记录
	//
	// 返回值:
	//   - error: 重置过程中发生的错误
	//     - nil: 重置成功
	//     - "failed to stop system": 系统停止失败
	//     - "failed to reinitialize subsystems": 子系统重新初始化失败
	//
	// 示例:
	//   // 重置系统
	//   if err := client.Reset(); err != nil {
	//       log.Printf("系统重置失败: %v", err)
	//       return err
	//   }
	//
	//   // 重置后需要重新启动系统
	//   if err := client.Start(); err != nil {
	//       log.Printf("系统启动失败: %v", err)
	//       return err
	//   }
	//
	//   fmt.Println("系统已成功重置和启动")
	// 重置运行状态
	c.isRunning = false

	// 委托给system处理
	return c.sys.Reset()
}

// Synchronize 同步系统状态
func (c *Client) Synchronize(params types.SyncParams) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 该方法用于在系统组件之间或与外部系统同步状态数据。
	// 支持多种同步模式，如增量同步、完全同步、单向同步等。
	//
	// 参数:
	//   - params: 同步参数，包含:
	//     - Mode: 同步模式(增量/完全/单向)
	//     - Target: 同步目标(组件名称或外部系统标识)
	//     - Source: 同步源(可选，默认为当前系统)
	//     - Options: 同步选项(如超时设置、冲突解决策略等)
	//     - Filter: 数据过滤器(指定同步哪些数据)
	//
	// 返回值:
	//   - error: 同步过程中的错误
	//
	// 示例:
	//   // 执行系统组件间同步
	//   params := types.SyncParams{
	//       Mode: types.SyncModeIncremental,
	//       Target: "quantum_state",
	//       Options: map[string]interface{}{
	//           "timeout": 30 * time.Second,
	//           "conflict_policy": "latest_wins",
	//       },
	//   }
	//
	//   if err := client.Synchronize(params); err != nil {
	//       log.Printf("同步失败: %v", err)
	//   }
	if !c.isRunning {
		return types.ErrNotRunning
	}

	return c.sys.Synchronize(params)
}

// Transform 执行系统转换
func (c *Client) Transform(pattern model.TransformPattern) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 执行特定模式的系统转换，会应用于所有注册的模型。
	// 转换会改变系统的能量分布、阴阳平衡度、和谐度等状态。
	//
	// 参数:
	//   - pattern: 转换模式，可选值:
	//     - model.PatternBalance: 平衡模式，使系统达到平衡状态
	//     - model.PatternForward: 前进模式，增加系统能量和活性
	//     - model.PatternReverse: 逆转模式，降低系统能量
	//     - model.PatternNatural: 自然模式，根据当前状态自动选择合适的转换
	//
	// 返回值:
	//   - error: 转换失败时返回错误
	//
	// 示例:
	//   // 执行平衡转换
	//   if err := client.Transform(model.PatternBalance); err != nil {
	//       log.Printf("系统平衡失败: %v", err)
	//   }
	//
	//   // 执行前进转换
	//   if err := client.Transform(model.PatternForward); err != nil {
	//       log.Printf("系统前进失败: %v", err)
	//   }
	if !c.isRunning {
		return types.ErrNotRunning
	}

	return c.sys.Transform(pattern)
}
