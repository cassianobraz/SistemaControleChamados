package domain_test

import (
	"testing"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResponsible(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		email   string
		wantErr error
	}{
		{name: "valid name", input: "Ana Souza", email: "ana@codificar.dev"},
		{name: "blank name is invalid", input: "   ", email: "ana@codificar.dev", wantErr: domain.ErrInvalidResponsibleName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			responsible, err := domain.NewResponsible(tt.input, tt.email)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, responsible)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, responsible)
			assert.Equal(t, tt.input, responsible.Name)
			assert.Equal(t, tt.email, responsible.Email)
			assert.True(t, responsible.Active)
		})
	}
}

func TestResponsible_UpdateDetails(t *testing.T) {
	t.Parallel()

	t.Run("applies new values", func(t *testing.T) {
		t.Parallel()

		responsible, err := domain.NewResponsible("Ana Souza", "ana@codificar.dev")
		require.NoError(t, err)

		err = responsible.UpdateDetails("Ana Paula", "ana.paula@codificar.dev", false)

		require.NoError(t, err)
		assert.Equal(t, "Ana Paula", responsible.Name)
		assert.Equal(t, "ana.paula@codificar.dev", responsible.Email)
		assert.False(t, responsible.Active)
	})

	t.Run("rejects a blank name", func(t *testing.T) {
		t.Parallel()

		responsible, err := domain.NewResponsible("Ana Souza", "ana@codificar.dev")
		require.NoError(t, err)

		err = responsible.UpdateDetails("   ", "ana@codificar.dev", true)

		require.ErrorIs(t, err, domain.ErrInvalidResponsibleName)
		assert.Equal(t, "Ana Souza", responsible.Name)
	})
}
