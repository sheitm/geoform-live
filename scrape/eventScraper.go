package scrape

import (
	"github.com/sheitm/ofever/contracts"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

type eventScraper struct {
	client *http.Client
}

func (s *eventScraper) Scrape(url string) (*contracts.Event, error) {
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
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