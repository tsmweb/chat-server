package profile

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/helper-go/middleware"
	"github.com/urfave/negroni"
)

const version string = "v1"

var resource string

func init() {
	resource = fmt.Sprintf("/%s/profile", version)
}

// Router for Profile end points.
type Router struct {
	auth *middleware.Auth
	controller Controller
}

// // NewRouter creates a router for Profile.
func NewRouter(a *middleware.Auth, c Controller) *Router {
	return &Router{
		auth: a,
		controller: c,
	}
}

// MakeRouters creates a router for Profile.
func (r *Router) MakeRouters(mr *mux.Router) {
	// GET /profile
	mr.Handle(resource, negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.Get())),
	).Methods("GET")

	// POST /profile
	mr.Handle(resource, r.controller.Create()).Methods("POST")

	// PUT /profile
	mr.Handle(resource, negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.Update())),
	).Methods("PUT")
}