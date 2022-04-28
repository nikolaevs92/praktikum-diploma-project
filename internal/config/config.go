package config

import (
	"time"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/accuralclient"
	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/authorization"
	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/database"
	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/rest"
	"github.com/spf13/viper"
)

type Config struct {
	Accural       accuralclient.Config
	DataBase      database.Config
	Rest          rest.Config
	Authorization authorization.Config
}

const (
	envAdress              = "RUN_ADDRESS"
	envDataBaseURI         = "DATABASE_URI"
	envAccuralSystemAdress = "ACCRUAL_SYSTEM_ADDRESS"
)

func NewConfigsWithDefaults(
	v *viper.Viper, adress string, databaseURI string, accuralSystemAdress string) *Config {

	v.SetDefault(envAdress, adress)
	v.SetDefault(envDataBaseURI, databaseURI)
	v.SetDefault(envAccuralSystemAdress, accuralSystemAdress)
	return &Config{
		Accural: accuralclient.Config{
			AccuralHost: v.GetString(envAccuralSystemAdress),
			Retries:     2,
			Timeout:     time.Second,
		},
		Rest:          rest.Config{Server: v.GetString(envAdress)},
		DataBase:      database.Config{DataBaseDSN: v.GetString(envDataBaseURI)},
		Authorization: authorization.Config{},
	}
}
