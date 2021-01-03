package types

import (
	"time"
)

// SeasonFetch the results after having attempted to get all events for a season.
type SeasonFetch struct {
	Series  string    `json:"series"`
	Year    int       `json:"year"`
	URL     string    `json:"url"`
	Results []*ScrapeResult `json:"results"`
	Error   string    `json:"error"`
}

// ScrapeResult is the result of a single event scraping session.
type ScrapeResult struct {
	// URL is the provided address to the page containing the results for the event.
	URL string `json:"url"`

	// Event is the details of the event.
	Event *Event `json:"event"`

	// Error details if any error occurred during execution.
	Error string `json:"error"`
}

// ScrapeEvent is used to signal that a season fetch should be persisted.
type ScrapeEvent struct {
	// DoneChan used to signal back that the event is persisted. If error is not nil, it failed.
	DoneChan chan<- error

	// Fetch is the fetch to persist.
	Fetch *SeasonFetch
}

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