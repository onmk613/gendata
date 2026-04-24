package driver

import (
	"time"

	"gorm.io/gorm"
)

type DefaultTableRow struct {
	ID           int       `gorm:"primaryKey;autoIncrement"`
	UserID       string    `gorm:"column:user_id;uniqueIndex;not null;size:50"`
	Name         string    `gorm:"column:name;not null;size:50"`
	Phone        string    `gorm:"column:phone;size:20"`
	Gender       string    `gorm:"column:gender;size:20"`
	Age          int       `gorm:"column:age;size:50"`
	Birthday     time.Time `gorm:"column:birthday;type:date"`
	Email        string    `gorm:"column:email;size:50"`
	Nationality  string    `gorm:"column:nationality;size:50"`
	State        string    `gorm:"column:state;size:255"`
	ZipCode      string    `gorm:"column:zip_code;size:20"`
	Height       float64   `gorm:"column:height;size:20"`
	Weight       float64   `gorm:"column:weight;size:20"`
	BloodType    string    `gorm:"column:blood_type;size:10"`
	Account      string    `gorm:"column:account;size:50"`
	AccountName  string    `gorm:"column:account_name;size:100"`
	Password     string    `gorm:"column:password;size:255"`
	OnlineStatus bool      `gorm:"column:online_status;type:boolean"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (DefaultTableRow) TableName() string {
	return SqlConf.Table
}

// PostgreSQL extended protocol 限制单个 SQL 最多 65535 个参数
// DefaultTableRow 有 20 个字段，所以最大安全批大小为 65535 / 20 = 3276
// 保守起见，我们使用 3000 作为上限
const PostgreSQLMaxParams = 65535
const DefaultTableRowColumns = 20
const SafeMaxBatchSizeForPostgreSQL = PostgreSQLMaxParams / DefaultTableRowColumns

func CreateTemplateDataBatchWithSize(templates []*DefaultTableRow, db *gorm.DB, batchSize int) error {
	actualBatchSize := batchSize
	if batchSize > SafeMaxBatchSizeForPostgreSQL {
		actualBatchSize = SafeMaxBatchSizeForPostgreSQL
	}

	result := db.CreateInBatches(templates, actualBatchSize)
	return result.Error
}
