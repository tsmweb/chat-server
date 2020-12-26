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
	controller *Controller
}

// NewRoutes creates a router for Login.
func NewRoutes(controller *Controller, auth middleware.Auth) *Router {
	return &Router{auth, controller}
}

// Route sets the end points for the Login.
func (r *Router) Route(router *mux.Router) {
	// POST /login
	router.HandleFunc(resource, r.controller.Login).Methods("POST")

	// PUT /login
	router.Handle(resource, negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.HandlerFunc(r.controller.Update),
	)).Methods("PUT")
}
