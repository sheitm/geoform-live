package storage

import (
	"encoding/json"
	"github.com/sheitm/ofever/scrape"
)

func computeSeason(js string) (*computedSeason, error) {
	var f scrape.SeasonFetch
	err := json.Unmarshal([]byte(js), &f)
	if err != nil {
		return nil, err
	}

	cs := &computedSeason{}
	cs.init(&f)

	return cs, nil
}

type computedSeason struct {
	year int
	athletes map[string]*Athlete
}

func (c *computedSeason) init(fetch *scrape.SeasonFetch) {
	athletes := map[string]*Athlete{}
	for _, result := range fetch.Results {
		if result.Event == nil {
			continue
		}
		if result.Event.Courses == nil {
			continue
		}
		eventName := result.Event.Name
		for _, course := range result.Event.Courses {
			if course.Results == nil {
				continue
			}
			for _, r := range course.Results {
				var athlete *Athlete
				var ok bool
				if athlete, ok = athletes[r.Athlete]; !ok {
					athlete = &Athlete{
						Name:    r.Athlete,
						Results: []AthleteResult{},
					}
					athletes[athlete.Name] = athlete
				}
				res := AthleteResult{
					Event:        eventName,
					Course:       course.Name,
					Disqualified: r.Disqualified,
					Placement:    r.Placement,
					ElapsedTime:  r.ElapsedTime,
					Points:       r.Points,
				}
				athlete.Results = append(athlete.Results, res)
			}
		}
	}
	c.athletes = athletes
}

