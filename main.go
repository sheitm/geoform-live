package main

import (
	"github.com/sheitm/ofever/scrape"
)

func main(){

	eventChan := make(chan *scrape.Result)
	doneChan := make(chan error)

	scrape.StartSeason("https://ilgeoform.no/rankinglop/", 2020, eventChan, doneChan)

	var results []*scrape.Result

	go func(ec <-chan *scrape.Result) {
		for  {
			r := <- ec
			results = append(results, r)
		}
	}(eventChan)

	<- doneChan

	x := 99
	_ = x
}
