package types

import "github.com/sheitm/ofever/contracts"

// SeasonFetch the results after having attempted to get all events for a season.
type SeasonFetch struct {
	Series  string    `json:"series"`
	Year    int       `json:"year"`
	URL     string    `json:"url"`
	Results []*Result `json:"results"`
	Error   string    `json:"error"`
}

// Result is the result of a single event scraping session.
type Result struct {
	// URL is the provided address to the page containing the results for the event.
	URL string `json:"url"`

	// Event is the details of the event.
	Event *contracts.Event `json:"event"`

	// Error details if any error occurred during execution.
	Error string `json:"error"`
}

// Event is used to signal that a season fetch should be persisted.
type Event struct {
	// DoneChan used to signal back that the event is persisted. If error is not nil, it failed.
	DoneChan chan<- error

	// Fetch is the fetch to persist.
	Fetch *SeasonFetch
}