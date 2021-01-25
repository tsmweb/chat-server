package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/tsmweb/auth-service/common"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewController(t *testing.T) {
	//t.Parallel()
	c := NewController(new(common.MockJWT), new(mockService))
	assert.NotNil(t, c)
}

func TestController_Get(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, resource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return(nil, errors.New("error")).
			Once()
		mService := new(mockService)
		mService.On("Get", "+5518999999999").
			Return(User{}, nil).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Get().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, resource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Get", "+5518999999999").
			Return(nil, errors.New("error")).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Get().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, resource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Get", "+5518999999999").
			Return(nil, ErrUserNotFound).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Get().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when controller return StatusOK", func(t *testing.T) {
		//t.Parallel()
		user := &User{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
		}

		p := &Presenter{}
		p.FromEntity(user)

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodGet, resource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Get", "+5518999999999").
			Return(user, nil).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Get().ServeHTTP(rec, req)

		//t.Log(rec.Body)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(pj), rec.Body.String())
	})
}

func TestController_Create(t *testing.T) {
	//t.Parallel()

	t.Run("when controller return StatusUnsupportedMediaType", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mService := new(mockService)

		ctrl := NewController(mJWT, mService)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when controller return StatusUnprocessableEntity", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader([]byte("{[}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mService := new(mockService)

		ctrl := NewController(mJWT, mService)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when controller return StatusBadRequest", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID:       "+5518999999999",
			LastName: "Jobs",
			Password: "123456",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mService := new(mockService)
		mService.On("Create", p.ID, p.Name, p.LastName, p.Password).
			Return(ErrNameValidateModel).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when controller return StatusConflict", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
			Password: "123456",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mService := new(mockService)
		mService.On("Create", p.ID, p.Name, p.LastName, p.Password).
			Return(ErrUserAlreadyExists).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
			Password: "123456",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mService := new(mockService)
		mService.On("Create", p.ID, p.Name, p.LastName, p.Password).
			Return(errors.New("error")).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusCreated", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
			Password: "123456",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mService := new(mockService)
		mService.On("Create", p.ID, p.Name, p.LastName, p.Password).
			Return(nil).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestController_Update(t *testing.T) {
	//t.Parallel()

	t.Run("when controller return StatusUnsupportedMediaType", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mService := new(mockService)

		ctrl := NewController(mJWT, mService)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return(nil, errors.New("jwt error")).
			Once()
		mService := new(mockService)

		ctrl := NewController(mJWT, mService)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusUnprocessableEntity", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader([]byte("{[}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)

		ctrl := NewController(mJWT, mService)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when controller return StatusUnauthorized", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID:       "+5518977777777",
			Name:     "Steve",
			LastName: "Jobs",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)

		ctrl := NewController(mJWT, mService)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when controller return StatusBadRequest", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID:       "+5518999999999",
			LastName: "Jobs",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Update", p.ToEntity()).
			Return(ErrNameValidateModel).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Update", p.ToEntity()).
			Return(errors.New("error")).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusOK", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Update", p.ToEntity()).
			Return(nil).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}