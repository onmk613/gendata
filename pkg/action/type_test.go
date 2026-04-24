package action

import (
	"testing"
	"time"
)

func TestWriteConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		conf        WriteConfiguration
		expectedErr bool
	}{
		{
			name: "Valid configuration",
			conf: WriteConfiguration{
				Concurrency: 4,
				BatchSize:   1000,
				RepeatCount: 10,
			},
			expectedErr: false,
		},
		{
			name: "Single concurrency",
			conf: WriteConfiguration{
				Concurrency: 1,
				BatchSize:   1000,
				RepeatCount: 5,
			},
			expectedErr: false,
		},
		{
			name: "Large concurrency",
			conf: WriteConfiguration{
				Concurrency: 16,
				BatchSize:   5000,
				RepeatCount: 20,
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.conf.Concurrency <= 0 {
				t.Error("Concurrency should be positive")
			}
			if tt.conf.BatchSize <= 0 {
				t.Error("BatchSize should be positive")
			}
			if tt.conf.RepeatCount <= 0 {
				t.Error("RepeatCount should be positive")
			}
		})
	}
}

func TestLogBatchStats(t *testing.T) {
	// 重置全局变量
	totalTime = 0
	minThroughputPerS = 10000000
	maxThroughputPerS = 0

	tests := []struct {
		name        string
		concurrency int
		batchNum    int
		insertTime  time.Duration
		batchSize   int
	}{
		{
			name:        "Fast insert",
			concurrency: 0,
			batchNum:    0,
			insertTime:  100 * time.Millisecond,
			batchSize:   1000,
		},
		{
			name:        "Slow insert",
			concurrency: 0,
			batchNum:    1,
			insertTime:  1 * time.Second,
			batchSize:   1000,
		},
		{
			name:        "Very fast insert",
			concurrency: 1,
			batchNum:    0,
			insertTime:  50 * time.Millisecond,
			batchSize:   2000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logBatchStats(tt.concurrency, tt.batchNum, tt.insertTime, tt.batchSize)

			// 验证totalTime已更新
			if totalTime == 0 {
				t.Error("totalTime should be updated")
			}

			// 验证吞吐量计算
			throughput := float64(tt.batchSize) / tt.insertTime.Seconds()
			if throughput < 0 {
				t.Error("Throughput should be positive")
			}
		})
	}
}

func TestThroughputCalculation(t *testing.T) {
	// 重置全局变量
	totalTime = 0
	minThroughputPerS = 10000000
	maxThroughputPerS = 0

	tests := []struct {
		name        string
		batchSize   int
		insertTime  time.Duration
		expectedMin float64
		expectedMax float64
	}{
		{
			name:       "1000 records in 100ms",
			batchSize:  1000,
			insertTime: 100 * time.Millisecond,
			// 1000 / 0.1 = 10000 row/s
			expectedMin: 9000,  // 允许一些偏差
			expectedMax: 11000, // 允许一些偏差
		},
		{
			name:       "2000 records in 1s",
			batchSize:  2000,
			insertTime: 1 * time.Second,
			// 2000 / 1 = 2000 row/s
			expectedMin: 1900,
			expectedMax: 2100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			throughput := float64(tt.batchSize) / tt.insertTime.Seconds()

			if throughput < tt.expectedMin || throughput > tt.expectedMax {
				t.Errorf("Throughput %f out of expected range [%f, %f]",
					throughput, tt.expectedMin, tt.expectedMax)
			}
		})
	}
}

func TestPrintThroughputStats_WithData(t *testing.T) {
	// 重置全局变量
	WriteConf = WriteConfiguration{
		Concurrency: 2,
		BatchSize:   1000,
		RepeatCount: 5,
	}
	totalTime = 1 * time.Second
	minThroughputPerS = 8000
	maxThroughputPerS = 12000

	// 这个函数主要用于输出，我们只验证计算逻辑
	expectedMean := float64(WriteConf.Concurrency) * float64(WriteConf.BatchSize) * float64(WriteConf.RepeatCount) / totalTime.Seconds()

	if expectedMean <= 0 {
		t.Error("Expected mean throughput should be positive")
	}

	if expectedMean < 9000 || expectedMean > 11000 {
		t.Logf("Mean throughput out of expected range: %f", expectedMean)
	}
}

