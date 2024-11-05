package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/zedawg/goaudiobook/config"
	"github.com/zedawg/goaudiobook/server"
	"github.com/zedawg/goaudiobook/sql"
)

func seedDB() {
	// sql.DeleteSources()
	sql.InsertSource("/users/admin/downloads/test_audiobooks")
	sql.InsertSource("/users/admin/downloads/audiobooks")
}

func main() {
	sql.Open(fmt.Sprintf("%v?_foreign_keys=on", config.DBName))

	if config.Dev {
		log.SetFlags(log.Ltime | log.Lshortfile)
		seedDB()
	}

	log.Println("database", strings.ToLower(config.DBName))
	log.Printf("browse audiobooks at http://localhost:%v", config.Port)

	errCh := make(chan error)

	go func() {
		errCh <- server.Listen()
	}()

	if err := buildLibrary(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\r")
	log.Println("library ready")

	log.Fatal(<-errCh)
}
