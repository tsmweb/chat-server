package profile

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/tsmweb/auth-service/helper/common"
	"github.com/tsmweb/go-helper-api/cerror"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewController(t *testing.T) {
	c := NewController(
			new(common.MockJWT),
			new(mockGetUseCase),
			new(mockCreateUseCase),
			new(mockUpdateUseCase))

	assert.NotNil(t, c)
}

func TestController_Get(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		mGet := new(mockGetUseCase)
		mGet.On("Execute", "+5518999999999").
			Return(Profile{}, nil).
			Once()
		mCreate := new(mockCreateUseCase)
		mUpdate := new(mockUpdateUseCase)

		req := httptest.NewRequest(http.MethodGet, resource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return(nil, errors.New("error")).
			Once()

		ctrl := NewController(mJWT, mGet, mCreate, mUpdate)
		ctrl.Get().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		mGet := new(mockGetUseCase)
		mGet.On("Execute", "+5518999999999").
			Return(nil, cerror.ErrInternalServer).
			Once()
		mCreate := new(mockCreateUseCase)
		mUpdate := new(mockUpdateUseCase)

		req := httptest.NewRequest(http.MethodGet, resource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()

		ctrl := NewController(mJWT, mGet, mCreate, mUpdate)
		ctrl.Get().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return ErrProfileNotFound", func(t *testing.T) {
		//t.Parallel()
		mGet := new(mockGetUseCase)
		mGet.On("Execute", "+5518999999999").
			Return(nil, ErrProfileNotFound).
			Once()
		mCreate := new(mockCreateUseCase)
		mUpdate := new(mockUpdateUseCase)

		req := httptest.NewRequest(http.MethodGet, resource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()

		ctrl := NewController(mJWT, mGet, mCreate, mUpdate)
		ctrl.Get().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when controller return StatusOK", func(t *testing.T) {
		//t.Parallel()
		profile := Profile{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
		}

		mGet := new(mockGetUseCase)
		mGet.On("Execute", "+5518999999999").
			Return(profile, nil).
			Once()
		mCreate := new(mockCreateUseCase)
		mUpdate := new(mockUpdateUseCase)

		req := httptest.NewRequest(http.MethodGet, resource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()

		ctrl := NewController(mJWT, mGet, mCreate, mUpdate)
		ctrl.Get().ServeHTTP(rec, req)

		//t.Log(rec.Body)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestController_Create(t *testing.T) {
	//t.Parallel()

	t.Run("when controller return StatusUnsupportedMediaType", func(t *testing.T) {
		//t.Parallel()
		vm := ViewModel{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
			Password: "123456",
		}

		mGet := new(mockGetUseCase)
		mCreate := new(mockCreateUseCase)
		mCreate.On("Execute", vm.ID, vm.Name, vm.LastName, vm.Password).
			Return(nil).
			Once()
		mUpdate := new(mockUpdateUseCase)

		json, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader(json))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)

		ctrl := NewController(mJWT, mGet, mCreate, mUpdate)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when controller return StatusBadRequest", func(t *testing.T) {
		//t.Parallel()
		vm := ViewModel{
			ID:       "+5518999999999",
			LastName: "Jobs",
			Password: "123456",
		}

		mGet := new(mockGetUseCase)
		mCreate := new(mockCreateUseCase)
		mCreate.On("Execute", vm.ID, "", vm.LastName, vm.Password).
			Return(ErrNameValidateModel).
			Once()
		mUpdate := new(mockUpdateUseCase)

		jsonBody, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)

		ctrl := NewController(mJWT, mGet, mCreate, mUpdate)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when controller return StatusUnprocessableEntity", func(t *testing.T) {
		//t.Parallel()
		vm := ViewModel{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
			Password: "123456",
		}

		mGet := new(mockGetUseCase)
		mCreate := new(mockCreateUseCase)
		mCreate.On("Execute", vm.ID, vm.Name, vm.LastName, vm.Password).
			Return(nil).
			Once()
		mUpdate := new(mockUpdateUseCase)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader([]byte("{[}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)

		ctrl := NewController(mJWT, mGet, mCreate, mUpdate)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when controller return StatusConflict", func(t *testing.T) {
		//t.Parallel()
		vm := ViewModel{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
			Password: "123456",
		}

		mGet := new(mockGetUseCase)
		mCreate := new(mockCreateUseCase)
		mCreate.On("Execute", vm.ID, vm.Name, vm.LastName, vm.Password).
			Return(cerror.ErrRecordAlreadyRegistered).
			Once()
		mUpdate := new(mockUpdateUseCase)

		jsonBody, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)

		ctrl := NewController(mJWT, mGet, mCreate, mUpdate)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		vm := ViewModel{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
			Password: "123456",
		}

		mGet := new(mockGetUseCase)
		mCreate := new(mockCreateUseCase)
		mCreate.On("Execute", vm.ID, vm.Name, vm.LastName, vm.Password).
			Return(cerror.ErrInternalServer).
			Once()
		mUpdate := new(mockUpdateUseCase)

		jsonBody, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)

		ctrl := NewController(mJWT, mGet, mCreate, mUpdate)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusCreated", func(t *testing.T) {
		//t.Parallel()
		vm := ViewModel{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
			Password: "123456",
		}

		mGet := new(mockGetUseCase)
		mCreate := new(mockCreateUseCase)
		mCreate.On("Execute", vm.ID, vm.Name, vm.LastName, vm.Password).
			Return(nil).
			Once()
		mUpdate := new(mockUpdateUseCase)

		jsonBody, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)

		ctrl := NewController(mJWT, mGet, mCreate, mUpdate)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestController_Update(t *testing.T) {
	t.Run("when controller return StatusUnsupportedMediaType", func(t *testing.T) {
		//t.Parallel()
		vm := ViewModel{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
		}

		mGet := new(mockGetUseCase)
		mCreate := new(mockCreateUseCase)
		mUpdate := new(mockUpdateUseCase)
		mUpdate.On("Execute", vm.ToEntity()).
			Return(nil).
			Once()

		json, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader(json))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()

		ctrl := NewController(mJWT, mGet, mCreate, mUpdate)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		vm := ViewModel{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
		}

		mGet := new(mockGetUseCase)
		mCreate := new(mockCreateUseCase)
		mUpdate := new(mockUpdateUseCase)
		mUpdate.On("Execute", vm.ToEntity()).
			Return(nil).
			Once()

		json, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader(json))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return(nil, errors.New("error")).
			Once()

		ctrl := NewController(mJWT, mGet, mCreate, mUpdate)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusUnprocessableEntity", func(t *testing.T) {
		//t.Parallel()
		vm := ViewModel{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
			Password: "123456",
		}

		mGet := new(mockGetUseCase)
		mCreate := new(mockCreateUseCase)
		mUpdate := new(mockUpdateUseCase)
		mUpdate.On("Execute", vm.ToEntity()).
			Return(nil).
			Once()

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader([]byte("{[}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()

		ctrl := NewController(mJWT, mGet, mCreate, mUpdate)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when controller return StatusUnauthorized", func(t *testing.T) {
		//t.Parallel()
		vm := ViewModel{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
		}

		mGet := new(mockGetUseCase)
		mCreate := new(mockCreateUseCase)
		mUpdate := new(mockUpdateUseCase)
		mUpdate.On("Execute", vm.ToEntity()).
			Return(nil).
			Once()

		jsonBody, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5511988888888", nil).
			Once()

		ctrl := NewController(mJWT, mGet, mCreate, mUpdate)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when controller return StatusBadRequest", func(t *testing.T) {
		//t.Parallel()
		vm := ViewModel{
			ID:       "+5518999999999",
			LastName: "Jobs",
		}

		mGet := new(mockGetUseCase)
		mCreate := new(mockCreateUseCase)
		mUpdate := new(mockUpdateUseCase)
		mUpdate.On("Execute", vm.ToEntity()).
			Return(ErrNameValidateModel).
			Once()

		jsonBody, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()

		ctrl := NewController(mJWT, mGet, mCreate, mUpdate)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		vm := ViewModel{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
		}

		mGet := new(mockGetUseCase)
		mCreate := new(mockCreateUseCase)
		mUpdate := new(mockUpdateUseCase)
		mUpdate.On("Execute", vm.ToEntity()).
			Return(cerror.ErrInternalServer).
			Once()

		jsonBody, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()

		ctrl := NewController(mJWT, mGet, mCreate, mUpdate)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusOK", func(t *testing.T) {
		//t.Parallel()
		vm := ViewModel{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
		}

		mGet := new(mockGetUseCase)
		mCreate := new(mockCreateUseCase)
		mUpdate := new(mockUpdateUseCase)
		mUpdate.On("Execute", vm.ToEntity()).
			Return(nil).
			Once()

		jsonBody, err := json.Marshal(vm)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()

		ctrl := NewController(mJWT, mGet, mCreate, mUpdate)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}