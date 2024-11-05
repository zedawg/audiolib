package sql

import (
	"database/sql"
	"embed"
	"encoding/json"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

var (
	//go:embed schema.sql
	Files embed.FS
	DB    *sql.DB
)

func init() {
	pathExt := func(p string) string {
		return strings.Trim(path.Ext(p), ".")
	}
	sql.Register("sqlite3_extended",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				err := conn.RegisterFunc("path_ext", pathExt, true)
				return err
			},
		},
	)
}

func Open(p string) (err error) {
	DB, err = sql.Open("sqlite3_extended", p)
	if err != nil {
		return err
	}

	b, _ := Files.ReadFile("schema.sql")
	DB.Exec(string(b))

	return DB.Ping()
}

func normalizePath(p string) string {
	return strings.ToLower(path.Clean(p))
}

func GetSources() []string {
	rows, err := DB.Query(`SELECT name FROM sources`)
	if err != nil {
		log.Println(err)
		return []string{}
	}
	defer rows.Close()
	sources := []string{}
	for rows.Next() {
		s := ""
		if err := rows.Scan(&s); err != nil {
			log.Println(err)
			return []string{}
		}
		sources = append(sources, s)
	}
	return sources
}

func DeleteSources() error {
	_, err := DB.Exec(`DELETE FROM sources`)
	return err
}

func InsertSource(source string) error {
	source = normalizePath(source)
	_, err := DB.Exec(`INSERT INTO sources (name) VALUES (?)`, source)
	return err
}

func InsertDir(tx *sql.Tx, dir, source string) error {
	dir = normalizePath(dir)
	source = normalizePath(source)
	_, err := tx.Exec(`INSERT INTO directories (name, source_id) VALUES (?, (SELECT id FROM sources WHERE name=?))`, dir, source)
	return err
}

func InsertFile(tx *sql.Tx, file, source string) error {
	file = normalizePath(file)
	source = normalizePath(source)
	_, err := tx.Exec(`INSERT INTO files (name, dir_id) VALUES
	(?, (SELECT id FROM directories WHERE name=? AND source_id=(SELECT id FROM sources WHERE name=?)))`, file, path.Dir(file), source)
	return err
}

func UpdateFileExt(tx *sql.Tx, file, source string) error {
	file = normalizePath(file)
	source = normalizePath(source)
	_, err := tx.Exec(`UPDATE files SET props=json_set(props, '$.ext', path_ext(name)) WHERE name=? AND dir_id=(SELECT id FROM directories WHERE name=? AND source_id=(SELECT id FROM sources WHERE name=?))`, file, path.Dir(file), source)
	return err
}

// returns sqlite files to extract file extension
func NoExtFiles() (files [][]string) {
	rows, err := DB.Query(`SELECT c.name, a.name FROM files a JOIN directories b ON b.id=a.dir_id JOIN sources c ON c.id=b.source_id WHERE a.props->>'$.ext' IS NULL`)
	if err != nil {
		log.Println(err)
		return files
	}
	defer rows.Close()
	for rows.Next() {
		var f, s string
		if err := rows.Scan(&s, &f); err != nil {
			log.Println(err)
			return files
		}
		files = append(files, []string{s, f})
	}
	return files
}

// returns sqlite files to ffprobe
func NoProbeDataFiles() (files [][]string) {
	rows, err := DB.Query(`SELECT c.name, a.name FROM files a JOIN directories b ON b.id=a.dir_id JOIN sources c ON c.id=b.source_id WHERE a.props->>'$.probe_data' IS NULL AND a.props->>'$.ext' IN ('mp3','m4a','m4b','jpg','jpeg','png')`)
	if err != nil {
		log.Println(err)
		return files
	}
	defer rows.Close()
	for rows.Next() {
		var f, s string
		if err := rows.Scan(&s, &f); err != nil {
			log.Println(err)
			return files
		}
		files = append(files, []string{s, f})
	}
	return files
}

type FFProbeData struct {
	File   string
	Source string
	Data   string
	Err    error
}

func InsertProbeData(tx *sql.Tx, d FFProbeData) error {
	d.File = normalizePath(d.File)
	d.Source = normalizePath(d.Source)
	_, err := tx.Exec(`UPDATE files SET props=json_set(props, '$.probe_data', JSON(?)) WHERE name=? AND dir_id=(SELECT id FROM directories WHERE name=? AND source_id=(SELECT id FROM sources WHERE name=?))`, d.Data, d.File, path.Dir(d.File), d.Source)
	return err
}

// lists database directories and files and removes them if os.Lstat errors
func CleanEntriesList() error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	rows, err := tx.Query(`SELECT b.name || '/' || a.name AS dir_path, a.id FROM directories a JOIN sources b ON b.id=a.source_id`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		p, id := "", 0
		if err := rows.Scan(&p, &id); err != nil {
			return err
		}
		if _, err := os.Lstat(p); err != nil {
			if _, err := tx.Exec(`DELETE FROM directories WHERE id=?`, id); err != nil {
				return err
			}
		}
	}
	if err := rows.Close(); err != nil {
		return err
	}
	rows, err = tx.Query(`SELECT c.name || '/' || a.name AS file_path, a.id FROM files a JOIN directories b ON b.id=a.dir_id JOIN sources c ON c.id=b.source_id`)
	if err != nil {
		return err
	}
	for rows.Next() {
		p, id := "", 0
		if err := rows.Scan(&p, &id); err != nil {
			return err
		}
		if _, err := os.Lstat(p); err != nil {
			if _, err := tx.Exec(`DELETE FROM files WHERE id=?`, id); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

type DirectoryData struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	SourceID   int    `json:"source_id"`
	MediaFiles map[string]struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Ext       string `json:"ext"`
		Duration  string `json:"duration"`
		BitRate   string `json:"bit_rate"`
		Size      string `json:"size"`
		Title     string `json:"title"`
		Artist    string `json:"artist"`
		Album     string `json:"album"`
		NBStreams int    `json:"nb_streams"`
	} `json:"media_files"`
	ImageFiles map[string]struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Ext    string `json:"ext"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"image_files"`
}

