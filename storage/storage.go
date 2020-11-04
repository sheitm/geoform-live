package storage

import (
	"github.com/sheitm/ofever/scrape"
	"log"
)

const folder = `C:\temp\ofever`

func Start(seasonChan <-chan *scrape.SeasonFetch) {
	go func(sc <-chan *scrape.SeasonFetch) {
		for {
			fetch := <- sc
			s := NewService(fetch.Year)
			err := s.Store(fetch)
			if err != nil {
				log.Printf("%v", err)
			}
		}
	}(seasonChan)
}

func NewService(year int) Service {
	return &service{
		folder: folder,
		year:   year,
	}
}
