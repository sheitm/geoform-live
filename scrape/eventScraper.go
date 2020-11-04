package scrape

import (
	"fmt"
	"github.com/sheitm/ofever/contracts"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

func startEventScrape(url string, resultChan chan<- *Result, client *http.Client) {
	go func(url string, resultChan chan<- *Result, client *http.Client) {
		res := &Result{URL: url}
		scraper := &eventScraper{client: &http.Client{}}
		event, err := scraper.Scrape(url)
		res.Event = event
		res.Error = err
		resultChan <- res
	}(url, resultChan, client)
}

type eventScraper struct {
	client *http.Client
}

func (s *eventScraper) Scrape(url string) (*contracts.Event, error) {
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

	event := &contracts.Event{}
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
				if err != nil {
					x := 22
					_ = x
				}
				course := &contracts.Course{
					Name:    nextCourseName,
					Info:    "",
					Results: results,
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