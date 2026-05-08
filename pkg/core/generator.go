package core

import (
	"gendata/pkg/generator"
	"math/rand/v2"
	"time"

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
	JsonData     string  `gorm:"column:json_data;type:jsonb"`
	MarkdownData string  `gorm:"column:markdown_data;type:text"`
	OnlineStatus bool    `gorm:"column:online_status;type:boolean"`
	CreatedAt    string  `gorm:"column:created_at"`
	UpdatedAt    string  `gorm:"column:updated_at"`
}

func (DefaultTableRow) TableName() string {
	return TableName
}

func GenerateDefaultTableData(count int) []*DefaultTableRow {
	rows := make([]*DefaultTableRow, count)
	for i := 0; i < count; i++ {
		generator := newDefaultTableGenerator()
		rows[i] = generator.generateRow()
	}
	return rows
}

type defaultTableGenerator struct {
}

func newDefaultTableGenerator() *defaultTableGenerator {
	return &defaultTableGenerator{}
}

func (dtg *defaultTableGenerator) generateRow() *DefaultTableRow {
	// 出生日期
	birthday := generator.Birthday()
	// 地址
	address := generator.Address()

	return &DefaultTableRow{
		UserID:       uuid.New().String(),
		Name:         generator.Name(),
		Phone:        generator.Phone(),
		Gender:       generator.Gender(),
		BloodType:    generator.BloodType(),
		Birthday:     birthday.Birthday,
		Age:          birthday.Age,
		Zodiac:       birthday.Zodiac,
		Ethnicity:    generator.Ethnicity(),
		Email:        generator.Email(),
		Currency:     generator.Currency(),
		Address:      address.Address,
		Country:      address.Country,
		State:        address.State,
		ZipCode:      address.Zip,
		Height:       generator.Height(),
		Weight:       generator.Weight(),
		OnlineStatus: rand.Float64() < 0.4,
		Emoji:        generator.Emoji().Emoji,
		JsonData:     generator.JSONString(),
		MarkdownData: generator.MarkdownString(),
		CreatedAt:    time.Now().Add(-time.Duration(rand.IntN(8784)) * time.Hour).Format("2006-01-02 15:04:05.000"),
		UpdatedAt:    time.Now().Format("2006-01-02 15:04:05.000"),
	}
}
