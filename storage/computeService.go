package storage

import (
	"github.com/sheitm/ofever/scrape"
	"sync"
)

type computeService interface {
	Start(seasonChan <-chan *scrape.SeasonFetch)
}

func newComputeService() computeService {
	return &computeServiceImpl{
		computes: map[int]*computedSeason{},
		mux:      &sync.Mutex{},
	}
}

type computeServiceImpl struct {
	computes map[int]*computedSeason
	mux      *sync.Mutex
}

func (c *computeServiceImpl) Start(seasonChan <-chan *scrape.SeasonFetch) {
	//go func(sc <-chan *scrape.SeasonFetch){
	//	for {
	//		fetch := <- sc
	//		cs, err := computedSeason{fetch}
	//	}
	//}(seasonChan)
}