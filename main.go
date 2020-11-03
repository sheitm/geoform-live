package main

import (
	"github.com/sheitm/ofever/scrape"
	"sync"
)

func main(){

	eventChan := make(chan *scrape.Result)

	scrape.StartScrape("https://ilgeoform.no/rankinglop/res2020-10-03.html", eventChan)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(ec <-chan *scrape.Result, wg *sync.WaitGroup) {
		r := <- ec
		e := r.Event
		_ = e
		wg.Done()
	}(eventChan, wg)

	wg.Wait()
}
