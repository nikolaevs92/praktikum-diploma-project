package gofermart

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/statuserror"
	"github.com/stretchr/testify/assert"
)

func TestGofermart(t *testing.T) {
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-cancelChan
		cancel()
	}()

	t.Run("test_push_order", func(t *testing.T) {
		accural := &AccuralInterfaceTest{AccuralOrders: []objects.AccuralOrder{}}
		gDB := &GofemartDBInterfaceTest{}
		gofemart := Gofemart{DB: gDB, Accural: accural}
		go gofemart.Run(ctx)

		err := gofemart.PushOrder("", "order")
		assert.Equal(t, 400, getStatus(err))
		err = gofemart.PushOrder("user", "")
		assert.Equal(t, 400, getStatus(err))
		err = gofemart.PushOrder("", "")
		assert.Equal(t, 400, getStatus(err))

		err = gofemart.PushOrder("user", "order")
		assert.Equal(t, 500, getStatus(err))

		accural.AddOrder(objects.AccuralOrder{Number: "order", Status: objects.OrderStatusNew, Accural: 10})
		err = gofemart.PushOrder("user", "order")
		assert.Equal(t, 202, getStatus(err))
		err = gofemart.PushOrder("user", "order")
		assert.Equal(t, 200, getStatus(err))

		err = gofemart.PushOrder("user2", "order")
		assert.Equal(t, 409, getStatus(err))

		accural.AddOrder(objects.AccuralOrder{Number: "order2", Status: objects.OrderStatusInvalid, Accural: 10})
		err = gofemart.PushOrder("user2", "order2")
		assert.Equal(t, 422, getStatus(err))

		gDB.Err = true
		accural.AddOrder(objects.AccuralOrder{Number: "order3", Status: objects.OrderStatusNew, Accural: 10})
		err = gofemart.PushOrder("user", "order3")
		assert.Equal(t, 500, getStatus(err))

		//test_get_orders
		accural = &AccuralInterfaceTest{AccuralOrders: []objects.AccuralOrder{}}
		gDB = &GofemartDBInterfaceTest{}
		gofemart = Gofemart{DB: gDB, Accural: accural}
		go gofemart.Run(ctx)

		orders, err := gofemart.GetOrders("user")
		assert.NoError(t, err)
		assert.Equal(t, 0, len(orders))

		accural.AddOrder(objects.AccuralOrder{Number: "order", Status: objects.OrderStatusNew, Accural: 10})
		gofemart.PushOrder("user", "order")
		orders, err = gofemart.GetOrders("user")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(orders))
		assert.Equal(t, orders[0].Accural, 10.0)
		assert.Equal(t, orders[0].Number, "order")
		assert.Equal(t, orders[0].Status, objects.OrderStatusNew)

		accural.UpdateStatus("order", objects.OrderStatusProcessed)
		orders, err = gofemart.GetOrders("user")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(orders))
		assert.Equal(t, orders[0].Accural, 10.0)
		assert.Equal(t, orders[0].Number, "order")
		assert.Equal(t, orders[0].Status, objects.OrderStatusProcessed)

		accural.AddOrder(objects.AccuralOrder{Number: "order2", Status: objects.OrderStatusProcessing, Accural: 20})
		accural.AddOrder(objects.AccuralOrder{Number: "order3", Status: objects.OrderStatusNew, Accural: 20})
		gofemart.PushOrder("user", "order2")
		gofemart.PushOrder("user2", "order3")
		gofemart.PushOrder("user", "order3")
		orders, err = gofemart.GetOrders("user")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(orders))
		assert.Equal(t, orders[0].Accural, 10.0)
		assert.Equal(t, orders[0].Number, "order")
		assert.Equal(t, orders[0].Status, objects.OrderStatusProcessed)
		assert.Equal(t, orders[1].Accural, 20.0)
		assert.Equal(t, orders[1].Number, "order2")
		assert.Equal(t, orders[1].Status, objects.OrderStatusProcessing)

		//test_get_balance_and_withdraw
		accural = &AccuralInterfaceTest{AccuralOrders: []objects.AccuralOrder{}}
		gDB = &GofemartDBInterfaceTest{}
		gofemart = Gofemart{DB: gDB, Accural: accural}
		go gofemart.Run(ctx)

		balance, err := gofemart.GetBalance("user")
		assert.NoError(t, err)
		assert.Equal(t, 0.0, balance.Current)
		assert.Equal(t, 0.0, balance.Withdraw)

		gDB.Err = true
		_, err = gofemart.GetBalance("user")
		assert.Error(t, err)

		accural.AddOrder(objects.AccuralOrder{Number: "order", Status: objects.OrderStatusProcessing, Accural: 20})
		gofemart.PushOrder("user", "order")
		balance, err = gofemart.GetBalance("user")
		assert.NoError(t, err)
		assert.Equal(t, 0.0, balance.Current)
		assert.Equal(t, 0.0, balance.Withdraw)

		accural.UpdateStatus("order", objects.OrderStatusProcessed)
		balance, err = gofemart.GetBalance("user")
		assert.NoError(t, err)
		assert.Equal(t, 20.0, balance.Current)
		assert.Equal(t, 0.0, balance.Withdraw)

		accural.AddOrder(objects.AccuralOrder{Number: "order2", Status: objects.OrderStatusProcessed, Accural: 30})
		gofemart.PushOrder("user", "order2")
		balance, err = gofemart.GetBalance("user")
		assert.NoError(t, err)
		assert.Equal(t, 50.0, balance.Current)
		assert.Equal(t, 0.0, balance.Withdraw)

		err = gofemart.Withdraw("user", objects.Withdraw{Order: "order3", Sum: 60})
		assert.Error(t, err)
		err = gofemart.Withdraw("user", objects.Withdraw{Order: "order3", Sum: 50})
		assert.NoError(t, err)
		balance, err = gofemart.GetBalance("user")
		assert.NoError(t, err)
		assert.Equal(t, 0.0, balance.Current)
		assert.Equal(t, 50.0, balance.Withdraw)
		err = gofemart.Withdraw("user", objects.Withdraw{Order: "order5", Sum: 10})
		assert.Error(t, err)

		//test_add_get_withdrawals
		accural = &AccuralInterfaceTest{AccuralOrders: []objects.AccuralOrder{}}
		gDB = &GofemartDBInterfaceTest{}
		gofemart = Gofemart{DB: gDB, Accural: accural}
		go gofemart.Run(ctx)

		withdrawals, err := gofemart.GetWithdrawals("user")
		assert.NoError(t, err)
		assert.Equal(t, 0, len(withdrawals))
		gDB.Err = true
		_, err = gofemart.GetWithdrawals("user")
		assert.Error(t, err)

		accural.AddOrder(objects.AccuralOrder{Number: "order", Status: objects.OrderStatusProcessed, Accural: 100})
		gofemart.PushOrder("user", "order")
		err = gofemart.Withdraw("user", objects.Withdraw{Order: "order2", Sum: 30})
		assert.NoError(t, err)
		withdrawals, err = gofemart.GetWithdrawals("user")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(withdrawals))

		err = gofemart.Withdraw("user", objects.Withdraw{Order: "order3", Sum: 40})
		assert.NoError(t, err)
		withdrawals, err = gofemart.GetWithdrawals("user")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(withdrawals))

		err = gofemart.Withdraw("user", objects.Withdraw{Order: "order4", Sum: 400})
		assert.Error(t, err)
		withdrawals, err = gofemart.GetWithdrawals("user")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(withdrawals))

		err = gofemart.Withdraw("user", objects.Withdraw{Order: "order5", Sum: 4})
		assert.NoError(t, err)
		withdrawals, err = gofemart.GetWithdrawals("user")
		assert.NoError(t, err)
		assert.Equal(t, 3, len(withdrawals))

		err = gofemart.Withdraw("user", objects.Withdraw{Order: "order6", Sum: 40})
		assert.Error(t, err)
		withdrawals, err = gofemart.GetWithdrawals("user")
		assert.NoError(t, err)
		assert.Equal(t, 3, len(withdrawals))

		err = gofemart.Withdraw("user", objects.Withdraw{Order: "order7", Sum: 24})
		assert.NoError(t, err)
		withdrawals, err = gofemart.GetWithdrawals("user")
		assert.NoError(t, err)
		assert.Equal(t, 4, len(withdrawals))

		assert.Equal(t, "order2", withdrawals[0].Order)
		assert.Equal(t, 30.0, withdrawals[0].Sum)
		assert.Equal(t, "order3", withdrawals[1].Order)
		assert.Equal(t, 40.0, withdrawals[1].Sum)
		assert.Equal(t, "order5", withdrawals[2].Order)
		assert.Equal(t, 4.0, withdrawals[2].Sum)
		assert.Equal(t, "order7", withdrawals[3].Order)
		assert.Equal(t, 24.0, withdrawals[3].Sum)

		balance, err = gofemart.GetBalance("user")
		assert.NoError(t, err)
		assert.Equal(t, 2.0, balance.Current)
		assert.Equal(t, 98.0, balance.Withdraw)
	})

}

func getStatus(err error) int {
	status, ok := err.(statuserror.StatusError)
	if ok {
		return status.Status
	} else {
		return 0
	}
}
