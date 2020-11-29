package storage

import (
	"net/http"
	"sort"
	"time"
)

type athletePersistFunc func([]*athlete)
type athleteFetchFunc func()([]*athlete, error)
type computedSeasonsFetchFunc func()([]*computedSeason, error)
type computedSeasonFetchFunc func(int)(*computedSeason, error)
type competitionPersistFunc func([]*competition)
type competitionFetchFunc func()([]*competition, error)

type athleteIDFunc func(string)string
type competitionByNamesFunc func(eventName, courseName string) (competitionAndCourse, error)

type placementByOfficialPoints []*computedAthlete

func (a placementByOfficialPoints) Len() int           { return len(a) }
func (a placementByOfficialPoints) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a placementByOfficialPoints) Less(i, j int) bool { return a[i].PointsOfficial > a[j].PointsOfficial }

type placementByTotalPoints []*computedAthlete

func (a placementByTotalPoints) Len() int           { return len(a) }
func (a placementByTotalPoints) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a placementByTotalPoints) Less(i, j int) bool { return a[i].PointsTotal > a[j].PointsTotal }

type httpHandler interface {
	Path() string
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type athlete struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	//Results []AthleteResult `json:"results"`
}

type byPoints []athleteResult

func (a byPoints) Len() int           { return len(a) }
func (a byPoints) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPoints) Less(i, j int) bool { return a[i].Points > a[j].Points }


type computedAthlete struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	Results           []athleteResult `json:"results"`
	BestResult        athleteResult   `json:"best_result"`
	WorstResult       athleteResult   `json:"worst_result"`
	PointsTotal       float64         `json:"points_total"`
	PlacementTotal    int             `json:"placement_total"`
	PointsOfficial    float64         `json:"points_official"`
	PlacementOfficial int             `json:"placement_official"`
}

func (a *computedAthlete) computePoints(officialCount int) {
	if a.Results == nil {
		return
	}

	sort.Sort(byPoints(a.Results))

	tot := 0.0
	poc := 0.0
	lim := officialCount
	if len(a.Results) < officialCount {
		lim = len(a.Results)
	}
	for i, result := range a.Results {
		tot += result.Points
		if i < lim {
			poc += result.Points
		}
	}

	a.PointsTotal = tot
	a.BestResult = a.Results[0]
	a.WorstResult = a.Results[len(a.Results)-1]
	a.PointsOfficial = poc
}

type athleteResult struct {
	Event        string        `json:"event"`
	Course       string        `json:"course"`
	Disqualified bool          `json:"disqualified"`
	Placement    int           `json:"placement"`
	ElapsedTime  time.Duration `json:"elapsed_time"`
	Points       float64       `json:"points"`
}

type competition struct {
	ID          string    `json:"id"`
	Number      int       `json:"number"`
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	Courses     []course  `json:"courses"`
	Info        string    `json:"info"`
	URL         string    `json:"url"`
	URLInvite   string    `json:"url_invite"`
	URLLiveLox  string    `json:"url_live_lox"`
	Place       string    `json:"place"`
	Organizer   string    `json:"organizer"`
	Responsible string    `json:"responsible"`
}

type course struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Info       string  `json:"info"`
	Length     float64 `json:"length"`
	CourseType string  `json:"course_type"` // long, medium, short, newbie
}

type competitionAndCourse struct {
	competition *competition
	course      *course
}