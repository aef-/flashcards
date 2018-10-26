package ui

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	t "github.com/aef-/flashcards/types"
	. "github.com/logrusorgru/aurora"
)

var correctSleepMs = flag.Int("correct-time", 500, "time (ms) until next question is shown after a correct question")
var savePath = flag.String("save-path", "/Users/207751/.flashcards", "directory where to save sessions")
var incorrectSleepMs = flag.Int("incorrect-time", 1500, "time (ms) until next question is shown after an incorrect question")
var strategy = flag.String("strategy", "leitner", "One of all|all-random|leitner")
var storage = flag.String("storage", "local", "One of local")
var flip = flag.String("flip", "front", "Prompt with either front, back or random")
var ignoreParenthesis = flag.Bool("ignore-paren", true, "Does not count text within parenthesis when comparing results")
var cardsToReviewPerDay = flag.Int("review-cards-per-day", 200, "Maximum number of cards to review each day, this does not include new cards")
var cardsToIntroducePerDay = flag.Int("new-cards-per-day", 2, "Maximum number of new cards to introduce each day, missed days do not accumulate")
var steps = flag.String("steps", "0,10,1440,4320,7200,43200", "Steps in minutes for different levels")
var modifier = flag.String("modifier", "1,10,1440,4320", "Steps for new cards")
var errorThreshold = flag.Int("error-threshold", 10, "Must be this percent wrong to be counted incorrect")

type Cli struct{}

func (c Cli) Setup() {}
func (c Cli) LoadConfig() t.Config {
	flag.Parse()

	args := flag.Args()

	correctSleep := time.Duration(time.Millisecond * time.Duration(*correctSleepMs))
	incorrectSleep := time.Duration(time.Millisecond * time.Duration(*incorrectSleepMs))
	stepsArr, err := sliceAtoi(strings.Split(*steps, ","))

	if err != nil {
		log.Fatal("Could not convert steps to ints")
	}

	return t.Config{
		FilePath:               args[0],
		SavePath:               *savePath,
		CorrectSleepMs:         correctSleep,
		IncorrectSleepMs:       incorrectSleep,
		StorageName:            *storage,
		StrategyName:           *strategy,
		CardsToReviewPerDay:    *cardsToReviewPerDay,
		CardsToIntroducePerDay: *cardsToIntroducePerDay,
		Steps:                  stepsArr,
		ErrorThreshold:         *errorThreshold,
	}
}

func (c Cli) Prompt(card *t.Card, config t.Config) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("> %s\n", card.Front)
	answer, _ := reader.ReadString('\n')
	if strings.TrimRight(answer, "\n") == card.Back {
		fmt.Print(Bold(Cyan("Correct!")))
		time.Sleep(config.CorrectSleepMs)
		clearPrompt()
		return true
	} else {
		fmt.Print(Bold(Red(card.Back)))
		time.Sleep(config.IncorrectSleepMs)
		clearPrompt()
		return false
	}
}

func clearPrompt() {
	clearLine()
	moveCursorUp()
	clearLine()
	moveCursorUp()
	clearLine()
}

func moveCursorUp() {
	fmt.Printf("\033[1A")
}
func clearLine() {
	fmt.Printf("\r\033[K")
}

func sliceAtoi(sa []string) ([]int, error) {
	si := make([]int, 0, len(sa))
	for _, a := range sa {
		i, err := strconv.Atoi(a)
		if err != nil {
			return si, err
		}
		si = append(si, i)
	}
	return si, nil
}
