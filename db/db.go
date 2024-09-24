package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

var (
	DB   *sql.DB
	Name string
)

func Connect(name string) (err error) {
	registerDriver()
	_, err = os.Stat(dbName)
	exists := err == nil
	// db.DB, err = sql.Open("sqlite3", db.Name)
	db.DB, err = sql.Open("sqlite3_with_hook_example", db.Name)
	if !exists && err == nil {
		log.Println("created and connected to", db.Name)
	} else if err != nil {
		log.Println(err)
	} else {
		log.Println("connected to", db.Name)
	}
	return err
}

func registerDriver() {
	sql.Register("sqlite3_with_hook_example", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			// sqlite3conn = append(sqlite3conn, conn)
			conn.RegisterUpdateHook(func(op int, db string, table string, rowid int64) {
				switch op {
				case sqlite3.SQLITE_INSERT:
					fmt.Println("Notified of insert on db", db, "table", table, "rowid", rowid)
				}
			})
			return nil
		},
	})
}

func (db *Database) Setup() error {
	s, _ := embedFS.ReadFile("embed/sql/schema.sql")
	_, err := db.Exec(string(s))
	return err
}

func (db *Database) Write(p []byte) (int, error) {
	_, err := db.Exec(`INSERT INTO logs (message) VALUES (?)`, string(bytes.TrimSpace(p)))
	return len(p), err
}
