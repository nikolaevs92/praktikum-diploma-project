package rest

import (
	"log"
	"net/http"
	"strings"
)

func GetAutification(a *AuthorizationInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := r.Cookie("auth-token")
			if err != nil {
				log.Println("error while get cookie: " + err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			reqToken := strings.TrimSpace(token.Value)
			user, err := (*a).GetUser(reqToken)
			if err != nil {
				log.Println("error while auth: " + err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			r.Header.Add("User", user)
			next.ServeHTTP(w, r)
		})
	}
}
