package utils

import (
	"math/rand"
	"time"

	t "github.com/aef-/flashcards/types"
)

var (
	timeNow = time.Now
)

func ShuffleCards(arr t.Cards) {
	r := rand.New(rand.NewSource(timeNow().Unix()))
	for n := len(arr); n > 0; n-- {
		randIndex := r.Intn(n)
		arr[n-1], arr[randIndex] = arr[randIndex], arr[n-1]
	}
}
