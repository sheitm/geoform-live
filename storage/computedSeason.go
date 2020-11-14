package storage

import (
	"github.com/sheitm/ofever/scrape"
	"sort"
)

type athleteIDFunc func(string)string

type placementByOfficialPoints []*computedAthlete

func (a placementByOfficialPoints) Len() int           { return len(a) }
func (a placementByOfficialPoints) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a placementByOfficialPoints) Less(i, j int) bool { return a[i].PointsOfficial > a[j].PointsOfficial }

type placementByTotalPoints []*computedAthlete

func (a placementByTotalPoints) Len() int           { return len(a) }
func (a placementByTotalPoints) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a placementByTotalPoints) Less(i, j int) bool { return a[i].PointsTotal > a[j].PointsTotal }

func computeSeasonForFetch(f *scrape.SeasonFetch, getID athleteIDFunc) (*computedSeason, error) {
	cs := &computedSeason{}
	cs.init(f, getID)
	cs.computePointsAndPlacements()
	return cs, nil
}

type computedSeason struct {
	Year             int                         `json:"year"`
	Athletes         map[string]*computedAthlete `json:"athletes"`
	EventsCount      int                         `json:"events_count"`
	ValidEventsCount int                         `json:"valid_events_count"`
}

func (c *computedSeason) computePointsAndPlacements() {
	officialEventCount := c.officialEventCount()
	for _, a := range c.Athletes {
		a.computePoints(officialEventCount)
	}

	as := c.athleteSlice()

	sort.Sort(placementByTotalPoints(as))
	for i, a := range as {
		a.PlacementTotal = i+1
	}

	sort.Sort(placementByOfficialPoints(as))
	for i, a := range as {
		a.PlacementOfficial = i+1
	}
}

func (c *computedSeason) athleteSlice() []*computedAthlete {
	var as []*computedAthlete
	for _, ath := range c.Athletes {
		as = append(as, ath)
	}
	return as
}

func (c *computedSeason) officialEventCount() int {
	oec := c.ValidEventsCount / 2
	if c.ValidEventsCount % 2 != 0 {
		oec++
	}
	return oec
}

func (c *computedSeason) init(fetch *scrape.SeasonFetch, getID athleteIDFunc) {
	validEventCount := 0
	c.Year = fetch.Year
	athletes := map[string]*computedAthlete{}
	for _, result := range fetch.Results {
		if result.Event == nil {
			continue
		}
		if result.Event.Courses == nil {
			continue
		}
		validEventCount++
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
	c.Athletes = athletes
	c.EventsCount = len(fetch.Results)
	c.ValidEventsCount = validEventCount
}
