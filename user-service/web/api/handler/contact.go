package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/go-helper-api/httputil"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/tsmweb/user-service/contact"
	"github.com/tsmweb/user-service/web/api/dto"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

// GetContact get a contact by contactID.
func GetContact(jwt auth.JWT, getUseCase contact.GetUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := jwt.GetDataToken(r, "id")
		if err != nil || data == nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		userID := data.(string)

		vars := mux.Vars(r)
		contactID := vars["id"]

		ct, err := getUseCase.Execute(r.Context(), userID, contactID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, contact.ErrContactNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		vm := &dto.Contact{}
		vm.FromEntity(ct)

		httputil.RespondWithJSON(w, http.StatusOK, vm)
	})
}

// GetAllContacts get all contacts from the profile.
func GetAllContacts(jwt auth.JWT, getAllUseCase contact.GetAllUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := jwt.GetDataToken(r, "id")
		if err != nil || data == nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		userID := data.(string)

		contacts, err := getAllUseCase.Execute(r.Context(), userID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, contact.ErrContactNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		vms := dto.EntityToContactDTO(contacts...)
		httputil.RespondWithJSON(w, http.StatusOK, vms)
	})
}

// GetContactPresence obtain the presence of the contact by contactID.
func GetContactPresence(jwt auth.JWT, getPresenceUseCase contact.GetPresenceUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := jwt.GetDataToken(r, "id")
		if err != nil || data == nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		userID := data.(string)

		vars := mux.Vars(r)
		contactID := vars["id"]

		presence, err := getPresenceUseCase.Execute(r.Context(), userID, contactID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, contact.ErrContactNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		p := &dto.Presence{
			ID: contactID,
			Presence: contact.PresenceTypeText(presence),
		}

		httputil.RespondWithJSON(w, http.StatusOK, p)
	})
}

// CreateContact creates a new contact.
func CreateContact(jwt auth.JWT, createUseCase contact.CreateUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !httputil.HasContentType(r, httputil.MimeApplicationJSON) {
			httputil.RespondWithError(w, http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType))
			return
		}

		data, err := jwt.GetDataToken(r, "id")
		if err != nil || data == nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		userID := data.(string)

		input := &dto.Contact{}
		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		err = createUseCase.Execute(r.Context(), input.ID, input.Name, input.LastName, userID)
		if err != nil {
			log.Println(err.Error())

			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				httputil.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			if errors.Is(err, contact.ErrUserNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			if errors.Is(err, contact.ErrContactAlreadyExists) {
				httputil.RespondWithError(w, http.StatusConflict, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

// UpdateContact updates contact data.
func UpdateContact(jwt auth.JWT, updateUseCase contact.UpdateUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !httputil.HasContentType(r, httputil.MimeApplicationJSON) {
			httputil.RespondWithError(w, http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType))
			return
		}

		data, err := jwt.GetDataToken(r, "id")
		if err != nil || data == nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		userID := data.(string)

		input := &dto.Contact{}
		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}
		input.UserID = userID

		err = updateUseCase.Execute(r.Context(), input.ToEntity())
		if err != nil {
			log.Println(err.Error())

			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				httputil.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			if errors.Is(err, contact.ErrContactNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// DeleteContact deletes a contact by contactID.
func DeleteContact(jwt auth.JWT, deleteUseCase contact.DeleteUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := jwt.GetDataToken(r, "id")
		if err != nil || data == nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		userID := data.(string)

		vars := mux.Vars(r)
		contactID := vars["id"]

		err = deleteUseCase.Execute(r.Context(), userID, contactID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, contact.ErrContactNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// BlockContact blocks a profile from receiving a message.
func BlockContact(jwt auth.JWT, blockUseCase contact.BlockUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !httputil.HasContentType(r, httputil.MimeApplicationJSON) {
			httputil.RespondWithError(w, http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType))
			return
		}

		data, err := jwt.GetDataToken(r, "id")
		if err != nil || data == nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		userID := data.(string)

		input := &dto.Contact{}
		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		err = blockUseCase.Execute(r.Context(), userID, input.ID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, contact.ErrUserNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			if errors.Is(err, contact.ErrContactAlreadyBlocked) {
				httputil.RespondWithError(w, http.StatusConflict, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// UnblockContact unblock a profile to receive message.
func UnblockContact(jwt auth.JWT, unblockUseCase contact.UnblockUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := jwt.GetDataToken(r, "id")
		if err != nil || data == nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		userID := data.(string)

		vars := mux.Vars(r)
		blockedUserID := vars["id"]

		err = unblockUseCase.Execute(r.Context(), userID, blockedUserID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, contact.ErrUserNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

const contactApiVersion string = "v1"

var contactResource string

func init() {
	contactResource = fmt.Sprintf("/%s/contact", contactApiVersion)
}

// MakeContactRouters creates a router for Contact.
func MakeContactRouters(
	r *mux.Router,
	jwt auth.JWT,
	auth middleware.Auth,
	getUseCase contact.GetUseCase,
	getAllUseCase contact.GetAllUseCase,
	getPresenceUseCase contact.GetPresenceUseCase,
	createUseCase contact.CreateUseCase,
	updateUseCase contact.UpdateUseCase,
	deleteUseCase contact.DeleteUseCase,
	blockUseCase contact.BlockUseCase,
	unblockUseCase contact.UnblockUseCase) {
	// contact/{id} [GET]
	r.Handle(fmt.Sprintf("%s/{id}", contactResource), negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(GetContact(jwt, getUseCase))),
	).Methods(http.MethodGet)

	// contact [GET]
	r.Handle(contactResource, negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(GetAllContacts(jwt, getAllUseCase))),
	).Methods(http.MethodGet)

	// contact/presence/{id} [GET]
	r.Handle(fmt.Sprintf("%s/presence/{id}", contactResource), negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(GetContactPresence(jwt, getPresenceUseCase))),
	).Methods(http.MethodGet)

	// contact [POST]
	r.Handle(contactResource, negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(CreateContact(jwt, createUseCase))),
	).Methods(http.MethodPost)

	// contact [PUT]
	r.Handle(contactResource, negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(UpdateContact(jwt, updateUseCase))),
	).Methods(http.MethodPut)

	// contact/{id} [DELETE]
	r.Handle(fmt.Sprintf("%s/{id}", contactResource), negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(DeleteContact(jwt, deleteUseCase))),
	).Methods(http.MethodDelete)

	// contact/block [POST]
	r.Handle(fmt.Sprintf("%s/block", contactResource), negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(BlockContact(jwt, blockUseCase))),
	).Methods(http.MethodPost)

	// contact/block/{id} [DELETE]
	r.Handle(fmt.Sprintf("%s/block/{id}", contactResource), negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(UnblockContact(jwt, unblockUseCase))),
	).Methods(http.MethodDelete)
}