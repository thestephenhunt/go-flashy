package utils

type IndexData struct {
	Bg         string
	FirstTerm  int
	SecondTerm int
	Operator   string
	Correct    bool
	Attempted  bool
	User       User
}

type User struct {
	name     string
	loggedIn bool
}

func NewIndexDataContext(data IndexData) *IndexData {
	return &IndexData{}
}
