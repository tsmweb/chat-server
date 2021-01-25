package contact

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/tsmweb/go-helper-api/auth"
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
		ctlr.NewController(jwt),
		service,
	}
}

// Get get a contact by contactID.
func (c *controller) Get() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		profileID, err := c.ExtractID(r)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		vars := mux.Vars(r)
		contactID := vars["id"]

		contact, err := c.service.Get(profileID, contactID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, ErrContactNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		vm := Presenter{}
		vm.FromEntity(contact)

		c.RespondWithJSON(w, http.StatusOK, vm)
	})
}

// GetAll get all contacts from the profile.
func (c *controller) GetAll() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		profileID, err := c.ExtractID(r)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		contacts, err := c.service.GetAll(profileID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, ErrContactNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		var vms []Presenter

		for _, contact := range contacts {
			vm := Presenter{}
			vm.FromEntity(contact)
			vms = append(vms, vm)
		}

		c.RespondWithJSON(w, http.StatusOK, vms)
	})
}

// GetPresence obtain the presence of the contact by contactID.
func (c *controller) GetPresence() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO
	})
}

// Create creates a new contact.
func (c *controller) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO
	})
}

// Update updates contact data.
func (c *controller) Update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO
	})
}

// Delete deletes a contact by contactID.
func (c *controller) Delete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO
	})
}

// Block blocks a profile from receiving a message.
func (c *controller) Block() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO
	})
}

// Unblock unblock a profile to receive message.
func (c *controller) Unblock() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO
	})
}
