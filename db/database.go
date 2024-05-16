package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func loginUser(user string) (username string, loggedIn bool) {
	stmt, _ := db.Prepare(`SELECT name FROM users WHERE name = 'user'`)
	stmt.Exec(user)
	username = user
	loggedIn = true
	//IData.User = user
	log.Printf("LOGGED IN %s", user)
	return username, loggedIn
}
func CreateUser(user, pw string) (string, bool) {
	log.Println("TRYING TO CREATE")
	stmt, err := db.Prepare(`INSERT INTO users (name, pw) VALUES (?, ?)`)
	if err != nil {
		log.Println(err)
	}
	log.Println(stmt)
	stmt.Exec(user, pw)
	log.Println("Inserted")
	name, logged := loginUser(user)
	return name, logged
}
func DbInit() error {
	var err error

	db, err = sql.Open("sqlite3", "db/data.db")
	if err != nil {
		log.Println(err)
	}
	log.Println(db)

	stmt, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			pw TEXT,
			created DATETIME DEFAULT CURRENT_TIMESTAMP
		)`)
	if err != nil {
		log.Println("Error in creating table")
	} else {
		log.Println("Successfully created table USERS")
	}
	stmt.Exec()

	return db.Ping()
}
