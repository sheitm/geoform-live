package scrape

import "github.com/sheitm/ofever/contracts"

type Scraper interface {
	Scrape(url string) (*contracts.Event, error)
}

func NewScraper() Scraper {
	return &eventScraper{}
}