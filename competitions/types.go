package competitions

import "time"

type comp struct {
	Series string `json:"series"`
	Season string `json:"season"`

	// Number is the numerical order of this event in the entire season.
	Number int `json:"number"`

	Name        string    `json:"name"`
	URLLiveLox  string    `json:"url_live_lox"`
	Courses     []*course `json:"courses"`
	WeekDay     string    `json:"week_day"`
	Date        time.Time `json:"date"`
	Place       string    `json:"place"`
	Organizer   string    `json:"organizer"`
	Responsible string    `json:"responsible"`
}

type course struct {
	Name       string    `json:"name"`
	Info       string    `json:"info"`
	Results    []*result `json:"results"`
}

type result struct {
	Placement          int     `json:"placement"`
	Disqualified       bool    `json:"disqualified"`
	AthleteID          string  `json:"athlete_id"`
	Athlete            string  `json:"name"`
	Club               string  `json:"club"`
	ElapsedTimeSeconds int     `json:"elapsed_time_seconds"`
	ElapsedTimeDisplay string  `json:"elapsed_time_display"`
	MissingControls    int     `json:"missing_controls"`
	Points             float64 `json:"points"`
}