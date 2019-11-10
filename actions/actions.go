package actions

import (
	"log"
	"os"
)

type Action struct {
	logger    *log.Logger
	tsaHost   string
	workerTag string
}

func New(tsaHost string, workerTag string) *Action {
	return &Action{
		tsaHost:   tsaHost,
		workerTag: workerTag,
		logger:    log.New(os.Stdout, "[ROUTING] ", log.LstdFlags),
	}
}
