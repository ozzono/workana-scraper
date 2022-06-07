package scraper

import (
	"net/http"
	"strings"
	"workana-tags/internals/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/pkg/errors"
)

const (
	domain     = "https://www.workana.com"
	pagePrefix = "/jobs?category=it-programming&language=en"
	jobsPrefix = "/job/"
	tagPrefix  = "/jobs?skills="
)

func GetJobPages() map[string]struct{} {
	c := colly.NewCollector()
	pages := map[string]struct{}{}

	// Find and visit all links
	c.OnHTML("a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasPrefix(href, domain+pagePrefix) {
			_, exist := pages[href]
			if !exist {
				pages[href] = struct{}{}
			}
			e.Request.Visit(href)
		}
	})

	c.OnRequest(func(r *colly.Request) {})

	c.Visit(domain + "/jobs?category=it-programming&language=en")
	return pages
}

type client struct {
	scrape *goquery.Document
	close  func()
	debug  bool
}

func newClient(url string, debug bool) (*client, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "http.Get")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "goquery.NewDocumentFromReader")
	}

	return &client{
		scrape: doc,
		close: func() {
			res.Body.Close()
		},
		debug: debug,
	}, nil
}

func GetWorkanaJobs(url string) ([]*models.Job, error) {

	client, err := newClient(domain+url, true)
	if err != nil {
		return nil, errors.Wrap(err, "newClient")
	}
	jobs := []*models.Job{}

	client.scrape.Find(".project-item").Each(func(i int, s *goquery.Selection) {
		j := models.Job{}
		j.Title = s.Find(".project-header a").Text()
		href, found := s.Find(".project-header a").Attr("href")
		if found {
			j.Path = href
		}

		jobs = append(jobs, &j)
	})

	return jobs, nil
}

func GetJobTags(url string) ([]string, error) {
	client, err := newClient(url, true)
	if err != nil {
		return nil, errors.Wrap(err, "newClient")
	}

	tags := []string{}

	client.scrape.Find(".skill").Each(func(i int, s *goquery.Selection) {
		tags = append(tags, s.Text())
	})

	return tags, nil
}
