package athletes

import (
	"encoding/json"
	"github.com/3lvia/telemetry-go"
	"net/http"
)

const handlerName = "athletes"

type handler struct {
	c           *cache
	logChannels telemetry.LogChans
}

func (h *handler) Handle(r *http.Request) telemetry.RoundTrip {
	all := h.c.all()
	b, err := json.Marshal(all)

	if err != nil {
		h.logChannels.ErrorChan <- err
		return telemetry.RoundTrip{
			HandlerName:      handlerName,
			HTTPResponseCode: 500,
			Contents:         nil,
		}
	}

	return telemetry.RoundTrip{
		HandlerName:      handlerName,
		HTTPResponseCode: 200,
		Contents:         b,
	}
}