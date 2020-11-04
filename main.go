package main

import (
	"github.com/sheitm/ofever/scrape"
	"github.com/sheitm/ofever/storage"
	"time"
)

func main(){

	seasonChan := make(chan *scrape.SeasonFetch)

	storage.Start(seasonChan)

	scrape.StartSeason("https://ilgeoform.no/rankinglop/", 2020, seasonChan)

	<- time.After(1 * time.Minute)
}
