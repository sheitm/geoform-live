package storage

import (
	"github.com/sheitm/ofever/types"
)

type seasonFetchDependency interface {
	Start(seasonSyncElement)
}

type sequentialSeasonHandlers struct {
	elements map[int]seasonSyncElement
}

func (s *sequentialSeasonHandlers) Start(seasonChan <-chan *types.SeasonFetch) {
	go func(sc <-chan *types.SeasonFetch){
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
		seasonChan: make(chan *types.SeasonFetch),
		doneChan:   make(chan struct{}),
	}
	s.elements[len(s.elements)] = element
	go dep.Start(element)
}

type seasonSyncElement struct {
	seasonChan chan *types.SeasonFetch
	doneChan   chan struct{}
}