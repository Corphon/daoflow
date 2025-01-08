// model/export.go

package model

// 对外暴露的工厂函数
var (
    // CreateModel 创建模型的工厂函数
    CreateModel = createModel
    
    // CreateIntegratedModel 创建集成模型的工厂函数
    CreateIntegratedModel = createIntegratedModel
)

// ExportedInterface 对外暴露的接口类型
type ExportedInterface interface {
    // FlowModel 基础模型接口
    FlowModel
    
    // GetModelType 获取模型类型
    GetModelType() ModelType
    
    // GetState 获取模型状态
    GetState() ModelState
    
    // Transform 执行模型转换
    Transform(pattern TransformPattern) error
}

// createModel 创建单个模型
func createModel(modelType ModelType, opts ...ConfigOption) (ExportedInterface, error) {
    config := NewConfig(modelType, opts...)
    if err := config.Validate(); err != nil {
        return nil, err
    }

    var model ExportedInterface
    switch modelType {
    case ModelYinYang:
        model = NewYinYangFlow()
    case ModelWuXing:
        model = NewWuXingFlow()
    case ModelBaGua:
        wx := NewWuXingFlow()
        model = NewBaGuaFlow(wx)
    case ModelGanZhi:
        wx := NewWuXingFlow()
        model = NewGanZhiFlow(wx)
    default:
        return nil, NewModelError(ErrCodeInitialization, "unsupported model type", nil)
    }

    return model, nil
}

// createIntegratedModel 创建集成模型
func createIntegratedModel(opts ...ConfigOption) (ExportedInterface, error) {
    config := NewConfig(ModelIntegrate, opts...)
    if err := config.Validate(); err != nil {
        return nil, err
    }

    return NewIntegrateFlow(), nil
}

// ModelStateSnapshot 模型状态快照
type ModelStateSnapshot struct {
    Type      ModelType             // 模型类型
    Phase     Phase                // 当前相位
    Energy    float64              // 当前能量
    Nature    Nature               // 当前性质
    Props     map[string]interface{} // 属性集合
}

// GetModelStateSnapshot 获取模型状态快照
func GetModelStateSnapshot(model ExportedInterface) *ModelStateSnapshot {
    if model == nil {
        return nil
    }

    state := model.GetState()
    return &ModelStateSnapshot{
        Type:   model.GetModelType(),
        Phase:  state.Phase,
        Energy: state.Energy,
        Nature: state.Nature,
        Props:  state.Properties,
    }
}

// TransformRequest 转换请求结构
type TransformRequest struct {
    ModelType ModelType        // 目标模型类型
    Pattern   TransformPattern // 转换模式
    Config    *ModelConfig    // 可选配置
}

// ExecuteTransform 执行模型转换
func ExecuteTransform(model ExportedInterface, req TransformRequest) error {
    if model == nil {
        return NewModelError(ErrCodeOperation, "model is nil", nil)
    }

    if model.GetModelType() != req.ModelType {
        return NewModelError(ErrCodeOperation, "model type mismatch", nil)
    }

    return model.Transform(req.Pattern)
}

// ModelRegistry 模型注册表
var ModelRegistry = struct {
    // 支持的模型类型
    SupportedTypes []ModelType

    // 模型类型描述
    TypeDescriptions map[ModelType]string

    // 默认配置
    DefaultConfigs map[ModelType]ModelConfig
}{
    SupportedTypes: []ModelType{
        ModelYinYang,
        ModelWuXing,
        ModelBaGua,
        ModelGanZhi,
        ModelIntegrate,
    },
    TypeDescriptions: map[ModelType]string{
        ModelYinYang:   "阴阳模型 - 基于阴阳理论的基础模型",
        ModelWuXing:    "五行模型 - 基于五行相生相克理论的模型",
        ModelBaGua:     "八卦模型 - 基于八卦变化理论的模型",
        ModelGanZhi:    "干支模型 - 基于天干地支理论的模型",
        ModelIntegrate: "集成模型 - 整合所有理论的综合模型",
    },
    DefaultConfigs: map[ModelType]ModelConfig{
        ModelYinYang:   YinYangConfig,
        ModelWuXing:    WuXingConfig,
        ModelBaGua:     BaGuaConfig,
        ModelGanZhi:    GanZhiConfig,
        ModelIntegrate: DefaultConfig,
    },
}

// Version 版本信息
const (
    Version     = "1.0.0"
    APIVersion  = "v1"
    CoreVersion = "1.0.0"
)

// 导出错误码
const (
    ErrInvalidModel      = ErrCodeInitialization
    ErrInvalidOperation  = ErrCodeOperation
    ErrInvalidState      = ErrCodeState
    ErrInvalidConfig     = ErrCodeConfiguration
)
