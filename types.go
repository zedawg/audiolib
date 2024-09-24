package main

import "time"

type LibraryEntry struct {
	ID            int
	Name          string
	ImportPath    string
	ConvertedPath string
	Created       time.Time
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

type UserEntry struct {
	ID           int
	Name         string
	PasswordHash string
	IsAdmin      bool
}

type ActivityEntry struct {
	Created time.Time
	Value   string
	Type    string
}

type SearchResultEntry struct {
	ID      int
	Name    string
	Details string
	Image   string
	Type    string
}
