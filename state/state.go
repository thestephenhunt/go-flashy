package state

import (
	"sync"

	"github.com/thestephenhunt/go-server/models"
)

type Manager struct {
	Bg         string
	FirstTerm  int
	SecondTerm int
	Operator   string
	Correct    bool
	Attempted  bool
	LoggedIn   bool
	User       models.User
}

var singleton *Manager
var once sync.Once

func SetGState(s ...*Manager) *Manager {
	once.Do(func() {
		if s != nil {
			singleton = s[0]
		}
		singleton = &Manager{}
	})
	return singleton
}

func (m *Manager) GetState() *Manager {
	return m
}

func (m *Manager) SetBg(bg string) {
	m.Bg = bg
}

func (m *Manager) SetFirstTerm(term int) {
	m.FirstTerm = term
}

func (m *Manager) SetSecondTerm(term int) {
	m.SecondTerm = term
}

func (m *Manager) SetOperator(op string) {
	m.Operator = op
}

func (m *Manager) SetCorrect(v bool) {
	m.Correct = v
}

func (m *Manager) SetAttempted(v bool) {
	m.Attempted = v
}

func (m *Manager) SetLoggedIn(v bool) {
	m.LoggedIn = v
}

func (m *Manager) SetUser(u models.User) {
	m.User = u
}
