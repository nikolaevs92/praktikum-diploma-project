package rest

import (
	"context"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
)

type GofemartInterface interface {
	Run(context.Context)
	PushOrder(string, string) error
	GetOrders(string) ([]objects.Order, error)
	GetBalance(string) (objects.Balance, error)
	Withdraw(string, objects.Withdraw) error
	GetWithdrawals(string) ([]objects.Withdraw, error)
}

type AuthorizationInterface interface {
	Run(context.Context)
	Registration(objects.RegisterMessage) (objects.TokenMessage, error)
	Login(objects.LoginMessage) (objects.TokenMessage, error)
	GetUser(string) (string, error)
}
