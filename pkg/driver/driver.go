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

func GetDriver(driverName string) error {
	var open gorm.Dialector
	dsn := SqlConf

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
	DB = db
	return nil
}

func CloseDB() {
	if DB != nil {
		sqlDB, _ := DB.DB()
		sqlDB.Close()
	}
}
