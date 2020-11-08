package storage

import (
	"encoding/json"
	"github.com/sheitm/ofever/scrape"
	"sync"
	"testing"
)

func Test_athleteServiceImpl_Start(t *testing.T) {
	// Arrange
	var f scrape.SeasonFetch
	err := json.Unmarshal([]byte(jsonSeason2019), &f)
	if err != nil {
		t.Errorf("unexpected error when unmarshaling json, %v", err)
	}

	var receivedAthletes []*Athlete
	wg := &sync.WaitGroup{}
	wg.Add(1)
	persist := func(a []*Athlete) {
		receivedAthletes = a
		wg.Done()
	}

	seasonChan := make(chan *scrape.SeasonFetch)
	service := newAthleteService(persist)

	// Act
	service.Start(seasonChan)
	seasonChan <- &f

	wg.Wait()

	// Assert
	if len(receivedAthletes) != 948 {
		t.Errorf("unexpetced number of athletes, got %d", len(receivedAthletes))
	}

	athletes, err := service.List()
	if err != nil {
		t.Errorf("unexpected error when listing athletes, %v", err)
	}
	if len(receivedAthletes) != len(athletes) {
		t.Errorf("count mismatch between received athletes (%d) and listed (%d)", len(receivedAthletes), len(athletes))
	}
}
