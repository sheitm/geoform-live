package storage

import (
	"github.com/sheitm/ofever/scrape"
	"sync"
	"testing"
	"time"
)

func Test_sequentialSeasonHandlers_Start(t *testing.T) {
	// Arrange
	seq := &sequentialSeasonHandlers{}
	c := 0
	wg := &sync.WaitGroup{}
	wg.Add(1)
	done := func(){
		c++
		if c >= 3 {
			wg.Done()
		}
	}
	d1 := &testDep{done: done}
	seq.Add(d1)
	d2 := &testDep{done: done}
	seq.Add(d2)
	d3 := &testDep{done: done}
	seq.Add(d3)

	seasonChan := make(chan *scrape.SeasonFetch)

	fetch := &scrape.SeasonFetch{Year: 2012}

	// Act
	seq.Start(seasonChan)
	seasonChan <- fetch

	wg.Wait()

	// Assert
	if d1.receivedFetch.Year != fetch.Year {
		t.Error("unexpected fetch received for d1")
	}
	if d2.receivedFetch.Year != fetch.Year {
		t.Error("unexpected fetch received for d2")
	}
	if d3.receivedFetch.Year != fetch.Year {
		t.Error("unexpected fetch received for d3")
	}
	if d1.when >= d2.when {
		t.Errorf("expected d1 before d2, %d -> %d", d1.when, d2.when)
	}
	if d2.when >= d3.when {
		t.Errorf("expected d2 before d3, %d -> %d", d1.when, d2.when)
	}
}

type testDep struct {
	receivedFetch *scrape.SeasonFetch
	when          int64
	done          func()
}

func (t *testDep) Start(element seasonSyncElement) {
	go func(sc <-chan *scrape.SeasonFetch, dc chan<- struct{}){
		for {
			fetch := <- sc
			t.receivedFetch = fetch
			t.when = time.Now().UnixNano()
			<- time.After(1 * time.Millisecond)
			dc <- struct{}{}
			t.done()
		}
	}(element.seasonChan, element.doneChan)
}
