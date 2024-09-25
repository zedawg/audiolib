package db

import (
	"database/sql"

	"github.com/mattn/go-sqlite3"
	"github.com/zedawg/audiolib/config"
)

var (
	DB *sql.DB
	C  = make(chan Message)
)

type Message struct {
	Op    int
	DB    string
	Table string
	RowID int64
}

func Connect() (err error) {
	registerDriver()
	DB, err = sql.Open("sqlite3_driver", config.C.Database)
	return err
}

func registerDriver() {
	sql.Register("sqlite3_driver", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			conn.RegisterUpdateHook(func(op int, db string, table string, rowid int64) {
				C <- Message{Op: op, DB: db, Table: table, RowID: rowid}
			})
			return nil
		},
	})
}

func Close() error {
	return DB.Close()
}
