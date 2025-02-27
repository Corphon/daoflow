// system/monitor/trace/recorder.go

package trace

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Corphon/daoflow/system/types"
)

// TraceRecord 追踪记录
type TraceRecord struct {
	ID        string            `json:"id"`
	Timestamp time.Time         `json:"timestamp"`
	TraceID   types.TraceID     `json:"trace_id"`
	SpanID    types.SpanID      `json:"span_id"`
	Type      string            `json:"type"`
	Data      interface{}       `json:"data"`
	Metadata  map[string]string `json:"metadata"`
}

// Recorder 记录器
type Recorder struct {
	mu sync.RWMutex

	// 配置
	config struct {
		StoragePath   string        // 存储路径
		RetentionDays time.Duration // 保留时间
		BatchSize     int           // 批处理大小
		FlushInterval time.Duration // 刷新间隔
		Compression   bool          // 是否压缩
		AsyncWrite    bool          // 异步写入
	}

	// 存储缓冲
	buffer struct {
		records []TraceRecord
		size    int64
	}

	// 存储统计
	stats struct {
		totalRecords int64
		totalSize    int64
		lastFlush    time.Time
		errors       []error
	}

	// 状态
	status struct {
		isRunning  bool
		isFlushing bool
	}

	// 通道
	recordChan chan TraceRecord
	flushChan  chan struct{}
}

// ----------------------------------------------------
// NewRecorder 创建新的记录器
func NewRecorder(config types.TraceConfig) *Recorder {
	r := &Recorder{
		recordChan: make(chan TraceRecord, config.BufferSize),
		flushChan:  make(chan struct{}, 1),
	}

	// 设置配置
	r.config.StoragePath = config.StoragePath
	r.config.RetentionDays = config.RetentionDays
	r.config.BatchSize = config.BatchSize
	r.config.FlushInterval = config.FlushInterval
	r.config.Compression = config.Compression
	r.config.AsyncWrite = config.AsyncWrite

	// 初始化缓冲
	r.buffer.records = make([]TraceRecord, 0, r.config.BatchSize)

	return r
}

// Start 启动记录器
func (r *Recorder) Start(ctx context.Context) error {
	r.mu.Lock()
	if r.status.isRunning {
		r.mu.Unlock()
		return types.NewSystemError(types.ErrRuntime, "recorder already running", nil)
	}
	r.status.isRunning = true
	r.mu.Unlock()

	// 启动处理循环
	go r.processLoop(ctx)

	return nil
}

// Stop 停止记录器
func (r *Recorder) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.status.isRunning {
		return nil
	}

	r.status.isRunning = false

	// 刷新剩余记录
	return r.flush()
}

// Record 记录追踪数据
func (r *Recorder) Record(record TraceRecord) error {
	if !r.status.isRunning {
		return types.NewSystemError(types.ErrRuntime, "recorder not running", nil)
	}

	select {
	case r.recordChan <- record:
		return nil
	default:
		return types.NewSystemError(types.ErrOverflow, "record buffer full", nil)
	}
}

// processLoop 处理循环
func (r *Recorder) processLoop(ctx context.Context) {
	ticker := time.NewTicker(r.config.FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			r.flush()
			return
		case record := <-r.recordChan:
			if err := r.processRecord(record); err != nil {
				r.recordError(err)
			}
		case <-ticker.C:
			if err := r.flush(); err != nil {
				r.recordError(err)
			}
		case <-r.flushChan:
			if err := r.flush(); err != nil {
				r.recordError(err)
			}
		}
	}
}

// processRecord 处理单条记录
func (r *Recorder) processRecord(record TraceRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 添加到缓冲
	r.buffer.records = append(r.buffer.records, record)
	r.buffer.size += r.estimateRecordSize(record)

	// 检查是否需要刷新
	if len(r.buffer.records) >= r.config.BatchSize {
		return r.flush()
	}

	return nil
}

// flush 刷新缓冲区
func (r *Recorder) flush() error {
	r.mu.Lock()
	if r.status.isFlushing || len(r.buffer.records) == 0 {
		r.mu.Unlock()
		return nil
	}
	r.status.isFlushing = true
	records := r.buffer.records
	r.buffer.records = make([]TraceRecord, 0, r.config.BatchSize)
	r.buffer.size = 0
	r.mu.Unlock()

	// 写入存储
	if err := r.writeRecords(records); err != nil {
		return err
	}

	r.mu.Lock()
	r.stats.totalRecords += int64(len(records))
	r.stats.lastFlush = time.Now()
	r.status.isFlushing = false
	r.mu.Unlock()

	return nil
}

