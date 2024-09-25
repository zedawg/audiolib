package db

import (
	"strings"
	"time"
)

type Book struct {
	ID       int
	Title    string
	Author   string
	Narrator string
	ISBN     string
	ASIN     string
	Genre    string
	Language string
	Chapters string
	Provider string
	Year     int
	Duration int
	Added    time.Time
}

func GetBooks(sort string, limit, offset int) (books []*Book, err error) {
	rows, err := DB.Query(`SELECT id, title, author, narrator, isbn, asin, genre, language, year, duration, chapters, provider, added FROM books ORDER BY ? LIMIT ? OFFSET ?`, sort, limit, offset)
	if err != nil {
		return
	}
	for rows.Next() {
		v := Book{}
		if err = rows.Scan(&v.ID, &v.Title, &v.Author, &v.Narrator, &v.ISBN, &v.ASIN, &v.Genre, &v.Language, &v.Year, &v.Duration, &v.Chapters, &v.Provider, &v.Added); err != nil {
			return
		}
		v.Title = strings.Title(v.Title)
		v.Author = strings.Title(v.Author)
		v.Narrator = strings.Title(v.Narrator)
		v.Genre = strings.Title(v.Genre)
		v.Language = strings.Title(v.Language)

		books = append(books, &v)
	}
	return
}

func InsertBook(title, author, narrator, isbn, asin, genre, language, chapters, provider string, year, duration int) (err error) {
	_, err = DB.Exec(`INSERT INTO books (title, author, narrator, isbn, asin, genre, language, chapters, provider, year, duration) VALUES (?,?,?,?,?,?,?,?,?,?,?)`, title, author, narrator, isbn, asin, genre, language, chapters, provider, year, duration)
	return
}

func SearchBooks(q string) (results []*Book, err error) {
	q = q + "%"
	rows, err := DB.Query(`SELECT id, title, author, narrator, isbn, asin, genre, language, chapters, provider, year, duration, added FROM books WHERE title LIKE ? OR author LIKE ? OR isbn LIKE ? OR asin LIKE ?`, q, q, q, q)
	if err != nil {
		return
	}
	for rows.Next() {
		v := Book{}
		err = rows.Scan(&v.ID, &v.Title, &v.Author, &v.Narrator, &v.ISBN, &v.ASIN, &v.Genre, &v.Language, &v.Chapters, &v.Provider, &v.Year, &v.Duration, &v.Added)
		if err != nil {
			return
		}
		results = append(results, &v)
	}
	return
}

type Task struct {
	ID       int
	Name     string
	Status   string
	Params   string
	Result   string
	Priority int
	Created  time.Time
}

func GetTasks(limit, offset int) (tasks []*Task, err error) {
	rows, err := DB.Query(`
SELECT id, name, priority, status, params, result, created
FROM tasks
ORDER BY created LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return
	}
	for rows.Next() {
		v := Task{}
		if err = rows.Scan(&v.ID, &v.Name, &v.Priority, &v.Status, &v.Params, &v.Result, &v.Created); err != nil {
			return
		}
		tasks = append(tasks, &v)
	}
	return
}

func InsertTask(name, status, params string, priority int) error {
	_, err := DB.Exec(`INSERT INTO tasks (name, status, params, priority) VALUES (?, ?, ?, ?)`, name, status, params, priority)
	return err
}

func UpdateTask(id int, status, result string) error {
	_, err := DB.Exec(`UPDATE tasks SET status=COALESCE(NULLIF(?, ''), status) AND result=COALESCE(NULLIF(?, '') WHERE id=?`, status, result, id)
	return err
}

type File struct {
	ID       int
	Name     string
	Created  time.Time
	Modified time.Time
}
