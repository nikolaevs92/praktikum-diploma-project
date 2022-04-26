package rest

import (
	"log"
	"net/http"
	"strings"
)

func GetAutification(a *AuthorizationInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			splitToken := strings.Split(auth, "Bearer")
			if len(splitToken) != 2 {
				log.Println("didnt get auth token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			reqToken := strings.TrimSpace(splitToken[1])
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
