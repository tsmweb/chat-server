package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/file-service/common"
	"github.com/tsmweb/file-service/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetUserFile(t *testing.T) {
	if err := config.Load("../../"); err != nil {
		t.Error(err)
	}

	t.Run("when handler.GetUserFile return StatusNotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/+5518900000000", userResource), nil)
		rec := httptest.NewRecorder()

		handler := GetUserFile()

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", userResource), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.GetUserFile return StatusOK", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/+5518977777777", userResource), nil)
		rec := httptest.NewRecorder()

		handler := GetUserFile()

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", userResource), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, rec.Header().Get("Content-Type"), "image/jpeg")
	})
}

func TestHandler_UploadUserFile(t *testing.T) {
	if err := config.Load("../../"); err != nil {
		t.Error(err)
	}

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, userResource, bytes.NewReader([]byte("")))
		req.Header.Set("Content-Type", "image/jpeg")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		UploadUserFile(mJWT).ServeHTTP(rec, req)

		//t.Log(rec.Result().Status)
		//t.Log(rec.Body.String())
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.UploadUserFile return with StatusBadRequest", func(t *testing.T) {
		config.SetMaxUploadSize(2) // KB

		contentType, content, err := createImageBuffer("jpg")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, userResource, content)
		req.Header.Add("Content-Type", contentType)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()

		UploadUserFile(mJWT).ServeHTTP(rec, req)

		t.Log(rec.Result().Status)
		t.Log(rec.Body.String())
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when handler.UploadUserFile return with StatusUnsupportedMediaType", func(t *testing.T) {
		config.SetMaxUploadSize(1024) // KB

		contentType, content, err := createImageBuffer("png")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, userResource, content)
		req.Header.Add("Content-Type", contentType)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()

		UploadUserFile(mJWT).ServeHTTP(rec, req)

		t.Log(rec.Result().Status)
		t.Log(rec.Body.String())
		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when handler.UploadUserFile return with StatusCreated", func(t *testing.T) {
		config.SetMaxUploadSize(1024) // KB

		contentType, content, err := createImageBuffer("jpg")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, userResource, content)
		req.Header.Add("Content-Type", contentType)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()

		UploadUserFile(mJWT).ServeHTTP(rec, req)

		//t.Log(rec.Result().Status)
		//t.Log(rec.Body.String())
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

}
