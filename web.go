package main

import (
	"fmt"
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

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))
	http.Handle("/", &HomeHandler{})
	http.Handle("/logs", &LogsHandler{})
	http.Handle("/tasks", &TasksHandler{})
	http.Handle("/settings", &SettingsHandler{})
	http.Handle("/user", &UserHandler{})

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Println(err)
	}
}

type HomeHandler struct {
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := struct{}{}
	var err error
	err = webTemplate.ExecuteTemplate(w, "page-home.html.tpl", data)
	if err != nil {
		fmt.Println(err)
	}
}

type LogsHandler struct{}

func (h *LogsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Logs []*LogEntry
	}{}
	var err error
	data.Logs, err = db.GetLogs(100, 0)
	if err != nil {
		fmt.Println(err)
	}
	err = webTemplate.ExecuteTemplate(w, "page-logs.html.tpl", data)
	if err != nil {
		fmt.Println(err)
	}
}

type TasksHandler struct{}

func (h *TasksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Tasks []*TaskEntry
	}{}
	var err error
	data.Tasks, err = db.GetTasks(100, 0)
	if err != nil {
		fmt.Println(err)
	}
	err = webTemplate.ExecuteTemplate(w, "page-tasks.html.tpl", data)
	if err != nil {
		fmt.Println(err)
	}
}

type SettingsHandler struct{}

func (h *SettingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := struct{}{}
	var err error
	err = webTemplate.ExecuteTemplate(w, "page-settings.html.tpl", data)
	if err != nil {
		fmt.Println(err)
	}
}

type UserHandler struct{}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := struct{}{}
	var err error
	err = webTemplate.ExecuteTemplate(w, "page-user.html.tpl", data)
	if err != nil {
		fmt.Println(err)
	}
}
