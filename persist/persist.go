package persist

import (
	"github.com/3lvia/hn-config-lib-go/vault"
	"github.com/3lvia/telemetry-go"
	"github.com/sheitm/ofever/types"
	"log"
)

// Start the internal functionality in a go routine.
func Start(v vault.SecretsManager, eventChan <-chan *types.ScrapeEvent, logChannels telemetry.LogChans) Persist {
	service, err := newStorageService(v, logChannels)
	if err != nil {
		logChannels.ErrorChan <- err
		log.Fatal(service)
	}

	go service.start(eventChan)

	return func(elements []*Element, c chan<- struct{}) {

	}
}

// PathFunc gets the relative path to which the data should be written.
type PathFunc func(interface{}) string

// Persist can be used by client who
type Persist func([]*Element, chan<- struct{})

// Element represents some instance to be persisted or fetched.
type Element struct {
	Series     string
	Data       interface{}
	PathGetter PathFunc
}
