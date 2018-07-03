package weightedrand // import "github.com/mroth/weightedrand"

import (
	"math/rand"
	"sort"
)

// Choice is a generic wrapper that can be used to add weights for any object
type Choice struct {
	Item   interface{}
	Weight int
}

// A Chooser caches the possible Choices, to improve performance on repeated calls
type Chooser struct {
	data   []Choice
	totals []int
	max    int
}

// NewChooser initializes a new Chooser, creating the caches
func NewChooser(cs ...Choice) Chooser {
	n := len(cs)
	sort.Slice(cs, func(i, j int) bool {
		return cs[i].Weight < cs[j].Weight
	})
	totals := make([]int, n, n)
	runningTotal := 0
	for i, c := range cs {
		runningTotal += c.Weight
		totals[i] = runningTotal
	}
	return Chooser{data: cs, totals: totals, max: runningTotal}
}

// Pick returns a single weighted random Choice.Item from the Chooser
func (chs Chooser) Pick() interface{} {
	r := rand.Intn(chs.max) + 1
	i := sort.SearchInts(chs.totals, r)
	return chs.data[i].Item
}
