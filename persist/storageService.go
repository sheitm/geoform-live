package persist

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/3lvia/hn-config-lib-go/vault"
	"github.com/3lvia/telemetry-go"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/sheitm/ofever/scrape"
	"net/url"
	"os"
	"strings"
)

const eventSeasonSaved = "Season_saved"

type saveFunc func(context.Context, map[string]string, *scrape.SeasonFetch) error

func newStorageService(v vault.SecretsManager, logChannels telemetry.LogChans) (*storageService, error) {
	// TODO: Add vault integration
	cs := os.Getenv("PERSIST_CONNECTIONSTRING")
	service := &storageService{
		connectionInfo: parseBlobConnectionString(cs),
		save:           saveFetch,
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

func (s *storageService) start(eventChan <-chan *scrape.Event) {
	for {
		e := <- eventChan
		fetch := e.Fetch
		err := s.save(context.Background(), s.connectionInfo, fetch)

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

func saveFetch(ctx context.Context, config map[string]string, fetch *scrape.SeasonFetch) error {
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
