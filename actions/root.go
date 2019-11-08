package actions

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	version = 1
)

func (a *Action) Root(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	a.logger.Printf("Received [%s] %s ...\n", r.Method, r.URL.Path)

	fmt.Fprintf(w, "Concourse Worker Manager: v%d\n", version)

	a.logger.Printf("Success\n")
}
