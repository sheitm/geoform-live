package main

import (
	"github.com/3lvia/hn-config-lib-go/vault"
	"github.com/3lvia/telemetry-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sheitm/ofever/persist"
	"github.com/sheitm/ofever/scrape"
	"github.com/sheitm/ofever/types"
	"log"
	"net/http"
)

func main(){
	port := "2112"

	v, err := vault.New()
	if err != nil {
		log.Fatal(err)
	}

	logChannels := telemetry.StartEmpty()
	eventChan := make(chan *types.ScrapeEvent)

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
