package login

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/tsmweb/auth-service/common"
	"github.com/tsmweb/auth-service/user"
	"github.com/tsmweb/go-helper-api/cerror"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewController(t *testing.T) {
	//t.Parallel()
	c := NewController(
		new(common.MockJWT),
		new(mockService))

	assert.NotNil(t, c)
}

func TestController_Login(t *testing.T) {
	//t.Parallel()

	t.Run("when controller return StatusUnsupportedMediaType", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mService := new(mockService)

		ctrl := NewController(mJWT, mService)
		ctrl.Login().ServeHTTP(rec, req)

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
		ctrl.Login().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when controller return StatusBadRequest", func(t *testing.T) {
		//t.Parallel()
		vm := &Presenter{
			ID: "+5518999999999",
			Password: "",
		}

		vmj, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader(vmj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mService := new(mockService)
		mService.On("Login", vm.ID, vm.Password).
			Return("", ErrPasswordValidateModel).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Login().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when controller return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		vm := &Presenter{
			ID: "+5518999999999",
			Password: "123456",
		}

		vmj, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader(vmj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mService := new(mockService)
		mService.On("Login", vm.ID, vm.Password).
			Return("", user.ErrUserNotFound).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Login().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when controller return StatusUnauthorized", func(t *testing.T) {
		//t.Parallel()
		vm := &Presenter{
			ID: "+5518999999999",
			Password: "123456",
		}

		vmj, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader(vmj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mService := new(mockService)
		mService.On("Login", vm.ID, vm.Password).
			Return("", cerror.ErrUnauthorized).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Login().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		vm := &Presenter{
			ID: "+5518999999999",
			Password: "123456",
		}

		vmj, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader(vmj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mService := new(mockService)
		mService.On("Login", vm.ID, vm.Password).
			Return("", errors.New("error")).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Login().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusOK", func(t *testing.T) {
		//t.Parallel()
		vm := &Presenter{
			ID: "+5518999999999",
			Password: "123456",
		}

		vmj, err := json.Marshal(vm)
		assert.Nil(t, err)

		token := TokenAuth{
			Token: "A1B2C3D4E5F6",
		}

		tokenj, err := json.Marshal(token)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader(vmj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mService := new(mockService)
		mService.On("Login", vm.ID, vm.Password).
			Return("A1B2C3D4E5F6", nil).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Login().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(tokenj), rec.Body.String())
		//t.Log(rec.Body.String())
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

	t.Run("when JWT fails with Error", func(t *testing.T) {
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
		vm := &Presenter{
			ID: "+5518977777777",
			Password: "123456",
		}

		vmj, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader(vmj))
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
		vm := &Presenter{
			ID: "+5518999999999",
			Password: "",
		}

		vmj, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader(vmj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Update", vm.ToEntity()).
			Return(ErrPasswordValidateModel).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when controller return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		vm := &Presenter{
			ID: "+5518999999999",
			Password: "123456",
		}

		vmj, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader(vmj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Update", vm.ToEntity()).
			Return(user.ErrUserNotFound).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		vm := &Presenter{
			ID: "+5518999999999",
			Password: "123456",
		}

		vmj, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader(vmj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Update", vm.ToEntity()).
			Return(errors.New("error")).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusOK", func(t *testing.T) {
		//t.Parallel()
		vm := &Presenter{
			ID: "+5518999999999",
			Password: "123456",
		}

		vmj, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader(vmj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Update", vm.ToEntity()).
			Return(nil).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
