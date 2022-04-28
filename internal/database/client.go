package database

import (
	"context"
	"database/sql"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type GofermartDB struct {
	Cfg Config
	Ctx context.Context
	DB  *sql.DB
}

func New(config Config) GofermartDB {
	return GofermartDB{Cfg: config}
}

func (gdb *GofermartDB) DropTables() error {
	log.Println("This is function test db only")
	if !strings.Contains(gdb.Cfg.DataBaseDSN, "test") {
		panic("try drop tables in db without test in name")
	}
	_, err := gdb.DB.ExecContext(gdb.Ctx, "Drop table if exists orders; Drop table if exists users; Drop table if exists withdrawals; Drop table if exists balance;")
	if err != nil {
		log.Println("table arent created")
		log.Println(err)
		return err
	}
	_, err = gdb.DB.ExecContext(gdb.Ctx, CreateTableIfNotExist)
	if err != nil {
		log.Println("table arent created")
		log.Println(err)
		return err
	}
	return nil
}

func (gdb *GofermartDB) Run(ctx context.Context) {
	// Run DB
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
