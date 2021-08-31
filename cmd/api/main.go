package main

import (
	"hipeople_task/internal/api"
	"hipeople_task/pkg/config"
)

func main() {
	app := api.New(config.New())
	app.Start()
}
