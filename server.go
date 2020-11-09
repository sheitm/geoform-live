package main

import (
	"fmt"
	"github.com/sheitm/ofever/scrape"
	"github.com/sheitm/ofever/storage"
	"log"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"strings"
	"time"
)

func startServer(port string, seasonChan chan<- *scrape.SeasonFetch){
	sHandler := &scrapeHandler{
		seasonChan: seasonChan,
		starter:    scrape.StartSeason,
	}
	http.Handle("/scrape/", sHandler)

	http.Handle("/metrics", promhttp.Handler())

	for path, handler := range storage.Handlers() {
		http.Handle(path, handler)
	}

	//fileServer := http.FileServer(http.Dir("./static"))
	//http.Handle("/", fileServer)
	//

	pp := ":" + port

	if err := http.ListenAndServe(pp, nil); err != nil {
		log.Fatal(err)
	}
}

type startScrapeFunc func(string, int, chan<- *scrape.SeasonFetch)

type scrapeHandler struct {
	seasonChan chan<- *scrape.SeasonFetch
	starter    startScrapeFunc
}

func (h *scrapeHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	arr := strings.Split(req.URL.Path, "/")
	year, err := strconv.Atoi(arr[len(arr)-1])
	if err != nil {
		//rw.Write([]byte(fmt.Sprintf("invalid request, %v", err)))
		rw.WriteHeader(500)
		return
	}
	thisYear := time.Now().Year()
	if thisYear < year {
		rw.WriteHeader(500)
		return
	}
	if year < 2009 {
		rw.WriteHeader(500)
		return
	}

	url := `https://ilgeoform.no/rankinglop/`
	if year < thisYear {
		url = fmt.Sprintf("https://ilgeoform.no/rankinglop/index-%d.html", year)
	}

	sc := make(chan *scrape.SeasonFetch)
	go h.starter(url, year, sc)

	fetch := <-sc
	h.seasonChan <- fetch

	x := 2
	_ = x
}

// https://ilgeoform.no/rankinglop/
// https://ilgeoform.no/rankinglop/index-2019.html
// https://ilgeoform.no/rankinglop/index-2009.html
