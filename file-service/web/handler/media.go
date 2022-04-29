package handler

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/file-service/app/media"
	"github.com/tsmweb/file-service/common/fileutil"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/httputil"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"strconv"
)

// GetMediaFile gets a media file by name.
func GetMediaFile(getUseCase media.GetUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fileName := vars["name"]

		fileBytes, err := getUseCase.Execute(fileName)
		if err != nil {
			httputil.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Length", strconv.Itoa(len(fileBytes)))
		w.Write(fileBytes)
	})
}

// UploadMediaFile uploads a media file.
func UploadMediaFile(jwt auth.JWT, uploadUseCase media.UploadUseCase) http.Handler {
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

		fileName, err := uploadUseCase.Execute(userID, file)
		if err != nil {
			log.Println(err.Error())

			if errors.Is(err, media.ErrFileTooBig) {
				httputil.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			if errors.Is(err, media.ErrUnsupportedMediaType) {
				httputil.RespondWithError(w, http.StatusUnsupportedMediaType, err.Error())
				return
			}

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
	getUseCase media.GetUseCase,
	uploadUseCase media.UploadUseCase,
) {
	// media/{name} [GET]
	r.Handle(fmt.Sprintf("%s/{name}", mediaResource), negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(GetMediaFile(getUseCase))),
	).Methods(http.MethodGet)

	// media [POST]
	r.Handle(mediaResource, negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.Wrap(UploadMediaFile(jwt, uploadUseCase))),
	).Methods(http.MethodPost)
}
