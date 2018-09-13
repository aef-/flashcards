package main

import (
	"fmt"
	"time"

	"github.com/aef-/flashcards/storage"
	"github.com/aef-/flashcards/strategy"
	t "github.com/aef-/flashcards/types"
	"github.com/aef-/flashcards/ui"
)

type Strategy interface {
	Name() string
	Correct(card *t.Card)
	Incorrect(card *t.Card)
	Sort(cards *t.Cards, config t.Config) t.Cards
}

type UI interface {
	Setup()
	LoadConfig(cards t.Cards)
	Prompt(card t.Card) bool // true - correct
}

type Storage interface {
	LoadDefinitions() (definitions *t.Definitions, session *t.Session)
	SaveSession(savePath string, fileName string, cards *t.Cards)
}

func main() {
	storage := storage.Local{}
	definitions, session := storage.LoadDefinitions()
	cards := collateDefinitionsWithSession(definitions, session)

	ui := ui.Cli{}
	ui.Setup()
	config := ui.LoadConfig()
	strategy := getStrategy(config.StrategyName)
	cards = strategy.Sort(&cards, config)
	for _, card := range cards {
		isCorrect := ui.Prompt(card, config)
		if isCorrect {
			strategy.Correct(card)
		} else {
			strategy.Incorrect(card)
		}
	}
	for _, card := range cards {
		fmt.Println(card)
	}
}

func getStrategy(strategyName string) Strategy {
	switch strategyName {
	case "leitner":
		return &strategy.Leitner{}
	}

	return nil
}

func collateDefinitionsWithSession(definitions t.Definitions, session t.Session) t.Cards {
	var cards t.Cards
	for front, back := range definitions {
		session := (session)[front]
		if session != nil {
			session.Front = front
			session.Back = back
			cards = append(cards, session)
		} else {
			cards = append(cards, &t.Card{
				Front:    front,
				Back:     back,
				Box:      0,
				LastSeen: time.Time{},
			})
		}
	}
	return cards
}
