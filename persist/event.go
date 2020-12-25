package persist

import "github.com/sheitm/ofever/scrape"

// Event is used to signal that a season fetch should be persisted.
type Event struct {
	// DoneChan used to signal back that the event is persisted. If error is not nil, it failed.
	DoneChan chan<- error

	// Fetch is the fetch to persist.
	Fetch *scrape.SeasonFetch
}
