package database_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/infrastructure/database"
)

func TestConnect_InvalidURI(t *testing.T) {
	t.Parallel()

	_, err := database.Connect(context.Background(), "not-a-valid-uri")

	require.Error(t, err)
}
