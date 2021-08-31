package models

const ServerErr = "SERVER_ERROR"

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Name    string `json:"-"`
	Error   error  `json:"-"`
}
