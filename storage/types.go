package storage

import (
	"github.com/sheitm/ofever/scrape"
	"time"
)

// Service provides an API for all events and athletes that have been active during a season. A season is defined by
// the year.
type Service interface {

	// Store stores the entire season in internal storage.
	Store(fetch *scrape.SeasonFetch) error
	Year() int
}

type Athlete struct {
	Name    string          `json:"name"`
	Results []AthleteResult `json:"results"`
}

type AthleteResult struct {
	Event        string    `json:"event"`
	Course       string    `json:"course"`
	Disqualified bool      `json:"disqualified"`
	Placement    int       `json:"placement"`
	ElapsedTime  time.Time `json:"elapsed_time"`
	Points       float64   `json:"points"`
}
