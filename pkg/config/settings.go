package config

import "flag"

type Settings struct {
	Host string
	Port string
}

func New() *Settings {

	host := flag.String("host", "127.0.0.1", "server hostname")
	port := flag.String("port", "4002", "server port")
	flag.Parse()

	return &Settings{
		Host: *host,
		Port: *port,
	}
}
