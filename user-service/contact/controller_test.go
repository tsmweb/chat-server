package contact

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/user-service/common"
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
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/+5518977777777", resource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mService := new(mockService)
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", resource), ctrl.Get()).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/+5518977777777", resource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Get", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, ErrContactNotFound).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", resource), ctrl.Get()).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/+5518977777777", resource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Get", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", resource), ctrl.Get()).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusOK", func(t *testing.T) {
		//t.Parallel()
		contact := &Contact{
			ID: "+5518977777777",
			Name: "Bill",
			LastName: "Gates",
		}

		p := Presenter{}
		p.FromEntity(contact)

		cj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/+5518977777777", resource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Get", mock.Anything, mock.Anything, mock.Anything).
			Return(contact, nil).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", resource), ctrl.Get()).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(cj), rec.Body.String())
	})
}

func TestController_GetAll(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, resource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mService := new(mockService)
		ctrl := NewController(mJWT, mService)
		ctrl.GetAll().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, resource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("GetAll", mock.Anything, mock.Anything).
			Return(nil, ErrContactNotFound).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.GetAll().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, resource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("GetAll", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.GetAll().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusOK", func(t *testing.T) {
		//t.Parallel()
		contacts := []*Contact {
			{
				ID: "+5518977777777",
				Name: "Bill",
				LastName: "Gates",
			},
			{
				ID: "+5518966666666",
				Name: "Elon",
				LastName: "Musk",
			},
		}

		p := EntityToPresenters(contacts...)
		cj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodGet, resource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("GetAll", mock.Anything, mock.Anything).
			Return(contacts, nil).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.GetAll().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(cj), rec.Body.String())
	})
}

func TestController_GetPresence(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/presence", resource)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/+5518977777777", path), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mService := new(mockService)
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", path), ctrl.GetPresence()).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/presence", resource)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/+5518977777777", path), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("GetPresence", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, ErrContactNotFound).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", path), ctrl.GetPresence()).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/presence", resource)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/+5518977777777", path), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("GetPresence", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", path), ctrl.GetPresence()).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusOK", func(t *testing.T) {
		//t.Parallel()
		var presence PresenceType = Online
		p := &Presence{
			ID: "+5518977777777",
			Presence: "online",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		path := fmt.Sprintf("%s/presence", resource)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/+5518977777777", path), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("GetPresence", mock.Anything, mock.Anything, mock.Anything).
			Return(presence, nil).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", path), ctrl.GetPresence()).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(pj), rec.Body.String())
		//t.Log(rec.Body.String())
	})

}

func TestController_Create(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mService := new(mockService)
		ctrl := NewController(mJWT, mService)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

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
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		ctrl := NewController(mJWT, mService)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when controller return StatusBadRequest", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID: "",
			Name: "Bill",
			LastName: "Gates",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(ErrIDValidateModel).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when controller return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID: "+5518977777777",
			Name: "Bill",
			LastName: "Gates",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(ErrUserNotFound).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when controller return StatusConflict", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID: "+5518977777777",
			Name: "Bill",
			LastName: "Gates",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(ErrContactAlreadyExists).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID: "+5518977777777",
			Name: "Bill",
			LastName: "Gates",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusCreated", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID: "+5518977777777",
			Name: "Bill",
			LastName: "Gates",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestController_Update(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mService := new(mockService)
		ctrl := NewController(mJWT, mService)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

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

	t.Run("when controller return StatusUnprocessableEntity", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader([]byte("{[}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		ctrl := NewController(mJWT, mService)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when controller return StatusBadRequest", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID: "",
			Name: "Bill",
			LastName: "Gates",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Update", mock.Anything, mock.Anything).
			Return(ErrIDValidateModel).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when controller return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID: "+5518977777777",
			Name: "Bill",
			LastName: "Gates",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Update", mock.Anything, mock.Anything).
			Return(ErrContactNotFound).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID: "+5518977777777",
			Name: "Bill",
			LastName: "Gates",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Update", mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusOK", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID: "+5518977777777",
			Name: "Bill",
			LastName: "Gates",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Update", mock.Anything, mock.Anything).
			Return(nil).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestController_Delete(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/+5518977777777", resource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mService := new(mockService)
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", resource), ctrl.Delete()).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/+5518977777777", resource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Delete", mock.Anything, mock.Anything, mock.Anything).
			Return(ErrContactNotFound).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", resource), ctrl.Delete()).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/+5518977777777", resource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Delete", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", resource), ctrl.Delete()).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusOK", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/+5518977777777", resource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Delete", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", resource), ctrl.Delete()).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestController_Block(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/block", resource)
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mService := new(mockService)
		ctrl := NewController(mJWT, mService)
		ctrl.Block().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusUnsupportedMediaType", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/block", resource)
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mService := new(mockService)
		ctrl := NewController(mJWT, mService)
		ctrl.Block().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when controller return StatusUnprocessableEntity", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/block", resource)
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader([]byte("{[}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		ctrl := NewController(mJWT, mService)
		ctrl.Block().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when controller return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID: "+5518977777777",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		path := fmt.Sprintf("%s/block", resource)
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Block", mock.Anything, mock.Anything, mock.Anything).
			Return(ErrUserNotFound).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.Block().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when controller return StatusConflict", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID: "+5518977777777",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		path := fmt.Sprintf("%s/block", resource)
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Block", mock.Anything, mock.Anything, mock.Anything).
			Return(ErrContactAlreadyBlocked).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.Block().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID: "+5518977777777",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		path := fmt.Sprintf("%s/block", resource)
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Block", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.Block().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusOK", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID: "+5518977777777",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		path := fmt.Sprintf("%s/block", resource)
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Block", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.Block().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestController_Unblock(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/block", resource)
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/+5518977777777", path), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mService := new(mockService)
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", path), ctrl.Unblock()).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/block", resource)
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/+5518977777777", path), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Unblock", mock.Anything, mock.Anything, mock.Anything).
			Return(ErrUserNotFound).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", path), ctrl.Unblock()).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/block", resource)
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/+5518977777777", path), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Unblock", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", path), ctrl.Unblock()).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusOK", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/block", resource)
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/+5518977777777", path), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Unblock", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", path), ctrl.Unblock()).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}