package tmpl

import (
	"gowiki/data"
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))

func RenderTemplate(file string, writer http.ResponseWriter, p *data.Page) {
	err := templates.ExecuteTemplate(writer, file + ".html", p)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
