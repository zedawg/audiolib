package main

import (
	"embed"
	"flag"
	"log"
)

var (
	//go:embed embed
	embedFS embed.FS
	db      *Database
	version = "0.0.1"
	port    string
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
	log.Printf("application started, version %v", version)
	//
	go runHTTP()
	//
	for {
		c := make(chan bool)
		<-c
	}

}
