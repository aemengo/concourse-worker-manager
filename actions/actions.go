package actions

import (
	"log"
	"os"
)

type Action struct {
	logger *log.Logger
}

func New() *Action {
	return &Action{
		logger: log.New(os.Stdout, "[ROUTING] ", log.LstdFlags),
	}
}
