package scrape

import (
	"fmt"
	"github.com/sheitm/ofever/types"
	"net/http"
	"testing"
	"time"
)

func Test_handler_Handle_noYear(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("GET", "/scrape", nil)
	if err != nil {
		t.Fatal(err)
	}

	h := &handler{}

	// Act
	roundTrip := h.Handle(req)

	// Assert
	if roundTrip.HandlerName != handlerName {
		t.Errorf("unexpected handler name, got %s", roundTrip.HandlerName)
	}
	if roundTrip.HTTPResponseCode != 500 {
		t.Errorf("unexpected http code 500, got %d", roundTrip.HTTPResponseCode)
	}
}

func Test_handler_Handle_future(t *testing.T) {
	// Arrange
	path := fmt.Sprintf("/scrape/%d", time.Now().Year() + 1)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	h := &handler{}

	// Act
	roundTrip := h.Handle(req)

	// Assert
	if roundTrip.HandlerName != handlerName {
		t.Errorf("unexpected handler name, got %s", roundTrip.HandlerName)
	}
	if roundTrip.HTTPResponseCode != 500 {
		t.Errorf("unexpected http code 500, got %d", roundTrip.HTTPResponseCode)
	}
}

func Test_handler_Handle_beforeBeginningOfTime(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("GET", "/scrape/2008", nil)
	if err != nil {
		t.Fatal(err)
	}

	h := &handler{}

	// Act
	roundTrip := h.Handle(req)

	// Assert
	if roundTrip.HandlerName != handlerName {
		t.Errorf("unexpected handler name, got %s", roundTrip.HandlerName)
	}
	if roundTrip.HTTPResponseCode != 500 {
		t.Errorf("unexpected http code 500, got %d", roundTrip.HTTPResponseCode)
	}
}

func Test_handler_Handle(t *testing.T) {
	// Arrange
	path := "/scrape/2020"
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	receivedURL := ""
	receivedYear := 0
	sentFetch := &types.SeasonFetch{
		Year:    2020,
		URL:     "",
	}
	starter := func(url string, year int, ch chan<- *types.SeasonFetch) {
		receivedURL = url
		receivedYear = year
		go func() {
			ch <- sentFetch
		}()
	}
	eventChan := make(chan *types.ScrapeEvent)
	sequenceTrigger := make(chan interface{})
	sequenceDone := make(chan struct{})
	h := &handler{
		eventChan:       eventChan,
		starter:         starter,
		sequenceTrigger: sequenceTrigger,
		finalDone:       sequenceDone,
	}

	// Act
	go func(ech <-chan *types.ScrapeEvent) {
		re := <- ech
		re.DoneChan <- nil
	}(eventChan)

	triggered := false
	go func(trigger <-chan interface{}, dc chan<- struct{}) {
		<-trigger
		triggered = true
		dc <- struct{}{}

	}(sequenceTrigger, sequenceDone)

	roundTrip := h.Handle(req)

	// Assert
	if roundTrip.HandlerName != handlerName {
		t.Errorf("unexpected handler name, got %s", roundTrip.HandlerName)
	}
	if roundTrip.HTTPResponseCode != 200 {
		t.Errorf("unexpected http code 200, got %d", roundTrip.HTTPResponseCode)
	}
	if receivedURL != "https://ilgeoform.no/rankinglop/" {
		t.Errorf("unexpected URL received, got %s", receivedURL)
	}
	if receivedYear != 2020 {
		t.Errorf("unexpected year received, got %d", receivedYear)
	}
	if !triggered {
		t.Error("sequence wasn't triggered")
	}
}