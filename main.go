package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aemengo/concourse-worker-manager/actions"
	"github.com/aemengo/concourse-worker-manager/config"
	"github.com/julienschmidt/httprouter"
)

var (
	logger     *log.Logger
	serverAddr = ":8787"
)

func main() {
	var (
		logger  = log.New(os.Stdout, "[MAIN] ", log.LstdFlags)
		homeDir = config.Homedir()
	)

	err := os.MkdirAll(homeDir, os.ModePerm)
	expectNoError(err)

	router := httprouter.New()
	handler := actions.New()
	router.GET("/", handler.Root)

	logger.Printf("Initializing cwm server on %s...", serverAddr)
	err = http.ListenAndServe(serverAddr, router)
	expectNoError(err)
}

func expectNoError(err error) {
	if err != nil {
		logger.Printf("Failed to initialize: %s\n", err)
		os.Exit(1)
	}
}
