
# DaoFlow Framework

[![Version](https://img.shields.io/badge/version-v2.0.0-blue.svg)](https://github.com/Corphon/daoflow)
[![Go Reference](https://pkg.go.dev/badge/github.com/Corphon/daoflow.svg)](https://pkg.go.dev/github.com/Corphon/daoflow)
[![License](https://img.shields.io/github/license/Corphon/daoflow.svg)](LICENSE)

DaoFlow 是一个基于 Go 语言的高性能自适应系统框架，融合了东方哲学中阴阳五行的思想，实现了一个能够自我演化、动态平衡的系统架构。

## 核心特性

- **自适应演化系统**: 基于复杂系统理论的自适应演化机制
- **动态能量平衡**: 智能的能量分配和调节系统
- **模式识别与涌现**: 支持复杂模式的识别和新特性的涌现
- **高性能事件处理**: 基于优先级的事件队列和动态缓冲区
- **实时监控和度量**: 全面的系统状态监控和性能指标收集
- **健康检查机制**: 多维度的系统健康状态评估

## 安装

```bash
go get github.com/Corphon/daoflow
```

## 快速开始

### 基础使用

```go
package main

import (
    "log"
    "github.com/Corphon/daoflow/api"
)

func main() {
    // 创建API客户端
    opts := &api.Options{
        SystemConfig: &system.SystemConfig{
            Capacity: 2000.0,
            Threshold: 0.7,
        },
        Debug: true,
    }
    
    client, err := api.NewDaoFlowAPI(opts)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // 初始化并启动系统
    if err := client.Lifecycle().Initialize(); err != nil {
        log.Fatal(err)
    }
    if err := client.Lifecycle().Start(); err != nil {
        log.Fatal(err)
    }

    // 监控系统事件
    events, _ := client.Events().Subscribe(api.EventFilter{
        Types: []api.EventType{api.EventSystemStartup, api.EventStateChange},
        Priority: api.PriorityHigh,
    })

    // 处理事件
    for event := range events {
        log.Printf("Received event: %s", event.Type)
    }
}
```

## 核心模块

### 系统核心 (System Core)

- 负责整体系统的协调和控制
- 提供生命周期管理
- 处理系统级事件
- 维护核心状态

### 能量系统 (Energy System)

```go
// 配置能量系统
energyConfig := api.EnergyConfig{
    MaxCapacity: 1000.0,
    Distribution: api.EnergyDistribution{
        Pattern: 0.3,    // 模式识别
        Evolution: 0.3,  // 演化过程
        Adaptation: 0.2, // 适应调整
        Reserve: 0.2,    // 能量储备
    },
}
```

### 演化系统 (Evolution System)

- 自适应演化机制
- 支持多维度演化
- 智能状态转换
- 演化链追踪

### 事件系统 (Events System)

- 支持优先级队列
- 动态缓冲区管理
- 灵活的事件订阅
- 事件历史追踪

### 监控系统 (Metrics System)

- 实时性能监控
- 多维度指标收集
- 自定义指标支持
- 指标聚合分析

## 高级特性

### 动态缓冲区

```go
bufferConfig := system.ResizePolicy{
    MinCapacity:    100,
    MaxCapacity:    10000,
    GrowthFactor:   2.0,
    ShrinkFactor:   0.5,
    ResizeInterval: time.Minute,
}
```

### 健康检查

```go
health, err := client.Health().GetSystemHealth()
if err != nil {
    log.Fatal(err)
}
log.Printf("System health score: %f", health.HealthScore)
```

### 配置管理

```go
err := client.Config().SetConfig("evolution.rate", 0.15, api.ScopeEvolution, nil)
if err != nil {
    log.Fatal(err)
}
```

## 最佳实践

1. **合理配置能量分配**
   - 根据系统负载调整能量分配比例
   - 保持适当的能量储备
   - 监控能量使用效率

2. **优化事件处理**
   - 使用合适的事件优先级
   - 配置合理的缓冲区大小
   - 及时处理关键事件

3. **监控系统健康**
   - 定期检查系统状态
   - 设置合理的告警阈值
   - 保持系统平衡

## 性能优化

- 使用动态缓冲区自动调整
- 实现高效的事件处理机制
- 优化能量分配算法
- 合理配置系统参数

## 贡献指南

欢迎贡献代码或提出建议，请参考我们的[贡献指南](CONTRIBUTING.md)。

## 许可证

本项目基于 [Apache License 2.0](LICENSE) 进行授权。

## 联系我们

- GitHub Issues: [https://github.com/Corphon/daoflow/issues](https://github.com/Corphon/daoflow/issues)
- Email: [contact@corphon.com](mailto:songkf@foxmail.com)

