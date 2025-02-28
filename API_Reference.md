# DaoFlow API Reference

## Table of Contents

1. [Introduction](#introduction)
2. [Core API](#core-api)
3. [System Components](#system-components)
4. [Flow Models](#flow-models)
5. [Monitoring & Metrics](#monitoring--metrics)
6. [Event System](#event-system)
7. [Configuration](#configuration)
8. [Error Handling](#error-handling)

## Introduction

DaoFlow API provides a comprehensive interface to interact with the quantum-inspired computational framework based on Eastern philosophical principles. This document outlines the available APIs, their usage, and common patterns for effective integration.

The API follows a modular design with specialized interfaces for different aspects of the system:
```MD
DaoFlowAPI
├── Lifecycle() → LifecycleAPI       // System initialization and control
├── Pattern() → PatternAPI           // Pattern recognition and analysis
├── Energy() → EnergyAPI             // Energy management and distribution
├── Evolution() → EvolutionAPI       // Adaptive evolution capabilities
├── Metrics() → MetricsAPI           // Performance monitoring
└── Events() → EventsAPI             // Event subscription and handling
```

## Core API Client

The main entry point for all DaoFlow operations is the `DaoFlowAPI` client.

### Creating a Client

```go
import "github.com/Corphon/daoflow/api"

// Create client with default configuration
client, err := api.NewDaoFlowAPI(nil)
if err != nil {
    // Handle initialization error
}
defer client.Close()

// Create client with custom configuration
client, err := api.NewDaoFlowAPI(&api.Options{
    SystemConfig: &api.SystemConfig{
        Capacity: 2000.0,
        Threshold: 0.7,
    },
    LogLevel: "debug",
})
```
### Client Options
```go
type Options struct {
    // Core system configuration
    SystemConfig *SystemConfig
    
    // Logging configuration
    LogLevel string     // "debug", "info", "warn", "error"
    
    // API behavior configuration
    MaxRetries int      // Max retry attempts for operations
    Timeout time.Duration  // Default operation timeout
}
```
### Lifecycle Management
LifecycleAPI handles system initialization, startup, shutdown, and status management.
```go
// Initialize the system
if err := client.Lifecycle().Initialize(); err != nil {
    // Handle initialization error
}

// Start the system
if err := client.Lifecycle().Start(); err != nil {
    // Handle startup error
}

// Check if system is running
if running := client.Lifecycle().IsRunning(); running {
    // System is running
}

// Get current system state
state, err := client.Lifecycle().GetState()
if err != nil {
    // Handle error
}
fmt.Printf("System state: %s, Energy level: %f\n", state.Status, state.Energy)

// Stop the system gracefully
if err := client.Lifecycle().Stop(); err != nil {
    // Handle shutdown error
}
```
### Pattern System
PatternAPI provides methods for pattern recognition, analysis, and management.
```go
// Configure pattern detection
config := api.PatternConfig{
    Sensitivity: 0.8,      // Recognition sensitivity (0.0-1.0)
    MinConfidence: 0.7,    // Minimum confidence threshold  
    MaxPatterns: 100,      // Maximum patterns to track
}

// Start pattern recognition
patterns, err := client.Pattern().StartRecognition(context.Background(), config)
if err != nil {
    // Handle error
}

// Process detected patterns
for pattern := range patterns {
    fmt.Printf("Pattern detected: %s (confidence: %.2f%%)\n", 
        pattern.ID, pattern.Strength * 100)
        
    // Access pattern properties
    for key, value := range pattern.Properties {
        fmt.Printf("- %s: %v\n", key, value)
    }
}

// Analyze a specific pattern
analysis, err := client.Pattern().AnalyzePattern(pattern)
if err != nil {
    // Handle error
}
fmt.Printf("Pattern stability: %.2f\n", analysis.Stability)
```
### Energy System
EnergyAPI manages the quantum-inspired energy distribution and flows within the system.
```go
// Configure energy distribution
distribution := api.EnergyDistribution{
    Pattern:   0.3,  // Energy allocated to pattern recognition
    Evolution: 0.3,  // Energy allocated to evolution processes
    Adaptation: 0.2, // Energy allocated to adaptation
    Reserve: 0.2,    // Energy kept in reserve
}

// Apply energy distribution
if err := client.Energy().Distribute(context.Background(), distribution); err != nil {
    // Handle error
}

// Get current energy metrics
metrics, err := client.Energy().GetMetrics(context.Background())
if err != nil {
    // Handle error
}
fmt.Printf("Total energy: %.2f\n", metrics.Total)
fmt.Printf("Efficiency: %.2f%%\n", metrics.Efficiency * 100)
fmt.Printf("Flow rate: %.2f\n", metrics.FlowRate)

// Adjust system energy level
if err := client.Energy().AdjustLevel(10.0); err != nil {
    // Handle error
}

// Create an energy flow between components
flow, err := client.Energy().CreateFlow("pattern", "evolution", 5.0)
if err != nil {
    // Handle error
}
```
### Evolution System
EvolutionAPI enables adaptive evolution and self-optimization capabilities.
```go
// Configure evolution parameters
evolutionConfig := api.EvolutionConfig{
    LearningRate: 0.1,     // Learning rate for adaptation
    Generations: 10,       // Number of evolution generations
    PopulationSize: 100,   // Size of the evolution population
    MutationRate: 0.05,    // Mutation probability
}

// Start evolution process
evolution := client.Evolution()
evolution.SetConfig(evolutionConfig)

// Subscribe to evolution progress
progress, err := evolution.Subscribe(context.Background())
if err != nil {
    // Handle error
}
for state := range progress {
    fmt.Printf("Generation %d: Fitness = %.2f\n", 
        state.Generation, state.Fitness)
    fmt.Printf("New features emerged: %v\n", state.EmergentProperties)
}

// Perform system optimization
params := api.OptimizationParams{
    MaxIterations: 100,
    Goals: api.OptimizationGoals{
        Targets: map[string]float64{
            "performance": 0.9,
            "stability": 0.8,
        },
        Weights: map[string]float64{
            "performance": 0.6,
            "stability": 0.4,
        },
    },
}
result, err := client.Evolution().Optimize(context.Background(), params)
if err != nil {
    // Handle error
}
fmt.Printf("Optimization result: %.2f\n", result.Score)
```
### Monitoring & Metrics
MetricsAPI provides access to system performance metrics and health information.
```go
// Get overall system metrics
metrics, err := client.Metrics().GetMetrics(context.Background())
if err != nil {
    // Handle error
}
fmt.Printf("Response time: %.2fms\n", metrics.ResponseTime.Average)
fmt.Printf("Throughput: %d events/sec\n", metrics.Throughput.EventsPerSecond)
fmt.Printf("CPU utilization: %.2f%%\n", metrics.ResourceUtilization.CPU)

// Get component-specific metrics
patternMetrics, err := client.Metrics().GetComponentMetrics(context.Background(), "pattern")
if err != nil {
    // Handle error
}
fmt.Printf("Pattern recognition rate: %.2f patterns/sec\n", 
    patternMetrics["recognition_rate"])

// Get system health
health, err := client.Metrics().GetSystemHealth(context.Background())
if err != nil {
    // Handle error
}
fmt.Printf("System health: %s (%.2f%%)\n", health.Status, health.Score * 100)

// Configure health checks
healthConfig := api.HealthConfig{
    CheckInterval: time.Second * 30,
    Thresholds: map[string]float64{
        "cpu_usage": 0.9,
        "memory_usage": 0.85,
        "error_rate": 0.01,
    },
}
client.Metrics().ConfigureHealthChecks(healthConfig)
```
### Event System
EventsAPI provides a real-time event notification system for system occurrences.
```go
// Configure event filtering
filter := api.EventFilter{
    Types: []api.EventType{
        api.EventPatternDetected,
        api.EventStateChange,
        api.EventEmergence,
    },
    Priority: api.PriorityHigh,
}

// Subscribe to events
events, err := client.Events().Subscribe(filter)
if err != nil {
    // Handle error
}

// Process events
for event := range events {
    switch event.Type {
    case api.EventPatternDetected:
        pattern := event.Payload.(api.Pattern)
        fmt.Printf("Pattern detected: %s\n", pattern.ID)
        
    case api.EventStateChange:
        state := event.Payload.(api.SystemState)
        fmt.Printf("System state changed to: %s\n", state.Phase)
        
    case api.EventEmergence:
        emergence := event.Payload.(api.EmergentPattern)
        fmt.Printf("Emergence detected with complexity: %.2f\n", emergence.Complexity)
    }
}

// Unsubscribe when done
client.Events().Unsubscribe(events)

// Register custom event handler
handler := func(event api.SystemEvent) {
    // Process event
    fmt.Printf("Event received: %s\n", event.Type)
}
client.Events().RegisterHandler(api.EventPatternDetected, handler)
```
### Configuration
DaoFlow configuration is handled through structured config objects that can be provided during client creation.
```go
// System configuration
type SystemConfig struct {
    // General system settings
    Capacity     float64         // System energy capacity
    Threshold    float64         // System activation threshold
    UpdateRate   time.Duration   // State update frequency
    
    // Component-specific settings
    Pattern      *PatternConfig  // Pattern recognition configuration
    Energy       *EnergyConfig   // Energy system configuration
    Evolution    *EvolutionConfig// Evolution system configuration
    Monitor      *MonitorConfig  // Monitoring configuration
}

// Create a new configuration
config := api.SystemConfig{
    Capacity: 5000.0,
    Threshold: 0.75,
    UpdateRate: time.Millisecond * 100,
    Pattern: &api.PatternConfig{
        Sensitivity: 0.8,
        MinConfidence: 0.7,
    },
    Energy: &api.EnergyConfig{
        DistributionRate: 0.1,
        FlowResistance: 0.05,
    },
}

// Create client with this configuration
client, err := api.NewDaoFlowAPI(&api.Options{
    SystemConfig: &config,
})
```
### Error Handling
DaoFlow uses a structured error system that provides detailed information about error types and contexts.
```go
// Error handling example
import "github.com/Corphon/daoflow/api"

func handleDaoFlowOperation() {
    // Attempt operation
    _, err := client.Pattern().StartRecognition(context.Background(), config)
    if err != nil {
        switch {
        case errors.Is(err, api.ErrNotInitialized):
            // System not initialized yet
            fmt.Println("Please initialize the system first")
            
        case errors.Is(err, api.ErrAlreadyRunning):
            // The operation is already running
            fmt.Println("Pattern recognition is already running")
            
        case errors.Is(err, api.ErrModelNotFound):
            // Model not found
            fmt.Println("Required model not found")
            
        default:
            // Handle other errors
            fmt.Printf("Operation failed: %v\n", err)
        }
    }
}
```

### Common error types:

| **Error Constant**         | **Description**                           |
|:---------------------------|:------------------------------------------|
| ⚠️ **ErrNotInitialized**   | System has not been initialized          |
| 🔄 **ErrAlreadyRunning**   | Operation is already in progress         |
| ⏹️ **ErrNotRunning**       | System is not running                    |
| 🔍 **ErrModelNotFound**    | Referenced model not found               |
| 🛑 **ErrInvalidConfig**    | Invalid configuration provided           |
| ⏳ **ErrTimeout**           | Operation timed out                      |
| 📉 **ErrResourceLimit**    | Resource limit exceeded                  |


### For extended error information, the error objects also contain:

 - Error codes for programmatic handling
 - Detailed messages for user feedback
 - Stack context for debugging

## Advanced Features
For more advanced usage scenarios, please refer to the Examples document and Performance Optimization guide. 
