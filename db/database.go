package database

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/thestephenhunt/go-server/models"
	"golang.org/x/crypto/bcrypt"
)

var Data *sql.DB

func getUser(u *models.User) (*models.User, error) {
	result := Data.QueryRow(`SELECT name, username, pw, logged FROM users WHERE username = ?`, u.Username)

	storedCreds := &models.User{}
	err := result.Scan(&storedCreds.Name, &storedCreds.Username, &storedCreds.Password, &storedCreds.Logged)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return storedCreds, nil
}

func FindUserByUsername(username string) (*models.User, error) {
	query, _ := Data.Prepare(`SELECT name, username, pw, logged FROM users WHERE username = ?`)
	storedCreds := &models.User{}
	result := query.QueryRow(username)
	err := result.Scan(&storedCreds.Name, &storedCreds.Username, &storedCreds.Password, &storedCreds.Logged)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return nil, err
		}
		log.Println(err)
		return nil, err
	}
	return storedCreds, nil
}

func LoginUser(u *models.User) (*models.User, error) {
	foundUser, err := getUser(u)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if !foundUser.Logged {
		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(u.Password))
		if err != nil {
			return nil, err
		}
		loggedIn := AddSession(foundUser)
		return loggedIn, nil
	}

	return foundUser, nil
}

func LogoutUser(username string) (*models.User, error) {
	currentUser, err := FindUserByUsername(username)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	updateQuery, _ := Data.Prepare(`UPDATE users SET logged = FALSE WHERE username = ? RETURNING name, username, pw, logged`)
	loggedOut := updateQuery.QueryRow(currentUser.Username)

	err = loggedOut.Scan(&currentUser.Name, &currentUser.Username, &currentUser.Password, &currentUser.Logged)
	if err != nil {
		log.Println(err)
	}
	return currentUser, nil
}

func CreateUser(u *models.User) (*models.User, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 10)

	query, _ := Data.Prepare(`INSERT INTO users(name, username, pw, logged) VALUES (?,?,?,?)`)
	_, err := query.Exec(u.Name, u.Username, string(hashedPassword), 1)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	foundUser, err := getUser(u)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return foundUser, nil
}

func AddSession(u *models.User) *models.User {
	queryStatement, _ := Data.Prepare(`UPDATE users SET logged = TRUE WHERE username = ? RETURNING name, username, pw, logged`)
	query := queryStatement.QueryRow(u.Username)

	err := query.Scan(&u.Name, &u.Username, &u.Password, &u.Logged)
	if err != nil {
		log.Println(err)
	}
	return u
}

func DeleteSession(s *models.Session) {
	_, err := Data.Exec(`DELETE FROM sessions WHERE sid = ?`, s.Sid)
	if err != nil {
		log.Println(err)
	}
}

func UpdateSession(s *models.Session) {
	stateJson, err := json.Marshal(s.State)
	if err != nil {
		log.Println(err)
	}
	nextJson, err := json.Marshal(s.Next)
	if err != nil {
		log.Println(err)
	}

	queryStatement := `UPDATE sessions SET timeAccessed = ?, state = ?, expiresAt = ?, next = ? WHERE sid = ?`

	query, _ := Data.Prepare(queryStatement)
	_, err = query.Exec(s.TimeAccessed, stateJson, s.ExpiresAt, nextJson, s.Sid)
	if err != nil {
		log.Println(err)
	}
}

func GetSession(sid string) (*models.Session, error) {
	result := Data.QueryRow(`SELECT sid, user, timeAccessed, state, expiresAt, next FROM sessions WHERE sid = ?`, sid)

	var state, next []byte
	var expires string
	storedCreds := &models.Session{}
	err := result.Scan(&storedCreds.Sid, &storedCreds.User, &storedCreds.TimeAccessed, &state, &expires, &next)
	json.Unmarshal([]byte(state), &storedCreds.State)
	json.Unmarshal([]byte(next), &storedCreds.Next)
	storedCreds.ExpiresAt, _ = time.Parse(time.RFC3339, expires)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return storedCreds, nil
}

func DbInit() error {
	var err error

	Data, err = sql.Open("sqlite3", "db/data.db")

	if err != nil {
		log.Println(err)
	}

	dataStmt, err := Data.Prepare(
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			username TEXT NOT NULL,
			pw TEXT,
			logged BOOLEAN, 
			created DATETIME DEFAULT CURRENT_TIMESTAMP
		)`)
	if err != nil {
		log.Println("Error in creating table")
	} else {
		log.Println("Successfully created table USERS")
	}
	dataStmt.Exec()

	sessStmt, err := Data.Prepare(
		`CREATE TABLE IF NOT EXISTS sessions (
			sid TEXT PRIMARY KEY,
			user TEXT,
			timeAccessed DATETIME NOT NULL,
			state JSON,
			expiresAt INTEGER,
			next JSON
		)`)
	if err != nil {
		log.Println("Error in creating table SESSIONS")
	} else {
		log.Println("Successfully created table SESSIONS")
	}
	sessStmt.Exec()

	return Data.Ping()
}
