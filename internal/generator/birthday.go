package generator

import (
	"time"
)

// 以下通过接口实现age birthday zodiac 相对应的数据
type BirthdayInfo struct {
	Birthday string
	Age      int
	Zodiac   string
}

func (g *Random) Birthday() *BirthdayInfo {
	age := 18 + g.IntN(63)
	year := time.Now().Year() - age
	month := time.Month(1 + g.IntN(12))
	maxDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	day := 1 + g.IntN(maxDay)

	birthday := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	md := int(birthday.Month())*100 + birthday.Day()
	var zodiac string
	switch {
	case md >= 321 && md <= 419:
		zodiac = "Aries"
	case md >= 420 && md <= 520:
		zodiac = "Taurus"
	case md >= 521 && md <= 621:
		zodiac = "Gemini"
	case md >= 622 && md <= 722:
		zodiac = "Cancer"
	case md >= 723 && md <= 822:
		zodiac = "Leo"
	case md >= 823 && md <= 922:
		zodiac = "Virgo"
	case md >= 923 && md <= 1023:
		zodiac = "Libra"
	case md >= 1024 && md <= 1122:
		zodiac = "Scorpio"
	case md >= 1123 && md <= 1221:
		zodiac = "Sagittarius"
	case md >= 1222 || md <= 119:
		zodiac = "Capricorn"
	case md >= 120 && md <= 218:
		zodiac = "Aquarius"
	default:
		zodiac = "Pisces"
	}

	return &BirthdayInfo{
		Age:      age,
		Birthday: birthday.Format("2006-01-02"),
		Zodiac:   zodiac,
	}
}
