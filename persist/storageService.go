package persist

import (
	"fmt"
	"github.com/3lvia/hn-config-lib-go/vault"
	"github.com/3lvia/telemetry-go"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/sheitm/ofever/scrape"
	"os"
	"strings"
)

const eventSeasonSaved = "Season_saved"

type saveFunc func(*scrape.SeasonFetch) error

func newStorageService(v vault.SecretsManager, logChannels telemetry.LogChans) (*storageService, error) {
	// TODO: Add vault integration
	cs := os.Getenv("PERSIST_CONNECTIONSTRING")
	service := &storageService{
		connectionInfo: parseBlobConnectionString(cs),
		logChannels:    logChannels,
	}
	return service, nil
}

// r.data["primary-connection-string"]

type storageService struct {
	connectionInfo map[string]string
	save           saveFunc
	logChannels    telemetry.LogChans
}

func (s *storageService) start(eventChan <-chan *Event) {
	for {
		e := <- eventChan
		fetch := e.Fetch
		err := s.save(fetch)

		if err != nil {
			s.logChannels.ErrorChan <- err
			e.DoneChan <- err
			continue
		}

		s.logChannels.EventChan <-telemetry.Event{
			Name: eventSeasonSaved,
			Data: map[string]string {
				"series": fetch.Series,
				"season": fmt.Sprintf("%d", fetch.Year),
			},
		}

		e.DoneChan <- nil
	}
}

func (s *storageService) xx(fetch *scrape.SeasonFetch) error {
	accountName := s.connectionInfo["AccountName"]
	accountKey := s.connectionInfo["AccountKey"]
	credentials, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return err
	}
	_ = credentials
	return nil
}

func parseBlobConnectionString(s string) map[string]string  {
	m := map[string]string{}
	split := strings.Split(s, ";")
	for _, s := range split {
		i := strings.Index(s, "=")
		m[s[0:i]] = s[i+1:]
	}

	return m
}
