package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/delivery/http/middleware"
	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	t.Parallel()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	middleware.Logging(next).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTeapot, rec.Code)
}
