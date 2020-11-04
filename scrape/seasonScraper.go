package scrape

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"sync"
)

// StartSeason scrapes all events that can be found via the given url which is assumed to be to a page with references
// to all events in a table.
func StartSeason(url string, resultChan chan<- *SeasonFetch) {
	client := &http.Client{}
	startSeason(url, resultChan, client)
}

func startSeason(url string, resultChan chan<- *SeasonFetch, client *http.Client) {
	scraper := &seasonScraper{client: client}
	go scraper.scrape(url, resultChan)
}

type seasonScraper struct {
	client *http.Client
}

func (s *seasonScraper) scrape(url string, resultChan chan<- *SeasonFetch) {
	fetch := &SeasonFetch{URL: url}
	eventURLs, err := s.getEventsURLs(url)
	if err != nil {
		fetch.Error = err
		resultChan <- fetch
		return
	}

	internalResultChan := make(chan *Result)
	for _, eventURL := range eventURLs {
		startEventScrape(eventURL, internalResultChan, s.client)
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(eventURLs))

	var results []*Result

	go func(in <-chan *Result, wg *sync.WaitGroup){
		for {
			r := <- in
			results = append(results, r)
			wg.Done()
		}
	}(internalResultChan, wg)

	wg.Wait()
	fetch.Results = results
	resultChan <- fetch
}

func (s *seasonScraper) getEventsURLs(url string) ([]string, error) {
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
	var rows []*tableRow

	var f func(*html.Node)
	f = func(n *html.Node){
		if n.Type == html.ElementNode && n.Data == "tr" {
			currentRow = newTableRow()
			rows = append(rows, currentRow)
		}
		if n.Type == html.ElementNode && n.Data == "td" {
			currentRow.add(n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	var links []string
	for _, row := range rows {
		eventURL := row.eventURL()
		if eventURL != "" {
			address := url + eventURL
			links = append(links, address)
		}
	}

	return links, nil
}

func newTableRow() *tableRow {
	return &tableRow{values: map[int]cellValue{}}
}

type tableRow struct {
	values map[int]cellValue
}

func (r *tableRow) eventURL() string {
	for _, value := range r.values {
		if value.typ == "a" && value.text == "Resultater" {
			return value.value
		}
	}

	return ""
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