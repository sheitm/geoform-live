package persist

import (
	"github.com/3lvia/hn-config-lib-go/vault"
	"github.com/3lvia/telemetry-go"
	"github.com/sheitm/ofever/types"
	"log"
)

// Start the internal functionality in a go routine.
func Start(
	v vault.SecretsManager,
	eventChan <-chan *types.ScrapeEvent,
	logChannels telemetry.LogChans) (Persist, ReadFunc, ReadContainersFunc) {
	pr := make(chan persistRequest)
	read := make(chan Read)
	readContainers := make(chan ReadContainers)
	service, err := newStorageService(v, pr, read, readContainers, logChannels)
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
	rcf := func(rc ReadContainers) {
		readContainers <- rc
	}
	return pf, rf, rcf
}

// PathFunc gets the relative path to which the data should be written.
type PathFunc func(interface{}) string

// Persist can be used by client who
type Persist func([]*Element, chan<- struct{})

// ReadFunc used by clients in order to read persisted data.
type ReadFunc func(Read)

// ReadContainersFunc gets all existing containers in the current storage account.
type ReadContainersFunc func(ReadContainers)

// Element represents some instance to be persisted.
type Element struct {
	Container  string
	Data       interface{}
	PathGetter PathFunc
}

// Read element used
type Read struct {
	Container string
	Regex     string
	Send      chan<- ReadResult
	Done      chan<- struct{}
}

type ReadResult struct {
	Path string
	Data []byte
}

// ReadContainers used when clients want to know which containers exist.
type ReadContainers struct {
	Send      chan<- []string
}