// parse directories data and directories files props to set directories props, to be used for display and matching
func SetDirProps() error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	rows, err := tx.Query(`SELECT A.id, A.name, A.source_id, COALESCE(A.media_files, '{}') AS media_files, COALESCE(B.image_files, '{}') AS image_files FROM (SELECT a.id, a.name, a.source_id, json_group_object(CAST(b.id AS TEXT), jsonb_object('id', b.id, 'name', b.name, 'ext', b.props->>'$.ext', 'duration', b.props->>'$.probe_data.format.duration', 'bit_rate', b.props->>'$.probe_data.format.bit_rate', 'size', b.props->>'$.probe_data.format.size', 'title', b.props->>'$.probe_data.format.tags.title', 'artist', b.props->>'$.probe_data.format.tags.artist', 'album', b.props->>'$.probe_data.format.tags.album', 'nb_streams', b.props->>'$.probe_data.format.nb_streams')) AS media_files FROM directories a JOIN files b ON b.dir_id=a.id JOIN sources c ON c.id=a.source_id WHERE b.props->>'$.ext' IN ('mp3','m4a','m4b') GROUP BY a.id, a.name, a.source_id) A LEFT JOIN (SELECT a.id, json_group_object(CAST(b.id AS TEXT), jsonb_object('id', b.id, 'name', b.name, 'ext', b.props->>'$.ext', 'height', b.props->>'$.probe_data.streams[0].height', 'width', b.props->>'$.probe_data.streams[0].width')) AS image_files FROM directories a JOIN files b ON b.dir_id=a.id WHERE b.props->>'$.ext' IN ('png','jpeg','jpg') GROUP BY a.id, a.name, a.source_id) B ON A.id=B.id`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		r, mediaFiles, imageFiles := DirectoryData{}, "", ""
		if err := rows.Scan(&r.ID, &r.Name, &r.SourceID, &mediaFiles, &imageFiles); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(mediaFiles), &r.MediaFiles); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(imageFiles), &r.ImageFiles); err != nil {
			return err
		}
		// compute keywords, duration, average bitrate, size, extensions
		keywordCounter, keywords := map[string]int{}, []string{}
		duration, bitrate, size := 0.0, 0, 0
		exts, extensions := map[string]struct{}{}, []string{}
		for _, f := range r.MediaFiles {
			keywordCounter[f.Title] = keywordCounter[f.Title] + 1
			keywordCounter[f.Artist] = keywordCounter[f.Artist] + 1
			keywordCounter[f.Album] = keywordCounter[f.Album] + 1

			d, err := strconv.ParseFloat(f.Duration, 32)
			if err != nil {
				d = 0.0
			}
			br, err := strconv.Atoi(f.BitRate)
			if err != nil {
				br = 0
			}
			s, err := strconv.Atoi(f.Size)
			if err != nil {
				s = 0
			}
			exts[f.Ext] = struct{}{}
			size = size + s
			bitrate = bitrate + br
			duration = duration + d
		}
		//
		for keyword, count := range keywordCounter {
			if len(strings.Trim(keyword, " ")) == 0 {
				continue
			}
			if count >= len(r.MediaFiles) {
				keywords = append(keywords, strings.Trim(keyword, " "))
			}
		}
		// if no keyword found, use directory name
		if len(keywords) == 0 {
			p := strings.ReplaceAll(r.Name, "_", " ")
			p = cases.Title(language.Und, cases.NoLower).String(p)
			keywords = append(keywords, strings.Split(p, "/")...)
		}
		keywordsBytes, err := json.Marshal(keywords)
		if err != nil {
			return err
		}
		//
		for e, _ := range exts {
			extensions = append(extensions, e)
		}
		//
		extensionsBytes, err := json.Marshal(extensions)
		if err != nil {
			extensionsBytes = []byte("[]")
		}
		//
		bitrate = bitrate / len(r.MediaFiles)

		if _, err := tx.Exec(`UPDATE directories
			SET props=json_set(
				json_set(
					json_set(
						json_set(
							json_set(props,
								'$.keywords', JSON(?)),
							'$.duration', ?),
						'$.size', ?
					), '$.bitrate', ?
				),
				'$.extensions', JSON(?)
			) WHERE id=?`,
			string(keywordsBytes),
			int(duration),
			size,
			bitrate,
			string(extensionsBytes),
			r.ID); err != nil {
			return err
		}

	}
	return tx.Commit()
}

// returns image absolute path, for FileServe use
func GetImageData(imageID int) (b []byte, err error) {
	if err = DB.QueryRow(`SELECT data FROM images WHERE id=?`, imageID).Scan(&b); err != nil {
		return nil, err
	}
	return b, nil
}

// returns a json array of all directories and associated files
func GetAudiobooks() ([]byte, error) {
	rows, err := DB.Query(`SELECT json_object('id', id, 'name', dir, 'files', files, 'image_id', image_id) FROM entries ORDER BY added DESC`)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	ms, dir := []map[string]interface{}{}, ""
	for rows.Next() {
		if err := rows.Scan(&dir); err != nil {
			continue
		}
		m := map[string]interface{}{}
		if err := json.Unmarshal([]byte(dir), &m); err != nil {
			log.Println(err)
			continue
		}
		ms = append(ms, m)
	}
	return json.Marshal(ms)
}
