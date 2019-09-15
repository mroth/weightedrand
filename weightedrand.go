// Package weightedrand contains a performant data structure and algorithm used
// to randomly select an element from some kind of list, where the chances of
// each element to be selected not being equal, but defined by relative
// "weights" (or probabilities). This is called weighted random selection.
//
// There is an existing Go library that has a generic implementation of this as
// github.com/jmcvetta/randutil, which optimizes for the single operation case.
// In contrast, this package creates a presorted cache optimized for binary
// search, allowing repeated selections from the same set to be significantly
// faster, especially for large data sets.
package weightedrand

import (
	"math/rand"
	"sort"
)

// Choice is a generic wrapper that can be used to add weights for any object
type Choice struct {
	Item   interface{}
	Weight uint
}

// A Chooser caches many possible Choices in a structure designed to improve
// performance on repeated calls for weighted random selection.
type Chooser struct {
	data   []Choice
	totals []int
	max    int
}

// NewChooser initializes a new Chooser consisting of the possible Choices.
func NewChooser(cs ...Choice) Chooser {
	sort.Slice(cs, func(i, j int) bool {
		return cs[i].Weight < cs[j].Weight
	})
	totals := make([]int, len(cs))
	runningTotal := 0
	for i, c := range cs {
		runningTotal += int(c.Weight)
		totals[i] = runningTotal
	}
	return Chooser{data: cs, totals: totals, max: runningTotal}
}

// Pick returns a single weighted random Choice.Item from the Chooser.
func (chs Chooser) Pick() interface{} {
	r := rand.Intn(chs.max) + 1
	i := sort.SearchInts(chs.totals, r)
	return chs.data[i].Item
}
