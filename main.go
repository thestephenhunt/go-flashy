package main

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	database "github.com/thestephenhunt/go-server/db"
	"github.com/thestephenhunt/go-server/middleware"
	"github.com/thestephenhunt/go-server/models"
	"github.com/thestephenhunt/go-server/state"
	"github.com/thestephenhunt/go-server/utils"
)

var data = &state.Manager{
	Bg:         "",
	FirstTerm:  0,
	SecondTerm: 0,
	Operator:   "+",
	Correct:    false,
	Attempted:  false,
	User: models.User{
		Name:   "",
		Logged: false,
	},
}

func main() {
	pageStack := middleware.MiddlewareStack(
		middleware.RefreshCookie,
	)
	log.Println("Server is running...")
	fs := http.FileServer(http.Dir("static"))
	state.SetGState(data)

	database.DbInit()

	mux := http.NewServeMux()
	pageRouter := http.NewServeMux()
	pageRouter.HandleFunc("/try-again", utils.TryAgainHandler)
	pageRouter.HandleFunc("/check-answer", utils.CheckAnswerHandler)
	pageRouter.HandleFunc("/new-equation", utils.NewEquationHandler)
	pageRouter.HandleFunc("/login", utils.LoginHandler)
	pageRouter.HandleFunc("/go-login", utils.GoLoginHandler)
	pageRouter.HandleFunc("/go-register", utils.GoRegisterHandler)
	pageRouter.HandleFunc("/register", utils.RegisterHandler)
	pageRouter.HandleFunc("/answer", utils.AnswerHandler)
	pageRouter.HandleFunc("/", utils.IndexHandler)

	mux.Handle("/static/", http.StripPrefix("/static", fs))
	mux.Handle("/", pageStack(pageRouter))

	log.Fatal(http.ListenAndServe(":8000", mux))
}
