package storage

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"

	t "github.com/aef-/flashcards/types"
)

type Local struct{}

// Collate new definitions with those already in a learning session
func (s Local) LoadDefinitions(config t.Config) (*t.Definitions, *t.Session) {
	definitionsFilename := config.FilePath

	definitionsYaml, err := ioutil.ReadFile(definitionsFilename)
	if err != nil {
		log.Fatal(err)
	}

	sessionPath := path.Join(config.SavePath, definitionsFilename)
	cards := make(t.Cards, 100)
	session := make(t.Session)
	_, err = os.Stat(sessionPath)
	if err == nil {
		sessionYaml, err := ioutil.ReadFile(sessionPath)
		if err != nil {
			log.Fatal(err)
		}

		err = yaml.Unmarshal(sessionYaml, &cards)

		if err != nil {
			log.Fatal(err)
		}

		for _, card := range cards {
			session[card.Front] = card
		}
	}

	definitions := make(t.Definitions)

	err = yaml.Unmarshal([]byte(definitionsYaml), &definitions)
	if err != nil {
		log.Fatal(err)
	}

	return &definitions, &session
}

func (s Local) SaveSession(cards *t.Cards, config t.Config) error {
	d, err := yaml.Marshal(*cards)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	_, err = os.Stat(config.SavePath)
	if err != nil {
		err = os.Mkdir(config.SavePath, os.ModePerm)
	}

	savePath := path.Join(config.SavePath, config.FilePath)
	fmt.Printf("Saving results to: %s\n", savePath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return ioutil.WriteFile(savePath, d, 0644)
}
