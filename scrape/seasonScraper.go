package scrape

import (
	"fmt"
	"github.com/sheitm/ofever/types"
	"golang.org/x/net/html"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// StartSeason scrapes all events that can be found via the given url which is assumed to be to a page with references
// to all events in a table.
func StartSeason(url string, year int, resultChan chan<- *types.SeasonFetch) {
	client := &http.Client{}
	startSeason(url, year, resultChan, client)
}

func startSeason(url string, year int, resultChan chan<- *types.SeasonFetch, client *http.Client) {
	scraper := &seasonScraper{
		client: client,
		year:   year,
	}
	go scraper.scrape(url, resultChan)
}

type seasonScraper struct {
	client *http.Client
	year   int
}

func (s *seasonScraper) scrape(url string, resultChan chan<- *types.SeasonFetch) {
	fetch := &types.SeasonFetch{
		URL:  url,
		Year: s.year,
	}

	rows, err := s.getEventRows(url)
	if err != nil {
		fetch.Error = err.Error()
		resultChan <- fetch
		return
	}

	internalResultChan := make(chan *types.ScrapeResult)
	for _, row := range rows {
		startEventScrape(row, internalResultChan, s.client)
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(rows))

	var results []*types.ScrapeResult

	go func(in <-chan *types.ScrapeResult, wg *sync.WaitGroup){
		for {
			r := <- in
			results = append(results, r)
			wg.Done()
		}
	}(internalResultChan, wg)

	wg.Wait()
	fetch.Results = results

	// TODO: Figure out how to set series.
	fetch.Series = "geoform"

	resultChan <- fetch
}

func (s *seasonScraper) getEventRows(url string) ([]*tableRow, error) {
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var currentRow *tableRow
	var headers *tableRow
	var rows []*tableRow

	var f func(*html.Node)
	f = func(n *html.Node){
		if n.Type == html.ElementNode && n.Data == "tr" {
			currentRow = newTableRow(url, s.year)
			if headers == nil {
				headers = currentRow
			} else {
				rows = append(rows, currentRow)
			}

		}
		if n.Type == html.ElementNode && n.Data == "td" {
			currentRow.add(n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	var validRows []*tableRow
	for _, row := range rows {
		if row.valid() {
			validRows = append(validRows, row)
		}
	}

	return validRows, nil
}

func newTableRow(baseURL string, year int) *tableRow {
	return &tableRow{
		baseURL: baseURL,
		values:  map[int]cellValue{},
		year:    year,
	}
}

type tableRow struct {
	baseURL string
	values  map[int]cellValue
	year    int
}

func (r *tableRow) valid() bool {
	if len(r.values) < 9 {
		return false
	}
	return r.rawEventURL() != ""
}

func (r *tableRow) number() int {
	if c, ok := r.values[0]; ok {
		nr, err := strconv.Atoi(c.text)
		if err != nil {
			return 0
		}
		return nr
	}
	return 0
}

func (r *tableRow) rawEventURL() string {
	index := len(r.values) - 3
	if c, ok := r.values[index]; ok {
		return c.value
	}
	return ""
}

func (r *tableRow) eventURL() string {
	b := r.baseURL
	if strings.Contains(b, ".html") {
		arr := strings.Split(b, "/")
		end := (len(b) - len(arr[len(arr)-1])) -1
		b = b[0:end]
	}
	return b + "/" + r.rawEventURL()
}

func (r *tableRow) urlInvite() string {
	if c, ok := r.values[3]; ok {
		return c.value
	}
	return ""
}

func (r *tableRow) place() string {
	if c, ok := r.values[3]; ok {
		return c.text
	}
	return ""
}

func (r *tableRow) urlLiveLox() string {
	index := len(r.values) - 2
	if c, ok := r.values[index]; ok {
		return c.value
	}
	return ""
}

func (r *tableRow) weekDay() string {
	if c, ok := r.values[2]; ok {
		return c.value
	}
	return ""
}

func (r *tableRow) organizer() string {
	if c, ok := r.values[5]; ok {
		return c.value
	}
	return ""
}

func (r *tableRow) date() time.Time {
	if c, ok := r.values[1]; ok {
		arr := strings.Split(c.value, ".")
		if len(arr) != 2 {
			return time.Time{}
		}
		d, err := strconv.Atoi(arr[0])
		if err != nil {
			return time.Time{}
		}
		m, err := strconv.Atoi(arr[1])
		if err != nil {
			return time.Time{}
		}
		return time.Date(r.year, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	}
	return time.Time{}
}


func (r *tableRow) add(n *html.Node) {
	c := n.FirstChild
	var cell cellValue
	if c != nil {
		switch c.Data {
		case "span":
			v := c.FirstChild.Data
			cell = cellValue{
				typ:   "span",
				text:  v,
				value: v,
			}
		case "a":
			v := c.FirstChild.Data
			href := ""
			for _, attribute := range c.Attr {
				if attribute.Key == "href" {
					href = attribute.Val
					break
				}
			}
			cell = cellValue{
				typ:   "a",
				text:  v,
				value: href,
			}
		case "strong":
			v := c.FirstChild.Data
			cell = cellValue{
				typ:   "strong",
				text:  v,
				value: v,
			}
		default:
			v := c.Data
			cell = cellValue{
				typ:   "text",
				text:  v,
				value: v,
			}
		}
	}

	r.values[len(r.values)] = cell
}

type cellValue struct {
	typ   string
	text  string
	value string
}