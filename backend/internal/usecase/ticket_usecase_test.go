package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/mocks"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newResponsible(t *testing.T, id, name string) *domain.Responsible {
	t.Helper()

	responsible, err := domain.NewResponsible(name, name+"@codificar.dev")
	require.NoError(t, err)
	responsible.ID = id

	return responsible
}

func TestTicketUseCase_Create(t *testing.T) {
	t.Parallel()

	validInput := usecase.CreateTicketInput{
		Title:       "Impressora não funciona",
		Description: "A impressora do 2º andar não liga.",
		Priority:    domain.PriorityMedium,
	}

	t.Run("creates a ticket with the manually chosen responsible", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		responsible := newResponsible(t, "resp-1", "Ana")
		input := validInput
		input.ResponsibleID = "resp-1"

		responsibleRepo.On("FindByID", mock.Anything, "resp-1").Return(responsible, nil).Once()
		ticketRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Ticket")).Return(nil).Once()

		ticket, err := uc.Create(context.Background(), input)

		require.NoError(t, err)
		require.NotNil(t, ticket)
		assert.Equal(t, "resp-1", ticket.ResponsibleID)
		assert.Equal(t, "Ana", ticket.ResponsibleName)
		ticketRepo.AssertExpectations(t)
		responsibleRepo.AssertExpectations(t)
	})

	t.Run("auto assigns to the responsible with the fewest open tickets", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		ana := newResponsible(t, "resp-1", "Ana")
		bruno := newResponsible(t, "resp-2", "Bruno")

		responsibleRepo.On("FindAll", mock.Anything).
			Return([]*domain.Responsible{ana, bruno}, nil).Once()
		ticketRepo.On("CountOpenByResponsible", mock.Anything, "resp-1").Return(int64(3), nil).Once()
		ticketRepo.On("CountOpenByResponsible", mock.Anything, "resp-2").Return(int64(1), nil).Once()
		ticketRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Ticket")).Return(nil).Once()

		ticket, err := uc.Create(context.Background(), validInput)

		require.NoError(t, err)
		assert.Equal(t, "resp-2", ticket.ResponsibleID)
		assert.Equal(t, "Bruno", ticket.ResponsibleName)
	})

	t.Run("skips inactive responsibles during auto assignment", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		ana := newResponsible(t, "resp-1", "Ana")
		ana.Active = false
		bruno := newResponsible(t, "resp-2", "Bruno")

		responsibleRepo.On("FindAll", mock.Anything).
			Return([]*domain.Responsible{ana, bruno}, nil).Once()
		ticketRepo.On("CountOpenByResponsible", mock.Anything, "resp-2").Return(int64(5), nil).Once()
		ticketRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Ticket")).Return(nil).Once()

		ticket, err := uc.Create(context.Background(), validInput)

		require.NoError(t, err)
		assert.Equal(t, "resp-2", ticket.ResponsibleID)
		ticketRepo.AssertNotCalled(t, "CountOpenByResponsible", mock.Anything, "resp-1")
	})

	t.Run("breaks ties between equally loaded responsibles instead of always picking the same one", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		ana := newResponsible(t, "resp-1", "Ana")
		bruno := newResponsible(t, "resp-2", "Bruno")
		carla := newResponsible(t, "resp-3", "Carla")

		responsibleRepo.On("FindAll", mock.Anything).Return([]*domain.Responsible{ana, bruno, carla}, nil)
		ticketRepo.On("CountOpenByResponsible", mock.Anything, "resp-1").Return(int64(0), nil)
		ticketRepo.On("CountOpenByResponsible", mock.Anything, "resp-2").Return(int64(0), nil)
		ticketRepo.On("CountOpenByResponsible", mock.Anything, "resp-3").Return(int64(0), nil)
		ticketRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Ticket")).Return(nil)

		picked := map[string]bool{}
		for range 200 {
			ticket, err := uc.Create(context.Background(), validInput)
			require.NoError(t, err)
			picked[ticket.ResponsibleID] = true
		}

		assert.Len(t, picked, 3, "expected every responsible tied for least loaded to be picked at least once")
	})

	t.Run("fails when there are no responsibles for auto assignment", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		responsibleRepo.On("FindAll", mock.Anything).Return([]*domain.Responsible{}, nil).Once()

		ticket, err := uc.Create(context.Background(), validInput)

		require.ErrorIs(t, err, domain.ErrNoResponsibleAvailable)
		assert.Nil(t, ticket)
	})

	t.Run("fails when the chosen responsible does not exist", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		input := validInput
		input.ResponsibleID = "resp-404"

		responsibleRepo.On("FindByID", mock.Anything, "resp-404").
			Return(nil, domain.ErrResponsibleNotFound).Once()

		ticket, err := uc.Create(context.Background(), input)

		require.ErrorIs(t, err, domain.ErrResponsibleNotFound)
		assert.Nil(t, ticket)
		ticketRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("fails when ticket fields are invalid", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		responsible := newResponsible(t, "resp-1", "Ana")
		responsibleRepo.On("FindByID", mock.Anything, "resp-1").Return(responsible, nil).Once()

		input := usecase.CreateTicketInput{Title: "", ResponsibleID: "resp-1"}

		ticket, err := uc.Create(context.Background(), input)

		require.ErrorIs(t, err, domain.ErrInvalidTitle)
		assert.Nil(t, ticket)
	})

	t.Run("wraps repository failures", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		responsible := newResponsible(t, "resp-1", "Ana")
		input := validInput
		input.ResponsibleID = "resp-1"

		responsibleRepo.On("FindByID", mock.Anything, "resp-1").Return(responsible, nil).Once()
		ticketRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Ticket")).
			Return(errors.New("conexão perdida")).Once()

		ticket, err := uc.Create(context.Background(), input)

		require.Error(t, err)
		assert.Nil(t, ticket)
	})
}

