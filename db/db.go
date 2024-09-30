package db

import (
	"database/sql"
	"embed"
	"log"

	"github.com/mattn/go-sqlite3"
	"github.com/zedawg/audiolib/config"
)

var (
	//go:embed schema.sql
	SQLFS embed.FS
	DB    *sql.DB
	M     chan Message
)

type Message struct {
	Op    int
	DB    string
	Table string
	RowID int64
}

func Init() {
	M = make(chan Message)

	sql.Register("sqlite3_driver", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			conn.RegisterUpdateHook(func(op int, db string, table string, rowid int64) {
				M <- Message{Op: op, DB: db, Table: table, RowID: rowid}
			})
			return nil
		},
	})
	var err error
	DB, err = sql.Open("sqlite3_driver", config.C.Database)
	if err != nil {
		log.Fatal(err)
	}

	b, _ := SQLFS.ReadFile("schema.sql")
	DB.Exec(string(b))
}

func Close() error {
	return DB.Close()
}
