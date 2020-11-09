package storage

import (
	"fmt"
	"github.com/sheitm/ofever/scrape"
	"log"
	"sync"
)

type computeService interface {
	Start(element seasonSyncElement)
	ComputedSeason(year int) (*computedSeason, error)
}

func newComputeService(getAthleteID athleteIDFunc) computeService {
	return &computeServiceImpl{
		computes:     map[int]*computedSeason{},
		mux:          &sync.Mutex{},
		getAthleteID: getAthleteID,
	}
}

type computeServiceImpl struct {
	computes     map[int]*computedSeason
	mux          *sync.Mutex
	getAthleteID athleteIDFunc
}

func (c *computeServiceImpl) Start(element seasonSyncElement) {
	go func(sc <-chan *scrape.SeasonFetch, dc chan<- struct{}){
		for {
			fetch := <- sc
			cs, err := computeSeasonForFetch(fetch, c.getAthleteID)
			if err != nil {
				log.Print(err)
				dc <- struct{}{}
				continue
			}
			c.computes[cs.Year] = cs
			dc <- struct{}{}
		}
	}(element.seasonChan, element.doneChan)
}

func (c *computeServiceImpl) ComputedSeason(year int) (*computedSeason, error) {
	if cs, ok := c.computes[year]; ok {
		return cs, nil
	}
	return nil, fmt.Errorf("missing computation for year %d", year)
}