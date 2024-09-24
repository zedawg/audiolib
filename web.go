package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	templateCache    *template.Template
	useTemplateCache = false
)

func runHTTP() {
	// staticFS, err := fs.Sub(embedFS, "embed/static")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))
	var err error
	templateCache, err = template.ParseFS(embedFS, "embed/templates/**")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	http.Handle("GET /static/", &FileHandler{})
	http.Handle("GET /{$}", &HomeHandler{})
	http.Handle("GET /activities", &ActivitiesHandler{})
	http.Handle("GET /settings", &SettingsHandler{})
	http.Handle("GET /user", &UserHandler{})
	http.Handle("POST /library", &CreateLibraryHandler{})
	http.Handle("POST /search", &SearchHandler{})

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Println(err)
	}
}

func executeTemplate(w http.ResponseWriter, name string, data any) error {
	if useTemplateCache && templateCache != nil {
		return templateCache.ExecuteTemplate(w, name, data)
	}
	t, err := template.ParseGlob("embed/templates/*")
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(w, name, data)
}

type FileHandler struct{}

func (h *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	p := filepath.Join("embed", r.URL.Path)
	http.ServeFile(w, r, p)
}

type HomeHandler struct{}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Name string
	}{
		Name: "home",
	}
	if err := executeTemplate(w, "pages.home", data); err != nil {
		fmt.Println(err)
	}
}

type ActivitiesHandler struct{}

func (h *ActivitiesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Name       string
		Activities []*ActivityEntry
	}{
		Name: "activities",
	}
	var err error
	data.Activities, err = db.GetActivities(100, 0)
	if err != nil {
		fmt.Println(err)
	}
	if err = executeTemplate(w, "pages.activities", data); err != nil {
		fmt.Println(err)
	}
}

type SettingsHandler struct{}

func (h *SettingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Name      string
		Libraries []*LibraryEntry
	}{
		Name: "settings",
	}
	var err error
	data.Libraries, err = db.GetLibraries()
	if err != nil {
		fmt.Println(err)
	}
	if err = executeTemplate(w, "pages.settings", data); err != nil {
		fmt.Println(err)
	}
}

type UserHandler struct{}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

type SearchHandler struct{}

func (h *SearchHandler) escape(q string) string {
	return q
}

func (h *SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Query         string
		SearchResults []*SearchResultEntry
	}{
		Query: h.escape(r.FormValue("search")),
	}
	var err error
	if len(data.Query) > 2 {
		data.SearchResults, err = db.Search(data.Query)
		if err != nil {
			fmt.Println(err)
		}
		for i := 0; i < 20; i++ {
			data.SearchResults = append(data.SearchResults, &SearchResultEntry{
				ID:      i,
				Name:    fmt.Sprintf("name %v", i),
				Type:    fmt.Sprintf("type %v", i),
				Details: fmt.Sprintf("details %v", i),
				Image:   "/static/nocover.jpg",
			})
		}
	}
	if err = executeTemplate(w, "search-results", data); err != nil {
		fmt.Println(err)
	}
}

type CreateLibraryHandler struct{}

func (h *CreateLibraryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := db.CreateLibrary(r.FormValue("name"), r.FormValue("import_path"), r.FormValue("converted_path")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	data := struct {
		Libraries []*LibraryEntry
	}{}
	var err error
	data.Libraries, err = db.GetLibraries()
	if err != nil {
		fmt.Println(err)
	}
	// w.Header().Set("Hx-Trigger", `{"closeModal": {}}`)
	w.Header().Set("Hx-Trigger", `closeModal`)
	if err = executeTemplate(w, "libraries", data.Libraries); err != nil {
		fmt.Println(err)
	}
}

func (h *CreateLibraryHandler) validateLibrary(r *http.Request) error {
	var (
		name          = r.FormValue("name")
		importPath    = r.FormValue("import_path")
		convertedPath = r.FormValue("converted_path")
	)
}
