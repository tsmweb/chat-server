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

func TestHandler_GetMediaFile(t *testing.T) {
	if err := config.Load("../../"); err != nil {
		t.Error(err)
	}

	t.Run("when handler.GetMediaFile return StatusNotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("%s/9f2093abeecac621b55489bf8cb0e08ee00d5fe6da6f30e77214a648e58bd91a.jpg",
				mediaResource), nil)
		rec := httptest.NewRecorder()

		handler := GetMediaFile()

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{name}", mediaResource), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.GetMediaFile return StatusOK", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("%s/9f2093abeecac621b55489bf8cb0e08ee00d5fe6da6f30e77214a648e58bd91b.jpg",
				mediaResource), nil)
		rec := httptest.NewRecorder()

		handler := GetMediaFile()

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{name}", mediaResource), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestHandler_UploadMediaFile(t *testing.T) {
	if err := config.Load("../../"); err != nil {
		t.Error(err)
	}

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, mediaResource, bytes.NewReader([]byte("")))
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

	t.Run("when handler.UploadMediaFile return with StatusBadRequest", func(t *testing.T) {
		config.SetMaxUploadSize(2) // KB

		contentType, content, err := createImageBuffer("jpg")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, mediaResource, content)
		req.Header.Add("Content-Type", contentType)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()

		UploadMediaFile(mJWT).ServeHTTP(rec, req)

		t.Log(rec.Result().Status)
		t.Log(rec.Body.String())
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when handler.UploadMediaFile return with StatusCreated", func(t *testing.T) {
		config.SetMaxUploadSize(1024) // KB

		contentType, content, err := createImageBuffer("jpg")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, mediaResource, content)
		req.Header.Add("Content-Type", contentType)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()

		UploadMediaFile(mJWT).ServeHTTP(rec, req)

		t.Log(rec.Result().Status)
		t.Log(rec.Body.String())
		t.Log(rec.Header())
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}
