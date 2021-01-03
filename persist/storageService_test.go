package persist

import (
	"context"
	"errors"
	"github.com/3lvia/telemetry-go"
	"github.com/sheitm/ofever/types"
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
	sentFetch := &types.SeasonFetch{
		Series:  "geoform",
		Year:    2020,
	}
	var receivedFetch *types.SeasonFetch
	service := &storageService{
		logChannels: logChannels,
		save: func(ctx context.Context, config map[string]string, fetch *types.SeasonFetch) error {
			receivedFetch = fetch
			return nil
		},
	}

	eChan := make(chan *types.ScrapeEvent)
	doneChan := make(chan error)

	// Act
	go service.start(eChan)
	e := &types.ScrapeEvent{
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
	sentFetch := &types.SeasonFetch{
		Series:  "geoform",
		Year:    2020,
	}
	service := &storageService{
		logChannels: logChannels,
		save: func(ctx context.Context, config map[string]string, fetch *types.SeasonFetch) error {
			return errors.New("some failure")
		},
	}

	eChan := make(chan *types.ScrapeEvent)
	doneChan := make(chan error)

	// Act
	go service.start(eChan)
	e := &types.ScrapeEvent{
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