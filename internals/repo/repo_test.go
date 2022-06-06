package db

import (
	"testing"
	"workana-tags/internals/models"

	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	db, err := NewDB(true)
	assert.NoError(t, err, "NewDB()")
	defer func() {
		db.Close()
		db.EraseDB()
	}()

	testJob := models.Job{
		Title: "test job",
		Path:  "test/path",
		Tags:  []string{"tag1", "tag2", "tag3"},
	}

	assert.NoError(t, db.InsertJob(testJob), "db.InsertJob(testJob)")
	jobs, err := db.GetJobs()
	assert.NoError(t, err, "db.GetJobs()")
	assert.Conditionf(t, func() (success bool) {
		return len(jobs) == 1
	}, "expeted: 1 job; found: %d jobs", len(jobs))
	for i := range jobs {
		jobs[i].Log()
	}
}
