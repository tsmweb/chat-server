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
	"github.com/tsmweb/file-service/group"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetGroupFile(t *testing.T) {
	if err := config.Load("../../"); err != nil {
		t.Error(err)
	}

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", groupResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mValidateUseCase := new(mockGroupValidateUseCase)

		handler := GetGroupFile(mJWT, mValidateUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", groupResource), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.GetGroupFile return with StatusNotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", groupResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()
		mValidateUseCase := new(mockGroupValidateUseCase)
		mValidateUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(group.ErrGroupNotFound).
			Once()

		handler := GetGroupFile(mJWT, mValidateUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", groupResource), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		t.Log(rec.Result().Status)
		t.Log(rec.Body.String())
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.GetGroupFile return with StatusUnauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", groupResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()
		mValidateUseCase := new(mockGroupValidateUseCase)
		mValidateUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(group.ErrOperationNotAllowed).
			Once()

		handler := GetGroupFile(mJWT, mValidateUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", groupResource), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		t.Log(rec.Result().Status)
		t.Log(rec.Body.String())
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when handler.GetGroupFile return with StatusOK", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", groupResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()
		mValidateUseCase := new(mockGroupValidateUseCase)
		mValidateUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		handler := GetGroupFile(mJWT, mValidateUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", groupResource), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, rec.Header().Get("Content-Type"), "image/jpeg")
	})
}

func TestHandler_UploadGroupFile(t *testing.T) {
	if err := config.Load("../../"); err != nil {
		t.Error(err)
	}

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, groupResource, bytes.NewReader([]byte("")))
		req.Header.Set("Content-Type", "image/jpeg")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mValidateUseCase := new(mockGroupValidateUseCase)

		UploadGroupFile(mJWT, mValidateUseCase).ServeHTTP(rec, req)

		//t.Log(rec.Result().Status)
		//t.Log(rec.Body.String())
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.UploadGroupFile return with StatusBadRequest", func(t *testing.T) {
		config.SetMaxUploadSize(2) // KB

		contentType, content, err := createImageBuffer("jpg")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, groupResource, content)
		req.Header.Add("Content-Type", contentType)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()
		mValidateUseCase := new(mockGroupValidateUseCase)

		UploadGroupFile(mJWT, mValidateUseCase).ServeHTTP(rec, req)

		t.Log(rec.Result().Status)
		t.Log(rec.Body.String())
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when handler.UploadGroupFile return with StatusNotFound", func(t *testing.T) {
		config.SetMaxUploadSize(1024) // KB

		contentType, content, err := createImageBuffer("jpg")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, groupResource, content)
		req.Header.Add("Content-Type", contentType)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()
		mValidateUseCase := new(mockGroupValidateUseCase)
		mValidateUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(group.ErrGroupNotFound).
			Once()

		UploadGroupFile(mJWT, mValidateUseCase).ServeHTTP(rec, req)

		t.Log(rec.Result().Status)
		t.Log(rec.Body.String())
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.UploadGroupFile return with StatusUnauthorized", func(t *testing.T) {
		config.SetMaxUploadSize(1024) // KB

		contentType, content, err := createImageBuffer("jpg")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, groupResource, content)
		req.Header.Add("Content-Type", contentType)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()
		mValidateUseCase := new(mockGroupValidateUseCase)
		mValidateUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(group.ErrOperationNotAllowed).
			Once()

		UploadGroupFile(mJWT, mValidateUseCase).ServeHTTP(rec, req)

		t.Log(rec.Result().Status)
		t.Log(rec.Body.String())
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when handler.UploadGroupFile return with StatusUnsupportedMediaType", func(t *testing.T) {
		config.SetMaxUploadSize(1024) // KB

		contentType, content, err := createImageBuffer("png")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, groupResource, content)
		req.Header.Add("Content-Type", contentType)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()
		mValidateUseCase := new(mockGroupValidateUseCase)
		mValidateUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		UploadGroupFile(mJWT, mValidateUseCase).ServeHTTP(rec, req)

		t.Log(rec.Result().Status)
		t.Log(rec.Body.String())
		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when handler.UploadGroupFile return with StatusCreated", func(t *testing.T) {
		config.SetMaxUploadSize(1024) // KB

		contentType, content, err := createImageBuffer("jpg")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, groupResource, content)
		req.Header.Add("Content-Type", contentType)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()
		mValidateUseCase := new(mockGroupValidateUseCase)
		mValidateUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		UploadGroupFile(mJWT, mValidateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}
