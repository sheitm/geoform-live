package persist

import (
	"bytes"
	"context"
	"fmt"
	"github.com/3lvia/telemetry-go"
	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"net/url"
	"sync"
)

type reader struct {
	connectionInfo map[string]string
	logChannels    telemetry.LogChans
}

func (r *reader) readAll(ctx context.Context, read Read) {
	credentials, accountName, err := credentials(r.connectionInfo)
	if err != nil {
		r.logChannels.ErrorChan <- err
		read.Done <- struct{}{}
		return
	}

	p := azblob.NewPipeline(credentials, azblob.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))

	serviceURL := azblob.NewServiceURL(*u, p)
	containerURL := serviceURL.NewContainerURL(read.Series)

	err = ensureContainer(ctx, containerURL)
	if err != nil {
		r.logChannels.ErrorChan <- err
		read.Done <- struct{}{}
		return
	}

	var blobReferences []string
	for marker := (azblob.Marker{}); marker.NotDone(); {
		listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
		if err != nil {
			r.logChannels.ErrorChan <- err
			read.Done <- struct{}{}
			return
		}

		marker = listBlob.NextMarker

		for _, blobInfo := range listBlob.Segment.BlobItems {
			bi := listBlob.ServiceEndpoint + listBlob.ContainerName + "/" + blobInfo.Name
			blobReferences = append(blobReferences, bi)
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(blobReferences))

	proxy := make(chan []byte)
	go func(rr <-chan []byte, sr chan<- []byte, wg *sync.WaitGroup){
		for {
			bb := <-rr
			sr <- bb
			wg.Done()
		}
	}(proxy, read.Send, wg)

	for _, reference := range blobReferences {
		go readBlob(ctx, reference, p, proxy)
	}

	wg.Wait()
	read.Done <- struct{}{}
}

func readBlob(ctx context.Context, ref string, p pipeline.Pipeline, ch chan<- []byte) {
	u, err := url.Parse(ref)
	if err != nil {
		//log.Fatal(err)
	}
	blobURL := azblob.NewBlobURL(*u, p)
	get, err := blobURL.Download(ctx, 0, 0, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		//log.Fatal(err)
	}
	responseBody := pipeline.NewResponseBodyProgress(get.Body(azblob.RetryReaderOptions{}),
		func(bytesTransferred int64) {
			//fmt.Printf("Read %d of %d bytes.", bytesTransferred, get.ContentLength())
		})
	defer responseBody.Close()
	downloadedData := &bytes.Buffer{}
	downloadedData.ReadFrom(responseBody)

	ch <- downloadedData.Bytes()
}