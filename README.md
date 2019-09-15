# weightedrand

[![Build Status](https://github.com/mroth/weightedrand/workflows/Test/badge.svg)](https://github.com/mroth/weightedrand/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/mroth/weightedrand)](https://goreportcard.com/report/github.com/mroth/weightedrand)
[![GoDoc](https://godoc.org/github.com/mroth/weightedrand?status.svg)](https://godoc.org/github.com/mroth/weightedrand)

Randomly select an element from some kind of list, with the chances of each
element to be selected not being equal, but defined by relative "weights" (or
probabilities). This is called weighted random selection.

The existing Go library that has a generic implementation of this is
[`github.com/jmcvetta/randutil`][1], which optimizes for the single operation
case. In contrast, this library creates a presorted cache optimized for binary
search, allowing repeated selections from the same set to be significantly
faster, especially for large data sets.

[1]: https://github.com/jmcvetta/randutil

## Usage

```go
import (
    /* ...snip... */
    wr "github.com/mroth/weightedrand"
)

func main() {
    rand.Seed(time.Now().UTC().UnixNano()) // always seed random!

    c := wr.NewChooser(
        wr.Choice{Item: "🍆", Weight: 0},
        wr.Choice{Item: "🍋", Weight: 1},
        wr.Choice{Item: "🍊", Weight: 1},
        wr.Choice{Item: "🍉", Weight: 3},
        wr.Choice{Item: "🥑", Weight: 5},
    )
    /* The following will print 🍋 and 🍊 with 0.1 probability, 🍉 with 0.3
    probability, and 🥑 with 0.5 probability. 🍆 will never be printed. (Note
    the weights don't have to add up to 10, that was just done here to make the
    example easier to read.) */
    result := c.Pick().(string)
    fmt.Println(result)
}
```

## Benchmarks
Comparison of this library versus `randutil.ChooseWeighted`. For large numbers
of samplings from large collections, `weightedrand` will be quicker.

| Num choices |    `randutil` | `weightedrand` |
| ----------: | ------------: | -------------: |
|          10 |     435 ns/op |       58 ns/op |
|         100 |     511 ns/op |       84 ns/op |
|       1,000 |    1297 ns/op |      112 ns/op |
|      10,000 |    7952 ns/op |      137 ns/op |
|     100,000 |   85142 ns/op |      173 ns/op |
|   1,000,000 | 2082248 ns/op |      312 ns/op |

Don't be mislead by these numbers into thinking `weightedrand` is always the
right choice! If you are only picking from the same distribution once,
`randutil` will be faster. `weightedrand` optimizes for repeated calls at the
expense of some setup time and memory storage.

## Caveats

Note this uses `math/rand` instead of `crypto/rand`, as it is optimized for
performance, not cryptographically secure implementation.

Relies on global rand for determinism, therefore, don't forget to seed random!

## Credits

The algorithm used in this library (as well as the one used in randutil) comes
from:
https://eli.thegreenplace.net/2010/01/22/weighted-random-generation-in-python/
