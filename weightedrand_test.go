package weightedrand

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

/******************************************************************************
*	Examples
*******************************************************************************/

// In this example, we create a Chooser to pick from amongst various emoji fruit
// runes. We assign a numeric weight to each choice. These weights are relative,
// not on any absolute scoring system. In this trivial case, we will assign a
// weight of 0 to all but one fruit, so that the output will be predictable.
func Example() {
	chooser := NewChooser(
		NewChoice('🍋', 0),
		NewChoice('🍊', 0),
		NewChoice('🍉', 0),
		NewChoice('🥑', 42),
	)
	fruit := chooser.Pick().(rune)
	fmt.Printf("%c", fruit)
	//Output: 🥑
}

/******************************************************************************
*	Tests
*******************************************************************************/

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

const (
	testChoices    = 10
	testIterations = 1000000
)

// TestChooser_Pick assembles a list of Choices, weighted 0-9, and tests that
// over the course of 1,000,000 calls to Pick() each choice is returned more
// often than choices with a lower weight.
func TestChooser_Pick(t *testing.T) {
	choices := mockFrequencyChoices(t, testChoices)
	chooser := NewChooser(choices...)
	t.Log("totals in chooser", chooser.totals)

	// run Pick() a million times, and record how often it returns each of the
	// possible choices.
	counts := make(map[int]int)
	for i := 0; i < testIterations; i++ {
		c := chooser.Pick()
		counts[c.(int)]++
	}

	verifyFrequencyCounts(t, counts, choices)
}

// TestChooser_PickSource is the same test methodology as TestChooser_Pick, but
// here we use the PickSource method and access the same chooser concurrently
// from multiple different goroutines, each providing its own source of
// randomness.
func TestChooser_PickSource(t *testing.T) {
	choices := mockFrequencyChoices(t, testChoices)
	chooser := NewChooser(choices...)
	t.Log("totals in chooser", chooser.totals)

	counts1 := make(map[int]int)
	counts2 := make(map[int]int)
	var wg sync.WaitGroup
	wg.Add(2)
	checker := func(counts map[int]int) {
		defer wg.Done()
		rs := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
		for i := 0; i < testIterations/2; i++ {
			c := chooser.PickSource(rs)
			counts[c.(int)]++
		}
	}
	go checker(counts1)
	go checker(counts2)
	wg.Wait()

	verifyFrequencyCounts(t, counts1, choices)
	verifyFrequencyCounts(t, counts2, choices)
}

// Similar to what is used in randutil test, but in randomized order to avoid
// any issues with algorithms that are accidentally dependant on presorted data.
func mockFrequencyChoices(t *testing.T, n int) []Choice {
	t.Helper()
	choices := make([]Choice, 0, n)
	list := rand.Perm(n)
	for _, v := range list {
		c := NewChoice(v, uint(v))
		choices = append(choices, c)
	}
	t.Log("mocked choices of", choices)
	return choices
}

func verifyFrequencyCounts(t *testing.T, counts map[int]int, choices []Choice) {
	t.Helper()

	// Ensure weight 0 results in no results
	if cczero := counts[0]; cczero != 0 {
		t.Error("Weight 0 results appeared nonzero times: ", cczero)
	}

	// Test that higher weighted choices were chosen more often than their lower
	// weighted peers.
	for i, c := range choices[0 : len(choices)-1] {
		next := choices[i+1]
		cw, nw := c.Weight, next.Weight
		if !(counts[int(cw)] < counts[int(nw)]) {
			t.Error("Value not lesser", cw, nw, counts[int(cw)], counts[int(nw)])
		}
	}
}

/******************************************************************************
*	Benchmarks
*******************************************************************************/

const BMMinChoices = 10
const BMMaxChoices = 1000000

func BenchmarkNewChooser(b *testing.B) {
	for n := BMMinChoices; n <= BMMaxChoices; n *= 10 {
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
	for n := BMMinChoices; n <= BMMaxChoices; n *= 10 {
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

func BenchmarkPickParallel(b *testing.B) {
	for n := BMMinChoices; n <= BMMaxChoices; n *= 10 {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			choices := mockChoices(n)
			chooser := NewChooser(choices...)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				rs := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
				for pb.Next() {
					chooser.PickSource(rs)
				}
			})
		})
	}
}

func mockChoices(n int) []Choice {
	choices := make([]Choice, 0, n)
	for i := 0; i < n; i++ {
		s := '🥑'
		w := rand.Intn(10)
		c := NewChoice(s, uint(w))
		choices = append(choices, c)
	}
	return choices
}
