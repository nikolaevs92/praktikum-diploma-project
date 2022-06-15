package database

type Config struct {
	DataBaseDSN string
}

func GetDefaultConfig() Config {
	return Config{DataBaseDSN: "postgres://postgres:postgres@localhost:5439/postgres?sslmode=disable"}
}
