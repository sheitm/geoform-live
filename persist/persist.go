package persist

import (
	"github.com/3lvia/hn-config-lib-go/vault"
	"github.com/3lvia/telemetry-go"
	"github.com/sheitm/ofever/scrape"
	"log"
)

// Start the internal functionality in a go routine.
func Start(v vault.SecretsManager, eventChan <-chan *scrape.Event, logChannels telemetry.LogChans) {
	service, err := newStorageService(v, logChannels)
	if err != nil {
		logChannels.ErrorChan <- err
		log.Fatal(service)
	}

	go service.start(eventChan)
}
