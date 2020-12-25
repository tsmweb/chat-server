package login

import (
	"encoding/json"
	"errors"
	"github.com/tsmweb/auth-service/profile"
	"github.com/tsmweb/helper-go/auth"
	"github.com/tsmweb/helper-go/cerror"
	"github.com/tsmweb/helper-go/controller"
	"net/http"
)

// Presenter provides the end point for the routers.
type Controller struct {
	*controller.Controller
	service Service
}

// NewRouter creates a new instance of Presenter.
func NewController(jwt *auth.JWT, service Service) *Controller {
	return &Controller{
		controller.NewController(jwt),
		service,
	}
}

// Login returns a token if ID and password are valid.
// Login end point for route .../login [POST method]
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	if !c.HasContentType(r, controller.MimeApplicationJSON) {
		c.RespondWithError(w, http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType))
		return
	}

	login := Login{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&login)
	if err != nil {
		c.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
		return
	}

	token, err := c.service.Login(login)
	if err != nil {
		var errValidateModel *cerror.ErrValidateModel
		if errors.As(err, &errValidateModel) {
			c.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, profile.ErrProfileNotFound) {
			c.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		if errors.Is(err, cerror.ErrUnauthorized) {
			c.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		c.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.RespondWithJSON(w, http.StatusOK, token)
}

// Update updates password in data base.
// Update end point for route .../login [PUT method]
func (c *Controller) Update(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !c.HasContentType(r, controller.MimeApplicationJSON) {
		c.RespondWithError(w, http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType))
		return
	}

	ID, err := c.ExtractID(r)
	if err != nil {
		c.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	login := Login{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&login)
	if err != nil {
		c.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
		return
	}

	//checks if the ID owns the data
	if login.ID != ID {
		c.RespondWithError(w, http.StatusUnauthorized, "You are not authorized to change the data")
		return
	}

	err = c.service.Update(login)
	if err != nil {
		var errValidateModel *cerror.ErrValidateModel
		if errors.As(err, &errValidateModel) {
			c.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, profile.ErrProfileNotFound) {
			c.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		c.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
