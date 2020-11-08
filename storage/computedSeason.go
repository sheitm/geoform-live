package storage

import (
	"github.com/sheitm/ofever/scrape"
)

type athleteIDFunc func(string)string

func computeSeasonForFetch(f *scrape.SeasonFetch, getID athleteIDFunc) (*computedSeason, error) {
	cs := &computedSeason{}
	cs.init(f, getID)

	return cs, nil
}

type computedSeason struct {
	year int
	athletes map[string]*ComputedAthlete
}

func (c *computedSeason) init(fetch *scrape.SeasonFetch, getID athleteIDFunc) {
	athletes := map[string]*ComputedAthlete{}
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
				var athlete *ComputedAthlete
				var ok bool
				if athlete, ok = athletes[r.Athlete]; !ok {
					athlete = &ComputedAthlete{
						Name:    r.Athlete,
						ID:      getID(r.Athlete),
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

