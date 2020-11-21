package storage

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func newComputedSeasonHandler(fetch computedSeasonFetchFunc) httpHandler {
	return &computedSeasonHandler{fetch: fetch}
}

type computedSeasonHandler struct {
	fetch computedSeasonFetchFunc
}

func (h *computedSeasonHandler) Path() string {
	return "/season/"
}

func (h *computedSeasonHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	p := splitPath(req.URL.Path)
	if len(p) <= 1 {
		rw.WriteHeader(500)
		return
	}
	y, err := strconv.Atoi(p[1])
	if err != nil {
		rw.WriteHeader(500)
		return
	}

	cs, err := h.fetch(y)
	if err != nil {
		rw.WriteHeader(500)
		return
	}

	b, err := json.Marshal(cs.dto())

	if err != nil {
		rw.WriteHeader(500)
		return
	}
	rw.Write(b)
}

func splitPath(p string) []string {
	a := strings.Split(p, "/")
	var res []string
	for _, s := range a {
		if s != "" {
			res = append(res, s)
		}
	}
	return res
}