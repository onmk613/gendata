package driver

import (
	"testing"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5"
)

func TestDSNSpecialCharacters(t *testing.T) {
	conf := SqlConfiguration{
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "us er",
		Password: "p@ss:word/&?",
		DBName:   "my db",
	}

	t.Run("mysql", func(t *testing.T) {
		cfg, err := mysqldriver.ParseDSN(conf.getMysqlDsn())
		if err != nil {
			t.Fatalf("parse mysql dsn: %v", err)
		}
		if cfg.User != conf.User || cfg.Passwd != conf.Password || cfg.DBName != conf.DBName {
			t.Fatalf("mysql dsn roundtrip mismatch: %+v", cfg)
		}
	})

	t.Run("postgres", func(t *testing.T) {
		cfg, err := pgx.ParseConfig(conf.getPostgresDsn())
		if err != nil {
			t.Fatalf("parse postgres dsn: %v", err)
		}
		if cfg.User != conf.User || cfg.Password != conf.Password || cfg.Database != conf.DBName {
			t.Fatalf("postgres dsn roundtrip mismatch: %+v", cfg)
		}
	})

	t.Run("clickhouse", func(t *testing.T) {
		opts, err := ch.ParseDSN(conf.getClickhouseDsn())
		if err != nil {
			t.Fatalf("parse clickhouse dsn: %v", err)
		}
		if opts.Auth.Username != conf.User || opts.Auth.Password != conf.Password || opts.Auth.Database != conf.DBName {
			t.Fatalf("clickhouse dsn roundtrip mismatch: %+v", opts.Auth)
		}
	})
}
