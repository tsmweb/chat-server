package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/file-service/config"
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

// GetUserFile gets the user image by ID.
func GetUserFile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["id"]

		path := filepath.Join(config.UserFilePath(), fmt.Sprintf("%s.jpg", userID))
		fileBytes, err := ioutil.ReadFile(path)
		if err != nil {
			httputil.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(fileBytes)))
		w.Write(fileBytes)
	})
}

// UploadUserFile uploads an image to the user.
func UploadUserFile(jwt auth.JWT) http.Handler {
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
		path := filepath.Join(config.UserFilePath(), fmt.Sprintf("%s.%s", userID, fileExtension))
		if err = copyFile(path, file); err != nil {
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
) {
	// user/{id} [GET]
	r.Handle(fmt.Sprintf("%s/{id}", userResource), negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(GetUserFile())),
	).Methods(http.MethodGet)

	// user [POST]
	r.Handle(userResource, negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(UploadUserFile(jwt))),
	).Methods(http.MethodPost)
}
