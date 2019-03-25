package main

import (
	logs "log"

	"github.com/Sultan-IH/wind/server"
)

func main() {
	s, err := server.NewServer()

	if err != nil {
		logs.Panicf("error starting server: %v", err)
	}
	s.Run()
}
