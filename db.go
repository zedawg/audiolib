package main

import (
	"bytes"
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sql.DB
	Name string
}

func (db *Database) Connect() (err error) {
	_, err = os.Stat(db.Name)
	exists := err == nil
	db.DB, err = sql.Open("sqlite3", db.Name)
	if !exists && err == nil {
		log.Println("created and connected to", db.Name)
	} else if err != nil {
		log.Println(err)
	} else {
		log.Println("connected to", db.Name)
	}
	return err
}

func (db *Database) Setup() error {
	_, err := db.Exec(db.schema())
	return err
}

func (db *Database) schema() string {
	s, _ := embedFS.ReadFile("embed/sql/schema.sql")
	return string(s)
}

func (db *Database) InsertTask(name, status, params string) (int, error) {
	var id int
	err := db.QueryRow(`INSERT INTO tasks (name, status, params) values (?,?,?) RETURNING id`, name, status, params).Scan(&id)
	return id, err
}

func (db *Database) Write(p []byte) (int, error) {
	_, err := db.Exec(`INSERT INTO logs (message) VALUES (?)`, string(bytes.TrimSpace(p)))
	return len(p), err
}

func (db *Database) GetLibrariesCount() (count int, err error) {
	err = db.QueryRow(`SELECT COUNT(*) FROM libraries`).Scan(&count)
	return count, err
}

func (db *Database) GetLibraries() (libraryEntries []*LibraryEntry, err error) {
	rows, err := db.Query(`SELECT id, name, paths, created FROM libraries ORDER BY created DESC`)
	if err != nil {
		return
	}
	for rows.Next() {
		l := LibraryEntry{}
		if err = rows.Scan(&l.ID, &l.Name, &l.Paths, &l.Created); err != nil {
			return
		}
		libraryEntries = append(libraryEntries, &l)
	}
	return
}

func (db *Database) GetLogsCount() (count int, err error) {
	err = db.QueryRow(`SELECT COUNT(*) FROM logs`).Scan(&count)
	return count, err
}

func (db *Database) GetLogs(limit, offset int) (logEntries []*LogEntry, err error) {
	rows, err := db.Query(`SELECT id, message, created FROM logs ORDER BY created DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return
	}
	for rows.Next() {
		l := LogEntry{}
		if err = rows.Scan(&l.ID, &l.Message, &l.Created); err != nil {
			return
		}
		logEntries = append(logEntries, &l)
	}
	return
}

func (db *Database) GetTasksCount() (count int, err error) {
	err = db.QueryRow(`SELECT COUNT(*) FROM tasks`).Scan(&count)
	return count, err
}

func (db *Database) GetTasks(limit, offset int) (taskEntries []*TaskEntry, err error) {
	rows, err := db.Query(`SELECT id, name, priority, status, params, result, created FROM tasks LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return
	}
	for rows.Next() {
		t := TaskEntry{}
		if err = rows.Scan(&t.ID, &t.Name, &t.Priority, &t.Status, &t.Params, &t.Result); err != nil {
			return
		}
		taskEntries = append(taskEntries, &t)
	}
	return
}
