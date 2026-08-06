package action

import (
	"math"
	"testing"

	"gendata/internal/core"
)

func TestCalcBatchCountStaysWithinParamLimit(t *testing.T) {
	rows := core.GenerateDefaultTableData(1, nil)

	for _, dbType := range []string{"mysql", "postgres"} {
		t.Run(dbType, func(t *testing.T) {
			got, err := calcBatchCount(rows, dbType)
			if err != nil {
				t.Fatalf("calcBatchCount: %v", err)
			}
			if got <= 0 {
				t.Fatalf("batch count should be > 0, got %d", got)
			}
			if got*24 > 65535 {
				t.Fatalf("batch count %d exceeds 65535 parameter limit (24 cols)", got)
			}
		})
	}
}

func TestCalcBatchCountClickHouse(t *testing.T) {
	rows := core.GenerateDefaultTableData(1, nil)
	got, err := calcBatchCount(rows, "clickhouse")
	if err != nil {
		t.Fatalf("calcBatchCount: %v", err)
	}
	if got != math.MaxInt32 {
		t.Fatalf("clickhouse batch count = %d, want %d", got, math.MaxInt32)
	}
}

func TestCalcBatchCountUnsupportedDB(t *testing.T) {
	rows := core.GenerateDefaultTableData(1, nil)
	if _, err := calcBatchCount(rows, "oracle"); err == nil {
		t.Fatal("expected error for unsupported database")
	}
}
