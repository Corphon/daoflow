// core/cycle.go

package core

import (
	"sync"
	"time"
)

// CycleState 周期状态
type CycleState struct {
	Index     int       // 当前周期索引
	Phase     float64   // 当前相位
	Energy    float64   // 周期能量
	Timestamp time.Time // 时间戳
}

// CycleManager 周期管理器
type CycleManager struct {
	mu sync.RWMutex

	// 基本属性
	length  int     // 周期长度
	current int     // 当前位置
	phase   float64 // 当前相位
	energy  float64 // 周期能量

	// 时间相关
	startTime  time.Time // 开始时间
	lastUpdate time.Time // 最后更新时间

	// 历史记录
	history []CycleState
}

// ---------------------------------------------
// NewCycleManager 创建新的周期管理器
func NewCycleManager(length int) *CycleManager {
	if length <= 0 {
		length = 60 // 默认六十周期
	}

	return &CycleManager{
		length:     length,
		current:    0,
		phase:      0,
		energy:     1.0,
		startTime:  time.Now(),
		lastUpdate: time.Now(),
		history:    make([]CycleState, 0),
	}
}

// Initialize 初始化周期管理器
func (cm *CycleManager) Initialize() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.current = 0
	cm.phase = 0
	cm.energy = 1.0
	cm.startTime = time.Now()
	cm.lastUpdate = time.Now()
	cm.history = make([]CycleState, 0)

	return nil
}

// Advance 推进周期
func (cm *CycleManager) Advance() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 更新位置
	cm.current = (cm.current + 1) % cm.length

	// 更新相位
	cm.phase = 2 * float64(cm.current) * 3.14159 / float64(cm.length)

	// 记录状态
	cm.recordState()

	cm.lastUpdate = time.Now()
	return nil
}

// GetCurrent 获取当前位置
func (cm *CycleManager) GetCurrent() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.current
}

// GetPhase 获取当前相位
func (cm *CycleManager) GetPhase() float64 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.phase
}

// GetEnergy 获取周期能量
func (cm *CycleManager) GetEnergy() float64 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.energy
}

// recordState 记录状态
func (cm *CycleManager) recordState() {
	state := CycleState{
		Index:     cm.current,
		Phase:     cm.phase,
		Energy:    cm.energy,
		Timestamp: time.Now(),
	}

	cm.history = append(cm.history, state)
}

// Close 关闭管理器
func (cm *CycleManager) Close() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.history = nil
	return nil
}
