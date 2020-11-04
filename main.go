package main

import (
	"github.com/sheitm/ofever/scrape"
)

func main(){

	seasonChan := make(chan *scrape.SeasonFetch)

	scrape.StartSeason("https://ilgeoform.no/rankinglop/", 2020, seasonChan)

	season := <- seasonChan

	x := season
	_ = x
}
