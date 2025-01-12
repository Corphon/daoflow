//system/evolution/mutation/handler.go

package mutation

import (
    "sync"
    "time"

    "github.com/Corphon/daoflow/system/common"
)

// MutationHandler 突变处理器
type MutationHandler struct {
    mu sync.RWMutex

    // 基础配置
    config struct {
        responseThreshold float64       // 响应阈值
        maxRetries       int           // 最大重试次数
        stabilityTarget  float64       // 稳定性目标
        adaptiveResponse bool          // 自适应响应
    }

    // 处理状态
    state struct {
        active      map[string]*MutationResponse  // 活跃响应
        history     []ResponseEvent              // 响应历史
        strategies  map[string]*ResponseStrategy  // 响应策略
    }

    // 依赖项
    detector common.MutationDetector
    analyzer common.PatternAnalyzer
}

// 确保实现接口
var _ common.MutationHandler = (*MutationHandler)(nil)

// MutationResponse 突变响应
type MutationResponse struct {
    ID          string                // 响应ID
    MutationID  string                // 对应突变ID
    Strategy    *ResponseStrategy     // 使用的策略
    Actions     []ResponseAction      // 响应动作
    Status      string                // 当前状态
    Progress    float64               // 进度
    StartTime   time.Time            // 开始时间
    LastUpdate  time.Time            // 最后更新时间
    Retries     int                  // 重试次数
}

// ResponseStrategy 响应策略
type ResponseStrategy struct {
    ID          string               // 策略ID
    Type        string               // 策略类型
    Conditions  []ResponseCondition  // 触发条件
    Actions     []ActionTemplate     // 动作模板
    Priority    int                  // 优先级
    Success     float64              // 成功率
}

// ResponseCondition 响应条件
type ResponseCondition struct {
    Type        string               // 条件类型
    Target      string               // 目标对象
    Operator    string               // 操作符
    Value       interface{}          // 比较值
    Weight      float64              // 权重
}

// ActionTemplate 动作模板
type ActionTemplate struct {
    Type        string               // 动作类型
    Parameters  map[string]interface{} // 参数模板
    Constraints []ActionConstraint    // 执行约束
    Timeout     time.Duration        // 超时时间
}

// ResponseAction 响应动作
type ResponseAction struct {
    ID          string               // 动作ID
    Type        string               // 动作类型
    Parameters  map[string]interface{} // 实际参数
    Status      string               // 执行状态
    Result      interface{}          // 执行结果
    StartTime   time.Time           // 开始时间
    EndTime     time.Time           // 结束时间
}

// ResponseEvent 响应事件
type ResponseEvent struct {
    Timestamp   time.Time
    ResponseID  string
    Type        string
    Status      string
    Details     map[string]interface{}
}

// NewMutationHandler 创建新的突变处理器
func NewMutationHandler(
    detector *MutationDetector,
    generator *pattern.PatternGenerator) *MutationHandler {
    
    mh := &MutationHandler{
        detector:  detector,
        generator: generator,
    }

    // 初始化配置
    mh.config.responseThreshold = 0.7
    mh.config.maxRetries = 3
    mh.config.stabilityTarget = 0.8
    mh.config.adaptiveResponse = true

    // 初始化状态
    mh.state.active = make(map[string]*MutationResponse)
    mh.state.history = make([]ResponseEvent, 0)
    mh.state.strategies = make(map[string]*ResponseStrategy)

    return mh
}

// Handle 处理突变
func (mh *MutationHandler) Handle() error {
    mh.mu.Lock()
    defer mh.mu.Unlock()

    // 获取待处理的突变
    mutations, err := mh.detector.GetActiveMutations()
    if err != nil {
        return err
    }

    // 为每个突变选择策略
    responses := mh.selectStrategies(mutations)

    // 执行响应
    if err := mh.executeResponses(responses); err != nil {
        return err
    }

    // 监控响应进度
    mh.monitorResponses()

    // 更新响应状态
    mh.updateResponseStatus()

    return nil
}