func TestTicketUseCase_GetByID(t *testing.T) {
	t.Parallel()

	t.Run("returns the ticket", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		expected := &domain.Ticket{ID: "ticket-1"}
		ticketRepo.On("FindByID", mock.Anything, "ticket-1").Return(expected, nil).Once()

		ticket, err := uc.GetByID(context.Background(), "ticket-1")

		require.NoError(t, err)
		assert.Same(t, expected, ticket)
	})

	t.Run("propagates not found", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		ticketRepo.On("FindByID", mock.Anything, "missing").
			Return(nil, domain.ErrTicketNotFound).Once()

		ticket, err := uc.GetByID(context.Background(), "missing")

		require.ErrorIs(t, err, domain.ErrTicketNotFound)
		assert.Nil(t, ticket)
	})
}

func TestTicketUseCase_List(t *testing.T) {
	t.Parallel()

	t.Run("normalizes pagination defaults and applies the default sort", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		expectedFilter := domain.TicketFilter{
			Sort:     domain.SortByCreatedAtDesc,
			Page:     1,
			PageSize: 20,
		}
		ticketRepo.On("FindAll", mock.Anything, expectedFilter).
			Return([]*domain.Ticket{{ID: "t1"}}, int64(1), nil).Once()

		output, err := uc.List(context.Background(), usecase.ListTicketsInput{Page: 0, PageSize: 0})

		require.NoError(t, err)
		assert.Equal(t, int64(1), output.Total)
		assert.Equal(t, 1, output.Page)
		assert.Equal(t, 20, output.PageSize)
	})

	t.Run("clamps an excessive page size", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		ticketRepo.On("FindAll", mock.Anything, mock.MatchedBy(func(f domain.TicketFilter) bool {
			return f.PageSize == 100
		})).Return([]*domain.Ticket{}, int64(0), nil).Once()

		_, err := uc.List(context.Background(), usecase.ListTicketsInput{PageSize: 999})

		require.NoError(t, err)
	})

	t.Run("wraps repository failures", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		ticketRepo.On("FindAll", mock.Anything, mock.Anything).
			Return(nil, int64(0), errors.New("timeout")).Once()

		output, err := uc.List(context.Background(), usecase.ListTicketsInput{})

		require.Error(t, err)
		assert.Nil(t, output)
	})
}

func TestTicketUseCase_Update(t *testing.T) {
	t.Parallel()

	validUpdate := usecase.UpdateTicketInput{
		Title:       "Novo título",
		Description: "Nova descrição",
		Priority:    domain.PriorityHigh,
		Status:      domain.StatusInProgress,
	}

	t.Run("updates an existing ticket", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		existing, err := domain.NewTicket("título", "descrição", domain.PriorityLow, "resp-1", "Ana")
		require.NoError(t, err)
		existing.ID = "ticket-1"

		ticketRepo.On("FindByID", mock.Anything, "ticket-1").Return(existing, nil).Once()
		ticketRepo.On("Update", mock.Anything, existing).Return(nil).Once()

		ticket, err := uc.Update(context.Background(), "ticket-1", validUpdate)

		require.NoError(t, err)
		assert.Equal(t, "Novo título", ticket.Title)
		assert.Equal(t, domain.StatusInProgress, ticket.Status)
	})

	t.Run("returns not found when the ticket does not exist", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		ticketRepo.On("FindByID", mock.Anything, "missing").
			Return(nil, domain.ErrTicketNotFound).Once()

		ticket, err := uc.Update(context.Background(), "missing", validUpdate)

		require.ErrorIs(t, err, domain.ErrTicketNotFound)
		assert.Nil(t, ticket)
	})

	t.Run("rejects invalid edits without touching the repository", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		existing, err := domain.NewTicket("título", "descrição", domain.PriorityLow, "resp-1", "Ana")
		require.NoError(t, err)
		existing.ID = "ticket-1"

		ticketRepo.On("FindByID", mock.Anything, "ticket-1").Return(existing, nil).Once()

		invalidUpdate := validUpdate
		invalidUpdate.Status = domain.Status("pausado")

		ticket, err := uc.Update(context.Background(), "ticket-1", invalidUpdate)

		require.ErrorIs(t, err, domain.ErrInvalidStatus)
		assert.Nil(t, ticket)
		ticketRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("wraps repository update failures", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		existing, err := domain.NewTicket("título", "descrição", domain.PriorityLow, "resp-1", "Ana")
		require.NoError(t, err)
		existing.ID = "ticket-1"

		ticketRepo.On("FindByID", mock.Anything, "ticket-1").Return(existing, nil).Once()
		ticketRepo.On("Update", mock.Anything, existing).Return(errors.New("timeout")).Once()

		ticket, err := uc.Update(context.Background(), "ticket-1", validUpdate)

		require.Error(t, err)
		assert.Nil(t, ticket)
	})
}

