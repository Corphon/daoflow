
# DaoFlow API Reference

## Table of Contents
- [Overview](#overview)
- [Core APIs](#core-apis)
  - [LifecycleAPI](#lifecycleapi)
  - [EnergyAPI](#energyapi)
  - [EvolutionAPI](#evolutionapi)
  - [PatternAPI](#patternapi)
  - [MetricsAPI](#metricsapi)
  - [ConfigAPI](#configapi)
  - [EventsAPI](#eventsapi)
  - [HealthAPI](#healthapi)
- [Common Types](#common-types)
- [Error Handling](#error-handling)
- [Best Practices](#best-practices)

## Overview

DaoFlow API version: v2.0.0

The DaoFlow framework provides a comprehensive set of APIs for building self-adaptive distributed systems based on Eastern philosophy principles.

### Installation

```bash
go get github.com/Corphon/daoflow
```

### Basic Usage

```go
import "github.com/Corphon/daoflow/api"
```

## Core APIs

### LifecycleAPI

Manages the system lifecycle and state transitions.

#### Methods

```go
type LifecycleAPI interface {
    // Initialize prepares the system for operation
    Initialize() error
    
    // Start begins system operation
    Start() error
    
    // Stop gracefully stops the system
    Stop() error
    
    // GetStatus returns current system status
    GetStatus() (*LifecycleStatus, error)
    
    // Subscribe to lifecycle events
    Subscribe() (<-chan LifecycleEvent, error)
}

type LifecycleStatus struct {
    State       LifecycleState     `json:"state"`
    Uptime      time.Duration      `json:"uptime"`
    StartTime   time.Time          `json:"start_time"`
    Components  map[string]Status  `json:"components"`
}
```

### EnergyAPI

Manages system energy distribution and balance.

#### Methods

```go
type EnergyAPI interface {
    // Configure sets energy system parameters
    Configure(ctx context.Context, config EnergyConfig) error
    
    // Distribute allocates energy across system components
    Distribute(ctx context.Context, dist EnergyDistribution) error
    
    // GetMetrics retrieves energy metrics
    GetMetrics(ctx context.Context) (*EnergyMetrics, error)
    
    // Balance performs energy rebalancing
    Balance(ctx context.Context) error
}

type EnergyDistribution struct {
    Pattern     float64  `json:"pattern"`     // Pattern recognition energy
    Evolution   float64  `json:"evolution"`   // Evolution process energy
    Adaptation  float64  `json:"adaptation"`  // Adaptation energy
    Reserve     float64  `json:"reserve"`     // Energy reserve
}
```

### EvolutionAPI

Controls system evolution and adaptation.

#### Methods

```go
type EvolutionAPI interface {
    // TriggerEvolution initiates an evolution cycle
    TriggerEvolution(mode string) error
    
    // GetEvolutionStatus retrieves current evolution state
    GetEvolutionStatus() (*EvolutionStatus, error)
    
    // Configure sets evolution parameters
    Configure(config EvolutionConfig) error
}

type EvolutionConfig struct {
    LearningRate    float64       `json:"learning_rate"`
    Generations     int           `json:"generations"`
    PopulationSize  int           `json:"population_size"`
    MutationRate    float64       `json:"mutation_rate"`
}
```

### PatternAPI

Handles pattern recognition and analysis.

#### Methods

```go
type PatternAPI interface {
    // StartRecognition begins pattern detection
    StartRecognition(config PatternConfig) (<-chan Pattern, error)
    
    // GetPatterns retrieves detected patterns
    GetPatterns() ([]Pattern, error)
    
    // AnalyzePattern performs deep pattern analysis
    AnalyzePattern(pattern Pattern) (*PatternAnalysis, error)
}

type Pattern struct {
    ID          string                 `json:"id"`
    Confidence  float64               `json:"confidence"`
    Properties  map[string]float64    `json:"properties"`
    Timestamp   time.Time             `json:"timestamp"`
}
```

### MetricsAPI

Provides system metrics and monitoring capabilities.

#### Methods

```go
type MetricsAPI interface {
    // RegisterMetric registers a new metric
    RegisterMetric(name string, typ MetricType, desc string) error
    
    // RecordMetric records a metric value
    RecordMetric(name string, value float64, labels map[string]string) error
    
    // QueryMetrics queries metrics based on criteria
    QueryMetrics(query map[string]interface{}) ([]*MetricSeries, error)
}

type MetricType string

const (
    TypeGauge     MetricType = "gauge"
    TypeCounter   MetricType = "counter"
    TypeHistogram MetricType = "histogram"
)
```

### ConfigAPI

Manages system configuration.

#### Methods

```go
type ConfigAPI interface {
    // SetConfig sets a configuration value
    SetConfig(key string, value interface{}, scope ConfigScope) error
    
    // GetConfig retrieves a configuration value
    GetConfig(key string) (*ConfigValue, error)
    
    // Subscribe to configuration changes
    Subscribe() (<-chan ConfigEvent, error)
}

type ConfigScope string

const (
    ScopeGlobal    ConfigScope = "global"
    ScopeComponent ConfigScope = "component"
)
```

### EventsAPI

Handles system event processing.

#### Methods

```go
type EventsAPI interface {
    // Publish publishes an event
    Publish(evt Event) error
    
    // Subscribe subscribes to events
    Subscribe(filter EventFilter) (*Subscription, error)
    
    // GetEvents retrieves historical events
    GetEvents(filter EventFilter) ([]*Event, error)
}

type EventPriority int

const (
    PriorityLow     EventPriority = 0
    PriorityNormal  EventPriority = 1
    PriorityHigh    EventPriority = 2
    PriorityCritical EventPriority = 3
)
```

### HealthAPI

Monitors system health.

#### Methods

```go
type HealthAPI interface {
    // GetSystemHealth retrieves overall system health
    GetSystemHealth() (*SystemHealth, error)
    
    // RegisterCheck registers a health check
    RegisterCheck(check *HealthCheck) error
    
    // StartCheck starts a specific health check
    StartCheck(checkID string) error
}

type HealthStatus string

const (
    StatusHealthy    HealthStatus = "healthy"
    StatusDegraded   HealthStatus = "degraded"
    StatusUnhealthy  HealthStatus = "unhealthy"
)
```

## Common Types

### Error Types

```go
var (
    ErrConfigNotFound     = errors.New("configuration not found")
    ErrInvalidConfig      = errors.New("invalid configuration")
    ErrMetricNotFound     = errors.New("metric not found")
    ErrSubscriptionFailed = errors.New("subscription failed")
)
```

### Context Usage

All long-running operations accept a context.Context parameter for cancellation and timeout control.

## Error Handling

DaoFlow uses error wrapping for detailed error information:

```go
if err != nil {
    var configErr *ConfigError
    if errors.As(err, &configErr) {
        // Handle configuration error
    }
    return fmt.Errorf("operation failed: %w", err)
}
```

## Best Practices

1. **Resource Management**
   - Always close resources using defer
   - Use context for operation control
   - Handle error cases appropriately

2. **Configuration**
   - Use appropriate scopes for configurations
   - Validate configurations before applying
   - Monitor configuration changes

3. **Event Handling**
   - Use appropriate event priorities
   - Implement event filtering
   - Handle subscription cleanup

4. **Health Monitoring**
   - Register health checks for critical components
   - Set appropriate check intervals
   - Monitor system health metrics

5. **Performance Optimization**
   - Use appropriate buffer sizes
   - Implement proper error handling
   - Monitor system metrics

For more detailed examples, please refer to the [Examples](./examples.md) document.
