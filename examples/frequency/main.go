package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/mroth/weightedrand"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) // always seed random!

	c, err := weightedrand.NewChooser(
		weightedrand.NewChoice('🍒', 0),
		weightedrand.NewChoice('🍋', 1),
		weightedrand.NewChoice('🍊', 1),
		weightedrand.NewChoice('🍉', 3),
		weightedrand.NewChoice('🥑', 5),
	)
	if err != nil {
		log.Fatal(err)
	}

	/* Let's pick a bunch of fruits so we can see the distribution in action! */
	fruits := make([]rune, 40*18)
	for i := 0; i < len(fruits); i++ {
		fruits[i] = c.Pick()
	}
	fmt.Println(string(fruits))

	/* That should have printed 🍋 and 🍊 with 0.1 probability, 🍉 with 0.3
	probability, and 🥑 with 0.5 probability. 🍒 should never be printed. (Note
	the weights don't have to add up to 10, that was just done here to make the
	example easier to read.) */
	freqs := make(map[rune]int)
	for _, f := range fruits {
		freqs[f]++
	}
	fmt.Printf("\n🍒: %d\t🍋: %d\t🍊: %d\t🍉: %d\t🥑: %d\n",
		freqs['🍒'], freqs['🍋'], freqs['🍊'], freqs['🍉'], freqs['🥑'])
}
