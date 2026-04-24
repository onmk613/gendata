package action

import (
	"fmt"
	"time"
)

// Mean Max	Min
var (
	WriteConf          WriteConfiguration
	totalTime          time.Duration
	meanThroughputPerS float64
	maxThroughputPerS  float64
	minThroughputPerS  float64
)

func init() {
	minThroughputPerS = 10000000
}

// 并发控制
type WriteConfiguration struct {
	Concurrency int
	BatchSize   int
	RepeatCount int
}

func logBatchStats(i, j int, insertTime time.Duration, batchSize int) {
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
	if totalTime.Seconds() == 0 {
		return
	}
	time.Sleep(2 * time.Second)
	meanThroughputPerS = float64(WriteConf.Concurrency) * float64(WriteConf.BatchSize) * float64(WriteConf.RepeatCount) / totalTime.Seconds()
	fmt.Printf("Total: Time = %.3fs, Mean = %.0f row/s, Max = %.0f row/s, Min = %.0f row/s\n",
		totalTime.Seconds(), meanThroughputPerS, maxThroughputPerS, minThroughputPerS)
}
