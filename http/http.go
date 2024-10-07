package http

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/zedawg/librarian/config"
	"github.com/zedawg/librarian/db"
)

var (
	//go:embed assets
	AssetsFS embed.FS
	//go:embed public
	PublicFS embed.FS
	//go:embed templates
	TemplateFS embed.FS
	tpl        *template.Template
)

func StartHTTP() error {
	if config.Dev {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, path.Join("./http", r.URL.Path))
		})
	} else {
		http.Handle("GET /assets/", http.FileServer(http.FS(AssetsFS)))
		http.Handle("GET /public/", http.FileServer(http.FS(PublicFS)))
	}
	http.HandleFunc("GET /{$}", serveHTML)
	http.HandleFunc("GET /books", serveBooks)
	http.HandleFunc("GET /tasks", serveTasks)
	http.HandleFunc("POST /search", search)
	http.Handle("/sock", &socker{})
	return http.ListenAndServe(config.Port(), nil)
}

func Middleware(hd http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// r.WithContext(r.Context(), "")

		hd.ServeHTTP(w, r)
	})
}

func executeTemplate(w http.ResponseWriter, name string, data any) (err error) {
	if config.Dev {
		t, err := template.ParseGlob("http/templates/*")
		if err != nil {
			return err
		}
		return t.ExecuteTemplate(w, name, data)
	}
	if tpl == nil {
		tpl, err = template.ParseFS(TemplateFS, "templates/*")
	}
	return tpl.ExecuteTemplate(w, name, data)
}

func getURLParam(r *http.Request, name, defaultValue string) string {
	v := r.URL.Query().Get(name)
	if len(v) == 0 {
		return defaultValue
	}
	return v
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	if err := executeTemplate(w, "html", nil); err != nil {
		log.Println(err)
	}
}

func serveBooks(w http.ResponseWriter, r *http.Request) {
	sort := getURLParam(r, "sort", "added")
	limit, _ := strconv.Atoi(getURLParam(r, "limit", "100"))
	offset, _ := strconv.Atoi(getURLParam(r, "offset", "0"))

	books, err := db.GetBooks(sort, limit, offset)
	if err != nil {
		log.Println(err)
	}
	if err := executeTemplate(w, "pages.books", books); err != nil {
		log.Println(err)
	}
}

func serveTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.GetTasks(100, 0)
	if err != nil {
		log.Println(err)
	}
	tasks = append(tasks, &db.Task{Name: "scan name", Status: "50%"})
	if err = executeTemplate(w, "tasks", tasks); err != nil {
		log.Println(err)
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	escape := func(q string) string {
		return q
	}
	q := escape(r.FormValue("search"))
	if len(q) <= 2 {
		if err := executeTemplate(w, "search-results", struct{}{}); err != nil {
			log.Println(err)
		}
		return
	}
	results, err := db.SearchBooks(q)
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < 20; i++ {
		results = append(results, &db.Book{ID: i, Title: fmt.Sprintf("title %v", i), Author: fmt.Sprintf("authors %v", i)})
	}

	if err = executeTemplate(w, "search-results", results); err != nil {
		log.Println(err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		// In production, you should verify the origin here.
		return true
	},
}

type socker struct{}

func (s *socker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		log.Println("websocket: listening")
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Read error: %v", err)
			}
			break
		}

		log.Printf("Received: %s", message)

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
