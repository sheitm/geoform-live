package storage

import (
	"github.com/sheitm/ofever/scrape"
	"log"
)

var storageDirectory string

func Start(storageFolder string, seasonChan <-chan *scrape.SeasonFetch) {
	storageDirectory = storageFolder

	currentCache = &cache{getter: getJSONsFromDirectory}
	currentCache.init()

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
		folder: storageDirectory,
		year:   year,
	}
}
