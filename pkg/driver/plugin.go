package driver

import (
	"gorm.io/gorm"
)

// BatchPlugin GORM 插件，注册后自动设置 CreateInBatches 大小
type BatchPlugin struct {
	batchSize int
}

// NewBatchPlugin 创建批量插入插件
func NewBatchPlugin(size int) *BatchPlugin {
	return &BatchPlugin{batchSize: size}
}

// Name 插件名称（实现 gorm.Plugin 接口）
func (bp *BatchPlugin) Name() string {
	return "BatchPlugin"
}

// Initialize 初始化插件（实现 gorm.Plugin 接口）
func (bp *BatchPlugin) Initialize(db *gorm.DB) error {
	return db.Callback().Create().Before("gorm:create").Register("batch_plugin:set_batch_size", func(tx *gorm.DB) {
		if tx.CreateBatchSize <= 0 {
			tx.CreateBatchSize = bp.batchSize
		}
	})
}
