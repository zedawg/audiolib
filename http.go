package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/zedawg/audiolib/db"
)

var (
	tpl *template.Template
)

func executeTemplate(w http.ResponseWriter, name string, data any) (err error) {
	if dev {
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

func StartHTTP() {
	if dev {
		http.HandleFunc("GET /assets/", http.HandlerFunc(ServeFileHandlerFunc))
		http.HandleFunc("GET /public/", http.HandlerFunc(ServeFileHandlerFunc))
	} else {
		http.Handle("GET /assets/", http.FileServer(http.FS(AssetsFS)))
		http.Handle("GET /public/", http.FileServer(http.FS(PublicFS)))
	}
	http.Handle("GET /{$}", http.HandlerFunc(AudiobooksHandlerFunc))
	http.Handle("GET /tasks", http.HandlerFunc(TasksHandlerFunc))
	http.Handle("GET /settings", http.HandlerFunc(SettingsHandlerFunc))
	http.Handle("GET /user", http.HandlerFunc(UserHandlerFunc))
	http.Handle("POST /libraries", http.HandlerFunc(CreateLibraryHandlerFunc))
	http.Handle("POST /search", http.HandlerFunc(SearchHandlerFunc))
	http.Handle("PUT /libraries", http.HandlerFunc(UpdateLibraryHandlerFunc))

	if err := http.ListenAndServe(Port, nil); err != nil {
		log.Println(err)
	}
}

func ServeFileHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeFile(w, r, path.Join(".", r.URL.Path))
}

func AudiobooksHandlerFunc(w http.ResponseWriter, r *http.Request) {
	audiobooks, err := db.GetAudiobooks()
	if err != nil {
		log.Println(err)
	}
	if err := executeTemplate(w, "pages.audiobooks", audiobooks); err != nil {
		fmt.Println(err)
	}
}

func TasksHandlerFunc(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.GetTasks(100, 0)
	if err != nil {
		log.Println(err)
	}
	if err = executeTemplate(w, "pages.tasks", tasks); err != nil {
		fmt.Println(err)
	}
}

func SettingsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	libraries, err := db.GetLibraries()
	if err != nil {
		log.Println(err)
	}
	if err = executeTemplate(w, "pages.settings", libraries); err != nil {
		log.Println(err)
	}
}

func UserHandlerFunc(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Name string
	}{
		Name: "user",
	}
	var err error
	if err = executeTemplate(w, "pages.user", data); err != nil {
		fmt.Println(err)
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
	results, err := db.Search(q)
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < 20; i++ {
		results = append(results, &db.Audiobook{ID: i, Title: fmt.Sprintf("title %v", i), Authors: fmt.Sprintf("authors %v", i)})
	}

	if err = executeTemplate(w, "search-results", results); err != nil {
		log.Println(err)
	}
}

func CreateLibraryHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if err := db.CreateLibrary(r.FormValue("name"), r.FormValue("import_path"), r.FormValue("converted_path")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	libraries, err := db.GetLibraries()
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Hx-Trigger", `closeModal`)
	if err = executeTemplate(w, "libraries", libraries); err != nil {
		log.Println(err)
	}
}

func UpdateLibraryHandlerFunc(w http.ResponseWriter, r *http.Request) {}
