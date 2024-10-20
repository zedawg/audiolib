package server

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/zedawg/librarian/config"
	"github.com/zedawg/librarian/sql"
)

var (
	//go:embed files
	FS       embed.FS
	upgrader = websocket.Upgrader{
		ReadBufferSize:    4096,
		WriteBufferSize:   4096,
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func Listen() error {
	http.Handle("/", FilesHandler())
	http.Handle("/socket", &WebSocketHandler{})
	http.HandleFunc("/images/", ImagesHandler)
	return http.ListenAndServe(":"+config.Port, nil)
}

func FilesHandler() http.Handler {
	if config.Dev {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, path.Join("./server/files", r.URL.Path))
		})
	} else {
		return http.FileServer(http.FS(FS))
	}
}

type WebsocketMessage struct {
	Action Action         `json:"action"`
	State  clientAppState `json:"state"`
	Data   []byte         `json:"data"`
	String string         `json:"string"`
}

type Action string

const (
	ActionUpdateState Action = "update-state"
	ActionLogin       Action = "login"
	ActionListDir     Action = "list-dir"
)

type clientAppState struct {
	Session string
	Sort    string
}

type WebSocketHandler struct{}

func (hd *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	m := WebsocketMessage{}
	for {
		err := conn.ReadJSON(&m)
		if err != nil {
			log.Println(err)
			break
		}

		switch m.Action {
		case ActionUpdateState:
		case ActionLogin:
		case ActionListDir:
			rows, err := sql.DB.Query(`SELECT json_object('id', id, 'name', name, 'entries', entries_details) FROM directories_details ORDER BY added DESC`)
			if err != nil {
				log.Println(err)
			}
			defer rows.Close()
			dirs, dir := []string{}, ""
			for rows.Next() {
				if err := rows.Scan(&dir); err != nil {
					continue
				}
				dirs = append(dirs, dir)
			}
			m.String = fmt.Sprintf("[%v]", strings.Join(dirs, ","))
		default:
		}

		err = conn.WriteJSON(m)
		if err != nil {
			log.Println(err)
			break
		}
	}
}

// func AppHandler(w http.ResponseWriter, r *http.Request) {
// 	if err := executeTemplate(w, "main", nil); err != nil {
// 		log.Println(err)
// 	}
// }

func ImagesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p := ""
	err = sql.DB.QueryRow(`SELECT a.root || '/' || a.name || '/' || b.name FROM directories a JOIN entries b ON a.id=b.directory_id WHERE b.id=? AND b.ext IN ('jpg','jpeg','png')`, id).Scan(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, path.Clean(p))
}

// func executeTemplate(w http.ResponseWriter, name string, data any) (err error) {
// 	if config.Dev {
// 		t, err := template.ParseGlob("server/templates/*")
// 		if err != nil {
// 			return err
// 		}
// 		return t.ExecuteTemplate(w, name, data)
// 	}
// 	return tpl.ExecuteTemplate(w, name, data)
// }

// func getURLParam(r *http.Request, name, defaultValue string) string {
// 	v := r.URL.Query().Get(name)
// 	if len(v) == 0 {
// 		return defaultValue
// 	}
// 	return v
// }

// func serveBooks(w http.ResponseWriter, r *http.Request) {
// 	var (
// 		sort      = getURLParam(r, "sort", "added")
// 		limit, _  = strconv.Atoi(getURLParam(r, "limit", "100"))
// 		offset, _ = strconv.Atoi(getURLParam(r, "offset", "0"))
// 	)

// 	books, err := sql.GetBooks(sort, limit, offset)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	if err := executeTemplate(w, "pages.books", books); err != nil {
// 		log.Println(err)
// 	}
// }

// func search(w http.ResponseWriter, r *http.Request) {
// 	escape := func(q string) string {
// 		return q
// 	}
// 	q := escape(r.FormValue("search"))
// 	if len(q) <= 2 {
// 		if err := executeTemplate(w, "search-results", struct{}{}); err != nil {
// 			log.Println(err)
// 		}
// 		return
// 	}
// 	results, err := sql.SearchBooks(q)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	for i := 0; i < 20; i++ {
// 		results = append(results, &sql.Book{ID: i, Title: fmt.Sprintf("title %v", i), Author: fmt.Sprintf("authors %v", i)})
// 	}

// 	if err = executeTemplate(w, "search-results", results); err != nil {
// 		log.Println(err)
// 	}
// }

// func watchFiles() {
// 	go func() {
// 		watcher, err := fsnotify.NewWatcher()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer watcher.Close()
// 		go func() {
// 			for {
// 				select {
// 				case event, ok := <-watcher.Events:
// 					if !ok {
// 						return
// 					}
// 					log.Println("event:", event)
// 					if event.Has(fsnotify.Write) {
// 						log.Println("modified file:", event.Name)
// 					}
// 				case err, ok := <-watcher.Errors:
// 					if !ok {
// 						return
// 					}
// 					log.Println("error:", err)
// 				}
// 			}
// 		}()
// 		err = watcher.Add("./server/files")
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		<-make(chan struct{})
// 	}()
// }

// func writeSessionCookie(w http.ResponseWriter, s sql.Session) {
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "sessionid",
// 		Value:    s.ID,
// 		Expires:  s.Expires,
// 		Path:     "/",
// 		HttpOnly: true,
// 	})
// }
