package generator

import (
	"gendata/internal/generator/data"
)

func (g *Random) Currency() string {
	return data.Currency[g.IntN(len(data.Currency))]
}
