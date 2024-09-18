package main

import "time"

type LogEntry struct {
	ID      int
	Message string
	Created time.Time
}

type TaskEntry struct {
	ID       int
	Name     string
	Priority int
	Status   string
	Params   [][]any
	Result   [][]any
	Created  time.Time
}

type LibraryEntry struct {
	ID      int
	Name    string
	Paths   []any
	Created time.Time
}

type FileEntry struct {
	ID       int
	Name     string
	Created  time.Time
	Modified time.Time
	TaskID   int
}

type AudiobookEntry struct {
	ID            int
	Title         string
	Subtitle      string
	Authors       string
	Narrator      string
	Genres        string
	ISBN          string
	ASIN          string
	Language      string
	Year          int
	Duration      int
	Chapters      [][]any
	Cover         []byte
	ConvertedFile string
	Created       time.Time
}

type PropertyEntry struct {
	Name  string
	Value string
}

type UserEntry struct {
	ID           int
	Name         string
	PasswordHash string
	IsAdmin      bool
}
