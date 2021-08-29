package main

import (
	"hipeople_task/internal/api"
)

func main() {
	app := api.New()
	app.Start()
}

