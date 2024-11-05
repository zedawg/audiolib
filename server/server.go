package server

import (
	"embed"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/zedawg/goaudiobook/config"
	"github.com/zedawg/goaudiobook/sql"
)

//go:embed files
var FS embed.FS

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

func ImagesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := sql.GetImageData(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Write(b)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct{}

type WebsocketMessage struct {
	Action string `json:"action"`
	State  struct {
		Session string `json:"session"`
		Sort    string `json:"sort"`
	} `json:"state"`
	Data   []byte `json:"data"`
	String string `json:"string"`
}

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
		case "update-state":
		case "login":
		case "list-audiobooks":
			b, err := sql.GetAudiobooks()
			if err != nil {
				log.Println(err)
			}
			m.String = string(b)
		default:
			log.Println(m.Action)
		}

		err = conn.WriteJSON(m)
		if err != nil {
			log.Println(err)
			break
		}
	}
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
