package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/delivery/http/dto"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/delivery/http/handler"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/mocks"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResponsibleHandler_List(t *testing.T) {
	t.Parallel()

	t.Run("returns 200 with the responsibles and their workload", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.ResponsibleUseCase)
		h := handler.NewResponsibleHandler(uc)

		responsible, err := domain.NewResponsible("Ana", "ana@codificar.dev")
		require.NoError(t, err)
		responsible.ID = "resp-1"

		uc.On("List", mock.Anything).Return([]*usecase.ResponsibleWithWorkload{
			{Responsible: responsible, OpenTickets: 3},
		}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/responsibles", nil)
		rec := httptest.NewRecorder()

		h.List(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var got []dto.ResponsibleResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
		require.Len(t, got, 1)
		assert.Equal(t, "resp-1", got[0].ID)
		assert.Equal(t, int64(3), got[0].OpenTickets)
	})

	t.Run("returns 500 when the use case fails", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.ResponsibleUseCase)
		h := handler.NewResponsibleHandler(uc)

		uc.On("List", mock.Anything).Return(nil, errors.New("boom")).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/responsibles", nil)
		rec := httptest.NewRecorder()

		h.List(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestResponsibleHandler_Create(t *testing.T) {
	t.Parallel()

	t.Run("returns 201 with the created responsible", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.ResponsibleUseCase)
		h := handler.NewResponsibleHandler(uc)

		responsible, err := domain.NewResponsible("Ana", "ana@codificar.dev")
		require.NoError(t, err)
		responsible.ID = "resp-1"

		uc.On("Create", mock.Anything, usecase.CreateResponsibleInput{Name: "Ana", Email: "ana@codificar.dev"}).
			Return(&usecase.ResponsibleWithWorkload{Responsible: responsible, OpenTickets: 0}, nil).Once()

		body := bytes.NewBufferString(`{"name":"Ana","email":"ana@codificar.dev"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/responsibles", body)
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code)

		var got dto.ResponsibleResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
		assert.Equal(t, "resp-1", got.ID)
	})

	t.Run("returns 400 for an invalid name", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.ResponsibleUseCase)
		h := handler.NewResponsibleHandler(uc)

		uc.On("Create", mock.Anything, mock.Anything).Return(nil, domain.ErrInvalidResponsibleName).Once()

		body := bytes.NewBufferString(`{"name":""}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/responsibles", body)
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestResponsibleHandler_Update(t *testing.T) {
	t.Parallel()

	t.Run("returns 200 with the updated responsible", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.ResponsibleUseCase)
		h := handler.NewResponsibleHandler(uc)

		responsible, err := domain.NewResponsible("Ana Paula", "ana.paula@codificar.dev")
		require.NoError(t, err)
		responsible.ID = "resp-1"
		responsible.Active = false

		uc.On("Update", mock.Anything, "resp-1", usecase.UpdateResponsibleInput{
			Name:   "Ana Paula",
			Email:  "ana.paula@codificar.dev",
			Active: false,
		}).Return(&usecase.ResponsibleWithWorkload{Responsible: responsible, OpenTickets: 1}, nil).Once()

		body := bytes.NewBufferString(`{"name":"Ana Paula","email":"ana.paula@codificar.dev","active":false}`)
		req := withURLParam(httptest.NewRequest(http.MethodPut, "/api/v1/responsibles/resp-1", body), "id", "resp-1")
		rec := httptest.NewRecorder()

		h.Update(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var got dto.ResponsibleResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
		assert.Equal(t, "Ana Paula", got.Name)
		assert.False(t, got.Active)
	})

	t.Run("returns 404 when the responsible does not exist", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.ResponsibleUseCase)
		h := handler.NewResponsibleHandler(uc)

		uc.On("Update", mock.Anything, "missing", mock.Anything).Return(nil, domain.ErrResponsibleNotFound).Once()

		body := bytes.NewBufferString(`{"name":"Ana"}`)
		req := withURLParam(httptest.NewRequest(http.MethodPut, "/api/v1/responsibles/missing", body), "id", "missing")
		rec := httptest.NewRecorder()

		h.Update(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestResponsibleHandler_Delete(t *testing.T) {
	t.Parallel()

	t.Run("returns 204 when the responsible is removed", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.ResponsibleUseCase)
		h := handler.NewResponsibleHandler(uc)

		uc.On("Delete", mock.Anything, "resp-1").Return(nil).Once()

		req := withURLParam(httptest.NewRequest(http.MethodDelete, "/api/v1/responsibles/resp-1", nil), "id", "resp-1")
		rec := httptest.NewRecorder()

		h.Delete(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("returns 404 when the responsible does not exist", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.ResponsibleUseCase)
		h := handler.NewResponsibleHandler(uc)

		uc.On("Delete", mock.Anything, "missing").Return(domain.ErrResponsibleNotFound).Once()

		req := withURLParam(httptest.NewRequest(http.MethodDelete, "/api/v1/responsibles/missing", nil), "id", "missing")
		rec := httptest.NewRecorder()

		h.Delete(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
