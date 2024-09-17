package main_test

import (
	"bytes"
	"os"
	"testing"
	template_ "text/template"
)

func TestWriteJS(t *testing.T) {
	tp, err := template_.ParseFiles("web/js/htmx.min.js")
	if err != nil {
		t.Error(err)
	}
	b := &bytes.Buffer{}
	if err := tp.Execute(b, "htmx.min.js"); err != nil {
		t.Error(err)
	}
	if b.Len() == 0 {
		t.Error("buffer not written")
	}
	t.Log(b.Len())

	bb, err := os.ReadFile("web/js/htmx.min.js")
	if err != nil {
		t.Error(err)
	}
	if len(bb) == 0 {
		t.Error("cannot read file")
	} else {
		t.Log("file not empty", len(bb))
	}
}
