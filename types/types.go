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
	Front    front     `yaml:"front"`
	Back     back      `yaml:"back"`
	Box      int       `yaml:"box"`
	LastSeen time.Time `yaml:"last_seen"`
	// Direction/magnitude of the last answer (e.g. box 5 -> 1 = -4)
	Velocity int `yaml:"velocity"`
}

func (c Card) String() string {
	return fmt.Sprintf("Front: %s\nBack: %s\nBox: %d\nLast Seen: %s\nVelocity: %d", c.Front, c.Back, c.Box, c.LastSeen.Format(time.UnixDate), c.Velocity)
}

type minutes = int
type Config struct {
	FilePath               string
	SavePath               string
	CorrectSleepMs         time.Duration
	IncorrectSleepMs       time.Duration
	StrategyName           string
	StorageName            string
	Flip                   string
	IgnoreParenthesis      bool
	CardsToIntroducePerDay int
	CardsToReviewPerDay    int
	Steps                  []minutes
	ErrorThreshold         int
}
