package config

import "flag"

type Settings struct {
	Port string
}

func New() *Settings {

	port := flag.String("port", "4002", "server port")
	flag.Parse()

	return &Settings{
		Port: *port,
	}
}
