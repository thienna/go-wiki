package data

import (
	"html/template"
	"io/ioutil"
)

type Page struct {
	Title string
	Body  []byte
	DisplayBody template.HTML
}

func (p *Page) Save() error {
	filename := "data/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}
