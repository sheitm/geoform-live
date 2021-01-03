package competitions

import (
	"encoding/json"
	"fmt"
	"github.com/3lvia/telemetry-go"
	"github.com/sheitm/ofever/athletes"
	"github.com/sheitm/ofever/persist"
	"github.com/sheitm/ofever/sequence"
	"github.com/sheitm/ofever/types"
	"strconv"
	"strings"
	"sync"
	"time"
)

const regexCompetitions = `[0-9]{4}/competitions/[0-9]*.json`

type impl struct {
	comps              map[string]*comp
	athleteID          athletes.AthleteIDFunc
	persistFunc        persist.Persist
	mux                *sync.Mutex
	logChannels        telemetry.LogChans
}

func (i *impl) start(eventChan <-chan *sequence.Event, reader persist.ReadFunc, readContainersFunc persist.ReadContainersFunc) {
	i.init(readContainersFunc, reader)
	for {
		e := <- eventChan
		fetch := e.Payload.(*types.SeasonFetch)

		i.logChannels.EventChan <- telemetry.Event{
			Name: "received-season",
			Data: map[string]string{
				"package": "competitions",
				"series": fetch.Series,
				"season": fmt.Sprintf("%d", fetch.Year),
			},
		}

		scraped := scrapedCompetitions(fetch)
		if scraped == nil || len(scraped) == 0 {
			e.DoneChan <- struct{}{}
			return
		}

		wg := &sync.WaitGroup{}
		wg.Add(len(scraped))

		compChan := make(chan *comp)
		for _, event := range scraped {
			go func(e *types.Event, f *types.SeasonFetch, cc chan<- *comp) {
				c := i.processScrapedComp(f, e)
				cc <- c
			}(event, fetch, compChan)
		}

		var persistElements []*persist.Element
		pf := func(obj interface{}) string {
			c := obj.(*comp)
			return fmt.Sprintf("%s/competitions/%d.json", c.Season, c.Number)
		}
		go func(f *types.SeasonFetch, cc <-chan *comp, pf persist.PathFunc, wg *sync.WaitGroup) {
			for {
				c := <-cc
				e := &persist.Element{
					Container:  f.Series,
					Data:       c,
					PathGetter: pf,
				}
				i.logChannels.DebugChan <- fmt.Sprintf("added perist element with number %d", c.Number)
				persistElements = append(persistElements, e)
				wg.Done()
			}
		}(fetch, compChan, pf, wg)

		wg.Wait()

		dc := make(chan struct{})
		i.persistFunc(persistElements, dc)

		<-dc

		e.DoneChan <- struct{}{}
	}
}

func (i *impl) init(readContainersFunc persist.ReadContainersFunc, reader persist.ReadFunc) {
	rcc := make(chan []string)
	readContainersFunc(persist.ReadContainers{Send: rcc})
	containers := <- rcc
	var readContainers []string
	for _, container := range containers {
		if container == "athletes" {
			continue
		}
		readContainers = append(readContainers, container)
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(readContainers))

	for _, container := range readContainers {
		go func(c string, rdr persist.ReadFunc, wg *sync.WaitGroup){
			send := make(chan persist.ReadResult)
			done := make(chan struct{})
			r := persist.Read{
				Container: c,
				Regex:     regexCompetitions,
				Send:      send,
				Done:      done,
			}
			go func(s <-chan persist.ReadResult) {
				for {
					b := <- s
					var co comp
					err := json.Unmarshal(b.Data, &co)
					if err != nil {
						i.logChannels.ErrorChan <- err
						continue
					}
					i.add(&co)
				}

			}(send)
			rdr(r)
			<- done
			wg.Done()
		}(container, reader, wg)
	}

	wg.Wait()

	i.logChannels.EventChan <- telemetry.Event{
		Name: "initialized",
		Data: map[string]string{
			"package": "competitions",
			"count": fmt.Sprintf("%d", len(i.comps)),
		},
	}
}

func (i *impl) processScrapedComp(fetch *types.SeasonFetch, sc *types.Event) *comp {
	c := &comp{
		Series:      fetch.Series,
		Season:      fmt.Sprintf("%d", fetch.Year),
		Number:      sc.Number,
		Name:        sc.Name,
		URLLiveLox:  sc.URLLiveLox,
		WeekDay:     sc.WeekDay,
		Date:        sc.Date,
		Place:       sc.Place,
		Organizer:   sc.Organizer,
		Responsible: sc.Responsible,
	}
	var courses []*course
	for _, scrapedCourse := range sc.Courses {
		cc := course{
			Name:       scrapedCourse.Name,
			Info:       scrapedCourse.Info,
			//Results:    nil,
		}
		var results []*result
		for _, r := range scrapedCourse.Results {
			secs, display := getElapsedTimeInfo(r.ElapsedTime)
			rr := result{
				Placement:          r.Placement,
				Disqualified:       r.Disqualified,
				AthleteID:          i.athleteID(r.Athlete, r.Club),
				Athlete:            r.Athlete,
				Club:               r.Club,
				ElapsedTimeSeconds: secs,
				ElapsedTimeDisplay: display,
				MissingControls:    r.MissingControls,
				Points:             r.Points,
			}
			results = append(results, &rr)
		}
		cc.Results = results
		courses = append(courses, &cc)
	}
	c.Courses = courses
	i.add(c)
	return c
}

func (i *impl) add(c *comp) {
	i.mux.Lock()
	defer i.mux.Unlock()

	k := compKey(c.Series, c.Season, c.Number)
	i.comps[k] = c
}

func (i *impl) get(series, season string, number int) *comp {
	k := compKey(series, season, number)
	if c, ok := i.comps[k]; ok {
		return c
	}
	return nil
}

func (i *impl) getAll(series, season string) []*compHeader {
	i.mux.Lock()
	defer i.mux.Unlock()

	pk := partialCompKey(series, season)
	var comps []*compHeader

	for k, c := range i.comps {
		if len(k) < len(pk) {
			continue
		}
		if k[0:len(pk)] == pk {
			ch := &compHeader{
				Series: c.Series,
				Season: c.Season,
				Number: c.Number,
				Name:   c.Name,
			}
			comps = append(comps, ch)
		}
	}
	return comps
}

func compKey(series, season string, number int) string {
	return fmt.Sprintf("%s%d", partialCompKey(series, season), number)
}

func partialCompKey(series, season string) string {
	return fmt.Sprintf("%s/%s/", strings.ToLower(series), strings.ToLower(season))
}

func  scrapedCompetitions(fetch *types.SeasonFetch) []*types.Event  {
	var events []*types.Event
	if fetch.Results == nil {
		return events
	}

	for _, scrapeResult := range fetch.Results {
		if scrapeResult.Event == nil || scrapeResult.Event.Courses == nil {
			continue
		}
		events = append(events, scrapeResult.Event)
	}

	return events
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