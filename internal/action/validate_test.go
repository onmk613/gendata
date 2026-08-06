package action

import (
	"testing"
	"time"
)

func withWriteConf(mode string, concurrency, readConcurrency, batchsize int, duration time.Duration) func() {
	old := WriteConf
	WriteConf = WriteConfiguration{
		Mode:            mode,
		Concurrency:     concurrency,
		ReadConcurrency: readConcurrency,
		BatchSize:       batchsize,
		Duration:        duration,
	}
	return func() { WriteConf = old }
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		mode    string
		conf    WriteConfiguration
		wantErr bool
	}{
		{name: "valid write", mode: ModeWrite, conf: WriteConfiguration{Mode: ModeWrite, Concurrency: 2, BatchSize: 1000}, wantErr: false},
		{name: "unknown mode", mode: "wite", conf: WriteConfiguration{Mode: "wite", Concurrency: 1, BatchSize: 1000}, wantErr: true},
		{name: "zero batchsize", mode: ModeWrite, conf: WriteConfiguration{Mode: ModeWrite, Concurrency: 1, BatchSize: 0}, wantErr: true},
		{name: "zero concurrency", mode: ModeWrite, conf: WriteConfiguration{Mode: ModeWrite, Concurrency: 0, BatchSize: 1000}, wantErr: true},
		{name: "read without duration", mode: ModeRead, conf: WriteConfiguration{Mode: ModeRead, Concurrency: 1, ReadConcurrency: 4, BatchSize: 1000}, wantErr: true},
		{name: "valid read", mode: ModeRead, conf: WriteConfiguration{Mode: ModeRead, Concurrency: 1, ReadConcurrency: 4, BatchSize: 1000, Duration: time.Minute}, wantErr: false},
		{name: "mixed without duration", mode: ModeMixed, conf: WriteConfiguration{Mode: ModeMixed, Concurrency: 2, ReadConcurrency: 4, BatchSize: 1000}, wantErr: true},
		{name: "mixed zero readers", mode: ModeMixed, conf: WriteConfiguration{Mode: ModeMixed, Concurrency: 2, ReadConcurrency: 0, BatchSize: 1000, Duration: time.Minute}, wantErr: true},
		{name: "valid mixed", mode: ModeMixed, conf: WriteConfiguration{Mode: ModeMixed, Concurrency: 2, ReadConcurrency: 4, BatchSize: 1000, Duration: time.Minute}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := withWriteConf(tt.mode, tt.conf.Concurrency, tt.conf.ReadConcurrency, tt.conf.BatchSize, tt.conf.Duration)
			defer restore()

			err := validateConfig()
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
