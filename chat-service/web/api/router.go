package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/urfave/negroni"
	"net/http"
)

const version string = "v1"

var resource string

func init() {
	resource = fmt.Sprintf("/%s/ws", version)
}

// Router for server end points.
type Router struct {
	auth   middleware.Auth
	handleWS http.Handler
}

// // NewRouter creates a router for Chat.
func NewRouter(a middleware.Auth, handleWS http.Handler) *Router {
	return &Router{
		auth:   a,
		handleWS: handleWS,
	}
}

// MakeRouter create a router for server.
func (r *Router) MakeRouters(mr *mux.Router) {
	// ws [GET]
	mr.Handle(resource, negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.handleWS)),
	).Methods(http.MethodGet)
}
