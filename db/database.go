package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func LoginUser(user, password string) (username string, loggedIn bool) {
	var u, pwhash string
	stmt, _ := db.Prepare(`SELECT name, pw FROM users WHERE name = ?`)
	err := stmt.QueryRow(user).Scan(&u, &pwhash)
	if err != nil {
		log.Println(err)
		return user, false
	}
	if CheckPasswordHash(password, pwhash) {
		loggedIn = true
	}
	log.Printf("USER: %s", u)
	return u, loggedIn
}

//	func CreateUser(user, pw string) (string, bool) {
//		log.Println("TRYING TO CREATE")
//		stmt, err := db.Prepare(`INSERT INTO users (name, pw) VALUES (?, ?)`)
//		if err != nil {
//			log.Println(err)
//		}
//		log.Println(stmt)
//		stmt.Exec(user, pw)
//		log.Println("Inserted")
//		//name, logged := LoginUser(user)
//		//return name, logged
//	}
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
