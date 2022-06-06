package models

import "log"

type Job struct {
	ID    int
	Title string
	Path  string
	Tags  []string
}

func (j Job) Log() {
	log.Printf("j.ID ----- %v", j.ID)
	log.Printf("j.Title -- %v", j.Title)
	log.Printf("j.Path --- %v", j.Path)
	log.Printf("j.Tags --- %v", j.Tags)
}
