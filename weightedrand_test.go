package weightedrand

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func mockChoices(n int) []Choice {
	choices := make([]Choice, 0, n)
	for i := 0; i < n; i++ {
		s := "⚽️"
		w := rand.Intn(10)
		c := Choice{Item: s, Weight: uint(w)}
		choices = append(choices, c)
	}
	return choices
}

// TestWeightedChoice assembles a list of choices, weighted 0-9, and tests that
// over the course of 1,000,000 calls to WeightedChoice() each choice is
// returned more often than choices with a lower weight.
func TestWeightedChoice(t *testing.T) {
	// Make weighted choices
	var choices []Choice
	chosenCount := make(map[int]int)

	/* Similar to what is used in randutil test, but in randomized order to
	avoid any issues with algorithms that are accidentally dependant on
	presorted data. */
	list := rand.Perm(10)
	for _, v := range list {
		c := Choice{Weight: uint(v), Item: v}
		choices = append(choices, c)
	}
	t.Log("FYI mocked choices of", choices)

	// Run WeightedChoice() a million times, and record how often it returns
	// each of the possible choices.
	chooser := NewChooser(choices...)
	t.Log("values in chooser", chooser.totals)
	for i := 0; i < 1000000; i++ {
		c := chooser.Pick()
		chosenCount[c.(int)]++
	}

	// Ensure weight 0 results in no results
	if cczero := chosenCount[0]; cczero != 0 {
		t.Error("Weight 0 results appeared nonzero times: ", cczero)
	}

	// Test that higher weighted choices were chosen more often than their lower
	// weighted peers.
	for i, c := range choices[0 : len(choices)-1] {
		next := choices[i+1]
		cw, nw := c.Weight, next.Weight
		if !(chosenCount[int(cw)] < chosenCount[int(nw)]) {
			t.Error("Value not lesser", cw, nw, chosenCount[int(cw)], chosenCount[int(nw)])
		}
	}

}

const BMminChoices = 10
const BMmaxChoices = 1000000

func BenchmarkNewChooser(b *testing.B) {
	for n := BMminChoices; n <= BMmaxChoices; n *= 10 {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			choices := mockChoices(n)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = NewChooser(choices...)
			}
		})
	}
}

func BenchmarkPick(b *testing.B) {
	for n := BMminChoices; n <= BMmaxChoices; n *= 10 {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			choices := mockChoices(n)
			chooser := NewChooser(choices...)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				chooser.Pick()
			}
		})
	}
}

// This following is a historic artifact from comparative benchmarking with
// randutil, however it is not critical to ongoing development.

// func BenchmarkRandutil(b *testing.B) {
// 	if testing.Short() {
// 		b.Skip()
// 	}
// 	for n := BMminChoices; n <= BMmaxChoices; n *= 10 {
// 		b.Run(strconv.Itoa(n), func(b *testing.B) {
// 			b.StopTimer()
// 			choices := mockChoices(n)
// 			choicesR := make([]randutil.Choice, len(choices), len(choices))
// 			for i, c := range choices {
// 				choicesR[i] = randutil.Choice{Weight: c.Weight, Item: c.Item}
// 			}
// 			b.StartTimer()

// 			for i := 0; i < b.N; i++ {
// 				randutil.WeightedChoice(choicesR)
// 			}
// 		})
// 	}
// }
