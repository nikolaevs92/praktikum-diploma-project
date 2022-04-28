package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/accuralclient"
	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/authorization"
	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/database"
	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/gofermart"
	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/rest"
)

func main() {
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-cancelChan
		cancel()
	}()

	accuralClient := accuralclient.New(accuralclient.GetDefaultConfig())
	gDB := database.New(database.GetDefaultConfig())
	gm := gofermart.New(&gDB, &accuralClient, gofermart.Config{})
	auth := authorization.New(&gDB, authorization.Config{})

	api := rest.New(&gm, &auth, rest.Config{})
	api.RunHTTPServer(ctx)
}
