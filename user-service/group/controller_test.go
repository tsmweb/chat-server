package group

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
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
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", resource), nil)
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
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", resource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Get", mock.Anything, mock.Anything).
			Return(nil, ErrGroupNotFound).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", resource), ctrl.Get()).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", resource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Get", mock.Anything, mock.Anything).
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
		group := &Group{
			ID: "be49afd2ee890805c21ddd55879db1387aec9751",
			Name: "Churrasco na Piscina",
			Description: "Amigos do churrasco.",
			Owner: "+5518999999999",
			Members: []*Member{
				{
					GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
					UserID: "+5518999999999",
					Admin: true,
				},
				{
					GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
					UserID: "+5518988888888",
					Admin: false,
				},
				{
					GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
					UserID: "+5518977777777",
					Admin: false,
				},
			},
		}

		p := Presenter{}
		p.FromEntity(group)

		gj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", resource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Get", mock.Anything, mock.Anything).
			Return(group, nil).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", resource), ctrl.Get()).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(gj), rec.Body.String())
		t.Log(string(gj))
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
			Return(nil, ErrGroupNotFound).
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
		groups := []*Group{
			{
				ID:          "be49afd2ee890805c21ddd55879db1387aec9751",
				Name:        "Churrasco na Piscina",
				Description: "Amigos do churrasco.",
				Owner:       "+5518999999999",
				Members: []*Member{
					{
						GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
						UserID:  "+5518999999999",
						Admin:   true,
					},
					{
						GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
						UserID:  "+5518988888888",
						Admin:   false,
					},
					{
						GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
						UserID:  "+5518977777777",
						Admin:   false,
					},
				},
			},
			{
				ID:          "be49afd2ee890805c21ddd55879db1387aec9752",
				Name:        "Grupo de Trabalho",
				Description: "Grupo de Trabalho (Exemplo)",
				Owner:       "+5518999999999",
				Members: []*Member{
					{
						GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
						UserID:  "+5518999999999",
						Admin:   true,
					},
					{
						GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
						UserID:  "+5518977777777",
						Admin:   false,
					},
				},
			},
		}

		p := EntityToPresenters(groups...)
		gj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodGet, resource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("GetAll", mock.Anything, mock.Anything).
			Return(groups, nil).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.GetAll().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(gj), rec.Body.String())
		//t.Log(rec.Body.String())
	})
}

