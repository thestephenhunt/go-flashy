package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type IndexData struct {
	Bg         string
	FirstTerm  int
	SecondTerm int
	Operator   string
	Correct    bool
}

var tmpl = template.Must(template.ParseGlob("web/templates/*"))
var availableBgs = [6]string{"wavy-bg", "red-square", "pink-green-swirl", "color-squares", "circle-halves", "circle-braid"}
var bgPattern = func() string {
	return availableBgs[rand.Intn(len(availableBgs))]
}
var iData = IndexData{
	Bg:         bgPattern(),
	FirstTerm:  NewTerm(),
	SecondTerm: NewTerm(),
	Operator:   "+",
	Correct:    false,
}

func reverseString(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}

func newEquationHandler(w http.ResponseWriter, r *http.Request) {
	iData.FirstTerm = NewTerm()
	iData.SecondTerm = NewTerm()
	tmpl.ExecuteTemplate(w, "_flash-card", iData)
}

func checkAnswerHandler(w http.ResponseWriter, r *http.Request) {
	answer, err := strconv.Atoi(reverseString(r.FormValue("answer")))
	fmt.Println(answer)
	if err != nil {
		return
	}
	if Solve(iData.FirstTerm, iData.SecondTerm, iData.Operator) == answer {
		iData.Correct = true
	}
}

func answerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ANSWER")
	fmt.Println(iData.Correct)
	tmpl.ExecuteTemplate(w, "_correct", nil)
	fmt.Println("EXECUTED")
}

func newTermHandler(w http.ResponseWriter, r *http.Request) {
	iData.FirstTerm = NewTerm()
	tmpl.ExecuteTemplate(w, "term", iData.FirstTerm)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("INDEX")
	iData.FirstTerm = NewTerm()
	iData.SecondTerm = NewTerm()
	iData.Bg = bgPattern()
	tmpl.ExecuteTemplate(w, "index", iData)
}
func main() {
	fmt.Println("Server is running...")

	fs := http.FileServer(http.Dir("static"))

	mux := http.NewServeMux()
	mux.HandleFunc("/new-term", newTermHandler)
	mux.HandleFunc("/check-answer", checkAnswerHandler)
	mux.HandleFunc("/new-equation", newEquationHandler)
	mux.HandleFunc("/answer", answerHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", indexHandler)

	log.Fatal(http.ListenAndServe(":8000", mux))
}
