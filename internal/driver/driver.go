package driver

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func GetDriver(driverName string, size, concurrency int) error {
	var open gorm.Dialector
	dsn := SqlConf

	// 配置dsn
	switch driverName {
	case "mysql":
		open = mysql.Open(dsn.getMysqlDsn())
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
	})
	if err != nil {
		return err
	}

	// 注册插入控制插件
	db.Use(NewBatchPlugin(size))

	// 配置连接池（ClickHouse 使用原生客户端，跳过标准连接池配置）
	if driverName != "clickhouse" {
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

// 结束数据库
func CloseDB() {
	if DB != nil {
		sqlDB, _ := DB.DB()
		sqlDB.Close()
	}
}
