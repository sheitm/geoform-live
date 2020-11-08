package storage

import (
	"encoding/json"
	"fmt"
	"github.com/sheitm/ofever/scrape"
	"io/ioutil"
	"log"
	"path"
	"sync"
)

type storageService interface {
	//Store(fetch *scrape.SeasonFetch) error
	Start(seasonChan <-chan *scrape.SeasonFetch)
	Store(obj interface{}, fn fileNameFunc) error
}

type fileNameFunc func(interface{}) string

func newStorageService(storageFolder string) storageService {
	return &storageServiceImpl{
		folder: storageFolder,
		mux:    &sync.Mutex{},
	}
}

type storageServiceImpl struct {
	folder string
	mux *sync.Mutex
}

func (s *storageServiceImpl) Start(seasonChan <-chan *scrape.SeasonFetch) {
	go func(sc <-chan *scrape.SeasonFetch) {
		for {
			fetch := <- seasonChan
			fn := func(obj interface{}) string {
				f := obj.(*scrape.SeasonFetch)
				return fmt.Sprintf("season_%d.json", f.Year)
			}
			err := s.Store(fetch, fn)
			if err != nil {
				log.Printf("%v", err)
			}
		}
	}(seasonChan)
}

func (s *storageServiceImpl) Store(obj interface{}, fn fileNameFunc) error {
	s.mux.Lock()
	s.mux.Unlock()

	filename := fn(obj)

	fp := path.Join(s.folder, filename)
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fp, b, 0644)
}
