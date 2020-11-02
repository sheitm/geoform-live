package contracts

import "time"

// RawResult is a structured version of a single line in a results table for a single event. With "Raw" we mean that
// references to the entities for athlete and club are just texts and not references to objects.
type RawResult struct {
	Placement    int           `json:"placement"`
	Disqualified bool          `json:"disqualified"`
	Athlete      string        `json:"name"`
	Club         string        `json:"club"`
	ElapsedTime  time.Duration `json:"elapsed_time"`
	Points       float64       `json:"points"`
}