package driver

import (
	"strings"
	"testing"
)

func TestGetMysqlDsn(t *testing.T) {
	tests := []struct {
		name     string
		conf     SqlConfiguration
		expected []string
		notFound []string
	}{
		{
			name: "With all parameters",
			conf: SqlConfiguration{
				Host:           "localhost",
				Port:           3306,
				User:           "root",
				Password:       "password123",
				DBName:         "testdb",
				AdditionalArgs: nil,
			},
			expected: []string{
				"root:password123",
				"localhost",
				"3306",
				"testdb",
				"charset=utf8",
			},
		},
		{
			name: "With default values",
			conf: SqlConfiguration{
				Host:           "",
				Port:           0,
				User:           "",
				Password:       "",
				DBName:         "",
				AdditionalArgs: nil,
			},
			expected: []string{
				"root",
				"127.0.0.1",
				"3306",
				"mysql",
				"charset=utf8",
			},
		},
		{
			name: "Without password",
			conf: SqlConfiguration{
				Host:           "db.example.com",
				Port:           3306,
				User:           "user",
				Password:       "",
				DBName:         "mydb",
				AdditionalArgs: nil,
			},
			expected: []string{
				"user@tcp",
				"db.example.com",
				"3306",
				"mydb",
			},
			notFound: []string{":password"},
		},
		{
			name: "With additional args",
			conf: SqlConfiguration{
				Host:     "localhost",
				Port:     3306,
				User:     "root",
				Password: "pass",
				DBName:   "db",
				AdditionalArgs: map[string]string{
					"charset": "utf8mb4",
				},
			},
			expected: []string{
				"utf8mb4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := tt.conf.getMysqlDsn()

			// 检查应该包含的子串
			for _, substr := range tt.expected {
				if !strings.Contains(dsn, substr) {
					t.Errorf("DSN should contain '%s'\nGot: %s", substr, dsn)
				}
			}

			// 检查不应该包含的子串
			for _, substr := range tt.notFound {
				if strings.Contains(dsn, substr) {
					t.Errorf("DSN should not contain '%s'\nGot: %s", substr, dsn)
				}
			}
		})
	}
}

func TestGetPostgresDsn(t *testing.T) {
	tests := []struct {
		name     string
		conf     SqlConfiguration
		expected []string
		notFound []string
	}{
		{
			name: "With all parameters",
			conf: SqlConfiguration{
				Host:           "localhost",
				Port:           5432,
				User:           "postgres",
				Password:       "pgpass",
				DBName:         "pgdb",
				AdditionalArgs: nil,
			},
			expected: []string{
				"user=postgres",
				"password=pgpass",
				"host=localhost",
				"port=5432",
				"dbname=pgdb",
				"sslmode=disable",
			},
		},
		{
			name: "With default values",
			conf: SqlConfiguration{
				Host:           "",
				Port:           0,
				User:           "",
				Password:       "",
				DBName:         "",
				AdditionalArgs: nil,
			},
			expected: []string{
				"user=postgres",
				"host=127.0.0.1",
				"port=5432",
				"dbname=postgres",
				"sslmode=disable",
			},
		},
		{
			name: "Without password",
			conf: SqlConfiguration{
				Host:           "pghost.com",
				Port:           5432,
				User:           "pguser",
				Password:       "",
				DBName:         "mydb",
				AdditionalArgs: nil,
			},
			expected: []string{
				"user=pguser",
				"host=pghost.com",
				"dbname=mydb",
			},
			notFound: []string{"password="},
		},
		{
			name: "With additional args",
			conf: SqlConfiguration{
				Host:     "localhost",
				Port:     5432,
				User:     "user",
				Password: "pass",
				DBName:   "db",
				AdditionalArgs: map[string]string{
					"sslmode": "require",
				},
			},
			expected: []string{
				"sslmode=require",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := tt.conf.getPostgresDsn()

			// 检查应该包含的子串
			for _, substr := range tt.expected {
				if !strings.Contains(dsn, substr) {
					t.Errorf("DSN should contain '%s'\nGot: %s", substr, dsn)
				}
			}

			// 检查不应该包含的子串
			for _, substr := range tt.notFound {
				if strings.Contains(dsn, substr) {
					t.Errorf("DSN should not contain '%s'\nGot: %s", substr, dsn)
				}
			}
		})
	}
}

