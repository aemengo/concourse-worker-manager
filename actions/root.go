package actions

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var (
	version = "0"
)

func (a *Action) Root(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	a.logger.Printf("Received [%s] %s ...\n", r.Method, r.URL.Path)

	fmt.Fprintf(w, "Concourse Worker Manager: v%s\n", version)

	a.logger.Printf("Success\n")
}
