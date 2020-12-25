package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/auth-service/helper/handler"
	"github.com/tsmweb/helper-go/cerror"
	"github.com/tsmweb/helper-go/middleware"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

const version string = "v1"

var resource string

func init() {
	resource = fmt.Sprintf("/%s/profile", version)
}

func getProfile(h *handler.Handler, c Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ID, err := h.ExtractID(r)
		if err != nil {
			log.Println(err.Error())
			h.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		p, err := c.Get(ID)
		if err != nil {
			if errors.Is(err, ErrProfileNotFound) {
				h.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			h.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		h.RespondWithJSON(w, http.StatusOK, p)
	})
}

func createProfile(h *handler.Handler, c Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !h.HasContentType(r, handler.MimeApplicationJSON) {
			h.RespondWithError(w, http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType))
			return
		}

		input := Presenter{}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			h.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		err = c.Create(input)
		if err != nil {
			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				h.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			if errors.Is(err, cerror.ErrRecordAlreadyRegistered) {
				h.RespondWithError(w, http.StatusConflict, err.Error())
				return
			}

			h.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

func updateProfile(h *handler.Handler, c Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !h.HasContentType(r, handler.MimeApplicationJSON) {
			h.RespondWithError(w, http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType))
			return
		}

		ID, err := h.ExtractID(r)
		if err != nil {
			log.Println(err.Error())
			h.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		input := Presenter{}
		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			h.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		//checks if the ID owns the data
		if input.ID != ID {
			h.RespondWithError(w, http.StatusUnauthorized, "You are not authorized to change the data")
			return
		}

		err = c.Update(input)
		if err != nil {
			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				h.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			if errors.Is(err, ErrProfileNotFound) {
				h.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			h.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// Router for Profile end points.
type Router struct {
	handler *handler.Handler
	auth *middleware.Auth
	controller Controller
}

// // NewRouter creates a router for Profile.
func NewRouter(h *handler.Handler, a *middleware.Auth, c Controller) *Router {
	return &Router{
		handler: h,
		auth: a,
		controller: c,
	}
}

// MakeRouters creates a router for Profile.
func (r *Router) MakeRouters(mr *mux.Router) {
	// GET /profile
	mr.Handle(resource, negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(getProfile(r.handler, r.controller))),
	).Methods("GET")

	// POST /profile
	mr.Handle(resource, createProfile(r.handler, r.controller)).Methods("POST")

	// PUT /profile
	mr.Handle(resource, negroni.New(
		negroni.HandlerFunc(r.auth.RequireTokenAuth),
		negroni.Wrap(updateProfile(r.handler, r.controller))),
	).Methods("PUT")
}