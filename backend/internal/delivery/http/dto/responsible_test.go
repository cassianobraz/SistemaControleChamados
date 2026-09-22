package dto_test

import (
	"testing"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/delivery/http/dto"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResponsibleResponseList(t *testing.T) {
	t.Parallel()

	t.Run("maps every responsible with their workload", func(t *testing.T) {
		t.Parallel()

		responsible, err := domain.NewResponsible("Ana", "ana@codificar.dev")
		require.NoError(t, err)
		responsible.ID = "resp-1"

		got := dto.NewResponsibleResponseList([]*usecase.ResponsibleWithWorkload{
			{Responsible: responsible, OpenTickets: 5},
		})

		require.Len(t, got, 1)
		assert.Equal(t, "resp-1", got[0].ID)
		assert.Equal(t, "Ana", got[0].Name)
		assert.Equal(t, int64(5), got[0].OpenTickets)
		assert.True(t, got[0].Active)
	})

	t.Run("returns an empty, non-nil slice for no responsibles", func(t *testing.T) {
		t.Parallel()

		got := dto.NewResponsibleResponseList(nil)

		assert.NotNil(t, got)
		assert.Empty(t, got)
	})
}
