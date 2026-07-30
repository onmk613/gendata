package action

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"sync"
	"time"

	"gendata/internal/core"
	mydriver "gendata/internal/driver"
)

// 读延迟样本上限，避免长时间运行时占用过多内存（QPS 仍按全量统计）
const maxReadSamples = 500000

// userIDPool 线程安全的 user_id 集合，供点查采样
type userIDPool struct {
	mu  sync.RWMutex
	ids []string
}

func newUserIDPool() *userIDPool {
	return &userIDPool{ids: make([]string, 0, 1024)}
}

func (p *userIDPool) add(id string) {
	p.mu.Lock()
	p.ids = append(p.ids, id)
	p.mu.Unlock()
}

func (p *userIDPool) random() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if len(p.ids) == 0 {
		return ""
	}
	return p.ids[rand.IntN(len(p.ids))]
}

func (p *userIDPool) size() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.ids)
}

// 全局读池（read / mixed 模式使用）
var idPool *userIDPool

// runReadOnly 纯读模式：先从表里加载 user_id，再并发做点查
func runReadOnly() error {
	if WriteConf.Duration <= 0 {
		return fmt.Errorf("--duration is required for read mode (e.g. --duration=60s)")
	}

	idPool = newUserIDPool()
	loadUserIDsFromDB(100000)
	if idPool.size() == 0 {
		return fmt.Errorf("no data in table to read; please run write mode first")
	}
	slog.Info("read_pool_loaded", slog.Int("user_ids", idPool.size()))

	ctx, cancel := runCtx()
	defer cancel()

	var (
		wg       sync.WaitGroup
		errMu    sync.Mutex
		firstErr error
	)
	for i := 0; i < WriteConf.ReadConcurrency; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			runReadWorker(ctx, i, &errMu, &firstErr)
		}(i)
	}
	wg.Wait()
	return firstErr
}

// runMixed 读写混合：先 seed 一批数据填充读池，然后并发读写
func runMixed() error {
	if WriteConf.Duration <= 0 {
		return fmt.Errorf("--duration is required for mixed mode (e.g. --duration=60s)")
	}

	idPool = newUserIDPool()

	// 预写一批数据，保证读池非空
	seed := WriteConf.BatchSize
	if seed < 1000 {
		seed = 1000
	}
	rows := core.GenerateDefaultTableData(seed)
	if err := mydriver.DB.Create(rows).Error; err != nil {
		return fmt.Errorf("mixed seed insert failed: %w", err)
	}
	for _, r := range rows {
		idPool.add(r.UserID)
	}
	slog.Info("mixed_seeded", slog.Int("rows", seed), slog.Int("pool", idPool.size()))

	ctx, cancel := runCtx()
	defer cancel()

	var (
		wg       sync.WaitGroup
		errMu    sync.Mutex
		firstErr error
	)

	// 写 worker
	for i := 0; i < WriteConf.Concurrency; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			runWriteWorker(ctx, i, &errMu, &firstErr, idPool)
		}(i)
	}
	// 读 worker
	for i := 0; i < WriteConf.ReadConcurrency; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			runReadWorker(ctx, i, &errMu, &firstErr)
		}(i)
	}
	wg.Wait()
	return firstErr
}

// runReadWorker 单个读 worker：循环按 user_id 点查
func runReadWorker(ctx context.Context, id int, errMu *sync.Mutex, firstErr *error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		uid := idPool.random()
		if uid == "" {
			// 池暂空，稍等
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Millisecond):
			}
			continue
		}

		var row core.DefaultTableRow
		start := time.Now()
		err := mydriver.DB.Where("user_id = ?", uid).First(&row).Error
		latency := time.Since(start)
		if err != nil {
			errMu.Lock()
			if *firstErr == nil {
				*firstErr = err
			}
			errMu.Unlock()
			return
		}
		logReadStats(id, latency)
	}
}

// loadUserIDsFromDB 从表中取一批现有 user_id 填充读池
func loadUserIDsFromDB(limit int) {
	var ids []string
	if err := mydriver.DB.Model(&core.DefaultTableRow{}).Limit(limit).Pluck("user_id", &ids).Error; err != nil {
		slog.Warn("load_user_ids_failed", slog.Any("error", err))
		return
	}
	for _, id := range ids {
		idPool.add(id)
	}
}
