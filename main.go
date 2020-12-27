package main

import (
	"github.com/3lvia/hn-config-lib-go/vault"
	"github.com/3lvia/telemetry-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sheitm/ofever/athletes"
	"github.com/sheitm/ofever/persist"
	"github.com/sheitm/ofever/scrape"
	"github.com/sheitm/ofever/sequence"
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

	// Start persistance
	persister, reader := persist.Start(v, eventChan, logChannels)

	// Start sequencer
	sequenceTrigger := make(chan interface{})
	sequenceDone := make(chan struct{})
	sequenceAdder := sequence.Start(sequenceTrigger, sequenceDone)

	// Start scraping
	scrapeHandler := scrape.Handler(eventChan, sequenceTrigger, sequenceDone)
	httpHandler := telemetry.Wrap(scrapeHandler, logChannels)
	http.Handle("/scrape/", httpHandler)

	// Start athletes
	athletesHandler := athletes.Start(sequenceAdder, persister, reader, logChannels)
	httpHandler = telemetry.Wrap(athletesHandler, logChannels)
	http.Handle("/athletes", httpHandler)

	// Start metrics
	http.Handle("/metrics", promhttp.Handler())



	pp := ":" + port

	if err := http.ListenAndServe(pp, nil); err != nil {
		log.Fatal(err)
	}
}
