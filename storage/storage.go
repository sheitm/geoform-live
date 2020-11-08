package storage

import (
	"github.com/sheitm/ofever/scrape"
)

//var currentStorageService storageService
//var currentAthleteService athleteService
//var currentComputeService computeService

// Start starts the internal functionality that reacts on incoming scraped seasons via the channel seasonChan. Also
// necessary internal wire up, sp that this method should be invoked before using this package.
func Start(storageFolder string, seasonChan <-chan *scrape.SeasonFetch) {
	sequenceHandler := &sequentialSeasonHandlers{}

	currentStorageService := newStorageService(storageFolder)
	sequenceHandler.Add(currentStorageService)

	athletePersist := func(athletes []*Athlete) {
		fn := func(interface{}) string { return "athletes.json" }
		currentStorageService.Store(athletes, fn)
	}
	athleteFetch := func()([]*Athlete, error) {
		var l []*Athlete
		fn := func(interface{}) string { return "athletes.json" }
		err := currentStorageService.Fetch(&l, fn)
		if err != nil {
			return l, err
		}
		return l, nil
	}
	currentAthleteService := newAthleteService(athletePersist, athleteFetch)
	sequenceHandler.Add(currentAthleteService)

	currentComputeService := newComputeService(currentAthleteService.ID)
	sequenceHandler.Add(currentComputeService)

	sequenceHandler.Start(seasonChan)
}

//func Handlers() map[string]http.Handler {
//
//}

// season/2020/athletes