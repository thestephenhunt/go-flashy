package utils

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	database "github.com/thestephenhunt/go-server/db"
	"github.com/thestephenhunt/go-server/state"
)

var tmpl = template.Must(template.ParseGlob("web/templates/*"))
var availableBgs = [6]string{"wavy-bg", "red-square", "pink-green-swirl", "color-squares", "circle-halves", "circle-braid"}
var bgPattern = func() string {
	return availableBgs[rand.Intn(len(availableBgs))]
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("INDEX")
	state := state.SetGState().GetState()
	state.FirstTerm = NewTerm()
	state.SecondTerm = NewTerm()
	state.Bg = bgPattern()
	state.Operator = "+"
	log.Println(state)
	tmpl.ExecuteTemplate(w, "index", state)
}

func GoLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := state.SetGState().GetState()
	tmpl.ExecuteTemplate(w, "login-form", state)
}

func GoRegisterHandler(w http.ResponseWriter, r *http.Request) {
	state := state.SetGState().GetState()
	tmpl.ExecuteTemplate(w, "register-form", state)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	state := state.SetGState().GetState()
	loggedUser, err := database.LoginUser(w, r)
	if err != nil {
		log.Println(err)
	}
	state.User.Name = loggedUser.Name
	state.User.Logged = true
	log.Println(state)
	tmpl.ExecuteTemplate(w, "_login", state)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	state := state.SetGState().GetState()
	database.CreateUser(w, r)
	log.Println("attempted register")
	log.Println(state)

}

func CheckAnswerHandler(w http.ResponseWriter, r *http.Request) {
	state := state.SetGState().GetState()
	answer, err := strconv.Atoi(ReverseString(r.FormValue("answer")))
	log.Println(answer)
	if err != nil {
		return
	}
	if Solve(state.FirstTerm, state.SecondTerm, state.Operator) == answer {
		state.Correct = true
	}
	log.Println(state.Correct)
}

func AnswerHandler(w http.ResponseWriter, r *http.Request) {
	state := state.SetGState().GetState()
	log.Println("ANSWER")
	log.Println(state.Correct)
	if !state.Correct {
		state.Attempted = true
	}
	tmpl.ExecuteTemplate(w, "_button", state)
	log.Println("RENDER NEW BUTTON")
}

func TryAgainHandler(w http.ResponseWriter, r *http.Request) {
	state := state.SetGState().GetState()
	log.Println("TRYING AGAIN")
	state.Attempted = false
	tmpl.ExecuteTemplate(w, "_button", state)
}

func NewEquationHandler(w http.ResponseWriter, r *http.Request) {
	state := state.SetGState().GetState()
	log.Println("NEW EQUATION")
	state.FirstTerm = NewTerm()
	state.SecondTerm = NewTerm()
	state.Correct = false
	state.Attempted = false
	tmpl.ExecuteTemplate(w, "_flash-card", state)
}
