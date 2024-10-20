package sql

import (
	"database/sql"
	"embed"

	_ "github.com/mattn/go-sqlite3"
)

var (
	//go:embed schema.sql
	Files embed.FS
	DB    *sql.DB
)

func Open(p string) (err error) {
	DB, err = sql.Open("sqlite3", p)
	if err != nil {
		return err
	}

	b, _ := Files.ReadFile("schema.sql")
	DB.Exec(string(b))

	return DB.Ping()
}
