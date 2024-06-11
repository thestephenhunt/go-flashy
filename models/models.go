package models

import "time"

type User struct {
	Name     string
	Username string
	Password string
	Logged   bool
}

type Session struct {
	Sid          string
	User         string
	TimeAccessed time.Time
	State        *SessionState
	ExpiresAt    time.Time
	Next         *Session
}

type SessionState struct {
	Id         int64
	Username   string
	Digits     int
	FirstTerm  int
	SecondTerm int
	Bg         string
	Operator   string
	Correct    bool
	Logged     bool
}

type Equation struct {
	Digits     int
	FirstTerm  int
	SecondTerm int
	Operator   string
}

type CtxKey string
