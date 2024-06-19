package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/thestephenhunt/go-server/internal/users"
	"github.com/thestephenhunt/go-server/models"
	"github.com/thestephenhunt/go-server/utils"
)

var tmpl = template.Must(template.ParseGlob("web/templates/*"))

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(models.CtxKey("user"))
	log.Println(user)
	var state models.SessionState

	if user == "guest" {
		state = models.SessionState{
			Username:   "guest",
			Digits:     2,
			FirstTerm:  utils.NewTerm(2),
			SecondTerm: utils.NewTerm(2),
			Bg:         utils.BgPattern(),
			Operator:   "+",
			Correct:    false,
			Logged:     false,
		}
	} else {
		state = models.SessionState{
			Username:   user.(string),
			Digits:     2,
			FirstTerm:  utils.NewTerm(2),
			SecondTerm: utils.NewTerm(2),
			Bg:         utils.BgPattern(),
			Operator:   "+",
			Correct:    false,
			Logged:     true,
		}
	}
	tmpl.ExecuteTemplate(w, "index", state)
	timer := time.NewTimer(1 * time.Second)
	go func() {
		<-timer.C
		stringState, err := json.Marshal(state)
		if err != nil {
			log.Println(err)
		}
		Bro.Messages <- string(stringState)
	}()
}

func GoLoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "login-form", nil)
}

func GoRegisterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "register-form", nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	loggedUser, err := users.LoginUser(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		go func() {
			Bro.Messages <- fmt.Sprintln(`{ "error" : "No user found." }`)
		}()
		noUser, err := template.New("nouser").Parse(`No user found.`)
		if err != nil {
			log.Println("DID NOT PARSE")
		}
		noUser.ExecuteTemplate(w, "nouser", nil)
		return
	}
	newToken := users.NewJwt(loggedUser.Username)
	http.SetCookie(w, &http.Cookie{
		Name:   "flashy_token",
		Value:  newToken,
		MaxAge: 0,
	})
	log.Println(loggedUser)
	tmpl.ExecuteTemplate(w, "_login", loggedUser)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	newUser, err := users.RegisterUser(r)
	if err != nil {
		log.Println(err)
		tmpl.ExecuteTemplate(w, "_login", err)
		return
	}
	newToken := users.NewJwt(newUser.Username)
	http.SetCookie(w, &http.Cookie{
		Name:  "flashy_token",
		Value: newToken,
	})
	tmpl.ExecuteTemplate(w, "_login", newUser)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	userExit, err := users.LogoutUser(r)
	if err != nil {
		log.Println(err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "flashy_token",
		Value:  "",
		MaxAge: -1,
	})
	log.Println(userExit.Logged)
	tmpl.ExecuteTemplate(w, "_login", userExit)
}

func AnswerHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("ANSWER")
	log.Println(r)
	firstTerm, err := strconv.Atoi(r.FormValue("firstTerm"))
	if err != nil {
		return
	}
	operator := r.FormValue("operator")
	secondTerm, err := strconv.Atoi(r.FormValue("secondTerm"))
	if err != nil {
		return
	}
	answer, err := strconv.Atoi(utils.ReverseString(r.FormValue("answer")))
	log.Println(answer)
	if err != nil {
		return
	}
	if utils.Solve(firstTerm, secondTerm, operator) == answer {
		log.Println("true")
		log.Println(operator)
		tmpl.ExecuteTemplate(w, "good_job", nil)
	} else {
		log.Println("false")
		tmpl.ExecuteTemplate(w, "try_again", nil)
	}

	log.Println("RENDER NEW BUTTON")
}

func TryAgainHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "_button", nil)
}

func NewEquationHandler(w http.ResponseWriter, r *http.Request) {
	operator := r.FormValue("operator")
	digits, _ := strconv.Atoi(r.FormValue("digits"))
	log.Println("NEW EQUATION")
	log.Println(operator)
	log.Println(digits)
	state := &models.Equation{
		Digits:     digits,
		FirstTerm:  utils.NewTerm(digits),
		SecondTerm: utils.NewTerm(digits),
		Operator:   operator,
	}
	tmpl.ExecuteTemplate(w, "_flash-card", state)
	timer := time.NewTimer(1 * time.Second)
	go func() {
		<-timer.C
		stringState, err := json.Marshal(state)
		if err != nil {
			log.Println(err)
		}
		Bro.Messages <- string(stringState)
	}()
}
