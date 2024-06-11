package state

import "github.com/thestephenhunt/go-server/models"

func Set(state *models.SessionState, key string, value interface{}) error {
	switch key {
	case "id":
		state.Id = value.(int64)
	case "username":
		state.Username = value.(string)
	case "FirstTerm":
		state.FirstTerm = value.(int)
	case "SecondTerm":
		state.SecondTerm = value.(int)
	case "Bg":
		state.Bg = value.(string)
	case "Operator":
		state.Operator = value.(string)
	}

	return nil
}

func Get(state *models.SessionState, key string) interface{} {
	var v interface{}
	switch key {
	case "id":
		v = state.Id
	case "username":
		v = state.Username
	case "FirstTerm":
		v = state.FirstTerm
	case "SecondTerm":
		v = state.SecondTerm
	case "Bg":
		v = state.Bg
	case "Operator":
		v = state.Operator
	}

	return v
}

func Delete(state *models.SessionState, key string) error {
	Set(state, key, nil)
	return nil
}
