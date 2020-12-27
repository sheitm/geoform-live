package persist

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/3lvia/hn-config-lib-go/vault"
	"github.com/3lvia/telemetry-go"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/sheitm/ofever/types"
	"net/url"
	"os"
	"strings"
)

const eventSeasonSaved = "Season_saved"

type saveFunc func(context.Context, map[string]string, *types.SeasonFetch) error

type persistRequest struct {
	elements []*Element
	doneChan chan<- struct{}
}

func newStorageService(v vault.SecretsManager, persistChan <-chan persistRequest, logChannels telemetry.LogChans) (*storageService, error) {
	// TODO: Add vault integration
	cs := os.Getenv("PERSIST_CONNECTIONSTRING")
	service := &storageService{
		connectionInfo: parseBlobConnectionString(cs),
		save:           saveFetch,
		persistChan:    persistChan,
		logChannels:    logChannels,
	}
	return service, nil
}

// r.data["primary-connection-string"]

type storageService struct {
	connectionInfo map[string]string
	save           saveFunc
	persistChan    <-chan persistRequest
	logChannels    telemetry.LogChans
}

func (s *storageService) start(eventChan <-chan *types.ScrapeEvent) {
	for {
		select {
		case e := <-eventChan:
			s.handleSeasonEvent(e)
			e.DoneChan <- nil
		case pr := <-s.persistChan:
			w := &writer{
				connectionInfo: s.connectionInfo,
				logChannels:    s.logChannels,
			}
			w.writeAll(context.Background(), pr.elements, pr.doneChan)
		}
	}
}

func (s *storageService) handleSeasonEvent(e *types.ScrapeEvent) {
	fetch := e.Fetch
	err := s.save(context.Background(), s.connectionInfo, fetch)

	if err != nil {
		s.logChannels.ErrorChan <- err
		e.DoneChan <- err
		return
	}

	s.logChannels.EventChan <-telemetry.Event{
		Name: eventSeasonSaved,
		Data: map[string]string {
			"series": fetch.Series,
			"season": fmt.Sprintf("%d", fetch.Year),
		},
	}
}

func saveFetch(ctx context.Context, config map[string]string, fetch *types.SeasonFetch) error {
	accountName := config["AccountName"]
	accountKey := config["AccountKey"]
	credentials, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return err
	}

	p := azblob.NewPipeline(credentials, azblob.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))

	serviceURL := azblob.NewServiceURL(*u, p)
	containerURL := serviceURL.NewContainerURL(fetch.Series)

	err = ensureContainer(ctx, containerURL)
	if err != nil {
		return err
	}

	b, err := json.Marshal(fetch)
	if err != nil {
		return err
	}

	blobURL := containerURL.NewBlockBlobURL(fmt.Sprintf("%d/fetch.json", fetch.Year))

	_, err = azblob.UploadBufferToBlockBlob(ctx, b, blobURL, azblob.UploadToBlockBlobOptions{})
	if err != nil {
		return err
	}

	return nil
}


func ensureContainer(ctx context.Context, containerURL azblob.ContainerURL) error {
	_, err := containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessContainer)
	if err != nil {
		if strings.Contains(err.Error(), "ServiceCode=ContainerAlreadyExists") {
			return nil
		}
		return err
	}
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
