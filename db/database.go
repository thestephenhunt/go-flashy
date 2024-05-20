package database

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/thestephenhunt/go-server/models"
	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	Username string
	Expiry   time.Time
}

var Sessions = map[string]Session{}

var Database *sql.DB

func LoginUser(w http.ResponseWriter, r *http.Request) (*models.User, error) {
	creds := &models.User{}
	creds.Username = r.FormValue("username")
	creds.Password = r.FormValue("password")

	result := Database.QueryRow(`SELECT name, username, pw FROM users WHERE username = ?`, creds.Username)

	storedCreds := &models.User{}
	err := result.Scan(&storedCreds.Name, &storedCreds.Username, &storedCreds.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return nil, err
		}

		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password))

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, err
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(60 * time.Second)

	Sessions[sessionToken] = Session{
		Username: storedCreds.Username,
		Expiry:   expiresAt,
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "flashy_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	return storedCreds, nil
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("TRYING TO CREATE")
	creds := &models.User{}
	creds.Name = r.FormValue("name")
	creds.Username = r.FormValue("username")
	creds.Password = r.FormValue("password")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(creds.Password), 10)

	stmt, err := Database.Exec(`INSERT INTO users(name, username, pw)VALUES(?,?,?)`, creds.Name, creds.Username, string(hashedPassword))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	log.Println(stmt)
	log.Println("Inserted")

	//return name, logged

}

func DbInit() error {
	var err error

	Database, err = sql.Open("sqlite3", "db/data.db")
	if err != nil {
		log.Println(err)
	}
	log.Println(Database)

	stmt, err := Database.Prepare(
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			username TEXT NOT NULL,
			pw TEXT,
			created DATETIME DEFAULT CURRENT_TIMESTAMP
		)`)
	if err != nil {
		log.Println("Error in creating table")
	} else {
		log.Println("Successfully created table USERS")
	}
	stmt.Exec()

	return Database.Ping()
}

func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}
