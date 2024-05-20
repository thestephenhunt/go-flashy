package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	database "github.com/thestephenhunt/go-server/db"
)

type Middleware func(http.Handler) http.Handler

type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func MiddlewareStack(ms ...Middleware) Middleware {
	return Middleware(func(next http.Handler) http.Handler {
		for i := len(ms) - 1; i >= 0; i-- {
			m := ms[i]
			next = m(next)
		}
		return next
	})
}

func RefreshCookie(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("ATTEMPTING TO REFRESH")
		wrapped := &wrappedResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		c, err := r.Cookie("flashy_token")
		if err != nil {
			log.Println("No cookie")
			wrapped.statusCode = http.StatusUnauthorized
			next.ServeHTTP(wrapped, r)
			return
		}
		sessionToken := c.Value
		// Check if there is a current session for the user
		userSession, exists := database.Sessions[sessionToken]
		if !exists {
			log.Println("Not in a session")

			next.ServeHTTP(wrapped, r)
			return
		}
		if userSession.IsExpired() {
			delete(database.Sessions, sessionToken)
			log.Println("Expired - unauthorized")
			wrapped.statusCode = http.StatusUnauthorized
			next.ServeHTTP(wrapped, r)
			return
		}

		newSessionToken := uuid.NewString()
		expiresAt := time.Now().Add(60 * time.Second)

		// Setup a session for the user
		database.Sessions[newSessionToken] = database.Session{
			Username: userSession.Username,
			Expiry:   expiresAt,
		}
		// Remove old session
		delete(database.Sessions, sessionToken)

		http.SetCookie(w, &http.Cookie{
			Name:    "flashy_token",
			Value:   newSessionToken,
			Expires: time.Now().Add(60 * time.Second),
		})
		log.Println(wrapped)
		next.ServeHTTP(wrapped, r)
	})
}
