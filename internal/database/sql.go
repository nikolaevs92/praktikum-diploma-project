package database

const (
	CreateTableIfNotExist = "CREATE TABLE IF NOT EXISTS users (userID text, login text, passwordHash test, CONSTRAINT id_pk PRIMARY KEY (userID), CONSTRAINT id_uq UNIQUE (userID), CONSTRAINT login_uq UNIQUE (login));"

	SelectUserByLogin             = "SELECT userId FROM users WHERE login = $1;"
	SelectUserByLoginPasswordHash = "SELECT userId FROM users WHERE login = $1 and passwordHash = $2;"

	InsertUser = "INSERT INTO users VALUES($1, $2, $3)"
)
