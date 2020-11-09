package storage

import (
	"encoding/json"
	"net/http"
)

func newAthleteHandler(fetch athleteFetchFunc) httpHandler {
	return &athleteHandler{fetch: fetch}
}

type athleteHandler struct {
	fetch athleteFetchFunc
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
	rw.Write(b)
}

