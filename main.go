package main

import (
	"log"
	"time"

	"github.com/zedawg/librarian/db"
	"github.com/zedawg/librarian/http"
)

var Errors = make(chan error)

var (
	RUN_HTTP_SRV = 0
	RUN_TASK_MGR = 0
)

func main() {
	defer db.Close()
	go func() { Errors <- runHTTPSrv() }()
	go func() { Errors <- runTaskMgr() }()

	for {
		select {
		case err := <-Errors:
			log.Println(err)
			time.Sleep(3 * time.Second)
			if RUN_HTTP_SRV == 0 {
				go func() { Errors <- runHTTPSrv() }()
			}
			if RUN_TASK_MGR == 0 {
				go func() { Errors <- runTaskMgr() }()
			}
		}
	}
}

func runTaskMgr() error {
	RUN_TASK_MGR = 1
	defer func() { RUN_TASK_MGR = 0 }()
	for {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func runHTTPSrv() error {
	RUN_HTTP_SRV = 1
	defer func() { RUN_HTTP_SRV = 0 }()

	return http.StartHTTP()
}
