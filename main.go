package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aemengo/concourse-worker-manager/actions"
	"github.com/aemengo/concourse-worker-manager/config"
	"github.com/julienschmidt/httprouter"
)

var (
	logger        *log.Logger
	serverAddr    = ":8787"
	tsaHostFlag   = flag.String("tsa-host", "", "TSA Host (HOST:PORT)")
	workerTagFlag = flag.String("tag", "", "Worker Tag")
)

func main() {
	flag.Parse()

	logger = log.New(os.Stdout, "[MAIN] ", log.LstdFlags)

	err := validate()
	expectNoError(err)

	err = os.MkdirAll(config.Homedir(), os.ModePerm)
	expectNoError(err)

	router := httprouter.New()
	handler := actions.New(*tsaHostFlag, *workerTagFlag)

	router.GET("/", handler.Root)
	router.GET("/install/:version", handler.Install)

	logger.Printf("Initializing cwm server on %s...", serverAddr)
	err = http.ListenAndServe(serverAddr, router)
	expectNoError(err)
}

func validate() error {
	var errs []string

	if *tsaHostFlag == "" {
		errs = append(errs, "- the '-tsa-host' flag must the passed in. See usage")
	}

	if *workerTagFlag == "" {
		errs = append(errs, "- the '-tag' flag must the passed in. See usage")
	}

	if doesNotExist(config.WorkerKeyPath()) {
		errs = append(errs, fmt.Sprintf("- worker_key file must be present at %q", config.WorkerKeyPath()))
	}

	if doesNotExist(config.TsaHostKeyPath()) {
		errs = append(errs, fmt.Sprintf("- tsa_host_key file must be present at %q", config.TsaHostKeyPath()))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New("the following must be configured before running:\n" + strings.Join(errs, "\n"))
}

func doesNotExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

func expectNoError(err error) {
	if err != nil {
		logger.Printf("Failed to initialize: %s\n", err)
		os.Exit(1)
	}
}
