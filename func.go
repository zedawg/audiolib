package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"math"
	"os"
	"path"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/zedawg/goaudiobook/sql"
)

func buildLibrary() error {
	if err := sql.CleanEntriesList(); err != nil {
		return err
	}
	if err := importEntries(); err != nil {
		log.Println(err)
		return err
	}
	log.Println("entries import ok")
	if err := indexEntries(); err != nil {
		log.Println(err)
		return err
	}
	log.Println("entries index ok")
	if err := copyDiskImages(); err != nil {
		log.Println(err)
		return err
	}
	log.Println("disk images copy ok")
	if err := copyEmbeddedImages(); err != nil {
		log.Println(err)
		return err
	}
	log.Println("embedded images copy ok")
	if err := fetchProviderData(); err != nil {
		log.Println(err)
		return err
	}
	log.Println("provider data fetch ok")
	return nil
}

func importEntries() error {
	tx, err := sql.DB.Begin()
	if err != nil {
		return err
	}
	n := 0
	for _, s := range sql.GetSources() {
		if err := fs.WalkDir(os.DirFS(s), ".", func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Println(err)
				return nil
			}
			switch {
			case p == ".":
			case d.IsDir():
				sql.InsertDir(tx, p, s)
			case path.Dir(p) == ".":
			default:
				sql.InsertFile(tx, p, s)
			}
			n = n + 1
			fmt.Printf("\rimporting entries (%v)", n)
			return nil
		}); err != nil {
			return err
		}
	}
	fmt.Printf("\r\x1b[2K")
	return tx.Commit()
}

func indexEntries() error {
	fmt.Printf("\rindexing entries...")
	if err := copyFileExts(); err != nil {
		return err
	}
	for n := 0; n < 3; n++ {
		if err := copyProbeData(); err != nil {
			log.Println(err)
			time.Sleep(time.Duration(3*n) * time.Second)
		} else {
			break
		}
	}
	if err := sql.SetDirProps(); err != nil {
		return err
	}

	return nil
}

// write file extension to files props->>'$.ext'
func copyFileExts() error {
	tx, err := sql.DB.Begin()
	if err != nil {
		return err
	}
	for _, ps := range sql.NoExtFiles() {
		sql.UpdateFileExt(tx, ps[1], ps[0])
	}
	return tx.Commit()
}

// write ffprobe data to files props->>'$.probe_data'
func copyProbeData() error {
	fmt.Printf("\rcopying ffprobe output...")
	tx, err := sql.DB.Begin()
	if err != nil {
		return err
	}
	probeCh := make(chan sql.FFProbeData)
	files := sql.NoProbeDataFiles()
	for _, ps := range files {
		go func(ps []string) {
			d := sql.FFProbeData{Source: ps[0], File: ps[1]}
			d.Data, d.Err = ffmpeg.Probe(path.Join(ps[0], ps[1]), nil)
			probeCh <- d
		}(ps)
	}
	for i := 0; i < len(files); i++ {
		err = sql.InsertProbeData(tx, <-probeCh)
		if err != nil {
			return err
		} else {
			fmt.Printf("\rcopying ffprobe output: %.2f%% (%v/%v)", float32(i)/float32(len(files))*100.0, i, len(files))
		}
	}
	fmt.Printf("\r\x1b[2K")
	return tx.Commit()
}

