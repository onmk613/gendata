package generator

import (
	"math/rand/v2"

	"gendata/internal/generator/data"
)

func Currency() string {
	return data.Currency[rand.IntN(len(data.Currency))]
}
