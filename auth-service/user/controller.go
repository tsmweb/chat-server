package user

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/tsmweb/auth-service/common"
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
	service Service
}

// NewController creates a new instance of Controller.
func NewController(
	jwt auth.JWT,
	service Service) Controller {

	return &controller{
		ctlr.NewController(jwt),
		service,
	}
}

// Get a user by ID.
func (c *controller) Get() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ID, err := c.ExtractID(r)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		user, err := c.service.Get(r.Context(), ID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, ErrUserNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		p := &Presenter{}
		p.FromEntity(user)

		c.RespondWithJSON(w, http.StatusOK, p)
	})
}

// Create a new user.
func (c *controller) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !c.HasContentType(r, ctlr.MimeApplicationJSON) {
			c.RespondWithError(w, http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType))
			return
		}

		input := &Presenter{}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		err = c.service.Create(r.Context(), input.ID, input.Name, input.LastName, input.Password)
		if err != nil {
			log.Println(err.Error())

			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				c.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			if errors.Is(err, ErrUserAlreadyExists) {
				c.RespondWithError(w, http.StatusConflict, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

// Update updates user data.
func (c *controller) Update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !c.HasContentType(r, ctlr.MimeApplicationJSON) {
			c.RespondWithError(w, http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType))
			return
		}

		userID, err := c.ExtractID(r)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		input := &Presenter{}
		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		ctx := context.WithValue(r.Context(), common.AuthContextKey, userID)

		err = c.service.Update(ctx, input.ToEntity())
		if err != nil {
			log.Println(err.Error())

			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				c.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			if errors.Is(err, ErrOperationNotAllowed) {
				c.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}

			if errors.Is(err, ErrUserNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}