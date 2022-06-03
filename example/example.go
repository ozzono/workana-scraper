package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

const (
	pagePrefix = "https://www.workana.com/jobs?category=it-programming&language=en"
	jobsPrefix = "https://www.workana.com/job/"
	tagPrefix  = "https://www.workana.com/jobs?skills="
)

type scrape struct {
	jobs  []job
	pages map[string]struct{}
}

type job struct {
	url  string
	tags []string
}

func main() {
	c := colly.NewCollector()
	scrape := scrape{[]job{}, map[string]struct{}{}}

	// Find and visit all links
	stop := false
	c.OnHTML("a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasPrefix(href, pagePrefix) && !stop {
			stop = true
			_, exist := scrape.pages[href]
			if !exist {
				scrape.pages[href] = struct{}{}
			}
			e.Request.Visit(href)
		}
	})

	c.OnRequest(func(r *colly.Request) {})

	c.Visit("https://www.workana.com/jobs?category=it-programming&language=en")
	scrape.GetJobs()
}

func (s *scrape) GetJobs() {
	for key := range s.pages {
		fmt.Println(getTags(key))
		break
	}
}

func getTags(url string) []string {
	c := colly.NewCollector()
	fmt.Println(url)
	out := []string{}
	c.OnHTML("[class='project-item  js-project']", func(e *colly.HTMLElement) {
		el := strings.Split(e.ChildText(".project-header"), "\n")
		fmt.Printf("%#v\n", trimSpaceSufix(el[len(el)-1]))
	})
	c.Visit(url)
	return out
}

func trimSpaceSufix(in string) string {
	for strings.HasPrefix(in, " ") {
		in = strings.TrimPrefix(in, " ")
	}
	return in
}
