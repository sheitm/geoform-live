package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sheitm/ofever/scrape"
	"github.com/sheitm/ofever/storage"
	"log"
	"net/http"
)

func startServer(port string, seasonChan chan<- *scrape.SeasonFetch){


	http.Handle("/metrics", promhttp.Handler())

	for path, handler := range storage.Handlers() {
		http.Handle(path, handler)
	}

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)

	pp := ":" + port

	if err := http.ListenAndServe(pp, nil); err != nil {
		log.Fatal(err)
	}
}

type startScrapeFunc func(string, int, chan<- *scrape.SeasonFetch)


// https://ilgeoform.no/rankinglop/
// https://ilgeoform.no/rankinglop/index-2019.html
// https://ilgeoform.no/rankinglop/index-2009.html
