package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/delivery/http/dto"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/delivery/http/handler"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/mocks"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func sampleTicket() *domain.Ticket {
	ticket, _ := domain.NewTicket("título", "descrição", domain.PriorityMedium, "resp-1", "Ana")
	ticket.ID = "ticket-1"

	return ticket
}

func withURLParam(r *http.Request, key, value string) *http.Request {
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add(key, value)

	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, routeCtx))
}

func TestTicketHandler_Create(t *testing.T) {
	t.Parallel()

	t.Run("returns 201 with the created ticket", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.TicketUseCase)
		h := handler.NewTicketHandler(uc)

		expected := sampleTicket()
		uc.On("Create", mock.Anything, mock.AnythingOfType("usecase.CreateTicketInput")).
			Return(expected, nil).Once()

		body := `{"title":"título","description":"descrição","priority":"media"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code)

		var got dto.TicketResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
		assert.Equal(t, expected.ID, got.ID)
	})

	t.Run("returns 400 for a malformed body", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.TicketUseCase)
		h := handler.NewTicketHandler(uc)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewBufferString(`{"title":`))
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		uc.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("returns 400 when the use case reports a validation error", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.TicketUseCase)
		h := handler.NewTicketHandler(uc)

		uc.On("Create", mock.Anything, mock.Anything).Return(nil, domain.ErrInvalidPriority).Once()

		body := `{"title":"título","description":"descrição","priority":"urgente"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 500 for unexpected use case errors", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.TicketUseCase)
		h := handler.NewTicketHandler(uc)

		uc.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("boom")).Once()

		body := `{"title":"título","description":"descrição","priority":"media"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestTicketHandler_Get(t *testing.T) {
	t.Parallel()

	t.Run("returns 200 with the ticket", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.TicketUseCase)
		h := handler.NewTicketHandler(uc)

		expected := sampleTicket()
		uc.On("GetByID", mock.Anything, "ticket-1").Return(expected, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/tickets/ticket-1", nil)
		req = withURLParam(req, "id", "ticket-1")
		rec := httptest.NewRecorder()

		h.Get(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("returns 404 when the ticket does not exist", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.TicketUseCase)
		h := handler.NewTicketHandler(uc)

		uc.On("GetByID", mock.Anything, "missing").Return(nil, domain.ErrTicketNotFound).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/tickets/missing", nil)
		req = withURLParam(req, "id", "missing")
		rec := httptest.NewRecorder()

		h.Get(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestTicketHandler_List(t *testing.T) {
	t.Parallel()

	t.Run("parses filters and pagination from the query string", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.TicketUseCase)
		h := handler.NewTicketHandler(uc)

		uc.On("List", mock.Anything, mock.MatchedBy(func(input usecase.ListTicketsInput) bool {
			return input.Status != nil && *input.Status == domain.StatusOpen &&
				input.Priority != nil && *input.Priority == domain.PriorityHigh &&
				input.ResponsibleID != nil && *input.ResponsibleID == "resp-1" &&
				input.Search != nil && *input.Search == "impressora" &&
				input.Sort == domain.SortByCreatedAtAsc &&
				input.Page == 2 && input.PageSize == 10
		})).Return(&usecase.ListTicketsOutput{
			Tickets:  []*domain.Ticket{sampleTicket()},
			Total:    11,
			Page:     2,
			PageSize: 10,
		}, nil).Once()

		url := "/api/v1/tickets?status=aberto&priority=alta&responsible_id=resp-1&search=impressora&sort=created_at_asc&page=2&page_size=10"
		req := httptest.NewRequest(http.MethodGet, url, nil)
		rec := httptest.NewRecorder()

		h.List(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var got dto.TicketListResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
		assert.Equal(t, int64(11), got.Total)
		assert.Equal(t, int64(2), got.TotalPages)
		assert.Len(t, got.Data, 1)
	})

	t.Run("ignores non-numeric page values", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.TicketUseCase)
		h := handler.NewTicketHandler(uc)

		uc.On("List", mock.Anything, mock.MatchedBy(func(input usecase.ListTicketsInput) bool {
			return input.Page == 0
		})).Return(&usecase.ListTicketsOutput{Tickets: []*domain.Ticket{}}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/tickets?page=abc", nil)
		rec := httptest.NewRecorder()

		h.List(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("returns an empty data array instead of null when there are no tickets", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.TicketUseCase)
		h := handler.NewTicketHandler(uc)

		uc.On("List", mock.Anything, mock.Anything).
			Return(&usecase.ListTicketsOutput{Tickets: []*domain.Ticket{}, Total: 0, Page: 1, PageSize: 20}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/tickets", nil)
		rec := httptest.NewRecorder()

		h.List(rec, req)

		assert.JSONEq(t, `{"data":[],"total":0,"page":1,"page_size":20,"total_pages":0}`, rec.Body.String())
	})

	t.Run("returns 500 when the use case fails", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.TicketUseCase)
		h := handler.NewTicketHandler(uc)

		uc.On("List", mock.Anything, mock.Anything).Return(nil, errors.New("boom")).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/tickets", nil)
		rec := httptest.NewRecorder()

		h.List(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestTicketHandler_Update(t *testing.T) {
	t.Parallel()

	t.Run("returns 200 with the updated ticket", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.TicketUseCase)
		h := handler.NewTicketHandler(uc)

		expected := sampleTicket()
		uc.On("Update", mock.Anything, "ticket-1", mock.AnythingOfType("usecase.UpdateTicketInput")).
			Return(expected, nil).Once()

		body := `{"title":"novo","description":"nova","priority":"alta","status":"em_andamento"}`
		req := httptest.NewRequest(http.MethodPut, "/api/v1/tickets/ticket-1", bytes.NewBufferString(body))
		req = withURLParam(req, "id", "ticket-1")
		rec := httptest.NewRecorder()

		h.Update(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("returns 400 for a malformed body", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.TicketUseCase)
		h := handler.NewTicketHandler(uc)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/tickets/ticket-1", bytes.NewBufferString(`{`))
		req = withURLParam(req, "id", "ticket-1")
		rec := httptest.NewRecorder()

		h.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 404 when the ticket does not exist", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.TicketUseCase)
		h := handler.NewTicketHandler(uc)

		uc.On("Update", mock.Anything, "missing", mock.Anything).
			Return(nil, domain.ErrTicketNotFound).Once()

		body := `{"title":"novo","description":"nova","priority":"alta","status":"em_andamento"}`
		req := httptest.NewRequest(http.MethodPut, "/api/v1/tickets/missing", bytes.NewBufferString(body))
		req = withURLParam(req, "id", "missing")
		rec := httptest.NewRecorder()

		h.Update(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestTicketHandler_Assign(t *testing.T) {
	t.Parallel()

	t.Run("assigns manually when a responsible id is given", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.TicketUseCase)
		h := handler.NewTicketHandler(uc)

		expected := sampleTicket()
		uc.On("AssignResponsible", mock.Anything, "ticket-1", "resp-2").Return(expected, nil).Once()

		body := `{"responsible_id":"resp-2"}`
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/tickets/ticket-1/assign", bytes.NewBufferString(body))
		req = withURLParam(req, "id", "ticket-1")
		req.ContentLength = int64(len(body))
		rec := httptest.NewRecorder()

		h.Assign(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("assigns automatically when the body is empty", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.TicketUseCase)
		h := handler.NewTicketHandler(uc)

		expected := sampleTicket()
		uc.On("AssignResponsible", mock.Anything, "ticket-1", "").Return(expected, nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/tickets/ticket-1/assign", nil)
		req = withURLParam(req, "id", "ticket-1")
		req.ContentLength = 0
		rec := httptest.NewRecorder()

		h.Assign(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("returns 400 for a malformed body", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.TicketUseCase)
		h := handler.NewTicketHandler(uc)

		body := `{`
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/tickets/ticket-1/assign", bytes.NewBufferString(body))
		req = withURLParam(req, "id", "ticket-1")
		req.ContentLength = int64(len(body))
		rec := httptest.NewRecorder()

		h.Assign(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 400 when there is no responsible available", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.TicketUseCase)
		h := handler.NewTicketHandler(uc)

		uc.On("AssignResponsible", mock.Anything, "ticket-1", "").
			Return(nil, domain.ErrNoResponsibleAvailable).Once()

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/tickets/ticket-1/assign", nil)
		req = withURLParam(req, "id", "ticket-1")
		rec := httptest.NewRecorder()

		h.Assign(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 404 when the ticket does not exist", func(t *testing.T) {
		t.Parallel()

		uc := new(mocks.TicketUseCase)
		h := handler.NewTicketHandler(uc)

		uc.On("AssignResponsible", mock.Anything, "missing", "").
			Return(nil, domain.ErrTicketNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/tickets/missing/assign", nil)
		req = withURLParam(req, "id", "missing")
		rec := httptest.NewRecorder()

		h.Assign(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
