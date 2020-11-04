package storage

import (
	"encoding/json"
	"fmt"
	"github.com/sheitm/ofever/scrape"
	"io/ioutil"
	"path"
)

type service struct {
	folder string
	year   int
}

func (s *service) Year() int {
	return s.year
}

func (s *service) Store(fetch *scrape.SeasonFetch) error {
	fp := path.Join(s.folder, fmt.Sprintf("season_%d.json", s.year))
	b, err := json.Marshal(fetch)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fp, b, 0644)
}
