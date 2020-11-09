package storage

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_athleteHandler_Path(t *testing.T) {
	type fields struct {
		fetch athleteFetchFunc
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"1", fields{}, "/athlete/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &athleteHandler{
				fetch: tt.fields.fetch,
			}
			if got := h.Path(); got != tt.want {
				t.Errorf("Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_athleteHandler_ServeHTTP(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("GET", "/athlete", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := &athleteHandler{
		fetch: func() ([]*athlete, error) {
			athletes := []*athlete{
				&athlete{
					ID:   "1",
					Name: "a",
				},
			}
			return athletes, nil
		},
	}

	// Act
	handler.ServeHTTP(rr, req)

	// Assert
	if rr.Code != 200 {
		t.Errorf("expected http code 200, got %d", rr.Code)
	}
}

func Test_athleteHandler_ServeHTTP_error(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("GET", "/athlete", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := &athleteHandler{
		fetch: func() ([]*athlete, error) {
			return nil, fmt.Errorf("")
		},
	}

	// Act
	handler.ServeHTTP(rr, req)

	// Assert
	if rr.Code != 500 {
		t.Errorf("expected http code 500, got %d", rr.Code)
	}
}