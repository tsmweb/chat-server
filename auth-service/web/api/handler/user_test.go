package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/auth-service/common"
	"github.com/tsmweb/auth-service/user"
	"github.com/tsmweb/auth-service/web/api/dto"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetUser(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, userResource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mGetUseCase := new(mockUserGetUseCase)
		mGetUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(user.User{}, nil).
			Once()

		GetUser(mJWT, mGetUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.GetUser return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, userResource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mGetUseCase := new(mockUserGetUseCase)
		mGetUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		GetUser(mJWT, mGetUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.GetUser return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, userResource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mGetUseCase := new(mockUserGetUseCase)
		mGetUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(nil, user.ErrUserNotFound).
			Once()

		GetUser(mJWT, mGetUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.GetUser return StatusOK", func(t *testing.T) {
		//t.Parallel()
		user := &user.User{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
		}

		userDto := &dto.User{}
		userDto.FromEntity(user)

		jUserDto, err := json.Marshal(userDto)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodGet, userResource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mGetUseCase := new(mockUserGetUseCase)
		mGetUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(user, nil).
			Once()

		GetUser(mJWT, mGetUseCase).ServeHTTP(rec, req)

		//t.Log(rec.Body)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(jUserDto), rec.Body.String())
	})
}

func TestHandler_CreateUser(t *testing.T) {
	//t.Parallel()

	t.Run("when handler.CreateUser return StatusUnsupportedMediaType", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPost, userResource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		mCreateUseCase := new(mockUserCreateUseCase)
		CreateUser(mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when handler.CreateUser return StatusUnprocessableEntity", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPost, userResource, bytes.NewReader([]byte("{[}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mCreateUseCase := new(mockUserCreateUseCase)
		CreateUser(mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when handler.CreateUser return StatusBadRequest", func(t *testing.T) {
		//t.Parallel()
		userDto := &dto.User{
			ID:       "+5518999999999",
			LastName: "Jobs",
			Password: "123456",
		}

		jUserDto, err := json.Marshal(userDto)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, userResource, bytes.NewReader(jUserDto))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mCreateUseCase := new(mockUserCreateUseCase)
		mCreateUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(user.ErrNameValidateModel).
			Once()

		CreateUser(mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when handler.CreateUser return StatusConflict", func(t *testing.T) {
		//t.Parallel()
		userDto := &dto.User{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
			Password: "123456",
		}

		jUserDto, err := json.Marshal(userDto)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, userResource, bytes.NewReader(jUserDto))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mCreateUseCase := new(mockUserCreateUseCase)
		mCreateUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(user.ErrUserAlreadyExists).
			Once()

		CreateUser(mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("when handler.CreateUser return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		userDto := &dto.User{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
			Password: "123456",
		}

		jUserDto, err := json.Marshal(userDto)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, userResource, bytes.NewReader(jUserDto))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mCreateUseCase := new(mockUserCreateUseCase)
		mCreateUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()

		CreateUser(mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.CreateUser return StatusCreated", func(t *testing.T) {
		//t.Parallel()
		userDto := &dto.User{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
			Password: "123456",
		}

		jUserDto, err := json.Marshal(userDto)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, userResource, bytes.NewReader(jUserDto))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mCreateUseCase := new(mockUserCreateUseCase)
		mCreateUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		CreateUser(mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestHandler_UpdateUser(t *testing.T) {
	//t.Parallel()

	t.Run("when handler.UpdateUser return StatusUnsupportedMediaType", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, userResource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mUpdateUseCase := new(mockUserUpdateUseCase)
		UpdateUser(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, userResource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("jwt error")).
			Once()
		mUpdateUseCase := new(mockUserUpdateUseCase)

		UpdateUser(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.UpdateUser return StatusUnprocessableEntity", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, userResource, bytes.NewReader([]byte("{[}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockUserUpdateUseCase)

		UpdateUser(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when controller return StatusBadRequest", func(t *testing.T) {
		//t.Parallel()
		userDto := &dto.User{
			ID:       "+5518999999999",
			LastName: "Jobs",
		}

		jUserDto, err := json.Marshal(userDto)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, userResource, bytes.NewReader(jUserDto))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockUserUpdateUseCase)
		mUpdateUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(user.ErrNameValidateModel).
			Once()

		UpdateUser(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when controller return StatusUnauthorized", func(t *testing.T) {
		//t.Parallel()
		userDto := &dto.User{
			ID:       "+5518977777777",
			Name:     "Steve",
			LastName: "Jobs",
		}

		jUserDto, err := json.Marshal(userDto)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, userResource, bytes.NewReader(jUserDto))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockUserUpdateUseCase)
		mUpdateUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(user.ErrOperationNotAllowed).
			Once()

		UpdateUser(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		userDto := &dto.User{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
		}

		jUserDto, err := json.Marshal(userDto)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, userResource, bytes.NewReader(jUserDto))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockUserUpdateUseCase)
		mUpdateUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()

		UpdateUser(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusOK", func(t *testing.T) {
		//t.Parallel()
		userDto := &dto.User{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
		}

		jUserDto, err := json.Marshal(userDto)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, userResource, bytes.NewReader(jUserDto))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockUserUpdateUseCase)
		mUpdateUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		UpdateUser(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}