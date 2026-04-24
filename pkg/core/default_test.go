package core

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestDefaultGeneratorUserID(t *testing.T) {
	gen := NewDefaultGenerator()
	userID := gen.GenerateUserID()

	if userID == "" {
		t.Error("Expected non-empty UserID")
	}

	// 验证UUID格式（36个字符，包含4个连字符）
	if len(userID) != 36 {
		t.Errorf("Expected UserID length 36, got %d", len(userID))
	}

	// 验证两次生成的ID不同
	userID2 := gen.GenerateUserID()
	if userID == userID2 {
		t.Error("Expected different UserIDs")
	}
}

func TestDefaultGeneratorName(t *testing.T) {
	gen := NewDefaultGenerator()

	for i := 0; i < 20; i++ {
		name := gen.GenerateName()

		if name == "" {
			t.Error("Expected non-empty Name")
		}

		// 名字应该包含至少一个空格
		if !strings.Contains(name, " ") {
			t.Errorf("Expected name to contain space, got: %s", name)
		}

		// 每个单词的首字母应该大写
		parts := strings.Fields(name)
		for _, part := range parts {
			if len(part) > 0 && (part[0] < 'A' || part[0] > 'Z') {
				t.Errorf("Expected capitalized name part, got: %s", part)
			}
		}
	}
}

func TestDefaultGeneratorPhone(t *testing.T) {
	gen := NewDefaultGenerator()
	phoneRegex := regexp.MustCompile(`^\+1-\d{3}-\d{3}-\d{4}$`)

	for i := 0; i < 20; i++ {
		phone := gen.GeneratePhone()

		if !phoneRegex.MatchString(phone) {
			t.Errorf("Invalid phone format: %s", phone)
		}
	}
}

func TestDefaultGeneratorGender(t *testing.T) {
	gen := NewDefaultGenerator()
	validGenders := map[string]bool{
		"female": true,
		"male":   true,
		"other":  true,
	}

	for i := 0; i < 50; i++ {
		gender := gen.GenerateGender()

		if !validGenders[gender] {
			t.Errorf("Invalid gender: %s", gender)
		}
	}
}

func TestDefaultGeneratorAge(t *testing.T) {
	gen := NewDefaultGenerator()

	for i := 0; i < 20; i++ {
		age := gen.GenerateAge()

		if age < 18 || age > 80 {
			t.Errorf("Age out of range: %d", age)
		}
	}
}

func TestDefaultGeneratorBirthday(t *testing.T) {
	gen := NewDefaultGenerator()
	now := time.Now()

	for i := 0; i < 20; i++ {
		birthday := gen.GenerateBirthday()

		// 检查日期是否有效
		if birthday.After(now) {
			t.Errorf("Birthday is in the future: %v", birthday)
		}

		// 计算年龄
		age := now.Year() - birthday.Year()
		if age < 18 || age > 80 {
			t.Errorf("Calculated age from birthday out of range: %d", age)
		}

		// 验证日期为UTC并且时间为00:00:00
		if birthday.Hour() != 0 || birthday.Minute() != 0 || birthday.Second() != 0 {
			t.Errorf("Birthday time should be 00:00:00, got %v", birthday)
		}
	}
}

