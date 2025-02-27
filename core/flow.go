// core/flow.go

package core

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// FlowState 定义流的状态
type FlowState uint8

const (
	FlowStateVoid         FlowState = iota // 虚无态：初始状态
	FlowStateInactive                      // 未激活态：已创建但未开始流动
	FlowStateFlowing                       // 流动态：正在进行能量流动
	FlowStateStatic                        // 静止态：暂时停止流动
	FlowStateTransforming                  // 转化态：正在进行状态转换
	FlowStateTerminated                    // 终止态：完全停止流动
)

// FlowDirection 定义流动方向
type FlowDirection struct {
	X     float64 // X轴方向分量
	Y     float64 // Y轴方向分量
	Z     float64 // Z轴方向分量
	Angle float64 // 方向角度(0-360)
}

// Flow 定义基本流结构
type Flow struct {
	ID        string         // 唯一标识
	Energy    float64        // 能量值(0-100)
	State     FlowState      // 当前状态
	Direction *FlowDirection // 流动方向

	created  time.Time // 创建时间
	modified time.Time // 最后修改时间
	//mu        sync.RWMutex   // 并发控制
}

// ------------------------------------------
// generateID 生成唯一标识符
func GenerateID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// 如果随机数生成失败，使用时间戳作为备选方案
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

// FlowSource 定义流的源接口
type FlowSource interface {
	// Initialize 初始化流
	Initialize(ctx context.Context) error

	// StartFlow 开始流动
	StartFlow(ctx context.Context) error

	// StopFlow 停止流动
	StopFlow(ctx context.Context) error

	// Transform 状态转换
	Transform(ctx context.Context, newState FlowState) error

	// AdjustEnergy 调节能量
	AdjustEnergy(delta float64) error

	// AdjustDirection 调整方向
	AdjustDirection(direction *FlowDirection) error

	// GetState 获取当前状态
	GetState() FlowState

	// GetEnergy 获取当前能量
	GetEnergy() float64
}

// NewFlow 创建新的流
func NewFlow() *Flow {
	now := time.Now()
	return &Flow{
		ID:     GenerateID(),
		Energy: 50.0, // 初始能量设为中等水平
		State:  FlowStateVoid,
		Direction: &FlowDirection{
			X:     0,
			Y:     0,
			Z:     0,
			Angle: 0,
		},
		created:  now,
		modified: now,
	}
}

// BaseFlow 提供FlowSource的基础实现
type BaseFlow struct {
	flow   *Flow
	config *FlowConfig
	ctx    context.Context
	cancel context.CancelFunc

	observers []FlowObserver
	mu        sync.RWMutex
}

// FlowConfig 流的配置
type FlowConfig struct {
	MinEnergy     float64
	MaxEnergy     float64
	FlowInterval  time.Duration
	MaxTransforms int
}

// NewBaseFlow 创建基础流实现
func NewBaseFlow(config *FlowConfig) *BaseFlow {
	ctx, cancel := context.WithCancel(context.Background())
	return &BaseFlow{
		flow:      NewFlow(),
		config:    config,
		ctx:       ctx,
		cancel:    cancel,
		observers: make([]FlowObserver, 0),
	}
}

// Initialize 实现初始化
func (bf *BaseFlow) Initialize(ctx context.Context) error {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	if bf.flow.State != FlowStateVoid {
		return fmt.Errorf("invalid state for initialization: %v", bf.flow.State)
	}

	bf.flow.State = FlowStateInactive
	bf.flow.modified = time.Now()
	bf.notifyObservers(FlowEventInitialize)

	return nil
}

// StartFlow 开始流动
func (bf *BaseFlow) StartFlow(ctx context.Context) error {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	if bf.flow.State != FlowStateInactive && bf.flow.State != FlowStateStatic {
		return fmt.Errorf("invalid state for starting flow: %v", bf.flow.State)
	}

	bf.flow.State = FlowStateFlowing
	bf.flow.modified = time.Now()
	bf.notifyObservers(FlowEventStart)

	return nil
}

