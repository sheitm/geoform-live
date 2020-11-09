package storage

import (
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
	currentAthleteService := newAthleteService(athletePersist, athleteFetch)
	sequenceHandler.Add(currentAthleteService)

	currentComputeService := newComputeService(currentAthleteService.ID)
	sequenceHandler.Add(currentComputeService)

	sequenceHandler.Start(seasonChan)

	// athlete
	handler := newAthleteHandler(currentAthleteService.List)
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

