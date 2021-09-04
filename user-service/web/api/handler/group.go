package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/go-helper-api/httputil"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/tsmweb/user-service/common"
	"github.com/tsmweb/user-service/group"
	"github.com/tsmweb/user-service/web/api/dto"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

// GetGroup get a group by groupID.
func GetGroup(jwt auth.JWT, getUseCase group.GetUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := jwt.GetDataToken(r, "id")
		if err != nil || data == nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		userID := data.(string)
		ctx := context.WithValue(r.Context(), common.AuthContextKey, userID)

		vars := mux.Vars(r)
		groupID := vars["id"]

		grp, err := getUseCase.Execute(ctx, groupID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, group.ErrGroupNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		groupDto := &dto.Group{}
		groupDto.FromEntity(grp)

		httputil.RespondWithJSON(w, http.StatusOK, groupDto)
	})
}

// GetAllGroups get all the groups that the user is a member of.
func GetAllGroups(jwt auth.JWT, getAllUseCase group.GetAllUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := jwt.GetDataToken(r, "id")
		if err != nil || data == nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		userID := data.(string)
		ctx := context.WithValue(r.Context(), common.AuthContextKey, userID)

		groups, err := getAllUseCase.Execute(ctx, userID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, group.ErrGroupNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		groupsDto := dto.EntityToGroupDTO(groups...)
		httputil.RespondWithJSON(w, http.StatusOK, groupsDto)
	})
}

// CreateGroup creates a new group.
func CreateGroup(jwt auth.JWT, createUseCase group.CreateUseCase) http.Handler {
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
		ctx := context.WithValue(r.Context(), common.AuthContextKey, userID)

		input := &dto.Group{}
		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		groupID, err := createUseCase.Execute(ctx, input.Name, input.Description, userID)
		if err != nil {
			log.Println(err.Error())

			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				httputil.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		headers := httputil.Headers{}
		headers["Location"] = fmt.Sprintf("%s/%s", groupResource, groupID)
		httputil.RespondWithHeader(w, http.StatusCreated, headers)
	})
}

// UpdateGroup updates group data.
func UpdateGroup(jwt auth.JWT, updateUseCase group.UpdateUseCase) http.Handler {
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
		ctx := context.WithValue(r.Context(), common.AuthContextKey, userID)

		input := &dto.Group{}
		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		err = updateUseCase.Execute(ctx, input.ToEntity())
		if err != nil {
			log.Println(err.Error())

			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				httputil.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			if errors.Is(err, group.ErrOperationNotAllowed) {
				httputil.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}

			if errors.Is(err, group.ErrGroupNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// DeleteGroup deletes a group by groupID.
func DeleteGroup(jwt auth.JWT, deleteUseCase group.DeleteUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := jwt.GetDataToken(r, "id")
		if err != nil || data == nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		userID := data.(string)
		ctx := context.WithValue(r.Context(), common.AuthContextKey, userID)

		vars := mux.Vars(r)
		groupID := vars["id"]

		err = deleteUseCase.Execute(ctx, groupID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, group.ErrOperationNotAllowed) {
				httputil.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}

			if errors.Is(err, group.ErrGroupNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// AddGroupMember add member to group.
func AddGroupMember(jwt auth.JWT, addMemberUseCase group.AddMemberUseCase) http.Handler {
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
		ctx := context.WithValue(r.Context(), common.AuthContextKey, userID)

		input := &dto.Member{}
		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		err = addMemberUseCase.Execute(ctx, input.GroupID, input.UserID, input.Admin)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, group.ErrOperationNotAllowed) {
				httputil.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}

			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				httputil.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			if errors.Is(err, group.ErrGroupNotFound) || errors.Is(err, group.ErrUserNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			if errors.Is(err, group.ErrMemberAlreadyExists) {
				httputil.RespondWithError(w, http.StatusConflict, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

// RemoveGroupMember removes a member from the group.
func RemoveGroupMember(jwt auth.JWT, removeMemberUseCase group.RemoveMemberUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := jwt.GetDataToken(r, "id")
		if err != nil || data == nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		userID := data.(string)
		ctx := context.WithValue(r.Context(), common.AuthContextKey, userID)

		vars := mux.Vars(r)
		groupID := vars["group"]
		memberID := vars["user"]

		err = removeMemberUseCase.Execute(ctx, groupID, memberID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, group.ErrOperationNotAllowed) || errors.Is(err, group.ErrGroupOwnerCannotRemoved) {
				httputil.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}

			if errors.Is(err, group.ErrMemberNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// SetGroupAdmin elevates a member to administrator status.
func SetGroupAdmin(jwt auth.JWT, setAdminUseCase group.SetAdminUseCase) http.Handler {
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
		ctx := context.WithValue(r.Context(), common.AuthContextKey, userID)

		input := &dto.Member{}
		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		err = setAdminUseCase.Execute(ctx, input.ToEntity())
		if err != nil {
			log.Println(err.Error())

			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				httputil.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			if errors.Is(err, group.ErrOperationNotAllowed) || errors.Is(err, group.ErrGroupOwnerCannotChanged) {
				httputil.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}

			if errors.Is(err, group.ErrMemberNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

const groupApiVersion string = "v1"

var (
	groupResource  string
	memberResource string
)

func init() {
	groupResource = fmt.Sprintf("/%s/group", groupApiVersion)
	memberResource = fmt.Sprintf("/%s/group/member", groupApiVersion)
}

// MakeGroupRouters creates a router for Group.
func MakeGroupRouters(
	r *mux.Router,
	jwt auth.JWT,
	auth middleware.Auth,
	getUseCase group.GetUseCase,
	getAllUseCase group.GetAllUseCase,
	createUseCase group.CreateUseCase,
	updateUseCase group.UpdateUseCase,
	deleteUseCase group.DeleteUseCase,
	addMemberUseCase group.AddMemberUseCase,
	removeMemberUseCase group.RemoveMemberUseCase,
	setAdminUseCase group.SetAdminUseCase) {
	// group/{id} [GET]
	r.Handle(fmt.Sprintf("%s/{id}", groupResource), negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(GetGroup(jwt, getUseCase))),
	).Methods(http.MethodGet)

	// group [GET]
	r.Handle(groupResource, negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(GetAllGroups(jwt, getAllUseCase))),
	).Methods(http.MethodGet)

	// group [POST]
	r.Handle(groupResource, negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(CreateGroup(jwt, createUseCase))),
	).Methods(http.MethodPost)

	// group [PUT]
	r.Handle(groupResource, negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(UpdateGroup(jwt, updateUseCase))),
	).Methods(http.MethodPut)

	// group/{id} [DELETE]
	r.Handle(fmt.Sprintf("%s/{id}", groupResource), negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(DeleteGroup(jwt, deleteUseCase))),
	).Methods(http.MethodDelete)

	// group/member [POST]
	r.Handle(memberResource, negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(AddGroupMember(jwt, addMemberUseCase))),
	).Methods(http.MethodPost)

	// group/member/{group}/{user} [DELETE]
	r.Handle(fmt.Sprintf("%s/{group}/{user}", memberResource), negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(RemoveGroupMember(jwt, removeMemberUseCase))),
	).Methods(http.MethodDelete)

	// group/member [PUT]
	r.Handle(memberResource, negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(SetGroupAdmin(jwt, setAdminUseCase))),
	).Methods(http.MethodPut)
}
