package main

import (
	"github.com/3lvia/hn-config-lib-go/vault"
	"github.com/3lvia/telemetry-go"
	"github.com/sheitm/ofever/persist"
	"github.com/sheitm/ofever/scrape"
	"github.com/sheitm/ofever/storage"
	"log"
	"os"
)

func main(){
	v, err := vault.New()
	if err != nil {
		log.Fatal(err)
	}

	logChannels := telemetry.StartEmpty()
	seasonChan := make(chan *scrape.SeasonFetch)

	persist.Start(v, seasonChan, logChannels)
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