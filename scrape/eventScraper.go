package scrape

import (
	"fmt"
	"github.com/sheitm/ofever/types"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

func startEventScrape(row *tableRow,  resultChan chan<- *types.ScrapeResult, client *http.Client) {
	go func(row *tableRow, resultChan chan<- *types.ScrapeResult, client *http.Client) {
		url := row.eventURL()
		res := &types.ScrapeResult{URL: url}
		scraper := &eventScraper{
			client: &http.Client{},
			row:    row,
		}
		event, err := scraper.Scrape(url)
		res.Event = event
		if err != nil {
			res.Error = err.Error()
		}
		resultChan <- res
	}(row, resultChan, client)
}

type eventScraper struct {
	client *http.Client
	row    *tableRow
}

func (s *eventScraper) Scrape(url string) (*types.Event, error) {
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	defer resp.Body.Close()


	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	event := &types.Event{
		URL:        url,
		Number:     s.row.number(),
		URLInvite:  s.row.urlInvite(),
		URLLiveLox: s.row.urlLiveLox(),
		WeekDay:    s.row.weekDay(),
		Date:       s.row.date(),
		Place:      s.row.place(),
		Organizer:  s.row.organizer(),
	}

	var nextCourseName string
	newCourseDetected := false

	resultsParser := &eventTableParser{}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			event.Name = n.FirstChild.Data
		}
		if n.Type == html.ElementNode && n.Data == "h3" && n.FirstChild.Data == "a"{
			anchor := n.FirstChild
			if len(anchor.Attr) == 1 && anchor.Attr[0].Key == "name" && strings.Contains(anchor.Attr[0].Val, "Res") {
				nextCourseName = anchor.FirstChild.Data
				newCourseDetected = true
			}
		}
		if n.Type == html.ElementNode && n.Data == "pre" {
			if n.FirstChild.Data == "b" && newCourseDetected {
				newCourseDetected = false
				columns := n.FirstChild.FirstChild.Data
				resString := n.FirstChild.NextSibling.Data
				results, err := resultsParser.parse(columns, resString)
				parserErrorText := ""
				if err != nil {
					parserErrorText = err.Error()
				}
				course := &types.Course{
					Name:       nextCourseName,
					Info:       "",
					Results:    results,
					ParseError: parserErrorText,
				}
				event.Courses = append(event.Courses, course)
			} else {
				event.Info = n.FirstChild.Data
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return event, nil
}