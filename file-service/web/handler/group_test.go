package handler

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/file-service/app/group"
	"github.com/tsmweb/file-service/common/appmock"
	"github.com/tsmweb/file-service/common/imageutil"
	"github.com/tsmweb/file-service/config"
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

		mJWT := new(appmock.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mRepo := new(appmock.MockGroupRepository)
		getUseCase := group.NewGetUseCase(mRepo)

		handler := GetGroupFile(mJWT, getUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", groupResource), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.GetGroupFile return with StatusNotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", groupResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(appmock.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()

		mRepo := new(appmock.MockGroupRepository)
		mRepo.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(false, nil).
			Once()
		getUseCase := group.NewGetUseCase(mRepo)

		handler := GetGroupFile(mJWT, getUseCase)

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

		mJWT := new(appmock.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()
		mRepo := new(appmock.MockGroupRepository)
		mRepo.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		mRepo.On("IsGroupMember", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()
		getUseCase := group.NewGetUseCase(mRepo)

		handler := GetGroupFile(mJWT, getUseCase)

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

		mJWT := new(appmock.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()
		mRepo := new(appmock.MockGroupRepository)
		mRepo.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		mRepo.On("IsGroupMember", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		getUseCase := group.NewGetUseCase(mRepo)

		handler := GetGroupFile(mJWT, getUseCase)

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
		config.SetMaxUploadSize(1024) // KB

		contentType, content, err := imageutil.CreateImageBuffer("jpg")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, groupResource, content)
		req.Header.Set("Content-Type", contentType)
		rec := httptest.NewRecorder()

		mJWT := new(appmock.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mRepo := new(appmock.MockGroupRepository)
		uploadUseCase := group.NewUploadUseCase(mRepo)

		UploadGroupFile(mJWT, uploadUseCase).ServeHTTP(rec, req)

		//t.Log(rec.Result().Status)
		//t.Log(rec.Body.String())
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.UploadGroupFile return with StatusNotFound", func(t *testing.T) {
		config.SetMaxUploadSize(1024) // KB

		contentType, content, err := imageutil.CreateImageBuffer("jpg")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, groupResource, content)
		req.Header.Add("Content-Type", contentType)
		rec := httptest.NewRecorder()

		mJWT := new(appmock.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()
		mRepo := new(appmock.MockGroupRepository)
		mRepo.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(false, nil).
			Once()
		uploadUseCase := group.NewUploadUseCase(mRepo)

		UploadGroupFile(mJWT, uploadUseCase).ServeHTTP(rec, req)

		t.Log(rec.Result().Status)
		t.Log(rec.Body.String())
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.UploadGroupFile return with StatusUnauthorized", func(t *testing.T) {
		config.SetMaxUploadSize(1024) // KB

		contentType, content, err := imageutil.CreateImageBuffer("jpg")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, groupResource, content)
		req.Header.Add("Content-Type", contentType)
		rec := httptest.NewRecorder()

		mJWT := new(appmock.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()
		mRepo := new(appmock.MockGroupRepository)
		mRepo.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		mRepo.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uploadUseCase := group.NewUploadUseCase(mRepo)
		UploadGroupFile(mJWT, uploadUseCase).ServeHTTP(rec, req)

		t.Log(rec.Result().Status)
		t.Log(rec.Body.String())
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when handler.UploadGroupFile return with StatusBadRequest", func(t *testing.T) {
		config.SetMaxUploadSize(2) // KB

		contentType, content, err := imageutil.CreateImageBuffer("jpg")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, groupResource, content)
		req.Header.Add("Content-Type", contentType)
		rec := httptest.NewRecorder()

		mJWT := new(appmock.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()
		mRepo := new(appmock.MockGroupRepository)
		mRepo.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		mRepo.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		uploadUseCase := group.NewUploadUseCase(mRepo)

		UploadGroupFile(mJWT, uploadUseCase).ServeHTTP(rec, req)

		t.Log(rec.Result().Status)
		t.Log(rec.Body.String())
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when handler.UploadGroupFile return with StatusUnsupportedMediaType", func(t *testing.T) {
		config.SetMaxUploadSize(1024) // KB

		contentType, content, err := imageutil.CreateImageBuffer("png")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, groupResource, content)
		req.Header.Add("Content-Type", contentType)
		rec := httptest.NewRecorder()

		mJWT := new(appmock.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()
		mRepo := new(appmock.MockGroupRepository)
		mRepo.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		mRepo.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		uploadUseCase := group.NewUploadUseCase(mRepo)
		UploadGroupFile(mJWT, uploadUseCase).ServeHTTP(rec, req)

		t.Log(rec.Result().Status)
		t.Log(rec.Body.String())
		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when handler.UploadGroupFile return with StatusCreated", func(t *testing.T) {
		config.SetMaxUploadSize(1024) // KB

		contentType, content, err := imageutil.CreateImageBuffer("jpg")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, groupResource, content)
		req.Header.Add("Content-Type", contentType)
		rec := httptest.NewRecorder()

		mJWT := new(appmock.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518977777777", nil).
			Once()
		mRepo := new(appmock.MockGroupRepository)
		mRepo.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		mRepo.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		uploadUseCase := group.NewUploadUseCase(mRepo)
		UploadGroupFile(mJWT, uploadUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}
