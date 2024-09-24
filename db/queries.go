package main

func (db *Database) GetLibraries() (libraryEntries []*LibraryEntry, err error) {
	rows, err := db.Query(`SELECT id, name, import_path, converted_path, created FROM libraries ORDER BY created DESC`)
	if err != nil {
		return
	}
	for rows.Next() {
		l := LibraryEntry{}
		if err = rows.Scan(&l.ID, &l.Name, &l.ImportPath, &l.ConvertedPath, &l.Created); err != nil {
			return
		}
		libraryEntries = append(libraryEntries, &l)
	}
	return
}

func (db *Database) GetActivities(limit, offset int) (activityEntries []*ActivityEntry, err error) {
	rows, err := db.Query(`
SELECT created, message AS value, 'log' AS type FROM logs UNION ALL
SELECT created, (name || ' ' || params || ' ' || result) AS value, 'task' AS type FROM tasks
ORDER BY created LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return
	}
	for rows.Next() {
		a := ActivityEntry{}
		if err = rows.Scan(&a.Created, &a.Value, &a.Type); err != nil {
			return
		}
		activityEntries = append(activityEntries, &a)
	}
	return

}

func (db *Database) Search(q string) (searchResults []*SearchResultEntry, err error) {
	rows, err := db.Query(`
SELECT 
	id, name, details, image, 'audiobook' AS type 
FROM audiobooks`)
	return
}

func (db *Database) CreateLibrary(name, importPath, convertedPath string) error {
	_, err := db.Exec(`
INSERT INTO libraries (name, import_path, converted_path)
VALUES (?,?,?)`, name, importPath, convertedPath)
	return err
}
