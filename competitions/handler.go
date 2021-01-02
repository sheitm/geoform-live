package competitions

import (
	"encoding/json"
	"github.com/3lvia/telemetry-go"
	"net/http"
	"strconv"
	"strings"
)

const handlerName = "competitions"

type getCompFunc func(series, season string, number int) *comp
type getAllCompsFunc func(series, season string) []*comp

type handler struct {
	get    getCompFunc
	getAll getAllCompsFunc
}

func (h *handler) Handle(r *http.Request) telemetry.RoundTrip {
	// /competitions/geoform/2020/1
	p := strings.Split(r.URL.Path, "/")
	if len(p) < 4 {
		return telemetry.RoundTrip{
			HandlerName:      handlerName,
			HTTPResponseCode: 500,
			Contents:         nil,
		}
	}
	var b []byte
	var err error
	if len(p) >= 5 {
		b, err = h.single(p[2], p[3], p[4])
	} else {
		b, err = h.all(p[2], p[3])
	}

	if err != nil {
		return telemetry.RoundTrip{
			HandlerName:      handlerName,
			HTTPResponseCode: 500,
			Contents:         []byte(err.Error()),
		}
	}

	return telemetry.RoundTrip{
		HandlerName:      handlerName,
		HTTPResponseCode: 200,
		Contents:         b,
	}
}

func (h *handler) all(series, season string) ([]byte, error) {
	l := h.getAll(series, season)
	if l == nil {
		return nil, nil
	}
	s, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (h *handler) single(series, season, number string) ([]byte, error) {
	n, err := strconv.Atoi(number)
	if err != nil {
		return nil, err
	}
	c := h.get(series, season, n)
	if c == nil {
		return nil, nil
	}
	s, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return s, nil
}

