package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
)

type IndexData struct {
	Bg         string
	FirstTerm  int64
	SecondTerm int64
	Operator   string
}

var tmpl = template.Must(template.ParseGlob("web/templates/*"))
var availableBgs = [6]string{"wavy-bg", "red-square", "pink-green-swirl", "color-squares", "circle-halves", "circle-braid"}
var bgPattern = availableBgs[rand.Intn(len(availableBgs))]
var iData = IndexData{
	Bg:         bgPattern,
	FirstTerm:  NewTerm(),
	SecondTerm: NewTerm(),
	Operator:   "+",
}

func newTermHandler(w http.ResponseWriter, r *http.Request) {
	iData.FirstTerm = NewTerm()
	tmpl.ExecuteTemplate(w, "term", iData.FirstTerm)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("INDEX")
	fmt.Println(tmpl.DefinedTemplates())
	fmt.Println(iData)
	tmpl.ExecuteTemplate(w, "index", iData)
}
func main() {
	fmt.Println("Server is running...")

	fs := http.FileServer(http.Dir("static"))

	mux := http.NewServeMux()
	mux.HandleFunc("/new-term", newTermHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", indexHandler)

	log.Fatal(http.ListenAndServe(":8000", mux))
}