func TestController_Create(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
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
			Name: "",
			Description: "Exemplo de Grupo",
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
		mService.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return("", ErrNameValidateModel).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			Name: "Grupo Teste",
			Description: "Exemplo de Grupo",
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
		mService.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return("", errors.New("error")).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusCreated", func(t *testing.T) {
		//t.Parallel()
		groupID := "be49afd2ee890805c21ddd55879db1387aec9751"
		location := fmt.Sprintf("%s/%s", resource, groupID)

		p := &Presenter{
			Name: "Grupo Teste",
			Description: "Exemplo de Grupo",
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
		mService.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(groupID, nil).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, location, rec.Header().Get("Location"))
		t.Log(rec.Header().Get("Location"))
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
			ID: "be49afd2ee890805c21ddd55879db1387aec9751",
			Name: "Grupo Teste",
			Description: "Exemplo de Grupo",
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

	t.Run("when controller return StatusUnauthorized", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID: "be49afd2ee890805c21ddd55879db1387aec9751",
			Name: "Grupo Teste",
			Description: "Exemplo de Grupo",
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
			Return(ErrOperationNotAllowed).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when controller return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID: "be49afd2ee890805c21ddd55879db1387aec9751",
			Name: "Grupo Teste",
			Description: "Exemplo de Grupo",
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
			Return(ErrGroupNotFound).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		p := &Presenter{
			ID: "be49afd2ee890805c21ddd55879db1387aec9751",
			Name: "Grupo Teste",
			Description: "Exemplo de Grupo",
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
			ID: "be49afd2ee890805c21ddd55879db1387aec9751",
			Name: "Grupo Teste",
			Description: "Exemplo de Grupo",
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
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", resource), nil)
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

	t.Run("when controller return StatusUnauthorized", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", resource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Delete", mock.Anything, mock.Anything).
			Return(ErrOperationNotAllowed).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", resource), ctrl.Delete()).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when controller return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", resource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Delete", mock.Anything, mock.Anything).
			Return(ErrGroupNotFound).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", resource), ctrl.Delete()).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", resource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Delete", mock.Anything, mock.Anything).
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
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", resource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("Delete", mock.Anything, mock.Anything).
			Return(nil).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", resource), ctrl.Delete()).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestController_AddMember(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, resourceMember, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mService := new(mockService)
		ctrl := NewController(mJWT, mService)
		ctrl.AddMember().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusUnsupportedMediaType", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPost, resourceMember, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mService := new(mockService)
		ctrl := NewController(mJWT, mService)
		ctrl.AddMember().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when controller return StatusUnprocessableEntity", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPost, resourceMember, bytes.NewReader([]byte("{[}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		ctrl := NewController(mJWT, mService)
		ctrl.AddMember().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when controller return StatusUnauthorized", func(t *testing.T) {
		//t.Parallel()
		p := &MemberPresenter{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID: "+5518988888888",
			Admin: false,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resourceMember, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("AddMember", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(ErrOperationNotAllowed).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.AddMember().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when controller return StatusBadRequest", func(t *testing.T) {
		//t.Parallel()
		p := &MemberPresenter{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID: "",
			Admin: false,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resourceMember, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("AddMember", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(ErrUserIDValidateModel).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.AddMember().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when controller return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		p := &MemberPresenter{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID: "+5518988888888",
			Admin: false,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resourceMember, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("AddMember", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(ErrGroupNotFound).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.AddMember().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when controller return StatusConflict", func(t *testing.T) {
		//t.Parallel()
		p := &MemberPresenter{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID: "+5518988888888",
			Admin: false,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resourceMember, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("AddMember", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(ErrMemberAlreadyExists).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.AddMember().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		p := &MemberPresenter{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID: "+5518988888888",
			Admin: false,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resourceMember, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("AddMember", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.AddMember().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusCreated", func(t *testing.T) {
		//t.Parallel()
		p := &MemberPresenter{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID: "+5518988888888",
			Admin: false,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, resourceMember, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("AddMember", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.AddMember().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestController_RemoveMember(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751/+5518988888888",
				resourceMember), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mService := new(mockService)
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{group}/{user}", resourceMember), ctrl.RemoveMember()).
			Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusUnauthorized", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751/+5518988888888",
				resourceMember), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("RemoveMember", mock.Anything, mock.Anything, mock.Anything).
			Return(ErrGroupOwnerCannotRemoved).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{group}/{user}", resourceMember), ctrl.RemoveMember()).
			Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when controller return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751/+5518988888888",
				resourceMember), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("RemoveMember", mock.Anything, mock.Anything, mock.Anything).
			Return(ErrMemberNotFound).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{group}/{user}", resourceMember), ctrl.RemoveMember()).
			Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751/+5518988888888",
				resourceMember), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("RemoveMember", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{group}/{user}", resourceMember), ctrl.RemoveMember()).
			Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusOK", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751/+5518988888888",
				resourceMember), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("RemoveMember", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		ctrl := NewController(mJWT, mService)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{group}/{user}", resourceMember), ctrl.RemoveMember()).
			Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestController_SetAdmin(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, resourceMember, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mService := new(mockService)
		ctrl := NewController(mJWT, mService)
		ctrl.SetAdmin().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusUnsupportedMediaType", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, resourceMember, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mService := new(mockService)
		ctrl := NewController(mJWT, mService)
		ctrl.SetAdmin().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when controller return StatusUnprocessableEntity", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, resourceMember, bytes.NewReader([]byte("{[}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		ctrl := NewController(mJWT, mService)
		ctrl.SetAdmin().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when controller return StatusBadRequest", func(t *testing.T) {
		//t.Parallel()
		p := &MemberPresenter{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID: "",
			Admin: true,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resourceMember, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("SetAdmin", mock.Anything, mock.Anything).
			Return(ErrUserIDValidateModel).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.SetAdmin().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when controller return StatusUnauthorized", func(t *testing.T) {
		//t.Parallel()
		p := &MemberPresenter{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID: "+5518977777777",
			Admin: true,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resourceMember, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("SetAdmin", mock.Anything, mock.Anything).
			Return(ErrOperationNotAllowed).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.SetAdmin().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when controller return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		p := &MemberPresenter{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID: "+5518977777777",
			Admin: true,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resourceMember, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("SetAdmin", mock.Anything, mock.Anything).
			Return(ErrMemberNotFound).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.SetAdmin().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when controller return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		p := &MemberPresenter{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID: "+5518977777777",
			Admin: true,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resourceMember, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("SetAdmin", mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.SetAdmin().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when controller return StatusOK", func(t *testing.T) {
		//t.Parallel()
		p := &MemberPresenter{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID: "+5518977777777",
			Admin: true,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, resourceMember, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mService := new(mockService)
		mService.On("SetAdmin", mock.Anything, mock.Anything).
			Return(nil).
			Once()
		ctrl := NewController(mJWT, mService)
		ctrl.SetAdmin().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}