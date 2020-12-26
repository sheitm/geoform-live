package main

import (
	"github.com/3lvia/hn-config-lib-go/vault"
	"github.com/3lvia/telemetry-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sheitm/ofever/persist"
	"github.com/sheitm/ofever/scrape"
	"github.com/sheitm/ofever/storage"
	"log"
	"net/http"
	"os"
)

func main(){
	port := "2112"

	v, err := vault.New()
	if err != nil {
		log.Fatal(err)
	}

	logChannels := telemetry.StartEmpty()
	eventChan := make(chan *scrape.Event)

	// Start scraping
	scrapeHandler := scrape.Handler(eventChan)
	httpHandler := telemetry.Wrap(scrapeHandler, logChannels)
	http.Handle("/scrape/", httpHandler)

	// Start metrics
	http.Handle("/metrics", promhttp.Handler())

	persist.Start(v, eventChan, logChannels)

	pp := ":" + port

	if err := http.ListenAndServe(pp, nil); err != nil {
		log.Fatal(err)
	}
}
//http://localhost:2112/scrape/2019

func old() {
	storageDirectory := os.Getenv("STORAGE_DIRECTORY")
	if storageDirectory == "" {
		log.Fatal("environment variable STORAGE_DIRECTORY must be set")
	}
	seasonChan := make(chan *scrape.SeasonFetch)

	storage.Start(storageDirectory, seasonChan)

	// startServer must be last line
	startServer("2112", seasonChan)
}