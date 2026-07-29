package generator

import (
	"math"
	"math/rand/v2"
	"strconv"

	"gendata/internal/generator/data"
)

func pick(key string) string {
	list := data.Name[key]
	return list[rand.IntN(len(list))]
}

func Name() string {
	if rand.IntN(2) == 0 {
		return pick("first") + " " + pick("middle") + " " + pick("last")
	}
	return pick("first") + " " + pick("last")
}

func Gender() string {
	genders := []string{"female", "male", "other"}
	return genders[rand.IntN(len(genders))]
}

func Phone() string {
	// 美国电话号码格式：+1-XXX-XXX-XXXX
	areaCode := 100 + rand.IntN(900)
	prefix := 100 + rand.IntN(900)
	lineNumber := 1000 + rand.IntN(9000)
	return "+1-" + strconv.Itoa(areaCode) + "-" + strconv.Itoa(prefix) + "-" + strconv.Itoa(lineNumber)
}

func Ethnicity() string {
	return data.Ethnicity[rand.IntN(len(data.Ethnicity))]
}

func Email() string {
	users := []string{pick("first"), pick("middle"), pick("last")}
	domains := "@" + data.Email[rand.IntN(len(data.Email))]
	return users[rand.IntN(3)] + domains
}

func round(f float64, n int) float64 {
	p := math.Pow(10, float64(n))
	return math.Round(f*p) / p
}

// GenerateHeight 生成身高（cm）
func Height() float64 {
	val := 150.0 + rand.Float64()*50.0
	return round(val, 0)
}

// GenerateWeight 生成体重（kg）
func Weight() float64 {
	val := 40.0 + rand.Float64()*80.0
	return round(val, 2)
}

func BloodType() string {
	bloodTypes := []string{"O", "A", "B", "AB", "other"}
	return bloodTypes[rand.IntN(len(bloodTypes))]
}
