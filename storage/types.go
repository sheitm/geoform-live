package storage

import (
	"net/http"
	"time"
)

type athletePersistFunc func([]*athlete)
type athleteFetchFunc func()([]*athlete, error)
type computedSeasonsFetchFunc func()([]*computedSeason, error)
type computedSeasonFetchFunc func(int)(*computedSeason, error)

type httpHandler interface {
	Path() string
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type athlete struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	//Results []AthleteResult `json:"results"`
}

type computedAthlete struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Results []athleteResult `json:"results"`
}

type athleteResult struct {
	Event        string        `json:"event"`
	Course       string        `json:"course"`
	Disqualified bool          `json:"disqualified"`
	Placement    int           `json:"placement"`
	ElapsedTime  time.Duration `json:"elapsed_time"`
	Points       float64       `json:"points"`
}
