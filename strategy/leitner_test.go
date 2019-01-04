package strategy

import (
	types "github.com/aef-/flashcards/types"
	"testing"
	"time"
)

/*
 * Need to test:
 *  various permutations of review/introduce cards, especially dates
 * 	steps
 * 	rest of config once implemented
 */
func TestSortsCardsCorrectly(t *testing.T) {
	timeNow = func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	}

	strategy := Leitner{}

	reviewCard1 := &types.Card{
		Box:      5,
		Velocity: 1,
		LastSeen: time.Date(2009, time.November, 9, 23, 0, 0, 0, time.UTC),
	}
	newCard1 := &types.Card{
		Box:      0,
		Velocity: 0,
		LastSeen: time.Date(2009, time.November, 9, 23, 0, 0, 0, time.UTC),
	}
	newCard2 := &types.Card{
		Box:      0,
		Velocity: 0,
	}
	reviewCard2 := &types.Card{
		Box:      1,
		Velocity: 1,
		LastSeen: time.Date(2009, time.November, 8, 23, 0, 0, 0, time.UTC),
	}
	reviewCard3 := &types.Card{
		Box:      1,
		Velocity: 1,
		LastSeen: time.Date(2009, time.November, 8, 23, 0, 0, 0, time.UTC),
	}

	cards := types.Cards{
		reviewCard1, newCard1, reviewCard2, newCard2, reviewCard3,
	}

	expectedCards := types.Cards{
		newCard1, reviewCard1, reviewCard2,
	}

	config := types.Config{
		CardsToIntroducePerDay: 1,
		CardsToReviewPerDay:    2,
		Steps:                  []int{0, 10, 1440, 4320, 7200, 43200},
	}

	cards = strategy.Sort(&cards, config)

	if len(cards) != len(expectedCards) {
		t.Error(
			"Number of cards mismatch, expected:", len(expectedCards),
			", but received:", len(cards),
		)
	}
	for index, card := range cards {
		if expectedCards[index] != card {
			t.Error(
				"At index", index, "\nExpected:\n", expectedCards[index],
				",\nReceived:\n", card,
			)
		}
	}
}

func TestMarksCardCorrect(t *testing.T) {
	lastSeen := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	timeNow = func() time.Time {
		return lastSeen
	}

	strategy := Leitner{}
	front := "front of card"
	back := "back of card"
	card := types.Card{
		Front:    front,
		Back:     back,
		Box:      3,
		Velocity: 2,
		LastSeen: time.Date(2009, time.November, 10, 21, 0, 0, 0, time.UTC),
	}

	expectedCard := types.Card{
		Front:    front,
		Back:     back,
		Box:      4,
		Velocity: 1,
		LastSeen: lastSeen,
	}

	strategy.Correct(&card)
	if card != expectedCard {
		t.Error(
			"\nEXPECTED\n", expectedCard,
			"\n\nRECEIVED\n", card,
		)
	}
}

func TestMarksCardIncorrect(t *testing.T) {
	lastSeen := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	timeNow = func() time.Time {
		return lastSeen
	}

	strategy := Leitner{}
	front := "front of card"
	back := "back of card"
	card := types.Card{
		Front:    front,
		Back:     back,
		Box:      3,
		Velocity: 2,
		LastSeen: time.Date(2009, time.November, 10, 21, 0, 0, 0, time.UTC),
	}

	expectedCard := types.Card{
		Front:    front,
		Back:     back,
		Box:      1,
		Velocity: -2,
		LastSeen: lastSeen,
	}

	strategy.Incorrect(&card)
	if card != expectedCard {
		t.Error(
			"\nEXPECTED\n", expectedCard,
			"\n\nRECEIVED\n", card,
		)
	}
}
