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
	resource = fmt.Sprintf("/%s/chat", version)
}

// Router for chat end points.
type Router struct {
	auth       middleware.Auth
	controller *Controller
}

// // NewRouter creates a router for Contact.
func NewRouter(auth middleware.Auth, ctrl *Controller) *Router {
	return &Router{
		auth: auth,
		controller: ctrl,
	}
}

// MakeRouters creates a router for chat.
func (r *Router) MakeRouters(mr *mux.Router) {
	// chat [GET]
	mr.Handle(resource, negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(r.controller.Connect())),
	).Methods(http.MethodGet)
}