package scrape

import (
	"github.com/sheitm/ofever/contracts"
	"net/http"
)

type scraper interface {
	Scrape(url string) (*contracts.Event, error)
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

func StartScrape(url string, resultChan chan<- *Result) {
	go func(url string, resultChan chan<- *Result) {
		res := &Result{URL: url}
		scraper := &eventScraper{client: &http.Client{}}
		event, err := scraper.Scrape(url)
		res.Event = event
		res.Error = err
		resultChan <- res
	}(url, resultChan)
}