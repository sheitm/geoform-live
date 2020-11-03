package contracts

import "time"

// Result is a structured version of a single line in a results table for a single event. With "Raw" we mean that
// references to the entities for athlete and club are just texts and not references to objects.
type Result struct {
	Placement       int           `json:"placement"`
	Disqualified    bool          `json:"disqualified"`
	Athlete         string        `json:"name"`
	Club            string        `json:"club"`
	ElapsedTime     time.Duration `json:"elapsed_time"`
	MissingControls int           `json:"missing_controls"`
	Points          float64       `json:"points"`
}

type Course struct {
	Name    string    `json:"name"`
	Info    string    `json:"info"`
	Results []*Result `json:"results"`
}

type Event struct {
	Name    string `json:"name"`
	Info    string `json:"info"`
	Courses []*Course
}