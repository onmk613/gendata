package generator

import (
	"math/rand/v2"
	"strconv"

	"gendata/internal/generator/data"
)

// 地址 街道 单元 城市 州/省 邮编 国家
type AddressInfo struct {
	Address string `json:"address" xml:"address"`
	Street  string `json:"street" xml:"street"`
	Unit    string `json:"unit" xml:"unit"`
	City    string `json:"city" xml:"city"`
	State   string `json:"state" xml:"state"`
	Zip     string `json:"zip" xml:"zip"`
	Country string `json:"country" xml:"country"`
}

func Address() *AddressInfo {
	street := Street()
	city := City()
	state := State()
	zip := Zip()

	// 30% chance to include a unit in the address
	var unitStr string
	var unitField string
	if 1+rand.IntN(10) <= 3 {
		unitStr = ", " + Unit()
		unitField = Unit()
	}

	addressStr := street + unitStr + ", " + city + ", " + state + " " + zip

	return &AddressInfo{
		Address: addressStr,
		Street:  street,
		Unit:    unitField,
		City:    city,
		State:   state,
		Zip:     zip,
		Country: Country(),
	}
}

func Street() string {
	var street string
	switch rand.IntN(2) {
	case 0:
		street = streetNumber() + " " + streetPrefix() + " " + streetName() + streetSuffix()
	case 1:
		street = streetNumber() + " " + streetName() + streetSuffix()
	}

	return street
}

func streetNumber() string {
	return strconv.Itoa(100 + rand.IntN(10000))
}

func streetPrefix() string {
	prefixs := data.Address["street_prefix"]
	return prefixs[rand.IntN(len(prefixs))]
}

func streetName() string {
	names := data.Address["street_name"]
	return names[rand.IntN(len(names))]
}

func streetSuffix() string {
	suffixs := data.Address["street_suffix"]
	return suffixs[rand.IntN(len(suffixs))]
}

func Unit() string {
	unitTypes := data.Address["unit"]
	return unitTypes[rand.IntN(len(unitTypes))] + " " + strconv.Itoa(100+rand.IntN(1000))
}

func City() string {
	citys := data.Address["city"]
	return citys[rand.IntN(len(citys))]
}

func State() string {
	states := data.Address["state"]
	return states[rand.IntN(len(states))]
}

func Zip() string {
	return strconv.Itoa(10000 + rand.IntN(100000))
}

func Country() string {
	return data.Country[rand.IntN(len(data.Country))]
}
