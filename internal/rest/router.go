package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/statuserror"
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
		r.Post("/login", LoginPostHandler(&a))
		r.Route("/", func(r chi.Router) {
			r.Use(GetAutification(&a))
			r.Post("/orders", OrdersPostHandler(&g))
			r.Get("/orders", OrdersGetHandler(&g))
			r.Route("/balance", func(r chi.Router) {
				r.Get("/", BalanceGetHandler(&g))
				r.Post("/withdraw", WithdrawGetHandler(&g))
				r.Get("/withdrawals", WithdrawalsGetHandler(&g))
			})
		})
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

		token, err := (*a).Registration(message)
		if err != nil {
			log.Println("error while registe: " + err.Error())
			w.WriteHeader(http.StatusConflict)
			return
		}
		cookie := http.Cookie{Name: "auth-token", Value: token.Token}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})
}

func LoginPostHandler(a *AuthorizationInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("got register request")
		w.Header().Set("content-type", "application/json")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("error while read body: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		message := objects.LoginMessage{}
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

		token, err := (*a).Login(message)
		if err != nil {
			log.Println("error while login: " + err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		cookie := http.Cookie{Name: "auth-token", Value: token.Token}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})
}

func OrdersGetHandler(g *GofemartInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("got post order request")
		w.Header().Set("content-type", "application/json")

		userID := r.Header.Get("User")
		if userID == "" {
			log.Println("no user Id, use Autification middleware")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		orders, err := (*g).GetOrders(userID)
		if err != nil {
			log.Println("error while get orders: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(orders) == 0 {
			log.Println("dount have any orders")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		resp, err := json.Marshal(orders)
		if err != nil {
			log.Println("error while marshal orders: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	})
}

func OrdersPostHandler(g *GofemartInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("got get orders request")
		w.Header().Set("content-type", "application/text")

		userID := r.Header.Get("User")
		if userID == "" {
			log.Println("no user Id, use Autification middleware")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("error while read body: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		order := string(body)
		if order == "" {
			log.Println("empty body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = (*g).PushOrder(userID, order)
		if err != nil {
			status, ok := err.(statuserror.StatusError)
			if ok {
				w.WriteHeader(status.Status)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			log.Println("error while push order: " + err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})
}

func BalanceGetHandler(g *GofemartInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("got get balance request")
		w.Header().Set("content-type", "application/json")

		userID := r.Header.Get("User")
		if userID == "" {
			log.Println("no user Id, use Autification middleware")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		balance, err := (*g).GetBalance(userID)
		if err != nil {
			log.Println("error while getting balance: " + err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		resp, err := json.Marshal(balance)
		if err != nil {
			log.Println("error while marhal response: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	})
}

func WithdrawGetHandler(g *GofemartInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("got get balance request")
		w.Header().Set("content-type", "application/json")

		userID := r.Header.Get("User")
		if userID == "" {
			log.Println("no user Id, use Autification middleware")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("error while read body: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		withdraw := objects.Withdraw{}
		if err := json.Unmarshal(body, &withdraw); err != nil {
			log.Println("error while unmarshal: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = (*g).Withdraw(userID, withdraw)
		if err != nil {
			status, ok := err.(statuserror.StatusError)
			if ok {
				w.WriteHeader(status.Status)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			log.Println("error while push order: " + err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})
}

func WithdrawalsGetHandler(g *GofemartInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("got get withdrawals request")
		w.Header().Set("content-type", "application/json")

		userID := r.Header.Get("User")
		if userID == "" {
			log.Println("no user Id, use Autification middleware")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		withdrawals, err := (*g).GetWithdrawals(userID)
		if err != nil {
			log.Println("error while getting balance: " + err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if len(withdrawals) == 0 {
			log.Println("dount have any withdrawls")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		log.Println(withdrawals[0].Order, withdrawals[0].Sum)
		resp, err := json.Marshal(withdrawals)
		if err != nil {
			log.Println("error while marhal response: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	})
}
