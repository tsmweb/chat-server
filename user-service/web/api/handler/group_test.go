package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/user-service/common"
	"github.com/tsmweb/user-service/group"
	"github.com/tsmweb/user-service/web/api/dto"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetGroup(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", groupResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mGetUseCase := new(mockGroupGetUseCase)

		handler := GetGroup(mJWT, mGetUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", groupResource), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.GetGroup return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", groupResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mGetUseCase := new(mockGroupGetUseCase)
		mGetUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(nil, group.ErrGroupNotFound).
			Once()

		handler := GetGroup(mJWT, mGetUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", groupResource), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.GetGroup return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", groupResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mGetUseCase := new(mockGroupGetUseCase)
		mGetUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		handler := GetGroup(mJWT, mGetUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", groupResource), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.GetGroup return StatusOK", func(t *testing.T) {
		//t.Parallel()
		group := &group.Group{
			ID:          "be49afd2ee890805c21ddd55879db1387aec9751",
			Name:        "Churrasco na Piscina",
			Description: "Amigos do churrasco.",
			Owner:       "+5518999999999",
			Members: []*group.Member{
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
		}

		p := dto.Group{}
		p.FromEntity(group)

		gj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", groupResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mGetUseCase := new(mockGroupGetUseCase)
		mGetUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(group, nil).
			Once()

		handler := GetGroup(mJWT, mGetUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", groupResource), handler).Methods(http.MethodGet)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(gj), rec.Body.String())
		t.Log(string(gj))
	})

}

func TestHandler_GetAllGroups(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, groupResource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mGetAllUseCase := new(mockGroupGetAllUseCase)

		GetAllGroups(mJWT, mGetAllUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.GetAllGroups return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, groupResource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mGetAllUseCase := new(mockGroupGetAllUseCase)
		mGetAllUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(nil, group.ErrGroupNotFound).
			Once()

		GetAllGroups(mJWT, mGetAllUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.GetAllGroups return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodGet, groupResource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mGetAllUseCase := new(mockGroupGetAllUseCase)
		mGetAllUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		GetAllGroups(mJWT, mGetAllUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.GetAllGroups return StatusOK", func(t *testing.T) {
		//t.Parallel()
		groups := []*group.Group{
			{
				ID:          "be49afd2ee890805c21ddd55879db1387aec9751",
				Name:        "Churrasco na Piscina",
				Description: "Amigos do churrasco.",
				Owner:       "+5518999999999",
			},
			{
				ID:          "be49afd2ee890805c21ddd55879db1387aec9752",
				Name:        "Grupo de Trabalho",
				Description: "Grupo de Trabalho (Exemplo)",
				Owner:       "+5518999999999",
			},
		}

		p := dto.EntityToGroupDTO(groups...)
		gj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodGet, groupResource, nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mGetAllUseCase := new(mockGroupGetAllUseCase)
		mGetAllUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(groups, nil).
			Once()

		GetAllGroups(mJWT, mGetAllUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(gj), rec.Body.String())
		//t.Log(rec.Body.String())
	})
}

func TestHandler_CreateGroup(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, groupResource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mCreateUseCase := new(mockGroupCreateUseCase)

		CreateGroup(mJWT, mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.CreateGroup return StatusUnsupportedMediaType", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPost, groupResource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mCreateUseCase := new(mockGroupCreateUseCase)

		CreateGroup(mJWT, mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when handler.CreateGroup return StatusUnprocessableEntity", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPost, groupResource, bytes.NewReader([]byte("{[}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mCreateUseCase := new(mockGroupCreateUseCase)

		CreateGroup(mJWT, mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when handler.CreateGroup return StatusBadRequest", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Group{
			Name:        "",
			Description: "Exemplo de Grupo",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, groupResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mCreateUseCase := new(mockGroupCreateUseCase)
		mCreateUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return("", group.ErrNameValidateModel).
			Once()

		CreateGroup(mJWT, mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when handler.CreateGroup return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Group{
			Name:        "Grupo Teste",
			Description: "Exemplo de Grupo",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, groupResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mCreateUseCase := new(mockGroupCreateUseCase)
		mCreateUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return("", errors.New("error")).
			Once()

		CreateGroup(mJWT, mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.CreateGroup return StatusCreated", func(t *testing.T) {
		//t.Parallel()
		groupID := "be49afd2ee890805c21ddd55879db1387aec9751"
		location := fmt.Sprintf("%s/%s", groupResource, groupID)

		p := &dto.Group{
			Name:        "Grupo Teste",
			Description: "Exemplo de Grupo",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, groupResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mCreateUseCase := new(mockGroupCreateUseCase)
		mCreateUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(groupID, nil).
			Once()

		CreateGroup(mJWT, mCreateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, location, rec.Header().Get("Location"))
		t.Log(rec.Header().Get("Location"))
	})
}

func TestHandler_UpdateGroup(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, groupResource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mUpdateUseCase := new(mockGroupUpdateUseCase)

		UpdateGroup(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.UpdateGroup return StatusUnsupportedMediaType", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, groupResource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mUpdateUseCase := new(mockGroupUpdateUseCase)

		UpdateGroup(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when handler.UpdateGroup return StatusUnprocessableEntity", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, groupResource, bytes.NewReader([]byte("{[}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockGroupUpdateUseCase)

		UpdateGroup(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when handler.UpdateGroup return StatusBadRequest", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Group{
			ID:          "be49afd2ee890805c21ddd55879db1387aec9751",
			Name:        "Grupo Teste",
			Description: "Exemplo de Grupo",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, groupResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockGroupUpdateUseCase)
		mUpdateUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(group.ErrIDValidateModel).
			Once()

		UpdateGroup(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when handler.UpdateGroup return StatusUnauthorized", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Group{
			ID:          "be49afd2ee890805c21ddd55879db1387aec9751",
			Name:        "Grupo Teste",
			Description: "Exemplo de Grupo",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, groupResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockGroupUpdateUseCase)
		mUpdateUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(group.ErrOperationNotAllowed).
			Once()

		UpdateGroup(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when handler.UpdateGroup return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Group{
			ID:          "be49afd2ee890805c21ddd55879db1387aec9751",
			Name:        "Grupo Teste",
			Description: "Exemplo de Grupo",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, groupResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockGroupUpdateUseCase)
		mUpdateUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(group.ErrGroupNotFound).
			Once()

		UpdateGroup(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.UpdateGroup return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Group{
			ID:          "be49afd2ee890805c21ddd55879db1387aec9751",
			Name:        "Grupo Teste",
			Description: "Exemplo de Grupo",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, groupResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockGroupUpdateUseCase)
		mUpdateUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()

		UpdateGroup(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.UpdateGroup return StatusOK", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Group{
			ID:          "be49afd2ee890805c21ddd55879db1387aec9751",
			Name:        "Grupo Teste",
			Description: "Exemplo de Grupo",
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, groupResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mUpdateUseCase := new(mockGroupUpdateUseCase)
		mUpdateUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		UpdateGroup(mJWT, mUpdateUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestHandler_DeleteGroup(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", groupResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mDeleteUseCase := new(mockGroupDeleteUseCase)

		handler := DeleteGroup(mJWT, mDeleteUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", groupResource), handler).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.DeleteGroup return StatusUnauthorized", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", groupResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mDeleteUseCase := new(mockGroupDeleteUseCase)
		mDeleteUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(group.ErrOperationNotAllowed).
			Once()

		handler := DeleteGroup(mJWT, mDeleteUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", groupResource), handler).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when handler.DeleteGroup return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", groupResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mDeleteUseCase := new(mockGroupDeleteUseCase)
		mDeleteUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(group.ErrGroupNotFound).
			Once()

		handler := DeleteGroup(mJWT, mDeleteUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", groupResource), handler).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.DeleteGroup return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", groupResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mDeleteUseCase := new(mockGroupDeleteUseCase)
		mDeleteUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()

		handler := DeleteGroup(mJWT, mDeleteUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", groupResource), handler).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.DeleteGroup return StatusOK", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751", groupResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mDeleteUseCase := new(mockGroupDeleteUseCase)
		mDeleteUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		handler := DeleteGroup(mJWT, mDeleteUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{id}", groupResource), handler).Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestHandler_AddGroupMember(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, memberResource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mAddMemberUseCase := new(mockGroupAddMemberUseCase)

		AddGroupMember(mJWT, mAddMemberUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.AddGroupMember return StatusUnsupportedMediaType", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPost, memberResource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mAddMemberUseCase := new(mockGroupAddMemberUseCase)

		AddGroupMember(mJWT, mAddMemberUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when handler.AddGroupMember return StatusUnprocessableEntity", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPost, memberResource, bytes.NewReader([]byte("{[}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mAddMemberUseCase := new(mockGroupAddMemberUseCase)

		AddGroupMember(mJWT, mAddMemberUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when handler.AddGroupMember return StatusUnauthorized", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Member{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID:  "+5518988888888",
			Admin:   false,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, memberResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mAddMemberUseCase := new(mockGroupAddMemberUseCase)
		mAddMemberUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(group.ErrOperationNotAllowed).
			Once()

		AddGroupMember(mJWT, mAddMemberUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when handler.AddGroupMember return StatusBadRequest", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Member{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID:  "",
			Admin:   false,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, memberResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mAddMemberUseCase := new(mockGroupAddMemberUseCase)
		mAddMemberUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(group.ErrUserIDValidateModel).
			Once()

		AddGroupMember(mJWT, mAddMemberUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when handler.AddGroupMember return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Member{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID:  "+5518988888888",
			Admin:   false,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, memberResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mAddMemberUseCase := new(mockGroupAddMemberUseCase)
		mAddMemberUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(group.ErrGroupNotFound).
			Once()

		AddGroupMember(mJWT, mAddMemberUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.AddGroupMember return StatusConflict", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Member{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID:  "+5518988888888",
			Admin:   false,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, memberResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mAddMemberUseCase := new(mockGroupAddMemberUseCase)
		mAddMemberUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(group.ErrMemberAlreadyExists).
			Once()

		AddGroupMember(mJWT, mAddMemberUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("when handler.AddGroupMember return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Member{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID:  "+5518988888888",
			Admin:   false,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, memberResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mAddMemberUseCase := new(mockGroupAddMemberUseCase)
		mAddMemberUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()

		AddGroupMember(mJWT, mAddMemberUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.AddGroupMember return StatusCreated", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Member{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID:  "+5518988888888",
			Admin:   false,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, memberResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mAddMemberUseCase := new(mockGroupAddMemberUseCase)
		mAddMemberUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		AddGroupMember(mJWT, mAddMemberUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestHandler_RemoveGroupMember(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751/+5518988888888",
				memberResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mRemoveMemberUseCase := new(mockGroupRemoveMemberUseCase)

		handler := RemoveGroupMember(mJWT, mRemoveMemberUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{group}/{user}", memberResource), handler).
			Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.RemoveGroupMember return StatusUnauthorized", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751/+5518988888888",
				memberResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mRemoveMemberUseCase := new(mockGroupRemoveMemberUseCase)
		mRemoveMemberUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(group.ErrGroupOwnerCannotRemoved).
			Once()

		handler := RemoveGroupMember(mJWT, mRemoveMemberUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{group}/{user}", memberResource), handler).
			Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when handler.RemoveGroupMember return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751/+5518988888888",
				memberResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mRemoveMemberUseCase := new(mockGroupRemoveMemberUseCase)
		mRemoveMemberUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(group.ErrMemberNotFound).
			Once()

		handler := RemoveGroupMember(mJWT, mRemoveMemberUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{group}/{user}", memberResource), handler).
			Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.RemoveGroupMember return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751/+5518988888888",
				memberResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mRemoveMemberUseCase := new(mockGroupRemoveMemberUseCase)
		mRemoveMemberUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()

		handler := RemoveGroupMember(mJWT, mRemoveMemberUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{group}/{user}", memberResource), handler).
			Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.RemoveGroupMember return StatusOK", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodDelete,
			fmt.Sprintf("%s/be49afd2ee890805c21ddd55879db1387aec9751/+5518988888888",
				memberResource), nil)
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mRemoveMemberUseCase := new(mockGroupRemoveMemberUseCase)
		mRemoveMemberUseCase.On("Execute", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		handler := RemoveGroupMember(mJWT, mRemoveMemberUseCase)

		router := mux.NewRouter()
		router.Handle(fmt.Sprintf("%s/{group}/{user}", memberResource), handler).
			Methods(http.MethodDelete)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestHandler_SetGroupAdmin(t *testing.T) {
	//t.Parallel()

	t.Run("when JWT fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, memberResource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mSetAdminUseCase := new(mockGroupSetAdminUseCase)

		SetGroupAdmin(mJWT, mSetAdminUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.SetGroupAdmin return StatusUnsupportedMediaType", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, memberResource, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mSetAdminUseCase := new(mockGroupSetAdminUseCase)

		SetGroupAdmin(mJWT, mSetAdminUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("when handler.SetGroupAdmin return StatusUnprocessableEntity", func(t *testing.T) {
		//t.Parallel()
		req := httptest.NewRequest(http.MethodPut, memberResource, bytes.NewReader([]byte("{[}")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mSetAdminUseCase := new(mockGroupSetAdminUseCase)

		SetGroupAdmin(mJWT, mSetAdminUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when handler.SetGroupAdmin return StatusBadRequest", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Member{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID:  "",
			Admin:   true,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, memberResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mSetAdminUseCase := new(mockGroupSetAdminUseCase)
		mSetAdminUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(group.ErrUserIDValidateModel).
			Once()

		SetGroupAdmin(mJWT, mSetAdminUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when handler.SetGroupAdmin return StatusUnauthorized", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Member{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID:  "+5518977777777",
			Admin:   true,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, memberResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mSetAdminUseCase := new(mockGroupSetAdminUseCase)
		mSetAdminUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(group.ErrOperationNotAllowed).
			Once()

		SetGroupAdmin(mJWT, mSetAdminUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when handler.SetGroupAdmin return StatusNotFound", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Member{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID:  "+5518977777777",
			Admin:   true,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, memberResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mSetAdminUseCase := new(mockGroupSetAdminUseCase)
		mSetAdminUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(group.ErrMemberNotFound).
			Once()

		SetGroupAdmin(mJWT, mSetAdminUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("when handler.SetGroupAdmin return StatusInternalServerError", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Member{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID:  "+5518977777777",
			Admin:   true,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, memberResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mSetAdminUseCase := new(mockGroupSetAdminUseCase)
		mSetAdminUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()

		SetGroupAdmin(mJWT, mSetAdminUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when handler.SetGroupAdmin return StatusOK", func(t *testing.T) {
		//t.Parallel()
		p := &dto.Member{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID:  "+5518977777777",
			Admin:   true,
		}

		pj, err := json.Marshal(p)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, memberResource, bytes.NewReader(pj))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mJWT := new(common.MockJWT)
		mJWT.On("GetDataToken", mock.Anything, mock.Anything).
			Return("+5518999999999", nil).
			Once()
		mSetAdminUseCase := new(mockGroupSetAdminUseCase)
		mSetAdminUseCase.On("Execute", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		SetGroupAdmin(mJWT, mSetAdminUseCase).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
