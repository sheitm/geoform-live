package main

import (
	"fmt"
	"github.com/sheitm/ofever/scrape"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func Test_scrapeHandler_ServeHTTP_noYear(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("GET", "/scrape", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	//seaconChan := make(chan *scrape.SeasonFetch)

	handler := &scrapeHandler{
	}

	// Act
	handler.ServeHTTP(rr, req)

	// Assert
	if rr.Code != 500 {
		t.Errorf("expected http status code 500, got %d", rr.Code)
	}
}

func Test_scrapeHandler_ServeHTTP_future(t *testing.T) {
	// Arrange
	path := fmt.Sprintf("/scrape/%d", time.Now().Year() + 1)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	//seaconChan := make(chan *scrape.SeasonFetch)

	handler := &scrapeHandler{
	}

	// Act
	handler.ServeHTTP(rr, req)

	// Assert
	if rr.Code != 500 {
		t.Errorf("expected http status code 500, got %d", rr.Code)
	}
}

func Test_scrapeHandler_ServeHTTP_beforeBeginningOfTime(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("GET", "/scrape/2008", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	//seaconChan := make(chan *scrape.SeasonFetch)

	handler := &scrapeHandler{
	}

	// Act
	handler.ServeHTTP(rr, req)

	// Assert
	if rr.Code != 500 {
		t.Errorf("expected http status code 500, got %d", rr.Code)
	}
}

func Test_scrapeHandler_ServeHTTP_thisYear(t *testing.T) {
	// Arrange
	thisYear := time.Now().Year()
	path := fmt.Sprintf("/scrape/%d", thisYear)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	seasonChan := make(chan *scrape.SeasonFetch)

	var receivedURL string
	var receivedYear int
	var fetch *scrape.SeasonFetch

	handler := &scrapeHandler{
		seasonChan: seasonChan,
		starter:    func(url string, year int, sc chan<- *scrape.SeasonFetch){
			receivedURL = url
			receivedYear = year
			sc <- &scrape.SeasonFetch{Year: year}
		},
	}

	// Act
	go handler.ServeHTTP(rr, req)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(wg *sync.WaitGroup){
		fetch = <- seasonChan
		wg.Done()
	}(wg)

	wg.Wait()

	// Assert
	if rr.Code != 200 {
		t.Errorf("expected http status code 200, got %d", rr.Code)
	}
	if receivedURL != `https://ilgeoform.no/rankinglop/` {
		t.Errorf("unexpected URL received, got %s", receivedURL)
	}
	if thisYear != thisYear {
		t.Errorf("unexpected year received, got %d", receivedYear)
	}
	if fetch.Year != receivedYear {
		t.Errorf("unexpected year in fetch, got %d", fetch.Year)
	}
}
