package main

import (
	"github.com/darabuchi/prism"
	"github.com/darabuchi/prism/pkg/geoip"
	"github.com/lazygophers/log"
)

// BatchInserter 批量插入器，减少单次插入的开销
type BatchInserter struct {
	writer    *geoip.Writer
	batch     []BatchItem
	batchSize int
	count     int
}

// BatchItem 批量插入项
type BatchItem struct {
	StartIP string
	EndIP   string
	Geo     *prism.GeoIP
}

// NewBatchInserter 创建批量插入器
func NewBatchInserter(writer *geoip.Writer, batchSize int) *BatchInserter {
	if batchSize <= 0 {
		batchSize = 1000 // 默认批大小
	}
	return &BatchInserter{
		writer:    writer,
		batch:     make([]BatchItem, 0, batchSize),
		batchSize: batchSize,
	}
}

// Add 添加到批次
func (b *BatchInserter) Add(startIP, endIP string, geo *prism.GeoIP) error {
	b.batch = append(b.batch, BatchItem{
		StartIP: startIP,
		EndIP:   endIP,
		Geo:     geo,
	})

	// 达到批大小时自动刷新
	if len(b.batch) >= b.batchSize {
		return b.Flush()
	}

	return nil
}

// Flush 刷新批次，执行实际插入
func (b *BatchInserter) Flush() error {
	if len(b.batch) == 0 {
		return nil
	}

	// 批量插入
	for _, item := range b.batch {
		if err := b.writer.InsertGeoIPRange(item.StartIP, item.EndIP, item.Geo); err != nil {
			log.Errorf("批量插入失败: %s-%s, 错误: %v", item.StartIP, item.EndIP, err)
			return err
		}
		b.count++
	}

	// 清空批次，复用底层数组
	b.batch = b.batch[:0]

	return nil
}

// Count 返回已插入的总数
func (b *BatchInserter) Count() int {
	return b.count
}
