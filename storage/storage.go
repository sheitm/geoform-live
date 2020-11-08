package storage

import (
	"github.com/sheitm/ofever/scrape"
)

var currentStorageService storageService
var currentAthleteService athleteService

func Start(storageFolder string, seasonChan <-chan *scrape.SeasonFetch) {
	storageSeasonChan := make(chan *scrape.SeasonFetch)
	athleteSeasonChan := make(chan *scrape.SeasonFetch)
	dispatches := []chan<- *scrape.SeasonFetch{
		storageSeasonChan,
		athleteSeasonChan,
	}

	go func(sc <-chan *scrape.SeasonFetch, dispatches []chan<- *scrape.SeasonFetch) {
		for {
			fetch := <- sc
			for _, dispatch := range dispatches {
				dispatch <- fetch
			}
		}
	}(seasonChan, dispatches)

	currentStorageService = newStorageService(storageFolder)
	currentStorageService.Start(storageSeasonChan)

	currentAthleteService = newAthleteService()
	currentAthleteService.Start(athleteSeasonChan)

	currentCache = &cache{getter: getJSONsFromDirectory}
	currentCache.init()


}
