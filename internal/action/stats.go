package action

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

// 运行模式
const (
	ModeWrite = "write" // 纯写（默认，向后兼容）
	ModeRead  = "read"  // 纯读
	ModeMixed = "mixed" // 读写混合
)

var (
	WriteConf WriteConfiguration
	statsMu   sync.Mutex

	// 写统计
	totalTime          time.Duration
	meanThroughputPerS float64
	maxThroughputPerS  float64
	minThroughputPerS  float64
	writeLatencies     []time.Duration
	writeRows          int

	// 读统计
	readTotalTime time.Duration
	readWallTime  time.Duration
	readCount     int64
	readLatencies []time.Duration
)

func init() {
	minThroughputPerS = math.MaxFloat64
	writeLatencies = make([]time.Duration, 0, 1024)
	readLatencies = make([]time.Duration, 0, 1024)
}

// 运行配置
type WriteConfiguration struct {
	Concurrency     int
	BatchSize       int
	RepeatCount     int
	Mode            string        // write | read | mixed
	ReadConcurrency int           // 读 worker 数（<=0 时在 Run 中置为 Concurrency）
	Duration        time.Duration // >0 时按持续时间运行，忽略 RepeatCount
}

// 记录一次写 batch 的耗时与吞吐
func logBatchStats(i, j int, insertTime time.Duration, batchSize int) {
	statsMu.Lock()
	defer statsMu.Unlock()

	totalTime += insertTime
	writeLatencies = append(writeLatencies, insertTime)
	writeRows += batchSize

	insertThroughputPerS := float64(batchSize) / insertTime.Seconds()

	if insertThroughputPerS < minThroughputPerS {
		minThroughputPerS = insertThroughputPerS
	}
	if insertThroughputPerS > maxThroughputPerS {
		maxThroughputPerS = insertThroughputPerS
	}

	fmt.Printf("Completed: Con = %d, Count = %d, InsertTime = %.3fs, PerS = %.0f row/s\n",
		i, j, insertTime.Seconds(), insertThroughputPerS)
}

// 记录一次读查询的耗时
func logReadStats(_ int, latency time.Duration) {
	statsMu.Lock()
	defer statsMu.Unlock()

	readTotalTime += latency
	readCount++
	// 超过样本上限后停止收集延迟（QPS 仍按全量统计），控制内存
	if len(readLatencies) < maxReadSamples {
		readLatencies = append(readLatencies, latency)
	}
}

// 输出汇总统计（写 + 读 + 延迟分位数）
func printThroughputStats() {
	statsMu.Lock()
	defer statsMu.Unlock()

	hasWrite := totalTime.Seconds() > 0
	hasRead := readCount > 0

	if hasWrite {
		meanThroughputPerS = float64(writeRows) / totalTime.Seconds()

		minVal := minThroughputPerS
		if minVal == math.MaxFloat64 {
			minVal = 0
		}

		fmt.Printf("\n=== Write ===\n")
		fmt.Printf("Total: Time = %.3fs, Rows = %d, Mean = %.0f row/s, Max = %.0f row/s, Min = %.0f row/s\n",
			totalTime.Seconds(), writeRows, meanThroughputPerS, maxThroughputPerS, minVal)
		printLatency("BatchLatency", writeLatencies)
	}

	if hasRead {
		wall := readWallTime.Seconds()
		if wall <= 0 {
			wall = readTotalTime.Seconds()
		}
		if wall <= 0 {
			wall = 1
		}
		qps := float64(readCount) / wall
		fmt.Printf("\n=== Read ===\n")
		fmt.Printf("Total: Time = %.3fs, Queries = %d, QPS = %.0f\n",
			wall, readCount, qps)
		printLatency("QueryLatency", readLatencies)
	}

	if !hasWrite && !hasRead {
		// 兼容旧逻辑：无数据直接返回
		return
	}
}

// printLatency 计算并打印 P50/P95/P99
func printLatency(label string, latencies []time.Duration) {
	if len(latencies) == 0 {
		return
	}
	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	pick := func(p float64) time.Duration {
		if len(sorted) == 1 {
			return sorted[0]
		}
		i := int(float64(len(sorted)-1) * p)
		if i < 0 {
			i = 0
		}
		if i >= len(sorted) {
			i = len(sorted) - 1
		}
		return sorted[i]
	}

	p50 := pick(0.50)
	p95 := pick(0.95)
	p99 := pick(0.99)

	fmt.Printf("%s: P50 = %.3fms, P95 = %.3fms, P99 = %.3fms (samples = %d)\n",
		label,
		float64(p50.Microseconds())/1000,
		float64(p95.Microseconds())/1000,
		float64(p99.Microseconds())/1000,
		len(sorted),
	)
}
