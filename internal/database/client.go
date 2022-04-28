package database

import (
	"context"
	"database/sql"
	"log"
)

type GofermartDB struct {
	Cfg Config
	Ctx context.Context
	DB  *sql.DB
}

func (gdb *GofermartDB) Run(ctx context.Context) {
	gdb.Ctx = ctx

	db, err := sql.Open("postgres", gdb.Cfg.DataBaseDSN)
	gdb.DB = db
	if err != nil {
		log.Println("sql arent opened")
		log.Println(err)
		return
	}
	defer db.Close()

	// create tables
	_, err = gdb.DB.ExecContext(gdb.Ctx, CreateTableIfNotExist)
	if err != nil {
		log.Println("table arent created")
		log.Println(err)
		return
	}
	<-gdb.Ctx.Done()
}
