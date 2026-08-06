package driver

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	"gorm.io/driver/clickhouse"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"gendata/internal/core"
)

var (
	DB         *gorm.DB
	driverName string
	chConn     ch.Conn

	// ClickHouse 原生批量写入的列元信息缓存
	chInsertColumns []string
	chInsertFields  []*schema.Field
)

func GetDriver(name string, size, concurrency int) error {
	driverName = name
	var open gorm.Dialector
	dsn := SqlConf

	// 配置dsn
	switch name {
	case "mysql":
		open = gormmysql.Open(dsn.getMysqlDsn())
	case "postgres":
		open = postgres.Open(dsn.getPostgresDsn())
	case "clickhouse":
		open = clickhouse.Open(dsn.getClickhouseDsn())
	default:
		return fmt.Errorf("Unsupported Database")
	}

	// 连接数据库
	db, err := gorm.Open(open, &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Silent,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
		// 让 GORM 在 Create 时真正按安全大小分批；
		// 压测场景下不需要把整批包在事务里。
		CreateBatchSize:        size,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return err
	}

	if driverName == "clickhouse" {
		// ClickHouse 没有唯一约束，用 user_id 作为排序键才能支撑点查压测。
		db = db.Set("gorm:table_options", "ENGINE=MergeTree() ORDER BY user_id")
		if err := openClickHouseNative(concurrency); err != nil {
			return err
		}
	} else {
		// 配置连接池
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		maxOpen := concurrency
		if maxOpen < 1 {
			maxOpen = 1
		}
		sqlDB.SetMaxOpenConns(maxOpen)
		sqlDB.SetMaxIdleConns(maxOpen)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	DB = db
	return nil
}

// InsertRows 按数据库类型写入一批数据：
// MySQL/PostgreSQL 走 GORM 的 CreateInBatches；
// ClickHouse 走原生 Batch API，避免 GORM 驱动逐行 Insert。
func InsertRows(ctx context.Context, rows []*core.DefaultTableRow) error {
	if driverName == "clickhouse" {
		return clickHouseInsertRows(ctx, rows)
	}
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}
	return DB.WithContext(ctx).Create(rows).Error
}

func openClickHouseNative(maxOpenConns int) error {
	opts, err := ch.ParseDSN(SqlConf.getClickhouseDsn())
	if err != nil {
		return fmt.Errorf("parse clickhouse dsn: %w", err)
	}
	if maxOpenConns < 1 {
		maxOpenConns = 1
	}
	opts.MaxOpenConns = maxOpenConns
	opts.MaxIdleConns = maxOpenConns

	chConn, err = ch.Open(opts)
	if err != nil {
		return fmt.Errorf("open clickhouse native connection: %w", err)
	}
	return nil
}

func clickHouseInsertRows(ctx context.Context, rows []*core.DefaultTableRow) error {
	if len(rows) == 0 {
		return nil
	}
	if err := ensureClickHouseInsertMeta(); err != nil {
		return err
	}

	query := "INSERT INTO " + quoteIdent(core.TableName) +
		" (" + strings.Join(quoteIdents(chInsertColumns), ", ") + ")"
	batch, err := chConn.PrepareBatch(ctx, query)
	if err != nil {
		return fmt.Errorf("clickhouse prepare batch: %w", err)
	}

	for _, row := range rows {
		values := make([]any, 0, len(chInsertFields))
		rv := reflect.ValueOf(row).Elem()
		for _, field := range chInsertFields {
			v, _ := field.ValueOf(ctx, rv)
			values = append(values, v)
		}
		if err := batch.Append(values...); err != nil {
			_ = batch.Abort()
			return fmt.Errorf("clickhouse append row: %w", err)
		}
	}
	if err := batch.Send(); err != nil {
		return fmt.Errorf("clickhouse send batch: %w", err)
	}
	return nil
}

func ensureClickHouseInsertMeta() error {
	if len(chInsertColumns) > 0 {
		return nil
	}
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	tx := DB.Session(&gorm.Session{})
	if err := tx.Statement.Parse(&core.DefaultTableRow{}); err != nil {
		return fmt.Errorf("parse table schema: %w", err)
	}
	stmt := tx.Statement
	if stmt.Schema == nil {
		return fmt.Errorf("table schema is nil")
	}
	for _, dbName := range stmt.Schema.DBNames {
		field := stmt.Schema.FieldsByDBName[dbName]
		if field == nil {
			continue
		}
		chInsertColumns = append(chInsertColumns, dbName)
		chInsertFields = append(chInsertFields, field)
	}
	return nil
}

func quoteIdent(s string) string {
	parts := strings.Split(s, ".")
	for i := range parts {
		parts[i] = "`" + strings.ReplaceAll(parts[i], "`", "``") + "`"
	}
	return strings.Join(parts, ".")
}

func quoteIdents(names []string) []string {
	out := make([]string, len(names))
	for i, name := range names {
		out[i] = quoteIdent(name)
	}
	return out
}

// 结束数据库
func CloseDB() {
	if chConn != nil {
		_ = chConn.Close()
		chConn = nil
	}
	if DB != nil {
		if sqlDB, err := DB.DB(); err == nil && sqlDB != nil {
			_ = sqlDB.Close()
		}
	}
}
