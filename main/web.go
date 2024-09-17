package main

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

func runHTTP() {
	var err error
	if staticFS, err = fs.Sub(embedFS, "embed/static"); err != nil {
		log.Fatal(err)
	}
	if webTemplate, err = template.ParseFS(embedFS, "embed/templates/**"); err != nil {
		log.Fatal(err)
	}
	if len(webTemplate.Templates()) == 0 {
		log.Fatal("web templates not found ")
	}

	http.Handle("/", &httpHandler{})
	http.Handle("/history", &HistoryHandler{})
	http.Handle("/settings", &SettingsHandler{})
	http.Handle("/user", &UserHandler{})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Println(err)
	}
}

type httpHandler struct {
}

type TemplateData struct {
	ToolbarItems []any
	LibraryItems []any
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{}
	webTemplate.ExecuteTemplate(w, "index.html.tpl", data)
}

type HistoryHandler struct{}

func (h *HistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{}
	webTemplate.ExecuteTemplate(w, "index.html.tpl", data)
}

type SettingsHandler struct{}

func (h *SettingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{}
	webTemplate.ExecuteTemplate(w, "index.html.tpl", data)
}

type UserHandler struct{}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{}
	webTemplate.ExecuteTemplate(w, "index.html.tpl", data)
}
