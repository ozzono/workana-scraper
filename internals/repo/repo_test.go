package db

import (
	"strings"
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
		Tags:  []string{"TAG1", "tAg2", "tag3"},
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

	expected := models.Job{
		Title: testJob.Title,
		Path:  testJob.Path,
		Tags:  strings.Split(strings.ToLower(strings.Join(testJob.Tags, ",")), ","),
	}
	returned := models.Job{
		Title: jobs[0].Title,
		Path:  jobs[0].Path,
		Tags:  jobs[0].Tags,
	}

	assert.Equalf(t, expected, returned, "exepecting: %v, found: %v", expected, returned)
}
