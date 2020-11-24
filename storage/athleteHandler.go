package storage

import (
	"encoding/json"
	"net/http"
)

func newAthleteHandler(fetch athleteFetchFunc, csFetch computedSeasonFetchFunc) httpHandler {
	return &athleteHandler{
		fetch:   fetch,
		csFetch: csFetch,
	}
}

type athleteHandler struct {
	fetch   athleteFetchFunc
	csFetch computedSeasonFetchFunc
}

func (h *athleteHandler) Path() string {
	return "/athlete/"
}

func (h *athleteHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	athletes, err := h.fetch()
	if err != nil {
		rw.WriteHeader(500)
		return
	}

	b, err := json.Marshal(athletes)
	if err != nil {
		rw.WriteHeader(500)
		return
	}
	rw.Header().Add("Access-Control-Allow-Origin", "*")
	rw.Write(b)
}

