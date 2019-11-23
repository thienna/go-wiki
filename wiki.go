package main

import (
	"encoding/json"
	"fmt"
	"gowiki/data"
	"gowiki/tmpl"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var linkRegexp = regexp.MustCompile("\\[([a-zA-Z0-9]+)\\]")
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func loadPage(title string) (*data.Page, error) {
	filename := "data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &data.Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)

	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	escapedBody := []byte(template.HTMLEscapeString(string(p.Body)))
	p.DisplayBody = template.HTML(linkRegexp.ReplaceAllFunc(escapedBody, func(str []byte) []byte {
		matched := linkRegexp.FindStringSubmatch(string(str))
		out := []byte("<a href=\"/view/" + matched[1] + "\">" + matched[1] + "</a>")
		return out
	}))

	tmpl.RenderTemplate("view", w, p)
}

func editHandler(writer http.ResponseWriter, request *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &data.Page{Title: title}
	}

	tmpl.RenderTemplate("edit", writer, p)

}

func saveHandler(writer http.ResponseWriter, request *http.Request, title string) {
	body := request.FormValue("body")
	p := &data.Page{
		Title: title,
		Body:  []byte(body),
	}

	err := p.Save()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(writer, request, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

type secondJSON struct {
	Lol string
	Abc string
}

type testJSON struct {
	Int64String int64
	TestStruct secondJSON
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	//http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
	test := testJSON{Int64String: 64, TestStruct: secondJSON{Lol: "1", Abc: "AAAAA"}}

	data, err := json.Marshal(test)
//	var jsonBlob = []byte(`[
//	{"Name": "Platypus", "Order": "Monotremata"},
//	{"Name": "Quoll",    "Order": "Dasyuromorphia"}
//]`)
//	type Animal struct {
//		Name  string
//		Order string
//	}
//	var animals []Animal
//	err = json.Unmarshal(jsonBlob, &animals)
//	if err != nil {
//		fmt.Println("error:", err)
//	}
//	fmt.Printf("%+v", animals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8081", nil))
}
