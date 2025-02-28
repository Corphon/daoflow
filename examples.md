# DaoFlow Examples

A collection of examples demonstrating how to harness the power of DaoFlow - where Eastern philosophical wisdom meets modern computational systems.

## Table of Contents

- [Getting Started](#getting-started)
- [Core Models](#core-models)
  - [Yin-Yang Flow](#yin-yang-flow)
  - [Wu Xing (Five Elements)](#wu-xing-five-elements)
  - [Ba Gua (Eight Trigrams)](#ba-gua-eight-trigrams)
  - [Gan Zhi (Celestial Stems and Terrestrial Branches)](#gan-zhi-celestial-stems-and-terrestrial-branches)
- [System Operations](#system-operations)
  - [Energy Management](#energy-management)
  - [State Transformations](#state-transformations)
  - [Event Handling](#event-handling)
- [Advanced Usage](#advanced-usage)
  - [Pattern Recognition](#pattern-recognition)
  - [System Optimization](#system-optimization)
  - [Monitoring](#monitoring)
- [Complete Applications](#complete-applications)

## Getting Started

Initialize and run the DaoFlow system with minimal configuration:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/Corphon/daoflow/api"
    "github.com/Corphon/daoflow/system"
)

func main() {
    // Create a new DaoFlow client with default configuration
    client, err := api.NewClient(nil)
    if err != nil {
        log.Fatalf("Failed to create DaoFlow client: %v", err)
    }
    defer client.Close()
    
    // Initialize with context
    ctx := context.Background()
    if err := client.Initialize(ctx); err != nil {
        log.Fatalf("Failed to initialize: %v", err)
    }
    
    // Start the system
    if err := client.Start(); err != nil {
        log.Fatalf("Failed to start: %v", err)
    }
    
    fmt.Println("DaoFlow system is running!")
    
    // Check system status
    status := client.GetSystemStatus()
    fmt.Printf("Current system status: %s\n", status)
    
    // Use the system for 5 seconds
    time.Sleep(5 * time.Second)
    
    // Graceful shutdown
    fmt.Println("Shutting down...")
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := client.Shutdown(shutdownCtx); err != nil {
        log.Printf("Shutdown error: %v", err)
    }
}
```

## Core Models

### Yin-Yang Flow

Example demonstrating the Yin-Yang model, which embodies the foundational duality principle:

```go
func YinYangExample() {
    client, _ := api.NewClient(nil)
    defer client.Close()
    
    // Initialize and start
    ctx := context.Background()
    client.Initialize(ctx)
    client.Start()
    
    // Get the YinYang flow model
    yinYang := client.GetYinYangFlow()
    
    // Get current YinYang energy values
    energy := yinYang.GetYinYangEnergy()
    fmt.Printf("Yin Energy: %.2f, Yang Energy: %.2f\n", energy.YinEnergy, energy.YangEnergy)
    fmt.Printf("Harmony: %.2f\n", energy.Harmony)
    
    // Transform to balance state
    fmt.Println("Balancing Yin and Yang energies...")
    err := client.TransformYinYang(model.PatternBalance)
    if err != nil {
        log.Printf("Transform error: %v", err)
        return
    }
    
    // Check new state after transformation
    energy = yinYang.GetYinYangEnergy()
    fmt.Printf("After balance - Yin: %.2f, Yang: %.2f, Harmony: %.2f\n", 
        energy.YinEnergy, energy.YangEnergy, energy.Harmony)
    
    // Transform to increase Yang energy
    fmt.Println("Increasing Yang energy...")
    client.TransformYinYang(model.PatternForward)
    
    // Check final state
    energy = yinYang.GetYinYangEnergy()
    fmt.Printf("Final state - Yin: %.2f, Yang: %.2f, Harmony: %.2f\n", 
        energy.YinEnergy, energy.YangEnergy, energy.Harmony)
}
```

### Wu Xing (Five Elements)

Example demonstrating the Five Elements (Wood, Fire, Earth, Metal, Water) relationships:

```go
func WuXingExample() {
    client, _ := api.NewClient(nil)
    defer client.Close()
    
    client.Initialize(context.Background())
    client.Start()
    
    // Get the WuXing flow model
    wuXing := client.GetWuXingFlow()
    
    // Display initial state
    fmt.Println("Initial Five Elements state:")
    displayWuXingState(wuXing)
    
    // Execute generating cycle (相生): Wood→Fire→Earth→Metal→Water→Wood
    fmt.Println("\nExecuting Generating Cycle (相生)...")
    client.Transform(model.PatternForward)
    
    fmt.Println("After Generating Cycle:")
    displayWuXingState(wuXing)
    
    // Execute controlling cycle (相克): Wood→Earth→Water→Fire→Metal→Wood
    fmt.Println("\nExecuting Controlling Cycle (相克)...")
    client.Transform(model.PatternReverse)
    
    fmt.Println("After Controlling Cycle:")
    displayWuXingState(wuXing)
}

func displayWuXingState(wuXing *model.WuXingFlow) {
    state := wuXing.GetState()
    
    // In a real implementation, the WuXingFlow would provide methods
    // to get individual element energies
    fmt.Printf("Wood: %.2f\n", state.Properties["wood_energy"].(float64))
    fmt.Printf("Fire: %.2f\n", state.Properties["fire_energy"].(float64))
    fmt.Printf("Earth: %.2f\n", state.Properties["earth_energy"].(float64))
    fmt.Printf("Metal: %.2f\n", state.Properties["metal_energy"].(float64))
    fmt.Printf("Water: %.2f\n", state.Properties["water_energy"].(float64))
    
    fmt.Printf("Overall Harmony: %.2f\n", state.Properties["harmony"].(float64))
}
```

### Ba Gua (Eight Trigrams)

Example demonstrating the Eight Trigrams pattern system:

```go
func BaGuaExample() {
    client, _ := api.NewClient(nil)
    defer client.Close()
    
    client.Initialize(context.Background())
    client.Start()
    
    // Get the BaGua flow model
    baGua := client.GetBaGuaFlow()
    
    // Display initial state
    state := baGua.GetState()
    fmt.Printf("Initial BaGua state - Energy: %.2f, Harmony: %.2f\n",
        state.Energy, state.Properties["harmony"].(float64))
    
    // Perform natural transformation
    fmt.Println("\nPerforming natural transformation...")
    client.Transform(model.PatternForward)
    
    // Check for pattern changes
    state = baGua.GetState()
    fmt.Printf("After transformation - Energy: %.2f, Harmony: %.2f\n",
        state.Energy, state.Properties["harmony"].(float64))
    
    // Achieve resonance between trigrams
    fmt.Println("\nCreating resonance between trigrams...")
    client.Transform(model.PatternReverse)
    
    state = baGua.GetState()
    fmt.Printf("After resonance - Energy: %.2f, Harmony: %.2f, Resonance: %.2f\n",
        state.Energy, 
        state.Properties["harmony"].(float64),
        state.Properties["resonance"].(float64))
}
```

### Gan Zhi (Celestial Stems and Terrestrial Branches)

Example demonstrating the cyclical time system:

```go
func GanZhiExample() {
    client, _ := api.NewClient(nil)
    defer client.Close()
    
    client.Initialize(context.Background())
    client.Start()
    
    // Get the GanZhi flow model
    ganZhi := client.GetGanZhiFlow()
    
    // Get initial state
    state := ganZhi.GetState()
    fmt.Printf("Initial Cycle State - Energy: %.2f\n", state.Energy)
    
    // Advance through several cycle positions
    fmt.Println("\nAdvancing through the cycle...")
    for i := 0; i < 5; i++ {
        client.Transform(model.PatternForward)
        
        // In a real implementation, GanZhiFlow would provide methods to get cycle position
        fmt.Printf("Cycle position %d\n", i+1)
        
        // Display energetic relationships between the elements at this position
        state = ganZhi.GetState()
        fmt.Printf("  Harmony: %.2f\n", state.Properties["harmony"].(float64))
        fmt.Printf("  WuXing Element: %s\n", state.Properties["current_element"].(string))
        fmt.Printf("  Polarity: %s\n", state.Properties["polarity"].(string))
    }
    
    // Reset to start of cycle
    fmt.Println("\nResetting cycle position...")
    client.Transform(model.PatternBalance)
    
    state = ganZhi.GetState()
    fmt.Printf("Reset complete - Energy: %.2f\n", state.Energy)
}
```

## System Operations

### Energy Management

Example showing how to manage system energy:

```go
func EnergyManagementExample() {
    client, _ := api.NewClient(nil)
    defer client.Close()
    
    client.Initialize(context.Background())
    client.Start()
    
    // Get current energy level
    energy := client.GetEnergy()
    fmt.Printf("Initial system energy: %.2f\n", energy)
    
    // Adjust energy levels
    fmt.Println("Increasing system energy by 20%...")
    client.AdjustEnergy(0.2)
    
    // Check new energy level
    energy = client.GetEnergy()
    fmt.Printf("New system energy: %.2f\n", energy)
    
    // Get detailed energy system
    energySystem := client.GetEnergySystem()
    
    // In a real implementation, EnergySystem would provide these methods
    fmt.Printf("Total energy: %.2f\n", energySystem.GetTotalEnergy())
    fmt.Printf("Potential energy: %.2f\n", energySystem.GetPotentialEnergy())
    fmt.Printf("Kinetic energy: %.2f\n", energySystem.GetKineticEnergy())
    
    // Reduce system energy
    fmt.Println("\nReducing system energy by 30%...")
    client.AdjustEnergy(-0.3)
    
    energy = client.GetEnergy()
    fmt.Printf("Final system energy: %.2f\n", energy)
}
```

### State Transformations

Example demonstrating system-wide transformations:

```go
func StateTransformationExample() {
    client, _ := api.NewClient(nil)
    defer client.Close()
    
    client.Initialize(context.Background())
    client.Start()
    
    // Get initial system state
    state := client.GetState()
    fmt.Printf("Initial system state:\n")
    fmt.Printf("  Energy: %.2f\n", state.Energy)
    fmt.Printf("  Entropy: %.2f\n", state.Entropy)
    fmt.Printf("  Harmony: %.2f\n", state.Harmony)
    fmt.Printf("  Balance: %.2f\n", state.Balance)
    
    // Execute system-wide balance transformation
    fmt.Println("\nExecuting system-wide balance transformation...")
    ctx := context.Background()
    err := client.TransformModel(ctx, model.PatternBalance)
    if err != nil {
        log.Printf("Transformation error: %v", err)
        return
    }
    
    // Check new system state
    state = client.GetState()
    fmt.Printf("\nSystem state after balance transformation:\n")
    fmt.Printf("  Energy: %.2f\n", state.Energy)
    fmt.Printf("  Entropy: %.2f\n", state.Entropy)
    fmt.Printf("  Harmony: %.2f\n", state.Harmony)
    fmt.Printf("  Balance: %.2f\n", state.Balance)
    
    // Execute forward transformation to increase system energy and activity
    fmt.Println("\nExecuting forward transformation (increasing energy)...")
    client.TransformModel(ctx, model.PatternForward)
    
    state = client.GetState()
    fmt.Printf("\nSystem state after forward transformation:\n")
    fmt.Printf("  Energy: %.2f\n", state.Energy)
    fmt.Printf("  Entropy: %.2f\n", state.Entropy)
    fmt.Printf("  Harmony: %.2f\n", state.Harmony)
    fmt.Printf("  Balance: %.2f\n", state.Balance)
}
```

### Event Handling

Example showing how to work with the event system:

```go
func EventHandlingExample() {
    client, _ := api.NewClient(nil)
    defer client.Close()
    
    client.Initialize(context.Background())
    client.Start()
    
    // Create event handlers
    stateChangeHandler := &customEventHandler{
        name: "StateChangeHandler",
        handleFunc: func(event types.SystemEvent) error {
            fmt.Printf("[State Change] %s: %v\n", 
                event.Timestamp.Format(time.RFC3339),
                event.Data)
            return nil
        },
    }
    
    transformHandler := &customEventHandler{
        name: "TransformHandler",
        handleFunc: func(event types.SystemEvent) error {
            fmt.Printf("[Transform] %s: %v\n", 
                event.Timestamp.Format(time.RFC3339),
                event.Data)
            return nil
        },
    }
    
    // Subscribe to events
    client.Subscribe("system.state_change", stateChangeHandler)
    client.Subscribe("system.transform", transformHandler)
    
    // Perform actions that will trigger events
    fmt.Println("Performing actions to generate events...")
    
    // Adjust energy (should trigger state change)
    client.AdjustEnergy(0.2)
    
    // Transform models (should trigger transform event)
    client.TransformYinYang(model.PatternForward)
    
    // Wait for events to be processed
    time.Sleep(500 * time.Millisecond)
    
    // Unsubscribe from events
    client.Unsubscribe("system.state_change", stateChangeHandler)
    client.Unsubscribe("system.transform", transformHandler)
}

// Custom event handler implementation
type customEventHandler struct {
    name string
    handleFunc func(types.SystemEvent) error
}

func (h *customEventHandler) HandleEvent(event types.SystemEvent) error {
    return h.handleFunc(event)
}

func (h *customEventHandler) Name() string {
    return h.name
}
```

## Advanced Usage

### Pattern Recognition

Example showing how to use pattern recognition features:

```go
func PatternRecognitionExample() {
    client, _ := api.NewClient(nil)
    defer client.Close()
    
    client.Initialize(context.Background())
    client.Start()
    
    // Generate some test data
    data := map[string]interface{}{
        "energy": []float64{0.5, 0.6, 0.7, 0.8, 0.7, 0.6},
        "phase": []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6},
        "timestamp": time.Now(),
    }
    
    // Detect patterns in the data
    pattern, err := client.DetectPattern(data)
    if err != nil {
        log.Printf("Pattern detection error: %v", err)
        return
    }
    
    // Display detected pattern
    fmt.Printf("Detected pattern: %s\n", pattern.Type)
    fmt.Printf("Pattern strength: %.2f\n", pattern.Strength)
    fmt.Printf("Duration: %v\n", pattern.Duration)
    fmt.Printf("Properties: %v\n", pattern.Properties)
    
    // Analyze the pattern further
    err = client.AnalyzePattern(pattern)
    if err != nil {
        log.Printf("Pattern analysis error: %v", err)
        return
    }
    
    // Pattern now contains additional analysis
    fmt.Printf("\nAfter analysis:\n")
    fmt.Printf("Pattern ID: %s\n", pattern.ID)
    fmt.Printf("Updated strength: %.2f\n", pattern.Strength)
    fmt.Printf("Updated properties: %v\n", pattern.Properties)
}
```

### System Optimization

Example showing how to use the optimization capabilities:

```go
func SystemOptimizationExample() {
    client, _ := api.NewClient(nil)
    defer client.Close()
    
    client.Initialize(context.Background())
    client.Start()
    
    // Create optimization parameters
    params := types.OptimizationParams{
        MaxIterations: 100,
        Goals: types.OptimizationGoals{
            Targets: map[string]float64{
                "performance": 0.9,
                "stability": 0.8,
                "energy_efficiency": 0.85,
            },
            Weights: map[string]float64{
                "performance": 0.5,
                "stability": 0.3,
                "energy_efficiency": 0.2,
            },
        },
    }
    
    // Run optimization process
    fmt.Println("Starting system optimization...")
    err := client.Optimize(params)
    if err != nil {
        log.Printf("Optimization error: %v", err)
        return
    }
    
    // Check system metrics after optimization
    metrics := client.GetSystemMetrics()
    fmt.Println("\nSystem metrics after optimization:")
    
    if performance, ok := metrics["performance"].(float64); ok {
        fmt.Printf("Performance: %.2f\n", performance)
    }
    
    if stability, ok := metrics["stability"].(float64); ok {
        fmt.Printf("Stability: %.2f\n", stability)
    }
    
    if efficiency, ok := metrics["energy_efficiency"].(float64); ok {
        fmt.Printf("Energy Efficiency: %.2f\n", efficiency)
    }
}
```

### Monitoring

Example showing how to monitor system status and performance:

```go
func MonitoringExample() {
    client, _ := api.NewClient(nil)
    defer client.Close()
    
    client.Initialize(context.Background())
    client.Start()
    
    // Get system status
    status := client.GetSystemStatus()
    fmt.Printf("System status: %s\n", status)
    
    // Get detailed metrics
    metrics := client.GetMetrics()
    fmt.Printf("\nSystem metrics:\n")
    fmt.Printf("Energy total: %.2f\n", metrics.Energy.Total)
    fmt.Printf("Energy average: %.2f\n", metrics.Energy.Average)
    fmt.Printf("Energy variance: %.2f\n", metrics.Energy.Variance)
    
    fmt.Printf("\nPerformance metrics:\n")
    fmt.Printf("Throughput: %.2f ops/min\n", metrics.Performance.Throughput)
    fmt.Printf("QPS: %.2f\n", metrics.Performance.QPS)
    fmt.Printf("Error rate: %.2f%%\n", metrics.Performance.ErrorRate*100)
    
    // Get system state
    state := client.GetState()
    fmt.Printf("\nSystem state:\n")
    fmt.Printf("Energy: %.2f\n", state.Energy)
    fmt.Printf("Entropy: %.2f\n", state.Entropy)
    fmt.Printf("Harmony: %.2f\n", state.Harmony)
    fmt.Printf("Balance: %.2f\n", state.Balance)
    fmt.Printf("Phase: %v\n", state.Phase)
    
    // Get model state 
    modelState := client.GetModelState()
    fmt.Printf("\nModel state:\n")
    fmt.Printf("Type: %d\n", modelState.Type)
    fmt.Printf("Energy: %.2f\n", modelState.Energy)
    fmt.Printf("Phase: %d\n", modelState.Phase)
    fmt.Printf("Health: %.2f\n", modelState.Health)
    fmt.Printf("Properties: %v\n", modelState.Properties)
}
```

## Complete Applications

### Dynamic Balancing System

A complete example that creates a self-balancing system:

```go
func DynamicBalancingSystem() {
    // Create client with custom configuration
    cfg := &system.Config{
        CoreConfig: &core.Config{
            InitialEnergy: 100,
            DecayRate: 0.01,
        },
        ModelConfig: &model.ModelConfig{
            EnableYinYang: true,
            EnableWuXing: true,
            EnableBaGua: true,
            EnableGanZhi: true,
        },
    }
    
    client, err := api.NewClient(cfg)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    defer client.Close()
    
    // Initialize and start
    ctx := context.Background()
    if err := client.Initialize(ctx); err != nil {
        log.Fatalf("Failed to initialize: %v", err)
    }
    
    if err := client.Start(); err != nil {
        log.Fatalf("Failed to start: %v", err)
    }
    
    // Create a monitoring goroutine
    done := make(chan bool)
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()
        
        for {
            select {
            case <-done:
                return
            case <-ticker.C:
                state := client.GetState()
                fmt.Printf("[%s] Energy: %.2f, Balance: %.2f, Harmony: %.2f\n",
                    time.Now().Format("15:04:05"),
                    state.Energy,
                    state.Balance,
                    state.Harmony)
                
                // Auto-balance if harmony drops too low
                if state.Harmony < 0.5 {
                    fmt.Println("Harmony is low, applying balance transformation...")
                    client.Transform(model.PatternBalance)
                }
                
                // Add energy if too low
                if state.Energy < 0.3 {
                    fmt.Println("Energy is low, adding energy...")
                    client.AdjustEnergy(0.2)
                }
            }
        }
    }()
    
    // Simulate external influences
    for i := 0; i < 5; i++ {
        time.Sleep(2 * time.Second)
        
        // Apply random transformations to simulate disruptions
        transformType := model.TransformPattern(rand.Intn(3) + 1) // 1, 2, or 3
        fmt.Printf("Applying external transformation: %d\n", transformType)
        client.Transform(transformType)
    }
    
    // Signal monitoring to stop
    close(done)
    
    // Final state
    state := client.GetState()
    fmt.Printf("\nFinal state - Energy: %.2f, Balance: %.2f, Harmony: %.2f\n",
        state.Energy,
        state.Balance,
        state.Harmony)
        
    // Shutdown
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    client.Shutdown(shutdownCtx)
}
```

### Pattern Detection Pipeline

An example showing a complete pattern detection and analysis pipeline:

```go
func PatternDetectionPipeline() {
    client, _ := api.NewClient(nil)
    defer client.Close()
    
    client.Initialize(context.Background())
    client.Start()
    
    // Create event handler for pattern detection
    patternHandler := &customEventHandler{
        name: "PatternHandler",
        handleFunc: func(event types.SystemEvent) error {
            if pattern, ok := event.Data["pattern"].(*model.FlowPattern); ok {
                fmt.Printf("Pattern detected: %s (strength: %.2f)\n", 
                    pattern.Type, pattern.Strength)
                
                // Analyze the pattern
                client.AnalyzePattern(pattern)
                
                fmt.Printf("Analysis results: %v\n", pattern.Properties["analysis"])
            }
            return nil
        },
    }
    
    // Subscribe to pattern events
    client.Subscribe("pattern.detected", patternHandler)
    
    // Generate sample data stream
    for i := 0; i < 10; i++ {
        // Create sample data point with some sine wave patterns
        t := float64(i) * 0.1
        data := map[string]interface{}{
            "value": math.Sin(t*2*math.Pi) + math.Sin(t*4*math.Pi),
            "time": time.Now().Add(time.Duration(i) * time.Second),
            "metadata": map[string]interface{}{
                "source": "sensor_1",
                "type": "waveform",
            },
        }
        
        // Process the data point
        client.DetectPattern(data)
        
        // Wait a bit between data points
        time.Sleep(500 * time.Millisecond)
    }
    
    // Unsubscribe
    client.Unsubscribe("pattern.detected", patternHandler)
}
```

These examples demonstrate the key features of the DaoFlow framework, from basic initialization to advanced pattern recognition and system optimization. The examples are designed to be both educational and practical, showing the integration of Eastern philosophical concepts with modern computational techniques.
