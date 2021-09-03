package api

import (
	"hipeople_task/internal/api/middleware"
	"hipeople_task/pkg/config"
	"hipeople_task/pkg/services"
	"log"
	"net/http"
)

type App struct {
	imgService services.IImageService
	settings   *config.Settings
}

func New(cfg *config.Settings) *App {
	return &App{
		imgService: services.New(),
		settings:   cfg,
	}
}

func (a App) configureRoutes() {
	//upload route
	http.Handle("/api/image", middleware.Log(a.Upload()))
	//get image route
	http.Handle("/api/image/", middleware.Log(a.GetImage()))
}

func (a App) Start() {
	a.configureRoutes()

	log.Printf("starting server at %s:%s\n", a.settings.Host, a.settings.Port)

	if err := http.ListenAndServe(":"+a.settings.Port, nil); err != nil {
		log.Fatal("error initiating web server:", err)
	}
}
