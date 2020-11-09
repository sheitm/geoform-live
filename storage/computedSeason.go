package storage

import (
	"github.com/sheitm/ofever/scrape"
)

type athleteIDFunc func(string)string

type computedSeason interface {
	Athletes() map[string]*computedAthlete
	Year() int
}

func computeSeasonForFetch(f *scrape.SeasonFetch, getID athleteIDFunc) (computedSeason, error) {
	cs := &computedSeasonImpl{}
	cs.init(f, getID)

	return cs, nil
}

type computedSeasonImpl struct {
	year int
	athletes map[string]*computedAthlete
}

func (c *computedSeasonImpl) init(fetch *scrape.SeasonFetch, getID athleteIDFunc) {
	athletes := map[string]*computedAthlete{}
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
				var athlete *computedAthlete
				var ok bool
				if athlete, ok = athletes[r.Athlete]; !ok {
					athlete = &computedAthlete{
						Name:    r.Athlete,
						ID:      getID(r.Athlete),
						Results: []athleteResult{},
					}
					athletes[athlete.Name] = athlete
				}
				res := athleteResult{
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

func (c *computedSeasonImpl) Athletes() map[string]*computedAthlete {
	return c.athletes
}

