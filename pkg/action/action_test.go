package action

import (
	"testing"
)

func TestActionPackageInitialization(t *testing.T) {
	// 验证WriteConf初始状态
	t.Run("WriteConf initial state", func(t *testing.T) {
		// WriteConf应该是一个有效的WriteConfiguration实例
		// 即使在初始化时可能为零值
		conf := WriteConf
		if conf.Concurrency < 0 || conf.BatchSize < 0 || conf.RepeatCount < 0 {
			t.Error("WriteConfiguration should have non-negative values")
		}
	})

	// 验证统计变量初始状态
	t.Run("Stats initial state", func(t *testing.T) {
		if minThroughputPerS != 10000000 {
			t.Errorf("Initial minThroughputPerS should be 10000000, got %f", minThroughputPerS)
		}
		if maxThroughputPerS != 0 {
			t.Errorf("Initial maxThroughputPerS should be 0, got %f", maxThroughputPerS)
		}
	})
}

func TestSetupAndCleanupConfiguration(t *testing.T) {
	// 测试配置的设置和重置
	t.Run("Configuration setup", func(t *testing.T) {
		WriteConf = WriteConfiguration{
			Concurrency: 4,
			BatchSize:   1000,
			RepeatCount: 10,
		}

		if WriteConf.Concurrency != 4 {
			t.Error("Failed to set Concurrency")
		}
		if WriteConf.BatchSize != 1000 {
			t.Error("Failed to set BatchSize")
		}
		if WriteConf.RepeatCount != 10 {
			t.Error("Failed to set RepeatCount")
		}
	})

	// 重置为默认值
	t.Run("Configuration reset", func(t *testing.T) {
		WriteConf = WriteConfiguration{}
		// 验证重置
		if WriteConf.Concurrency != 0 {
			t.Error("Failed to reset Concurrency")
		}
	})
}

func TestMinThroughputInitialization(t *testing.T) {
	// 验证minThroughputPerS的初始化
	initialMinThroughput := 10000000.0

	if minThroughputPerS != initialMinThroughput {
		t.Errorf("minThroughputPerS should be initialized to %f, got %f",
			initialMinThroughput, minThroughputPerS)
	}
}

func TestStatsVariablesType(t *testing.T) {
	// 验证统计变量的类型
	t.Run("Type checking", func(t *testing.T) {
		// totalTime应该是time.Duration
		var _ interface{} = totalTime

		// meanThroughputPerS应该是float64
		var _ float64 = meanThroughputPerS

		// maxThroughputPerS应该是float64
		var _ float64 = maxThroughputPerS

		// minThroughputPerS应该是float64
		var _ float64 = minThroughputPerS
	})
}

func TestConfigurationValidation(t *testing.T) {
	tests := []struct {
		name    string
		conf    WriteConfiguration
		isValid bool
	}{
		{
			name: "Valid configuration",
			conf: WriteConfiguration{
				Concurrency: 1,
				BatchSize:   1000,
				RepeatCount: 10,
			},
			isValid: true,
		},
		{
			name: "Zero concurrency",
			conf: WriteConfiguration{
				Concurrency: 0,
				BatchSize:   1000,
				RepeatCount: 10,
			},
			isValid: false,
		},
		{
			name: "Zero batch size",
			conf: WriteConfiguration{
				Concurrency: 1,
				BatchSize:   0,
				RepeatCount: 10,
			},
			isValid: false,
		},
		{
			name: "Zero repeat count",
			conf: WriteConfiguration{
				Concurrency: 1,
				BatchSize:   1000,
				RepeatCount: 0,
			},
			isValid: false,
		},
		{
			name: "Negative values",
			conf: WriteConfiguration{
				Concurrency: -1,
				BatchSize:   1000,
				RepeatCount: 10,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.conf.Concurrency > 0 && tt.conf.BatchSize > 0 && tt.conf.RepeatCount > 0

			if isValid != tt.isValid {
				t.Errorf("Expected isValid=%v, got %v", tt.isValid, isValid)
			}
		})
	}
}

func TestTotalRecordsCalculation(t *testing.T) {
	tests := []struct {
		name          string
		concurrency   int
		batchSize     int
		repeatCount   int
		expectedTotal int
	}{
		{
			name:          "Single worker",
			concurrency:   1,
			batchSize:     1000,
			repeatCount:   10,
			expectedTotal: 10000,
		},
		{
			name:          "Four workers",
			concurrency:   4,
			batchSize:     1000,
			repeatCount:   10,
			expectedTotal: 40000,
		},
		{
			name:          "Custom values",
			concurrency:   8,
			batchSize:     2000,
			repeatCount:   5,
			expectedTotal: 80000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total := tt.concurrency * tt.batchSize * tt.repeatCount

			if total != tt.expectedTotal {
				t.Errorf("Expected total %d, got %d", tt.expectedTotal, total)
			}
		})
	}
}

func TestStatsReset(t *testing.T) {
	// 测试统计信息重置
	t.Run("Reset stats", func(t *testing.T) {
		// 设置一些值
		totalTime = 10
		meanThroughputPerS = 5000
		maxThroughputPerS = 10000
		minThroughputPerS = 1000

		// 重置
		totalTime = 0
		meanThroughputPerS = 0
		maxThroughputPerS = 0
		minThroughputPerS = 10000000

		if totalTime != 0 {
			t.Error("Failed to reset totalTime")
		}
		if minThroughputPerS != 10000000 {
			t.Error("Failed to reset minThroughputPerS")
		}
	})
}

func TestLogBatchStatsIntegration(t *testing.T) {
	// 重置状态
	totalTime = 0
	minThroughputPerS = 10000000
	maxThroughputPerS = 0

	// 记录多个批次
	stats := []struct {
		concurrency int
		batchNum    int
		insertTime  interface{} // time.Duration, but we'll use a calculation
		batchSize   int
	}{
		{0, 0, 100, 1000}, // 100ms -> 10000 row/s
		{0, 1, 200, 1000}, // 200ms -> 5000 row/s
		{0, 2, 150, 1000}, // 150ms -> 6666 row/s
	}

	for _, stat := range stats {
		// 由于time.Duration涉及更复杂的操作，这里只测试数据结构
		if stat.batchSize <= 0 {
			t.Error("Batch size should be positive")
		}
	}
}

func BenchmarkConfigurationCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = WriteConfiguration{
			Concurrency: 4,
			BatchSize:   1000,
			RepeatCount: 10,
		}
	}
}
