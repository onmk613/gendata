package action

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"gendata/internal/core"
	mydriver "gendata/internal/driver"
)

func Run(driverName string) error {
	// 计算表列数
	b := core.GenerateDefaultTableData(1)
	count, err := calcBatchCount(b, driverName)
	if err != nil {
		return err
	}

	// 连接数据库
	if err := mydriver.GetDriver(driverName, count, WriteConf.Concurrency); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	if mydriver.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	slog.Info("Configuration",
		slog.Int("concurrency", WriteConf.Concurrency),
		slog.Int("batch_size", WriteConf.BatchSize),
		slog.Int("repeat_count", WriteConf.RepeatCount),
		slog.Int("batch_size", count),
	)

	// 自动创建表
	if err := mydriver.DB.AutoMigrate(&core.DefaultTableRow{}); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// 更新数据
	err = genDefaultDataAndInsert()
	return err
}

func genDefaultDataAndInsert() error {
	var wg sync.WaitGroup
	var errMu sync.Mutex
	var firstErr error

	// 并发控制
	for i := 0; i < WriteConf.Concurrency; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// 写入次数控制
			for j := 0; j < WriteConf.RepeatCount; j++ {
				// 生成数据
				a := time.Now()
				rows := core.GenerateDefaultTableData(WriteConf.BatchSize)
				b := time.Since(a)
				slog.Info("gen_data_time", slog.Any("time", b))

				startInsert := time.Now()
				if err := mydriver.DB.Create(rows).Error; err != nil {
					errMu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					errMu.Unlock()
					return
				}
				insertTime := time.Since(startInsert)

				// 统计
				logBatchStats(i, j, insertTime, WriteConf.BatchSize)
			}
		}(i)
	}

	wg.Wait()
	return firstErr
}

func CloseGendata() {
	mydriver.CloseDB()
	printThroughputStats()
}
