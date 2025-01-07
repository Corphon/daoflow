# DaoFlow API Structure

## 🔄 LifecycleAPI
├── Initialize()
├── Start()
├── Stop()
├── GetStatus()
└── Subscribe()

## ⚡ EnergyAPI
├── Configure()
├── Distribute()
│   ├── Pattern Energy
│   ├── Evolution Energy
│   ├── Adaptation Energy
│   └── Reserve Energy
├── GetMetrics()
├── Balance()
└── Subscribe()

## 🧬 EvolutionAPI
├── TriggerEvolution()
├── GetEvolutionStatus()
├── WuXingScheduler
│   ├── Optimize()
│   └── GetMetrics()
└── YinYangBalance
    ├── Adjust()
    └── GetStatus()

## 📊 PatternAPI
├── StartRecognition()
├── BaGuaDetector
│   ├── Configure()
│   └── Subscribe()
├── GetPatterns()
└── AnalyzePattern()

## 📈 MetricsAPI
├── RegisterMetric()
├── RecordMetric()
├── GetMetric()
├── GetMetricValue()
├── QueryMetrics()
└── Subscribe()

## ⚙️ ConfigAPI
├── SetConfig()
├── GetConfig()
├── DeleteConfig()
├── GetConfigsByScope()
├── SetValidation()
├── GetHistory()
├── Subscribe()
├── Export()
└── Import()

## 🔔 EventsAPI
├── Publish()
├── Subscribe()
├── Unsubscribe()
├── GetEvents()
└── GetStats()

## 💓 HealthAPI
├── RegisterCheck()
├── StartCheck()
├── StopCheck()
├── GetComponentHealth()
├── GetSystemHealth()
└── Subscribe()

## 🌐 SystemAdapter
├── Core Components
│   ├── SystemCore
│   └── LegacySupport
└── Mappers
    ├── EventMapper
    └── BufferMapper

## 💾 BufferSystem
├── DynamicBuffer
│   ├── Capacity Management
│   └── Threshold Control
└── BufferMetrics
    ├── Utilization
    ├── DropRate
    └── Latency
