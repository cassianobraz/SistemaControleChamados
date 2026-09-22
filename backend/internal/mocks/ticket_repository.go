// Package mocks contains hand-written testify mocks for the domain and
// usecase interfaces, used to unit test consumers in isolation.
package mocks

import (
	"context"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

// TicketRepository is a testify mock implementing domain.TicketRepository.
type TicketRepository struct {
	mock.Mock
}

func (m *TicketRepository) Create(ctx context.Context, ticket *domain.Ticket) error {
	args := m.Called(ctx, ticket)
	return args.Error(0)
}

func (m *TicketRepository) FindByID(ctx context.Context, id string) (*domain.Ticket, error) {
	args := m.Called(ctx, id)

	ticket, _ := args.Get(0).(*domain.Ticket)

	return ticket, args.Error(1)
}

func (m *TicketRepository) FindAll(ctx context.Context, filter domain.TicketFilter) ([]*domain.Ticket, int64, error) {
	args := m.Called(ctx, filter)

	tickets, _ := args.Get(0).([]*domain.Ticket)

	return tickets, args.Get(1).(int64), args.Error(2)
}

func (m *TicketRepository) Update(ctx context.Context, ticket *domain.Ticket) error {
	args := m.Called(ctx, ticket)
	return args.Error(0)
}

func (m *TicketRepository) CountOpenByResponsible(ctx context.Context, responsibleID string) (int64, error) {
	args := m.Called(ctx, responsibleID)
	return args.Get(0).(int64), args.Error(1)
}
