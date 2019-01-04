package strategy

import (
	t "github.com/aef-/flashcards/types"
	"time"
)

var (
	timeNow = time.Now
)

type Leitner struct{}

func (l *Leitner) Name() string {
	return "Leitner"
}

func (l *Leitner) Correct(card *t.Card) {
	card.LastSeen = timeNow()
	card.Box = card.Box + 1
	card.Velocity = 1
}

func (l *Leitner) Incorrect(card *t.Card) {
	card.LastSeen = timeNow()
	card.Velocity = 1 - card.Box
	card.Box = 1

}

func (l *Leitner) Sort(cards *t.Cards, config t.Config) t.Cards {
	newCards := make(t.Cards, 0, config.CardsToIntroducePerDay)
	reviewCards := make(t.Cards, 0, config.CardsToReviewPerDay)
	now := timeNow()

	for _, card := range *cards {
		if len(newCards) == config.CardsToIntroducePerDay && len(reviewCards) > config.CardsToReviewPerDay {
			break
		}

		if card.Box == 0 && len(newCards) < config.CardsToIntroducePerDay {
			newCards = append(newCards, card)
			continue
		}
		if card.Box > 0 &&
			len(reviewCards) < config.CardsToReviewPerDay &&
			now.After(card.LastSeen.Add(time.Duration(config.Steps[card.Box]))) {
			reviewCards = append(reviewCards, card)
			continue
		}
	}

	studyCards := append(newCards, reviewCards...)
	// utils.ShuffleCards(studyCards)
	return studyCards
}
