package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
)

func MakeRouter(g GofemartInterface, a AuthorizationInterface) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(GzipHandle)

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", RegisterPostHandler(&a))
		// r.Post("/login", LoginPostHandler(a))
		// r.Post("/orders", OrdersPostPostHandler(g, a))
		// r.Get("/orders", RegisterGetHandler(g, a))
		// r.Route("/balance", func(r chi.Router) {
		// 	r.Get("/", BalanceGetHandler(g, a))
		// 	r.Post("/withdraw", WithdrawGetHandler(g, a))
		// 	r.Get("/withdrawals", WithdrawalsGetHandler(g, a))
		// })
	})

	return r
}

func RegisterPostHandler(a *AuthorizationInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("got register request")
		w.Header().Set("content-type", "application/json")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("error while read body: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		message := objects.RegisterMessage{}
		if err := json.Unmarshal(body, &message); err != nil {
			log.Println("error while unmarshal: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if message.Login == "" {
			log.Println("empty login")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if message.Password == "" {
			log.Println("empty password")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = (*a).Regist(message)
		if err != nil {
			log.Println("error while registe: " + err.Error())
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})
}

func LoginPostHandler(a *AuthorizationInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func OrdersPostPostHandler(g *GofemartInterface, a *AuthorizationInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func RegisterGetHandler(g *GofemartInterface, a *AuthorizationInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func BalanceGetHandler(g *GofemartInterface, a *AuthorizationInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func WithdrawGetHandler(g *GofemartInterface, a *AuthorizationInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func WithdrawalsGetHandler(g *GofemartInterface, a *AuthorizationInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}
