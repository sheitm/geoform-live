package storage

import (
	"github.com/google/uuid"
	"github.com/sheitm/ofever/scrape"
	"log"
	"sync"
)

type athletePersistFunc func([]*Athlete)
type athleteFetchFunc func()([]*Athlete, error)

type athleteService interface {
	Start(element seasonSyncElement)
	List() ([]*Athlete, error)
	ID(name string) string
}

func newAthleteService(persist athletePersistFunc, fetch athleteFetchFunc) athleteService {
	impl := &athleteServiceImpl{
		byName:  map[string]*Athlete{},
		byID:    map[string]*Athlete{},
		persist: persist,
		mux:     &sync.Mutex{},
	}
	impl.init(fetch)
	return impl
}

type athleteServiceImpl struct {
	byName  map[string]*Athlete
	byID    map[string]*Athlete
	persist athletePersistFunc
	mux     *sync.Mutex
}

func (a *athleteServiceImpl) Start(element seasonSyncElement){
	go func(sc <-chan *scrape.SeasonFetch, dc chan<- struct{}){
		for {
			fetch := <-sc
			anyChange := false
			if fetch.Results == nil {
				continue
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
						if r.Athlete == "" {
							continue
						}
						if _, ok := a.byName[r.Athlete]; !ok {
							a.newAthlete(r.Athlete)
							anyChange = true
						}
					}
				}
			}
			if anyChange {
				l, err := a.List()
				if err != nil {
					log.Printf("%v", err)
					continue
				}
				a.persist(l)
			}
			dc <- struct{}{}
		}
	}(element.seasonChan, element.doneChan)
}

func (a *athleteServiceImpl) ID(name string) string{
	if athlete, ok := a.byName[name]; ok {
		return athlete.ID
	}

	a.newAthlete(name)
	return a.byName[name].ID
}

func (a *athleteServiceImpl) List() ([]*Athlete, error) {
	var res []*Athlete
	for _, athlete := range a.byID {
		res = append(res, athlete)
	}
	return res, nil
}

func (a *athleteServiceImpl) init(fetch athleteFetchFunc) {
	l, err := fetch()
	if err != nil{
		return
	}

	for _, athlete := range l {
		a.byName[athlete.Name] = athlete
		a.byID[athlete.ID] = athlete
	}
}

func (a *athleteServiceImpl) newAthlete(name string) {
	a.mux.Lock()
	defer a.mux.Unlock()

	guid := uuid.New()
	athlete := &Athlete{
		ID:   guid.String(),
		Name: name,
	}
	a.byName[athlete.Name] = athlete
	a.byID[athlete.ID] = athlete
}