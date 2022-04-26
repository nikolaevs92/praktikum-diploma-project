package rest

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/statuserror"
)

type AuthorizationTest struct {
}

func (a AuthorizationTest) Registration(message objects.RegisterMessage) error {
	if message.Login == "exist_user" {
		return errors.New("user already exist")
	}
	return nil
}

func (a AuthorizationTest) Login(message objects.LoginMessage) (objects.TokenMessage, error) {
	if message.Login == "exist_user" && message.Password == "correctpassword" {
		return objects.TokenMessage{Token: "token"}, nil
	}
	return objects.TokenMessage{}, errors.New("user doesnt exist")
}

func (a AuthorizationTest) GetUser(token string) (string, error) {
	if token == "token" {
		return "existed_user", nil
	}
	if token == "token3" {
		return "existed_user3", nil
	}
	return "", errors.New("wrong token")
}

type GofemartTest struct {
}

func (g GofemartTest) Run(end context.Context) {
}

func (g GofemartTest) PushOrder(user string, order string) error {
	if user != "existed_user" {
		log.Println("incorrect client", user)
		return errors.New("incorrect user")
	}
	if order == "order200" {
		return statuserror.NewStatusError("", 200)
	}
	if order == "order202" {
		return statuserror.NewStatusError("", 202)
	}
	if order == "order400" {
		return statuserror.NewStatusError("", 400)
	}
	if order == "order401" {
		return statuserror.NewStatusError("", 401)
	}
	if order == "order409" {
		return statuserror.NewStatusError("", 409)
	}
	if order == "order422" {
		return statuserror.NewStatusError("", 422)
	}
	if order == "order500" {
		return statuserror.NewStatusError("", 500)
	}
	return nil
}

func (g GofemartTest) GetOrders(user string) ([]objects.Order, error) {
	if user == "existed_user" {
		return []objects.Order{{Number: "124", Status: objects.OrderStatusNew, Accural: 111.0, UploudedAt: time.Now()}}, nil
	}
	if user == "existed_user3" {
		return []objects.Order{}, nil
	}
	return nil, nil
}

func (g GofemartTest) GetBalance(user string) (objects.Balance, error) {
	if user == "existed_user" {
		return objects.Balance{Current: 134, Withdraw: 34}, nil
	}
	return objects.Balance{}, nil
}

func (g GofemartTest) Withdraw(user string, withdraw objects.Withdraw) error {
	if withdraw.Order != "134" {
		return statuserror.NewStatusError("not existed order", 422)
	}
	if withdraw.Sum > 100 {
		return statuserror.NewStatusError("not enough balance", 402)
	}
	return nil
}

func (g GofemartTest) GetWithdrawals(user string) ([]objects.Withdraw, error) {
	if user == "existed_user" {
		return []objects.Withdraw{{Order: "134", Sum: 40, ProcessedAt: time.Now()}}, nil
	}
	if user == "existed_user3" {
		return []objects.Withdraw{}, nil
	}
	return nil, nil
}
