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

func TestHandleWS(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, chatResource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		HandleWS(mJWT, nil).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when Router return StatusUnauthorized", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, chatResource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(MockJWT)
		mJWT.On("ExtractToken", mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		router := mux.NewRouter()
		MakeChatRouter(router, mJWT, middleware.NewAuth(mJWT), nil)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
