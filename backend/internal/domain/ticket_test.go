package domain_test

import (
	"strings"
	"testing"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTicket(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		title         string
		description   string
		priority      domain.Priority
		responsibleID string
		wantErr       error
	}{
		{
			name:          "valid ticket",
			title:         "Impressora não funciona",
			description:   "A impressora do 2º andar não liga.",
			priority:      domain.PriorityMedium,
			responsibleID: "resp-1",
		},
		{
			name:        "blank title is invalid",
			title:       "   ",
			description: "descrição válida",
			priority:    domain.PriorityLow,
			wantErr:     domain.ErrInvalidTitle,
		},
		{
			name:        "title over 200 chars is invalid",
			title:       strings.Repeat("a", 201),
			description: "descrição válida",
			priority:    domain.PriorityLow,
			wantErr:     domain.ErrTitleTooLong,
		},
		{
			name:        "blank description is invalid",
			title:       "título válido",
			description: "   ",
			priority:    domain.PriorityLow,
			wantErr:     domain.ErrInvalidDescription,
		},
		{
			name:        "invalid priority",
			title:       "título válido",
			description: "descrição válida",
			priority:    domain.Priority("urgente"),
			wantErr:     domain.ErrInvalidPriority,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ticket, err := domain.NewTicket(tt.title, tt.description, tt.priority, tt.responsibleID, "Fulano")

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, ticket)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, ticket)
			assert.Equal(t, strings.TrimSpace(tt.title), ticket.Title)
			assert.Equal(t, tt.description, ticket.Description)
			assert.Equal(t, tt.priority, ticket.Priority)
			assert.Equal(t, domain.StatusOpen, ticket.Status)
			assert.Equal(t, tt.responsibleID, ticket.ResponsibleID)
			assert.False(t, ticket.CreatedAt.IsZero())
			assert.Equal(t, ticket.CreatedAt, ticket.UpdatedAt)
		})
	}
}

func TestTicket_UpdateDetails(t *testing.T) {
	t.Parallel()

	newValidTicket := func(t *testing.T) *domain.Ticket {
		t.Helper()
		ticket, err := domain.NewTicket("título", "descrição", domain.PriorityLow, "resp-1", "Fulano")
		require.NoError(t, err)
		return ticket
	}

	t.Run("applies a valid edition", func(t *testing.T) {
		t.Parallel()

		ticket := newValidTicket(t)
		previousUpdatedAt := ticket.UpdatedAt

		err := ticket.UpdateDetails("novo título", "nova descrição", domain.PriorityHigh, domain.StatusInProgress)

		require.NoError(t, err)
		assert.Equal(t, "novo título", ticket.Title)
		assert.Equal(t, "nova descrição", ticket.Description)
		assert.Equal(t, domain.PriorityHigh, ticket.Priority)
		assert.Equal(t, domain.StatusInProgress, ticket.Status)
		assert.False(t, ticket.UpdatedAt.Before(previousUpdatedAt))
	})

	tests := []struct {
		name        string
		title       string
		description string
		priority    domain.Priority
		status      domain.Status
		wantErr     error
	}{
		{
			name:        "invalid title",
			title:       "",
			description: "descrição",
			priority:    domain.PriorityLow,
			status:      domain.StatusOpen,
			wantErr:     domain.ErrInvalidTitle,
		},
		{
			name:        "invalid description",
			title:       "título",
			description: "",
			priority:    domain.PriorityLow,
			status:      domain.StatusOpen,
			wantErr:     domain.ErrInvalidDescription,
		},
		{
			name:        "invalid priority",
			title:       "título",
			description: "descrição",
			priority:    domain.Priority("crítica"),
			status:      domain.StatusOpen,
			wantErr:     domain.ErrInvalidPriority,
		},
		{
			name:        "invalid status",
			title:       "título",
			description: "descrição",
			priority:    domain.PriorityLow,
			status:      domain.Status("pausado"),
			wantErr:     domain.ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ticket := newValidTicket(t)
			err := ticket.UpdateDetails(tt.title, tt.description, tt.priority, tt.status)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestTicket_AssignResponsible(t *testing.T) {
	t.Parallel()

	ticket, err := domain.NewTicket("título", "descrição", domain.PriorityLow, "resp-1", "Fulano")
	require.NoError(t, err)
	previousUpdatedAt := ticket.UpdatedAt

	ticket.AssignResponsible("resp-2", "Ciclana")

	assert.Equal(t, "resp-2", ticket.ResponsibleID)
	assert.Equal(t, "Ciclana", ticket.ResponsibleName)
	assert.False(t, ticket.UpdatedAt.Before(previousUpdatedAt))
}

func TestTicket_IsOpen(t *testing.T) {
	t.Parallel()

	ticket, err := domain.NewTicket("título", "descrição", domain.PriorityLow, "resp-1", "Fulano")
	require.NoError(t, err)

	assert.True(t, ticket.IsOpen())

	require.NoError(t, ticket.UpdateDetails(ticket.Title, ticket.Description, ticket.Priority, domain.StatusClosed))
	assert.False(t, ticket.IsOpen())
}
