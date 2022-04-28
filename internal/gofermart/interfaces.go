package gofermart

import (
	"context"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
)

type AccuralInterface interface {
	GetOrder(string) (objects.AccuralOrder, error)
}

type GofemartDBInterface interface {
	Run(context.Context)
	GetOrderIfExist(string) (objects.OrderRow, bool)
	InsertOrder(objects.OrderRow) error
	GetOrders(string) ([]objects.OrderRow, error)
	GetBalance(string) (objects.BalanceRow, error)
	UpdateOrders([]objects.OrderRow) error
	InsertWithdraw(objects.WithdrawRow) error
	UpdateBalance(objects.BalanceRow) error
	UpdateOrdersAndBalance([]objects.OrderRow, objects.BalanceRow) error
	GetWithdrawals(string) ([]objects.WithdrawRow, error)
}