// RegisterStrategy 注册响应策略
func (mh *MutationHandler) RegisterStrategy(strategy *ResponseStrategy) error {
    if strategy == nil {
        return model.WrapError(nil, model.ErrCodeValidation, "nil strategy")
    }

    mh.mu.Lock()
    defer mh.mu.Unlock()

    // 验证策略
    if err := mh.validateStrategy(strategy); err != nil {
        return err
    }

    // 存储策略
    mh.state.strategies[strategy.ID] = strategy

    return nil
}

// selectStrategies 选择响应策略
func (mh *MutationHandler) selectStrategies(
    mutations []*Mutation) []*MutationResponse {
    
    responses := make([]*MutationResponse, 0)

    for _, mutation := range mutations {
        // 检查是否已有响应
        if mh.hasActiveResponse(mutation.ID) {
            continue
        }

        // 选择最佳策略
        strategy := mh.selectBestStrategy(mutation)
        if strategy == nil {
            continue
        }

        // 创建响应
        response := mh.createResponse(mutation, strategy)
        responses = append(responses, response)
    }

    return responses
}

// executeResponses 执行响应
func (mh *MutationHandler) executeResponses(responses []*MutationResponse) error {
    for _, response := range responses {
        // 准备执行环境
        context := mh.prepareExecutionContext(response)

        // 执行响应动作
        if err := mh.executeResponseActions(response, context); err != nil {
            // 记录错误但继续执行其他响应
            mh.recordResponseEvent(response, "execution_error", map[string]interface{}{
                "error": err.Error(),
            })
            continue
        }

        // 更新响应状态
        response.LastUpdate = time.Now()
        mh.state.active[response.ID] = response
    }

    return nil
}

// monitorResponses 监控响应进度
func (mh *MutationHandler) monitorResponses() {
    for id, response := range mh.state.active {
        // 检查响应是否超时
        if mh.isResponseTimeout(response) {
            mh.handleResponseTimeout(response)
            continue
        }

        // 更新进度
        progress := mh.calculateResponseProgress(response)
        response.Progress = progress

        // 检查完成状态
        if mh.isResponseComplete(response) {
            mh.finalizeResponse(response)
            delete(mh.state.active, id)
        }
    }
}

// updateResponseStatus 更新响应状态
func (mh *MutationHandler) updateResponseStatus() {
    currentTime := time.Now()

    for _, response := range mh.state.active {
        // 检查动作状态
        allComplete := true
        for _, action := range response.Actions {
            if action.Status != "completed" {
                allComplete = false
                break
            }
        }

        // 更新响应状态
        if allComplete {
            response.Status = "completed"
        } else if response.Retries >= mh.config.maxRetries {
            response.Status = "failed"
        }

        // 记录状态变更
        mh.recordResponseEvent(response, "status_update", map[string]interface{}{
            "status": response.Status,
            "progress": response.Progress,
        })

        response.LastUpdate = currentTime
    }
}

// 辅助函数

func (mh *MutationHandler) validateStrategy(strategy *ResponseStrategy) error {
    if strategy.ID == "" {
        return model.WrapError(nil, model.ErrCodeValidation, "empty strategy ID")
    }

    // 验证条件
    for _, condition := range strategy.Conditions {
        if err := mh.validateCondition(condition); err != nil {
            return err
        }
    }

    // 验证动作模板
    for _, template := range strategy.Actions {
        if err := mh.validateActionTemplate(template); err != nil {
            return err
        }
    }

    return nil
}

func (mh *MutationHandler) recordResponseEvent(
    response *MutationResponse,
    eventType string,
    details map[string]interface{}) {
    
    event := ResponseEvent{
        Timestamp:  time.Now(),
        ResponseID: response.ID,
        Type:      eventType,
        Status:    response.Status,
        Details:   details,
    }

    mh.state.history = append(mh.state.history, event)

    // 限制历史记录长度
    if len(mh.state.history) > maxHistoryLength {
        mh.state.history = mh.state.history[1:]
    }
}

const (
    maxHistoryLength = 1000
)
