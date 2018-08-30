package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
	"gopkg.in/yaml.v2"
)

type Cards = []*Card
type Card struct {
	Front    string `yaml:"-"`
	Back     string `yaml:"-"`
	Box      int    `yaml:"box"`
	LastSeen int64  `yaml:"last_seen"`
}

func (c Card) String() string {
	return fmt.Sprintf("Front: %s\nBack: %s\nBox: %d\nLast Seen: %d", c.Front, c.Back, c.Box, c.LastSeen)
}

type Strategy interface {
	Name() string
	Correct(card *Card)
	Incorrect(card *Card)
	Sort(cards *Cards)
}

var correctSleepMs = flag.Int("correct-time", 500, "time (ms) until next question is shown after a correct question")
var incorrectSleepMs = flag.Int("incorrect-time", 1500, "time (ms) until next question is shown after an incorrect question")
var savePath = flag.String("save-path", "/Users/207751/.gocards", "where to store ")
var strategy = flag.String("strategy", "all", "One of all|all-random|leitner")
var flip = flag.String("flip", "front", "Prompt with either front, back or random")
var untilCorrect = flag.Bool("until-correct", true, "Repeat words until they're correct")
var ignoreParenthesis = flag.Bool("ignore-paren", true, "Does not count text within parenthesis when comparing results")

func main() {

	log.SetPrefix("wait: ")
	log.SetFlags(0)

	flag.Parse()
	args := flag.Args()
	cardsFileName := args[0]

	cardsData, err := ioutil.ReadFile(cardsFileName)
	if err != nil {
		log.Fatal(err)
	}

	sessionPath := path.Join(*savePath, cardsFileName)
	sessionYaml := make(map[string]*Card)
	_, err = os.Stat(sessionPath)
	if err == nil {
		sessionData, err := ioutil.ReadFile(sessionPath)
		if err != nil {
			log.Fatal(err)
		}

		err = yaml.Unmarshal(sessionData, &sessionYaml)

		if err != nil {
			log.Fatal(err)
		}
	}

	cardsYaml := make(map[string]string)

	err = yaml.Unmarshal([]byte(cardsData), &cardsYaml)
	if err != nil {
		log.Fatal(err)
	}

	cards := loadCards(&cardsYaml, &sessionYaml)

	var strategy Strategy = strategy.Leitner{}

	cmd := exec.Command("tput", "civis")
	cmd.Stdout = os.Stdout
	cmd.Run()

	correctSleep := time.Duration(time.Millisecond * time.Duration(*correctSleepMs))
	incorrectSleep := time.Duration(time.Millisecond * time.Duration(*incorrectSleepMs))
	for _, card := range cards {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("> %s\n", card.Front)
		answer, _ := reader.ReadString('\n')
		if strings.TrimRight(answer, "\n") == card.Back {
			strategy.Correct(card)
			fmt.Print(Bold(Cyan("Correct!")))
			time.Sleep(correctSleep)
		} else {
			strategy.Incorrect(card)
			fmt.Print(Bold(Red(card.Back)))
			time.Sleep(incorrectSleep)
		}

		clearLine()
		moveCursorUp()
		clearLine()
		moveCursorUp()
		clearLine()
	}

	saveResults(*savePath, cardsFileName, &cards)
}

func loadCards(cardMap *map[string]string, sessionCards *map[string]*Card) Cards {
	var cards Cards
	for front, back := range *cardMap {
		session := (*sessionCards)[front]
		if session != nil {
			session.Front = front
			session.Back = back
			cards = append(cards, session)
		} else {
			cards = append(cards, &Card{
				Front:    front,
				Back:     back,
				Box:      0,
				LastSeen: 0,
			})
		}
	}
	return cards
}

func saveResults(savePath string, fileName string, cards *Cards) {
	d, err := yaml.Marshal(*cards)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("Saving results to: %s\n", path.Join(savePath, fileName))
	_ = os.Mkdir(savePath, os.ModePerm)
	ioutil.WriteFile(path.Join(savePath, fileName), d, 0644)
}

func moveCursorUp() {
	fmt.Printf("\033[1A")
}
func clearLine() {
	fmt.Printf("\r\033[K")
}
