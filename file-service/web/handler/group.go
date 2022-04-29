package handler

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/file-service/app/group"
	"github.com/tsmweb/file-service/common/fileutil"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/httputil"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"strconv"
)

// GetGroupFile gets the group image by ID.
func GetGroupFile(jwt auth.JWT, getUseCase group.GetUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user id.
		data, err := jwt.GetDataToken(r, "id")
		if err != nil || data == nil {
			log.Println(err.Error())
			httputil.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		userID := data.(string)

		vars := mux.Vars(r)
		groupID := vars["id"]

		fileBytes, err := getUseCase.Execute(r.Context(), groupID, userID)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, group.ErrGroupNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			if errors.Is(err, group.ErrOperationNotAllowed) {
				httputil.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}

			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", httputil.MimeTypeText(httputil.MimeImageJPEG))
		w.Header().Set("Content-Length", strconv.Itoa(len(fileBytes)))
		w.Write(fileBytes)
	})
}

// UploadGroupFile uploads an image to the group.
func UploadGroupFile(jwt auth.JWT, uploadUseCase group.UploadUseCase) http.Handler {
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

		groupID := r.FormValue("id")
		file, _, err := r.FormFile("file")
		if err != nil {
			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer file.Close()

		if err = uploadUseCase.Execute(r.Context(), file, groupID, userID); err != nil {
			log.Println(err.Error())

			if errors.Is(err, group.ErrGroupNotFound) {
				httputil.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}
			if errors.Is(err, group.ErrOperationNotAllowed) {
				httputil.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}
			if errors.Is(err, group.ErrFileTooBig) {
				httputil.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			if errors.Is(err, group.ErrUnsupportedMediaType) {
				httputil.RespondWithError(w, http.StatusUnsupportedMediaType, err.Error())
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

const groupApiVersion string = "v1"

var groupResource string

func init() {
	groupResource = fmt.Sprintf("/%s/group", groupApiVersion)
}

func MakeGroupHandlers(
	r *mux.Router,
	jwt auth.JWT,
	auth middleware.Auth,
	getUseCase group.GetUseCase,
	uploadUseCase group.UploadUseCase,
) {
	// group/{id} [GET]
	r.Handle(fmt.Sprintf("%s/{id}", groupResource), negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(GetGroupFile(jwt, getUseCase))),
	).Methods(http.MethodGet)

	// group [POST]
	r.Handle(groupResource, negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(UploadGroupFile(jwt, uploadUseCase))),
	).Methods(http.MethodPost)
}
