package storage

import (
	"encoding/json"
	"github.com/sheitm/ofever/types"
	"testing"
	"time"
)

func Test_makeCompetitionID(t *testing.T) {
	// Arrange
	dt := time.Date(2020, 11, 28, 9, 52, 1, 0, time.UTC)

	// Act
	id := makeCompetitionID(dt)

	// Assert
	if len(id) != 13 {
		t.Errorf("unexpected length for id %s", id)
	}
	if id[0:9] != "20201128-" {
		t.Errorf("unexpected id, got %s", id)
	}
}

func Test_competitionServiceImpl_Start(t *testing.T) {
	// Arrange
	var f types.SeasonFetch
	err := json.Unmarshal([]byte(jsonSeason2019), &f)
	if err != nil {
		t.Errorf("unexpected error when unmarshaling json, %v", err)
	}

	var persistedCompetitions []*competition
	persistFunc := func(c []*competition) {
		persistedCompetitions = c
	}

	fetchFunc := func()([]*competition, error) {
		return nil, nil
	}

	service := newCompetitionService(persistFunc, fetchFunc)

	seasonChan := make(chan *types.SeasonFetch)
	doneChan := make(chan struct{})
	syncElement := seasonSyncElement{
		seasonChan: seasonChan,
		doneChan:   doneChan,
	}

	// Act
	service.Start(syncElement)
	seasonChan <- &f
	<- doneChan

	// Assert
	if len(persistedCompetitions) != 27 {
		t.Errorf("expected 27 persisted competitions, got %d", len(persistedCompetitions))
	}

	//for _, result := range f.Results {
	//	e := result.Event
	//	fmt.Println(e.Name)
	//	for _, c := range e.Courses {
	//		fmt.Printf("   %s\n", c.Name)
	//	}
	//}
}

func Test_getCourseLength(t *testing.T) {
	type args struct {
		n string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"comma", args{"Resultater Lang (5,4 km)"}, 5.4},
		{"dot", args{"Resultater Kort (2.4 km)"}, 2.4},
		{"invalid", args{"Resultater Kort (2.4 kkm)"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCourseLength(tt.args.n); got != tt.want {
				t.Errorf("getCourseLength() = %v, want %v", got, tt.want)
			}
		})
	}
}