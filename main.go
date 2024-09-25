package main

import (
	"embed"
	"flag"
	"log"
	"os"
	"path"

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
	//go:embed db/schema.sql
	SQLFS embed.FS
)

var (
	VERSION = "0.0.1"
	DEV     bool
)

func init() {
	flag.BoolVar(&DEV, "dev", false, "development mode")
	flag.StringVar(&config.Name, "config", "audiolib.config", "config json")
	flag.Parse()

	wd, _ := os.Getwd()
	if !path.IsAbs(config.Name) {
		config.Name = path.Join(wd, config.Name)
	}
	if _, err := os.Lstat(config.Name); err != nil {
		if err = config.UseDefault(); err != nil {
			log.Fatal(err)
		}
	}
	if err := config.Parse(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	b, _ := SQLFS.ReadFile("db/schema.sql")
	db.DB.Exec(string(b))
	//
	log.SetFlags(0)
	log.Println("config:", config.Name)
	log.Println("database:", config.C.Database)
	log.Println("app port:", config.C.Port)
	log.Println("version:", VERSION)
	log.Println("dev:", DEV)
	if DEV {
		log.SetFlags(log.Lshortfile)
	}
	//
	go StartHTTP()
	//
	for {
		m := <-db.C
		log.Println(m)
	}

}