// StopFlow 停止流动
func (bf *BaseFlow) StopFlow(ctx context.Context) error {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	if bf.flow.State != FlowStateFlowing {
		return fmt.Errorf("invalid state for stopping flow: %v", bf.flow.State)
	}

	bf.flow.State = FlowStateStatic
	bf.flow.modified = time.Now()
	bf.notifyObservers(FlowEventStop)

	return nil
}

// Transform 状态转换
func (bf *BaseFlow) Transform(ctx context.Context, newState FlowState) error {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	if !isValidFlowTransition(bf.flow.State, newState) {
		return fmt.Errorf("invalid state transition from %v to %v", bf.flow.State, newState)
	}

	bf.flow.State = newState
	bf.flow.modified = time.Now()
	bf.notifyObservers(FlowEventTransform)

	return nil
}

// AdjustEnergy 调节能量
func (bf *BaseFlow) AdjustEnergy(delta float64) error {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	newEnergy := bf.flow.Energy + delta
	if newEnergy < bf.config.MinEnergy || newEnergy > bf.config.MaxEnergy {
		return fmt.Errorf("energy level %f out of range [%f, %f]",
			newEnergy, bf.config.MinEnergy, bf.config.MaxEnergy)
	}

	bf.flow.Energy = newEnergy
	bf.flow.modified = time.Now()
	bf.notifyObservers(FlowEventEnergyChange)

	return nil
}

// AdjustDirection 调整方向
func (bf *BaseFlow) AdjustDirection(direction *FlowDirection) error {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	bf.flow.Direction = direction
	bf.flow.modified = time.Now()
	bf.notifyObservers(FlowEventDirectionChange)

	return nil
}

// GetState 获取当前状态
func (bf *BaseFlow) GetState() FlowState {
	bf.mu.RLock()
	defer bf.mu.RUnlock()
	return bf.flow.State
}

// GetEnergy 获取当前能量
func (bf *BaseFlow) GetEnergy() float64 {
	bf.mu.RLock()
	defer bf.mu.RUnlock()
	return bf.flow.Energy
}

// FlowEvent 流事件类型
type FlowEvent uint8

const (
	FlowEventInitialize FlowEvent = iota
	FlowEventStart
	FlowEventStop
	FlowEventTransform
	FlowEventEnergyChange
	FlowEventDirectionChange
)

// FlowObserver 观察者接口
type FlowObserver interface {
	OnFlowEvent(event FlowEvent, flow *Flow)
}

// AddObserver 添加观察者
func (bf *BaseFlow) AddObserver(observer FlowObserver) {
	bf.mu.Lock()
	defer bf.mu.Unlock()
	bf.observers = append(bf.observers, observer)
}

// notifyObservers 通知所有观察者
func (bf *BaseFlow) notifyObservers(event FlowEvent) {
	for _, observer := range bf.observers {
		observer.OnFlowEvent(event, bf.flow)
	}
}

// isValidFlowTransition 检查状态转换是否有效
func isValidFlowTransition(current, new FlowState) bool {
	transitions := map[FlowState][]FlowState{
		FlowStateVoid:         {FlowStateInactive},
		FlowStateInactive:     {FlowStateFlowing, FlowStateTerminated},
		FlowStateFlowing:      {FlowStateStatic, FlowStateTransforming, FlowStateTerminated},
		FlowStateStatic:       {FlowStateFlowing, FlowStateTransforming, FlowStateTerminated},
		FlowStateTransforming: {FlowStateFlowing, FlowStateStatic, FlowStateTerminated},
		FlowStateTerminated:   {FlowStateVoid}, // 允许重新初始化
	}

	validStates, exists := transitions[current]
	if !exists {
		return false
	}

	for _, validState := range validStates {
		if new == validState {
			return true
		}
	}
	return false
}
