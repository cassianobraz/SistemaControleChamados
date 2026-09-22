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

func TestResponsibleUseCase_List(t *testing.T) {
	t.Parallel()

	t.Run("returns responsibles with their open ticket count", func(t *testing.T) {
		t.Parallel()

		responsibleRepo := new(mocks.ResponsibleRepository)
		ticketRepo := new(mocks.TicketRepository)
		uc := usecase.NewResponsibleUseCase(responsibleRepo, ticketRepo)

		ana := newResponsible(t, "resp-1", "Ana")
		bruno := newResponsible(t, "resp-2", "Bruno")

		responsibleRepo.On("FindAll", mock.Anything).Return([]*domain.Responsible{ana, bruno}, nil).Once()
		ticketRepo.On("CountOpenByResponsible", mock.Anything, "resp-1").Return(int64(4), nil).Once()
		ticketRepo.On("CountOpenByResponsible", mock.Anything, "resp-2").Return(int64(0), nil).Once()

		result, err := uc.List(context.Background())

		require.NoError(t, err)
		require.Len(t, result, 2)
		assert.Equal(t, int64(4), result[0].OpenTickets)
		assert.Equal(t, int64(0), result[1].OpenTickets)
	})

	t.Run("wraps failures listing responsibles", func(t *testing.T) {
		t.Parallel()

		responsibleRepo := new(mocks.ResponsibleRepository)
		ticketRepo := new(mocks.TicketRepository)
		uc := usecase.NewResponsibleUseCase(responsibleRepo, ticketRepo)

		responsibleRepo.On("FindAll", mock.Anything).Return(nil, errors.New("timeout")).Once()

		result, err := uc.List(context.Background())

		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("wraps failures counting open tickets", func(t *testing.T) {
		t.Parallel()

		responsibleRepo := new(mocks.ResponsibleRepository)
		ticketRepo := new(mocks.TicketRepository)
		uc := usecase.NewResponsibleUseCase(responsibleRepo, ticketRepo)

		ana := newResponsible(t, "resp-1", "Ana")
		responsibleRepo.On("FindAll", mock.Anything).Return([]*domain.Responsible{ana}, nil).Once()
		ticketRepo.On("CountOpenByResponsible", mock.Anything, "resp-1").
			Return(int64(0), errors.New("timeout")).Once()

		result, err := uc.List(context.Background())

		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestResponsibleUseCase_Create(t *testing.T) {
	t.Parallel()

	t.Run("creates a responsible with zero open tickets", func(t *testing.T) {
		t.Parallel()

		responsibleRepo := new(mocks.ResponsibleRepository)
		ticketRepo := new(mocks.TicketRepository)
		uc := usecase.NewResponsibleUseCase(responsibleRepo, ticketRepo)

		responsibleRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Responsible")).Return(nil).Once()

		result, err := uc.Create(context.Background(), usecase.CreateResponsibleInput{Name: "Ana", Email: "ana@codificar.dev"})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "Ana", result.Responsible.Name)
		assert.Equal(t, int64(0), result.OpenTickets)
	})

	t.Run("rejects an invalid name before touching the repository", func(t *testing.T) {
		t.Parallel()

		responsibleRepo := new(mocks.ResponsibleRepository)
		ticketRepo := new(mocks.TicketRepository)
		uc := usecase.NewResponsibleUseCase(responsibleRepo, ticketRepo)

		result, err := uc.Create(context.Background(), usecase.CreateResponsibleInput{Name: "   "})

		require.ErrorIs(t, err, domain.ErrInvalidResponsibleName)
		assert.Nil(t, result)
		responsibleRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})
}

func TestResponsibleUseCase_Update(t *testing.T) {
	t.Parallel()

	t.Run("updates the responsible and returns their current workload", func(t *testing.T) {
		t.Parallel()

		responsibleRepo := new(mocks.ResponsibleRepository)
		ticketRepo := new(mocks.TicketRepository)
		uc := usecase.NewResponsibleUseCase(responsibleRepo, ticketRepo)

		ana := newResponsible(t, "resp-1", "Ana")

		responsibleRepo.On("FindByID", mock.Anything, "resp-1").Return(ana, nil).Once()
		responsibleRepo.On("Update", mock.Anything, ana).Return(nil).Once()
		ticketRepo.On("CountOpenByResponsible", mock.Anything, "resp-1").Return(int64(2), nil).Once()

		result, err := uc.Update(context.Background(), "resp-1", usecase.UpdateResponsibleInput{
			Name:   "Ana Paula",
			Email:  "ana.paula@codificar.dev",
			Active: false,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "Ana Paula", result.Responsible.Name)
		assert.False(t, result.Responsible.Active)
		assert.Equal(t, int64(2), result.OpenTickets)
	})

	t.Run("returns an error when the responsible does not exist", func(t *testing.T) {
		t.Parallel()

		responsibleRepo := new(mocks.ResponsibleRepository)
		ticketRepo := new(mocks.TicketRepository)
		uc := usecase.NewResponsibleUseCase(responsibleRepo, ticketRepo)

		responsibleRepo.On("FindByID", mock.Anything, "missing").Return(nil, domain.ErrResponsibleNotFound).Once()

		result, err := uc.Update(context.Background(), "missing", usecase.UpdateResponsibleInput{Name: "Ana"})

		require.ErrorIs(t, err, domain.ErrResponsibleNotFound)
		assert.Nil(t, result)
	})
}

func TestResponsibleUseCase_Delete(t *testing.T) {
	t.Parallel()

	t.Run("removes the responsible", func(t *testing.T) {
		t.Parallel()

		responsibleRepo := new(mocks.ResponsibleRepository)
		ticketRepo := new(mocks.TicketRepository)
		uc := usecase.NewResponsibleUseCase(responsibleRepo, ticketRepo)

		responsibleRepo.On("Delete", mock.Anything, "resp-1").Return(nil).Once()

		err := uc.Delete(context.Background(), "resp-1")

		require.NoError(t, err)
	})

	t.Run("wraps the not-found error from the repository", func(t *testing.T) {
		t.Parallel()

		responsibleRepo := new(mocks.ResponsibleRepository)
		ticketRepo := new(mocks.TicketRepository)
		uc := usecase.NewResponsibleUseCase(responsibleRepo, ticketRepo)

		responsibleRepo.On("Delete", mock.Anything, "missing").Return(domain.ErrResponsibleNotFound).Once()

		err := uc.Delete(context.Background(), "missing")

		require.ErrorIs(t, err, domain.ErrResponsibleNotFound)
	})
}
