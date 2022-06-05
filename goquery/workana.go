package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type job struct {
	title string
	path  string
}

func Workana() {
	// Request the HTML page.
	res, err := http.Get("https://www.workana.com/jobs?category=it-programming&language=en&page=2")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	repeat := true
	doc.Find(".project-item").Each(func(i int, s *goquery.Selection) {
		if repeat {
			j := job{}
			j.title = s.Find(".project-header a").Text()
			href, found := s.Find(".project-header a").Attr("href")
			if found {
				j.path = href
			}
			j.Log()

			s.Find(".project-body").Text()
			repeat = false
		}
	})
}

func (j job) Log() {
	log.Printf("j.title -- %v", j.title)
	log.Printf("j.path --- %v", j.path)
}

func main() {
	Workana()
}
