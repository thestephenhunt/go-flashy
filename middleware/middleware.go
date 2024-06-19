package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	database "github.com/thestephenhunt/go-server/db"
	"github.com/thestephenhunt/go-server/internal/users"
	"github.com/thestephenhunt/go-server/models"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func SessionMiddleware() Middleware {

	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var ctx context.Context
			key := models.CtxKey("user")
			c, err := r.Cookie("flashy_token")
			if err != nil {
				log.Println(err)
				http.SetCookie(w, &http.Cookie{
					Name:    "flashy_token",
					Value:   "",
					Expires: time.Now().AddDate(0, 0, -1),
					MaxAge:  -1,
				})
				ctx = context.WithValue(r.Context(), key, "guest")
				newReq := r.Clone(ctx)
				f(w, newReq)
				return
			}

			token := c.Value
			tkn, err := users.CheckJwt(token)
			if err != nil {
				database.LogoutUser(tkn)
				http.SetCookie(w, &http.Cookie{
					Name:    "flashy_token",
					Value:   "",
					Expires: time.Now().AddDate(0, 0, -1),
					MaxAge:  -1,
				})
				ctx = context.WithValue(r.Context(), key, "guest")
				newReq := r.Clone(ctx)
				f(w, newReq)
			}

			if tkn != "" {
				newToken := users.NewJwt(tkn)
				http.SetCookie(w, &http.Cookie{
					Name:   "flashy_token",
					Value:  newToken,
					MaxAge: 0,
				})
				ctx = context.WithValue(r.Context(), key, tkn)
				newReq := r.Clone(ctx)
				f(w, newReq)
			}
		}
	}
}
