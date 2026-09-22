package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
)

func TestNewTicketModel(t *testing.T) {
	t.Parallel()

	ticket, err := domain.NewTicket("título", "descrição", domain.PriorityHigh, "resp-1", "Ana")
	require.NoError(t, err)

	model := newTicketModel(ticket)

	assert.Equal(t, ticket.ID, model.ID)
	assert.Equal(t, "título", model.Title)
	assert.Equal(t, "alta", model.Priority)
	assert.Equal(t, 3, model.PriorityRank)
}

func TestTicketModel_ToDomain(t *testing.T) {
	t.Parallel()

	model := ticketModel{
		ID:              "018f8f3a-1234-7abc-8def-0123456789ab",
		Title:           "título",
		Description:     "descrição",
		Priority:        "media",
		PriorityRank:    2,
		Status:          "em_andamento",
		ResponsibleID:   "resp-1",
		ResponsibleName: "Ana",
	}

	ticket := model.toDomain()

	assert.Equal(t, model.ID, ticket.ID)
	assert.Equal(t, domain.PriorityMedium, ticket.Priority)
	assert.Equal(t, domain.StatusInProgress, ticket.Status)
	assert.Equal(t, "Ana", ticket.ResponsibleName)
}

func TestPriorityRank(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		priority domain.Priority
		want     int
	}{
		{name: "alta ranks highest", priority: domain.PriorityHigh, want: 3},
		{name: "media ranks middle", priority: domain.PriorityMedium, want: 2},
		{name: "baixa ranks lowest", priority: domain.PriorityLow, want: 1},
		{name: "unknown ranks zero", priority: domain.Priority("desconhecida"), want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, priorityRank(tt.priority))
		})
	}
}
