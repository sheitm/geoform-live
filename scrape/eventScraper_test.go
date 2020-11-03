package scrape

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_eventScraper_Scrape(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(htmlEventResponse))
	}))
	defer server.Close()


	url := server.URL + "/res2020-10-24.html"
	scraper := eventScraper{client: server.Client()}

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
}
