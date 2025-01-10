// system/backpressure.go
package system

import (
	"errors"
	"math/rand/v2"
	"time"
)

type BackpressureHandler struct {
	dropPolicy    DropPolicy
	retryStrategy RetryStrategy
	metrics       *QueueMetrics
}

type DropPolicy int

const (
	DropOldest DropPolicy = iota
	DropNewest
	DropRandom
)

type RetryStrategy struct {
	MaxRetries    int
	RetryInterval time.Duration
	BackoffFactor float64
}

type QueueMetrics struct {
	DroppedCount   int64
	RetryCount     int64
	SuccessRate    float64
	AverageLatency time.Duration
}

func NewBackpressureHandler(policy DropPolicy, strategy RetryStrategy) *BackpressureHandler {
	return &BackpressureHandler{
		dropPolicy:    policy,
		retryStrategy: strategy,
		metrics:       &QueueMetrics{},
	}
}

func (bh *BackpressureHandler) handleOverflow(event interface{}) error {
	switch bh.dropPolicy {
	case DropOldest:
		return bh.dropOldestAndRetry(event)
	case DropNewest:
		return bh.dropNewestAndRetry(event)
	default:
		return bh.dropRandomAndRetry(event)
	}
}

// dropOldestAndRetry 丢弃最旧的事件并重试当前事件
func (bh *BackpressureHandler) dropOldestAndRetry(event interface{}) error {
	retries := 0
	interval := bh.retryStrategy.RetryInterval

	for retries < bh.retryStrategy.MaxRetries {
		// 尝试处理事件
		if err := bh.tryProcessEvent(event); err == nil {
			bh.metrics.SuccessRate = float64(bh.metrics.RetryCount) / float64(bh.metrics.DroppedCount+bh.metrics.RetryCount)
			return nil
		}

		// 增加重试计数
		retries++
		bh.metrics.RetryCount++

		// 使用退避策略计算下次重试间隔
		interval = time.Duration(float64(interval) * bh.retryStrategy.BackoffFactor)

		// 等待后重试
		time.Sleep(interval)
	}

	// 达到最大重试次数，丢弃最旧的事件
	bh.metrics.DroppedCount++
	return errors.New("max retries exceeded, dropped oldest event")
}

// dropNewestAndRetry 丢弃最新的事件并保留当前处理中的事件
func (bh *BackpressureHandler) dropNewestAndRetry(event interface{}) error {
	// 直接丢弃新事件
	bh.metrics.DroppedCount++

	// 更新指标
	startTime := time.Now()
	bh.metrics.AverageLatency = time.Since(startTime)

	return errors.New("queue full, dropped newest event")
}

// dropRandomAndRetry 随机丢弃事件并重试
func (bh *BackpressureHandler) dropRandomAndRetry(event interface{}) error {
	// 生成随机数决定是否丢弃当前事件
	if rand.Float64() < 0.5 {
		bh.metrics.DroppedCount++
		return errors.New("randomly selected to drop event")
	}

	// 尝试重试当前事件
	retries := 0
	interval := bh.retryStrategy.RetryInterval

	for retries < bh.retryStrategy.MaxRetries {
		if err := bh.tryProcessEvent(event); err == nil {
			bh.metrics.SuccessRate = float64(bh.metrics.RetryCount) / float64(bh.metrics.DroppedCount+bh.metrics.RetryCount)
			return nil
		}

		retries++
		bh.metrics.RetryCount++
		interval = time.Duration(float64(interval) * bh.retryStrategy.BackoffFactor)
		time.Sleep(interval)
	}

	bh.metrics.DroppedCount++
	return errors.New("max retries exceeded in random drop strategy")
}

// tryProcessEvent 尝试处理事件
func (bh *BackpressureHandler) tryProcessEvent(event interface{}) error {
	// 这里添加实际的事件处理逻辑
	// 例如：将事件发送到队列或处理管道
	startTime := time.Now()

	// TODO: 实现实际的事件处理逻辑

	// 更新处理延迟指标
	bh.metrics.AverageLatency = time.Since(startTime)

	return nil
}
