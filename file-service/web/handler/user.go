package handler

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/file-service/app/user"
	"github.com/tsmweb/file-service/common/fileutil"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/httputil"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"strconv"
)

// GetUserFile gets the user image by ID.
func GetUserFile(getUseCase user.GetUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["id"]

		fileBytes, err := getUseCase.Execute(userID)
		if err != nil {
			httputil.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", httputil.MimeTypeText(httputil.MimeImageJPEG))
		w.Header().Set("Content-Length", strconv.Itoa(len(fileBytes)))
		w.Write(fileBytes)
	})
}

// UploadUserFile uploads an image to the user.
func UploadUserFile(jwt auth.JWT, uploadUseCase user.UploadUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user id.
		data, err := jwt.GetDataToken(r, "id")
		if err != nil || data == nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		userID := data.(string)

		// Validate file size.
		if err = fileutil.ValidateFileSize(w, r); err != nil {
			httputil.RespondWithError(w, http.StatusBadRequest, "The uploaded file is too big.")
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer file.Close()

		if err = uploadUseCase.Execute(userID, file); err != nil {
			log.Println(err.Error())

			if errors.Is(err, user.ErrFileTooBig) {
				httputil.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			if errors.Is(err, user.ErrUnsupportedMediaType) {
				httputil.RespondWithError(w, http.StatusUnsupportedMediaType, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		headers := httputil.Headers{}
		headers["Location"] = fmt.Sprintf("%s/%s", userResource, userID)
		httputil.RespondWithHeader(w, http.StatusCreated, headers)
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
	uploadUseCase user.UploadUseCase,
) {
	// user/{id} [GET]
	r.Handle(fmt.Sprintf("%s/{id}", userResource), negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(GetUserFile(getUseCase))),
	).Methods(http.MethodGet)

	// user [POST]
	r.Handle(userResource, negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(UploadUserFile(jwt, uploadUseCase))),
	).Methods(http.MethodPost)
}
