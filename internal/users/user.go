package users

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	database "github.com/thestephenhunt/go-server/db"
	"github.com/thestephenhunt/go-server/models"
)

func NewJwt(u string) string {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   u,
	}
	tokenSecret := os.Getenv("token")
	secretKey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("Token not generated: %s", err)
	}
	return s
}

func CheckJwt(u string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	secretKey := []byte(os.Getenv("SECRET"))
	claims := &jwt.RegisteredClaims{}

	_, err = jwt.ParseWithClaims(u, claims, func(tkn *jwt.Token) (any, error) {
		return secretKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("TOKEN MISMATCH")
			return "", err
		}
		log.Printf("PARSE ERROR: %s", err)
		user, _ := claims.GetSubject()
		return user, err
	}

	user, err := claims.GetSubject()
	if err != nil {
		log.Println(err)
		return "", err
	}

	return user, err
}

func RegisterUser(r *http.Request) (*models.User, error) {
	creds := models.User{}
	creds.Name = r.FormValue("name")
	creds.Username = r.FormValue("username")
	creds.Password = r.FormValue("password")
	newUser, err := database.CreateUser(&creds)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return newUser, nil
}

func LoginUser(r *http.Request) (*models.User, error) {
	creds := models.User{}
	creds.Username = r.FormValue("username")
	creds.Password = r.FormValue("password")
	logged, err := database.LoginUser(&creds)
	if err != nil {
		log.Println(err)
	}
	return logged, err
}

func LogoutUser(r *http.Request) (*models.User, error) {
	log.Println("LOGGING OUT IN USERS")
	c, err := r.Cookie("flashy_token")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Println("No cookie")
			return nil, err
		}
		return nil, err
	}

	tokenStr := c.Value
	user, err := CheckJwt(tokenStr)
	if err != nil {
		log.Println(err)
	}
	loggedOut, err := database.LogoutUser(user)
	if err != nil {
		return nil, err
	}
	return loggedOut, nil
}

// func RefreshToken(t string) (string, error) {
// 	user, err := CheckJwt(t)
// 	if err != nil {
// 		log.Println(err)
// 		database.LogoutUser(user)
// 		return "", err
// 	}
// 	currentUser, err := database.FindUserByUsername(user)
// 	if err != nil {
// 		return "", err
// 	}
// 	currentUser.Token = newJwt(currentUser.Username)
// 	currentUser = database.AddSession(currentUser)
// 	return currentUser.Username, nil
// }
