package persist

import (
	"github.com/3lvia/hn-config-lib-go/vault"
	"github.com/3lvia/telemetry-go"
	"github.com/sheitm/ofever/types"
	"log"
)

// Start the internal functionality in a go routine.
func Start(v vault.SecretsManager, eventChan <-chan *types.ScrapeEvent, logChannels telemetry.LogChans) (Persist, ReadFunc) {
	pr := make(chan persistRequest)
	read := make(chan Read)
	service, err := newStorageService(v, pr, read, logChannels)
	if err != nil {
		logChannels.ErrorChan <- err
		log.Fatal(service)
	}

	go service.start(eventChan)

	pf := func(elements []*Element, c chan<- struct{}) {
		pRequest := persistRequest {
			elements: elements,
			doneChan: c,
		}
		pr <- pRequest
	}
	rf := func(r Read) {
		read <- r
	}
	return pf, rf
}

// PathFunc gets the relative path to which the data should be written.
type PathFunc func(interface{}) string

// Persist can be used by client who
type Persist func([]*Element, chan<- struct{})

// ReadFunc used by clients in order to read persisted data.
type ReadFunc func(Read)

// Element represents some instance to be persisted.
type Element struct {
	Container  string
	Data       interface{}
	PathGetter PathFunc
}

// Read element used
type Read struct {
	Container string
	Path      string
	Send      chan<- []byte
	Done      chan<- struct{}
}
