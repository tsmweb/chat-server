package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/auth-service/common"
	"github.com/tsmweb/auth-service/login"
	"github.com/tsmweb/auth-service/web/api/dto"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/go-helper-api/httputil"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

// Login returns a token if ID and password are valid.
func Login(loginUseCase login.LoginUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !httputil.HasContentType(r, httputil.MimeApplicationJSON) {
			httputil.RespondWithError(w, http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType))
			return
		}

		input := dto.Login{}
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&input); err != nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		token, err := loginUseCase.Execute(r.Context(), input.ID, input.Password)
		if err != nil {
			log.Println(err.Error())
			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				httputil.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			if errors.Is(err, cerror.ErrUnauthorized) {
				httputil.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		httputil.RespondWithJSON(w, http.StatusOK, &dto.TokenAuth{Token: token})
	})
}

// UpdatePassword updates password in data base.
func UpdatePassword(jwt auth.JWT, updateUseCase login.UpdateUseCase) http.Handler {
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

		input := dto.Login{}
		decoder := json.NewDecoder(r.Body)

		if err = decoder.Decode(&input); err != nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusUnprocessableEntity, "Malformed JSON")
			return
		}

		ctx := context.WithValue(r.Context(), common.AuthContextKey, userID)

		if err = updateUseCase.Execute(ctx, input.ToEntity()); err != nil {
			log.Println(err.Error())
			var errValidateModel *cerror.ErrValidateModel
			if errors.As(err, &errValidateModel) {
				httputil.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			if errors.Is(err, login.ErrOperationNotAllowed) {
				httputil.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}

			if errors.Is(err, login.ErrUserNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

const loginApiVersion string = "v1"

var loginResource string

func init() {
	loginResource = fmt.Sprintf("/%s/login", loginApiVersion)
}

func MakeLoginHandlers(
	r *mux.Router,
	jwt auth.JWT,
	auth middleware.Auth,
	loginUseCase login.LoginUseCase,
	updateUseCase login.UpdateUseCase) {
	// login [POST]
	r.Handle(loginResource, Login(loginUseCase)).
		Methods(http.MethodPost)

	// login [PUT]
	r.Handle(loginResource, negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(UpdatePassword(jwt, updateUseCase)),
	)).Methods(http.MethodPut)
}