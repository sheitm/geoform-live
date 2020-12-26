package persist

import (
	"github.com/3lvia/hn-config-lib-go/vault"
	"github.com/3lvia/telemetry-go"
	"github.com/sheitm/ofever/types"
	"log"
)

// Start the internal functionality in a go routine.
func Start(v vault.SecretsManager, eventChan <-chan *types.ScrapeEvent, logChannels telemetry.LogChans) {
	service, err := newStorageService(v, logChannels)
	if err != nil {
		logChannels.ErrorChan <- err
		log.Fatal(service)
	}

	go service.start(eventChan)
}
