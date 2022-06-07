package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"sync"
	"time"
	"workana-tags/internals/models"
	db "workana-tags/internals/repo"
	"workana-tags/scraper"

	"github.com/pkg/errors"
)

var (
	mu sync.Mutex
	wg sync.WaitGroup
)

type jobs []*models.Job

func main() {
	pages := scraper.GetJobPages()

	var (
		mu       = sync.Mutex{}
		jobs     = &jobs{}
		errPages = []string{}
	)

	fmt.Println("pages", len(pages))

	for key := range pages {
		wg.Add(1)
		go func(key string) {
			try, err := getJobs(key, 0, jobs)
			if err != nil {
				fmt.Println(err)
				mu.Lock()
				errPages = append(errPages, key)
				mu.Unlock()
			}
			if try > 0 {
				fmt.Printf("%s got %d tries\n", key, try)
			}
		}(key)
	}

	wg.Wait()
	if len(errPages) > 0 {
		sort.Strings(errPages)
		fmt.Println("err pages", len(errPages))
	}
	for i := range errPages {
		fmt.Println(errPages[i])
	}

	fmt.Println("len", len(*jobs))

	tagsErr := []error{}

	db, err := db.NewDB(false)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	for _, job := range *jobs {
		go func(job *models.Job) {
			tags, err := scraper.GetJobTags(job.Path)
			if err != nil {
				mu.Lock()
				tagsErr = append(tagsErr, err)
				mu.Unlock()
				return
			}
			job.Tags = tags
			if err := db.InsertJob(*job); err != nil {
				log.Printf("failed to insert job %v\n", job)
			}
		}(job)
	}
	if len(tagsErr) > 0 {
		fmt.Println("len tag err", len(tagsErr))
		for i := range tagsErr {
			fmt.Println(i, "tag err", tagsErr[i])
		}
	}
}

func getJobs(url string, try int, jobs *jobs) (int, error) {
	defer wg.Done()
	if try == 10 {
		return try, errors.Errorf("overtried %s", url)
	}
	j, err := scraper.GetWorkanaJobs(url)
	if err != nil {
		if try > 0 {
			time.Sleep(time.Duration(rand.Intn(time.Now().Nanosecond()+try)%100+(try*100)) * time.Millisecond)
		}
		wg.Add(1)
		try++
		return getJobs(url, try, jobs)
	}
	mu.Lock()
	*jobs = append(*jobs, j...)
	mu.Unlock()
	return try, nil
}
