package main

import (
	"embed"
	"flag"
	"html/template"
	"io/fs"
	"log"
)

var (
	//go:embed embed
	embedFS     embed.FS
	staticFS    fs.FS
	webTemplate *template.Template
	db          *Database
	taskCh      = make(chan Task)
	version     = "0.0.1"
	port        string
)

func init() {
	//
	db = &Database{}
	//
	flag.StringVar(&port, "port", ":8080", "http port")
	flag.StringVar(&port, "p", ":8080", "short for [port]")
	flag.StringVar(&db.Name, "dbname", "audiolib.db", "path to sqlite database file")
	flag.StringVar(&db.Name, "d", "audiolib.db", "short for [dbname]")
	flag.Parse()
}

func main() {
	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Setup()
	//
	log.SetFlags(log.Lshortfile)
	log.SetFlags(0)
	log.SetOutput(db)
	log.Println("**application started**")
	log.Printf("**version %v**", version)
	//
	go runHTTP()
	//
	for {
		t := <-taskCh
		switch t.Name {
		case "scan":
		case "match":
		case "shutdown":
			log.Println("application shutdown")
			return
		default:

		}
	}
}
