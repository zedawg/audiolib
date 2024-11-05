package config

import (
	"flag"
	"log"
	"os"
	"path"
	"strings"
)

var (
	home, _ = os.UserHomeDir()
	Dev     = false                          // default, flag
	Port    = "8000"                         // default, flag
	dir     = path.Join(home, "goaudiobook") // default, flag
	DBName  = "db.sqlite"                    // ${dir}/db.sqlite < parseFlags()
)

func init() {
	parseFlags()
	validateFlags()
}

func parseFlags() {
	flag.StringVar(&dir, "dir", dir, "data directory")
	flag.StringVar(&Port, "port", Port, "port")
	flag.BoolVar(&Dev, "dev", Dev, "development mode")

	testMode := false
	// detect if run from `go test`
	for _, f := range os.Args[1:] {
		if strings.Index(f, "test.") > -1 {
			testMode = true
		}
	}
	if !testMode {
		flag.Parse()
	}
	DBName = path.Join(dir, DBName)
}

func validateFlags() {
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
