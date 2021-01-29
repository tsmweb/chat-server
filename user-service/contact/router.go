package contact

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
	resource = fmt.Sprintf("/%s/contact", version)
}

// Router for Contact end points.
type Router struct {
	auth middleware.Auth
	controller Controller
}

// // NewRouter creates a router for Contact.
func NewRouter(a middleware.Auth, c Controller) *Router {
	return &Router{
		auth: a,
		controller: c,
	}
}

// MakeRouters creates a router for Contact.
func (r *Router) MakeRouters(mr *mux.Router) {
	// contact/{id} [GET]
	mr.Handle(fmt.Sprintf("%s/{id}", resource), negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.Get())),
	).Methods(http.MethodGet)

	// contact [GET]
	mr.Handle(resource, negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.GetAll())),
	).Methods(http.MethodGet)

	// contact/presence/{id} [GET]
	mr.Handle(fmt.Sprintf("%s/presence/{id}", resource), negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.GetPresence())),
	).Methods(http.MethodGet)

	// contact [POST]
	mr.Handle(resource, negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.Create())),
	).Methods(http.MethodPost)

	// contact [PUT]
	mr.Handle(resource, negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.Update())),
	).Methods(http.MethodPut)

	// contact/{id} [DELETE]
	mr.Handle(fmt.Sprintf("%s/{id}", resource), negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.Delete())),
	).Methods(http.MethodDelete)

	// contact/block [POST]
	mr.Handle(fmt.Sprintf("%s/block", resource), negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.Block())),
	).Methods(http.MethodPost)

	// contact/block/{id} [DELETE]
	mr.Handle(fmt.Sprintf("%s/block/{id}", resource), negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.Unblock())),
	).Methods(http.MethodDelete)
}