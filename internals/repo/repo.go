package db

import (
	"database/sql"
	"log"
	"os"
	"strings"
	"workana-tags/internals/models"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

const (
	dbFileName = "sqlite.db"
)

type DB struct {
	SQL     *sql.DB
	Debug   bool
	Close   func()
	EraseDB func()
}

func NewDB(debug bool) (*DB, error) {
	if debug {
		log.Printf("removing old %s file", dbFileName)
	}
	os.Remove(dbFileName)

	if debug {
		log.Printf("creating %s file", dbFileName)
	}
	file, err := os.Create(dbFileName)
	if err != nil {
		return nil, errors.Wrap(err, "os.Create")
	}
	file.Close()
	if debug {
		log.Println(dbFileName, "created")
	}

	sqliteDatabase, err := sql.Open("sqlite3", "./"+dbFileName)
	if err != nil {
		return nil, errors.Wrap(err, "sql.Open")
	}

	db := &DB{
		SQL:   sqliteDatabase,
		Debug: debug,
		Close: func() {
			if debug {
				log.Printf("closing %s db", dbFileName)
			}
			sqliteDatabase.Close()
		},
		EraseDB: func() {
			if debug {
				log.Printf("erasing %s db file", dbFileName)
			}
			os.Remove("./" + dbFileName)
		},
	}

	if err := db.createTables(); err != nil {
		return nil, errors.Wrap(err, "db.createTables")
	}
	return db, nil
}

func (db *DB) createTables() error {
	if db.Debug {
		log.Println("create job table")
	}
	statement, err := db.SQL.Prepare(createJobsTable)
	if err != nil {
		return errors.Wrap(err, "db.SQL.Prepare")
	}
	defer statement.Close()

	statement.Exec()
	if db.Debug {
		log.Println("job table created")
	}

	if db.Debug {
		log.Println("create Tag table")
	}
	statement, err = db.SQL.Prepare(createTagsTable)
	if err != nil {
		return errors.Wrap(err, "db.SQL.Prepare")
	}
	statement.Exec()
	if db.Debug {
		log.Println("tag table created")
	}

	if db.Debug {
		log.Println("create jobtag table")
	}
	statement, err = db.SQL.Prepare(createJobTagTable)
	if err != nil {
		return errors.Wrap(err, "db.Prepare")
	}
	statement.Exec()
	if db.Debug {
		log.Println("jobtag table created")
	}

	return nil
}

func (db *DB) InsertJob(job models.Job) error {
	if db.Debug {
		log.Println("inserting job record")
	}
	statement, err := db.SQL.Prepare(insertJobQuery)
	if err != nil {
		return errors.Wrap(err, "db.SQL.Prepare(insertJobQuery)")
	}
	defer statement.Close()

	if err := statement.QueryRow(job.Title, job.Path).Scan(&job.ID); err != nil {
		return errors.Wrapf(err, "statement.QueryRow().Scan()")
	}

	if err := db.insertTags(job); err != nil {
		return errors.Wrap(err, "db.insertTags")
	}
	return nil
}

func (db *DB) insertTags(job models.Job) error {
	if db.Debug {
		log.Println("inserting tags record")
	}
	statement, err := db.SQL.Prepare(insertTagQuery)
	if err != nil {
		return errors.Wrap(err, "db.SQL.Prepare(insertTagQuery)")
	}
	defer statement.Close()

	tagIDs := []int{}
	for i := range job.Tags {
		var tagID int
		err = statement.QueryRow(strings.ToLower(job.Tags[i])).Scan(&tagID)
		if err != nil {
			return errors.Wrapf(err, "[%d]statement.Exec()", i)
		}
		tagIDs = append(tagIDs, tagID)
	}

	statement, err = db.SQL.Prepare(insertJobTagQuery)
	if err != nil {
		return errors.Wrap(err, "db.SQL.Prepare(insertTagQuery)")
	}

	for i := range tagIDs {
		_, err = statement.Exec(job.ID, tagIDs[i])
		if err != nil {
			return errors.Wrapf(err, "[%d]statement.Exec()", i)
		}
	}

	return nil
}

func (db *DB) GetJobs() ([]models.Job, error) {
	row, err := db.SQL.Query(getJobs)
	if err != nil {
		return nil, errors.Wrap(err, "db.Query")
	}
	defer row.Close()

	jobs := []models.Job{}
	for row.Next() {
		job := models.Job{}
		var tags string
		row.Scan(&job.ID, &job.Title, &job.Path, &tags)
		job.Tags = strings.Split(tags, ",")
		jobs = append(jobs, job)
	}
	return jobs, nil
}
