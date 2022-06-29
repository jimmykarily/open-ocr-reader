package controllers

import (
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"github.com/pkg/errors"
)

// RenderWithLayout renders a layout with an additional page layout
func RenderWithLayout(pageTemplate string, w http.ResponseWriter, data interface{}) error {
	var err error
	mainlayout, err := os.Open("views/main_layout.html")
	if err != nil {
		return errors.Wrap(err, "opening the main layout file")
	}
	page, err := os.Open("views/" + pageTemplate + ".html")
	if err != nil {
		return errors.Wrap(err, "opening the page layout file")
	}

	var tmpl *template.Template
	layoutContent, err := ioutil.ReadAll(mainlayout)
	if err != nil {
		return errors.Wrap(err, "reading the main layout file")
	}
	tmpl = template.New(pageTemplate)
	tmpl = tmpl.Delims("[[", "]]")
	tmpl.Parse(string(layoutContent))
	if err != nil {
		panic(err)
	}

	if pageTemplate != "" {
		pageContent, err := ioutil.ReadAll(page)
		if err != nil {
			return errors.Wrap(err, "reading the page layout file")
		}
		tmpl, err = tmpl.Parse(string(pageContent))
		if err != nil {
			return errors.Wrap(err, "parsing the page layout")
		}
	}

	//w.WriteHeader(http.StatusOK)
	tmpl.ExecuteTemplate(w, "main_layout", data)

	return nil
}

func Home(w http.ResponseWriter, r *http.Request) {
	err := RenderWithLayout("home", w, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
