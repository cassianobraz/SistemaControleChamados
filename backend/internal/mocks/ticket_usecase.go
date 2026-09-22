package mocks

import (
	"context"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/usecase"
	"github.com/stretchr/testify/mock"
)

// TicketUseCase is a testify mock implementing usecase.TicketUseCase.
type TicketUseCase struct {
	mock.Mock
}

func (m *TicketUseCase) Create(ctx context.Context, input usecase.CreateTicketInput) (*domain.Ticket, error) {
	args := m.Called(ctx, input)

	ticket, _ := args.Get(0).(*domain.Ticket)

	return ticket, args.Error(1)
}

func (m *TicketUseCase) GetByID(ctx context.Context, id string) (*domain.Ticket, error) {
	args := m.Called(ctx, id)

	ticket, _ := args.Get(0).(*domain.Ticket)

	return ticket, args.Error(1)
}

func (m *TicketUseCase) List(ctx context.Context, input usecase.ListTicketsInput) (*usecase.ListTicketsOutput, error) {
	args := m.Called(ctx, input)

	output, _ := args.Get(0).(*usecase.ListTicketsOutput)

	return output, args.Error(1)
}

func (m *TicketUseCase) Update(ctx context.Context, id string, input usecase.UpdateTicketInput) (*domain.Ticket, error) {
	args := m.Called(ctx, id, input)

	ticket, _ := args.Get(0).(*domain.Ticket)

	return ticket, args.Error(1)
}

func (m *TicketUseCase) AssignResponsible(ctx context.Context, id string, responsibleID string) (*domain.Ticket, error) {
	args := m.Called(ctx, id, responsibleID)

	ticket, _ := args.Get(0).(*domain.Ticket)

	return ticket, args.Error(1)
}
