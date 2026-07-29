package action

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// Mean Max	Min
var (
	WriteConf          WriteConfiguration
	statsMu            sync.Mutex
	totalTime          time.Duration
	meanThroughputPerS float64
	maxThroughputPerS  float64
	minThroughputPerS  float64
)

func init() {
	minThroughputPerS = math.MaxFloat64
}

// 并发控制
type WriteConfiguration struct {
	Concurrency int
	BatchSize   int
	RepeatCount int
	RowSize     int
}

func logBatchStats(i, j int, insertTime time.Duration, batchSize int) {
	statsMu.Lock()
	defer statsMu.Unlock()

	totalTime += insertTime
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

func printThroughputStats() {
	statsMu.Lock()
	defer statsMu.Unlock()

	if totalTime.Seconds() == 0 {
		return
	}

	meanThroughputPerS = float64(WriteConf.Concurrency) * float64(WriteConf.BatchSize) * float64(WriteConf.RepeatCount) / totalTime.Seconds()

	minVal := minThroughputPerS
	if minVal == math.MaxFloat64 {
		minVal = 0
	}

	fmt.Printf("Total: Time = %.3fs, Mean = %.0f row/s, Max = %.0f row/s, Min = %.0f row/s\n",
		totalTime.Seconds(), meanThroughputPerS, maxThroughputPerS, minVal)
}
