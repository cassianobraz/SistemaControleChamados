package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/delivery/http/handler"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.Health(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
}
