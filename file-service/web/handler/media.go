package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/file-service/config"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/httputil"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/tsmweb/go-helper-api/util/hashutil"
	"github.com/urfave/negroni"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

// GetMediaFile gets a media file by name.
func GetMediaFile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fileName := vars["name"]

		path := filepath.Join(config.MediaFilePath(), fileName)
		fileBytes, err := ioutil.ReadFile(path)
		if err != nil {
			httputil.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		//w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(fileBytes)))
		w.Write(fileBytes)
	})
}

// UploadMediaFile uploads a media file.
func UploadMediaFile(jwt auth.JWT) http.Handler {
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
		if err != nil {
			httputil.RespondWithError(w, http.StatusUnsupportedMediaType,
				http.StatusText(http.StatusUnsupportedMediaType))
			return
		}

		fileNameHash, _ := hashutil.HashSHA1(fmt.Sprintf("%s%v", userID, time.Now().UnixNano()))
		fileName := fmt.Sprintf("%s.%s", fileNameHash, fileExtension)

		// Creates the file on the local file system.
		path := filepath.Join(config.MediaFilePath(), fileName)
		if err = copyFile(path, file); err != nil {
			httputil.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		headers := httputil.Headers{}
		headers["Location"] = fmt.Sprintf("%s/%s", mediaResource, fileName)
		httputil.RespondWithHeader(w, http.StatusCreated, headers)
	})
}

const mediaApiVersion string = "v1"

var mediaResource string

func init() {
	mediaResource = fmt.Sprintf("/%s/media", mediaApiVersion)
}

func MakeMediaHandlers(
	r *mux.Router,
	jwt auth.JWT,
	auth middleware.Auth,
) {
	// media/{name} [GET]
	r.Handle(fmt.Sprintf("%s/{name}", mediaResource), negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(GetMediaFile())),
	).Methods(http.MethodGet)

	// media [POST]
	r.Handle(mediaResource, negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(UploadMediaFile(jwt))),
	).Methods(http.MethodPost)
}
