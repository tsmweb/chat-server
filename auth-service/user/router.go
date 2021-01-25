package user

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
	resource = fmt.Sprintf("/%s/user", version)
}

// Router for User end points.
type Router struct {
	auth middleware.Auth
	controller Controller
}

// // NewRouter creates a router for User.
func NewRouter(a middleware.Auth, c Controller) *Router {
	return &Router{
		auth: a,
		controller: c,
	}
}

// MakeRouters creates a router for User.
func (r *Router) MakeRouters(mr *mux.Router) {
	// user [GET]
	mr.Handle(resource, negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.Get())),
	).Methods(http.MethodGet)

	// user [POST]
	mr.Handle(resource, r.controller.Create()).Methods(http.MethodPost)

	// user [PUT]
	mr.Handle(resource, negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.Update())),
	).Methods(http.MethodPut)
}