// Package compbench is used to generate informal comparative benchmarks vs
// randutil.
package compbench

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/jmcvetta/randutil"
	"github.com/mroth/weightedrand/v2"
)

const BMMinChoices = 10
const BMMaxChoices = 10_000_000

func BenchmarkMultiple(b *testing.B) {
	for n := BMMinChoices; n <= BMMaxChoices; n *= 10 {
		b.Run(fmt.Sprintf("size=%s", fmt1eN(n)), func(b *testing.B) {
			wr_choices := mockChoices(b, n)
			ru_choices := convertChoices(b, wr_choices)

			b.Run("concurrency=single", func(b *testing.B) {
				b.Run("lib=randutil", func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						randutil.WeightedChoice(ru_choices)
					}
				})

				b.Run("lib=weightedrand", func(b *testing.B) {
					chs, err := weightedrand.NewChooser(wr_choices...)
					if err != nil {
						b.Fatal(err)
					}
					b.ResetTimer()

					for i := 0; i < b.N; i++ {
						chs.Pick()
					}
				})
			})

			b.Run("concurrency=parallel", func(b *testing.B) {
				b.Run("lib=weightedrand", func(b *testing.B) {
					chs, err := weightedrand.NewChooser(wr_choices...)
					if err != nil {
						b.Fatal(err)
					}
					b.ResetTimer()

					b.RunParallel(func(pb *testing.PB) {
						for pb.Next() {
							chs.Pick()
						}
					})
				})

				b.Run("lib=randutil", func(b *testing.B) {
					b.RunParallel(func(pb *testing.PB) {
						for pb.Next() {
							randutil.WeightedChoice(ru_choices)
						}
					})
				})
			})
		})
	}
}

// THE SINGLE USAGE CASE IS AN ANTI-PATTERN FOR THE INTENDED USAGE OF THIS
// LIBRARY. Provide some optional benchmarks for that to illustrate the point.
func BenchmarkSingle(b *testing.B) {
	if testing.Short() {
		b.Skip()
	}

	for n := BMMinChoices; n <= BMMaxChoices; n *= 10 {
		b.Run(fmt.Sprintf("size=%s", fmt1eN(n)), func(b *testing.B) {
			wr_choices := mockChoices(b, n)
			ru_choices := convertChoices(b, wr_choices)

			b.Run("lib=randutil", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					randutil.WeightedChoice(ru_choices)
				}
			})

			b.Run("lib=weightedrand", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					// never actually do this, this is not how the library is used
					chs, _ := weightedrand.NewChooser(wr_choices...)
					chs.Pick()
				}
			})
		})
	}
}

func mockChoices(tb testing.TB, n int) []weightedrand.Choice[rune, uint] {
	tb.Helper()
	choices := make([]weightedrand.Choice[rune, uint], 0, n)
	for i := 0; i < n; i++ {
		s := '🥑'
		w := rand.Intn(10)
		c := weightedrand.NewChoice(s, uint(w))
		choices = append(choices, c)
	}
	return choices
}

func convertChoices(tb testing.TB, cs []weightedrand.Choice[rune, uint]) []randutil.Choice {
	tb.Helper()
	res := make([]randutil.Choice, len(cs))
	for i, c := range cs {
		res[i] = randutil.Choice{Weight: int(c.Weight), Item: c.Item}
	}
	return res
}

// fmt1eN returns simplified order of magnitude scientific notation for n,
// e.g. "1e2" for 100, "1e7" for 10 million.
func fmt1eN(n int) string {
	return fmt.Sprintf("1e%d", int(math.Log10(float64(n))))
}