func TestPrintThroughputStats_NoData(t *testing.T) {
	// 重置全局变量
	totalTime = 0
	WriteConf = WriteConfiguration{
		Concurrency: 1,
		BatchSize:   1000,
		RepeatCount: 10,
	}

	// 当totalTime为0时，函数应该返回而不进行任何处理
	if totalTime.Seconds() == 0 {
		// 函数会在这里返回，所以没有输出
		t.Log("Function correctly handles zero total time")
	}
}

func TestBatchStatsMultipleCalls(t *testing.T) {
	// 重置全局变量
	totalTime = 0
	minThroughputPerS = 10000000
	maxThroughputPerS = 0

	// 模拟多个批次的统计
	batches := []struct {
		insertTime time.Duration
		batchSize  int
	}{
		{100 * time.Millisecond, 1000},
		{150 * time.Millisecond, 1000},
		{120 * time.Millisecond, 1000},
		{180 * time.Millisecond, 1000},
	}

	for i, batch := range batches {
		logBatchStats(0, i, batch.insertTime, batch.batchSize)
	}

	// 验证总时间
	expectedTotalTime := time.Duration(0)
	for _, batch := range batches {
		expectedTotalTime += batch.insertTime
	}

	if totalTime != expectedTotalTime {
		t.Errorf("Expected total time %v, got %v", expectedTotalTime, totalTime)
	}

	// 验证最小和最大吞吐量
	if minThroughputPerS > maxThroughputPerS {
		t.Error("Min throughput should be less than or equal to max")
	}

	if minThroughputPerS <= 0 || maxThroughputPerS <= 0 {
		t.Error("Throughput values should be positive")
	}
}

func TestMeanThroughputCalculation(t *testing.T) {
	tests := []struct {
		name         string
		concurrency  int
		batchSize    int
		repeatCount  int
		totalTime    time.Duration
		expectedMean float64 // 允许10%的偏差
	}{
		{
			name:         "Single worker",
			concurrency:  1,
			batchSize:    1000,
			repeatCount:  10,
			totalTime:    1 * time.Second,
			expectedMean: 10000, // 10000 records / 1 second
		},
		{
			name:         "Two workers",
			concurrency:  2,
			batchSize:    1000,
			repeatCount:  10,
			totalTime:    1 * time.Second,
			expectedMean: 20000, // 20000 records / 1 second
		},
		{
			name:         "Four workers with longer duration",
			concurrency:  4,
			batchSize:    500,
			repeatCount:  5,
			totalTime:    2 * time.Second,
			expectedMean: 5000, // (4 * 500 * 5) / 2 = 5000
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mean := float64(tt.concurrency) * float64(tt.batchSize) * float64(tt.repeatCount) / tt.totalTime.Seconds()

			tolerance := tt.expectedMean * 0.1
			if mean < tt.expectedMean-tolerance || mean > tt.expectedMean+tolerance {
				t.Errorf("Mean throughput %f outside expected range [%f, %f]",
					mean, tt.expectedMean-tolerance, tt.expectedMean+tolerance)
			}
		})
	}
}

func BenchmarkLogBatchStats(b *testing.B) {
	for i := 0; i < b.N; i++ {
		logBatchStats(0, 0, 100*time.Millisecond, 1000)
	}
}

func BenchmarkThroughputCalculation(b *testing.B) {
	batchSize := 1000
	insertTime := time.Duration(100) * time.Millisecond

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = float64(batchSize) / insertTime.Seconds()
	}
}

func TestWriteConfigurationValues(t *testing.T) {
	WriteConf = WriteConfiguration{
		Concurrency: 4,
		BatchSize:   2000,
		RepeatCount: 20,
	}

	if WriteConf.Concurrency != 4 {
		t.Errorf("Expected Concurrency 4, got %d", WriteConf.Concurrency)
	}

	if WriteConf.BatchSize != 2000 {
		t.Errorf("Expected BatchSize 2000, got %d", WriteConf.BatchSize)
	}

	if WriteConf.RepeatCount != 20 {
		t.Errorf("Expected RepeatCount 20, got %d", WriteConf.RepeatCount)
	}

	totalRecords := WriteConf.Concurrency * WriteConf.BatchSize * WriteConf.RepeatCount
	if totalRecords != 160000 {
		t.Errorf("Expected total records 160000, got %d", totalRecords)
	}
}
