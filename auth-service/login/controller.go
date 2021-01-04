package login

import (
	"encoding/json"
	"errors"
	"github.com/tsmweb/auth-service/profile"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/cerror"
	"log"
	"net/http"

	ctlr "github.com/tsmweb/go-helper-api/controller"
)

// Controller provides the end point for the routers.
type Controller interface {
	Login() http.Handler
	Update() http.Handler
}

type controller struct {
	*ctlr.Controller
	loginUseCase LoginUseCase
	updateUseCase UpdateUseCase
}

// NewController creates a new instance of Controller.
func NewController(
	jwt auth.JWT,
	loginUseCase LoginUseCase,
	updateUseCase UpdateUseCase) Controller {
	return &controller{
		ctlr.NewController(jwt),
		loginUseCase,
		updateUseCase,
	}
}

// Login returns a token if ID and password are valid.
func (c *controller) Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !c.HasContentType(r, ctlr.MimeApplicationJSON) {
			c.RespondWithError(w, http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType))
			return
		}

		input := ViewModel{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&input)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		token, err := c.loginUseCase.Execute(input.ID, input.Password)
		if err != nil {
			log.Println(err.Error())
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

		c.RespondWithJSON(w, http.StatusOK, &TokenAuth{Token: token})
	})
}

// Update updates password in data base.
func (c *controller) Update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !c.HasContentType(r, ctlr.MimeApplicationJSON) {
			c.RespondWithError(w, http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType))
			return
		}

		ID, err := c.ExtractID(r)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		input := ViewModel{}
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&input)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		// checks if the ID owns the data
		if input.ID != ID {
			c.RespondWithError(w, http.StatusUnauthorized, "You are not authorized to change the data")
			return
		}

		err = c.updateUseCase.Execute(input.ToEntity())
		if err != nil {
			log.Println(err.Error())
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
	})
}
