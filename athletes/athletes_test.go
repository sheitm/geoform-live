package athletes

import (
	"github.com/3lvia/telemetry-go"
	"github.com/sheitm/ofever/persist"
	"github.com/sheitm/ofever/sequence"
	"github.com/sheitm/ofever/types"
	"testing"
)

func TestStart(t *testing.T) {
	// Arrange
	doneChan := make(chan struct{})
	logChannels := telemetry.StartEmpty()
	fetch := seasonFetch()
	sequenceAdder := func(ch chan<- *sequence.Event) {
		go func() {
			ch <- &sequence.Event{
				Payload:  fetch,
				DoneChan: doneChan,
			}
		}()
	}
	var persistedElements []*persist.Element
	persister := func(elements []*persist.Element, dc chan<- struct{}) {
		persistedElements = elements
		go func() {dc <- struct{}{}}()
	}

	// Act
	Start(sequenceAdder, persister, logChannels)
	<- doneChan

	// Assert
	if len(persistedElements) != 3 {
		t.Errorf("expected 3 persisted elements, got %d", len(persistedElements))
	}

	var ee *persist.Element
	for _, element := range persistedElements {
		a := element.Data.(*athleteWithID)
		if a.Name == "Johan Mygland" {
			ee = element
		}
	}
	p := ee.PathGetter(ee.Data)
	// athletes/4d0523b4-9ad3-72d4-74f4-cd3a627e909c.json
	if len(p) != 50 {
		t.Errorf("unexpected path, got %s len(%d)", p, len(p))
	}
}

func seasonFetch() *types.SeasonFetch {
	return &types.SeasonFetch{
		Series: "geoform",
		Year:   2020,
		URL:    "",
		Results: []*types.ScrapeResult{
			{
				Event: &types.Event{
					Courses: []*types.Course{
						{
							Results: []*types.Result{
								{Athlete: "Kåre Fyn"},
								{Athlete: "Benny Carlsson", Club: "Sturla"},
								{Athlete: "Johan Mygland", Club: "Geoform"},
							},
						},
					},
				},
			},
		},
		Error: "",
	}
}