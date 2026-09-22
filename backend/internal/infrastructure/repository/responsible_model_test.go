package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
)

func TestNewResponsibleModel(t *testing.T) {
	t.Parallel()

	responsible, err := domain.NewResponsible("Ana", "ana@codificar.dev")
	require.NoError(t, err)

	model := newResponsibleModel(responsible)

	assert.Equal(t, responsible.ID, model.ID)
	assert.Equal(t, "Ana", model.Name)
	assert.True(t, model.Active)
}

func TestResponsibleModel_ToDomain(t *testing.T) {
	t.Parallel()

	model := responsibleModel{
		ID:     "018f8f3a-1234-7abc-8def-0123456789ab",
		Name:   "Ana",
		Email:  "ana@codificar.dev",
		Active: true,
	}

	responsible := model.toDomain()

	assert.Equal(t, model.ID, responsible.ID)
	assert.Equal(t, "Ana", responsible.Name)
	assert.True(t, responsible.Active)
}
