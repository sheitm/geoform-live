package storage

import "github.com/sheitm/ofever/scrape"

type seasonFetchDependency interface {
	Start(seasonSyncElement)
}

type sequentialSeasonHandlers struct {
	elements map[int]seasonSyncElement
}

func (s *sequentialSeasonHandlers) Start(seasonChan <-chan *scrape.SeasonFetch) {
	go func(sc <-chan *scrape.SeasonFetch){
		for {
			fetch := <- sc
			for i := 0; i < len(s.elements); i++ {
				d := s.elements[i]
				d.seasonChan <- fetch
				<- d.doneChan
			}
		}
	}(seasonChan)
}

func (s *sequentialSeasonHandlers) Add(dep seasonFetchDependency) {
	if s.elements == nil {
		s.elements = map[int]seasonSyncElement{}
	}

	element := seasonSyncElement{
		seasonChan: make(chan *scrape.SeasonFetch),
		doneChan:   make(chan struct{}),
	}
	s.elements[len(s.elements)] = element
	go dep.Start(element)
}

type seasonSyncElement struct {
	seasonChan chan *scrape.SeasonFetch
	doneChan   chan struct{}
}