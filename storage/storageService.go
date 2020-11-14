package storage

import (
	"encoding/json"
	"fmt"
	"github.com/sheitm/ofever/scrape"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"sync"
)

type storageService interface {
	Start(element seasonSyncElement)
	Store(obj interface{}, fn fileNameFunc) error
	Fetch(obj interface{}, fn fileNameFunc) error
	ReadFolder(folder, filePattern string) ([]string, error)
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

func (s *storageServiceImpl) ReadFolder(folder, filePattern string) ([]string, error) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	r := regexp.MustCompile(filePattern)
	var res []string
	for _, file := range files {
		n := file.Name()
		if !r.Match([]byte(n)) {
			continue
		}
		fn := path.Join(folder, n)
		content, err := ioutil.ReadFile(fn)
		if err != nil {
			return nil, err
		}
		res = append(res, string(content))
	}
	return res, nil
}

func (s *storageServiceImpl) Start(element seasonSyncElement) {
	go func(sc <-chan *scrape.SeasonFetch, dc chan<- struct{}) {
		for {
			fetch := <- sc
			fn := func(obj interface{}) string {
				f := obj.(*scrape.SeasonFetch)
				return fmt.Sprintf("season_%d.json", f.Year)
			}
			err := s.Store(fetch, fn)
			if err != nil {
				log.Printf("%v", err)
			}
			dc <- struct{}{}
		}
	}(element.seasonChan, element.doneChan)
}

func (s *storageServiceImpl) Fetch(obj interface{}, fn fileNameFunc) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	filename := fn(obj)
	fp := path.Join(s.folder, filename)
	jsonFile, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(byteValue, obj)
}

func (s *storageServiceImpl) Store(obj interface{}, fn fileNameFunc) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	filename := fn(obj)

	fp := path.Join(s.folder, filename)
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fp, b, 0644)
}
