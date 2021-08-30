package api

import (
	"flag"
	"fmt"
	"hipeople_task/pkg/services"
	"log"
	"net/http"
)

type App struct {
	imgService services.IImageService
	addr *string
}

func New() *App {
	return &App{
		imgService: services.NewImageService(),
		addr: flag.String("addr", ":4002", "http service address"),
	}
}

func (a App) configureRoutes() {
	//upload route
	http.Handle("/api/image", a.Upload())
	//get image route
	http.Handle("/api/image/", a.GetImage())
}

func (a App) Start() {
	a.configureRoutes()

	err := http.ListenAndServe(*a.addr, nil)
	if err != nil {
		log.Fatal("error initiating web server:", err)
	}
	fmt.Println("Listening at port 4002")
}