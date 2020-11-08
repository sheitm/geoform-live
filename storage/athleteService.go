package storage

import (
	"github.com/google/uuid"
	"github.com/sheitm/ofever/scrape"
	"log"
)

type athletePersistFunc func([]*Athlete)

type athleteService interface {
	Start(seasonChan <-chan *scrape.SeasonFetch)
	List() ([]*Athlete, error)
}

func newAthleteService(persistArgs ...athletePersistFunc) athleteService {
	var persist athletePersistFunc
	if len(persistArgs) > 0 {
		persist = persistArgs[0]
	} else {
		persist = func(athletes []*Athlete) {
			fn := func(interface{}) string { return "athletes.json" }
			currentStorageService.Store(athletes, fn)
		}
	}

	return & athleteServiceImpl{
		byName:  map[string]*Athlete{},
		byID:    map[string]*Athlete{},
		persist: persist,
	}
}

type athleteServiceImpl struct {
	byName  map[string]*Athlete
	byID    map[string]*Athlete
	persist athletePersistFunc
}

func (a *athleteServiceImpl) Start(seasonChan <-chan *scrape.SeasonFetch){
	go func(sc <-chan *scrape.SeasonFetch){
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
		}
	}(seasonChan)
}

func (a *athleteServiceImpl) List() ([]*Athlete, error) {
	var res []*Athlete
	for _, athlete := range a.byID {
		res = append(res, athlete)
	}
	return res, nil
}

func (a *athleteServiceImpl) newAthlete(name string) {
	guid := uuid.New()
	athlete := &Athlete{
		ID:   guid.String(),
		Name: name,
	}
	a.byName[athlete.Name] = athlete
	a.byID[athlete.ID] = athlete
}