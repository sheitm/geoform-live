package competitions

import (
	"bytes"
	"github.com/3lvia/telemetry-go"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handler_Handle_single(t *testing.T) {
	// Arrange
	expectedBody := `{"series":"geoform","season":"2020","number":1,"name":"Name","url_live_lox":"","courses":null,"week_day":"","date":"0001-01-01T00:00:00Z","place":"","organizer":"","responsible":""}`
	buf := new(bytes.Buffer)
	logChannels := telemetry.StartEmpty(telemetry.WithWriter(buf))
	req, err := http.NewRequest("GET", "/competitions/geoform/2020/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	get := func(series, season string, number int) *comp {
		if series == "geoform" && season == "2020" && number == 1 {
			return &comp{
				Series:      series,
				Season:      season,
				Number:      1,
				Name:        "Name",
			}
		}
		return nil
	}

	rr := httptest.NewRecorder()
	handler := &handler{get: get}
	wrapper := telemetry.Wrap(handler, logChannels)

	// Act
	wrapper.ServeHTTP(rr, req)

	// Assert
	if rr.Code != 200 {
		t.Errorf("expected http code 200, got %d", rr.Code)
	}
	if rr.Body.String() != expectedBody {
		t.Errorf("unexpected body, got %s", rr.Body.String())
	}

	logs := buf.String()
	if logs == "" {
		t.Error("expected logs, got nothing")
	}
}

func Test_handler_Handle_single_malformedURL(t *testing.T) {
	// Arrange
	expectedBody := `strconv.Atoi: parsing "NaN": invalid syntax`
	buf := new(bytes.Buffer)
	logChannels := telemetry.StartEmpty(telemetry.WithWriter(buf))
	req, err := http.NewRequest("GET", "/competitions/geoform/2020/NaN", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := &handler{}
	wrapper := telemetry.Wrap(handler, logChannels)

	// Act
	wrapper.ServeHTTP(rr, req)

	// Assert
	if rr.Code != 500 {
		t.Errorf("expected http code 500, got %d", rr.Code)
	}
	if rr.Body.String() != expectedBody {
		t.Errorf("unexpected body, got %s", rr.Body.String())
	}

	logs := buf.String()
	if logs == "" {
		t.Error("expected logs, got nothing")
	}
}

func Test_handler_Handle_all(t *testing.T) {
	// Arrange
	expectedBody := `[{"series":"geoform","season":"2020","number":1,"name":"Name"}]`
	buf := new(bytes.Buffer)
	logChannels := telemetry.StartEmpty(telemetry.WithWriter(buf))
	req, err := http.NewRequest("GET", "/competitions/geoform/2020", nil)
	if err != nil {
		t.Fatal(err)
	}

	getAll := func(series, season string) []*compHeader {
		if series == "geoform" && season == "2020" {
			return []*compHeader{
				&compHeader{
					Series: series,
					Season: season,
					Number: 1,
					Name:   "Name",
				},
			}
		}
		return nil
	}

	rr := httptest.NewRecorder()
	handler := &handler{getAll: getAll}
	wrapper := telemetry.Wrap(handler, logChannels)

	// Act
	wrapper.ServeHTTP(rr, req)

	// Assert
	if rr.Code != 200 {
		t.Errorf("expected http code 200, got %d", rr.Code)
	}
	if rr.Body.String() != expectedBody {
		t.Errorf("unexpected body, got %s", rr.Body.String())
	}

	logs := buf.String()
	if logs == "" {
		t.Error("expected logs, got nothing")
	}
}