package main

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	database "github.com/thestephenhunt/go-server/db"
	"github.com/thestephenhunt/go-server/middleware"
)

var Bro = &Broker{
	Clients:        make(map[chan string]bool),
	NewClients:     make(chan (chan string)),
	ClosingClients: make(chan (chan string)),
	Messages:       make(chan string),
}

func main() {

	Bro.Start()

	log.Println("Server is running...")
	fs := http.FileServer(http.Dir("static"))

	database.DbInit()

	mux := http.NewServeMux()
	mux.HandleFunc("/try-again", TryAgainHandler)
	mux.HandleFunc("/new-equation", NewEquationHandler)
	mux.HandleFunc("/login", LoginHandler)
	mux.HandleFunc("/logout", LogoutHandler)
	mux.HandleFunc("/go-login", GoLoginHandler)
	mux.HandleFunc("/go-register", GoRegisterHandler)
	mux.HandleFunc("/register", RegisterHandler)
	mux.HandleFunc("/answer", AnswerHandler)
	mux.HandleFunc("/", middleware.Chain(IndexHandler, middleware.SessionMiddleware()))

	mux.Handle("/events", Bro)
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	log.Fatal(http.ListenAndServe(":80", mux))
}