func TestDefaultGeneratorEmail(t *testing.T) {
	gen := NewDefaultGenerator()
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9]+@[a-z.]+$`)
	validDomains := map[string]bool{
		"gmail.com":     true,
		"outlook.com":   true,
		"proton.me":     true,
		"sfmail.cc":     true,
		"xdander.com":   true,
		"yahoo.com":     true,
		"icloud.com":    true,
		"aol.com":       true,
		"zoho.com":      true,
		"mail.com":      true,
		"gmx.com":       true,
		"yandex.com":    true,
		"tutanota.com":  true,
		"mailfence.com": true,
		"teleworm.us":   true,
	}

	for i := 0; i < 20; i++ {
		email := gen.GenerateEmail()

		if !emailRegex.MatchString(email) {
			t.Errorf("Invalid email format: %s", email)
		}

		// 检查域名是否有效
		parts := strings.Split(email, "@")
		if len(parts) != 2 {
			t.Errorf("Email should have exactly one @: %s", email)
		}
		if !validDomains[parts[1]] {
			t.Errorf("Invalid email domain: %s", parts[1])
		}
	}
}

func TestDefaultGeneratorNationality(t *testing.T) {
	gen := NewDefaultGenerator()
	validNationalities := map[string]bool{
		"US":      true,
		"Canada":  true,
		"UK":      true,
		"France":  true,
		"Germany": true,
		"Italy":   true,
		"Spain":   true,
	}

	for i := 0; i < 50; i++ {
		nationality := gen.GenerateNationality()

		if !validNationalities[nationality] {
			t.Errorf("Invalid nationality: %s", nationality)
		}
	}
}

func TestDefaultGeneratorState(t *testing.T) {
	gen := NewDefaultGenerator()

	for i := 0; i < 50; i++ {
		state := gen.GenerateState()

		if state == "" {
			t.Error("Expected non-empty state")
		}
	}
}

func TestDefaultGeneratorZipCode(t *testing.T) {
	gen := NewDefaultGenerator()
	zipRegex := regexp.MustCompile(`^\d{5,7}$`)

	for i := 0; i < 20; i++ {
		zipCode := gen.GenerateZipCode()

		if !zipRegex.MatchString(zipCode) {
			t.Errorf("Invalid zip code format: %s", zipCode)
		}
	}
}

func TestDefaultGeneratorHeight(t *testing.T) {
	gen := NewDefaultGenerator()

	for i := 0; i < 20; i++ {
		height := gen.GenerateHeight()

		if height < 150 || height > 200 {
			t.Errorf("Height out of range: %f", height)
		}

		// 高度应该是整数（四舍五入到0位小数）
		if height != float64(int(height)) {
			t.Errorf("Height should be an integer: %f", height)
		}
	}
}

func TestDefaultGeneratorWeight(t *testing.T) {
	gen := NewDefaultGenerator()

	for i := 0; i < 20; i++ {
		weight := gen.GenerateWeight()

		if weight < 40 || weight > 120 {
			t.Errorf("Weight out of range: %f", weight)
		}
	}
}

func TestDefaultGeneratorBloodType(t *testing.T) {
	gen := NewDefaultGenerator()
	validBloodTypes := map[string]bool{
		"O":  true,
		"A":  true,
		"B":  true,
		"AB": true,
	}

	for i := 0; i < 50; i++ {
		bloodType := gen.GenerateBloodType()

		if !validBloodTypes[bloodType] {
			t.Errorf("Invalid blood type: %s", bloodType)
		}
	}
}

func TestDefaultGeneratorAccount(t *testing.T) {
	gen := NewDefaultGenerator()
	accountRegex := regexp.MustCompile(`^account_[a-zA-Z0-9]{8}$`)

	for i := 0; i < 20; i++ {
		account := gen.GenerateAccount()

		if !accountRegex.MatchString(account) {
			t.Errorf("Invalid account format: %s", account)
		}
	}
}

func TestDefaultGeneratorAccountName(t *testing.T) {
	gen := NewDefaultGenerator()

	for i := 0; i < 20; i++ {
		accountName := gen.GenerateAccountName()

		if accountName == "" {
			t.Error("Expected non-empty account name")
		}

		// 首字母应该大写
		if accountName[0] < 'A' || accountName[0] > 'Z' {
			t.Errorf("Expected capitalized account name, got: %s", accountName)
		}
	}
}

func TestDefaultGeneratorPassword(t *testing.T) {
	gen := NewDefaultGenerator()
	passwordRegex := regexp.MustCompile(`^[a-zA-Z0-9]{16}$`)

	for i := 0; i < 20; i++ {
		password := gen.GeneratePassword()

		if !passwordRegex.MatchString(password) {
			t.Errorf("Invalid password format: %s", password)
		}
	}
}

func TestDefaultGeneratorOnlineStatus(t *testing.T) {
	gen := NewDefaultGenerator()
	trueCount := 0
	totalCount := 100

	for i := 0; i < totalCount; i++ {
		if gen.GenerateOnlineStatus() {
			trueCount++
		}
	}

	// 应该大约60%为真
	percentage := float64(trueCount) / float64(totalCount)
	if percentage < 0.4 || percentage > 0.8 {
		t.Logf("Online status distribution: %f (expected ~0.6)", percentage)
	}
}

func TestDefaultGeneratorDateTime(t *testing.T) {
	gen := NewDefaultGenerator()
	now := time.Now()

	for i := 0; i < 20; i++ {
		dateTime := gen.GenerateDateTime()

		if dateTime.After(now) {
			t.Errorf("DateTime should be in the past: %v", dateTime)
		}

		// 应该在最近一年内
		yearAgo := now.Add(-365 * 24 * time.Hour)
		if dateTime.Before(yearAgo) {
			t.Logf("DateTime is more than a year ago: %v", dateTime)
		}
	}
}

func TestDefaultGeneratorRandomAlphaNum(t *testing.T) {
	gen := NewDefaultGenerator()

	for _, length := range []int{5, 10, 20} {
		str := gen.randomAlphaNum(length)

		if len(str) != length {
			t.Errorf("Expected length %d, got %d", length, len(str))
		}

		// 检查字符集
		for _, ch := range str {
			if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')) {
				t.Errorf("Invalid character in alphanumeric string: %c", ch)
			}
		}
	}
}

func TestDefaultGeneratorRandomNamePart(t *testing.T) {
	gen := NewDefaultGenerator()

	for i := 0; i < 20; i++ {
		namePart := gen.randomNamePart(3, 8)

		if len(namePart) < 3 || len(namePart) > 8 {
			t.Errorf("Name part length out of range: %d", len(namePart))
		}

		// 首字母应该大写
		if namePart[0] < 'A' || namePart[0] > 'Z' {
			t.Errorf("Expected capitalized name part, got: %s", namePart)
		}

		// 其他字母应该是小写
		for _, ch := range namePart[1:] {
			if ch < 'a' || ch > 'z' {
				t.Errorf("Expected lowercase letters, got: %c", ch)
			}
		}
	}
}
