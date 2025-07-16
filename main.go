package main

import (
	"log/slog"
	"os"

	"github.com/fernandobalieirof/cachydb/internal/server"
)

func main() {
	s := server.NewServer(server.DefaultConfig)
	if err := s.Start(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
