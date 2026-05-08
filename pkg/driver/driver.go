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

func GetDriver(driverName string, size int) error {
	var open gorm.Dialector
	dsn := SqlConf

	// 配置dsn
	switch driverName {
	case "mysql":
		open = mysql.Open(dsn.getMysqlDsn())
	case "postgres":
		open = postgres.Open(dsn.getPostgresDsn())
	case "clickhouse":
		clickhouse.Open(dsn.getClickhouseDsn())
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
