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

// Course is the details for a single course of the event.
type Course struct {
	Name       string    `json:"name"`
	Info       string    `json:"info"`
	Results    []*Result `json:"results"`
	ParseError string    `json:"parse_error"`
}

// Event is the details for a single event.
type Event struct {
	// Number is the numerical order of this event in the entire season.
	Number int `json:"number"`

	// Name the name of the event.
	Name string `json:"name"`

	// Info contains textual information from the event results page.
	Info string `json:"info"`

	// URL is the address to the main event page.
	URL string `json:"url"`

	// URLInvite is the address to the invitation to the event.
	URLInvite string `json:"url_invite"`

	// URLLiveLox is the address to the event in LiveLox.
	URLLiveLox string `json:"url_live_lox"`

	Courses       []*Course `json:"courses"`
	WeekDay       string    `json:"week_day"`
	Date          time.Time `json:"date"`
	Place         string    `json:"place"`
	Organizer     string    `json:"organizer"`
	Responsible   string    `json:"responsible"`
}