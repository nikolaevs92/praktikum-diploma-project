package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/accuralclient"
	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/authorization"
	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/config"
	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/database"
	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/gofermart"
	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/rest"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	adress := pflag.StringP("adress", "a", "", "")
	databaseURI := pflag.StringP("db-uri", "d", "", "")
	accuralSystemAdress := pflag.StringP("accural", "r", "", "")
	pflag.Parse()

	v := viper.New()
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	cfg := config.NewConfigsWithDefaults(v, *adress, *databaseURI, *accuralSystemAdress)
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-cancelChan
		cancel()
	}()

	accuralClient := accuralclient.New(cfg.Accural)
	gDB := database.New(cfg.DataBase)
	gm := gofermart.New(&gDB, &accuralClient, gofermart.Config{})
	auth := authorization.New(&gDB, authorization.Config{})

	api := rest.New(&gm, &auth, cfg.Rest)
	api.RunHTTPServer(ctx)
}
