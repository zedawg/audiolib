package main

import "time"

type Task struct {
	ID      int
	Name    string
	Status  string
	Params  string
	Result  string
	Created time.Time
}
