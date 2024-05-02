package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
)

var tmpl = template.Must(template.ParseGlob("web/templates/*"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	availableBgs := [6]string{"wavy-bg", "red-square", "pink-green-swirl", "color-squares", "circle-halves", "circle-braid"}
	bgPattern := availableBgs[rand.Intn(len(availableBgs))]
	tmpl.ExecuteTemplate(w, "index", bgPattern)
}
func main() {
	fmt.Println("Server is running...")

	fs := http.FileServer(http.Dir("static"))

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", indexHandler)

	log.Fatal(http.ListenAndServe(":8000", mux))
}