func copyDiskImages() error {
	fmt.Printf("\rcopying disk images...")
	tx, err := sql.DB.Begin()
	if err != nil {
		return err
	}
	// select files from directories that don't have images
	rows, err := tx.Query(`SELECT c.name || '/' || a.name AS full_path, a.dir_id, a.props->>'$.ext' AS ext, a.props->>'$.probe_data.streams[0].width' AS width, a.props->>'$.probe_data.streams[0].height' AS height FROM files a JOIN directories b ON b.id=a.dir_id JOIN sources c ON c.id=b.source_id WHERE a.props->>'$.ext' IN ('png','jpg','jpeg') AND a.dir_id NOT IN (SELECT dir_id FROM images)`)
	if err != nil {
		return err
	}
	defer rows.Close()
	type imageData struct {
		FullPath string
		DirID    int
		Ext      string
		Width    int
		Height   int
	}
	dirImages := map[int]imageData{}
	for rows.Next() {
		d := imageData{}
		if err := rows.Scan(&d.FullPath, &d.DirID, &d.Ext, &d.Width, &d.Height); err != nil {
			return err
		}
		img, ok := dirImages[d.DirID]
		if ok {
			r1 := float64(img.Width) / float64(img.Height)
			r2 := float64(d.Width) / float64(d.Height)
			if (math.Abs(r1) - 1) <= (math.Abs(r2) - 1) {
				dirImages[d.DirID] = d
				continue
			}
		}
		imageBytes, err := os.ReadFile(d.FullPath)
		if err != nil {
			return err
		}
		propsBytes, err := json.Marshal(map[string]interface{}{
			"width":  d.Width,
			"height": d.Height,
			"ext":    d.Ext,
		})
		// if existing image exists
		if img.Width*img.Height > 0 {
			if _, err := tx.Exec(`UPDATE images SET data=?, props=JSON(?) WHERE dir_id=?`, imageBytes, string(propsBytes), d.DirID); err != nil {
				return err
			}
		} else {
			if _, err := tx.Exec(`INSERT INTO images (data, dir_id, props) VALUES (?, ?, JSON(?))`, imageBytes, d.DirID, string(propsBytes)); err != nil {
				return err
			}
		}
		dirImages[d.DirID] = d
	}
	fmt.Printf("\r\x1b[2K")
	return tx.Commit()
}

func copyEmbeddedImages() error {
	fmt.Printf("\rcopying ffmpeg cover art...")
	tx, err := sql.DB.Begin()
	if err != nil {
		return err
	}
	// select directories probe_data and files list
	rows, err := tx.Query(`SELECT a.id AS dir_id, json_group_array(jsonb_object('id', c.id, 'name', b.name || '/' || c.name, 'streams', c.props->'$.probe_data.streams')) AS files FROM directories a JOIN sources b ON b.id=a.source_id JOIN files c ON c.dir_id=a.id LEFT JOIN images d ON d.dir_id=a.id WHERE d.id IS NULL AND c.props->>'$.ext' IN ('mp3','m4a','m4b') GROUP BY a.id, a.name HAVING count(c.id) > 0;`)
	if err != nil {
		return err
	}
	defer rows.Close()
	type media struct {
		DirID int
		Files []struct {
			ID      int    `json:"id"`
			Name    string `json:"name"`
			Streams []struct {
				Index       int    `json:"index"`
				CodecName   string `json:"codec_name"`
				CodecType   string `json:"codec_type"`
				Height      int    `json:"height"`
				Width       int    `json:"width"`
				Disposition struct {
					AttachedPic int `json:"attached_pic"`
				} `json:"disposition"`
			} `json:"streams"`
		}
	}
	medias := []media{}
	for rows.Next() {
		m, files := media{}, ""
		if err := rows.Scan(&m.DirID, &files); err != nil {
			log.Println(err)
			return err
		}
		if err := json.Unmarshal([]byte(files), &m.Files); err != nil {
			log.Println(err)
			return err
		}
		medias = append(medias, m)
	}
	// extract cover art
	type coverData struct {
		DirID int
		Bytes []byte
		Props []byte
	}
	coverCh := make(chan coverData)
	for _, m := range medias {
		go func(m media) {
			for _, f := range m.Files {
				for _, s := range f.Streams {
					if s.CodecType == "video" || s.Disposition.AttachedPic > 0 {
						buf := bytes.NewBuffer(nil)
						if err := ffmpeg.Input(f.Name).Output("pipe:", ffmpeg.KwArgs{"map": fmt.Sprintf("0:%d", s.Index), "c": "copy", "f": "image2pipe"}).WithOutput(buf).Silent(true).Run(); err != nil {
							log.Println(err)
							continue
						}
						imageBytes := buf.Bytes()

						ext := ""
						switch s.CodecName {
						case "mjpeg":
							ext = "jpeg"
						case "png":
							ext = "png"
						case "webp":
							ext = "webp"
						default:
							ext = ""
						}
						propsBytes, err := json.Marshal(map[string]interface{}{
							"width":  s.Width,
							"height": s.Height,
							"ext":    ext,
						})
						if err != nil {
							log.Println(err)
							continue
						}
						coverCh <- coverData{
							DirID: m.DirID,
							Bytes: imageBytes,
							Props: propsBytes,
						}
						return
					}
				}
			}
			coverCh <- coverData{}
		}(m)
	}
	for i := 0; i < len(medias); i++ {
		c := <-coverCh
		if c.DirID == 0 {
			continue
		}
		if _, err := tx.Exec(`INSERT INTO images (dir_id, data, props) VALUES (?, ?, JSON(?))`, c.DirID, c.Bytes, c.Props); err != nil {
			log.Println(err)
			continue
		}
		fmt.Printf("\rcopying ffmpeg cover art: %.2f%% (%v/%v)", float32(i)/float32(len(medias))*100.0, i, len(medias))
	}
	fmt.Printf("\r\x1b[2K")
	return tx.Commit()
}

