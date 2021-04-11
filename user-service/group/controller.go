package group

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/cerror"
	ctlr "github.com/tsmweb/go-helper-api/controller"
	"github.com/tsmweb/user-service/common"
	"log"
	"net/http"
)

// Controller provides the end point for the routers.
type Controller interface {
	Get() http.Handler
	GetAll() http.Handler
	Create() http.Handler
	Update() http.Handler
	Delete() http.Handler
	AddMember() http.Handler
	RemoveMember() http.Handler
	SetAdmin() http.Handler
}

type controller struct {
	*ctlr.Controller
	service Service
}

// NewController creates a new instance of Controller.
func NewController(jwt auth.JWT, service Service) Controller {
	return &controller{
		ctlr.New(jwt),
		service,
	}
}

// Get get a group by groupID.
func (c *controller) Get() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := c.ExtractID(r)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		ctx := context.WithValue(r.Context(), common.AuthContextKey, userID)

		vars := mux.Vars(r)
		groupID := vars["id"]

		group, err := c.service.Get(ctx, groupID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, ErrGroupNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		vm := &Presenter{}
		vm.FromEntity(group)

		c.RespondWithJSON(w, http.StatusOK, vm)
	})
}

// GetAll get all the groups that the user is a member of.
func (c *controller) GetAll() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := c.ExtractID(r)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		ctx := context.WithValue(r.Context(), common.AuthContextKey, userID)

		groups, err := c.service.GetAll(ctx, userID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, ErrGroupNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		vms := EntityToPresenters(groups...)
		c.RespondWithJSON(w, http.StatusOK, vms)
	})
}

// Create creates a new group.
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
		ctx := context.WithValue(r.Context(), common.AuthContextKey, userID)

		input := &Presenter{}
		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		groupID, err := c.service.Create(ctx, input.Name, input.Description, userID)
		if err != nil {
			log.Println(err.Error())

			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				c.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		headers := ctlr.Headers{}
		headers["Location"] = fmt.Sprintf("%s/%s", resource, groupID)
		c.RespondWithHeader(w, http.StatusCreated, headers)
	})
}

// Update updates group data.
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
		ctx := context.WithValue(r.Context(), common.AuthContextKey, userID)

		input := &Presenter{}
		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

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

			if errors.Is(err, ErrGroupNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// Delete deletes a group by groupID.
func (c *controller) Delete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := c.ExtractID(r)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		ctx := context.WithValue(r.Context(), common.AuthContextKey, userID)

		vars := mux.Vars(r)
		groupID := vars["id"]

		err = c.service.Delete(ctx, groupID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, ErrOperationNotAllowed) {
				c.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}

			if errors.Is(err, ErrGroupNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// AddMember add member to group.
func (c *controller) AddMember() http.Handler {
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
		ctx := context.WithValue(r.Context(), common.AuthContextKey, userID)

		input := &MemberPresenter{}
		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		err = c.service.AddMember(ctx, input.GroupID, input.UserID, input.Admin)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, ErrOperationNotAllowed) {
				c.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}

			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				c.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			if errors.Is(err, ErrGroupNotFound) || errors.Is(err, ErrUserNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			if errors.Is(err, ErrMemberAlreadyExists) {
				c.RespondWithError(w, http.StatusConflict, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

// RemoveMember removes a member from the group.
func (c *controller) RemoveMember() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := c.ExtractID(r)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		ctx := context.WithValue(r.Context(), common.AuthContextKey, userID)

		vars := mux.Vars(r)
		groupID := vars["group"]
		memberID := vars["user"]

		err = c.service.RemoveMember(ctx, groupID, memberID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, ErrOperationNotAllowed) || errors.Is(err, ErrGroupOwnerCannotRemoved) {
				c.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}

			if errors.Is(err, ErrMemberNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// SetAdmin elevates a member to administrator status.
func (c *controller) SetAdmin() http.Handler {
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
		ctx := context.WithValue(r.Context(), common.AuthContextKey, userID)

		input := &MemberPresenter{}
		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			c.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		err = c.service.SetAdmin(ctx, input.ToEntity())
		if err != nil {
			log.Println(err.Error())

			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				c.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			if errors.Is(err, ErrOperationNotAllowed) || errors.Is(err, ErrGroupOwnerCannotChanged) {
				c.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}

			if errors.Is(err, ErrMemberNotFound) {
				c.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			c.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
