package actions

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	version = 1
)

func Root(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Concourse Worker Manager: v%d\n", version)
}
