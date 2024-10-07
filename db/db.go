package db

import (
	"database/sql"
	"embed"
	"log"

	"github.com/mattn/go-sqlite3"
	"github.com/zedawg/librarian/config"
)

var (
	//go:embed sql
	FS     embed.FS
	DB     *sql.DB
	Events = make(chan Event)
)

type Event struct {
	Op    int
	DB    string
	Table string
	RowID int64
}

func init() {
	sql.Register("sqlite3_driver", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			conn.RegisterUpdateHook(updateHookFunc)
			return nil
		},
	})

	// note: must import package config before db for this to make sense
	var err error
	DB, err = sql.Open("sqlite3_driver", config.DatabasePath())
	if err != nil {
		log.Fatal(err)
	}

	if err := executeScript("sql/schema.sql"); err != nil {
		log.Println(err)
	}
}

func Close() error {
	return DB.Close()
}

func updateHookFunc(op int, db string, table string, rowid int64) {
	Events <- Event{op, db, table, rowid}
}

func executeScript(name string) error {
	b, _ := FS.ReadFile(name)
	_, err := DB.Exec(string(b))
	return err
}
