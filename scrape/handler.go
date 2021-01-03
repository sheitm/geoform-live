package scrape

import (
	"fmt"
	"github.com/3lvia/telemetry-go"
	"github.com/sheitm/ofever/types"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const handlerName = "scrape"

func Handler(eventChan chan<- *types.ScrapeEvent, sequenceTrigger chan<- interface{}, finalDone <-chan struct{}) telemetry.RequestHandler {
	return &handler{
		eventChan:       eventChan,
		starter:         StartSeason,
		sequenceTrigger: sequenceTrigger,
		finalDone:       finalDone,
	}
}

type startScrapeFunc func(string, int, chan<- *types.SeasonFetch)

type handler struct {
	eventChan       chan<- *types.ScrapeEvent
	starter         startScrapeFunc
	sequenceTrigger chan<- interface{}
	finalDone       <-chan struct{}
}

func (h *handler) Handle(r *http.Request) telemetry.RoundTrip {
	arr := strings.Split(r.URL.Path, "/")
	year, err := strconv.Atoi(arr[len(arr)-1])
	if err != nil {
		return telemetry.RoundTrip{
			HandlerName:      handlerName,
			HTTPResponseCode: 500,
		}
	}
	thisYear := time.Now().Year()
	if thisYear < year {
		return telemetry.RoundTrip{
			HandlerName:      handlerName,
			HTTPResponseCode: 500,
		}
	}
	if year < 2009 {
		return telemetry.RoundTrip{
			HandlerName:      handlerName,
			HTTPResponseCode: 500,
		}
	}

	url := `https://ilgeoform.no/rankinglop/`
	if year < thisYear {
		url = fmt.Sprintf("https://ilgeoform.no/rankinglop/index-%d.html", year)
	}

	sc := make(chan *types.SeasonFetch)
	go h.starter(url, year, sc)

	fetch := <-sc
	doneChan := make(chan error)
	ev := &types.ScrapeEvent{
		DoneChan: doneChan,
		Fetch:    fetch,
	}
	h.eventChan <- ev
	err = <- doneChan
	if err != nil {
		return telemetry.RoundTrip{
			HandlerName:      handlerName,
			HTTPResponseCode: 500,
		}
	}

	h.sequenceTrigger <- fetch
	<- h.finalDone

	return telemetry.RoundTrip{
		HandlerName:      handlerName,
		HTTPResponseCode: 200,
	}
}