package login

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/urfave/negroni"
)

const version string = "v1"

var resource string

func init() {
	resource = fmt.Sprintf("/%s/login", version)
}

// Router for Login end points.
type Router struct {
	auth middleware.Auth
	controller Controller
}

// NewRoutes creates a router for Login.
func NewRoutes(a middleware.Auth, c Controller) *Router {
	return &Router{
		auth: a,
		controller: c,
	}
}

// MakeRouters creates a router for Login.
func (r *Router) MakeRouters(mr *mux.Router) {
	// POST /login
	mr.Handle(resource, r.controller.Login()).Methods("POST")

	// PUT /login
	mr.Handle(resource, negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.Update()),
	)).Methods("PUT")
}
