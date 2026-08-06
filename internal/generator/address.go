package generator

import (
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

func (g *Random) Address() *AddressInfo {
	street := g.Street()
	city := g.City()
	state := g.State()
	zip := g.Zip()

	// 30% chance to include a unit in the address
	var unitStr string
	var unitField string
	if 1+g.IntN(10) <= 3 {
		unitStr = ", " + g.Unit()
		unitField = g.Unit()
	}

	addressStr := street + unitStr + ", " + city + ", " + state + " " + zip

	return &AddressInfo{
		Address: addressStr,
		Street:  street,
		Unit:    unitField,
		City:    city,
		State:   state,
		Zip:     zip,
		Country: g.Country(),
	}
}

func (g *Random) Street() string {
	var street string
	switch g.IntN(2) {
	case 0:
		street = g.streetNumber() + " " + g.streetPrefix() + " " + g.streetName() + g.streetSuffix()
	case 1:
		street = g.streetNumber() + " " + g.streetName() + g.streetSuffix()
	}

	return street
}

func (g *Random) streetNumber() string {
	return strconv.Itoa(100 + g.IntN(10000))
}

func (g *Random) streetPrefix() string {
	prefixs := data.Address["street_prefix"]
	return prefixs[g.IntN(len(prefixs))]
}

func (g *Random) streetName() string {
	names := data.Address["street_name"]
	return names[g.IntN(len(names))]
}

func (g *Random) streetSuffix() string {
	suffixs := data.Address["street_suffix"]
	return suffixs[g.IntN(len(suffixs))]
}

func (g *Random) Unit() string {
	unitTypes := data.Address["unit"]
	return unitTypes[g.IntN(len(unitTypes))] + " " + strconv.Itoa(100+g.IntN(1000))
}

func (g *Random) City() string {
	citys := data.Address["city"]
	return citys[g.IntN(len(citys))]
}

func (g *Random) State() string {
	states := data.Address["state"]
	return states[g.IntN(len(states))]
}

func (g *Random) Zip() string {
	return strconv.Itoa(10000 + g.IntN(100000))
}

func (g *Random) Country() string {
	return data.Country[g.IntN(len(data.Country))]
}
