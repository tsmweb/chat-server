package profile

import (
	"encoding/json"
	"errors"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/cerror"
	ctlr "github.com/tsmweb/go-helper-api/controller"
	"log"
	"net/http"
)

// Controller provides the end point for the routers.
type Controller interface {
	Get() http.Handler
	Create() http.Handler
	Update() http.Handler
}

type controller struct {
	*ctlr.Controller
	getUseCase GetUseCase
	createUseCase CreateUseCase
	updateUseCase UpdateUseCase
}

// NewController creates a new instance of Controller.
func NewController(
	jwt auth.JWT,
	getUseCase GetUseCase,
	createUseCase CreateUseCase,
	updateUseCase UpdateUseCase) Controller {

	return &controller{
		ctlr.NewController(jwt),
		getUseCase,
		createUseCase,
		updateUseCase,
	}
}

// Get a profile by ID.
func (c *controller) Get() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ID, err := c.ExtractID(r)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		p, err := c.getUseCase.Execute(ID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, ErrProfileNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		vm := ViewModel{}
		vm.FromEntity(p)

		c.RespondWithJSON(w, http.StatusOK, vm)
	})
}

// Create a new profile.
func (c *controller) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !c.HasContentType(r, ctlr.MimeApplicationJSON) {
			c.RespondWithError(w, http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType))
			return
		}

		input := ViewModel{}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		err = c.createUseCase.Execute(input.ID, input.Name, input.LastName, input.Password)
		if err != nil {
			log.Println(err.Error())

			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				c.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			if errors.Is(err, cerror.ErrRecordAlreadyRegistered) {
				c.RespondWithError(w, http.StatusConflict, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

// Update updates profile data.
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
		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		//checks if the ID owns the data
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

			if errors.Is(err, ErrProfileNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}