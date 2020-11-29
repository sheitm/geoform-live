package storage

import (
	"encoding/json"
	"github.com/sheitm/ofever/scrape"
	"net/http"
)

var handlers []httpHandler

// Start starts the internal functionality that reacts on incoming scraped seasons via the channel seasonChan. Also
// necessary internal wire up, sp that this method should be invoked before using this package.
func Start(storageFolder string, seasonChan <-chan *scrape.SeasonFetch) {
	sequenceHandler := &sequentialSeasonHandlers{}

	currentStorageService := newStorageService(storageFolder)
	sequenceHandler.Add(currentStorageService)

	competitionFetch, competitionPersist := getCompetitionFunctions(currentStorageService)
	currentCompetitionService := newCompetitionService(competitionPersist, competitionFetch)
	sequenceHandler.Add(currentCompetitionService)

	athleteFetch, athletePersist := getAthleteFunctions(currentStorageService)
	currentAthleteService := newAthleteService(athletePersist, athleteFetch)
	sequenceHandler.Add(currentAthleteService)

	computedSeasonFetch := getComputedSeasonsFunctions(currentStorageService, storageFolder, currentAthleteService.ID)
	currentComputeService := newComputeService(computedSeasonFetch, currentAthleteService.ID)
	sequenceHandler.Add(currentComputeService)

	sequenceHandler.Start(seasonChan)

	// athlete
	handler := newAthleteHandler(currentAthleteService.List, currentComputeService.ComputedSeason)
	handlers = append(handlers, handler)
	// season
	handler = newComputedSeasonHandler(currentComputeService.ComputedSeason)
	handlers = append(handlers, handler)
}

// Handlers returns a map of http.Handlers where the keys of the map are the paths to map the handlers to.
func Handlers() map[string]http.Handler {
	h := map[string]http.Handler{}
	for _, handler := range handlers {
		h[handler.Path()] = handler
	}
	return h
}

func getComputedSeasonsFunctions(currentStorageService storageService, folder string, getID athleteIDFunc) computedSeasonsFetchFunc {
	fetch := func()([]*computedSeason, error){
		contents, err := currentStorageService.ReadFolder(folder, `season_\d{4}.json`)
		if err != nil {
			return nil, err
		}
		var res []*computedSeason
		for _, c := range contents {
			var scrapedFetch scrape.SeasonFetch
			err := json.Unmarshal([]byte(c), &scrapedFetch)
			if err != nil {
				return nil, err
			}
			computed, err := computeSeasonForFetch(&scrapedFetch, getID)
			if err != nil {
				return nil, err
			}
			res = append(res, computed)
		}
		return res, nil
	}
	return fetch
}

func getAthleteFunctions(currentStorageService storageService) (athleteFetchFunc, athletePersistFunc)  {
	athletePersist := func(athletes []*athlete) {
		fn := func(interface{}) string { return "athletes.json" }
		currentStorageService.Store(athletes, fn)
	}
	athleteFetch := func()([]*athlete, error) {
		var l []*athlete
		fn := func(interface{}) string { return "athletes.json" }
		err := currentStorageService.Fetch(&l, fn)
		if err != nil {
			return l, err
		}
		return l, nil
	}

	return athleteFetch, athletePersist
}

func getCompetitionFunctions(currentStorageService storageService) (competitionFetchFunc, competitionPersistFunc) {
	competitionPersist := func(competitions []*competition) {
		fn := func(interface{}) string { return "competitions.json"}
		currentStorageService.Store(competitions, fn)
	}
	competitionFetch := func() ([]*competition, error) {
		var c []*competition
		fn := func(interface{}) string { return "competitions.json"}
		err := currentStorageService.Fetch(&c, fn)
		if err != nil {
			return c, err
		}
		return c, nil
	}
	return competitionFetch, competitionPersist
}