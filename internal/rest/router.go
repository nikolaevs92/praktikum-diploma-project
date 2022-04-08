package gofermart

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/nikolaevs92/practikumDiplomaProject/src/utils"
)

func MakeRouter(g *GofemartInterface) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(utils.GzipHandle)

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", RegisterPostHandler(g))
		r.Post("/login", LoginPostHandler(g))
		r.Post("/orders", OrdersPostPostHandler(g))
		r.Get("/orders", RegisterGetHandler(g))
		r.Route("/balance", func(r chi.Router) {
			r.Get("/", BalanceGetHandler(g))
			r.Post("/withdraw", WithdrawGetHandler(g))
			r.Get("/withdrawals", WithdrawalsGetHandler(g))
		})
	})

	return r
}

func RegisterPostHandler(g *GofemartInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func LoginPostHandler(g *GofemartInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func OrdersPostPostHandler(g *GofemartInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func RegisterGetHandler(g *GofemartInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func BalanceGetHandler(g *GofemartInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func WithdrawGetHandler(g *GofemartInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func WithdrawalsGetHandler(g *GofemartInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}
