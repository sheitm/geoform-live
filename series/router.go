package series

import (
	"fmt"
	"github.com/sheitm/ofever/sequence"
	"github.com/sheitm/ofever/types"
)

type router struct {
	seriesMap map[string]chan<- *types.SeasonFetch
}

func (r *router) start(eventChan <-chan *sequence.Event) {
	for {
		e := <- eventChan
		fetch := e.Payload.(*types.SeasonFetch)
		key := fmt.Sprintf("%s-%d", fetch.Series, fetch.Year)
		if ch, ok := r.seriesMap[key]; ok {
			ch <- fetch
			e.DoneChan <- struct{}{}
			continue
		}

		ch := make(chan *types.SeasonFetch)
		ss := &seriesSeason{}
		go ss.start(ch)
		ch <- fetch
		e.DoneChan <- struct{}{}
	}
}
