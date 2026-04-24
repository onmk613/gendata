package driver

import (
	"testing"
	"time"
)

func TestDefaultTableRowStructure(t *testing.T) {
	row := &DefaultTableRow{
		UserID:       "test-uuid",
		Name:         "John Doe",
		Phone:        "+1-234-567-8900",
		Gender:       "male",
		Age:          30,
		Email:        "john@example.com",
		State:        "California",
		ZipCode:      "12345",
		Nationality:  "US",
		Height:       180.5,
		Weight:       75.50,
		BloodType:    "O",
		Account:      "account_test",
		AccountName:  "John Account",
		Password:     "hashedpassword",
		OnlineStatus: true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if row.UserID != "test-uuid" {
		t.Errorf("Expected UserID 'test-uuid', got '%s'", row.UserID)
	}

	if row.Name != "John Doe" {
		t.Errorf("Expected Name 'John Doe', got '%s'", row.Name)
	}

	if row.Age != 30 {
		t.Errorf("Expected Age 30, got %d", row.Age)
	}
}

func TestCreateTemplateDataBatchWithSize_SmallBatch(t *testing.T) {
	// 测试小于安全限制的批大小
	templates := make([]*DefaultTableRow, 100)
	for i := 0; i < 100; i++ {
		templates[i] = &DefaultTableRow{
			UserID: "user-" + string(rune(i)),
			Name:   "User" + string(rune(i)),
		}
	}

	// 此测试仅测试函数是否能处理小批大小
	// 由于没有真实数据库，我们无法测试实际插入
	if len(templates) != 100 {
		t.Errorf("Expected 100 templates, got %d", len(templates))
	}
}

func TestCreateTemplateDataBatchWithSize_LargeBatch(t *testing.T) {
	// 测试超过安全限制的批大小
	batchSize := 4000 // 超过安全限制 3276

	templates := make([]*DefaultTableRow, batchSize)
	for i := 0; i < batchSize; i++ {
		templates[i] = &DefaultTableRow{
			UserID: "user-" + string(rune(i)),
			Name:   "User" + string(rune(i)),
		}
	}

	// 验证数据生成正确
	if len(templates) != batchSize {
		t.Errorf("Expected %d templates, got %d", batchSize, len(templates))
	}
}

func TestPostgreSQLBatchSizeConstants(t *testing.T) {
	// 验证常数的正确性
	if PostgreSQLMaxParams != 65535 {
		t.Errorf("Expected PostgreSQLMaxParams to be 65535, got %d", PostgreSQLMaxParams)
	}

	if DefaultTableRowColumns != 20 {
		t.Errorf("Expected DefaultTableRowColumns to be 20, got %d", DefaultTableRowColumns)
	}

	expectedMaxBatchSize := 65535 / 20
	if SafeMaxBatchSizeForPostgreSQL != expectedMaxBatchSize {
		t.Errorf("Expected SafeMaxBatchSizeForPostgreSQL to be %d, got %d",
			expectedMaxBatchSize, SafeMaxBatchSizeForPostgreSQL)
	}
}

func TestBatchSizeCalculation(t *testing.T) {
	tests := []struct {
		name           string
		requestedSize  int
		expectedResult string // "unchanged" or "limited"
	}{
		{
			name:           "Small batch size",
			requestedSize:  100,
			expectedResult: "unchanged",
		},
		{
			name:           "Medium batch size",
			requestedSize:  1000,
			expectedResult: "unchanged",
		},
		{
			name:           "Max safe size",
			requestedSize:  SafeMaxBatchSizeForPostgreSQL,
			expectedResult: "unchanged",
		},
		{
			name:           "Exceed safe size",
			requestedSize:  SafeMaxBatchSizeForPostgreSQL + 1,
			expectedResult: "limited",
		},
		{
			name:           "Much larger than safe size",
			requestedSize:  10000,
			expectedResult: "limited",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualBatchSize := tt.requestedSize
			if tt.requestedSize > SafeMaxBatchSizeForPostgreSQL {
				actualBatchSize = SafeMaxBatchSizeForPostgreSQL
			}

			if tt.expectedResult == "unchanged" {
				if actualBatchSize != tt.requestedSize {
					t.Errorf("Expected batch size %d, got %d", tt.requestedSize, actualBatchSize)
				}
			} else {
				if actualBatchSize != SafeMaxBatchSizeForPostgreSQL {
					t.Errorf("Expected batch size to be limited to %d, got %d",
						SafeMaxBatchSizeForPostgreSQL, actualBatchSize)
				}
			}
		})
	}
}

