package scrape

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_eventScraper_Scrape_404(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
	}))
	defer server.Close()

	url := server.URL + "/res2020-10-24.html"
	scraper := eventScraper{client: server.Client()}

	// Act
	event, err := scraper.Scrape(url)

	// Assert
	if event != nil {
		t.Error("expected nil event")
	}
	if err == nil {
		t.Error("expected error, got none")
	}
}

func Test_eventScraper_Scrape(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/res2020-10-24.html" {
			t.Errorf("unexpected path, got %s", req.URL.Path)
		}
		rw.Write([]byte(htmlEventResponse))
	}))
	defer server.Close()


	url := server.URL + "/res2020-10-24.html"
	scraper := eventScraper{
		client: server.Client(),
		row: &tableRow{
			baseURL: "http://base",
			values: map[int]cellValue{
				0: { text:  "1", value: "1" },
				1: { text:  "03.06", value: "03.06" },
				2: { text:  "ons", value: "ons" },
				3: { text:  "Skansebakken", value: "https://eventor.no/innbydelse.html" },
				4: { text:  "Bergendal-Lyse", value: "Bergendal-Lyse" },
				5: { text:  "Fossum IF", value: "Fossum IF" },
				6: { text:  "Resultater", value: "res2020-10-24.html" },
				7: { text:  "LL", value: "https://LiveLox/res2020-10-24" },
				8: { text:  "res.zip", value: "results.zip" },
			},
			year:    2020,
		},
	}

	// Act
	event, err := scraper.Scrape(url)

	// Assert
	if err != nil {
		t.Errorf("unexpected error, %v", err)
	}

	if event.Name != "Rankingløp nr 16" {
		t.Errorf("unexpected event name, got %s", event.Name)
	}
	if len(event.Courses) != 3 {
		t.Errorf("expected 3 courses, got %d", len(event.Courses))
	}
	expectedCount := 0
	for _, course := range event.Courses {
		if course.Name == "Resultater Lang (5.7 km)" || course.Name == "Resultater Mellom (4.4 km)" || course.Name == "Resultater Kort (3.1 km)" {
			expectedCount++
		}
	}
	if expectedCount != 3 {
		t.Error("at least one course has unexpected name")
	}
	if fmt.Sprintf("%v", event.Date) != "2020-06-03 00:00:00 +0000 UTC" {
		t.Errorf("unexpected event date, got %v", event.Date)
	}
	if event.URLInvite != "https://eventor.no/innbydelse.html" {
		t.Errorf("unexpected URLInvite, got %s", event.URLInvite)
	}
	if event.URLLiveLox != "https://LiveLox/res2020-10-24" {
		t.Errorf("unexpected URLLiveLox, got %s", event.URLLiveLox)
	}
	if event.WeekDay != "ons" {
		t.Errorf("unexpected WeekDay, got %s", event.WeekDay)
	}
	if event.Place != "Skansebakken" {
		t.Errorf("unexpected Place, got %s", event.Place)
	}
	if event.Organizer != "Fossum IF" {
		t.Errorf("unexpected Organizer, got %s", event.Organizer)
	}
}
