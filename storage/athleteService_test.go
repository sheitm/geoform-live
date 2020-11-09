package storage

import (
	"encoding/json"
	"github.com/sheitm/ofever/scrape"
	"testing"
)

func Test_athleteServiceImpl_Start(t *testing.T) {
	// Arrange
	var f scrape.SeasonFetch
	err := json.Unmarshal([]byte(jsonSeason2019), &f)
	if err != nil {
		t.Errorf("unexpected error when unmarshaling json, %v", err)
	}

	var receivedAthletes []*athlete
	persist := func(a []*athlete) {
		receivedAthletes = a
	}
	fetch := func() ([]*athlete, error) {
		return []*athlete{}, nil
	}

	service := newAthleteService(persist, fetch)

	element := seasonSyncElement{
		seasonChan: make(chan *scrape.SeasonFetch),
		doneChan:   make(chan struct{}),
	}

	// Act
	service.Start(element)
	element.seasonChan <- &f

	<- element.doneChan

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
