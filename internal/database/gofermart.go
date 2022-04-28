package database

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
)

func (gdb *GofermartDB) GetOrderIfExist(order string) (objects.OrderRow, bool) {
	orderRow := objects.OrderRow{}
	row := gdb.DB.QueryRowContext(gdb.Ctx, SelectOrderByNumber, order)
	var ts int64
	err := row.Scan(&orderRow.Number, &orderRow.UserID, &orderRow.Status, &orderRow.Accural, &ts)
	orderRow.UploudedAt = time.Unix(ts, 0)
	if err != nil {
		log.Println("Error while scan select from orders table: " + err.Error())
		return orderRow, false
	}
	return orderRow, true
}

func (gdb *GofermartDB) InsertOrder(orderRow objects.OrderRow) error {
	_, err := gdb.DB.ExecContext(gdb.Ctx, InsertOrder, orderRow.Number, orderRow.UserID, orderRow.Status, orderRow.Accural, orderRow.UploudedAt.Unix())
	if err != nil {
		log.Println("Error while insert to orders table " + err.Error())
		return errors.New("GofermartDB.InsertOrder: " + err.Error())
	}
	return nil
}

func (gdb *GofermartDB) GetOrders(user string) ([]objects.OrderRow, error) {
	stmt, err := gdb.DB.PrepareContext(gdb.Ctx, SelectOrderByUser)
	if err != nil {
		log.Printf("Error %s while preparing SQL statement\n", err)
		return []objects.OrderRow{}, fmt.Errorf("error %s while preparing sql statement", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(gdb.Ctx, user)
	if err != nil {
		log.Printf("Error %s while exec SQL statement\n", err)
		return []objects.OrderRow{}, fmt.Errorf("error %s while exec sql statement", err)
	}
	defer rows.Close()

	orderRows := []objects.OrderRow{}
	for rows.Next() {
		orderRow := objects.OrderRow{}
		var ts int64
		if err := rows.Scan(&orderRow.Number, &orderRow.UserID, &orderRow.Status, &orderRow.Accural, &ts); err != nil {
			orderRow.UploudedAt = time.Unix(ts, 0)
			log.Printf("Error %s while scan rows\n", err)
			return []objects.OrderRow{}, fmt.Errorf("rrror %s while scan rows", err)
		}
		orderRows = append(orderRows, orderRow)
	}

	return orderRows, nil
}

func (gdb *GofermartDB) GetBalance(user string) (objects.BalanceRow, error) {
	balanceRow := objects.BalanceRow{}
	row := gdb.DB.QueryRowContext(gdb.Ctx, SelectBalanceByUser, user)
	err := row.Scan(&balanceRow.UserID, &balanceRow.Current, &balanceRow.Withdraw)
	if err != nil {
		log.Println("No balance row in table: " + err.Error())
		gdb.UpdateBalance(objects.BalanceRow{UserID: user, Current: 0, Withdraw: 0})
		return objects.BalanceRow{UserID: user, Current: 0, Withdraw: 0}, nil
	}
	return balanceRow, nil
}

func (gdb *GofermartDB) InsertWithdraw(withdrawRow objects.WithdrawRow) error {
	_, err := gdb.DB.ExecContext(gdb.Ctx, InsertWithdraw, withdrawRow.Order, withdrawRow.UserID, withdrawRow.Sum, withdrawRow.ProcessedAt.Unix())
	if err != nil {
		log.Println("Error while insert to withdrawals table " + err.Error())
		return errors.New("GofermartDB.InsertWithdraw: " + err.Error())
	}
	return nil
}

func (gdb *GofermartDB) UpdateBalance(balanceRow objects.BalanceRow) error {
	_, err := gdb.DB.ExecContext(gdb.Ctx, UpsertBalance, balanceRow.UserID, balanceRow.Current, balanceRow.Withdraw)
	if err != nil {
		log.Println("Error while insert to balsnce table " + err.Error())
		return errors.New("GofermartDB.UpdateBalance: " + err.Error())
	}
	return nil
}

func (gdb *GofermartDB) UpdateOrders(orderRows []objects.OrderRow) error {
	tx, err := gdb.DB.Begin()
	if err != nil {
		log.Println("Transaxtion didnt started: " + err.Error())
		return errors.New("Transaxtion didnt started: " + err.Error())
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(gdb.Ctx, UsertOrders)
	if err != nil {
		log.Println("Context didnt prepared: " + err.Error())
		return errors.New("Context didnt prepared: " + err.Error())
	}
	defer stmt.Close()

	for _, orderRow := range orderRows {
		if _, err = stmt.ExecContext(gdb.Ctx, orderRow.Number, orderRow.UserID, orderRow.Status, orderRow.Accural, orderRow.UploudedAt.Unix()); err != nil {
			log.Println("Error while upsert in orders: " + err.Error())
			return errors.New("Error while upsert in orders: " + err.Error())
		}
	}

	if err = tx.Commit(); err != nil {
		log.Println("Error while commit upsert in orders: " + err.Error())
		return errors.New("Error while commit upsert in orders: " + err.Error())
	}

	return nil
}

func (gdb *GofermartDB) UpdateOrdersAndBalance(orderRows []objects.OrderRow, balanceRow objects.BalanceRow) error {
	tx, err := gdb.DB.Begin()
	if err != nil {
		log.Println("Transaxtion didnt started: " + err.Error())
		return errors.New("Transaxtion didnt started: " + err.Error())
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(gdb.Ctx, UsertOrders)
	if err != nil {
		log.Println("Context didnt prepared: " + err.Error())
		return errors.New("Context didnt prepared: " + err.Error())
	}
	defer stmt.Close()

	for _, orderRow := range orderRows {
		if _, err = stmt.ExecContext(gdb.Ctx, orderRow.Number, orderRow.UserID, orderRow.Status, orderRow.Accural, orderRow.UploudedAt.Unix()); err != nil {
			log.Println("Error while upsert in orders: " + err.Error())
			return errors.New("Error while upsert in orders: " + err.Error())
		}
	}

	_, err = gdb.DB.ExecContext(gdb.Ctx, UpsertBalance, balanceRow.UserID, balanceRow.Current, balanceRow.Withdraw)
	if err != nil {
		log.Println("Error while insert to balsnce table " + err.Error())
		return errors.New("GofermartDB.UpdateOrdersAndBalance: " + err.Error())
	}

	if err = tx.Commit(); err != nil {
		log.Println("Error while commit upsert in orders: " + err.Error())
		return errors.New("Error while commit upsert in orders: " + err.Error())
	}

	return nil
}

func (gdb *GofermartDB) GetWithdrawals(user string) ([]objects.WithdrawRow, error) {
	stmt, err := gdb.DB.PrepareContext(gdb.Ctx, SelectWithdrawalsByUser)
	if err != nil {
		log.Printf("Error %s while preparing SQL statement\n", err)
		return []objects.WithdrawRow{}, fmt.Errorf("error %s while preparing sql statement", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(gdb.Ctx, user)
	if err != nil {
		log.Printf("Error %s while exec SQL statement\n", err)
		return []objects.WithdrawRow{}, fmt.Errorf("error %s while exec sql statement", err)
	}
	defer rows.Close()

	withdrawRows := []objects.WithdrawRow{}
	for rows.Next() {
		withdrawRow := objects.WithdrawRow{}
		var ts int64
		if err := rows.Scan(&withdrawRow.Order, &withdrawRow.UserID, &withdrawRow.Sum, &ts); err != nil {
			withdrawRow.ProcessedAt = time.Unix(ts, 0)
			log.Printf("Error %s while scan rows\n", err)
			return []objects.WithdrawRow{}, fmt.Errorf("error %s while scan rows", err)
		}
		withdrawRows = append(withdrawRows, withdrawRow)
	}

	return withdrawRows, nil
}
