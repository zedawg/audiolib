package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

var (
	Dev  = false // default
	port = 8000  // default
	data = ""
	name = "db.sqlite"
)

var (
	ErrBadConfig = errors.New("bad config: execute --help")
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	flag.StringVar(&data, "data", data, "application data")
	flag.StringVar(&data, "d", data, "short for --data")
	flag.IntVar(&port, "port", port, "port")
	flag.IntVar(&port, "p", port, "short for --port")
	flag.BoolVar(&Dev, "dev", Dev, "development mode")
	flag.BoolVar(&Dev, "D", Dev, "short for --dev")
	flag.Parse()

	if len(data) == 0 {
		log.Fatal(ErrBadConfig)
	}

	if !path.IsAbs(data) {
		wd, _ := os.Getwd()
		data = path.Join(wd, data)
	}

	if _, err := os.Lstat(data); err != nil {
		os.MkdirAll(data, 0700)
	}

	log.Println("data path:", data)
	log.Println("port:", port)
	log.Println("dev mode:", Dev)
}

func DatabasePath() string {
	return path.Join(data, name)
}

func Port() string {
	return fmt.Sprintf(":%v", port)
}
