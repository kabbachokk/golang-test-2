package app

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func check(hash, s []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, s)
	return err == nil
}

func (app *App) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			expectedUsernameHash := []byte(app.auth.username)
			expectedPasswordHash := []byte(app.auth.password)

			usernameMatch := check(expectedUsernameHash, []byte(username))
			passwordMatch := check(expectedPasswordHash, []byte(password))

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
