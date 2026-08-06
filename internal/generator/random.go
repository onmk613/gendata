package generator

import "math/rand/v2"

// Random 是带独立随机数源的数据生成器。
// 每个写入 worker 持有一个实例，避免所有 worker 竞争全局随机源。
type Random struct {
	r *rand.Rand
}

func NewRandom(seed uint64) *Random {
	return &Random{r: rand.New(rand.NewPCG(seed, 0))}
}

func (g *Random) IntN(n int) int {
	if n <= 0 {
		return 0
	}
	return g.r.IntN(n)
}

func (g *Random) Float64() float64 {
	return g.r.Float64()
}
