package api

import (
	"fmt"
	"hipeople_task/pkg/config"
	"hipeople_task/pkg/services"
	"log"
	"net/http"
)

type App struct {
	imgService services.IImageService
	settings *config.Settings
}

func New(cfg *config.Settings) *App {
	return &App{
		imgService: services.New(),
		settings: cfg,
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

	err := http.ListenAndServe(":" + a.settings.Port, nil)
	if err != nil {
		log.Fatal("error initiating web server:", err)
	}
	fmt.Printf("%s %s\n", "Listening at port", a.settings.Port)
}