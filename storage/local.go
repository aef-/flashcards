package storage

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"

	t "github.com/aef-/flashcards/types"
)

type Local struct{}

var savePath = flag.String("save-path", "/Users/207751/.gocards", "where to store ")

// Collate new definitions with those already in a learning session
func (s Local) LoadDefinitions() (t.Definitions, t.Session) {
	flag.Parse()
	args := flag.Args()
	definitionsFilename := args[0]

	definitionsYaml, err := ioutil.ReadFile(definitionsFilename)
	if err != nil {
		log.Fatal(err)
	}

	sessionPath := path.Join(*savePath, definitionsFilename)
	session := make(t.Session)
	_, err = os.Stat(sessionPath)
	if err == nil {
		sessionYaml, err := ioutil.ReadFile(sessionPath)
		if err != nil {
			log.Fatal(err)
		}

		err = yaml.Unmarshal(sessionYaml, &session)

		if err != nil {
			log.Fatal(err)
		}
	}

	definitions := make(t.Definitions)

	err = yaml.Unmarshal([]byte(definitionsYaml), &definitions)
	if err != nil {
		log.Fatal(err)
	}

	return definitions, session
}

func (s Local) saveSession(savePath string, fileName string, cards *t.Cards) {
	d, err := yaml.Marshal(*cards)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("Saving results to: %s\n", path.Join(savePath, fileName))
	_ = os.Mkdir(savePath, os.ModePerm)
	ioutil.WriteFile(path.Join(savePath, fileName), d, 0644)
}
