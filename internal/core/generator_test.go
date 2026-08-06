package core

import (
	"testing"

	"gendata/internal/generator"
)

func TestGenerateDefaultTableData(t *testing.T) {
	rows := GenerateDefaultTableData(10, nil)
	if len(rows) != 10 {
		t.Fatalf("got %d rows, want 10", len(rows))
	}
	for i, row := range rows {
		if row.UserID == "" {
			t.Fatalf("row %d has empty UserID", i)
		}
		if row.Name == "" || row.Email == "" || row.Address == "" {
			t.Fatalf("row %d has empty generated fields: %+v", i, row)
		}
	}
}

func TestGenerateDefaultTableDataDeterministic(t *testing.T) {
	rng := generator.NewRandom(42)
	a := GenerateDefaultTableData(3, rng)
	rng = generator.NewRandom(42)
	b := GenerateDefaultTableData(3, rng)
	for i := range a {
		if a[i].Name != b[i].Name || a[i].Email != b[i].Email || a[i].JsonData != b[i].JsonData {
			t.Fatalf("row %d random content mismatch: %+v vs %+v", i, a[i], b[i])
		}
	}
}
