package main

import (
	"embed"
	"flag"
	"log"

	"github.com/mattn/go-sqlite3"
	"github.com/zedawg/audiolib/config"
	"github.com/zedawg/audiolib/db"
)

var (
	//go:embed assets
	AssetsFS embed.FS
	//go:embed public
	PublicFS embed.FS
	//go:embed templates
	TemplateFS embed.FS
)

var (
	CONFIG_FILE string
	DEV         bool
)

func init() {
	flag.StringVar(&CONFIG_FILE, "config", "audiolib.config", "config json")
	flag.BoolVar(&DEV, "dev", false, "development mode")
	flag.Parse()
	config.Parse(CONFIG_FILE)
	db.Init()
}

func main() {
	log.SetFlags(0)
	log.Printf("dev=%v", DEV)
	log.Println(config.C)
	if DEV {
		log.SetFlags(log.Lshortfile)
	}
	//
	defer db.Close()
	go StartHTTP()
	//
	for {
		m := <-db.M
		switch m.Table {
		case "tasks":
			switch m.Op {
			case sqlite3.SQLITE_INSERT:
			case sqlite3.SQLITE_DELETE:
			case sqlite3.SQLITE_UPDATE:
			default:
			}
		case "books":
			switch m.Op {
			case sqlite3.SQLITE_INSERT:
			case sqlite3.SQLITE_DELETE:
			case sqlite3.SQLITE_UPDATE:
			default:
			}

		case "files":
			switch m.Op {
			case sqlite3.SQLITE_INSERT:
			case sqlite3.SQLITE_DELETE:
			case sqlite3.SQLITE_UPDATE:
			default:
			}

		case "images":
			switch m.Op {
			case sqlite3.SQLITE_INSERT:
			case sqlite3.SQLITE_DELETE:
			case sqlite3.SQLITE_UPDATE:
			default:
			}

		default:
		}

		log.Println(m)
	}

}
