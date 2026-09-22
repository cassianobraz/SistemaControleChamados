package dto_test

import (
	"testing"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/delivery/http/dto"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTicketResponse(t *testing.T) {
	t.Parallel()

	ticket, err := domain.NewTicket("título", "descrição", domain.PriorityHigh, "resp-1", "Ana")
	require.NoError(t, err)
	ticket.ID = "ticket-1"

	got := dto.NewTicketResponse(ticket)

	assert.Equal(t, "ticket-1", got.ID)
	assert.Equal(t, "título", got.Title)
	assert.Equal(t, "alta", got.Priority)
	assert.Equal(t, "aberto", got.Status)
	assert.Equal(t, "resp-1", got.ResponsibleID)
	assert.Equal(t, "Ana", got.ResponsibleName)
}

func TestNewTicketResponseList(t *testing.T) {
	t.Parallel()

	t.Run("maps every ticket", func(t *testing.T) {
		t.Parallel()

		ticket, err := domain.NewTicket("título", "descrição", domain.PriorityLow, "resp-1", "Ana")
		require.NoError(t, err)

		got := dto.NewTicketResponseList([]*domain.Ticket{ticket})

		require.Len(t, got, 1)
		assert.Equal(t, "título", got[0].Title)
	})

	t.Run("returns an empty, non-nil slice for no tickets", func(t *testing.T) {
		t.Parallel()

		got := dto.NewTicketResponseList(nil)

		assert.NotNil(t, got)
		assert.Empty(t, got)
	})
}
