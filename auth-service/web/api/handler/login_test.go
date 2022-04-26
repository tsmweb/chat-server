package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/auth-service/app/login"
	"github.com/tsmweb/auth-service/common"
	"github.com/tsmweb/auth-service/web/api/dto"
	"github.com/tsmweb/go-helper-api/cerror"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Login(t *testing.T) {
	//t.Parallel()

	t.Run("when handler.Login return StatusUnsupportedMediaType", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPost, loginResource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		mLoginUseCase := new(mockLoginUseCase)
		Login(mLoginUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when handler.Login return StatusUnprocessableEntity", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPost, loginResource, bytes.NewReader([]byte("{[}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mLoginUseCase := new(mockLoginUseCase)
		Login(mLoginUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when handler.Login return StatusBadRequest", func(t *testing.T) {
		//t.Parallel()
		loginDto := &dto.Login{
			ID:       "+5518999999999",
			Password: "",
		}

		jLoginDto, err := json.Marshal(loginDto)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, loginResource, bytes.NewReader(jLoginDto))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mLoginUseCase := new(mockLoginUseCase)
		mLoginUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return("", login.ErrPasswordValidateModel).
			Once()

		Login(mLoginUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when handler.Login return StatusUnauthorized", func(t *testing.T) {
		//t.Parallel()
		loginDto := &dto.Login{
			ID:       "+5518999999999",
			Password: "123456",
		}

		jLoginDto, err := json.Marshal(loginDto)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, loginResource, bytes.NewReader(jLoginDto))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mLoginUseCase := new(mockLoginUseCase)
		mLoginUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return("", cerror.ErrUnauthorized).
			Once()

		Login(mLoginUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when handler.Login return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		loginDto := &dto.Login{
			ID:       "+5518999999999",
			Password: "123456",
		}

		jLoginDto, err := json.Marshal(loginDto)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, loginResource, bytes.NewReader(jLoginDto))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mLoginUseCase := new(mockLoginUseCase)
		mLoginUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return("", errors.New("error")).
			Once()

		Login(mLoginUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.Login return StatusOK", func(t *testing.T) {
		//t.Parallel()
		loginDto := &dto.Login{
			ID:       "+5518999999999",
			Password: "123456",
		}

		jLoginDto, err := json.Marshal(loginDto)
		assert.Nil(t, err)

		token := dto.TokenAuth{
			Token: "A1B2C3D4E5F6",
		}

		jToken, err := json.Marshal(token)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, loginResource, bytes.NewReader(jLoginDto))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mLoginUseCase := new(mockLoginUseCase)
		mLoginUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return("A1B2C3D4E5F6", nil).
			Once()

		Login(mLoginUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(jToken), rec.Body.String())
		//t.Log(rec.Body.String())
	})
}

func TestHandler_UpdatePassword(t *testing.T) {
	//t.Parallel()

	t.Run("when handler.UpdatePassword return StatusUnsupportedMediaType", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, loginResource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mUpdateUseCase := new(mockLoginUpdateUseCase)
		UpdatePassword(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when JWT fails with Error", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, loginResource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("jwt error")).
			Once()
		mUpdateUseCase := new(mockLoginUpdateUseCase)
		UpdatePassword(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.UpdatePassword return StatusUnprocessableEntity", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, loginResource, bytes.NewReader([]byte("{[}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockLoginUpdateUseCase)
		UpdatePassword(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when handler.UpdatePassword return StatusUnauthorized", func(t *testing.T) {
		//t.Parallel()
		loginDto := &dto.Login{
			ID:       "+5518977777777",
			Password: "123456",
		}

		jLoginDto, err := json.Marshal(loginDto)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, loginResource, bytes.NewReader(jLoginDto))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockLoginUpdateUseCase)
		mUpdateUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(login.ErrOperationNotAllowed).
			Once()

		UpdatePassword(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when handler.UpdatePassword return StatusBadRequest", func(t *testing.T) {
		//t.Parallel()
		loginDto := &dto.Login{
			ID:       "+5518999999999",
			Password: "",
		}

		jLoginDto, err := json.Marshal(loginDto)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, loginResource, bytes.NewReader(jLoginDto))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockLoginUpdateUseCase)
		mUpdateUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(login.ErrPasswordValidateModel).
			Once()

		UpdatePassword(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when handler.UpdatePassword return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		loginDto := &dto.Login{
			ID:       "+5518999999999",
			Password: "123456",
		}

		jLoginDto, err := json.Marshal(loginDto)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, loginResource, bytes.NewReader(jLoginDto))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockLoginUpdateUseCase)
		mUpdateUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(login.ErrUserNotFound).
			Once()

		UpdatePassword(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.UpdatePassword return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		loginDto := &dto.Login{
			ID:       "+5518999999999",
			Password: "123456",
		}

		jLoginDto, err := json.Marshal(loginDto)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, loginResource, bytes.NewReader(jLoginDto))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockLoginUpdateUseCase)
		mUpdateUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()

		UpdatePassword(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.UpdatePassword return StatusOK", func(t *testing.T) {
		//t.Parallel()
		loginDto := &dto.Login{
			ID:       "+5518999999999",
			Password: "123456",
		}

		jLoginDto, err := json.Marshal(loginDto)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, loginResource, bytes.NewReader(jLoginDto))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockLoginUpdateUseCase)
		mUpdateUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		UpdatePassword(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
