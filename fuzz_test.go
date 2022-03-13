//go:build go1.18
// +build go1.18

package weightedrand

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

// Fuzz testing does not support slices as a corpus type in go1.18, thus we
// write a bunch of boilerplate here to allow us to encode []uint64 as []byte
// for kicks.

func bEncodeSlice(xs []uint64) []byte {
	bs := make([]byte, len(xs)*8)
	for i, x := range xs {
		n := i * 8
		binary.LittleEndian.PutUint64(bs[n:], x)
	}
	return bs
}

func bDecodeSlice(bs []byte) []uint64 {
	n := len(bs) / 8
	xs := make([]uint64, 0, n)
	for i := 0; i < n; i++ {
		x := binary.LittleEndian.Uint64(bs[8*i:])
		xs = append(xs, x)
	}
	return xs
}

// test our own encoder to make sure we didn't introduce errors.
func Test_bEncodeSlice(t *testing.T) {
	var testcases = [][]uint64{
		{},
		{1},
		{42},
		{912346},
		{1, 2},
		{1, 1, 1},
		{1, 2, 3},
		{1, 1000000},
		{1, 2, 3, 4, 5, 6, 7, 8, 9},
	}
	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			before := tc
			encoded := bEncodeSlice(before)
			if want, got := len(before)*8, len(encoded); want != got {
				t.Errorf("encoded length not as expected: want %d got %d", want, got)
			}
			decoded := bDecodeSlice(encoded)
			if !reflect.DeepEqual(before, decoded) {
				t.Errorf("want %v got %v", before, decoded)
			}
		})
	}
}

func FuzzNewChooser(f *testing.F) {
	var fuzzcases = [][]uint64{
		{},
		{0},
		{1},
		{1, 1},
		{1, 2, 3},
		{0, 1, 2},
	}
	for _, tc := range fuzzcases {
		f.Add(bEncodeSlice(tc))
	}

	f.Fuzz(func(t *testing.T, encodedWeights []byte) {
		weights := bDecodeSlice(encodedWeights)
		const sentinel = 1

		cs := make([]Choice[int, uint64], 0, len(weights))
		for _, w := range weights {
			cs = append(cs, Choice[int, uint64]{Item: sentinel, Weight: w})
		}

		// fuzz for error or panic on NewChooser
		c, err := NewChooser(cs...)
		if err != nil && !errors.Is(err, errNoValidChoices) && !errors.Is(err, errWeightOverflow) {
			t.Fatal(err)
		}

		if err == nil {
			result := c.Pick()      // fuzz for panic on Panic
			if result != sentinel { // fuzz for returned value unexpected (just use same non-zero sentinel value for all choices)
				t.Fatalf("expected %v got %v", sentinel, result)
			}
		}
	})
}
