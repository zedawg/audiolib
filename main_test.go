package main

import (
	"encoding/json"
	"path"
	"testing"

	"github.com/zedawg/librarian/config"
	"github.com/zedawg/librarian/sql"
)

func TestNothing(t *testing.T) {
	//
	sql.Open(config.DatabasePath())

	if sql.DB == nil {
		t.Error("db is nil")
	}

	var (
		values []string
		val    string
	)
	rows, err := sql.DB.Query(`SELECT json_object('id',id,'name',name,'root',root,'entries_details',entries_details) FROM directories_details LIMIT 3`)
	if err != nil {
		t.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&val); err != nil {
			t.Error(err)
		}
		values = append(values, val)
	}

	type entry struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Ext        string `json:"ext"`
		DetailsStr string `json:"details"`
		Details    map[string]interface{}
	}
	type dirdetails struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Root       string `json:"root"`
		EntriesStr string `json:"entries_details"`
		Entries    []entry
	}

	var (
		dirs []dirdetails
		dir  dirdetails
	)

	for _, val := range values {
		if err := json.Unmarshal([]byte(val), &dir); err != nil {
			t.Error(err)
		}
		if err := json.Unmarshal([]byte(dir.EntriesStr), &dir.Entries); err != nil {
			t.Error(err)
		}
		for i := 0; i < len(dir.Entries); i++ {
			if err := json.Unmarshal([]byte(dir.Entries[i].DetailsStr), &dir.Entries[i].Details); err != nil {
				t.Error(err)
			}
		}
		dirs = append(dirs, dir)
	}

	for _, dir := range dirs {
		t.Logf("\t%s", path.Clean(path.Join(dir.Root, dir.Name)))
		for _, file := range dir.Entries {
			t.Logf("\t\t%s", file.Name)
		}
	}
	t.Log(len(dirs), dirs[1].Entries[0].Ext)
}
