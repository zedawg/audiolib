package config

import (
	"flag"
	"log"
	"os"
	"path"
	"strings"
)

var (
	Dev     = false  // default
	Port    = "8000" // default
	dbname  = "db.sqlite"
	home, _ = os.UserHomeDir()
	dir     = path.Join(home, "librarian")
)

func init() {
	log.SetFlags(log.Ltime)
	defineFlags()
	parseFlags()
	validate()
}

func defineFlags() {
	flag.StringVar(&dir, "dir", dir, "data directory")
	flag.StringVar(&Port, "port", Port, "port")
	flag.BoolVar(&Dev, "dev", Dev, "development mode")
}

func parseFlags() {
	testMode := false
	for _, f := range os.Args[1:] {
		if strings.Index(f, "test.") > -1 {
			testMode = true
		}
	}
	if !testMode {
		flag.Parse()
	}
}

func validate() {
	if !path.IsAbs(dir) {
		wd, _ := os.Getwd()
		dir = path.Join(wd, dir)
	}
	os.MkdirAll(dir, 0700)
	if _, err := os.Lstat(dir); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func DatabasePath() string {
	return path.Join(dir, dbname)
}
