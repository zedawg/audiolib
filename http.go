package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/zedawg/audiolib/config"
	"github.com/zedawg/audiolib/db"
)

var (
	tpl *template.Template
)

func StartHTTP() {
	if DEV {
		http.HandleFunc("GET /assets/", http.HandlerFunc(DevFileHandlerFunc))
		http.HandleFunc("GET /public/", http.HandlerFunc(DevFileHandlerFunc))
	} else {
		http.Handle("GET /assets/", http.FileServer(http.FS(AssetsFS)))
		http.Handle("GET /public/", http.FileServer(http.FS(PublicFS)))
	}
	http.Handle("GET /{$}", http.HandlerFunc(HTMLHandlerFunc))
	http.Handle("GET /books", http.HandlerFunc(BooksHandlerFunc))
	http.Handle("GET /tasks", http.HandlerFunc(TasksHandlerFunc))
	http.Handle("GET /config", http.HandlerFunc(ConfigHandlerFunc))
	http.Handle("POST /search", http.HandlerFunc(SearchHandlerFunc))

	http.HandleFunc("/ws", websocketHandler)

	if err := http.ListenAndServe(config.C.Port, nil); err != nil {
		log.Println(err)
	}
}

func executeTemplate(w http.ResponseWriter, name string, data any) (err error) {
	if DEV {
		t, err := template.ParseGlob("templates/*")
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

func getUrlparam(r *http.Request, name, defaultValue string) string {
	v := r.URL.Query().Get(name)
	if len(v) == 0 {
		return defaultValue
	}
	return v
}

func DevFileHandlerFunc(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join(".", r.URL.Path))
}

func HTMLHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if err := executeTemplate(w, "html", config.C); err != nil {
		log.Println(err)
	}
}

func BooksHandlerFunc(w http.ResponseWriter, r *http.Request) {
	sort := getUrlparam(r, "sort", "added")
	limit, _ := strconv.Atoi(getUrlparam(r, "limit", "100"))
	offset, _ := strconv.Atoi(getUrlparam(r, "offset", "0"))

	books, err := db.GetBooks(sort, limit, offset)
	if err != nil {
		log.Println(err)
	}
	if err := executeTemplate(w, "pages.books", books); err != nil {
		log.Println(err)
	}
}

func TasksHandlerFunc(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.GetTasks(100, 0)
	if err != nil {
		log.Println(err)
	}
	tasks = append(tasks, &db.Task{Name: "scan name", Status: "50%"})
	if err = executeTemplate(w, "tasks", tasks); err != nil {
		log.Println(err)
	}
}

func SearchHandlerFunc(w http.ResponseWriter, r *http.Request) {
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

func CreateLibraryHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Hx-Trigger", `closeModal`)
	if err := executeTemplate(w, "libraries", struct{}{}); err != nil {
		log.Println(err)
	}
}

func UpdateLibraryHandlerFunc(w http.ResponseWriter, r *http.Request) {}

func ConfigHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if err := executeTemplate(w, "config", config.C); err != nil {
		log.Println(err)
	}

}
