package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"

	"github.com/zedawg/librarian/config"
	"github.com/zedawg/librarian/server"
	"github.com/zedawg/librarian/sql"

	"github.com/dhowden/tag"
)

var parseDirs = []string{
	"/users/admin/downloads/audiobooks-1",
	"/users/admin/downloads/audiobooks-2",
}

func main() {
	sql.Open(config.DatabasePath())
	log.Println("database", strings.ToLower(config.DatabasePath()))
	log.Printf("browse audiobooks at http://localhost:%v", config.Port)

	go func() {
		parse()
		extractPicturesAll()
		writeMetadata() // copy id3, mp4 tag data to sqlite
	}()

	if err := server.Listen(); err != nil {
		log.Fatal(err)
	}
}

func parse() {
	for _, p := range parseDirs {
		parseDir(p)
	}
}

func parseDir(dir string) error {
	dir = strings.ToLower(path.Clean(dir))
	tx, _ := sql.DB.Begin()
	err := fs.WalkDir(os.DirFS(dir), ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		p = strings.ToLower(path.Clean(p))
		if d.IsDir() {
			tx.Exec(`INSERT INTO directories (name, root) VALUES (?, ?)`, p, dir)
		} else {
			ext := path.Ext(p)
			switch ext {
			case ".jpg":
			case ".jpeg":
			case ".png":
			case ".mp3":
			case ".m4a":
			case ".m4b":
			default:
				return nil
			}
			tx.Exec(`INSERT INTO entries (name, ext, directory_id) VALUES (?, ?, (SELECT id FROM directories WHERE name = ?))`, path.Base(p), strings.Trim(ext, "."), path.Dir(p))
		}
		return nil
	})
	if err != nil {
		return errors.Join(tx.Rollback(), err)
	}
	return tx.Commit()
}

func extractPicturesAll() error {
	tx, _ := sql.DB.Begin()
	rows, err := tx.Query(`SELECT json_object('name', name, 'root', root, 'children', entries) FROM directories_details`)
	if err != nil {
		return err
	}
	defer rows.Close()

	dat, dir := "", struct {
		Name        string `json:"name"`
		Root        string `json:"root"`
		ChildrenStr string `json:"children"`
		Children    []string
	}{}
	for rows.Next() {
		if err := rows.Scan(&dat); err != nil {
			log.Println(err)
			continue
		}
		if err := json.Unmarshal([]byte(dat), &dir); err != nil {
			log.Println(err)
			continue
		}
		if err := json.Unmarshal([]byte(dir.ChildrenStr), &dir.Children); err != nil {
			log.Println(err)
			continue
		}
		sizes := map[int]struct{}{}
		for _, name := range dir.Children {
			f, err := os.Open(path.Join(dir.Root, dir.Name, name))
			if err != nil {
				continue
			}
			defer f.Close()
			t, err := tag.ReadFrom(f)
			if err != nil {
				continue
			}
			p := t.Picture()
			if p == nil {
				continue
			}
			if _, ok := sizes[len(p.Data)]; ok {
				continue
			}
			if p.Ext == "" || p.Ext == "jpg" {
				p.Ext = "jpeg"
			}
			os.WriteFile(path.Join(dir.Root, dir.Name, fmt.Sprintf("~%v.%v", path.Base(name), p.Ext)), p.Data, 0644)
			sizes[len(p.Data)] = struct{}{}
		}
	}

	return tx.Commit()
}

func writeMetadata() error {
	// TODO for computing chapters in ffmpeg encoding
	rows, err := sql.DB.Query(`SELECT name, root, entries_details FROM directories_details`)
	if err != nil {
		return err
	}
	defer rows.Close()
	var (
		Name           string
		Root           string
		entriesDetails string
		EntriesDetails []struct {
			ID      int
			Name    string
			Ext     string
			Details string
		}
	)
	for rows.Next() {
		if err := rows.Scan(&Name, &Root, &entriesDetails); err != nil {
			log.Println(err)
			continue
		}
		if err := json.Unmarshal([]byte(entriesDetails), &EntriesDetails); err != nil {
			log.Println(err)
			continue
		}
		for _, e := range EntriesDetails {
			if e.Ext == "mp3" || e.Ext == "m4a" || e.Ext == "m4b" {
				// p := path.Join(Root, Name, e.Name)
				// m := extractMetadata(p)

			}

		}
	}
	return nil
}
