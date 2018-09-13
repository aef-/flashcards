package types

import (
	"fmt"
	"time"
)

type Cards = []*Card
type front = string
type back = string

// Learning Session
type Session = map[front]*Card

type Definitions = map[front]back
type Card struct {
	Front    front     `yaml:"-"`
	Back     back      `yaml:"-"`
	Box      int       `yaml:"box"`
	LastSeen time.Time `yaml:"last_seen"`
	// Direction/magnitude of the last answer (e.g. box 5 -> 1 = -4)
	Velocity int `yaml:"velocity"`
}

func (c Card) String() string {
	return fmt.Sprintf("Front: %s\nBack: %s\nBox: %d\nLast Seen: %d", c.Front, c.Back, c.Box, c.LastSeen.Format(time.UnixDate))
}

type minutes = int
type Config struct {
	CorrectSleepMs         time.Duration
	IncorrectSleepMs       time.Duration
	StrategyName           string
	Flip                   string
	IgnoreParenthesis      bool
	CardsToIntroducePerDay int
	CardsToReviewPerDay    int
	Steps                  []minutes
	ErrorThreshold         int
}
