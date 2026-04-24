package action

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"gendata/pkg/core"
	mydriver "gendata/pkg/driver"
)

func Run(driverName string) error {
	if err := mydriver.GetDriver(driverName); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	if mydriver.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	slog.Info("Starting data generation and insertion")
	slog.Info("Configuration",
		slog.Int("concurrency", WriteConf.Concurrency),
		slog.Int("batch_size", WriteConf.BatchSize),
		slog.Int("repeat_Count", WriteConf.RepeatCount),
	)

	if mydriver.SqlConf.Table == "" {
		mydriver.SqlConf.Table = "gendata_table"
	}

	mydriver.DB.AutoMigrate(&mydriver.DefaultTableRow{})

	err := genDefaultDataAndInsert()
	return err
}

func genDefaultDataAndInsert() error {
	var wg sync.WaitGroup

	for i := 0; i < WriteConf.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < WriteConf.RepeatCount; j++ {
				rows := core.GenerateDefaultTableData(WriteConf.BatchSize)

				startInsert := time.Now()
				err := mydriver.CreateTemplateDataBatchWithSize(rows, mydriver.DB, WriteConf.BatchSize)
				insertTime := time.Since(startInsert)

				if err != nil {
					slog.Warn("Write batch failed", slog.Any("err", err))
					continue
				}

				go logBatchStats(i, j, insertTime, WriteConf.BatchSize)
			}
		}()
	}

	wg.Wait()
	return nil
}

func CloseGendata() {
	mydriver.CloseDB()
	printThroughputStats()
}
