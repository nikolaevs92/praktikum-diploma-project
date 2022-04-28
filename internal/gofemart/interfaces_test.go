package gofermart

import (
	"context"
	"errors"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
)

type AccuralInterfaceTest struct {
	AccuralOrders []objects.AccuralOrder
}

func (a *AccuralInterfaceTest) AddOrder(accuralOrder objects.AccuralOrder) {
	a.AccuralOrders = append(a.AccuralOrders, accuralOrder)
}

func (a *AccuralInterfaceTest) UpdateStatus(order string, status string) {
	for i := range a.AccuralOrders {
		if a.AccuralOrders[i].Number == order {
			a.AccuralOrders[i].Status = status
		}
	}
}

func (a AccuralInterfaceTest) GetOrder(order string) (objects.AccuralOrder, error) {
	for _, accuralOrder := range a.AccuralOrders {
		if accuralOrder.Number == order {
			return accuralOrder, nil
		}
	}
	return objects.AccuralOrder{}, errors.New("no order")
}

type GofemartDBInterfaceTest struct {
	Err          bool
	OrderRows    []objects.OrderRow
	WithdrawRows []objects.WithdrawRow
	Users        []objects.User
	Balances     []objects.BalanceRow
}

func (g *GofemartDBInterfaceTest) Run(ctx context.Context) {
	g.Err = false
	g.Balances = make([]objects.BalanceRow, 0)
	g.WithdrawRows = make([]objects.WithdrawRow, 0)
	g.Users = make([]objects.User, 0)
	g.OrderRows = make([]objects.OrderRow, 0)
}

func (g *GofemartDBInterfaceTest) GetOrderIfExist(order string) (objects.OrderRow, bool) {
	for _, orderRow := range g.OrderRows {
		if orderRow.Number == order {
			return orderRow, true
		}
	}
	return objects.OrderRow{}, false
}

func (g *GofemartDBInterfaceTest) InsertOrder(orderRow objects.OrderRow) error {
	if g.Err {
		g.Err = false
		return errors.New("err")
	}
	g.OrderRows = append(g.OrderRows, orderRow)
	return nil
}

func (g *GofemartDBInterfaceTest) GetOrders(user string) ([]objects.OrderRow, error) {
	if g.Err {
		g.Err = false
		return []objects.OrderRow{}, errors.New("err")
	}
	resultOrders := make([]objects.OrderRow, 0)
	for _, orderRow := range g.OrderRows {
		if orderRow.UserID == user {
			resultOrders = append(resultOrders, orderRow)
		}
	}
	return resultOrders, nil
}

func (g *GofemartDBInterfaceTest) GetBalance(user string) (objects.BalanceRow, error) {
	if g.Err {
		g.Err = false
		return objects.BalanceRow{}, errors.New("err")
	}
	for _, balanceRow := range g.Balances {
		if balanceRow.UserID == user {
			return balanceRow, nil
		}
	}
	g.Balances = append(g.Balances, objects.BalanceRow{UserID: user, Withdraw: 0, Current: 0})
	return objects.BalanceRow{UserID: user, Withdraw: 0, Current: 0}, nil
}

func (g *GofemartDBInterfaceTest) UpdateOrders(ordersUpdate []objects.OrderRow) error {
	if g.Err {
		g.Err = false
		return errors.New("err")
	}
	for _, rowUpdate := range ordersUpdate {
		for i, row := range g.OrderRows {
			if rowUpdate.Number == row.Number {
				g.OrderRows[i] = rowUpdate
			}
		}
	}
	return nil
}

func (g *GofemartDBInterfaceTest) InsertWithdraw(withdraw objects.WithdrawRow) error {
	if g.Err {
		g.Err = false
		return errors.New("err")
	}
	g.WithdrawRows = append(g.WithdrawRows, withdraw)
	return nil
}

func (g *GofemartDBInterfaceTest) UpdateBalance(balanceUpdate objects.BalanceRow) error {
	if g.Err {
		g.Err = false
		return errors.New("err")
	}
	for i, balanceRow := range g.Balances {
		if balanceRow.UserID == balanceUpdate.UserID {
			g.Balances[i].Current = balanceUpdate.Current
			g.Balances[i].Withdraw = balanceUpdate.Withdraw
			return nil
		}
	}
	g.Balances = append(g.Balances, balanceUpdate)
	return nil
}

func (g *GofemartDBInterfaceTest) UpdateOrdersAndBalance(ordersUpdate []objects.OrderRow, balanceUpdate objects.BalanceRow) error {
	if g.Err {
		g.Err = false
		return errors.New("err")
	}
	err := g.UpdateBalance(balanceUpdate)
	if err != nil {
		return err
	}
	return g.UpdateOrders(ordersUpdate)
}

func (g *GofemartDBInterfaceTest) GetWithdrawals(user string) ([]objects.WithdrawRow, error) {
	if g.Err {
		g.Err = false
		return []objects.WithdrawRow{}, errors.New("err")
	}
	resultWithdrawals := make([]objects.WithdrawRow, 0)
	for _, withdrawRow := range g.WithdrawRows {
		if withdrawRow.UserID == user {
			resultWithdrawals = append(resultWithdrawals, withdrawRow)
		}
	}
	return resultWithdrawals, nil
}
