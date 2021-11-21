package handler

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/file-service/config"
	"github.com/tsmweb/file-service/group"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/httputil"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/urfave/negroni"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

// GetGroupFile gets the group image by ID.
func GetGroupFile(jwt auth.JWT, validateUseCase group.ValidateUseCase) http.Handler {
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

		if err = validateUseCase.Execute(r.Context(), groupID, userID, false); err != nil {
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

		path := filepath.Join(config.GroupFilePath(), fmt.Sprintf("%s.jpg", groupID))
		fileBytes, err := ioutil.ReadFile(path)
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

// UploadUserFile uploads an image to the group.
func UploadGroupFile(jwt auth.JWT, validateUseCase group.ValidateUseCase) http.Handler {
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
		if err = validateFileSize(w, r); err != nil {
			httputil.RespondWithError(w, http.StatusBadRequest, "The uploaded file is too big.")
			return
		}

		groupID := r.FormValue("id")

		if err = validateUseCase.Execute(r.Context(), groupID, userID,true); err != nil {
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

		file, _, err := r.FormFile("file")
		if err != nil {
			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer file.Close()

		// Get content type.
		_, fileExtension, err := getContentType(file)
		if err != nil || fileExtension != "jpg" {
			httputil.RespondWithError(w, http.StatusUnsupportedMediaType,
				http.StatusText(http.StatusUnsupportedMediaType))
			return
		}

		// Creates the file on the local file system.
		path := filepath.Join(config.GroupFilePath(), fmt.Sprintf("%s.%s", groupID, fileExtension))
		if err = copyFile(path, file); err != nil {
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
	validateUseCase group.ValidateUseCase,
) {
	// group/{id} [GET]
	r.Handle(fmt.Sprintf("%s/{id}", groupResource), negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(GetGroupFile(jwt, validateUseCase))),
	).Methods(http.MethodGet)

	// group [POST]
	r.Handle(groupResource, negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(UploadGroupFile(jwt, validateUseCase))),
	).Methods(http.MethodPost)
}
