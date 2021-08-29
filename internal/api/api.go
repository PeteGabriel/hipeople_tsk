package api

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

type App struct {
	addr *string
}

func New() *App {
	return &App{
		addr: flag.String("addr", ":4002", "http service address"),
	}
}

func (a App) configureRoutes() {
	//upload route
	http.Handle("/api/image", a.Upload())
	//get image route
	http.Handle("/api/image/", a.GetImage())
}

//TODO write an handler that redirects to upload or getImage handlers.
//keep it "local" because it has no interest outside this function.

func (a App) Start() {
	a.configureRoutes()

	err := http.ListenAndServe(*a.addr, nil)
	if err != nil {
		log.Fatal("error initiating web server:", err)
	}
	fmt.Println("Listening at port 4002")
}