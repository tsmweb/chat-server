/*
Package nandler implements a Handler with utility methods.
The Handler is used in the composition of other more specific handlers.

Package handler also provides the MimeType type which represents the mime type as
"application/json", "text/plain", "image/jpeg", ...

*/

package handler

import (
	"encoding/json"
	"github.com/tsmweb/helper-go/auth"
	"mime"
	"net/http"
	"strings"
)

type errorMessage struct {
	ErrorMessage string `json:"error_message"`
}

type Headers map[string]string

// Handler base handler with utility methods.
type Handler struct {
	jwt *auth.JWT
}

// NewHandler returns an instance of the Handler.
func NewHandler(jwt *auth.JWT) *Handler {
	return &Handler{jwt}
}

// ExtractID extracts the JWT token id.
func (h *Handler) ExtractID(r *http.Request) (string, error) {
	ID, err := h.jwt.MapClaims(r, "sub") //sub = id
	if err != nil || ID == nil {
		return "", err
	}

	return ID.(string), nil
}

// HasContentType validates the content type, return true for valid and false for invalid.
func (h *Handler) HasContentType(r *http.Request, mimetype MimeType) bool {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return mimetype == MimeApplicationOctetStream
	}

	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}

		if t == mimetype.String() {
			return true
		}
	}

	return false
}

func (h *Handler) RespondWithHeader(writer http.ResponseWriter, status int, header Headers) {
	for k, v := range header {
		writer.Header().Set(k, v)
	}

	writer.WriteHeader(status)
}

func (h *Handler) RespondWithError(w http.ResponseWriter, status int, message string) {
	h.RespondWithHeader(w, status, Headers{"Content-Type": MimeTypeText(MimeApplicationJSON)})
	errorMessage, err := json.Marshal(errorMessage{ErrorMessage: message})
	if err == nil {
		w.Write(errorMessage)
	}
}

func (h *Handler) RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	h.RespondWithHeader(w, status, Headers{"Content-Type": MimeTypeText(MimeApplicationJSON)})
	jsonData, err := json.Marshal(data)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	} else {
		w.Write(jsonData)
	}
}
