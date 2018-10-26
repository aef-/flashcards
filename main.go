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
	LoadConfig() t.Config
	Prompt(card *t.Card, config t.Config) bool // true - correct
}

type Storage interface {
	LoadDefinitions(config t.Config) (definitions *t.Definitions, session *t.Session)
	SaveSession(cards *t.Cards, config t.Config) error
}

func main() {
	ui := getUI("cli")
	ui.Setup()
	config := ui.LoadConfig()

	storage := getStorage(config.StorageName)
	definitions, session := storage.LoadDefinitions(config)
	cards := collateDefinitionsWithSession(definitions, session)

	strategy := getStrategy(config.StrategyName)
	cards = strategy.Sort(&cards, config)
	for _, card := range cards {
		isCorrect := ui.Prompt(card, config)
		if isCorrect {
			strategy.Correct(card)
		} else {
			strategy.Incorrect(card)
			cards = append(cards, card)
		}
	}
	storage.SaveSession(&cards, config)
	for _, card := range cards {
		fmt.Println(card)
	}
}

func getUI(uiName string) UI {
	switch uiName {
	case "cli":
		return &ui.Cli{}
	}

	return nil
}

func getStorage(storageName string) Storage {
	switch storageName {
	case "local":
		return &storage.Local{}
	}

	return nil
}

func getStrategy(strategyName string) Strategy {
	switch strategyName {
	case "leitner":
		return &strategy.Leitner{}
	}

	return nil
}

func collateDefinitionsWithSession(definitions *t.Definitions, session *t.Session) t.Cards {
	var cards t.Cards
	for front, back := range *definitions {
		session := (*session)[front]
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
