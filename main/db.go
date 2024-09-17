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