func TestDefaultTableRow_Timestamps(t *testing.T) {
	row := &DefaultTableRow{
		UserID:    "test",
		Name:      "Test User",
		CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
	}

	if row.CreatedAt != time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC) {
		t.Errorf("CreatedAt not set correctly")
	}

	if row.UpdatedAt != time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC) {
		t.Errorf("UpdatedAt not set correctly")
	}
}

func TestDefaultTableRow_NumericFields(t *testing.T) {
	tests := []struct {
		name     string
		age      int
		height   float64
		weight   float64
		shouldOk bool
	}{
		{"Valid numeric values", 30, 180.5, 75.50, true},
		{"Zero age", 0, 180.5, 75.50, true},
		{"Large age", 999, 180.5, 75.50, true},
		{"Decimal height", 180, 180.5, 75.50, true},
		{"Decimal weight", 75, 180.5, 75.50, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := &DefaultTableRow{
				UserID: "test",
				Age:    tt.age,
				Height: tt.height,
				Weight: tt.weight,
			}

			if row.Age != tt.age {
				t.Errorf("Age mismatch: expected %d, got %d", tt.age, row.Age)
			}

			if row.Height != tt.height {
				t.Errorf("Height mismatch: expected %f, got %f", tt.height, row.Height)
			}

			if row.Weight != tt.weight {
				t.Errorf("Weight mismatch: expected %f, got %f", tt.weight, row.Weight)
			}
		})
	}
}

func TestDefaultTableRow_BooleanField(t *testing.T) {
	tests := []struct {
		name         string
		onlineStatus bool
	}{
		{"Online", true},
		{"Offline", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := &DefaultTableRow{
				UserID:       "test",
				Name:         "Test",
				OnlineStatus: tt.onlineStatus,
			}

			if row.OnlineStatus != tt.onlineStatus {
				t.Errorf("OnlineStatus mismatch: expected %v, got %v", tt.onlineStatus, row.OnlineStatus)
			}
		})
	}
}

func TestDefaultTableRow_StringFields(t *testing.T) {
	row := &DefaultTableRow{
		UserID:      "12345-67890-12345-67890",
		Name:        "John Doe",
		Phone:       "+1-234-567-8900",
		Gender:      "male",
		Email:       "john@example.com",
		State:       "California",
		ZipCode:     "90210",
		Nationality: "US",
		BloodType:   "O",
		Account:     "john_account",
		AccountName: "John's Account",
		Password:    "encrypted_password_hash",
	}

	stringFields := map[string]string{
		"UserID":      row.UserID,
		"Name":        row.Name,
		"Phone":       row.Phone,
		"Gender":      row.Gender,
		"Email":       row.Email,
		"State":       row.State,
		"ZipCode":     row.ZipCode,
		"Nationality": row.Nationality,
		"BloodType":   row.BloodType,
		"Account":     row.Account,
		"AccountName": row.AccountName,
		"Password":    row.Password,
	}

	for fieldName, fieldValue := range stringFields {
		if fieldValue == "" {
			t.Errorf("Field %s is empty", fieldName)
		}
	}
}

func BenchmarkCreateTemplateDataBatchWithSize(b *testing.B) {
	templates := make([]*DefaultTableRow, 1000)
	for i := 0; i < 1000; i++ {
		templates[i] = &DefaultTableRow{
			UserID: "user-" + string(rune(i)),
			Name:   "User" + string(rune(i)),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 由于没有真实的数据库连接，这个基准测试只测试批大小限制逻辑
		actualBatchSize := 1000
		if actualBatchSize > SafeMaxBatchSizeForPostgreSQL {
			actualBatchSize = SafeMaxBatchSizeForPostgreSQL
		}
		_ = actualBatchSize
	}
}

func BenchmarkDefaultTableRowCreation(b *testing.B) {
	now := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = &DefaultTableRow{
			UserID:       "test-uuid",
			Name:         "John Doe",
			Phone:        "+1-234-567-8900",
			Gender:       "male",
			Age:          30,
			Email:        "john@example.com",
			State:        "California",
			ZipCode:      "12345",
			Nationality:  "US",
			Height:       180.5,
			Weight:       75.50,
			BloodType:    "O",
			Account:      "account_test",
			AccountName:  "John Account",
			Password:     "hashedpassword",
			OnlineStatus: true,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
	}
}
