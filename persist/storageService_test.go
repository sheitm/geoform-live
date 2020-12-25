package persist

import (
	"errors"
	"github.com/3lvia/telemetry-go"
	"github.com/sheitm/ofever/scrape"
	"reflect"
	"testing"
)

func Test_parseBlobConnectionString(t *testing.T) {
	cs := `DefaultEndpointsProtocol=https;AccountName=ofeverdevelopment;AccountKey=SECRET_KEY==;EndpointSuffix=core.windows.net`
	expected := map[string]string {
		"DefaultEndpointsProtocol": "https",
		"AccountName":              "ofeverdevelopment",
		"AccountKey":               "SECRET_KEY==",
		"EndpointSuffix":           "core.windows.net",
	}
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{"1", args{s: cs}, expected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseBlobConnectionString(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseBlobConnectionString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_storageService_start(t *testing.T) {
	// Arrange
	logChannels := telemetry.StartEmpty()
	sentFetch := &scrape.SeasonFetch{
		Series:  "geoform",
		Year:    2020,
	}
	var receivedFetch *scrape.SeasonFetch
	service := &storageService{
		logChannels: logChannels,
		save: func(fetch *scrape.SeasonFetch) error {
			receivedFetch = fetch
			return nil
		},
	}

	eChan := make(chan *Event)
	doneChan := make(chan error)

	// Act
	go service.start(eChan)
	e := &Event{
		DoneChan: doneChan,
		Fetch:    sentFetch,
	}
	eChan <- e
	err := <- doneChan

	// Assert
	if err != nil {
		t.Errorf("unexpected error, %v", err)
	}
	if receivedFetch.Year != 2020 {
		t.Errorf("unexpected received fetch, got %v", receivedFetch)
	}
}

func Test_storageService_start_error(t *testing.T) {
	// Arrange
	logChannels := telemetry.StartEmpty()
	sentFetch := &scrape.SeasonFetch{
		Series:  "geoform",
		Year:    2020,
	}
	service := &storageService{
		logChannels: logChannels,
		save: func(fetch *scrape.SeasonFetch) error {
			return errors.New("some failure")
		},
	}

	eChan := make(chan *Event)
	doneChan := make(chan error)

	// Act
	go service.start(eChan)
	e := &Event{
		DoneChan: doneChan,
		Fetch:    sentFetch,
	}
	eChan <- e
	err := <- doneChan

	// Assert
	if err == nil {
		t.Error("expected error, got none")
	}
}