package core

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type DefaultGenerator struct {
	rng *rand.Rand
}

func NewDefaultGenerator() *DefaultGenerator {
	return &DefaultGenerator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (dg *DefaultGenerator) randomAlphaNum(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[dg.rng.Intn(len(charset))]
	}
	return string(b)
}

// 随机生成一个类似名字的字符串，如 “Lirena”, “Tomir”, “Avelyn”
func (dg *DefaultGenerator) randomNamePart(minLen, maxLen int) string {
	letters := "abcdefghijklmnopqrstuvwxyz"
	length := minLen + dg.rng.Intn(maxLen-minLen+1)
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(letters[dg.rng.Intn(len(letters))])
	}
	name := sb.String()

	caser := cases.Title(language.English)
	return caser.String(name)
}

// 用户ID
func (dg *DefaultGenerator) GenerateUserID() string {
	id := uuid.New()
	return id.String()
}

// 人名
func (dg *DefaultGenerator) GenerateName() string {
	firstName := dg.randomNamePart(3, 7)
	lastName := dg.randomNamePart(4, 8)
	useMiddle := dg.rng.Intn(2) == 0
	if useMiddle {
		middleName := dg.randomNamePart(2, 6)
		return fmt.Sprintf("%s %s %s", firstName, middleName, lastName)
	}
	return fmt.Sprintf("%s %s", firstName, lastName)
}

// 生成电话号码，确保格式正确
func (dg *DefaultGenerator) GeneratePhone() string {
	// 美国电话号码格式：+1-XXX-XXX-XXXX
	areaCode := 100 + dg.rng.Intn(900)
	prefix := 100 + dg.rng.Intn(900)
	lineNumber := 1000 + dg.rng.Intn(9000)
	return fmt.Sprintf("+1-%03d-%03d-%04d", areaCode, prefix, lineNumber)
}

// GenerateGender 生成性别
func (dg *DefaultGenerator) GenerateGender() string {
	genders := []string{"female", "male", "other"}
	return genders[dg.rng.Intn(len(genders))]
}

// GenerateAge 生成年龄（18-80）
func (dg *DefaultGenerator) GenerateAge() int {
	return 18 + dg.rng.Intn(63)
}

// GenerateBirthday 生成出生日期
func (dg *DefaultGenerator) GenerateBirthday() time.Time {
	year := 2026 - (18 + dg.rng.Intn(63))
	month := 1 + dg.rng.Intn(12)
	day := 1 + dg.rng.Intn(28)
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// GenerateEmail 生成电子邮件，确保格式正确
func (dg *DefaultGenerator) GenerateEmail() string {
	domains := []string{"gmail.com", "outlook.com", "proton.me", "sfmail.cc", "xdander.com", "yahoo.com", "icloud.com", "aol.com", "zoho.com", "mail.com", "gmx.com", "yandex.com", "tutanota.com", "mailfence.com", "teleworm.us"}
	username := dg.randomAlphaNum(8)
	domain := domains[dg.rng.Intn(len(domains))]
	return fmt.Sprintf("%s@%s", username, domain)
}

// GenerateNationality 生成国籍
func (dg *DefaultGenerator) GenerateNationality() string {
	nationalities := []string{"US", "Canada", "UK", "France", "Germany", "Italy", "Spain", "France"}
	return nationalities[dg.rng.Intn(len(nationalities))]
}

// GenerateState 生成州/省
func (dg *DefaultGenerator) GenerateState() string {
	statesByCountry := map[string][]string{
		"US": {
			"California", "Texas", "New York", "Florida", "Illinois",
			"Ohio", "Pennsylvania", "Georgia", "Michigan", "Washington",
		},
		"Canada": {
			"Ontario", "Quebec", "British Columbia", "Alberta",
			"Manitoba", "Saskatchewan", "Nova Scotia", "New Brunswick",
		},
		"UK": {
			"England", "Scotland", "Wales", "Northern Ireland",
		},
		"France": {
			"Île-de-France", "Provence-Alpes-Côte d’Azur", "Occitanie", "Nouvelle-Aquitaine",
			"Grand Est", "Auvergne-Rhône-Alpes", "Brittany",
		},
		"Germany": {
			"Bavaria", "Berlin", "Brandenburg", "Hamburg", "Hesse",
			"Lower Saxony", "North Rhine-Westphalia", "Saxony", "Thuringia",
		},
		"Italy": {
			"Lazio", "Lombardy", "Sicily", "Tuscany", "Veneto", "Piedmont", "Emilia-Romagna",
		},
		"Spain": {
			"Madrid", "Catalonia", "Andalusia", "Valencia", "Galicia", "Basque Country",
		},
	}
	country := dg.GenerateNationality()
	states := statesByCountry[country]
	return states[dg.rng.Intn(len(states))]
}

// GenerateZipCode 生成邮政编码
func (dg *DefaultGenerator) GenerateZipCode() string {
	return fmt.Sprintf("%d", 10000+dg.rng.Intn(1000000))
}

func round(f float64, n int) float64 {
	p := math.Pow(10, float64(n))
	return math.Round(f*p) / p
}

// GenerateHeight 生成身高（cm）
func (dg *DefaultGenerator) GenerateHeight() float64 {
	val := 150.0 + dg.rng.Float64()*50.0
	return round(val, 0)
}

// GenerateWeight 生成体重（kg）
func (dg *DefaultGenerator) GenerateWeight() float64 {
	val := 40.0 + dg.rng.Float64()*80.0
	return round(val, 2)
}

// GenerateBloodType 生成血型
func (dg *DefaultGenerator) GenerateBloodType() string {
	bloodTypes := []string{"O", "A", "B", "AB"}
	return bloodTypes[dg.rng.Intn(len(bloodTypes))]
}

// GenerateAccount 生成账户
func (dg *DefaultGenerator) GenerateAccount() string {
	return "account_" + dg.randomAlphaNum(8)
}

// GenerateAccountName 生成账户名称
func (dg *DefaultGenerator) GenerateAccountName() string {
	return dg.randomNamePart(5, 10)
}

// GeneratePassword 生成密码（哈希后的），这里生成随机字符串表示
func (dg *DefaultGenerator) GeneratePassword() string {
	return dg.randomAlphaNum(16)
}

// GenerateOnlineStatus 生成在线状态
func (dg *DefaultGenerator) GenerateOnlineStatus() bool {
	return rand.Float64() < 0.6
}

// GenerateDateTime 生成日期时间
func (dg *DefaultGenerator) GenerateDateTime() time.Time {
	return time.Now().Add(-time.Duration(dg.rng.Intn(8784)) * time.Hour)
}
