package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	database "github.com/thestephenhunt/go-server/db"
	"github.com/thestephenhunt/go-server/utils"
)

var data = utils.IndexData{
	Bg:         "",
	FirstTerm:  0,
	SecondTerm: 0,
	Operator:   "+",
	Correct:    false,
	Attempted:  false,
	User:       utils.User{},
}

func main() {
	fmt.Println("Server is running...")
	fs := http.FileServer(http.Dir("static"))
	appCtx := utils.NewIndexDataContext(data)

	database.DbInit()

	mux := http.NewServeMux()
	mux.HandleFunc("/try-again", appCtx.TryAgainHandler)
	mux.HandleFunc("/check-answer", appCtx.CheckAnswerHandler)
	mux.HandleFunc("/new-equation", appCtx.NewEquationHandler)
	mux.HandleFunc("/login", appCtx.LoginHandler)
	mux.HandleFunc("/answer", appCtx.AnswerHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", appCtx.IndexHandler)

	log.Fatal(http.ListenAndServe(":8000", mux))
}
