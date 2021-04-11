package group

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/urfave/negroni"
	"net/http"
)

const version string = "v1"

var (
	resource string
	resourceMember string
)

func init() {
	resource = fmt.Sprintf("/%s/group", version)
	resourceMember = fmt.Sprintf("/%s/group/member", version)
}

// Router for Contact end points.
type Router struct {
	auth middleware.Auth
	controller Controller
}

// NewRouter creates a router for Group.
func NewRouter(a middleware.Auth, c Controller) *Router {
	return &Router{
		auth: a,
		controller: c,
	}
}

// MakeRouters creates a router for Group.
func (r *Router) MakeRouters(mr *mux.Router) {
	// group/{id} [GET]
	mr.Handle(fmt.Sprintf("%s/{id}", resource), negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.Get())),
	).Methods(http.MethodGet)

	// group [GET]
	mr.Handle(resource, negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.GetAll())),
	).Methods(http.MethodGet)

	// group [POST]
	mr.Handle(resource, negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.Create())),
	).Methods(http.MethodPost)

	// group [PUT]
	mr.Handle(resource, negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.Update())),
	).Methods(http.MethodPut)

	// group/{id} [DELETE]
	mr.Handle(fmt.Sprintf("%s/{id}", resource), negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.Delete())),
	).Methods(http.MethodDelete)

	// group/member [POST]
	mr.Handle(resourceMember, negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.AddMember())),
	).Methods(http.MethodPost)

	// group/member/{group}/{user} [DELETE]
	mr.Handle(fmt.Sprintf("%s/{group}/{user}", resourceMember), negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.RemoveMember())),
	).Methods(http.MethodDelete)

	// group/member [PUT]
	mr.Handle(resourceMember, negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.SetAdmin())),
	).Methods(http.MethodPut)
}