func TestGetClickhouseDsn(t *testing.T) {
	tests := []struct {
		name     string
		conf     SqlConfiguration
		expected []string
		notFound []string
	}{
		{
			name: "With all parameters",
			conf: SqlConfiguration{
				Host:           "localhost",
				Port:           9000,
				User:           "default",
				Password:       "ckpass",
				DBName:         "ckdb",
				AdditionalArgs: nil,
			},
			expected: []string{
				"tcp://",
				"localhost:9000",
				"username=default",
				"password=ckpass",
				"database=ckdb",
				"write_timeout=20",
			},
		},
		{
			name: "With default values",
			conf: SqlConfiguration{
				Host:           "",
				Port:           0,
				User:           "",
				Password:       "",
				DBName:         "",
				AdditionalArgs: nil,
			},
			expected: []string{
				"tcp://127.0.0.1:9000",
				"username=default",
				"database=default",
				"write_timeout=20",
			},
		},
		{
			name: "Without password",
			conf: SqlConfiguration{
				Host:           "ckhost.com",
				Port:           9000,
				User:           "ckuser",
				Password:       "",
				DBName:         "mydb",
				AdditionalArgs: nil,
			},
			expected: []string{
				"tcp://ckhost.com:9000",
				"username=ckuser",
				"database=mydb",
			},
			notFound: []string{"password="},
		},
		{
			name: "With additional args",
			conf: SqlConfiguration{
				Host:     "localhost",
				Port:     9000,
				User:     "user",
				Password: "pass",
				DBName:   "db",
				AdditionalArgs: map[string]string{
					"read_timeout": "30",
				},
			},
			expected: []string{
				"read_timeout=30",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := tt.conf.getClickhouseDsn()

			// 检查应该包含的子串
			for _, substr := range tt.expected {
				if !strings.Contains(dsn, substr) {
					t.Errorf("DSN should contain '%s'\nGot: %s", substr, dsn)
				}
			}

			// 检查不应该包含的子串
			for _, substr := range tt.notFound {
				if strings.Contains(dsn, substr) {
					t.Errorf("DSN should not contain '%s'\nGot: %s", substr, dsn)
				}
			}
		})
	}
}

func TestDefaultTableRow_TableName(t *testing.T) {
	// 保存原始表名
	originalTable := SqlConf.Table

	tests := []struct {
		name      string
		tableName string
	}{
		{"Default table name", "gendata_table"},
		{"Custom table name", "users"},
		{"Another table name", "test_data"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SqlConf.Table = tt.tableName
			row := &DefaultTableRow{}

			if row.TableName() != tt.tableName {
				t.Errorf("Expected table name '%s', got '%s'", tt.tableName, row.TableName())
			}
		})
	}

	// 恢复原始表名
	SqlConf.Table = originalTable
}

func BenchmarkGetMysqlDsn(b *testing.B) {
	conf := SqlConfiguration{
		Host:           "localhost",
		Port:           3306,
		User:           "root",
		Password:       "password",
		DBName:         "testdb",
		AdditionalArgs: nil,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conf.getMysqlDsn()
	}
}

func BenchmarkGetPostgresDsn(b *testing.B) {
	conf := SqlConfiguration{
		Host:           "localhost",
		Port:           5432,
		User:           "postgres",
		Password:       "password",
		DBName:         "testdb",
		AdditionalArgs: nil,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conf.getPostgresDsn()
	}
}

func BenchmarkGetClickhouseDsn(b *testing.B) {
	conf := SqlConfiguration{
		Host:           "localhost",
		Port:           9000,
		User:           "default",
		Password:       "password",
		DBName:         "testdb",
		AdditionalArgs: nil,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conf.getClickhouseDsn()
	}
}
