package athletes

import (
	"fmt"
	"github.com/3lvia/telemetry-go"
	"github.com/sheitm/ofever/persist"
	"github.com/sheitm/ofever/sequence"
	"github.com/sheitm/ofever/types"
	"sync"
)

func Start(sequenceAdder sequence.Adder, persister persist.Persist, logChannels telemetry.LogChans) telemetry.RequestHandler {
	seasonChan := make(chan *sequence.Event)
	sequenceAdder(seasonChan)
	c := &cache{
		competitorsBySHA:  map[string]*athleteWithID{},
		competitorsByGuid: map[string]*athleteWithID{},
		mux:               sync.Mutex{},
	}
	i := &impl{
		cache:       c,
		seasonChan:  seasonChan,
		persister:   persister,
		logChannels: logChannels,
	}

	go i.start()

	return &handler{}
}

type impl struct {
	cache       *cache
	persister   persist.Persist
	seasonChan  <-chan *sequence.Event
	logChannels telemetry.LogChans
}

func (a *impl) start() {
	for {
		e := <- a.seasonChan
		fetch := e.Payload.(*types.SeasonFetch)

		results := athleteResults(fetch)
		if results == nil {
			e.DoneChan <- struct{}{}
			continue
		}

		var newAthletes []*persist.Element
		for _, result := range results {
			a, existed := a.cache.competitor(result.Athlete, result.Club)
			if !existed {
				element := &persist.Element{
					Series:     "athletes",
					Data:       a,
					PathGetter: athletePath,
				}
				newAthletes = append(newAthletes, element)
			}
		}

		written := make(chan struct{})
		a.persister(newAthletes, written)

		<- written
		e.DoneChan <- struct{}{}
	}
}

func athletePath(e interface{}) string {
	a := e.(*athleteWithID)
	return fmt.Sprintf("%s.json", a.ID)
}

func athleteResults(fetch *types.SeasonFetch) []*types.Result {
	var results []*types.Result
	if fetch.Results == nil {
		return results
	}

	for _, result := range fetch.Results {
		if result.Event == nil || result.Event.Courses == nil {
			continue
		}

		for _, course := range result.Event.Courses {
			if course.Results == nil {
				continue
			}
			for _, r := range course.Results {
				results = append(results, r)
			}
		}
	}

	return results
}