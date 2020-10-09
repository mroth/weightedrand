// Package compbench is used to generate informal comparative benchmarks vs
// randutil.
package compbench

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/jmcvetta/randutil"
	"github.com/mroth/weightedrand"
)

const BMMinChoices = 10
const BMMaxChoices = 1000000

func BenchmarkMultiple(b *testing.B) {
	b.Run("jmc_randutil", func(b *testing.B) {
		for n := BMMinChoices; n <= BMMaxChoices; n *= 10 {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				choices := convertChoices(b, mockChoices(b, n))
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					randutil.WeightedChoice(choices)
				}
			})
		}
	})

	b.Run("weightedrand", func(b *testing.B) {
		for n := BMMinChoices; n <= BMMaxChoices; n *= 10 {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				choices := mockChoices(b, n)
				chs := weightedrand.NewChooser(choices...)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					chs.Pick()
				}
			})
		}
	})

	b.Run("wr-parallel", func(b *testing.B) {
		for n := BMMinChoices; n <= BMMaxChoices; n *= 10 {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				choices := mockChoices(b, n)
				chs := weightedrand.NewChooser(choices...)
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					rs := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
					for pb.Next() {
						chs.PickSource(rs)
					}
				})
			})
		}
	})
}

// The single usage case is an anti-pattern for the intended usage of this
// library. Might as well keep some optional benchmarks for that to illustrate
// the point.
func BenchmarkSingle(b *testing.B) {
	if testing.Short() {
		b.Skip()
	}

	b.Run("jmc_randutil", func(b *testing.B) {
		for n := BMMinChoices; n <= BMMaxChoices; n *= 10 {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				choices := convertChoices(b, mockChoices(b, n))
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					randutil.WeightedChoice(choices)
				}
			})
		}
	})

	b.Run("weightedrand", func(b *testing.B) {
		for n := BMMinChoices; n <= BMMaxChoices; n *= 10 {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				choices := mockChoices(b, n)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					chs := weightedrand.NewChooser(choices...)
					chs.Pick()
				}
			})
		}
	})
}

func mockChoices(tb testing.TB, n int) []weightedrand.Choice {
	tb.Helper()
	choices := make([]weightedrand.Choice, 0, n)
	for i := 0; i < n; i++ {
		s := '🥑'
		w := rand.Intn(10)
		c := weightedrand.NewChoice(s, uint(w))
		choices = append(choices, c)
	}
	return choices
}

func convertChoices(tb testing.TB, cs []weightedrand.Choice) []randutil.Choice {
	tb.Helper()
	res := make([]randutil.Choice, len(cs))
	for i, c := range cs {
		res[i] = randutil.Choice{Weight: int(c.Weight), Item: c.Item}
	}
	return res
}
