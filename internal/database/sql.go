package database

const (
	CreateTableIfNotExist = "CREATE TABLE IF NOT EXISTS withdrawals ( orderID text, userID text, sum double precision, processedAt bigserial);" +
		"CREATE TABLE IF NOT EXISTS users ( userID text, login text, passwordHash text, CONSTRAINT uid_pk PRIMARY KEY (userID), CONSTRAINT id_uq UNIQUE (userID), CONSTRAINT login_uq UNIQUE (login));" +
		"CREATE TABLE IF NOT EXISTS orders ( number text, userID text, status text, accural double precision, uploadedAt bigserial, CONSTRAINT number_pk PRIMARY KEY (number), CONSTRAINT number_uq UNIQUE (number));" +
		"CREATE TABLE IF NOT EXISTS balance ( userID text, current double precision, withdraw double precision, CONSTRAINT cid_pk PRIMARY KEY (userID), CONSTRAINT id_uq UNIQUE (userID));"

	SelectUserByLogin             = "SELECT userId FROM users WHERE login = $1;"
	SelectUserByLoginPasswordHash = "SELECT userId FROM users WHERE login = $1 and passwordHash = $2;"
	InsertUser                    = "INSERT INTO users VALUES($1, $2, $3);"

	SelectOrderByNumber = "SELECT number, userID, status, accural, uploadedAt FROM orders WHERE number = $1;"
	SelectOrderByUser   = "SELECT number, userID, status, accural, uploadedAt FROM orders WHERE userID = $1;"
	InsertOrder         = "INSERT INTO orders VALUES($1, $2, $3, $4, $5);"
	UsertOrders         = "INSERT INTO orders VALUES($1, $2, $3, $4, $5) ON CONFLICT (number) DO UPDATE SET status = $3"

	SelectWithdrawalsByUser = "SELECT orderID, userID, sum, processedAt FROM withdrawals WHERE userID = $1;"
	InsertWithdraw          = "INSERT INTO withdrawals VALUES($1, $2, $3, $4);"

	SelectBalanceByUser = "SELECT userID, current, withdraw FROM balance WHERE userID = $1;"
	UpsertBalance       = "INSERT INTO balance VALUES($1, $2, $3) ON CONFLICT (userID) DO UPDATE SET current = $2, withdraw = $3;"
)