// writeRecords 写入记录到存储
func (r *Recorder) writeRecords(records []TraceRecord) error {
	// 按日期组织文件路径
	path := r.generateStoragePath(time.Now())

	// 序列化记录
	data, err := r.serializeRecords(records)
	if err != nil {
		return err
	}

	// 压缩数据
	if r.config.Compression {
		data, err = r.compressData(data)
		if err != nil {
			return err
		}
	}

	// 异步写入
	if r.config.AsyncWrite {
		go func() {
			if err := r.writeToStorage(path, data); err != nil {
				r.recordError(err)
			}
		}()
		return nil
	}

	// 同步写入
	return r.writeToStorage(path, data)
}

// serializeRecords 序列化记录
func (r *Recorder) serializeRecords(records []TraceRecord) ([]byte, error) {
	return json.Marshal(records)
}

// compressData 压缩数据
func (r *Recorder) compressData(data []byte) ([]byte, error) {
	// 创建缓冲区
	var buf bytes.Buffer

	// 创建gzip writer
	gw := gzip.NewWriter(&buf)

	// 写入数据
	if _, err := gw.Write(data); err != nil {
		return nil, err
	}

	// 关闭writer确保所有数据都被写入
	if err := gw.Close(); err != nil {
		return nil, err
	}

	// 返回压缩后的数据
	return buf.Bytes(), nil
}

// writeToStorage 写入存储实现
func (r *Recorder) writeToStorage(path string, data []byte) error {
	// 创建目录
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory failed: %v", err)
	}

	// 创建或打开文件
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("open file failed: %v", err)
	}
	defer file.Close()

	// 写入数据
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("write data failed: %v", err)
	}

	// 强制刷新到磁盘
	if err := file.Sync(); err != nil {
		return fmt.Errorf("sync data failed: %v", err)
	}

	return nil
}

// generateStoragePath 生成存储路径
func (r *Recorder) generateStoragePath(t time.Time) string {
	// 基于时间生成存储路径
	return fmt.Sprintf("%s/%s/%s.trace",
		r.config.StoragePath,
		t.Format("2006/01/02"),
		t.Format("15-04-05"))
}

// estimateRecordSize 估算记录大小
func (r *Recorder) estimateRecordSize(record TraceRecord) int64 {
	// 简单估算记录大小
	data, _ := json.Marshal(record)
	return int64(len(data))
}

// cleanOldRecords 清理旧记录
func (r *Recorder) cleanOldRecords() error {
	cutoff := time.Now().Add(-r.config.RetentionDays)

	// 遍历存储目录
	err := filepath.Walk(r.config.StoragePath, func(path string, info os.FileInfo, err error) error {
		// 错误处理
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 检查文件后缀
		if !strings.HasSuffix(info.Name(), ".trace") {
			return nil
		}

		// 解析文件日期
		fileTime, err := time.Parse("2006/01/02/15-04-05.trace",
			strings.TrimPrefix(path, r.config.StoragePath+"/"))
		if err != nil {
			return nil // 跳过无法解析的文件
		}

		// 删除过期文件
		if fileTime.Before(cutoff) {
			if err := os.Remove(path); err != nil {
				r.recordError(fmt.Errorf("failed to remove old trace file %s: %v", path, err))
				return nil // 继续处理其他文件
			}

			r.mu.Lock()
			r.stats.totalRecords-- // 更新统计
			r.mu.Unlock()
		}

		return nil
	})

	if err != nil {
		return types.NewSystemError(types.ErrStorage, "failed to clean old records", err)
	}

	return nil
}

// recordError 记录错误
func (r *Recorder) recordError(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.stats.errors = append(r.stats.errors, err)
}

// GetStats 获取统计信息
func (r *Recorder) GetStats() struct {
	TotalRecords int64
	TotalSize    int64
	LastFlush    time.Time
	ErrorCount   int
} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return struct {
		TotalRecords int64
		TotalSize    int64
		LastFlush    time.Time
		ErrorCount   int
	}{
		TotalRecords: r.stats.totalRecords,
		TotalSize:    r.stats.totalSize,
		LastFlush:    r.stats.lastFlush,
		ErrorCount:   len(r.stats.errors),
	}
}

// GetRecords 获取记录数据
func (r *Recorder) GetRecords() []TraceRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 返回当前缓冲区的记录
	records := make([]TraceRecord, len(r.buffer.records))
	copy(records, r.buffer.records)
	return records
}
