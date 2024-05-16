package utils

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	database "github.com/thestephenhunt/go-server/db"
)

var tmpl = template.Must(template.ParseGlob("web/templates/*"))
var availableBgs = [6]string{"wavy-bg", "red-square", "pink-green-swirl", "color-squares", "circle-halves", "circle-braid"}
var bgPattern = func() string {
	return availableBgs[rand.Intn(len(availableBgs))]
}

func (ctx *IndexData) IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("INDEX")
	ctx.FirstTerm = NewTerm()
	ctx.SecondTerm = NewTerm()
	ctx.Bg = bgPattern()
	ctx.Operator = "+"
	log.Println(ctx)
	tmpl.ExecuteTemplate(w, "index", ctx)
}

func (ctx *IndexData) LoginHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("username")
	pw := r.FormValue("password")
	fmt.Println(user)
	fmt.Println(pw)
	name, logged := database.LoginUser(user, pw)
	ctx.User.Logged = logged
	ctx.User.Name = name
	log.Println(ctx)
	tmpl.ExecuteTemplate(w, "_login", ctx)
}

func (ctx *IndexData) CheckAnswerHandler(w http.ResponseWriter, r *http.Request) {
	answer, err := strconv.Atoi(ReverseString(r.FormValue("answer")))
	fmt.Println(answer)
	if err != nil {
		return
	}
	if Solve(ctx.FirstTerm, ctx.SecondTerm, ctx.Operator) == answer {
		ctx.Correct = true
	}
	log.Println(ctx.Correct)
}

func (ctx *IndexData) AnswerHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("ANSWER")
	fmt.Println(ctx.Correct)
	if !ctx.Correct {
		ctx.Attempted = true
	}
	tmpl.ExecuteTemplate(w, "_button", ctx)
	fmt.Println("RENDER NEW BUTTON")
}

func (ctx *IndexData) TryAgainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TRYING AGAIN")
	ctx.Attempted = false
	tmpl.ExecuteTemplate(w, "_button", ctx)
}

func (ctx *IndexData) NewEquationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("NEW EQUATION")
	ctx.FirstTerm = NewTerm()
	ctx.SecondTerm = NewTerm()
	ctx.Correct = false
	ctx.Attempted = false
	tmpl.ExecuteTemplate(w, "_flash-card", ctx)
}
