package weightedrand

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
)

/******************************************************************************
*	Examples
*******************************************************************************/

// In this example, we create a Chooser to pick from amongst various emoji fruit
// runes. We assign a numeric weight to each choice. These weights are relative,
// not on any absolute scoring system. In this trivial case, we will assign a
// weight of 0 to all but one fruit, so that the output will be predictable.
func Example() {
	chooser, _ := NewChooser(
		NewChoice('🍋', 0),
		NewChoice('🍊', 0),
		NewChoice('🍉', 0),
		NewChoice('🥑', 42),
	)
	fruit := chooser.Pick()
	fmt.Printf("%c", fruit)
	//Output: 🥑
}

/******************************************************************************
*	Tests
*******************************************************************************/

const (
	testChoices    = 10
	testIterations = 1000000
)

func TestNewChooser(t *testing.T) {
	tests := []struct {
		name    string
		cs      []Choice[rune, int64]
		wantErr error
	}{
		{
			name:    "zero choices",
			cs:      []Choice[rune, int64]{},
			wantErr: errNoValidChoices,
		},
		{
			name:    "no choices with positive weight",
			cs:      []Choice[rune, int64]{{Item: 'a', Weight: 0}, {Item: 'b', Weight: 0}},
			wantErr: errNoValidChoices,
		},
		{
			name:    "choice with weight equals 1",
			cs:      []Choice[rune, int64]{{Item: 'a', Weight: 1}},
			wantErr: nil,
		},
		{
			name: "weight overflow",
			cs: []Choice[rune, int64]{
				{Item: 'a', Weight: math.MaxInt64/2 + 1},
				{Item: 'b', Weight: math.MaxInt64/2 + 1},
				{Item: 'c', Weight: math.MaxInt64/2 + 1},
				{Item: 'd', Weight: math.MaxInt64/2 + 1},
			},
			wantErr: errWeightOverflow,
		},
		{
			name:    "nominal case",
			cs:      []Choice[rune, int64]{{Item: 'a', Weight: 1}, {Item: 'b', Weight: 2}},
			wantErr: nil,
		},
		{
			name:    "negative weight case",
			cs:      []Choice[rune, int64]{{Item: 'a', Weight: 3}, {Item: 'b', Weight: -2}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewChooser(tt.cs...)
			if err != tt.wantErr {
				t.Errorf("NewChooser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil { // run a few Picks to make sure there are no panics
				for i := 0; i < 10; i++ {
					_ = c.Pick()
				}
			}
		})
	}

	u64tests := []struct {
		name    string
		cs      []Choice[rune, uint64]
		wantErr error
	}{
		{
			name:    "single uint64 equalling MaxUint64",
			cs:      []Choice[rune, uint64]{{Item: 'a', Weight: math.MaxUint64}},
			wantErr: errWeightOverflow,
		},
		{
			name: "single uint64 equalling MaxUint64 and a zero weight",
			cs: []Choice[rune, uint64]{
				{Item: 'a', Weight: math.MaxUint64},
				{Item: 'b', Weight: 0},
			},
			wantErr: errWeightOverflow,
		},
		{
			name: "multiple uint64s with sum MaxUint64",
			cs: []Choice[rune, uint64]{
				{Item: 'a', Weight: math.MaxUint64/2 + 1},
				{Item: 'b', Weight: math.MaxUint64/2 + 1},
			},
			wantErr: errWeightOverflow,
		},
	}
	for _, tt := range u64tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewChooser(tt.cs...)
			if err != tt.wantErr {
				t.Errorf("NewChooser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil { // run a few Picks to make sure there are no panics
				for i := 0; i < 10; i++ {
					_ = c.Pick()
				}
			}
		})
	}
}

// TestChooser_Pick assembles a list of Choices, weighted 0-9, and tests that
// over the course of 1,000,000 calls to Pick() each choice is returned more
// often than choices with a lower weight.
func TestChooser_Pick(t *testing.T) {
	choices := mockFrequencyChoices(t, testChoices)
	chooser, err := NewChooser(choices...)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("totals in chooser", chooser.totals)

	// run Pick() a million times, and record how often it returns each of the
	// possible choices.
	counts := make(map[int]int)
	for i := 0; i < testIterations; i++ {
		c := chooser.Pick()
		counts[c]++
	}

	verifyFrequencyCounts(t, counts, choices)
}

// Similar to what is used in randutil test, but in randomized order to avoid
// any issues with algorithms that are accidentally dependant on presorted data.
func mockFrequencyChoices(t *testing.T, n int) []Choice[int, int] {
	t.Helper()
	choices := make([]Choice[int, int], 0, n)
	list := rand.Perm(n)
	for _, v := range list {
		c := NewChoice(v, v)
		choices = append(choices, c)
	}
	t.Log("mocked choices of", choices)
	return choices
}

func verifyFrequencyCounts(t *testing.T, counts map[int]int, choices []Choice[int, int]) {
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
const BMMaxChoices = 10_000_000

func BenchmarkNewChooser(b *testing.B) {
	for n := BMMinChoices; n <= BMMaxChoices; n *= 10 {
		b.Run(fmt.Sprintf("size=%s", fmt1eN(n)), func(b *testing.B) {
			choices := mockChoices(n)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _ = NewChooser(choices...)
			}
		})
	}
}

func BenchmarkPick(b *testing.B) {
	for n := BMMinChoices; n <= BMMaxChoices; n *= 10 {
		b.Run(fmt.Sprintf("size=%s", fmt1eN(n)), func(b *testing.B) {
			choices := mockChoices(n)
			chooser, err := NewChooser(choices...)
			if err != nil {
				b.Fatal(err)
			}
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = chooser.Pick()
			}
		})
	}
}

func BenchmarkPickParallel(b *testing.B) {
	for n := BMMinChoices; n <= BMMaxChoices; n *= 10 {
		b.Run(fmt.Sprintf("size=%s", fmt1eN(n)), func(b *testing.B) {
			choices := mockChoices(n)
			chooser, err := NewChooser(choices...)
			if err != nil {
				b.Fatal(err)
			}
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_ = chooser.Pick()
				}
			})
		})
	}
}

func mockChoices(n int) []Choice[rune, int] {
	choices := make([]Choice[rune, int], 0, n)
	for i := 0; i < n; i++ {
		s := '🥑'
		w := rand.IntN(10)
		c := NewChoice(s, w)
		choices = append(choices, c)
	}
	return choices
}

// fmt1eN returns simplified order of magnitude scientific notation for n,
// e.g. "1e2" for 100, "1e7" for 10 million.
func fmt1eN(n int) string {
	return fmt.Sprintf("1e%d", int(math.Log10(float64(n))))
}