func TestTicketUseCase_AssignResponsible(t *testing.T) {
	t.Parallel()

	t.Run("assigns manually to the given responsible", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		existing, err := domain.NewTicket("título", "descrição", domain.PriorityLow, "resp-1", "Ana")
		require.NoError(t, err)
		existing.ID = "ticket-1"

		bruno := newResponsible(t, "resp-2", "Bruno")

		ticketRepo.On("FindByID", mock.Anything, "ticket-1").Return(existing, nil).Once()
		responsibleRepo.On("FindByID", mock.Anything, "resp-2").Return(bruno, nil).Once()
		ticketRepo.On("Update", mock.Anything, existing).Return(nil).Once()

		ticket, err := uc.AssignResponsible(context.Background(), "ticket-1", "resp-2")

		require.NoError(t, err)
		assert.Equal(t, "resp-2", ticket.ResponsibleID)
		assert.Equal(t, "Bruno", ticket.ResponsibleName)
	})

	t.Run("assigns automatically when no responsible is given", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		existing, err := domain.NewTicket("título", "descrição", domain.PriorityLow, "resp-1", "Ana")
		require.NoError(t, err)
		existing.ID = "ticket-1"

		ana := newResponsible(t, "resp-1", "Ana")
		bruno := newResponsible(t, "resp-2", "Bruno")

		ticketRepo.On("FindByID", mock.Anything, "ticket-1").Return(existing, nil).Once()
		responsibleRepo.On("FindAll", mock.Anything).Return([]*domain.Responsible{ana, bruno}, nil).Once()
		ticketRepo.On("CountOpenByResponsible", mock.Anything, "resp-1").Return(int64(2), nil).Once()
		ticketRepo.On("CountOpenByResponsible", mock.Anything, "resp-2").Return(int64(0), nil).Once()
		ticketRepo.On("Update", mock.Anything, existing).Return(nil).Once()

		ticket, err := uc.AssignResponsible(context.Background(), "ticket-1", "")

		require.NoError(t, err)
		assert.Equal(t, "resp-2", ticket.ResponsibleID)
	})

	t.Run("returns not found when the ticket does not exist", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		ticketRepo.On("FindByID", mock.Anything, "missing").
			Return(nil, domain.ErrTicketNotFound).Once()

		ticket, err := uc.AssignResponsible(context.Background(), "missing", "resp-1")

		require.ErrorIs(t, err, domain.ErrTicketNotFound)
		assert.Nil(t, ticket)
	})

	t.Run("returns not found when the chosen responsible does not exist", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		existing, err := domain.NewTicket("título", "descrição", domain.PriorityLow, "resp-1", "Ana")
		require.NoError(t, err)
		existing.ID = "ticket-1"

		ticketRepo.On("FindByID", mock.Anything, "ticket-1").Return(existing, nil).Once()
		responsibleRepo.On("FindByID", mock.Anything, "resp-404").
			Return(nil, domain.ErrResponsibleNotFound).Once()

		ticket, err := uc.AssignResponsible(context.Background(), "ticket-1", "resp-404")

		require.ErrorIs(t, err, domain.ErrResponsibleNotFound)
		assert.Nil(t, ticket)
	})

	t.Run("wraps repository update failures", func(t *testing.T) {
		t.Parallel()

		ticketRepo := new(mocks.TicketRepository)
		responsibleRepo := new(mocks.ResponsibleRepository)
		uc := usecase.NewTicketUseCase(ticketRepo, responsibleRepo)

		existing, err := domain.NewTicket("título", "descrição", domain.PriorityLow, "resp-1", "Ana")
		require.NoError(t, err)
		existing.ID = "ticket-1"

		bruno := newResponsible(t, "resp-2", "Bruno")

		ticketRepo.On("FindByID", mock.Anything, "ticket-1").Return(existing, nil).Once()
		responsibleRepo.On("FindByID", mock.Anything, "resp-2").Return(bruno, nil).Once()
		ticketRepo.On("Update", mock.Anything, existing).Return(errors.New("timeout")).Once()

		ticket, err := uc.AssignResponsible(context.Background(), "ticket-1", "resp-2")

		require.Error(t, err)
		assert.Nil(t, ticket)
	})
}
