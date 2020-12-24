package storage

import (
	"fmt"
	"github.com/sheitm/ofever/scrape"
	"log"
	"sort"
	"strconv"
	"time"
)

func computeSeasonForFetch(f *scrape.SeasonFetch, getID athleteIDFunc, getCompetition competitionByNamesFunc) (*computedSeason, error) {
	cs := &computedSeason{}
	cs.init(f, getID, getCompetition)
	cs.computePointsAndPlacements()
	return cs, nil
}

type computedSeasonDTO struct {
	Year         int                `json:"year"`
	Athletes     []*computedAthlete `json:"athletes"`
	Competitions []*competition     `json:"competitions"`
	Statistics   *seasonStatistics  `json:"statistics"`
}

type computedSeason struct {
	Year             int                         `json:"year"`
	Athletes         map[string]*computedAthlete `json:"athletes"`
	Competitions     []*competition              `json:"competitions"`
	Statistics       *seasonStatistics           `json:"statistics"`
}

func (c *computedSeason) dto() *computedSeasonDTO {
	return &computedSeasonDTO{
		Year:         c.Year,
		Athletes:     c.athleteSlice(),
		Statistics:   c.Statistics,
		Competitions: c.Competitions,
	}
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
	oec := c.Statistics.ValidEventsCount / 2
	if c.Statistics.ValidEventsCount % 2 != 0 {
		oec++
	}
	return oec
}

func (c *computedSeason) init(fetch *scrape.SeasonFetch, getID athleteIDFunc, getCompetition competitionByNamesFunc) {
	validEventCount := 0
	c.Year = fetch.Year
	competitions := map[string]*competition{}
	athletes := map[string]*computedAthlete{}
	totalStarts := 0
	totalDuration := 0 * time.Second
	for _, result := range fetch.Results {
		if result.Event == nil {
			continue
		}
		if result.Event.Courses == nil {
			continue
		}
		validEventCount++
		for _, course := range result.Event.Courses {
			compAndCourse, err := getCompetition(result.Event.Name, course.Name)
			if err != nil {
				log.Fatal(err)
			}
			competitions[compAndCourse.competition.ID] = compAndCourse.competition
			if course.Results == nil {
				continue
			}
			for _, r := range course.Results {
				totalStarts++
				totalDuration = totalDuration + r.ElapsedTime
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
				seconds, displayString := getElapsedTimeInfo(r.ElapsedTime)
				res := athleteResult{
					Event:              compAndCourse.competition.ID,
					Course:             compAndCourse.course.ID,
					Disqualified:       r.Disqualified,
					Placement:          r.Placement,
					ElapsedTimeSeconds: seconds,
					ElapsedTimeDisplay: displayString,
					Points:             r.Points,
				}
				athlete.Results = append(athlete.Results, res)
			}
		}
	}
	c.Athletes = athletes

	secs, displayTime := getElapsedTimeInfo(totalDuration)
	c.Statistics = &seasonStatistics{
		EventsCount:             len(fetch.Results),
		ValidEventsCount:        validEventCount,
		StartsCount:             totalStarts,
		AthletesCount:           len(athletes),
		TotalElapsedTimeSeconds: secs,
		TotalElapsedTimeDisplay: displayTime,
	}

	var comps []*competition
	for _, comp := range competitions {
		comps = append(comps, comp)
	}
	c.Competitions = comps
}

func getElapsedTimeInfo(d time.Duration) (int, string) {
	ss := fmt.Sprintf("%.0f", d.Seconds())
	totalSeconds, _ := strconv.Atoi(ss)

	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second


	return totalSeconds, fmt.Sprintf("%d:%02d:%02d", h, m, s)
}


