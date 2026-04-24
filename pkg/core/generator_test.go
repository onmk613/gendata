package core

import (
	"testing"

	mydriver "gendata/pkg/driver"
)

func TestGenerateDefaultTableData(t *testing.T) {
	tests := []struct {
		name  string
		count int
	}{
		{"Single row", 1},
		{"Ten rows", 10},
		{"Hundred rows", 100},
		{"Zero rows", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows := GenerateDefaultTableData(tt.count)

			if len(rows) != tt.count {
				t.Errorf("Expected %d rows, got %d", tt.count, len(rows))
			}

			// 验证每一行都有有效的数据
			for i, row := range rows {
				if row == nil {
					t.Errorf("Row %d is nil", i)
					continue
				}

				// 验证必填字段非空
				if row.UserID == "" {
					t.Errorf("Row %d: UserID is empty", i)
				}
				if row.Name == "" {
					t.Errorf("Row %d: Name is empty", i)
				}
				if row.Account == "" {
					t.Errorf("Row %d: Account is empty", i)
				}

				// 验证年龄范围
				if row.Age < 18 || row.Age > 80 {
					t.Errorf("Row %d: Age out of range: %d", i, row.Age)
				}
			}
		})
	}
}

func TestGenerateDefaultTableDataUniqueness(t *testing.T) {
	rows := GenerateDefaultTableData(100)

	userIDSet := make(map[string]bool)
	accountSet := make(map[string]bool)

	for i, row := range rows {
		// 验证UserID唯一性（大概率）
		if userIDSet[row.UserID] {
			t.Logf("Row %d: Duplicate UserID found (可能出现但不太可能): %s", i, row.UserID)
		}
		userIDSet[row.UserID] = true

		// 验证Account唯一性（大概率）
		if accountSet[row.Account] {
			t.Logf("Row %d: Duplicate Account found: %s", i, row.Account)
		}
		accountSet[row.Account] = true
	}

	// 验证至少有一些唯一的值
	if len(userIDSet) < 90 {
		t.Errorf("Expected most UserIDs to be unique, got %d unique out of 100", len(userIDSet))
	}
}

func TestDefaultTableGeneratorGenerateDefaultTableRow(t *testing.T) {
	gen := newDefaultTableGenerator()

	for i := 0; i < 10; i++ {
		row := gen.generateDefaultTableRow()

		if row == nil {
			t.Error("Generated row is nil")
			continue
		}

		// 验证行的类型
		if _, ok := interface{}(row).(*mydriver.DefaultTableRow); !ok {
			t.Error("Generated row is not of type DefaultTableRow")
		}

		// 验证所有字段都被填充
		if row.UserID == "" {
			t.Error("UserID is empty")
		}
		if row.Name == "" {
			t.Error("Name is empty")
		}
		if row.Phone == "" {
			t.Error("Phone is empty")
		}
		if row.Gender == "" {
			t.Error("Gender is empty")
		}
		if row.Email == "" {
			t.Error("Email is empty")
		}
		if row.Account == "" {
			t.Error("Account is empty")
		}
		if row.Password == "" {
			t.Error("Password is empty")
		}
	}
}

func TestDefaultTableGeneratorDataRandomness(t *testing.T) {
	gen := newDefaultTableGenerator()

	row1 := gen.generateDefaultTableRow()
	row2 := gen.generateDefaultTableRow()

	// UserID应该不同（UUID）
	if row1.UserID == row2.UserID {
		t.Error("UserIDs should be different")
	}

	// 名字可能不同（随机生成）
	if row1.Name == row2.Name {
		t.Logf("Names are the same: %s (可能出现但概率较低)", row1.Name)
	}

	// 年龄可能不同
	if row1.Age == row2.Age {
		t.Logf("Ages are the same: %d (可能出现但概率较低)", row1.Age)
	}
}

func BenchmarkGenerateDefaultTableData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateDefaultTableData(1000)
	}
}

func BenchmarkGenerateDefaultTableRow(b *testing.B) {
	gen := newDefaultTableGenerator()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gen.generateDefaultTableRow()
	}
}
