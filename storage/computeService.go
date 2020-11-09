package storage

import (
	"github.com/sheitm/ofever/scrape"
	"log"
	"sync"
)

type computeService interface {
	Start(element seasonSyncElement)
}

func newComputeService(getAthleteID athleteIDFunc) computeService {
	return &computeServiceImpl{
		computes: map[int]computedSeason{},
		mux:      &sync.Mutex{},
	}
}

type computeServiceImpl struct {
	computes     map[int]computedSeason
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
			c.computes[cs.year] = cs
			dc <- struct{}{}
		}
	}(element.seasonChan, element.doneChan)
}