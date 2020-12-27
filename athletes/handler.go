package athletes

import (
	"github.com/3lvia/telemetry-go"
	"net/http"
)

const handlerName = "athletes"

type handler struct {}

func (h *handler) Handle(r *http.Request) telemetry.RoundTrip {
	return telemetry.RoundTrip{
		HandlerName:      handlerName,
		HTTPResponseCode: 200,
		Contents:         nil,
	}
}