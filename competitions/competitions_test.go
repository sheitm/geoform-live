package competitions

import (
	"fmt"
	"github.com/3lvia/telemetry-go"
	"github.com/sheitm/ofever/persist"
	"github.com/sheitm/ofever/sequence"
	"github.com/sheitm/ofever/types"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	// Arrange
	series := "geoform"
	year := 2020
	var trigger chan<- *sequence.Event
	sequenceAdder := func(ch chan<- *sequence.Event) {
		trigger = ch
	}
	athleteID := func(name, club string) string {
		return fmt.Sprintf("id-%s-%s", name, club)
	}
	var receivedElements []*persist.Element
	persistFunc := func(elements []*persist.Element, doneChan chan<- struct{}){
		receivedElements = elements
		go func(){doneChan <- struct{}{}}()
	}
	logChannels := telemetry.StartEmpty()

	// Act
	Start(sequenceAdder, athleteID, persistFunc, logChannels)
	doneChan := make(chan struct{})
	e := &sequence.Event{
		Payload:  fetch(series, year),
		DoneChan: doneChan,
	}
	trigger <- e
	<- doneChan

	// Assert
	if len(receivedElements) != 1 {
		t.Errorf("expected 1 received element, got %d", len(receivedElements))
	}
	re := receivedElements[0]
	if re.Container != series {
		t.Errorf("unexpected container, got %s", re.Container)
	}
	p := re.PathGetter(re.Data)
	if p != "2020/competitions/1.json" {
		t.Errorf("unexpected path, got %s", p)
	}
}

func fetch(s string, y int) *types.SeasonFetch {
	return &types.SeasonFetch{
		Series: s,
		Year:   y,
		Results: []*types.ScrapeResult{
			{
				Event: &types.Event{
					Number: 1,
					Name:   "O-løp!",
					Courses: []*types.Course{
						&types.Course{
							Name:       "Lang",
							Results:    []*types.Result {
								&types.Result{
									Placement:       1,
									Disqualified:    false,
									Athlete:         "Kåre Bentsen",
									Club:            "Haslum IL",
									ElapsedTime:     1 * time.Hour,
									MissingControls: 0,
									Points:          150.12,
								},
							},
						},
					},
				},
			},
		},
	}
}