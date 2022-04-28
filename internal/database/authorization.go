package database

import (
	"errors"
	"log"
)

func (gdb *GofermartDB) IsLoginExist(login string) bool {
	row := gdb.DB.QueryRowContext(gdb.Ctx, SelectUserByLogin, login)
	var userID string
	err := row.Scan(&userID)
	return err == nil
}

func (gdb *GofermartDB) CheckLoginPasswordHash(login string, paswordHash string) (bool, string) {
	row := gdb.DB.QueryRowContext(gdb.Ctx, SelectUserByLoginPasswordHash, login, paswordHash)
	var userID string
	err := row.Scan(&userID)
	return err == nil, userID
}

func (gdb *GofermartDB) CreateUser(login string, paswordHash string) error {
	// TODO smth with userID?
	_, err := gdb.DB.ExecContext(gdb.Ctx, InsertUser, login, login, paswordHash)
	if err != nil {
		log.Println("Error while insert to User table " + err.Error())
		return errors.New("GofermartDB.CreateUser: " + err.Error())
	}
	return nil
}
