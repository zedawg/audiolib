package main

import (
	"embed"
	"flag"
	"log"

	"github.com/zedawg/audiolib/db"
)

var (
	//go:embed assets
	AssetsFS embed.FS
	//go:embed public
	PublicFS embed.FS
	//go:embed templates
	TemplateFS embed.FS
	//go:embed db/schema.sql
	SQLFS embed.FS
)

var (
	Port    string
	Version = "0.0.1"
	dev     bool
)

func init() {
	//
	flag.StringVar(&Port, "port", ":8080", "client http port")
	flag.StringVar(&db.Name, "db", "audiolib.db", "sqlite database file path")
	flag.BoolVar(&dev, "dev", false, "development mode")
	flag.Parse()
}

func main() {
	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	b, _ := SQLFS.ReadFile("db/schema.sql")
	db.DB.Exec(string(b))
	//
	log.Println("version:", Version)
	log.Println("dev:", dev)
	log.SetFlags(log.Lshortfile)
	//
	go StartHTTP()
	//
	for {
		m := <-db.C
		log.Println(m)
	}

}
