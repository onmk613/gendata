package generator

import (
	"gendata/pkg/generator/data"
	"math/rand/v2"
)

func Currency() string {
	return data.Currency[rand.IntN(len(data.Currency))]
}
