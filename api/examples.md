# api.go
```go
func Example() {
    // 创建API客户端
    opts := &Options{
        SystemConfig: &system.SystemConfig{
            Capacity: 2000.0,
            Threshold: 0.7,
        },
        Debug: true,
    }
    
    client, err := NewDaoFlowAPI(opts)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // 初始化系统
    if err := client.Lifecycle().Initialize(); err != nil {
        log.Fatal(err)
    }

    // 启动系统
    if err := client.Lifecycle().Start(); err != nil {
        log.Fatal(err)
    }

    // 监控演化事件
    eventChan, err := client.Events().Subscribe([]string{"evolution"})
    if err != nil {
        log.Fatal(err)
    }

    go func() {
        for event := range eventChan {
            log.Printf("Evolution event: %+v", event)
        }
    }()

    // 触发演化
    if err := client.Evolution().TriggerEvolution("optimize"); err != nil {
        log.Fatal(err)
    }

    // 获取系统状态
    status, err := client.Lifecycle().GetStatus()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("System status: %+v", status)
}
···
