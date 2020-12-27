package persist

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/3lvia/telemetry-go"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"net/url"
	"sync"
)

type writer struct {
	connectionInfo map[string]string
	logChannels    telemetry.LogChans
}

func (w *writer) writeAll(ctx context.Context, elements []*Element, doneChan chan<- struct{}) {
	if elements == nil || len(elements) == 0 {
		doneChan <- struct{}{}
		return
	}

	credentials, accountName, err := w.credentials()
	if err != nil {
		w.logChannels.ErrorChan <- err
		doneChan <- struct{}{}
		return
	}

	p := azblob.NewPipeline(credentials, azblob.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))

	serviceURL := azblob.NewServiceURL(*u, p)
	containerURL := serviceURL.NewContainerURL(elements[0].Series)

	wg := &sync.WaitGroup{}
	wg.Add(len(elements))

	for _, element := range elements {
		go func(ctx context.Context, e *Element, cURL azblob.ContainerURL, wg *sync.WaitGroup, lc telemetry.LogChans){
			b, err := json.Marshal(e.Data)
			if err != nil {
				lc.ErrorChan <- err
				wg.Done()
				return
			}
			blobURL := cURL.NewBlockBlobURL(e.PathGetter(e.Data))
			_, err = azblob.UploadBufferToBlockBlob(ctx, b, blobURL, azblob.UploadToBlockBlobOptions{})
			if err != nil {
				lc.ErrorChan <- err
				wg.Done()
				return
			}
			wg.Done()
		}(ctx, element, containerURL, wg, w.logChannels)
	}

	wg.Wait()
	doneChan <- struct{}{}
}

func (w *writer) credentials() (*azblob.SharedKeyCredential, string, error) {
	accountName := w.connectionInfo["AccountName"]
	accountKey := w.connectionInfo["AccountKey"]
	credentials, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, "", err
	}
	return credentials, accountName, nil
}