func fetchProviderData() error {
	fmt.Printf("\rdownloading provider data...")
	str, dirs := "", []struct {
		ID       int      `json:"id"`
		Keywords []string `json:"keywords"`
	}{}
	tx, err := sql.DB.Begin()
	if err != nil {
		return err
	}
	// get directories id's, with mp3/m4a/m4b files, no image, and no provider_search
	if err := tx.QueryRow(`SELECT json_group_array(jsonb_object('id', id, 'keywords', json(keywords))) FROM (SELECT a.id, a.props->'$.keywords' AS keywords FROM directories a JOIN files b ON b.dir_id=a.id LEFT JOIN images c ON c.dir_id=a.id WHERE b.props->>'$.ext' IN ('mp3','m4a','m4b') AND a.id NOT IN (SELECT dir_id FROM images) AND a.id NOT IN (SELECT dir_id FROM provider_searches) GROUP BY a.id)`).Scan(&str); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(str), &dirs); err != nil {
		return err
	}
	fmt.Printf(" found %v missing art(s)", len(dirs))
	resCh, n := make(chan []AudiobookResult), 0
	for _, d := range dirs {
		for _, keyword := range d.Keywords {
			n = n + 1
			go func(id int, keyword string) {
				res := searchKeyword(keyword)
				for i := 0; i < len(res); i++ {
					res[i].Query = keyword
					res[i].DirID = id
				}
				resCh <- res
			}(d.ID, keyword)
		}
	}
	for i := 0; i < n; i++ {
		res := <-resCh
		if len(res) > 0 {
			b, err := json.Marshal(res)
			if err != nil {
				log.Println(err)
				continue
			}
			if _, err := tx.Exec(`INSERT INTO provider_searches (query, results, dir_id) VALUES (?, JSON(?), ?) ON CONFLICT(query, dir_id) DO UPDATE SET results = excluded.results, updated = CURRENT_TIMESTAMP;`, res[0].Query, string(b), res[0].DirID); err != nil {
				continue
			}
		}
		fmt.Printf("\rdownloading provider data: %.2f%% (%v/%v)", float32(i)/float32(n)*100.0, i, n)
	}
	fmt.Printf("\r\x1b[2K")
	return tx.Commit()
}

type AudiobookResult struct {
	Title       string   `json:"title"`
	Authors     []string `json:"authors"`
	Description string   `json:"description"`
	InfoLink    string   `json:"info"`
	CoverImage  string   `json:"cover"`
	Query       string   `json:"-"`
	DirID       int      `json:"-"`
}

func searchKeyword(query string) []AudiobookResult {
	results, resultsCh := []AudiobookResult{}, make(chan []AudiobookResult)
	go func() {
		res, _ := searchOpenLibrary(query)
		resultsCh <- res
	}()
	go func() {
		res, _ := searchITunes(query)
		resultsCh <- res
	}()
	go func() {
		res, _ := searchLibriVox(query)
		resultsCh <- res
	}()
	go func() {
		res, _ := searchArchiveOrg(query)
		resultsCh <- res
	}()
	for i := 0; i < 4; i++ {
		res := <-resultsCh
		if res == nil {
			continue
		}
		results = append(results, res...)
	}
	return results
}
