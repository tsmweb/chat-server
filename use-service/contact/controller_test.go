package contact

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/use-service/common"
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
		qry := req.URL.Query()
		qry.Add("id", "+5518977777777")
		req.URL.RawQuery = qry.Encode()
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return(nil, errors.New("error")).
			Once()
		mService := new(mockService)
		mService.On("Get", mock.Anything, mock.Anything).
			Return(&Contact{}, nil).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Get().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, resource, nil)
		qry := req.URL.Query()
		qry.Add("id", "+5518977777777")
		req.URL.RawQuery = qry.Encode()
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Get", mock.Anything, mock.Anything).
			Return(nil, ErrContactNotFound).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Get().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, resource, nil)
		qry := req.URL.Query()
		qry.Add("id", "+5518977777777")
		req.URL.RawQuery = qry.Encode()
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Get", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Get().ServeHTTP(rec, req)

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

		req := httptest.NewRequest(http.MethodGet, resource, nil)
		qry := req.URL.Query()
		qry.Add("id", "+5518977777777")
		req.URL.RawQuery = qry.Encode()
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Get", mock.Anything, mock.Anything).
			Return(contact, nil).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.Get().ServeHTTP(rec, req)

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
		mJWT.On("GetDataToken", req, "id").
			Return(nil, errors.New("error")).
			Once()
		mService := new(mockService)
		mService.On("GetAll", mock.Anything).
			Return([]Contact{}, nil).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.GetAll().ServeHTTP(rec, req)

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
		mService.On("GetAll", mock.Anything).
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
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("GetAll", mock.Anything).
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

		var ps []Presenter
		for _, contact := range contacts {
			p := Presenter{}
			p.FromEntity(contact)
			ps = append(ps, p)
		}

		cj, err := json.Marshal(ps)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodGet, resource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", req, "id").
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("GetAll", mock.Anything).
			Return(contacts, nil).
			Once()

		ctrl := NewController(mJWT, mService)
		ctrl.GetAll().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(cj), rec.Body.String())
	})
}
