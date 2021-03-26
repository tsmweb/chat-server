package contact

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/cerror"
	ctlr "github.com/tsmweb/go-helper-api/controller"
	"log"
	"net/http"
)

// Controller provides the end point for the routers.
type Controller interface {
	Get() http.Handler
	GetAll() http.Handler
	GetPresence() http.Handler
	Create() http.Handler
	Update() http.Handler
	Delete() http.Handler
	Block() http.Handler
	Unblock() http.Handler
}

type controller struct {
	*ctlr.Controller
	service Service
}

// NewController creates a new instance of Controller.
func NewController(jwt auth.JWT, service Service) Controller {
	return &controller {
		ctlr.New(jwt),
		service,
	}
}

// Get get a contact by contactID.
func (c *controller) Get() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := c.ExtractID(r)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		vars := mux.Vars(r)
		contactID := vars["id"]

		contact, err := c.service.Get(r.Context(), userID, contactID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, ErrContactNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		vm := &Presenter{}
		vm.FromEntity(contact)

		c.RespondWithJSON(w, http.StatusOK, vm)
	})
}

// GetAll get all contacts from the profile.
func (c *controller) GetAll() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := c.ExtractID(r)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		contacts, err := c.service.GetAll(r.Context(), userID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, ErrContactNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		var vms []*Presenter

		for _, contact := range contacts {
			vm := &Presenter{}
			vm.FromEntity(contact)
			vms = append(vms, vm)
		}

		c.RespondWithJSON(w, http.StatusOK, vms)
	})
}

// GetPresence obtain the presence of the contact by contactID.
func (c *controller) GetPresence() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := c.ExtractID(r)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		vars := mux.Vars(r)
		contactID := vars["id"]

		presence, err := c.service.GetPresence(r.Context(), userID, contactID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, ErrContactNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		p := &Presence{
			ID: contactID,
			Presence: PresenceTypeText(presence),
		}

		c.RespondWithJSON(w, http.StatusOK, p)
	})
}

// Create creates a new contact.
func (c *controller) Create() http.Handler {
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

		err = c.service.Create(r.Context(), input.ID, input.Name, input.LastName, userID)
		if err != nil {
			log.Println(err.Error())

			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				c.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			if errors.Is(err, ErrUserNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			if errors.Is(err, ErrContactAlreadyExists) {
				c.RespondWithError(w, http.StatusConflict, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

// Update updates contact data.
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
		input.UserID = userID

		err = c.service.Update(r.Context(), input.ToEntity())
		if err != nil {
			log.Println(err.Error())

			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				c.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			if errors.Is(err, ErrContactNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// Delete deletes a contact by contactID.
func (c *controller) Delete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := c.ExtractID(r)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		vars := mux.Vars(r)
		contactID := vars["id"]

		err = c.service.Delete(r.Context(), userID, contactID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, ErrContactNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// Block blocks a profile from receiving a message.
func (c *controller) Block() http.Handler {
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

		err = c.service.Block(r.Context(), userID, input.ID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, ErrUserNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			if errors.Is(err, ErrContactAlreadyBlocked) {
				c.RespondWithError(w, http.StatusConflict, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// Unblock unblock a profile to receive message.
func (c *controller) Unblock() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := c.ExtractID(r)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		vars := mux.Vars(r)
		blockedUserID := vars["id"]

		err = c.service.Unblock(r.Context(), userID, blockedUserID)
		if err != nil {
			log.Println(err.Error())

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
