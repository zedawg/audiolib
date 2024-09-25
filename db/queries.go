package db

import (
	"log"
	"time"
)

type Library struct {
	ID            int
	Name          string
	ScanPath      string
	ConvertedPath string
	Created       time.Time
}

func GetLibraries() (libraries []*Library, err error) {
	rows, err := DB.Query(`SELECT id, name, import_path, converted_path, created FROM libraries ORDER BY created DESC`)
	if err != nil {
		return
	}
	for rows.Next() {
		l := Library{}
		if err = rows.Scan(&l.ID, &l.Name, &l.ImportPath, &l.ConvertedPath, &l.Created); err != nil {
			return
		}
		libraries = append(libraries, &l)
	}
	return
}

func CreateLibrary(name, importPath, convertedPath string) error {
	_, err := DB.Exec(`
INSERT INTO libraries (name, import_path, converted_path)
VALUES (?,?,?)`, name, importPath, convertedPath)
	return err
}

type Audiobook struct {
	ID       int
	Title    string
	Subtitle string
	Authors  string
	Narrator string
	Genre    string
	ISBN     string
	ASIN     string
	Language string
	Year     int
	Duration int
	Chapters string
	Created  time.Time
}

func GetAudiobooks() (audiobooks []*Audiobook, err error) {
	rows, err := DB.Query(`
SELECT id, title, subtitle, authors, narrator, genre, isbn, asin, language, year, duration, chapters, created
FROM audiobooks`)
	if err != nil {
		return
	}
	for rows.Next() {
		v := Audiobook{}
		if err := rows.Scan(&v.ID, &v.Title, &v.Subtitle, &v.Authors, &v.Narrator, &v.Genre, &v.ISBN, &v.ASIN, &v.Language, &v.Year, &v.Duration, &v.Chapters, &v.Created); err != nil {
			log.Println(err)
		}
		audiobooks = append(audiobooks, &v)
	}
	return
}

type Task struct {
	ID       int
	Name     string
	Priority int
	Status   string
	Params   string
	Result   string
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

func Search(q string) (results []*Audiobook, err error) {
	q = q + "%"
	rows, err := DB.Query(`
SELECT id, title, subtitle, authors, narrator, genre, isbn, asin, language, year, duration, chapters, created
FROM audiobooks
WHERE title LIKE ? OR subtitle LIKE ? OR authors LIKE ? OR isbn LIKE ? OR asin LIKE ?`, q, q, q, q, q)
	if err != nil {
		return
	}
	for rows.Next() {
		v := Audiobook{}
		if err := rows.Scan(&v.ID, &v.Title, &v.Subtitle, &v.Authors, &v.Narrator, &v.Genre, &v.ISBN, &v.ASIN, &v.Language, &v.Year, &v.Duration, &v.Chapters, &v.Created); err != nil {
			log.Println(err)
		}
		results = append(results, &v)
	}
	return
}

type File struct {
	ID       int
	Name     string
	Created  time.Time
	Modified time.Time
	TaskID   int
}

type User struct {
	ID           int
	Name         string
	PasswordHash string
	IsAdmin      bool
}
