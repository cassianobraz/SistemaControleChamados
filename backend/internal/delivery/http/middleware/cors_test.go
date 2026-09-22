package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/delivery/http/middleware"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	t.Parallel()

	t.Run("sets permissive headers and forwards the request", func(t *testing.T) {
		t.Parallel()

		called := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/tickets", nil)
		rec := httptest.NewRecorder()

		middleware.CORS(next).ServeHTTP(rec, req)

		assert.True(t, called)
		assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Headers"))
		assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), "DELETE")
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("short-circuits preflight requests", func(t *testing.T) {
		t.Parallel()

		called := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
		})

		req := httptest.NewRequest(http.MethodOptions, "/api/v1/tickets", nil)
		rec := httptest.NewRecorder()

		middleware.CORS(next).ServeHTTP(rec, req)

		assert.False(t, called)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}
