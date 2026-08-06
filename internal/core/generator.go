package core

import (
	"time"

	"gendata/internal/generator"

	"github.com/google/uuid"
)

var TableName string

// 默认的数据结构
type DefaultTableRow struct {
	ID           int     `gorm:"primaryKey;autoIncrement"`
	UserID       string  `gorm:"column:user_id;uniqueIndex"`
	Name         string  `gorm:"column:name"`
	Phone        string  `gorm:"column:phone"`
	Gender       string  `gorm:"column:gender"`
	BloodType    string  `gorm:"column:blood_type"`
	Birthday     string  `gorm:"column:birthday"`
	Age          int     `gorm:"column:age"`
	Zodiac       string  `gorm:"column:zodiac"`
	Ethnicity    string  `gorm:"column:ethnicity"`
	Email        string  `gorm:"column:email"`
	Currency     string  `gorm:"column:currency"`
	Address      string  `gorm:"column:address"`
	Country      string  `gorm:"column:country"`
	State        string  `gorm:"column:state"`
	ZipCode      string  `gorm:"column:zip_code"`
	Height       float64 `gorm:"column:height"`
	Weight       float64 `gorm:"column:weight"`
	Emoji        string  `gorm:"column:emoji"`
	JsonData     string  `gorm:"column:json_data"`
	MarkdownData string  `gorm:"column:markdown_data"`
	OnlineStatus bool    `gorm:"column:online_status"`
	CreatedAt    string  `gorm:"column:created_at"`
	UpdatedAt    string  `gorm:"column:updated_at"`
}

func (DefaultTableRow) TableName() string {
	return TableName
}

func GenerateDefaultTableData(count int, rng *generator.Random) []*DefaultTableRow {
	if rng == nil {
		rng = generator.NewRandom(uint64(time.Now().UnixNano()))
	}
	rows := make([]*DefaultTableRow, count)
	for i := 0; i < count; i++ {
		rows[i] = (&defaultTableGenerator{rand: rng}).generateRow()
	}
	return rows
}

type defaultTableGenerator struct {
	rand *generator.Random
}

func (dtg *defaultTableGenerator) generateRow() *DefaultTableRow {
	// 出生日期
	birthday := dtg.rand.Birthday()
	// 地址
	address := dtg.rand.Address()

	return &DefaultTableRow{
		UserID:       uuid.New().String(),
		Name:         dtg.rand.Name(),
		Phone:        dtg.rand.Phone(),
		Gender:       dtg.rand.Gender(),
		BloodType:    dtg.rand.BloodType(),
		Birthday:     birthday.Birthday,
		Age:          birthday.Age,
		Zodiac:       birthday.Zodiac,
		Ethnicity:    dtg.rand.Ethnicity(),
		Email:        dtg.rand.Email(),
		Currency:     dtg.rand.Currency(),
		Address:      address.Address,
		Country:      address.Country,
		State:        address.State,
		ZipCode:      address.Zip,
		Height:       dtg.rand.Height(),
		Weight:       dtg.rand.Weight(),
		OnlineStatus: dtg.rand.Float64() < 0.4,
		Emoji:        dtg.rand.Emoji().Emoji,
		JsonData:     dtg.rand.JSONString(),
		MarkdownData: dtg.rand.MarkdownString(),
		CreatedAt:    time.Now().Add(-time.Duration(dtg.rand.IntN(8784)) * time.Hour).Format("2006-01-02 15:04:05.000"),
		UpdatedAt:    time.Now().Format("2006-01-02 15:04:05.000"),
	}
}
