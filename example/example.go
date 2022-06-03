package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

const (
	pagesDomain = "https://www.workana.com/jobs?category=it-programming"
	jobsDomain  = "https://www.workana.com/job/"
)

func main() {
	c := colly.NewCollector(
	// colly.AllowedDomains("workana.com/"),
	)
	// c.Async = true
	// Find and visit all links
	urlList := []string{}
	c.OnHTML("a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasPrefix(href, pagesDomain) || strings.HasPrefix(href, jobsDomain) {
			e.Request.Visit(e.Attr("href"))
			urlList = append(urlList, href)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.workana.com/jobs?category=it-programming&language=en")
}
