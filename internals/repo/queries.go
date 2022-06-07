package db

const (
	createJobsTable = `CREATE TABLE jobs (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"title" TEXT,
		"path" TEXT NOT NULL UNIQUE
	  );`

	createTagsTable = `create table tags(
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"label" TEXT NOT NULL UNIQUE
		);`

	createJobTagTable = `create table jobtags (
		"job_id"  INTEGER REFERENCES job(id) ON UPDATE CASCADE,
		"tag_id"  INTEGER REFERENCES tag(id) ON UPDATE CASCADE
	);`

	insertJobQuery    = `INSERT OR IGNORE INTO jobs(title, path) VALUES (?, ?) returning id;`
	insertTagQuery    = `INSERT OR IGNORE INTO tags(label) VALUES(?) returning id;`
	insertJobTagQuery = `INSERT OR IGNORE INTO jobtags(job_id,tag_id) VALUES(?,?);`

	getJobs = `SELECT 
		j.id,
		j.title,
		j.path,
		group_concat(t.label,",")
	FROM jobs j
	join jobtags jt on j.id=jt.job_id
	join tags t     on t.id=jt.tag_id
	group by j.id
	`
)
