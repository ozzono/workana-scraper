package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

const (
	domain     = "https://www.workana.com"
	pagePrefix = "/jobs?category=it-programming&language=en"
	jobsPrefix = "/job/"
	tagPrefix  = "/jobs?skills="
)

type scrape struct {
	jobs  []job
	pages map[string]struct{}
}

type job struct {
	url   string
	title string
	tags  []string
}

func main() {
	c := colly.NewCollector()
	scrape := scrape{[]job{}, map[string]struct{}{}}

	// Find and visit all links
	stop := false
	c.OnHTML("a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasPrefix(href, domain+pagePrefix) && !stop {
			stop = true
			_, exist := scrape.pages[href]
			if !exist {
				scrape.pages[href] = struct{}{}
			}
			e.Request.Visit(href)
		}
	})

	c.OnRequest(func(r *colly.Request) {})

	c.Visit(domain + "/jobs?category=it-programming&language=en")
	scrape.GetJobs()
}

func (s *scrape) GetJobs() {
	for key := range s.pages {
		fmt.Println(getTags(key))
		break
	}
}

func getTags(url string) job {
	c := colly.NewCollector()
	j := job{}

	c.OnHTML("[class='project-item  js-project']", func(e *colly.HTMLElement) {
		el := strings.Split(e.ChildText(".project-header"), "\n")
		j.title = trimSpaceSufix(el[len(el)-1])
		j.url = domain + e.ChildAttr("a", "href")
		e.ForEach(".skills", func(i int, h *colly.HTMLElement) {
			fmt.Println(h.ChildText(".skill"))
		})
	})

	c.Visit(url)
	return j
}

func trimSpaceSufix(in string) string {
	for strings.HasPrefix(in, " ") {
		in = strings.TrimPrefix(in, " ")
	}
	return in
}
