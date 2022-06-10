package gofermart

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/statuserror"
)

type Gofemart struct {
	DB      GofemartDBInterface
	Accural AccuralInterface
	Cfg     Config
}

func New(gDB GofemartDBInterface, accural AccuralInterface, config Config) Gofemart {
	return Gofemart{
		DB:      gDB,
		Accural: accural,
		Cfg:     config,
	}
}

func (g Gofemart) Run(end context.Context) {
	// Run Gofermart
	BDEndCtx, BDCancel := context.WithCancel(end)
	defer BDCancel()
	go g.DB.Run(BDEndCtx)
	<-end.Done()
}

func (g Gofemart) PushOrder(user string, orderID string) error {
	if orderID == "" || user == "" {
		log.Println("user and order must be not empty")
		return statuserror.NewStatusError("err", 400)
	}

	// check order in GofermartDB
	orderRow, exist := g.DB.GetOrderIfExist(orderID)
	if exist {
		if orderRow.UserID == user {
			log.Println("user already upload this order")
			return statuserror.NewStatusError("", 200)
		}
		log.Println("another user already upload this order")
		return statuserror.NewStatusError("another user already upload this order", 409)
	}

	// Start Procced new order
	err := g.DB.InsertOrder(objects.OrderRow{UserID: user, Number: orderID, UploudedAt: time.Now()})
	if err != nil {
		log.Println("Error while insert order: " + err.Error())
		return statuserror.NewStatusError("Error while insert order", 500)
	}
	return statuserror.NewStatusError("", 202)

}

func (g Gofemart) GetOrders(user string) ([]objects.Order, error) {
	err := g.UpdateUserOrdersAndBalance(user)
	// Maybe shold stop on the error
	if err != nil {
		log.Println("error while update orders: " + err.Error())
	}

	orderRows, err := g.DB.GetOrders(user)
	if err != nil {
		log.Println("error while get orders: " + err.Error())
		return []objects.Order{}, errors.New("error while get orders: " + err.Error())
	}
	orders := make([]objects.Order, len(orderRows))
	for i, orderRow := range orderRows {
		orders[i] = objects.Order{Number: orderRow.Number, Status: orderRow.Status, Accural: orderRow.Accural, UploudedAt: orderRow.UploudedAt}
	}
	return orders, nil
}

func (g Gofemart) UpdateUserOrdersAndBalance(user string) error {
	ordersRows, err := g.DB.GetOrders(user)
	if err != nil {
		log.Println("error while get orders: " + err.Error())
		return errors.New("error while get orders: " + err.Error())
	}

	balance, err := g.DB.GetBalance(user)
	if err != nil {
		// No balance its ok, will insert it in the end
		balance.Current = 0
		balance.UserID = user
		balance.Withdraw = 0
		log.Println("error while get balance from db: " + err.Error())
	}

	balanceChange := 0.0
	for i, orderRow := range ordersRows {
		accuralOrder, err := g.Accural.GetOrder(orderRow.Number)
		// Maybe shold stop on the error
		if err != nil {
			continue
		}
		orderRow.Accural = accuralOrder.Accural
		orderRow.Status = accuralOrder.Status
		if orderRow.Status != accuralOrder.Status && accuralOrder.Status == objects.OrderStatusProcessed {
			balanceChange += accuralOrder.Accural
		}
		ordersRows[i].Status = accuralOrder.Status
	}

	balance.Current += balanceChange
	// because one transaction
	err = g.DB.UpdateOrdersAndBalance(ordersRows, balance)
	if err != nil {
		log.Println("error while update orders and balance table: " + err.Error())
		return errors.New("error while update orders and balance table: " + err.Error())
	}

	return nil
}

func (g Gofemart) GetBalance(user string) (objects.Balance, error) {
	err := g.UpdateUserOrdersAndBalance(user)
	if err != nil {
		log.Println("error while update orders and balance: " + err.Error())
		return objects.Balance{}, errors.New("error while update orders and balance: " + err.Error())
	}

	balanceRow, err := g.DB.GetBalance(user)
	if err != nil {
		log.Println("error while get balance from db: " + err.Error())
		return objects.Balance{}, errors.New("error while get balance from db: " + err.Error())
	}

	return objects.Balance{Current: balanceRow.Current, Withdraw: balanceRow.Withdraw}, nil
}

func (g Gofemart) Withdraw(user string, withdraw objects.Withdraw) error {
	err := g.UpdateUserOrdersAndBalance(user)
	if err != nil {
		log.Println("error while update orders and balance: " + err.Error())
		return errors.New("error while update orders and balance: " + err.Error())
	}

	balanceRow, err := g.DB.GetBalance(user)
	if err != nil {
		log.Println("error while get balance from db: " + err.Error())
		return errors.New("error while get balance from db: " + err.Error())
	}

	if balanceRow.Current < withdraw.Sum {
		log.Printf("not enought on balance: %f to withdraw sum %f\n", balanceRow.Current, withdraw.Sum)
		return fmt.Errorf("not enought on balance: %f to withdraw sum %f", balanceRow.Current, withdraw.Sum)
	}

	err = g.DB.InsertWithdraw(objects.WithdrawRow{UserID: user, Order: withdraw.Order, Sum: withdraw.Sum, ProcessedAt: time.Now()})
	if err != nil {
		log.Println("error while insert withdraw: " + err.Error())
		return errors.New("error while insert withdraw: " + err.Error())
	}

	err = g.DB.UpdateBalance(objects.BalanceRow{UserID: user, Current: balanceRow.Current - withdraw.Sum, Withdraw: balanceRow.Withdraw + withdraw.Sum})
	if err != nil {
		// Not probelm if we insertt withdraw and dount change balance,
		log.Println("error while uda: " + err.Error())
		return errors.New("error while insert withdraw: " + err.Error())
	}

	return nil
}

func (g Gofemart) GetWithdrawals(user string) ([]objects.Withdraw, error) {
	withdrawalsRows, err := g.DB.GetWithdrawals(user)
	if err != nil {
		log.Println("error while get withdrawals: " + err.Error())
		return []objects.Withdraw{}, errors.New("error while get withdrawals: " + err.Error())
	}
	withdrawals := make([]objects.Withdraw, len(withdrawalsRows))
	for i, withdrawRow := range withdrawalsRows {
		withdrawals[i] = objects.Withdraw{Order: withdrawRow.Order, Sum: withdrawRow.Sum, ProcessedAt: withdrawRow.ProcessedAt}
	}
	return withdrawals, nil
}
