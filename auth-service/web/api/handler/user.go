package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/auth-service/app/user"
	"github.com/tsmweb/auth-service/common"
	"github.com/tsmweb/auth-service/web/api/dto"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/go-helper-api/httputil"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

// GetUser a user by ID.
func GetUser(jwt auth.JWT, getUseCase user.GetUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := jwt.GetDataToken(r, "id")
		if err != nil || data == nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		userID := data.(string)

		u, err := getUseCase.Execute(r.Context(), userID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, user.ErrUserNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		userDto := &dto.User{}
		userDto.FromEntity(u)

		httputil.RespondWithJSON(w, http.StatusOK, userDto)
	})
}

// CreateUser a new user.
func CreateUser(createUseCase user.CreateUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !httputil.HasContentType(r, httputil.MimeApplicationJSON) {
			httputil.RespondWithError(w, http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType))
			return
		}

		userDto := &dto.User{}
		err := json.NewDecoder(r.Body).Decode(&userDto)
		if err != nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		err = createUseCase.Execute(r.Context(), userDto.ID, userDto.Name, userDto.LastName, userDto.Password)
		if err != nil {
			log.Println(err.Error())

			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				httputil.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			if errors.Is(err, user.ErrUserAlreadyExists) {
				httputil.RespondWithError(w, http.StatusConflict, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

// UpdateUser updates user data.
func UpdateUser(jwt auth.JWT, updateUseCase user.UpdateUseCase) http.Handler {
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

		userDto := &dto.User{}
		err = json.NewDecoder(r.Body).Decode(&userDto)
		if err != nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		ctx := context.WithValue(r.Context(), common.AuthContextKey, userID)

		err = updateUseCase.Execute(ctx, userDto.ToEntity())
		if err != nil {
			log.Println(err.Error())

			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				httputil.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			if errors.Is(err, user.ErrOperationNotAllowed) {
				httputil.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}

			if errors.Is(err, user.ErrUserNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

const userApiVersion string = "v1"

var userResource string

func init() {
	userResource = fmt.Sprintf("/%s/user", userApiVersion)
}

func MakeUserHandlers(
	r *mux.Router,
	jwt auth.JWT,
	auth middleware.Auth,
	getUseCase user.GetUseCase,
	createUseCase user.CreateUseCase,
	updateUseCase user.UpdateUseCase) {

	// user [GET]
	r.Handle(userResource, negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(GetUser(jwt, getUseCase))),
	).Methods(http.MethodGet)

	// user [POST]
	r.Handle(userResource, CreateUser(createUseCase)).
		Methods(http.MethodPost)

	// user [PUT]
	r.Handle(userResource, negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(UpdateUser(jwt, updateUseCase))),
	).Methods(http.MethodPut)
}
