package main

import (
	"math/rand"
)

type RandomQuoteStatic struct {
	randomizerIntner
	quotes []string
}

func NewRandomQuoteStatic(randomizer randomizerIntner, quotes []string) *RandomQuoteStatic {
	return &RandomQuoteStatic{randomizerIntner: randomizer, quotes: quotes}
}

func (r *RandomQuoteStatic) getRandom() string {
	return r.quotes[rand.Intn(len(r.quotes))] //no blocking because it's all static
}
