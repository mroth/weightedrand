package main

import (
	"fmt"
	"math/rand"
	"time"

	wr "github.com/mroth/weightedrand"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) // always seed random!

	c := wr.NewChooser(
		wr.Choice{Item: '🍆', Weight: 0},
		wr.Choice{Item: '🍋', Weight: 1},
		wr.Choice{Item: '🍊', Weight: 1},
		wr.Choice{Item: '🍉', Weight: 3},
		wr.Choice{Item: '🥑', Weight: 5},
	)

	/* Let's pick a bunch of fruits so we can see the distribution in action! */
	fruits := make([]rune, 40*18)
	for i := 0; i < len(fruits); i++ {
		fruits[i] = c.Pick().(rune)
	}
	fmt.Println(string(fruits))

	/* That should have printed 🍋 and 🍊 with 0.1 probability, 🍉 with 0.3
	probability, and 🥑 with 0.5 probability. 🍆 should never be printed. (Note
	the weights don't have to add up to 10, that was just done here to make the
	example easier to read.) */
	freqs := make(map[rune]int)
	for _, f := range fruits {
		freqs[f]++
	}
	fmt.Printf("\n🍆: %d\t🍋: %d\t🍊: %d\t🍉: %d\t🥑: %d\n",
		freqs['🍆'], freqs['🍋'], freqs['🍊'], freqs['🍉'], freqs['🥑'])
}
