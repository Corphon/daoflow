// model/integrate.go

package model

import (
    "sync"
    "time"

    "github.com/Corphon/daoflow/core"
)

// IntegrateFlow 集成模型 - 实现系统层所需的主要接口
type IntegrateFlow struct {
    *BaseFlowModel // 继承基础模型

    // 子模型 - 对外隐藏实现细节
    yinyang *YinYangFlow
    wuxing  *WuXingFlow
    bagua   *BaGuaFlow
    ganzhi  *GanZhiFlow

    // 集成状态
    state struct {
        yinYangEnergy  float64
        wuxingEnergy   float64
        baguaEnergy    float64
        ganzhiEnergy   float64
    }

    // 同步管理器 - 内部使用
    sync *syncManager
}

// syncManager 同步管理器
type syncManager struct {
    mu sync.RWMutex

    // 内部使用的核心组件
    unifiedField   *core.Field
    unifiedState   *core.QuantumState
    energyNetwork  *core.EnergyNetwork

    // 同步参数
    lastSync      time.Time
    syncInterval  time.Duration
    syncThreshold float64
}

// NewIntegrateFlow 创建集成模型
func NewIntegrateFlow() *IntegrateFlow {
    // 创建基础模型
    base := NewBaseFlowModel(ModelIntegrate, 1000.0)

    // 创建集成模型
    flow := &IntegrateFlow{
        BaseFlowModel: base,
    }

    // 初始化子模型
    flow.initializeSubModels()

    // 初始化同步管理器
    flow.initializeSyncManager()

    return flow
}

// initializeSubModels 初始化子模型
func (f *IntegrateFlow) initializeSubModels() {
    f.yinyang = NewYinYangFlow()
    f.wuxing = NewWuXingFlow()
    f.bagua = NewBaGuaFlow()
    f.ganzhi = NewGanZhiFlow()
}

// initializeSyncManager 初始化同步管理器
func (f *IntegrateFlow) initializeSyncManager() {
    f.sync = &syncManager{
        unifiedField:   core.NewField(),
        unifiedState:   core.NewQuantumState(),
        energyNetwork:  core.NewEnergyNetwork(),
        syncInterval:   time.Second * 5,
        syncThreshold:  0.8,
    }
}

// Start 启动集成模型
func (f *IntegrateFlow) Start() error {
    if err := f.BaseFlowModel.Start(); err != nil {
        return err
    }

    // 启动子模型
    if err := f.startSubModels(); err != nil {
        return WrapError(err, ErrCodeOperation, "failed to start sub-models")
    }

    // 初始化统一场
    if err := f.sync.unifiedField.Initialize(); err != nil {
        return WrapError(err, ErrCodeOperation, "failed to initialize unified field")
    }

    return nil
}

// Stop 停止集成模型
func (f *IntegrateFlow) Stop() error {
    // 停止子模型
    if err := f.stopSubModels(); err != nil {
        return WrapError(err, ErrCodeOperation, "failed to stop sub-models")
    }

    return f.BaseFlowModel.Stop()
}

// Transform 执行集成转换
func (f *IntegrateFlow) Transform(pattern TransformPattern) error {
    if err := f.BaseFlowModel.Transform(pattern); err != nil {
        return err
    }

    // 转换子模型
    if err := f.transformSubModels(pattern); err != nil {
        return WrapError(err, ErrCodeTransform, "sub-model transform failed")
    }

    // 同步子模型状态
    if err := f.synchronizeModels(); err != nil {
        return WrapError(err, ErrCodeSync, "model synchronization failed")
    }

    // 更新集成状态
    return f.updateIntegratedState()
}

// GetSystemState 获取系统状态 - 实现系统层所需的接口
func (f *IntegrateFlow) GetSystemState() SystemState {
    // 获取基础状态
    state := f.BaseFlowModel.GetSystemState()

    // 添加子系统状态
    f.mu.RLock()
    state.YinYang = f.state.yinYangEnergy
    state.WuXingEnergy = f.state.wuxingEnergy
    state.BaGuaEnergy = f.state.baguaEnergy
    state.GanZhiEnergy = f.state.ganzhiEnergy

    // 更新系统详情
    state.System.WuXingEnergy = f.state.wuxingEnergy
    state.System.BaGuaEnergy = f.state.baguaEnergy
    state.System.GanZhiEnergy = f.state.ganzhiEnergy
    f.mu.RUnlock()

    return state
}

// 以下是内部实现方法

func (f *IntegrateFlow) startSubModels() error {
    if err := f.yinyang.Start(); err != nil {
        return err
    }
    if err := f.wuxing.Start(); err != nil {
        return err
    }
    if err := f.bagua.Start(); err != nil {
        return err
    }
    if err := f.ganzhi.Start(); err != nil {
        return err
    }
    return nil
}

func (f *IntegrateFlow) stopSubModels() error {
    if err := f.yinyang.Stop(); err != nil {
        return err
    }
    if err := f.wuxing.Stop(); err != nil {
        return err
    }
    if err := f.bagua.Stop(); err != nil {
        return err
    }
    if err := f.ganzhi.Stop(); err != nil {
        return err
    }
    return nil
}

func (f *IntegrateFlow) transformSubModels(pattern TransformPattern) error {
    if err := f.yinyang.Transform(pattern); err != nil {
        return err
    }
    if err := f.wuxing.Transform(pattern); err != nil {
        return err
    }
    if err := f.bagua.Transform(pattern); err != nil {
        return err
    }
    if err := f.ganzhi.Transform(pattern); err != nil {
        return err
    }
    return nil
}

func (f *IntegrateFlow) synchronizeModels() error {
    f.sync.mu.Lock()
    defer f.sync.mu.Unlock()

    // 检查是否需要同步
    if time.Since(f.sync.lastSync) < f.sync.syncInterval {
        return nil
    }

    // 更新统一场
    quantum, field, energy := f.BaseFlowModel.getInternalState()
    if err := f.sync.unifiedField.Merge(field); err != nil {
        return err
    }

    // 同步量子态
    if err := f.sync.unifiedState.Entangle(quantum); err != nil {
        return err
    }

    // 同步能量网络
    if err := f.sync.energyNetwork.Synchronize(energy); err != nil {
        return err
    }

    f.sync.lastSync = time.Now()
    return nil
}

func (f *IntegrateFlow) updateIntegratedState() error {
    f.mu.Lock()
    defer f.mu.Unlock()

    // 更新子模型能量状态
    f.state.yinYangEnergy = f.yinyang.GetState().Energy
    f.state.wuxingEnergy = f.wuxing.GetState().Energy
    f.state.baguaEnergy = f.bagua.GetState().Energy
    f.state.ganzhiEnergy = f.ganzhi.GetState().Energy

    return nil
}

// Close 关闭集成模型
func (f *IntegrateFlow) Close() error {
    // 关闭子模型
    if err := f.stopSubModels(); err != nil {
        return err
    }

    // 清理同步管理器
    f.sync.unifiedField = nil
    f.sync.unifiedState = nil
    f.sync.energyNetwork = nil

    return f.BaseFlowModel.Close()
}
