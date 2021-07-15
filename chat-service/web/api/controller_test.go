package api

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/go-helper-api/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestController_Connect(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, resource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		ctrl := NewController(mJWT, nil)
		ctrl.Connect().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when Router return StatusUnauthorized", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, resource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(MockJWT)
		mJWT.On("ExtractToken", mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		ctrl := NewController(mJWT, nil)

		router := mux.NewRouter()
		server := NewRouter(middleware.NewAuth(mJWT), ctrl)
		server.MakeRouters(router)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
