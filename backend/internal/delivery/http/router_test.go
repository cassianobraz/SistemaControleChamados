package httpapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httpapi "github.com/cassianobraz/SistemaControleChamados/backend/internal/delivery/http"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewRouter(t *testing.T) {
	t.Parallel()

	ticketUC := new(mocks.TicketUseCase)
	responsibleUC := new(mocks.ResponsibleUseCase)

	router := httpapi.NewRouter(httpapi.Dependencies{
		TicketUseCase:      ticketUC,
		ResponsibleUseCase: responsibleUC,
		OpenAPISpec:        []byte(`{"openapi":"3.0.0"}`),
	})

	t.Run("routes the health check", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("routes ticket listing to the ticket use case", func(t *testing.T) {
		t.Parallel()

		ticketUC.On("List", mock.Anything, mock.Anything).
			Return(nil, domain.ErrTicketNotFound).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/tickets", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		ticketUC.AssertExpectations(t)
	})

	t.Run("routes responsible listing to the responsible use case", func(t *testing.T) {
		t.Parallel()

		responsibleUC.On("List", mock.Anything).Return(nil, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/responsibles", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		responsibleUC.AssertExpectations(t)
	})

	t.Run("serves the raw OpenAPI spec", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"openapi":"3.0.0"}`, rec.Body.String())
	})

	t.Run("serves the swagger UI", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("applies open CORS headers to every response", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodOptions, "/api/v1/tickets", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
	})
}
