package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/user-service/app/contact"
	"github.com/tsmweb/user-service/common"
	"github.com/tsmweb/user-service/web/api/dto"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetContact(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/+5518977777777", contactResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mGetUseCase := new(mockContactGetUseCase)

		handler := GetContact(mJWT, mGetUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", contactResource), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.GetContacts return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/+5518977777777", contactResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mGetUseCase := new(mockContactGetUseCase)
		mGetUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, contact.ErrContactNotFound).
			Once()

		handler := GetContact(mJWT, mGetUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", contactResource), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.GetContacts return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/+5518977777777", contactResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mGetUseCase := new(mockContactGetUseCase)
		mGetUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		handler := GetContact(mJWT, mGetUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", contactResource), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.GetContacts return StatusOK", func(t *testing.T) {
		//t.Parallel()
		contact := &contact.Contact{
			ID:       "+5518977777777",
			Name:     "Bill",
			LastName: "Gates",
		}

		p := dto.Contact{}
		p.FromEntity(contact)

		cj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/+5518977777777", contactResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mGetUseCase := new(mockContactGetUseCase)
		mGetUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(contact, nil).
			Once()

		handler := GetContact(mJWT, mGetUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", contactResource), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(cj), rec.Body.String())
	})
}

func TestHandler_GetAllContacts(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, contactResource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mGetAllUseCase := new(mockContactGetAllUseCase)

		GetAllContacts(mJWT, mGetAllUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.GetAllContacts return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, contactResource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mGetAllUseCase := new(mockContactGetAllUseCase)
		mGetAllUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(nil, contact.ErrContactNotFound).
			Once()

		GetAllContacts(mJWT, mGetAllUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.GetAllContacts return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, contactResource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mGetAllUseCase := new(mockContactGetAllUseCase)
		mGetAllUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		GetAllContacts(mJWT, mGetAllUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.GetAllContacts return StatusOK", func(t *testing.T) {
		//t.Parallel()
		contacts := []*contact.Contact{
			{
				ID:       "+5518977777777",
				Name:     "Bill",
				LastName: "Gates",
			},
			{
				ID:       "+5518966666666",
				Name:     "Elon",
				LastName: "Musk",
			},
		}

		p := dto.EntityToContactDTO(contacts...)
		cj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodGet, contactResource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mGetAllUseCase := new(mockContactGetAllUseCase)
		mGetAllUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(contacts, nil).
			Once()

		GetAllContacts(mJWT, mGetAllUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(cj), rec.Body.String())
	})
}

func TestHandler_GetContactPresence(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/presence", contactResource)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/+5518977777777", path), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mGetPresenceUseCase := new(mockContactGetPresenceUseCase)

		handler := GetContactPresence(mJWT, mGetPresenceUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", path), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.GetContactPresence return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/presence", contactResource)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/+5518977777777", path), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mGetPresenceUseCase := new(mockContactGetPresenceUseCase)
		mGetPresenceUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, contact.ErrContactNotFound).
			Once()

		handler := GetContactPresence(mJWT, mGetPresenceUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", path), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.GetContactPresence return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/presence", contactResource)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/+5518977777777", path), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mGetPresenceUseCase := new(mockContactGetPresenceUseCase)
		mGetPresenceUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		handler := GetContactPresence(mJWT, mGetPresenceUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", path), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.GetContactPresence return StatusOK", func(t *testing.T) {
		//t.Parallel()
		var presence contact.PresenceType = contact.Online
		p := &dto.Presence{
			ID:       "+5518977777777",
			Presence: "online",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		path := fmt.Sprintf("%s/presence", contactResource)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/+5518977777777", path), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mGetPresenceUseCase := new(mockContactGetPresenceUseCase)
		mGetPresenceUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(presence, nil).
			Once()

		handler := GetContactPresence(mJWT, mGetPresenceUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", path), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(pj), rec.Body.String())
		//t.Log(rec.Body.String())
	})

}

func TestHandler_CreateContact(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPost, contactResource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		mCreateUseCase := new(mockContactCreateUseCase)

		CreateContact(mJWT, mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.CreateContact return StatusUnsupportedMediaType", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPost, contactResource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mCreateUseCase := new(mockContactCreateUseCase)

		CreateContact(mJWT, mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when handler.CreateContact return StatusUnprocessableEntity", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPost, contactResource, bytes.NewReader([]byte("{[}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mCreateUseCase := new(mockContactCreateUseCase)

		CreateContact(mJWT, mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when handler.CreateContact return StatusBadRequest", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Contact{
			ID:       "",
			Name:     "Bill",
			LastName: "Gates",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, contactResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mCreateUseCase := new(mockContactCreateUseCase)
		mCreateUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(contact.ErrIDValidateModel).
			Once()

		CreateContact(mJWT, mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when handler.CreateContact return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Contact{
			ID:       "+5518977777777",
			Name:     "Bill",
			LastName: "Gates",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, contactResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mCreateUseCase := new(mockContactCreateUseCase)
		mCreateUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(contact.ErrUserNotFound).
			Once()

		CreateContact(mJWT, mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.CreateContact return StatusConflict", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Contact{
			ID:       "+5518977777777",
			Name:     "Bill",
			LastName: "Gates",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, contactResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mCreateUseCase := new(mockContactCreateUseCase)
		mCreateUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(contact.ErrContactAlreadyExists).
			Once()

		CreateContact(mJWT, mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("when handler.CreateContact return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Contact{
			ID:       "+5518977777777",
			Name:     "Bill",
			LastName: "Gates",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, contactResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mCreateUseCase := new(mockContactCreateUseCase)
		mCreateUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()

		CreateContact(mJWT, mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.CreateContact return StatusCreated", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Contact{
			ID:       "+5518977777777",
			Name:     "Bill",
			LastName: "Gates",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, contactResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mCreateUseCase := new(mockContactCreateUseCase)
		mCreateUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		CreateContact(mJWT, mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestHandler_UpdateContact(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, contactResource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mUpdateUseCase := new(mockContactUpdateUseCase)

		UpdateContact(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.UpdateContact return StatusUnsupportedMediaType", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, contactResource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mUpdateUseCase := new(mockContactUpdateUseCase)

		UpdateContact(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when handler.UpdateContact return StatusUnprocessableEntity", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, contactResource, bytes.NewReader([]byte("{[}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockContactUpdateUseCase)

		UpdateContact(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when handler.UpdateContact return StatusBadRequest", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Contact{
			ID:       "",
			Name:     "Bill",
			LastName: "Gates",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, contactResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockContactUpdateUseCase)
		mUpdateUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(contact.ErrIDValidateModel).
			Once()

		UpdateContact(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when handler.UpdateContact return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Contact{
			ID:       "+5518977777777",
			Name:     "Bill",
			LastName: "Gates",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, contactResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockContactUpdateUseCase)
		mUpdateUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(contact.ErrContactNotFound).
			Once()

		UpdateContact(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.UpdateContact return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Contact{
			ID:       "+5518977777777",
			Name:     "Bill",
			LastName: "Gates",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, contactResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockContactUpdateUseCase)
		mUpdateUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()

		UpdateContact(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.UpdateContact return StatusOK", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Contact{
			ID:       "+5518977777777",
			Name:     "Bill",
			LastName: "Gates",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, contactResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockContactUpdateUseCase)
		mUpdateUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		UpdateContact(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestHandler_DeleteContact(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/+5518977777777", contactResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mDeleteUseCase := new(mockContactDeleteUseCase)

		handler := DeleteContact(mJWT, mDeleteUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", contactResource), handler).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.DeleteContact return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/+5518977777777", contactResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mDeleteUseCase := new(mockContactDeleteUseCase)
		mDeleteUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(contact.ErrContactNotFound).
			Once()

		handler := DeleteContact(mJWT, mDeleteUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", contactResource), handler).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.DeleteContact return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/+5518977777777", contactResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mDeleteUseCase := new(mockContactDeleteUseCase)
		mDeleteUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()

		handler := DeleteContact(mJWT, mDeleteUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", contactResource), handler).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.DeleteContact return StatusOK", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/+5518977777777", contactResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mDeleteUseCase := new(mockContactDeleteUseCase)
		mDeleteUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		handler := DeleteContact(mJWT, mDeleteUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", contactResource), handler).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestHandler_BlockContact(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/block", contactResource)
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mBlockUseCase := new(mockContactBlockUseCase)

		BlockContact(mJWT, mBlockUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.BlockContact return StatusUnsupportedMediaType", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/block", contactResource)
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mBlockUseCase := new(mockContactBlockUseCase)

		BlockContact(mJWT, mBlockUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when handler.BlockContact return StatusUnprocessableEntity", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/block", contactResource)
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader([]byte("{[}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mBlockUseCase := new(mockContactBlockUseCase)

		BlockContact(mJWT, mBlockUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when handler.BlockContact return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Contact{
			ID: "+5518977777777",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		path := fmt.Sprintf("%s/block", contactResource)
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mBlockUseCase := new(mockContactBlockUseCase)
		mBlockUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(contact.ErrUserNotFound).
			Once()

		BlockContact(mJWT, mBlockUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.BlockContact return StatusConflict", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Contact{
			ID: "+5518977777777",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		path := fmt.Sprintf("%s/block", contactResource)
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mBlockUseCase := new(mockContactBlockUseCase)
		mBlockUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(contact.ErrContactAlreadyBlocked).
			Once()

		BlockContact(mJWT, mBlockUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("when handler.BlockContact return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Contact{
			ID: "+5518977777777",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		path := fmt.Sprintf("%s/block", contactResource)
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mBlockUseCase := new(mockContactBlockUseCase)
		mBlockUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()

		BlockContact(mJWT, mBlockUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.BlockContact return StatusOK", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Contact{
			ID: "+5518977777777",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		path := fmt.Sprintf("%s/block", contactResource)
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mBlockUseCase := new(mockContactBlockUseCase)
		mBlockUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		BlockContact(mJWT, mBlockUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestHandler_UnblockContact(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/block", contactResource)
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/+5518977777777", path), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mUnblockUseCase := new(mockContactUnblockUseCase)

		handler := UnblockContact(mJWT, mUnblockUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", path), handler).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.UnblockContact return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/block", contactResource)
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/+5518977777777", path), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUnblockUseCase := new(mockContactUnblockUseCase)
		mUnblockUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(contact.ErrUserNotFound).
			Once()

		handler := UnblockContact(mJWT, mUnblockUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", path), handler).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.UnblockContact return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/block", contactResource)
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/+5518977777777", path), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUnblockUseCase := new(mockContactUnblockUseCase)
		mUnblockUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()

		handler := UnblockContact(mJWT, mUnblockUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", path), handler).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.UnblockContact return StatusOK", func(t *testing.T) {
		//t.Parallel()
		path := fmt.Sprintf("%s/block", contactResource)
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/+5518977777777", path), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUnblockUseCase := new(mockContactUnblockUseCase)
		mUnblockUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		handler := UnblockContact(mJWT, mUnblockUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", path), handler).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
