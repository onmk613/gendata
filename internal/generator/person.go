package generator

import (
	"math"
	"strconv"

	"gendata/internal/generator/data"
)

func (g *Random) pick(key string) string {
	list := data.Name[key]
	return list[g.IntN(len(list))]
}

func (g *Random) Name() string {
	if g.IntN(2) == 0 {
		return g.pick("first") + " " + g.pick("middle") + " " + g.pick("last")
	}
	return g.pick("first") + " " + g.pick("last")
}

func (g *Random) Gender() string {
	genders := []string{"female", "male", "other"}
	return genders[g.IntN(len(genders))]
}

func (g *Random) Phone() string {
	// 美国电话号码格式：+1-XXX-XXX-XXXX
	areaCode := 100 + g.IntN(900)
	prefix := 100 + g.IntN(900)
	lineNumber := 1000 + g.IntN(9000)
	return "+1-" + strconv.Itoa(areaCode) + "-" + strconv.Itoa(prefix) + "-" + strconv.Itoa(lineNumber)
}

func (g *Random) Ethnicity() string {
	return data.Ethnicity[g.IntN(len(data.Ethnicity))]
}

func (g *Random) Email() string {
	users := []string{g.pick("first"), g.pick("middle"), g.pick("last")}
	domains := "@" + data.Email[g.IntN(len(data.Email))]
	return users[g.IntN(3)] + domains
}

func round(f float64, n int) float64 {
	p := math.Pow(10, float64(n))
	return math.Round(f*p) / p
}

// GenerateHeight 生成身高（cm）
func (g *Random) Height() float64 {
	val := 150.0 + g.Float64()*50.0
	return round(val, 0)
}

// GenerateWeight 生成体重（kg）
func (g *Random) Weight() float64 {
	val := 40.0 + g.Float64()*80.0
	return round(val, 2)
}

func (g *Random) BloodType() string {
	bloodTypes := []string{"O", "A", "B", "AB", "other"}
	return bloodTypes[g.IntN(len(bloodTypes))]
}
