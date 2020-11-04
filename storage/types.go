package storage

import "github.com/sheitm/ofever/scrape"

// Season provides an API for all events and athletes that have been active during a season. A season is defined by
// the year.
type Season interface {

	// Store stores the entire season in internal storage.
	Store(fetch *scrape.SeasonFetch) error
	Year() int
}


