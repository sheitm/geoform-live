package storage

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_splitPath(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"justYear", args{p:"/season/2020"}, []string{"season", "2020"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitPath(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_computedSeasonHandler_ServeHTTP(t *testing.T) {
	// No year
	req, err := http.NewRequest("GET", "/season", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := &computedSeasonHandler{}

	handler.ServeHTTP(rr, req)

	if rr.Code != 500 {
		t.Errorf("expected http code 500, got %d", rr.Code)
	}

	// Invalid year
	req, err = http.NewRequest("GET", "/season/202o", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != 500 {
		t.Errorf("expected http code 500, got %d", rr.Code)
	}

	// Computed season error
	fetch := func(y int)(*computedSeason, error) {
		return nil, fmt.Errorf("could not process year %d", y)
	}
	req, err = http.NewRequest("GET", "/season/2020", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()

	handler = &computedSeasonHandler{fetch: fetch}

	handler.ServeHTTP(rr, req)

	if rr.Code != 500 {
		t.Errorf("expected http code 500, got %d", rr.Code)
	}

	// Happy days
	fetch = func(y int)(*computedSeason, error) {
		return &computedSeason{Year: y}, nil
	}

	rr = httptest.NewRecorder()

	handler = &computedSeasonHandler{fetch: fetch}

	handler.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("expected http code 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	if body != `{"year":2020,"athletes":null,"competitions":null,"statistics":null}` {
		t.Errorf("unexpected body, got %s", body)
	}
}