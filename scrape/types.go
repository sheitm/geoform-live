package scrape

import "github.com/sheitm/ofever/contracts"

// SeasonFetch the results after having attempted to get all events for a season.
type SeasonFetch struct {
	URL string

	Results []*Result

	Error error
}

// Result is the result of a single event scraping session.
type Result struct {
	// URL is the provided address to the page containing the results for the event.
	URL   string

	// Event is the details of the event.
	Event *contracts.Event

	// Error details if any error occurred during execution.
	Error error
}