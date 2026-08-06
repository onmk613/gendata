package action

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"gendata/internal/core"
	mydriver "gendata/internal/driver"
	"gendata/internal/generator"
)

func Run(driverName string) error {
	// 规范化配置
	if WriteConf.Mode == "" {
		WriteConf.Mode = ModeWrite
	}
	if WriteConf.ReadConcurrency <= 0 {
		WriteConf.ReadConcurrency = WriteConf.Concurrency
	}
	if err := validateConfig(); err != nil {
		return err
	}

	// 计算表列数
	b := core.GenerateDefaultTableData(1, nil)
	count, err := calcBatchCount(b, driverName)
	if err != nil {
		return err
	}

	// 连接数据库（连接池大小 = 读+写 worker 总数）
	if err := mydriver.GetDriver(driverName, count, totalConnections()); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	if mydriver.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	slog.Info("Configuration",
		slog.String("mode", WriteConf.Mode),
		slog.Int("concurrency", WriteConf.Concurrency),
		slog.Int("read_concurrency", WriteConf.ReadConcurrency),
		slog.Int("batch_size", WriteConf.BatchSize),
		slog.Int("repeat_count", WriteConf.RepeatCount),
		slog.Duration("duration", WriteConf.Duration),
		slog.Int("calc_batch_size", count),
	)

	// 自动创建表（纯读模式只查询，不建表、不改表）
	if WriteConf.Mode != ModeRead {
		if err := mydriver.DB.AutoMigrate(&core.DefaultTableRow{}); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	// 按模式分流
	switch WriteConf.Mode {
	case ModeRead:
		return runReadOnly()
	case ModeMixed:
		return runMixed()
	default:
		return runWrite()
	}
}

// validateConfig 在连接数据库之前校验参数，避免拼写错误被静默当成 write 执行。
func validateConfig() error {
	switch WriteConf.Mode {
	case ModeWrite, ModeRead, ModeMixed:
	default:
		return fmt.Errorf("unsupported --mode %q (must be write, read or mixed)", WriteConf.Mode)
	}
	if WriteConf.BatchSize <= 0 {
		return fmt.Errorf("--batchsize must be > 0, got %d", WriteConf.BatchSize)
	}
	if WriteConf.Concurrency <= 0 {
		return fmt.Errorf("--concurrency must be > 0, got %d", WriteConf.Concurrency)
	}
	if WriteConf.Mode == ModeRead || WriteConf.Mode == ModeMixed {
		if WriteConf.Duration <= 0 {
			return fmt.Errorf("--duration is required for %s mode (e.g. --duration=60s)", WriteConf.Mode)
		}
		if WriteConf.ReadConcurrency <= 0 {
			return fmt.Errorf("--readconcurrency must be > 0, got %d", WriteConf.ReadConcurrency)
		}
	}
	if WriteConf.Mode == ModeWrite && WriteConf.Duration <= 0 && WriteConf.RepeatCount < 0 {
		return fmt.Errorf("--repeatcount must be >= 0, got %d", WriteConf.RepeatCount)
	}
	return nil
}

// totalConnections 按模式返回实际需要的连接池大小（保证基准可比）
func totalConnections() int {
	var n int
	switch WriteConf.Mode {
	case ModeRead:
		n = WriteConf.ReadConcurrency
	case ModeMixed:
		n = WriteConf.Concurrency + WriteConf.ReadConcurrency
	default: // write
		n = WriteConf.Concurrency
	}
	if n < 1 {
		n = 1
	}
	return n
}

// runCtx 构造运行 context；duration>0 时带超时
func runCtx() (context.Context, context.CancelFunc) {
	if WriteConf.Duration > 0 {
		return context.WithTimeout(context.Background(), WriteConf.Duration)
	}
	return context.WithCancel(context.Background())
}

// runWrite 纯写模式（向后兼容：用 repeatcount 或 duration 控制）
func runWrite() error {
	ctx, cancel := runCtx()
	defer cancel()

	var (
		wg       sync.WaitGroup
		errMu    sync.Mutex
		firstErr error
	)
	for i := 0; i < WriteConf.Concurrency; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			runWriteWorker(ctx, i, &errMu, &firstErr, nil)
		}(i)
	}
	wg.Wait()
	return firstErr
}

// runWriteWorker 单个写 worker 的主循环
// pool 非 nil 时，把写入的 user_id 收集进池（供 mixed 模式读使用）
func runWriteWorker(ctx context.Context, id int, errMu *sync.Mutex, firstErr *error, pool *userIDPool) {
	rng := generator.NewRandom(uint64(time.Now().UnixNano()) + uint64(id)*2654435761)
	for j := 0; ; j++ {
		select {
		case <-ctx.Done():
			return
		default:
		}
		// duration 未设置时用 repeatcount 控制次数
		if WriteConf.Duration <= 0 && j >= WriteConf.RepeatCount {
			return
		}

		a := time.Now()
		rows := core.GenerateDefaultTableData(WriteConf.BatchSize, rng)
		slog.Info("gen_data_time", slog.Any("time", time.Since(a)))

		startInsert := time.Now()
		if err := mydriver.InsertRows(ctx, rows); err != nil {
			errMu.Lock()
			if *firstErr == nil {
				*firstErr = err
			}
			errMu.Unlock()
			return
		}
		insertTime := time.Since(startInsert)

		if pool != nil {
			userIDs := make([]string, len(rows))
			for i, r := range rows {
				userIDs[i] = r.UserID
			}
			pool.addAll(userIDs)
		}

		logBatchStats(id, j, insertTime, WriteConf.BatchSize)
	}
}

func CloseGendata() {
	mydriver.CloseDB()
	printThroughputStats()
}
