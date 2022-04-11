package rest

import (
	"context"
	"log"
	"net/http"
)

type RestAPI struct {
	Gofermart     GofemartInterface
	Authorization AuthorizationInterface
	Config
}

func (g *RestAPI) Init() {
}

func New(cfg Config) *RestAPI {
	server := new(RestAPI)
	server.Config = cfg
	server.Init()
	return server
}

func (g *RestAPI) RunHTTPServer(end context.Context) {
	r := MakeRouter(g.Gofermart, g.Authorization)

	server := &http.Server{
		Addr:    g.Server,
		Handler: r,
	}

	go func() {
		<-end.Done()
		log.Println("Shutting down the HTTP server...")
		if err := server.Shutdown(end); err != nil {
			panic(err)
		}
	}()

	log.Fatal(server.ListenAndServe())
}

func (g *RestAPI) Run(end context.Context) {
	log.Println("Server started")

	DBClientCtx, DBClientCancel := context.WithCancel(end)
	defer DBClientCancel()
	go g.Gofermart.Run(DBClientCtx)

	httpServerEndCtx, httpServerCancel := context.WithCancel(end)
	defer httpServerCancel()
	go g.RunHTTPServer(httpServerEndCtx)

	<-end.Done()
	log.Println("Server stoped")
}
